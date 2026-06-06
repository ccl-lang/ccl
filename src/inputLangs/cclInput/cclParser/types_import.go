package cclParser

import "github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"

type importGraphResolver struct {
	visitedFiles map[string]bool
	activeFiles  map[string]bool
	fileStack    []string
	fileAsts     map[string]*cclAst.CCLFileAST
	orderedAsts  []*cclAst.CCLFileAST
}
