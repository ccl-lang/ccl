package tsGen_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type RunTSOptions struct {
	TargetPath      string
	RunnerContent   string
	RequiredModules []string
}

// RunTSProject runs a TypeScript project located at targetPath with the provided runnerContent script.
func RunTSProject(opts *RunTSOptions) (string, error) {
	// Write the runner script
	if err := os.WriteFile(filepath.Join(opts.TargetPath, "verify.ts"), []byte(opts.RunnerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write verify.ts: %v", err)
	}

	// Pin the local toolchain: ts-node requires the JavaScript compiler API absent in TypeScript 7.
	//TODO: eventually move to ts7 without relying on ts-node at all
	packageJson := `{
  "name": "ts-test",
  "private": true,
  "devDependencies": {
    "ts-node": "10.9.2",
    "typescript": "5.9.3"
  }
}`
	if err := os.WriteFile(filepath.Join(opts.TargetPath, "package.json"), []byte(packageJson), 0644); err != nil {
		return "", fmt.Errorf("failed to write package.json: %v", err)
	}

	installArgs := append([]string{"install", "--include=dev", "--no-package-lock", "--no-audit", "--no-fund"}, opts.RequiredModules...)
	installCmd := exec.Command("npm", installArgs...)
	installCmd.Dir = opts.TargetPath
	if output, err := installCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to install TypeScript test dependencies: %v\nOutput:\n%s", err, output)
	}

	// Invoke the local runner directly so global packages cannot override the pinned versions.
	cmd := exec.Command("node", filepath.Join("node_modules", "ts-node", "dist", "bin.js"), "--compiler-options", `{"module":"commonjs"}`, "verify.ts")
	cmd.Dir = opts.TargetPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run verification script: %v\nOutput:\n%s", err, output)
	}

	return string(output), nil
}
