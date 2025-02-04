package goGenerator

var (
	// LanguageAliases is a list of aliases for the Go language.
	LanguageAliases = []string{
		"go",
		"golang",
	}

	CCLTypesToGoTypes = map[string]string{
		"int":      "int",
		"uint":     "uint",
		"int8":     "int8",
		"uint8":    "uint8",
		"int16":    "int16",
		"uint16":   "uint16",
		"int32":    "int32",
		"uint32":   "uint32",
		"int64":    "int64",
		"uint64":   "uint64",
		"float":    "float64",
		"float64":  "float64",
		"float32":  "float32",
		"string":   "string",
		"bool":     "bool",
		"datetime": "time.Time",
		"bytes":    "[]byte",
	}
)
