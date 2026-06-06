package cclValues

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
)

// RegisterSourceCodeDefinition registers one source file definition in the context.
func (c *CCLCodeContext) RegisterSourceCodeDefinition(def *SourceCodeDefinition) SourceFileId {
	if def == nil {
		return 0
	}

	c.sourceFilesLock.Lock()
	defer c.sourceFilesLock.Unlock()

	if def.SourceFileId == 0 {
		if def.FilePath != "" {
			if existingId, exists := c.sourceFileIdsByPath[def.FilePath]; exists {
				def.SourceFileId = existingId
			}
		}

		if def.SourceFileId == 0 {
			def.SourceFileId = c.nextSourceFileId
			c.nextSourceFileId++
		}
	}

	c.sourceFiles[def.SourceFileId] = def
	if def.FilePath != "" {
		c.sourceFileIdsByPath[def.FilePath] = def.SourceFileId
	}

	return def.SourceFileId
}

// GetSourceCodeDefinition returns a registered source file definition.
func (c *CCLCodeContext) GetSourceCodeDefinition(id SourceFileId) *SourceCodeDefinition {
	c.sourceFilesLock.RLock()
	defer c.sourceFilesLock.RUnlock()

	return c.sourceFiles[id]
}

// GetSourceFileIdByPath returns the source file id registered for a path.
func (c *CCLCodeContext) GetSourceFileIdByPath(path string) SourceFileId {
	c.sourceFilesLock.RLock()
	defer c.sourceFilesLock.RUnlock()

	return c.sourceFileIdsByPath[path]
}

// GetSourceDefinitions returns all source definitions registered in this context.
func (c *CCLCodeContext) GetSourceDefinitions() []*SourceCodeDefinition {
	c.sourceFilesLock.RLock()
	defer c.sourceFilesLock.RUnlock()

	definitions := make([]*SourceCodeDefinition, 0, len(c.sourceFiles))
	for _, def := range c.sourceFiles {
		definitions = append(definitions, def)
	}
	return definitions
}

// RegisterScopedAttributes registers all scoped attributes from one source file.
func (c *CCLCodeContext) RegisterScopedAttributes(def *SourceCodeDefinition) {
	if def == nil {
		return
	}

	c.scopedAttributesLock.Lock()
	defer c.scopedAttributesLock.Unlock()

	c.contextGlobalAttributes = append(c.contextGlobalAttributes, def.GlobalAttributes...)
	for _, attr := range def.NamespaceAttributes {
		if attr == nil {
			continue
		}

		namespace := attr.Namespace
		if namespace == "" {
			namespace = def.Namespace
		}
		if namespace == "" {
			namespace = gValues.DefaultMainNamespace
		}

		c.contextNamespaceAttributes[namespace] = append(
			c.contextNamespaceAttributes[namespace],
			attr,
		)
	}
}

// RegisterGlobalAttribute registers one compilation-wide global attribute.
func (c *CCLCodeContext) RegisterGlobalAttribute(attr *AttributeUsageInfo) {
	if attr == nil {
		return
	}

	c.scopedAttributesLock.Lock()
	defer c.scopedAttributesLock.Unlock()

	c.contextGlobalAttributes = append(c.contextGlobalAttributes, attr)
}

// FindContextGlobalAttributes returns context-wide global attributes.
func (c *CCLCodeContext) FindContextGlobalAttributes(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
) []*AttributeUsageInfo {
	c.scopedAttributesLock.RLock()
	defer c.scopedAttributesLock.RUnlock()

	return filterAttributes(c.contextGlobalAttributes, targetLang, name)
}

// FindContextGlobalAttribute returns the first context-wide global attribute.
func (c *CCLCodeContext) FindContextGlobalAttribute(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
) *AttributeUsageInfo {
	attrs := c.FindContextGlobalAttributes(targetLang, name)
	if len(attrs) == 0 {
		return nil
	}

	return attrs[0]
}

// FindNamespaceAttributes returns namespace-scoped attributes.
func (c *CCLCodeContext) FindNamespaceAttributes(
	namespace string,
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
) []*AttributeUsageInfo {
	c.scopedAttributesLock.RLock()
	defer c.scopedAttributesLock.RUnlock()

	return filterAttributes(c.contextNamespaceAttributes[namespace], targetLang, name)
}

// FindNamespaceAttribute returns the first namespace-scoped attribute.
func (c *CCLCodeContext) FindNamespaceAttribute(
	namespace string,
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
) *AttributeUsageInfo {
	attrs := c.FindNamespaceAttributes(namespace, targetLang, name)
	if len(attrs) == 0 {
		return nil
	}

	return attrs[0]
}
