package cclSanitizer

//---------------------------------------------------------

// fieldNameValidator is responsible for validating field names in
// model declarations to ensure they do not conflict with built-in
// type names or other model names within the same namespace.
type fieldNameValidator struct {
	modelNamesByNamespace map[string]map[string]string
}

//---------------------------------------------------------
