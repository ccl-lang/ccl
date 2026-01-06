package globalValues

const (
	CurrentCCLVersion = "v0.0.3"
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

// namespaces constants
const (
	DefaultMainNamespace = "main"
)
