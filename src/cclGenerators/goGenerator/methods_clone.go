package goGenerator

//---------------------------------------------------------

func (c *GoGenerationContext) generateCloneMethods(currentModel *CCLModel) error {
	c.MethodsCode.NewLine()
	c.MethodsCode.WriteLine("func (m *" + currentModel.Name + ") CloneEmpty() *" + currentModel.Name + " {")
	c.MethodsCode.Indent()
	c.MethodsCode.WriteLine("if m == nil {")
	c.MethodsCode.Indent()
	c.MethodsCode.WriteLine("return nil")
	c.MethodsCode.Unindent()
	c.MethodsCode.WriteLine("}")
	c.MethodsCode.WriteLine("return &" + currentModel.Name + "{}")
	c.MethodsCode.Unindent()
	c.MethodsCode.WriteLine("}")

	if c.Options.CCLDefinition.HasGlobalAttribute("AddSerializableInterface") {
		c.MethodsCode.NewLine()
		c.MethodsCode.WriteLine("func (m *" + currentModel.Name + ") CloneEmptySerializable() Serializable {")
		c.MethodsCode.Indent()
		c.MethodsCode.WriteLine("if m == nil {")
		c.MethodsCode.Indent()
		c.MethodsCode.WriteLine("return nil")
		c.MethodsCode.Unindent()
		c.MethodsCode.WriteLine("}")
		c.MethodsCode.WriteLine("return &" + currentModel.Name + "{}")
		c.MethodsCode.Unindent()
		c.MethodsCode.WriteLine("}")
	}
	return nil
}

//---------------------------------------------------------
