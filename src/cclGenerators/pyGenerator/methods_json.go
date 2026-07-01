package pyGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *PythonGenerationContext) generateSerializeJsonMethods(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.DoImport("base64", "import base64").
		DoImport("json", "import json")

	if err := c.generateSerializeJsonDictMethod(model, builder); err != nil {
		return err
	}
	c.generateSerializeJsonMethod(builder)
	if err := c.generateDeserializeJsonDictMethod(model, builder); err != nil {
		return err
	}
	c.generateDeserializeJsonMethod(model, builder)

	return nil
}

func (c *PythonGenerationContext) generateSerializeJsonDictMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("def serialize_json_dict(self) -> dict:").
		Indent().
		WriteLine("data = {}").
		NewLine()

	for _, field := range model.Fields {
		jsonName, err := c.GetJsonFieldName(CurrentLanguage, model, field)
		if err != nil {
			return err
		}

		if field.IsArray() {
			if err = c.generateArraySerializeJson(field, jsonName, builder); err != nil {
				return err
			}
			continue
		}

		if err = c.generateFieldSerializeJson(field, jsonName, builder); err != nil {
			return err
		}
	}

	builder.WriteLine("return data").
		UnindentLine().
		NewLine()

	return nil
}

func (c *PythonGenerationContext) generateSerializeJsonMethod(
	builder *codeBuilder.CodeBuilder,
) {
	builder.WriteLine("def serialize_json(self) -> str:").
		Indent().
		WriteLine("return json.dumps(self.serialize_json_dict(), separators=(\",\", \":\"))").
		UnindentLine().
		NewLine()
}

func (c *PythonGenerationContext) generateDeserializeJsonDictMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.MapVarPairs(
		"model", model.Name,
	)
	defer builder.UnmapVar("model")

	builder.WriteLine("@staticmethod").
		LineD(`def deserialize_json_dict(data: dict) -> "$model" | None:`).
		Indent().
		WriteLine("if not isinstance(data, dict):").
		Indent().
		WriteLine("return None").
		UnindentLine().
		LineD("model_result = $model()").
		NewLine()

	for _, field := range model.Fields {
		jsonName, err := c.GetJsonFieldName(CurrentLanguage, model, field)
		if err != nil {
			return err
		}

		if field.IsArray() {
			if err = c.generateArrayDeserializeJson(field, jsonName, builder); err != nil {
				return err
			}
			continue
		}

		if err = c.generateFieldDeserializeJson(field, jsonName, builder); err != nil {
			return err
		}
	}

	builder.WriteLine("return model_result").
		UnindentLine().
		NewLine()

	return nil
}

func (c *PythonGenerationContext) generateDeserializeJsonMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) {
	builder.MapVarPairs(
		"model", model.Name,
	)
	defer builder.UnmapVar("model")

	builder.WriteLine("@staticmethod").
		LineD(`def deserialize_json(json_text: str) -> "$model" | None:`).
		Indent().
		WriteLine("if not json_text:").
		Indent().
		WriteLine("return None").
		UnindentLine().
		WriteLine("try:").
		Indent().
		WriteLine("data = json.loads(json_text)").
		Unindent().
		WriteLine("except (TypeError, ValueError):").
		Indent().
		WriteLine("return None").
		UnindentLine().
		LineD("return $model.deserialize_json_dict(data)").
		UnindentLine().
		NewLine()
}

func (c *PythonGenerationContext) generateFieldSerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldRawName := caseUtils.ToSnakeCase(field.Name)
	fieldName := "self." + fieldRawName
	builder.MapVarPairs(
		"field", fieldName,
		"jsonName", jsonName,
	)
	defer builder.UnmapVar(
		"field",
		"jsonName",
	)

	switch field.Type.GetName() {
	case cclValues.TypeNameBytes:
		builder.LineD(`data["$jsonName"] = base64.b64encode($field).decode("ascii")`)
	default:
		if field.IsCustomTypeModel() {
			builder.LineD(`data["$jsonName"] = $field.serialize_json_dict() if $field is not None else None`)
		} else if c.isPythonJsonPrimitive(field.Type) {
			builder.LineD(`data["$jsonName"] = $field`)
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       field.Type.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}

	builder.NewLine()
	return nil
}

