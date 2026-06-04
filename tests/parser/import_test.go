package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func TestParseCCLSourceFileImportsRelativeToSourceFile(t *testing.T) {
	sourceDir := t.TempDir()
	otherDir := t.TempDir()

	nestedDir := filepath.Join(sourceDir, "schemas")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested source dir: %v", err)
	}

	importedPath := filepath.Join(nestedDir, "shared.ccl")
	if err := os.WriteFile(importedPath, []byte(`
model SharedThing {
	Value: string;
}
`), 0644); err != nil {
		t.Fatalf("Failed to write imported source file: %v", err)
	}

	mainPath := filepath.Join(sourceDir, "main.ccl")
	if err := os.WriteFile(mainPath, []byte(`
import "schemas/shared.ccl";

model LocalThing {
	Shared: SharedThing;
}
`), 0644); err != nil {
		t.Fatalf("Failed to write main source file: %v", err)
	}

	t.Chdir(otherDir)

	ctx := cclValues.NewCCLCodeContext()
	definition, err := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: mainPath,
		CodeContext:    ctx,
	})
	if err != nil {
		t.Fatalf("Failed to parse source file with import: %v", err)
	}

	models := definition.GetAllModels()
	if len(models) != 2 {
		t.Fatalf("Expected 2 models, got %d", len(models))
	}

	sharedModel := definition.GetModelByName("SharedThing")
	if sharedModel == nil {
		t.Fatalf("Expected imported model SharedThing to be available")
	}

	localModel := definition.GetModelByName("LocalThing")
	if localModel == nil {
		t.Fatalf("Expected local model LocalThing to be available")
	}

	sharedField := localModel.GetFieldByName("Shared")
	if sharedField == nil {
		t.Fatalf("Expected LocalThing.Shared field")
	}

	fieldType := sharedField.Type.GetDefinition()
	if fieldType.IsIncomplete() {
		t.Fatalf("Expected imported type usage to be complete")
	}

	if fieldType.GetModelDefinition() != sharedModel {
		t.Fatalf("Expected LocalThing.Shared to reference imported SharedThing")
	}

	if fieldType.GetSourceFileId() == 0 || fieldType.GetSourceFileId() != sharedModel.SourceFileId {
		t.Fatalf("Expected imported type definition to keep the imported source file id")
	}

	if sharedModel.SourceFileId == 0 || localModel.SourceFileId == 0 {
		t.Fatalf("Expected models to have source file ids")
	}

	if sharedModel.SourceFileId == localModel.SourceFileId {
		t.Fatalf("Expected imported and local models to have different source file ids")
	}

	if len(ctx.GetSourceDefinitions()) != 2 {
		t.Fatalf("Expected 2 registered source definitions, got %d", len(ctx.GetSourceDefinitions()))
	}
}
