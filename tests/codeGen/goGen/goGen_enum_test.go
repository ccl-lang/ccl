package goGen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

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

func TestGoGeneratorEnumMemberNamePrefix(t *testing.T) {
	tmpDir := t.TempDir()

	parsedDefinitions, parseErr := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
		SourceContent: `
model DefaultContainer {
	enum SType {
		Hello,
	}

	Value: SType = SType.Hello;
}

[EnumMemberNamePrefix("Abc")]
model ParentPrefixContainer {
	enum SType {
		Hello,
	}

	Value: SType = SType.Hello;
}

model ExplicitPrefixContainer {
	[EnumMemberNamePrefix("")]
	enum EmptyType {
		Hello,
	}

	[EnumMemberNamePrefix(null)]
	enum NullType {
		Hello,
	}

	[EnumMemberNamePrefix("Abc")]
	enum CustomType {
		Hello,
	}

	EmptyValue: EmptyType = EmptyType.Hello;
	NullValue: NullType = NullType.Hello;
	CustomValue: CustomType = CustomType.Hello;
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
		TargetLanguage: "go",
	})
	if err != nil {
		t.Fatalf("Failed to generate Go enum prefix code: %v", err)
	}
	if result == nil {
		t.Fatalf("Expected generation result")
	}

	constantsContent := readGeneratedGoFile(t, filepath.Join(tmpDir, "models", "constants.go"))
	for _, expectedSnippet := range []string{
		"DefaultContainerSTypeHello",
		"AbcSTypeHello",
		"EmptyTypeHello",
		"NullTypeHello",
		"AbcCustomTypeHello",
	} {
		if !strings.Contains(constantsContent, expectedSnippet) {
			t.Fatalf("Generated constants are missing %q.\nGenerated:\n%s", expectedSnippet, constantsContent)
		}
	}
	for _, forbiddenSnippet := range []string{
		"ExplicitPrefixContainerEmptyTypeHello",
		"ExplicitPrefixContainerNullTypeHello",
		"ExplicitPrefixContainerCustomTypeHello",
		"ParentPrefixContainerSTypeHello",
	} {
		if strings.Contains(constantsContent, forbiddenSnippet) {
			t.Fatalf("Generated constants contain forbidden %q.\nGenerated:\n%s", forbiddenSnippet, constantsContent)
		}
	}

	output, err := RunGoProject(&RunGoOptions{
		TargetPath: tmpDir,
		RunnerContent: `package main

import (
	"fmt"

	"ccl_test_generated/models"
)

func main() {
	defaultContainer := models.NewDefaultContainer()
	if defaultContainer.Value != models.DefaultContainerSTypeHello {
		panic("wrong default nested enum prefix")
	}

	parentPrefixContainer := models.NewParentPrefixContainer()
	if parentPrefixContainer.Value != models.AbcSTypeHello {
		panic("wrong parent enum prefix")
	}

	explicitPrefixContainer := models.NewExplicitPrefixContainer()
	if explicitPrefixContainer.EmptyValue != models.EmptyTypeHello {
		panic("wrong empty enum prefix")
	}
	if explicitPrefixContainer.NullValue != models.NullTypeHello {
		panic("wrong null enum prefix")
	}
	if explicitPrefixContainer.CustomValue != models.AbcCustomTypeHello {
		panic("wrong custom enum prefix")
	}

	fmt.Println("enum prefix ok")
}
`,
	})
	if err != nil {
		t.Fatalf("Failed to run generated enum prefix code: %v\nOutput:\n%s", err, output)
	}
}
