package jsGenerator

import "github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"

func (c *JavaScriptGenerationContext) generateEnumMappings(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	mappings := enumDef.Mappings[LanguageName]
	if len(mappings) == 0 {
		return nil
	}
	sourceType, err := c.getJavaScriptEnumTypeName(enumDef)
	if err != nil {
		return err
	}
	builder.MapVarPairs("mapHelpers", sourceType+"Mappings")
	builder.NewLine()
	if enumDef.IsNested() {
		builder.LineD("static $mapHelpers = class {")
	} else {
		builder.LineD("export class $mapHelpers {")
	}
	builder.Indent()
	for _, mapping := range mappings {
		targetType, err := c.getJavaScriptEnumTypeName(mapping.Target)
		if err != nil {
			return err
		}
		methodName := "to" + targetType
		if mapping.Target.OwnedBy != nil {
			methodName = "to" + mapping.Target.OwnedBy.Name + targetType
		}
		if !c.IsSingleFile && mapping.Target != enumDef && (enumDef.OwnedBy == nil || mapping.Target.OwnedBy != enumDef.OwnedBy) {
			importName := targetType
			fileName := mapping.Target.Name
			if mapping.Target.OwnedBy != nil {
				importName = mapping.Target.OwnedBy.Name
				fileName = importName
			}
			builder.DoImport(importName, "import { "+importName+" } from './"+fileName+".js';")
		}
		fallback, err := c.getJavaScriptEnumReference(mapping.Target, mapping.Fallback)
		if err != nil {
			return err
		}
		builder.MapVarPairs(
			"mapMethod", methodName,
			"mapFallback", fallback,
		)
		builder.LineD("static $mapMethod(value) {").
			Indent().
			WriteLine("switch (value) {").
			Indent()
		for _, branch := range mapping.Cases {
			source, err := c.getJavaScriptEnumReference(enumDef, branch.Source)
			if err != nil {
				return err
			}
			target, err := c.getJavaScriptEnumReference(mapping.Target, branch.Target)
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
		"mapHelpers",
		"mapMethod",
		"mapFallback",
		"mapCase",
		"mapResult",
	)
	return nil
}
