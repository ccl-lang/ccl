package goGen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
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
		TargetLanguage: "go",
	})
	if err != nil {
		t.Fatal(err)
	}
	runnerBytes, err := os.ReadFile("contents/enum_mapping_runner.txt")
	if err != nil {
		t.Fatal(err)
	}
	runner := string(runnerBytes)

	output, err := RunGoProject(&RunGoOptions{
		TargetPath:    targetPath,
		RunnerContent: runner,
	})
	if err != nil {
		t.Fatalf("generated enum mapping runtime failed: %v\n%s", err, output)
	}
}
