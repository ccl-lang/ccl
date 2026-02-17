package cclParser

import (
	"os"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

// ParseCCLSourceFileAsAST reads a CCL source file and parses it into a CCLFileAST.
func ParseCCLSourceFileAsAST(options *CCLParseOptions) (*cclAst.CCLFileAST, error) {
	content, err := os.ReadFile(options.SourceFilePath)
	if err != nil {
		return nil, err
	}

	options.SourceContent = string(content)
	return ParseCCLSourceContentAsAST(options)
}

// ParseCCLSourceContentAsAST parses the source content into a CCLFileAST (syntax only).
func ParseCCLSourceContentAsAST(options *CCLParseOptions) (*cclAst.CCLFileAST, error) {
	allTokens, err := cclLexer.Lex(options.SourceContent)
	if err != nil {
		return nil, err
	}

	if options.CodeContext == nil {
		options.CodeContext = cclValues.NewCCLCodeContext()
	}

	return ParseCCLAst(options.CodeContext, allTokens, options)
}

func ParseCCLAst(
	ctx *cclValues.CCLCodeContext,
	tokens []*cclLexer.CCLToken,
	options *CCLParseOptions,
) (*cclAst.CCLFileAST, error) {
	theParser := &CCLParser{
		Options: options,
		tokens:  tokens,
		ctx:     ctx,
	}
	return theParser.ParseAsAST()
}
