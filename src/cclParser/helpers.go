package cclParser

import (
	"os"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func ParseCCLSourceFile_OLD(options *CCLParseOptions) (*cclValues.SourceCodeDefinition, error) {
	content, err := os.ReadFile(options.Source)
	if err != nil {
		return nil, err
	}

	srcDefinition := &cclValues.SourceCodeDefinition{}
	modelMatches := modelRegex.FindAllStringSubmatch(string(content), -1)
	parsedModels := []*cclValues.ModelDefinition{}
	definedModels := map[string]bool{}

	for _, modelMatch := range modelMatches {
		modelName := modelMatch[1]
		fieldMatches := fieldRegex.FindAllStringSubmatch(modelMatch[2], -1)
		definedFields := map[string]bool{}
		currentModel := &cclValues.ModelDefinition{
			ModelId: srcDefinition.GetNextModelId(),
			Name:    modelName,
			Fields:  []*cclValues.FieldDefinition{},
		}

		for _, fieldMatch := range fieldMatches {
			fieldName := fieldMatch[1]
			fieldType := fieldMatch[2]
			extraOperators := fieldMatch[3]
			if _, exists := definedFields[fieldName]; exists {
				return nil, &cclErrors.DuplicateFieldError{
					ModelName: modelName,
					FieldName: fieldName,
				}
			}

			currentModel.Fields = append(currentModel.Fields, &cclValues.FieldDefinition{
				OwnedBy: currentModel,
				Name:    fieldName,
				Type:    cclValues.NewTypeInfoWithOperators(fieldType, extraOperators),
			})
			definedFields[fieldName] = true
		}

		if _, exists := definedModels[modelName]; exists {
			return nil, &cclErrors.DuplicateModelError{
				ModelName: modelName,
			}
		}
		parsedModels = append(parsedModels, currentModel)
		definedModels[modelName] = true
	}

	srcDefinition.Models = parsedModels
	return srcDefinition, nil
}

func ParseCCLSourceFile(options *CCLParseOptions) (*cclValues.SourceCodeDefinition, error) {
	allTokens, err := cclLexer.Lex(options.Source)
	if err != nil {
		return nil, err
	}

	theParser := &CCLParser{
		Options: options,
		tokens:  allTokens,
	}
	return theParser.ParseAsCCL()
}
