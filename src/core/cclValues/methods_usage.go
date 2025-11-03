package cclValues

//---------------------------------------------------------

func (u *CCLTypeUsage) String() string {
	return "Usage of: " + u.definition.GetName()
}

func (u *CCLTypeUsage) IsBuiltIn() bool {
	return u.definition.IsBuiltIn()
}

func (u *CCLTypeUsage) GetDefinition() *CCLTypeDefinition {
	return u.definition
}

func (u *CCLTypeUsage) SetDefinition(definition *CCLTypeDefinition) {
	u.definition = definition
}

func (u *CCLTypeUsage) GetGenericArgs() []*CCLTypeUsage {
	return u.genericArgs
}

func (u *CCLTypeUsage) IsImmutable() bool {
	return u.definition.IsImmutable()
}

func (u *CCLTypeUsage) GetName() string {
	return u.definition.GetName()
}

func (u CCLTypeUsage) IsArray() bool {
	return u.definition.IsArray()
}

//---------------------------------------------------------

// ChangeValueType changes the value type of the parameter.
func (p *ModelFieldDefinition) ChangeValueType(typeInfo *CCLTypeUsage) {
	p.Type = typeInfo
}

// ChangeValue sets the value of the parameter.
func (p *ModelFieldDefinition) ChangeValue(value any) {
	p.value = value
}

//---------------------------------------------------------
