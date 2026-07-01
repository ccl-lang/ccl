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

	constantsContent := readGeneratedGoFile(t, filepath.Join(tmpDir, "models", "constants.go"))
	assertGoConstantsAreGrouped(t, constantsContent)

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

func readGeneratedGoFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read generated Go file %s: %v", path, err)
	}

	return string(data)
}

func assertGoConstantsAreGrouped(t *testing.T, constantsContent string) {
	t.Helper()

	constBlocks := extractGoConstBlocks(t, constantsContent)
	if len(constBlocks) != 4 {
		t.Fatalf("Expected one model ID const block and three enum const blocks, got %d.\nGenerated:\n%s",
			len(constBlocks), constantsContent)
	}

	modelIdBlock := constBlocks[0]
	requiredModelIds := []string{
		"ModelIdImportedGdThing",
		"ModelIdPositionInfo",
		"ModelIdStrictUserData",
	}
	for _, snippet := range requiredModelIds {
		if !strings.Contains(modelIdBlock, snippet) {
			t.Fatalf("Generated model ID const block is missing %q.\nBlock:\n%s", snippet, modelIdBlock)
		}
	}
	for _, enumSnippet := range []string{"MyCoolEnumElement1", "SkinTypeBasic", "UserInfoUserTypeUnknown"} {
		if strings.Contains(modelIdBlock, enumSnippet) {
			t.Fatalf("Generated model ID const block contains enum member %q.\nBlock:\n%s", enumSnippet, modelIdBlock)
		}
	}

	assertSingleGoEnumConstBlock(t, constBlocks[1], "MyCoolEnumElement1", []string{"SkinTypeBasic", "UserInfoUserTypeUnknown"})
	assertSingleGoEnumConstBlock(t, constBlocks[2], "SkinTypeBasic", []string{"MyCoolEnumElement1", "UserInfoUserTypeUnknown"})
	assertSingleGoEnumConstBlock(t, constBlocks[3], "UserInfoUserTypeUnknown", []string{"MyCoolEnumElement1", "SkinTypeBasic"})
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

func assertSingleGoEnumConstBlock(
	t *testing.T,
	constBlock string,
	requiredSnippet string,
	forbiddenSnippets []string,
) {
	t.Helper()

	if !strings.Contains(constBlock, requiredSnippet) {
		t.Fatalf("Generated enum const block is missing %q.\nBlock:\n%s", requiredSnippet, constBlock)
	}
	if strings.Contains(constBlock, "ModelId") {
		t.Fatalf("Generated enum const block contains model IDs.\nBlock:\n%s", constBlock)
	}
	for _, forbiddenSnippet := range forbiddenSnippets {
		if strings.Contains(constBlock, forbiddenSnippet) {
			t.Fatalf("Generated enum const block mixes enum member %q.\nBlock:\n%s", forbiddenSnippet, constBlock)
		}
	}
}

func TestGoGeneratorEnums(t *testing.T) {
	tmpDir, err := filepath.Abs("ccl_go_enum_test")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		t.Fatalf("Failed to remove existing dir: %v", err)
	}

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	parsedDefinitions, parseErr := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
		SourceContent: `
#[SerializationType("binary")]
#[SerializationType("json")]

enum UserType {
	Unknown,
	Admin = 10,
	Guest,
}

model ApiRequestEnvelop {
	enum RequestType: int32 {
		Unknown,
		Login,
		GetUserData = 10,
		GetUserAvatar,
	}

	RequestId: RequestType;
	OtherField: RequestType = RequestType.Login;
	UserType: UserType = UserType.Admin;
}
`,
	})
	if parseErr != nil {
		t.Fatalf("Failed to parse enum source: %v", parseErr)
	}

	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:    parsedDefinitions.CodeContext,
		OutputPath:     filepath.Join(tmpDir, "models"),
		TargetLanguage: "go",
	})
	if err != nil {
		t.Fatalf("Failed to generate Go enum code: %v", err)
	}
	if result == nil {
		t.Fatalf("Expected generation result")
	}

	output, err := RunGoProject(&RunGoOptions{
		TargetPath: tmpDir,
		RunnerContent: `package main

import (
	"fmt"

	"ccl_test_generated/models"
)

func main() {
	req := models.NewApiRequestEnvelop()
	if req.OtherField != models.ApiRequestEnvelopRequestTypeLogin {
		panic("wrong nested enum default")
	}
	if req.UserType != models.UserTypeAdmin {
		panic("wrong top enum default")
	}

	req.RequestId = models.ApiRequestEnvelopRequestTypeGetUserData
	binaryData, err := req.SerializeBinary()
	if err != nil {
		panic(err)
	}
	var binaryDecoded models.ApiRequestEnvelop
	if err := binaryDecoded.DeserializeBinary(binaryData); err != nil {
		panic(err)
	}
	if binaryDecoded.RequestId != models.ApiRequestEnvelopRequestTypeGetUserData {
		panic("wrong binary enum roundtrip")
	}

	jsonText, err := req.SerializeJSON()
	if err != nil {
		panic(err)
	}
	var jsonDecoded models.ApiRequestEnvelop
	if err := jsonDecoded.DeserializeJSON(jsonText); err != nil {
		panic(err)
	}
	if jsonDecoded.UserType != models.UserTypeAdmin {
		panic("wrong json enum roundtrip")
	}

	fmt.Println("enum ok")
}
`,
	})
	if err != nil {
		t.Fatalf("Failed to run generated enum code: %v\nOutput:\n%s", err, output)
	}
}
