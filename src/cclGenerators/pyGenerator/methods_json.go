package pyGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
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
	builder.WriteLine("@staticmethod").
		WriteLine("def deserialize_json_dict(data: dict) -> \"" + model.Name + "\" | None:").
		Indent().
		WriteLine("if not isinstance(data, dict):").
		Indent().
		WriteLine("return None").
		UnindentLine().
		WriteLine("model_result = " + model.Name + "()").
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
	builder.WriteLine("@staticmethod").
		WriteLine("def deserialize_json(json_text: str) -> \"" + model.Name + "\" | None:").
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
		WriteLine("return " + model.Name + ".deserialize_json_dict(data)").
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

	switch field.Type.GetName() {
	case cclValues.TypeNameBytes:
		builder.WriteLine("data[\"" + jsonName + "\"] = base64.b64encode(" + fieldName + ").decode(\"ascii\")")
	default:
		if field.IsCustomTypeModel() {
			builder.WriteLine("data[\"" + jsonName + "\"] = " + fieldName + ".serialize_json_dict() if " + fieldName + " is not None else None")
		} else if c.isPythonJsonPrimitive(field.Type) {
			builder.WriteLine("data[\"" + jsonName + "\"] = " + fieldName)
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

	builder.WriteLine(listName + " = []").
		WriteLine("if " + fieldName + " is not None:").
		Indent().
		WriteLine("for item in " + fieldName + ":").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameBytes:
		builder.WriteLine(listName + ".append(base64.b64encode(item).decode(\"ascii\"))")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine(listName + ".append(item.serialize_json_dict() if item is not None else None)")
		} else if c.isPythonJsonPrimitive(targetFieldType) {
			builder.WriteLine(listName + ".append(item)")
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
		WriteLine("data[\"" + jsonName + "\"] = " + listName).
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

	builder.WriteLine(valueName + " = data.get(\"" + jsonName + "\", None)").
		WriteLine("if " + valueName + " is not None:").
		Indent()

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine(resultField + " = str(" + valueName + ")")
	case cclValues.TypeNameBool:
		builder.WriteLine("if isinstance(" + valueName + ", bool):").
			Indent().
			WriteLine(resultField + " = " + valueName).
			Unindent().
			WriteLine("elif isinstance(" + valueName + ", str):").
			Indent().
			WriteLine(resultField + " = " + valueName + ".lower() not in (\"\", \"0\", \"false\")").
			Unindent().
			WriteLine("else:").
			Indent().
			WriteLine(resultField + " = bool(" + valueName + ")").
			Unindent()
	case cclValues.TypeNameBytes:
		builder.WriteLine("if not isinstance(" + valueName + ", str):").
			Indent().
			WriteLine("return None").
			Unindent().
			WriteLine("try:").
			Indent().
			WriteLine(resultField + " = base64.b64decode(" + valueName + ")").
			Unindent().
			WriteLine("except ValueError:").
			Indent().
			WriteLine("return None").
			Unindent()
	default:
		if c.isPythonJsonInteger(field.Type) {
			builder.WriteLine("try:").
				Indent().
				WriteLine(resultField + " = int(" + valueName + ")").
				Unindent().
				WriteLine("except (TypeError, ValueError):").
				Indent().
				WriteLine("return None").
				Unindent()
		} else if c.isPythonJsonFloat(field.Type) {
			builder.WriteLine("try:").
				Indent().
				WriteLine(resultField + " = float(" + valueName + ")").
				Unindent().
				WriteLine("except (TypeError, ValueError):").
				Indent().
				WriteLine("return None").
				Unindent()
		} else if field.IsCustomTypeModel() {
			builder.WriteLine(resultField + " = " + field.Type.GetName() + ".deserialize_json_dict(" + valueName + ")").
				WriteLine("if " + resultField + " is None:").
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

	builder.WriteLine(valueName + " = data.get(\"" + jsonName + "\", None)").
		WriteLine("if " + valueName + " is None:").
		Indent().
		WriteLine(resultField + " = []").
		Unindent().
		WriteLine("elif not isinstance(" + valueName + ", list):").
		Indent().
		WriteLine("return None").
		Unindent().
		WriteLine("else:").
		Indent().
		WriteLine(resultField + " = []").
		WriteLine("for item in " + valueName + ":").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine(resultField + ".append(str(item))")
	case cclValues.TypeNameBool:
		builder.WriteLine("if isinstance(item, bool):").
			Indent().
			WriteLine(itemName + " = item").
			Unindent().
			WriteLine("elif isinstance(item, str):").
			Indent().
			WriteLine(itemName + " = item.lower() not in (\"\", \"0\", \"false\")").
			Unindent().
			WriteLine("else:").
			Indent().
			WriteLine(itemName + " = bool(item)").
			Unindent().
			WriteLine(resultField + ".append(" + itemName + ")")
	case cclValues.TypeNameBytes:
		builder.WriteLine("if not isinstance(item, str):").
			Indent().
			WriteLine("return None").
			Unindent().
			WriteLine("try:").
			Indent().
			WriteLine(itemName + " = base64.b64decode(item)").
			Unindent().
			WriteLine("except ValueError:").
			Indent().
			WriteLine("return None").
			Unindent().
			WriteLine(resultField + ".append(" + itemName + ")")
	default:
		if c.isPythonJsonInteger(targetFieldType) {
			builder.WriteLine("try:").
				Indent().
				WriteLine(itemName + " = int(item)").
				Unindent().
				WriteLine("except (TypeError, ValueError):").
				Indent().
				WriteLine("return None").
				Unindent().
				WriteLine(resultField + ".append(" + itemName + ")")
		} else if c.isPythonJsonFloat(targetFieldType) {
			builder.WriteLine("try:").
				Indent().
				WriteLine(itemName + " = float(item)").
				Unindent().
				WriteLine("except (TypeError, ValueError):").
				Indent().
				WriteLine("return None").
				Unindent().
				WriteLine(resultField + ".append(" + itemName + ")")
		} else if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("if item is None:").
				Indent().
				WriteLine(resultField + ".append(None)").
				Unindent().
				WriteLine("else:").
				Indent().
				WriteLine(itemName + " = " + targetFieldType.GetName() + ".deserialize_json_dict(item)").
				WriteLine("if " + itemName + " is None:").
				Indent().
				WriteLine("return None").
				Unindent().
				WriteLine(resultField + ".append(" + itemName + ")").
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
	switch targetType.GetName() {
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
