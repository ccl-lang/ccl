package cclValues

func NewBuiltinTypeUsage(name string) *CCLTypeUsage {
	if !IsBuiltinTypeName(name) {
		return nil
	}

	return NewTypeUsage(NewTypeDefinition(
		&SimpleTypeName{
			TypeName:  name,
			Namespace: NamespaceBuiltin,
		},
		TypeFlagBuiltIn, // flags
		0,               // length
	))
}

// NewCustomTypeUsage returns a custom type usage for the given full type name.
func NewCustomTypeUsage(name *SimpleTypeName) *CCLTypeUsage {
	typeDefinitionsLock.Lock()
	defer typeDefinitionsLock.Unlock()

	typeDef := getTypeDefinition(name)
	if typeDef == nil {
		typeDef = getOrNewIncompleteTypeDef(name)
	}

	return NewTypeUsage(typeDef)
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
func NewArrayTypeUsage(elementType *CCLTypeUsage, arrayLength int) *CCLTypeUsage {
	return &CCLTypeUsage{
		definition:     getArrayDefinition(arrayLength),
		underlyingType: elementType,
	}
}
