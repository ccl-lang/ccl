package cclLexer

import (
	"unicode"

	"github.com/ALiwoto/ccl/src/core/cclValues"
)

func Lex(input string) ([]*CCLToken, error) {
	var tokens []*CCLToken
	line := 1
	column := 1
	runes := []rune(input)
	pos := 0
	totalRunesLen := len(runes)

	for pos < totalRunesLen {
		r := runes[pos]

		if isWhitespace(r) {
			isNextLF := pos+1 < totalRunesLen && runes[pos+1] == '\n'
			if r == '\n' || r == '\r' {
				// New Line
				line++
				column = 1
				if isNextLF {
					pos++
				}
			} else {
				column++
			}
			pos++
			continue
		}

		if r == '/' && pos+1 < totalRunesLen && runes[pos+1] == '/' {
			// Comment
			comment := "//"
			pos += 2
			column += 2

			for pos < totalRunesLen && runes[pos] != '\n' {
				comment += string(runes[pos])
				pos++
				column++
			}

			tokens = append(tokens, &CCLToken{
				Type:   TokenTypeComment,
				Value:  comment,
				Line:   line,
				Column: column - len([]rune(comment)),
			})
			continue
		}

		// simple token means it's a single character token
		// like : ; { } [ ]
		theTokenType, isSimpleToken := oneCharSimpleTokens[r]
		if isSimpleToken {
			tokens = append(tokens, &CCLToken{
				Type:   theTokenType,
				Value:  string(r),
				Line:   line,
				Column: column,
			})
			pos++
			column++
			continue
		}

		if r == '"' {
			// String Literal
			literal, newPos, newCol, err := parseStringLiteral(runes, pos, line, column)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, literal)
			pos = newPos
			column = newCol
			continue
		}

		// Identifier or Keyword or DataType
		identifier := ""
		startCol := column
		for pos < totalRunesLen && isAlphaNumeric(runes[pos]) {
			identifier += string(runes[pos])
			pos++
			column++
		}

		if identifier != "" {
			tokenType := TokenTypeIdentifier
			if cclValues.IsKeywordName(identifier) {
				tokenType = TokenTypeKeywordModel
				identifier = cclValues.GetNormalizedKeywordName(identifier)
			} else if cclValues.IsTypeName(identifier) {
				tokenType = TokenTypeDataType
				identifier = cclValues.GetNormalizedTypeName(identifier)
			}

			tokens = append(tokens, &CCLToken{
				Type:   tokenType,
				Value:  identifier,
				Line:   line,
				Column: startCol,
			})
			continue
		}

		return nil, &UnexpectedCharacterError{
			Character: r,
			Line:      line,
			Column:    column,
		}
	}

	return tokens, nil
}

func parseStringLiteral(runes []rune, pos, line, column int) (*CCLToken, int, int, error) {
	literal := "\""
	pos++
	column++
	startCol := column - 1
	for pos < len(runes) && runes[pos] != '"' {
		literal += string(runes[pos])
		pos++
		column++
	}

	if pos >= len(runes) {
		return nil, 0, 0, &UnexpectedEndOfStringLiteralError{
			Line:   line,
			Column: startCol,
		}
	}
	literal += "\""
	pos++
	column++

	return &CCLToken{
		Type:   TokenTypeStringLiteral,
		Value:  literal,
		Line:   line,
		Column: startCol,
	}, pos, column, nil
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
