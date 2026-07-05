package rsGenerator

var (
	LanguageAliases = []string{
		"rs",
		"rust",
	}

	CCLTypesToRustTypes = map[string]string{
		"int":      "i32",
		"uint":     "u32",
		"int8":     "i8",
		"uint8":    "u8",
		"int16":    "i16",
		"uint16":   "u16",
		"int32":    "i32",
		"uint32":   "u32",
		"int64":    "i64",
		"uint64":   "u64",
		"float":    "f64",
		"float64":  "f64",
		"float32":  "f32",
		"string":   "String",
		"bool":     "bool",
		"datetime": "i64",
		"bytes":    "Vec<u8>",
	}
)
