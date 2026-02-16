package codeBuilder

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

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
		panic("CodeBuilder: illegal usage of CodeBuilder without initiating a section")
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
	c.builders[SectionImports].WriteString(importLine + c.newLineStr)
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
	c.builders[SectionCommentHeaders].WriteString(comment)
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
	c.builders[SectionHeaders].WriteString(strValue)
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
	c.builders[SectionDeclareNamespace].WriteString(strValue)
	return c
}

// BeginSection begins a new section or switches to an existing one.
// IMPORTANT NOTE: As long as you are in this section, there is a lock held on this object,
// to avoid deadlocks, you SHOULD call EndSection() when you are done with this section.
func (c *CodeBuilder) BeginSection(section string) *CodeBuilder {
	if section == "" {
		panic("CodeBuilder: section name cannot be empty")
	}

	if c.currentSection != "" {
		if c.currentSection == section {
			return c
		}
		panic("CodeBuilder: cannot begin a new section without ending the previous one")
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
		panic("CodeBuilder: cannot end a section when no section is active")
	}

	c.currentSection = ""
	c.mut.Unlock()
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
	c.builders[c.currentSection].WriteString(c.newLineStr)
	return c
}

// WriteStr writes a string with the current indentation.
// If you don't want indentation, use AppendStr instead.
func (c *CodeBuilder) WriteStr(s string) *CodeBuilder {
	c.checkSection()

	c.addDebugInfo(2)
	c.writeIndentation()
	c.builders[c.currentSection].WriteString(s)
	return c
}

// AppendStr appends a string without adding indentation.
func (c *CodeBuilder) AppendStr(s string) *CodeBuilder {
	c.checkSection()

	c.addDebugInfo(2)
	c.checkSection()
	c.builders[c.currentSection].WriteString(s)
	return c
}

// WriteLine writes a line with the current indentation.
// If you don't want indentation, use AppendLine instead.
func (c *CodeBuilder) WriteLine(s string) *CodeBuilder {
	c.checkSection()

	c.addDebugInfo(2)
	c.writeIndentation()
	c.builders[c.currentSection].WriteString(s)
	c.builders[c.currentSection].WriteString(c.newLineStr)
	return c
}

// AppendLine appends a line without adding indentation.
func (c *CodeBuilder) AppendLine(s string) *CodeBuilder {
	c.addDebugInfo(2)
	c.checkSection()
	c.builders[c.currentSection].WriteString(s)
	c.builders[c.currentSection].WriteString(c.newLineStr)
	return c
}

// WriteLinef writes a formatted line with the current indentation.
func (c *CodeBuilder) WriteLinef(format string, args ...any) *CodeBuilder {
	c.checkSection()

	c.addDebugInfo(2)
	c.writeIndentation()
	fmt.Fprintf(c.builders[c.currentSection], format, args...)
	c.builders[c.currentSection].WriteString(c.newLineStr)
	return c
}

// NewLine adds a new line.
func (c *CodeBuilder) NewLine() *CodeBuilder {
	c.addDebugInfo(2)
	c.checkSection()
	c.builders[c.currentSection].WriteString(c.newLineStr)
	return c
}

// writeIndentation writes the current indentation to the string builder.
// NOTE: since this is an internal method, it does NOT check for section or
// debug info; so the caller must ensure those are handled appropriately.
func (c *CodeBuilder) writeIndentation() *CodeBuilder {
	for i := 0; i < c.indentations[c.currentSection]; i++ {
		c.builders[c.currentSection].WriteString(c.indentationStr)
	}
	return c
}

// String returns the built string.
func (c *CodeBuilder) String(orderedKeys []string) string {
	return c.Build(orderedKeys).Code
}

// Build builds the code and returns the result.
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
