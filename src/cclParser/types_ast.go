package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

// CCLAstParser parses tokens into a syntax-only AST.
type CCLAstParser struct {
	*CCLParser
}

func newCCLAstParser(
	ctx *cclValues.CCLCodeContext,
	tokens []*cclLexer.CCLToken,
	options *CCLParseOptions,
) *CCLAstParser {
	return &CCLAstParser{
		CCLParser: &CCLParser{
			Options: options,
			tokens:  tokens,
			ctx:     ctx,
		},
	}
}
