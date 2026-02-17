package cclAst

import "github.com/ccl-lang/ccl/src/core/cclUtils"

// GetAttributeName returns the name of the attribute.
func (a *AttributeNode) GetAttributeName() string {
	if a == nil {
		return ""
	}
	return a.Name
}

// GetAttributeParams returns the parameters of the attribute.
func (a *AttributeNode) GetAttributeParams() []*AttributeParamNode {
	if a == nil {
		return nil
	}
	return a.Params
}

// GetSourcePosition returns the source position of the attribute.
func (a *AttributeNode) GetSourcePosition() *cclUtils.SourceCodePosition {
	if a == nil {
		return nil
	}
	return a.SourcePosition
}

// GetAttributeName returns the name of the global attribute.
func (g *GlobalAttributeNode) GetAttributeName() string {
	if g == nil {
		return ""
	}
	return g.Name
}

// GetAttributeParams returns the parameters of the global attribute.
func (g *GlobalAttributeNode) GetAttributeParams() []*AttributeParamNode {
	if g == nil {
		return nil
	}
	return g.Params
}

// GetSourcePosition returns the source position of the global attribute.
func (g *GlobalAttributeNode) GetSourcePosition() *cclUtils.SourceCodePosition {
	if g == nil {
		return nil
	}
	return g.SourcePosition
}

// GetAttributeValueKind returns the kind of this literal value expression.
func (l *LiteralValueExpression) GetAttributeValueKind() AttributeValueKind {
	return AttributeValueKindLiteral
}

// GetSourcePosition returns the source position of this literal value expression.
func (l *LiteralValueExpression) GetSourcePosition() *cclUtils.SourceCodePosition {
	if l == nil {
		return nil
	}
	return l.SourcePosition
}

// GetAttributeValueKind returns the kind of this identifier value expression.
func (i *IdentifierValueExpression) GetAttributeValueKind() AttributeValueKind {
	return AttributeValueKindIdentifier
}

// GetSourcePosition returns the source position of this identifier value expression.
func (i *IdentifierValueExpression) GetSourcePosition() *cclUtils.SourceCodePosition {
	if i == nil {
		return nil
	}
	return i.SourcePosition
}
