package goGenerator

import (
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *GoGenerationContext) ensureJsonHelpers() {
	if c.JsonHelpersGenerated {
		return
	}
	c.JsonHelpersGenerated = true

	c.MethodsCode.WriteLine("func cclReadJSONMap(data string) (map[string]json.RawMessage, error) {").
		Indent().
		WriteLine("data = strings.TrimSpace(data)").
		WriteLine("if data == \"\" || data == \"null\" {").
		Indent().
		WriteLine("return nil, nil").
		Unindent().
		WriteLine("}").
		WriteLine("var result map[string]json.RawMessage").
		WriteLine("if err := json.Unmarshal([]byte(data), &result); err != nil {").
		Indent().
		WriteLine("return nil, err").
		Unindent().
		WriteLine("}").
		WriteLine("return result, nil").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("func cclReadJSONString(data json.RawMessage) (string, error) {").
		Indent().
		WriteLine("var result string").
		WriteLine("if err := json.Unmarshal(data, &result); err != nil {").
		Indent().
		WriteLine("return \"\", err").
		Unindent().
		WriteLine("}").
		WriteLine("return result, nil").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("func cclReadJSONInt(data json.RawMessage) (int64, error) {").
		Indent().
		WriteLine("value := strings.TrimSpace(string(data))").
		WriteLine("if len(value) >= 2 && value[0] == '\"' && value[len(value)-1] == '\"' {").
		Indent().
		WriteLine("unquoted, err := strconv.Unquote(value)").
		WriteLine("if err != nil {").
		Indent().
		WriteLine("return 0, err").
		Unindent().
		WriteLine("}").
		WriteLine("value = unquoted").
		Unindent().
		WriteLine("}").
		WriteLine("return strconv.ParseInt(value, 10, 64)").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("func cclReadJSONUint(data json.RawMessage) (uint64, error) {").
		Indent().
		WriteLine("value := strings.TrimSpace(string(data))").
		WriteLine("if len(value) >= 2 && value[0] == '\"' && value[len(value)-1] == '\"' {").
		Indent().
		WriteLine("unquoted, err := strconv.Unquote(value)").
		WriteLine("if err != nil {").
		Indent().
		WriteLine("return 0, err").
		Unindent().
		WriteLine("}").
		WriteLine("value = unquoted").
		Unindent().
		WriteLine("}").
		WriteLine("return strconv.ParseUint(value, 10, 64)").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("func cclReadJSONFloat(data json.RawMessage) (float64, error) {").
		Indent().
		WriteLine("value := strings.TrimSpace(string(data))").
		WriteLine("if len(value) >= 2 && value[0] == '\"' && value[len(value)-1] == '\"' {").
		Indent().
		WriteLine("unquoted, err := strconv.Unquote(value)").
		WriteLine("if err != nil {").
		Indent().
		WriteLine("return 0, err").
		Unindent().
		WriteLine("}").
		WriteLine("value = unquoted").
		Unindent().
		WriteLine("}").
		WriteLine("return strconv.ParseFloat(value, 64)").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("func cclReadJSONBool(data json.RawMessage) (bool, error) {").
		Indent().
		WriteLine("value := strings.TrimSpace(string(data))").
		WriteLine("if len(value) >= 2 && value[0] == '\"' && value[len(value)-1] == '\"' {").
		Indent().
		WriteLine("unquoted, err := strconv.Unquote(value)").
		WriteLine("if err != nil {").
		Indent().
		WriteLine("return false, err").
		Unindent().
		WriteLine("}").
		WriteLine("value = strings.ToLower(unquoted)").
		Unindent().
		WriteLine("}").
		WriteLine("switch value {").
		Indent().
		WriteLine("case \"true\", \"1\":").
		Indent().
		WriteLine("return true, nil").
		Unindent().
		WriteLine("case \"false\", \"0\", \"\":").
		Indent().
		WriteLine("return false, nil").
		Unindent().
		WriteLine("default:").
		Indent().
		WriteLine("return false, strconv.ErrSyntax").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("func cclReadJSONBytes(data json.RawMessage) ([]byte, error) {").
		Indent().
		WriteLine("value, err := cclReadJSONString(data)").
		WriteLine("if err != nil {").
		Indent().
		WriteLine("return nil, err").
		Unindent().
		WriteLine("}").
		WriteLine("return base64.StdEncoding.DecodeString(value)").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("func cclReadJSONArray(data json.RawMessage) ([]json.RawMessage, error) {").
		Indent().
		WriteLine("var result []json.RawMessage").
		WriteLine("if err := json.Unmarshal(data, &result); err != nil {").
		Indent().
		WriteLine("return nil, err").
		Unindent().
		WriteLine("}").
		WriteLine("return result, nil").
		Unindent().
		WriteLine("}").
		NewLine()
}

func (c *GoGenerationContext) isGoJsonSignedInteger(targetType *cclValues.CCLTypeUsage) bool {
	switch targetType.GetName() {
	case cclValues.TypeNameInt, cclValues.TypeNameInt8, cclValues.TypeNameInt16,
		cclValues.TypeNameInt32, cclValues.TypeNameInt64, cclValues.TypeNameDateTime:
		return true
	default:
		return false
	}
}

func (c *GoGenerationContext) isGoJsonUnsignedInteger(targetType *cclValues.CCLTypeUsage) bool {
	switch targetType.GetName() {
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

func (c *GoGenerationContext) goJsonIntegerCast(targetType *cclValues.CCLTypeUsage) string {
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
