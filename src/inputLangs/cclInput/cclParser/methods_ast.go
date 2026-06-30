package cclParser

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser/cclLexer"
)

func (p *CCLAstParser) ParseAsAST() (*cclAst.CCLFileAST, error) {
	if err := p.initializeParsing(); err != nil {
		return nil, err
	}

	fileAst := &cclAst.CCLFileAST{
		FilePath:  p.Options.SourceFilePath,
		Namespace: gValues.DefaultMainNamespace,
	}

	var currentPendingAttributes []*cclAst.AttributeNode
	currentNamespace := gValues.DefaultMainNamespace

	for !p.IsAtEnd() {
		if p.current.Type == cclLexer.TokenTypeHash {
			scopedAttr, err := p.parseGlobalAttributeNode()
			if err != nil {
				return nil, err
			}

			switch scopedAttr.Scope {
			case cclAst.AttributeScopeGlobal:
				fileAst.GlobalAttributes = append(fileAst.GlobalAttributes, scopedAttr)
			case cclAst.AttributeScopeFile:
				fileAst.FileAttributes = append(fileAst.FileAttributes, scopedAttr)
			case cclAst.AttributeScopeNamespace:
				scopedAttr.Namespace = currentNamespace
				fileAst.NamespaceAttributes = append(fileAst.NamespaceAttributes, scopedAttr)
			default:
				return nil, &InvalidSyntaxError{
					Language:       gValues.LanguageCCL,
					HintMessage:    "Expected attribute scope to be file, global, or namespace.",
					SourcePosition: scopedAttr.SourcePosition,
				}
			}
			continue
		}

		if p.current.Type == cclLexer.TokenTypeKeywordNamespace {
			if len(currentPendingAttributes) > 0 {
				lastAttr := currentPendingAttributes[len(currentPendingAttributes)-1]
				return nil, &InvalidAttributeUsageError{
					SourcePosition: lastAttr.SourcePosition,
				}
			}

			namespace, err := p.parseNamespaceDeclAst()
			if err != nil {
				return nil, err
			}
			currentNamespace = namespace
			fileAst.Namespace = namespace
			continue
		}

		if p.current.Type == cclLexer.TokenTypeKeywordImport {
			if len(currentPendingAttributes) > 0 {
				lastAttr := currentPendingAttributes[len(currentPendingAttributes)-1]
				return nil, &InvalidAttributeUsageError{
					SourcePosition: lastAttr.SourcePosition,
				}
			}

			importDecl, err := p.parseImportDeclAst()
			if err != nil {
				return nil, err
			}

			fileAst.Imports = append(fileAst.Imports, importDecl)
			continue
		}

		if p.isCurrentAttribute() {
			afterAttribute := p.peekAfterAttribute()
			if afterAttribute == cclLexer.TokenTypeEOF {
				return nil, &UnexpectedEndOfAttributeError{
					SourcePosition: p.getSourcePosition(),
				}
			}

			attrs, err := p.parseAttributeNodes()
			if err != nil {
				return nil, err
			}

			currentPendingAttributes = append(currentPendingAttributes, attrs...)
			continue
		}

		if p.current.Type == cclLexer.TokenTypeKeywordModel {
			model, err := p.parseModelDeclAst(currentNamespace)
			if err != nil {
				return nil, err
			}

			if len(currentPendingAttributes) > 0 {
				model.Attributes = append(model.Attributes, currentPendingAttributes...)
				currentPendingAttributes = nil
			}

			fileAst.Models = append(fileAst.Models, model)
			continue
		}

		if p.current.Type == cclLexer.TokenTypeKeywordEnum {
			enumDecl, err := p.parseEnumDeclAst(currentNamespace)
			if err != nil {
				return nil, err
			}

			if len(currentPendingAttributes) > 0 {
				enumDecl.Attributes = append(enumDecl.Attributes, currentPendingAttributes...)
				currentPendingAttributes = nil
			}

			fileAst.Enums = append(fileAst.Enums, enumDecl)
			continue
		}

		if p.current.Type == cclLexer.TokenTypeComment {
			p.advance()
			continue
		}

		return nil, &UnexpectedTokenError{
			Expected:       cclLexer.TokenTypeKeywordModel,
			Actual:         p.current.Type,
			SourcePosition: p.getSourcePosition(),
		}
	}

	if len(currentPendingAttributes) > 0 {
		lastAttr := currentPendingAttributes[len(currentPendingAttributes)-1]
		return nil, &InvalidAttributeUsageError{
			SourcePosition: lastAttr.SourcePosition,
		}
	}

	return fileAst, nil
}
