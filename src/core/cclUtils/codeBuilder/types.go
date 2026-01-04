package codeBuilder

import (
	"strings"
	"sync"
)

// CodeBuilder is a utility struct for building code strings with proper indentation
// and formatting.
// It provides methods to append lines, manage indentation levels, and is chainable.
type CodeBuilder struct {
	mut      *sync.Mutex
	builders map[string]*strings.Builder

	// indentations keeps track of the current indentation level for each section.
	indentations map[string]int

	// importedKeys is a map to track keys of imports, to avoid having duplicate imports.
	// The key is the import key, the value is always true.
	importedKeys map[string]bool

	currentSection string

	indentationStr string
	newLineStr     string

	// debugInfos keeps track of debug information for each section.
	// The key is the section name.
	debugInfos map[string][]*DebugInfo

	enableDebugInfo bool
}

// CodeBuilderOptions holds configuration options for the CodeBuilder.
type CodeBuilderOptions struct {
	IndentationStr  string
	NewLineStr      string
	EnableDebugInfo bool
}

// CodeBuildResult holds the result of the code building process.
type CodeBuildResult struct {
	Code      string
	DebugInfo string
}

// DebugInfo holds information about the source code that generated a specific part of the output code.
type DebugInfo struct {
	// SourceFile is the path to the source file that generated this code.
	SourceFile string `json:"source_file"`

	// SourceLine is the line number in the source file.
	SourceLine int `json:"source_line"`

	// GeneratedLine is the line number in the generated code (relative to the section start).
	// This will be adjusted during the final build process.
	GeneratedLine int `json:"generated_line"`

	// SectionOffset is the byte offset in the section's buffer where this debug info starts.
	SectionOffset int `json:"-"`
}
