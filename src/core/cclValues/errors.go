package cclValues

import "errors"

var (
	ErrGenericParamCantBeNil = errors.New("generic-param cannot be nil")
	ErrCircularGenericType   = errors.New("generic-param type can't be same as current type")
)
