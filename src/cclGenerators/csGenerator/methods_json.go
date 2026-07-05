package csGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
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
	builder.WriteLine("public JsonObject SerializeJsonObject()").
		WriteLine("{").
		Indent().
		WriteLine("var data = new JsonObject();").
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

	builder.WriteLine("return data;").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateSerializeJsonMethod(builder *codeBuilder.CodeBuilder) {
	builder.WriteLine("public string SerializeJson()").
		WriteLine("{").
		Indent().
		WriteLine("return SerializeJsonObject().ToJsonString();").
		Unindent().
		WriteLine("}").
		NewLine()
}

func (c *CSharpGenerationContext) generateDeserializeJsonObjectMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.MapVarPairs(
		"model", model.Name,
	)
	defer builder.UnmapVar(
		"model",
	)

	builder.LineD("public static $model DeserializeJsonObject(JsonNode node)").
		WriteLine("{").
		Indent().
		WriteLine("try").
		WriteLine("{").
		Indent().
		WriteLine("if (node == null || node is not JsonObject data) return null;").
		LineD("var result = new $model();").
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

	builder.WriteLine("return result;").
		Unindent().
		WriteLine("}").
		WriteLine("catch").
		WriteLine("{").
		Indent().
		WriteLine("return null;").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateDeserializeJsonMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) {
	builder.MapVarPairs(
		"model", model.Name,
	)
	defer builder.UnmapVar(
		"model",
	)

	builder.LineD("public static $model DeserializeJson(string jsonText)").
		WriteLine("{").
		Indent().
		WriteLine("if (string.IsNullOrEmpty(jsonText)) return null;").
		WriteLine("try").
		WriteLine("{").
		Indent().
		WriteLine("return DeserializeJsonObject(JsonNode.Parse(jsonText));").
		Unindent().
		WriteLine("}").
		WriteLine("catch").
		WriteLine("{").
		Indent().
		WriteLine("return null;").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
}

