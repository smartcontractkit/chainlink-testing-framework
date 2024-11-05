package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/git"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/golang"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/utils"
	"github.com/spf13/cobra"
)

var FindTestsCmd = &cobra.Command{
	Use:   "find",
	Long:  "Analyzes Golang project repository for changed files against a specified base reference and determines the test packages that are potentially impacted",
	Short: "Find test packages that may be affected by changes",
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, _ := cmd.Flags().GetString("project-path")
		verbose, _ := cmd.Flags().GetBool("verbose")
		jsonOutput, _ := cmd.Flags().GetBool("json")
		filterEmptyTests, _ := cmd.Flags().GetBool("filter-empty-tests")
		baseRef, _ := cmd.Flags().GetString("base-ref")
		excludes, _ := cmd.Flags().GetStringSlice("excludes")
		levels, _ := cmd.Flags().GetInt("levels")
		findByTestFilesDiff, _ := cmd.Flags().GetBool("find-by-test-files-diff")
		findByAffected, _ := cmd.Flags().GetBool("find-by-affected-packages")
		onlyShowChangedTestFiles, _ := cmd.Flags().GetBool("only-show-changed-test-files")

		// Find all changes in test files and get their package names
		var changedTestPkgs []string
		if findByTestFilesDiff {
			changedTestFiles, err := git.FindChangedFiles(projectPath, baseRef, "grep '_test\\.go$'")
			if err != nil {
				log.Fatalf("Error finding changed test files: %v", err)
			}
			if onlyShowChangedTestFiles {
				outputResults(changedTestFiles, jsonOutput)
				return
			}
			if verbose {
				fmt.Println("Changed test files:", changedTestFiles)
			}
			changedTestPkgs, err = golang.GetFilePackages(changedTestFiles)
			if err != nil {
				log.Fatalf("Error getting package names for test files: %v", err)
			}
		}

		// Find all affected test packages
		var affectedTestPkgs []string
		if findByAffected {
			if verbose {
				fmt.Println("Finding affected packages...")
			}
			affectedTestPkgs = findAffectedPackages(baseRef, projectPath, excludes, levels)
		}

		// Combine and deduplicate test package names
		testPkgs := append(changedTestPkgs, affectedTestPkgs...)
		testPkgs = utils.Deduplicate(testPkgs)

		// Filter out packages that do not have tests
		if filterEmptyTests {
			if verbose {
				fmt.Println("Filtering packages without tests...")
			}
			testPkgs = golang.FilterPackagesWithTests(testPkgs)
		}

		outputResults(testPkgs, jsonOutput)
	},
}

func init() {
	FindTestsCmd.Flags().StringP("project-path", "r", ".", "The path to the Go project. Default is the current directory. Useful for subprojects.")
	FindTestsCmd.Flags().String("base-ref", "", "Git base reference (branch, tag, commit) for comparing changes. Required.")
	FindTestsCmd.Flags().BoolP("verbose", "v", false, "Enable verbose mode")
	FindTestsCmd.Flags().Bool("json", false, "Output the results in JSON format")
	FindTestsCmd.Flags().Bool("filter-empty-tests", false, "Filter out test packages with no actual test functions. Can be very slow for large projects.")
	FindTestsCmd.Flags().StringSlice("excludes", []string{}, "List of paths to exclude. Useful for repositories with multiple Go projects within.")
	FindTestsCmd.Flags().IntP("levels", "l", 2, "The number of levels of recursion to search for affected packages. Default is 2. 0 is unlimited.")
	FindTestsCmd.Flags().Bool("find-by-test-files-diff", true, "Enable the mode to find test packages by changes in test files.")
	FindTestsCmd.Flags().Bool("find-by-affected-packages", true, "Enable the mode to find test packages that may be affected by changes in any of the project packages.")
	FindTestsCmd.Flags().Bool("only-show-changed-test-files", false, "Only show the changed test files and exit")

	if err := FindTestsCmd.MarkFlagRequired("base-ref"); err != nil {
		fmt.Println("Error marking base-ref as required:", err)
	}
}

func findAffectedPackages(baseRef, projectPath string, excludes []string, levels int) []string {
	goList, err := golang.GoList()
	if err != nil {
		log.Fatalf("Error getting go list: %v\nStdErr: %s", err, goList.Stderr.String())
	}
	gitDiff, err := git.Diff(baseRef)
	if err != nil {
		log.Fatalf("Error getting the git diff: %v\nStdErr: %s", err, gitDiff.Stderr.String())
	}
	gitModDiff, err := git.ModDiff(baseRef, projectPath)
	if err != nil {
		log.Fatalf("Error getting the git mod diff: %v\nStdErr: %s", err, gitModDiff.Stderr.String())
	}

	packages, err := golang.ParsePackages(goList.Stdout)
	if err != nil {
		log.Fatalf("Error parsing packages: %v", err)
	}

	fileMap := golang.GetGoFileMap(packages, true)

	var changedPackages []string
	changedPackages, err = git.GetChangedGoPackagesFromDiff(gitDiff.Stdout, projectPath, excludes, fileMap)
	if err != nil {
		log.Fatalf("Error getting changed packages: %v", err)
	}

	changedModPackages, err := git.GetGoModChangesFromDiff(gitModDiff.Stdout)
	if err != nil {
		log.Fatalf("Error getting go.mod changes: %v", err)
	}

	depMap := golang.GetGoDepMap(packages)

	// Find affected packages
	// use map to make handling duplicates simpler
	affectedPkgs := map[string]bool{}

	// loop through packages changed via file changes
	for _, pkg := range changedPackages {
		p := golang.FindAffectedPackages(pkg, depMap, false, levels)
		for _, p := range p {
			affectedPkgs[p] = true
		}
	}

	// loop through packages changed via go.mod changes
	for _, pkg := range changedModPackages {
		p := golang.FindAffectedPackages(pkg, depMap, true, levels)
		for _, p := range p {
			affectedPkgs[p] = true
		}
	}

	// convert map to array
	pkgs := []string{}
	for k := range affectedPkgs {
		pkgs = append(pkgs, k)
	}

	return pkgs
}

func outputResults(packages []string, jsonOutput bool) {
	if jsonOutput {
		data, err := json.Marshal(packages)
		if err != nil {
			log.Fatalf("Error marshaling test files to JSON: %v", err)
		}
		fmt.Println(string(data))
	} else {
		for _, pkg := range packages {
			fmt.Print(pkg, " ")
		}
	}
}
