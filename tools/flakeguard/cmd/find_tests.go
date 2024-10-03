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

		if jsonOutput {
			data, err := json.MarshalIndent(testPackages, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshaling test files to JSON: %s\n", err)
				os.Exit(1)
			}
			fmt.Println(string(data))
		} else {
			fmt.Println("Changed test packages:")
			for _, file := range testPackages {
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

// TODO: currently this prints package names that import modified packages
// Refactor to return the unique list of packages to run based on list of modified packages
// TODO: create another function to print all test names to run based on the list of packages go test -v github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip -list .
func findPackageWithDependencies(repoPath string, packages []string) {

	// Find all Go test packages in the current module
	cmd := exec.Command("go", "list", "-f", `{{if .TestGoFiles}}{{.ImportPath}} {{join .Imports " "}} {{join .TestImports " "}}{{end}}`, "./...")
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running go list: %s\n", err)
		return
	}

	// Scan through each package and check if it imports any of the modified packages
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			continue // Skip entries without enough parts
		}
		packageName := parts[0]
		// filePath := parts[1] + parts[2]

		// cmd := exec.Command("go", "list", "-f", `{{join .Imports " "}} {{join .TestImports " "}}`, packageName)
		// cmd.Dir = repoPath
		// importsOutput, err := cmd.CombinedOutput()
		dependencies, err := getDependencies(packageName, repoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running go list for imports: %s\n", err)
			continue
		}

		for _, p := range packages {
			for _, imp := range dependencies {
				if imp == p {
					fmt.Printf("Package %s depends on package %s\n", packageName, p)
					break
				}
			}
		}
	}
}

func getDependencies(packageName, dir string) ([]string, error) {
	cmd := exec.Command("go", "list", "-f", `{{join .Deps " "}}`, packageName)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return strings.Fields(out.String()), nil
}

func findTestPackagesToRun(repoPath string) ([]string, error) {
	changedTestPackages, err := FindChangedTestFiles(repoPath)
	if err != nil {
		return nil, err
	}

	affectedTestPackages, err := findChangedNonTestPackages(repoPath)
	if err != nil {
		return nil, err
	}

	// Combine and deduplicate test package names
	allTestPackages := append(changedTestPackages, affectedTestPackages...)
	return deduplicate(allTestPackages), nil
}

// Determine affected test packages by analyzing package dependencies
func getAffectedTestPackages(changedPackages []string, allPackages map[string][]string) []string {
	testPackages := []string{}
	// This is a simple and non-optimal way to check dependencies
	for pkg, deps := range allPackages {
		for _, dep := range deps {
			for _, changed := range changedPackages {
				if dep == changed {
					testPackages = append(testPackages, pkg)
					break
				}
			}
		}
	}
	return testPackages
}

func findChangedNonTestPackages(repoPath string) ([]string, error) {
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

	uniqueDirs := uniqueDirectories(files)
	return getPackageNames(uniqueDirs, repoPath), nil
}

func getTestPackageNamesFromDirs(dirs []string, repoPath string) []string {
	var packageNames []string
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		cmd := exec.Command("go", "list", "-f", "{{if .TestGoFiles}}{{.ImportPath}}{{end}}", "-test", "./"+dir)
		cmd.Dir = repoPath
		out, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error getting test packages for directory %s: %s\n", dir, err)
			continue
		}
		packageName := strings.TrimSpace(string(out))
		if packageName != "" {
			packageNames = append(packageNames, packageName)
		}
	}
	return packageNames
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
