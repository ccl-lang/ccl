package goGenerator

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *GoGenerationContext) getGoEnumTypeName(enumDef *CCLEnum) string {
	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		return enumDef.OwnedBy.Name + enumDef.Name
	}

	return enumDef.Name
}

func (c *GoGenerationContext) getGoEnumMemberName(
	enumDef *CCLEnum,
	member *cclValues.EnumMemberDefinition,
) (string, error) {
	defaultPrefix := ""
	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		defaultPrefix = enumDef.OwnedBy.Name
	}

	prefix, err := c.GetEnumMemberNamePrefix(
		CurrentLanguage,
		enumDef,
		defaultPrefix,
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

	return prefix + enumDef.Name + style.ApplyStyle(member.Name), nil
}

func (c *GoGenerationContext) getGoEnumBaseType(enumDef *CCLEnum) string {
	if mappedType, ok := CCLTypesToGoTypes[enumDef.BaseType.GetName()]; ok {
		return mappedType
	}

	return "int32"
}

func (c *GoGenerationContext) getGoEnumOutputFileGroup(enumDef *CCLEnum) (string, error) {
	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		return c.getModelOutputFileGroup(enumDef.OwnedBy)
	}

	return c.GetEnumOutputFileGroup(CurrentLanguage, enumDef)
}

func (c *GoGenerationContext) getGoTypeForField(field *CCLField) (string, error) {
	targetType := field.Type
	if targetType.IsArray() {
		targetType = targetType.GetUnderlyingType()
	}

	goType, err := c.getGoTypeForUsage(targetType, field)
	if err != nil {
		return "", err
	}

	if field.IsArray() {
		goType = "[]" + goType
	}

	return goType, nil
}

func (c *GoGenerationContext) getGoTypeForUsage(
	targetType *cclValues.CCLTypeUsage,
	field *CCLField,
) (string, error) {
	if mappedType, ok := CCLTypesToGoTypes[targetType.GetName()]; ok {
		return mappedType, nil
	}

	if targetType.IsCustomTypeEnum() {
		return c.getGoEnumTypeName(targetType.GetDefinition().GetEnumDefinition()), nil
	}

	if targetType.IsCustomTypeModel() {
		customModel := targetType.GetDefinition().GetModelDefinition()
		if customModel != nil {
			return "*" + customModel.Name, nil
		}
	}

	modelName := ""
	fieldName := ""
	if field != nil {
		modelName = field.GetModelFullName()
		fieldName = field.Name
	}

	return "", unsupportedGoFieldType(targetType.GetName(), fieldName, modelName)
}
