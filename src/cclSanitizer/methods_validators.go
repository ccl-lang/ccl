package cclSanitizer

import (
	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

//---------------------------------------------------------

func (v *fieldNameValidator) ValidateFieldName(
	namespace string,
	modelName string,
	fieldAst *cclAst.FieldDecl,
) error {
	if fieldAst == nil {
		return &AstSanitizationError{
			Message: "nil field declaration in AST",
		}
	}

	normalizedFieldName := normalizeName(fieldAst.Name)
	if builtinName, ok := builtinTypeNamesLower[normalizedFieldName]; ok {
		return &cclErrors.FieldNameConflictError{
			ModelName:      modelName,
			FieldName:      fieldAst.Name,
			ConflictName:   builtinName,
			Kind:           cclErrors.ConflictKindBuiltinType,
			Namespace:      cclValues.NamespaceBuiltin,
			SourcePosition: fieldAst.SourcePosition,
		}
	}

	if v == nil {
		return nil
	}

	modelNames := v.modelNamesByNamespace[namespace]
	if modelNames != nil {
		if conflictName, ok := modelNames[normalizedFieldName]; ok {
			return &cclErrors.FieldNameConflictError{
				ModelName:      modelName,
				FieldName:      fieldAst.Name,
				ConflictName:   conflictName,
				Kind:           cclErrors.ConflictKindModel,
				Namespace:      namespace,
				SourcePosition: fieldAst.SourcePosition,
			}
		}
	}

	return nil
}

//---------------------------------------------------------

func (u *fieldTypeUsageCheck) validateTypeUsageCompleteness(
	typeUsage *cclValues.CCLTypeUsage,
) error {
	if typeUsage == nil || typeUsage.GetDefinition() == nil {
		return &AstSanitizationError{
			Message:        "field type resolved to nil",
			SourcePosition: u.sourcePosition,
		}
	}

	typeDef := typeUsage.GetDefinition()
	if typeDef.IsIncomplete() {
		return &AstSanitizationError{
			Message: "unknown type '" + typeDef.GetFullName() +
				"' for field '" + u.fieldName + "' in model '" + u.modelName + "'",
			SourcePosition: u.sourcePosition,
		}
	}

	underlying := typeUsage.GetUnderlyingType()
	if underlying != nil {
		if err := u.validateTypeUsageCompleteness(underlying); err != nil {
			return err
		}
	}

	for _, arg := range typeUsage.GetGenericArgs() {
		if err := u.validateTypeUsageCompleteness(arg); err != nil {
			return err
		}
	}

	return nil
}

//---------------------------------------------------------
