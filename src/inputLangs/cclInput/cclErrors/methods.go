package cclErrors

//---------------------------------------------------------

func (v *ValidationError) Error() string {
	return v.Message
}

//---------------------------------------------------------

func (d *DuplicateFieldError) Error() string {
	message := "Duplicate field: " + d.ModelName + "." + d.FieldName
	return d.SourcePosition.FormatError(message)
}

//---------------------------------------------------------

func (d *DuplicateModelError) Error() string {
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
	message := "Unsupported file naming style: " + u.StyleName +
		" for model " + u.ModelName +
		". Supported styles are: [" + u.SupportedStyles + "]" +
		" when compiling to " + u.TargetLanguage

	return u.SourcePosition.FormatError(message)
}

//---------------------------------------------------------

func (e *FieldNameConflictError) Error() string {
	message := "Field name '" + e.FieldName +
		"' in model '" + e.ModelName +
		"' conflicts with " + string(e.Kind) +
		" name '" + e.ConflictName + "'"

	if e.Namespace != "" {
		message += " in namespace '" + e.Namespace + "'"
	}

	return e.SourcePosition.FormatError(message)
}

//---------------------------------------------------------

func (e *InvalidAttributeUsageError) Error() string {
	message := "Invalid attribute usage for '" +
		e.AttrName.ToString() + "'"

	if e.Message != "" {
		message += ": " + e.Message
	}

	return e.SourcePosition.FormatError(message)
}
