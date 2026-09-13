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
	for _, mode := range []string{"long", "short"} {
		t.Run(mode, func(t *testing.T) {
			targetPath := filepath.Join(t.TempDir(), "ccl_enum_mapping")
			definition, err := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
				SourceFilePath: filepath.Join("..", "enum_mapping.ccl"),
			})
			if err != nil {
				t.Fatal(err)
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
			if mode == "short" {
				runner = strings.NewReplacer(
					"GameItem.game_item_type_to_player_item_player_item_type", "GameItem.to_player_item_type",
					"GameItem.game_item_type_to_saved_archive_type", "GameItem.to_archive_type",
					"PlayerItem.player_item_type_to_game_item_game_item_type", "PlayerItem.to_game_item_type",
				).Replace(runner)
			}

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
		})
	}
}
