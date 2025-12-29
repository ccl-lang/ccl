package codeBuilder

import "strings"

type CodeBuilder struct {
	sb          strings.Builder
	indentation int
}
