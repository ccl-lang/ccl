package goGen_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

var (
	goSource1 = filepath.Join("..", "definitions1.ccl")

	//go:embed contents/main_go_content1.txt
	mainGoContent1 string
)

func TestGoGenerator1(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_go_test_1")
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

	fmt.Printf("Generating Go code to: %s\n", tmpDir)

	// Parse CCL
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: goSource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", goSource1, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:       parsedDefinitions.CodeContext,
		OutputPath:        filepath.Join(tmpDir, "models"),
		TargetLanguage:    "go",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	output, err := RunGoProject(&RunGoOptions{
		TargetPath:    tmpDir,
		RunnerContent: mainGoContent1,
	})
	if err != nil {
		t.Fatalf("Failed to run generated code: %v\nOutput:\n%s", err, output)
		return
	}
	fmt.Printf("Output:\n%s\n", output)
}

func TestGoGeneratorOutputFileGroup(t *testing.T) {
	tmpDir := t.TempDir()
	usersPath := filepath.Join(tmpDir, "users.ccl")
	if err := os.WriteFile(usersPath, []byte(`
namespace main.users;
#namespace:[$go:OutputFileGroup("users")]
model NamespaceGrouped {
	Id: int;
}
`), 0644); err != nil {
		t.Fatalf("Failed to write users CCL source: %v", err)
	}

	fileGroupPath := filepath.Join(tmpDir, "file_group.ccl")
	if err := os.WriteFile(fileGroupPath, []byte(`
#file:[$go:OutputFileGroup("filegroup")]
model FileGrouped {
	Id: int;
}
`), 0644); err != nil {
		t.Fatalf("Failed to write file-group CCL source: %v", err)
	}

	sourcePath := filepath.Join(tmpDir, "api_types.ccl")
	if err := os.WriteFile(sourcePath, []byte(`
import "users.ccl";
import "file_group.ccl";

#[OutputFileGroup("global")]

model GlobalGrouped {
	Id: int;
}

[$go:OutputFileGroup("direct")]
model DirectGrouped {
	Id: int;
}
`), 0644); err != nil {
		t.Fatalf("Failed to write CCL source: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "models")
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: sourcePath,
	})
	if parseErr != nil {
		t.Fatalf("Failed to parse grouped source: %v", parseErr)
	}

	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:       parsedDefinitions.CodeContext,
		OutputPath:        outputPath,
		TargetLanguage:    "go",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Failed to generate grouped Go code: %v", err)
	}

	expectedFiles := []string{
		"constants.go",
		"vars.go",
		"helpers.go",
		"types_users.go",
		"methods_users.go",
		"types_filegroup.go",
		"methods_filegroup.go",
		"types_global.go",
		"methods_global.go",
		"types_direct.go",
		"methods_direct.go",
	}
	for _, fileName := range expectedFiles {
		if _, err := os.Stat(filepath.Join(outputPath, fileName)); err != nil {
			t.Fatalf("Expected generated file %s: %v", fileName, err)
		}
	}

	unexpectedFiles := []string{
		"types.go",
		"methods.go",
	}
	for _, fileName := range unexpectedFiles {
		if _, err := os.Stat(filepath.Join(outputPath, fileName)); !os.IsNotExist(err) {
			t.Fatalf("Did not expect generated file %s", fileName)
		}
	}

	if len(result.OutputFiles) != len(expectedFiles) {
		t.Fatalf("Expected %d output files, got %d: %v", len(expectedFiles), len(result.OutputFiles), result.OutputFiles)
	}
}
