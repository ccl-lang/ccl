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
}

// CodeBuilderOptions holds configuration options for the CodeBuilder.
type CodeBuilderOptions struct {
	IndentationStr string
	NewLineStr     string
}
