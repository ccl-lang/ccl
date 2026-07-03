package tsGen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

func TestTSGeneratorEnums(t *testing.T) {
	tmpDir, err := filepath.Abs("ccl_ts_enum_test")
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
		TargetLanguage: "ts",
	})
	if err != nil {
		t.Fatalf("Failed to generate TypeScript enum code: %v", err)
	}
	if result == nil {
		t.Fatalf("Expected generation result")
	}

	enumContent := readGeneratedTSFile(t, filepath.Join(tmpDir, "models", "UserType.ts"))
	for _, snippet := range []string{
		"export enum UserType {",
		"Admin = 10,",
		"Guest = 11,",
	} {
		if !strings.Contains(enumContent, snippet) {
			t.Fatalf("Generated top-level enum is missing %q.\nGenerated:\n%s", snippet, enumContent)
		}
	}

	modelContent := readGeneratedTSFile(t, filepath.Join(tmpDir, "models", "ApiRequestEnvelop.ts"))
	for _, snippet := range []string{
		"import { UserType } from './UserType';",
		"public requestId: ApiRequestEnvelop.RequestType;",
		"this.otherField = ApiRequestEnvelop.RequestType.Login;",
		"this.userType = UserType.Admin;",
		"export namespace ApiRequestEnvelop {",
		"export enum RequestType {",
		"GetUserData = 10,",
	} {
		if !strings.Contains(modelContent, snippet) {
			t.Fatalf("Generated enum model is missing %q.\nGenerated:\n%s", snippet, modelContent)
		}
	}
}

func TestTSGeneratorEnumMemberNamePrefix(t *testing.T) {
	tmpDir := t.TempDir()

	parsedDefinitions, parseErr := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
		SourceContent: `
[EnumMemberNamePrefix("Api")]
enum TopType {
	Hello,
}

[EnumMemberNamePrefix("Api")]
model ApiRequestEnvelop {
	enum RequestType {
		Login,
	}

	OtherField: RequestType = RequestType.Login;
	TopField: TopType = TopType.Hello;
}
`,
	})
	if parseErr != nil {
		t.Fatalf("Failed to parse enum prefix source: %v", parseErr)
	}

	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:    parsedDefinitions.CodeContext,
		OutputPath:     filepath.Join(tmpDir, "models"),
		TargetLanguage: "ts",
	})
	if err != nil {
		t.Fatalf("Failed to generate TypeScript enum prefix code: %v", err)
	}
	if result == nil {
		t.Fatalf("Expected generation result")
	}

	topEnumContent := readGeneratedTSFile(t, filepath.Join(tmpDir, "models", "TopType.ts"))
	if !strings.Contains(topEnumContent, "ApiHello = 0,") {
		t.Fatalf("Generated top-level enum is missing prefixed member.\nGenerated:\n%s", topEnumContent)
	}

	modelContent := readGeneratedTSFile(t, filepath.Join(tmpDir, "models", "ApiRequestEnvelop.ts"))
	for _, snippet := range []string{
		"ApiLogin = 0,",
		"this.otherField = ApiRequestEnvelop.RequestType.ApiLogin;",
		"this.topField = TopType.ApiHello;",
	} {
		if !strings.Contains(modelContent, snippet) {
			t.Fatalf("Generated enum prefix model is missing %q.\nGenerated:\n%s", snippet, modelContent)
		}
	}
}
