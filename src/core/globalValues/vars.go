package globalValues

var (
	langsToShortName = map[LanguageType]string{
		LanguageUnknown: "unknown",
		LanguageCCL:     "ccl",
		LanguageGd:      "gd",
		LanguageGo:      "go",
		LanguageCS:      "cs",
		LanguagePy:      "py",
		LanguageJS:      "js",
		LanguageTS:      "ts",
	}

	langsAliasNames = map[string]LanguageType{
		"unknown":    LanguageUnknown,
		"ccl":        LanguageCCL,
		"go":         LanguageGo,
		"golang":     LanguageGo,
		"cs":         LanguageCS,
		"csharp":     LanguageCS,
		"py":         LanguagePy,
		"python":     LanguagePy,
		"python3":    LanguagePy,
		"js":         LanguageJS,
		"javascript": LanguageJS,
		"ts":         LanguageTS,
		"typescript": LanguageTS,
	}
)
