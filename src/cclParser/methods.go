package cclParser

import (
	"errors"
	"strings"

	"slices"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclUtils"
)

//---------------------------------------------------------

// peekHasAssignment returns true if there is an assignment token in front of the current token,
// ignoring comments.
func (p *CCLAstParser) peekHasAssignment() bool {
	// keep doing a loop until we hit an assignment.
	// if we hit something other than assignment, we should return false
	// if we hit comment, we should skip it
	// if we hit EOF, we should return false
	currentPos := p.pos + 1
	for currentPos < len(p.tokens) {
		if p.tokens[currentPos].Type == cclLexer.TokenTypeComment {
			currentPos++
			continue
		}

		if p.tokens[currentPos].Type == cclLexer.TokenTypeAssignment {
			return true
		}

		return false
	}

	return false
}

func (p *CCLAstParser) initializeParsing() error {
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
func (p *CCLAstParser) advance() {
	p.pos++
	if p.pos < len(p.tokens) {
		p.current = p.tokens[p.pos]
	} else {
		p.current = &cclLexer.CCLToken{
			Type: cclLexer.TokenTypeEOF,
		}
	}
}

// readUntilSemicolon reads tokens until it hits a semicolon or the end of the input.
func (p *CCLAstParser) readUntilSemicolon() []*cclLexer.CCLToken {
	startPos := p.pos

	for !p.isCurrentType(cclLexer.TokenTypeSemicolon) && !p.IsAtEnd() {
		p.advance()
	}

	return p.tokens[startPos:p.pos]
}

// GetCurrent returns the current token being parsed.
// Please note that this method is exported mostly for tests.
func (p *CCLAstParser) GetCurrent() *cclLexer.CCLToken {
	return p.current
}

// IsAtEnd checks if the parser has reached the end of the input.
// Please note that this method is exported mostly for tests.
func (p *CCLAstParser) IsAtEnd() bool {
	return p.current == nil || p.current.Type == cclLexer.TokenTypeEOF
}

// isCurrentType checks if the current token is any of the specified type.
// You can pass multiple token types to check for, and if the current token
// is any of them, it will return true.
func (p *CCLAstParser) isCurrentType(tokenTypes ...cclLexer.CCLTokenType) bool {
	return slices.Contains(tokenTypes, p.current.Type)
}

// isNextType checks if the next token is any of the specified type.
// You can pass multiple token types to check for, and if the next token
// is any of them, it will return true.
func (p *CCLAstParser) isNextType(tokenTypes ...cclLexer.CCLTokenType) bool {
	if p.pos+1 >= len(p.tokens) {
		return false
	}
	return slices.Contains(tokenTypes, p.tokens[p.pos+1].Type)
}

// isCurrentComment checks if the current token is a comment.
func (p *CCLAstParser) isCurrentComment() bool {
	return p.current.Type == cclLexer.TokenTypeComment
}

func (p *CCLAstParser) consume(tokenType cclLexer.CCLTokenType) error {
	if p.current.Type == tokenType {
		p.advance()
		return nil
	}

	return &UnexpectedTokenError{
		Expected:       tokenType,
		Actual:         p.current.Type,
		SourcePosition: p.getSourcePosition(),
	}
}

// getCurrentSourceLine returns the source line for the given line number.
// Please note that calling this method frequently is expensive, since
// each time it will split the source code into lines.
// The reason we are not caching the lines is that we would only need
// to call this method in case of an error, and we don't want to keep
// the lines in memory for a long time.
func (p *CCLAstParser) getCurrentSourceLine(lineNum int) string {
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

// getSourcePosition returns the current source position.
func (p *CCLAstParser) getSourcePosition() *cclUtils.SourceCodePosition {
	return &cclUtils.SourceCodePosition{
		Line:       p.current.Line,
		Column:     p.current.Column,
		SourceLine: p.getCurrentSourceLine(p.current.Line),
	}
}

func (p *CCLAstParser) getSourcePositionForToken(token *cclLexer.CCLToken) *cclUtils.SourceCodePosition {
	if token == nil {
		return nil
	}

	return &cclUtils.SourceCodePosition{
		Line:       token.Line,
		Column:     token.Column,
		SourceLine: p.getCurrentSourceLine(token.Line),
	}
}

// IsCurrentValue checks if the current token is a literal value token.
func (p *CCLAstParser) IsCurrentLiteralValue() bool {
	return p.current.IsTokenLiteralValue()
}

// IsCurrentReservedLiteral checks if the current token is a reserved literal.
func (p *CCLAstParser) IsCurrentReservedLiteral() bool {
	return p.current.IsReservedLiteral()
}

// FindTokenPattern peeks in front of the current token to see if the provided
// patterns match the tokens from now on.
// Comments are ignored. Tokens which are -1 in the specified arguments, will
// make this function to accept any token.
func (p *CCLAstParser) FindTokenPattern(tokens []cclLexer.CCLTokenType) bool {
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

// IsCurrentValueOrIdentifier checks if the current token is a literal value
// token or an identifier.
func (p *CCLAstParser) IsCurrentValueOrIdentifier() bool {
	return p.current.IsTokenLiteralValue() || p.current.IsIdentifier()
}

//---------------------------------------------------------
