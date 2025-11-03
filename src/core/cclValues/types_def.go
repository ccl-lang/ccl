package cclValues

import "github.com/ccl-lang/ccl/src/core/cclUtils"

// SourceCodeDefinition is a struct that represents a source code definition
// and all the information about a cll source file.
// This struct is NOT thread-safe.
type SourceCodeDefinition struct {
	// Models is an array of model definitions.
	Models []*ModelDefinition

	// GlobalAttributes is an array of attribute definitions which are applied
	// to the whole source code.
	GlobalAttributes []*AttributeUsageInfo

	// modelIdCounter is a counter that is used to generate unique model IDs.
	modelIdCounter int64
}

// CCLTypeDefinition is a struct that represents a CCL *type*.
// Now this type can be a built-in type or a custom type defined by the user.
// NOTE: This struct only holds general info about a type, NOT its values, etc.
type CCLTypeDefinition struct {
	// name is the name of the type.
	name string

	// namespace is the namespace of the type.
	// For built-in types, this field is an empty string.
	// For custom types, this field is the namespace where the type is defined;
	// if no namespace is defined, it will be set to "global".
	namespace string

	// isIncomplete indicates whether the type definition is incomplete.
	isIncomplete bool

	// model is a reference to the model definition if the type is a custom model
	// type defined by the user.
	// For built-in and alias types, this field is nil.
	// Also take note that the typeFlags must contain TypeFlagCustomModel.
	model *ModelDefinition

	// typeFlags contains flags applied to the type.
	// To work with this field, you should use the TypeFlag* constants.
	typeFlags cclTypeFlag

	// genericParams is the parameters passed as generic-types to the
	// current type.
	// E.g. MyType<Type1, Type2> (where each of thees could recursively contain
	// other generic params in them).
	genericParams []*CCLTypeDefinition

	// SourcePosition is the position of the type in the source code.
	// Please note that for built-in types, this field is nil.
	SourcePosition *cclUtils.SourceCodePosition
}

// VariableDefinition is a struct that represents a variable definition.
// e.g. global variables, local variables, etc.
// for example, in the following code:
// var myVar int = 10
// the variable definition would be:
// Name: "myVar"
// Type: CCLTypeDefinition for "int"
// Value: 10
type VariableDefinition struct {
	// Name is the name of the variable that is being defined.
	Name string

	// Type is the type of the variable.
	Type *CCLTypeUsage

	// Value is the value of the variable.
	value any

	// isAutomaticVariable is a flag that indicates if the variable is an
	// automatic variable set by ccl compiler itself and cannot be overridden by the user.
	isAutomaticVariable bool

	// SourcePosition is the position of the variable in the source code.
	SourcePosition *cclUtils.SourceCodePosition
}

// ModelDefinition is a struct that represents a model definition.
type ModelDefinition struct {
	// ModelId is the unique identifier of the model.
	ModelId int64

	// Name is the name of the model.
	Name string

	// Fields is an array of field definitions.
	Fields []*ModelFieldDefinition

	// Attributes is an array of attribute definitions which are applied
	// to the model.
	Attributes []*AttributeUsageInfo

	// SourcePosition is the position of the model in the source code.
	SourcePosition *cclUtils.SourceCodePosition
}

// ModelFieldDefinition is a struct that represents a field definition.
type ModelFieldDefinition struct {
	// OwnedBy is a reference to the model that owns the field.
	OwnedBy *ModelDefinition

	// Name is the name of the field.
	Name string

	// Type is the type of the field.
	Type *CCLTypeUsage

	// Attributes is an array of attribute definitions which are applied
	// to this field.
	Attributes []*AttributeUsageInfo

	// value is the current value of this field.
	value any
}
