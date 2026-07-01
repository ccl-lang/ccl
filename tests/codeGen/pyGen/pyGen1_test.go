package pyGen_test

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
	pySource1 = filepath.Join("..", "definitions1.ccl")

	//go:embed contents/main_py_content1.txt
	mainPyContent1 string
)

func TestPythonGenerator1(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_py_test_1")
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

	fmt.Printf("Generating Python code to: %s\n", tmpDir)

	// Parse CCL
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: pySource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", pySource1, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:       parsedDefinitions.CodeContext,
		OutputPath:        filepath.Join(tmpDir, "models"),
		TargetLanguage:    "python",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	debugInfoPath := filepath.Join(tmpDir, "models", "auth_guest_answer.py.cclinfo")
	if _, err := os.Stat(debugInfoPath); err != nil {
		t.Fatalf("Expected Python debug info file to generate %s: %v", debugInfoPath, err)
	}

	output, err := RunPythonProject(&RunPythonOptions{
		TargetPath:    tmpDir,
		RunnerContent: mainPyContent1,
	})
	if err != nil {
		t.Fatalf("Failed to run generated code: %v\nOutput:\n%s", err, output)
		return
	}
	fmt.Printf("Output:\n%s\n", output)
}

func TestPythonGeneratorEnums(t *testing.T) {
	tmpDir, err := filepath.Abs("ccl_py_enum_test")
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
		TargetLanguage: "python",
	})
	if err != nil {
		t.Fatalf("Failed to generate Python enum code: %v", err)
	}
	if result == nil {
		t.Fatalf("Expected generation result")
	}

	enumContent := readGeneratedPythonFile(t, filepath.Join(tmpDir, "models", "user_type.py"))
	for _, snippet := range []string{
		"from enum import IntEnum",
		"class UserType(IntEnum):",
		"ADMIN = 10",
		"GUEST = 11",
	} {
		if !strings.Contains(enumContent, snippet) {
			t.Fatalf("Generated top-level enum is missing %q.\nGenerated:\n%s", snippet, enumContent)
		}
	}

	modelContent := readGeneratedPythonFile(t, filepath.Join(tmpDir, "models", "api_request_envelop.py"))
	for _, snippet := range []string{
		"from enum import IntEnum",
		"from .user_type import UserType",
		"class RequestType(IntEnum):",
		"GET_USER_DATA = 10",
		"self.other_field: ApiRequestEnvelop.RequestType = ApiRequestEnvelop.RequestType.LOGIN",
		"self.user_type: UserType = UserType.ADMIN",
	} {
		if !strings.Contains(modelContent, snippet) {
			t.Fatalf("Generated enum model is missing %q.\nGenerated:\n%s", snippet, modelContent)
		}
	}

	output, err := RunPythonProject(&RunPythonOptions{
		TargetPath: tmpDir,
		RunnerContent: `from models import ApiRequestEnvelop, UserType

req = ApiRequestEnvelop()
assert req.other_field == ApiRequestEnvelop.RequestType.LOGIN
assert req.user_type == UserType.ADMIN

req.request_id = ApiRequestEnvelop.RequestType.GET_USER_DATA
json_decoded = ApiRequestEnvelop.deserialize_json(req.serialize_json())
assert json_decoded.request_id == ApiRequestEnvelop.RequestType.GET_USER_DATA
assert json_decoded.user_type == UserType.ADMIN

binary_decoded = ApiRequestEnvelop.deserialize_binary(req.serialize_binary())
assert binary_decoded.request_id == ApiRequestEnvelop.RequestType.GET_USER_DATA
assert binary_decoded.user_type == UserType.ADMIN

print("enum ok")
`,
	})
	if err != nil {
		t.Fatalf("Failed to run generated enum code: %v\nOutput:\n%s", err, output)
	}
}

func readGeneratedPythonFile(t *testing.T, filePath string) string {
	t.Helper()

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read generated file %s: %v", filePath, err)
	}
	return string(data)
}
