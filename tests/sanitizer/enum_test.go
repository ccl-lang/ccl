package sanitizer_test

import (
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

func TestEnumRejectsNegativeUnsignedValue(t *testing.T) {
	_, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
		SourceContent: `
enum BadEnum: uint32 {
	Bad = -1,
}
`,
	})
	if err == nil {
		t.Fatalf("Expected unsigned enum value validation error")
	}

	if !strings.Contains(err.Error(), "outside the base type range") {
		t.Fatalf("Expected enum range error, got: %v", err)
	}
}
