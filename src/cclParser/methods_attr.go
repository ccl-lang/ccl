package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	"github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

func (p *CCLParser) ParseGlobalAttribute() (*cclValues.AttributeUsageInfo, error) {
	if err := p.consume(cclLexer.TokenTypeHash); err != nil {
		return nil, err
	}

	return p.parseSingleAttribute()
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
			attribute, err := p.parseSingleAttribute()
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

func (p *CCLParser) parseSingleAttribute() (*cclValues.AttributeUsageInfo, error) {
	// cache the starting token, since we are going to need some
	startingToken := p.current
	if err := p.consume(cclLexer.TokenTypeLeftBracket); err != nil {
		return nil, err
	}

	name := p.current.GetIdentifier()
	if err := p.consume(cclLexer.TokenTypeIdentifier); err != nil {
		return nil, err
	}

	attrParams := []*cclValues.ParameterInstance{}

	if p.current.Type == cclLexer.TokenTypeLeftParenthesis {
		p.advance()
		var currentParam *cclValues.ParameterInstance
		for !p.isCurrentType(cclLexer.TokenTypeRightParenthesis) && !p.IsAtEnd() {
			// if we are having a ',' token, we should have a parameter before it
			if p.current.Type == cclLexer.TokenTypeComma {
				if currentParam == nil {
					return nil, &UnexpectedTokenError{
						Expected: cclLexer.TokenTypeIdentifier,
						Actual:   cclLexer.TokenTypeComma,
						Line:     p.current.Line,
						Column:   p.current.Column,
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
			if p.current.Type == cclLexer.TokenTypeIdentifier {
				if currentParam != nil {
					return nil, &UnexpectedTokenAfterParameterError{
						Line:       p.current.Line,
						Column:     p.current.Column,
						SourceLine: p.getCurrentSourceLine(p.current.Line),
						ParamName:  currentParam.Name,
						TokenValue: p.current.GetIdentifier(),
					}
				}

				if !p.peekHasAssignment() {
					// We have a variable usage here
					// without assigning any param name
					targetIdentifier := p.current.GetIdentifier()
					targetVariable := cclValues.GetGlobalVariable(targetIdentifier)
					if targetVariable == nil {
						return nil, &UndefinedIdentifierError{
							Line:             p.current.Line,
							Column:           p.current.Column,
							SourceLine:       p.getCurrentSourceLine(p.current.Line),
							TargetIdentifier: targetIdentifier,
							Language:         globalValues.LanguageCCL,
						}
					}
					currentParam = &cclValues.ParameterInstance{}
					if targetVariable.IsAutomatic() {
						currentParam.ChangeValueType(cclValues.NewPointerTypeInfo(targetVariable.Type))
						currentParam.ChangeValue(&cclValues.VariableUsageInstance{
							Name:       p.current.GetIdentifier(),
							Definition: targetVariable,
						})
					} else {
						// since the variable is not an automatic variable, we don't
						// have to *point* to it.
						currentParam.ChangeValueType(targetVariable.Type)
						currentParam.ChangeValue(targetVariable.GetValue())
					}
					p.advance()
					continue
				}

				gotAssignment := false
				currentParam = &cclValues.ParameterInstance{
					Name: p.current.GetIdentifier(),
				}
				p.advance()

				// this loop will only break in these cases:
				// 1. we reach end of current param by ',' or ')'
				// 2. we reach EOF
				for {
					// early bound-check
					if p.IsAtEnd() {
						return nil, &UnexpectedEOFError{
							Line:       p.current.Line,
							Column:     p.current.Column,
							SourceLine: p.getCurrentSourceLine(p.current.Line),
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
						if currentParam.ValueType == nil {
							// we got an identifier, but no value type
							// we were expecting some kind of value after this
							return nil, &ExpectedValueError{
								SourceLine: p.getCurrentSourceLine(p.current.Line),
								ParamName:  currentParam.Name,
								Line:       p.current.Line,
								Column:     p.current.Column,
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
								Line:       p.current.Line,
								Column:     p.current.Column,
								SourceLine: p.getCurrentSourceLine(p.current.Line),
								ParamName:  currentParam.Name,
								TokenValue: p.current.GetIdentifier(),
							}
						}
						gotAssignment = true
						p.advance()
						continue
					}

					if p.IsCurrentValueOrIdentifier() {
						if currentParam.ValueType != nil {
							return nil, &UnexpectedTokenAfterParameterError{
								Line:       p.current.Line,
								Column:     p.current.Column,
								SourceLine: p.getCurrentSourceLine(p.current.Line),
								ParamName:  currentParam.Name,
								TokenValue: p.current.GetIdentifier(),
							}
						}

						// is it a literal value?
						if p.IsCurrentValue() {
							currentParam.ChangeValueType(p.current.GetLiteralTypeInfo())
							currentParam.ChangeValue(p.current.GetLiteralValue())
							p.advance()
							continue
						}

						// We have a variable usage here
						// and the parameter name has previously been assigned
						targetIdentifier := p.current.GetIdentifier()
						targetVariable := cclValues.GetGlobalVariable(targetIdentifier)
						if targetVariable == nil {
							return nil, &UndefinedIdentifierError{
								Line:             p.current.Line,
								Column:           p.current.Column,
								SourceLine:       p.getCurrentSourceLine(p.current.Line),
								TargetIdentifier: targetIdentifier,
								Language:         globalValues.LanguageCCL,
							}
						}
						if targetVariable.IsAutomatic() {
							currentParam.ChangeValueType(cclValues.NewPointerTypeInfo(targetVariable.Type))
							currentParam.ChangeValue(&cclValues.VariableUsageInstance{
								Name:       p.current.GetIdentifier(),
								Definition: targetVariable,
							})
						} else {
							// since the variable is not an automatic variable, we don't
							// have to *point* to it.
							currentParam.ChangeValueType(targetVariable.Type)
							currentParam.ChangeValue(targetVariable.GetValue()) // copy the value
						}
						p.advance()
						continue
					}
				}

				continue
			}

			if p.IsCurrentValue() {
				if currentParam != nil {
					return nil, &UnexpectedTokenAfterParameterError{
						Line:       p.current.Line,
						Column:     p.current.Column,
						SourceLine: p.getCurrentSourceLine(p.current.Line),
						ParamName:  currentParam.Name,
						TokenValue: p.current.GetIdentifier(),
					}
				}

				// if we are here, then we have a value without a parameter name
				currentParam = &cclValues.ParameterInstance{
					ValueType: p.current.GetLiteralTypeInfo(),
				}
				currentParam.ChangeValue(p.current.GetLiteralValue())
				p.advance()
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

	return &cclValues.AttributeUsageInfo{
		Name:       name,
		Parameters: attrParams,
		Line:       startingToken.Line,
		Column:     startingToken.Column,
	}, nil
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
			if p.tokens[tempPos].Type == cclLexer.TokenTypeLeftBracket {
				bracketCount++
			} else if p.tokens[tempPos].Type == cclLexer.TokenTypeRightBracket {
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
