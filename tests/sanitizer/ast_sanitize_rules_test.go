package sanitizer_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
	"github.com/ccl-lang/ccl/src/cclSanitizer"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
)

const DuplicateModelsInput = `
model Room {
	id: int;
}

model Room {
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

const DuplicateFieldsInput = `
model Room {
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

const FieldNameConflictModelInput = `
model Room {
	id: int;
}

model UserInfo {
	Room: string;
}
`

func TestSanitizeFieldNameConflictWithModel(t *testing.T) {
	astFile, err := cclParser.ParseCCLSourceContentAsAST(&cclParser.CCLParseOptions{
		SourceContent: FieldNameConflictModelInput,
	})
	if err != nil {
		t.Fatalf("Failed to parse AST: %v", err)
	}

	_, err = cclSanitizer.SanitizeCCLAst(nil, astFile)
	if err == nil {
		t.Fatalf("Expected field name conflict error, got nil")
	}

	var conflictErr *cclErrors.FieldNameConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("Expected FieldNameConflictError, got %T", err)
	}

	if conflictErr.ConflictKind != "model" {
		t.Fatalf("Expected model conflict, got %s", conflictErr.ConflictKind)
	}
}

const FieldNameConflictBuiltinInput = `
model UserInfo {
	INT32: string;
}
`

func TestSanitizeFieldNameConflictWithBuiltin(t *testing.T) {
	astFile, err := cclParser.ParseCCLSourceContentAsAST(&cclParser.CCLParseOptions{
		SourceContent: FieldNameConflictBuiltinInput,
	})
	if err != nil {
		t.Fatalf("Failed to parse AST: %v", err)
	}

	_, err = cclSanitizer.SanitizeCCLAst(nil, astFile)
	fmt.Println(err)
	if err == nil {
		t.Fatalf("Expected field name conflict error, got nil")
	}

	var conflictErr *cclErrors.FieldNameConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("Expected FieldNameConflictError, got %T", err)
	}

	if conflictErr.ConflictKind != "builtin-type" {
		t.Fatalf("Expected builtin-type conflict, got %s", conflictErr.ConflictKind)
	}
}
