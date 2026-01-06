package gdGen_test

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
	gdSource1 = "definitions1.ccl"
)

var (
	//go:embed contents/main_gd_content1.txt
	mainGdContent1 string
)

func TestGdGenerator1(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_gd_test_1")
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

	fmt.Printf("Generating GDScript code to: %s\n", tmpDir)

	// Parse CCL
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: gdSource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", gdSource1, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CCLDefinition:     parsedDefinitions,
		OutputPath:        filepath.Join(tmpDir, "models"),
		TargetLanguage:    "gd",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	fmt.Printf("Running GDScript code from: %s\n", tmpDir)

	output, err := RunGodotProject(&RunGodotOptions{
		TargetPath:    tmpDir,
		RunnerContent: mainGdContent1,
	})
	if err != nil {
		t.Fatalf("Failed to run Godot: %v\nOutput:\n%s", err, output)
		return
	}

	fmt.Printf("Output:\n%s\n", output)
}
