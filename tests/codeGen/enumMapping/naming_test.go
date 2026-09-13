package enumMapping_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

func TestShortEnumMappingNamesRejectAmbiguousDestinations(t *testing.T) {
	cclLoader.LoadGenerators()
	for _, language := range []string{"go", "gd", "cs", "js", "ts", "py", "rust"} {
		t.Run(language, func(t *testing.T) {
			definition, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
				SourceContent: `
#[EnumMapUseShortMethodName(true)]
model Source {
    [EnumMapSameMembers(First.Kind)]
    [EnumMapSameMembers(Second.Kind)]
    enum Kind { Unknown }
}
model First {
    [EnumMapUnknownMember(Kind.Unknown)]
    enum Kind { Unknown }
}
model Second {
    [EnumMapUnknownMember(Kind.Unknown)]
    enum Kind { Unknown }
}
`,
			})
			if err != nil {
				t.Fatal(err)
			}
			_, err = cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
				CodeContext:    definition.CodeContext,
				OutputPath:     t.TempDir(),
				TargetLanguage: language,
			})
			var diagnostic *cclErrors.InvalidAttributeUsageError
			if !errors.As(err, &diagnostic) || diagnostic.SourcePosition == nil || !strings.Contains(err.Error(), "collides") {
				t.Fatalf("expected a positioned naming collision diagnostic, got %v", err)
			}
		})
	}
}

func TestEnumMappingNamePreferenceRequiresBoolean(t *testing.T) {
	cclLoader.LoadGenerators()
	for _, argument := range []string{"", `"true"`, "true, false"} {
		t.Run(argument, func(t *testing.T) {
			definition, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
				SourceContent: `
#[EnumMapUseShortMethodName(` + argument + `)]
[EnumMapSameMembers(Target)]
enum Source { Gold }
[EnumMapUnknownMember(Target.Unknown)]
enum Target { Unknown, Gold }
`,
			})
			if err != nil {
				t.Fatal(err)
			}
			_, err = cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
				CodeContext:    definition.CodeContext,
				OutputPath:     t.TempDir(),
				TargetLanguage: "go",
			})
			var diagnostic *cclErrors.InvalidAttributeUsageError
			if !errors.As(err, &diagnostic) || !strings.Contains(err.Error(), "boolean parameter") {
				t.Fatalf("expected a boolean argument diagnostic, got %v", err)
			}
		})
	}
}
