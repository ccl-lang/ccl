package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/globalValues"
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

		if p.isCurrentTokenFieldOfModel() {
			field, err := p.parseModelFieldAst(currentNamespace)
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
		SourcePosition: modelPosition,
	}, nil
}

func (p *CCLAstParser) parseModelFieldAst(currentNamespace string) (*cclAst.FieldDecl, error) {
	field := &cclAst.FieldDecl{}
	gotColon := false
	gotAssignment := false

	for {
		isDataType := p.isCurrentType(cclLexer.TokenTypeDataType)
		if p.IsAtEnd() {
			return nil, &UnexpectedEOFError{
				SourcePosition: p.getSourcePosition(),
			}
		}

		if p.isCurrentComment() {
			p.advance()
			continue
		}

		if p.isCurrentType(cclLexer.TokenTypeSemicolon) {
			if (!gotColon && !gotAssignment) || (field.Type == nil && field.Value == nil) {
				return nil, &InvalidSyntaxError{
					Language:       globalValues.LanguageCCL,
					SourcePosition: p.getSourcePosition(),
				}
			}

			p.advance()
			break
		}

		if isDataType || p.isCurrentType(cclLexer.TokenTypeIdentifier) {
			if field.Name == "" {
				if isDataType {
					return nil, p.ErrInvalidSyntax("Cannot use built-in data-types as field names")
				}

				field.Name = p.current.GetIdentifier()
				field.SourcePosition = p.getSourcePosition()
				p.advance()
				continue
			}

			if !gotColon && !gotAssignment {
				return nil, &InvalidSyntaxError{
					Language:       globalValues.LanguageCCL,
					SourcePosition: p.getSourcePosition(),
				}
			}

			if field.Type == nil && field.Value == nil {
				if gotColon && !gotAssignment {
					typeExpr, err := p.parseCurrentTypeExpression()
					if err != nil {
						return nil, err
					} else if typeExpr == nil {
						return nil, p.ErrInvalidSyntax("Invalid type usage in field definition")
					}
					field.Type = typeExpr
					continue
				} else if gotAssignment {
					if isDataType {
						return nil, p.ErrInvalidSyntax(
							"Don't use built-in type names in field assignments. " +
								"Use generics for that.")
					}

					if !p.current.IsIdentifier() {
						return nil, p.ErrInvalidSyntax("Expected identifier in field assignment")
					}

					field.Value = &cclAst.IdentifierValueExpression{
						Name:           p.current.GetIdentifier(),
						SourcePosition: p.getSourcePosition(),
					}
					p.advance()
					continue
				}

				return nil, p.ErrInvalidSyntax("Impossible scenario reached")
			}

			return nil, p.ErrInvalidSyntax("")
		}

		if p.isCurrentType(cclLexer.TokenTypeColon) {
			gotColon = true
			p.advance()
			continue
		}

		if p.isCurrentType(cclLexer.TokenTypeAssignment) {
			gotAssignment = true
			p.advance()
			continue
		}

		return nil, p.ErrInvalidSyntax("")
	}

	return field, nil
}
