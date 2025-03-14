package cclLexer

//---------------------------------------------------------

// CCLTokenType represents the type of a token.
type CCLTokenType int

// CCLToken represents a token.
type CCLToken struct {
	Type   CCLTokenType
	Value  string
	Line   int
	Column int
}

//---------------------------------------------------------

// UnexpectedCharacterError represents an error when an unexpected character is found.
type UnexpectedCharacterError struct {
	Character rune
	Line      int
	Column    int
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
