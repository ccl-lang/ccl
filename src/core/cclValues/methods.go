package cclValues

import (
	"fmt"

	"github.com/ALiwoto/ssg/ssg"
)

//---------------------------------------------------------

// GetNextModelId returns the next model ID.
func (d *SourceCodeDefinition) GetNextModelId() int64 {
	d.modelIdCounter++
	return d.modelIdCounter
}

// GetModelByName returns the model definition by the given name.
func (d *SourceCodeDefinition) GetModelByName(name string) *ModelDefinition {
	for _, model := range d.Models {
		if model.Name == name || model.DoesAliasMatch(name) {
			return model
		}
	}

	return nil
}

func (d *SourceCodeDefinition) IsCustomType(typeName string) bool {
	return d.GetModelByName(typeName) != nil
}

//---------------------------------------------------------

func (m *ModelDefinition) String() string {
	return "Model " + m.Name + " (ID: " + ssg.ToBase10(m.ModelId) + ")"
}

func (m *ModelDefinition) DoesAliasMatch(targetAlias string) bool {
	// TODO: implement model aliases in the future
	return false
}

func (m *ModelDefinition) GetFieldByName(name string) *FieldDefinition {
	for _, field := range m.Fields {
		if field.Name == name {
			return field
		}
	}

	return nil
}

//---------------------------------------------------------

// IsArray returns true if the field is an array.
// If the field does not have any type field assigned to it,
// it will result in a panic. So be careful before using this
// method.
func (f *FieldDefinition) IsArray() bool {
	return f.Type.IsArray()
}

// HasNoType returns true when the field's type field is not
// assigned to any value.
func (f *FieldDefinition) HasNoType() bool {
	return f.Type == nil
}

//---------------------------------------------------------

func (p *ParameterInstance) String() string {
	return fmt.Sprintf("Parameter %v (%s)", p.value, p.ValueType)
}

// ChangeValueType changes the value type of the parameter.
func (p *ParameterInstance) ChangeValueType(typeInfo *CCLTypeInfo) {
	p.ValueType = typeInfo
}

// ChangeValue sets the value of the parameter.
func (p *ParameterInstance) ChangeValue(value any) {
	p.value = value
}

// CompareValue compares the value of the parameter with the given value.
// It returns true if both values are equal, otherwise false.
func (p *ParameterInstance) CompareValue(value any) bool {
	if p.value == nil && value == nil {
		return true
	}

	if p.value == nil || value == nil {
		return false
	}

	return fmt.Sprintf("%v", p.value) == fmt.Sprintf("%v", value)
}

// HasBuiltInType returns true if the parameter is a built-in type.
func (p *ParameterInstance) IsBuiltInType() bool {
	return p.ValueType.IsBuiltIn()
}

// GetInt returns the integer value of the parameter.
// If the parameter is not an integer, it returns 0.
// Before using this method, it's highly recommended to get the value type
// and making sure the value of this parameter is in fact an integer.
func (p *ParameterInstance) GetInt() int {
	result, ok := p.value.(int)
	if !ok {
		return 0
	}

	return result
}

// GetString returns the string value of the parameter.
// If the parameter is not a string, it returns an empty string.
// Before using this method, it's highly recommended to get the value type
// and making sure the value of this parameter is in fact a string.
func (p *ParameterInstance) GetString() string {
	result, ok := p.value.(string)
	if !ok {
		return ""
	}

	return result
}

// GetValue returns the raw value of the parameter.
func (p *ParameterInstance) GetValue() any {
	return p.value
}

//---------------------------------------------------------

// IsBuiltIn returns true if the type is a built-in type.
func (t *CCLTypeInfo) IsBuiltIn() bool {
	return t.typeFlags&TypeFlagBuiltIn != 0
}

// IsArray returns true if the type is an array.
func (t *CCLTypeInfo) IsArray() bool {
	return t.typeFlags&TypeFlagArray != 0
}

// IsMap returns true if the type is a map.
func (t *CCLTypeInfo) IsMap() bool {
	return t.typeFlags&TypeFlagMap != 0
}

// String returns the string representation of the type info.
func (t *CCLTypeInfo) String() string {
	return t.name + " (flags: " + ssg.ToBase10(t.typeFlags) + ")"
}

// GetUnderlyingType returns the underlying type of the current type.
func (t *CCLTypeInfo) GetUnderlyingType() *CCLTypeInfo {
	return t.underlyingType
}

// GetName returns the name of the type.
func (t *CCLTypeInfo) GetName() string {
	return t.name
}

//---------------------------------------------------------

func (d *VariableDefinition) String() string {
	return d.Name + ": " + d.Type.GetName()
}

func (d *VariableDefinition) SetValue(value any) {
	d.value = value
}

func (d *VariableDefinition) GetValue() any {
	return d.value
}

func (d *VariableDefinition) IsAutomatic() bool {
	return d.isAutomaticVariable
}

//---------------------------------------------------------
