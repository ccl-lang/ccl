package cclSanitizer

import "github.com/ccl-lang/ccl/src/core/cclUtils"

// TypeUsageResolutionError represents a failure to resolve a type expression into a type usage.
type TypeUsageResolutionError struct {
	Message        string
	SourcePosition *cclUtils.SourceCodePosition
}

// AttributeResolutionError represents a failure to resolve an attribute usage.
type AttributeResolutionError struct {
	Message        string
	SourcePosition *cclUtils.SourceCodePosition
}

// AstSanitizationError represents a general failure while converting AST to IR.
type AstSanitizationError struct {
	Message        string
	SourcePosition *cclUtils.SourceCodePosition
}
