package cclGenerators

import "errors"

var (
	ErrLanguageNotSupported = errors.New("language not supported")
	ErrMissingCodeContext   = errors.New("missing code context")
)
