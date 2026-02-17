package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclAst"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

func (p *CCLParser) ParseAsAST() (*cclAst.CCLFileAST, error) {
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
			if p.isAttributeAt(p.pos + 1) {
				globalAttr, err := p.parseGlobalAttributeNode()
				if err != nil {
					return nil, err
				}
				fileAst.GlobalAttributes = append(fileAst.GlobalAttributes, globalAttr)
			}
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
