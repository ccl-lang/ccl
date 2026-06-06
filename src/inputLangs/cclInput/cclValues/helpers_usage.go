package cclValues

// NewTypeUsage creates a new type usage for the given type definition.
func NewTypeUsage(definition *CCLTypeDefinition) *CCLTypeUsage {
	return &CCLTypeUsage{
		definition: definition,
	}
}
