package cclParser

import (
	"fmt"
	"strings"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

type UnexpectedTokenError struct {
	Expected   cclLexer.CCLTokenType
	Actual     cclLexer.CCLTokenType
	Line       int
	Column     int
	SourceLine string
}

func (e *UnexpectedTokenError) Error() string {
	if e.SourceLine == "" {
		return fmt.Sprintf(
			"cclParser: expected token type %s, got %s at line %d, column %d",
			e.Expected,
			e.Actual,
			e.Line,
			e.Column,
		)
	}

	message := fmt.Sprintf(
		"Error: expected token type %s, got %s\n  at line %d, column %d\n",
		e.Expected,
		e.Actual,
		e.Line,
		e.Column,
	)

	// Add the source code line
	message += "  " + e.SourceLine + "\n"
	// Add the pointer to the exact position
	pointerIndent := "  " + strings.Repeat(" ", e.Column)
	message += pointerIndent + "^ Expected token '" + e.Expected.String() + "' here\n"

	return message
}

//---------------------------------------------------------

type UnexpectedEOFError struct {
	Line       int
	Column     int
	SourceLine string
}

func (e *UnexpectedEOFError) Error() string {
	if e.SourceLine == "" {
		return fmt.Sprintf(
			"cclParser: Unexpected EOF at line %d, column %d",
			e.Line,
			e.Column,
		)
	}

	message := fmt.Sprintf(
		"cclParser: Unexpected EOF at line %d, column %d",
		e.Line,
		e.Column,
	)

	// Add the source code line
	message += "  " + e.SourceLine + "\n"
	// Add the pointer to the exact position
	pointerIndent := "  " + strings.Repeat(" ", e.Column)
	message += pointerIndent + "^ Unexpected EOF\n"

	return message
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

	targetParam := "Param2"
	if e.ParamName != "" {
		targetParam = e.ParamName
	}

	message += "\n\nHint: After assignment, there should be a literal value or variable like this:"
	message += "\n  [MyAttribute(Param1 = \"value\", " + targetParam + " = ConstantValue)]"

	return message
}

//---------------------------------------------------------

type UnexpectedEndOfAttributeError struct {
	Line       int
	Column     int
	SourceLine string
}

func (e *UnexpectedEndOfAttributeError) Error() string {
	if e.SourceLine == "" {
		return fmt.Sprintf(
			"cclParser: unexpected end of attribute at line %d, column %d",
			e.Line,
			e.Column,
		)
	}

	message := fmt.Sprintf(
		"Error: Unexpected end of attribute\n  at line %d, column %d\n",
		e.Line,
		e.Column,
	)

	// Add the source code line
	message += "  " + e.SourceLine + "\n"

	// Add the pointer to the exact position
	pointerIndent := "  " + strings.Repeat(" ", e.Column)
	message += pointerIndent + "^ Expected end of attribute and an entity after this"

	message += "\n\nHint: After attributes, we expect an entity (such as a model, or a field) like this:"
	message += "\n  [MyAttribute(Param1 = \"value\", )]"
	message += "\n  model MyModel { [MaxLength(10)] myField: string; }"

	return message
}

//---------------------------------------------------------

type InvalidAttributeUsageError struct {
	Line       int
	Column     int
	SourceLine string
}

func (e *InvalidAttributeUsageError) Error() string {
	if e.SourceLine == "" {
		return fmt.Sprintf(
			"cclParser: invalid attribute usage at line %d, column %d",
			e.Line,
			e.Column,
		)
	}

	message := fmt.Sprintf(
		"Error: Invalid attribute usage\n  at line %d, column %d\n",
		e.Line,
		e.Column,
	)

	// Add the source code line
	message += "  " + e.SourceLine + "\n"

	// Add the pointer to the exact position
	pointerIndent := "  " + strings.Repeat(" ", e.Column)
	message += pointerIndent + "^ Expected a valid attribute usage here"

	message += "\n\nHint: Attributes can only be applied on models, fields, etc..."
	message += "\n  If you don't want to do that, consider converting this to a global attribute."

	return message
}

//---------------------------------------------------------

type InvalidSyntaxError struct {
	Line        int
	Column      int
	SourceLine  string
	Language    globalValues.LanguageType
	HintMessage string
}

func (e *InvalidSyntaxError) Error() string {
	if e.SourceLine == "" {
		return fmt.Sprintf(
			"cclParser: invalid "+e.Language.String()+" syntax at line %d, column %d",
			e.Line,
			e.Column,
		)
	}

	message := fmt.Sprintf(
		"Error: Invalid "+e.Language.String()+" syntax\n  at line %d, column %d\n",
		e.Line,
		e.Column,
	)

	// Add the source code line
	message += "  " + e.SourceLine + "\n"

	// Add the pointer to the exact position
	pointerIndent := "  " + strings.Repeat(" ", e.Column)
	message += pointerIndent + "^ Invalid " + e.Language.String() + " syntax"

	if e.HintMessage != "" {
		message += "\n\nHint: " + e.HintMessage
	}

	return message
}

//---------------------------------------------------------

type UndefinedIdentifierError struct {
	Line             int
	Column           int
	SourceLine       string
	TargetIdentifier string
	Language         globalValues.LanguageType
}

func (e *UndefinedIdentifierError) Error() string {
	message := fmt.Sprintf(
		e.Language.String()+"Parser error: undefined identifier '"+e.TargetIdentifier+"' at line %d, column %d",
		e.Line,
		e.Column,
	)

	if e.SourceLine == "" {
		return message
	}

	// Add the source code line
	message += "  " + e.SourceLine + "\n"

	// Add the pointer to the exact position
	pointerIndent := "  " + strings.Repeat(" ", e.Column)
	message += pointerIndent + "^ Undefined identifier '" + e.TargetIdentifier + "' "

	return message
}

//---------------------------------------------------------
