package globalValues

import "strings"

func GetNormalizedLanguageName(lang string) NormalizedLangName {
	return GetLanguageTypeFromName(lang).GetShortName()
}

func GetLanguageTypeFromName(name string) LanguageType {
	return langsAliasNames[strings.ToLower(name)]
}
