package cclValues

import "errors"

var (
	ErrGenericParamCantBeNil = errors.New("generic-param cannot be nil")
	ErrCircularGenericType   = errors.New("generic-param type can't be same as current type")
)

const (
	StrErrCannotOverrideBuiltInType = "cannot override built-in type: %s in namespace %s"
	StrErrTypeAlreadyDefined        = "type already defined: %s in namespace %s"
)
