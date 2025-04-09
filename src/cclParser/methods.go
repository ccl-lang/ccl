package cclParser

import (
	"fmt"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

//---------------------------------------------------------

func (p *CCLParser) ParseAsCCL() (*cclValues.SourceCodeDefinition, error) {
	program := &cclValues.SourceCodeDefinition{}

	for p.current.Type != cclLexer.TokenTypeEOF {
		if p.current.Type == cclLexer.TokenTypeHash {
			// in the future, the '#' token can be used for other things as well
			attribute, err := p.ParseGlobalAttribute()
			if err != nil {
				return nil, err
			}

			program.GlobalAttributes = append(program.GlobalAttributes, attribute)
			continue
		}

		if p.current.Type == cclLexer.TokenTypeKeywordModel {
			model, err := p.ParseModel()
			if err != nil {
				return nil, err
			}

			program.Models = append(program.Models, model)
			continue
		}

		if p.isCurrentAttribute() {
			// an attribute should come before a model (or some other keyword)
			afterAttribute := p.peekAfterAttribute()
			if afterAttribute == cclLexer.TokenTypeEOF {
				return nil, &UnexpectedEndOfAttributeError{
					Line:   p.current.Line,
					Column: p.current.Column,
				}
			}

			if afterAttribute == cclLexer.TokenTypeKeywordModel {
				allAttributes, err := p.ParseAttributes()
				if err != nil {
					return nil, err
				}

				model, err := p.ParseModel()
				if err != nil {
					return nil, err
				}

				model.Attributes = allAttributes
				program.Models = append(program.Models, model)
			} else {
				// parse other definitions alongside of attributes
				p.advance()
			}

			continue
		}

		if p.current.Type == cclLexer.TokenTypeComment {
			// skip comments
			// TODO: in future, add the ability to have document comments, in a way
			// that they can be attached to the model or attribute
			p.advance()
			continue
		}

		// if we reach here, then we have an unexpected token
		return nil, &UnexpectedTokenError{
			Expected: cclLexer.TokenTypeKeywordModel,
			Actual:   p.current.Type,
			Line:     p.current.Line,
			Column:   p.current.Column,
		}
	}
	return nil, nil
}

func (p *CCLParser) ParseGlobalAttribute() (*cclValues.AttributeDefinition, error) {
	if err := p.consume(cclLexer.TokenTypeHash); err != nil {
		return nil, err
	}

	if err := p.consume(cclLexer.TokenTypeLeftBracket); err != nil {
		return nil, err
	}

	name := p.current.GetIdentifier()
	if err := p.consume(cclLexer.TokenTypeIdentifier); err != nil {
		return nil, err
	}

	attrParams := []*cclValues.ParameterDefinition{}

	if p.current.Type == cclLexer.TokenTypeLeftParenthesis {
		p.advance()
		var currentParam *cclValues.ParameterDefinition
		for p.current.Type != cclLexer.TokenTypeRightParenthesis {
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

			// I WAS HERE

		}
	}

	return &cclValues.AttributeDefinition{
		Name:       name,
		Parameters: attrParams,
	}, nil
}

// ParseAttributes Keeps parsing all of the available attributes in the current position
// until it hits something other than attribute.
func (p *CCLParser) ParseAttributes() ([]*cclValues.AttributeDefinition, error) {
	// TODO
	return nil, nil
}

func (p *CCLParser) ParseModel() (*cclValues.ModelDefinition, error) {
	// TODO
	return nil, nil
}

func (p *CCLParser) advance() {
	p.pos++
	if p.pos < len(p.tokens) {
		p.current = p.tokens[p.pos]
	} else {
		p.current = &cclLexer.CCLToken{
			Type: cclLexer.TokenTypeEOF,
		}
	}
}

func (p *CCLParser) consume(tokenType cclLexer.CCLTokenType) error {
	if p.current.Type == tokenType {
		p.advance()
		return nil
	}

	return &UnexpectedTokenError{
		Expected: tokenType,
		Actual:   p.current.Type,
		Line:     p.current.Line,
		Column:   p.current.Column,
	}
}

func (p *CCLParser) isCurrentAttribute() bool {
	// if we are currently parsing an attribute, we should be having these:
	// 1. a left bracket
	// 2. an identifier
	// 3. a parenthesis

	// first, length safety check
	if p.pos+2 >= len(p.tokens) {
		return false
	}

	return p.tokens[p.pos].Type == cclLexer.TokenTypeLeftBracket &&
		p.tokens[p.pos+1].Type == cclLexer.TokenTypeIdentifier &&
		p.tokens[p.pos+2].Type == cclLexer.TokenTypeLeftParenthesis
}

func (p *CCLParser) peekAfterAttribute() cclLexer.CCLTokenType {
	tempPos := p.pos

	// Skip all attributes
	for p.tokens[tempPos].Type == cclLexer.TokenTypeLeftBracket {
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

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf(
		"cclParser: expected token type %s, got %s at line %d, column %d",
		e.Expected,
		e.Actual,
		e.Line,
		e.Column,
	)
}

func (e *UnexpectedEndOfAttributeError) Error() string {
	return fmt.Sprintf(
		"cclParser: unexpected end of attribute at line %d, column %d",
		e.Line,
		e.Column,
	)
}

//---------------------------------------------------------
