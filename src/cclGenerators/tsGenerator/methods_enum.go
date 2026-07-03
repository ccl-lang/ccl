package tsGenerator

import (
	"github.com/ALiwoto/ssg/ssg"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *TypeScriptGenerationContext) getTypeScriptEnumLocalTypeName(enumDef *CCLEnum) (string, error) {
	prefix, err := c.GetEnumTypeNamePrefix(
		CurrentLanguage,
		enumDef,
		"",
	)
	if err != nil {
		return "", err
	}

	return prefix + enumDef.Name, nil
}

func (c *TypeScriptGenerationContext) getTypeScriptEnumTypeName(enumDef *CCLEnum) (string, error) {
	enumTypeName, err := c.getTypeScriptEnumLocalTypeName(enumDef)
	if err != nil {
		return "", err
	}

	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		return enumDef.OwnedBy.Name + "." + enumTypeName, nil
	}

	return enumTypeName, nil
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

	enumTypeName, err := c.getTypeScriptEnumTypeName(enumDef)
	if err != nil {
		return "", err
	}

	return enumTypeName + "." + memberName, nil
}

func (c *TypeScriptGenerationContext) generateEnumDeclaration(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	enumTypeName, err := c.getTypeScriptEnumLocalTypeName(enumDef)
	if err != nil {
		return err
	}

	builder.WriteLine("export enum " + enumTypeName + " {").
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
