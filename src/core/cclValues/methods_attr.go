package cclValues

import gValues "github.com/ccl-lang/ccl/src/core/globalValues"

//---------------------------------------------------------

// IsForLanguageStr checks if the attribute usage info is for the given language.
func (a *AttributeUsageInfo) IsForLanguageStr(lang string) bool {
	return a.IsForLanguage(gValues.GetLanguageTypeFromName(lang))
}

// IsForLanguage checks if the attribute usage info is for the given language.
func (a *AttributeUsageInfo) IsForLanguage(lang gValues.LanguageType) bool {
	if a.Language == 0 {
		// This attribute can be used for all target languages.
		return true
	} else if lang.IsUnsupported() {
		// The target language is unsupported, so we cannot match it.
		return false
	}

	return a.Language&lang != 0
}

// GetParamAt returns the parameter at the specified index.
func (a *AttributeUsageInfo) GetParamAt(index int) *ParameterInstance {
	if index < 0 || index >= len(a.Parameters) {
		return nil
	}

	return a.Parameters[index]
}

//---------------------------------------------------------

// GetParamsAt returns all parameters at the specified index from all attributes
// in the collection.
func (c *AttributesCollection) GetParamsAt(index int) []*ParameterInstance {
	var params []*ParameterInstance
	for _, attr := range c.Attrs {
		param := attr.GetParamAt(index)
		if param != nil {
			params = append(params, param)
		}
	}

	return params
}

// GetParamsAtAsStrings returns all parameters at the specified index from all attributes
// in the collection as strings.
func (c *AttributesCollection) GetParamsAtAsStrings(index int) []string {
	var params []string

	for _, attr := range c.Attrs {
		param := attr.GetParamAt(index)
		if param != nil {
			params = append(params, param.GetAsString())
		}
	}

	return params
}

//---------------------------------------------------------
