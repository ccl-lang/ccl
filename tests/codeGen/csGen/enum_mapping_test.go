package csGen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

func TestEnumMappingsRuntime(t *testing.T) {
	for _, mode := range []string{"separate", "single", "short"} {
		t.Run(mode, func(t *testing.T) {
			targetPath := filepath.Join(t.TempDir(), "ccl_enum_mapping")
			definition, err := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
				SourceFilePath: filepath.Join("..", "enum_mapping.ccl"),
			})
			if err != nil {
				t.Fatal(err)
			}
			if mode == "single" {
				_, err = cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
					CodeContext:   definition.CodeContext,
					SourceContent: `#[GenerateSingleFile(true, "Mappings.cs")]`,
				})
				if err != nil {
					t.Fatal(err)
				}
			}
			if mode == "short" {
				_, err = cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
					CodeContext:   definition.CodeContext,
					SourceContent: `#[EnumMapUseShortMethodName(true)]`,
				})
				if err != nil {
					t.Fatal(err)
				}
			}
			cclLoader.LoadGenerators()
			_, err = cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
				CodeContext:    definition.CodeContext,
				OutputPath:     filepath.Join(targetPath, "models"),
				TargetLanguage: "cs",
			})
			if err != nil {
				t.Fatal(err)
			}
			runnerBytes, err := os.ReadFile("contents/enum_mapping_runner.txt")
			if err != nil {
				t.Fatal(err)
			}
			runner := string(runnerBytes)
			if mode == "short" {
				runner = strings.NewReplacer(
					"GameItemTypeMappings.ToPlayerItemPlayerItemType", "GameItemTypeMappings.ToPlayerItemType",
					"GameItemTypeMappings.ToSavedArchiveType", "GameItemTypeMappings.ToArchiveType",
					"PlayerItemTypeMappings.ToGameItemGameItemType", "PlayerItemTypeMappings.ToGameItemType",
				).Replace(runner)
			}
			projectBytes, err := os.ReadFile("contents/main_csproj_content1_1.txt")
			if err != nil {
				t.Fatal(err)
			}
			project := string(projectBytes)
			output, err := RunCSProject(&RunCSOptions{
				TargetPath:    targetPath,
				RunnerContent: runner,
				ProjectName:   "EnumMappings",
				CSProjContent: project,
			})
			if err != nil {
				t.Fatalf("generated enum mapping runtime failed: %v\n%s", err, output)
			}
		})
	}
}
