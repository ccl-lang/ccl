package jsGen_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type RunJSOptions struct {
	TargetPath    string
	RunnerContent string
}

// RunJSProject runs a JavaScript project located at targetPath with the provided runnerContent script.
func RunJSProject(opts *RunJSOptions) (string, error) {
	// Write the runner script
	if err := os.WriteFile(filepath.Join(opts.TargetPath, "verify.js"), []byte(opts.RunnerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write verify.js: %v", err)
	}

	// Create package.json
	packageJson := `{
  "name": "js-test",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "type": "module",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "keywords": [],
  "author": "",
  "license": "ISC"
}`
	if err := os.WriteFile(filepath.Join(opts.TargetPath, "package.json"), []byte(packageJson), 0644); err != nil {
		return "", fmt.Errorf("failed to write package.json: %v", err)
	}

	// Run the verification script using node
	cmd := exec.Command("node", "verify.js")
	cmd.Dir = opts.TargetPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run verification script: %v\nOutput:\n%s", err, output)
	}

	return string(output), nil
}
