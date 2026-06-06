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

	// CodeByPath maps relative output file paths to their builders.
	CodeByPath map[string]*codeBuilder.CodeBuilder

	JsonHelpersGenerated bool
}
