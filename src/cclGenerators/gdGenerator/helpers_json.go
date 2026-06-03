package gdGenerator

import (
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclUtils"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

func findFieldAttribute(
	field *CCLField,
	targetLang gValues.LanguageType,
	name string,
) *cclValues.AttributeUsageInfo {
	if field == nil {
		return nil
	}

	for _, attr := range field.Attributes {
		if attr.Name != name {
			continue
		}
		if attr.IsForLanguage(targetLang) {
			return attr
		}
	}

	return nil
}

func applyJsonNamingStrategy(fieldName string, strategy string) string {
	switch strategy {
	case gValues.StyleCamelCase:
		return cclUtils.ToCamelCase(fieldName)
	case gValues.StylePascalCase:
		return cclUtils.ToPascalCase(fieldName)
	case gValues.StyleSnakeCase:
		return cclUtils.ToSnakeCase(fieldName)
	case gValues.StyleKebabCase:
		return strings.ReplaceAll(cclUtils.ToSnakeCase(fieldName), "_", "-")
	default:
		return fieldName
	}
}
