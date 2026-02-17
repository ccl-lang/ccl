package sanitizer_test

import (
	"errors"
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
	"github.com/ccl-lang/ccl/src/cclSanitizer"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
)

const DuplicateModelsInput = `
model RoomInfo {
	id: int;
}

model RoomInfo {
	name: string;
}
`

func TestSanitizeDuplicateModels(t *testing.T) {
	astFile, err := cclParser.ParseCCLSourceContentAsAST(&cclParser.CCLParseOptions{
		SourceContent: DuplicateModelsInput,
	})
	if err != nil {
		t.Fatalf("Failed to parse AST: %v", err)
	}

	_, err = cclSanitizer.SanitizeCCLAst(nil, astFile)
	if err == nil {
		t.Fatalf("Expected duplicate model error, got nil")
	}

	var dupErr *cclErrors.DuplicateModelError
	if !errors.As(err, &dupErr) {
		t.Fatalf("Expected DuplicateModelError, got %T", err)
	}
}
