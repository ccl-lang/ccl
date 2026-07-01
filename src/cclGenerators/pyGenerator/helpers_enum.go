package pyGenerator

import (
	"fmt"
	"strconv"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func pythonEnumFileName(enumDef *CCLEnum) string {
	return gValues.StyleSnakeCase.ApplyStyle(enumDef.Name)
}

func pythonStorageTypeName(typeUsage *cclValues.CCLTypeUsage) string {
	if typeUsage.IsCustomTypeEnum() {
		return typeUsage.GetEnumBaseTypeName()
	}

	return typeUsage.GetName()
}

func pythonDefaultLiteral(value any) string {
	switch typedValue := value.(type) {
	case string:
		return strconv.Quote(typedValue)
	case bool:
		if typedValue {
			return "True"
		}
		return "False"
	default:
		return fmt.Sprintf("%v", typedValue)
	}
}
