package cclAst

import "github.com/ccl-lang/ccl/src/core/cclUtils"

// AttributeLiteralKind represents the literal kind of an attribute value.
type AttributeLiteralKind int

// AttributeValueKind represents the kind of a value inside an attribute parameter.
type AttributeValueKind int

// TypeExpressionKind represents the kind of a type expression node.
type TypeExpressionKind int

// SimpleTypeName represents a simple (unqualified) type name with an optional namespace.
type SimpleTypeName struct {
	Name      string
	Namespace string
}

// TypeExpression represents a type expression node in the AST.
type TypeExpression interface {
	GetTypeExpressionKind() TypeExpressionKind
	GetSourcePosition() *cclUtils.SourceCodePosition
}

// SimpleTypeExpression represents a simple type usage, such as "string" or "User".
type SimpleTypeExpression struct {
	TypeName       SimpleTypeName
	IsBuiltinToken bool
	SourcePosition *cclUtils.SourceCodePosition
}

// ArrayTypeExpression represents an array type usage, such as "string[]" or "int[10]".
// Length is -1 for dynamic-length arrays.
type ArrayTypeExpression struct {
	ElementType    TypeExpression
	Length         int
	SourcePosition *cclUtils.SourceCodePosition
}
