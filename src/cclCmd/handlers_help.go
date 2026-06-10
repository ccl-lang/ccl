package cclCmd

import "fmt"

func HandleHelpCommand() {
	fmt.Print(`ccl - Common Code Language code generator

CCL is a small IDL for defining models once and generating source code for
multiple target languages.

Usage:
  ccl help
  ccl --help
  ccl version
  ccl generate --source <file.ccl> --language <go|cs|gd|py|js|ts> --output <dir> [--generate-debug-info]
  ccl info --file <generated-file> --line <line>

Commands:
  generate    Generate code from a CCL source file.
  info        Map a generated-file line back to the source CCL line. Requires debug info.
  version     Show the installed CCL version.
  help        Show this help message.

Supported language values:
  go, cs, gd, py, js, ts
  Also accepted: golang, csharp, c#, godot, gdscript, python, python3,
  javascript, typescript.

CCL syntax basics:
  // Line comment
  import "other.ccl";
  namespace game.auth;

  #[CCLVersion("1.0.0")]
  #[SerializationType("binary")]
  #[SerializationType("json")]
  #namespace:[$go:OutputFileGroup("go_models")]

  model UserInfo {
      [JsonPropertyName("user_id")]
      Id: string;
      Username: string;
      Scores: int32[];
      ProfileImage: bytes;
      CreatedAt: datetime;
  }

Rules to remember:
  - A source file contains imports, optional namespace declarations, scoped
    attributes, and model definitions.
  - Model fields use "FieldName: Type;".
  - Built-in types include string, bytes, bool, int/int8/int16/int32/int64,
    uint/uint8/uint16/uint32/uint64, float/float32/float64, and datetime.
  - Arrays use Type[] or Type[10]. Custom model types use the model name.
  - CCL is intentionally minimal: enums, maps, unions, and generic types are not
    supported yet.
  - Local attributes use [Name(...)] before a model or field.
  - Global scoped attributes use #[Name(...)]; #[...] is global by default.
  - Scoped prefixes are #global:, #file:, and #namespace:.
  - Language selectors go inside the brackets before the attribute name, for
    example [$go:JsonPropertyName("id")] or
    #namespace:[$cs:OutputFileGroup("models")].

Useful attributes:
  - CCLVersion("1.0.0"): declares the minimum CCL version this source targets.
    Stable compilers must support sources targeting their version or older.
  - SerializationType("binary") and SerializationType("json"): enable binary
    or JSON serialization globally or on a model.
  - JsonPropertyName("name"): overrides the JSON name for a field.
  - OutputFileGroup("group_name"): routes generated files for supported
    generators. Use letters, digits, and underscores.
  - BinarySerializationEndian("little"|"big"): chooses binary byte order where
    supported.
  - AddCloneMethods(true): asks supported generators to emit clone helpers.
  - GenerateSingleFile(true, "filename"): asks supported generators to emit one
    file.

Examples:
  ccl generate --source definitions.ccl --language go --output ./gen/go
  ccl generate --source definitions.ccl --language ts --output ./gen/ts --generate-debug-info
  ccl info --file ./gen/go/models.go --line 42
`)
}
