package cclLexer

import (
	"fmt"

	"github.com/ccl-lang/ccl/src/core/cclValues"
)

//---------------------------------------------------------

// ToString returns the string representation of the token type.
func (t CCLTokenType) ToString() string {
	return tokenTypesToNames[t]
}

// String returns the string representation of the token type.
func (t CCLTokenType) String() string {
	return tokenTypesToNames[t]
}

//---------------------------------------------------------

// String returns the string representation of the token.
func (t *CCLToken) String() string {
	return t.FormatValueAsString() + " -> " + t.Type.String()
}

// FormatValueAsString returns the formatted value of the token as a string.
func (t *CCLToken) FormatValueAsString() string {
	return fmt.Sprintf("%v", t.value)
}

// GetStringLiteral returns the string literal value of the token.
// If the token is not a string literal, it returns an empty string.
func (t *CCLToken) GetStringLiteral() string {
	if t.Type == TokenTypeStringLiteral {
		return t.value.(string)
	}

	return ""
}

// GetReservedLiteral returns the reserved literal value of the token.
// If the token is not a reserved literal, it returns an empty string.
// Some examples of reserved literals are: true, false, null, nil, self, super, this.
func (t *CCLToken) GetReservedLiteral() string {
	if t.Type == TokenTypeReservedLiteral {
		return t.value.(string)
	}

	return ""
}

// GetComment returns the comment value of the token.
// If the token is not a comment, it returns an empty string.
func (t *CCLToken) GetComment() string {
	if t.Type == TokenTypeComment {
		return t.value.(string)
	}

	return ""
}

// GetIntLiteral returns the integer literal value of the token.
// If the token is not an integer literal, it returns 0.
func (t *CCLToken) GetIntLiteral() int {
	if t.Type == TokenTypeIntLiteral {
		return toIntegerValueT[int](t.value)
	}

	return 0
}

// GetFloatLiteral returns the float literal value of the token.
// If the token is not a float literal, it returns 0.
func (t *CCLToken) GetFloatLiteral() float64 {
	if t.Type == TokenTypeFloatLiteral {
		return t.value.(float64)
	}

	return 0
}

// GetIdentifier returns the identifier value of the token.
// If the token is not an identifier, it returns an empty string.
func (t *CCLToken) GetIdentifier() string {
	if t.Type == TokenTypeIdentifier {
		return t.value.(string)
	}

	return ""
}

// GetDataTypeName returns the built-in data type name if the token is a data type.
// If the token is not a data type, it returns an empty string.
func (t *CCLToken) GetDataTypeName() string {
	if t.Type == TokenTypeDataType {
		name, ok := t.value.(string)
		if ok {
			return name
		}
	}

	return ""
}

// IsTokenValue returns true if the token is a value token.
func (t *CCLToken) IsTokenLiteralValue() bool {
	return valueTokens[t.Type]
}

// IsIdentifier returns true if the token is an identifier.
func (t *CCLToken) IsIdentifier() bool {
	return t.Type == TokenTypeIdentifier
}

// IsReservedLiteral returns true if the token is a reserved literal.
// e.g. true, false, null, nil, self, super, this
func (t *CCLToken) IsReservedLiteral() bool {
	return t.Type == TokenTypeReservedLiteral
}

// IsBuiltinDataType returns true if the token is a built-in data type.
// e.g. int, float, string, bool, etc.
func (t *CCLToken) IsBuiltinDataType() bool {
	return t.Type == TokenTypeDataType
}

// GetLiteralTypeInfo returns type info of the current literal token.
// a literal token is for example: string literal, int literal, float literal, bool literal, etc.
// so for a "hello"" token, it will return the type info for string type.
func (t *CCLToken) GetLiteralTypeInfo(ctx *cclValues.CCLCodeContext) *cclValues.CCLTypeUsage {
	if t.IsTokenLiteralValue() {
		// we have a value token, so we need to return the type info
		// based on the token type
		generator := literalValueTokenToTypeUsage[t.Type]
		if generator == nil {
			return nil
		}
		return generator(ctx)
	} else if t.IsReservedLiteral() {
		// reserved literals also have type info
		generator := reservedLiteralTokenToTypeUsage[t.GetReservedLiteral()]
		if generator == nil {
			return nil
		}
		return generator(ctx)
	}

	return nil
}

// GetBuiltInDataTypeUsage returns the built-in data type usage if the token is a built-in data type.
func (t *CCLToken) GetBuiltInDataTypeUsage(ctx *cclValues.CCLCodeContext) *cclValues.CCLTypeUsage {
	if t.IsBuiltinDataType() {
		nameStr, ok := t.value.(string)
		if !ok {
			return nil
		}

		return ctx.NewBuiltinTypeUsage(nameStr)
	}

	return nil
}

// GetCustomTypeUsage returns the custom type usage if the token is an identifier
// representing a custom type, in the specified ccl code context and namespace.
func (t *CCLToken) GetCustomTypeUsage(
	ctx *cclValues.CCLCodeContext,
	currentNamespace string,
) *cclValues.CCLTypeUsage {
	if t.IsIdentifier() {
		return ctx.NewCustomTypeUsage(&cclValues.SimpleTypeName{
			TypeName:  t.GetIdentifier(),
			Namespace: currentNamespace,
		})
	}

	// maybe some other thing? I don't know
	return nil
}

// GetLiteralValue returns the literal value of the token.
func (t *CCLToken) GetLiteralValue() any {
	if t.IsReservedLiteral() {
		return reservedLiteralTokenToInternalValue[t.value.(string)]
	}

	if !t.IsTokenLiteralValue() {
		// we don't have a value token
		return nil
	}

	// we have a value token, so we need to return the value
	return t.value
}

//---------------------------------------------------------
