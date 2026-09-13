package csGenerator

import (
	"strings"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
)

func (c *CSharpGenerationContext) generateEnumMappings(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	mappings := enumDef.Mappings[CurrentLanguage]
	if len(mappings) == 0 {
		return nil
	}
	sourceType, err := c.getCSharpEnumTypeName(enumDef)
	if err != nil {
		return err
	}
	builder.MapVarPairs(
		"mapSource", sourceType,
		"mapHelpers", sourceType+"Mappings",
	)
	builder.NewLine().
		LineD("public static class $mapHelpers").
		WriteLine("{").
		Indent()
	for _, mapping := range mappings {
		targetType, err := c.getCSharpEnumTypeReference(mapping.Target, nil)
		if err != nil {
			return err
		}
		fallback, err := c.getCSharpEnumReference(mapping.Target, mapping.Fallback, nil)
		if err != nil {
			return err
		}
		methodName, err := c.GetEnumMappingMethodName(CurrentLanguage, enumDef, mapping.Target, "To"+strings.ReplaceAll(targetType, ".", ""))
		if err != nil {
			return err
		}
		builder.MapVarPairs(
			"mapTarget", targetType,
			"mapMethod", methodName,
			"mapFallback", fallback,
		)
		builder.LineD("public static $mapTarget $mapMethod($mapSource value)").
			WriteLine("{").
			Indent().
			WriteLine("switch (value)").
			WriteLine("{").
			Indent()
		for _, branch := range mapping.Cases {
			source, err := c.getCSharpEnumReference(enumDef, branch.Source, enumDef.OwnedBy)
			if err != nil {
				return err
			}
			target, err := c.getCSharpEnumReference(mapping.Target, branch.Target, nil)
			if err != nil {
				return err
			}
			builder.MapVarPairs(
				"mapCase", source,
				"mapResult", target,
			)
			builder.LineD("case $mapCase:").
				Indent().
				LineD("return $mapResult;").
				Unindent()
		}
		builder.WriteLine("default:").
			Indent().
			LineD("return $mapFallback;").
			Unindent().
			Unindent().
			WriteLine("}").
			Unindent().
			WriteLine("}")
	}
	builder.Unindent().
		WriteLine("}")
	builder.UnmapVar(
		"mapSource",
		"mapHelpers",
		"mapTarget",
		"mapMethod",
		"mapFallback",
		"mapCase",
		"mapResult",
	)
	return nil
}
