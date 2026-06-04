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
	StylePascalCase NamingStyle = "PascalCase"
	StyleCamelCase  NamingStyle = "camelCase"
	StyleSnakeCase  NamingStyle = "snake_case"
	StyleKebabCase  NamingStyle = "kebab-case"
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
