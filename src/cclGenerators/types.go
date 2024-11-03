package cclGenerators

import "github.com/ALiwoto/ccl/src/core/cclValues"

type GenerateCode func(options *CodeGenerationOptions) (*CodeGenerationResult, error)

type CodeGenerationOptions struct {
	CCLDefinition  *cclValues.SourceCodeDefinition
	OutputPath     string
	TargetLanguage string
}

type CodeGenerationResult struct {
}
