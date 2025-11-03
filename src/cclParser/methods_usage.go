package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

// parseCurrentType tries to parse a ccl type-info from the current lexer tokens.
// No need to call advance after calling this method, as it will handle all the necessary
// advance calls by itself.
func (p *CCLParser) parseCurrentTypeUsage(currentNamespace string) (*cclValues.CCLTypeUsage, error) {
	allTokens := p.readUntilSemicolon()
	print(allTokens) //TODO: properly parse all the tokens here

	isDataType := p.isCurrentType(cclLexer.TokenTypeDataType)
	isIdentifier := p.isCurrentType(cclLexer.TokenTypeIdentifier)
	if !isDataType && !isIdentifier {
		return nil, p.ErrInvalidSyntax("Expected built-in data-type or an identifier as first token")
	}

	// I WAS HERE
	if isDataType {
		return p.current.GetBuiltInDataTypeUsage(), nil
	}

	return p.current.GetCustomTypeUsage(currentNamespace), nil
}
