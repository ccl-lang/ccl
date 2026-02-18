package cclUtils_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/core/cclUtils"
)

func TestSnakeCase1(t *testing.T) {
	fieldName := cclUtils.ToSnakeCase("RpcId")
	if fieldName != "rpc_id" {
		t.Errorf("Expected 'rpc_id', got '%s'", fieldName)
	}

	fieldName = cclUtils.ToSnakeCase("Rpc-Id")
	if fieldName != "rpc_id" {
		t.Errorf("Expected 'rpc_id', got '%s'", fieldName)
	}

	fieldName = cclUtils.ToSnakeCase("some_property")
	if fieldName != "some_property" {
		t.Errorf("Expected 'some_property', got '%s'", fieldName)
	}
}

func TestPascalCase1(t *testing.T) {
	fieldName := cclUtils.ToPascalCase("rpc_id")
	if fieldName != "RpcId" {
		t.Errorf("Expected 'RpcId', got '%s'", fieldName)
	}

	fieldName = cclUtils.ToPascalCase("MySkin")
	if fieldName != "MySkin" {
		t.Errorf("Expected 'MySkin', got '%s'", fieldName)
	}

	fieldName = cclUtils.ToPascalCase("someProperty")
	if fieldName != "SomeProperty" {
		t.Errorf("Expected 'SomeProperty', got '%s'", fieldName)
	}

	fieldName = cclUtils.ToPascalCase("Rpc-Id")
	if fieldName != "RpcId" {
		t.Errorf("Expected 'RpcId', got '%s'", fieldName)
	}
}

func TestCamelCase1(t *testing.T) {
	fieldName := cclUtils.ToCamelCase("rpc_id")
	if fieldName != "rpcId" {
		t.Errorf("Expected 'rpcId', got '%s'", fieldName)
	}

	fieldName = cclUtils.ToCamelCase("MySkin")
	if fieldName != "mySkin" {
		t.Errorf("Expected 'mySkin', got '%s'", fieldName)
	}

	fieldName = cclUtils.ToCamelCase("some_property")
	if fieldName != "someProperty" {
		t.Errorf("Expected 'someProperty', got '%s'", fieldName)
	}

	fieldName = cclUtils.ToCamelCase("SomeProperty")
	if fieldName != "someProperty" {
		t.Errorf("Expected 'someProperty', got '%s'", fieldName)
	}
}
