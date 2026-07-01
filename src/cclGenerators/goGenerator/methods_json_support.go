package goGenerator

import (
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *GoGenerationContext) isGoJsonSignedInteger(targetType *cclValues.CCLTypeUsage) bool {
	typeName := targetType.GetName()
	if targetType.IsCustomTypeEnum() {
		typeName = targetType.GetEnumBaseTypeName()
	}

	switch typeName {
	case cclValues.TypeNameInt, cclValues.TypeNameInt8, cclValues.TypeNameInt16,
		cclValues.TypeNameInt32, cclValues.TypeNameInt64, cclValues.TypeNameDateTime:
		return true
	default:
		return false
	}
}

func (c *GoGenerationContext) isGoJsonUnsignedInteger(targetType *cclValues.CCLTypeUsage) bool {
	typeName := targetType.GetName()
	if targetType.IsCustomTypeEnum() {
		typeName = targetType.GetEnumBaseTypeName()
	}

	switch typeName {
	case cclValues.TypeNameUint, cclValues.TypeNameUint8, cclValues.TypeNameUint16,
		cclValues.TypeNameUint32, cclValues.TypeNameUint64:
		return true
	default:
		return false
	}
}

func (c *GoGenerationContext) isGoJsonFloat(targetType *cclValues.CCLTypeUsage) bool {
	switch targetType.GetName() {
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		return true
	default:
		return false
	}
}

func (c *GoGenerationContext) writeGoJsonStringRead(
	builder *codeBuilder.CodeBuilder,
	valueName string,
	rawName string,
) {
	registerGoImport(builder, "encoding/json")
	builder.MapVarPairs(
		"value", valueName,
		"raw", rawName,
	)
	defer builder.UnmapVar("value", "raw")

	builder.LineD("var $value string").
		LineD("if err := json.Unmarshal($raw, &$value); err != nil {").
		Indent().
		WriteLine("return err").
		Unindent().
		WriteLine("}")
}

func (c *GoGenerationContext) writeGoJsonIntRead(
	builder *codeBuilder.CodeBuilder,
	valueName string,
	rawName string,
) {
	registerGoImport(builder, "strconv")
	registerGoImport(builder, "strings")
	valueTextName := valueName + "Text"
	builder.MapVarPairs(
		"value", valueName,
		"raw", rawName,
		"text", valueTextName,
	)
	defer builder.UnmapVar("value", "raw", "text")

	builder.LineD("$text := strings.TrimSpace(string($raw))")
	c.writeGoJsonQuotedValueUnwrap(builder, false)
	builder.LineD("$value, err := strconv.ParseInt($text, 10, 64)").
		WriteLine("if err != nil {").
		Indent().
		WriteLine("return err").
		Unindent().
		WriteLine("}")
}

func (c *GoGenerationContext) writeGoJsonUintRead(
	builder *codeBuilder.CodeBuilder,
	valueName string,
	rawName string,
) {
	registerGoImport(builder, "strconv")
	registerGoImport(builder, "strings")
	valueTextName := valueName + "Text"
	builder.MapVarPairs(
		"value", valueName,
		"raw", rawName,
		"text", valueTextName,
	)
	defer builder.UnmapVar("value", "raw", "text")

	builder.LineD("$text := strings.TrimSpace(string($raw))")
	c.writeGoJsonQuotedValueUnwrap(builder, false)
	builder.LineD("$value, err := strconv.ParseUint($text, 10, 64)").
		WriteLine("if err != nil {").
		Indent().
		WriteLine("return err").
		Unindent().
		WriteLine("}")
}

func (c *GoGenerationContext) writeGoJsonFloatRead(
	builder *codeBuilder.CodeBuilder,
	valueName string,
	rawName string,
) {
	registerGoImport(builder, "strconv")
	registerGoImport(builder, "strings")
	valueTextName := valueName + "Text"
	builder.MapVarPairs(
		"value", valueName,
		"raw", rawName,
		"text", valueTextName,
	)
	defer builder.UnmapVar("value", "raw", "text")

	builder.LineD("$text := strings.TrimSpace(string($raw))")
	c.writeGoJsonQuotedValueUnwrap(builder, false)
	builder.LineD("$value, err := strconv.ParseFloat($text, 64)").
		WriteLine("if err != nil {").
		Indent().
		WriteLine("return err").
		Unindent().
		WriteLine("}")
}

