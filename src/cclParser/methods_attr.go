package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/cclSanitizer"
	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

//---------------------------------------------------------

func (p *CCLParser) ParseGlobalAttribute() (*cclValues.AttributeUsageInfo, error) {
	if err := p.consume(cclLexer.TokenTypeHash); err != nil {
		return nil, err
	}

	attrNode, err := p.parseSingleAttributeNode()
	if err != nil {
		return nil, err
	}

	globalNode := &cclAst.GlobalAttributeNode{
		Name:           attrNode.Name,
		Params:         attrNode.Params,
		SourcePosition: attrNode.SourcePosition,
	}

	return cclSanitizer.ResolveAttributeUsage(p.ctx, globalNode)
}

func (p *CCLParser) parseGlobalAttributeNode() (*cclAst.GlobalAttributeNode, error) {
	if err := p.consume(cclLexer.TokenTypeHash); err != nil {
		return nil, err
	}

	attrNode, err := p.parseSingleAttributeNode()
	if err != nil {
		return nil, err
	}

	return &cclAst.GlobalAttributeNode{
		Name:           attrNode.Name,
		Params:         attrNode.Params,
		SourcePosition: attrNode.SourcePosition,
	}, nil
}

// ParseAttributes Keeps parsing all of the available attributes in the current position
// until it hits something other than attribute.
func (p *CCLParser) ParseAttributes() ([]*cclValues.AttributeUsageInfo, error) {
	// in here, we will keep scanning until we hit something which is not an attribute
	allAttributes := []*cclValues.AttributeUsageInfo{}
	for !p.IsAtEnd() && p.isCurrentAttribute() {
		if p.isCurrentType(cclLexer.TokenTypeComment) {
			p.advance()
			continue
		}

		if p.isCurrentType(cclLexer.TokenTypeLeftBracket) {
			// we have an attribute here
			attributeNode, err := p.parseSingleAttributeNode()
			if err != nil {
				return nil, err
			}

			attribute, err := cclSanitizer.ResolveAttributeUsage(p.ctx, attributeNode)
			if err != nil {
				return nil, err
			}

			if attribute != nil {
				allAttributes = append(allAttributes, attribute)
			}
			continue
		}
	}

	return allAttributes, nil
}

func (p *CCLParser) parseAttributeNodes() ([]*cclAst.AttributeNode, error) {
	allAttributes := []*cclAst.AttributeNode{}
	for !p.IsAtEnd() && p.isCurrentAttribute() {
		if p.isCurrentType(cclLexer.TokenTypeComment) {
			p.advance()
			continue
		}

		if p.isCurrentType(cclLexer.TokenTypeLeftBracket) {
			attributeNode, err := p.parseSingleAttributeNode()
			if err != nil {
				return nil, err
			}

			if attributeNode != nil {
				allAttributes = append(allAttributes, attributeNode)
			}
			continue
		}
	}

	return allAttributes, nil
}

