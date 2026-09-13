package cclSanitizer

import (
	"fmt"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

// ResolveEnumMappings resolves attributes after all declarations in the source graph exist.
func ResolveEnumMappings(ctx *cclValues.CCLCodeContext) error {
	enums := []*cclValues.EnumDefinition{}
	for _, definition := range ctx.GetGenerationTypeDefinitions() {
		if !definition.IsCustomEnum() {
			continue
		}
		enumDef := definition.GetEnumDefinition()
		enums = append(enums, enumDef)
		enumDef.MappingFallbacks = make(map[gValues.LanguageType]*cclValues.EnumMemberDefinition)
		enumDef.Mappings = make(map[gValues.LanguageType][]*cclValues.EnumMapping)
		for _, attr := range enumDef.Attributes {
			if !isEnumMappingAttribute(attr.Name) {
				continue
			}
			if err := resolveEnumMappingSymbol(ctx, enumDef, attr); err != nil {
				return err
			}
			if attr.Name != cclAttr.AttrEnumMapUnknownMember {
				continue
			}
			member := attr.GetParamAt(0).GetValue().(*cclValues.EnumMemberReference).Member
			for _, language := range enumMappingLanguages(attr) {
				previous := enumDef.MappingFallbacks[language]
				if previous != nil && previous != member {
					return enumMappingError(attr, "conflicting fallback members for "+enumDef.GetFullName()+" in "+language.String())
				}
				enumDef.MappingFallbacks[language] = member
			}
		}
	}
	for _, source := range enums {
		for _, attr := range source.Attributes {
			if attr.Name != cclAttr.AttrEnumMapSameMembers {
				continue
			}
			target := attr.GetParamAt(0).GetValue().(*cclValues.EnumDefinition)
			for _, language := range enumMappingLanguages(attr) {
				fallback := target.MappingFallbacks[language]
				if fallback == nil {
					return enumMappingError(attr, "destination enum '"+target.GetFullName()+"' requires EnumMapUnknownMember for "+language.String())
				}
				mapping, err := buildEnumMapping(source, target, fallback, attr)
				if err != nil {
					return err
				}
				duplicate := false
				for _, previous := range source.Mappings[language] {
					if previous.Target == target {
						duplicate = true
						break
					}
				}
				if !duplicate {
					source.Mappings[language] = append(source.Mappings[language], mapping)
				}
			}
		}
	}
	return nil
}

func enumMappingLanguages(attr *cclValues.AttributeUsageInfo) []gValues.LanguageType {
	if len(attr.Languages) != 0 {
		return attr.Languages
	}
	return []gValues.LanguageType{
		gValues.LanguageGo, gValues.LanguageGd, gValues.LanguageCS, gValues.LanguagePy,
		gValues.LanguageJS, gValues.LanguageTS, gValues.LanguageRust,
	}
}

func buildEnumMapping(
	source, target *cclValues.EnumDefinition,
	fallback *cclValues.EnumMemberDefinition,
	attr *cclValues.AttributeUsageInfo,
) (*cclValues.EnumMapping, error) {
	mapping := &cclValues.EnumMapping{Target: target, Fallback: fallback}
	values := map[int64]*cclValues.EnumMappingCase{}
	for _, member := range source.Members {
		destination := target.GetMemberByName(member.Name)
		if destination == nil {
			destination = fallback
		}
		if previous := values[member.Value]; previous != nil {
			if previous.Target.Value != destination.Value {
				return nil, enumMappingError(attr, fmt.Sprintf(
					"source aliases %s.%s and %s share value %d but require different results in %s (%s and %s)",
					source.GetFullName(), previous.Source.Name, member.Name, member.Value,
					target.GetFullName(), previous.Target.Name, destination.Name))
			}
			continue
		}
		mappingCase := &cclValues.EnumMappingCase{Source: member, Target: destination}
		values[member.Value] = mappingCase
		// Unmatched members are already handled by the default arm.
		if destination.Value != fallback.Value {
			mapping.Cases = append(mapping.Cases, mappingCase)
		}
	}
	return mapping, nil
}
