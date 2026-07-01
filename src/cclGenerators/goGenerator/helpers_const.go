package goGenerator

import "github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"

func beginGoConstBlock(builder *codeBuilder.CodeBuilder) {
	builder.WriteLine("const (").
		Indent()
}

func endGoConstBlock(builder *codeBuilder.CodeBuilder) {
	builder.Unindent().
		WriteLine(")")
}

func getGoConstantGroup(
	constantGroups map[string]*goConstantGroup,
	groupOrder *[]string,
	group string,
) *goConstantGroup {
	constantGroup := constantGroups[group]
	if constantGroup != nil {
		return constantGroup
	}

	constantGroup = &goConstantGroup{}
	constantGroups[group] = constantGroup
	*groupOrder = append(*groupOrder, group)
	return constantGroup
}
