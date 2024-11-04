package cclValues

// SourceCodeDefinition is a struct that represents a source code definition
// and all the information about a cll source file.
// This struct is NOT thread-safe.
type SourceCodeDefinition struct {
	// Models is an array of model definitions.
	Models []*ModelDefinition

	// modelIdCounter is a counter that is used to generate unique model IDs.
	modelIdCounter int64
}

// ModelDefinition is a struct that represents a model definition.
type ModelDefinition struct {
	// ModelId is the unique identifier of the model.
	ModelId int64

	// Name is the name of the model.
	Name string

	// Fields is an array of field definitions.
	Fields []*FieldDefinition
}

// FieldDefinition is a struct that represents a field definition.
type FieldDefinition struct {
	// OwnedBy is a reference to the model that owns the field.
	OwnedBy *ModelDefinition

	// Name is the name of the field.
	Name string

	// Type is the type of the field.
	Type string

	// ExtraOperators is a string that contains extra operators for the field.
	ExtraOperators string
}
