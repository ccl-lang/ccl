package pyGenerator

import gValues "github.com/ccl-lang/ccl/src/core/globalValues"

var (
	supportedFileNameStyles = []gValues.NamingStyle{
		gValues.StylePascalCase,
		gValues.StyleSnakeCase,
		gValues.StyleCamelCase,
	}
)

var (
	LanguageAliases = []string{
		"py",
		"python",
		"python3",
	}
)
