package parser_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
)

const ModelTestInput1 = `
[MyAttribute1("MyParam")]
model MyModel {
	[MyAttribute2("Param1", "Param2")]
	[MyAttribute3("Param1", "Param2", 1234)]
	myField: string[];
}
`

// TODO
var ModelTestInput1Expected = []attributeInfo{}

func TestModelTestInput1(t *testing.T) {
	cclSource, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
		SourceContent: ModelTestInput1,
	})
	if err != nil {
		t.Fatalf("Failed to parse CCL source: %v", err)
		return
	}

	if len(cclSource.GlobalAttributes) != 1 {
		t.Fatalf("Expected 1 global attribute, got %d", len(cclSource.GlobalAttributes))
		return
	}

	if len(cclSource.Models) != 1 {
		t.Fatalf("Expected 1 model, got %d", len(cclSource.Models))
		return
	}

	model := cclSource.Models[0]
	if model.Name != "MyModel" {
		t.Errorf("Expected model name MyModel, got %s", model.Name)
		return
	}

	if len(model.Fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(model.Fields))
		return
	}

	field := model.Fields[0]
	if field.Name != "myField" {
		t.Errorf("Expected field name myField, got %s", field.Name)
		return
	}
}