func (c *PythonGenerationContext) generateArraySerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldRawName := caseUtils.ToSnakeCase(field.Name)
	fieldName := "self." + fieldRawName
	listName := fieldRawName + "_json_list"
	builder.MapVarPairs(
		"field", fieldName,
		"jsonName", jsonName,
		"list", listName,
	)
	defer builder.UnmapVar(
		"field",
		"jsonName",
		"list",
	)

	builder.LineD("$list = []").
		LineD("if $field is not None:").
		Indent().
		LineD("for item in $field:").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameBytes:
		builder.LineD(`$list.append(base64.b64encode(item).decode("ascii"))`)
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.LineD("$list.append(item.serialize_json_dict() if item is not None else None)")
		} else if c.isPythonJsonPrimitive(targetFieldType) {
			builder.LineD("$list.append(item)")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldType.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}

	builder.Unindent().
		Unindent().
		LineD(`data["$jsonName"] = $list`).
		NewLine()

	return nil
}

func (c *PythonGenerationContext) generateFieldDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldRawName := caseUtils.ToSnakeCase(field.Name)
	resultField := "model_result." + fieldRawName
	valueName := fieldRawName + "_value"
	builder.MapVarPairs(
		"field", resultField,
		"jsonName", jsonName,
		"type", field.Type.GetName(),
		"value", valueName,
	)
	defer builder.UnmapVar(
		"field",
		"jsonName",
		"type",
		"value",
	)

	builder.LineD(`$value = data.get("$jsonName", None)`).
		LineD("if $value is not None:").
		Indent()

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.LineD("$field = str($value)")
	case cclValues.TypeNameBool:
		builder.LineD("if isinstance($value, bool):").
			Indent().
			LineD("$field = $value").
			Unindent().
			LineD("elif isinstance($value, str):").
			Indent().
			LineD(`$field = $value.lower() not in ("", "0", "false")`).
			Unindent().
			WriteLine("else:").
			Indent().
			LineD("$field = bool($value)").
			Unindent()
	case cclValues.TypeNameBytes:
		builder.LineD("if not isinstance($value, str):").
			Indent().
			WriteLine("return None").
			Unindent().
			WriteLine("try:").
			Indent().
			LineD("$field = base64.b64decode($value)").
			Unindent().
			WriteLine("except ValueError:").
			Indent().
			WriteLine("return None").
			Unindent()
	default:
		if field.Type.IsCustomTypeEnum() {
			builder.WriteLine("try:").
				Indent().
				LineD("$field = " + c.pythonEnumCastExpression(field.Type, "int($value)")).
				Unindent().
				WriteLine("except (TypeError, ValueError):").
				Indent().
				WriteLine("return None").
				Unindent()
		} else if c.isPythonJsonInteger(field.Type) {
			builder.WriteLine("try:").
				Indent().
				LineD("$field = int($value)").
				Unindent().
				WriteLine("except (TypeError, ValueError):").
				Indent().
				WriteLine("return None").
				Unindent()
		} else if c.isPythonJsonFloat(field.Type) {
			builder.WriteLine("try:").
				Indent().
				LineD("$field = float($value)").
				Unindent().
				WriteLine("except (TypeError, ValueError):").
				Indent().
				WriteLine("return None").
				Unindent()
		} else if field.IsCustomTypeModel() {
			builder.LineD("$field = $type.deserialize_json_dict($value)").
				LineD("if $field is None:").
				Indent().
				WriteLine("return None").
				Unindent()
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       field.Type.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}

	builder.UnindentLine().
		NewLine()

	return nil
}

