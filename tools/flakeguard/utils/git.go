package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// FindChangedFiles executes a git diff and pipes the output through a user-defined grep command or sequence.
// The grepCmd parameter should include the full command to be executed after git diff, such as "grep '_test.go$'" or "grep -v '_test.go$' | sort".
func FindChangedFiles(repoPath, grepCmd string) ([]string, error) {
	cmdString := fmt.Sprintf("git diff --name-only --diff-filter=AM develop...HEAD | %s", grepCmd)
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
