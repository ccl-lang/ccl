package goGenerator

import "github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"

func registerGoImport(builder *codeBuilder.CodeBuilder, importPath string) {
	if !builder.IsImported("__go_import_group_open") {
		builder.DoImport("__go_import_group_open", "import (")
	}

	builder.DoImport(importPath, "\t\""+importPath+"\"")
}

func closeGoImportGroup(builder *codeBuilder.CodeBuilder) {
	if builder.IsImported("__go_import_group_open") && !builder.IsImported("__go_import_group_close") {
		builder.DoImport("__go_import_group_close", ")")
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
