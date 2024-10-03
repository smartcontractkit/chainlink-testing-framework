package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// FindChangedFiles executes a git diff against a specified base reference and pipes the output through a user-defined grep command or sequence.
// The baseRef parameter specifies the base git reference for comparison (e.g., "main", "develop").
// The grepCmd parameter should include the full command to be executed after git diff, such as "grep '_test.go$'" or "grep -v '_test.go$' | sort".
func FindChangedFiles(repoPath, baseRef, grepCmd string) ([]string, error) {
	cmdString := fmt.Sprintf("git diff --name-only --diff-filter=AM %s...HEAD | %s", baseRef, grepCmd)
	cmd := exec.Command("bash", "-c", cmdString)
	cmd.Dir = repoPath
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error executing command: %w", err)
	}
	files := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}
	return files, nil
}
