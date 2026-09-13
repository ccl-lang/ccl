package rsGenerator

import (
	"strconv"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
)

func (c *RustGenerationContext) generateEnumMappings(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	for _, mapping := range enumDef.Mappings[CurrentLanguage] {
		sourceType, err := c.getRustEnumTypeName(enumDef)
		if err != nil {
			return err
		}
		targetType, err := c.getRustEnumTypeName(mapping.Target)
		if err != nil {
			return err
		}
		fallback, err := c.getRustEnumMemberName(mapping.Target, mapping.Fallback)
		if err != nil {
			return err
		}
		methodName := "to_" + gValues.StyleSnakeCase.ApplyStyle(targetType)
		builder.MapVarPairs(
			"mapSource", sourceType,
			"mapTarget", "crate::"+targetType,
			"mapMethod", methodName,
			"mapRawMethod", "map_raw_"+methodName,
			"mapBase", c.getRustEnumBaseType(enumDef),
			"mapFallback", fallback,
		)
		builder.NewLine().
			LineD("impl $mapSource {").
			Indent().
			LineD("pub fn $mapMethod(self) -> $mapTarget {").
			Indent().
			LineD("Self::$mapRawMethod(self as $mapBase)").
			Unindent().
			WriteLine("}").
			NewLine().
			LineD("pub fn $mapRawMethod(value: $mapBase) -> $mapTarget {").
			Indent().
			WriteLine("match value {").
			Indent()
		for _, branch := range mapping.Cases {
			target, err := c.getRustEnumMemberName(mapping.Target, branch.Target)
			if err != nil {
				return err
			}
			builder.MapVarPairs(
				"mapCase", strconv.FormatInt(branch.Source.Value, 10),
				"mapResult", target,
			)
			builder.LineD("$mapCase => $mapTarget::$mapResult,")
		}
		builder.LineD("_ => $mapTarget::$mapFallback,").
			Unindent().
			WriteLine("}").
			Unindent().
			WriteLine("}").
			Unindent().
			WriteLine("}").
			NewLine()
		builder.UnmapVar(
			"mapSource",
			"mapTarget",
			"mapMethod",
			"mapRawMethod",
			"mapBase",
			"mapFallback",
			"mapCase",
			"mapResult",
		)
	}
	return nil
}
