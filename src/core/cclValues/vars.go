package cclValues

var (
	typeNamesToNormalizedValues = map[string]string{
		"string": TypeNameString,
		"String": TypeNameString,
		"bytes":  TypeNameBytes,
		"int":    TypeNameInt,
		"int8":   TypeNameInt8,
		// "byte":     TypeNameByte,
		"int16":    TypeNameInt16,
		"int32":    TypeNameInt32,
		"int64":    TypeNameInt64,
		"uint":     TypeNameUint,
		"uint8":    TypeNameUint8,
		"uint16":   TypeNameUint16,
		"uint32":   TypeNameUint32,
		"uint64":   TypeNameUint64,
		"float":    TypeNameFloat,
		"float32":  TypeNameFloat32,
		"float64":  TypeNameFloat64,
		"bool":     TypeNameBool,
		"datetime": TypeNameDateTime,
	}

	keywordNamesToNormalizedValues = map[string]string{
		"model": KeywordNameModel,
	}
)

// Type infos for built-in types.
var (
	BuiltInTypeInfos = func() map[string]*CCLTypeInfo {
		typeInfos := make(map[string]*CCLTypeInfo)
		for _, normalizedName := range typeNamesToNormalizedValues {
			typeInfos[normalizedName] = &CCLTypeInfo{
				name:      normalizedName,
				typeFlags: TypeFlagBuiltIn,
			}
		}
		return typeInfos
	}()
)

// All the global variables used in the source code.
var (
	cclGlobalVariables = make(map[string]*VariableDefinition)

	cclAutomaticVariables = func() map[string]*VariableDefinition {
		return make(map[string]*VariableDefinition)
	}()
)
