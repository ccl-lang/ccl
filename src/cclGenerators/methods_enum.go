package cclGenerators

import (
	"fmt"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

// GetGlobalOrEnumAttributes retrieves all attributes with the specified name
// from an enum or its scoped fallbacks.
func (c *CodeGenerationBase) GetGlobalOrEnumAttributes(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
	currentEnum *cclValues.EnumDefinition,
) *cclValues.AttributesCollection {
	return cclValues.NewAttrsCollection(
		c.getCodeContext().ResolveAttributes(
			targetLang,
			name,
			&cclValues.AttributeResolutionSubject{
				Enum: currentEnum,
			},
			nil,
		),
	)
}

// GetEnumMemberNamingStyle returns the member naming style for an enum.
func (c *CodeGenerationBase) GetEnumMemberNamingStyle(
	targetLang gValues.LanguageType,
	currentEnum *cclValues.EnumDefinition,
	defaultStyle gValues.NamingStyle,
) (gValues.NamingStyle, error) {
	attr := c.GetGlobalOrEnumAttributes(
		targetLang,
		cclAttr.AttrEnumMemberNamingStyle,
		currentEnum,
	).GetLast()
	if attr == nil {
		return defaultStyle, nil
	}

	param := attr.GetParamAt(0)
	if param == nil || param.GetAsString() == "" {
		return "", &cclErrors.InvalidAttributeUsageError{
			AttrName:       attr.Name,
			Message:        "requires a non-empty string first-parameter",
			SourcePosition: attr.SourcePosition,
		}
	}

	style := gValues.NamingStyle(param.GetAsString())
	if !style.IsValid() {
		return "", &cclErrors.UnsupportedFileNamingStyleError{
			StyleName:      param.GetAsString(),
			TargetLanguage: targetLang.String(),
			SourcePosition: attr.SourcePosition,
		}
	}

	return style, nil
}

// GetEnumStorageType returns the integer storage type for enum fields.
func (c *CodeGenerationBase) GetEnumStorageType(
	typeUsage *cclValues.CCLTypeUsage,
) *cclValues.CCLTypeUsage {
	if typeUsage == nil || !typeUsage.IsCustomTypeEnum() {
		return typeUsage
	}

	enumDef := typeUsage.GetDefinition().GetEnumDefinition()
	if enumDef == nil {
		return typeUsage
	}

	return enumDef.BaseType
}

// GetEnumDefaultReference returns the typed enum default reference if present.
func (c *CodeGenerationBase) GetEnumDefaultReference(
	field *cclValues.ModelFieldDefinition,
) *cclValues.EnumMemberReference {
	if field == nil {
		return nil
	}

	defaultValue := field.GetDefaultValue()
	ref, ok := defaultValue.(*cclValues.EnumMemberReference)
	if !ok {
		return nil
	}

	return ref
}

// FormatPrimitiveDefault formats primitive default values in a C-like syntax.
func (c *CodeGenerationBase) FormatPrimitiveDefault(value any) string {
	switch typedValue := value.(type) {
	case string:
		return fmt.Sprintf("%q", typedValue)
	case bool:
		if typedValue {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", typedValue)
	}
}
