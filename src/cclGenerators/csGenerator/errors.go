package csGenerator

import "fmt"

type NamespaceDetectionError struct {
	Path string
	Err  error
}

func (e *NamespaceDetectionError) Error() string {
	return fmt.Sprintf("failed to detect namespace from path %s: %v", e.Path, e.Err)
}
