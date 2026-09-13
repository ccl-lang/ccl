package rsGenerator

import (
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *RustGenerationContext) generateEnumAliases(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
	enumTypeName string,
) error {
	values := map[int64]*cclValues.EnumMemberDefinition{}
	opened := false
	for _, member := range enumDef.Members {
		canonical := values[member.Value]
		if canonical == nil {
			values[member.Value] = member
			continue
		}
		aliasName, err := c.getRustEnumMemberName(enumDef, member)
		if err != nil {
			return err
		}
		canonicalName, err := c.getRustEnumMemberName(enumDef, canonical)
		if err != nil {
			return err
		}
		if !opened {
			builder.MapVarPairs("aliasType", enumTypeName)
			builder.LineD("impl $aliasType {").
				Indent()
			opened = true
		}
		builder.MapVarPairs(
			"aliasName", aliasName,
			"aliasValue", canonicalName,
		)
		builder.WriteLine("#[allow(non_upper_case_globals)]").
			LineD("pub const $aliasName: Self = Self::$aliasValue;")
	}
	if opened {
		builder.Unindent().
			WriteLine("}").
			NewLine()
	}
	builder.UnmapVar("aliasType", "aliasName", "aliasValue")
	return nil
}
