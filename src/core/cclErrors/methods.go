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
