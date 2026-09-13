package cclValues

// EnumMapping is a resolved directional conversion, with one case per source value.
type EnumMapping struct {
	Target   *EnumDefinition
	Fallback *EnumMemberDefinition
	Cases    []*EnumMappingCase
}

type EnumMappingCase struct {
	Source *EnumMemberDefinition
	Target *EnumMemberDefinition
}