func (c *PythonGenerationContext) generateArrayDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldRawName := caseUtils.ToSnakeCase(field.Name)
	resultField := "model_result." + fieldRawName
	valueName := fieldRawName + "_value"
	itemName := fieldRawName + "_item"
	builder.MapVarPairs(
		"field", resultField,
		"item", itemName,
		"jsonName", jsonName,
		"type", targetFieldType.GetName(),
		"value", valueName,
	)
	defer builder.UnmapVar(
		"field",
		"item",
		"jsonName",
		"type",
		"value",
	)

	builder.LineD(`$value = data.get("$jsonName", None)`).
		LineD("if $value is None:").
		Indent().
		LineD("$field = []").
		Unindent().
		LineD("elif not isinstance($value, list):").
		Indent().
		WriteLine("return None").
		Unindent().
		WriteLine("else:").
		Indent().
		LineD("$field = []").
		LineD("for item in $value:").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.LineD("$field.append(str(item))")
	case cclValues.TypeNameBool:
		builder.WriteLine("if isinstance(item, bool):").
			Indent().
			LineD("$item = item").
			Unindent().
			WriteLine("elif isinstance(item, str):").
			Indent().
			LineD(`$item = item.lower() not in ("", "0", "false")`).
			Unindent().
			WriteLine("else:").
			Indent().
			LineD("$item = bool(item)").
			Unindent().
			LineD("$field.append($item)")
	case cclValues.TypeNameBytes:
		builder.WriteLine("if not isinstance(item, str):").
			Indent().
			WriteLine("return None").
			Unindent().
			WriteLine("try:").
			Indent().
			LineD("$item = base64.b64decode(item)").
			Unindent().
			WriteLine("except ValueError:").
			Indent().
			WriteLine("return None").
			Unindent().
			LineD("$field.append($item)")
	default:
		if targetFieldType.IsCustomTypeEnum() {
			builder.WriteLine("try:").
				Indent().
				LineD("$item = " + c.pythonEnumCastExpression(targetFieldType, "int(item)")).
				Unindent().
				WriteLine("except (TypeError, ValueError):").
				Indent().
				WriteLine("return None").
				Unindent().
				LineD("$field.append($item)")
		} else if c.isPythonJsonInteger(targetFieldType) {
			builder.WriteLine("try:").
				Indent().
				LineD("$item = int(item)").
				Unindent().
				WriteLine("except (TypeError, ValueError):").
				Indent().
				WriteLine("return None").
				Unindent().
				LineD("$field.append($item)")
		} else if c.isPythonJsonFloat(targetFieldType) {
			builder.WriteLine("try:").
				Indent().
				LineD("$item = float(item)").
				Unindent().
				WriteLine("except (TypeError, ValueError):").
				Indent().
				WriteLine("return None").
				Unindent().
				LineD("$field.append($item)")
		} else if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("if item is None:").
				Indent().
				LineD("$field.append(None)").
				Unindent().
				WriteLine("else:").
				Indent().
				LineD("$item = $type.deserialize_json_dict(item)").
				LineD("if $item is None:").
				Indent().
				WriteLine("return None").
				Unindent().
				LineD("$field.append($item)").
				Unindent()
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldType.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}

	builder.Unindent().
		UnindentLine().
		NewLine()

	return nil
}

func (c *PythonGenerationContext) isPythonJsonPrimitive(targetType *cclValues.CCLTypeUsage) bool {
	return c.isPythonJsonInteger(targetType) ||
		c.isPythonJsonFloat(targetType) ||
		targetType.GetName() == cclValues.TypeNameString ||
		targetType.GetName() == cclValues.TypeNameBool
}

func (c *PythonGenerationContext) isPythonJsonInteger(targetType *cclValues.CCLTypeUsage) bool {
	switch pythonStorageTypeName(targetType) {
	case cclValues.TypeNameInt, cclValues.TypeNameInt8, cclValues.TypeNameInt16,
		cclValues.TypeNameInt32, cclValues.TypeNameInt64,
		cclValues.TypeNameUint, cclValues.TypeNameUint8, cclValues.TypeNameUint16,
		cclValues.TypeNameUint32, cclValues.TypeNameUint64, cclValues.TypeNameDateTime:
		return true
	default:
		return false
	}
}

func (c *PythonGenerationContext) isPythonJsonFloat(targetType *cclValues.CCLTypeUsage) bool {
	switch targetType.GetName() {
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		return true
	default:
		return false
	}
}
