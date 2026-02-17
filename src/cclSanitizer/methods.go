package cclSanitizer

import (
	"fmt"
)

// Error returns the string representation of the type usage resolution error.
func (e *TypeUsageResolutionError) Error() string {
	if e == nil {
		return "cclSanitizer: unknown type usage resolution error"
	}

	if e.SourcePosition == nil {
		return "cclSanitizer: " + e.Message
	}

	return fmt.Sprintf(
		"cclSanitizer: %s (line %d, col %d)",
		e.Message,
		e.SourcePosition.Line,
		e.SourcePosition.Column,
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

	return fmt.Sprintf(
		"cclSanitizer: %s (line %d, col %d)",
		e.Message,
		e.SourcePosition.Line,
		e.SourcePosition.Column,
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

	return fmt.Sprintf(
		"cclSanitizer: %s (line %d, col %d)",
		e.Message,
		e.SourcePosition.Line,
		e.SourcePosition.Column,
	)
}
