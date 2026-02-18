package gdGenerator

import (
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclUtils"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

func (c *GDScriptGenerationContext) getJsonNamingStrategy(model *CCLModel) (string, error) {
	if model == nil {
		return "", nil
	}

	attr := model.FindAttribute("JsonPropertyNamingStrategy")
	if attr != nil && !attr.IsForLanguage(gValues.LanguageGd) {
		attr = nil
	}
	if attr == nil {
		attr = c.GetGlobalAttribute(gValues.LanguageGd, "JsonPropertyNamingStrategy")
	}
	if attr == nil {
		return "", nil
	}

	param := attr.GetParamAt(0)
	if param == nil || param.GetAsString() == "" {
		return "", &cclErrors.ValidationError{
			Message: "JsonPropertyNamingStrategy requires a non-empty string parameter",
		}
	}

	strategy := param.GetAsString()
	switch strategy {
	case gValues.StyleCamelCase,
		gValues.StylePascalCase,
		gValues.StyleSnakeCase,
		gValues.StyleKebabCase:
		return strategy, nil
	default:
		return "", &cclErrors.ValidationError{
			Message: "Unsupported JsonPropertyNamingStrategy value: " + strategy,
		}
	}
}

func (c *GDScriptGenerationContext) getJsonFieldName(
	model *CCLModel,
	field *CCLField,
) (string, error) {
	if field == nil {
		return "", &cclErrors.ValidationError{Message: "Field is nil when generating JSON name"}
	}

	attr := findFieldAttribute(field, gValues.LanguageGd, "JsonPropertyName")
	if attr != nil {
		param := attr.GetParamAt(0)
		if param == nil || param.GetAsString() == "" {
			return "", &cclErrors.ValidationError{
				Message: "JsonPropertyName requires a non-empty string parameter for field " + field.Name,
			}
		}
		return param.GetAsString(), nil
	}

	strategy, err := c.getJsonNamingStrategy(model)
	if err != nil {
		return "", err
	}
	if strategy == "" {
		return field.Name, nil
	}

	return applyJsonNamingStrategy(field.Name, strategy), nil
}

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
