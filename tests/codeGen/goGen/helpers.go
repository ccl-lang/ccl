package goGen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunGoProject runs a Go project located at targetPath with the provided runnerContent script.
func RunGoProject(opts *RunGoOptions) (string, error) {
	// Create go.mod
	goModContent := `module ccl_test_generated

go 1.25
`
	if err := os.WriteFile(filepath.Join(opts.TargetPath, "go.mod"), []byte(goModContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write go.mod: %v", err)
	}

	if err := os.WriteFile(filepath.Join(opts.TargetPath, "main.go"), []byte(opts.RunnerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write main.go: %v", err)
	}

	// Run the generated code
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = opts.TargetPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run generated code: %v\nOutput:\n%s", err, output)
	}

	return string(output), nil
}
