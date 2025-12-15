package codeGen_test

import (
	"fmt"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/cclParser"
)

const (
	gdSource1 = "../../examples/definitions.ccl"
	gdOutput1 = "../../examples/gd_gen1"
)

func TestGdGenerator1(t *testing.T) {
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: gdSource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", gdSource1, parseErr)
		return
	}

	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CCLDefinition:  parsedDefinitions,
		OutputPath:     gdOutput1,
		TargetLanguage: "gd",
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
