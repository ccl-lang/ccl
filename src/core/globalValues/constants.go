package globalValues

const (
	CurrentCCLVersion = "v0.0.2"
)

const (
	LanguageUnknown LanguageType = 1 << iota
	LanguageCCL
	LanguageGo
	LanguageGd
	LanguageCS
	LanguagePy
	LanguageJS
	LanguageTS
)

const (
	DefaultMainNamespace = "main"
)
