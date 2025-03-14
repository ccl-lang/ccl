package cclLexer

const (
	TokenTypeReservedForFuture CCLTokenType = -1
)

const (
	TokenTypeError CCLTokenType = iota
	TokenTypeComment
	TokenTypeHash
	TokenTypeKeywordModel
	TokenTypeIdentifier
	TokenTypeColon
	TokenTypeSemicolon
	TokenTypeDataType
	TokenTypeLeftBrace
	TokenTypeRightBrace
	TokenTypeLeftBracket
	TokenTypeRightBracket
	TokenTypeStringLiteral
	TokenTypeWhitespace
	TokenTypeLeftParenthesis
	TokenTypeRightParenthesis
	TokenTypeDot
	TokenTypeComma
)
