package cclGenerators

import (
	"strings"
	"unicode"
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

func isValidOutputFileGroup(group string) bool {
	if group == "" {
		return true
	}

	for _, currentRune := range group {
		if unicode.IsLetter(currentRune) ||
			unicode.IsDigit(currentRune) ||
			currentRune == '_' {
			continue
		}

		return false
	}

	return true
}
