package cclLexer

var (
	tokenTypesToNames = map[CCLTokenType]string{
		TokenTypeReservedForFuture: "RESERVED_FOR_FUTURE",
		TokenTypeError:             "ERROR",
		TokenTypeComment:           "COMMENT",
		TokenTypeHash:              "HASH",
		TokenTypeKeywordModel:      "KEYWORD_MODEL",
		TokenTypeIdentifier:        "IDENTIFIER",
		TokenTypeColon:             "COLON",
		TokenTypeSemicolon:         "SEMICOLON",
		TokenTypeDataType:          "DATA_TYPE",
		TokenTypeLeftBrace:         "LEFT_BRACE",
		TokenTypeRightBrace:        "RIGHT_BRACE",
		TokenTypeLeftBracket:       "LEFT_BRACKET",
		TokenTypeRightBracket:      "RIGHT_BRACKET",
		TokenTypeStringLiteral:     "STRING_LITERAL",
		TokenTypeWhitespace:        "WHITESPACE",
		TokenTypeLeftParenthesis:   "LEFT_PARENTHESIS",
		TokenTypeRightParenthesis:  "RIGHT_PARENTHESIS",
		TokenTypeDot:               "DOT",
		TokenTypeComma:             "COMMA",
	}

	oneCharSimpleTokens = map[rune]CCLTokenType{
		':': TokenTypeColon,
		';': TokenTypeSemicolon,
		'{': TokenTypeLeftBrace,
		'}': TokenTypeRightBrace,
		'[': TokenTypeLeftBracket,
		']': TokenTypeRightBracket,
		'(': TokenTypeLeftParenthesis,
		')': TokenTypeRightParenthesis,
		',': TokenTypeComma,
		'.': TokenTypeDot,
		'#': TokenTypeHash,
	}
)
