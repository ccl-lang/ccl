package parser_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
)

// ---------------------------------------------------------
const ModelTestInput1 = `
#[MyAttributeGlobal("GlobalParam")]

[MyAttribute1("MyParam")]
model MyModel {
	[MyAttribute2("Param1", "Param2")]
	[MyAttribute3("Param1", "Param2", 1234)]
	myField: string[];
}
	`

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

	models := cclSource.GetAllModels()

	if len(models) != 1 {
		t.Fatalf("Expected 1 model, got %d", len(models))
		return
	}

	model := models[0]
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

	if len(field.Attributes) != 2 {
		t.Fatalf("Expected 2 attributes on field, got %d", len(field.Attributes))
		return
	}

	if field.Attributes[0].Name != "MyAttribute2" {
		t.Errorf("Expected first attribute name MyAttribute2, got %s", field.Attributes[0].Name)
		return
	}

	if field.Attributes[1].Name != "MyAttribute3" {
		t.Errorf("Expected second attribute name MyAttribute3, got %s", field.Attributes[1].Name)
		return
	}
}

// ---------------------------------------------------------

const ModelTestInput2 = `
#[MyAttributeGlobal("GlobalParam")]
#[CCLVersion("1.0.0")]

[MyAttribute1("MyParam")]
[MyAttribute2("MyParam2")]
model MyModel1 {
	[MyAttribute2("Param1", "Param2")]
	[MyAttribute3("Param1", "Param2", 1234)]
	myField1: string[];


	myField2: MyModel2[];
}

model MyModel2 {
	fieldA: int;
	fieldB: float;
}
`

func TestModelTestInput2(t *testing.T) {
	cclSource, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
		SourceContent: ModelTestInput2,
	})
	if err != nil {
		t.Fatalf("Failed to parse CCL source: %v", err)
		return
	}

	if len(cclSource.GlobalAttributes) != 2 {
		t.Fatalf("Expected 2 global attribute, got %d", len(cclSource.GlobalAttributes))
		return
	}

	models := cclSource.GetAllModels()

	if len(models) != 2 {
		t.Fatalf("Expected 2 models, got %d", len(models))
		return
	}

	model1 := models[0]
	if model1.Name != "MyModel1" {
		t.Errorf("Expected model name MyModel1, got %s", model1.Name)
		return
	}

	model2 := models[1]
	if model2.Name != "MyModel2" {
		t.Errorf("Expected model name MyModel2, got %s", model2.Name)
		return
	}

	// because myField2 is of type MyModel2[], hence an array,
	// we need to get the underlying type
	field2TypeDef := model1.GetFieldByName("myField2").Type.GetUnderlyingType().GetDefinition()
	if field2TypeDef.IsIncomplete() {
		t.Errorf("Expected myField2 type to be complete, but it is incomplete")
		return
	}

	field2Model := field2TypeDef.GetModelDefinition()
	if model2 != field2Model {
		t.Errorf("Expected myField2 to be of type MyModel2, got %s", field2Model.Name)
		return
	}
}

// ---------------------------------------------------------
// ---------------------------------------------------------
// ---------------------------------------------------------
// ---------------------------------------------------------
// ---------------------------------------------------------
