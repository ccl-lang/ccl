package cclValues

var (
	keywordNamesToNormalizedValues = map[string]string{
		"model": KeywordNameModel,
	}
)

var (
	reservedLiteralsToNormalizedValues = map[string]string{
		"null":     ReservedLiteralNull,
		"true":     ReservedLiteralTrue,
		"false":    ReservedLiteralFalse,
		"nil":      ReservedLiteralNil,
		"none":     ReservedLiteralNone,
		"self":     ReservedLiteralSelf,
		"super":    ReservedLiteralSuper,
		"this":     ReservedLiteralThis,
		"nan":      ReservedLiteralNaN,
		"infinity": ReservedLiteralInfinity,
	}
)
