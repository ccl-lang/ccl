package gdGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *GDScriptGenerationContext) generateFieldDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldRawName := caseUtils.ToSnakeCase(field.Name)
	targetFieldTypeName := gdStorageTypeName(field.Type)
	resultField := modelResultName + "." + fieldRawName
	valueName := fieldRawName + "_value"
	modelName := field.OwnedBy.GetName()
	enumCastSuffix, err := c.getGDScriptEnumCastSuffix(field.Type)
	if err != nil {
		return err
	}

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"jsonName", jsonName,
		"value", valueName,
		"field", resultField,
		"fieldT", targetFieldTypeName,
		"enumCast", enumCastSuffix,
	)
	defer builder.UnmapVar(
		"jsonName",
		"value",
		"field",
		"fieldT",
		"enumCast",
	)

	// TODO #21: maybe we can have default value (from ccl) instead of null here?
	builder.
		LineD(`var $value = data.get("$jsonName", null)`).
		LineD("if $value == null:").
		Indent().
		LineD("pass").
		Unindent()

	switch targetFieldTypeName {
	case cclValues.TypeNameString:
		builder.LineD("elif $value is String:").
			Indent().
			LineD("$field = $value").
			Unindent().
			LineD("else:").
			Indent().
			LineD("$field = str($value)").
			Unindent()
	case cclValues.TypeNameInt, cclValues.TypeNameInt8, cclValues.TypeNameInt16,
		cclValues.TypeNameInt32, cclValues.TypeNameInt64,
		cclValues.TypeNameUint, cclValues.TypeNameUint8, cclValues.TypeNameUint16,
		cclValues.TypeNameUint32, cclValues.TypeNameUint64, cclValues.TypeNameDateTime:
		builder.WriteLine("elif (").
			Indent().
			LineD("$value is int or $value is float or").
			LineD("(($value is String or $value is StringName) and").
			LineD("$value.is_valid_int())").
			Unindent().
			WriteLine("):").
			Indent().
			LineD("$field = int($value)$enumCast").
			Unindent().
			WriteLine("else:").
			Indent().
			LineD(`push_error("Expected number for field $field in " +`).
			Indent().
			LineD(`"$model, but got ", $value)`).
			Unindent().
			WriteLine("return null").
			Unindent()
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteLine("elif (").
			Indent().
			LineD("$value is int or $value is float or").
			LineD("(($value is String or $value is StringName) and").
			LineD("$value.is_valid_float())").
			Unindent().
			WriteLine("):").
			Indent().
			LineD("$field = float($value)").
			Unindent().
			WriteLine("else:").
			Indent().
			LineD(`push_error("Expected float for field $field in " +`).
			Indent().
			LineD(`"$model, but got ", $value)`).
			Unindent().
			WriteLine("return null").
			Unindent()
	case cclValues.TypeNameBool:
		builder.LineD("elif $value is bool:").
			Indent().
			LineD("$field = $value").
			Unindent().
			LineD("elif $value is String or $value is StringName:").
			Indent().
			LineD("$field = (").
			Indent().
			LineD("!$value.is_empty() and").
			LineD(`$value != "0" and`).
			LineD(`$value.to_lower() != "false"`).
			Unindent().
			WriteLine(")").
			Unindent().
			LineD("elif $value is int or $value is float:").
			Indent().
			LineD("$field = bool($value)").
			Unindent().
			WriteLine("else:").
			Indent().
			LineD(`push_error("Expected bool for field $field in " +`).
			Indent().
			LineD(`"$model, but got ", $value)`).
			Unindent().
			WriteLine("return null").
			Unindent()
	case cclValues.TypeNameBytes:
		builder.LineD("elif $value is String:").
			Indent().
			LineD("$field = Marshalls.base64_to_raw($value)").
			Unindent().
			LineD("elif $value is PackedByteArray:").
			Indent().
			LineD("$field = $value").
			Unindent().
			WriteLine("else:").
			Indent().
			LineD(`push_error("Expected base64 string for field $field in " +`).
			Indent().
			LineD(`"$model, but got ", $value)`).
			Unindent().
			WriteLine("return null").
			Unindent()
	default:
		if field.IsCustomTypeModel() {
			builder.LineD("elif $value is Dictionary:").
				Indent().
				LineD("$field = $fieldT.deserialize_json_dict($value)").
				Unindent().
				WriteLine("else:").
				Indent().
				LineD(`push_error("Expected json object (Dictionary) for field $field in " +`).
				Indent().
				LineD(`"$model, but got ", $value)`).
				Unindent().
				WriteLine("return null").
				Unindent()
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldTypeName,
				ModelName:      modelName,
				FieldName:      fieldRawName,
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}

	builder.NewLine()
	return nil
}

