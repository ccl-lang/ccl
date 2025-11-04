package cclGenerators

import (
	"slices"

	"github.com/ccl-lang/ccl/src/core/cclValues"
)

//---------------------------------------------------------

func (c *CodeGenerationBase) GetGlobalAttribute(name string) *cclValues.AttributeUsageInfo {
	return c.Options.CCLDefinition.FindGlobalAttribute(name)
}

func (c *CodeGenerationBase) GetGlobalAttributes(name string) []*cclValues.AttributeUsageInfo {
	return c.Options.CCLDefinition.FindGlobalAttributes(name)
}

//---------------------------------------------------------

func (c *CodeGenerationBase) NeedsCloneMethods(
	currentModel *cclValues.ModelDefinition,
) bool {
	attr := c.GetGlobalAttribute("AddCloneMethods")
	if attr == nil {
		attr = currentModel.FindAttribute("AddCloneMethods")
	}

	if attr != nil {
		return attr.GetParamAt(0).GetAsBool()
	}
	return false
}

func (c *CodeGenerationBase) NeedsBinarySerialization(
	currentModel *cclValues.ModelDefinition,
) bool {
	return c.NeedsSerializationType(currentModel, "binary")
}

func (c *CodeGenerationBase) NeedsSerializationType(
	currentModel *cclValues.ModelDefinition,
	sType string,
) bool {
	attrs := c.GetGlobalAttributes("SerializationType")
	if len(attrs) == 0 {
		attrs = append(attrs, currentModel.FindAttributes("SerializationType")...)
	}

	collection := cclValues.NewAttrsCollection(attrs)
	return slices.Contains(collection.GetParamsAtAsStrings(0), sType)
}

//---------------------------------------------------------
