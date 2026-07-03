package pyGenerator

import (
	"github.com/ALiwoto/ssg/ssg"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *PythonGenerationContext) getPythonEnumLocalTypeName(enumDef *CCLEnum) (string, error) {
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

func (c *PythonGenerationContext) getPythonEnumTypeName(enumDef *CCLEnum) (string, error) {
	enumTypeName, err := c.getPythonEnumLocalTypeName(enumDef)
	if err != nil {
		return "", err
	}

	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		return enumDef.OwnedBy.Name + "." + enumTypeName, nil
	}

	return enumTypeName, nil
}

func (c *PythonGenerationContext) getPythonEnumMemberName(
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
		gValues.StyleUpperCase,
	)
	if err != nil {
		return "", err
	}

	return prefix + style.ApplyStyle(member.Name), nil
}

func (c *PythonGenerationContext) getPythonEnumReference(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	memberName, err := c.getPythonEnumMemberName(enumDef, member)
	if err != nil {
		return "", err
	}

	enumTypeName, err := c.getPythonEnumTypeName(enumDef)
	if err != nil {
		return "", err
	}

	return enumTypeName + "." + memberName, nil
}

func (c *PythonGenerationContext) generateEnumClass(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	enumTypeName, err := c.getPythonEnumLocalTypeName(enumDef)
	if err != nil {
		return err
	}

	builder.WriteLine("class " + enumTypeName + "(IntEnum):").
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

	enumTypeName, err := c.getPythonEnumLocalTypeName(enumDef)
	if err != nil {
		return "", err
	}

	return "from ." + pythonEnumFileName(enumDef) + " import " + enumTypeName, nil
}

func (c *PythonGenerationContext) pythonEnumCastExpression(
	typeUsage *cclValues.CCLTypeUsage,
	valueExpression string,
) (string, error) {
	if !typeUsage.IsCustomTypeEnum() {
		return valueExpression, nil
	}

	enumTypeName, err := c.getPythonEnumTypeName(typeUsage.GetDefinition().GetEnumDefinition())
	if err != nil {
		return "", err
	}

	return enumTypeName + "(" + valueExpression + ")", nil
}
