package cclAst

import "github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils"

// CCLFileAST represents a parsed CCL source file (syntax only).
type CCLFileAST struct {
	FilePath            string
	Namespace           string
	Imports             []*ImportDecl
	GlobalAttributes    []*GlobalAttributeNode
	FileAttributes      []*GlobalAttributeNode
	NamespaceAttributes []*GlobalAttributeNode
	Models              []*ModelDecl
	Enums               []*EnumDecl
	SourcePosition      *cclUtils.SourceCodePosition
}

// ImportDecl represents a source-file import declaration in CCL.
type ImportDecl struct {
	Path           string
	SourcePosition *cclUtils.SourceCodePosition
}

// ModelDecl represents a model declaration in CCL.
type ModelDecl struct {
	Name           string
	Namespace      string
	Fields         []*FieldDecl
	Enums          []*EnumDecl
	Attributes     []*AttributeNode
	SourcePosition *cclUtils.SourceCodePosition
}

// FieldDecl represents a field declaration inside a model.
type FieldDecl struct {
	Name           string
	Type           TypeExpression
	Value          ValueExpression
	Attributes     []*AttributeNode
	SourcePosition *cclUtils.SourceCodePosition
}

// EnumDecl represents an enum declaration in CCL.
type EnumDecl struct {
	Name           string
	Namespace      string
	BaseType       TypeExpression
	Members        []*EnumMemberDecl
	Attributes     []*AttributeNode
	SourcePosition *cclUtils.SourceCodePosition
}

// EnumMemberDecl represents one enum member declaration.
type EnumMemberDecl struct {
	Name           string
	Value          *int64
	SourcePosition *cclUtils.SourceCodePosition
}
