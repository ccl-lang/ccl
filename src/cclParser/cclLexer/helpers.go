package cclLexer

import (
	"strconv"
	"unicode"

	"github.com/ALiwoto/ssg/ssg/rangeValues"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func Lex(input string) ([]*CCLToken, error) {
	var tokens []*CCLToken
	line := 1
	column := 1
	runes := []rune(input)
	pos := 0
	totalRunesLen := len(runes)

	for pos < totalRunesLen {
		if isWhitespace(runes[pos]) {
			isNextLF := pos+1 < totalRunesLen && runes[pos+1] == '\n'
			if runes[pos] == '\n' || runes[pos] == '\r' {
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

		if runes[pos] == '/' && pos+1 < totalRunesLen && runes[pos+1] == '/' {
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
				value:  comment,
				Line:   line,
				Column: column - len([]rune(comment)),
			})
			continue
		}

		// two character token
		if pos+1 < totalRunesLen {
			secondMap, isTwoCharToken := twoCharsSimpleTokens[runes[pos]]
			if isTwoCharToken {
				theTokenType, isTwoCharToken := secondMap[runes[pos+1]]
				if isTwoCharToken {
					tokens = append(tokens, &CCLToken{
						Type:   theTokenType,
						value:  string(runes[pos]) + string(runes[pos+1]),
						Line:   line,
						Column: column,
					})
					pos += 2
					column += 2
					continue
				}
			}
		}

		// simple token means it's a single character token
		// like : ; { } [ ]
		theTokenType, isSimpleToken := oneCharSimpleTokens[runes[pos]]
		if isSimpleToken {
			tokens = append(tokens, &CCLToken{
				Type:   theTokenType,
				value:  string(runes[pos]),
				Line:   line,
				Column: column,
			})
			pos++
			column++
			continue
		}

		if runes[pos] == '"' {
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

		// integer literal, since an identifier cannot start with a digit
		// so we can safely assume that this is an integer literal
		if unicode.IsDigit(runes[pos]) {
			literal := string(runes[pos])
			startCol := column
			isFloat := false
			integerBase := 10
			pos++
			column++
			for pos < totalRunesLen {
				// check for digits
				if unicode.IsDigit(runes[pos]) {
					literal += string(runes[pos])
					pos++
					column++
					continue
				}

				if isWhitespace(runes[pos]) {
					// fast break to avoid unnecessary checks
					break
				}

				// ok so it's not digit, maybe a dot?
				if runes[pos] == '.' {
					// float literal
					if isFloat || integerBase != 10 {
						return nil, &UnexpectedCharacterError{
							Character: runes[pos],
							Line:      line,
							Column:    column,
						}
					}
					isFloat = true
					literal += "."
					pos++
					column++
					continue
				}

				// not dot and digit, maybe a base?
				if runes[pos] == 'x' || runes[pos] == 'X' {
					if integerBase != 10 || isFloat || len(literal) < 1 || literal[0] != '0' {
						// hmmm we have already switched bases?
						// yeah, that's not allowed
						return nil, &UnexpectedCharacterError{
							Character: runes[pos],
							Line:      line,
							Column:    column,
						}
					}

					for pos < totalRunesLen && isHexDigit(runes[pos]) {
						literal += string(runes[pos])
						pos++
						column++
					}

					currentInteger, err := strconv.ParseInt(literal, 10, 64)
					if err != nil || currentInteger != 0 {
						// why is currentInteger != 0?
						// because before X or x, we should only have 0
						// basically like 0x1234
						return nil, &UnexpectedCharacterError{
							Character:  runes[pos],
							Line:       line,
							Column:     column,
							InnerError: err,
						}
					}

					integerBase = 16
					literal += string(runes[pos])
					pos++
					column++
					continue
				}

				// maybe a binary base (0b or 0B)?
				if runes[pos] == 'b' || runes[pos] == 'B' {
					if integerBase != 10 || isFloat || len(literal) < 1 || literal[0] != '0' {
						// hmmm we have already switched bases?
						// yeah, that's not allowed
						return nil, &UnexpectedCharacterError{
							Character: runes[pos],
							Line:      line,
							Column:    column,
						}
					}

					for pos < totalRunesLen && (runes[pos] == '0' || runes[pos] == '1') {
						literal += string(runes[pos])
						pos++
						column++
					}

					currentInteger, err := strconv.ParseInt(literal, 10, 64)
					if err != nil || currentInteger != 0 {
						// why is currentInteger != 0?
						// because before B or b, we should only have 0
						// basically like 0b1010
						return nil, &UnexpectedCharacterError{
							Character:  runes[pos],
							Line:       line,
							Column:     column,
							InnerError: err,
						}
					}

					integerBase = 2
					literal += string(runes[pos])
					pos++
					column++
					continue
				}

				// maybe an octal base (0o or 0O)?
				if runes[pos] == 'o' || runes[pos] == 'O' {
					if integerBase != 10 || isFloat || len(literal) < 1 || literal[0] != '0' {
						// hmmm we have already switched bases?
						// yeah, that's not allowed
						return nil, &UnexpectedCharacterError{
							Character: runes[pos],
							Line:      line,
							Column:    column,
						}
					}

					for pos < totalRunesLen && runes[pos] >= '0' && runes[pos] <= '7' {
						literal += string(runes[pos])
						pos++
						column++
					}

					currentInteger, err := strconv.ParseInt(literal, 10, 64)
					if err != nil || currentInteger != 0 {
						// why is currentInteger != 0?
						// because before o or O, we should only have 0
						// basically like 0o1234
						return nil, &UnexpectedCharacterError{
							Character:  runes[pos],
							Line:       line,
							Column:     column,
							InnerError: err,
						}
					}

					integerBase = 8
					literal += string(runes[pos])
					pos++
					column++
					continue
				}

				break
			}

			if isFloat {
				correctFloat, err := strconv.ParseFloat(literal, 64)
				if err != nil {
					// maybe change this error to invalid float or something?
					return nil, &UnexpectedCharacterError{
						Character:  runes[pos],
						Line:       line,
						Column:     column,
						InnerError: err,
					}
				}
				tokens = append(tokens, &CCLToken{
					Type:   TokenTypeFloatLiteral,
					value:  correctFloat,
					Line:   line,
					Column: startCol,
				})
				continue
			}

			if integerBase != 10 {
				// remove the first two characters (0x, 0b, 0o)
				literal = literal[2:]
			}

			correctInteger, err := strconv.ParseInt(literal, integerBase, 64)
			if err != nil {
				// maybe change this error to invalid integer or something?
				return nil, &UnexpectedCharacterError{
					Character:  runes[pos],
					Line:       line,
					Column:     column,
					InnerError: err,
				}
			}
			tokens = append(tokens, &CCLToken{
				Type:   TokenTypeIntLiteral,
				value:  correctInteger,
				Line:   line,
				Column: startCol,
			})
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
				value:  identifier,
				Line:   line,
				Column: startCol,
			})
			continue
		}

		return nil, &UnexpectedCharacterError{
			Character: runes[pos],
			Line:      line,
			Column:    column,
		}
	}

	return tokens, nil
}

func parseStringLiteral(runes []rune, pos, line, column int) (*CCLToken, int, int, error) {
	literal := ""
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
	pos++
	column++

	return &CCLToken{
		Type:   TokenTypeStringLiteral,
		value:  literal,
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

func isHexDigit(r rune) bool {
	return unicode.IsDigit(r) ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}

// toIntegerValueT converts the given value to the integer value of type T.
// If the value is not an integer, it returns 0.
// NOTE: This function should be moved to ssg package in future.
func toIntegerValueT[T rangeValues.Integer](value any) T {
	switch correctValue := value.(type) {
	case T:
		return correctValue
	case int8:
		return T(correctValue)
	case int:
		return T(correctValue)
	case int16:
		return T(correctValue)
	case int32:
		return T(correctValue)
	case int64:
		return T(correctValue)
	case uint8:
		return T(correctValue)
	case uint:
		return T(correctValue)
	case uint16:
		return T(correctValue)
	case uint32:
		return T(correctValue)
	case uint64:
		return T(correctValue)
	}

	return 0
}
