package cclGenerators

import (
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclUtils"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

// IsCustomType returns true if the given CCL type is a custom type.
func (c *CodeGenerationBase) IsCustomType(cclType string) bool {
	return c.Options.CCLDefinition.IsCustomType(cclType)
}

//---------------------------------------------------------

func (c *CodeGenerationBase) GetFileNamingStyle(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
	defaultStyle string,
) string {
	collection := c.GetGlobalOrModelAttributes(
		targetLang,
		"FileNamingStyle",
		currentModel,
	)
	if collection.IsEmpty() {
		return defaultStyle
	}
	return collection.GetParamsAtAsStrings(0)[0]
}

// GetFileNameForModel returns the file name for the given model based on the naming style.
// NOTE: this is only the file name, it does not apply any file extension or base path.
func (c *CodeGenerationBase) GetFileNameForModel(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
	defaultStyle string,
	supportedStyles []string,
) (string, error) {
	namingStyle := c.GetFileNamingStyle(gValues.LanguagePy, currentModel, defaultStyle)
	fileName := ""
	switch namingStyle {
	case gValues.StylePascalCase:
		fileName = cclUtils.ToPascalCase(currentModel.Name)
	case gValues.StyleSnakeCase:
		fileName = cclUtils.ToSnakeCase(currentModel.Name)
	case gValues.StyleCamelCase:
		fileName = cclUtils.ToCamelCase(currentModel.Name)
	default:
		return "", &cclErrors.UnsupportedFileNamingStyleError{
			ModelName:       currentModel.GetFullName(),
			StyleName:       namingStyle,
			SupportedStyles: supportedStyles,
			TargetLanguage:  targetLang.String(),
		}
	}
	return fileName, nil
}

//---------------------------------------------------------
