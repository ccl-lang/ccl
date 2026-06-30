package cclValues

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils"
)

// AttributeUsageInfo is a struct that represents an attribute definition
// in the source code with its parameters.
type AttributeUsageInfo struct {
	// Name is the name of the attribute.
	Name cclAttr.CCLAttributeName

	// SourceFileId is the source file where this attribute is defined.
	SourceFileId SourceFileId

	// Namespace is the namespace this attribute belongs to when it is
	// namespace-scoped.
	Namespace string

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

// AttributeResolutionSubject describes the declaration or file for scoped
// attribute resolution.
type AttributeResolutionSubject struct {
	SourceFileId SourceFileId
	Namespace    string
	Enum         *EnumDefinition
	Model        *ModelDefinition
	Field        *ModelFieldDefinition
}

// AttributeResolutionOptions controls which fallback levels are used.
type AttributeResolutionOptions struct {
	AllowModelFallback  bool
	AllowScopedFallback bool
	AllowGlobalFallback bool
}
