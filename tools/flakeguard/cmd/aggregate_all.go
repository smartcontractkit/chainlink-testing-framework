package cmd

import (
	"log"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var AggregateAllCmd = &cobra.Command{
	Use:   "aggregate-all",
	Short: "Aggregate all test results and output them to a file",
	Run: func(cmd *cobra.Command, args []string) {
		resultsFolderPath, _ := cmd.Flags().GetString("results-path")
		outputResultsPath, _ := cmd.Flags().GetString("output-results")
		outputLogsPath, _ := cmd.Flags().GetString("output-logs")

		// Aggregate all test results
		allResults, err := reports.AggregateTestResults(resultsFolderPath)
		if err != nil {
			log.Fatalf("Error aggregating results: %v", err)
		}

		// Output all results to JSON files
		if len(allResults) > 0 {
			reports.SaveFilteredResultsAndLogs(outputResultsPath, outputLogsPath, allResults)
		}
	},
}

func init() {
	AggregateAllCmd.Flags().String("results-path", "", "Path to the folder containing JSON test result files")
	AggregateAllCmd.Flags().String("output-results", "./failed_tests.json", "Path to output the filtered failed test results in JSON format")
	AggregateAllCmd.Flags().String("output-logs", "", "Path to output the filtered failed test logs in JSON format")
}
