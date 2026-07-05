package rsGenerator

import (
	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

type CCLModel = cclValues.ModelDefinition
type CCLField = cclValues.ModelFieldDefinition
type CCLEnum = cclValues.EnumDefinition

type RustGenerationContext struct {
	*gen.CodeGenerationBase

	CodeByPath map[string]*codeBuilder.CodeBuilder
}
