package tsGen_test

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

var (
	//go:embed contents/main_ts_content1_1.txt
	mainTSContent1_1 string

	//go:embed contents/main_ts_content1_2.txt
	mainTSContent1_2 string
)

const (
	tsSource1_1 = "definitions1_1.ccl"
	tsSource1_2 = "definitions1_2.ccl"
)

// TestTSGenerator1_1 tests the TypeScript code generator with multiple output files.
func TestTSGenerator1_1(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_ts_test_1_1")
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

	fmt.Printf("Generating TS code to: %s\n", tmpDir)

	// Parse CCL
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: tsSource1_1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", tsSource1_1, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CCLDefinition:     parsedDefinitions,
		OutputPath:        filepath.Join(tmpDir, "models"),
		TargetLanguage:    "ts",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	output, err := RunTSProject(&RunTSOptions{
		TargetPath:    tmpDir,
		RunnerContent: mainTSContent1_1,
	})
	if err != nil {
		t.Fatalf("Failed to run generated code: %v\nOutput:\n%s", err, output)
		return
	}
	fmt.Printf("Output:\n%s\n", output)
}

// TestTSGenerator1_2 tests the TypeScript code generator with single output file.
func TestTSGenerator1_2(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_ts_test_1_2")
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

	fmt.Printf("Generating TS code to: %s\n", tmpDir)

	// Parse CCL
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: tsSource1_2,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", tsSource1_2, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CCLDefinition:     parsedDefinitions,
		OutputPath:        filepath.Join(tmpDir, "models"),
		TargetLanguage:    "ts",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	output, err := RunTSProject(&RunTSOptions{
		TargetPath:    tmpDir,
		RunnerContent: mainTSContent1_2,
	})
	if err != nil {
		t.Fatalf("Failed to run generated code: %v\nOutput:\n%s", err, output)
		return
	}
	fmt.Printf("Output:\n%s\n", output)
}
