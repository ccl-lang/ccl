package cclLexer

import (
	"fmt"
)

//---------------------------------------------------------

// ToString returns the string representation of the token type.
func (t CCLTokenType) ToString() string {
	return tokenTypesToNames[t]
}

// String returns the string representation of the token type.
func (t CCLTokenType) String() string {
	return tokenTypesToNames[t]
}

//---------------------------------------------------------

// String returns the string representation of the token.
func (t *CCLToken) String() string {
	return t.FormatValueAsString() + " -> " + t.Type.String()
}

// FormatValueAsString returns the formatted value of the token as a string.
func (t *CCLToken) FormatValueAsString() string {
	return fmt.Sprintf("%v", t.value)
}

// GetStringLiteral returns the string literal value of the token.
// If the token is not a string literal, it returns an empty string.
func (t *CCLToken) GetStringLiteral() string {
	if t.Type == TokenTypeStringLiteral {
		return t.value.(string)
	}

	return ""
}

// GetComment returns the comment value of the token.
// If the token is not a comment, it returns an empty string.
func (t *CCLToken) GetComment() string {
	if t.Type == TokenTypeComment {
		return t.value.(string)
	}

	return ""
}

// GetIntLiteral returns the integer literal value of the token.
// If the token is not an integer literal, it returns 0.
func (t *CCLToken) GetIntLiteral() int {
	if t.Type == TokenTypeIntLiteral {
		return toIntegerValueT[int](t.value)
	}

	return 0
}

// GetFloatLiteral returns the float literal value of the token.
// If the token is not a float literal, it returns 0.
func (t *CCLToken) GetFloatLiteral() float64 {
	if t.Type == TokenTypeFloatLiteral {
		return t.value.(float64)
	}

	return 0
}

// GetIdentifier returns the identifier value of the token.
// If the token is not an identifier, it returns an empty string.
func (t *CCLToken) GetIdentifier() string {
	if t.Type == TokenTypeIdentifier {
		return t.value.(string)
	}

	return ""
}

//---------------------------------------------------------
//---------------------------------------------------------

// Error returns the string representation of the unexpected character error.
func (e *UnexpectedCharacterError) Error() string {
	errStr := fmt.Sprintf("cclParser: unexpected character '%c' at line %d, column %d", e.Character,
		e.Line, e.Column)

	if e.InnerError != nil {
		errStr += fmt.Sprintf(": %s", e.InnerError.Error())
	}

	return errStr
}

// Error returns the string representation of the unexpected end of attribute error.
func (e *UnexpectedEndOfAttributeError) Error() string {
	return fmt.Sprintf("unexpected end of attribute at line %d, column %d", e.Line, e.Column)
}

// Error returns the string representation of the unexpected end of string literal error.
func (e *UnexpectedEndOfStringLiteralError) Error() string {
	return fmt.Sprintf("unexpected end of string literal at line %d, column %d", e.Line, e.Column)
}

//---------------------------------------------------------
