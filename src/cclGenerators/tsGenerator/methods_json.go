package tsGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *TypeScriptGenerationContext) generateSerializeJsonMethods(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	c.generateJsonByteHelpers(builder)
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

func (c *TypeScriptGenerationContext) generateJsonByteHelpers(builder *codeBuilder.CodeBuilder) {
	builder.WriteLine("private static _bytesToBase64(bytes: Uint8Array | null | undefined): string {").
		Indent().
		WriteLine("if (!bytes) return \"\";").
		WriteLine("const bufferClass = (globalThis as any).Buffer;").
		WriteLine("if (bufferClass) return bufferClass.from(bytes).toString(\"base64\");").
		WriteLine("let binary = \"\";").
		WriteLine("for (const byte of bytes) binary += String.fromCharCode(byte);").
		WriteLine("return btoa(binary);").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("private static _base64ToBytes(value: unknown): Uint8Array | null {").
		Indent().
		WriteLine("if (typeof value !== \"string\") return null;").
		WriteLine("const bufferClass = (globalThis as any).Buffer;").
		WriteLine("if (bufferClass) return new Uint8Array(bufferClass.from(value, \"base64\"));").
		WriteLine("const binary = atob(value);").
		WriteLine("const bytes = new Uint8Array(binary.length);").
		WriteLine("for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i);").
		WriteLine("return bytes;").
		Unindent().
		WriteLine("}").
		NewLine()
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
		UnindentLine().
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
		UnindentLine().
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
		UnindentLine().
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
		UnindentLine().
		WriteLine("}").
		NewLine()
}

func (c *TypeScriptGenerationContext) generateFieldSerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldName := "this." + caseUtils.ToCamelCase(field.Name)
	helperOwner := field.OwnedBy.GetName()
	builder.MapVarPairs(
		"field", fieldName,
		"helperOwner", helperOwner,
		"jsonName", jsonName,
	)
	defer builder.UnmapVar(
		"field",
		"helperOwner",
		"jsonName",
	)

	switch field.Type.GetName() {
	case cclValues.TypeNameBytes:
		builder.LineD(`data["$jsonName"] = $helperOwner._bytesToBase64($field);`)
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
	helperOwner := field.OwnedBy.GetName()
	builder.MapVarPairs(
		"field", fieldName,
		"helperOwner", helperOwner,
		"jsonName", jsonName,
		"list", listName,
	)
	defer builder.UnmapVar(
		"field",
		"helperOwner",
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
		builder.LineD("$list.push($helperOwner._bytesToBase64(item));")
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
	helperOwner := field.OwnedBy.GetName()
	builder.MapVarPairs(
		"field", resultField,
		"helperOwner", helperOwner,
		"jsonName", jsonName,
		"type", field.Type.GetName(),
		"value", valueName,
	)
	defer builder.UnmapVar(
		"field",
		"helperOwner",
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
		builder.LineD("const bytesValue = $helperOwner._base64ToBytes($value);").
			WriteLine("if (bytesValue === null) return null;").
			LineD("$field = bytesValue;")
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

	builder.UnindentLine().
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
	helperOwner := field.OwnedBy.GetName()
	builder.MapVarPairs(
		"field", resultField,
		"helperOwner", helperOwner,
		"item", itemName,
		"jsonName", jsonName,
		"type", targetFieldType.GetName(),
		"value", valueName,
	)
	defer builder.UnmapVar(
		"field",
		"helperOwner",
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
		builder.LineD("const $item = $helperOwner._base64ToBytes(item);").
			LineD("if ($item === null) return null;").
			LineD("$field.push($item);")
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
