package goGenerator

//---------------------------------------------------------

func (c *GoGenerationContext) generateCloneMethods(currentModel *CCLModel) error {
	c.MethodsCode.NewLine().
		WriteLine("func (m *" + currentModel.Name + ") CloneEmpty() *" + currentModel.Name + " {").
		Indent().
		WriteLine("if m == nil {").
		Indent().
		WriteLine("return nil").
		Unindent().
		WriteLine("}").
		WriteLine("return &" + currentModel.Name + "{}").
		Unindent().
		WriteLine("}")

	if c.Options.CCLDefinition.HasGlobalAttribute("AddSerializableInterface") {
		c.MethodsCode.NewLine().
			WriteLine("func (m *" + currentModel.Name + ") CloneEmptySerializable() Serializable {").
			Indent().
			WriteLine("if m == nil {").
			Indent().
			WriteLine("return nil").
			Unindent().
			WriteLine("}").
			WriteLine("return &" + currentModel.Name + "{}").
			Unindent().
			WriteLine("}")
	}
	return nil
}

//---------------------------------------------------------
