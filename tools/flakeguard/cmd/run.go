package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
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
		sendToLoki, _ := cmd.Flags().GetBool("send-to-loki")

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
		if outputPath != "" && len(testResults) > 0 {
			jsonData, err := json.MarshalIndent(testResults, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling test results to JSON: %v", err)
			}
			if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
				log.Fatalf("Error writing test results to file: %v", err)
			}
			fmt.Printf("All test results saved to %s\n", outputPath)
		}

		// Send test results to Loki
		if sendToLoki {
			lc, err := newLokiClient()
			if err != nil {
				log.Fatalf("Error creating Loki client: %v", err)
			}
			reports.SendResultsToLoki(lc, testResults)
			lc.StopNow()
		}

		if len(failedTests) > 0 {
			os.Exit(1)
		} else if len(testResults) == 0 {
			fmt.Printf("No tests were run for the specified packages.\n")
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
	RunTestsCmd.Flags().Bool("send-to-loki", false, "Send test results to Loki")
	RunTestsCmd.Flags().String("loki-endpoint", "", "Loki endpoint")
	RunTestsCmd.Flags().String("loki-tenant-id", "", "Loki tenant ID")
	RunTestsCmd.Flags().String("loki-basic-auth-login", "", "Loki basic auth login")
}

func newLokiClient() (*client.LokiPromtailClient, error) {
	endpoint := os.Getenv("LOKI_ENDPOINT")
	tenant := os.Getenv("LOKI_TENANT_ID")
	basicAuth := os.Getenv("LOKI_BASIC_AUTH")
	token := os.Getenv("LOKI_TOKEN")

	config := client.NewLokiConfig(&endpoint, &tenant, &basicAuth, &token)

	return client.NewLokiPromtailClient(config)
}
