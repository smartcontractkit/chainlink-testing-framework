package cmd

import (
	"log"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var AggregateResultsCmd = &cobra.Command{
	Use:   "aggregate-results",
	Short: "Aggregate test results and optionally filter failed tests based on a threshold",
	Run: func(cmd *cobra.Command, args []string) {
		resultsFolderPath, _ := cmd.Flags().GetString("results-path")
		outputResultsPath, _ := cmd.Flags().GetString("output-results")
		outputLogsPath, _ := cmd.Flags().GetString("output-logs")
		threshold, _ := cmd.Flags().GetFloat64("threshold")
		minPassRatio, _ := cmd.Flags().GetFloat64("min-pass-ratio")
		filterFailed, _ := cmd.Flags().GetBool("filter-failed")

		// Aggregate all test results
		allResults, err := reports.AggregateTestResults(resultsFolderPath)
		if err != nil {
			log.Fatalf("Error aggregating results: %v", err)
		}

		var resultsToSave []reports.TestResult

		if filterFailed {
			// Filter to only include failed tests based on threshold and minPassRatio
			for _, result := range allResults {
				if result.PassRatio < threshold && result.PassRatio > minPassRatio && !result.Skipped {
					resultsToSave = append(resultsToSave, result)
				}
			}
		} else {
			resultsToSave = allResults
		}

		// Output results to JSON files
		if len(resultsToSave) > 0 {
			reports.SaveFilteredResultsAndLogs(outputResultsPath, outputLogsPath, resultsToSave)
		}
	},
}

func init() {
	AggregateResultsCmd.Flags().String("results-path", "", "Path to the folder containing JSON test result files")
	AggregateResultsCmd.Flags().String("output-results", "./results.json", "Path to output the aggregated or filtered test results in JSON format")
	AggregateResultsCmd.Flags().String("output-logs", "", "Path to output the filtered test logs in JSON format")
	AggregateResultsCmd.Flags().Float64("threshold", 0.8, "Threshold for considering a test as failed (used with --filter-failed)")
	AggregateResultsCmd.Flags().Float64("min-pass-ratio", 0.001, "Minimum pass ratio for considering a test as flaky (used with --filter-failed)")
	AggregateResultsCmd.Flags().Bool("filter-failed", false, "If true, filter and output only failed tests based on the threshold")
}
