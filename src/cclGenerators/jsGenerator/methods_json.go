package jsGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *JavaScriptGenerationContext) generateSerializeJsonMethods(
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

func (c *JavaScriptGenerationContext) generateSerializeJsonObjectMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("serializeJsonObject() {").
		Indent().
		WriteLine("const data = {};").
		NewLine()

	for _, field := range model.Fields {
		jsonName, err := c.GetJsonFieldName(LanguageName, model, field)
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

func (c *JavaScriptGenerationContext) generateSerializeJsonMethod(
	builder *codeBuilder.CodeBuilder,
) {
	builder.WriteLine("serializeJson() {").
		Indent().
		WriteLine("return JSON.stringify(this.serializeJsonObject());").
		Unindent().
		WriteLine("}").
		NewLine()
}

func (c *JavaScriptGenerationContext) generateDeserializeJsonObjectMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.MapVarPairs(
		"model", model.Name,
	)
	defer builder.UnmapVar("model")

	builder.WriteLine("static deserializeJsonObject(data) {").
		Indent().
		WriteLine("if (data === null || typeof data !== \"object\" || Array.isArray(data)) return null;").
		LineD("const result = new $model();").
		NewLine()

	for _, field := range model.Fields {
		jsonName, err := c.GetJsonFieldName(LanguageName, model, field)
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
		NewLine()

	return nil
}

func (c *JavaScriptGenerationContext) generateDeserializeJsonMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) {
	builder.MapVarPairs(
		"model", model.Name,
	)
	defer builder.UnmapVar("model")

	builder.WriteLine("static deserializeJson(jsonText) {").
		Indent().
		WriteLine("if (!jsonText) return null;").
		WriteLine("try {").
		Indent().
		LineD("return $model.deserializeJsonObject(JSON.parse(jsonText));").
		Unindent().
		WriteLine("} catch (_) {").
		Indent().
		WriteLine("return null;").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
}

func (c *JavaScriptGenerationContext) generateFieldSerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldName := "this." + caseUtils.ToCamelCase(field.Name)
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
		c.writeJavaScriptJsonBytesSerialize(builder, fieldName, `data["`+jsonName+`"]`, false)
	default:
		if field.IsCustomTypeModel() {
			builder.LineD(`data["$jsonName"] = $field ? $field.serializeJsonObject() : null;`)
		} else if c.isJavaScriptJsonPrimitive(field.Type) {
			builder.LineD(`data["$jsonName"] = $field;`)
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       field.Type.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: LanguageName.String(),
			}
		}
	}

	builder.NewLine()
	return nil
}

func (c *JavaScriptGenerationContext) generateArraySerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldName := "this." + caseUtils.ToCamelCase(field.Name)
	listName := caseUtils.ToCamelCase(field.Name) + "JsonList"
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

	builder.LineD("const $list = [];").
		LineD("if ($field) {").
		Indent().
		LineD("for (const item of $field) {").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameBytes:
		c.writeJavaScriptJsonBytesSerialize(builder, "item", listName+"Value", true)
		builder.MapVarPairs("listValue", listName+"Value")
		builder.LineD("$list.push($listValue);")
		builder.UnmapVar("listValue")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.LineD("$list.push(item ? item.serializeJsonObject() : null);")
		} else if c.isJavaScriptJsonPrimitive(targetFieldType) {
			builder.LineD("$list.push(item);")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldType.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: LanguageName.String(),
			}
		}
	}

	builder.Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		LineD(`data["$jsonName"] = $list;`).
		NewLine()

	return nil
}

func (c *JavaScriptGenerationContext) generateFieldDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldName := caseUtils.ToCamelCase(field.GetName())
	resultField := "result." + fieldName
	valueName := fieldName + "Value"
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

	builder.LineD(`const $value = data["$jsonName"];`).
		LineD("if ($value !== undefined && $value !== null) {").
		Indent()

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.LineD("$field = String($value);")
	case cclValues.TypeNameBool:
		builder.LineD(`if (typeof $value === "boolean") {`).
			Indent().
			LineD("$field = $value;").
			Unindent().
			LineD(`} else if (typeof $value === "string") {`).
			Indent().
			LineD(`$field = $value !== "" && $value !== "0" && $value.toLowerCase() !== "false";`).
			Unindent().
			WriteLine("} else {").
			Indent().
			LineD("$field = Boolean($value);").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameBytes:
		c.writeJavaScriptJsonBytesDeserialize(builder, valueName, "bytesValue")
		builder.LineD("$field = bytesValue;")
	default:
		if c.isJavaScriptJsonNumber(field.Type) {
			builder.LineD("const numberValue = Number($value);").
				WriteLine("if (!Number.isFinite(numberValue)) return null;").
				LineD("$field = numberValue;")
		} else if field.IsCustomTypeModel() {
			builder.LineD("$field = $type.deserializeJsonObject($value);").
				LineD("if ($field === null) return null;")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       field.Type.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: LanguageName.String(),
			}
		}
	}

	builder.Unindent().
		WriteLine("}").
		NewLine()

	return nil
}

