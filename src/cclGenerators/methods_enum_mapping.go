package cclGenerators

import (
	"fmt"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

// GetEnumMappingMethodName resolves naming preferences and reserves the generated method name.
func (c *CodeGenerationBase) GetEnumMappingMethodName(
	language gValues.LanguageType,
	source, target *cclValues.EnumDefinition,
	longName string,
) (string, error) {
	attr := c.GetGlobalOrEnumAttributes(
		language,
		cclAttr.AttrEnumMapUseShortMethodName,
		source,
	).GetLast()
	methodName := longName
	if attr != nil {
		if len(attr.Parameters) != 1 {
			return "", &cclErrors.InvalidAttributeUsageError{
				AttrName:       attr.Name,
				Message:        "requires exactly one boolean parameter (true or false)",
				SourcePosition: attr.SourcePosition,
			}
		}
		useShortName, valid := attr.GetParamAt(0).GetValue().(bool)
		if !valid {
			return "", &cclErrors.InvalidAttributeUsageError{
				AttrName:       attr.Name,
				Message:        "requires a boolean parameter (true or false)",
				SourcePosition: attr.SourcePosition,
			}
		}
		if useShortName {
			switch language {
			case gValues.LanguageGo, gValues.LanguageCS:
				methodName = "To" + target.Name
			case gValues.LanguageJS, gValues.LanguageTS:
				methodName = "to" + target.Name
			default:
				methodName = "to_" + gValues.StyleSnakeCase.ApplyStyle(target.Name)
			}
		}
	}

	scope := source.GetFullName()
	if language == gValues.LanguageGd && source.OwnedBy != nil {
		scope = source.OwnedBy.GetFullName()
		for _, field := range source.OwnedBy.Fields {
			if gValues.StyleSnakeCase.ApplyStyle(field.Name) == methodName {
				return "", c.enumMappingNameConflict(source, methodName, scope)
			}
		}
	}
	if c.enumMappingMethodNames == nil {
		c.enumMappingMethodNames = make(map[string]map[string]bool)
	}
	if c.enumMappingMethodNames[scope] == nil {
		c.enumMappingMethodNames[scope] = make(map[string]bool)
	}
	if c.enumMappingMethodNames[scope][methodName] {
		return "", c.enumMappingNameConflict(source, methodName, scope)
	}
	c.enumMappingMethodNames[scope][methodName] = true
	return methodName, nil
}

func (c *CodeGenerationBase) enumMappingNameConflict(
	source *cclValues.EnumDefinition,
	methodName, scope string,
) error {
	return &cclErrors.InvalidAttributeUsageError{
		AttrName: cclAttr.AttrEnumMapUseShortMethodName,
		Message: fmt.Sprintf(
			"mapping method %q collides in %s; set EnumMapUseShortMethodName(false) on %s or use distinct enum names",
			methodName, scope, source.GetFullName(),
		),
		SourcePosition: source.SourcePosition,
	}
}
