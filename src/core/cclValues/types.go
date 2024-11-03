package cclValues

// ModelsMap is a map of model names to model definitions.
type ModelsMap map[string]*ModelDefinition
type FieldsMap map[string]*FieldDefinition

// SourceCodeDefinition is a struct that represents a source code definition
// and all the information about a cll source file.
// This struct is NOT thread-safe.
type SourceCodeDefinition struct {
	Models ModelsMap
}

type ModelDefinition struct {
	Name   string
	Fields FieldsMap
}

type FieldDefinition struct {
	Name           string
	Type           string
	ExtraOperators string
}
