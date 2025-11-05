package codeGen_test

import (
	"fmt"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/cclParser"
)

const (
	goSource1 = "../../examples/definitions.ccl"
	goOutput1 = "../../examples/go_gen3"
)

func TestGoGenerator1(t *testing.T) {
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: goSource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", goSource1, parseErr)
		return
	}

	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CCLDefinition:  parsedDefinitions,
		OutputPath:     goOutput1,
		TargetLanguage: "go",
	})
	if err != nil {
		fmt.Printf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		fmt.Println("Unknown error: failed to generate code")
		return
	}

	fmt.Println("\nCode generation completed successfully")
}
