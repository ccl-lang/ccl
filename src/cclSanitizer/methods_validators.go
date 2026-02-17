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
