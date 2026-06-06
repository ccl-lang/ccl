package goGenerator

import gValues "github.com/ccl-lang/ccl/src/core/globalValues"

const (
	CurrentLanguage = gValues.LanguageGo
)

const (
	ConstantsFileName = "constants.go"
	MethodsFileName   = "methods.go"
	TypesFileName     = "types.go"
	HelpersFileName   = "helpers.go"
	VarsFileName      = "vars.go"
)

const (
	openImportKey  = "__go_import_group_open"
	closeImportKey = "__go_import_group_close"
)
