package goGenerator

import (
	"github.com/ALiwoto/ssg/ssg"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
)

func (c *GoGenerationContext) GenerateConstants() error {
	constantGroups := map[string]*goConstantGroup{}
	groupOrder := []string{}
	models := []*CCLModel{}

	for _, currentTypeDef := range c.GetGenerationTypeDefinitions() {
		if currentTypeDef.IsCustomModel() {
			models = append(models, currentTypeDef.GetModelDefinition())
			continue
		}

		if currentTypeDef.IsCustomEnum() {
			enumDef := currentTypeDef.GetEnumDefinition()
			group, err := c.getGoEnumOutputFileGroup(enumDef)
			if err != nil {
				return err
			}
			constantGroup := getGoConstantGroup(constantGroups, &groupOrder, group)
			constantGroup.Enums = append(constantGroup.Enums, enumDef)
			continue
		}

		return &cclErrors.UnsupportedTypeDefinitionError{
			TypeName:       currentTypeDef.GetFullName(),
			TargetLanguage: CurrentLanguage.String(),
		}
	}

	if len(models) != 0 {
		builder := c.getConstantsCodeBuilder("")
		beginGoConstBlock(builder)
		for _, currentModel := range models {
			if err := c.generateConstantsForModel(builder, currentModel); err != nil {
				return err
			}
		}
		endGoConstBlock(builder)
	}

	for _, group := range groupOrder {
		constantGroup := constantGroups[group]
		builder := c.getConstantsCodeBuilder(group)
		wroteConstBlock := len(models) != 0 && group == ""

		for _, enumDef := range constantGroup.Enums {
			if len(enumDef.Members) == 0 {
				continue
			}

			if wroteConstBlock {
				builder.NewLine()
			}
			beginGoConstBlock(builder)
			if err := c.generateConstantsForEnum(builder, enumDef); err != nil {
				return err
			}
			endGoConstBlock(builder)
			wroteConstBlock = true
		}
	}

	return nil
}

func (c *GoGenerationContext) getConstantsCodeBuilder(group string) *codeBuilder.CodeBuilder {
	return c.getCodeBuilder(getGoCategoryFileName("constants", group), "constants")
}

func (c *GoGenerationContext) generateConstantsForModel(
	builder *codeBuilder.CodeBuilder,
	currentModel *CCLModel,
) error {
	builder.WriteLine("ModelId" +
		currentModel.Name + " = " + ssg.ToBase10(currentModel.ModelId),
	)
	return nil
}

func (c *GoGenerationContext) generateConstantsForEnum(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	enumTypeName, err := c.getGoEnumTypeName(enumDef)
	if err != nil {
		return err
	}

	for _, member := range enumDef.Members {
		memberName, err := c.getGoEnumMemberName(enumDef, member)
		if err != nil {
			return err
		}
		builder.WriteLine(memberName + " " +
			enumTypeName + " = " + ssg.ToBase10(member.Value),
		)
	}
	return nil
}

//---------------------------------------------------------
