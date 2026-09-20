package cclSanitizer

import (
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

// ResolveSkipFileAttribute validates the directive before the parser skips a file.
func ResolveSkipFileAttribute(node cclAst.AttributeNodeBase) (*cclValues.AttributeUsageInfo, error) {
	scopedAttr, isScoped := node.(*cclAst.GlobalAttributeNode)
	if !isScoped || scopedAttr.Scope != cclAst.AttributeScopeFile {
		return nil, &cclErrors.InvalidAttributeUsageError{
			AttrName:       cclAttr.AttrSkipFile,
			Message:        "is file-only; use #file:[$gd:SkipFile()] to skip a file for GDScript",
			SourcePosition: node.GetSourcePosition(),
		}
	}
	if len(node.GetAttributeParams()) != 0 {
		return nil, &cclErrors.InvalidAttributeUsageError{
			AttrName:       cclAttr.AttrSkipFile,
			Message:        "does not accept parameters; select languages before the attribute name, e.g. #file:[$gd:SkipFile()]",
			SourcePosition: node.GetSourcePosition(),
		}
	}

	languages, err := resolveAttributeLanguages(node.GetAttributeLanguages(), node.GetSourcePosition())
	if err != nil {
		return nil, err
	}
	return &cclValues.AttributeUsageInfo{
		Name:           cclAttr.AttrSkipFile,
		Languages:      languages,
		SourcePosition: node.GetSourcePosition(),
	}, nil
}
