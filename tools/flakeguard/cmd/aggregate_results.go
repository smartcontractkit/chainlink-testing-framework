package cmd

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var AggregateResultsCmd = &cobra.Command{
	Use:   "aggregate-results",
	Short: "Aggregate test results into a single report, with optional filtering and code owners mapping",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		aggregateResultsPath, _ := cmd.Flags().GetString("results-path")
		aggregateOutputPath, _ := cmd.Flags().GetString("output-path")
		includeOutputs, _ := cmd.Flags().GetBool("include-outputs")
		includePackageOutputs, _ := cmd.Flags().GetBool("include-package-outputs")
		filterFailed, _ := cmd.Flags().GetBool("filter-failed")
		maxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio")
		codeOwnersPath, _ := cmd.Flags().GetString("codeowners-path")
		repoPath, _ := cmd.Flags().GetString("repo-path")

		// Load test reports from JSON files
		testReports, err := reports.LoadReports(aggregateResultsPath)
		if err != nil {
			return fmt.Errorf("error loading test reports: %w", err)
		}

		// Aggregate the reports
		aggregatedReport, err := reports.Aggregate(testReports...)
		if err != nil {
			return fmt.Errorf("error aggregating test reports: %w", err)
		}

		// Map test results to test paths
		err = reports.MapTestResultsToPaths(aggregatedReport, repoPath)
		if err != nil {
			return fmt.Errorf("error mapping test results to paths: %w", err)
		}

		// Map test results to code owners if codeOwnersPath is provided
		if codeOwnersPath != "" {
			err = reports.MapTestResultsToOwners(aggregatedReport, codeOwnersPath)
			if err != nil {
				return fmt.Errorf("error mapping test results to code owners: %w", err)
			}
		}

		// Filter results if needed
		if filterFailed {
			aggregatedReport.Results = reports.FilterTests(aggregatedReport.Results, func(tr reports.TestResult) bool {
				return !tr.Skipped && tr.PassRatio < maxPassRatio
			})
		}

		// Process the aggregated results based on the flags
		if !includeOutputs || !includePackageOutputs {
			for i := range aggregatedReport.Results {
				if !includeOutputs {
					aggregatedReport.Results[i].Outputs = nil
				}
				if !includePackageOutputs {
					aggregatedReport.Results[i].PackageOutputs = nil
				}
			}
		}

		// Save the aggregated report
		if err := reports.SaveReport(reports.OSFileSystem{}, aggregateOutputPath, *aggregatedReport); err != nil {
			return fmt.Errorf("error saving aggregated report: %w", err)
		}

		fmt.Printf("Aggregated report saved to %s\n", aggregateOutputPath)
		return nil
	},
}

func init() {
	AggregateResultsCmd.Flags().StringP("results-path", "p", "", "Path to the folder containing JSON test result files (required)")
	AggregateResultsCmd.Flags().StringP("output-path", "o", "./aggregated-results.json", "Path to output the aggregated test results")
	AggregateResultsCmd.Flags().Bool("include-outputs", false, "Include test outputs in the aggregated test results")
	AggregateResultsCmd.Flags().Bool("include-package-outputs", false, "Include test package outputs in the aggregated test results")
	AggregateResultsCmd.Flags().Bool("filter-failed", false, "If true, filter and output only failed tests based on the max-pass-ratio threshold")
	AggregateResultsCmd.Flags().Float64("max-pass-ratio", 1.0, "The maximum pass ratio threshold for a test to be considered flaky. Any tests below this pass rate will be considered flaky.")
	AggregateResultsCmd.Flags().String("codeowners-path", "", "Path to the CODEOWNERS file")
	AggregateResultsCmd.Flags().String("repo-path", ".", "The path to the root of the repository/project")
	AggregateResultsCmd.MarkFlagRequired("results-path")
	AggregateResultsCmd.MarkFlagRequired("repo-path")
}
