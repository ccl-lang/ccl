package cclValues

import (
	"strings"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
)

func filterAttributes(
	attrs []*AttributeUsageInfo,
	targetLang gValues.LanguageType,
	name cclAttr.CCLAttributeName,
) []*AttributeUsageInfo {
	result := []*AttributeUsageInfo{}
	for _, attr := range attrs {
		if attr == nil {
			continue
		}

		if attr.Name == name && attr.IsForLanguage(targetLang) {
			result = append(result, attr)
		}
	}

	return result
}

func parentNamespace(namespace string) string {
	index := strings.LastIndex(namespace, ".")
	if index < 0 {
		return ""
	}

	return namespace[:index]
}
