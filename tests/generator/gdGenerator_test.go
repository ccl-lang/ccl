package generator_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators/gdGenerator"
)

func TestGdGenerator1(t *testing.T) {
	fieldName := gdGenerator.ToSnakeCase("RpcId")
	print(fieldName)

	fieldName = gdGenerator.ToPascalCase("rpc_id")
	print(fieldName)

	fieldName = gdGenerator.ToPascalCase("MySkin")
	print(fieldName)
}
