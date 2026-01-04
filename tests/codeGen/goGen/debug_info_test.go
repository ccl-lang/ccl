package goGen_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/cclParser"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
)

func TestGoGeneratorDebugInfo(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_go_test_debug")
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

	fmt.Printf("Generating Go code with debug info to: %s\n", tmpDir)

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
		CCLDefinition:     parsedDefinitions,
		OutputPath:        filepath.Join(tmpDir, "models"),
		TargetLanguage:    "go",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	// Verify .cclinfo files exist
	filesToCheck := []string{"methods.go", "types.go", "constants.go"}
	for _, file := range filesToCheck {
		infoPath := filepath.Join(tmpDir, "models", file+".cclinfo")
		if _, err := os.Stat(infoPath); os.IsNotExist(err) {
			t.Errorf("Expected debug info file %s to exist, but it does not", infoPath)
		} else {
			// Verify content
			data, err := os.ReadFile(infoPath)
			if err != nil {
				t.Errorf("Failed to read debug info file %s: %v", infoPath, err)
				continue
			}

			var debugInfos []*codeBuilder.DebugInfo
			err = json.Unmarshal(data, &debugInfos)
			if err != nil {
				t.Errorf("Failed to parse debug info file %s: %v", infoPath, err)
				continue
			}

			if len(debugInfos) == 0 {
				t.Errorf("Debug info file %s is empty", infoPath)
			}

			// Check some entries
			for i, info := range debugInfos {
				if info.SourceFile == "" {
					t.Errorf("Debug info entry %d in %s has empty SourceFile", i, infoPath)
				}
				if info.SourceLine == 0 {
					t.Errorf("Debug info entry %d in %s has 0 SourceLine", i, infoPath)
				}
				if info.GeneratedLine == 0 {
					t.Errorf("Debug info entry %d in %s has 0 GeneratedLine", i, infoPath)
				}

				// Basic sanity check: SourceFile should contain "goGenerator"
				if !strings.Contains(info.SourceFile, "goGenerator") {
					// It might be acceptable if it points elsewhere, but for now let's log it
					t.Logf("Warning: Debug info entry %d in %s points to %s", i, infoPath, info.SourceFile)
				}
			}
		}
	}
}
