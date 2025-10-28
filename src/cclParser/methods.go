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

	var currentPendingAttributes []*cclValues.AttributeUsageInfo

	for !p.IsAtEnd() {
		if p.current.Type == cclLexer.TokenTypeHash {
			// in the future, the '#' token can be used for other things as well
			if p.isAttributeAt(p.pos + 1) {
				attribute, err := p.ParseGlobalAttribute()
				if err != nil {
					return nil, err
				}

				p.codeDefinition.GlobalAttributes = append(
					p.codeDefinition.GlobalAttributes,
					attribute,
				)
			}
			continue
		}

		if p.isCurrentAttribute() {
			// an attribute should come before a model (or some other keyword)
			afterAttribute := p.peekAfterAttribute()
			if afterAttribute == cclLexer.TokenTypeEOF {
				return nil, &UnexpectedEndOfAttributeError{
					Line:       p.current.Line,
					Column:     p.current.Column,
					SourceLine: p.getCurrentSourceLine(p.current.Line),
				}
			}

			allAttributes, err := p.ParseAttributes()
			if err != nil {
				return nil, err
			}
			// since we don't want to make our parser too complex, we will just set
			// pending attributes here and let other parts of parser handle this.
			// E.g. if after this, we get a model, we will set the attributes to the model;
			// or if we get a field, we will set the attributes to the field, etc...
			currentPendingAttributes = append(currentPendingAttributes, allAttributes...)
			continue
		}

		if p.current.Type == cclLexer.TokenTypeKeywordModel {
			model, err := p.ParseModelDefinition()
			if err != nil {
				return nil, err
			}

			if len(currentPendingAttributes) > 0 {
				// we have some pending attributes, we should set them to the model
				model.Attributes = append(
					model.Attributes,
					currentPendingAttributes...,
				)
				currentPendingAttributes = nil
			}
			p.codeDefinition.Models = append(p.codeDefinition.Models, model)
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

	if len(currentPendingAttributes) > 0 {
		// we have unused *normal* attributes...which is a compiler error
		lastToken := currentPendingAttributes[len(currentPendingAttributes)-1]
		return nil, &InvalidAttributeUsageError{
			Line:       lastToken.Line,
			Column:     lastToken.Column,
			SourceLine: p.getCurrentSourceLine(lastToken.Line),
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

// isCurrentComment checks if the current token is a comment.
func (p *CCLParser) isCurrentComment() bool {
	return p.current.Type == cclLexer.TokenTypeComment
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
	result := ""
	lines := strings.Split(p.Options.SourceContent, "\n")
	if lineNum > 0 && lineNum <= len(lines) {
		result = lines[lineNum-1]
	}

	if len(result) > MaxShownSourceLineLen {
		result = result[:MaxShownSourceLineLen]
	}

	return result
}

// IsCurrentValue checks if the current token is a literal value token.
func (p *CCLParser) IsCurrentValue() bool {
	return p.current.IsTokenValue()
}

// FindTokenPattern peeks in front of the current token to see if the provided
// patterns match the tokens from now on.
// Comments are ignored. Tokens which are -1 in the specified arguments, will
// make this function to accept any token.
func (p *CCLParser) FindTokenPattern(tokens []cclLexer.CCLTokenType) bool {
	currentTargetIndex := 0
	currentPos := p.pos - 1
	if len(p.tokens)-p.pos < len(tokens) {
		// general bound-checking before entering the loop
		return false
	}
	for {
		currentPos++
		if currentPos >= len(p.tokens) {
			return false
		}

		if p.tokens[currentPos].Type == cclLexer.TokenTypeComment {
			// can be safely ignored
			continue
		}

		if tokens[currentTargetIndex] != cclLexer.TokenTypeReservedForFuture &&
			p.tokens[currentPos].Type != tokens[currentTargetIndex] {
			return false
		}

		currentTargetIndex++
		if currentTargetIndex >= len(tokens) {
			// everything is matched
			return true
		}
	}
}

// IsCurrentValue checks if the current token is a value token or an identifier.
func (p *CCLParser) IsCurrentValueOrIdentifier() bool {
	return p.current.IsTokenValue() ||
		p.current.Type == cclLexer.TokenTypeIdentifier
}

// parseCurrentType tries to parse a ccl type-info from the current lexer tokens.
// No need to call advance after calling this method, as it will handle all the necessary
// advance calls by itself.
func (p *CCLParser) parseCurrentType() (*cclValues.CCLTypeInfo, error) {
	isDataType := p.isCurrentType(cclLexer.TokenTypeDataType)
	isIdentifier := p.isCurrentType(cclLexer.TokenTypeIdentifier)
	if !isDataType && !isIdentifier {
		return nil, p.ErrInvalidSyntax("Expected built-in data-type or an identifier as first token")
	}

	typeName := ""
	if isDataType {
		typeName = p.current.GetLiteralTypeInfo().GetName()
	} else {

	}
	if typeName == "" {
		return nil, p.ErrInvalidSyntax("Unexpected empty type-name")
	}

	return nil, nil
}

//---------------------------------------------------------
