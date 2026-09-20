package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

func TestSkipFileLanguageSelection(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		directive string
		language  gValues.LanguageType
		skipped   bool
	}{
		{"gd", "#file:[$gd:SkipFile()]", gValues.LanguageGd, true},
		{"go", "#file:[$gd:SkipFile()]", gValues.LanguageGo, false},
		{"ts", "#file:[$gd:SkipFile()]", gValues.LanguageTS, false},
		{"no_target", "#file:[$gd:SkipFile()]", gValues.LanguageUnknown, false},
		{"multiple_languages", "#file:[$go,$gd:SkipFile()]", gValues.LanguageGd, true},
		{"multiple_languages_no_match", "#file:[$go,$gd:SkipFile()]", gValues.LanguageTS, false},
		{"alias", "#file:[$GDScript:SkipFile()]", gValues.LanguageGd, true},
		{"all_languages", "#file:[SkipFile()]", gValues.LanguageGo, true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			definition, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
				SourceContent: `#file:[OutputFileGroup("admins")]` + "\n" + testCase.directive + `
model AdminData { Value: string; }
`,
				TargetLanguage: testCase.language,
			})
			if err != nil {
				t.Fatalf("Failed to parse file: %v", err)
			}
			if gotSkipped := definition.GetModelByName("AdminData") == nil; gotSkipped != testCase.skipped {
				t.Fatalf("Expected skipped=%v, got %v", testCase.skipped, gotSkipped)
			}
		})
	}
}

func TestSkipFileStopsParsingAndDiscardsFile(t *testing.T) {
	source := `
#[SomeAttribute(UnknownVariable)]
#file:[OutputFileGroup("admins")]
import "missing_before.ccl";
model BeforeSkip { Value: UnknownType; }
enum BeforeSkipEnum { Value }
#file:[$gd:SkipFile()]
import "missing_after.ccl";
model InvalidBody { Value string;
`
	definition, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
		SourceContent:  source,
		TargetLanguage: gValues.LanguageGd,
	})
	if err != nil {
		t.Fatalf("Skipped content must not be parsed, resolved, or sanitized: %v", err)
	}
	if len(definition.TypeDefinitions) != 0 || len(definition.GlobalAttributes) != 0 ||
		len(definition.FileAttributes) != 0 || len(definition.NamespaceAttributes) != 0 {
		t.Fatal("Skipped file contributed declarations or attributes")
	}
	if _, err := cclParser.ParseCCLSourceContentAsAST(&cclParser.CCLParseOptions{
		SourceContent:  source,
		TargetLanguage: gValues.LanguageGo,
	}); err == nil {
		t.Fatal("Expected the invalid body to fail parsing for Go")
	}
}

func TestSkipFileInImportGraph(t *testing.T) {
	sourceDir := t.TempDir()
	mainPath := filepath.Join(sourceDir, "main.ccl")
	for name, content := range map[string]string{
		"main.ccl": `import "admins.ccl";
import "admins.ccl";
model ClientData { Value: string; }`,
		"admins.ccl": `#file:[OutputFileGroup("admins")]
#file:[$gd:SkipFile()]
import "admin_details.ccl";
model AdminData { Detail: AdminDetail; }`,
		"admin_details.ccl": `model AdminDetail { Value: string; }`,
	} {
		if err := os.WriteFile(filepath.Join(sourceDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	for _, language := range []gValues.LanguageType{gValues.LanguageGd, gValues.LanguageGo, gValues.LanguageTS} {
		t.Run(language.String(), func(t *testing.T) {
			options := &cclParser.CCLParseOptions{SourceFilePath: mainPath, TargetLanguage: language}
			definition, err := cclParser.ParseCCLSourceFile(options)
			if err != nil {
				t.Fatalf("Failed to parse import graph: %v", err)
			}
			ast, err := cclParser.ParseCCLSourceFileAsAST(options)
			if err != nil {
				t.Fatalf("Failed to parse import graph as AST: %v", err)
			}
			expectedModels := 3
			if language == gValues.LanguageGd {
				expectedModels = 1
			}
			if len(definition.GetAllModels()) != expectedModels || len(ast.Models) != expectedModels {
				t.Fatalf("Expected %d models in both IR and AST, got %d and %d",
					expectedModels, len(definition.GetAllModels()), len(ast.Models))
			}
			if definition.GetModelByName("ClientData") == nil {
				t.Fatal("Skipping an import must not skip the importing file")
			}
		})
	}
}

func TestSkipFileRejectsInvalidUsage(t *testing.T) {
	for _, testCase := range []struct {
		name    string
		source  string
		message string
	}{
		{"global", "#[SkipFile()]", "file-only"},
		{"namespace", "#namespace:[SkipFile()]", "file-only"},
		{"model", "[SkipFile()] model AdminData {}", "file-only"},
		{"parameters", "#file:[$gd:SkipFile(true)]", "does not accept parameters"},
		{"unknown_language", "#file:[$gd,$invalid:SkipFile()]", "unsupported attribute language"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
				SourceContent:  testCase.source,
				TargetLanguage: gValues.LanguageGd,
			})
			if err == nil || !strings.Contains(err.Error(), testCase.message) {
				t.Fatalf("Expected error containing %q, got %v", testCase.message, err)
			}
		})
	}
}
