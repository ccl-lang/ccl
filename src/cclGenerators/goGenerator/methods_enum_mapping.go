package goGenerator

func (c *GoGenerationContext) generateEnumMappings(enumDef *CCLEnum) error {
	for _, mapping := range enumDef.Mappings[CurrentLanguage] {
		builder, err := c.getEnumCodeBuilder("methods", enumDef)
		if err != nil {
			return err
		}
		sourceType, err := c.getGoEnumTypeName(enumDef)
		if err != nil {
			return err
		}
		targetType, err := c.getGoEnumTypeName(mapping.Target)
		if err != nil {
			return err
		}
		fallback, err := c.getGoEnumMemberName(mapping.Target, mapping.Fallback)
		if err != nil {
			return err
		}
		builder.MapVarPairs(
			"mapSource", sourceType,
			"mapTarget", targetType,
			"mapFallback", fallback,
		)
		builder.NewLine().
			LineD("func (value $mapSource) To$mapTarget() $mapTarget {").
			Indent().
			WriteLine("switch value {")
		for _, branch := range mapping.Cases {
			source, err := c.getGoEnumMemberName(enumDef, branch.Source)
			if err != nil {
				return err
			}
			target, err := c.getGoEnumMemberName(mapping.Target, branch.Target)
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
		builder.WriteLine("default:").
			Indent().
			LineD("return $mapFallback").
			Unindent().
			WriteLine("}").
			Unindent().
			WriteLine("}").
			NewLine()
		builder.UnmapVar(
			"mapSource",
			"mapTarget",
			"mapFallback",
			"mapCase",
			"mapResult",
		)
	}
	return nil
}
