package git

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v72/github"
	"github.com/shurcooL/githubv4"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/utils"
	"golang.org/x/oauth2"
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

// MakeSignedCommit adds all changes to a repo and creates a signed commit for GitHub
func MakeSignedCommit(repoPath, commitMessage, branch, githubToken string) (string, error) {
	tok := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	token := oauth2.NewClient(context.Background(), tok)
	graphqlClient := githubv4.NewClient(token)

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	// Inspired by https://github.com/planetscale/ghcommit/tree/main

	// process added / modified files:
	worktree, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	// Get the status of all files in the worktree
	status, err := worktree.Status()
	if err != nil {
		return "", err
	}

	additions := []githubv4.FileAddition{}
	deletions := []githubv4.FileDeletion{}

	// Process each file based on its status
	for filePath, fileStatus := range status {
		switch fileStatus.Staging {
		case git.Added, git.Modified:
			// File is added or modified - add to additions
			enc, err := base64EncodeFile(filepath.Join(repoPath, filePath))
			if err != nil {
				return "", err
			}
			additions = append(additions, githubv4.FileAddition{
				Path:     githubv4.String(filePath),
				Contents: githubv4.Base64String(enc),
			})
		case git.Deleted:
			// File is deleted - add to deletions
			deletions = append(deletions, githubv4.FileDeletion{
				Path: githubv4.String(filePath),
			})
		}

		// Also check worktree status (unstaged changes)
		switch fileStatus.Worktree {
		case git.Modified:
			// Only add if not already processed from staging
			if fileStatus.Staging != git.Added && fileStatus.Staging != git.Modified {
				enc, err := base64EncodeFile(filepath.Join(repoPath, filePath))
				if err != nil {
					return "", err
				}
				additions = append(additions, githubv4.FileAddition{
					Path:     githubv4.String(filePath),
					Contents: githubv4.Base64String(enc),
				})
			}
		case git.Deleted:
			// Only add if not already processed from staging
			if fileStatus.Staging != git.Deleted {
				deletions = append(deletions, githubv4.FileDeletion{
					Path: githubv4.String(filePath),
				})
			}
		}
	}

	var m struct {
		CreateCommitOnBranch struct {
			Commit struct {
				URL       string `graphql:"url"`
				OID       string `graphql:"oid"`
				Additions int    `graphql:"additions"`
				Deletions int    `graphql:"deletions"`
			}
		} `graphql:"createCommitOnBranch(input:$input)"`
	}

	splitMsg := strings.SplitN(commitMessage, "\n", 2)
	headline := splitMsg[0]
	body := ""
	if len(splitMsg) > 1 {
		body = splitMsg[1]
	}

	owner, repoName, _, err := GetOwnerRepoDefaultBranchFromLocalRepo(repoPath)
	if err != nil {
		return "", err
	}

	// Get HEAD reference to get the current commit hash
	headRef, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD reference: %w", err)
	}
	expectedHeadOid := headRef.Hash().String()
	// create the $input struct for the graphQL createCommitOnBranch mutation request:
	input := githubv4.CreateCommitOnBranchInput{
		Branch: githubv4.CommittableBranch{
			RepositoryNameWithOwner: githubv4.NewString(githubv4.String(fmt.Sprintf("%s/%s", owner, repoName))),
			BranchName:              githubv4.NewString(githubv4.String(branch)),
		},
		Message: githubv4.CommitMessage{
			Headline: githubv4.String(headline),
			Body:     githubv4.NewString(githubv4.String(body)),
		},
		FileChanges: &githubv4.FileChanges{
			Additions: &additions,
			Deletions: &deletions,
		},
		ExpectedHeadOid: githubv4.GitObjectID(expectedHeadOid),
	}

	if err := graphqlClient.Mutate(context.Background(), &m, input, nil); err != nil {
		return "", err
	}

	return m.CreateCommitOnBranch.Commit.OID, nil
}

func base64EncodeFile(path string) (string, error) {
	in, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer in.Close() // nolint: errcheck

	buf := bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)

	if _, err := io.Copy(encoder, in); err != nil {
		return "", err
	}
	if err := encoder.Close(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GitHubFileInfo represents file information from GitHub API
type GitHubFileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"` // "file" or "dir"
	DownloadURL string `json:"download_url"`
	SHA         string `json:"sha"`
}

// GitHubRepoStructure represents the analyzed repository structure
type GitHubRepoStructure struct {
	GoModDirs    []string            // Directories containing go.mod files
	TestFiles    map[string][]string // Package path -> list of test files
	PackageFiles map[string][]string // Package path -> list of go files
}

// DiscoverRepoStructureViaGitHub analyzes repository structure using GitHub API
// This replaces the need for local filesystem operations
func DiscoverRepoStructureViaGitHub(owner, repo, ref, githubToken string) (*GitHubRepoStructure, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	structure := &GitHubRepoStructure{
		GoModDirs:    []string{},
		TestFiles:    make(map[string][]string),
		PackageFiles: make(map[string][]string),
	}

	// Recursively walk the repository structure
	err := walkGitHubDirectory(ctx, client, owner, repo, ref, "", structure)
	if err != nil {
		return nil, fmt.Errorf("failed to walk repository structure: %w", err)
	}

	return structure, nil
}

