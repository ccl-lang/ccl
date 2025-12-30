package cclValues

import (
	"github.com/ALiwoto/ssg/ssg"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

// IsBuiltIn returns true if the type is a built-in type.
func (t *CCLTypeDefinition) IsBuiltIn() bool {
	return t.typeFlags&TypeFlagBuiltIn != 0
}

// IsArray returns true if the type is an array.
func (t *CCLTypeDefinition) IsArray() bool {
	return t.typeFlags&TypeFlagArray != 0
}

// IsMap returns true if the type is a map.
func (t *CCLTypeDefinition) IsMap() bool {
	return t.typeFlags&TypeFlagMap != 0
}

// String returns the string representation of the type info.
func (t *CCLTypeDefinition) String() string {
	return t.name + " (flags: " + ssg.ToBase10(t.typeFlags) + ")"
}

// GetName returns the name of the type.
func (t *CCLTypeDefinition) GetName() string {
	return t.name
}

// GetShortModelName returns the short name of the model type.
func (t *CCLTypeDefinition) GetShortModelName() string {
	if t.model != nil {
		return t.model.Name
	}

	return ""
}

func (t *CCLTypeDefinition) AddGenericParam(targetType *CCLTypeDefinition) error {
	if targetType == nil {
		return ErrGenericParamCantBeNil
	}

	// avoid entirely circular generic type
	if t == targetType {
		return ErrCircularGenericType
	}

	t.genericParams = append(t.genericParams, targetType)
	return nil
}

// GetNamespace returns a non-empty namespace for the type.
// If the type is built-in, it returns "builtin".
// If the type has no namespace, it returns "main".
// Otherwise, it returns the assigned namespace.
func (t *CCLTypeDefinition) GetNamespace() string {
	if t.IsBuiltIn() {
		return NamespaceBuiltin
	}

	if t.namespace == "" {
		return gValues.DefaultMainNamespace
	}

	return t.namespace
}

// GetFullName returns the full name of the type, including its namespace.
// This should be the one used for type comparisons and lookups.
func (t *CCLTypeDefinition) GetFullName() string {
	return t.GetNamespace() + "." + t.name
}

// GetModelDefinition returns the model definition of the type.
// It returns nil if the type is not a custom model type.
func (t *CCLTypeDefinition) GetModelDefinition() *ModelDefinition {
	return t.model
}

// IsCustomModel returns true if the type is a custom model.
// It does NOT necessarily check that model field is set.
func (t *CCLTypeDefinition) IsCustomModel() bool {
	return t != nil && t.typeFlags&TypeFlagCustomModel != 0
}

// IsCustomType simply calls IsCustomModel (for now).
func (t *CCLTypeDefinition) IsCustomType() bool {
	return t.IsCustomModel()
}

func (t *CCLTypeDefinition) HasModelField() bool {
	return t.model != nil
}

func (t *CCLTypeDefinition) IsImmutable() bool {
	return t.typeFlags&TypeFlagImmutable != 0
}

func (t *CCLTypeDefinition) IsIncomplete() bool {
	return t.isIncomplete
}

// GetLength returns the length field of the type definition.
// The meaning of this field depends on the type.
// For example, for array types, it represents the length of the array.
// If the type does not use length, it returns 0.
// If the length is dynamic or not supposed to be considered, it returns -1.
func (t *CCLTypeDefinition) GetLength() int {
	return t.length
}

//---------------------------------------------------------

func (v *VariableDefinition) IsAutomatic() bool {
	return v.isAutomaticVariable
}

//---------------------------------------------------------
