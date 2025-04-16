package cclLexer

import "github.com/ccl-lang/ccl/src/core/cclValues"

var (
	tokenTypesToNames = map[CCLTokenType]string{
		TokenTypeReservedForFuture:  "RESERVED_FOR_FUTURE",
		TokenTypeError:              "ERROR",
		TokenTypeComment:            "COMMENT",
		TokenTypeHash:               "HASH",
		TokenTypeKeywordModel:       "KEYWORD_MODEL",
		TokenTypeIdentifier:         "IDENTIFIER",
		TokenTypeColon:              "COLON",
		TokenTypeSemicolon:          "SEMICOLON",
		TokenTypeDataType:           "DATA_TYPE",
		TokenTypeLeftBrace:          "LEFT_BRACE",
		TokenTypeRightBrace:         "RIGHT_BRACE",
		TokenTypeLeftBracket:        "LEFT_BRACKET",
		TokenTypeRightBracket:       "RIGHT_BRACKET",
		TokenTypeStringLiteral:      "STRING_LITERAL",
		TokenTypeIntLiteral:         "INT_LITERAL",
		TokenTypeFloatLiteral:       "FLOAT_LITERAL",
		TokenTypeWhitespace:         "WHITESPACE",
		TokenTypeLeftParenthesis:    "LEFT_PARENTHESIS",
		TokenTypeRightParenthesis:   "RIGHT_PARENTHESIS",
		TokenTypeDot:                "DOT",
		TokenTypeComma:              "COMMA",
		TokenTypePlus:               "PLUS",
		TokenTypeMinus:              "MINUS",
		TokenTypeMultiply:           "MULTIPLY",
		TokenTypePower:              "POWER",
		TokenTypeDivide:             "DIVIDE",
		TokenTypeModulo:             "MODULO",
		TokenTypeAssignment:         "ASSIGNMENT",
		TokenTypePlusAssignment:     "PLUS_ASSIGNMENT",
		TokenTypeMinusAssignment:    "MINUS_ASSIGNMENT",
		TokenTypeMultiplyAssignment: "MULTIPLY_ASSIGNMENT",
		TokenTypeDivideAssignment:   "DIVIDE_ASSIGNMENT",
		tokenTypeEqualOperator:      "EQUAL_OPERATOR",
		TokenTypeNotEqualOperator:   "NOT_EQUAL_OPERATOR",
		TokenTypeAmpersand:          "AMPERSAND",
		TokenTypeAnd:                "AND",
		TokenTypePipe:               "PIPE",
		TokenTypeOr:                 "OR",
	}

	oneCharSimpleTokens = singleTokenContainer{
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
		'+': TokenTypePlus,
		'-': TokenTypeMinus,
		'*': TokenTypeMultiply,
		'/': TokenTypeDivide,
		'%': TokenTypeModulo,
		'=': TokenTypeAssignment,
		'&': TokenTypeAmpersand,
		'|': TokenTypePipe,
	}

	twoCharsSimpleTokens = doubleTokenContainer{
		'/': {
			'/': TokenTypeComment, // handled separately
			'=': TokenTypeDivideAssignment,
		},
		'+': {
			'=': TokenTypePlusAssignment,
			'+': TokenTypePlusPlus,
		},
		'-': {
			'=': TokenTypeMinusAssignment,
			'-': TokenTypeMinusMinus,
		},
		'*': {
			'=': TokenTypeMultiplyAssignment,
			'*': TokenTypePower,
		},
		'=': {
			'=': tokenTypeEqualOperator,
		},
		'!': {
			'=': TokenTypeNotEqualOperator,
		},
		'&': {
			'&': TokenTypeAnd,
		},
		'|': {
			'|': TokenTypeOr,
		},
	}
)

// maps related to literal value tokens
var (
	// valueTokens is a map of token types that are considered value tokens.
	valueTokens = map[CCLTokenType]bool{
		TokenTypeStringLiteral: true,
		TokenTypeIntLiteral:    true,
		TokenTypeFloatLiteral:  true,
	}

	literalValueTypeInfos = map[CCLTokenType]*cclValues.CCLTypeInfo{
		TokenTypeStringLiteral: cclValues.BuiltInTypeInfos[cclValues.TypeNameString],
		TokenTypeIntLiteral:    cclValues.BuiltInTypeInfos[cclValues.TypeNameInt],
		TokenTypeFloatLiteral:  cclValues.BuiltInTypeInfos[cclValues.TypeNameFloat],
	}
)
