package tsGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *TypeScriptGenerationContext) generateSerializeJsonMethods(
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

func (c *TypeScriptGenerationContext) generateSerializeJsonObjectMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("public serializeJsonObject(): Record<string, unknown> {").
		Indent().
		WriteLine("const data: Record<string, unknown> = {};").
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

func (c *TypeScriptGenerationContext) generateSerializeJsonMethod(
	builder *codeBuilder.CodeBuilder,
) {
	builder.WriteLine("public serializeJson(): string {").
		Indent().
		WriteLine("return JSON.stringify(this.serializeJsonObject());").
		Unindent().
		WriteLine("}").
		NewLine()
}

func (c *TypeScriptGenerationContext) generateDeserializeJsonObjectMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.MapVarPairs(
		"model", model.Name,
	)
	defer builder.UnmapVar("model")

	builder.LineD("public static deserializeJsonObject(data: unknown): $model | null {").
		Indent().
		WriteLine("if (data === null || typeof data !== \"object\" || Array.isArray(data)) return null;").
		WriteLine("const source = data as Record<string, unknown>;").
		LineD("const result = new $model();").
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
		NewLine()

	return nil
}

func (c *TypeScriptGenerationContext) generateDeserializeJsonMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) {
	builder.MapVarPairs(
		"model", model.Name,
	)
	defer builder.UnmapVar("model")

	builder.LineD("public static deserializeJson(jsonText: string): $model | null {").
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
		WriteLine("}")
}

func (c *TypeScriptGenerationContext) generateFieldSerializeJson(
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
		c.writeTypeScriptJsonBytesSerialize(builder, fieldName, `data["`+jsonName+`"]`, false)
	default:
		if field.IsCustomTypeModel() {
			builder.LineD(`data["$jsonName"] = $field ? $field.serializeJsonObject() : null;`)
		} else if c.isTypeScriptJsonPrimitive(field.Type) {
			builder.LineD(`data["$jsonName"] = $field;`)
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

func (c *TypeScriptGenerationContext) generateArraySerializeJson(
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

	builder.LineD("const $list: unknown[] = [];").
		LineD("if ($field) {").
		Indent().
		LineD("for (const item of $field) {").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameBytes:
		c.writeTypeScriptJsonBytesSerialize(builder, "item", listName+"Value", true)
		builder.MapVarPairs("listValue", listName+"Value")
		builder.LineD("$list.push($listValue);")
		builder.UnmapVar("listValue")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.LineD("$list.push(item ? item.serializeJsonObject() : null);")
		} else if c.isTypeScriptJsonPrimitive(targetFieldType) {
			builder.LineD("$list.push(item);")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldType.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}

	builder.UnindentLine().
		WriteLine("}").
		UnindentLine().
		WriteLine("}").
		LineD(`data["$jsonName"] = $list;`).
		NewLine()

	return nil
}

func (c *TypeScriptGenerationContext) generateFieldDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldName := caseUtils.ToCamelCase(field.Name)
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

	builder.LineD(`const $value = source["$jsonName"];`).
		LineD("if ($value !== undefined && $value !== null) {").
		Indent()

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.LineD("$field = String($value);")
	case cclValues.TypeNameBool:
		builder.LineD(`if (typeof $value === "boolean") {`).
			Indent().
			LineD("$field = $value;").
			UnindentLine().
			LineD(`} else if (typeof $value === "string") {`).
			Indent().
			LineD(`$field = $value !== "" && $value !== "0" && $value.toLowerCase() !== "false";`).
			UnindentLine().
			WriteLine("} else {").
			Indent().
			LineD("$field = Boolean($value);").
			UnindentLine().
			WriteLine("}")
	case cclValues.TypeNameBytes:
		c.writeTypeScriptJsonBytesDeserialize(builder, valueName, "bytesValue")
		builder.LineD("$field = bytesValue;")
	default:
		if c.isTypeScriptJsonNumber(field.Type) {
			builder.LineD("const numberValue = Number($value);").
				WriteLine("if (!Number.isFinite(numberValue)) return null;").
				LineD("$field = numberValue;")
		} else if field.IsCustomTypeModel() {
			builder.LineD("const nestedValue = $type.deserializeJsonObject($value);").
				WriteLine("if (nestedValue === null) return null;").
				LineD("$field = nestedValue as any;")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       field.Type.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}

	builder.Unindent().
		WriteLine("}").
		NewLine()

	return nil
}

