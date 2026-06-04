package cclGenerators

import (
	"strings"
)

func DoGenerateCode(options *CodeGenerationOptions) (*CodeGenerationResult, error) {
	if options == nil || options.CodeContext == nil {
		return nil, ErrMissingCodeContext
	}

	generatorFunc := CodeGenerators[strings.ToLower(options.TargetLanguage)]
	if generatorFunc == nil {
		return nil, ErrLanguageNotSupported
	}

	return generatorFunc(options)
}
