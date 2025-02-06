package gdGenerator

var (
	LanguageAliases = []string{
		"gd",
		"godot",
		"gdscript",
	}

	CCLTypesToGdTypes = map[string]string{
		"int":      "int",
		"uint":     "int",
		"int8":     "int",
		"uint8":    "int",
		"int16":    "int",
		"uint16":   "int",
		"int32":    "int",
		"uint32":   "int",
		"int64":    "int",
		"uint64":   "int",
		"float":    "float",
		"float64":  "float",
		"float32":  "float",
		"string":   "String",
		"bool":     "bool",
		"datetime": "int",
		"bytes":    "PackedByteArray",
	}
)
