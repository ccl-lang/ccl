package goGenerator

import (
	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

// CCLModel is a type alias for the model definition type from the CCL library.
type CCLModel = cclValues.ModelDefinition

// CCLField is a type alias for the field definition type from the CCL library.
type CCLField = cclValues.ModelFieldDefinition

// CCLEnum is a type alias for the enum definition type from the CCL library.
type CCLEnum = cclValues.EnumDefinition

type GoGenerationContext struct {
	*gen.CodeGenerationBase

	// CodeByPath maps relative output file paths to their builders.
	CodeByPath map[string]*codeBuilder.CodeBuilder
}
