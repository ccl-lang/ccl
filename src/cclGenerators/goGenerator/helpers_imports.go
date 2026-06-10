package goGenerator

import "github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"

func registerGoImport(builder *codeBuilder.CodeBuilder, importPath string) {
	if !builder.IsImported(openImportKey) {
		builder.DoImport(openImportKey, "import (")
	}

	indentStr := builder.GetIndentationStr()
	builder.DoImport(importPath, indentStr+"\""+importPath+"\"")
}

func closeGoImportGroup(builder *codeBuilder.CodeBuilder) {
	if builder.IsImported(openImportKey) && !builder.IsImported(closeImportKey) {
		builder.DoImport(closeImportKey, ")")
	}
}
