package jsGenerator

import (
	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
)

func GenerateCode(options *gen.CodeGenerationOptions) (*gen.CodeGenerationResult, error) {
	ctx := &JavaScriptGenerationContext{
		CodeGenerationBase: &gen.CodeGenerationBase{
			Options: options,
		},
		ModelClasses: make(map[string]*codeBuilder.CodeBuilder),
	}

	return ctx.generateCode()
}
