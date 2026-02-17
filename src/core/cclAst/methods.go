package cclAst

import "github.com/ccl-lang/ccl/src/core/cclUtils"

// GetTypeExpressionKind returns the kind of this simple type expression.
func (s *SimpleTypeExpression) GetTypeExpressionKind() TypeExpressionKind {
	return TypeExpressionKindSimpleName
}

// GetSourcePosition returns the source position of this simple type expression.
func (s *SimpleTypeExpression) GetSourcePosition() *cclUtils.SourceCodePosition {
	if s == nil {
		return nil
	}
	return s.SourcePosition
}

// GetTypeExpressionKind returns the kind of this array type expression.
func (a *ArrayTypeExpression) GetTypeExpressionKind() TypeExpressionKind {
	return TypeExpressionKindArray
}

// GetSourcePosition returns the source position of this array type expression.
func (a *ArrayTypeExpression) GetSourcePosition() *cclUtils.SourceCodePosition {
	if a == nil {
		return nil
	}
	return a.SourcePosition
}
