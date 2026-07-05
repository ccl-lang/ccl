package rsGenerator

import (
	"github.com/ALiwoto/ssg/ssg"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *RustGenerationContext) generateEnum(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	enumTypeName, err := c.getRustEnumTypeName(enumDef)
	if err != nil {
		return err
	}
	baseType := c.getRustEnumBaseType(enumDef)

	builder.WriteLine("#[repr(" + baseType + ")]").
		WriteLine("#[derive(Debug, Clone, Copy, PartialEq, Eq)]").
		WriteLine("pub enum " + enumTypeName + " {").
		Indent()
	for _, member := range enumDef.Members {
		memberName, err := c.getRustEnumMemberName(enumDef, member)
		if err != nil {
			return err
		}
		builder.WriteLine(memberName + " = " + ssg.ToBase10(member.Value) + ",")
	}
	builder.Unindent().
		WriteLine("}").
		NewLine()

	c.generateEnumConversions(builder, enumDef, enumTypeName, baseType)
	return nil
}

func (c *RustGenerationContext) generateEnumConversions(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
	enumTypeName string,
	baseType string,
) {
	builder.WriteLine("impl " + enumTypeName + " {").
		Indent().
		WriteLine("pub fn from_raw(value: " + baseType + ") -> Option<Self> {").
		Indent().
		WriteLine("match value {").
		Indent()
	for _, member := range enumDef.Members {
		memberName, _ := c.getRustEnumMemberName(enumDef, member)
		builder.WriteLine(ssg.ToBase10(member.Value) + " => Some(Self::" + memberName + "),")
	}
	builder.WriteLine("_ => None,").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()

	builder.WriteLine("impl serde::Serialize for " + enumTypeName + " {").
		Indent().
		WriteLine("fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>").
		WriteLine("where").
		Indent().
		WriteLine("S: serde::Serializer,").
		Unindent().
		WriteLine("{").
		Indent().
		WriteLine("serializer.serialize_" + rustSerdeIntegerMethod(baseType) + "(*self as " + baseType + ")").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()

	builder.WriteLine("impl<'de> serde::Deserialize<'de> for " + enumTypeName + " {").
		Indent().
		WriteLine("fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>").
		WriteLine("where").
		Indent().
		WriteLine("D: serde::Deserializer<'de>,").
		Unindent().
		WriteLine("{").
		Indent().
		WriteLine("let value = <" + baseType + " as serde::Deserialize>::deserialize(deserializer)?;").
		WriteLine("Self::from_raw(value).ok_or_else(|| serde::de::Error::custom(\"invalid enum value\"))").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}")
}

func (c *RustGenerationContext) getRustEnumTypeName(enumDef *CCLEnum) (string, error) {
	defaultPrefix := ""
	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		defaultPrefix = enumDef.OwnedBy.Name
	}

	prefix, err := c.GetEnumTypeNamePrefix(CurrentLanguage, enumDef, defaultPrefix)
	if err != nil {
		return "", err
	}
	return prefix + enumDef.Name, nil
}

func (c *RustGenerationContext) getRustEnumTypeReference(
	enumDef *CCLEnum,
	currentModel *CCLModel,
) (string, error) {
	return c.getRustEnumTypeName(enumDef)
}

func (c *RustGenerationContext) getRustEnumMemberName(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	prefix, err := c.GetEnumMemberNamePrefix(CurrentLanguage, enumDef, "")
	if err != nil {
		return "", err
	}
	style, err := c.GetEnumMemberNamingStyle(CurrentLanguage, enumDef, gValues.StylePascalCase)
	if err != nil {
		return "", err
	}
	return prefix + style.ApplyStyle(member.Name), nil
}

func (c *RustGenerationContext) getRustEnumBaseType(enumDef *CCLEnum) string {
	if mappedType, ok := CCLTypesToRustTypes[enumDef.BaseType.GetName()]; ok {
		return mappedType
	}
	return "i32"
}
