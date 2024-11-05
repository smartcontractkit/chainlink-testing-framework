package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var AggregateFailedTestsCmd = &cobra.Command{
	Use:   "aggregate-failed",
	Short: "Aggregate test results and output only failed tests based on a threshold",
	Run: func(cmd *cobra.Command, args []string) {
		folderPath, _ := cmd.Flags().GetString("folder-path")
		outputPath, _ := cmd.Flags().GetString("output-json")
		threshold, _ := cmd.Flags().GetFloat64("threshold")

		// Aggregate and merge results from the specified folder
		allResults, err := reports.AggregateTestResults(folderPath)
		if err != nil {
			log.Fatalf("Error aggregating results: %v", err)
		}

		// Filter failed tests based on threshold
		failedTests := reports.FilterFailedTests(allResults, threshold)

		// Format PassRatio as a percentage for display
		for i := range failedTests {
			failedTests[i].PassRatioPercentage = fmt.Sprintf("%.0f%%", failedTests[i].PassRatio*100)
		}

		// Output failed tests to JSON file
		if outputPath != "" && len(failedTests) > 0 {
			if err := saveResults(outputPath, failedTests); err != nil {
				log.Fatalf("Error writing failed tests to file: %v", err)
			}
			fmt.Printf("Aggregated failed test results saved to %s\n", outputPath)
		} else {
			fmt.Println("No failed tests found based on the specified threshold.")
		}
	},
}

func init() {
	AggregateFailedTestsCmd.Flags().String("folder-path", "testresult/", "Path to the folder containing JSON test result files")
	AggregateFailedTestsCmd.Flags().String("output-json", "failed_tests.json", "Path to output the aggregated failed test results in JSON format")
	AggregateFailedTestsCmd.Flags().Float64("threshold", 0.8, "Threshold for considering a test as failed")
}

// Helper function to save results to JSON file
func saveResults(filePath string, results []reports.TestResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling results: %v", err)
	}
	return os.WriteFile(filePath, data, 0644)
}
