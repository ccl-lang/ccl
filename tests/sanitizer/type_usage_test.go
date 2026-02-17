package sanitizer_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/cclSanitizer"
	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func TestResolveTypeUsageBuiltin(t *testing.T) {
	ctx := cclValues.NewCCLCodeContext()
	expr := &cclAst.SimpleTypeExpression{
		TypeName: cclAst.SimpleTypeName{
			Name: "string",
		},
		IsBuiltinToken: true,
	}

	usage, err := cclSanitizer.ResolveTypeUsage(ctx, "main", expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if usage == nil || !usage.IsBuiltIn() {
		t.Fatalf("expected built-in type usage, got %v", usage)
	}
}

func TestResolveTypeUsageArray(t *testing.T) {
	ctx := cclValues.NewCCLCodeContext()
	expr := &cclAst.ArrayTypeExpression{
		ElementType: &cclAst.SimpleTypeExpression{
			TypeName: cclAst.SimpleTypeName{
				Name: "int",
			},
			IsBuiltinToken: true,
		},
		Length: -1,
	}

	usage, err := cclSanitizer.ResolveTypeUsage(ctx, "main", expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if usage == nil || !usage.IsArray() {
		t.Fatalf("expected array type usage, got %v", usage)
	}

	if usage.GetUnderlyingType() == nil || !usage.GetUnderlyingType().IsBuiltIn() {
		t.Fatalf("expected array element to be built-in, got %v", usage.GetUnderlyingType())
	}
}

func TestResolveTypeUsageCustomType(t *testing.T) {
	ctx := cclValues.NewCCLCodeContext()
	expr := &cclAst.SimpleTypeExpression{
		TypeName: cclAst.SimpleTypeName{
			Name: "MyModel",
		},
		IsBuiltinToken: false,
	}

	usage, err := cclSanitizer.ResolveTypeUsage(ctx, "main", expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if usage == nil {
		t.Fatalf("expected custom type usage, got nil")
	}

	if usage.GetDefinition() == nil || !usage.GetDefinition().IsIncomplete() {
		t.Fatalf("expected incomplete custom type definition, got %v", usage.GetDefinition())
	}
}
