package sourcePos

import (
	"fmt"
	"strings"
)

// FormatError formats an error message with source position context.
// It expands tabs to spaces for consistent caret alignment.
func (p *SourceCodePosition) FormatError(message string) string {
	if p == nil {
		return message
	}

	if p.SourceLine == "" {
		return fmt.Sprintf(
			"%s at line %d, column %d",
			message,
			p.Line,
			p.Column,
		)
	}

	expandedLine := ExpandTabs(p.SourceLine, SourceTabWidth)
	visualColumn := CalculateVisualColumn(p.SourceLine, p.Column, SourceTabWidth)
	trimmedLine, caretColumn := buildErrorLineSnippet(expandedLine, visualColumn)

	result := fmt.Sprintf(
		"Error: %s\n  at line %d, column %d\n",
		message,
		p.Line,
		p.Column,
	)

	result += "  " + trimmedLine + "\n"
	pointerIndent := "  " + strings.Repeat(" ", caretColumn)
	result += pointerIndent + "^ " + message

	return result
}
