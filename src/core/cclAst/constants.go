package cclAst

// TypeExpressionKind represents the kind of a type expression node.
type TypeExpressionKind int

const (
	TypeExpressionKindUnknown TypeExpressionKind = iota
	TypeExpressionKindSimpleName
	TypeExpressionKindArray
)
