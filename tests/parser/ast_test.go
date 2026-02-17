package parser_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
	"github.com/ccl-lang/ccl/src/core/cclAst"
)

const AstModelTestInput = `
#[MyAttributeGlobal("GlobalParam")]

[MyAttribute1("MyParam")]
model MyModel {
	[MyAttribute2("Param1", "Param2")]
	[MyAttribute3("Param1", "Param2", 1234)]
	myField: string[];
}
`

func TestParseAsAST(t *testing.T) {
	astFile, err := cclParser.ParseCCLSourceContentAsAST(&cclParser.CCLParseOptions{
		SourceContent: AstModelTestInput,
	})
	if err != nil {
		t.Fatalf("Failed to parse AST: %v", err)
	}

	if len(astFile.GlobalAttributes) != 1 {
		t.Fatalf("Expected 1 global attribute, got %d", len(astFile.GlobalAttributes))
	}

	if len(astFile.Models) != 1 {
		t.Fatalf("Expected 1 model, got %d", len(astFile.Models))
	}

	model := astFile.Models[0]
	if model.Name != "MyModel" {
		t.Fatalf("Expected model name MyModel, got %s", model.Name)
	}

	if len(model.Fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(model.Fields))
	}

	field := model.Fields[0]
	if field.Name != "myField" {
		t.Fatalf("Expected field name myField, got %s", field.Name)
	}

	arrayType, ok := field.Type.(*cclAst.ArrayTypeExpression)
	if !ok {
		t.Fatalf("Expected array type expression, got %T", field.Type)
	}

	simpleType, ok := arrayType.ElementType.(*cclAst.SimpleTypeExpression)
	if !ok {
		t.Fatalf("Expected simple type element, got %T", arrayType.ElementType)
	}

	if simpleType.TypeName.Name != "string" {
		t.Fatalf("Expected element type string, got %s", simpleType.TypeName.Name)
	}

	if len(field.Attributes) != 2 {
		t.Fatalf("Expected 2 field attributes, got %d", len(field.Attributes))
	}
}
