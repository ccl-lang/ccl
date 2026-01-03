package tsGen_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RunTSOptions struct {
	TargetPath      string
	RunnerContent   string
	RequiredModules []string
}

var baseRequiredModules = []string{
	"ts-node",
	"typescript",
}

// checkAndInstallGlobalModules checks if required global npm modules are installed, and installs them if missing.
func checkAndInstallGlobalModules(requiredModules []string) error {
	cmd := exec.Command("npm", "list", "-g", "--depth=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to list global npm packages: %v", err)
	}

	// Check if required modules are installed
	for _, module := range requiredModules {
		if !strings.Contains(string(output), " "+module+"@") {
			// Install the missing module
			cmd := exec.Command("npm", "install", "-g", module)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install global npm module %q: %v", module, err)
			}
		}
	}

	return nil
}

// RunTSProject runs a TypeScript project located at targetPath with the provided runnerContent script.
func RunTSProject(opts *RunTSOptions) (string, error) {
	// Write the runner script
	if err := os.WriteFile(filepath.Join(opts.TargetPath, "verify.ts"), []byte(opts.RunnerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write verify.ts: %v", err)
	}

	// Create package.json
	packageJson := `{
  "name": "ts-test",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
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

	opts.RequiredModules = append(opts.RequiredModules, baseRequiredModules...)
	// Install dependencies if they aren't installed already
	if err := checkAndInstallGlobalModules(opts.RequiredModules); err != nil {
		return "", fmt.Errorf("failed to install global npm modules: %v", err)
	}

	// Run the verification script using npx ts-node
	cmd := exec.Command("npx", "ts-node", "--compiler-options", `{"module":"commonjs"}`, "verify.ts")
	cmd.Dir = opts.TargetPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run verification script: %v\nOutput:\n%s", err, output)
	}

	return string(output), nil
}
