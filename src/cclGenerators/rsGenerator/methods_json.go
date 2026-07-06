package rsGenerator

import (
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
)

func (c *RustGenerationContext) generateModelJsonMethods(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
) {
	builder.WriteLine("impl " + model.Name + " {").
		Indent().
		WriteLine("pub fn serialize_json(&self) -> Result<String, serde_json::Error> {").
		Indent().
		WriteLine("serde_json::to_string(self)").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("pub fn deserialize_json(json_text: &str) -> Result<Self, serde_json::Error> {").
		Indent().
		WriteLine("serde_json::from_str(json_text)").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
}

func (c *RustGenerationContext) generateModelJsonAdapters(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
) {
	hasBytesField := false
	hasBytesArrayField := false
	for _, field := range model.Fields {
		if isRustBytesType(field.Type) {
			hasBytesField = true
		}
		if isRustBytesArrayType(field.Type) {
			hasBytesArrayField = true
		}
	}
	if hasBytesField {
		generateRustBytesJsonAdapter(builder, "ccl_bytes_json")
	}
	if hasBytesArrayField {
		generateRustBytesArrayJsonAdapter(builder, "ccl_bytes_array_json")
	}
}
