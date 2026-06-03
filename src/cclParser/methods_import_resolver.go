package cclParser

import (
	"os"
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclAst"
)

func (r *importGraphResolver) parseSourceFileAsAST(options *CCLParseOptions) (*cclAst.CCLFileAST, error) {
	if options == nil {
		return nil, &ImportResolutionError{
			Message: "missing parse options",
		}
	}

	sourceFilePath, err := getAbsoluteSourceFilePath(options.SourceFilePath)
	if err != nil {
		return nil, &ImportResolutionError{
			Message:    "invalid source file path",
			InnerError: err,
		}
	}

	astFile, sourceContent, err := r.parseAbsoluteSourceFileAsAST(sourceFilePath, nil)
	if err != nil {
		return nil, err
	}

	options.SourceFilePath = sourceFilePath
	options.SourceContent = sourceContent
	return astFile, nil
}

func (r *importGraphResolver) parseAbsoluteSourceFileAsAST(
	sourceFilePath string,
	importDecl *cclAst.ImportDecl,
) (*cclAst.CCLFileAST, string, error) {
	if r.visitedFiles[sourceFilePath] {
		return &cclAst.CCLFileAST{
			FilePath: sourceFilePath,
		}, "", nil
	}

	if r.activeFiles[sourceFilePath] {
		return nil, "", &ImportResolutionError{
			ImportPath:     getImportDeclPath(importDecl),
			ResolvedPath:   sourceFilePath,
			Message:        "import cycle detected: " + r.formatImportCycle(sourceFilePath),
			SourcePosition: getImportDeclSourcePosition(importDecl),
		}
	}

	sourceContentBytes, err := os.ReadFile(sourceFilePath)
	if err != nil {
		return nil, "", &ImportResolutionError{
			ImportPath:     getImportDeclPath(importDecl),
			ResolvedPath:   sourceFilePath,
			Message:        "failed to read CCL source file",
			SourcePosition: getImportDeclSourcePosition(importDecl),
			InnerError:     err,
		}
	}

	sourceContent := string(sourceContentBytes)
	sourceAst, err := ParseCCLSourceContentAsAST(&CCLParseOptions{
		SourceFilePath: sourceFilePath,
		SourceContent:  sourceContent,
	})
	if err != nil {
		if importDecl == nil {
			return nil, "", err
		}

		return nil, "", &ImportResolutionError{
			ImportPath:     importDecl.Path,
			ResolvedPath:   sourceFilePath,
			Message:        "failed to parse imported CCL source file",
			SourcePosition: importDecl.SourcePosition,
			InnerError:     err,
		}
	}

	r.activeFiles[sourceFilePath] = true
	r.fileStack = append(r.fileStack, sourceFilePath)
	defer func() {
		delete(r.activeFiles, sourceFilePath)
		r.fileStack = r.fileStack[:len(r.fileStack)-1]
	}()

	mergedAst := &cclAst.CCLFileAST{
		FilePath:       sourceAst.FilePath,
		Namespace:      sourceAst.Namespace,
		Imports:        append([]*cclAst.ImportDecl{}, sourceAst.Imports...),
		SourcePosition: sourceAst.SourcePosition,
	}

	for _, currentImport := range sourceAst.Imports {
		importedFilePath, err := resolveImportPath(sourceFilePath, currentImport)
		if err != nil {
			return nil, "", err
		}

		importedAst, _, err := r.parseAbsoluteSourceFileAsAST(importedFilePath, currentImport)
		if err != nil {
			return nil, "", err
		}

		mergedAst.GlobalAttributes = append(mergedAst.GlobalAttributes, importedAst.GlobalAttributes...)
		mergedAst.Models = append(mergedAst.Models, importedAst.Models...)
	}

	mergedAst.GlobalAttributes = append(mergedAst.GlobalAttributes, sourceAst.GlobalAttributes...)
	mergedAst.Models = append(mergedAst.Models, sourceAst.Models...)
	r.visitedFiles[sourceFilePath] = true

	return mergedAst, sourceContent, nil
}

func (r *importGraphResolver) formatImportCycle(repeatedPath string) string {
	cyclePaths := append([]string{}, r.fileStack...)
	cyclePaths = append(cyclePaths, repeatedPath)
	for index, currentPath := range cyclePaths {
		if currentPath == repeatedPath {
			cyclePaths = cyclePaths[index:]
			break
		}
	}

	return strings.Join(cyclePaths, " -> ")
}
