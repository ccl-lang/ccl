package cclValues

// GetTypeDefinition returns the type definition with the given full name from the cache.
func GetTypeDefinition(name *SimpleTypeName) *CCLTypeDefinition {
	typeDefinitionsLock.RLock()
	defer typeDefinitionsLock.RUnlock()

	return getTypeDefinition(name)
}

func getTypeDefinition(name *SimpleTypeName) *CCLTypeDefinition {
	if typeDef, exists := typeDefinitionsCache[name.FullName()]; exists {
		return typeDef
	}

	return nil
}

func getOrNewIncompleteTypeDef(name *SimpleTypeName) *CCLTypeDefinition {
	if typeDef, exists := incompleteTypeDefinitionsCache[name.FullName()]; exists {
		return typeDef
	}

	typeDef := &CCLTypeDefinition{
		name:         name.TypeName,
		namespace:    name.Namespace,
		isIncomplete: true,
	}

	incompleteTypeDefinitionsCache[name.FullName()] = typeDef
	return typeDef
}

func getIncompleteTypeDefinition(name *SimpleTypeName) *CCLTypeDefinition {
	if typeDef, exists := incompleteTypeDefinitionsCache[name.FullName()]; exists {
		return typeDef
	}

	return nil
}

func deleteIncompleteTypeDefinition(name *SimpleTypeName) {
	delete(incompleteTypeDefinitionsCache, name.FullName())
}

// CacheTypeDefinition caches the given type definition with its full name.
func CacheTypeDefinition(typeDef *CCLTypeDefinition) {
	typeDefinitionsLock.Lock()
	defer typeDefinitionsLock.Unlock()

	cacheTypeDefinition(typeDef)
}

func cacheTypeDefinition(typeDef *CCLTypeDefinition) {
	typeDefinitionsCache[typeDef.GetFullName()] = typeDef
}

// NewTypeDefinition creates a new type info.
func NewTypeDefinition(name *SimpleTypeName, flags cclTypeFlag) *CCLTypeDefinition {
	typeDefinitionsLock.Lock()
	defer typeDefinitionsLock.Unlock()

	var typeDef *CCLTypeDefinition

	if IsBuiltinTypeName(name.TypeName) {
		if flags&TypeFlagBuiltIn == 0 {
			flags |= TypeFlagBuiltIn
		}

		name.Namespace = NamespaceBuiltin
	} else {
		typeDef = getIncompleteTypeDefinition(name)
		if typeDef != nil {
			return markAsCompleteTypeDefinition(name, typeDef, flags)
		}
	}

	typeDef = getTypeDefinition(name)
	if typeDef != nil {
		// just return the cached type definition
		return typeDef
	}

	typeDef = &CCLTypeDefinition{
		name:      name.TypeName,
		namespace: name.Namespace,
		typeFlags: flags,
	}
	cacheTypeDefinition(typeDef)
	return typeDef
}

func markAsCompleteTypeDefinition(
	name *SimpleTypeName,
	typeDef *CCLTypeDefinition,
	flags cclTypeFlag,
) *CCLTypeDefinition {
	// complete the incomplete type definition
	typeDef.isIncomplete = false
	typeDef.typeFlags = flags
	deleteIncompleteTypeDefinition(name)

	// cache the completed type definition
	cacheTypeDefinition(typeDef)

	return typeDef
}

// getPointerDefinition returns the built-in pointer type definition.
// NOTE: If you want to set underlying type, you have to use a type usage.
func getPointerDefinition() *CCLTypeDefinition {
	return NewTypeDefinition(
		&SimpleTypeName{
			TypeName:  TypeNamePointer,
			Namespace: NamespaceBuiltin,
		},
		TypeFlagBuiltIn|TypeFlagPointer,
	)
}

func getArrayDefinition() *CCLTypeDefinition {
	return NewTypeDefinition(
		&SimpleTypeName{
			TypeName:  TypeNameArray,
			Namespace: NamespaceBuiltin,
		},
		TypeFlagBuiltIn|TypeFlagArray,
	)
}
