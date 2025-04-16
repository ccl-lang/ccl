package cclParser

import (
	"fmt"
	"strings"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
)

//---------------------------------------------------------

type UnexpectedTokenError struct {
	Expected cclLexer.CCLTokenType
	Actual   cclLexer.CCLTokenType
	Line     int
	Column   int
}

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf(
		"cclParser: expected token type %s, got %s at line %d, column %d",
		e.Expected,
		e.Actual,
		e.Line,
		e.Column,
	)
}

//---------------------------------------------------------

type ExpectedValueError struct {
	SourceLine string
	ParamName  string
	Line       int
	Column     int
}

func (e *ExpectedValueError) Error() string {
	// Create the main error message
	message := fmt.Sprintf(
		"cclParser: Expected a value after '=' in attribute parameter\n  at line %d, column %d\n",
		e.Line,
		e.Column,
	)

	// Add the source code line
	message += "  " + e.SourceLine + "\n"

	// Add the pointer to the exact position
	// First, calculate how many spaces to add before the ^
	pointerIndent := "  " + strings.Repeat(" ", e.Column)
	message += pointerIndent + "^ Missing value here"

	return message
}

//---------------------------------------------------------

type UnexpectedTokenAfterParameterError struct {
	Line       int
	Column     int
	SourceLine string
	ParamName  string
	TokenValue string
}

func (e *UnexpectedTokenAfterParameterError) Error() string {
	message := fmt.Sprintf(
		"Error: Unexpected token '%s' after parameter value\n  at line %d, column %d\n",
		e.TokenValue,
		e.Line,
		e.Column,
	)

	// Add the source code line
	message += "  " + e.SourceLine + "\n"

	// Add the pointer to the exact position
	pointerIndent := "  " + strings.Repeat(" ", e.Column)
	message += pointerIndent + "^ Expected ',' or ')' here"

	message += "\n\nHint: Parameter values should be separated by commas like:"
	message += "\n  [MyAttribute(Param1 = \"value\", Param2 = 123)]"

	return message
}

// ---------------------------------------------------------

type UnexpectedTokenAfterAssignmentError struct {
	Line       int
	Column     int
	SourceLine string
	ParamName  string
	TokenValue string
}

func (e *UnexpectedTokenAfterAssignmentError) Error() string {
	message := fmt.Sprintf(
		"Error: Unexpected token '%s' after assignment\n  at line %d, column %d\n",
		e.TokenValue,
		e.Line,
		e.Column,
	)

	// Add the source code line
	message += "  " + e.SourceLine + "\n"

	// Add the pointer to the exact position
	pointerIndent := "  " + strings.Repeat(" ", e.Column)
	message += pointerIndent + "^ Expected literal value or variable here"

	message += "\n\nHint: After assignment, there should be a literal value or variable like this:"
	message += "\n  [MyAttribute(Param1 = \"value\", Param2 = ConstantValue)]"

	return message
}

//---------------------------------------------------------

type UnexpectedEndOfAttributeError struct {
	Line   int
	Column int
}

func (e *UnexpectedEndOfAttributeError) Error() string {
	return fmt.Sprintf(
		"cclParser: unexpected end of attribute at line %d, column %d",
		e.Line,
		e.Column,
	)
}

//---------------------------------------------------------
