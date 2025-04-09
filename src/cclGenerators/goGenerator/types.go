package goGenerator

import (
	"strings"

	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

// CCLModel is a type alias for the model definition type from the CCL library.
type CCLModel = cclValues.ModelDefinition

// CCLField is a type alias for the field definition type from the CCL library.
type CCLField = cclValues.FieldDefinition

type GoGenerationContext struct {
	Options *gen.CodeGenerationOptions

	// MethodsCode is a string builder that contains the generated Go code for the methods.
	MethodsCode   *strings.Builder
	TypesCode     *strings.Builder
	HelpersCode   *strings.Builder
	VarsCode      *strings.Builder
	ConstantsCode *strings.Builder
}
