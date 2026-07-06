package rsGenerator

import "github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"

func (c *RustGenerationContext) generateModelDeepClone(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
) {
	builder.WriteLine("impl " + model.Name + " {").
		Indent().
		WriteLine("pub fn deep_clone(&self) -> Self {").
		Indent().
		WriteLine("self.clone()").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
}
