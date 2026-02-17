package sanitizer_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
	"github.com/ccl-lang/ccl/src/cclSanitizer"
)

const AstSanitizeInput = `
#[MyAttributeGlobal("GlobalParam")]

model ModelA {
	fieldB: ModelB;
}

model ModelB {
	id: int;
}
`

func TestSanitizeCCLAst(t *testing.T) {
	astFile, err := cclParser.ParseCCLSourceContentAsAST(&cclParser.CCLParseOptions{
		SourceContent: AstSanitizeInput,
	})
	if err != nil {
		t.Fatalf("Failed to parse AST: %v", err)
	}

	definition, err := cclSanitizer.SanitizeCCLAst(nil, astFile)
	if err != nil {
		t.Fatalf("Failed to sanitize AST: %v", err)
	}

	if len(definition.GlobalAttributes) != 1 {
		t.Fatalf("Expected 1 global attribute, got %d", len(definition.GlobalAttributes))
	}

	models := definition.GetAllModels()
	if len(models) != 2 {
		t.Fatalf("Expected 2 models, got %d", len(models))
	}

	modelA := definition.GetModelByName("ModelA")
	modelB := definition.GetModelByName("ModelB")
	if modelA == nil || modelB == nil {
		t.Fatalf("Expected both ModelA and ModelB to exist")
	}

	fieldB := modelA.GetFieldByName("fieldB")
	if fieldB == nil {
		t.Fatalf("Expected fieldB on ModelA")
	}

	fieldTypeDef := fieldB.Type.GetDefinition()
	if fieldTypeDef == nil || fieldTypeDef.IsIncomplete() {
		t.Fatalf("Expected fieldB type to be complete, got %v", fieldTypeDef)
	}

	if fieldTypeDef.GetModelDefinition() != modelB {
		t.Fatalf("Expected fieldB type to resolve to ModelB")
	}
}
