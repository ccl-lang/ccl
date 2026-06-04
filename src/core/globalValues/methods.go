package globalValues

import (
	"strings"

	"github.com/ALiwoto/ssg/ssg/caseUtils"
)

//---------------------------------------------------------

// String returns the string representation of this language type.
func (l LanguageType) String() string {
	return l.GetShortName()
}

// GetShortName returns the short name of the current language
// type.
func (l LanguageType) GetShortName() NormalizedLangName {
	return langsToShortName[l]
}

// IsUnsupported returns true if the current language type is unsupported.
func (l LanguageType) IsUnsupported() bool {
	_, ok := langsToShortName[l]
	return !ok
}

//---------------------------------------------------------

// IsValid returns true if the current naming style is recognized as
// valid and supported by us.
func (n NamingStyle) IsValid() bool {
	switch n {
	case StyleCamelCase,
		StylePascalCase,
		StyleSnakeCase,
		StyleKebabCase:
		return true
	default:
		return false
	}
}

// ApplyStyle applies the specific naming style to the passed value.
func (n NamingStyle) ApplyStyle(value string) string {
	switch n {
	case StyleCamelCase:
		return caseUtils.ToCamelCase(value)
	case StylePascalCase:
		return caseUtils.ToPascalCase(value)
	case StyleSnakeCase:
		return caseUtils.ToSnakeCase(value)
	case StyleKebabCase:
		return strings.ReplaceAll(caseUtils.ToSnakeCase(value), "_", "-")
	default:
		return value
	}
}

// ToString returns the string representation of this naming style.
// please note that this is only for representing value, if you just
// simply want to convert this to a string type in a way that you can
// revert it later in future, prefer using string() instead.
func (n NamingStyle) ToString() string {
	return "NamingStyle(" + string(n) + ")"
}
