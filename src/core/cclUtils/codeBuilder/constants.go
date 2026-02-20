package codeBuilder

const (
	panicStatement = `
This error is caused by a developer misuse and is originated from the codeBuilder package.
It's not something users should fix or workaround.
Please open a GitHub issue with the description of this error, stacktrace and reproduction steps:
https://github.com/ccl-lang/ccl/issues/new
`
)

const (
	// SectionCommentHeaders is the section for file comment headers (at the very top of the file)
	SectionCommentHeaders = "comment_headers"

	// SectionHeaders is the section for file headers (like package declaration, file comments, etc)
	SectionHeaders = "headers"

	// SectionImports is the section for importing other types/modules
	SectionImports = "imports"

	// SectionDeclareNamespace is the section for declaring namespaces, used in languages such as C#,
	// where the namespace declaration MUST BE after the imports and before any type declarations.
	SectionDeclareNamespace = "declare_namespace"
)

const (
	varIndicator = '$'
)
