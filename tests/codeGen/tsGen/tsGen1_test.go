package tsGen_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

var (
	tsSource1 = filepath.Join("..", "definitions1.ccl")

	//go:embed contents/main_ts_content1_1.txt
	mainTSContent1_1 string

	//go:embed contents/main_ts_content1_2.txt
	mainTSContent1_2 string
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
		SourceFilePath: tsSource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", tsSource1, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:       parsedDefinitions.CodeContext,
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
		SourceFilePath: tsSource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", tsSource1, parseErr)
		return
	}
	addSingleFileGenerationAttribute(parsedDefinitions, "generated.ts")

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:       parsedDefinitions.CodeContext,
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

	generatedContent := readGeneratedTSFile(t, filepath.Join(tmpDir, "models", "generated.ts"))
	for _, snippet := range []string{
		"export namespace UserInfo {",
		"export enum UserType {",
	} {
		if !strings.Contains(generatedContent, snippet) {
			t.Fatalf("Generated single-file TypeScript is missing %q.\nGenerated:\n%s", snippet, generatedContent)
		}
	}
	if strings.Contains(generatedContent, "export namespace  {") {
		t.Fatalf("Generated single-file TypeScript contains a blank namespace.\nGenerated:\n%s", generatedContent)
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

func addSingleFileGenerationAttribute(
	definition *cclValues.SourceCodeDefinition,
	fileName string,
) {
	enabledParam := &cclValues.ParameterInstance{}
	enabledParam.ChangeValue(true)

	fileNameParam := &cclValues.ParameterInstance{}
	fileNameParam.ChangeValue(fileName)

	definition.GlobalAttributes = append(definition.GlobalAttributes, &cclValues.AttributeUsageInfo{
		Name: cclAttr.AttrGenerateSingleFile,
		Parameters: []*cclValues.ParameterInstance{
			enabledParam,
			fileNameParam,
		},
	})
	definition.CodeContext.RegisterGlobalAttribute(definition.GlobalAttributes[len(definition.GlobalAttributes)-1])
}

func readGeneratedTSFile(t *testing.T, filePath string) string {
	t.Helper()

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read generated file %s: %v", filePath, err)
	}
	return string(data)
}
