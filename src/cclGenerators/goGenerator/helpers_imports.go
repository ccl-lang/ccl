package goGenerator

import "github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"

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

func registerGoImports(builder *codeBuilder.CodeBuilder, imports map[string]bool) {
	orderedImports := [...]string{
		"bytes",
		"encoding/base64",
		"encoding/binary",
		"encoding/json",
		"strconv",
		"strings",
		"time",
	}

	for _, importPath := range orderedImports {
		if imports[importPath] {
			registerGoImport(builder, importPath)
		}
	}
}
