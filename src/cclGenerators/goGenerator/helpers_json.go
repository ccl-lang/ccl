package goGenerator

import "strings"

func lowerGoVarName(value string) string {
	if value == "" {
		return value
	}
	return strings.ToLower(value[:1]) + value[1:]
}
