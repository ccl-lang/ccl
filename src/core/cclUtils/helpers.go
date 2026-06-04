package cclUtils

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
