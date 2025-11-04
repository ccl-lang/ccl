package cclValues

//---------------------------------------------------------

// GetGlobalVariable returns the global variable with the given name.
func (c *CCLCodeContext) GetGlobalVariable(name string) *VariableDefinition {
	c.globalVarsLock.RLock()
	defer c.globalVarsLock.RUnlock()

	return c.getGlobalVariable(name)
}

func (c *CCLCodeContext) getGlobalVariable(name string) *VariableDefinition {
	if variable, exists := c.cclAutomaticVariables[name]; exists {
		return variable
	}

	if variable, exists := c.cclGlobalVariables[name]; exists {
		return variable
	}

	return nil
}

//---------------------------------------------------------

func (c *CCLCodeContext) newStringAutomaticVariable(name string, value string) *VariableDefinition {
	return &VariableDefinition{
		Name:                name,
		value:               value,
		isAutomaticVariable: true,
		Type:                c.NewBuiltinTypeUsage(TypeNameString),
	}
}

//---------------------------------------------------------
