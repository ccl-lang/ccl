package cclParser

import (
	"os"
	"regexp"

	"github.com/ALiwoto/ccl/src/core/cclErrors"
	"github.com/ALiwoto/ccl/src/core/cclValues"
)

func ParseCCLSourceFile(options *CCLParseOptions) (*cclValues.SourceCodeDefinition, error) {
	content, err := os.ReadFile(options.Source)
	if err != nil {
		return nil, err
	}

	modelRegex := regexp.MustCompile(`(?ms)model\s+(\w+)\s*\{(.*?)\}`)
	fieldRegex := regexp.MustCompile(`(\w+)\s*:\s*([\w]+)\s*(\[\s*\])?\s*;`)

	modelMatches := modelRegex.FindAllStringSubmatch(string(content), -1)
	parsedModels := cclValues.ModelsMap{}

	for _, modelMatch := range modelMatches {
		modelName := modelMatch[1]
		fieldMatches := fieldRegex.FindAllStringSubmatch(modelMatch[2], -1)
		parsedFields := cclValues.FieldsMap{}

		for _, fieldMatch := range fieldMatches {
			fieldName := fieldMatch[1]
			fieldType := fieldMatch[2]
			extraOperators := fieldMatch[3]
			if _, exists := parsedFields[fieldName]; exists {
				return nil, &cclErrors.DuplicateFieldError{
					ModelName: modelName,
					FieldName: fieldName,
				}
			}

			parsedFields[fieldName] = &cclValues.FieldDefinition{
				Name:           fieldName,
				Type:           fieldType,
				ExtraOperators: extraOperators,
			}
		}

		if _, exists := parsedModels[modelName]; exists {
			return nil, &cclErrors.DuplicateModelError{
				ModelName: modelName,
			}
		}
		parsedModels[modelName] = &cclValues.ModelDefinition{
			Name:   modelName,
			Fields: parsedFields,
		}
	}

	return &cclValues.SourceCodeDefinition{
		Models: parsedModels,
	}, nil
}
