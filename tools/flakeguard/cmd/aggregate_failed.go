package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var AggregateFailedCmd = &cobra.Command{
	Use:   "aggregate-failed",
	Short: "Aggregate all test results, then filter and output only failed tests based on a threshold",
	Run: func(cmd *cobra.Command, args []string) {
		resultsFolderPath, _ := cmd.Flags().GetString("results-path")
		outputPath, _ := cmd.Flags().GetString("output-json")
		threshold, _ := cmd.Flags().GetFloat64("threshold")

		// Aggregate all test results
		allResults, err := reports.AggregateTestResults(resultsFolderPath)
		if err != nil {
			log.Fatalf("Error aggregating results: %v", err)
		}

		// Filter to only include failed tests based on threshold
		var failedResults []reports.TestResult
		for _, result := range allResults {
			if result.PassRatio < threshold && !result.Skipped {
				failedResults = append(failedResults, result)
			}
		}

		// Output failed results to JSON file
		if outputPath != "" && len(failedResults) > 0 {
			if err := saveResults(outputPath, failedResults); err != nil {
				log.Fatalf("Error writing failed results to file: %v", err)
			}
			fmt.Printf("Filtered failed test results saved to %s\n", outputPath)
		} else {
			fmt.Println("No failed tests found based on the specified threshold.")
		}
	},
}

func init() {
	AggregateFailedCmd.Flags().String("results-path", "testresult/", "Path to the folder containing JSON test result files")
	AggregateFailedCmd.Flags().String("output-json", "failed_tests.json", "Path to output the filtered failed test results in JSON format")
	AggregateFailedCmd.Flags().Float64("threshold", 0.8, "Threshold for considering a test as failed")
}

// Helper function to save results to JSON file
func saveResults(filePath string, results []reports.TestResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling results: %v", err)
	}
	return os.WriteFile(filePath, data, 0644)
}
