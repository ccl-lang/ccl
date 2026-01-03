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
	TokenTypeKeywordRequired                      // required
	TokenTypeKeywordPublic                        // public
	TokenTypeKeywordPrivate                       // private
	TokenTypeKeywordInternal                      // internal
	TokenTypeKeywordProtected                     // protected
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
	TokenTypeReservedLiteral                      // a reserved literal: true, false, null, nil, self, super, this, etc...
	TokenTypeWhitespace                           // a whitespace: ' ', '\t', '\n', '\r'
	TokenTypeLeftParenthesis                      // (
	TokenTypeRightParenthesis                     // )
	TokenTypeDot                                  // .
	TokenTypeComma                                // ,
	TokenTypePlus                                 // +
	TokenTypePlusPlus                             // ++
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
