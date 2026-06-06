package goGenerator

import "github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"

//---------------------------------------------------------

func (c *GoGenerationContext) generateCloneMethods(builder *codeBuilder.CodeBuilder) error {
	builder.ExpectMappedVars(
		"model",
		"modelName",
	)

	builder.NewLine().
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

	if c.GetGlobalAttribute(CurrentLanguage, "AddSerializableInterface") != nil {
		builder.NewLine().
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