func (c *CSharpGenerationContext) generateFieldSerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldName := "this." + caseUtils.ToPascalCase(field.Name)
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
		builder.LineD(`data["$jsonName"] = Convert.ToBase64String($field ?? new byte[0]);`)
	default:
		if field.IsCustomTypeModel() {
			builder.LineD(`data["$jsonName"] = $field != null ? $field.SerializeJsonObject() : null;`)
		} else if field.Type.IsCustomTypeEnum() {
			fieldWrite, err := c.csharpBinaryWriteExpression(field.Type, "$field")
			if err != nil {
				return err
			}
			builder.LineD(`data["$jsonName"] = ` + fieldWrite + `;`)
		} else if c.isCSharpJsonPrimitive(field.Type) {
			builder.LineD(`data["$jsonName"] = $field;`)
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

	builder.MapVarPairs(
		"array", arrayName,
		"field", fieldName,
		"jsonName", jsonName,
	)
	defer builder.UnmapVar(
		"array",
		"field",
		"jsonName",
	)

	builder.LineD("var $array = new JsonArray();").
		LineD("if ($field != null)").
		WriteLine("{").
		Indent().
		LineD("foreach (var item in $field)").
		WriteLine("{").
		Indent()

	switch targetType.GetName() {
	case cclValues.TypeNameBytes:
		builder.LineD("$array.Add(Convert.ToBase64String(item ?? new byte[0]));")
	default:
		if targetType.IsCustomTypeModel() {
			builder.LineD("$array.Add(item != null ? item.SerializeJsonObject() : null);")
		} else if targetType.IsCustomTypeEnum() {
			itemWrite, err := c.csharpBinaryWriteExpression(targetType, "item")
			if err != nil {
				return err
			}
			builder.LineD("$array.Add(" + itemWrite + ");")
		} else if c.isCSharpJsonPrimitive(targetType) {
			builder.LineD("$array.Add(item);")
		} else {
			return c.unsupportedCSharpJsonField(field)
		}
	}

	builder.Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		LineD(`data["$jsonName"] = $array;`).
		NewLine()
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
	fieldType, err := c.getCSharpType(field)
	if err != nil {
		return err
	}

	builder.MapVarPairs(
		"field", resultField,
		"fieldType", fieldType,
		"jsonName", jsonName,
		"node", nodeName,
		"type", field.Type.GetName(),
	)
	defer builder.UnmapVar(
		"field",
		"fieldType",
		"jsonName",
		"node",
		"type",
	)

	builder.LineD(`if (data.TryGetPropertyValue("$jsonName", out var $node) && $node != null)`).
		WriteLine("{").
		Indent()

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.LineD("$field = $node.GetValue<string>();")
	case cclValues.TypeNameBool:
		builder.LineD("$field = $node.GetValue<bool>();")
	case cclValues.TypeNameBytes:
		builder.LineD("$field = Convert.FromBase64String($node.GetValue<string>());")
	default:
		if field.Type.IsCustomTypeEnum() {
			baseType := c.getCSharpEnumBaseType(field.Type.GetDefinition().GetEnumDefinition())
			enumType, err := c.getCSharpEnumTypeReference(
				field.Type.GetDefinition().GetEnumDefinition(),
				field.OwnedBy,
			)
			if err != nil {
				return err
			}
			builder.LineD("$field = (" + enumType + ")$node.GetValue<" + baseType + ">();")
		} else if c.isCSharpJsonNumber(field.Type) {
			builder.LineD("$field = $node.GetValue<$fieldType>();")
		} else if field.IsCustomTypeModel() {
			builder.LineD("$field = $type.DeserializeJsonObject($node);").
				LineD("if ($field == null) return null;")
		} else {
			return c.unsupportedCSharpJsonField(field)
		}
	}

	builder.Unindent().
		WriteLine("}").
		NewLine()
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
	arrayType, err := c.getCSharpType(field)
	if err != nil {
		return err
	}
	itemType, err := c.getCSharpArrayItemType(targetType, field.OwnedBy)
	if err != nil {
		return err
	}

	builder.MapVarPairs(
		"array", arrayName,
		"arrayType", arrayType,
		"field", resultField,
		"itemType", itemType,
		"jsonName", jsonName,
		"node", nodeName,
		"type", targetType.GetName(),
	)
	defer builder.UnmapVar(
		"array",
		"arrayType",
		"field",
		"itemType",
		"jsonName",
		"node",
		"type",
	)

	builder.LineD(`if (data.TryGetPropertyValue("$jsonName", out var $node) && $node != null)`).
		WriteLine("{").
		Indent().
		LineD("if ($node is not JsonArray $array) return null;").
		LineD("$field = new $arrayType();").
		LineD("foreach (var item in $array)").
		WriteLine("{").
		Indent().
		WriteLine("if (item == null) return null;")

	switch targetType.GetName() {
	case cclValues.TypeNameString:
		builder.LineD("$field.Add(item.GetValue<string>());")
	case cclValues.TypeNameBool:
		builder.LineD("$field.Add(item.GetValue<bool>());")
	case cclValues.TypeNameBytes:
		builder.LineD("$field.Add(Convert.FromBase64String(item.GetValue<string>()));")
	default:
		if targetType.IsCustomTypeEnum() {
			baseType := c.getCSharpEnumBaseType(targetType.GetDefinition().GetEnumDefinition())
			enumType, err := c.getCSharpEnumTypeReference(
				targetType.GetDefinition().GetEnumDefinition(),
				field.OwnedBy,
			)
			if err != nil {
				return err
			}
			builder.LineD("$field.Add((" + enumType + ")item.GetValue<" + baseType + ">());")
		} else if c.isCSharpJsonNumber(targetType) {
			builder.LineD("$field.Add(item.GetValue<$itemType>());")
		} else if targetType.IsCustomTypeModel() {
			builder.LineD("var itemValue = $type.DeserializeJsonObject(item);").
				WriteLine("if (itemValue == null) return null;").
				LineD("$field.Add(itemValue);")
		} else {
			return c.unsupportedCSharpJsonField(field)
		}
	}

	builder.Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
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

func (c *CSharpGenerationContext) getCSharpArrayItemType(
	targetType *cclValues.CCLTypeUsage,
	currentModel *CCLModel,
) (string, error) {
	field := &CCLField{
		Type:    targetType,
		OwnedBy: currentModel,
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