func (p *CCLParser) parseSingleAttributeNode() (*cclAst.AttributeNode, error) {
	// cache the starting token, since we are going to need some info from it
	startingToken := p.current
	if err := p.consume(cclLexer.TokenTypeLeftBracket); err != nil {
		return nil, err
	}

	name := p.current.GetIdentifier()
	if err := p.consume(cclLexer.TokenTypeIdentifier); err != nil {
		return nil, err
	}

	attrParams := []*cclAst.AttributeParamNode{}

	if p.current.Type == cclLexer.TokenTypeLeftParenthesis {
		p.advance()
		var currentParam *cclAst.AttributeParamNode
		for !p.isCurrentType(cclLexer.TokenTypeRightParenthesis) && !p.IsAtEnd() {
			// if we are having a ',' token, we should have a parameter before it
			if p.current.Type == cclLexer.TokenTypeComma {
				if currentParam == nil {
					return nil, &UnexpectedTokenError{
						Expected:       cclLexer.TokenTypeIdentifier,
						Actual:         cclLexer.TokenTypeComma,
						SourcePosition: p.getSourcePosition(),
					}
				}

				attrParams = append(attrParams, currentParam)
				currentParam = nil
				p.advance()
				continue
			}

			// This can actually bring us two possibilities:
			// 1. ParameterName = Value
			// 2. VariableName (which is defined somewhere else)
			if p.current.IsIdentifier() {
				if currentParam != nil {
					return nil, &UnexpectedTokenAfterParameterError{
						ParamName:      currentParam.Name,
						TokenValue:     p.current.GetIdentifier(),
						SourcePosition: p.getSourcePosition(),
					}
				}

				if !p.peekHasAssignment() {
					valueExpr := p.parseAttributeValueExpression()
					currentParam = &cclAst.AttributeParamNode{
						Value:          valueExpr,
						SourcePosition: valueExpr.GetSourcePosition(),
					}
					continue
				}

				gotAssignment := false
				paramName := p.current.GetIdentifier()
				currentParam = &cclAst.AttributeParamNode{
					Name:           paramName,
					SourcePosition: p.getSourcePosition(),
				}
				p.advance()

				// this loop will only break in these cases:
				// 1. we reach end of current param by ',' or ')'
				// 2. we reach EOF
				for {
					// early bound-check
					if p.IsAtEnd() {
						return nil, &UnexpectedEOFError{
							SourcePosition: p.getSourcePosition(),
						}
					}

					if p.isCurrentType(cclLexer.TokenTypeComment) {
						// just skip comments
						p.advance()
						continue
					}

					if p.isCurrentType(
						cclLexer.TokenTypeRightParenthesis,
						cclLexer.TokenTypeComma,
					) {
						if currentParam.Value == nil {
							// we got an identifier, but no value type
							// we were expecting some kind of value after this
							return nil, &ExpectedValueError{
								ParamName:      currentParam.Name,
								SourcePosition: p.getSourcePosition(),
							}
						}

						attrParams = append(attrParams, currentParam)
						currentParam = nil
						p.advance()
						break
					}

					if p.isCurrentType(cclLexer.TokenTypeAssignment) {
						if gotAssignment {
							return nil, &UnexpectedTokenAfterAssignmentError{
								ParamName:      currentParam.Name,
								TokenValue:     p.current.GetIdentifier(),
								SourcePosition: p.getSourcePosition(),
							}
						}
						gotAssignment = true
						p.advance()
						continue
					}

					if p.IsCurrentValueOrIdentifier() {
						if currentParam.Value != nil {
							return nil, &UnexpectedTokenAfterParameterError{
								ParamName:      currentParam.Name,
								TokenValue:     p.current.GetIdentifier(),
								SourcePosition: p.getSourcePosition(),
							}
						}

						valueExpr := p.parseAttributeValueExpression()
						currentParam.Value = valueExpr
						continue
					}

					// we got invalid token here
					return nil, &UnexpectedTokenAfterAssignmentError{
						ParamName:      currentParam.Name,
						TokenValue:     p.current.GetIdentifier(),
						SourcePosition: p.getSourcePosition(),
					}
				}

				continue
			}

			if p.IsCurrentLiteralValue() || p.IsCurrentReservedLiteral() {
				if currentParam != nil {
					return nil, &UnexpectedTokenAfterParameterError{
						ParamName:      currentParam.Name,
						TokenValue:     p.current.FormatValueAsString(),
						SourcePosition: p.getSourcePosition(),
					}
				}

				// if we are here, then we have a value without a parameter name
				valueExpr := p.parseAttributeValueExpression()
				currentParam = &cclAst.AttributeParamNode{
					Value:          valueExpr,
					SourcePosition: valueExpr.GetSourcePosition(),
				}
				continue
			}

			// if we reach here, then we have an unexpected token
			return nil, &UnexpectedTokenError{
				Expected:       cclLexer.TokenTypeIdentifier,
				Actual:         p.current.Type,
				SourcePosition: p.getSourcePosition(),
			}
		}

		if currentParam != nil {
			// we have unused current param here... we have to append it
			attrParams = append(attrParams, currentParam)
			// we set this to nil to to avoid misusing it in future
			currentParam = nil
		}

		if err := p.consume(cclLexer.TokenTypeRightParenthesis); err != nil {
			return nil, err
		}
	}

	if err := p.consume(cclLexer.TokenTypeRightBracket); err != nil {
		return nil, err
	}

	return &cclAst.AttributeNode{
		Name:           name,
		Params:         attrParams,
		SourcePosition: p.getSourcePositionForToken(startingToken),
	}, nil
}

func (p *CCLParser) parseAttributeValueExpression() cclAst.AttributeValueExpression {
	token := p.current
	sourcePos := p.getSourcePositionForToken(token)

	if token.IsIdentifier() {
		expr := &cclAst.IdentifierValueExpression{
			Name:           token.GetIdentifier(),
			SourcePosition: sourcePos,
		}
		p.advance()
		return expr
	}

	literalKind := cclAst.AttributeLiteralKindUnknown
	switch token.Type {
	case cclLexer.TokenTypeStringLiteral:
		literalKind = cclAst.AttributeLiteralKindString
	case cclLexer.TokenTypeIntLiteral:
		literalKind = cclAst.AttributeLiteralKindInt
	case cclLexer.TokenTypeFloatLiteral:
		literalKind = cclAst.AttributeLiteralKindFloat
	default:
		if token.IsReservedLiteral() {
			literalKind = cclAst.AttributeLiteralKindReserved
		}
	}

	literalExpr := &cclAst.LiteralValueExpression{
		LiteralKind:     literalKind,
		Value:           token.GetLiteralValue(),
		ReservedLiteral: token.GetReservedLiteral(),
		SourcePosition:  sourcePos,
	}
	p.advance()
	return literalExpr
}

func (p *CCLParser) isCurrentAttribute() bool {
	return p.isAttributeAt(p.pos)
}

func (p *CCLParser) isAttributeAt(targetPos int) bool {
	// if we are currently parsing an attribute, we should be having these:
	// 1. a left bracket
	// 2. an identifier
	// 3. a parenthesis

	// first, length safety check
	if targetPos+2 >= len(p.tokens) {
		return false
	}

	return p.tokens[targetPos].Type == cclLexer.TokenTypeLeftBracket &&
		p.tokens[targetPos+1].Type == cclLexer.TokenTypeIdentifier
}

func (p *CCLParser) peekAfterAttribute() cclLexer.CCLTokenType {
	tempPos := p.pos

	// Skip all attributes
	for p.isAttributeAt(tempPos) {
		// Find the matching right bracket
		bracketCount := 1
		tempPos++
		for tempPos < len(p.tokens) && bracketCount > 0 {
			switch p.tokens[tempPos].Type {
			case cclLexer.TokenTypeLeftBracket:
				bracketCount++
			case cclLexer.TokenTypeRightBracket:
				bracketCount--
			}
			tempPos++
		}

		if tempPos >= len(p.tokens) {
			return cclLexer.TokenTypeEOF
		}
	}

	if tempPos < len(p.tokens) {
		return p.tokens[tempPos].Type
	}
	return cclLexer.TokenTypeEOF
}

//---------------------------------------------------------
