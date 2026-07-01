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
