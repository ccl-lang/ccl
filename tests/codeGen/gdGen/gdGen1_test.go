package gdGen_test

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
	gdSource1 = filepath.Join("..", "definitions1.ccl")

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
		CodeContext:       parsedDefinitions.CodeContext,
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

	for _, outputFile := range result.OutputFiles {
		data, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read generated file %s: %v", outputFile, err)
		}
		content := string(data)
		expectedContent := strings.TrimRight(content, " \t\r\n") + "\n"
		if content != expectedContent {
			t.Fatalf("Generated file %s has trailing EOF whitespace", outputFile)
		}
	}

	importedOutputPath := filepath.Join(tmpDir, "models", "imported_gd_thing.gd")
	if _, err := os.Stat(importedOutputPath); err != nil {
		t.Fatalf("Expected imported CCL model to generate %s: %v", importedOutputPath, err)
	}

	skinTypeContent := readGeneratedGdModel(t, result.OutputFiles, "class_name SkinType")
	if !strings.Contains(skinTypeContent, "enum SkinTypeEnum {") {
		t.Fatalf("Generated SkinType enum did not use named top-level enum.\nGenerated:\n%s", skinTypeContent)
	}

	skinInfoContent := readGeneratedGdModel(t, result.OutputFiles, "class_name SkinInfo")
	skinInfoSnippets := []string{
		"var type: SkinType.SkinTypeEnum",
		"model_result.type = buffer.get_u8() as SkinType.SkinTypeEnum",
		"model_result.type = int(type_value) as SkinType.SkinTypeEnum",
	}
	for _, snippet := range skinInfoSnippets {
		if !strings.Contains(skinInfoContent, snippet) {
			t.Fatalf("Generated SkinInfo model is missing enum snippet %q.\nGenerated:\n%s", snippet, skinInfoContent)
		}
	}

	gameShopItemContent := readGeneratedGdModel(t, result.OutputFiles, "class_name GameShopItem")
	if !strings.Contains(gameShopItemContent, "var featured_type: GameItem.GameItemType") {
		t.Fatalf("Generated GameShopItem model did not qualify external nested enum type.\nGenerated:\n%s", gameShopItemContent)
	}

	gameItemContent := readGeneratedGdModel(t, result.OutputFiles, "class_name GameItem")
	if !strings.Contains(gameItemContent, "var item_type: GameItemType") {
		t.Fatalf("Generated GameItem model did not keep same-owner nested enum type local.\nGenerated:\n%s", gameItemContent)
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
	if !strings.Contains(output, "Runtime verification successful!") {
		t.Fatalf("Godot runtime verification did not report success.\nOutput:\n%s", output)
		return
	}

	fmt.Printf("Output:\n%s\n", output)
}

