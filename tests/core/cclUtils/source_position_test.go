package cclUtils_test

import (
	"strings"
	"testing"

	"github.com/ccl-lang/ccl/src/core/cclUtils"
)

func TestFormatErrorShowsCulpritForLongLine(t *testing.T) {
	longLine := "model UserInfo { " +
		strings.Repeat("a", 90) +
		" INT32: string; " +
		strings.Repeat("b", 90) +
		" }"

	column := strings.Index(longLine, "INT32") + 1
	if column <= 0 {
		t.Fatalf("failed to locate field name in test line")
	}

	pos := &cclUtils.SourceCodePosition{
		Line:       1,
		Column:     column,
		SourceLine: longLine,
	}

	formatted := pos.FormatError("Field name 'INT32' conflicts")

	if !strings.Contains(formatted, "INT32: string") {
		t.Fatalf("expected formatted error to include culprit snippet, got:\n%s", formatted)
	}

	if !strings.Contains(formatted, cclUtils.SourceErrorEllipsis) {
		t.Fatalf("expected formatted error to include truncation ellipsis, got:\n%s", formatted)
	}
}
