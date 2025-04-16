package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

//---------------------------------------------------------

// CCLParseOptions is a struct that contains the options for
// parsing a CCL source file.
type CCLParseOptions struct {
	SourceFilePath string
	SourceContent  string
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

	// codeDefinition is the definition that this parser is building.
	codeDefinition *cclValues.SourceCodeDefinition
}

//---------------------------------------------------------
