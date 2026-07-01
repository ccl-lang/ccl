package tsGenerator

import (
	"fmt"
	"strconv"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func typeScriptStorageTypeName(typeUsage *cclValues.CCLTypeUsage) string {
	if typeUsage.IsCustomTypeEnum() {
		return typeUsage.GetEnumBaseTypeName()
	}

	return typeUsage.GetName()
}

func typeScriptDefaultLiteral(value any) string {
	switch typedValue := value.(type) {
	case string:
		return strconv.Quote(typedValue)
	case bool:
		if typedValue {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", typedValue)
	}
}
