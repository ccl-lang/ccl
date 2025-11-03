package cclValues

//---------------------------------------------------------

// GetGlobalVariable returns the global variable with the given name.
func GetGlobalVariable(name string) *VariableDefinition {
	if variable, exists := cclAutomaticVariables[name]; exists {
		return variable
	}

	if variable, exists := cclGlobalVariables[name]; exists {
		return variable
	}

	return nil
}

//---------------------------------------------------------

func newStringAutomaticVariable(name string, value string) *VariableDefinition {
	return &VariableDefinition{
		Name:                name,
		value:               value,
		isAutomaticVariable: true,
		Type:                NewBuiltinTypeUsage(TypeNameString),
	}
}

//---------------------------------------------------------
