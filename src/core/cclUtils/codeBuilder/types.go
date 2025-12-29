package codeBuilder

import "strings"

// CodeBuilder is a utility struct for building code strings with proper indentation
// and formatting.
// It provides methods to append lines, manage indentation levels, and is chainable.
type CodeBuilder struct {
	sb             strings.Builder
	indentation    int
	indentationStr string
}

// CodeBuilderOptions holds configuration options for the CodeBuilder.
type CodeBuilderOptions struct {
	IndentationStr string
}
