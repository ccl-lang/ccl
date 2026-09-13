package gdGenerator

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
)

func (c *GDScriptGenerationContext) generateEnumMappings(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	for _, mapping := range enumDef.Mappings[CurrentLanguage] {
		sourceType, err := c.getGDScriptEnumTypeName(enumDef)
		if err != nil {
			return err
		}
		targetType, err := c.getGDScriptEnumTypeName(mapping.Target)
		if err != nil {
			return err
		}
		targetDeclaration, err := c.getGDScriptEnumDeclarationName(mapping.Target)
		if err != nil {
			return err
		}
		if mapping.Target.IsNested() {
			targetType = mapping.Target.OwnedBy.Name
		}
		targetReference := targetType + "." + targetDeclaration
		methodTarget := targetType
		if mapping.Target.IsNested() {
			methodTarget += "_" + targetDeclaration
		}
		fallback, err := c.getGDScriptEnumReference(mapping.Target, mapping.Fallback, enumDef.OwnedBy)
		if err != nil {
			return err
		}
		methodName := gValues.StyleSnakeCase.ApplyStyle(sourceType + "_to_" + methodTarget)
		methodName, err = c.GetEnumMappingMethodName(CurrentLanguage, enumDef, mapping.Target, methodName)
		if err != nil {
			return err
		}
		builder.MapVarPairs(
			"mapMethod", methodName,
			"mapTarget", targetReference,
			"mapFallback", fallback,
		)
		builder.NewLine().
			LineD("static func $mapMethod(value: int) -> $mapTarget:").
			Indent()
		if len(mapping.Cases) > 0 {
			builder.WriteLine("match value:").
				Indent()
			for _, branch := range mapping.Cases {
				source, err := c.getGDScriptEnumReference(enumDef, branch.Source, enumDef.OwnedBy)
				if err != nil {
					return err
				}
				target, err := c.getGDScriptEnumReference(mapping.Target, branch.Target, enumDef.OwnedBy)
				if err != nil {
					return err
				}
				builder.MapVarPairs(
					"mapCase", source,
					"mapResult", target,
				)
				builder.LineD("$mapCase:").
					Indent().
					LineD("return $mapResult").
					Unindent()
			}
			builder.Unindent()
		}
		builder.LineD("return $mapFallback").
			Unindent().
			NewLine()
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
