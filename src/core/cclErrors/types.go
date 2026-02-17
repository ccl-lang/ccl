package cclErrors

import "github.com/ccl-lang/ccl/src/core/cclUtils"

// ValidationError is an error that is returned when a validation
// error occurs.
type ValidationError struct {
	Message string
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
	ModelName       string
	StyleName       string
	SupportedStyles []string
	TargetLanguage  string
}

// FieldNameConflictError is an error that is returned when a field name
// conflicts with a reserved or existing name (such as model or builtin type names).
type FieldNameConflictError struct {
	ModelName      string
	FieldName      string
	ConflictName   string
	ConflictKind   string
	Namespace      string
	SourcePosition *cclUtils.SourceCodePosition
}
