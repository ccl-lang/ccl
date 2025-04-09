package cclValues

type cclTypeFlag int

// SourceCodeDefinition is a struct that represents a source code definition
// and all the information about a cll source file.
// This struct is NOT thread-safe.
type SourceCodeDefinition struct {
	// Models is an array of model definitions.
	Models []*ModelDefinition

	// GlobalAttributes is an array of attribute definitions which are applied
	// to the whole source code.
	GlobalAttributes []*AttributeDefinition

	// modelIdCounter is a counter that is used to generate unique model IDs.
	modelIdCounter int64
}

// AttributeDefinition is a struct that represents an attribute definition
// in the source code with its parameters.
type AttributeDefinition struct {
	// Name is the name of the attribute.
	Name string

	// Parameters is the list of parameters for the attribute.
	Parameters []*ParameterDefinition
}

// CCLTypeInfo is a struct that represents a CCL *type*.
// Now this type can be a built-in type or a custom type defined by the user.
// NOTE: This struct only holds general info about a type, NOT its values, etc.
type CCLTypeInfo struct {
	// name is the name of the type.
	name string

	// typeFlags contains flags applied to the type.
	// To work with this field, you should use the TypeFlag* constants.
	typeFlags cclTypeFlag
}

// ParameterDefinition is a struct that represents a parameter definition.
// This parameter is for when the user is passing a parameter to a function or any
// other place in the source code, such as attributes.
type ParameterDefinition struct {
	// value is the value of the parameter, specified in the source code.
	// Please note that this field is not exported, you should use the
	// methods to get or set this field.
	value any

	// ValueType is the type of the parameter.
	ValueType *CCLTypeInfo
}

// ModelDefinition is a struct that represents a model definition.
type ModelDefinition struct {
	// ModelId is the unique identifier of the model.
	ModelId int64

	// Name is the name of the model.
	Name string

	// Fields is an array of field definitions.
	Fields []*FieldDefinition

	// Attributes is an array of attribute definitions which are applied
	// to the model.
	Attributes []*AttributeDefinition
}

// FieldDefinition is a struct that represents a field definition.
type FieldDefinition struct {
	// OwnedBy is a reference to the model that owns the field.
	OwnedBy *ModelDefinition

	// Name is the name of the field.
	Name string

	// Type is the type of the field.
	Type *CCLTypeInfo
}
