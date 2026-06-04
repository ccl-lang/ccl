package csGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *CSharpGenerationContext) generateSerializeJsonMethods(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	if err := c.generateSerializeJsonObjectMethod(model, builder); err != nil {
		return err
	}
	c.generateSerializeJsonMethod(builder)
	if err := c.generateDeserializeJsonObjectMethod(model, builder); err != nil {
		return err
	}
	c.generateDeserializeJsonMethod(model, builder)
	return nil
}

func (c *CSharpGenerationContext) generateSerializeJsonObjectMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("public JsonObject SerializeJsonObject()")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("var data = new JsonObject();")
	builder.NewLine()

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

	builder.WriteLine("return data;")
	builder.Unindent()
	builder.WriteLine("}")
	builder.NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateSerializeJsonMethod(builder *codeBuilder.CodeBuilder) {
	builder.WriteLine("public string SerializeJson()")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("return SerializeJsonObject().ToJsonString();")
	builder.Unindent()
	builder.WriteLine("}")
	builder.NewLine()
}

func (c *CSharpGenerationContext) generateDeserializeJsonObjectMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("public static " + model.Name + " DeserializeJsonObject(JsonNode node)")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("try")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("if (node == null || node is not JsonObject data) return null;")
	builder.WriteLine("var result = new " + model.Name + "();")
	builder.NewLine()

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

	builder.WriteLine("return result;")
	builder.Unindent()
	builder.WriteLine("}")
	builder.WriteLine("catch")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("return null;")
	builder.Unindent()
	builder.WriteLine("}")
	builder.Unindent()
	builder.WriteLine("}")
	builder.NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateDeserializeJsonMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) {
	builder.WriteLine("public static " + model.Name + " DeserializeJson(string jsonText)")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("if (string.IsNullOrEmpty(jsonText)) return null;")
	builder.WriteLine("try")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("return DeserializeJsonObject(JsonNode.Parse(jsonText));")
	builder.Unindent()
	builder.WriteLine("}")
	builder.WriteLine("catch")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("return null;")
	builder.Unindent()
	builder.WriteLine("}")
	builder.Unindent()
	builder.WriteLine("}")
	builder.NewLine()
}

