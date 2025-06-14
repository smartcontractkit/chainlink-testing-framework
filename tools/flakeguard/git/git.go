package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/utils"
)

// FindChangedFiles executes a git diff against a specified base reference and pipes the output through a user-defined grep command or sequence.
// The baseRef parameter specifies the base git reference for comparison (e.g., "main", "develop").
// The filterCmd parameter should include the full command to be executed after git diff, such as "grep '_test.go$'" or "grep -v '_test.go$' | sort".
func FindChangedFiles(rootGoModPath, baseRef, filterCmd string) ([]string, error) {
	// Find directories containing a go.mod file and build an exclusion string
	excludeStr, err := buildExcludeStringForGoModDirs(rootGoModPath)
	if err != nil {
		return nil, fmt.Errorf("error finding go.mod directories: %w", err)
	}

	// First command to list files changed between the baseRef and HEAD, excluding specified paths
	diffCmdStr := fmt.Sprintf("git diff --name-only --diff-filter=AM %s...HEAD -- %s %s", baseRef, rootGoModPath, excludeStr)
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

	// Check if there are any files listed; if not, return an empty slice
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

// buildExcludeStringForGoModDirs searches the given root directory for subdirectories
// containing a go.mod file and returns a formatted string to exclude those directories
// (except the root directory if it contains a go.mod file) from git diff.
func buildExcludeStringForGoModDirs(rootGoModPath string) (string, error) {
	var excludeStr string

	err := filepath.Walk(rootGoModPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "go.mod" {
			dir := filepath.Dir(path)
			// Skip excluding the root directory if go.mod is found there
			if dir != rootGoModPath {
				excludeStr += fmt.Sprintf("':(exclude)%s/**' ", dir)
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return excludeStr, nil
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

// GetOwnerRepoDefaultBranchFromLocalRepo returns the owner, repo name, and default branch of a local git repository.
// It uses the origin remote URL to determine the owner and repo name, and the default branch is determined from the
// refs/remotes/origin/HEAD reference.
func GetOwnerRepoDefaultBranchFromLocalRepo(repoPath string) (owner, repoName, defaultBranch string, err error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", "", "", err
	}

	// Get remote URL (origin)
	remotes, err := repo.Remotes()
	if err != nil {
		return "", "", "", err
	}
	var originURL string
	for _, remote := range remotes {
		if remote.Config().Name == "origin" && len(remote.Config().URLs) > 0 {
			originURL = remote.Config().URLs[0]
			break
		}
	}
	if originURL == "" {
		return "", "", "", fmt.Errorf("origin remote not found")
	}

	// Parse owner and repo from URL
	originURL = strings.TrimSuffix(originURL, ".git")
	var path string
	if strings.Contains(originURL, "@github.com:") {
		parts := strings.SplitN(originURL, ":", 2)
		if len(parts) == 2 {
			path = parts[1]
		}
	} else if strings.HasPrefix(originURL, "https://") {
		parts := strings.SplitN(originURL, "github.com/", 2)
		if len(parts) == 2 {
			path = parts[1]
		}
	}
	if path == "" {
		return "", "", "", fmt.Errorf("could not parse remote URL: %s", originURL)
	}
	segments := strings.Split(path, "/")
	if len(segments) != 2 {
		return "", "", "", fmt.Errorf("unexpected path format: %s", path)
	}
	owner, repoName = segments[0], segments[1]

	// Find default branch from refs/remotes/origin/HEAD
	refs, err := repo.References()
	if err != nil {
		return "", "", "", err
	}
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsRemote() && ref.Name().String() == "refs/remotes/origin/HEAD" {
			target := ref.Target().String()
			parts := strings.Split(target, "/")
			if len(parts) > 0 {
				defaultBranch = parts[len(parts)-1]
			}
		}
		return nil
	})
	if err != nil {
		return "", "", "", err
	}
	if defaultBranch == "" {
		return "", "", "", fmt.Errorf("could not determine default branch")
	}

	return owner, repoName, defaultBranch, nil
}
