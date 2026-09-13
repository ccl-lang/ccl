package pyGenerator

import (
	"strings"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
)

func (c *PythonGenerationContext) generateEnumMappings(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	for _, mapping := range enumDef.Mappings[CurrentLanguage] {
		targetType, err := c.getPythonEnumTypeName(mapping.Target)
		if err != nil {
			return err
		}
		fallback, err := c.getPythonEnumReference(mapping.Target, mapping.Fallback)
		if err != nil {
			return err
		}
		methodName := "to_" + gValues.StyleSnakeCase.ApplyStyle(strings.ReplaceAll(targetType, ".", "_"))
		builder.MapVarPairs(
			"mapMethod", methodName,
			"mapTarget", targetType,
			"mapFallback", fallback,
		)
		builder.NewLine().
			WriteLine("@staticmethod").
			LineD("def $mapMethod(value: int) -> '$mapTarget':").
			Indent()
		// Resolve cross-file imports on calls so reciprocal mappings do not create import cycles.
		if mapping.Target != enumDef && (enumDef.OwnedBy == nil || mapping.Target.OwnedBy != enumDef.OwnedBy) {
			importLine, err := c.getImportLineForEnum(mapping.Target)
			if err != nil {
				return err
			}
			builder.WriteLine(importLine)
		}
		if len(mapping.Cases) > 0 {
			builder.WriteLine("match value:").
				Indent()
			for _, branch := range mapping.Cases {
				source, err := c.getPythonEnumReference(enumDef, branch.Source)
				if err != nil {
					return err
				}
				target, err := c.getPythonEnumReference(mapping.Target, branch.Target)
				if err != nil {
					return err
				}
				builder.MapVarPairs(
					"mapCase", source,
					"mapResult", target,
				)
				builder.LineD("case $mapCase:").
					Indent().
					LineD("return $mapResult").
					Unindent()
			}
			builder.Unindent()
		}
		builder.LineD("return $mapFallback").
			Unindent()
		builder.UnmapVar(
			"mapMethod",
			"mapTarget",
			"mapFallback",
			"mapCase",
			"mapResult",
		)
	}
	return nil
}
