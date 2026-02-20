package gdGenerator

import (
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclUtils"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *GDScriptGenerationContext) generateSerializeJsonMethods(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	if err := c.generateSerializeJsonDictMethod(model, builder); err != nil {
		return err
	}
	if err := c.generateSerializeJsonMethod(builder); err != nil {
		return err
	}
	if err := c.generateDeserializeJsonDictMethod(model, builder); err != nil {
		return err
	}
	if err := c.generateDeserializeJsonMethod(builder); err != nil {
		return err
	}

	return nil
}

func (c *GDScriptGenerationContext) generateSerializeJsonDictMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("func serialize_json_dict() -> Dictionary:").
		Indent().
		WriteLine("var data: Dictionary = {}").
		NewLine()

	for _, field := range model.Fields {
		jsonName, err := c.getJsonFieldName(model, field)
		if err != nil {
			return err
		}

		if field.IsArray() {
			err = c.generateArraySerializeJson(field, jsonName, builder)
			if err != nil {
				return err
			}
			continue
		}

		c.generateFieldSerializeJson(field, jsonName, builder)
	}

	builder.WriteLine("return data").
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateSerializeJsonMethod(
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("func serialize_json() -> String:").
		Indent().
		WriteLine("return JSON.stringify(self.serialize_json_dict())").
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateDeserializeJsonDictMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"modelResult", modelResultName,
	)
	defer builder.UnmapVar(
		"modelResult",
	)

	builder.LineD("static func deserialize_json_dict(data: Dictionary) -> $model:").
		Indent().
		WriteLine("if not data:").
		Indent().
		WriteLine("return null").
		UnindentLine().
		LineD("var $modelResult = $model.new()").
		NewLine()

	for _, field := range model.Fields {
		jsonName, err := c.getJsonFieldName(model, field)
		if err != nil {
			return err
		}

		if field.IsArray() {
			err = c.generateArrayDeserializeJson(field, jsonName, builder)
			if err != nil {
				return err
			}
			continue
		}

		err = c.generateFieldDeserializeJson(field, jsonName, builder)
		if err != nil {
			return err
		}
	}

	builder.LineD("return $modelResult").
		UnindentLine()

	return nil
}

// generateDeserializeJsonMethod generates code for deserializing json method.
// The method expects `model` mapped-var to be present in the builder.
func (c *GDScriptGenerationContext) generateDeserializeJsonMethod(
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars(
		"model",
	)

	builder.LineD("static func deserialize_json(json_text: String) -> $model:").
		Indent().
		WriteLine(`if not json_text:`).
		Indent().
		WriteLine("return null").
		UnindentLine().
		WriteLine("var json = JSON.new()").
		WriteLine("var err = json.parse(json_text)").
		WriteLine("if err != OK:").
		Indent().
		LineD(`push_error("Failed to parse JSON for $model: " + json.get_error_message())`).
		WriteLine("return null").
		UnindentLine().
		WriteLine("if not (json.data is Dictionary):").
		Indent().
		LineD(`push_error("Expected JSON object for $model, got ", json.data)`).
		WriteLine("return null").
		UnindentLine().
		LineD(`return $model.deserialize_json_dict(json.data)`).
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateFieldSerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) {
	fieldRawName := cclUtils.ToSnakeCase(field.Name)
	resultField := "self." + fieldRawName

	builder.MapVarPairs(
		"field", resultField,
		"jsonName", jsonName,
	)
	defer builder.UnmapVar(
		"field",
		"jsonName",
	)

	switch field.Type.GetName() {
	case cclValues.TypeNameBytes:
		builder.LineD(`data["$jsonName"] = Marshalls.raw_to_base64($field)`)
	default:
		if field.IsCustomTypeModel() {
			builder.WriteLine(`@warning_ignore("incompatible_ternary")`).
				LineD(`data["$jsonName"] = $field.serialize_json_dict() if $field else null`)
		} else {
			builder.LineD(`data["$jsonName"] = $field`)
		}
	}

	builder.NewLine()
}

func (c *GDScriptGenerationContext) generateArraySerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldRawName := cclUtils.ToSnakeCase(field.Name)
	resultField := "self." + fieldRawName
	listName := fieldRawName + "_list"

	builder.MapVarPairs(
		"list", listName,
		"field", resultField,
		"jsonName", jsonName,
	)
	defer builder.UnmapVar(
		"list",
		"field",
		"jsonName",
	)

	builder.LineD("var $list = []").
		LineD("if $field != null:").
		Indent().
		LineD("for item in $field:").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameBytes:
		builder.LineD("$list.append(Marshalls.raw_to_base64(item))")
	default:
		if targetFieldType.IsCustomTypeModel() {
			// this will definitely break for < godot 4.0 (e.g. 3.x)
			// but then we are using lots of other features which are not available
			// for 3.x and lower anyway...
			builder.WriteLine(`@warning_ignore("incompatible_ternary")`).
				LineD("$list.append(item.serialize_json_dict() if item else null)")
		} else {
			// are we sure about this plain else? are all "non-custom types" out there
			// supported by json schema (or at least Godot's json serializer)?
			// well, if that's not the case, this is the place where it should be added
			// in future.
			builder.LineD("$list.append(item)")
		}
	}

	builder.Unindent().
		LineD(`data["$jsonName"] = $list`).
		Unindent().
		WriteLine("else:").
		Indent().
		LineD(`data["$jsonName"] = null`).
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateFieldDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldRawName := cclUtils.ToSnakeCase(field.Name)
	targetFieldTypeName := field.Type.GetName()
	resultField := modelResultName + "." + fieldRawName
	valueName := fieldRawName + "_value"
	modelName := field.OwnedBy.GetName()

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"jsonName", jsonName,
		"value", valueName,
		"field", resultField,
		"fieldT", targetFieldTypeName,
	)
	defer builder.UnmapVar(
		"jsonName",
		"value",
		"field",
		"fieldT",
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
			LineD("$field = int($value)").
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
	targetFieldTypeName := targetFieldType.GetName()
	fieldTargetLangType := c.getGDScriptType(field)
	fieldRawName := cclUtils.ToSnakeCase(field.Name)
	resultField := modelResultName + "." + fieldRawName
	valueName := fieldRawName + "_value"
	listName := fieldRawName + "_list"
	modelName := field.OwnedBy.GetName()

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"field", resultField,
		"fieldT", targetFieldTypeName,
		"fieldTargetT", fieldTargetLangType,
		"value", valueName,
		"jsonName", jsonName,
		"list", listName,
	)
	defer builder.UnmapVar(
		"field",
		"fieldT",
		"fieldTargetT",
		"value",
		"jsonName",
		"list",
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
			LineD("$list.append(int(item))").
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
