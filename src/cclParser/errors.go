package cclParser

import (
	"fmt"
	"strings"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclUtils"
	"github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

type UnexpectedTokenError struct {
	Expected cclLexer.CCLTokenType
	Actual   cclLexer.CCLTokenType

	SourcePosition *cclUtils.SourceCodePosition
}

func (e *UnexpectedTokenError) Error() string {
	if e.SourcePosition == nil {
		return fmt.Sprintf(
			"cclParser: expected token type %s, got %s",
			e.Expected,
			e.Actual,
		)
	}

	if e.SourcePosition.SourceLine == "" {
		return fmt.Sprintf(
			"cclParser: expected token type %s, got %s at line %d, column %d",
			e.Expected,
			e.Actual,
			e.SourcePosition.Line,
			e.SourcePosition.Column,
		)
	}

	message := fmt.Sprintf(
		"Error: expected token type %s, got %s\n  at line %d, column %d\n",
		e.Expected,
		e.Actual,
		e.SourcePosition.Line,
		e.SourcePosition.Column,
	)

	// Add the source code line
	message += "  " + e.SourcePosition.SourceLine + "\n"
	// Add the pointer to the exact position
	pointerIndent := "  " + strings.Repeat(" ", e.SourcePosition.Column)
	message += pointerIndent + "^ Expected token '" + e.Expected.String() + "' here\n"

	return message
}

//---------------------------------------------------------

type UnexpectedEOFError struct {
	SourcePosition *cclUtils.SourceCodePosition
}

func (e *UnexpectedEOFError) Error() string {
	if e.SourcePosition == nil {
		return "cclParser: Unexpected EOF"
	}

	if e.SourcePosition.SourceLine == "" {
		return fmt.Sprintf(
			"cclParser: Unexpected EOF at line %d, column %d",
			e.SourcePosition.Line,
			e.SourcePosition.Column,
		)
	}

	message := fmt.Sprintf(
		"cclParser: Unexpected EOF at line %d, column %d",
		e.SourcePosition.Line,
		e.SourcePosition.Column,
	)

	// Add the source code line
	message += "  " + e.SourcePosition.SourceLine + "\n"
	// Add the pointer to the exact position
	pointerIndent := "  " + strings.Repeat(" ", e.SourcePosition.Column)
	message += pointerIndent + "^ Unexpected EOF\n"

	return message
}

//---------------------------------------------------------

type ExpectedValueError struct {
	ParamName string

	SourcePosition *cclUtils.SourceCodePosition
}

func (e *ExpectedValueError) Error() string {
	if e.SourcePosition == nil {
		return "cclParser: Expected a value after '=' in attribute parameter"
	}

	// Create the main error message
	return e.SourcePosition.FormatError("Missing value here")
}

//---------------------------------------------------------

type UnexpectedTokenAfterParameterError struct {
	ParamName  string
	TokenValue string

	SourcePosition *cclUtils.SourceCodePosition
}

func (e *UnexpectedTokenAfterParameterError) Error() string {
	if e.SourcePosition == nil {
		return fmt.Sprintf(
			"cclParser: unexpected token '%s' after parameter value",
			e.TokenValue,
		)
	}

	message := fmt.Sprintf(
		"Error: Unexpected token '%s' after parameter value\n  at line %d, column %d\n",
		e.TokenValue,
		e.SourcePosition.Line,
		e.SourcePosition.Column,
	)

	message += "\n\nHint: Parameter values should be separated by commas like:"
	message += "\n  [MyAttribute(Param1 = \"value\", Param2 = 123)]"

	return e.SourcePosition.FormatError(message)
}

// ---------------------------------------------------------

type UnexpectedTokenAfterAssignmentError struct {
	ParamName  string
	TokenValue string

	SourcePosition *cclUtils.SourceCodePosition
}

func (e *UnexpectedTokenAfterAssignmentError) Error() string {
	if e.SourcePosition == nil {
		return fmt.Sprintf(
			"cclParser: unexpected token '%s' after assignment",
			e.TokenValue,
		)
	}

	message := fmt.Sprintf(
		"Error: Unexpected token '%s' after assignment\n  at line %d, column %d\n",
		e.TokenValue,
		e.SourcePosition.Line,
		e.SourcePosition.Column,
	)

	targetParam := "Param2"
	if e.ParamName != "" {
		targetParam = e.ParamName
	}

	message += "\n\nHint: After assignment, there should be a literal value or variable like this:"
	message += "\n  [MyAttribute(Param1 = \"value\", " + targetParam + " = ConstantValue)]"

	return e.SourcePosition.FormatError(message)
}

//---------------------------------------------------------

type UnexpectedEndOfAttributeError struct {
	SourcePosition *cclUtils.SourceCodePosition
}

func (e *UnexpectedEndOfAttributeError) Error() string {
	if e.SourcePosition == nil {
		return "cclParser: unexpected end of attribute"
	}

	if e.SourcePosition.SourceLine == "" {
		return fmt.Sprintf(
			"cclParser: unexpected end of attribute at line %d, column %d",
			e.SourcePosition.Line,
			e.SourcePosition.Column,
		)
	}

	message := fmt.Sprintf(
		"Error: Unexpected end of attribute\n  at line %d, column %d\n",
		e.SourcePosition.Line,
		e.SourcePosition.Column,
	)

	message += "\n\nHint: After attributes, we expect an entity (such as a model, or a field) like this:"
	message += "\n  [MyAttribute(Param1 = \"value\", )]"
	message += "\n  model MyModel { [MaxLength(10)] myField: string; }"

	return e.SourcePosition.FormatError(message)
}

//---------------------------------------------------------

type InvalidAttributeUsageError struct {
	SourcePosition *cclUtils.SourceCodePosition
}

func (e *InvalidAttributeUsageError) Error() string {
	if e.SourcePosition == nil {
		return "cclParser: invalid attribute usage"
	}

	if e.SourcePosition.SourceLine == "" {
		return fmt.Sprintf(
			"cclParser: invalid attribute usage at line %d, column %d",
			e.SourcePosition.Line,
			e.SourcePosition.Column,
		)
	}

	message := fmt.Sprintf(
		"Error: Invalid attribute usage\n  at line %d, column %d\n",
		e.SourcePosition.Line,
		e.SourcePosition.Column,
	)

	message += "\n\nHint: Attributes can only be applied on models, fields, etc..."
	message += "\n  If you don't want to do that, consider converting this to a global attribute."

	return e.SourcePosition.FormatError(message)
}

//---------------------------------------------------------

type InvalidSyntaxError struct {
	Language    globalValues.LanguageType
	HintMessage string

	SourcePosition *cclUtils.SourceCodePosition
}

func (e *InvalidSyntaxError) Error() string {
	if e.SourcePosition == nil {
		return fmt.Sprintf(
			"cclParser: invalid " + e.Language.String() + " syntax",
		)
	}

	if e.SourcePosition.SourceLine == "" {
		return fmt.Sprintf(
			"cclParser: invalid "+e.Language.String()+" syntax at line %d, column %d",
			e.SourcePosition.Line,
			e.SourcePosition.Column,
		)
	}

	message := fmt.Sprintf(
		"Error: Invalid "+e.Language.String()+" syntax\n  at line %d, column %d\n",
		e.SourcePosition.Line,
		e.SourcePosition.Column,
	)

	if e.HintMessage != "" {
		message += "\n\nHint: " + e.HintMessage
	}

	return e.SourcePosition.FormatError(message)
}

//---------------------------------------------------------

type UndefinedIdentifierError struct {
	TargetIdentifier string
	Language         globalValues.LanguageType

	SourcePosition *cclUtils.SourceCodePosition
}

func (e *UndefinedIdentifierError) Error() string {
	if e.SourcePosition == nil {
		return fmt.Sprintf(
			"cclParser: undefined identifier '%s'",
			e.TargetIdentifier,
		)
	}

	message := fmt.Sprintf(
		e.Language.String()+"Parser error: undefined identifier '"+e.TargetIdentifier+"' at line %d, column %d",
		e.SourcePosition.Line,
		e.SourcePosition.Column,
	)

	if e.SourcePosition.SourceLine == "" {
		return message
	}

	return e.SourcePosition.FormatError(message)
}

//---------------------------------------------------------
