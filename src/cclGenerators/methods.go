package cclGenerators

import (
	"github.com/ALiwoto/ssg/ssg"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

// WriteCodeFile writes the generated code and debug info (if available) to the specified path.
func (c *CodeGenerationBase) WriteCodeFile(path string, result *codeBuilder.CodeBuildResult) error {
	err := ssg.WriteFileStr(path, result.Code)
	if err != nil {
		return err
	}

	if result.DebugInfo != "" {
		return ssg.WriteFileStr(path+".cclinfo", result.DebugInfo)
	}
	return nil
}

//---------------------------------------------------------

func (c *CodeGenerationBase) GetFileNamingStyle(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
	defaultStyle gValues.NamingStyle,
) gValues.NamingStyle {
	collection := c.GetGlobalOrModelAttributes(
		targetLang,
		"FileNamingStyle",
		currentModel,
	)
	if collection.IsEmpty() {
		return defaultStyle
	}

	return gValues.NamingStyle(collection.GetParamsAtAsStrings(0)[0])
}

// GetFileNameForModel returns the file name for the given model based on the naming style.
// NOTE: this is only the file name, it does not apply any file extension or base path.
func (c *CodeGenerationBase) GetFileNameForModel(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
	defaultStyle gValues.NamingStyle,
	supportedStyles []gValues.NamingStyle,
) (string, error) {
	namingStyle := c.GetFileNamingStyle(targetLang, currentModel, defaultStyle)
	fileName := currentModel.GetName()
	if !namingStyle.IsValid() {
		return "", &cclErrors.UnsupportedFileNamingStyleError{
			ModelName:       currentModel.GetFullName(),
			StyleName:       string(namingStyle),
			SupportedStyles: ssg.JoinStr(supportedStyles, ","),
			TargetLanguage:  targetLang.String(),
		}
	}

	return namingStyle.ApplyStyle(fileName), nil
}

//---------------------------------------------------------
