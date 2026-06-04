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
	builder.WriteLine("public static deserializeJsonObject(data: unknown): " + model.Name + " | null {").
		Indent().
		WriteLine("if (data === null || typeof data !== \"object\" || Array.isArray(data)) return null;").
		WriteLine("const source = data as Record<string, unknown>;").
		WriteLine("const result = new " + model.Name + "();").
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
	builder.WriteLine("public static deserializeJson(jsonText: string): " + model.Name + " | null {").
		Indent().
		WriteLine("if (!jsonText) return null;").
		WriteLine("try {").
		Indent().
		WriteLine("return " + model.Name + ".deserializeJsonObject(JSON.parse(jsonText));").
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

	switch field.Type.GetName() {
	case cclValues.TypeNameBytes:
		builder.WriteLine("data[\"" + jsonName + "\"] = " + helperOwner + "._bytesToBase64(" + fieldName + ");")
	default:
		if field.IsCustomTypeModel() {
			builder.WriteLine("data[\"" + jsonName + "\"] = " + fieldName + " ? " + fieldName + ".serializeJsonObject() : null;")
		} else if c.isTypeScriptJsonPrimitive(field.Type) {
			builder.WriteLine("data[\"" + jsonName + "\"] = " + fieldName + ";")
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

	builder.WriteLine("const " + listName + ": unknown[] = [];").
		WriteLine("if (" + fieldName + ") {").
		Indent().
		WriteLine("for (const item of " + fieldName + ") {").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameBytes:
		builder.WriteLine(listName + ".push(" + helperOwner + "._bytesToBase64(item));")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine(listName + ".push(item ? item.serializeJsonObject() : null);")
		} else if c.isTypeScriptJsonPrimitive(targetFieldType) {
			builder.WriteLine(listName + ".push(item);")
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
		WriteLine("data[\"" + jsonName + "\"] = " + listName + ";").
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

	builder.WriteLine("const " + valueName + " = source[\"" + jsonName + "\"];").
		WriteLine("if (" + valueName + " !== undefined && " + valueName + " !== null) {").
		Indent()

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine(resultField + " = String(" + valueName + ");")
	case cclValues.TypeNameBool:
		builder.WriteLine("if (typeof " + valueName + " === \"boolean\") {").
			Indent().
			WriteLine(resultField + " = " + valueName + ";").
			UnindentLine().
			WriteLine("} else if (typeof " + valueName + " === \"string\") {").
			Indent().
			WriteLine(resultField + " = " + valueName + " !== \"\" && " + valueName + " !== \"0\" && " + valueName + ".toLowerCase() !== \"false\";").
			UnindentLine().
			WriteLine("} else {").
			Indent().
			WriteLine(resultField + " = Boolean(" + valueName + ");").
			UnindentLine().
			WriteLine("}")
	case cclValues.TypeNameBytes:
		builder.WriteLine("const bytesValue = " + helperOwner + "._base64ToBytes(" + valueName + ");").
			WriteLine("if (bytesValue === null) return null;").
			WriteLine(resultField + " = bytesValue;")
	default:
		if c.isTypeScriptJsonNumber(field.Type) {
			builder.WriteLine("const numberValue = Number(" + valueName + ");").
				WriteLine("if (!Number.isFinite(numberValue)) return null;").
				WriteLine(resultField + " = numberValue;")
		} else if field.IsCustomTypeModel() {
			builder.WriteLine("const nestedValue = " + field.Type.GetName() + ".deserializeJsonObject(" + valueName + ");").
				WriteLine("if (nestedValue === null) return null;").
				WriteLine(resultField + " = nestedValue as any;")
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

	builder.WriteLine("const " + valueName + " = source[\"" + jsonName + "\"];").
		WriteLine("if (" + valueName + " === undefined || " + valueName + " === null) {").
		Indent().
		WriteLine(resultField + " = [];").
		UnindentLine().
		WriteLine("} else {").
		Indent().
		WriteLine("if (!Array.isArray(" + valueName + ")) return null;").
		WriteLine(resultField + " = [];").
		WriteLine("for (const item of " + valueName + ") {").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine(resultField + ".push(String(item));")
	case cclValues.TypeNameBool:
		builder.WriteLine("let " + itemName + ": boolean;").
			WriteLine("if (typeof item === \"boolean\") {").
			Indent().
			WriteLine(itemName + " = item;").
			UnindentLine().
			WriteLine("} else if (typeof item === \"string\") {").
			Indent().
			WriteLine(itemName + " = item !== \"\" && item !== \"0\" && item.toLowerCase() !== \"false\";").
			UnindentLine().
			WriteLine("} else {").
			Indent().
			WriteLine(itemName + " = Boolean(item);").
			UnindentLine().
			WriteLine("}").
			WriteLine(resultField + ".push(" + itemName + ");")
	case cclValues.TypeNameBytes:
		builder.WriteLine("const " + itemName + " = " + helperOwner + "._base64ToBytes(item);").
			WriteLine("if (" + itemName + " === null) return null;").
			WriteLine(resultField + ".push(" + itemName + ");")
	default:
		if c.isTypeScriptJsonNumber(targetFieldType) {
			builder.WriteLine("const " + itemName + " = Number(item);").
				WriteLine("if (!Number.isFinite(" + itemName + ")) return null;").
				WriteLine(resultField + ".push(" + itemName + ");")
		} else if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("if (item === null) {").
				Indent().
				WriteLine(resultField + ".push(null);").
				UnindentLine().
				WriteLine("} else {").
				Indent().
				WriteLine("const " + itemName + " = " + targetFieldType.GetName() + ".deserializeJsonObject(item);").
				WriteLine("if (" + itemName + " === null) return null;").
				WriteLine(resultField + ".push(" + itemName + ");").
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
