package rsGenerator

import "github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"

func generateRustBytesJsonAdapter(builder *codeBuilder.CodeBuilder, moduleName string) {
	builder.WriteLine("mod " + moduleName + " {").
		Indent().
		WriteLine("use base64::{engine::general_purpose, Engine as _};").
		WriteLine("use serde::{Deserialize, Deserializer, Serializer};").
		NewLine().
		WriteLine("pub fn serialize<S>(value: &Vec<u8>, serializer: S) -> Result<S::Ok, S::Error>").
		WriteLine("where").
		Indent().
		WriteLine("S: Serializer,").
		Unindent().
		WriteLine("{").
		Indent().
		WriteLine("serializer.serialize_str(&general_purpose::STANDARD.encode(value))").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("pub fn deserialize<'de, D>(deserializer: D) -> Result<Vec<u8>, D::Error>").
		WriteLine("where").
		Indent().
		WriteLine("D: Deserializer<'de>,").
		Unindent().
		WriteLine("{").
		Indent().
		WriteLine("let text = String::deserialize(deserializer)?;").
		WriteLine("general_purpose::STANDARD.decode(text).map_err(serde::de::Error::custom)").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
}

func generateRustBytesArrayJsonAdapter(builder *codeBuilder.CodeBuilder, moduleName string) {
	builder.WriteLine("mod " + moduleName + " {").
		Indent().
		WriteLine("use base64::{engine::general_purpose, Engine as _};").
		WriteLine("use serde::{Deserialize, Deserializer, Serialize, Serializer};").
		NewLine().
		WriteLine("pub fn serialize<S>(value: &Vec<Vec<u8>>, serializer: S) -> Result<S::Ok, S::Error>").
		WriteLine("where").
		Indent().
		WriteLine("S: Serializer,").
		Unindent().
		WriteLine("{").
		Indent().
		WriteLine("let encoded: Vec<String> = value.iter()").
		Indent().
		WriteLine(".map(|item| general_purpose::STANDARD.encode(item))").
		WriteLine(".collect();").
		Unindent().
		WriteLine("encoded.serialize(serializer)").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("pub fn deserialize<'de, D>(deserializer: D) -> Result<Vec<Vec<u8>>, D::Error>").
		WriteLine("where").
		Indent().
		WriteLine("D: Deserializer<'de>,").
		Unindent().
		WriteLine("{").
		Indent().
		WriteLine("let encoded = Vec::<String>::deserialize(deserializer)?;").
		WriteLine("encoded.into_iter()").
		Indent().
		WriteLine(".map(|item| general_purpose::STANDARD.decode(item).map_err(serde::de::Error::custom))").
		WriteLine(".collect()").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
}
