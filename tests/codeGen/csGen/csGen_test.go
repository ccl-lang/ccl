package csGen_test

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
	csSource1_1 = "definitions1_1.ccl"
)

var (
	//go:embed contents/main_cs_content1_1.txt
	mainCSContent1_1 string

	//go:embed contents/main_csproj_content1_1.txt
	mainCsProjContent1_1 string
)

// TestCSGenerator1_1 tests the C# code generator with multiple output files.
func TestCSGenerator1_1(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_cs_test_1_1")
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

	fmt.Printf("Generating C# code to: %s\n", tmpDir)

	// Parse CCL
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: csSource1_1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", csSource1_1, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CCLDefinition:     parsedDefinitions,
		OutputPath:        tmpDir,
		TargetLanguage:    "cs",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	output, err := RunCSProject(&RunCSOptions{
		TargetPath:    tmpDir,
		RunnerContent: mainCSContent1_1,
		CSProjContent: mainCsProjContent1_1,
		ProjectName:   "csTest",
	})
	if err != nil {
		t.Fatalf("Failed to run generated code: %v\nOutput:\n%s", err, output)
		return
	}
	fmt.Printf("Output:\n%s\n", output)
}
