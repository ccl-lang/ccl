package tsGenerator

import (
	"github.com/ALiwoto/ssg/ssg"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *TypeScriptGenerationContext) getTypeScriptEnumTypeName(enumDef *CCLEnum) string {
	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		return enumDef.OwnedBy.Name + "." + enumDef.Name
	}

	return enumDef.Name
}

func (c *TypeScriptGenerationContext) getTypeScriptEnumMemberName(
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

func (c *TypeScriptGenerationContext) getTypeScriptEnumReference(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	memberName, err := c.getTypeScriptEnumMemberName(enumDef, member)
	if err != nil {
		return "", err
	}

	return c.getTypeScriptEnumTypeName(enumDef) + "." + memberName, nil
}

func (c *TypeScriptGenerationContext) generateEnumDeclaration(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	builder.WriteLine("export enum " + enumDef.Name + " {").
		Indent()
	for _, member := range enumDef.Members {
		memberName, err := c.getTypeScriptEnumMemberName(enumDef, member)
		if err != nil {
			return err
		}
		builder.WriteLine(memberName + " = " + ssg.ToBase10(member.Value) + ",")
	}
	builder.Unindent().
		WriteLine("}")
	return nil
}

func (c *TypeScriptGenerationContext) generateNestedEnumNamespace(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
) error {
	if len(model.Enums) == 0 {
		return nil
	}

	builder.WriteLine("export namespace " + model.Name + " {").
		Indent()
	for enumIndex, enumDef := range model.Enums {
		if enumIndex > 0 {
			builder.NewLine()
		}
		if err := c.generateEnumDeclaration(builder, enumDef); err != nil {
			return err
		}
	}
	builder.Unindent().
		WriteLine("}")
	return nil
}
