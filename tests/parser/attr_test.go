package parser_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser"
)

type attributeParamInfo struct {
	Name  string
	Value any
}

type attributeInfo struct {
	Name       string
	Parameters []attributeParamInfo
}

const AttrTestInput1 = `
#[MyAttribute1("Hello")]
#[MyAttribute2("Param1", "Param2")]
#[MyAttribute3("Param1", "Param2", 1234)]
`

var AttrTestInput1Expected = []attributeInfo{
	{
		Name: "MyAttribute1",
		Parameters: []attributeParamInfo{
			{
				Name:  "",
				Value: "Hello",
			},
		},
	},
	{
		Name: "MyAttribute2",
		Parameters: []attributeParamInfo{
			{
				Name:  "",
				Value: "Param1",
			},
			{
				Name:  "",
				Value: "Param2",
			},
		},
	},
	{
		Name: "MyAttribute3",
		Parameters: []attributeParamInfo{
			{
				Name:  "",
				Value: "Param1",
			},
			{
				Name:  "",
				Value: "Param2",
			},
			{
				Name:  "",
				Value: 1234,
			},
		},
	},
}

func TestAttributeParse1(t *testing.T) {
	cclSource, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{
		SourceContent: AttrTestInput1,
	})
	if err != nil {
		t.Fatalf("Failed to parse CCL source: %v", err)
		return
	}

	if len(cclSource.GlobalAttributes) != 3 {
		t.Fatalf("Expected 3 global attribute, got %d", len(cclSource.GlobalAttributes))
		return
	}

	for i, attr := range cclSource.GlobalAttributes {
		if attr.Name != AttrTestInput1Expected[i].Name {
			t.Errorf("Expected attribute name %s, got %s", AttrTestInput1Expected[i].Name, attr.Name)
		}

		if len(attr.Parameters) != len(AttrTestInput1Expected[i].Parameters) {
			t.Errorf("Expected %d parameters, got %d", len(AttrTestInput1Expected[i].Parameters), len(attr.Parameters))
			continue
		}

		for j, param := range attr.Parameters {
			if param.Name != AttrTestInput1Expected[i].Parameters[j].Name {
				t.Errorf("Expected parameter name %s, got %s", AttrTestInput1Expected[i].Parameters[j].Name, param.Name)
			}
			if !param.CompareValue(AttrTestInput1Expected[i].Parameters[j].Value) {
				t.Errorf(
					"Expected parameter value %v, got %v",
					AttrTestInput1Expected[i].Parameters[j].Value,
					param.GetValue(),
				)
			}
		}
	}
}
