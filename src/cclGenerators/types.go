package cclGenerators

import "github.com/ccl-lang/ccl/src/core/cclValues"

type GenerateCode func(options *CodeGenerationOptions) (*CodeGenerationResult, error)

type CodeGenerationOptions struct {
	CCLDefinition  *cclValues.SourceCodeDefinition
	OutputPath     string
	TargetLanguage string
	PackageName    string
}

type CodeGenerationResult struct {
}

type CodeGenerationBase struct {
	Options *CodeGenerationOptions
}
