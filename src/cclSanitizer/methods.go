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
