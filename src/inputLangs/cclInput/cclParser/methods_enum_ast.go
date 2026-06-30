package cclParser

import (
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser/cclLexer"
)

func (p *CCLAstParser) parseEnumDeclAst(currentNamespace string) (*cclAst.EnumDecl, error) {
	if err := p.consume(cclLexer.TokenTypeKeywordEnum); err != nil {
		return nil, err
	}

	if !p.isCurrentType(cclLexer.TokenTypeIdentifier) {
		return nil, &UnexpectedTokenError{
			Expected:       cclLexer.TokenTypeIdentifier,
			Actual:         p.current.Type,
			SourcePosition: p.getSourcePosition(),
		}
	}

	enumDecl := &cclAst.EnumDecl{
		Name:           p.current.GetIdentifier(),
		Namespace:      currentNamespace,
		SourcePosition: p.getSourcePosition(),
	}
	p.advance()

	if p.isCurrentType(cclLexer.TokenTypeColon) {
		p.advance()
		baseType, err := p.parseTypeExpressionUntil(cclLexer.TokenTypeLeftBrace)
		if err != nil {
			return nil, err
		}
		enumDecl.BaseType = baseType
	}

	if err := p.consume(cclLexer.TokenTypeLeftBrace); err != nil {
		return nil, err
	}

	for !p.IsAtEnd() {
		if p.isCurrentType(cclLexer.TokenTypeComment, cclLexer.TokenTypeComma, cclLexer.TokenTypeSemicolon) {
			p.advance()
			continue
		}

		if p.isCurrentType(cclLexer.TokenTypeRightBrace) {
			p.advance()
			return enumDecl, nil
		}

		member, err := p.parseEnumMemberDeclAst()
		if err != nil {
			return nil, err
		}
		enumDecl.Members = append(enumDecl.Members, member)
	}

	return nil, &UnexpectedEOFError{
		SourcePosition: p.getSourcePosition(),
	}
}

func (p *CCLAstParser) parseEnumMemberDeclAst() (*cclAst.EnumMemberDecl, error) {
	if !p.isCurrentType(cclLexer.TokenTypeIdentifier) {
		return nil, &UnexpectedTokenError{
			Expected:       cclLexer.TokenTypeIdentifier,
			Actual:         p.current.Type,
			SourcePosition: p.getSourcePosition(),
		}
	}

	member := &cclAst.EnumMemberDecl{
		Name:           p.current.GetIdentifier(),
		SourcePosition: p.getSourcePosition(),
	}
	p.advance()

	if !p.isCurrentType(cclLexer.TokenTypeAssignment) {
		return member, nil
	}

	p.advance()
	value, err := p.parseEnumMemberValue()
	if err != nil {
		return nil, err
	}
	member.Value = &value
	return member, nil
}

func (p *CCLAstParser) parseEnumMemberValue() (int64, error) {
	isNegative := false
	if p.isCurrentType(cclLexer.TokenTypeMinus) {
		isNegative = true
		p.advance()
	}

	if !p.isCurrentType(cclLexer.TokenTypeIntLiteral) {
		return 0, &UnexpectedTokenError{
			Expected:       cclLexer.TokenTypeIntLiteral,
			Actual:         p.current.Type,
			SourcePosition: p.getSourcePosition(),
		}
	}

	value := p.current.GetInt64Literal()
	if isNegative {
		value = -value
	}
	p.advance()
	return value, nil
}
