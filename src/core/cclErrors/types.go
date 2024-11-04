package cclErrors

// ValidationError is an error that is returned when a validation
// error occurs.
type ValidationError struct {
	Message string
}

// DuplicateFieldError is an error that is returned when a field
// with the same name already exists in a model.
type DuplicateFieldError struct {
	ModelName string
	FieldName string
}

// DuplicateModelError is an error that is returned when a model
// with the same name already exists in a ccl source file / project.
type DuplicateModelError struct {
	ModelName string
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
