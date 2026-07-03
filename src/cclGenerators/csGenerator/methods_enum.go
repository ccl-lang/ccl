package csGenerator

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *CSharpGenerationContext) getCSharpEnumTypeName(enumDef *CCLEnum) string {
	return enumDef.Name
}

func (c *CSharpGenerationContext) getCSharpEnumMemberName(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	prefix, err := c.GetEnumMemberNamePrefix(
		CurrentLanguage,
		enumDef,
		"",
	)
	if err != nil {
		return "", err
	}

	style, err := c.GetEnumMemberNamingStyle(
		CurrentLanguage,
		enumDef,
		gValues.StylePascalCase,
	)
	if err != nil {
		return "", err
	}

	return prefix + style.ApplyStyle(member.Name), nil
}

func (c *CSharpGenerationContext) getCSharpEnumReference(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	memberName, err := c.getCSharpEnumMemberName(enumDef, member)
	if err != nil {
		return "", err
	}

	return c.getCSharpEnumTypeName(enumDef) + "." + memberName, nil
}

func (c *CSharpGenerationContext) getCSharpEnumBaseType(enumDef *CCLEnum) string {
	return c.getCSharpBuiltinTypeName(enumDef.BaseType.GetName())
}

func (c *CSharpGenerationContext) getCSharpBuiltinTypeName(typeName string) string {
	switch typeName {
	case cclValues.TypeNameString:
		return "string"
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		return "int"
	case cclValues.TypeNameInt8:
		return "sbyte"
	case cclValues.TypeNameInt16:
		return "short"
	case cclValues.TypeNameInt64:
		return "long"
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		return "uint"
	case cclValues.TypeNameUint8:
		return "byte"
	case cclValues.TypeNameUint16:
		return "ushort"
	case cclValues.TypeNameUint64:
		return "ulong"
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		return "float"
	case cclValues.TypeNameFloat64:
		return "double"
	case cclValues.TypeNameBool:
		return "bool"
	case cclValues.TypeNameBytes:
		return "byte[]"
	case cclValues.TypeNameDateTime:
		return "long"
	default:
		return ""
	}
}

func (c *CSharpGenerationContext) csharpBinaryWriteExpression(
	typeUsage *cclValues.CCLTypeUsage,
	expression string,
) string {
	if typeUsage.IsCustomTypeEnum() {
		return "(" + c.getCSharpEnumBaseType(typeUsage.GetDefinition().GetEnumDefinition()) +
			")" + expression
	}

	return expression
}

func (c *CSharpGenerationContext) csharpReaderMethod(typeName string) string {
	switch typeName {
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		return "ReadInt32"
	case cclValues.TypeNameInt8:
		return "ReadSByte"
	case cclValues.TypeNameInt16:
		return "ReadInt16"
	case cclValues.TypeNameInt64:
		return "ReadInt64"
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		return "ReadUInt32"
	case cclValues.TypeNameUint8:
		return "ReadByte"
	case cclValues.TypeNameUint16:
		return "ReadUInt16"
	case cclValues.TypeNameUint64:
		return "ReadUInt64"
	default:
		return ""
	}
}
