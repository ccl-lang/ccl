package gdGen_test

import (
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

func TestEnumMappingsSharingModelScope(t *testing.T) {
	const source = `
#[$gd:EnumMapUseShortMethodName(true)]
[EnumMapUnknownMember(Target.Unknown)]
enum Target { Unknown = 99, Gold = 2 }
model Mapper {
    [EnumMapSameMembers(Target)]
    enum First { Gold = 7 }

    [EnumMapSameMembers(Target)]
    enum Second { Gold = 8 }
}
`
	cclLoader.LoadGenerators()
	for _, override := range []bool{false, true} {
		content := source
		if override {
			content = strings.Replace(content, "enum Second", "[EnumMapUseShortMethodName(false)] enum Second", 1)
		}
		definition, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
			SourceContent: content,
		})
		if err != nil {
			t.Fatal(err)
		}
		targetPath := t.TempDir()
		_, err = cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
			CodeContext:    definition.CodeContext,
			OutputPath:     targetPath,
			TargetLanguage: "gd",
		})
		if !override {
			if err == nil || !strings.Contains(err.Error(), "collides in main.Mapper") {
				t.Fatalf("expected shared model scope collision, got %v", err)
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		output, err := RunGodotProject(&RunGodotOptions{
			TargetPath: targetPath,
			RunnerContent: `extends SceneTree
func _init():
    if Mapper.to_target(7) != Target.TargetEnum.GOLD or Mapper.second_to_target(8) != Target.TargetEnum.GOLD:
        quit(1)
        return
    if Mapper.to_target(8) != Target.TargetEnum.UNKNOWN or Mapper.second_to_target(7) != Target.TargetEnum.UNKNOWN:
        quit(1)
        return
    print("mapping scope ok")
    quit(0)
`,
		})
		if err != nil || !strings.Contains(output, "mapping scope ok") {
			t.Fatalf("generated mappings with local override failed: %v\n%s", err, output)
		}
	}
}
