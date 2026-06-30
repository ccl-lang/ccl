package cclSanitizer

import (
	"strings"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
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

	definition := &cclValues.SourceCodeDefinition{
		CodeContext: ctx,
		FilePath:    ast.FilePath,
		Namespace:   fileNamespace,
	}
	sourceFileId := ctx.RegisterSourceCodeDefinition(definition)
	nameValidator := newFieldNameValidator(ast, fileNamespace)
	fieldTypeUsages := []*fieldTypeUsageCheck{}

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
		attrUsage.SourceFileId = sourceFileId
		definition.GlobalAttributes = append(definition.GlobalAttributes, attrUsage)
	}

	for _, fileAttr := range ast.FileAttributes {
		if fileAttr == nil {
			return nil, &AstSanitizationError{
				Message: "nil file attribute in AST",
			}
		}

		attrUsage, err := ResolveAttributeUsage(ctx, fileAttr)
		if err != nil {
			return nil, err
		}
		attrUsage.SourceFileId = sourceFileId
		definition.FileAttributes = append(definition.FileAttributes, attrUsage)
	}

	for _, namespaceAttr := range ast.NamespaceAttributes {
		if namespaceAttr == nil {
			return nil, &AstSanitizationError{
				Message: "nil namespace attribute in AST",
			}
		}

		attrUsage, err := ResolveAttributeUsage(ctx, namespaceAttr)
		if err != nil {
			return nil, err
		}
		attrUsage.SourceFileId = sourceFileId
		attrUsage.Namespace = namespaceAttr.Namespace
		if attrUsage.Namespace == "" {
			attrUsage.Namespace = fileNamespace
		}
		definition.NamespaceAttributes = append(definition.NamespaceAttributes, attrUsage)
	}
	ctx.RegisterScopedAttributes(definition)

	for _, enumAst := range ast.Enums {
		if enumAst == nil {
			return nil, &AstSanitizationError{
				Message: "nil enum declaration in AST",
			}
		}

		enumNamespace := enumAst.Namespace
		if enumNamespace == "" {
			enumNamespace = fileNamespace
		}

		enumDef, enumTypeDef, err := sanitizeEnumDeclaration(
			ctx,
			definition,
			sourceFileId,
			enumNamespace,
			nil,
			enumAst,
		)
		if err != nil {
			return nil, err
		}
		_ = enumDef
		definition.TypeDefinitions = append(definition.TypeDefinitions, enumTypeDef)
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
			SourceFileId:   sourceFileId,
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
			attrUsage.SourceFileId = sourceFileId
			modelDef.Attributes = append(modelDef.Attributes, attrUsage)
		}

		for _, enumAst := range modelAst.Enums {
			if enumAst == nil {
				return nil, &AstSanitizationError{
					Message:        "nil nested enum declaration in AST",
					SourcePosition: modelAst.SourcePosition,
				}
			}

			enumDef, enumTypeDef, err := sanitizeEnumDeclaration(
				ctx,
				definition,
				sourceFileId,
				modelNamespace,
				modelDef,
				enumAst,
			)
			if err != nil {
				return nil, err
			}
			modelDef.Enums = append(modelDef.Enums, enumDef)
			definition.TypeDefinitions = append(definition.TypeDefinitions, enumTypeDef)
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
				attrUsage.SourceFileId = sourceFileId
				fieldDef.Attributes = append(fieldDef.Attributes, attrUsage)
			}

			if fieldAst.Type != nil {
				typeUsage, err := ResolveTypeUsageForModel(
					ctx,
					modelNamespace,
					modelDef.Name,
					fieldAst.Type,
				)
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
				fieldTypeUsages = append(fieldTypeUsages, &fieldTypeUsageCheck{
					modelName:      modelDef.Name,
					fieldName:      fieldDef.Name,
					typeUsage:      typeUsage,
					sourcePosition: fieldAst.SourcePosition,
				})
			} else {
				return nil, &AstSanitizationError{
					Message:        "field has no type",
					SourcePosition: fieldAst.SourcePosition,
				}
			}

			if fieldAst.Value != nil {
				value, err := resolveFieldDefaultValue(
					ctx,
					fieldDef.Type,
					modelNamespace,
					modelDef.Name,
					fieldAst.Value,
				)
				if err != nil {
					return nil, err
				}
				fieldDef.ChangeDefaultValue(value)
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
		modelTypeDef.ChangeSourceFileId(sourceFileId)

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
