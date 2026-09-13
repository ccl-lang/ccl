package tsGenerator

import (
	"strings"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
)

func (c *TypeScriptGenerationContext) generateEnumMappings(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	mappings := enumDef.Mappings[CurrentLanguage]
	if len(mappings) == 0 {
		return nil
	}
	sourceType, err := c.getTypeScriptEnumLocalTypeName(enumDef)
	if err != nil {
		return err
	}
	builder.MapVarPairs(
		"mapSource", sourceType,
		"mapHelpers", sourceType+"Mappings",
	)
	builder.NewLine().
		LineD("export namespace $mapHelpers {").
		Indent()
	for _, mapping := range mappings {
		targetType, err := c.getTypeScriptEnumTypeName(mapping.Target)
		if err != nil {
			return err
		}
		if !c.IsSingleFile && mapping.Target != enumDef && (enumDef.OwnedBy == nil || mapping.Target.OwnedBy != enumDef.OwnedBy) {
			importName, err := c.getTypeScriptEnumLocalTypeName(mapping.Target)
			if err != nil {
				return err
			}
			fileName := mapping.Target.Name
			if mapping.Target.OwnedBy != nil {
				importName = mapping.Target.OwnedBy.Name
				fileName = importName
			}
			builder.DoImport(importName, "import { "+importName+" } from './"+fileName+"';")
		}
		fallback, err := c.getTypeScriptEnumReference(mapping.Target, mapping.Fallback)
		if err != nil {
			return err
		}
		builder.MapVarPairs(
			"mapTarget", targetType,
			"mapMethod", "to"+strings.ReplaceAll(targetType, ".", ""),
			"mapFallback", fallback,
		)
		builder.LineD("export function $mapMethod(value: $mapSource): $mapTarget {").
			Indent().
			WriteLine("switch (value) {").
			Indent()
		for _, branch := range mapping.Cases {
			source, err := c.getTypeScriptEnumReference(enumDef, branch.Source)
			if err != nil {
				return err
			}
			target, err := c.getTypeScriptEnumReference(mapping.Target, branch.Target)
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
