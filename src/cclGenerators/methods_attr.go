package cclGenerators

import (
	"slices"

	"github.com/ccl-lang/ccl/src/core/cclErrors"
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

// NeedsBinarySerialization returns true if the current model or global attributes
// indicate that binary serialization is needed.
func (c *CodeGenerationBase) NeedsBinarySerialization(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
) bool {
	return c.NeedsSerializationType(
		targetLang,
		currentModel,
		"binary",
	)
}

// NeedsJsonSerialization returns true if the current model or global attributes
// indicate that JSON serialization is needed.
func (c *CodeGenerationBase) NeedsJsonSerialization(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
) bool {
	return c.NeedsSerializationType(
		targetLang,
		currentModel,
		"json",
	)
}

// NeedsSerializationType checks if the specified serialization type is needed
// based on global attributes and model-specific attributes.
func (c *CodeGenerationBase) NeedsSerializationType(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
	sType string,
) bool {
	collection := c.GetGlobalOrModelAttributes(
		targetLang,
		"SerializationType",
		currentModel,
	)
	return slices.Contains(collection.GetParamsAtAsStrings(0), sType)
}

// GetBinarySerializationEndian returns the requested endianness for binary
// serialization. Supported values are "big" and "small" (little-endian).
// Defaults to "small" to preserve backwards compatibility.
func (c *CodeGenerationBase) GetBinarySerializationEndian(
	targetLang gValues.LanguageType,
	currentModel *cclValues.ModelDefinition,
) (string, error) {
	attr := c.GetModelOrGlobalAttribute(targetLang, "BinarySerializationEndian", currentModel)
	if attr == nil {
		return gValues.EndianLittle, nil
	}

	param := attr.GetParamAt(0)
	if param == nil || param.GetAsString() == "" {
		return "", &cclErrors.ValidationError{
			Message: "BinarySerializationEndian requires a non-empty string parameter",
		}
	}

	endian := param.GetAsString()
	switch endian {
	case gValues.EndianBig:
		return endian, nil
	case gValues.EndianLittle:
		return endian, nil
	case "small": // considered alias to "little"
		return gValues.EndianLittle, nil
	default:
		return "", &cclErrors.ValidationError{
			Message: "Unsupported BinarySerializationEndian value: " + endian,
		}
	}
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
