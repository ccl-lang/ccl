package cclParser

import (
	"strings"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
)

func (p *CCLAstParser) parseNamespaceDeclAst() (string, error) {
	if err := p.consume(cclLexer.TokenTypeKeywordNamespace); err != nil {
		return "", err
	}

	parts := []string{}
	for {
		if p.current.Type != cclLexer.TokenTypeIdentifier {
			return "", &UnexpectedTokenError{
				Expected:       cclLexer.TokenTypeIdentifier,
				Actual:         p.current.Type,
				SourcePosition: p.getSourcePosition(),
			}
		}

		parts = append(parts, p.current.GetIdentifier())
		p.advance()

		if !p.isCurrentType(cclLexer.TokenTypeDot) {
			break
		}

		p.advance()
	}

	if err := p.consume(cclLexer.TokenTypeSemicolon); err != nil {
		return "", err
	}

	return strings.Join(parts, "."), nil
}
