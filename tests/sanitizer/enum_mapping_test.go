package sanitizer_test

import (
	"errors"
	"strings"
	"testing"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

func TestEnumMappingDiagnostics(t *testing.T) {
	for _, scenario := range []struct{ name, source, diagnostic string }{
		{"unknown target", `[EnumMapSameMembers(Missing)] enum Source { Value }`, "Missing"},
		{"model target", `[EnumMapSameMembers(Target)] enum Source { Value } model Target {}`, "must resolve to an enum"},
		{"literal target", `[EnumMapSameMembers("Target")] enum Source { Value }`, "symbol reference"},
		{"missing argument", `[EnumMapSameMembers()] enum Source { Value }`, "exactly one"},
		{"extra arguments", `[EnumMapSameMembers(Target, Target)] enum Source { Value }`, "exactly one"},
		{"unqualified fallback", `[EnumMapUnknownMember(Value)] enum Source { Value }`, "qualified member"},
		{"foreign fallback", `[EnumMapUnknownMember(Target.Unknown)] enum Source { Value } enum Target { Unknown }`, "belong to the annotated enum"},
		{"missing member", `[EnumMapUnknownMember(Source.Missing)] enum Source { Value }`, "unknown enum member"},
		{"missing fallback", `[EnumMapSameMembers(Target)] enum Source { Value } enum Target { Value }`, "requires EnumMapUnknownMember"},
		{"conflicting fallbacks", `[EnumMapUnknownMember(Source.First)] [$go:EnumMapUnknownMember(Source.Second)] enum Source { First, Second }`, "conflicting fallback"},
		{"conflicting aliases", `[EnumMapSameMembers(Target)] enum Source { First = 1, Second = 1 } [EnumMapUnknownMember(Target.Unknown)] enum Target { Unknown, First, Second }`, "source aliases"},
		{"unmatched alias", `[EnumMapSameMembers(Target)] enum Source { First = 1, Other = 1 } [EnumMapUnknownMember(Target.Unknown)] enum Target { Unknown, First }`, "source aliases"},
		{"wrong scope", `[EnumMapSameMembers(Target)] model Source {}`, "only be applied to an enum"},
		{"selector without fallback", `[$go,$py:EnumMapSameMembers(Target)] enum Source { Value } [$go:EnumMapUnknownMember(Target.Unknown)] enum Target { Unknown }`, "requires EnumMapUnknownMember for py"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			_, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{SourceContent: scenario.source})
			var diagnostic *cclErrors.InvalidAttributeUsageError
			if !errors.As(err, &diagnostic) || diagnostic.SourcePosition == nil || !strings.Contains(err.Error(), scenario.diagnostic) {
				t.Fatalf("expected positioned attribute error containing %q, got %v", scenario.diagnostic, err)
			}
		})
	}
}

func TestEnumMappingDirectionsAndSelectors(t *testing.T) {
	definition, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{SourceContent: `
[$go,$py:EnumMapSameMembers(Target)]
[$go:EnumMapSameMembers(Target)]
enum Source { NameChanger = 7, Alias = 7, nameChanger = 8 }
[$go:EnumMapUnknownMember(Target.GoUnknown)]
[$py:EnumMapUnknownMember(Target.PythonUnknown)]
enum Target { GoUnknown = 99, PythonUnknown = 98, NameChanger = 1, Alias = 1 }
`})
	if err != nil {
		t.Fatal(err)
	}
	enums := definition.CodeContext.GetGenerationTypeDefinitions()
	source, target := enums[0].GetEnumDefinition(), enums[1].GetEnumDefinition()
	if len(source.Mappings) != 2 || len(target.Mappings) != 0 {
		t.Fatal("mappings must follow explicit directions and language selectors")
	}
	for language, fallback := range map[gValues.LanguageType]int64{gValues.LanguageGo: 99, gValues.LanguagePy: 98} {
		mappings := source.Mappings[language]
		if len(mappings) != 1 || mappings[0].Target != target || mappings[0].Fallback.Value != fallback {
			t.Fatalf("wrong destination or fallback for %s: %+v", language, mappings)
		}
		cases := mappings[0].Cases
		if len(cases) != 1 || cases[0].Source.Value != 7 || cases[0].Target.Value != 1 {
			t.Fatal("compatible aliases must share a case; member matching must be case sensitive")
		}
	}
	if len(source.Members) != 3 || len(target.Members) != 4 {
		t.Fatal("mapping resolution must preserve enum members")
	}
}
