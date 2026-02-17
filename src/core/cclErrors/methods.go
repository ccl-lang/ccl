package cclErrors

import "strings"

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
	return d.SourcePosition.FormatError(message)
}

//---------------------------------------------------------

func (d *DuplicateModelError) Error() string {
	if d == nil {
		return "Duplicate model"
	}

	message := "Duplicate model: " + d.ModelName
	return d.SourcePosition.FormatError(message)
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
		"' conflicts with " + string(e.Kind) +
		" name '" + e.ConflictName + "'"

	if e.Namespace != "" {
		message += " in namespace '" + e.Namespace + "'"
	}

	return e.SourcePosition.FormatError(message)
}
