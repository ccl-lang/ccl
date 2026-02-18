package cclValues

import (
	"fmt"
	"slices"

	"github.com/ALiwoto/ssg/ssg"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

// GetNextModelId returns the next model ID.
func (d *SourceCodeDefinition) GetNextModelId() int64 {
	d.modelIdCounter++
	return d.modelIdCounter
}

// FindGlobalAttributes returns all global attributes with the given name.
func (d *SourceCodeDefinition) FindGlobalAttributes(
	targetLang gValues.LanguageType,
	name string,
) []*AttributeUsageInfo {
	attributes := []*AttributeUsageInfo{}

	for _, attr := range d.GlobalAttributes {
		if attr.Name == name && attr.IsForLanguage(targetLang) {
			attributes = append(attributes, attr)
		}
	}

	return attributes
}

// FindGlobalAttributes returns the first global attribute with the given name.
func (d *SourceCodeDefinition) FindGlobalAttribute(
	targetLang gValues.LanguageType,
	name string,
) *AttributeUsageInfo {
	for _, attr := range d.GlobalAttributes {
		if attr.Name == name && attr.IsForLanguage(targetLang) {
			return attr
		}
	}

	return nil
}

// GetModelByName returns the model definition by the given name.
func (d *SourceCodeDefinition) GetModelByName(name string) *ModelDefinition {
	for _, typeDef := range d.TypeDefinitions {
		if typeDef.IsCustomModel() {
			model := typeDef.GetModelDefinition()
			if model.Name == name || model.DoesAliasMatch(name) {
				return model
			}
		}
	}

	return nil
}

// GetAllModels returns all model definitions defined in the source code.
func (d *SourceCodeDefinition) GetAllModels() []*ModelDefinition {
	models := []*ModelDefinition{}
	for _, typeDef := range d.TypeDefinitions {
		if typeDef.IsCustomModel() {
			models = append(models, typeDef.GetModelDefinition())
		}
	}
	return models
}

// HasGlobalAttribute returns true if the source code definition
// has at least one of the given global attributes.
func (d *SourceCodeDefinition) HasGlobalAttribute(attributeName ...string) bool {
	for _, attr := range d.GlobalAttributes {
		if slices.Contains(attributeName, attr.Name) {
			return true
		}
	}

	return false
}

// IsCustomType returns true if the given type name is a custom type
// defined in the source code.
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

// GetNamespace returns a non-empty namespace for the model.
func (m *ModelDefinition) GetNamespace() string {
	if m.Namespace == "" {
		return gValues.DefaultMainNamespace
	}

	return m.Namespace
}

// GetFullName returns the full name of the model, including its namespace.
func (m *ModelDefinition) GetFullName() string {
	return m.GetNamespace() + "." + m.Name
}

// GetFullName returns the name of the model, NOT including its namespace.
func (m *ModelDefinition) GetName() string {
	return m.Name
}

// HasAttribute returns true if the model definition has at least one of the
// given attributes.
func (m *ModelDefinition) HasAttribute(attributeName ...string) bool {
	for _, attr := range m.Attributes {
		if slices.Contains(attributeName, attr.Name) {
			return true
		}
	}

	return false
}

// GetFieldByName returns the field definition by the given name.
func (m *ModelDefinition) GetFieldByName(name string) *ModelFieldDefinition {
	for _, field := range m.Fields {
		if field.Name == name {
			return field
		}
	}

	return nil
}

// FindAttributes returns all attributes with the given name.
func (m *ModelDefinition) FindAttributes(
	targetLang gValues.LanguageType,
	name string,
) []*AttributeUsageInfo {
	attributes := []*AttributeUsageInfo{}
	for _, attr := range m.Attributes {
		if attr.Name == name && attr.IsForLanguage(targetLang) {
			attributes = append(attributes, attr)
		}
	}

	return attributes
}

// FindAttribute returns the first attribute with the given name.
func (m *ModelDefinition) FindAttribute(name string) *AttributeUsageInfo {
	for _, attr := range m.Attributes {
		if attr.Name == name {
			return attr
		}
	}
	return nil
}

//---------------------------------------------------------

// IsArray returns true if the field is an array.
// If the field does not have any type field assigned to it,
// it will result in a panic. So be careful before using this
// method.
func (f *ModelFieldDefinition) IsArray() bool {
	return f.Type.IsArray()
}

// IsCustomTypeModel returns true if the field's type is a custom model.
func (f *ModelFieldDefinition) IsCustomTypeModel() bool {
	return f != nil && f.Type.IsCustomTypeModel()
}

// IsNullable returns true if the field's type is nullable.
func (c *ModelFieldDefinition) IsNullable() bool {
	// for now, we just return custom type models as nullable
	return c.IsCustomTypeModel()
}

// GetFullTypeName returns the full type name of the field's type.
func (f *ModelFieldDefinition) GetFullTypeName() string {
	return f.Type.GetDefinition().GetFullName()
}

// GetName returns the name of the field.
func (f *ModelFieldDefinition) GetName() string {
	return f.Name
}

// HasNoType returns true when the field's type field is not
// assigned to any value.
func (f *ModelFieldDefinition) HasNoType() bool {
	return f.Type == nil
}

//---------------------------------------------------------

func (p *ParameterInstance) String() string {
	return fmt.Sprintf("Parameter %v (%s)", p.value, p.ValueType)
}

// ChangeValueType changes the value type of the parameter.
func (p *ParameterInstance) ChangeValueType(typeUsage *CCLTypeUsage) {
	p.ValueType = typeUsage
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

// GetAsBool tries to convert the parameter value to a boolean.
// It supports bool, int, and string types for conversion.
// For int, any non-zero value is considered true.
// For string, "true" and "1" are considered true; all other values are false.
func (p *ParameterInstance) GetAsBool() bool {
	if p == nil || p.value == nil {
		return false
	}

	switch v := p.value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case string:
		return v == "true" || v == "1"
	default:
		return false
	}
}

// GetAsString tries to convert the parameter value to a string.
// It supports string types and types implementing the fmt.Stringer interface.
// For other types, it uses fmt.Sprintf to convert the value to a string.
func (p *ParameterInstance) GetAsString() string {
	if p == nil || p.value == nil {
		return ""
	}

	switch v := p.value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", p.value)
	}
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

func (d *VariableDefinition) String() string {
	return d.Name + ": " + d.Type.GetName()
}

func (d *VariableDefinition) SetValue(value any) {
	d.value = value
}

func (d *VariableDefinition) GetValue() any {
	return d.value
}

// HasImmutableType returns true if the variable's type is immutable.
// Immutable types are types that their value cannot be changed after
// being initialized.
// Such types include:
//   - Built-in types like int, float, string, bool
//   - User-defined models that are marked as immutable (if such feature is implemented)
func (d *VariableDefinition) HasImmutableType() bool {
	return d.Type.IsImmutable()
}

//---------------------------------------------------------

func (n *SimpleTypeName) FullName() string {
	return n.Namespace + "." + n.TypeName
}

//---------------------------------------------------------
