package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/utils"
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

		changedTestFiles, err := utils.FindChangedFiles(repoPath, "grep '_test\\.go$'")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding changed test files: %s\n", err)
			os.Exit(1)
		}

		testPackages, err := utils.GetFilePackages(repoPath, changedTestFiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting package names for test files: %s\n", err)
			os.Exit(1)
		}

		// Find all changes in non-test files

		changedFiles, err := utils.FindChangedFiles(repoPath, "grep -v '_test\\.go$'")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding changed non-test packages: %s\n", err)
			os.Exit(1)
		}

		changedPackages, err := utils.GetFilePackages(repoPath, changedFiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting package names for non-test files: %s\n", err)
			os.Exit(1)
		}

		dependentTestPackages, err := utils.FindDependentPackages(repoPath, changedPackages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding dependent test packages: %s\n", err)
			os.Exit(1)
		}

		// Combine and deduplicate test package names
		allTestPackages := append(testPackages, dependentTestPackages...)
		allTestPackages = utils.Deduplicate(allTestPackages)

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
