package cclSanitizer

import (
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclUtils"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

// SanitizeCCLAst converts a syntax-only AST into IR (cclValues).
func SanitizeCCLAst(
	ctx *cclValues.CCLCodeContext,
	ast *cclAst.CCLFileAST,
) (*cclValues.SourceCodeDefinition, error) {
	if ast == nil {
		return nil, &AstSanitizationError{
			Message: "missing CCL AST",
		}
	}

	if ctx == nil {
		ctx = cclValues.NewCCLCodeContext()
	}

	fileNamespace := ast.Namespace
	if fileNamespace == "" {
		fileNamespace = gValues.DefaultMainNamespace
	}

	definition := &cclValues.SourceCodeDefinition{}
	nameValidator := newFieldNameValidator(ast, fileNamespace)
	fieldTypeUsages := []fieldTypeUsageCheck{}

	for _, globalAttr := range ast.GlobalAttributes {
		if globalAttr == nil {
			return nil, &AstSanitizationError{
				Message: "nil global attribute in AST",
			}
		}

		attrUsage, err := ResolveAttributeUsage(ctx, globalAttr)
		if err != nil {
			return nil, err
		}
		definition.GlobalAttributes = append(definition.GlobalAttributes, attrUsage)
	}

	for _, modelAst := range ast.Models {
		if modelAst == nil {
			return nil, &AstSanitizationError{
				Message: "nil model declaration in AST",
			}
		}

		modelNamespace := modelAst.Namespace
		if modelNamespace == "" {
			modelNamespace = fileNamespace
		}

		if definition.GetModelByName(modelAst.Name) != nil {
			return nil, &cclErrors.DuplicateModelError{
				ModelName:      modelAst.Name,
				SourcePosition: modelAst.SourcePosition,
			}
		}

		modelDef := &cclValues.ModelDefinition{
			ModelId:        definition.GetNextModelId(),
			Name:           modelAst.Name,
			Namespace:      modelNamespace,
			SourcePosition: modelAst.SourcePosition,
		}

		for _, attrAst := range modelAst.Attributes {
			if attrAst == nil {
				return nil, &AstSanitizationError{
					Message:        "nil model attribute in AST",
					SourcePosition: modelAst.SourcePosition,
				}
			}

			attrUsage, err := ResolveAttributeUsage(ctx, attrAst)
			if err != nil {
				return nil, err
			}
			modelDef.Attributes = append(modelDef.Attributes, attrUsage)
		}

		for _, fieldAst := range modelAst.Fields {
			if fieldAst == nil {
				return nil, &AstSanitizationError{
					Message:        "nil field declaration in AST",
					SourcePosition: modelAst.SourcePosition,
				}
			}

			if err := nameValidator.ValidateFieldName(
				modelNamespace,
				modelDef.Name,
				fieldAst,
			); err != nil {
				return nil, err
			}

			if modelDef.GetFieldByName(fieldAst.Name) != nil {
				return nil, &cclErrors.DuplicateFieldError{
					ModelName:      modelDef.Name,
					FieldName:      fieldAst.Name,
					SourcePosition: fieldAst.SourcePosition,
				}
			}

			fieldDef := &cclValues.ModelFieldDefinition{
				OwnedBy: modelDef,
				Name:    fieldAst.Name,
			}

			for _, attrAst := range fieldAst.Attributes {
				if attrAst == nil {
					return nil, &AstSanitizationError{
						Message:        "nil field attribute in AST",
						SourcePosition: fieldAst.SourcePosition,
					}
				}

				attrUsage, err := ResolveAttributeUsage(ctx, attrAst)
				if err != nil {
					return nil, err
				}
				fieldDef.Attributes = append(fieldDef.Attributes, attrUsage)
			}

			if fieldAst.Type != nil && fieldAst.Value != nil {
				return nil, &AstSanitizationError{
					Message:        "field has both type and value",
					SourcePosition: fieldAst.SourcePosition,
				}
			}

			if fieldAst.Type != nil {
				typeUsage, err := ResolveTypeUsage(ctx, modelNamespace, fieldAst.Type)
				if err != nil {
					return nil, err
				}
				if typeUsage == nil {
					return nil, &AstSanitizationError{
						Message:        "field type resolved to nil",
						SourcePosition: fieldAst.SourcePosition,
					}
				}
				fieldDef.ChangeValueType(typeUsage)
				fieldTypeUsages = append(fieldTypeUsages, fieldTypeUsageCheck{
					modelName:      modelDef.Name,
					fieldName:      fieldDef.Name,
					typeUsage:      typeUsage,
					sourcePosition: fieldAst.SourcePosition,
				})
			} else if fieldAst.Value != nil {
				valueType, value, err := resolveFieldAssignment(ctx, fieldAst.Value, fieldAst.SourcePosition)
				if err != nil {
					return nil, err
				}
				fieldDef.ChangeValueType(valueType)
				fieldDef.ChangeDefaultValue(value)
				fieldTypeUsages = append(fieldTypeUsages, fieldTypeUsageCheck{
					modelName:      modelDef.Name,
					fieldName:      fieldDef.Name,
					typeUsage:      valueType,
					sourcePosition: fieldAst.SourcePosition,
				})
			} else {
				return nil, &AstSanitizationError{
					Message:        "field has no type or value",
					SourcePosition: fieldAst.SourcePosition,
				}
			}

			modelDef.Fields = append(modelDef.Fields, fieldDef)
		}

		modelTypeDef, err := ctx.NewModelTypeDefinition(&cclValues.SimpleTypeName{
			TypeName:  modelDef.Name,
			Namespace: modelNamespace,
		}, modelDef)
		if err != nil {
			return nil, err
		}

		definition.TypeDefinitions = append(definition.TypeDefinitions, modelTypeDef)
	}

	if err := validateFieldTypeUsages(fieldTypeUsages); err != nil {
		return nil, err
	}

	return definition, nil
}

func newFieldNameValidator(ast *cclAst.CCLFileAST, defaultNamespace string) *fieldNameValidator {
	modelNamesByNamespace := map[string]map[string]string{}
	for _, modelAst := range ast.Models {
		if modelAst == nil {
			continue
		}

		namespace := modelAst.Namespace
		if namespace == "" {
			namespace = defaultNamespace
		}

		normalized := normalizeName(modelAst.Name)
		if modelNamesByNamespace[namespace] == nil {
			modelNamesByNamespace[namespace] = map[string]string{}
		}
		modelNamesByNamespace[namespace][normalized] = modelAst.Name
	}

	return &fieldNameValidator{
		modelNamesByNamespace: modelNamesByNamespace,
	}
}

func normalizeName(name string) string {
	return strings.ToLower(name)
}

func resolveFieldAssignment(
	ctx *cclValues.CCLCodeContext,
	value cclAst.ValueExpression,
	sourcePos *cclUtils.SourceCodePosition,
) (*cclValues.CCLTypeUsage, any, error) {
	switch expr := value.(type) {
	case *cclAst.IdentifierValueExpression:
		if expr == nil {
			return nil, nil, &AstSanitizationError{
				Message:        "invalid identifier expression",
				SourcePosition: sourcePos,
			}
		}

		targetVariable := ctx.GetGlobalVariable(expr.Name)
		if targetVariable == nil {
			return nil, nil, &AstSanitizationError{
				Message:        "undefined identifier '" + expr.Name + "'",
				SourcePosition: expr.SourcePosition,
			}
		}

		if targetVariable.IsAutomatic() {
			return ctx.NewPointerTypeUsage(targetVariable.Type), &cclValues.VariableUsageInstance{
				Name:       expr.Name,
				Definition: targetVariable,
			}, nil
		}

		return targetVariable.Type, targetVariable.GetValue(), nil
	default:
		return nil, nil, &AstSanitizationError{
			Message:        "unsupported field assignment value",
			SourcePosition: sourcePos,
		}
	}
}
