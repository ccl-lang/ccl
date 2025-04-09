package cclLexer

//---------------------------------------------------------

// CCLTokenType represents the type of a token.
type CCLTokenType int

// singleTokenContainer is a map of runes to CCLTokenType.
type singleTokenContainer = map[rune]CCLTokenType

// doubleTokenContainer is a map of runes to another map of runes to CCLTokenType.
type doubleTokenContainer = map[rune]singleTokenContainer

// CCLToken represents a token.
type CCLToken struct {
	// Type is the type of the token.
	Type CCLTokenType

	// Value is the value that this token is holding.
	value any

	// Line is the line number where this token is found.
	Line int

	// Column is the column number where this token is found.
	Column int
}

//---------------------------------------------------------

// UnexpectedCharacterError represents an error when an unexpected character is found.
type UnexpectedCharacterError struct {
	Character  rune
	Line       int
	Column     int
	InnerError error
}

// UnexpectedEndOfAttributeError represents an error when an unexpected end of
// attribute is found.
type UnexpectedEndOfAttributeError struct {
	Line   int
	Column int
}

// UnexpectedEndOfStringLiteralError represents an error when an unexpected end of
// string literal is found.
type UnexpectedEndOfStringLiteralError struct {
	Line   int
	Column int
}

//---------------------------------------------------------
