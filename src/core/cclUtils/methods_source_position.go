package cclUtils

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

// ExpandTabs replaces tabs with spaces using a fixed tab width.
func ExpandTabs(line string, tabWidth int) string {
	if tabWidth <= 0 {
		tabWidth = SourceTabWidth
	}

	if !strings.Contains(line, "\t") {
		return line
	}

	builder := strings.Builder{}
	builder.Grow(len(line))

	visualColumn := 0
	for _, r := range line {
		if r == '\t' {
			spaces := tabWidth - (visualColumn % tabWidth)
			if spaces == 0 {
				spaces = tabWidth
			}
			builder.WriteString(strings.Repeat(" ", spaces))
			visualColumn += spaces
			continue
		}

		builder.WriteRune(r)
		visualColumn++
	}

	return builder.String()
}

// CalculateVisualColumn returns the visual column index (0-based) after tab expansion.
func CalculateVisualColumn(line string, column int, tabWidth int) int {
	if column <= 0 {
		return 0
	}

	if tabWidth <= 0 {
		tabWidth = SourceTabWidth
	}

	visualColumn := 0
	currentColumn := 0
	for _, r := range line {
		if currentColumn >= column {
			break
		}

		if r == '\t' {
			spaces := tabWidth - (visualColumn % tabWidth)
			if spaces == 0 {
				spaces = tabWidth
			}
			visualColumn += spaces
		} else {
			visualColumn++
		}

		currentColumn++
	}

	if visualColumn < 0 {
		return 0
	}

	return visualColumn
}
