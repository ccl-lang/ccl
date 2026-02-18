package gdGenerator

import (
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
	if err := c.generateDeserializeJsonMethod(model, builder); err != nil {
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
			c.generateArraySerializeJson(field, jsonName, builder)
			continue
		}

		c.generateFieldSerializeJson(field, jsonName, builder)
	}

	builder.WriteLine("return data").
		UnindentLine().
		NewLine()
	return nil
}

func (c *GDScriptGenerationContext) generateSerializeJsonMethod(
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("func serialize_json() -> String:").
		Indent().
		WriteLine("return JSON.stringify(serialize_json_dict())").
		UnindentLine().
		NewLine()
	return nil
}

func (c *GDScriptGenerationContext) generateDeserializeJsonDictMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("static func deserialize_json_dict(data: Dictionary) -> " + model.Name + ":").
		Indent().
		WriteLine("if data == null:").
		Indent().
		WriteLine("return null").
		UnindentLine().
		WriteLine("var model_result = " + model.Name + ".new()").
		NewLine()

	for _, field := range model.Fields {
		jsonName, err := c.getJsonFieldName(model, field)
		if err != nil {
			return err
		}

		if field.IsArray() {
			c.generateArrayDeserializeJson("model_result", model.Name, field, jsonName, builder)
			continue
		}

		c.generateFieldDeserializeJson("model_result", model.Name, field, jsonName, builder)
	}

	builder.WriteLine("return model_result").
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateDeserializeJsonMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("static func deserialize_json(json_text: String) -> " + model.Name + ":").
		Indent().
		WriteLine("if json_text == \"\":").
		Indent().
		WriteLine("return null").
		UnindentLine().
		WriteLine("var json = JSON.new()").
		WriteLine("var err = json.parse(json_text)").
		WriteLine("if err != OK:").
		Indent().
		WriteLine("push_error(\"Failed to parse JSON for " + model.Name + ": \" + json.get_error_message())").
		WriteLine("return null").
		UnindentLine().
		WriteLine("if typeof(json.data) != TYPE_DICTIONARY:").
		Indent().
		WriteLine("push_error(\"Expected JSON object for " + model.Name + "\")").
		WriteLine("return null").
		UnindentLine().
		WriteLine("return " + model.Name + ".deserialize_json_dict(json.data)").
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateFieldSerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) {
	fieldRawName := cclUtils.ToSnakeCase(field.Name)
	fieldName := "self." + fieldRawName

	switch field.Type.GetName() {
	case cclValues.TypeNameBytes:
		builder.WriteLine("data[\"" + jsonName + "\"] = Marshalls.raw_to_base64(" + fieldName + ")")
	default:
		if field.Type.IsCustomTypeModel() {
			builder.WriteLine("data[\"" + jsonName + "\"] = " +
				fieldName + ".serialize_json_dict() if " + fieldName + " else null")
		} else {
			builder.WriteLine("data[\"" + jsonName + "\"] = " + fieldName)
		}
	}

	builder.NewLine()
}

func (c *GDScriptGenerationContext) generateArraySerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldRawName := cclUtils.ToSnakeCase(field.Name)
	fieldName := "self." + fieldRawName
	listName := fieldRawName + "_list"

	builder.WriteLine("var " + listName + " = []")
	builder.WriteLine("if " + fieldName + " != null:")
	builder.Indent()
	builder.WriteLine("for item in " + fieldName + ":")
	builder.Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameBytes:
		builder.WriteLine(listName + ".append(Marshalls.raw_to_base64(item))")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine(listName + ".append(item.serialize_json_dict() if item else null)")
		} else {
			builder.WriteLine(listName + ".append(item)")
		}
	}

	builder.Unindent()
	builder.Unindent()
	builder.WriteLine("data[\"" + jsonName + "\"] = " + listName)
	builder.NewLine()
}

func (c *GDScriptGenerationContext) generateFieldDeserializeJson(
	resultName string,
	modelName string,
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) {
	fieldRawName := cclUtils.ToSnakeCase(field.Name)
	resultField := resultName + "." + fieldRawName
	valueName := fieldRawName + "_value"

	builder.MapVar("jsonName", jsonName).
		MapVar("value", valueName).
		MapVar("field", resultField).
		MapVar("fieldT", field.Type.GetName()).
		MapVar("model", modelName).
		// TODO #21: maybe we can have default value (from ccl) instead of null here?
		LineD("var $value = data.get(\"$jsonName\", null)").
		LineD("if $value == null:").
		Indent().
		// TODO #21: maybe we can have default value (from ccl) instead of pass here?
		LineD("pass").
		Unindent()

	defer func() {
		builder.UnmapVar(
			"jsonName",
			"value",
			"field",
			"fieldT",
			"model",
		)
	}()

	switch field.Type.GetName() {
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
		builder.WriteLine("if (").
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
			LineD("push_error(\"Expected number for field $field in \" +").
			Indent().
			LineD("\"$model, but got \", $value)").
			Unindent().
			WriteLine("return null").
			Unindent()
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteLine("if (").
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
			LineD("push_error(\"Expected float for field $field in \" +").
			Indent().
			LineD("\"$model, but got \", $value)").
			Unindent().
			WriteLine("return null").
			Unindent()
	case cclValues.TypeNameBool:
		builder.LineD("if $value is bool:").
			Indent().
			LineD("$field = $value").
			Unindent().
			LineD("elif $value is String or $value is StringName:").
			Indent().
			LineD("$field =  (").
			Indent().
			LineD("!$value.is_empty() and").
			LineD("$value != \"0\" and").
			LineD("$value.to_lower() != \"false\"").
			Unindent().
			WriteLine(")").
			Unindent().
			LineD("elif $value is int or $value is float:").
			Indent().
			LineD("$field = bool($value)").
			Unindent().
			WriteLine("else:").
			Indent().
			LineD("push_error(\"Expected bool for field $field in \" +").
			Indent().
			LineD("\"$model, but got \", $value)").
			Unindent().
			WriteLine("return null").
			Unindent()
	case cclValues.TypeNameBytes:
		builder.LineD("if $value is String:").
			Indent().
			LineD("$field = Marshalls.base64_to_raw($value)").
			Unindent().
			LineD("elif $value is PackedByteArray:").
			Indent().
			LineD("$field = $value").
			Unindent().
			WriteLine("else:").
			Indent().
			LineD("push_error(\"Expected base64 string for field $field in \" +").
			Indent().
			LineD("\"$model, but got \", $value)").
			Unindent().
			WriteLine("return null").
			Unindent()
	default:
		if field.IsCustomTypeModel() {
			builder.LineD("if $value is Dictionary:").
				Indent().
				LineD("$field = $fieldT.deserialize_json_dict($value)").
				Unindent().
				WriteLine("else:").
				Indent().
				LineD("push_error(\"Expected base64 string for field $field in \" +").
				Indent().
				LineD("\"$model, but got \", $value)").
				Unindent().
				WriteLine("return null").
				Unindent()
		}
	}
}

