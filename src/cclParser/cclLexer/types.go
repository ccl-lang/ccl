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
