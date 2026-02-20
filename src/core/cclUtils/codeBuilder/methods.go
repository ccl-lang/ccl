package codeBuilder

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/ALiwoto/ssg/ssg"
)

// addDebugInfo captures the current source location and adds it to the debug info for the current section.
func (c *CodeBuilder) addDebugInfo(skip int) {
	if !c.enableDebugInfo {
		return
	}
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return
	}

	fileParts := strings.Split(file, "/src/")
	file = "src/" + fileParts[len(fileParts)-1]

	c.debugInfos[c.currentSection] = append(c.debugInfos[c.currentSection], &DebugInfo{
		SourceFile:    file,
		SourceLine:    line,
		SectionOffset: c.builders[c.currentSection].Len(),
	})
}

// checkSection checks if the current section is set.
func (c *CodeBuilder) checkSection() {
	if c.currentSection == "" {
		// Developer misuse: this is an invariant violation, not a user error.
		// Panicking here is intentional to surface incorrect builder usage quickly.
		panic("CodeBuilder: illegal usage of CodeBuilder without initiating a section" + panicStatement)
	}
}

// IsImported checks if the given key has been marked as imported.
func (c *CodeBuilder) IsImported(key string) bool {
	c.checkSection()
	_, exists := c.importedKeys[key]
	return exists
}

// DoImport marks the given key as imported and appends the import line if not already imported.
// The import line is written to the "imports" section; this is a special section that does not
// need to be manually begun; HOWEVER, it does require that a section is active when this method is called.
func (c *CodeBuilder) DoImport(key, importLine string) *CodeBuilder {
	c.checkSection()
	if c.IsImported(key) {
		return c
	}

	c.importedKeys[key] = true
	c.write(SectionImports, importLine)
	c.write(SectionImports, c.newLineStr)
	return c
}

// AddCommentHeader adds a comment header to the "comment_headers" section with a new line.
// This is useful for comment headers such as license headers, "DO NOT EDIT" warnings, etc.
func (c *CodeBuilder) AddCommentHeader(comment string) *CodeBuilder {
	return c.AppendCommentHeader(comment + c.newLineStr)
}

// AppendCommentHeader adds a comment header to the "comment_headers" section.
// This is useful for comment headers such as license headers, "DO NOT EDIT" warnings, etc.
func (c *CodeBuilder) AppendCommentHeader(comment string) *CodeBuilder {
	c.checkSection()
	c.write(SectionCommentHeaders, comment)
	return c
}

// AddHeader adds a header to the "headers" section with a new line.
// This is useful for file headers such as package declarations, general file comments, etc.
// This will come after comment headers but before imports.
func (c *CodeBuilder) AddHeader(strValue string) *CodeBuilder {
	return c.AppendHeader(strValue + c.newLineStr)
}

// AppendHeader adds a header to the "headers" section.
// This is useful for file headers such as package declarations, general file comments, etc.
// This will come after comment headers but before imports.
func (c *CodeBuilder) AppendHeader(strValue string) *CodeBuilder {
	c.checkSection()
	c.write(SectionHeaders, strValue)
	return c
}

// AddNamespaceDeclaration adds a namespace declaration to the "namespace_declarations" section with a new line.
// This is useful for organizing code into different namespaces in languages such as C#.
// NOTE: this comes after imports and before type declarations.
func (c *CodeBuilder) AddNamespaceDeclaration(strValue string) *CodeBuilder {
	return c.AppendNamespaceDeclaration(strValue + c.newLineStr)
}

// AppendNamespaceDeclaration adds a namespace declaration to the "namespace_declarations" section.
// This is useful for organizing code into different namespaces.
func (c *CodeBuilder) AppendNamespaceDeclaration(strValue string) *CodeBuilder {
	c.checkSection()
	c.write(SectionDeclareNamespace, strValue)
	return c
}

