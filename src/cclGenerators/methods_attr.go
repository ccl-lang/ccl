package cclGenerators

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

//---------------------------------------------------------

// GetGlobalAttribute retrieves a global attribute with the specified name.
func (c *CodeGenerationBase) GetGlobalAttribute(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
) *cclValues.AttributeUsageInfo {
	return c.getCodeContext().FindContextGlobalAttribute(targetLang, name)
}

// GetGlobalAttributes retrieves all global attributes with the specified name.
func (c *CodeGenerationBase) GetGlobalAttributes(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
) []*cclValues.AttributeUsageInfo {
	return c.getCodeContext().FindContextGlobalAttributes(targetLang, name)
}

// GetGlobalOrModelAttribute retrieves an attribute with the specified name
// from global attributes or the current model.
func (c *CodeGenerationBase) GetGlobalOrModelAttribute(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributeUsageInfo {
	return c.getCodeContext().ResolveAttribute(
		targetLang,
		name,
		&cclValues.AttributeResolutionSubject{
			Model: currentModel,
		},
		nil,
	)
}

// GetGlobalOrModelAttributes retrieves all attributes with the specified name
// from global attributes or the current model.
func (c *CodeGenerationBase) GetGlobalOrModelAttributes(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributesCollection {
	return cclValues.NewAttrsCollection(
		c.getCodeContext().ResolveAttributes(
			targetLang,
			name,
			&cclValues.AttributeResolutionSubject{
				Model: currentModel,
			},
			nil,
		),
	)
}

// GetGlobalAndModelAttributes retrieves all attributes with the specified name
// from both global attributes and the current model.
func (c *CodeGenerationBase) GetGlobalAndModelAttributes(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
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
	name cclAttr.CCLAttributeName,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributeUsageInfo {
	return c.getCodeContext().ResolveAttribute(
		targetLang,
		name,
		&cclValues.AttributeResolutionSubject{
			Model: currentModel,
		},
		nil,
	)
}

// GetModelAndGlobalAttributes retrieves all attributes with the specified name
// from both global attributes and the current model.
func (c *CodeGenerationBase) GetModelAndGlobalAttributes(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
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
	name cclAttr.CCLAttributeName,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributesCollection {
	return cclValues.NewAttrsCollection(
		c.getCodeContext().ResolveAttributes(
			targetLang,
			name,
			&cclValues.AttributeResolutionSubject{
				Model: currentModel,
			},
			nil,
		),
	)
}

// GetOutputFileGroup returns the generated output file group for a model.
func (c *CodeGenerationBase) GetOutputFileGroup(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
) (string, error) {
	attr := c.GetModelOrGlobalAttributes(
		targetLang,
		cclAttr.AttrOutputFileGroup,
		currentModel,
	).GetLast()
	if attr == nil {
		return "", nil
	}

	param := attr.GetParamAt(0)
	if param == nil || param.GetAsString() == "" {
		return "", &cclErrors.InvalidAttributeUsageError{
			AttrName:       attr.Name,
			Message:        "requires a non-empty string first-parameter",
			SourcePosition: attr.SourcePosition,
		}
	}

	group := param.GetAsString()
	if !isValidOutputFileGroup(group) {
		return "", &cclErrors.InvalidAttributeUsageError{
			AttrName: attr.Name,
			Message: " value '" + group +
				"' is not valid; only letters, digits, and underscores are allowed",
			SourcePosition: attr.SourcePosition,
		}
	}

	return group, nil
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
		cclAttr.AttrAddCloneMethods,
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
	name cclAttr.CCLAttributeName,
) *cclValues.AttributeUsageInfo {
	if field == nil {
		return nil
	}

	return c.getCodeContext().ResolveAttribute(
		targetLang,
		name,
		&cclValues.AttributeResolutionSubject{
			Field: field,
		},
		nil,
	)
}

// getCodeContext returns the current code context.
// from ccl version 0.0.4; having code context is strictly mandatory, hence
// all the code paths using this method can easily omit all the nil-checks on
// the returned pointer of this method, because it will be fine to panic if
// the returned pointer is nil, since that means something is *VERY* wrong with
// our current operation and it's better to be halted on panic instead of adding
// lots of hacks and unreliable fallbacks (and in most cases, those fallbacks will
// just bring code duplicates for us, which is why they shouldn't exist in multiple
// places at all).
func (c *CodeGenerationBase) getCodeContext() *cclValues.CCLCodeContext {
	return c.Options.CodeContext
}

//---------------------------------------------------------

// IsSingleFileMode checks if the code generation should be done in single file mode
// and returns the single file name if applicable.
func (c *CodeGenerationBase) IsSingleFileMode(targetLang gValues.LanguageType) (bool, string) {
	attr := c.GetGlobalAttribute(targetLang, cclAttr.AttrGenerateSingleFile)
	if attr != nil {
		return attr.GetParamAt(0).GetAsBool(), attr.GetParamAt(1).GetAsString()
	}
	return false, ""
}
