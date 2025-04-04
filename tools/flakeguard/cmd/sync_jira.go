package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/jirautils"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/localdb"
	"github.com/spf13/cobra"
)

var (
	syncJiraSearchLabel string
	syncTestDBPath      string
	syncDryRun          bool
	updateAssignees     bool
)

var SyncJiraCmd = &cobra.Command{
	Use:   "sync-jira",
	Short: "Sync Jira tickets with local database",
	Long: `Searches for all flaky test tickets in Jira that aren't yet tracked in the local database.
This command will:
1. Search Jira for all tickets with the flaky_test label
2. Compare against the local database
3. Add any missing tickets to the local database
4. Update assignee information from Jira
5. Optionally show which tickets were added/updated (--dry-run)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Set default label if not provided
		if syncJiraSearchLabel == "" {
			syncJiraSearchLabel = "flaky_test"
		}

		// 2) Load local DB
		db, err := localdb.LoadDBWithPath(syncTestDBPath)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to load local DB; continuing with empty DB.")
			db = localdb.NewDB()
		}

		// 3) Get Jira client
		client, err := jirautils.GetJiraClient()
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Jira client")
			return err
		}

		// 4) Search for all flaky test tickets in Jira
		jql := fmt.Sprintf(`labels = "%s" ORDER BY created DESC`, syncJiraSearchLabel)
		var startAt int
		var allIssues []jira.Issue

		for {
			issues, resp, err := client.Issue.SearchWithContext(context.Background(), jql, &jira.SearchOptions{
				StartAt:    startAt,
				MaxResults: 50,                              // Fetch in batches of 50
				Fields:     []string{"summary", "assignee"}, // Include assignee field
			})
			if err != nil {
				return fmt.Errorf("error searching Jira: %w (resp: %v)", err, resp)
			}
			if len(issues) == 0 {
				break
			}
			allIssues = append(allIssues, issues...)
			startAt += len(issues)
		}

		// 5) Process each issue
		var added int
		var updated int
		var skipped int
		var assigneeUpdated int

		// Get all entries for efficient updating
		entries := db.GetAllEntries()
		entriesMap := make(map[string]*localdb.Entry) // map by Jira ticket key
		for i := range entries {
			entriesMap[entries[i].JiraTicket] = &entries[i]
		}

		for _, issue := range allIssues {
			// Extract test name from summary
			summary := issue.Fields.Summary
			testName := extractTestName(summary)
			if testName == "" {
				log.Warn().Msgf("Could not extract test name from summary: %s", summary)
				skipped++
				continue
			}

			// Get assignee ID if available
			var assigneeID string
			if issue.Fields.Assignee != nil {
				assigneeID = issue.Fields.Assignee.AccountID
			}

			// Check if this ticket is already in the local DB
			if entry, exists := entriesMap[issue.Key]; exists {
				if assigneeID != "" && entry.AssigneeID != assigneeID {
					if !syncDryRun {
						entry.AssigneeID = assigneeID
						// Update the entry in the DB
						db.UpdateEntry(*entry)
					}
					assigneeUpdated++
					log.Info().Msgf("Updated assignee for ticket %s to %s", issue.Key, assigneeID)
				}
				updated++
			} else {
				if !syncDryRun {
					// Create new entry with assignee information
					entry := localdb.Entry{
						TestPackage: "", // Empty as we can't reliably extract it
						TestName:    testName,
						JiraTicket:  issue.Key,
						AssigneeID:  assigneeID,
					}
					db.AddEntry(entry)
				}
				added++
				log.Info().Msgf("Added ticket %s for test %s (assignee: %s)", issue.Key, testName, assigneeID)
			}
		}

		// 6) Save DB if not in dry run mode
		if !syncDryRun && (added > 0 || assigneeUpdated > 0) {
			if err := db.Save(); err != nil {
				log.Error().Err(err).Msg("Failed to save local DB")
				return err
			}
			log.Info().Msgf("Local DB has been updated at: %s", db.FilePath())
		}

		// 7) Print summary
		fmt.Printf("\nSummary:\n")
		fmt.Printf("Total Jira tickets found: %d\n", len(allIssues))
		fmt.Printf("New tickets added to DB: %d\n", added)
		fmt.Printf("Existing tickets found: %d\n", updated)
		fmt.Printf("Assignees updated: %d\n", assigneeUpdated)
		fmt.Printf("Skipped (could not parse): %d\n", skipped)
		if syncDryRun {
			fmt.Printf("\nThis was a dry run. No changes were made to the local database.\n")
		}

		return nil
	},
}

func init() {
	SyncJiraCmd.Flags().StringVar(&syncJiraSearchLabel, "jira-search-label", "", "Jira label to filter existing tickets (default: flaky_test)")
	SyncJiraCmd.Flags().StringVar(&syncTestDBPath, "test-db-path", "", "Path to the flaky test JSON database (default: ~/.flaky_test_db.json)")
	SyncJiraCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "If true, only show what would be added without making changes")
	InitCommonFlags(SyncJiraCmd)
}

// extractTestName attempts to extract the test name from a ticket summary
func extractTestName(summary string) string {
	// Expected format: "Fix Flaky Test: TestName (X% flake rate)"
	parts := strings.Split(summary, ": ")
	if len(parts) != 2 {
		return ""
	}
	testPart := parts[1]
	// Split on " (" to remove the flake rate part
	testName := strings.Split(testPart, " (")[0]
	return strings.TrimSpace(testName)
}
