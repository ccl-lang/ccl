package pyGenerator

import (
	"github.com/ALiwoto/ssg/ssg"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *PythonGenerationContext) getPythonEnumTypeName(enumDef *CCLEnum) string {
	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		return enumDef.OwnedBy.Name + "." + enumDef.Name
	}

	return enumDef.Name
}

func (c *PythonGenerationContext) getPythonEnumMemberName(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	style, err := c.GetEnumMemberNamingStyle(
		CurrentLanguage,
		enumDef,
		gValues.StyleUpperCase,
	)
	if err != nil {
		return "", err
	}

	return style.ApplyStyle(member.Name), nil
}

func (c *PythonGenerationContext) getPythonEnumReference(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	memberName, err := c.getPythonEnumMemberName(enumDef, member)
	if err != nil {
		return "", err
	}

	return c.getPythonEnumTypeName(enumDef) + "." + memberName, nil
}

func (c *PythonGenerationContext) generateEnumClass(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	builder.WriteLine("class " + enumDef.Name + "(IntEnum):").
		Indent()
	if len(enumDef.Members) == 0 {
		builder.WriteLine("pass")
	} else {
		for _, member := range enumDef.Members {
			memberName, err := c.getPythonEnumMemberName(enumDef, member)
			if err != nil {
				return err
			}
			builder.WriteLine(memberName + " = " + ssg.ToBase10(member.Value))
		}
	}
	builder.UnindentLine()
	return nil
}

func (c *PythonGenerationContext) getImportLineForEnum(enumDef *CCLEnum) (string, error) {
	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		fileName, err := c.GetFileNameForModel(
			CurrentLanguage,
			enumDef.OwnedBy,
			DefaultFileNamingStyle,
			supportedFileNameStyles,
		)
		if err != nil {
			return "", err
		}

		return "from ." + fileName + " import " + enumDef.OwnedBy.Name, nil
	}

	return "from ." + pythonEnumFileName(enumDef) + " import " + enumDef.Name, nil
}

func (c *PythonGenerationContext) pythonEnumCastExpression(
	typeUsage *cclValues.CCLTypeUsage,
	valueExpression string,
) string {
	if !typeUsage.IsCustomTypeEnum() {
		return valueExpression
	}

	return c.getPythonEnumTypeName(typeUsage.GetDefinition().GetEnumDefinition()) +
		"(" + valueExpression + ")"
}
