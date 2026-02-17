package sanitizer_test

import (
	"errors"
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
	"github.com/ccl-lang/ccl/src/cclSanitizer"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
)

const DuplicateFieldsInput = `
model RoomInfo {
	id: int;
	id: int;
}
`

func TestSanitizeDuplicateFields(t *testing.T) {
	astFile, err := cclParser.ParseCCLSourceContentAsAST(&cclParser.CCLParseOptions{
		SourceContent: DuplicateFieldsInput,
	})
	if err != nil {
		t.Fatalf("Failed to parse AST: %v", err)
	}

	_, err = cclSanitizer.SanitizeCCLAst(nil, astFile)
	if err == nil {
		t.Fatalf("Expected duplicate field error, got nil")
	}

	var dupErr *cclErrors.DuplicateFieldError
	if !errors.As(err, &dupErr) {
		t.Fatalf("Expected DuplicateFieldError, got %T", err)
	}
}
