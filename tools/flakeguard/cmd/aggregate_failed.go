package cmd

import (
	"log"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var AggregateFailedCmd = &cobra.Command{
	Use:   "aggregate-failed",
	Short: "Aggregate all test results, then filter and output only failed tests based on a threshold",
	Run: func(cmd *cobra.Command, args []string) {
		resultsFolderPath, _ := cmd.Flags().GetString("results-path")
		outputResultsPath, _ := cmd.Flags().GetString("output-results")
		outputLogsPath, _ := cmd.Flags().GetString("output-logs")
		threshold, _ := cmd.Flags().GetFloat64("threshold")
		minPassRatio, _ := cmd.Flags().GetFloat64("min-pass-ratio")

		// Aggregate all test results
		allResults, err := reports.AggregateTestResults(resultsFolderPath)
		if err != nil {
			log.Fatalf("Error aggregating results: %v", err)
		}

		// Filter to only include failed tests based on threshold and minPassRatio
		var failedResults []reports.TestResult
		for _, result := range allResults {
			if result.PassRatio < threshold && result.PassRatio > minPassRatio && !result.Skipped {
				failedResults = append(failedResults, result)
			}
		}

		// Output results to JSON files
		if len(failedResults) > 0 {
			reports.SaveFilteredResultsAndLogs(outputResultsPath, outputLogsPath, failedResults)
		}
	},
}

func init() {
	AggregateFailedCmd.Flags().String("results-path", "testresult/", "Path to the folder containing JSON test result files")
	AggregateFailedCmd.Flags().String("output-results", "failed_tests.json", "Path to output the filtered failed test results in JSON format")
	AggregateFailedCmd.Flags().String("output-logs", "failed_logs.json", "Path to output the filtered failed test logs in JSON format")
	AggregateFailedCmd.Flags().Float64("threshold", 0.8, "Threshold for considering a test as failed")
	AggregateFailedCmd.Flags().Float64("min-pass-ratio", 0.001, "Minimum pass ratio for considering a test as flaky. Used to distinguish between tests that are truly flaky (with inconsistent results) and those that are consistently failing.")
}
