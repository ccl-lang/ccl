package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclAst"
)

func (p *CCLAstParser) parseGlobalAttributeNode() (*cclAst.GlobalAttributeNode, error) {
	if err := p.consume(cclLexer.TokenTypeHash); err != nil {
		return nil, err
	}

	attrNode, err := p.parseSingleAttributeNode()
	if err != nil {
		return nil, err
	}

	return &cclAst.GlobalAttributeNode{
		Name:           attrNode.Name,
		Params:         attrNode.Params,
		SourcePosition: attrNode.SourcePosition,
	}, nil
}

func (p *CCLAstParser) parseAttributeNodes() ([]*cclAst.AttributeNode, error) {
	allAttributes := []*cclAst.AttributeNode{}
	for !p.IsAtEnd() && p.isCurrentAttribute() {
		if p.isCurrentType(cclLexer.TokenTypeComment) {
			p.advance()
			continue
		}

		if p.isCurrentType(cclLexer.TokenTypeLeftBracket) {
			attributeNode, err := p.parseSingleAttributeNode()
			if err != nil {
				return nil, err
			}

			if attributeNode != nil {
				allAttributes = append(allAttributes, attributeNode)
			}
			continue
		}
	}

	return allAttributes, nil
}
