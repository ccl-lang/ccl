package codeBuilder

import (
	"strings"
	"sync"
)

// NewCodeBuilder creates a new CodeBuilder instance with default options.
func NewCodeBuilder() *CodeBuilder {
	return NewCodeBuilderWithOptions(GetDefaultCodeBuilderOptions())
}

// NewCodeBuilderWithOptions creates a new CodeBuilder instance with the given options.
func NewCodeBuilderWithOptions(opts *CodeBuilderOptions) *CodeBuilder {
	return &CodeBuilder{
		mut: &sync.Mutex{},
		builders: map[string]*strings.Builder{
			SectionCommentHeaders: {},
			SectionHeaders:        {},
			SectionImports:        {},
		},
		indentations:   map[string]int{},
		importedKeys:   map[string]bool{},
		currentSection: "",
		indentationStr: opts.IndentationStr,
		newLineStr:     opts.NewLineStr,
	}
}

// GetDefaultCodeBuilderOptions returns the default options for CodeBuilder.
func GetDefaultCodeBuilderOptions() *CodeBuilderOptions {
	return &CodeBuilderOptions{
		IndentationStr: "\t",
		NewLineStr:     "\n",
	}
}

// GetDefaultOrderedSections returns the default ordered sections for output.
func GetDefaultOrderedSections() []string {
	return []string{SectionCommentHeaders, SectionHeaders, SectionImports}
}
