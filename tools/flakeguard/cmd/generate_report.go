package cmd

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

var (
	reportInputPath      string
	reportOutputPath     string
	reportFormat         string
	reportCodeOwnersPath string
	reportProjectPath    string
)

var GenerateReportCmd = &cobra.Command{
	Use:   "generate-report",
	Short: "Generate reports from test results",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load the test results
		testReport, err := reports.LoadReport(reportInputPath)
		if err != nil {
			return fmt.Errorf("error loading test report: %w", err)
		}

		// Generate the report
		if err := generateReport(testReport, reportFormat, reportOutputPath); err != nil {
			return fmt.Errorf("error generating report: %w", err)
		}

		fmt.Printf("Report generated at %s\n", reportOutputPath)
		return nil
	},
}

func init() {
	GenerateReportCmd.Flags().StringVarP(&reportInputPath, "aggregated-report-path", "i", "", "Path to the aggregated test results file (required)")
	GenerateReportCmd.Flags().StringVarP(&reportOutputPath, "output-path", "o", "./report", "Path to output the generated report (without extension)")
	GenerateReportCmd.Flags().StringVarP(&reportFormat, "format", "f", "markdown", "Format of the report (markdown, json)")
	GenerateReportCmd.MarkFlagRequired("aggregated-report-path")
}

func generateReport(report *reports.TestReport, format, outputPath string) error {
	fs := reports.OSFileSystem{}
	switch strings.ToLower(format) {
	case "markdown":
		mdFileName := outputPath + ".md"
		mdFile, err := fs.Create(mdFileName)
		if err != nil {
			return fmt.Errorf("error creating markdown file: %w", err)
		}
		defer mdFile.Close()
		reports.GenerateMarkdownSummary(mdFile, report, 1.0)
		fmt.Printf("Markdown report saved to %s\n", mdFileName)
	case "json":
		jsonFileName := outputPath + ".json"
		if err := reports.SaveReportNoLogs(fs, jsonFileName, *report); err != nil {
			return fmt.Errorf("error saving JSON report: %w", err)
		}
		fmt.Printf("JSON report saved to %s\n", jsonFileName)
	default:
		return fmt.Errorf("unsupported report format: %s", format)
	}

	// Generate summary JSON
	summaryData := reports.GenerateSummaryData(report.Results, 1.0)
	summaryFileName := outputPath + "-summary.json"
	if err := reports.SaveSummaryAsJSON(fs, summaryFileName, summaryData); err != nil {
		return fmt.Errorf("error saving summary JSON: %w", err)
	}
	fmt.Printf("Summary JSON saved to %s\n", summaryFileName)

	return nil
}
