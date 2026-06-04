package jsGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *JavaScriptGenerationContext) generateSerializeJsonMethods(
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

func (c *JavaScriptGenerationContext) generateJsonByteHelpers(builder *codeBuilder.CodeBuilder) {
	builder.WriteLine("static _bytesToBase64(bytes) {").
		Indent().
		WriteLine("if (!bytes) return \"\";").
		WriteLine("if (typeof Buffer !== \"undefined\") return Buffer.from(bytes).toString(\"base64\");").
		WriteLine("let binary = \"\";").
		WriteLine("for (const byte of bytes) binary += String.fromCharCode(byte);").
		WriteLine("return btoa(binary);").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("static _base64ToBytes(value) {").
		Indent().
		WriteLine("if (typeof value !== \"string\") return null;").
		WriteLine("if (typeof Buffer !== \"undefined\") return new Uint8Array(Buffer.from(value, \"base64\"));").
		WriteLine("const binary = atob(value);").
		WriteLine("const bytes = new Uint8Array(binary.length);").
		WriteLine("for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i);").
		WriteLine("return bytes;").
		Unindent().
		WriteLine("}").
		NewLine()
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
	builder.WriteLine("static deserializeJsonObject(data) {").
		Indent().
		WriteLine("if (data === null || typeof data !== \"object\" || Array.isArray(data)) return null;").
		WriteLine("const result = new " + model.Name + "();").
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
	builder.WriteLine("static deserializeJson(jsonText) {").
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

	switch field.Type.GetName() {
	case cclValues.TypeNameBytes:
		builder.WriteLine("data[\"" + jsonName + "\"] = " + field.OwnedBy.GetName() + "._bytesToBase64(" + fieldName + ");")
	default:
		if field.IsCustomTypeModel() {
			builder.WriteLine("data[\"" + jsonName + "\"] = " + fieldName + " ? " + fieldName + ".serializeJsonObject() : null;")
		} else if c.isJavaScriptJsonPrimitive(field.Type) {
			builder.WriteLine("data[\"" + jsonName + "\"] = " + fieldName + ";")
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

	builder.WriteLine("const " + listName + " = [];").
		WriteLine("if (" + fieldName + ") {").
		Indent().
		WriteLine("for (const item of " + fieldName + ") {").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameBytes:
		builder.WriteLine(listName + ".push(" + field.OwnedBy.GetName() + "._bytesToBase64(item));")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine(listName + ".push(item ? item.serializeJsonObject() : null);")
		} else if c.isJavaScriptJsonPrimitive(targetFieldType) {
			builder.WriteLine(listName + ".push(item);")
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
		WriteLine("data[\"" + jsonName + "\"] = " + listName + ";").
		NewLine()

	return nil
}

func (c *JavaScriptGenerationContext) generateFieldDeserializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldName := caseUtils.ToCamelCase(field.Name)
	resultField := "result." + fieldName
	valueName := fieldName + "Value"

	builder.WriteLine("const " + valueName + " = data[\"" + jsonName + "\"];").
		WriteLine("if (" + valueName + " !== undefined && " + valueName + " !== null) {").
		Indent()

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine(resultField + " = String(" + valueName + ");")
	case cclValues.TypeNameBool:
		builder.WriteLine("if (typeof " + valueName + " === \"boolean\") {").
			Indent().
			WriteLine(resultField + " = " + valueName + ";").
			Unindent().
			WriteLine("} else if (typeof " + valueName + " === \"string\") {").
			Indent().
			WriteLine(resultField + " = " + valueName + " !== \"\" && " + valueName + " !== \"0\" && " + valueName + ".toLowerCase() !== \"false\";").
			Unindent().
			WriteLine("} else {").
			Indent().
			WriteLine(resultField + " = Boolean(" + valueName + ");").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameBytes:
		builder.WriteLine("const bytesValue = " + field.OwnedBy.GetName() + "._base64ToBytes(" + valueName + ");").
			WriteLine("if (bytesValue === null) return null;").
			WriteLine(resultField + " = bytesValue;")
	default:
		if c.isJavaScriptJsonNumber(field.Type) {
			builder.WriteLine("const numberValue = Number(" + valueName + ");").
				WriteLine("if (!Number.isFinite(numberValue)) return null;").
				WriteLine(resultField + " = numberValue;")
		} else if field.IsCustomTypeModel() {
			builder.WriteLine(resultField + " = " + field.Type.GetName() + ".deserializeJsonObject(" + valueName + ");").
				WriteLine("if (" + resultField + " === null) return null;")
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

	builder.WriteLine("const " + valueName + " = data[\"" + jsonName + "\"];").
		WriteLine("if (" + valueName + " === undefined || " + valueName + " === null) {").
		Indent().
		WriteLine(resultField + " = [];").
		Unindent().
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
		builder.WriteLine("let " + itemName + ";").
			WriteLine("if (typeof item === \"boolean\") {").
			Indent().
			WriteLine(itemName + " = item;").
			Unindent().
			WriteLine("} else if (typeof item === \"string\") {").
			Indent().
			WriteLine(itemName + " = item !== \"\" && item !== \"0\" && item.toLowerCase() !== \"false\";").
			Unindent().
			WriteLine("} else {").
			Indent().
			WriteLine(itemName + " = Boolean(item);").
			Unindent().
			WriteLine("}").
			WriteLine(resultField + ".push(" + itemName + ");")
	case cclValues.TypeNameBytes:
		builder.WriteLine("const " + itemName + " = " + field.OwnedBy.GetName() + "._base64ToBytes(item);").
			WriteLine("if (" + itemName + " === null) return null;").
			WriteLine(resultField + ".push(" + itemName + ");")
	default:
		if c.isJavaScriptJsonNumber(targetFieldType) {
			builder.WriteLine("const " + itemName + " = Number(item);").
				WriteLine("if (!Number.isFinite(" + itemName + ")) return null;").
				WriteLine(resultField + ".push(" + itemName + ");")
		} else if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("if (item === null) {").
				Indent().
				WriteLine(resultField + ".push(null);").
				Unindent().
				WriteLine("} else {").
				Indent().
				WriteLine("const " + itemName + " = " + targetFieldType.GetName() + ".deserializeJsonObject(item);").
				WriteLine("if (" + itemName + " === null) return null;").
				WriteLine(resultField + ".push(" + itemName + ");").
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
