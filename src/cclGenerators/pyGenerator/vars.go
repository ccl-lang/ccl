package pyGenerator

import gValues "github.com/ccl-lang/ccl/src/core/globalValues"

var (
	supportedStyles = []string{
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
