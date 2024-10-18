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
		testPackagesJson, _ := cmd.Flags().GetString("test-packages-json")
		testPackagesArg, _ := cmd.Flags().GetStringSlice("test-packages")
		runCount, _ := cmd.Flags().GetInt("run-count")
		useRace, _ := cmd.Flags().GetBool("race")
		failFast, _ := cmd.Flags().GetBool("fail-fast")
		outputPath, _ := cmd.Flags().GetString("output-json")
		threshold, _ := cmd.Flags().GetFloat64("threshold")

		var testPackages []string
		if testPackagesJson != "" {
			if err := json.Unmarshal([]byte(testPackagesJson), &testPackages); err != nil {
				log.Fatalf("Error decoding test packages JSON: %v", err)
			}
		} else if len(testPackagesArg) > 0 {
			testPackages = testPackagesArg
		} else {
			log.Fatalf("Error: must specify either --test-packages-json or --test-packages")
		}

		runner := runner.Runner{
			Verbose:  true,
			RunCount: runCount,
			UseRace:  useRace,
			FailFast: failFast,
		}

		testResults, err := runner.RunTests(testPackages)
		if err != nil {
			fmt.Printf("Error running tests: %v\n", err)
			os.Exit(1)
		}

		// Filter out failed tests based on the threshold
		failedTests := reports.FilterFailedTests(testResults, threshold)
		if len(failedTests) > 0 {
			jsonData, err := json.MarshalIndent(failedTests, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling test results to JSON: %v", err)
			}
			fmt.Printf("Threshold for flaky tests: %.2f\n%d failed tests:\n%s\n", threshold, len(failedTests), string(jsonData))
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
			fmt.Printf("All test results saved to %s\n", outputPath)
		}

		if len(failedTests) > 0 {
			os.Exit(1)
		} else {
			fmt.Printf("All %d tests passed.\n", len(testResults))
		}
	},
}

func init() {
	RunTestsCmd.Flags().String("test-packages-json", "", "JSON-encoded string of test packages")
	RunTestsCmd.Flags().StringSlice("test-packages", nil, "Comma-separated list of test packages to run")
	RunTestsCmd.Flags().IntP("run-count", "c", 1, "Number of times to run the tests")
	RunTestsCmd.Flags().Bool("race", false, "Enable the race detector")
	RunTestsCmd.Flags().Bool("fail-fast", false, "Stop on the first test failure")
	RunTestsCmd.Flags().String("output-json", "", "Path to output the test results in JSON format")
	RunTestsCmd.Flags().Float64("threshold", 0.8, "Threshold for considering a test as flaky")
}
