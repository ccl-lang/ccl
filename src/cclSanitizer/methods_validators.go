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

func validateFieldTypeUsages(typeUsages []fieldTypeUsageCheck) error {
	for _, usage := range typeUsages {
		if err := validateTypeUsageCompleteness(usage, usage.typeUsage); err != nil {
			return err
		}
	}

	return nil
}

func validateTypeUsageCompleteness(
	usageInfo fieldTypeUsageCheck,
	typeUsage *cclValues.CCLTypeUsage,
) error {
	if typeUsage == nil || typeUsage.GetDefinition() == nil {
		return &AstSanitizationError{
			Message:        "field type resolved to nil",
			SourcePosition: usageInfo.sourcePosition,
		}
	}

	typeDef := typeUsage.GetDefinition()
	if typeDef.IsIncomplete() {
		return &AstSanitizationError{
			Message:        "unknown type '" + typeDef.GetFullName() + "' for field '" + usageInfo.fieldName + "' in model '" + usageInfo.modelName + "'",
			SourcePosition: usageInfo.sourcePosition,
		}
	}

	underlying := typeUsage.GetUnderlyingType()
	if underlying != nil {
		if err := validateTypeUsageCompleteness(usageInfo, underlying); err != nil {
			return err
		}
	}

	for _, arg := range typeUsage.GetGenericArgs() {
		if err := validateTypeUsageCompleteness(usageInfo, arg); err != nil {
			return err
		}
	}

	return nil
}

//---------------------------------------------------------
