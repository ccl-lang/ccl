package globalValues

const (
	CurrentCCLVersion = "v0.0.4"
)

const (
	LanguageUnknown LanguageType = iota
	LanguageCCL
	LanguageGo
	LanguageGd
	LanguageCS
	LanguagePy
	LanguageJS
	LanguageTS
)

// naming styles
const (
	StylePascalCase = "PascalCase"
	StyleCamelCase  = "camelCase"
	StyleSnakeCase  = "snake_case"
	StyleKebabCase  = "kebab-case"
)

// binary serialization endian values
const (
	EndianLittle = "little"
	EndianBig    = "big"
)

// namespaces constants
const (
	DefaultMainNamespace = "main"
)
