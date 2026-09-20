package cclParser

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
)

type importGraphResolver struct {
	targetLanguage gValues.LanguageType
	visitedFiles   map[string]bool
	activeFiles    map[string]bool
	fileStack      []string
	fileAsts       map[string]*cclAst.CCLFileAST
	orderedAsts    []*cclAst.CCLFileAST
}
