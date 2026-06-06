package cclValues

// IsBuiltinTypeName checks if the given value is a valid type name.
// Please note that this function only checks if the type is a valid
// built-in type name, this package has no responsibility for checking
// custom type names.
func IsBuiltinTypeName(value string) bool {
	_, exists := builtinNamesToNormalizedValues[value]
	return exists
}
