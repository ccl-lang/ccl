package goGenerator

import (
	"strings"

	gen "github.com/ALiwoto/ccl/src/cclGenerators"
)

type GoGenerationContext struct {
	Options *gen.CodeGenerationOptions

	// MethodsCode is a string builder that contains the generated Go code for the methods.
	MethodsCode   *strings.Builder
	TypesCode     *strings.Builder
	HelpersCode   *strings.Builder
	VarsCode      *strings.Builder
	ConstantsCode *strings.Builder
}
