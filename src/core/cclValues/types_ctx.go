package cclValues

import "sync"

type CCLCodeContext struct {
	// typeDefinitionsCache is a cache for type definitions.
	// This is used to avoid creating multiple instances of the same type definition.
	// Key: full-name of the type definition (including namespace).
	// Value: type definition.
	typeDefinitionsCache map[string]*CCLTypeDefinition
	typeDefinitionsLock  *sync.RWMutex

	// incompleteTypeDefinitionsCache is a cache for incomplete type definitions.
	// This is used to avoid creating multiple instances of the same incomplete
	// type definition each time it is referenced.
	// Key: full-name of the type definition (including namespace).
	// Value: type definition.
	// Later on, when the type definition is completed, it will be moved
	// to the typeDefinitionsCache and marked as complete.
	incompleteTypeDefinitionsCache map[string]*CCLTypeDefinition

	// The map of all global variable definitions.
	// Key: variable name.
	// Value: variable definition.
	// Please do NOT use this map directly, use GetGlobalVariable instead.
	cclGlobalVariables map[string]*VariableDefinition

	// cclAutomaticVariables is the map of all automatic variable definitions.
	// Key: variable name.
	// Value: variable definition.
	// Please do NOT use this map directly, use GetGlobalVariable instead.
	cclAutomaticVariables map[string]*VariableDefinition

	// globalVarsLock is the lock for global (and automatic) variables.
	globalVarsLock *sync.RWMutex
}
