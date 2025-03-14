package cclLexer

import "fmt"

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

func (t *CCLToken) String() string {
	return t.Value + " -> " + t.Type.String()
}

//---------------------------------------------------------
//---------------------------------------------------------

// Error returns the string representation of the unexpected character error.
func (e *UnexpectedCharacterError) Error() string {
	return fmt.Sprintf("unexpected character '%c' at line %d, column %d", e.Character,
		e.Line, e.Column)
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
