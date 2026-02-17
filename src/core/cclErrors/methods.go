package cclErrors

import (
	"fmt"
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclUtils"
)

//---------------------------------------------------------

func (v *ValidationError) Error() string {
	return v.Message
}

//---------------------------------------------------------

func (d *DuplicateFieldError) Error() string {
	if d == nil {
		return "Duplicate field"
	}

	message := "Duplicate field: " + d.ModelName + "." + d.FieldName
	return formatErrorWithSourcePosition(message, d.SourcePosition)
}

//---------------------------------------------------------

func (d *DuplicateModelError) Error() string {
	if d == nil {
		return "Duplicate model"
	}

	message := "Duplicate model: " + d.ModelName
	return formatErrorWithSourcePosition(message, d.SourcePosition)
}

//---------------------------------------------------------

func (u *UnsupportedFieldTypeError) Error() string {
	return "Unsupported field type: " + u.TypeName +
		" for field " + u.FieldName +
		" in model " + u.ModelName +
		" when compiling to " + u.TargetLanguage
}

//---------------------------------------------------------

func (u *UnsupportedTypeDefinitionError) Error() string {
	return "Unsupported type definition: " + u.TypeName +
		" when compiling to " + u.TargetLanguage
}

//---------------------------------------------------------

func (u *UnsupportedFileNamingStyleError) Error() string {
	return "Unsupported file naming style: " + u.StyleName +
		" for model " + u.ModelName +
		". Supported styles are: [" + strings.Join(u.SupportedStyles, ", ") + "]" +
		" when compiling to " + u.TargetLanguage
}

//---------------------------------------------------------

func (e *FieldNameConflictError) Error() string {
	if e == nil {
		return "Field name conflicts with reserved name"
	}

	message := "Field name '" + e.FieldName +
		"' in model '" + e.ModelName +
		"' conflicts with " + e.ConflictKind +
		" name '" + e.ConflictName + "'"

	if e.Namespace != "" {
		message += " in namespace '" + e.Namespace + "'"
	}

	return formatErrorWithSourcePosition(message, e.SourcePosition)
}

//---------------------------------------------------------

func formatErrorWithSourcePosition(message string, pos *cclUtils.SourceCodePosition) string {
	if pos == nil {
		return message
	}

	if pos.SourceLine == "" {
		return fmt.Sprintf(
			"%s at line %d, column %d",
			message,
			pos.Line,
			pos.Column,
		)
	}

	result := fmt.Sprintf(
		"Error: %s\n  at line %d, column %d\n",
		message,
		pos.Line,
		pos.Column,
	)

	result += "  " + pos.SourceLine + "\n"
	pointerIndent := "  " + strings.Repeat(" ", pos.Column)
	result += pointerIndent + "^ " + message

	return result
}
