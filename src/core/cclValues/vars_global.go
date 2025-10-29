package cclValues

// All the global variables used in the source code.
var (
	// The map of all global variable definitions.
	// Key: variable name.
	// Value: variable definition.
	// Please do NOT use this map directly, use GetGlobalVariable instead.
	cclGlobalVariables = make(map[string]*VariableDefinition)

	cclAutomaticVariables = func() map[string]*VariableDefinition {
		return make(map[string]*VariableDefinition)
	}()
)
