package cclParser

import (
	"errors"
	"strings"

	"slices"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

//---------------------------------------------------------

func (p *CCLParser) ParseAsCCL() (*cclValues.SourceCodeDefinition, error) {
	p.codeDefinition = &cclValues.SourceCodeDefinition{}

	if err := p.initializeParsing(); err != nil {
		return nil, err
	}

	for !p.IsAtEnd() {
		if p.current.Type == cclLexer.TokenTypeHash {
			// in the future, the '#' token can be used for other things as well
			attribute, err := p.ParseGlobalAttribute()
			if err != nil {
				return nil, err
			}

			p.codeDefinition.GlobalAttributes = append(
				p.codeDefinition.GlobalAttributes,
				attribute,
			)
			continue
		}

		if p.current.Type == cclLexer.TokenTypeKeywordModel {
			model, err := p.ParseModel()
			if err != nil {
				return nil, err
			}

			p.codeDefinition.Models = append(p.codeDefinition.Models, model)
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
				p.codeDefinition.Models = append(p.codeDefinition.Models, model)
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

	return p.codeDefinition, nil
}

func (p *CCLParser) peekHasAssignment() bool {
	// keep doing a loop until we hit an assignment.
	// if we hit something other than assignment, we should return false
	// if we hit comment, we should skip it
	// if we hit EOF, we should return false
	currentPos := p.pos + 1
	for currentPos < len(p.tokens) {

		if p.tokens[currentPos].Type == cclLexer.TokenTypeComment {
			continue
		}

		if p.tokens[currentPos].Type == cclLexer.TokenTypeAssignment {
			return true
		}
	}

	return false
}

func (p *CCLParser) initializeParsing() error {
	if len(p.tokens) == 0 {
		return errors.New("cclParser: no tokens to parse")
	}

	if p.current != nil {
		return errors.New("cclParser: parser already initialized")
	}

	p.pos = 0
	p.current = p.tokens[0]

	return nil
}

// advance advances the parser to the next token.
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

// GetCurrent returns the current token being parsed.
// Please note that this method is exported mostly for tests.
func (p *CCLParser) GetCurrent() *cclLexer.CCLToken {
	return p.current
}

// IsAtEnd checks if the parser has reached the end of the input.
// Please note that this method is exported mostly for tests.
func (p *CCLParser) IsAtEnd() bool {
	return p.current == nil || p.current.Type == cclLexer.TokenTypeEOF
}

// isCurrentType checks if the current token is of the specified type.
// You can return multiple token types to check for, and if the current token
// is any of them, it will return true.
func (p *CCLParser) isCurrentType(tokenTypes ...cclLexer.CCLTokenType) bool {
	return slices.Contains(tokenTypes, p.current.Type)
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

// getCurrentSourceLine returns the source line for the given line number.
// Please note that calling this method frequently is expensive, since
// each time it will split the source code into lines.
// The reason we are not caching the lines is that we would only need
// to call this method in case of an error, and we don't want to keep
// the lines in memory for a long time.
func (p *CCLParser) getCurrentSourceLine(lineNum int) string {
	lines := strings.Split(p.Options.SourceContent, "\n")
	if lineNum > 0 && lineNum <= len(lines) {
		return lines[lineNum-1]
	}
	return ""
}

// IsCurrentValue checks if the current token is a value token.
func (p *CCLParser) IsCurrentValue() bool {
	return p.current.IsTokenValue()
}

// IsCurrentValue checks if the current token is a value token or an identifier.
func (p *CCLParser) IsCurrentValueOrIdentifier() bool {
	return p.current.IsTokenValue() ||
		p.current.Type == cclLexer.TokenTypeIdentifier
}

//---------------------------------------------------------
