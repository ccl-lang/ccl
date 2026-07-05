package rsGenerator

import "github.com/ALiwoto/ssg/ssg/caseUtils"

func rustModuleFileName(name string) string {
	return caseUtils.ToSnakeCase(name) + ".rs"
}

func rustModuleName(name string) string {
	return caseUtils.ToSnakeCase(name)
}

func rustFieldName(name string) string {
	return rustSafeIdentifier(caseUtils.ToSnakeCase(name))
}

func rustSafeIdentifier(name string) string {
	switch name {
	case "as", "break", "const", "continue", "crate", "else", "enum", "extern",
		"false", "fn", "for", "if", "impl", "in", "let", "loop", "match",
		"mod", "move", "mut", "pub", "ref", "return", "self", "Self",
		"static", "struct", "super", "trait", "true", "type", "unsafe",
		"use", "where", "while", "async", "await", "dyn":
		return "r#" + name
	default:
		return name
	}
}
