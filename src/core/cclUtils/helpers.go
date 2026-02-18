package cclUtils

import (
	"strings"

	"github.com/ALiwoto/ssg/ssg"
)

// SnakeToTitle converts a snake_case string to TitleCase.
func SnakeToTitle(str string) string {
	str = ToSnakeCase(str)
	allStrs := strings.Split(str, "_")
	builder := strings.Builder{}

	for _, current := range allStrs {
		builder.WriteString(ssg.Title(current))
	}

	return builder.String()
}

// ToCamelCase converts a string to camelCase.
func ToCamelCase(s string) string {
	title := SnakeToTitle(s)

	return strings.ToLower(title[:1]) + title[1:]
}

// ToPascalCase converts a string to PascalCase.
func ToPascalCase(str string) string {
	title := SnakeToTitle(str)

	return strings.ToUpper(title[:1]) + title[1:]
}

// ToSnakeCase converts a CamelCase string to snake_case.
func ToSnakeCase(str string) string {
	var b strings.Builder
	runes := []rune(str)
	str_len := len(runes)

	// Track if the last character written was an underscore
	// Initialize true to prevent leading underscores
	lastUnderscore := true

	for current_index := range str_len {
		current_rune := runes[current_index]
		nextIsLower := false
		if current_index+1 < str_len && runes[current_index+1] >= 'a' && runes[current_index+1] <= 'z' {
			nextIsLower = true
		}

		// 1. Handle delimiters (-, ., space, :)
		// Convert them all to underscores for processing
		if current_rune == '-' || current_rune == '.' || current_rune == ' ' ||
			current_rune == ':' || current_rune == '_' {
			if !lastUnderscore {
				b.WriteRune('_')
				lastUnderscore = true
			}
			continue
		}

		// 2. Handle Uppercase
		if current_rune >= 'A' && current_rune <= 'Z' {
			// Check if we need to insert an underscore before this capital
			if !lastUnderscore {
				prev := runes[current_index-1]
				isPrevLower := (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9')

				// Insert _ if:
				// a. Previous was lowercase/digit (camelCase -> camel_case)
				// b. Previous was Upper but next is Lower (JSONId -> json_id)
				if isPrevLower || (current_rune >= 'A' && current_rune <= 'Z' && nextIsLower) {
					b.WriteRune('_')
					lastUnderscore = true
				}
			}

			// Convert to lowercase
			current_rune = current_rune + 32
		}

		// 3. Write the character
		b.WriteRune(current_rune)
		lastUnderscore = false
	}

	return b.String()
}

// GetSourceLineByNumber returns the line of source code at the given line number (1-based).
func GetSourceLineByNumber(content string, lineNum int) string {
	if lineNum <= 0 || content == "" {
		return ""
	}

	runes := []rune(content)
	currentLine := 1
	lineStart := 0

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r != '\n' && r != '\r' {
			continue
		}

		if currentLine == lineNum {
			return string(runes[lineStart:i])
		}

		if r == '\r' && i+1 < len(runes) && runes[i+1] == '\n' {
			i++
		}

		currentLine++
		lineStart = i + 1
	}

	if currentLine == lineNum {
		return string(runes[lineStart:])
	}

	return ""
}
