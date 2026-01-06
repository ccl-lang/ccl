package cclValues

import (
	"github.com/ccl-lang/ccl/src/core/cclUtils"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

// AttributeUsageInfo is a struct that represents an attribute definition
// in the source code with its parameters.
type AttributeUsageInfo struct {
	// Name is the name of the attribute.
	Name string

	// Languages specifies which languages this attribute should
	// get applied to.
	// it can be an empty string to indicate all languages.
	Languages []gValues.LanguageType

	// Parameters is the list of parameters for the attribute.
	Parameters []*ParameterInstance

	// SourcePosition is the position of the attribute in the source code.
	SourcePosition *cclUtils.SourceCodePosition
}

// AttributesCollection is a collection of attributes that holds some
// attributes usage info under itself and has some useful methods for
// manipulating them.
type AttributesCollection struct {
	Attrs []*AttributeUsageInfo
}