func (c *GDScriptGenerationContext) generateArrayDeserializeJson(
	resultName string,
	modelName string,
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldRawName := cclUtils.ToSnakeCase(field.Name)
	resultField := resultName + "." + fieldRawName
	valueName := fieldRawName + "_value"
	listName := fieldRawName + "_list"

	builder.WriteLine("if data.has(\"" + jsonName + "\"):")
	builder.Indent()
	builder.WriteLine("var " + valueName + " = data[\"" + jsonName + "\"]")
	builder.WriteLine("if " + valueName + " != null:")
	builder.Indent()
	builder.WriteLine("if typeof(" + valueName + ") != TYPE_ARRAY:")
	builder.Indent()
	builder.WriteLine("push_error(\"Expected array for field " + field.Name + " in " + modelName + "\")")
	builder.WriteLine("return null")
	builder.Unindent()
	builder.WriteLine("var " + listName + " = [] as " + c.getGDScriptType(field))
	builder.WriteLine("for item in " + valueName + ":")
	builder.Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("if typeof(item) != TYPE_STRING:")
		builder.Indent()
		builder.WriteLine("push_error(\"Expected string items for field " + field.Name + " in " + modelName + "\")")
		builder.WriteLine("return null")
		builder.Unindent()
		builder.WriteLine(listName + ".append(item)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt8, cclValues.TypeNameInt16,
		cclValues.TypeNameInt32, cclValues.TypeNameInt64,
		cclValues.TypeNameUint, cclValues.TypeNameUint8, cclValues.TypeNameUint16,
		cclValues.TypeNameUint32, cclValues.TypeNameUint64, cclValues.TypeNameDateTime:
		builder.WriteLine("if typeof(item) != TYPE_INT and typeof(item) != TYPE_FLOAT:")
		builder.Indent()
		builder.WriteLine("push_error(\"Expected number items for field " + field.Name + " in " + modelName + "\")")
		builder.WriteLine("return null")
		builder.Unindent()
		builder.WriteLine(listName + ".append(int(item))")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteLine("if typeof(item) != TYPE_INT and typeof(item) != TYPE_FLOAT:")
		builder.Indent()
		builder.WriteLine("push_error(\"Expected number items for field " + field.Name + " in " + modelName + "\")")
		builder.WriteLine("return null")
		builder.Unindent()
		builder.WriteLine(listName + ".append(float(item))")
	case cclValues.TypeNameBool:
		builder.WriteLine("if typeof(item) != TYPE_BOOL:")
		builder.Indent()
		builder.WriteLine("push_error(\"Expected bool items for field " + field.Name + " in " + modelName + "\")")
		builder.WriteLine("return null")
		builder.Unindent()
		builder.WriteLine(listName + ".append(item)")
	case cclValues.TypeNameBytes:
		builder.WriteLine("if typeof(item) == TYPE_STRING:")
		builder.Indent()
		builder.WriteLine(listName + ".append(Marshalls.base64_to_raw(item))")
		builder.Unindent()
		builder.WriteLine("elif typeof(item) == TYPE_PACKED_BYTE_ARRAY:")
		builder.Indent()
		builder.WriteLine(listName + ".append(item)")
		builder.Unindent()
		builder.WriteLine("else:")
		builder.Indent()
		builder.WriteLine("push_error(\"Expected base64 string items for field " + field.Name + " in " + modelName + "\")")
		builder.WriteLine("return null")
		builder.Unindent()
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("if item == null:")
			builder.Indent()
			builder.WriteLine(listName + ".append(null)")
			builder.Unindent()
			builder.WriteLine("elif typeof(item) != TYPE_DICTIONARY:")
			builder.Indent()
			builder.WriteLine("push_error(\"Expected object items for field " + field.Name + " in " + modelName + "\")")
			builder.WriteLine("return null")
			builder.Unindent()
			builder.WriteLine("else:")
			builder.Indent()
			builder.WriteLine(listName + ".append(" + targetFieldType.GetName() + ".deserialize_json_dict(item))")
			builder.Unindent()
		}
	}

	builder.Unindent()
	builder.WriteLine(resultField + " = " + listName)
	builder.Unindent()
	builder.Unindent()
	builder.NewLine()
}
