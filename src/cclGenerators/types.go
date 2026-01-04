package cclGenerators

import (
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

type GenerateCode func(options *CodeGenerationOptions) (*CodeGenerationResult, error)

type CodeGenerationOptions struct {
	CCLDefinition     *cclValues.SourceCodeDefinition
	OutputPath        string
	TargetLanguage    string
	PackageName       string
	GenerateDebugInfo bool
}

// CodeGenerationResult holds the result of a code generation process.
type CodeGenerationResult struct {
	OutputFiles []string

	// SourceLanguage is the language from which the code was generated.
	SourceLanguage gValues.NormalizedLangName

	// TargetLanguage is the language to which the code was generated.
	TargetLanguage gValues.NormalizedLangName
}

type CodeGenerationBase struct {
	Options *CodeGenerationOptions
}
