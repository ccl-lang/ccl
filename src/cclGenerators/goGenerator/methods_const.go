package goGenerator

import (
	"github.com/ALiwoto/ssg/ssg"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
)

func (c *GoGenerationContext) GenerateConstants() error {
	builder := c.getCodeBuilder(ConstantsFileName, "constants")
	models := []*CCLModel{}
	enums := []*CCLEnum{}

	for _, currentTypeDef := range c.GetGenerationTypeDefinitions() {
		if currentTypeDef.IsCustomModel() {
			models = append(models, currentTypeDef.GetModelDefinition())
			continue
		}

		if currentTypeDef.IsCustomEnum() {
			enums = append(enums, currentTypeDef.GetEnumDefinition())
			continue
		}

		return &cclErrors.UnsupportedTypeDefinitionError{
			TypeName:       currentTypeDef.GetFullName(),
			TargetLanguage: CurrentLanguage.String(),
		}
	}

	wroteConstBlock := false
	if len(models) != 0 {
		beginGoConstBlock(builder)
		for _, currentModel := range models {
			if err := c.generateConstantsForModel(builder, currentModel); err != nil {
				return err
			}
		}
		endGoConstBlock(builder)
		wroteConstBlock = true
	}

	for _, enumDef := range enums {
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

	return nil
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
	for _, member := range enumDef.Members {
		memberName, err := c.getGoEnumMemberName(enumDef, member)
		if err != nil {
			return err
		}
		builder.WriteLine(memberName + " " +
			c.getGoEnumTypeName(enumDef) + " = " + ssg.ToBase10(member.Value),
		)
	}
	return nil
}

//---------------------------------------------------------
