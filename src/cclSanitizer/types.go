package cclSanitizer

import "github.com/ccl-lang/ccl/src/core/cclUtils"

// TypeUsageResolutionError represents a failure to resolve a type expression into a type usage.
type TypeUsageResolutionError struct {
	Message        string
	SourcePosition *cclUtils.SourceCodePosition
}
