package cclGenerators

import (
	"slices"

	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclValues"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

// GetJsonNamingStrategy returns the json naming strategy that's applied to the current model.
func (c *CodeGenerationBase) GetJsonNamingStrategy(
	targetLang gValues.LanguageType,
	model *CCLModel,
) (gValues.NamingStyle, error) {
	if model == nil {
		return "", nil
	}

	attr := c.GetModelOrGlobalAttribute(targetLang, "JsonPropertyNamingStrategy", model)
	if attr == nil {
		return "", nil
	}

	param := attr.GetParamAt(0)
	if param == nil {
		return "", &cclErrors.ValidationError{
			Message: "JsonPropertyNamingStrategy requires a non-empty string parameter",
		}
	}

	strategy := gValues.NamingStyle(param.GetAsString())
	if strategy.IsValid() {
		return strategy, nil
	}

	return "", &cclErrors.ValidationError{
		Message: "Unsupported JsonPropertyNamingStrategy value: " + strategy.ToString(),
	}
}

// GetJsonFieldName returns the correct "field name" used for json serializing/deserializing
// a specific field from a specific model in a specific language.
func (c *CodeGenerationBase) GetJsonFieldName(
	targetLang gValues.LanguageType,
	model *CCLModel,
	field *CCLField,
) (string, error) {
	if field == nil {
		return "", &cclErrors.ValidationError{
			Message: "Field is nil when generating JSON name",
		}
	}

	attr := c.FindFieldAttribute(field, targetLang, "JsonPropertyName")
	if attr != nil {
		param := attr.GetParamAt(0)
		if param == nil || param.GetAsString() == "" {
			return "", &cclErrors.ValidationError{
				Message: "JsonPropertyName requires a non-empty string parameter for field " + field.Name,
			}
		}
		return param.GetAsString(), nil
	}

	strategy, err := c.GetJsonNamingStrategy(targetLang, model)
	if err != nil {
		return "", err
	}
	if strategy == "" {
		return field.Name, nil
	}

	return strategy.ApplyStyle(field.GetName()), nil
}

//---------------------------------------------------------

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
