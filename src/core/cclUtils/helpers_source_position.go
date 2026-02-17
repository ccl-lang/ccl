package cclUtils

import "unicode"

func buildErrorLineSnippet(line string, caretColumn int) (string, int) {
	if line == "" {
		return line, 0
	}

	lineRunes := []rune(line)
	lineLen := len(lineRunes)

	if caretColumn < 0 {
		caretColumn = 0
	} else if caretColumn > lineLen {
		caretColumn = lineLen
	}

	if lineLen <= MaxErrorSourceLineBeforeLen+MaxErrorSourceLineAfterLen+1 {
		return line, caretColumn
	}

	focusIndex := caretColumn
	if focusIndex >= lineLen {
		focusIndex = lineLen - 1
	}

	if focusIndex > 0 &&
		!isIdentifierRune(lineRunes[focusIndex]) &&
		isIdentifierRune(lineRunes[focusIndex-1]) {
		focusIndex--
	}

	focusStart, focusEnd := findIdentifierSpan(lineRunes, focusIndex)
	if focusStart == focusEnd {
		focusStart = focusIndex
		focusEnd = focusIndex + 1
		if focusEnd > lineLen {
			focusEnd = lineLen
		}
	}

	beforeLen := MaxErrorSourceLineBeforeLen
	afterLen := MaxErrorSourceLineAfterLen

	if beforeLen < 0 {
		beforeLen = 0
	}
	if afterLen < 0 {
		afterLen = 0
	}

	start := focusStart - beforeLen
	if start < 0 {
		start = 0
	}
	end := focusEnd + afterLen
	if end > lineLen {
		end = lineLen
	}

	ellipsisRunes := []rune(SourceErrorEllipsis)
	ellipsisLen := len(ellipsisRunes)
	needsPrefix := start > 0
	needsSuffix := end < lineLen

	snippetRunes := lineRunes[start:end]
	caretPos := caretColumn - start

	if needsPrefix {
		snippetRunes = append(ellipsisRunes, snippetRunes...)
		caretPos += ellipsisLen
	}
	if needsSuffix {
		snippetRunes = append(snippetRunes, ellipsisRunes...)
	}

	if caretPos < 0 {
		caretPos = 0
	} else if caretPos > len(snippetRunes) {
		caretPos = len(snippetRunes)
	}

	return string(snippetRunes), caretPos
}

func findIdentifierSpan(line []rune, index int) (int, int) {
	if index < 0 || index >= len(line) {
		return 0, 0
	}

	if !isIdentifierRune(line[index]) {
		return 0, 0
	}

	start := index
	for start > 0 && isIdentifierRune(line[start-1]) {
		start--
	}

	end := index + 1
	for end < len(line) && isIdentifierRune(line[end]) {
		end++
	}

	return start, end
}

func isIdentifierRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
