package cclUtils_test

import (
	"testing"

	"github.com/ccl-lang/ccl/src/core/cclUtils"
)

func TestSnakeCase1(t *testing.T) {
	fieldName := cclUtils.ToSnakeCase("RpcId")
	print(fieldName)

	fieldName = cclUtils.ToPascalCase("rpc_id")
	print(fieldName)

	fieldName = cclUtils.ToPascalCase("MySkin")
	print(fieldName)
}
