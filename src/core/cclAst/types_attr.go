package cclAst

import "github.com/ccl-lang/ccl/src/core/cclUtils"

// AttributeNodeBase is the common interface for attribute nodes.
type AttributeNodeBase interface {
	GetAttributeName() string
	GetAttributeParams() []*AttributeParamNode
	GetSourcePosition() *cclUtils.SourceCodePosition
}

// AttributeNode represents a normal attribute usage, such as [MyAttribute(...)].
type AttributeNode struct {
	Name           string
	Params         []*AttributeParamNode
	SourcePosition *cclUtils.SourceCodePosition
}

// GlobalAttributeNode represents a global attribute usage, such as #[MyAttribute(...)].
type GlobalAttributeNode struct {
	Name           string
	Params         []*AttributeParamNode
	SourcePosition *cclUtils.SourceCodePosition
}

// AttributeParamNode represents a parameter passed to an attribute.
// Name is optional for unnamed parameters.
type AttributeParamNode struct {
	Name           string
	Value          AttributeValueExpression
	SourcePosition *cclUtils.SourceCodePosition
}

// AttributeValueExpression represents a value inside an attribute parameter.
type AttributeValueExpression interface {
	GetAttributeValueKind() AttributeValueKind
	GetSourcePosition() *cclUtils.SourceCodePosition
}

// LiteralValueExpression represents a literal value in an attribute parameter.
type LiteralValueExpression struct {
	LiteralKind     AttributeLiteralKind
	Value           any
	ReservedLiteral string
	SourcePosition  *cclUtils.SourceCodePosition
}

// IdentifierValueExpression represents a variable usage inside an attribute parameter.
type IdentifierValueExpression struct {
	Name           string
	SourcePosition *cclUtils.SourceCodePosition
}
