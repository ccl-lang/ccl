package cclAst

const (
	TypeExpressionKindUnknown TypeExpressionKind = iota
	TypeExpressionKindSimpleName
	TypeExpressionKindArray
)

const (
	AttributeValueKindUnknown AttributeValueKind = iota
	AttributeValueKindLiteral
	AttributeValueKindIdentifier
)

const (
	AttributeLiteralKindUnknown AttributeLiteralKind = iota
	AttributeLiteralKindString
	AttributeLiteralKindInt
	AttributeLiteralKindFloat
	AttributeLiteralKindReserved
)

const (
	AttributeScopeUnknown AttributeScope = iota
	AttributeScopeGlobal
	AttributeScopeFile
	AttributeScopeNamespace
)
