package cclGenerators

func DoGenerateCode(options *CodeGenerationOptions) (*CodeGenerationResult, error) {
	generatorFunc := CodeGenerators[options.TargetLanguage]
	if generatorFunc == nil {
		return nil, ErrLanguageNotSupported
	}

	return generatorFunc(options)
}
