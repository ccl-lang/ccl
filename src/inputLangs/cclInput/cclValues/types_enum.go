package cclValues

import "github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils"

// EnumDefinition is a sanitized enum declaration.
type EnumDefinition struct {
	// SourceFileId is the source file where this enum is defined.
	SourceFileId SourceFileId

	// Name is the enum name.
	Name string

	// Namespace is the namespace where this enum type is registered.
	Namespace string

	// OwnedBy is set when this enum is declared inside a model.
	OwnedBy *ModelDefinition

	// BaseType is the integer type used to store enum values.
	BaseType *CCLTypeUsage

	// Members is the ordered list of enum members.
	Members []*EnumMemberDefinition

	// Attributes is an array of attribute definitions applied to this enum.
	Attributes []*AttributeUsageInfo

	// SourcePosition is the position of the enum in the source code.
	SourcePosition *cclUtils.SourceCodePosition
}

// EnumMemberDefinition is a sanitized enum member declaration.
type EnumMemberDefinition struct {
	// OwnedBy is the enum that owns this member.
	OwnedBy *EnumDefinition

	// Name is the enum member name.
	Name string

	// Value is the explicit integer value of this enum member.
	Value int64

	// SourcePosition is the position of the member in the source code.
	SourcePosition *cclUtils.SourceCodePosition
}

// EnumMemberReference is a field default value that references an enum member.
type EnumMemberReference struct {
	Enum   *EnumDefinition
	Member *EnumMemberDefinition
}
