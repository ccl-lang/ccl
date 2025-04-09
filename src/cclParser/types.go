package cclParser

import "github.com/ccl-lang/ccl/src/cclParser/cclLexer"

//---------------------------------------------------------

// CCLParseOptions is a struct that contains the options for
// parsing a CCL source file.
type CCLParseOptions struct {
	Source string
}

type CCLParser struct {
	// Options is the options for the parser.
	Options *CCLParseOptions

	// tokens is the list of tokens that the parser will parse.
	tokens []*cclLexer.CCLToken

	// pos is the current position of the parser in the tokens list.
	pos int

	// current is the current token that the parser is parsing.
	current *cclLexer.CCLToken
}

//---------------------------------------------------------

type UnexpectedTokenError struct {
	Expected cclLexer.CCLTokenType
	Actual   cclLexer.CCLTokenType
	Line     int
	Column   int
}

type UnexpectedEndOfAttributeError struct {
	Line   int
	Column int
}

//---------------------------------------------------------
