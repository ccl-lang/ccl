package cclValues

import "fmt"

// GetTypeDefinition returns the type definition with the given full name from the cache.
func GetTypeDefinition(name *SimpleTypeName) *CCLTypeDefinition {
	typeDefinitionsLock.RLock()
	defer typeDefinitionsLock.RUnlock()

	return getTypeDefinition(name)
}

// getTypeDefinition is the internal version of GetTypeDefinition.
// This function checks both complete and incomplete type definitions caches and
// it does NOT lock the mutex.
func getTypeDefinition(name *SimpleTypeName) *CCLTypeDefinition {
	// 1. check complete type definitions cache
	if typeDef, exists := typeDefinitionsCache[name.FullName()]; exists {
		return typeDef
	}

	// 2. check incomplete type definitions cache
	if typeDef, exists := incompleteTypeDefinitionsCache[name.FullName()]; exists {
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
func NewTypeDefinition(
	name *SimpleTypeName,
	flags cclTypeFlag,
	tLength int,
) *CCLTypeDefinition {
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
		length:    tLength,
	}
	cacheTypeDefinition(typeDef)
	return typeDef
}

// NewCustomTypeDefinition creates a new custom type definition.
func NewCustomTypeDefinition(
	name *SimpleTypeName,
	flags cclTypeFlag,
) (*CCLTypeDefinition, error) {
	typeDefinitionsLock.Lock()
	defer typeDefinitionsLock.Unlock()

	return newCustomTypeDefinition(name, flags)
}

// newCustomTypeDefinition is the internal version of NewCustomTypeDefinition
func newCustomTypeDefinition(
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

	typeDef := getTypeDefinition(name)
	if typeDef != nil {
		if typeDef.isIncomplete {
			return markAsCompleteTypeDefinition(name, typeDef, flags), nil
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
	cacheTypeDefinition(typeDef)
	return typeDef, nil
}

func NewModelTypeDefinition(
	name *SimpleTypeName,
	modelDef *ModelDefinition,
) (*CCLTypeDefinition, error) {
	typeDefinitionsLock.Lock()
	defer typeDefinitionsLock.Unlock()

	typeDef, err := newCustomTypeDefinition(name, TypeFlagCustomModel)
	if err != nil {
		return nil, err
	}

	typeDef.model = modelDef
	return typeDef, nil
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
		0, // length is unused for pointer types
	)
}

func getArrayDefinition(arrayLength int) *CCLTypeDefinition {
	return NewTypeDefinition(
		&SimpleTypeName{
			TypeName:  TypeNameArray,
			Namespace: NamespaceBuiltin,
		},
		TypeFlagBuiltIn|TypeFlagArray, // flags
		arrayLength,                   // length
	)
}
