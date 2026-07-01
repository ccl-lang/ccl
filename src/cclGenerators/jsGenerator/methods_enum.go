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
	style, err := c.GetEnumMemberNamingStyle(
		LanguageName,
		enumDef,
		gValues.StyleUpperCase,
	)
	if err != nil {
		return "", err
	}

	return style.ApplyStyle(member.Name), nil
}

func (c *JavaScriptGenerationContext) getJavaScriptEnumReference(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	memberName, err := c.getJavaScriptEnumMemberName(enumDef, member)
	if err != nil {
		return "", err
	}

	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		return enumDef.OwnedBy.Name + "." + enumDef.Name + "." + memberName, nil
	}

	return enumDef.Name + "." + memberName, nil
}

func (c *JavaScriptGenerationContext) generateEnumObjectDeclaration(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
	prefix string,
) error {
	builder.WriteLine(prefix + enumDef.Name + " = Object.freeze({").
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
