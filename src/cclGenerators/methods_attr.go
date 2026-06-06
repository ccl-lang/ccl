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
	if ctx := c.getCodeContext(); ctx != nil {
		return ctx.FindContextGlobalAttribute(targetLang, name)
	}

	return nil
}

// GetGlobalAttributes retrieves all global attributes with the specified name.
func (c *CodeGenerationBase) GetGlobalAttributes(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
) []*cclValues.AttributeUsageInfo {
	if ctx := c.getCodeContext(); ctx != nil {
		return ctx.FindContextGlobalAttributes(targetLang, name)
	}

	return nil
}

// GetGlobalOrModelAttribute retrieves an attribute with the specified name
// from global attributes or the current model.
func (c *CodeGenerationBase) GetGlobalOrModelAttribute(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributeUsageInfo {
	if ctx := c.getCodeContext(); ctx != nil {
		return ctx.ResolveAttribute(
			targetLang,
			name,
			&cclValues.AttributeResolutionSubject{
				Model: currentModel,
			},
			nil,
		)
	}

	return c.GetModelOrGlobalAttribute(targetLang, name, currentModel)
}

// GetGlobalOrModelAttributes retrieves all attributes with the specified name
// from global attributes or the current model.
func (c *CodeGenerationBase) GetGlobalOrModelAttributes(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
	currentModel *cclValues.ModelDefinition,
) *cclValues.AttributesCollection {
	if ctx := c.getCodeContext(); ctx != nil {
		return cclValues.NewAttrsCollection(ctx.ResolveAttributes(
			targetLang,
			name,
			&cclValues.AttributeResolutionSubject{
				Model: currentModel,
			},
			nil,
		))
	}

	attrs := currentModel.FindAttributes(targetLang, name)
	if len(attrs) == 0 {
		attrs = c.GetGlobalAttributes(targetLang, name)
	}
	return cclValues.NewAttrsCollection(attrs)
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
	if ctx := c.getCodeContext(); ctx != nil {
		return ctx.ResolveAttribute(
			targetLang,
			name,
			&cclValues.AttributeResolutionSubject{
				Model: currentModel,
			},
			nil,
		)
	}

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
	if ctx := c.getCodeContext(); ctx != nil {
		return cclValues.NewAttrsCollection(ctx.ResolveAttributes(
			targetLang,
			name,
			&cclValues.AttributeResolutionSubject{
				Model: currentModel,
			},
			nil,
		))
	}

	attrs := currentModel.FindAttributes(targetLang, name)
	if len(attrs) == 0 {
		attrs = c.GetGlobalAttributes(targetLang, name)
	}
	return cclValues.NewAttrsCollection(attrs)
}

// GetOutputFileGroup returns the generated output file group for a model.
func (c *CodeGenerationBase) GetOutputFileGroup(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
) (string, error) {
	attrs := c.GetModelOrGlobalAttributes(
		targetLang,
		AttributeOutputFileGroup,
		currentModel,
	)
	if attrs.IsEmpty() {
		return "", nil
	}

	// wtf is this?
	// if len(attrs.Attrs) > 1 {
	// 	return "", &cclErrors.ValidationError{
	// 		Message: AttributeOutputFileGroup + " must be defined at most once for model " +
	// 			currentModel.GetFullName(),
	// 	}
	// }

	param := attrs.Attrs[0].GetParamAt(0)
	if param == nil || param.GetAsString() == "" {
		return "", &cclErrors.ValidationError{
			Message: AttributeOutputFileGroup + " requires a non-empty string parameter for model " +
				currentModel.GetFullName(),
		}
	}

	group := param.GetAsString()
	if !isValidOutputFileGroup(group) {
		return "", &cclErrors.ValidationError{
			Message: AttributeOutputFileGroup + " value '" + group +
				"' is not valid for model " + currentModel.GetFullName() +
				"; use only letters, digits, and underscores",
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
	name cclAttr.CCLAttributeName,
) *cclValues.AttributeUsageInfo {
	if field == nil {
		return nil
	}

	if ctx := c.getCodeContext(); ctx != nil {
		return ctx.ResolveAttribute(
			targetLang,
			name,
			&cclValues.AttributeResolutionSubject{
				Field: field,
			},
			nil,
		)
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

func (c *CodeGenerationBase) getCodeContext() *cclValues.CCLCodeContext {
	if c == nil || c.Options == nil {
		return nil
	}

	if c.Options.CodeContext != nil {
		return c.Options.CodeContext
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
