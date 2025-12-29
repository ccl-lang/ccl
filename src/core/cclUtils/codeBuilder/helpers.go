package codeBuilder

// NewCodeBuilder creates a new CodeBuilder instance with default options.
func NewCodeBuilder() *CodeBuilder {
	return NewCodeBuilderWithOptions(GetDefaultCodeBuilderOptions())
}

// NewCodeBuilderWithOptions creates a new CodeBuilder instance with the given options.
func NewCodeBuilderWithOptions(opts *CodeBuilderOptions) *CodeBuilder {
	return &CodeBuilder{
		indentationStr: opts.IndentationStr,
	}
}

// GetDefaultCodeBuilderOptions returns the default options for CodeBuilder.
func GetDefaultCodeBuilderOptions() *CodeBuilderOptions {
	return &CodeBuilderOptions{
		IndentationStr: "\t",
	}
}
