package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/cclSanitizer"
	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

// parseCurrentType tries to parse a ccl type-info from the current lexer tokens.
// No need to call advance after calling this method, as it will handle all the necessary
// advance calls by itself.
func (p *CCLParser) parseCurrentTypeUsage(currentNamespace string) (*cclValues.CCLTypeUsage, error) {
	typeExpr, err := p.parseCurrentTypeExpression()
	if err != nil {
		return nil, err
	}

	return cclSanitizer.ResolveTypeUsage(p.ctx, currentNamespace, typeExpr)
}

func (p *CCLParser) parseCurrentTypeExpression() (cclAst.TypeExpression, error) {
	allTokens := p.readUntilSemicolon()
	if len(allTokens) == 0 {
		return nil, p.ErrInvalidSyntax("Missing type expression")
	}

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

	baseToken := allTokens[0]
	basePosition := p.getSourcePositionForToken(baseToken)

	var baseTypeExpr *cclAst.SimpleTypeExpression
	switch baseToken.Type {
	case cclLexer.TokenTypeDataType:
		baseTypeExpr = &cclAst.SimpleTypeExpression{
			TypeName: cclAst.SimpleTypeName{
				Name: baseToken.GetDataTypeName(),
			},
			IsBuiltinToken: true,
			SourcePosition: basePosition,
		}
	case cclLexer.TokenTypeIdentifier:
		baseTypeExpr = &cclAst.SimpleTypeExpression{
			TypeName: cclAst.SimpleTypeName{
				Name: baseToken.GetIdentifier(),
			},
			IsBuiltinToken: false,
			SourcePosition: basePosition,
		}
	default:
		return nil, p.ErrInvalidSyntax("Expected builtin data-type or an identifier as first token")
	}

	if isArray {
		return &cclAst.ArrayTypeExpression{
			ElementType:    baseTypeExpr,
			Length:         arrayLength,
			SourcePosition: basePosition,
		}, nil
	}

	return baseTypeExpr, nil
}