func (c *GoGenerationContext) writeGoJsonBoolRead(
	builder *codeBuilder.CodeBuilder,
	valueName string,
	rawName string,
) {
	registerGoImport(builder, "strconv")
	registerGoImport(builder, "strings")
	valueTextName := valueName + "Text"
	builder.MapVarPairs(
		"value", valueName,
		"raw", rawName,
		"text", valueTextName,
	)
	defer builder.UnmapVar("value", "raw", "text")

	builder.LineD("$text := strings.TrimSpace(string($raw))")
	c.writeGoJsonQuotedValueUnwrap(builder, true)
	builder.LineD("var $value bool").
		LineD("switch $text {").
		Indent().
		WriteLine("case \"true\", \"1\":").
		Indent().
		LineD("$value = true").
		Unindent().
		WriteLine("case \"false\", \"0\", \"\":").
		Indent().
		LineD("$value = false").
		Unindent().
		WriteLine("default:").
		Indent().
		WriteLine("return strconv.ErrSyntax").
		Unindent().
		Unindent().
		WriteLine("}")
}

func (c *GoGenerationContext) writeGoJsonBytesRead(
	builder *codeBuilder.CodeBuilder,
	valueName string,
	rawName string,
) {
	registerGoImport(builder, "encoding/base64")
	registerGoImport(builder, "encoding/json")
	textName := valueName + "Text"
	builder.MapVarPairs(
		"value", valueName,
		"raw", rawName,
		"text", textName,
	)
	defer builder.UnmapVar("value", "raw", "text")

	builder.LineD("var $text string").
		LineD("if err := json.Unmarshal($raw, &$text); err != nil {").
		Indent().
		WriteLine("return err").
		Unindent().
		WriteLine("}").
		LineD("$value, err := base64.StdEncoding.DecodeString($text)").
		WriteLine("if err != nil {").
		Indent().
		WriteLine("return err").
		Unindent().
		WriteLine("}")
}

func (c *GoGenerationContext) writeGoJsonArrayRead(
	builder *codeBuilder.CodeBuilder,
	valueName string,
	rawName string,
) {
	registerGoImport(builder, "encoding/json")
	builder.MapVarPairs(
		"value", valueName,
		"raw", rawName,
	)
	defer builder.UnmapVar("value", "raw")

	builder.LineD("var $value []json.RawMessage").
		LineD("if err := json.Unmarshal($raw, &$value); err != nil {").
		Indent().
		WriteLine("return err").
		Unindent().
		WriteLine("}")
}

func (c *GoGenerationContext) writeGoJsonQuotedValueUnwrap(
	builder *codeBuilder.CodeBuilder,
	lowercase bool,
) {
	builder.LineD("if len($text) >= 2 && $text[0] == '\"' && $text[len($text)-1] == '\"' {").
		Indent().
		LineD("unquoted, err := strconv.Unquote($text)").
		WriteLine("if err != nil {").
		Indent().
		WriteLine("return err").
		Unindent().
		WriteLine("}")
	if lowercase {
		builder.LineD("$text = strings.ToLower(unquoted)")
	} else {
		builder.LineD("$text = unquoted")
	}
	builder.Unindent().
		WriteLine("}")
}

func (c *GoGenerationContext) goJsonIntegerCast(targetType *cclValues.CCLTypeUsage) string {
	if targetType.IsCustomTypeEnum() {
		return c.getGoEnumTypeName(targetType.GetDefinition().GetEnumDefinition())
	}

	if mappedType, ok := CCLTypesToGoTypes[targetType.GetName()]; ok {
		return mappedType
	}
	return "int64"
}

func (c *GoGenerationContext) goJsonFloatCast(targetType *cclValues.CCLTypeUsage) string {
	if targetType.GetName() == cclValues.TypeNameFloat32 {
		return "float32"
	}
	return "float64"
}

func (c *GoGenerationContext) getGoJsonArrayItemType(targetType *cclValues.CCLTypeUsage) string {
	if targetType.IsCustomTypeModel() {
		return "*" + targetType.GetName()
	}
	if targetType.IsCustomTypeEnum() {
		return c.getGoEnumTypeName(targetType.GetDefinition().GetEnumDefinition())
	}
	if mappedType, ok := CCLTypesToGoTypes[targetType.GetName()]; ok {
		return mappedType
	}
	return ""
}

func (c *GoGenerationContext) unsupportedGoJsonField(field *CCLField) error {
	return &cclErrors.UnsupportedFieldTypeError{
		TypeName:       field.Type.GetName(),
		FieldName:      field.Name,
		ModelName:      field.GetModelFullName(),
		TargetLanguage: CurrentLanguage.String(),
	}
}
