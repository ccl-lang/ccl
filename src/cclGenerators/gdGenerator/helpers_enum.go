package gdGenerator

import (
	"fmt"
	"strconv"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func gdStorageTypeName(typeUsage *cclValues.CCLTypeUsage) string {
	if typeUsage.IsCustomTypeEnum() {
		return typeUsage.GetEnumBaseTypeName()
	}

	return typeUsage.GetName()
}

func gdEnumDeclarationName(enumDef *CCLEnum) string {
	if enumDef.IsNested() {
		return enumDef.Name
	}

	return enumDef.Name + "Enum"
}

func gdEnumTypeReference(typeUsage *cclValues.CCLTypeUsage) string {
	enumDef := typeUsage.GetDefinition().GetEnumDefinition()
	if enumDef.IsNested() {
		return enumDef.Name
	}

	return enumDef.Name + "." + gdEnumDeclarationName(enumDef)
}

func gdEnumCastSuffix(typeUsage *cclValues.CCLTypeUsage) string {
	if !typeUsage.IsCustomTypeEnum() {
		return ""
	}

	return " as " + gdEnumTypeReference(typeUsage)
}

func gdDefaultLiteral(value any) string {
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
