package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner"
	"github.com/spf13/cobra"
)

var RunTestsCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests to check if they are flaky",
	Run: func(cmd *cobra.Command, args []string) {
		repoPath, _ := cmd.Flags().GetString("repo-path")
		testPackagesJson, _ := cmd.Flags().GetString("test-packages-json")
		testPackage, _ := cmd.Flags().GetString("test-package")
		count, _ := cmd.Flags().GetInt("count")
		useRace, _ := cmd.Flags().GetBool("race")
		failFast, _ := cmd.Flags().GetBool("fail-fast")
		outputPath, _ := cmd.Flags().GetString("output-json")
		threshold, _ := cmd.Flags().GetFloat64("threshold")

		var testPackages []string
		if testPackagesJson != "" {
			if err := json.Unmarshal([]byte(testPackagesJson), &testPackages); err != nil {
				log.Fatalf("Error decoding test packages JSON: %v", err)
			}
		} else if testPackage != "" {
			testPackages = append(testPackages, testPackage)
		} else {
			log.Fatalf("Error: must specify either --test-packages-json or --test-package")
		}

		runner := runner.Runner{
			Verbose:  true,
			Dir:      repoPath,
			Count:    count,
			UseRace:  useRace,
			FailFast: failFast,
		}

		testResults, _ := runner.RunTests(testPackages)
		// TODO: Handle error

		// Filter out failed tests based on the threshold
		failedTests := reports.FilterFailedTests(testResults, threshold)
		if len(failedTests) > 0 {
			jsonData, err := json.MarshalIndent(failedTests, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling test results to JSON: %v", err)
			}
			fmt.Printf("Threshold for flaky tests: %.2f\nFailed tests:\n%s\n", threshold, string(jsonData))
			fmt.Println(string(jsonData))
		}

		// Save the test results in JSON format
		if outputPath != "" {
			jsonData, err := json.MarshalIndent(testResults, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling test results to JSON: %v", err)
			}
			if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
				log.Fatalf("Error writing test results to file: %v", err)
			}
			fmt.Printf("Test results saved to %s\n", outputPath)
		}

		if len(failedTests) > 0 {
			os.Exit(1)
		} else {
			fmt.Println("All tests passed.")
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
	RunTestsCmd.Flags().Float64("threshold", 0.8, "Threshold for considering a test as flaky")
}
