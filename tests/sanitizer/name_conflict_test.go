package sanitizer_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
	"github.com/ccl-lang/ccl/src/cclSanitizer"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
)

const FieldNameConflictModelInput = `
model RoomInfo {
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

	if conflictErr.Kind != cclErrors.ConflictKindModel {
		t.Fatalf("Expected model conflict, got %s", conflictErr.Kind)
	}
}

const FieldNameConflictBuiltinInput = `

#[CCLVersion("0.0.4")]
model UserInfo {
	blah_blah: string;	hall_hallo: string;      some_other_field:int; INT32: string; cool_stuff: int;
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

	if conflictErr.SourcePosition == nil {
		t.Fatalf("Expected SourcePosition to be set on conflict error")
	}

	if !strings.Contains(conflictErr.SourcePosition.SourceLine, "INT32: string") {
		t.Fatalf(
			"Expected SourceLine to include culprit, got: %q",
			conflictErr.SourcePosition.SourceLine,
		)
	}

	if conflictErr.Kind != cclErrors.ConflictKindBuiltinType {
		t.Fatalf("Expected builtin-type conflict, got %s", conflictErr.Kind)
	}
}
