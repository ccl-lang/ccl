package cclValues

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
)

// GetNamespace returns a non-empty namespace for the enum.
func (e *EnumDefinition) GetNamespace() string {
	if e.Namespace == "" {
		return gValues.DefaultMainNamespace
	}

	return e.Namespace
}

// GetFullName returns the full enum type name.
func (e *EnumDefinition) GetFullName() string {
	return e.GetNamespace() + "." + e.Name
}

// GetName returns the enum name.
func (e *EnumDefinition) GetName() string {
	return e.Name
}

// IsNested returns true when this enum is declared inside a model.
func (e *EnumDefinition) IsNested() bool {
	return e != nil && e.OwnedBy != nil
}

// GetMemberByName returns the enum member with the given name.
func (e *EnumDefinition) GetMemberByName(name string) *EnumMemberDefinition {
	for _, member := range e.Members {
		if member.Name == name {
			return member
		}
	}

	return nil
}

// FindAttributes returns all enum attributes with the given name.
func (e *EnumDefinition) FindAttributes(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
) []*AttributeUsageInfo {
	attributes := []*AttributeUsageInfo{}
	for _, attr := range e.Attributes {
		if attr.Name == name && attr.IsForLanguage(targetLang) {
			attributes = append(attributes, attr)
		}
	}

	return attributes
}
