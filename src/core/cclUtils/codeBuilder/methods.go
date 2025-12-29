package codeBuilder

import "fmt"

func NewCodeBuilder() *CodeBuilder {
	return &CodeBuilder{}
}

func (c *CodeBuilder) Indent() {
	c.indentation++
}

func (c *CodeBuilder) Unindent() {
	if c.indentation > 0 {
		c.indentation--
	}
}

// WriteStr writes a string with the current indentation.
// If you don't want indentation, use AppendStr instead.
func (c *CodeBuilder) WriteStr(s string) {
	c.writeIndentation()
	c.sb.WriteString(s)
}

// AppendStr appends a string without adding indentation.
func (c *CodeBuilder) AppendStr(s string) {
	c.sb.WriteString(s)
}

// WriteLine writes a line with the current indentation.
// If you don't want indentation, use AppendLine instead.
func (c *CodeBuilder) WriteLine(s string) {
	c.writeIndentation()
	c.sb.WriteString(s)
	c.sb.WriteString("\n")
}

// AppendLine appends a line without adding indentation.
func (c *CodeBuilder) AppendLine(s string) {
	c.sb.WriteString(s)
	c.sb.WriteString("\n")
}

func (c *CodeBuilder) WriteLinef(format string, args ...any) {
	c.writeIndentation()
	fmt.Fprintf(&c.sb, format, args...)
	c.sb.WriteString("\n")
}

func (c *CodeBuilder) NewLine() {
	c.sb.WriteString("\n")
}

func (c *CodeBuilder) writeIndentation() {
	for i := 0; i < c.indentation; i++ {
		c.sb.WriteString("\t")
	}
}

func (c *CodeBuilder) String() string {
	return c.sb.String()
}
