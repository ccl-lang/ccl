package cclValues

// NewAttrsCollection returns collection
func NewAttrsCollection(attrs []*AttributeUsageInfo) *AttributesCollection {
	return &AttributesCollection{
		Attrs: attrs,
	}
}