// BeginSection begins a new section or switches to an existing one.
// IMPORTANT NOTE: As long as you are in this section, there is a lock held on this object,
// to avoid deadlocks, you SHOULD call EndSection() when you are done with this section.
func (c *CodeBuilder) BeginSection(section string) *CodeBuilder {
	if section == "" {
		panic("CodeBuilder: section name cannot be empty" + panicStatement)
	}

	if c.currentSection != "" {
		if c.currentSection == section {
			return c
		}
		panic("CodeBuilder: cannot begin a new section without ending the previous one" + panicStatement)
	}

	c.mut.Lock()

	c.currentSection = section
	if _, exists := c.builders[section]; !exists {
		c.builders[section] = &strings.Builder{}
		c.indentations[section] = 0
	}
	return c
}

// EndSection ends the current section and releases the lock.
func (c *CodeBuilder) EndSection() *CodeBuilder {
	if c.currentSection == "" {
		panic("CodeBuilder: cannot end a section when no section is active" + panicStatement)
	}

	c.currentSection = ""
	c.mut.Unlock()
	return c
}

// MapVar will basically map a variable name to its true real name
func (c *CodeBuilder) MapVar(varName, realVarName string) *CodeBuilder {
	c.checkSection()
	c.mapVar(varName, realVarName)
	return c
}

// mapVar internal method for map var.
func (c *CodeBuilder) mapVar(varName, realVarName string) *CodeBuilder {
	c.mappedVars.setVar(c.currentSection, varName, realVarName)
	return c
}

// MapVarPairs tries to call MapVar method using pairs as (name, realName) from the passed argument.
// It's worthy to note that the passed args' length MUST BE even, otherwise this method will panic.
// NOTE: these are NOT global vars. these vars will be added inside of the current section.
func (c *CodeBuilder) MapVarPairs(values ...string) *CodeBuilder {
	if len(values)%2 != 0 {
		panic("CodeBuilder: MapVarPairs method was called with invalid length" + panicStatement)
	}

	c.checkSection()
	for i := 0; i < len(values); i += 2 {
		c.mapVar(values[i], values[i+1])
	}

	return c
}

// ExpectMappedVars will make sure all the passed varNames exist as a mapped-var inside
// of the code builder and if it doesn't exist, it will panic.
func (c *CodeBuilder) ExpectMappedVars(varNames ...string) *CodeBuilder {
	for _, current := range varNames {
		if !c.mappedVars.valueExists(c.currentSection, current) {
			panic("CodeBuilder: mapped var does not exist: " + current + panicStatement)
		}
	}

	return c
}

// UnmapVar will undefine the varName that was defined using the Define method.
func (c *CodeBuilder) UnmapVar(varNames ...string) *CodeBuilder {
	c.checkSection()
	c.mappedVars.deleteVar(c.currentSection, varNames...)
	return c
}

// MapGlobalVar will map a global varName to its realVarName that will be replaced *globally*.
func (c *CodeBuilder) MapGlobalVar(varName, realVarName string) *CodeBuilder {
	c.checkSection()
	c.mappedVars.setGlobal(varName, realVarName)
	return c
}

// UnmapGlobalVar will remove the varName from the global map.
func (c *CodeBuilder) UnmapGlobalVar(varName string) *CodeBuilder {
	c.checkSection()
	c.mappedVars.deleteGlobalVar(varName)
	return c
}

// Indent increases the indentation level.
func (c *CodeBuilder) Indent() *CodeBuilder {
	c.checkSection()
	c.indentations[c.currentSection]++
	return c
}

// Unindent decreases the indentation level.
func (c *CodeBuilder) Unindent() *CodeBuilder {
	c.checkSection()

	if c.indentations[c.currentSection] > 0 {
		c.indentations[c.currentSection]--
	}
	return c
}

// UnindentLine decreases the indentation level and adds a new line.
func (c *CodeBuilder) UnindentLine() *CodeBuilder {
	c.Unindent()

	c.addDebugInfo(2)
	c.checkSection()
	c.write(c.currentSection, c.newLineStr)
	return c
}

// WriteStr writes a string with the current indentation.
// If you don't want indentation, use AppendStr instead.
func (c *CodeBuilder) WriteStr(s string) *CodeBuilder {
	c.checkSection()

	c.addDebugInfo(2)
	c.writeIndentation()
	c.write(c.currentSection, s)
	return c
}

// AppendStr appends a string without adding indentation.
func (c *CodeBuilder) AppendStr(s string) *CodeBuilder {
	c.checkSection()

	c.addDebugInfo(2)
	c.checkSection()
	c.write(c.currentSection, s)
	return c
}

