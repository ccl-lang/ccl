package cclParser

import (
	"fmt"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils"
)

// ImportResolutionError describes a failure while resolving CCL import declarations.
type ImportResolutionError struct {
	ImportPath     string
	ResolvedPath   string
	Message        string
	SourcePosition *cclUtils.SourceCodePosition
	InnerError     error
}

func (e *ImportResolutionError) Error() string {
	if e == nil {
		return "cclParser: import resolution failed"
	}

	message := e.Message
	if message == "" {
		message = "import resolution failed"
	}

	if e.ImportPath != "" {
		message += fmt.Sprintf(" for %q", e.ImportPath)
	}

	if e.ResolvedPath != "" {
		message += " resolved to " + e.ResolvedPath
	}

	if e.InnerError != nil {
		message += ": " + e.InnerError.Error()
	}

	if e.SourcePosition == nil {
		return "cclParser: " + message
	}

	if e.SourcePosition.FilePath != "" {
		message = e.SourcePosition.FilePath + ": " + message
	}

	return e.SourcePosition.FormatError(message)
}

func (e *ImportResolutionError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.InnerError
}
