package goGenerator

//---------------------------------------------------------

func (c *GoGenerationContext) generateCloneMethods(currentModel *CCLModel) error {
	c.MethodsCode.WriteString("\nfunc (m *" + currentModel.Name + ") CloneEmpty() *")
	c.MethodsCode.WriteString(currentModel.Name + " {\n")
	c.MethodsCode.WriteString("\tif m == nil {\n")
	c.MethodsCode.WriteString("\t\treturn nil\n")
	c.MethodsCode.WriteString("\t}\n")
	c.MethodsCode.WriteString("\treturn &" + currentModel.Name + "{}\n")
	c.MethodsCode.WriteString("}\n")

	if c.Options.CCLDefinition.HasGlobalAttribute("AddSerializableInterface") {
		c.MethodsCode.WriteString("\nfunc (m *" + currentModel.Name + ") CloneEmptySerializable() Serializable {\n")
		c.MethodsCode.WriteString("\tif m == nil {\n")
		c.MethodsCode.WriteString("\t\treturn nil\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\treturn &" + currentModel.Name + "{}\n")
		c.MethodsCode.WriteString("}\n")
	}
	return nil
}

//---------------------------------------------------------