// WriteLine writes a line with the current indentation.
// If you don't want indentation, use AppendLine instead.
func (c *CodeBuilder) WriteLine(s string) *CodeBuilder {
	c.checkSection()

	c.addDebugInfo(2)
	c.writeIndentation()
	c.write(c.currentSection, s)
	c.write(c.currentSection, c.newLineStr)
	return c
}

// LineD Dynamically writes the string value.
// Calling this method is more expensive than normally calling WriteLine, because
// it will try to *resolve* all the variable names inside of the string.
func (c *CodeBuilder) LineD(s string) *CodeBuilder {
	c.checkSection()

	c.addDebugInfo(2)
	c.writeIndentation()
	c.writeDynamic(c.currentSection, s)
	c.write(c.currentSection, c.newLineStr)
	return c
}

// AppendD will dynamically append the string value without any new line char.
// Calling this method is more expensive than normally calling Append, because
// it will try to *resolve* all the variable names inside of the string.
func (c *CodeBuilder) AppendD(s string) *CodeBuilder {
	c.checkSection()

	c.addDebugInfo(2)
	c.writeIndentation()
	c.writeDynamic(c.currentSection, s)
	return c
}

// AppendLine appends a line without adding indentation.
func (c *CodeBuilder) AppendLine(s string) *CodeBuilder {
	c.addDebugInfo(2)
	c.checkSection()
	c.write(c.currentSection, s)
	c.write(c.currentSection, c.newLineStr)
	return c
}

// WriteLinef writes a formatted line with the current indentation.
func (c *CodeBuilder) WriteLinef(format string, args ...any) *CodeBuilder {
	c.checkSection()

	c.addDebugInfo(2)
	c.writeIndentation()
	c.writeF(c.currentSection, format, args...)
	c.write(c.currentSection, c.newLineStr)
	return c
}

// NewLine adds a new line.
func (c *CodeBuilder) NewLine() *CodeBuilder {
	c.addDebugInfo(2)
	c.checkSection()
	c.write(c.currentSection, c.newLineStr)
	return c
}

// writeIndentation writes the current indentation to the string builder.
// NOTE: since this is an internal method, it does NOT check for section or
// debug info; so the caller must ensure those are handled appropriately.
func (c *CodeBuilder) writeIndentation() *CodeBuilder {
	for i := 0; i < c.indentations[c.currentSection]; i++ {
		c.write(c.currentSection, c.indentationStr)
	}
	return c
}

// expandStr expands the string using the variables defined in the code builder.
func (c *CodeBuilder) expandStr(inputValue string) string {
	var builder strings.Builder
	// good heuristic; expanded value may exceed this, but still helps
	builder.Grow(len(inputValue))

	inVar := false
	// byte index of var name start (after indicator)
	varStart := 0
	indicatorLen := utf8.RuneLen(varIndicator)

	for index, current := range inputValue {
		if !inVar {
			if current == varIndicator {
				inVar = true
				varStart = index + indicatorLen
				continue
			}
			builder.WriteRune(current)
			continue
		}

		// We are parsing a variable name.
		// Logic change: Stop on space OR punctuation/symbols.
		// We only allow letters, digits, and underscores in var names here.
		if !isVarChar(current) {
			name := inputValue[varStart:index]
			val := c.mappedVars.getValue(c.currentSection, name)
			builder.WriteString(val)
			builder.WriteRune(current) // Write the char that ended the var
			inVar = false
		}
	}

	// Flush if string ended while still in var mode
	if inVar {
		name := inputValue[varStart:]
		val := c.mappedVars.getValue(c.currentSection, name)
		builder.WriteString(val)
	}

	return builder.String()
}

// writeDynamic
func (c *CodeBuilder) writeDynamic(targetSection, value string) {
	c.write(targetSection, c.expandStr(value))
}

// write writes to the target section's builder with the provided string value.
// This is an internal method and is not supposed to be used by any code outside
// of this package.
func (c *CodeBuilder) write(targetSection, value string) {
	c.builders[targetSection].WriteString(value)
}

func (c *CodeBuilder) writeF(targetSection, format string, args ...any) {
	fmt.Fprintf(c.builders[targetSection], format, args...)
}

