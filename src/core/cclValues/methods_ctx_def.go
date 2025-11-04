package cclValues

import "fmt"

// GetTypeDefinition returns the type definition with the given full name from the cache.
func (c *CCLCodeContext) GetTypeDefinition(name *SimpleTypeName) *CCLTypeDefinition {
	c.typeDefinitionsLock.RLock()
	defer c.typeDefinitionsLock.RUnlock()

	return c.getTypeDefinition(name)
}

// getTypeDefinition is the internal version of GetTypeDefinition.
// This function checks both complete and incomplete type definitions caches and
// it does NOT lock the mutex.
func (c *CCLCodeContext) getTypeDefinition(name *SimpleTypeName) *CCLTypeDefinition {
	// 1. check complete type definitions cache
	if typeDef, exists := c.typeDefinitionsCache[name.FullName()]; exists {
		return typeDef
	}

	// 2. check incomplete type definitions cache
	if typeDef, exists := c.incompleteTypeDefinitionsCache[name.FullName()]; exists {
		return typeDef
	}

	return nil
}

// getOrNewIncompleteTypeDef returns the incomplete type definition with the given full name from the cache.
// If it does not exist, it creates a new incomplete type definition, caches it, and returns it.
func (c *CCLCodeContext) getOrNewIncompleteTypeDef(name *SimpleTypeName) *CCLTypeDefinition {
	if typeDef, exists := c.incompleteTypeDefinitionsCache[name.FullName()]; exists {
		return typeDef
	}

	typeDef := &CCLTypeDefinition{
		name:         name.TypeName,
		namespace:    name.Namespace,
		isIncomplete: true,
	}

	c.incompleteTypeDefinitionsCache[name.FullName()] = typeDef
	return typeDef
}

// getIncompleteTypeDefinition returns the incomplete type definition with the given full name from the cache.
// If it does not exist, it returns nil.
func (c *CCLCodeContext) getIncompleteTypeDefinition(name *SimpleTypeName) *CCLTypeDefinition {
	if typeDef, exists := c.incompleteTypeDefinitionsCache[name.FullName()]; exists {
		return typeDef
	}

	return nil
}

// deleteIncompleteTypeDefinition deletes the incomplete type definition with the given full name from the cache.
func (c *CCLCodeContext) deleteIncompleteTypeDefinition(name *SimpleTypeName) {
	delete(c.incompleteTypeDefinitionsCache, name.FullName())
}

// CacheTypeDefinition caches the given type definition with its full name.
func (c *CCLCodeContext) CacheTypeDefinition(typeDef *CCLTypeDefinition) {
	c.typeDefinitionsLock.Lock()
	defer c.typeDefinitionsLock.Unlock()

	c.cacheTypeDefinition(typeDef)
}

// cacheTypeDefinition is the internal version of CacheTypeDefinition
// that does NOT lock the mutex.
func (c *CCLCodeContext) cacheTypeDefinition(typeDef *CCLTypeDefinition) {
	c.typeDefinitionsCache[typeDef.GetFullName()] = typeDef
}

// NewTypeDefinition creates a new type info.
func (c *CCLCodeContext) NewTypeDefinition(
	name *SimpleTypeName,
	flags cclTypeFlag,
	tLength int,
) *CCLTypeDefinition {
	c.typeDefinitionsLock.Lock()
	defer c.typeDefinitionsLock.Unlock()

	var typeDef *CCLTypeDefinition

	if IsBuiltinTypeName(name.TypeName) {
		if flags&TypeFlagBuiltIn == 0 {
			flags |= TypeFlagBuiltIn
		}

		name.Namespace = NamespaceBuiltin
	} else {
		typeDef = c.getIncompleteTypeDefinition(name)
		if typeDef != nil {
			return c.markAsCompleteTypeDefinition(name, typeDef, flags)
		}
	}

	typeDef = c.getTypeDefinition(name)
	if typeDef != nil {
		// just return the cached type definition
		return typeDef
	}

	typeDef = &CCLTypeDefinition{
		name:      name.TypeName,
		namespace: name.Namespace,
		typeFlags: flags,
		length:    tLength,
	}
	c.cacheTypeDefinition(typeDef)
	return typeDef
}

