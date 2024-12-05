package cmd

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"

	"github.com/spf13/cobra"
)

var (
	filterInputPath    string
	filterOutputPath   string
	filterMaxPassRatio float64
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter aggregated test results based on criteria",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load the aggregated report
		aggregatedReport, err := reports.LoadReport(filterInputPath)
		if err != nil {
			return fmt.Errorf("error loading aggregated report: %w", err)
		}

		// Filter the test results
		filteredReport := reports.FilterResults(aggregatedReport, filterMaxPassRatio)

		// Save the filtered report
		if err := reports.SaveReport(reports.OSFileSystem{}, filterOutputPath, *filteredReport); err != nil {
			return fmt.Errorf("error saving filtered report: %w", err)
		}

		fmt.Printf("Filtered report saved to %s\n", filterOutputPath)
		return nil
	},
}

func init() {
	filterCmd.Flags().StringVarP(&filterInputPath, "input-path", "i", "", "Path to the aggregated test results file (required)")
	filterCmd.Flags().StringVarP(&filterOutputPath, "output-path", "o", "./filtered-results.json", "Path to output the filtered test results")
	filterCmd.Flags().Float64VarP(&filterMaxPassRatio, "max-pass-ratio", "m", 1.0, "Maximum pass ratio threshold for filtering tests")
	filterCmd.MarkFlagRequired("input-path")
}
