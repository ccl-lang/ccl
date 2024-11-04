package cclValues

import "github.com/ALiwoto/ssg/ssg"

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

func (f *FieldDefinition) IsArray() bool {
	// TODO: handle this in a better way to have support for
	// more complex types like arrays, maps, etc.
	return f.ExtraOperators == "[]"
}

//---------------------------------------------------------
//---------------------------------------------------------