// String returns the built string.
func (c *CodeBuilder) String(orderedKeys []string) string {
	return c.Build(orderedKeys).Code
}

// Build builds the code and returns the result.
// You can specify which sections to build (and by what order) by passing `orderedKeys`
// arg to this method; pass nil for building all sections.
func (c *CodeBuilder) Build(orderedKeys []string) *CodeBuildResult {
	if len(orderedKeys) == 0 || (len(orderedKeys) == 1 && orderedKeys[0] == "*") {
		orderedKeys = GetDefaultOrderedSections()
		for key := range c.builders {
			orderedKeys = ssg.AppendUnique(orderedKeys, key)
		}
	} else if len(orderedKeys) == 1 {
		targetBuilder, exists := c.builders[orderedKeys[0]]
		if !exists || targetBuilder == nil {
			return &CodeBuildResult{}
		}
		// If only one section, we still need to process debug info if enabled
		if !c.enableDebugInfo {
			return &CodeBuildResult{Code: targetBuilder.String()}
		}
		// Fallthrough to general processing for debug info
	}

	var result strings.Builder
	var finalDebugInfos []*DebugInfo
	currentLine := 1

	for _, key := range orderedKeys {
		currentSb, exists := c.builders[key]
		if !exists || currentSb == nil || currentSb.Len() == 0 {
			continue
		}
		sectionContent := currentSb.String()
		result.WriteString(sectionContent)
		result.WriteString(c.newLineStr)

		if c.enableDebugInfo {
			sectionDebugInfos := c.debugInfos[key]
			for _, info := range sectionDebugInfos {
				// Calculate the line number relative to the start of the section
				// We need to count newlines in the section content up to the offset
				// This is a bit expensive, but necessary for accuracy
				// Optimization: we could cache line offsets for the section

				// Ensure offset is within bounds (it should be)
				if info.SectionOffset > len(sectionContent) {
					info.SectionOffset = len(sectionContent)
				}

				precedingContent := sectionContent[:info.SectionOffset]
				linesBefore := strings.Count(precedingContent, c.newLineStr)

				info.GeneratedLine = currentLine + linesBefore
				finalDebugInfos = append(finalDebugInfos, info)
			}
		}

		// Update currentLine for the next section
		// Count lines in the section content + 1 for the newline added after the section
		currentLine += strings.Count(sectionContent, c.newLineStr) + 1
	}

	res := &CodeBuildResult{
		Code: result.String(),
	}

	if c.enableDebugInfo {
		jsonData, err := json.Marshal(finalDebugInfos)
		if err == nil {
			res.DebugInfo = string(jsonData)
		}
	}

	return res
}

//---------------------------------------------------------

// valueExists returns true if the value exists in any place.
// this method only cares about the existence, it won't check if
// the value is empty or not.
func (v *codeBuilderVars) valueExists(section, name string) bool {
	if section == "" {
		_, exists := v.globalVars[name]
		return exists
	}

	sectionMap, ok := v.perSections[section]
	if !ok {
		return false
	}

	_, ok = sectionMap[name]
	return ok
}

// getValue returns value of a var by its name either from the section name OR from
// the global vars. Pass empty section name to force getting from global vars.
func (v *codeBuilderVars) getValue(section, name string) string {
	if section == "" {
		return v.globalVars[name]
	}

	sectionMap, ok := v.perSections[section]
	if ok && sectionMap != nil {
		value, ok := sectionMap[name]
		if ok {
			return value
		}
	}

	return v.globalVars[name]
}

func (v *codeBuilderVars) setVar(section, name, value string) {
	sectionMap := v.perSections[section]
	if sectionMap == nil {
		sectionMap = map[string]string{}
		v.perSections[section] = sectionMap
	}

	sectionMap[name] = value
}

func (v *codeBuilderVars) setGlobal(name, value string) {
	v.globalVars[name] = value
}

func (v *codeBuilderVars) deleteVar(section string, names ...string) {
	sectionMap, ok := v.perSections[section]
	if ok && sectionMap != nil {
		for _, current := range names {
			delete(sectionMap, current)
		}
	}
}

func (v *codeBuilderVars) deleteGlobalVar(name string) {
	delete(v.globalVars, name)
}
