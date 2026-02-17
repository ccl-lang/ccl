package cclAst

// TypeExpressionKind represents the kind of a type expression node.
type TypeExpressionKind int

const (
	TypeExpressionKindUnknown TypeExpressionKind = iota
	TypeExpressionKindSimpleName
	TypeExpressionKindArray
)

// AttributeValueKind represents the kind of a value inside an attribute parameter.
type AttributeValueKind int

const (
	AttributeValueKindUnknown AttributeValueKind = iota
	AttributeValueKindLiteral
	AttributeValueKindIdentifier
)

// AttributeLiteralKind represents the literal kind of an attribute value.
type AttributeLiteralKind int

const (
	AttributeLiteralKindUnknown AttributeLiteralKind = iota
	AttributeLiteralKindString
	AttributeLiteralKindInt
	AttributeLiteralKindFloat
	AttributeLiteralKindReserved
)
