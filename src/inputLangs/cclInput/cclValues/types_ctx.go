package cclValues

import "sync"

type CCLCodeContext struct {
	// typeDefinitionsCache is a cache for type definitions.
	// This is used to avoid creating multiple instances of the same type definition.
	// Key: full-name of the type definition (including namespace).
	// Value: type definition.
	typeDefinitionsCache map[string]*CCLTypeDefinition
	typeDefinitionsLock  *sync.RWMutex

	// typeDefinitionsOrder stores completed custom type definitions in
	// deterministic registration order for code generation.
	typeDefinitionsOrder []*CCLTypeDefinition

	// incompleteTypeDefinitionsCache is a cache for incomplete type definitions.
	// This is used to avoid creating multiple instances of the same incomplete
	// type definition each time it is referenced.
	// Key: full-name of the type definition (including namespace).
	// Value: type definition.
	// Later on, when the type definition is completed, it will be moved
	// to the typeDefinitionsCache and marked as complete.
	incompleteTypeDefinitionsCache map[string]*CCLTypeDefinition

	// The map of all global variable definitions.
	// Key: variable name.
	// Value: variable definition.
	// Please do NOT use this map directly, use GetGlobalVariable instead.
	cclGlobalVariables map[string]*VariableDefinition

	// cclAutomaticVariables is the map of all automatic variable definitions.
	// Key: variable name.
	// Value: variable definition.
	// Please do NOT use this map directly, use GetGlobalVariable instead.
	cclAutomaticVariables map[string]*VariableDefinition

	// globalVarsLock is the lock for global (and automatic) variables.
	globalVarsLock *sync.RWMutex

	// sourceFiles contains all sanitized source files registered in this context.
	// Key: context-assigned source file id.
	// Value: source code definition for that file.
	sourceFiles map[SourceFileId]*SourceCodeDefinition

	// sourceFileIdsByPath maps canonical absolute file paths to source file ids.
	sourceFileIdsByPath map[string]SourceFileId

	// nextSourceFileId is the next context source file id to assign.
	nextSourceFileId SourceFileId

	// sourceFilesLock is the lock for source file indexes.
	sourceFilesLock *sync.RWMutex

	// contextGlobalAttributes contains compilation-wide attributes.
	contextGlobalAttributes []*AttributeUsageInfo

	// contextNamespaceAttributes contains namespace-scoped attributes.
	// Key: namespace.
	contextNamespaceAttributes map[string][]*AttributeUsageInfo

	// scopedAttributesLock is the lock for global and namespace attributes.
	scopedAttributesLock *sync.RWMutex

	// modelIdCounter is a counter that is used to generate unique model IDs.
	// needs a lock to be put on typeDefinitionsLock.
	modelIdCounter int64
}
