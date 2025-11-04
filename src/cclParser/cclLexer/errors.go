package cclLexer

import "fmt"

//---------------------------------------------------------

// UnexpectedCharacterError represents an error when an unexpected character is found.
type UnexpectedCharacterError struct {
	Character  rune
	Line       int
	Column     int
	InnerError error
}

// Error returns the string representation of the unexpected character error.
func (e *UnexpectedCharacterError) Error() string {
	errStr := fmt.Sprintf("cclParser: unexpected character '%c' at line %d, column %d", e.Character,
		e.Line, e.Column)

	if e.InnerError != nil {
		errStr += fmt.Sprintf(": %s", e.InnerError.Error())
	}

	return errStr
}

//---------------------------------------------------------

// UnexpectedEndOfAttributeError represents an error when an unexpected end of
// attribute is found.
type UnexpectedEndOfAttributeError struct {
	Line   int
	Column int
}

// Error returns the string representation of the unexpected end of attribute error.
func (e *UnexpectedEndOfAttributeError) Error() string {
	return fmt.Sprintf("unexpected end of attribute at line %d, column %d", e.Line, e.Column)
}

//---------------------------------------------------------

// UnexpectedEndOfStringLiteralError represents an error when an unexpected end of
// string literal is found.
type UnexpectedEndOfStringLiteralError struct {
	Line   int
	Column int
}

// Error returns the string representation of the unexpected end of string literal error.
func (e *UnexpectedEndOfStringLiteralError) Error() string {
	return fmt.Sprintf("unexpected end of string literal at line %d, column %d", e.Line, e.Column)
}

//---------------------------------------------------------
