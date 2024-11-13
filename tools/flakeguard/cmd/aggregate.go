package cmd

import (
	"fmt"
	"log"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var AggregateAllCmd = &cobra.Command{
	Use:   "aggregate-all",
	Short: "Aggregate all test results and output them to a file",
	Run: func(cmd *cobra.Command, args []string) {
		resultsFolderPath, _ := cmd.Flags().GetString("results-path")
		outputPath, _ := cmd.Flags().GetString("output-json")

		// Aggregate all test results
		allResults, err := reports.AggregateTestResults(resultsFolderPath)
		if err != nil {
			log.Fatalf("Error aggregating results: %v", err)
		}

		// Output all results to JSON file
		if outputPath != "" && len(allResults) > 0 {
			if err := saveResults(outputPath, allResults); err != nil {
				log.Fatalf("Error writing aggregated results to file: %v", err)
			}
			fmt.Printf("Aggregated test results saved to %s\n", outputPath)
		} else {
			fmt.Println("No test results found.")
		}
	},
}

func init() {
	AggregateAllCmd.Flags().String("results-path", "testresult/", "Path to the folder containing JSON test result files")
	AggregateAllCmd.Flags().String("output-json", "all_tests.json", "Path to output the aggregated test results in JSON format")
}
