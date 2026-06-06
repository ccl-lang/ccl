package cclValues

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
)

// ResolveAttribute resolves one attribute using local, file, namespace, and
// global fallback rules.
func (c *CCLCodeContext) ResolveAttribute(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
	subject *AttributeResolutionSubject,
	options *AttributeResolutionOptions,
) *AttributeUsageInfo {
	attrs := c.ResolveAttributes(targetLang, name, subject, options)
	if len(attrs) == 0 {
		return nil
	}

	return attrs[0]
}

// ResolveAttributes resolves attributes using local, file, namespace, and
// global fallback rules. It stops at the first scoped level with matches.
func (c *CCLCodeContext) ResolveAttributes(
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
	subject *AttributeResolutionSubject,
	options *AttributeResolutionOptions,
) []*AttributeUsageInfo {
	if subject == nil {
		// no specific subject to work on, just assume we want global
		return c.FindContextGlobalAttributes(targetLang, name)
	}

	if options == nil {
		options = &AttributeResolutionOptions{
			AllowModelFallback:  true,
			AllowScopedFallback: true,
			AllowGlobalFallback: true,
		}
	}

	if subject.Field != nil {
		attrs := filterAttributes(subject.Field.Attributes, targetLang, name)
		if len(attrs) > 0 {
			return attrs
		}

		if options.AllowModelFallback && subject.Field.OwnedBy != nil && subject.Model == nil {
			subject.Model = subject.Field.OwnedBy
		}
	}

	if subject.Model != nil {
		attrs := subject.Model.FindAttributes(targetLang, name)
		if len(attrs) > 0 {
			return attrs
		}
	}

	if !options.AllowScopedFallback {
		return nil
	}

	sourceFileId := subject.SourceFileId
	namespace := subject.Namespace
	if subject.Model != nil {
		if sourceFileId == 0 {
			sourceFileId = subject.Model.SourceFileId
		}
		if namespace == "" {
			namespace = subject.Model.GetNamespace()
		}
	}
	if subject.Field != nil && subject.Field.OwnedBy != nil {
		if sourceFileId == 0 {
			sourceFileId = subject.Field.OwnedBy.SourceFileId
		}
		if namespace == "" {
			namespace = subject.Field.OwnedBy.GetNamespace()
		}
	}
	if namespace == "" {
		namespace = gValues.DefaultMainNamespace
	}

	if sourceFileId != 0 {
		if def := c.GetSourceCodeDefinition(sourceFileId); def != nil {
			attrs := def.FindFileAttributes(targetLang, name)
			if len(attrs) > 0 {
				return attrs
			}
		}
	}

	for currentNamespace := namespace; currentNamespace != ""; currentNamespace = parentNamespace(currentNamespace) {
		attrs := c.FindNamespaceAttributes(currentNamespace, targetLang, name)
		if len(attrs) > 0 {
			return attrs
		}
	}

	if !options.AllowGlobalFallback {
		return nil
	}

	return c.FindContextGlobalAttributes(targetLang, name)
}
