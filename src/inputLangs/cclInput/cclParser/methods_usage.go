package cclParser

import (
	"strings"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser/cclLexer"
)

func (p *CCLAstParser) parseCurrentTypeExpression() (cclAst.TypeExpression, error) {
	allTokens := p.readUntilSemicolon()
	return p.parseTypeExpressionFromTokens(allTokens)
}

func (p *CCLAstParser) parseTypeExpressionUntil(
	stopTokens ...cclLexer.CCLTokenType,
) (cclAst.TypeExpression, error) {
	allTokens := []*cclLexer.CCLToken{}
	for !p.IsAtEnd() && !p.isCurrentType(stopTokens...) {
		if p.isCurrentType(cclLexer.TokenTypeComment) {
			p.advance()
			continue
		}

		allTokens = append(allTokens, p.current)
		p.advance()
	}

	return p.parseTypeExpressionFromTokens(allTokens)
}

func (p *CCLAstParser) parseTypeExpressionFromTokens(
	allTokens []*cclLexer.CCLToken,
) (cclAst.TypeExpression, error) {
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

	baseTypeExpr, err := p.parseSimpleTypeExpressionFromTokens(allTokens[:last+1])
	if err != nil {
		return nil, err
	}

	if isArray {
		return &cclAst.ArrayTypeExpression{
			ElementType:    baseTypeExpr,
			Length:         arrayLength,
			SourcePosition: baseTypeExpr.GetSourcePosition(),
		}, nil
	}

	return baseTypeExpr, nil
}

func (p *CCLAstParser) parseSimpleTypeExpressionFromTokens(
	tokens []*cclLexer.CCLToken,
) (*cclAst.SimpleTypeExpression, error) {
	if len(tokens) == 0 {
		return nil, p.ErrInvalidSyntax("Missing type name")
	}

	baseToken := tokens[0]
	basePosition := p.getSourcePositionForToken(baseToken)
	if baseToken.Type == cclLexer.TokenTypeDataType {
		if len(tokens) != 1 {
			return nil, p.ErrInvalidSyntax("Built-in data-types cannot be qualified")
		}

		return &cclAst.SimpleTypeExpression{
			TypeName: cclAst.SimpleTypeName{
				Name: baseToken.GetDataTypeName(),
			},
			IsBuiltinToken: true,
			SourcePosition: basePosition,
		}, nil
	}

	parts := []string{}
	expectIdentifier := true
	for _, token := range tokens {
		if expectIdentifier {
			if token.Type != cclLexer.TokenTypeIdentifier {
				return nil, p.ErrInvalidSyntax("Expected identifier in type name")
			}
			parts = append(parts, token.GetIdentifier())
			expectIdentifier = false
			continue
		}

		if token.Type != cclLexer.TokenTypeDot {
			return nil, p.ErrInvalidSyntax("Expected dot in qualified type name")
		}
		expectIdentifier = true
	}

	if expectIdentifier || len(parts) == 0 {
		return nil, p.ErrInvalidSyntax("Incomplete qualified type name")
	}

	return &cclAst.SimpleTypeExpression{
		TypeName: cclAst.SimpleTypeName{
			Name:      parts[len(parts)-1],
			Namespace: strings.Join(parts[:len(parts)-1], "."),
		},
		IsBuiltinToken: false,
		SourcePosition: basePosition,
	}, nil
}
