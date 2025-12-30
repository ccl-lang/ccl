package pyGen_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type RunPythonOptions struct {
	TargetPath    string
	RunnerContent string
}

// RunPythonProject runs a Python project located at targetPath with the provided runnerContent script.
func RunPythonProject(opts *RunPythonOptions) (string, error) {
	// Write the runner script
	if err := os.WriteFile(filepath.Join(opts.TargetPath, "verify.py"), []byte(opts.RunnerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write verify.py: %v", err)
	}

	// Run the verification script
	cmd := exec.Command("python", "verify.py")
	cmd.Dir = opts.TargetPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run verification script: %v\nOutput:\n%s", err, output)
	}

	return string(output), nil
}