// walkGitHubDirectory recursively walks through GitHub repository directories
func walkGitHubDirectory(ctx context.Context, client *github.Client, owner, repo, ref, path string, structure *GitHubRepoStructure) error {
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		return fmt.Errorf("failed to get directory contents for %s: %w", path, err)
	}

	var goFiles []string
	var testFiles []string

	for _, content := range directoryContent {
		if content.GetType() == "file" {
			fileName := content.GetName()
			filePath := content.GetPath()

			// Check for go.mod files
			if fileName == "go.mod" {
				structure.GoModDirs = append(structure.GoModDirs, path)
			}

			// Check for Go files
			if strings.HasSuffix(fileName, ".go") {
				if strings.HasSuffix(fileName, "_test.go") {
					testFiles = append(testFiles, filePath)
				} else {
					goFiles = append(goFiles, filePath)
				}
			}
		} else if content.GetType() == "dir" {
			// Recursively walk subdirectories
			err := walkGitHubDirectory(ctx, client, owner, repo, ref, content.GetPath(), structure)
			if err != nil {
				return err
			}
		}
	}

	// If this directory has Go files, determine the package path
	if len(goFiles) > 0 || len(testFiles) > 0 {
		// For simplicity, use the directory path as package identifier
		// In a real implementation, you might want to parse the package declaration
		packagePath := path
		if packagePath == "" {
			packagePath = "." // root package
		}

		if len(goFiles) > 0 {
			structure.PackageFiles[packagePath] = goFiles
		}
		if len(testFiles) > 0 {
			structure.TestFiles[packagePath] = testFiles
		}
	}

	return nil
}

// GetFileContentsFromGitHub fetches file contents from GitHub
func GetFileContentsFromGitHub(owner, repo, ref, path, githubToken string) (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		return "", fmt.Errorf("failed to get file contents for %s: %w", path, err)
	}

	if fileContent == nil {
		return "", fmt.Errorf("file content is nil for %s", path)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to decode file content for %s: %w", path, err)
	}

	return content, nil
}

// FindPackageForTest finds which package a test belongs to using GitHub API
func FindPackageForTest(owner, repo, ref, testPackageImportPath, testName, githubToken string, structure *GitHubRepoStructure) (string, []string, error) {
	// Convert import path to directory path
	// This is a simplified approach - you might need more sophisticated logic
	packageDir := strings.ReplaceAll(testPackageImportPath, "/", "/")

	// Find test files in the package directory
	testFiles, exists := structure.TestFiles[packageDir]
	if !exists {
		// Try alternative mappings or search through all test files
		for dir, files := range structure.TestFiles {
			// Check if this directory might contain the package we're looking for
			if strings.Contains(dir, packageDir) || strings.HasSuffix(testPackageImportPath, filepath.Base(dir)) {
				testFiles = files
				packageDir = dir
				break
			}
		}
	}

	if len(testFiles) == 0 {
		return "", nil, fmt.Errorf("no test files found for package %s", testPackageImportPath)
	}

	return packageDir, testFiles, nil
}

// CreateBranchOnGitHub creates a new branch on GitHub from the base branch
func CreateBranchOnGitHub(owner, repo, branchName, baseBranch, githubToken string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get the reference of the base branch
	baseRef, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+baseBranch)
	if err != nil {
		return fmt.Errorf("failed to get base branch reference: %w", err)
	}

	// Create new branch reference
	newRef := &github.Reference{
		Ref: github.Ptr("refs/heads/" + branchName),
		Object: &github.GitObject{
			SHA: baseRef.Object.SHA,
		},
	}

	_, _, err = client.Git.CreateRef(ctx, owner, repo, newRef)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

// CommitFilesToGitHub commits multiple files to a GitHub branch using GraphQL API
func CommitFilesToGitHub(owner, repo, branchName string, files map[string]string, commitMsg, githubToken string) (string, error) {
	tok := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	token := oauth2.NewClient(context.Background(), tok)
	graphqlClient := githubv4.NewClient(token)

	// Prepare file additions
	additions := []githubv4.FileAddition{}
	for filePath, content := range files {
		// Base64 encode the content
		encoded := base64.StdEncoding.EncodeToString([]byte(content))
		additions = append(additions, githubv4.FileAddition{
			Path:     githubv4.String(filePath),
			Contents: githubv4.Base64String(encoded),
		})
	}

	var m struct {
		CreateCommitOnBranch struct {
			Commit struct {
				URL string `graphql:"url"`
				OID string `graphql:"oid"`
			}
		} `graphql:"createCommitOnBranch(input:$input)"`
	}

	splitMsg := strings.SplitN(commitMsg, "\n", 2)
	headline := splitMsg[0]
	body := ""
	if len(splitMsg) > 1 {
		body = splitMsg[1]
	}

	// Get current HEAD SHA of the branch
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	restClient := github.NewClient(tc)

	branchRef, _, err := restClient.Git.GetRef(ctx, owner, repo, "refs/heads/"+branchName)
	if err != nil {
		return "", fmt.Errorf("failed to get branch reference: %w", err)
	}
	expectedHeadOid := branchRef.Object.GetSHA()

	// Create the GraphQL input
	input := githubv4.CreateCommitOnBranchInput{
		Branch: githubv4.CommittableBranch{
			RepositoryNameWithOwner: githubv4.NewString(githubv4.String(fmt.Sprintf("%s/%s", owner, repo))),
			BranchName:              githubv4.NewString(githubv4.String(branchName)),
		},
		Message: githubv4.CommitMessage{
			Headline: githubv4.String(headline),
			Body:     githubv4.NewString(githubv4.String(body)),
		},
		FileChanges: &githubv4.FileChanges{
			Additions: &additions,
		},
		ExpectedHeadOid: githubv4.GitObjectID(expectedHeadOid),
	}

	if err := graphqlClient.Mutate(context.Background(), &m, input, nil); err != nil {
		return "", fmt.Errorf("failed to create commit: %w", err)
	}

	return m.CreateCommitOnBranch.Commit.OID, nil
}
