package cclLexer

const (
	TokenTypeReservedForFuture CCLTokenType = -1
)

const (
	TokenTypeError            CCLTokenType = iota // error occurred
	TokenTypeEOF                                  // end of file
	TokenTypeComment                              // a comment
	TokenTypeHash                                 // #
	TokenTypeKeywordModel                         // model
	TokenTypeIdentifier                           // an identifier
	TokenTypeColon                                // :
	TokenTypeSemicolon                            // ;
	TokenTypeDataType                             // data type
	TokenTypeLeftBrace                            // {
	TokenTypeRightBrace                           // }
	TokenTypeLeftBracket                          // [
	TokenTypeRightBracket                         // ]
	TokenTypeStringLiteral                        // a string literal
	TokenTypeIntLiteral                           // an integer literal
	TokenTypeFloatLiteral                         // a float literal
	TokenTypeWhitespace                           // a whitespace: ' ', '\t', '\n', '\r'
	TokenTypeLeftParenthesis                      // (
	TokenTypeRightParenthesis                     // )
	TokenTypeDot                                  // .
	TokenTypeComma                                // ,
	TokenTypePlus                                 // +
	TOkenTypePlusPlus                             // ++
	TokenTypeMinus                                // -
	TokenTypeMinusMinus                           // --
	TokenTypeMultiply                             // *
	TokenTypePower                                // **
	TokenTypeDivide                               // /
	TokenTypeModulo                               // %
	TokenTypeAmpersand                            // &
	TokenTypePipe                                 // |
	TokenTypeAnd                                  // &&
	TokenTypeOr                                   // ||

	/* assignment-related operators */

	TokenTypeAssignment         // =
	TokenTypePlusAssignment     // +=
	TokenTypeMinusAssignment    // -=
	TokenTypeMultiplyAssignment // *=
	TokenTypeDivideAssignment   // /=
	tokenTypeEqualOperator      // ==
	TokenTypeNotEqualOperator   // !=
)
