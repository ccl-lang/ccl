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

	filePath := p.FilePath
	if filePath == "" {
		// we will make it become like this for two reasons:
		// 1. we don't get empty weird lines such as ":10:20"
		// 2. we check the logs, find out in some places we are not passing the correct
		// 	source position, so we can easily fix it
		filePath = "unknown_file"
	}
	location := fmt.Sprintf("%s:%d:%d", filePath, p.Line, p.Column)

	if p.SourceLine == "" {
		return fmt.Sprintf(
			"%s: Error: %s",
			location,
			message,
		)
	}

	expandedLine := ExpandTabs(p.SourceLine, SourceTabWidth)
	visualColumn := CalculateVisualColumn(p.SourceLine, p.Column, SourceTabWidth)
	trimmedLine, caretColumn := buildErrorLineSnippet(expandedLine, visualColumn)

	result := fmt.Sprintf(
		"Error: %s\n at %s \n",
		message,
		location,
	)

	result += "  " + trimmedLine + "\n"
	pointerIndent := "  " + strings.Repeat(" ", caretColumn)
	result += pointerIndent + "^ " + message

	return result
}
