package csGenerator

import "regexp"

var (
	LanguageAliases = []string{
		"cs",
		"csharp",
		"c#",
	}
)

var (
	csNamespaceRegex = regexp.MustCompile(`namespace\s+([\w\.]+)`)
)