func (c *JavaScriptGenerationContext) generateArrayDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldName := caseUtils.ToCamelCase(field.Name)
	resultField := "result." + fieldName
	valueName := fieldName + "Value"
	itemName := fieldName + "Item"
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

	builder.LineD(`const $value = data["$jsonName"];`).
		LineD("if ($value === undefined || $value === null) {").
		Indent().
		LineD("$field = [];").
		Unindent().
		WriteLine("} else {").
		Indent().
		LineD("if (!Array.isArray($value)) return null;").
		LineD("$field = [];").
		LineD("for (const item of $value) {").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.LineD("$field.push(String(item));")
	case cclValues.TypeNameBool:
		builder.LineD("let $item;").
			WriteLine("if (typeof item === \"boolean\") {").
			Indent().
			LineD("$item = item;").
			Unindent().
			WriteLine("} else if (typeof item === \"string\") {").
			Indent().
			LineD(`$item = item !== "" && item !== "0" && item.toLowerCase() !== "false";`).
			Unindent().
			WriteLine("} else {").
			Indent().
			LineD("$item = Boolean(item);").
			Unindent().
			WriteLine("}").
			LineD("$field.push($item);")
	case cclValues.TypeNameBytes:
		c.writeJavaScriptJsonBytesDeserialize(builder, "item", itemName)
		builder.LineD("$field.push($item);")
	default:
		if c.isJavaScriptJsonNumber(targetFieldType) {
			builder.LineD("const $item = Number(item);").
				LineD("if (!Number.isFinite($item)) return null;").
				LineD("$field.push($item);")
		} else if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("if (item === null) {").
				Indent().
				LineD("$field.push(null);").
				Unindent().
				WriteLine("} else {").
				Indent().
				LineD("const $item = $type.deserializeJsonObject(item);").
				LineD("if ($item === null) return null;").
				LineD("$field.push($item);").
				Unindent().
				WriteLine("}")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldType.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: LanguageName.String(),
			}
		}
	}

	builder.Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()

	return nil
}

func (c *JavaScriptGenerationContext) isJavaScriptJsonPrimitive(targetType *cclValues.CCLTypeUsage) bool {
	return c.isJavaScriptJsonNumber(targetType) ||
		targetType.GetName() == cclValues.TypeNameString ||
		targetType.GetName() == cclValues.TypeNameBool
}

func (c *JavaScriptGenerationContext) isJavaScriptJsonNumber(targetType *cclValues.CCLTypeUsage) bool {
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

func (c *JavaScriptGenerationContext) writeJavaScriptJsonBytesSerialize(
	builder *codeBuilder.CodeBuilder,
	valueName string,
	targetName string,
	declareTarget bool,
) {
	builder.MapVarPairs(
		"value", valueName,
		"target", targetName,
	)
	defer builder.UnmapVar("value", "target")

	if declareTarget {
		builder.LineD("let $target = \"\";")
	} else {
		builder.LineD("$target = \"\";")
	}
	builder.LineD("if ($value) {").
		Indent().
		WriteLine("if (typeof Buffer !== \"undefined\") {").
		Indent().
		LineD("$target = Buffer.from($value).toString(\"base64\");").
		Unindent().
		WriteLine("} else {").
		Indent().
		WriteLine("let binary = \"\";").
		LineD("for (const byte of $value) binary += String.fromCharCode(byte);").
		LineD("$target = btoa(binary);").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}")
}

func (c *JavaScriptGenerationContext) writeJavaScriptJsonBytesDeserialize(
	builder *codeBuilder.CodeBuilder,
	valueName string,
	targetName string,
) {
	builder.MapVarPairs(
		"value", valueName,
		"target", targetName,
	)
	defer builder.UnmapVar("value", "target")

	builder.LineD("if (typeof $value !== \"string\") return null;").
		LineD("let $target;").
		WriteLine("if (typeof Buffer !== \"undefined\") {").
		Indent().
		LineD("$target = new Uint8Array(Buffer.from($value, \"base64\"));").
		Unindent().
		WriteLine("} else {").
		Indent().
		LineD("const binary = atob($value);").
		LineD("$target = new Uint8Array(binary.length);").
		LineD("for (let i = 0; i < binary.length; i++) $target[i] = binary.charCodeAt(i);").
		Unindent().
		WriteLine("}")
}
