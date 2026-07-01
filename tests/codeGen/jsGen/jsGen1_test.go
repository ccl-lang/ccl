package jsGen_test

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
	jsSource1 = filepath.Join("..", "definitions1.ccl")

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
		SourceFilePath: jsSource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", jsSource1, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:       parsedDefinitions.CodeContext,
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
		SourceFilePath: jsSource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", jsSource1, parseErr)
		return
	}
	addSingleFileGenerationAttribute(parsedDefinitions, "generated.js")

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:       parsedDefinitions.CodeContext,
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

func TestJSGeneratorEnums(t *testing.T) {
	tmpDir, err := filepath.Abs("ccl_js_enum_test")
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
		TargetLanguage: "js",
	})
	if err != nil {
		t.Fatalf("Failed to generate JavaScript enum code: %v", err)
	}
	if result == nil {
		t.Fatalf("Expected generation result")
	}

	enumContent := readGeneratedJSFile(t, filepath.Join(tmpDir, "models", "UserType.js"))
	for _, snippet := range []string{
		"export const UserType = Object.freeze({",
		"ADMIN: 10,",
		"GUEST: 11,",
	} {
		if !strings.Contains(enumContent, snippet) {
			t.Fatalf("Generated top-level enum is missing %q.\nGenerated:\n%s", snippet, enumContent)
		}
	}

	modelContent := readGeneratedJSFile(t, filepath.Join(tmpDir, "models", "ApiRequestEnvelop.js"))
	for _, snippet := range []string{
		"import { UserType } from './UserType.js';",
		"static RequestType = Object.freeze({",
		"GET_USER_DATA: 10,",
		"this.otherField = ApiRequestEnvelop.RequestType.LOGIN;",
		"this.userType = UserType.ADMIN;",
	} {
		if !strings.Contains(modelContent, snippet) {
			t.Fatalf("Generated enum model is missing %q.\nGenerated:\n%s", snippet, modelContent)
		}
	}
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

func readGeneratedJSFile(t *testing.T, filePath string) string {
	t.Helper()

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read generated file %s: %v", filePath, err)
	}
	return string(data)
}
