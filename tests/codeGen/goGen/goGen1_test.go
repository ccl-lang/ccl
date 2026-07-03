package goGen_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	modelIdConstantsContent := readGeneratedGoFile(t, filepath.Join(tmpDir, "models", "constants.go"))
	assertGoConstantFile(t, modelIdConstantsContent, []string{
		"ModelIdImportedGdThing",
		"ModelIdPositionInfo",
		"ModelIdUserInfo",
		"ModelIdStrictUserData",
	}, nil, []string{
		"MyCoolEnumElement1",
		"SkinTypeBasic",
		"UserInfoUserTypeUnknown",
	})

	commonConstantsContent := readGeneratedGoFile(t, filepath.Join(tmpDir, "models", "constants_common.go"))
	assertGoConstantFile(t, commonConstantsContent, nil, []string{"MyCoolEnumElement1"}, []string{
		"ModelIdImportedGdThing",
		"SkinTypeBasic",
		"UserInfoUserTypeUnknown",
	})

	skinsConstantsContent := readGeneratedGoFile(t, filepath.Join(tmpDir, "models", "constants_skins.go"))
	assertGoConstantFile(t, skinsConstantsContent, nil, []string{"SkinTypeBasic", "SkinTypeTexture"}, []string{
		"ModelIdUserInfo",
		"UserInfoUserTypeUnknown",
		"MyCoolEnumElement1",
	})

	usersConstantsContent := readGeneratedGoFile(t, filepath.Join(tmpDir, "models", "constants_users.go"))
	assertGoConstantFile(t, usersConstantsContent, nil, []string{
		"UserInfoUserTypeUnknown",
		"UserInfoUserTypePlayer",
	}, []string{
		"ModelIdUserInfo",
		"SkinTypeBasic",
		"MyCoolEnumElement1",
	})

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
		"constants_users.go",
		"constants_filegroup.go",
		"constants_global.go",
		"constants_direct.go",
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

func readGeneratedGoFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read generated Go file %s: %v", path, err)
	}

	return string(data)
}

func assertGoConstantFile(
	t *testing.T,
	constantsContent string,
	requiredModelIdSnippets []string,
	requiredEnumSnippets []string,
	forbiddenSnippets []string,
) {
	t.Helper()

	constBlocks := extractGoConstBlocks(t, constantsContent)
	if len(requiredModelIdSnippets) != 0 && len(requiredEnumSnippets) != 0 {
		if len(constBlocks) != 2 {
			t.Fatalf("Expected one model ID const block and one enum const block, got %d.\nGenerated:\n%s",
				len(constBlocks), constantsContent)
		}
	} else if len(constBlocks) != 1 {
		t.Fatalf("Expected one const block, got %d.\nGenerated:\n%s", len(constBlocks), constantsContent)
	}

	if len(requiredModelIdSnippets) != 0 {
		modelIdBlock := constBlocks[0]
		for _, snippet := range requiredModelIdSnippets {
			if !strings.Contains(modelIdBlock, snippet) {
				t.Fatalf("Generated model ID const block is missing %q.\nBlock:\n%s", snippet, modelIdBlock)
			}
		}
		for _, enumSnippet := range requiredEnumSnippets {
			if strings.Contains(modelIdBlock, enumSnippet) {
				t.Fatalf("Generated model ID const block contains enum member %q.\nBlock:\n%s", enumSnippet, modelIdBlock)
			}
		}
	}

	if len(requiredEnumSnippets) != 0 {
		enumBlock := constBlocks[len(constBlocks)-1]
		if strings.Contains(enumBlock, "ModelId") {
			t.Fatalf("Generated enum const block contains model IDs.\nBlock:\n%s", enumBlock)
		}
		for _, snippet := range requiredEnumSnippets {
			if !strings.Contains(enumBlock, snippet) {
				t.Fatalf("Generated enum const block is missing %q.\nBlock:\n%s", snippet, enumBlock)
			}
		}
	}

	for _, snippet := range forbiddenSnippets {
		if !strings.Contains(constantsContent, snippet) {
			continue
		}
		t.Fatalf("Generated constants file contains forbidden snippet %q.\nGenerated:\n%s", snippet, constantsContent)
	}
}

func extractGoConstBlocks(t *testing.T, constantsContent string) []string {
	t.Helper()

	constBlockParts := strings.Split(constantsContent, "const (")
	constBlocks := make([]string, 0, len(constBlockParts)-1)
	for _, currentPart := range constBlockParts[1:] {
		endIndex := strings.Index(currentPart, "\n)")
		if endIndex == -1 {
			t.Fatalf("Generated constants.go has an unclosed const block.\nGenerated:\n%s", constantsContent)
		}
		constBlocks = append(constBlocks, strings.TrimSpace(currentPart[:endIndex]))
	}

	return constBlocks
}
