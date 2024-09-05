package breakingchanges

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectBreakingChanges(t *testing.T) {
	t.Skip("This test describe what should be done to test the changes, however, it is infeasible to test because you need a published commit for gorelease to work")
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "testrepo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Initialize a Git repository
	runCommand(t, tempDir, "git", "init")
	runCommand(t, tempDir, "git", "config", "user.name", "test")
	runCommand(t, tempDir, "git", "config", "user.email", "test@example.com")

	// Configure Git to not use any signing key
	runCommand(t, tempDir, "git", "config", "commit.gpgSign", "false")
	runCommand(t, tempDir, "git", "config", "tag.gpgSign", "false")

	// Create a simple Go program with an external method
	mainGo := `package main

import "fmt"

func ExternalMethod(a int, b int) {
	fmt.Println(a, b)
}

func main() {
	ExternalMethod(1, 2)
}
`
	writeFile(t, tempDir, "main.go", mainGo)

	// Initialize a Go module within the scope of the test
	runCommand(t, tempDir, "go", "mod", "init", "github.com/testtest/breaking_changes")

	// Add and commit the initial version
	runCommand(t, tempDir, "git", "add", ".")
	runCommand(t, tempDir, "git", "commit", "-m", "Initial commit")
	initialTag := "github.com/testtest/breaking_changes/v1.0.0"
	runCommand(t, tempDir, "git", "tag", initialTag)

	// Modify the ExternalMethod to introduce a breaking change
	mainGoV2 := `package main

import "fmt"

func ExternalMethod(a int, b int, c int) {
	fmt.Println(a, b, c)
}

func main() {
	ExternalMethod(1, 2, 3)
}
`
	writeFile(t, tempDir, "main.go", mainGoV2)

	// Add and commit the breaking change
	runCommand(t, tempDir, "git", "add", ".")
	runCommand(t, tempDir, "git", "commit", "-m", "Introduce breaking change")
	runCommand(t, tempDir, "git", "tag", "github.com/testtest/breaking_changes/v1.1.0")

	// Capture the output of DetectBreakingChanges
	var output bytes.Buffer
	runDetectBreakingChanges(tempDir, initialTag, &output)

	// Scan output lines for "incompatible changes"
	scanner := bufio.NewScanner(&output)
	foundIncompatibleChanges := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "incompatible changes") {
			foundIncompatibleChanges = true
			break
		}
	}

	if !foundIncompatibleChanges {
		t.Errorf("Expected output to contain 'incompatible changes', but got:\n%s", output.String())
	}
}

// Helper function to run commands
func runCommand(t *testing.T, dir string, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command %s %v failed with error: %v and output: %s", name, args, err, string(out))
	}
}

// Helper function to write files
func writeFile(t *testing.T, dir, filename, content string) {
	filePath := filepath.Join(dir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write file %s: %v", filePath, err)
	}
}

// Helper function to run DetectBreakingChanges and capture output
func runDetectBreakingChanges(dir, baseTag string, output *bytes.Buffer) {
	// Temporarily change working directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(dir)

	// Run the gorelease command with the base tag
	cmd := exec.Command("gorelease", "-base", baseTag)
	cmd.Dir = dir
	cmd.Stdout = output
	cmd.Stderr = output
	cmd.Run()
}
