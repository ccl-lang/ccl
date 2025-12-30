package goGenerator

import (
	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

// CCLModel is a type alias for the model definition type from the CCL library.
type CCLModel = cclValues.ModelDefinition

// CCLField is a type alias for the field definition type from the CCL library.
type CCLField = cclValues.ModelFieldDefinition

type GoGenerationContext struct {
	*gen.CodeGenerationBase

	// MethodsCode is a string builder that contains the generated Go code for the methods.
	MethodsCode   *codeBuilder.CodeBuilder
	TypesCode     *codeBuilder.CodeBuilder
	HelpersCode   *codeBuilder.CodeBuilder
	VarsCode      *codeBuilder.CodeBuilder
	ConstantsCode *codeBuilder.CodeBuilder
}
