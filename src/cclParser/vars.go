package cclParser

import "regexp"

var (
	modelRegex = regexp.MustCompile(`(?ms)model\s+(\w+)\s*\{(.*?)\}`)
	fieldRegex = regexp.MustCompile(`(\w+)\s*:\s*([\w]+)\s*(\[\s*\])?\s*;`)
)
