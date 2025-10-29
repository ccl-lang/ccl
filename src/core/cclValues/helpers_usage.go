package cclValues

func NewBuiltinTypeUsage(typeName string) *CCLTypeUsage {
	if !IsBuiltinTypeName(typeName) {
		return nil
	}

	return NewTypeUsage(NewTypeDefinition(
		typeName,
		NamespaceBuiltin,
		TypeFlagBuiltIn,
	))
}

// NewTypeUsage creates a new type usage for the given type definition.
func NewTypeUsage(definition *CCLTypeDefinition) *CCLTypeUsage {
	return &CCLTypeUsage{
		definition: definition,
	}
}

// NewPointerTypeUsage creates a new pointer type usage that points to the given target type.
func NewPointerTypeUsage(targetType *CCLTypeUsage) *CCLTypeUsage {
	return &CCLTypeUsage{
		definition:     getPointerDefinition(),
		underlyingType: targetType,
	}
}

// NewArrayTypeUsage creates a new array type usage that holds elements of the given element type.
func NewArrayTypeUsage(elementType *CCLTypeUsage) *CCLTypeUsage {
	return &CCLTypeUsage{
		definition:     getArrayDefinition(),
		underlyingType: elementType,
	}
}
