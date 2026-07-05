package gdGenerator

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *GDScriptGenerationContext) getGDScriptEnumMemberName(
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

func (c *GDScriptGenerationContext) getGDScriptEnumTypeName(enumDef *CCLEnum) (string, error) {
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

func (c *GDScriptGenerationContext) getGDScriptEnumDeclarationName(enumDef *CCLEnum) (string, error) {
	enumTypeName, err := c.getGDScriptEnumTypeName(enumDef)
	if err != nil {
		return "", err
	}

	if enumDef.IsNested() {
		return enumTypeName, nil
	}

	return enumTypeName + "Enum", nil
}

func (c *GDScriptGenerationContext) getGDScriptEnumTypeReference(
	typeUsage *cclValues.CCLTypeUsage,
	currentModel *CCLModel,
) (string, error) {
	enumDef := typeUsage.GetDefinition().GetEnumDefinition()
	enumTypeName, err := c.getGDScriptEnumTypeName(enumDef)
	if err != nil {
		return "", err
	}

	enumDeclarationName, err := c.getGDScriptEnumDeclarationName(enumDef)
	if err != nil {
		return "", err
	}

	if enumDef.IsNested() {
		if enumDef.OwnedBy == currentModel {
			return enumDeclarationName, nil
		}

		return enumDef.OwnedBy.Name + "." + enumDeclarationName, nil
	}

	return enumTypeName + "." + enumDeclarationName, nil
}

func (c *GDScriptGenerationContext) getGDScriptEnumCastSuffix(
	typeUsage *cclValues.CCLTypeUsage,
	currentModel *CCLModel,
) (string, error) {
	if !typeUsage.IsCustomTypeEnum() {
		return "", nil
	}

	enumTypeReference, err := c.getGDScriptEnumTypeReference(typeUsage, currentModel)
	if err != nil {
		return "", err
	}

	return " as " + enumTypeReference, nil
}

func (c *GDScriptGenerationContext) getGDScriptEnumReference(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
	currentModel *CCLModel,
) (string, error) {
	memberName, err := c.getGDScriptEnumMemberName(enumDef, member)
	if err != nil {
		return "", err
	}

	enumTypeName, err := c.getGDScriptEnumTypeName(enumDef)
	if err != nil {
		return "", err
	}

	if enumDef.IsNested() {
		if enumDef.OwnedBy == currentModel {
			return enumTypeName + "." + memberName, nil
		}

		return enumDef.OwnedBy.Name + "." + enumTypeName + "." + memberName, nil
	}

	enumDeclarationName, err := c.getGDScriptEnumDeclarationName(enumDef)
	if err != nil {
		return "", err
	}

	return enumTypeName + "." + enumDeclarationName + "." + memberName, nil
}