// NewCustomTypeDefinition creates a new custom type definition.
func (c *CCLCodeContext) NewCustomTypeDefinition(
	name *SimpleTypeName,
	flags cclTypeFlag,
) (*CCLTypeDefinition, error) {
	c.typeDefinitionsLock.Lock()
	defer c.typeDefinitionsLock.Unlock()

	return c.newCustomTypeDefinition(name, flags)
}

// newCustomTypeDefinition is the internal version of NewCustomTypeDefinition
func (c *CCLCodeContext) newCustomTypeDefinition(
	name *SimpleTypeName,
	flags cclTypeFlag,
) (*CCLTypeDefinition, error) {
	if IsBuiltinTypeName(name.TypeName) {
		return nil, fmt.Errorf(
			StrErrCannotOverrideBuiltInType,
			name.TypeName,
			name.Namespace,
		)
	}

	typeDef := c.getTypeDefinition(name)
	if typeDef != nil {
		if typeDef.isIncomplete {
			return c.markAsCompleteTypeDefinition(name, typeDef, flags), nil
		}

		return nil, fmt.Errorf(
			StrErrTypeAlreadyDefined,
			name.TypeName,
			name.Namespace,
		)
	}

	typeDef = &CCLTypeDefinition{
		name:      name.TypeName,
		namespace: name.Namespace,
		typeFlags: flags,
		length:    0, // should custom types really be able to use this?
	}
	c.cacheTypeDefinition(typeDef)
	return typeDef, nil
}

// NewModelTypeDefinition creates a new model type definition.
func (c *CCLCodeContext) NewModelTypeDefinition(
	name *SimpleTypeName,
	modelDef *ModelDefinition,
) (*CCLTypeDefinition, error) {
	c.typeDefinitionsLock.Lock()
	defer c.typeDefinitionsLock.Unlock()

	typeDef, err := c.newCustomTypeDefinition(name, TypeFlagCustomModel)
	if err != nil {
		return nil, err
	}

	typeDef.model = modelDef
	return typeDef, nil
}

// markAsCompleteTypeDefinition marks the given incomplete type definition
// as complete and updates its flags.
// It also removes it from the incomplete type definitions cache
// and adds it to the complete type definitions cache.
func (c *CCLCodeContext) markAsCompleteTypeDefinition(
	name *SimpleTypeName,
	typeDef *CCLTypeDefinition,
	flags cclTypeFlag,
) *CCLTypeDefinition {
	// complete the incomplete type definition
	typeDef.isIncomplete = false
	typeDef.typeFlags = flags
	c.deleteIncompleteTypeDefinition(name)

	// cache the completed type definition
	c.cacheTypeDefinition(typeDef)

	return typeDef
}

// getPointerDefinition returns the built-in pointer type definition.
// NOTE: If you want to set underlying type, you have to use a type usage.
func (c *CCLCodeContext) getPointerDefinition() *CCLTypeDefinition {
	return c.NewTypeDefinition(
		&SimpleTypeName{
			TypeName:  TypeNamePointer,
			Namespace: NamespaceBuiltin,
		},
		TypeFlagBuiltIn|TypeFlagPointer,
		0, // length is unused for pointer types (at least for now?)
	)
}

// getArrayDefinition returns the built-in array type definition.
// NOTE: If you want to set underlying type, you have to use a type usage.
func (c *CCLCodeContext) getArrayDefinition(arrayLength int) *CCLTypeDefinition {
	return c.NewTypeDefinition(
		&SimpleTypeName{
			TypeName:  TypeNameArray,
			Namespace: NamespaceBuiltin,
		},
		TypeFlagBuiltIn|TypeFlagArray, // flags
		arrayLength,                   // length
	)
}
