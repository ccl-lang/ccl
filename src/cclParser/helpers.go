package cclParser

import (
	"os"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

// ParseCCLSourceFile reads a CCL source file and parses it into a
// SourceCodeDefinition value. It uses the CCL lexer to tokenize the source content and then
// parses the tokens using the CCLParser. The function returns a pointer to a SourceCodeDefinition
// value and an error if any occurred during the parsing process.
func ParseCCLSourceFile(options *CCLParseOptions) (*cclValues.SourceCodeDefinition, error) {
	content, err := os.ReadFile(options.SourceFilePath)
	if err != nil {
		return nil, err
	}

	options.SourceContent = string(content)
	return ParseCCLSourceContent(options)
}

// ParseCCLSourceContent takes a CCLParseOptions struct and parses the
// SourceContent field into a SourceCodeDefinition value. It uses the CCL lexer to tokenize
// the source content and then parses the tokens using the CCLParser. The function returns a
// pointer to a SourceCodeDefinition value and an error if any occurred during the parsing process.
func ParseCCLSourceContent(options *CCLParseOptions) (*cclValues.SourceCodeDefinition, error) {
	allTokens, err := cclLexer.Lex(options.SourceContent)
	if err != nil {
		return nil, err
	}

	if options.CodeContext == nil {
		options.CodeContext = cclValues.NewCCLCodeContext()
	}
	return ParseCCL(options.CodeContext, allTokens, options)
}

func ParseCCL(
	ctx *cclValues.CCLCodeContext,
	tokens []*cclLexer.CCLToken,
	options *CCLParseOptions,
) (*cclValues.SourceCodeDefinition, error) {
	theParser := &CCLParser{
		Options: options,
		tokens:  tokens,
		ctx:     ctx,
	}
	return theParser.ParseAsCCL()
}
