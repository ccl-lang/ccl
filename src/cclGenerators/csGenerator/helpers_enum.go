package csGenerator

import (
	"fmt"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func csharpStorageTypeName(typeUsage *cclValues.CCLTypeUsage) string {
	if typeUsage.IsCustomTypeEnum() {
		return typeUsage.GetEnumBaseTypeName()
	}

	return typeUsage.GetName()
}

func enumBinarySize(typeName string) string {
	switch typeName {
	case cclValues.TypeNameInt8, cclValues.TypeNameUint8:
		return "1"
	case cclValues.TypeNameInt16, cclValues.TypeNameUint16:
		return "2"
	case cclValues.TypeNameInt64, cclValues.TypeNameUint64:
		return "8"
	default:
		return "4"
	}
}

func csharpDefaultLiteral(value any) string {
	switch typedValue := value.(type) {
	case string:
		return fmt.Sprintf("%q", typedValue)
	case bool:
		if typedValue {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", typedValue)
	}
}
