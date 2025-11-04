package cclValues

// NewBuiltinTypeUsage returns a built-in type usage for the given type name
// in the current code context. If the name is not a built-in type name, it returns nil.
func (c *CCLCodeContext) NewBuiltinTypeUsage(name string) *CCLTypeUsage {
	if !IsBuiltinTypeName(name) {
		return nil
	}

	return NewTypeUsage(c.NewTypeDefinition(
		&SimpleTypeName{
			TypeName:  name,
			Namespace: NamespaceBuiltin,
		},
		TypeFlagBuiltIn, // flags
		0,               // length
	))
}

// NewCustomTypeUsage returns a custom type usage for the given full type name.
func (c *CCLCodeContext) NewCustomTypeUsage(name *SimpleTypeName) *CCLTypeUsage {
	c.typeDefinitionsLock.Lock()
	defer c.typeDefinitionsLock.Unlock()

	typeDef := c.getTypeDefinition(name)
	if typeDef == nil {
		typeDef = c.getOrNewIncompleteTypeDef(name)
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
func (c *CCLCodeContext) NewPointerTypeUsage(targetType *CCLTypeUsage) *CCLTypeUsage {
	return &CCLTypeUsage{
		definition:     c.getPointerDefinition(),
		underlyingType: targetType,
	}
}

// NewArrayTypeUsage creates a new array type usage that holds elements of the given element type.
func (c *CCLCodeContext) NewArrayTypeUsage(elementType *CCLTypeUsage, arrayLength int) *CCLTypeUsage {
	return &CCLTypeUsage{
		definition:     c.getArrayDefinition(arrayLength),
		underlyingType: elementType,
	}
}
