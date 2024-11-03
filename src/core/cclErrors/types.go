package cclErrors

type ValidationError struct {
	Message string
}

type DuplicateFieldError struct {
	ModelName string
	FieldName string
}

type DuplicateModelError struct {
	ModelName string
}
