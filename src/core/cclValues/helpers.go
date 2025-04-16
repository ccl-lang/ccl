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

// NewTypeInfo creates a new type info.
func NewTypeInfo(name string) *CCLTypeInfo {
	if IsTypeName(name) {
		return &CCLTypeInfo{
			name:      name,
			typeFlags: TypeFlagBuiltIn,
		}
	}

	return &CCLTypeInfo{
		name:      name,
		typeFlags: 0b0,
	}
}

// NewTypeInfoWithOperators creates a new type info with operators.
// TODO: refactor this to handle operators in a better way.
func NewTypeInfoWithOperators(name string, operators ...string) *CCLTypeInfo {
	var flags cclTypeFlag = 0b0
	if IsTypeName(name) {
		flags |= TypeFlagBuiltIn
	}

	for _, currentOperator := range operators {
		if currentOperator == "[]" {
			flags |= TypeFlagArray
		}
	}

	return &CCLTypeInfo{
		name:      name,
		typeFlags: flags,
	}
}

func NewPointerTypeInfo(targetType *CCLTypeInfo) *CCLTypeInfo {
	return &CCLTypeInfo{
		name:           targetType.name,
		typeFlags:      TypeFlagPointer,
		underlyingType: targetType,
	}
}

// NewArrayTypeInfo creates a new array type info.
func NewArrayTypeInfo(name string) *CCLTypeInfo {
	if IsTypeName(name) {
		return &CCLTypeInfo{
			name:      name,
			typeFlags: TypeFlagBuiltIn | TypeFlagArray,
		}
	}

	return &CCLTypeInfo{
		name:      name,
		typeFlags: TypeFlagArray,
	}
}

//---------------------------------------------------------

// GetGlobalVariable returns the global variable with the given name.
func GetGlobalVariable(name string) *VariableDefinition {
	if variable, exists := cclAutomaticVariables[name]; exists {
		return variable
	}

	if variable, exists := cclGlobalVariables[name]; exists {
		return variable
	}

	return nil
}

//---------------------------------------------------------
