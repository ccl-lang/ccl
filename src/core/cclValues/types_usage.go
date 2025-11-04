package cclValues

// SpecificCCLTypeUsageGenerator is a function type that generates a specific CCLTypeUsage
// given a CCLCodeContext.
type SpecificCCLTypeUsageGenerator func(c *CCLCodeContext) *CCLTypeUsage

// CCLTypeUsage is a struct that represents a CCL type usage instance.
// This type is for when the user is using a type in the source code,
// such as in variable definitions, function parameters, etc.
type CCLTypeUsage struct {
	// definition is a reference to the type definition.
	// NOTE: The type definition SHOULD be pooled/shared between all usages
	// of the same type, to save memory.
	definition *CCLTypeDefinition

	// underlyingType is the underlying type of the current type.
	// This field is used when the current type is a complex type,
	// such as arrays, maps, pointers, etc.
	// For simple types, this field is set to nil.
	underlyingType *CCLTypeUsage

	// genericArgs is the list of generic arguments passed to the type.
	genericArgs []*CCLTypeUsage
}

// VariableUsageInstance is a struct that represents a variable usage instance.
// This variable is for when the user is using a variable in the source code,
// such as in function calls, attribute calls, etc.
type VariableUsageInstance struct {
	// Name is the name of the variable that is being used
	Name string

	// Definition points to the specified variable definition.
	// This could be a global variable or a local variable, or
	// could be set to nil, so always validate it before using it.
	Definition *VariableDefinition
}
