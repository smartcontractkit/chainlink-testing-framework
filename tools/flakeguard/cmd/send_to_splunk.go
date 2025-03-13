package cmd

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/spf13/cobra"
)

// SendToSplunkCmd sends a TestReport to a Splunk.
var SendToSplunkCmd = &cobra.Command{
	Use:   "send-to-splunk",
	Short: "Send TestReport data to Splunk",
	Run: func(cmd *cobra.Command, args []string) {
		reportPath, _ := cmd.Flags().GetString("report-path")
		splunkURL, _ := cmd.Flags().GetString("splunk-url")
		splunkToken, _ := cmd.Flags().GetString("splunk-token")
		splunkEvent, _ := cmd.Flags().GetString("splunk-event")
		failLogsURL, _ := cmd.Flags().GetString("failed-logs-url")
		repoURL, _ := cmd.Flags().GetString("repo-url")
		branchName, _ := cmd.Flags().GetString("branch-name")
		headSHA, _ := cmd.Flags().GetString("head-sha")
		baseSHA, _ := cmd.Flags().GetString("base-sha")
		githubWorkflowName, _ := cmd.Flags().GetString("github-workflow-name")
		githubWorkflowRunURL, _ := cmd.Flags().GetString("github-workflow-run-url")
		reportID, _ := cmd.Flags().GetString("report-id")
		genReportID, _ := cmd.Flags().GetBool("gen-report-id")
		repoPath, _ := cmd.Flags().GetString("repo-path")
		codeownersPath, _ := cmd.Flags().GetString("codeowners-path")

		// Read the report file.
		data, err := os.ReadFile(reportPath)
		if err != nil {
			log.Error().Err(err).Str("path", reportPath).Msg("Error reading report file")
			os.Exit(1)
		}

		// Unmarshal JSON data into the TestReport struct.
		var testReport reports.TestReport
		if err := json.Unmarshal(data, &testReport); err != nil {
			log.Error().Err(err).Msg("Error unmarshalling report JSON")
			os.Exit(1)
		}
		testReport.GenerateSummaryData()

		// Override report fields with flags if provided.
		if repoURL != "" {
			testReport.RepoURL = repoURL
		}
		if branchName != "" {
			testReport.BranchName = branchName
		}
		if headSHA != "" {
			testReport.HeadSHA = headSHA
		}
		if baseSHA != "" {
			testReport.BaseSHA = baseSHA
		}
		if githubWorkflowName != "" {
			testReport.GitHubWorkflowName = githubWorkflowName
		}
		if githubWorkflowRunURL != "" {
			testReport.GitHubWorkflowRunURL = githubWorkflowRunURL
		}
		if reportID != "" {
			testReport.SetReportID(reportID)
		}
		if genReportID {
			testReport.SetRandomReportID()
		}
		if repoPath != "" {
			err = reports.MapTestResultsToPaths(&testReport, repoPath)
			if err != nil {
				log.Error().Err(err).Msg("Error mapping test results to paths")
				os.Exit(1)
			}
		}
		if codeownersPath != "" && repoPath != "" {
			reports.MapTestResultsToOwners(&testReport, codeownersPath)
		}
		if failLogsURL != "" {
			testReport.FailedLogsURL = failLogsURL
		}

		// Send the test report to Splunk.
		err = reports.SendTestReportToSplunk(splunkURL, splunkToken, splunkEvent, testReport)
		if err != nil {
			log.Error().Err(err).Msg("Error sending test report to Splunk")
			os.Exit(1)
		}

		log.Info().Msg("Successfully sent test report to Splunk")
	},
}

func init() {
	// Define flags for the new command.
	SendToSplunkCmd.Flags().String("report-path", "", "Path to the test report JSON file (required)")
	SendToSplunkCmd.Flags().String("failed-logs-url", "", "Optional URL linking to additional logs for failed tests")
	SendToSplunkCmd.Flags().String("splunk-url", "", "Optional URL to send the test results to Splunk")
	SendToSplunkCmd.Flags().String("splunk-token", "", "Optional Splunk HEC token to send the test results")
	SendToSplunkCmd.Flags().String("splunk-event", "", "Optional Splunk event to send as the triggering event for the test results")
	SendToSplunkCmd.Flags().String("repo-url", "", "The repository URL")
	SendToSplunkCmd.Flags().String("branch-name", "", "Branch name for the test report")
	SendToSplunkCmd.Flags().String("head-sha", "", "Head commit SHA for the test report")
	SendToSplunkCmd.Flags().String("base-sha", "", "Base commit SHA for the test report")
	SendToSplunkCmd.Flags().String("github-workflow-name", "", "GitHub workflow name for the test report")
	SendToSplunkCmd.Flags().String("github-workflow-run-url", "", "GitHub workflow run URL for the test report")
	SendToSplunkCmd.Flags().String("report-id", "", "Optional identifier for the test report. Will be generated if not provided")
	SendToSplunkCmd.Flags().StringP("repo-path", "", ".", "The path to the root of the repository/project")
	SendToSplunkCmd.Flags().StringP("codeowners-path", "", "", "Path to the CODEOWNERS file")
	SendToSplunkCmd.Flags().Bool("gen-report-id", false, "Generate a random report ID")

	// Mark required flags.
	if err := SendToSplunkCmd.MarkFlagRequired("report-path"); err != nil {
		log.Fatal().Err(err).Msg("Error marking report-path flag as required")
	}
	if err := SendToSplunkCmd.MarkFlagRequired("splunk-url"); err != nil {
		log.Fatal().Err(err).Msg("Error marking splunk-url flag as required")
	}
	if err := SendToSplunkCmd.MarkFlagRequired("splunk-token"); err != nil {
		log.Fatal().Err(err).Msg("Error marking splunk-token flag as required")
	}
}
