package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/utils"
)

// FindChangedFiles executes a git diff against a specified base reference and pipes the output through a user-defined grep command or sequence.
// The baseRef parameter specifies the base git reference for comparison (e.g., "main", "develop").
// The filterCmd parameter should include the full command to be executed after git diff, such as "grep '_test.go$'" or "grep -v '_test.go$' | sort".
func FindChangedFiles(baseRef, filterCmd string, excludePaths []string) ([]string, error) {
	// Constructing the exclusion part of the git command
	excludeStr := ""
	for _, path := range excludePaths {
		excludeStr += fmt.Sprintf("':(exclude)%s' ", path)
	}

	// First command to list files changed between the baseRef and HEAD, excluding specified paths
	diffCmdStr := fmt.Sprintf("git diff --name-only --diff-filter=AM %s...HEAD %s", baseRef, excludeStr)
	diffCmd := exec.Command("bash", "-c", diffCmdStr)

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

func Diff(baseBranch string) (*utils.CmdOutput, error) {
	return utils.ExecuteCmd("git", "diff", "--name-only", baseBranch)
}

func ModDiff(baseBranch, projectPath string) (*utils.CmdOutput, error) {
	return utils.ExecuteCmd("git", "diff", baseBranch, "--unified=0", "--", filepath.Join(projectPath, "go.mod"))
}

func GetGoModChangesFromDiff(lines bytes.Buffer) ([]string, error) {
	changedLines := strings.Split(lines.String(), "\n")

	// Filter out lines that do not indicate package changes
	var packages []string
	for _, line := range changedLines {
		if strings.HasPrefix(line, "+") {
			// ignore comments or empty lines (e.g., not relevant)
			if strings.HasPrefix(line, "+ ") || strings.HasPrefix(line, "+++ ") {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) > 1 {
				// The second field should contains the module path
				packages = append(packages, fields[1])
			}
		}
	}

	return packages, nil
}

// GetChangedGoPackagesFromDiff identifies the Go packages affected by changes in a Git diff output.
// It analyzes a buffer containing the output of a 'git diff' command, filtering the list of changed
// files to determine which Go packages have been affected based on the project's file map.
//
// Parameters:
//   - out: A bytes.Buffer containing the 'git diff' command output. This output should list the
//     changed files, one per line.
//   - projectPath: The root directory of the project within the repository. This parameter is
//     used to filter files based on their paths. [Note: This functionality is currently commented out,
//     ensure to uncomment the related line if you decide to use it.]
//   - excludes: A slice of strings representing paths to exclude from the analysis. This can be useful
//     to ignore changes in certain directories or files that are not relevant to the package analysis.
//   - fileMap: A map where the key is a file path and the value is a slice of strings representing the
//     Go import paths of the packages that file belongs to. This map is used to map changed files
//     to their respective packages.
//
// Returns:
//   - A slice of strings representing the unique Go packages that have changes. These packages are
//     identified by their import paths.
//   - An error, which is nil in the current implementation but can be used to return errors encountered
//     during the execution of the function.
func GetChangedGoPackagesFromDiff(out bytes.Buffer, projectPath string, excludes []string, fileMap map[string][]string) ([]string, error) {
	changedFiles := strings.Split(out.String(), "\n")

	// Filter out non-Go files and directories and embeds
	changedPackages := make(map[string]bool)
	for _, file := range changedFiles {
		if file == "" || shouldExclude(excludes, file) {
			continue
		}

		// if the changed file is in the fileMap then we add it to the changed packages
		for _, importPath := range fileMap[file] {
			changedPackages[importPath] = true
		}
	}

	// Convert map keys to slice
	var packages []string
	for pkg := range changedPackages {
		packages = append(packages, pkg)
	}

	return packages, nil
}

func shouldExclude(excludes []string, item string) bool {
	for _, v := range excludes {
		if strings.HasPrefix(item, v) {
			return true
		}
	}
	return false
}
