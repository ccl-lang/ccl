package cclValues

// IsKeywordName checks if the given value is a valid keyword name.
func IsKeywordName(value string) bool {
	_, exists := keywordNamesToNormalizedValues[value]
	return exists
}

// GetNormalizedTypeName returns the normalized value of the given type name.
// If the given value is not a valid type name, an empty string will be returned.
func GetNormalizedTypeName(value string) string {
	return builtinNamesToNormalizedValues[value]
}

// GetNormalizedKeywordName returns the normalized value of the given keyword name.
// If the given value is not a valid keyword name, an empty string will be returned.
func GetNormalizedKeywordName(value string) string {
	return keywordNamesToNormalizedValues[value]
}

// NewTypeInfoWithOperators_OLD creates a new type info with operators.
// TODO: refactor this to handle operators in a better way.
func NewTypeInfoWithOperators_OLD(name string, operators ...string) *CCLTypeDefinition {
	var flags cclTypeFlag = 0b0
	if IsBuiltinTypeName(name) {
		flags |= TypeFlagBuiltIn
	}

	for _, currentOperator := range operators {
		if currentOperator == "[]" {
			flags |= TypeFlagArray
		}
	}

	return &CCLTypeDefinition{
		name:      name,
		typeFlags: flags,
	}
}

//---------------------------------------------------------
