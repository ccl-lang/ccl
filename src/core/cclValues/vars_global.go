package cclValues

import gValues "github.com/ccl-lang/ccl/src/core/globalValues"

// All the global variables used in the source code.
var (
	// The map of all global variable definitions.
	// Key: variable name.
	// Value: variable definition.
	// Please do NOT use this map directly, use GetGlobalVariable instead.
	cclGlobalVariables = make(map[string]*VariableDefinition)

	cclAutomaticVariables = map[string]*VariableDefinition{
		"__ccl_version": newStringAutomaticVariable("__ccl_version", gValues.CurrentCCLVersion),
	}
)
