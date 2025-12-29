package codeBuilder

import "fmt"

// Indent increases the indentation level.
func (c *CodeBuilder) Indent() *CodeBuilder {
	c.indentation++
	return c
}

// Unindent decreases the indentation level.
func (c *CodeBuilder) Unindent() *CodeBuilder {
	if c.indentation > 0 {
		c.indentation--
	}
	return c
}

// UnindentLine decreases the indentation level and adds a new line.
func (c *CodeBuilder) UnindentLine() *CodeBuilder {
	c.Unindent()
	c.NewLine()
	return c
}

// WriteStr writes a string with the current indentation.
// If you don't want indentation, use AppendStr instead.
func (c *CodeBuilder) WriteStr(s string) *CodeBuilder {
	c.writeIndentation()
	c.sb.WriteString(s)
	return c
}

// AppendStr appends a string without adding indentation.
func (c *CodeBuilder) AppendStr(s string) *CodeBuilder {
	c.sb.WriteString(s)
	return c
}

// WriteLine writes a line with the current indentation.
// If you don't want indentation, use AppendLine instead.
func (c *CodeBuilder) WriteLine(s string) *CodeBuilder {
	c.writeIndentation()
	c.sb.WriteString(s)
	c.sb.WriteString("\n")
	return c
}

// AppendLine appends a line without adding indentation.
func (c *CodeBuilder) AppendLine(s string) *CodeBuilder {
	c.sb.WriteString(s)
	c.sb.WriteString("\n")
	return c
}

// WriteLinef writes a formatted line with the current indentation.
func (c *CodeBuilder) WriteLinef(format string, args ...any) *CodeBuilder {
	c.writeIndentation()
	fmt.Fprintf(&c.sb, format, args...)
	c.sb.WriteString("\n")
	return c
}

// NewLine adds a new line.
func (c *CodeBuilder) NewLine() *CodeBuilder {
	c.sb.WriteString("\n")
	return c
}

// writeIndentation writes the current indentation to the string builder.
func (c *CodeBuilder) writeIndentation() *CodeBuilder {
	for i := 0; i < c.indentation; i++ {
		c.sb.WriteString("\t")
	}
	return c
}

// String returns the built string.
func (c *CodeBuilder) String() string {
	return c.sb.String()
}
