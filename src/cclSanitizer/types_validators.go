package cclSanitizer

import (
	"github.com/ccl-lang/ccl/src/core/cclUtils"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

//---------------------------------------------------------

// fieldNameValidator is responsible for validating field names in
// model declarations to ensure they do not conflict with built-in
// type names or other model names within the same namespace.
type fieldNameValidator struct {
	modelNamesByNamespace map[string]map[string]string
}

//---------------------------------------------------------

type fieldTypeUsageCheck struct {
	modelName      string
	fieldName      string
	typeUsage      *cclValues.CCLTypeUsage
	sourcePosition *cclUtils.SourceCodePosition
}

//---------------------------------------------------------
