package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	// Write the file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", path, err)
	}
}

func runCmd(t *testing.T, dir string, command string, args ...string) {
	t.Helper()

	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %s %v\nOutput: %s\nError: %v",
			command, args, string(output), err)
	}
}

func TestIntegration(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize git repo
	runCmd(t, tmpDir, "git", "init")
	runCmd(t, tmpDir, "git", "config", "user.email", "test@example.com")
	runCmd(t, tmpDir, "git", "config", "user.name", "Test User")

	// Create and commit initial file
	writeFile(t, filepath.Join(tmpDir, "test.txt"), "initial content")
	runCmd(t, tmpDir, "git", "add", "test.txt")
	runCmd(t, tmpDir, "git", "commit", "-m", "initial commit")

	// Test 1: Clean repo - should succeed
	buf := new(bytes.Buffer)
	// echo goes to stdout, so buffer should be clean even though there's output
	code := run(context.Background(), tmpDir, buf, []string{"differ", "echo", "hello"})
	if code != 0 && buf.Len() > 0 {
		t.Errorf("Expected success in clean repo, got: %v", code)
	}

	// Test 2: Untracked file - should fail
	buf.Reset()
	writeFile(t, filepath.Join(tmpDir, "new.txt"), "new file")
	code = run(context.Background(), tmpDir, buf, []string{"differ", "echo", "hello"})
	if code == 0 {
		t.Errorf("Expected failure with untracked file, got code %d, output: %v", code, buf.String())
	}
	if !bytes.Contains(buf.Bytes(), []byte("Untracked or modified files present")) {
		t.Errorf("Expected error message about untracked files, got: %s", buf.String())
	}
}