func (c *TypeScriptGenerationContext) generateArrayDeserializeJson(
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

	builder.LineD(`const $value = source["$jsonName"];`).
		LineD("if ($value === undefined || $value === null) {").
		Indent().
		LineD("$field = [];").
		UnindentLine().
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
		builder.LineD("let $item: boolean;").
			WriteLine("if (typeof item === \"boolean\") {").
			Indent().
			LineD("$item = item;").
			UnindentLine().
			WriteLine("} else if (typeof item === \"string\") {").
			Indent().
			LineD(`$item = item !== "" && item !== "0" && item.toLowerCase() !== "false";`).
			UnindentLine().
			WriteLine("} else {").
			Indent().
			LineD("$item = Boolean(item);").
			UnindentLine().
			WriteLine("}").
			LineD("$field.push($item);")
	case cclValues.TypeNameBytes:
		c.writeTypeScriptJsonBytesDeserialize(builder, "item", itemName)
		builder.LineD("$field.push($item);")
	default:
		if c.isTypeScriptJsonNumber(targetFieldType) {
			builder.LineD("const $item = Number(item);").
				LineD("if (!Number.isFinite($item)) return null;").
				LineD("$field.push($item);")
		} else if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("if (item === null) {").
				Indent().
				LineD("$field.push(null);").
				UnindentLine().
				WriteLine("} else {").
				Indent().
				LineD("const $item = $type.deserializeJsonObject(item);").
				LineD("if ($item === null) return null;").
				LineD("$field.push($item);").
				UnindentLine().
				WriteLine("}")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldType.GetName(),
				FieldName:      field.Name,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}

	builder.UnindentLine().
		WriteLine("}").
		UnindentLine().
		WriteLine("}").
		NewLine()

	return nil
}

func (c *TypeScriptGenerationContext) isTypeScriptJsonPrimitive(targetType *cclValues.CCLTypeUsage) bool {
	return c.isTypeScriptJsonNumber(targetType) ||
		targetType.GetName() == cclValues.TypeNameString ||
		targetType.GetName() == cclValues.TypeNameBool
}

func (c *TypeScriptGenerationContext) isTypeScriptJsonNumber(targetType *cclValues.CCLTypeUsage) bool {
	typeName := targetType.GetName()
	if targetType.IsCustomTypeEnum() {
		typeName = targetType.GetEnumBaseTypeName()
	}

	switch typeName {
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

func (c *TypeScriptGenerationContext) writeTypeScriptJsonBytesSerialize(
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
		builder.LineD(`let $target = "";`)
	} else {
		builder.LineD(`$target = "";`)
	}
	builder.LineD("if ($value) {").
		Indent().
		WriteLine("const bufferClass = (globalThis as any).Buffer;").
		WriteLine("if (bufferClass) {").
		Indent().
		LineD(`$target = bufferClass.from($value).toString("base64");`).
		UnindentLine().
		WriteLine("} else {").
		Indent().
		WriteLine(`let binary = "";`).
		LineD("for (const byte of $value) binary += String.fromCharCode(byte);").
		LineD("$target = btoa(binary);").
		UnindentLine().
		WriteLine("}").
		UnindentLine().
		WriteLine("}")
}

func (c *TypeScriptGenerationContext) writeTypeScriptJsonBytesDeserialize(
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
		WriteLine("const bufferClass = (globalThis as any).Buffer;").
		LineD("let $target: Uint8Array;").
		WriteLine("if (bufferClass) {").
		Indent().
		LineD("$target = new Uint8Array(bufferClass.from($value, \"base64\"));").
		UnindentLine().
		WriteLine("} else {").
		Indent().
		LineD("const binary = atob($value);").
		LineD("$target = new Uint8Array(binary.length);").
		LineD("for (let i = 0; i < binary.length; i++) $target[i] = binary.charCodeAt(i);").
		UnindentLine().
		WriteLine("}")
}
