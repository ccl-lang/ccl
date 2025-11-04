package globalValues

// String returns the string representation of this language type.
func (l LanguageType) String() string {
	return l.GetShortName()
}

// GetShortName returns the short name of the current language
// type.
func (l LanguageType) GetShortName() NormalizedLangName {
	return langsToShortName[l]
}