func TestGdGeneratorAddAnnotation(t *testing.T) {
	tmpDir, err := filepath.Abs("ccl_gd_annotation_test")
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
#[AddAnnotation("@global_stuff")]

[AddAnnotation("@stuff")]
[AddAnnotation("@stuff2")]
[AddAnnotation("@stuff(23)")]
model AnnotatedThing {
	Id: string;
}

model GloballyAnnotatedThing {
	Id: string;
}
`,
	})
	if parseErr != nil {
		t.Fatalf("Failed to parse annotation source: %v", parseErr)
	}

	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:    parsedDefinitions.CodeContext,
		OutputPath:     filepath.Join(tmpDir, "models"),
		TargetLanguage: "gd",
	})
	if err != nil {
		t.Fatalf("Failed to generate GDScript code: %v", err)
	}
	if result == nil || len(result.OutputFiles) != 2 {
		t.Fatalf("Expected two generated output files, got %#v", result)
	}

	content := readGeneratedGdModel(t, result.OutputFiles, "class_name AnnotatedThing")
	expectedHeader := strings.Join([]string{
		"# THIS FILE IS AUTOGENERATED BY A CCL TOOL. DO NOT EDIT.",
		"",
		"@stuff",
		"@stuff2",
		"@stuff(23)",
		"class_name AnnotatedThing",
	}, "\n")
	if !strings.Contains(content, expectedHeader) {
		t.Fatalf("Generated file did not include raw annotations before class_name.\nExpected header:\n%s\nGenerated:\n%s", expectedHeader, content)
	}

	globalContent := readGeneratedGdModel(t, result.OutputFiles, "class_name GloballyAnnotatedThing")
	expectedGlobalHeader := strings.Join([]string{
		"# THIS FILE IS AUTOGENERATED BY A CCL TOOL. DO NOT EDIT.",
		"",
		"@global_stuff",
		"class_name GloballyAnnotatedThing",
	}, "\n")
	if !strings.Contains(globalContent, expectedGlobalHeader) {
		t.Fatalf("Generated file did not include global raw annotation before class_name.\nExpected header:\n%s\nGenerated:\n%s", expectedGlobalHeader, globalContent)
	}
}

func TestGdGeneratorUseWGodot(t *testing.T) {
	tmpDir, err := filepath.Abs("ccl_gd_wgodot_test")
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
#[$gd:UseWGodot(false)]

model ApiCardInfo {
	Name: string;
}

[$gd:UseWGodot(true)]
model ApiCardEnvelope {
	Card: ApiCardInfo;
	Cards: ApiCardInfo[];
	Payload: bytes;
}
`,
	})
	if parseErr != nil {
		t.Fatalf("Failed to parse WGodot source: %v", parseErr)
	}

	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:    parsedDefinitions.CodeContext,
		OutputPath:     filepath.Join(tmpDir, "models"),
		TargetLanguage: "gd",
	})
	if err != nil {
		t.Fatalf("Failed to generate GDScript code: %v", err)
	}

	cardInfoContent := readGeneratedGdModel(t, result.OutputFiles, "class_name ApiCardInfo")
	if !strings.Contains(cardInfoContent, "var name_len := buffer.get_u32()") {
		t.Fatalf("Expected global UseWGodot(false) to use inferred length declaration.\nGenerated:\n%s", cardInfoContent)
	}
	if !strings.Contains(cardInfoContent, "model_result.name = buffer.get_data(name_len)[1].get_string_from_utf8()") {
		t.Fatalf("Expected global UseWGodot(false) to keep inline string get_data read.\nGenerated:\n%s", cardInfoContent)
	}

	envelopeContent := readGeneratedGdModel(t, result.OutputFiles, "class_name ApiCardEnvelope")
	expectedSnippets := []string{
		"var card_len := buffer.get_u32()",
		"buffer.get_data_bytes",
	}
	for _, snippet := range expectedSnippets {
		if !strings.Contains(envelopeContent, snippet) {
			t.Fatalf("Generated WGodot model is missing snippet %q.\nGenerated:\n%s", snippet, envelopeContent)
		}
	}
}

func TestGdGeneratorEnums(t *testing.T) {
	tmpDir, err := filepath.Abs("ccl_gd_enum_test")
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
		TargetLanguage: "gd",
	})
	if err != nil {
		t.Fatalf("Failed to generate GDScript enum code: %v", err)
	}
	if result == nil || len(result.OutputFiles) != 2 {
		t.Fatalf("Expected two generated files, got %#v", result)
	}

	enumContent := readGeneratedGdModel(t, result.OutputFiles, "class_name UserType")
	enumSnippets := []string{
		"enum UserTypeEnum {",
		"UNKNOWN = 0,",
		"ADMIN = 10,",
		"GUEST = 11,",
	}
	for _, snippet := range enumSnippets {
		if !strings.Contains(enumContent, snippet) {
			t.Fatalf("Generated top-level enum is missing %q.\nGenerated:\n%s", snippet, enumContent)
		}
	}

	modelContent := readGeneratedGdModel(t, result.OutputFiles, "class_name ApiRequestEnvelop")
	modelSnippets := []string{
		"enum RequestType {",
		"GET_USER_DATA = 10,",
		"GET_USER_AVATAR = 11,",
		"var request_id: RequestType",
		"var other_field: RequestType = RequestType.LOGIN",
		"var user_type: UserType.UserTypeEnum = UserType.UserTypeEnum.ADMIN",
	}
	for _, snippet := range modelSnippets {
		if !strings.Contains(modelContent, snippet) {
			t.Fatalf("Generated enum model is missing %q.\nGenerated:\n%s", snippet, modelContent)
		}
	}
}

func readGeneratedGdModel(t *testing.T, outputFiles []string, classNameLine string) string {
	t.Helper()

	for _, outputFile := range outputFiles {
		data, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read generated file %s: %v", outputFile, err)
		}

		content := string(data)
		if strings.Contains(content, classNameLine) {
			return content
		}
	}

	t.Fatalf("Generated files did not contain %s", classNameLine)
	return ""
}
