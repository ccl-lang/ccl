package cclParser

import "github.com/ccl-lang/ccl/src/core/globalValues"

// ErrInvalidSyntax returns an invalid syntax error with an optional hint
// message (can be empty).
func (p *CCLParser) ErrInvalidSyntax(hint string) *InvalidSyntaxError {
	return &InvalidSyntaxError{
		Language:       globalValues.LanguageCCL,
		HintMessage:    hint,
		SourcePosition: p.getSourcePosition(),
	}
}
