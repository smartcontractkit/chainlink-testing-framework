package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner"
	"github.com/spf13/cobra"
)

var RunTestsCmd = &cobra.Command{
	Use:   "run-tests",
	Short: "Run tests to find flaky ones",
	Run: func(cmd *cobra.Command, args []string) {
		repoPath, _ := cmd.Flags().GetString("repo-path")
		testPackagesJson, _ := cmd.Flags().GetString("test-packages-json")
		testPackage, _ := cmd.Flags().GetString("test-package")
		count, _ := cmd.Flags().GetInt("count")
		useRace, _ := cmd.Flags().GetBool("race")
		failFast, _ := cmd.Flags().GetBool("fail-fast")
		outputPath, _ := cmd.Flags().GetString("output-json")

		var testPackages []string
		if testPackagesJson != "" {
			if err := json.Unmarshal([]byte(testPackagesJson), &testPackages); err != nil {
				fmt.Fprintf(os.Stderr, "Error decoding test packages JSON: %s\n", err)
				os.Exit(1)
			}
		} else if testPackage != "" {
			testPackages = append(testPackages, testPackage)
		} else {
			fmt.Fprintf(os.Stderr, "Error: must specify either --test-packages-json or --test-package\n")
			os.Exit(1)
		}

		runner := runner.Runner{
			Verbose:  true,
			Dir:      repoPath,
			Count:    count,
			UseRace:  useRace,
			FailFast: failFast,
		}

		testResults, err := runner.RunTests(testPackages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running tests: %s\n", err)
			os.Exit(1)
		}

		// Save the test results in JSON format
		if outputPath != "" {
			jsonData, err := json.MarshalIndent(testResults, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshaling test results to JSON: %s\n", err)
				os.Exit(1)
			}
			if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing test results to file: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("Test results saved to %s\n", outputPath)
		}
	},
}

func init() {
	RunTestsCmd.Flags().StringP("repo-path", "r", ".", "Path to the Git repository")
	RunTestsCmd.Flags().String("test-packages-json", "", "JSON-encoded string of test packages")
	RunTestsCmd.Flags().String("test-package", "", "Single test package to run")
	RunTestsCmd.Flags().IntP("count", "c", 1, "Number of times to run the tests")
	RunTestsCmd.Flags().Bool("race", false, "Enable the race detector")
	RunTestsCmd.Flags().Bool("fail-fast", false, "Stop on the first test failure")
	RunTestsCmd.Flags().String("output-json", "", "Path to output the test results in JSON format")
}
