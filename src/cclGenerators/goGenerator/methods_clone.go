package goGenerator

//---------------------------------------------------------

func (c *GoGenerationContext) generateCloneMethods() error {
	c.MethodsCode.ExpectMappedVars(
		"model",
		"modelName",
	)

	c.MethodsCode.NewLine().
		LineD("func (m $model) CloneEmpty() $model {").
		Indent().
		WriteLine("if m == nil {").
		Indent().
		WriteLine("return nil").
		Unindent().
		WriteLine("}").
		LineD("return &$modelName{}").
		Unindent().
		WriteLine("}")

	if c.Options.CCLDefinition.HasGlobalAttribute("AddSerializableInterface") {
		c.MethodsCode.NewLine().
			LineD("func (m $model) CloneEmptySerializable() Serializable {").
			Indent().
			WriteLine("if m == nil {").
			Indent().
			WriteLine("return nil").
			Unindent().
			WriteLine("}").
			LineD("return &$modelName{}").
			Unindent().
			WriteLine("}")
	}
	return nil
}

//---------------------------------------------------------
