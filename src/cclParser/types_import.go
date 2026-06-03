package cclParser

type importGraphResolver struct {
	visitedFiles map[string]bool
	activeFiles  map[string]bool
	fileStack    []string
}
