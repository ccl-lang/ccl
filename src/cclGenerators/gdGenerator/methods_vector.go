package gdGenerator

import "github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"

func (c *GDScriptGenerationContext) generateVector3JsonValidation(
	builder *codeBuilder.CodeBuilder,
	value string,
) {
	builder.MapVarPairs("vectorValue", value)
	defer builder.UnmapVar("vectorValue")
	builder.LineD("if not ($vectorValue is Array) or $vectorValue.size() != 3:").
		Indent().
		LineD(`push_error("Expected three numbers for Vector3 field $field in $model")`).
		WriteLine("return null").
		Unindent().
		LineD("for component in $vectorValue:").
		Indent().
		WriteLine("if not (component is int or component is float):").
		Indent().
		LineD(`push_error("Expected numeric Vector3 components for field $field in $model")`).
		WriteLine("return null").
		Unindent().Unindent()
}
