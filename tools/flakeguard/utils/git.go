package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// FindChangedFiles executes a git diff against a specified base reference and pipes the output through a user-defined grep command or sequence.
// The baseRef parameter specifies the base git reference for comparison (e.g., "main", "develop").
// The filterCmd parameter should include the full command to be executed after git diff, such as "grep '_test.go$'" or "grep -v '_test.go$' | sort".
func FindChangedFiles(repoPath, baseRef, filterCmd string) ([]string, error) {
	// First command to list files changed between the baseRef and HEAD
	diffCmdStr := fmt.Sprintf("git diff --name-only --diff-filter=AM %s...HEAD", baseRef)
	diffCmd := exec.Command("bash", "-c", diffCmdStr)
	diffCmd.Dir = repoPath

	// Using a buffer to capture stdout and a separate buffer for stderr
	var out bytes.Buffer
	var errBuf bytes.Buffer
	diffCmd.Stdout = &out
	diffCmd.Stderr = &errBuf

	// Running the diff command
	if err := diffCmd.Run(); err != nil {
		return nil, fmt.Errorf("error executing git diff command: %s; error: %w; stderr: %s", diffCmdStr, err, errBuf.String())
	}

	// Check if there are any files listed, if not, return an empty slice
	diffOutput := strings.TrimSpace(out.String())
	if diffOutput == "" {
		return []string{}, nil
	}

	// Second command to filter files using grepCmd
	grepCmdStr := fmt.Sprintf("echo '%s' | %s", diffOutput, filterCmd)
	grepCmd := exec.Command("bash", "-c", grepCmdStr)
	grepCmd.Dir = repoPath

	// Reset buffers for reuse
	out.Reset()
	errBuf.Reset()
	grepCmd.Stdout = &out
	grepCmd.Stderr = &errBuf

	// Running the grep command
	if err := grepCmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 1 {
					// Exit status 1 for grep means no lines matched, which is not an error in this context
					return []string{}, nil
				}
			}
		}
		return nil, fmt.Errorf("error executing grep command: %s; error: %w; stderr: %s", grepCmdStr, err, errBuf.String())
	}

	// Preparing the final list of files
	files := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}

	return files, nil
}
