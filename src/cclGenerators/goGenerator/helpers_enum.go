package goGenerator

import (
	"fmt"
	"strconv"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func goIntegerWriteCast(targetType *cclValues.CCLTypeUsage, expression string) string {
	typeName := targetType.GetName()
	if targetType.IsCustomTypeEnum() {
		typeName = targetType.GetEnumBaseTypeName()
	}

	switch typeName {
	case cclValues.TypeNameInt:
		return "int32(" + expression + ")"
	case cclValues.TypeNameUint:
		return "uint32(" + expression + ")"
	default:
		return expression
	}
}

func goBaseIntegerTypeForRead(targetType *cclValues.CCLTypeUsage) string {
	typeName := targetType.GetName()
	if targetType.IsCustomTypeEnum() {
		typeName = targetType.GetEnumBaseTypeName()
	}

	switch typeName {
	case cclValues.TypeNameInt:
		return "int32"
	case cclValues.TypeNameUint:
		return "uint32"
	default:
		if mappedType, ok := CCLTypesToGoTypes[typeName]; ok {
			return mappedType
		}
		return "int32"
	}
}

func goDefaultLiteral(value any) string {
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

func goModelHasDefaults(model *CCLModel) bool {
	for _, field := range model.Fields {
		if field.HasDefaultValue() {
			return true
		}
	}

	return false
}

func unsupportedGoFieldType(typeName, fieldName, modelName string) error {
	return &cclErrors.UnsupportedFieldTypeError{
		TypeName:       typeName,
		FieldName:      fieldName,
		ModelName:      modelName,
		TargetLanguage: CurrentLanguage.String(),
	}
}
