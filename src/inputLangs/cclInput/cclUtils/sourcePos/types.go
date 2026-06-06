package sourcePos

// SourceCodePointer is a struct that represents a certain position
// in the source code.
type SourceCodePosition struct {
	FilePath   string
	Line       int
	Column     int
	SourceLine string
}
