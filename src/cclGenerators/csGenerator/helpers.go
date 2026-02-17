package csGenerator

import (
	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
)

// GenerateCode generates C# code from the provided CCL source file.
func GenerateCode(options *gen.CodeGenerationOptions) (*gen.CodeGenerationResult, error) {
	context := &CSharpGenerationContext{
		CodeGenerationBase: &gen.CodeGenerationBase{
			Options: options,
		},
		ModelClasses: make(map[string]*codeBuilder.CodeBuilder),
	}

	return context.generateCode()
}
