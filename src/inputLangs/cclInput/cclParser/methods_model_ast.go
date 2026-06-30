package cclParser

import (
	"github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser/cclLexer"
)

func (p *CCLAstParser) parseModelDeclAst(currentNamespace string) (*cclAst.ModelDecl, error) {
	if err := p.consume(cclLexer.TokenTypeKeywordModel); err != nil {
		return nil, err
	}

	if !p.isCurrentType(cclLexer.TokenTypeIdentifier) {
		return nil, &UnexpectedTokenError{
			Expected:       cclLexer.TokenTypeIdentifier,
			Actual:         p.current.Type,
			SourcePosition: p.getSourcePosition(),
		}
	}

	modelName := p.current.GetIdentifier()
	modelPosition := p.getSourcePosition()
	p.advance()

	openBraceCount := 0
	currentPendingAttributes := []*cclAst.AttributeNode{}
	fields := []*cclAst.FieldDecl{}
	enums := []*cclAst.EnumDecl{}

	for !p.IsAtEnd() {
		if p.isCurrentComment() {
			p.advance()
			continue
		}

		if p.isCurrentType(cclLexer.TokenTypeLeftBrace) {
			if openBraceCount > 0 {
				return nil, &UnexpectedTokenError{
					Expected:       cclLexer.TokenTypeRightBrace,
					Actual:         p.current.Type,
					SourcePosition: p.getSourcePosition(),
				}
			}

			openBraceCount++
			p.advance()
			continue
		} else if p.isCurrentType(cclLexer.TokenTypeRightBrace) {
			openBraceCount--
			if openBraceCount < 0 {
				return nil, &UnexpectedTokenError{
					Expected:       cclLexer.TokenTypeEOF,
					Actual:         p.current.Type,
					SourcePosition: p.getSourcePosition(),
				}
			}

			p.advance()

			if openBraceCount == 0 {
				break
			}
			continue
		}

		if p.isCurrentAttribute() {
			attrs, err := p.parseAttributeNodes()
			if err != nil {
				return nil, err
			}
			currentPendingAttributes = append(currentPendingAttributes, attrs...)
			continue
		}

		if p.isCurrentType(cclLexer.TokenTypeKeywordEnum) {
			enumDecl, err := p.parseEnumDeclAst(currentNamespace)
			if err != nil {
				return nil, err
			}

			if len(currentPendingAttributes) > 0 {
				enumDecl.Attributes = append(enumDecl.Attributes, currentPendingAttributes...)
				currentPendingAttributes = nil
			}

			enums = append(enums, enumDecl)
			continue
		}

		if p.isCurrentTokenFieldOfModel() {
			field, err := p.parseModelFieldAst()
			if err != nil {
				return nil, err
			}

			if len(currentPendingAttributes) > 0 {
				field.Attributes = append(field.Attributes, currentPendingAttributes...)
				currentPendingAttributes = nil
			}

			fields = append(fields, field)
			continue
		}

		return nil, &UnexpectedTokenError{
			Expected:       cclLexer.TokenTypeIdentifier,
			Actual:         p.current.Type,
			SourcePosition: p.getSourcePosition(),
		}
	}

	return &cclAst.ModelDecl{
		Name:           modelName,
		Namespace:      currentNamespace,
		Fields:         fields,
		Enums:          enums,
		SourcePosition: modelPosition,
	}, nil
}

func (p *CCLAstParser) parseModelFieldAst() (*cclAst.FieldDecl, error) {
	if !p.isCurrentType(cclLexer.TokenTypeIdentifier) {
		return nil, &InvalidSyntaxError{
			Language:       globalValues.LanguageCCL,
			SourcePosition: p.getSourcePosition(),
		}
	}

	field := &cclAst.FieldDecl{
		Name:           p.current.GetIdentifier(),
		SourcePosition: p.getSourcePosition(),
	}
	p.advance()

	if err := p.consume(cclLexer.TokenTypeColon); err != nil {
		return nil, err
	}

	typeExpr, err := p.parseTypeExpressionUntil(
		cclLexer.TokenTypeAssignment,
		cclLexer.TokenTypeSemicolon,
	)
	if err != nil {
		return nil, err
	}
	field.Type = typeExpr

	if p.isCurrentType(cclLexer.TokenTypeAssignment) {
		p.advance()
		valueExpr, err := p.parseValueExpressionUntil(cclLexer.TokenTypeSemicolon)
		if err != nil {
			return nil, err
		}
		field.Value = valueExpr
	}

	if err := p.consume(cclLexer.TokenTypeSemicolon); err != nil {
		return nil, err
	}

	return field, nil
}

// isCurrentTokenFieldOfModel returns true only when the current token is at the
// beginning of field of a model.
func (p *CCLAstParser) isCurrentTokenFieldOfModel() bool {
	// case1: identifier followed by colon
	if p.isCurrentType(cclLexer.TokenTypeIdentifier) {
		// lookahead for colon
		if p.isNextType(cclLexer.TokenTypeColon) {
			return true
		}
	}

	// maybe add more cases in future
	return false
}
