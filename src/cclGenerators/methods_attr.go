package cclGenerators

import (
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

// GetGlobalAttribute retrieves a global attribute with the specified name.
func (c *CodeGenerationBase) GetGlobalAttribute(
	targetLang gValues.LanguageType,
	name string,
) *cclValues.AttributeUsageInfo {
	return c.Options.CCLDefinition.FindGlobalAttribute(targetLang, name)
}

// GetGlobalAttributes retrieves all global attributes with the specified name.
func (c *CodeGenerationBase) GetGlobalAttributes(
	targetLang gValues.LanguageType,
	name string,
) []*cclValues.AttributeUsageInfo {
	return c.Options.CCLDefinition.FindGlobalAttributes(targetLang, name)
}

// GetGlobalOrModelAttribute retrieves an attribute with the specified name
// from global attributes or the current model.
func (c *CodeGenerationBase) GetGlobalOrModelAttribute(
	targetLang gValues.LanguageType,
	name string,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributeUsageInfo {
	attr := c.GetGlobalAttribute(targetLang, name)
	if attr == nil {
		attr = currentModel.FindAttribute(name)
	}

	return attr
}

// GetGlobalOrModelAttributes retrieves all attributes with the specified name
// from global attributes or the current model.
func (c *CodeGenerationBase) GetGlobalOrModelAttributes(
	targetLang gValues.LanguageType,
	name string,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributesCollection {
	attrs := c.GetGlobalAttributes(targetLang, name)
	if len(attrs) == 0 {
		attrs = currentModel.FindAttributes(targetLang, name)
	}
	return cclValues.NewAttrsCollection(attrs)
}

// GetGlobalAndModelAttributes retrieves all attributes with the specified name
// from both global attributes and the current model.
func (c *CodeGenerationBase) GetGlobalAndModelAttributes(
	targetLang gValues.LanguageType,
	name string,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributesCollection {
	attrs := c.GetGlobalAttributes(targetLang, name)
	attrs = append(attrs, currentModel.FindAttributes(targetLang, name)...)
	return cclValues.NewAttrsCollection(attrs)
}

// GetModelOrGlobalAttribute retrieves an attribute with the specified name
// from the current model or global attributes.
func (c *CodeGenerationBase) GetModelOrGlobalAttribute(
	targetLang gValues.LanguageType,
	name string,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributeUsageInfo {
	attr := currentModel.FindAttribute(name)
	if attr == nil {
		attr = c.GetGlobalAttribute(targetLang, name)
	}

	return attr
}

// GetModelAndGlobalAttributes retrieves all attributes with the specified name
// from both global attributes and the current model.
func (c *CodeGenerationBase) GetModelAndGlobalAttributes(
	targetLang gValues.LanguageType,
	name string,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributesCollection {
	attrs := c.GetGlobalAttributes(targetLang, name)
	attrs = append(attrs, currentModel.FindAttributes(targetLang, name)...)
	return cclValues.NewAttrsCollection(attrs)
}

// GetModelOrGlobalAttributes retrieves all attributes with the specified name
// from the current model or global attributes.
func (c *CodeGenerationBase) GetModelOrGlobalAttributes(
	targetLang gValues.LanguageType,
	name string,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributesCollection {
	attrs := currentModel.FindAttributes(targetLang, name)
	if len(attrs) == 0 {
		attrs = c.GetGlobalAttributes(targetLang, name)
	}
	return cclValues.NewAttrsCollection(attrs)
}

//---------------------------------------------------------

// NeedsCloneMethods returns true if the current model or global attributes
// indicate that clone methods are needed.
func (c *CodeGenerationBase) NeedsCloneMethods(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
) bool {
	attr := c.GetModelOrGlobalAttribute(
		targetLang,
		"AddCloneMethods",
		currentModel,
	)
	if attr != nil {
		return attr.GetParamAt(0).GetAsBool()
	}
	return false
}

//---------------------------------------------------------

func (c *CodeGenerationBase) FindFieldAttribute(
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

//---------------------------------------------------------

// IsSingleFileMode checks if the code generation should be done in single file mode
// and returns the single file name if applicable.
func (c *CodeGenerationBase) IsSingleFileMode(targetLang gValues.LanguageType) (bool, string) {
	attr := c.GetGlobalAttribute(targetLang, "CCLGenerateSingleFile")
	if attr != nil {
		return attr.GetParamAt(0).GetAsBool(), attr.GetParamAt(1).GetAsString()
	}
	return false, ""
}
