package cclValues

//---------------------------------------------------------

// String returns a string representation of the CCLTypeUsage.
func (u *CCLTypeUsage) String() string {
	return "Usage of: " + u.definition.GetName()
}

// IsBuiltIn returns true if the type usage is of a built-in type.
func (u *CCLTypeUsage) IsBuiltIn() bool {
	return u.definition.IsBuiltIn()
}

// GetDefinition returns the type definition of the type usage.
func (u *CCLTypeUsage) GetDefinition() *CCLTypeDefinition {
	return u.definition
}

// IsCustomTypeModel returns true if the type usage is of a custom model type.
func (u *CCLTypeUsage) IsCustomTypeModel() bool {
	return u != nil && u.definition != nil && u.definition.IsCustomModel()
}

// GetUnderlyingType returns the underlying type of the type usage.
func (u *CCLTypeUsage) GetUnderlyingType() *CCLTypeUsage {
	return u.underlyingType
}

// SetDefinition sets the type definition of the type usage.
func (u *CCLTypeUsage) SetDefinition(definition *CCLTypeDefinition) {
	u.definition = definition
}

// GetGenericArgs returns the generic arguments of the type usage.
func (u *CCLTypeUsage) GetGenericArgs() []*CCLTypeUsage {
	return u.genericArgs
}

// IsImmutable returns true if the type usage is of an immutable type.
func (u *CCLTypeUsage) IsImmutable() bool {
	return u.definition.IsImmutable()
}

// GetName returns the name of the type usage.
func (u *CCLTypeUsage) GetName() string {
	return u.definition.GetName()
}

// IsArray returns true if the type usage is of an array type.
func (u *CCLTypeUsage) IsArray() bool {
	return u.definition.IsArray()
}

//---------------------------------------------------------

// ChangeValueType changes the value type of the parameter.
func (p *ModelFieldDefinition) ChangeValueType(typeInfo *CCLTypeUsage) {
	p.Type = typeInfo
}

// ChangeDefaultValue sets the value of the parameter.
func (p *ModelFieldDefinition) ChangeDefaultValue(value any) {
	p.defaultValue = value
}

//---------------------------------------------------------
