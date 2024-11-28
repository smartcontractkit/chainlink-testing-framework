package cmd

import (
	"log"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var (
	resultsFolderPath string
	outputResultsPath string
	outputLogsPath    string
	maxPassRatio      float64
	filterFailed      bool
)

var AggregateResultsCmd = &cobra.Command{
	Use:   "aggregate-results",
	Short: "Aggregate test results and optionally filter failed tests based on a threshold",
	RunE: func(cmd *cobra.Command, args []string) error {
		allReport, err := reports.AggregateTestResults(resultsFolderPath)
		if err != nil {
			log.Fatalf("Error aggregating results: %v", err)
		}

		var resultsToSave []reports.TestResult

		if filterFailed {
			// Filter to only include tests that failed below the threshold
			for _, result := range allReport.Results {
				if result.PassRatio < maxPassRatio && !result.Skipped {
					resultsToSave = append(resultsToSave, result)
				}
			}
		} else {
			resultsToSave = allReport.Results
		}
		allReport.Results = resultsToSave

		// Output results to JSON files
		if len(resultsToSave) > 0 {
			return reports.SaveFilteredResultsAndLogs(outputResultsPath, outputLogsPath, allReport)
		}
		return nil
	},
}

func init() {
	AggregateResultsCmd.Flags().StringVarP(&resultsFolderPath, "results-path", "p", "", "Path to the folder containing JSON test result files")
	AggregateResultsCmd.Flags().StringVarP(&outputResultsPath, "output-results", "o", "./results", "Path to output the aggregated or filtered test results in JSON and markdown format")
	AggregateResultsCmd.Flags().StringVarP(&outputLogsPath, "output-logs", "l", "", "Path to output the filtered test logs in JSON format")
	AggregateResultsCmd.Flags().Float64VarP(&maxPassRatio, "max-pass-ratio", "m", 1.0, "The maximum (non-inclusive) pass ratio threshold for a test to be considered a failure. Any tests below this pass rate will be considered flaky.")
	AggregateResultsCmd.Flags().BoolVarP(&filterFailed, "filter-failed", "f", false, "If true, filter and output only failed tests based on the max-pass-ratio threshold")
}
