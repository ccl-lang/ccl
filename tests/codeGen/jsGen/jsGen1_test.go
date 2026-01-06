package jsGen_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/cclParser"
)

const (
	jsSource1_1 = "definitions1_1.ccl"
	jsSource1_2 = "definitions1_2.ccl"
)

var (
	//go:embed contents/main_js_content1_1.txt
	mainJSContent1_1 string

	//go:embed contents/main_js_content1_2.txt
	mainJSContent1_2 string
)

// TestJSGenerator1_1 tests the JavaScript code generator with multiple output files.
func TestJSGenerator1_1(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_js_test_1_1")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Clean up previous run
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Fatalf("Failed to remove existing dir: %v", err)
	}

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	fmt.Printf("Generating JS code to: %s\n", tmpDir)

	// Parse CCL
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: jsSource1_1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", jsSource1_1, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CCLDefinition:     parsedDefinitions,
		OutputPath:        filepath.Join(tmpDir, "models"),
		TargetLanguage:    "js",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	output, err := RunJSProject(&RunJSOptions{
		TargetPath:    tmpDir,
		RunnerContent: mainJSContent1_1,
	})
	if err != nil {
		t.Fatalf("Failed to run generated code: %v\nOutput:\n%s", err, output)
		return
	}
	fmt.Printf("Output:\n%s\n", output)
}

// TestJSGenerator1_2 tests the JavaScript code generator with single output file.
func TestJSGenerator1_2(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_js_test_1_2")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Clean up previous run
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Fatalf("Failed to remove existing dir: %v", err)
	}

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	fmt.Printf("Generating JS code to: %s\n", tmpDir)

	// Parse CCL
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: jsSource1_2,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", jsSource1_2, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CCLDefinition:     parsedDefinitions,
		OutputPath:        filepath.Join(tmpDir, "models"),
		TargetLanguage:    "js",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	output, err := RunJSProject(&RunJSOptions{
		TargetPath:    tmpDir,
		RunnerContent: mainJSContent1_2,
	})
	if err != nil {
		t.Fatalf("Failed to run generated code: %v\nOutput:\n%s", err, output)
		return
	}
	fmt.Printf("Output:\n%s\n", output)
}
