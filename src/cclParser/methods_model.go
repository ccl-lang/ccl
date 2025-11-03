package cclParser

import (
	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	"github.com/ccl-lang/ccl/src/core/globalValues"
)

// ParseModelDefinition parses a model definition from the current position in the source code.
func (p *CCLParser) ParseModelDefinition(currentNamespace string) (*cclValues.ModelDefinition, error) {
	// TODO: optionally add public or private keyword here in future
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
	p.advance()

	openBraceCount := 0
	currentPendingAttributes := []*cclValues.AttributeUsageInfo{}
	myFields := []*cclValues.ModelFieldDefinition{}

	for !p.IsAtEnd() {
		if p.isCurrentComment() {
			p.advance()
			continue
		}

		if p.isCurrentType(cclLexer.TokenTypeLeftBrace) {
			// maybe reconsider this in future?
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
				// end of model definition
				break
			}

			// can we do anything else here?
			continue
		}

		if p.isCurrentAttribute() {
			// we can have multiple attributes before an entity inside of the model
			// e.g. a field, a method, etc...
			allAttributes, err := p.ParseAttributes()
			if err != nil {
				return nil, err
			}

			// since we don't want to make our parser too complex, we will just set
			// pending attributes here and let other parts of parser handle this.
			// E.g. if after this, we get a field, we will set the attributes to the field;
			// or if we get a method, we will set the attributes to the method, etc...
			currentPendingAttributes = append(currentPendingAttributes, allAttributes...)
			continue
		}

		if p.isCurrentTokenFieldOfModel() {
			currentField, err := p.ParseModelField(currentNamespace)
			if err != nil {
				return nil, err
			}

			if len(currentPendingAttributes) > 0 {
				currentField.Attributes = append(currentField.Attributes, currentPendingAttributes...)
				currentPendingAttributes = nil
			}

			myFields = append(myFields, currentField)
			continue
		}

		// TODO: in future we can have methods inside of the model and other stuff
		return nil, &UnexpectedTokenError{
			Expected:       cclLexer.TokenTypeIdentifier,
			Actual:         p.current.Type,
			SourcePosition: p.getSourcePosition(),
		}
	}

	return &cclValues.ModelDefinition{
		Name:    modelName,
		ModelId: p.codeDefinition.GetNextModelId(),
		Fields:  myFields,
	}, nil
}

func (p *CCLParser) ParseModelField(currentNamespace string) (*cclValues.ModelFieldDefinition, error) {
	theField := &cclValues.ModelFieldDefinition{}
	gotColon := false
	gotAssignment := false
	for {
		isDataType := p.isCurrentType(cclLexer.TokenTypeDataType)
		if p.IsAtEnd() {
			// we still haven't got any semicolon...
			return nil, &UnexpectedEOFError{
				SourcePosition: p.getSourcePosition(),
			}
		}

		if p.isCurrentComment() {
			p.advance()
			continue
		}

		if p.isCurrentType(cclLexer.TokenTypeSemicolon) {
			if (!gotColon && !gotAssignment) || theField.Type == nil {
				// return error here
				return nil, &InvalidSyntaxError{
					Language:       globalValues.LanguageCCL,
					SourcePosition: p.getSourcePosition(),
				}
			}

			p.advance()
			// validate other stuff here in future
			break
		}

		if isDataType || p.isCurrentType(cclLexer.TokenTypeIdentifier) {
			// first identifier has to be the name of the field
			if theField.Name == "" {
				if isDataType {
					return nil, p.ErrInvalidSyntax("Cannot use built-in data-types as field names")
				}

				theField.Name = p.current.GetIdentifier()
				p.advance()
				continue
			}

			// this is our second (or more) time getting an identifier here.
			// we have three possibilities
			// 1. if we got colon before: it means this is a type name
			// 2. if we got assignment before: it means this has to be a constant or
			// 		an automatic variable.
			// 3. if we got both colo and assignment before: we will still accept it
			// 		and just ignore the colon.
			// if we haven't got any of these...it means it's a syntax error, we don't
			// accept go-style syntax for fields in ccl.
			if !gotColon && !gotAssignment {
				return nil, &InvalidSyntaxError{
					Language:       globalValues.LanguageCCL,
					SourcePosition: p.getSourcePosition(),
				}
			}

			if theField.Type == nil {
				if gotColon && !gotAssignment {
					fieldType, err := p.parseCurrentTypeUsage(currentNamespace)
					if err != nil {
						return nil, err
					} else if fieldType == nil {
						return nil, p.ErrInvalidSyntax("Invalid type usage in field definition")
					}
					theField.ChangeValueType(fieldType)
					continue
				} else if gotAssignment {
					// We have a variable usage here
					// and the parameter name has previously been specified
					if isDataType {
						return nil, p.ErrInvalidSyntax(
							"Don't use built-in type names in field assignments. " +
								"Use generics for that.")
					}

					targetIdentifier := p.current.GetIdentifier()
					targetVariable := cclValues.GetGlobalVariable(targetIdentifier)
					if targetVariable == nil {
						return nil, &UndefinedIdentifierError{
							TargetIdentifier: targetIdentifier,
							Language:         globalValues.LanguageCCL,
							SourcePosition:   p.getSourcePosition(),
						}
					}
					if targetVariable.IsAutomatic() {
						theField.ChangeValueType(cclValues.NewPointerTypeUsage(
							targetVariable.Type,
						))
						theField.ChangeValue(&cclValues.VariableUsageInstance{
							Name:       targetIdentifier,
							Definition: targetVariable,
						})
					} else {
						// since the variable is not an automatic variable, we don't
						// have to *point* to it.
						theField.ChangeValueType(targetVariable.Type)
						theField.ChangeValue(targetVariable.GetValue()) // copy the value
					}
					p.advance()
					continue
				}

				// just in case
				return nil, p.ErrInvalidSyntax("Impossible scenario reached")
			}

			// maybe in the future we want to do something else here?
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

	return theField, nil
}

// isCurrentTokenFieldOfModel returns true only when the current token is at the
// beginning of field of a model.
func (p *CCLParser) isCurrentTokenFieldOfModel() bool {
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
