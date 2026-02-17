package sanitizer_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/cclSanitizer"
	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func TestResolveAttributeUsageLiterals(t *testing.T) {
	ctx := cclValues.NewCCLCodeContext()
	node := &cclAst.AttributeNode{
		Name: "TestAttr",
		Params: []*cclAst.AttributeParamNode{
			{
				Value: &cclAst.LiteralValueExpression{
					LiteralKind: cclAst.AttributeLiteralKindString,
					Value:       "hello",
				},
			},
			{
				Name: "Count",
				Value: &cclAst.LiteralValueExpression{
					LiteralKind: cclAst.AttributeLiteralKindInt,
					Value:       3,
				},
			},
			{
				Value: &cclAst.LiteralValueExpression{
					LiteralKind:     cclAst.AttributeLiteralKindReserved,
					ReservedLiteral: cclValues.ReservedLiteralTrue,
					Value:           true,
				},
			},
		},
	}

	usage, err := cclSanitizer.ResolveAttributeUsage(ctx, node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if usage == nil || len(usage.Parameters) != 3 {
		t.Fatalf("expected 3 parameters, got %v", usage)
	}

	if usage.Parameters[0].ValueType == nil || usage.Parameters[0].ValueType.GetName() != cclValues.TypeNameString {
		t.Fatalf("expected string value type, got %v", usage.Parameters[0].ValueType)
	}

	if usage.Parameters[1].Name != "Count" {
		t.Fatalf("expected named parameter Count, got %s", usage.Parameters[1].Name)
	}

	if usage.Parameters[1].ValueType == nil || usage.Parameters[1].ValueType.GetName() != cclValues.TypeNameInt {
		t.Fatalf("expected int value type, got %v", usage.Parameters[1].ValueType)
	}

	if usage.Parameters[2].ValueType == nil || usage.Parameters[2].ValueType.GetName() != cclValues.TypeNameBool {
		t.Fatalf("expected bool value type, got %v", usage.Parameters[2].ValueType)
	}
}

func TestResolveAttributeUsageIdentifier(t *testing.T) {
	ctx := cclValues.NewCCLCodeContext()
	node := &cclAst.AttributeNode{
		Name: "TestAttr",
		Params: []*cclAst.AttributeParamNode{
			{
				Value: &cclAst.IdentifierValueExpression{
					Name: "__ccl_version",
				},
			},
		},
	}

	usage, err := cclSanitizer.ResolveAttributeUsage(ctx, node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if usage == nil || len(usage.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %v", usage)
	}

	if usage.Parameters[0].ValueType == nil || usage.Parameters[0].ValueType.GetName() != cclValues.TypeNameString {
		t.Fatalf("expected string value type, got %v", usage.Parameters[0].ValueType)
	}

	if _, ok := usage.Parameters[0].GetValue().(*cclValues.VariableUsageInstance); !ok {
		t.Fatalf("expected variable usage instance, got %T", usage.Parameters[0].GetValue())
	}
}
