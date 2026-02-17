package cclSanitizer

// Error returns the string representation of the type usage resolution error.
func (e *TypeUsageResolutionError) Error() string {
	if e == nil {
		return "cclSanitizer: unknown type usage resolution error"
	}

	if e.SourcePosition == nil {
		return "cclSanitizer: " + e.Message
	}

	return e.SourcePosition.FormatError("Type error: " + e.Message)
}

// Error returns the string representation of the attribute resolution error.
func (e *AttributeResolutionError) Error() string {
	if e == nil {
		return "cclSanitizer: unknown attribute resolution error"
	}

	if e.SourcePosition == nil {
		return "cclSanitizer: " + e.Message
	}

	return e.SourcePosition.FormatError("Attribute error: " + e.Message)
}

// Error returns the string representation of the AST sanitization error.
func (e *AstSanitizationError) Error() string {
	if e == nil {
		return "cclSanitizer: unknown AST sanitization error"
	}

	if e.SourcePosition == nil {
		return "cclSanitizer: " + e.Message
	}

	return e.SourcePosition.FormatError("Sanitizer error: " + e.Message)
}
