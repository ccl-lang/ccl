package cclValues

import (
	"sync"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

func (c *CCLCodeContext) initialize() {
	c.typeDefinitionsCache = map[string]*CCLTypeDefinition{}
	c.typeDefinitionsOrder = []*CCLTypeDefinition{}
	c.typeDefinitionsLock = &sync.RWMutex{}
	c.incompleteTypeDefinitionsCache = map[string]*CCLTypeDefinition{}
	c.globalVarsLock = &sync.RWMutex{}
	c.sourceFiles = map[SourceFileId]*SourceCodeDefinition{}
	c.sourceFileIdsByPath = map[string]SourceFileId{}
	c.nextSourceFileId = 1
	c.sourceFilesLock = &sync.RWMutex{}
	c.contextGlobalAttributes = []*AttributeUsageInfo{}
	c.contextNamespaceAttributes = map[string][]*AttributeUsageInfo{}
	c.scopedAttributesLock = &sync.RWMutex{}

	// initialize automatic variables

	c.initializeAutoVars()
	c.initializeGlobalVars()
}

func (c *CCLCodeContext) initializeAutoVars() {
	c.cclAutomaticVariables = map[string]*VariableDefinition{
		"__ccl_version": c.newStringAutomaticVariable("__ccl_version", gValues.CurrentCCLVersion),
	}
}

func (c *CCLCodeContext) initializeGlobalVars() {
	c.cclGlobalVariables = map[string]*VariableDefinition{}
}
