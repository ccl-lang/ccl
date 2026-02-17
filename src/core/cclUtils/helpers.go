package cclUtils

import (
	"strings"

	"github.com/ALiwoto/ssg/ssg"
)

// SnakeToTitle converts a snake_case string to TitleCase.
func SnakeToTitle(s string) string {
	bd := strings.Builder{}

	for _, split := range strings.Split(s, "_") {
		bd.WriteString(ssg.Title(split))
	}

	return bd.String()
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
func ToSnakeCase(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		// A is 65, a is 97
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		// v is capital letter here
		// disregard first letter
		// add underscore if last letter is capital letter
		// add underscore when previous letter is lowercase
		// add underscore when next letter is lowercase
		if (i != 0 || i == l-1) && (          // head and tail
		(i > 0 && rune(camel[i-1]) >= 'a') || // pre
			(i < l-1 && rune(camel[i+1]) >= 'a')) { //next
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
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
