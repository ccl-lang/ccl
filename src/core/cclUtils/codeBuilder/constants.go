package codeBuilder

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
