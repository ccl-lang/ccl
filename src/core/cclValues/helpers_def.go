package cclValues

// GetTypeDefinition returns the type definition with the given full name from the cache.
func GetTypeDefinition(fullName string) *CCLTypeDefinition {
	typeDefinitionsLock.RLock()
	defer typeDefinitionsLock.RUnlock()

	if typeDef, exists := typeDefinitionsCache[fullName]; exists {
		return typeDef
	}

	return nil
}

// CacheTypeDefinition caches the given type definition with its full name.
func CacheTypeDefinition(typeDef *CCLTypeDefinition) {
	typeDefinitionsLock.Lock()
	defer typeDefinitionsLock.Unlock()

	typeDefinitionsCache[typeDef.GetFullName()] = typeDef
}

// NewTypeDefinition creates a new type info.
func NewTypeDefinition(name, namespace string, flags cclTypeFlag) *CCLTypeDefinition {
	if IsBuiltinTypeName(name) {
		if flags&TypeFlagBuiltIn == 0 {
			flags |= TypeFlagBuiltIn
		}

		namespace = NamespaceBuiltin
	}

	typeDef := GetTypeDefinition(namespace + "." + name)
	if typeDef != nil {
		// just return the cached type definition
		return typeDef
	}

	typeDef = &CCLTypeDefinition{
		name:      name,
		namespace: namespace,
		typeFlags: flags,
	}
	CacheTypeDefinition(typeDef)
	return typeDef
}

// getPointerDefinition returns the built-in pointer type definition.
// NOTE: If you want to set underlying type, you have to use a type usage.
func getPointerDefinition() *CCLTypeDefinition {
	return NewTypeDefinition(
		TypeNamePointer,
		NamespaceBuiltin,
		TypeFlagBuiltIn|TypeFlagPointer,
	)
}

func getArrayDefinition() *CCLTypeDefinition {
	return NewTypeDefinition(
		TypeNameArray,
		NamespaceBuiltin,
		TypeFlagBuiltIn|TypeFlagArray,
	)
}
