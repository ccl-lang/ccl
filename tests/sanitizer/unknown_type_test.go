package sanitizer_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
	"github.com/ccl-lang/ccl/src/cclSanitizer"
)

const MissingTypeInput = `
model ModelA {
	fieldB: MissingType;
}
`

const MissingArrayTypeInput = `
model ModelA {
	fieldB: MissingType[];
}
`

func TestSanitizeCCLAstMissingType(t *testing.T) {
	err := sanitizeInput(t, MissingTypeInput)
	assertMissingTypeError(t, err, "MissingType")
}

func TestSanitizeCCLAstMissingArrayElementType(t *testing.T) {
	err := sanitizeInput(t, MissingArrayTypeInput)
	assertMissingTypeError(t, err, "MissingType")
}

func sanitizeInput(t *testing.T, input string) error {
	t.Helper()

	astFile, err := cclParser.ParseCCLSourceContentAsAST(&cclParser.CCLParseOptions{
		SourceContent: input,
	})
	if err != nil {
		t.Fatalf("Failed to parse AST: %v", err)
	}

	_, err = cclSanitizer.SanitizeCCLAst(nil, astFile)
	return err
}

func assertMissingTypeError(t *testing.T, err error, missingType string) {
	t.Helper()

	if err == nil {
		t.Fatalf("Expected error for missing type %q, got nil", missingType)
	}

	var sanitizeErr *cclSanitizer.AstSanitizationError
	if !errors.As(err, &sanitizeErr) {
		t.Fatalf("Expected AstSanitizationError, got %T", err)
	}

	message := err.Error()
	if !strings.Contains(message, "unknown type") {
		t.Fatalf("Expected missing type error message, got: %v", err)
	}

	if !strings.Contains(message, missingType) {
		t.Fatalf("Expected missing type %q in message, got: %v", missingType, err)
	}
}
