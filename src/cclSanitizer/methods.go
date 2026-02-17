package cclSanitizer

import (
	"fmt"
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclUtils"
)

// Error returns the string representation of the type usage resolution error.
func (e *TypeUsageResolutionError) Error() string {
	if e == nil {
		return "cclSanitizer: unknown type usage resolution error"
	}

	if e.SourcePosition == nil {
		return "cclSanitizer: " + e.Message
	}

	return formatSanitizerErrorWithSourcePosition(
		"Type error: "+e.Message,
		e.SourcePosition,
	)
}

// Error returns the string representation of the attribute resolution error.
func (e *AttributeResolutionError) Error() string {
	if e == nil {
		return "cclSanitizer: unknown attribute resolution error"
	}

	if e.SourcePosition == nil {
		return "cclSanitizer: " + e.Message
	}

	return formatSanitizerErrorWithSourcePosition(
		"Attribute error: "+e.Message,
		e.SourcePosition,
	)
}

// Error returns the string representation of the AST sanitization error.
func (e *AstSanitizationError) Error() string {
	if e == nil {
		return "cclSanitizer: unknown AST sanitization error"
	}

	if e.SourcePosition == nil {
		return "cclSanitizer: " + e.Message
	}

	return formatSanitizerErrorWithSourcePosition(
		"Sanitizer error: "+e.Message,
		e.SourcePosition,
	)
}

func formatSanitizerErrorWithSourcePosition(
	message string,
	pos *cclUtils.SourceCodePosition,
) string {
	if pos == nil {
		return "cclSanitizer: " + message
	}

	if pos.SourceLine == "" {
		return fmt.Sprintf(
			"cclSanitizer: %s at line %d, column %d",
			message,
			pos.Line,
			pos.Column,
		)
	}

	result := fmt.Sprintf(
		"Error: %s\n  at line %d, column %d\n",
		message,
		pos.Line,
		pos.Column,
	)

	result += "  " + pos.SourceLine + "\n"
	pointerIndent := "  " + strings.Repeat(" ", pos.Column)
	result += pointerIndent + "^ " + message

	return result
}
