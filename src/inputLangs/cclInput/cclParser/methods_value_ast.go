package cclParser

import (
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser/cclLexer"
)

func (p *CCLAstParser) parseValueExpressionUntil(
	stopTokens ...cclLexer.CCLTokenType,
) (cclAst.ValueExpression, error) {
	tokens := []*cclLexer.CCLToken{}
	for !p.IsAtEnd() && !p.isCurrentType(stopTokens...) {
		if p.isCurrentType(cclLexer.TokenTypeComment) {
			p.advance()
			continue
		}

		tokens = append(tokens, p.current)
		p.advance()
	}

	return p.parseValueExpressionFromTokens(tokens)
}

func (p *CCLAstParser) parseValueExpressionFromTokens(
	tokens []*cclLexer.CCLToken,
) (cclAst.ValueExpression, error) {
	if len(tokens) == 0 {
		return nil, p.ErrInvalidSyntax("Missing value expression")
	}

	if len(tokens) == 1 {
		token := tokens[0]
		if token.IsTokenLiteralValue() || token.IsReservedLiteral() {
			return p.literalExpressionFromToken(token), nil
		}

		if token.Type == cclLexer.TokenTypeIdentifier {
			return &cclAst.IdentifierValueExpression{
				Name:           token.GetIdentifier(),
				SourcePosition: p.getSourcePositionForToken(token),
			}, nil
		}
	}

	if len(tokens) == 2 && tokens[0].Type == cclLexer.TokenTypeMinus &&
		(tokens[1].Type == cclLexer.TokenTypeIntLiteral ||
			tokens[1].Type == cclLexer.TokenTypeFloatLiteral) {
		return p.negativeLiteralExpressionFromToken(tokens[1])
	}

	return p.qualifiedIdentifierExpressionFromTokens(tokens)
}

func (p *CCLAstParser) literalExpressionFromToken(
	token *cclLexer.CCLToken,
) *cclAst.LiteralValueExpression {
	sourcePos := p.getSourcePositionForToken(token)
	literalExpr := &cclAst.LiteralValueExpression{
		SourcePosition: sourcePos,
	}

	switch token.Type {
	case cclLexer.TokenTypeStringLiteral:
		literalExpr.LiteralKind = cclAst.AttributeLiteralKindString
		literalExpr.Value = token.GetStringLiteral()
	case cclLexer.TokenTypeIntLiteral:
		literalExpr.LiteralKind = cclAst.AttributeLiteralKindInt
		literalExpr.Value = token.GetInt64Literal()
	case cclLexer.TokenTypeFloatLiteral:
		literalExpr.LiteralKind = cclAst.AttributeLiteralKindFloat
		literalExpr.Value = token.GetFloatLiteral()
	case cclLexer.TokenTypeReservedLiteral:
		literalExpr.LiteralKind = cclAst.AttributeLiteralKindReserved
		literalExpr.ReservedLiteral = token.GetReservedLiteral()
		literalExpr.Value = token.GetLiteralValue()
	}

	return literalExpr
}

func (p *CCLAstParser) negativeLiteralExpressionFromToken(
	token *cclLexer.CCLToken,
) (*cclAst.LiteralValueExpression, error) {
	sourcePos := p.getSourcePositionForToken(token)
	switch token.Type {
	case cclLexer.TokenTypeIntLiteral:
		return &cclAst.LiteralValueExpression{
			LiteralKind:    cclAst.AttributeLiteralKindInt,
			Value:          -token.GetInt64Literal(),
			SourcePosition: sourcePos,
		}, nil
	case cclLexer.TokenTypeFloatLiteral:
		return &cclAst.LiteralValueExpression{
			LiteralKind:    cclAst.AttributeLiteralKindFloat,
			Value:          -token.GetFloatLiteral(),
			SourcePosition: sourcePos,
		}, nil
	default:
		return nil, p.ErrInvalidSyntax("Expected numeric literal after '-'")
	}
}

func (p *CCLAstParser) qualifiedIdentifierExpressionFromTokens(
	tokens []*cclLexer.CCLToken,
) (cclAst.ValueExpression, error) {
	parts := []string{}
	expectIdentifier := true
	for _, token := range tokens {
		if expectIdentifier {
			if token.Type != cclLexer.TokenTypeIdentifier {
				return nil, p.ErrInvalidSyntax("Expected identifier in value expression")
			}
			parts = append(parts, token.GetIdentifier())
			expectIdentifier = false
			continue
		}

		if token.Type != cclLexer.TokenTypeDot {
			return nil, p.ErrInvalidSyntax("Expected dot in qualified value expression")
		}
		expectIdentifier = true
	}

	if expectIdentifier || len(parts) < 2 {
		return nil, p.ErrInvalidSyntax("Expected qualified identifier in value expression")
	}

	return &cclAst.QualifiedIdentifierValueExpression{
		Parts:          parts,
		SourcePosition: p.getSourcePositionForToken(tokens[0]),
	}, nil
}
