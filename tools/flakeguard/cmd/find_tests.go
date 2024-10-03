package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// findtestsCmd represents the findtests command
var FindTestsCmd = &cobra.Command{
	Use:   "find-tests",
	Short: "Find tests based on changed Go files",
	Run: func(cmd *cobra.Command, args []string) {
		repoPath, _ := cmd.Flags().GetString("repo")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		// Find all changes in test files and get their package names

		testFiles, err := FindChangedTestFiles(repoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding changed test files: %s\n", err)
			os.Exit(1)
		}

		testPackages, err := GetFilePackages(repoPath, testFiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting package names for test files: %s\n", err)
			os.Exit(1)
		}

		// Find all changes in non-test files

		changedPackages, err := FindChangedNonTestPackages(repoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding changed non-test packages: %s\n", err)
			os.Exit(1)
		}

		dependentTestPackages, err := FindDependentPackages(repoPath, changedPackages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding dependent test packages: %s\n", err)
			os.Exit(1)
		}

		// Combine and deduplicate test package names
		allTestPackages := append(testPackages, dependentTestPackages...)
		allTestPackages = deduplicate(allTestPackages)

		if jsonOutput {
			data, err := json.MarshalIndent(allTestPackages, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshaling test files to JSON: %s\n", err)
				os.Exit(1)
			}
			fmt.Println(string(data))
		} else {
			fmt.Println("Changed test packages:")
			for _, file := range allTestPackages {
				fmt.Println(file)
			}
		}
	},
}

func init() {
	FindTestsCmd.Flags().StringP("repo", "r", ".", "Path to the Git repository")
	FindTestsCmd.Flags().Bool("json", false, "Output the results in JSON format")
}

func FindChangedTestFiles(repoPath string) ([]string, error) {
	cmd := exec.Command("bash", "-c", "git diff --name-only --diff-filter=AM develop...HEAD | grep '_test\\.go$'")
	cmd.Dir = repoPath
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error executing git diff command: %w", err)
	}
	testFiles := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(testFiles) == 1 && testFiles[0] == "" {
		return []string{}, nil
	}
	return testFiles, nil
}

func GetFilePackages(repoPath string, files []string) ([]string, error) {
	uniqueDirs := uniqueDirectories(files)
	return getPackageNames(uniqueDirs, repoPath), nil
}

func FindChangedNonTestPackages(repoPath string) ([]string, error) {
	cmd := exec.Command("bash", "-c", "git diff --name-only develop...HEAD | grep -v '_test\\.go$'")
	cmd.Dir = repoPath
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error executing git diff command: %w", err)
	}
	files := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}
	return GetFilePackages(repoPath, files)
}

// FindDependentPackages returns a list of packages that depend on any of the specified packages.
func FindDependentPackages(repoPath string, targetPackages []string) ([]string, error) {
	dependentPackages := make([]string, 0)

	// Execute 'go list' to find all Go test packages with their dependencies in the current module
	cmd := exec.Command("go", "list", "-f", `{{if .TestGoFiles}}{{.ImportPath}} {{join .Imports " "}} {{join .TestImports " "}}{{end}}`, "./...")
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running 'go list': %w", err)
	}

	// Scan each line to determine if it imports any of the target packages
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			continue // Skip if there aren't enough parts to include imports
		}
		packageName := parts[0]
		imports := strings.Fields(parts[1] + " " + parts[2])

		for _, target := range targetPackages {
			for _, imp := range imports {
				if imp == target {
					dependentPackages = append(dependentPackages, packageName)
					break
				}
			}
		}
	}

	return dependentPackages, nil
}

func uniqueDirectories(files []string) []string {
	dirSet := make(map[string]struct{})
	for _, file := range files {
		dirname := filepath.Dir(file)
		dirSet[dirname] = struct{}{}
	}
	var dirs []string
	for dir := range dirSet {
		dirs = append(dirs, dir)
	}
	return dirs
}

func deduplicate(items []string) []string {
	seen := make(map[string]struct{})
	var uniqueItems []string
	for _, item := range items {
		if _, found := seen[item]; !found {
			seen[item] = struct{}{}
			uniqueItems = append(uniqueItems, item)
		}
	}
	return uniqueItems
}

func getPackageNames(dirs []string, repoPath string) []string {
	var packageNames []string
	for _, dir := range dirs {
		cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}", "./"+dir)
		cmd.Dir = repoPath
		out, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error getting package name for directory %s: %s\n", dir, err)
			continue
		}
		packageName := strings.TrimSpace(string(out))
		if packageName != "" {
			packageNames = append(packageNames, packageName)
		}
	}
	return packageNames
}
