package cclParser

import (
	"path/filepath"

	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/cclUtils"
)

func newImportGraphResolver() *importGraphResolver {
	return &importGraphResolver{
		visitedFiles: map[string]bool{},
		activeFiles:  map[string]bool{},
	}
}

func getAbsoluteSourceFilePath(sourceFilePath string) (string, error) {
	if sourceFilePath == "" {
		return "", &ImportResolutionError{
			Message: "source file path is required when resolving imports",
		}
	}

	absolutePath, err := filepath.Abs(sourceFilePath)
	if err != nil {
		return "", err
	}

	return filepath.Clean(absolutePath), nil
}

func resolveImportPath(
	importingFilePath string,
	importDecl *cclAst.ImportDecl,
) (string, error) {
	if importDecl == nil {
		return "", &ImportResolutionError{
			Message: "missing import declaration",
		}
	}

	if importDecl.Path == "" {
		return "", &ImportResolutionError{
			Message:        "import path cannot be empty",
			SourcePosition: importDecl.SourcePosition,
		}
	}

	targetPath := importDecl.Path
	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(filepath.Dir(importingFilePath), targetPath)
	}

	absolutePath, err := filepath.Abs(targetPath)
	if err != nil {
		return "", &ImportResolutionError{
			ImportPath:     importDecl.Path,
			Message:        "failed to resolve import path",
			SourcePosition: importDecl.SourcePosition,
			InnerError:     err,
		}
	}

	return filepath.Clean(absolutePath), nil
}

func getImportDeclPath(importDecl *cclAst.ImportDecl) string {
	if importDecl == nil {
		return ""
	}

	return importDecl.Path
}

func getImportDeclSourcePosition(importDecl *cclAst.ImportDecl) *cclUtils.SourceCodePosition {
	if importDecl == nil {
		return nil
	}

	return importDecl.SourcePosition
}
