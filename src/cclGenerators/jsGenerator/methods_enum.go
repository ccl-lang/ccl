package jsGenerator

import (
	"github.com/ALiwoto/ssg/ssg"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *JavaScriptGenerationContext) getJavaScriptEnumMemberName(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	prefix, err := c.GetEnumMemberNamePrefix(
		LanguageName,
		enumDef,
		"",
	)
	if err != nil {
		return "", err
	}

	style, err := c.GetEnumMemberNamingStyle(
		LanguageName,
		enumDef,
		gValues.StyleUpperCase,
	)
	if err != nil {
		return "", err
	}

	return prefix + style.ApplyStyle(member.Name), nil
}

func (c *JavaScriptGenerationContext) getJavaScriptEnumTypeName(enumDef *CCLEnum) (string, error) {
	prefix, err := c.GetEnumTypeNamePrefix(
		LanguageName,
		enumDef,
		"",
	)
	if err != nil {
		return "", err
	}

	return prefix + enumDef.Name, nil
}

func (c *JavaScriptGenerationContext) getJavaScriptEnumReference(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	memberName, err := c.getJavaScriptEnumMemberName(enumDef, member)
	if err != nil {
		return "", err
	}
	enumTypeName, err := c.getJavaScriptEnumTypeName(enumDef)
	if err != nil {
		return "", err
	}

	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		return enumDef.OwnedBy.Name + "." + enumTypeName + "." + memberName, nil
	}

	return enumTypeName + "." + memberName, nil
}

func (c *JavaScriptGenerationContext) generateEnumObjectDeclaration(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
	prefix string,
) error {
	enumTypeName, err := c.getJavaScriptEnumTypeName(enumDef)
	if err != nil {
		return err
	}

	builder.WriteLine(prefix + enumTypeName + " = Object.freeze({").
		Indent()
	for _, member := range enumDef.Members {
		memberName, err := c.getJavaScriptEnumMemberName(enumDef, member)
		if err != nil {
			return err
		}
		builder.WriteLine(memberName + ": " + ssg.ToBase10(member.Value) + ",")
	}
	builder.Unindent().
		WriteLine("});")
	return nil
}