func (c *CSharpGenerationContext) generateFieldSerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldName := "this." + caseUtils.ToPascalCase(field.Name)
	switch field.Type.GetName() {
	case cclValues.TypeNameBytes:
		builder.WriteLine("data[\"" + jsonName + "\"] = Convert.ToBase64String(" + fieldName + " ?? new byte[0]);")
	default:
		if field.IsCustomTypeModel() {
			builder.WriteLine("data[\"" + jsonName + "\"] = " + fieldName + " != null ? " + fieldName + ".SerializeJsonObject() : null;")
		} else if c.isCSharpJsonPrimitive(field.Type) {
			builder.WriteLine("data[\"" + jsonName + "\"] = " + fieldName + ";")
		} else {
			return c.unsupportedCSharpJsonField(field)
		}
	}
	builder.NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateArraySerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	targetType := field.Type.GetUnderlyingType()
	fieldName := "this." + caseUtils.ToPascalCase(field.Name)
	arrayName := caseUtils.ToCamelCase(field.Name) + "JsonArray"

	builder.WriteLine("var " + arrayName + " = new JsonArray();")
	builder.WriteLine("if (" + fieldName + " != null)")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("foreach (var item in " + fieldName + ")")
	builder.WriteLine("{")
	builder.Indent()

	switch targetType.GetName() {
	case cclValues.TypeNameBytes:
		builder.WriteLine(arrayName + ".Add(Convert.ToBase64String(item ?? new byte[0]));")
	default:
		if targetType.IsCustomTypeModel() {
			builder.WriteLine(arrayName + ".Add(item != null ? item.SerializeJsonObject() : null);")
		} else if c.isCSharpJsonPrimitive(targetType) {
			builder.WriteLine(arrayName + ".Add(item);")
		} else {
			return c.unsupportedCSharpJsonField(field)
		}
	}

	builder.Unindent()
	builder.WriteLine("}")
	builder.Unindent()
	builder.WriteLine("}")
	builder.WriteLine("data[\"" + jsonName + "\"] = " + arrayName + ";")
	builder.NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateFieldDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldName := caseUtils.ToPascalCase(field.Name)
	resultField := "result." + fieldName
	nodeName := caseUtils.ToCamelCase(field.Name) + "Node"

	builder.WriteLine("if (data.TryGetPropertyValue(\"" + jsonName + "\", out var " + nodeName + ") && " + nodeName + " != null)")
	builder.WriteLine("{")
	builder.Indent()

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine(resultField + " = " + nodeName + ".GetValue<string>();")
	case cclValues.TypeNameBool:
		builder.WriteLine(resultField + " = " + nodeName + ".GetValue<bool>();")
	case cclValues.TypeNameBytes:
		builder.WriteLine(resultField + " = Convert.FromBase64String(" + nodeName + ".GetValue<string>());")
	default:
		if c.isCSharpJsonNumber(field.Type) {
			builder.WriteLine(resultField + " = " + nodeName + ".GetValue<" + c.getCSharpType(field) + ">();")
		} else if field.IsCustomTypeModel() {
			builder.WriteLine(resultField + " = " + field.Type.GetName() + ".DeserializeJsonObject(" + nodeName + ");")
			builder.WriteLine("if (" + resultField + " == null) return null;")
		} else {
			return c.unsupportedCSharpJsonField(field)
		}
	}

	builder.Unindent()
	builder.WriteLine("}")
	builder.NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateArrayDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	targetType := field.Type.GetUnderlyingType()
	fieldName := caseUtils.ToPascalCase(field.Name)
	resultField := "result." + fieldName
	nodeName := caseUtils.ToCamelCase(field.Name) + "Node"
	arrayName := caseUtils.ToCamelCase(field.Name) + "Array"

	builder.WriteLine("if (data.TryGetPropertyValue(\"" + jsonName + "\", out var " + nodeName + ") && " + nodeName + " != null)")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("if (" + nodeName + " is not JsonArray " + arrayName + ") return null;")
	builder.WriteLine(resultField + " = new " + c.getCSharpType(field) + "();")
	builder.WriteLine("foreach (var item in " + arrayName + ")")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("if (item == null) return null;")

	switch targetType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine(resultField + ".Add(item.GetValue<string>());")
	case cclValues.TypeNameBool:
		builder.WriteLine(resultField + ".Add(item.GetValue<bool>());")
	case cclValues.TypeNameBytes:
		builder.WriteLine(resultField + ".Add(Convert.FromBase64String(item.GetValue<string>()));")
	default:
		if c.isCSharpJsonNumber(targetType) {
			fieldType := c.getCSharpArrayItemType(targetType)
			builder.WriteLine(resultField + ".Add(item.GetValue<" + fieldType + ">());")
		} else if targetType.IsCustomTypeModel() {
			builder.WriteLine("var itemValue = " + targetType.GetName() + ".DeserializeJsonObject(item);")
			builder.WriteLine("if (itemValue == null) return null;")
			builder.WriteLine(resultField + ".Add(itemValue);")
		} else {
			return c.unsupportedCSharpJsonField(field)
		}
	}

	builder.Unindent()
	builder.WriteLine("}")
	builder.Unindent()
	builder.WriteLine("}")
	builder.NewLine()
	return nil
}

func (c *CSharpGenerationContext) isCSharpJsonPrimitive(targetType *cclValues.CCLTypeUsage) bool {
	return c.isCSharpJsonNumber(targetType) ||
		targetType.GetName() == cclValues.TypeNameString ||
		targetType.GetName() == cclValues.TypeNameBool
}

func (c *CSharpGenerationContext) isCSharpJsonNumber(targetType *cclValues.CCLTypeUsage) bool {
	switch targetType.GetName() {
	case cclValues.TypeNameInt, cclValues.TypeNameInt8, cclValues.TypeNameInt16,
		cclValues.TypeNameInt32, cclValues.TypeNameInt64,
		cclValues.TypeNameUint, cclValues.TypeNameUint8, cclValues.TypeNameUint16,
		cclValues.TypeNameUint32, cclValues.TypeNameUint64,
		cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64,
		cclValues.TypeNameDateTime:
		return true
	default:
		return false
	}
}

func (c *CSharpGenerationContext) getCSharpArrayItemType(targetType *cclValues.CCLTypeUsage) string {
	field := &CCLField{
		Type: targetType,
	}
	return c.getCSharpType(field)
}

func (c *CSharpGenerationContext) unsupportedCSharpJsonField(field *CCLField) error {
	return &cclErrors.UnsupportedFieldTypeError{
		TypeName:       field.Type.GetName(),
		FieldName:      field.Name,
		ModelName:      field.GetModelFullName(),
		TargetLanguage: CurrentLanguage.String(),
	}
}
