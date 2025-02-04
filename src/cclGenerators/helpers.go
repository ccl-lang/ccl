package cclGenerators

import "strings"

func DoGenerateCode(options *CodeGenerationOptions) (*CodeGenerationResult, error) {
	generatorFunc := CodeGenerators[strings.ToLower(options.TargetLanguage)]
	if generatorFunc == nil {
		return nil, ErrLanguageNotSupported
	}

	return generatorFunc(options)
}
