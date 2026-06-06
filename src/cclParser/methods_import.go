package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
)

func (p *CCLAstParser) parseImportDeclAst() (*cclAst.ImportDecl, error) {
	sourcePosition := p.getSourcePosition()
	if err := p.consume(cclLexer.TokenTypeKeywordImport); err != nil {
		return nil, err
	}

	if p.current.Type != cclLexer.TokenTypeStringLiteral {
		return nil, &UnexpectedTokenError{
			Expected:       cclLexer.TokenTypeStringLiteral,
			Actual:         p.current.Type,
			SourcePosition: p.getSourcePosition(),
		}
	}

	importPath := p.current.GetStringLiteral()
	p.advance()

	if err := p.consume(cclLexer.TokenTypeSemicolon); err != nil {
		return nil, err
	}

	return &cclAst.ImportDecl{
		Path:           importPath,
		SourcePosition: sourcePosition,
	}, nil
}
