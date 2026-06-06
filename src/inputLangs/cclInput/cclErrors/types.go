package cclErrors

import (
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils"
)

// ConflictKind represents the kind of conflict that occurred, such as "model", "builtin type", "reserved name", etc.
type ConflictKind string

// ValidationError is an error that is returned when a validation
// error occurs.
type ValidationError struct {
	Message string
}

// InvalidAttributeUsageError represents an attribute usage with wrong set of params
// or used in wrong context or on wrong entity.
type InvalidAttributeUsageError struct {
	Message        string
	AttrName       cclAttr.CCLAttributeName
	SourcePosition *cclUtils.SourceCodePosition
}

// DuplicateFieldError is an error that is returned when a field
// with the same name already exists in a model.
type DuplicateFieldError struct {
	ModelName      string
	FieldName      string
	SourcePosition *cclUtils.SourceCodePosition
}

// DuplicateModelError is an error that is returned when a model
// with the same name already exists in a ccl source file / project.
type DuplicateModelError struct {
	ModelName      string
	SourcePosition *cclUtils.SourceCodePosition
}

// UnsupportedFieldTypeError is an error that is returned when an
// unsupported field type is encountered in a certain model, for a
// certain field when compiling to a certain target language.
type UnsupportedFieldTypeError struct {
	TypeName       string
	FieldName      string
	ModelName      string
	TargetLanguage string
}

// UnsupportedTypeDefinitionError is an error that is returned when an
// unsupported type definition is encountered when compiling to a
// certain target language.
type UnsupportedTypeDefinitionError struct {
	TypeName       string
	TargetLanguage string
}

// UnsupportedFileNamingStyleError is an error that is returned when an
// unsupported file naming style is encountered for a certain model
// when compiling to a certain target language.
type UnsupportedFileNamingStyleError struct {
	// ModelName is the model name that the style was used in.
	ModelName string

	// StyleName is the specific style name that the user has provided but is not
	// supported by us.
	StyleName string

	// SupportedStyles is a representation of a list of supported styles.
	// You can use JoinStr to append them together.
	SupportedStyles string

	// TargetLanguage is the target language in which we where trying to generate
	// code.
	TargetLanguage string
}

// FieldNameConflictError is an error that is returned when a field name
// conflicts with a reserved or existing name (such as model or builtin type names).
type FieldNameConflictError struct {
	ModelName      string
	FieldName      string
	ConflictName   string
	Kind           ConflictKind
	Namespace      string
	SourcePosition *cclUtils.SourceCodePosition
}
