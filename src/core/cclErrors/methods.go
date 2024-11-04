package cclErrors

//---------------------------------------------------------

func (v *ValidationError) Error() string {
	return v.Message
}

//---------------------------------------------------------

func (d *DuplicateFieldError) Error() string {
	return "Duplicate field: " + d.ModelName + "." + d.FieldName
}

//---------------------------------------------------------

func (d *DuplicateModelError) Error() string {
	return "Duplicate model: " + d.ModelName
}

//---------------------------------------------------------

func (u *UnsupportedFieldTypeError) Error() string {
	return "Unsupported field type: " + u.TypeName +
		" for field " + u.FieldName +
		" in model " + u.ModelName +
		" when compiling to " + u.TargetLanguage
}

//---------------------------------------------------------
