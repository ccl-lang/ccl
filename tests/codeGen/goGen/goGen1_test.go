package goGen_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/cclParser"
)

const (
	goSource1      = "definitions1.ccl"
	mainGoContent1 = `
package main

import (
	"ccl_test_generated/models"
	"fmt"
)

func main() {
	// Test creating a User
	user := models.User{
		Id:       "user123",
		Username: "test_user",
		Skin: &models.SkinInfo{
			Type: 0,
			Basic: &models.BasicSkin{
				R: 255,
				G: 0,
				B: 0,
			},
		},
	}

	userBinary, err := user.SerializeBinary()
	if err != nil {
		panic(fmt.Sprintf("Failed to serialize user: %v", err))
	}

	deserializedUser := &models.User{}
	err = deserializedUser.DeserializeBinary(userBinary)
	if err != nil {
		panic(fmt.Sprintf("Failed to deserialize user: %v", err))
	}

	if user.Username != "test_user" {
		panic("Username mismatch")
	} else if deserializedUser.Username != "test_user" {
		panic("Deserialized username mismatch")
	}

	// Test creating a Position
	pos := models.Position{
		X: 10,
		Y: 20,
		Z: 30,
	}

	if pos.X != 10 {
		panic("Position X mismatch")
	}

	positionBinary, err := pos.SerializeBinary()
	if err != nil {
		panic(fmt.Sprintf("Failed to serialize position: %v", err))
	}

	deserializedPosition := &models.Position{}
	err = deserializedPosition.DeserializeBinary(positionBinary)
	if err != nil {
		panic(fmt.Sprintf("Failed to deserialize position: %v", err))
	}

	if pos.X != 10 {
		panic("Position X mismatch")
	}

	fmt.Println("Runtime verification successful!")
}
`
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
		CCLDefinition:  parsedDefinitions,
		OutputPath:     filepath.Join(tmpDir, "models"),
		TargetLanguage: "go",
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
