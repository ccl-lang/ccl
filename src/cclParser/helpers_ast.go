package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
)

// ParseCCLSourceFileAsAST reads a CCL source file, resolves imports, and parses
// the combined source graph into a CCLFileAST.
func ParseCCLSourceFileAsAST(options *CCLParseOptions) (*cclAst.CCLFileAST, error) {
	return newImportGraphResolver().parseSourceFileAsAST(options)
}

// ParseCCLSourceContentAsAST parses standalone source content into a CCLFileAST
// (syntax only).
func ParseCCLSourceContentAsAST(options *CCLParseOptions) (*cclAst.CCLFileAST, error) {
	allTokens, err := cclLexer.Lex(options.SourceContent)
	if err != nil {
		return nil, err
	}

	return ParseCCLAst(allTokens, options)
}

func ParseCCLAst(
	tokens []*cclLexer.CCLToken,
	options *CCLParseOptions,
) (*cclAst.CCLFileAST, error) {
	theParser := newCCLAstParser(tokens, options)
	return theParser.ParseAsAST()
}

func newCCLAstParser(
	tokens []*cclLexer.CCLToken,
	options *CCLParseOptions,
) *CCLAstParser {
	return &CCLAstParser{
		Options: options,
		tokens:  tokens,
	}
}
