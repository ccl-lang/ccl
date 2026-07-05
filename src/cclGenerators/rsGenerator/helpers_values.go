package rsGenerator

import (
	"fmt"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func rustPrimitiveDefaultLiteral(value any) string {
	switch typedValue := value.(type) {
	case string:
		return fmt.Sprintf("%q.to_string()", typedValue)
	case bool:
		if typedValue {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", typedValue)
	}
}

func rustStorageTypeName(targetType *cclValues.CCLTypeUsage) string {
	if targetType.IsCustomTypeEnum() {
		return targetType.GetEnumBaseTypeName()
	}
	return targetType.GetName()
}

func isRustBytesType(targetType *cclValues.CCLTypeUsage) bool {
	return targetType.GetName() == cclValues.TypeNameBytes
}

func isRustBytesArrayType(targetType *cclValues.CCLTypeUsage) bool {
	return targetType.IsArray() &&
		targetType.GetUnderlyingType().GetName() == cclValues.TypeNameBytes
}
