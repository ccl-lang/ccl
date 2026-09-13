package gdGen_test

import (
	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnumMappingsRuntime(t *testing.T) {
	targetPath := filepath.Join(t.TempDir(), "ccl_enum_mapping")
	definition, err := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: filepath.Join("..", "enum_mapping.ccl"),
	})
	if err != nil {
		t.Fatal(err)
	}
	cclLoader.LoadGenerators()
	_, err = cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CodeContext:    definition.CodeContext,
		OutputPath:     filepath.Join(targetPath, "models"),
		TargetLanguage: "gd",
	})
	if err != nil {
		t.Fatal(err)
	}
	runnerBytes, err := os.ReadFile("contents/enum_mapping_runner.txt")
	if err != nil {
		t.Fatal(err)
	}
	runner := string(runnerBytes)

	output, err := RunGodotProject(&RunGodotOptions{
		TargetPath:    targetPath,
		RunnerContent: runner,
	})
	if err != nil {
		t.Fatalf("generated enum mapping runtime failed: %v\n%s", err, output)
	}
	if !strings.Contains(output, "enum mapping ok") {
		t.Fatalf("Godot did not complete the runner:\n%s", output)
	}
}
