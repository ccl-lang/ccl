package cclParser

import (
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclAst"
)

func parseAttributeScope(scope string) cclAst.AttributeScope {
	switch strings.ToLower(scope) {
	case "", "global":
		return cclAst.AttributeScopeGlobal
	case "file":
		return cclAst.AttributeScopeFile
	case "namespace":
		return cclAst.AttributeScopeNamespace
	default:
		return cclAst.AttributeScopeUnknown
	}
}
