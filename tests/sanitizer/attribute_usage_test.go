package sanitizer_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclSanitizer"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
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

func TestScopedAttributeResolution(t *testing.T) {
	ctx := cclValues.NewCCLCodeContext()
	definition, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
		CodeContext: ctx,
		SourceContent: `
#[SomeAttribute("global")]
#file:[$go,$js:SomeAttribute("file")]
namespace main;
#namespace:[$go:SomeAttribute("namespace-root")]
namespace main.users;
#namespace:[$go:SomeAttribute("namespace-users")]
model User {
	Id: int;
}
`,
	})
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	model := definition.GetModelByName("User")
	if model == nil {
		t.Fatalf("Expected User model")
	}

	attr := ctx.ResolveAttribute(
		gValues.LanguageGo,
		"SomeAttribute",
		&cclValues.AttributeResolutionSubject{Model: model},
		nil,
	)
	if attr == nil {
		t.Fatalf("Expected scoped attribute")
	}

	if attr.GetParamAt(0).GetAsString() != "file" {
		t.Fatalf("Expected file scoped attribute to win first, got %s", attr.GetParamAt(0).GetAsString())
	}

	attrs := ctx.GetSourceCodeDefinition(model.SourceFileId).FileAttributes
	ctx.GetSourceCodeDefinition(model.SourceFileId).FileAttributes = nil
	defer func() {
		ctx.GetSourceCodeDefinition(model.SourceFileId).FileAttributes = attrs
	}()

	attr = ctx.ResolveAttribute(
		gValues.LanguageGo,
		"SomeAttribute",
		&cclValues.AttributeResolutionSubject{Model: model},
		nil,
	)
	if attr == nil || attr.GetParamAt(0).GetAsString() != "namespace-users" {
		t.Fatalf("Expected child namespace override, got %v", attr)
	}
}

func TestAttributeLanguageSelectorResolution(t *testing.T) {
	ctx := cclValues.NewCCLCodeContext()
	definition, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
		CodeContext: ctx,
		SourceContent: `
#file:[$go,$js:SomeAttribute("go-js")]
#file:[$csharp:SomeAttribute("cs")]
model Item {
	Id: int;
}
`,
	})
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	model := definition.GetModelByName("Item")
	if model == nil {
		t.Fatalf("Expected Item model")
	}

	goAttr := ctx.ResolveAttribute(
		gValues.LanguageGo,
		"SomeAttribute",
		&cclValues.AttributeResolutionSubject{Model: model},
		nil,
	)
	if goAttr == nil || goAttr.GetParamAt(0).GetAsString() != "go-js" {
		t.Fatalf("Expected Go attribute, got %v", goAttr)
	}

	csAttr := ctx.ResolveAttribute(
		gValues.LanguageCS,
		"SomeAttribute",
		&cclValues.AttributeResolutionSubject{Model: model},
		nil,
	)
	if csAttr == nil || csAttr.GetParamAt(0).GetAsString() != "cs" {
		t.Fatalf("Expected CSharp attribute, got %v", csAttr)
	}
}
