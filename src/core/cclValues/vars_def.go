package cclValues

import "sync"

var (
	// typeDefinitionsCache is a cache for type definitions.
	// This is used to avoid creating multiple instances of the same type definition.
	// Key: full-name of the type definition (including namespace).
	// Value: type definition.
	typeDefinitionsCache = make(map[string]*CCLTypeDefinition)
	typeDefinitionsLock  = &sync.RWMutex{}
)
