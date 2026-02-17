package sourcePos

const (
	SourceTabWidth = 4

	// MaxErrorSourceLineBeforeLen controls how many characters to show
	// before the culprit in formatted errors.
	MaxErrorSourceLineBeforeLen = 20

	// MaxErrorSourceLineAfterLen controls how many characters to show
	// after the culprit in formatted errors.
	MaxErrorSourceLineAfterLen = 40

	// SourceErrorEllipsis is the marker used when truncating long source lines.
	SourceErrorEllipsis = "..."
)
