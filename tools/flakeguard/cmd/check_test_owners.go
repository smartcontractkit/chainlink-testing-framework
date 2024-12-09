package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/codeowners"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var (
	codeownersPath     string
	printTestFunctions bool
)

// CheckTestOwnersCmd checks which tests lack code owners
var CheckTestOwnersCmd = &cobra.Command{
	Use:   "check-test-owners",
	Short: "Check which tests in the project do not have code owners",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath, _ := cmd.Flags().GetString("project-path")

		// Scan project for test functions
		testFileMap, err := reports.ScanTestFiles(projectPath)
		if err != nil {
			log.Fatalf("Error scanning test files: %v", err)
		}

		// Parse CODEOWNERS file
		codeOwnerPatterns, err := codeowners.Parse(codeownersPath)
		if err != nil {
			log.Fatalf("Error parsing CODEOWNERS file: %v", err)
		}

		// Check for tests without code owners
		testsWithoutOwners := make(map[string]string) // Map of test names to their relative paths
		for testName, filePath := range testFileMap {
			relFilePath, err := filepath.Rel(projectPath, filePath)
			if err != nil {
				fmt.Printf("Error getting relative path for test %s: %v\n", testName, err)
				continue
			}
			// Convert to Unix-style path for matching
			relFilePath = filepath.ToSlash(relFilePath)

			owners := codeowners.FindOwners(relFilePath, codeOwnerPatterns)
			if len(owners) == 0 {
				testsWithoutOwners[testName] = relFilePath
			}
		}

		// Calculate percentages
		totalTests := len(testFileMap)
		totalWithoutOwners := len(testsWithoutOwners)
		totalWithOwners := totalTests - totalWithoutOwners

		percentageWithOwners := float64(totalWithOwners) / float64(totalTests) * 100
		percentageWithoutOwners := float64(totalWithoutOwners) / float64(totalTests) * 100

		// Report results
		fmt.Printf("Total Test functions found: %d\n", totalTests)
		fmt.Printf("Test functions with owners: %d (%.2f%%)\n", totalWithOwners, percentageWithOwners)
		fmt.Printf("Test functions without owners: %d (%.2f%%)\n", totalWithoutOwners, percentageWithoutOwners)

		if printTestFunctions {

			if totalWithoutOwners > 0 {
				fmt.Println("\nTest functions without owners:")
				for testName, relPath := range testsWithoutOwners {
					fmt.Printf("- %s (%s)\n", testName, relPath)
				}
			} else {
				fmt.Println("All Test functions have code owners!")
			}
		}

		return nil
	},
}

func init() {
	CheckTestOwnersCmd.Flags().StringP("project-path", "p", ".", "Path to the root of the project")
	CheckTestOwnersCmd.Flags().StringVarP(&codeownersPath, "codeowners-path", "c", ".github/CODEOWNERS", "Path to the CODEOWNERS file")
	CheckTestOwnersCmd.Flags().BoolVarP(&printTestFunctions, "print-test-functions", "t", false, "Print all test functions without owners")
}
