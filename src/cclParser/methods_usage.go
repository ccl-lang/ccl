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
	last := len(allTokens) - 1

	isArray := allTokens[last].Type == cclLexer.TokenTypeRightBracket
	arrayLength := -1
	if isArray {
		// we have an array type
		// check if it's a fixed-length array
		if last-2 > 0 && allTokens[last-1].Type == cclLexer.TokenTypeIntLiteral &&
			allTokens[last-2].Type == cclLexer.TokenTypeLeftBracket {
			// fixed-length array
			arrayLength = allTokens[last-1].GetIntLiteral()
			// adjust the last index to point before the array tokens
			last -= 3
		} else if last-1 >= 0 && allTokens[last-1].Type == cclLexer.TokenTypeLeftBracket {
			// dynamic-length array
			// adjust the last index to point before the array tokens
			last -= 2
		} else {
			return nil, p.ErrInvalidSyntax("Invalid array type syntax")
		}
	}

	if last != 0 {
		// there are extra tokens we don't support for now
		return nil, p.ErrInvalidSyntax("Unsupported type syntax with extra tokens")
	}

	var baseTypeUsage *cclValues.CCLTypeUsage
	// I WAS HERE
	switch allTokens[0].Type {
	case cclLexer.TokenTypeDataType:
		baseTypeUsage = allTokens[0].GetBuiltInDataTypeUsage()
	case cclLexer.TokenTypeIdentifier:
		baseTypeUsage = allTokens[0].GetCustomTypeUsage(currentNamespace)
	default:
		return nil, p.ErrInvalidSyntax("Expected builtin data-type or an identifier as first token")
	}

	if isArray {
		return cclValues.NewArrayTypeUsage(baseTypeUsage, arrayLength), nil
	}
	return baseTypeUsage, nil
}