func (c *GDScriptGenerationContext) generateArrayDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	targetFieldTypeName := gdStorageTypeName(targetFieldType)
	fieldTargetLangType, err := c.getGDScriptType(field)
	if err != nil {
		return err
	}
	fieldRawName := caseUtils.ToSnakeCase(field.Name)
	resultField := modelResultName + "." + fieldRawName
	valueName := fieldRawName + "_value"
	listName := fieldRawName + "_list"
	modelName := field.OwnedBy.GetName()
	enumCastSuffix, err := c.getGDScriptEnumCastSuffix(targetFieldType)
	if err != nil {
		return err
	}

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"field", resultField,
		"fieldT", targetFieldTypeName,
		"fieldTargetT", fieldTargetLangType,
		"value", valueName,
		"jsonName", jsonName,
		"list", listName,
		"enumCast", enumCastSuffix,
	)
	defer builder.UnmapVar(
		"field",
		"fieldT",
		"fieldTargetT",
		"value",
		"jsonName",
		"list",
		"enumCast",
	)

	builder.LineD(`var $value = data.get("$jsonName", null)`).
		LineD("if $value == null:").
		Indent().
		WriteLine("# we can't iterate over it later if it stays null").
		LineD("$value = []").
		Unindent().
		LineD("elif not ($value is Array):").
		Indent().
		LineD(`push_error("Expected array type for field $field in " +`).
		Indent().
		LineD(`"$model, but got ", $value)`).
		Unindent().
		WriteLine("return null").
		Unindent()

	builder.LineD("var $list = [] as $fieldTargetT").
		LineD("for item in $value:").
		Indent()

	switch targetFieldTypeName {
	case cclValues.TypeNameString:
		builder.LineD("if item is String:").
			Indent().
			LineD("$list.append(item)").
			Unindent().
			LineD("else:").
			Indent().
			LineD("$list.append(str(item))").
			Unindent()
	case cclValues.TypeNameInt, cclValues.TypeNameInt8, cclValues.TypeNameInt16,
		cclValues.TypeNameInt32, cclValues.TypeNameInt64,
		cclValues.TypeNameUint, cclValues.TypeNameUint8, cclValues.TypeNameUint16,
		cclValues.TypeNameUint32, cclValues.TypeNameUint64, cclValues.TypeNameDateTime:
		builder.WriteLine("if (").
			Indent().
			LineD("item is int or item is float or").
			LineD("((item is String or item is StringName) and").
			LineD("item.is_valid_int())").
			Unindent().
			WriteLine("):").
			Indent().
			LineD("$list.append(int(item)$enumCast)").
			Unindent().
			WriteLine("else:").
			Indent().
			LineD(`push_error("Expected number for every element of field $field in " +`).
			Indent().
			LineD(`"$model, but got ", item)`).
			Unindent().
			WriteLine("return null").
			Unindent()
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteLine("if (").
			Indent().
			WriteLine("item is int or item is float or").
			WriteLine("((item is String or item is StringName) and").
			WriteLine("item.is_valid_float())").
			Unindent().
			WriteLine("):").
			Indent().
			LineD("$list.append(float(item))").
			Unindent().
			WriteLine("else:").
			Indent().
			LineD(`push_error("Expected float for every element of field $field in " +`).
			Indent().
			LineD(`"$model, but got ", item)`).
			Unindent().
			WriteLine("return null").
			Unindent()
	case cclValues.TypeNameBool:
		builder.WriteLine("if item is bool:").
			Indent().
			LineD("$list.append(item)").
			Unindent().
			WriteLine("elif item is String or item is StringName:").
			Indent().
			LineD("$list.append(").
			Indent().
			WriteLine("!item.is_empty() and").
			WriteLine(`item != "0" and`).
			WriteLine(`item.to_lower() != "false"`).
			Unindent().
			WriteLine(")").
			Unindent().
			LineD("elif item is int or item is float:").
			Indent().
			LineD("$list.append(bool(item))").
			Unindent().
			LineD("elif item == null:").
			Indent().
			LineD("$list.append(false)").
			Unindent().
			WriteLine("else:").
			Indent().
			LineD(`push_error("Expected bool for every element of field $field in " +`).
			Indent().
			LineD(`"$model, but got ", item)`).
			Unindent().
			WriteLine("return null").
			Unindent()
	case cclValues.TypeNameBytes:
		builder.WriteLine("if item is String:").
			Indent().
			LineD("$list.append(Marshalls.base64_to_raw(item))").
			Unindent().
			WriteLine("elif item is PackedByteArray:").
			Indent().
			LineD("$list.append(item)").
			Unindent().
			WriteLine("else:").
			Indent().
			LineD(`push_error("Expected base64 string for every element of field $field in " +`).
			Indent().
			LineD(`"$model, but got ", item)`).
			Unindent().
			WriteLine("return null").
			Unindent()
	default:
		if targetFieldType.IsCustomTypeModel() {
			// in case you are asking yourself why this condition isn't before the switch:
			// because null can't be assigned to ALL of the primitive types in gdscript
			builder.WriteLine("if item == null:").
				Indent().
				LineD("$list.append(null)").
				Unindent().
				WriteLine("elif not (item is Dictionary):").
				Indent().
				LineD(`push_error("Expected object items for every element of field $field in " +`).
				Indent().
				LineD(`"$model, but got ", item)`).
				Unindent().
				WriteLine("return null").
				Unindent().
				WriteLine("else:").
				Indent().
				LineD("$list.append($fieldT.deserialize_json_dict(item))").
				Unindent()
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldTypeName,
				ModelName:      modelName,
				FieldName:      fieldRawName,
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}

	builder.Unindent().
		LineD("$field = $list").
		NewLine()

	return nil
}
