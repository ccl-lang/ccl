package cclValues

// IsTypeName checks if the given value is a valid type name.
// Please note that this function only checks if the type is a valid
// built-in type name, this package has no responsibility for checking
// custom type names.
func IsTypeName(value string) bool {
	_, exists := typeNamesToNormalizedValues[value]
	return exists
}

// IsKeywordName checks if the given value is a valid keyword name.
func IsKeywordName(value string) bool {
	_, exists := keywordNamesToNormalizedValues[value]
	return exists
}

// GetNormalizedTypeName returns the normalized value of the given type name.
// If the given value is not a valid type name, an empty string will be returned.
func GetNormalizedTypeName(value string) string {
	return typeNamesToNormalizedValues[value]
}

// GetNormalizedKeywordName returns the normalized value of the given keyword name.
// If the given value is not a valid keyword name, an empty string will be returned.
func GetNormalizedKeywordName(value string) string {
	return keywordNamesToNormalizedValues[value]
}
