package rsGen_test

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

var (
	rustSource1 = filepath.Join("..", "definitions1.ccl")

	//go:embed contents/main_rs_content1.txt
	mainRustContent1 string

	//go:embed contents/go_compat_content1.txt
	goCompatContent1 string

	//go:embed contents/rs_compat_content1.txt
	rustCompatContent1 string
)

func TestRustGenerator1(t *testing.T) {
	tmpDir, err := filepath.Abs("ccl_rs_test_1")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Fatalf("Failed to remove existing dir: %v", err)
	}
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: rustSource1,
	})
	if parseErr != nil {
		t.Fatalf("Failed to parse source file %s: %v", rustSource1, parseErr)
	}

	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:       parsedDefinitions.CodeContext,
		OutputPath:        tmpDir,
		TargetLanguage:    "rust",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Failed to generate Rust code: %v", err)
	}
	if result == nil {
		t.Fatalf("Failed to generate Rust code")
	}

	output, err := RunRustProject(&RunRustOptions{
		TargetPath:    tmpDir,
		RunnerContent: mainRustContent1,
	})
	if err != nil {
		t.Fatalf("Failed to run generated Rust code: %v\nOutput:\n%s", err, output)
	}
	fmt.Printf("Rust output:\n%s\n", output)
}

func TestRustGoBinaryCompatibility(t *testing.T) {
	tmpDir := t.TempDir()
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: rustSource1,
	})
	if parseErr != nil {
		t.Fatalf("Failed to parse source file %s: %v", rustSource1, parseErr)
	}

	cclLoader.LoadGenerators()
	goPath := filepath.Join(tmpDir, "go")
	rustPath := filepath.Join(tmpDir, "ccl_rs_test_1")
	if _, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:       parsedDefinitions.CodeContext,
		OutputPath:        filepath.Join(goPath, "models"),
		TargetLanguage:    "go",
		GenerateDebugInfo: true,
	}); err != nil {
		t.Fatalf("Failed to generate Go code: %v", err)
	}
	if _, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:       parsedDefinitions.CodeContext,
		OutputPath:        rustPath,
		TargetLanguage:    "rust",
		GenerateDebugInfo: true,
	}); err != nil {
		t.Fatalf("Failed to generate Rust code: %v", err)
	}

	goOutput, err := runGoCompatProject(goPath, goCompatContent1)
	if err != nil {
		t.Fatalf("Failed to run Go compatibility producer: %v\nOutput:\n%s", err, goOutput)
	}
	expectedHex := strings.TrimSpace(goOutput)
	if expectedHex == "" {
		t.Fatalf("Go compatibility producer returned empty output")
	}

	rustRunner := strings.ReplaceAll(rustCompatContent1, "{{GO_HEX}}", expectedHex)
	rustOutput, err := RunRustProject(&RunRustOptions{
		TargetPath:    rustPath,
		RunnerContent: rustRunner,
	})
	if err != nil {
		t.Fatalf("Failed to run Rust compatibility test: %v\nOutput:\n%s", err, rustOutput)
	}
}

func runGoCompatProject(targetPath string, runnerContent string) (string, error) {
	goModContent := `module ccl_test_generated

go 1.25
`
	if err := os.WriteFile(filepath.Join(targetPath, "go.mod"), []byte(goModContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(targetPath, "main.go"), []byte(runnerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write main.go: %v", err)
	}

	output, err := runCommand(targetPath, "go", "run", ".")
	return output, err
}

func runCommand(targetPath string, name string, args ...string) (string, error) {
	cmd := execCommand(name, args...)
	cmd.Dir = targetPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %v", err)
	}
	return string(output), nil
}

var execCommand = func(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}
