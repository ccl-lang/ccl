package gdGenerator

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *GDScriptGenerationContext) getGDScriptEnumMemberName(
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

func (c *GDScriptGenerationContext) getGDScriptEnumReference(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	memberName, err := c.getGDScriptEnumMemberName(enumDef, member)
	if err != nil {
		return "", err
	}

	if enumDef.IsNested() {
		return enumDef.Name + "." + memberName, nil
	}

	return enumDef.Name + "." + memberName, nil
}
