package cclAst

import "github.com/ccl-lang/ccl/src/core/cclUtils"

// CCLFileAST represents a parsed CCL source file (syntax only).
type CCLFileAST struct {
	FilePath         string
	Namespace        string
	Imports          []*ImportDecl
	GlobalAttributes []*GlobalAttributeNode
	Models           []*ModelDecl
	SourcePosition   *cclUtils.SourceCodePosition
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
