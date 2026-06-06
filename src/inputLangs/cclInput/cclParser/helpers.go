package cclParser

import (
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser/cclLexer"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclSanitizer"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

// ParseCCLSourceFile reads a CCL source file and parses it into a
// SourceCodeDefinition value. It resolves source imports, parses the combined
// AST, and sanitizes the AST into IR.
// The function returns a pointer to a SourceCodeDefinition value and an error if any occurred.
func ParseCCLSourceFile(options *CCLParseOptions) (*cclValues.SourceCodeDefinition, error) {
	if options == nil {
		return nil, &ImportResolutionError{
			Message: "missing parse options",
		}
	}

	if options.CodeContext == nil {
		options.CodeContext = cclValues.NewCCLCodeContext()
	}

	astFiles, rootAst, _, err := newImportGraphResolver().parseSourceFileGraphAsAST(options)
	if err != nil {
		return nil, err
	}

	return sanitizeFileGraph(options.CodeContext, astFiles, rootAst)
}

// ParseCCLSourceContent takes a CCLParseOptions struct and parses standalone
// SourceContent into a SourceCodeDefinition value by going through AST -> IR.
// The function returns a pointer to a SourceCodeDefinition value and an error if any occurred.
func ParseCCLSourceContent(options *CCLParseOptions) (*cclValues.SourceCodeDefinition, error) {
	allTokens, err := cclLexer.Lex(options.SourceContent)
	if err != nil {
		return nil, err
	}

	if options.CodeContext == nil {
		options.CodeContext = cclValues.NewCCLCodeContext()
	}
	return ParseCCL(options.CodeContext, allTokens, options)
}

func ParseCCL(
	ctx *cclValues.CCLCodeContext,
	tokens []*cclLexer.CCLToken,
	options *CCLParseOptions,
) (*cclValues.SourceCodeDefinition, error) {
	astFile, err := ParseCCLAst(tokens, options)
	if err != nil {
		return nil, err
	}

	if len(astFile.Imports) > 0 {
		firstImport := astFile.Imports[0]
		return nil, &ImportResolutionError{
			ImportPath:     firstImport.Path,
			Message:        "imports require parsing from a source file path",
			SourcePosition: firstImport.SourcePosition,
		}
	}

	return cclSanitizer.SanitizeCCLAst(ctx, astFile)
}

func sanitizeFileGraph(
	ctx *cclValues.CCLCodeContext,
	astFiles []*cclAst.CCLFileAST,
	rootAst *cclAst.CCLFileAST,
) (*cclValues.SourceCodeDefinition, error) {
	aggregate := &cclValues.SourceCodeDefinition{
		CodeContext: ctx,
	}

	for _, astFile := range astFiles {
		definition, err := cclSanitizer.SanitizeCCLAst(ctx, astFile)
		if err != nil {
			return nil, err
		}

		if rootAst != nil && astFile.FilePath == rootAst.FilePath {
			aggregate.SourceFileId = definition.SourceFileId
			aggregate.FilePath = definition.FilePath
			aggregate.Namespace = definition.Namespace
		}

		aggregate.TypeDefinitions = append(aggregate.TypeDefinitions, definition.TypeDefinitions...)
		aggregate.GlobalAttributes = append(aggregate.GlobalAttributes, definition.GlobalAttributes...)
		aggregate.FileAttributes = append(aggregate.FileAttributes, definition.FileAttributes...)
		aggregate.NamespaceAttributes = append(aggregate.NamespaceAttributes, definition.NamespaceAttributes...)
	}

	return aggregate, nil
}
