package cmd

import (
	"context"
	"fmt"
	"strings"
	"time" // Import time package

	"github.com/andygrunwald/go-jira"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/jirautils"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/localdb" // Use updated localdb
	"github.com/spf13/cobra"
)

// Flags specific to sync-jira command
var (
	syncJiraSearchLabel string
	syncTestDBPath      string
	syncDryRun          bool
)

var SyncJiraCmd = &cobra.Command{
	Use:   "sync-jira",
	Short: "Sync Jira tickets with local database",
	Long: `Scans Jira for flaky test tickets and ensures they exist in the local database.

This command performs the following actions:
1. Searches Jira for all tickets matching the specified label (default: flaky_test).
2. Fetches ticket summary and assignee information.
3. Compares the found Jira tickets against the local database (by Jira Key).
4. Adds entries to the local database for any Jira tickets not found locally.
   - Note: TestPackage will be empty for newly added entries as it cannot be reliably determined from Jira.
5. Updates the Assignee ID in the local database if it differs from the one in Jira.
6. Use --dry-run to preview changes without modifying the local database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Set default label if not provided
		if syncJiraSearchLabel == "" {
			syncJiraSearchLabel = "flaky_test" // Default label
		}
		log.Info().Str("label", syncJiraSearchLabel).Msg("Using Jira search label")

		// 2) Load local DB
		// LoadDBWithPath now returns an empty DB if file not found, but returns error on decode failure etc.
		db, err := localdb.LoadDBWithPath(syncTestDBPath)
		if err != nil {
			// If error is critical (not just file not found), log and return
			log.Error().Err(err).Str("path", syncTestDBPath).Msg("Failed to load or decode local DB")
			return fmt.Errorf("failed to load local DB from %s: %w", syncTestDBPath, err)
			// If we wanted to continue with an empty DB even on decode errors:
			// log.Warn().Err(err).Msg("Failed to load local DB; continuing with empty DB.")
			// db = localdb.NewDBWithPath(syncTestDBPath) // Ensure path is set correctly on new DB
		}

		// 3) Get Jira client
		client, err := jirautils.GetJiraClient()
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Jira client")
			return err // Cannot proceed without Jira client
		}
		log.Info().Msg("Jira client created successfully.")

		// 4) Search for all flaky test tickets in Jira
		jql := fmt.Sprintf(`labels = "%s" ORDER BY created DESC`, syncJiraSearchLabel)
		var startAt int
		var allIssues []jira.Issue
		totalTicketsInJira := 0 // Keep track even before appending

		log.Info().Msg("Searching Jira for tickets...")
		for {
			// Fetch required fields: summary (for test name), assignee
			issues, resp, searchErr := client.Issue.SearchWithContext(context.Background(), jql, &jira.SearchOptions{
				StartAt:    startAt,
				MaxResults: 50, // Fetch in batches
				Fields:     []string{"summary", "assignee"},
			})
			if searchErr != nil {
				// Attempt to read response for more detailed error
				errMsg := jirautils.ReadJiraErrorResponse(resp)
				log.Error().Err(searchErr).Str("jql", jql).Str("response", errMsg).Msg("Error searching Jira")
				return fmt.Errorf("error searching Jira: %w (response: %s)", searchErr, errMsg)
			}

			if resp != nil {
				totalTicketsInJira = resp.Total // Get total count from the response
			}

			if len(issues) == 0 {
				break // No more issues found
			}
			allIssues = append(allIssues, issues...)
			startAt += len(issues)
			log.Debug().Int("fetched", len(issues)).Int("total_fetched", startAt).Int("jira_total", totalTicketsInJira).Msg("Fetched batch of issues from Jira")

			// Safety break if StartAt exceeds total, although Jira should handle this.
			if totalTicketsInJira > 0 && startAt >= totalTicketsInJira {
				break
			}
		}
		log.Info().Int("count", len(allIssues)).Msg("Finished fetching all matching Jira tickets.")

		// 5) Process each issue
		var addedCount int
		var updatedCount int // Counts tickets found in DB (assignee might or might not be updated)
		var skippedCount int
		var assigneeUpdatedCount int
		dbModified := false // Flag to track if save is needed

		// Get all current DB entries for efficient lookup by Jira Key
		// Note: This map uses JiraKey as the map key, NOT the internal pkg::name key.
		// This is specific to this command's logic for comparing against Jira search results.
		existingEntries := db.GetAllEntries()
		entryMapByJiraKey := make(map[string]localdb.Entry, len(existingEntries))
		for _, entry := range existingEntries {
			if entry.JiraTicket != "" {
				entryMapByJiraKey[entry.JiraTicket] = entry
			}
		}
		log.Debug().Int("count", len(entryMapByJiraKey)).Msg("Created map of existing DB entries by Jira Key.")

		for _, issue := range allIssues {
			// Extract test name from summary (using the existing helper)
			summary := issue.Fields.Summary
			testName := extractTestName(summary)
			if testName == "" {
				log.Warn().Str("summary", summary).Str("key", issue.Key).Msg("Could not extract test name from summary, skipping.")
				skippedCount++
				continue
			}

			// Get assignee ID from Jira issue fields
			var assigneeID string
			if issue.Fields.Assignee != nil {
				assigneeID = issue.Fields.Assignee.AccountID // Use AccountID
				if assigneeID == "" {
					log.Warn().Str("key", issue.Key).Str("assigneeName", issue.Fields.Assignee.Name).Msg("Assignee found but AccountID is empty, trying Name.")
					assigneeID = issue.Fields.Assignee.Name // Fallback? Check what Jira requires. AccountID is preferred.
				}
			}

			// Check if this ticket (by Jira Key) is already in our local DB map
			if entry, exists := entryMapByJiraKey[issue.Key]; exists {
				// Ticket exists in DB, check if assignee needs update
				updatedCount++ // Increment count of tickets found in DB
				if entry.AssigneeID != assigneeID {
					log.Info().Str("key", issue.Key).Str("old_assignee", entry.AssigneeID).Str("new_assignee", assigneeID).Msg("Assignee mismatch found.")
					if !syncDryRun {
						// Update the entry using UpsertEntry, preserving existing fields
						errUpsert := db.UpsertEntry(entry.TestPackage, entry.TestName, entry.JiraTicket, entry.SkippedAt, assigneeID) // Pass new assigneeID
						if errUpsert != nil {
							log.Error().Err(errUpsert).Str("key", issue.Key).Msg("Failed to update assignee in local DB")
							// Continue processing other tickets? Yes.
						} else {
							assigneeUpdatedCount++ // Count successful updates
							dbModified = true      // Mark DB as modified
							log.Info().Str("key", issue.Key).Str("new_assignee", assigneeID).Msg("Successfully updated assignee in local DB.")
						}
					} else {
						// In dry run, just log the potential update and increment counter
						assigneeUpdatedCount++
						log.Info().Str("key", issue.Key).Str("new_assignee", assigneeID).Msg("[Dry Run] Would update assignee.")
					}
				} else {
					// Assignee matches, nothing to do for this entry
					log.Debug().Str("key", issue.Key).Msg("Existing ticket found in DB, assignee matches.")
				}
			} else {
				// Ticket NOT found in DB, add it
				log.Info().Str("key", issue.Key).Str("test", testName).Str("assignee", assigneeID).Msg("New ticket found in Jira, adding to DB.")
				if !syncDryRun {
					// Add new entry using UpsertEntry
					// TestPackage is unknown, SkippedAt is zero time for new entries
					errUpsert := db.UpsertEntry("", testName, issue.Key, time.Time{}, assigneeID)
					if errUpsert != nil {
						log.Error().Err(errUpsert).Str("key", issue.Key).Msg("Failed to add new ticket to local DB")
						// Continue processing other tickets? Yes.
					} else {
						addedCount++      // Count successful additions
						dbModified = true // Mark DB as modified
						log.Info().Str("key", issue.Key).Msg("Successfully added new ticket to local DB.")
					}
				} else {
					// In dry run, just log the potential addition and increment counter
					addedCount++
					log.Info().Str("key", issue.Key).Str("test", testName).Str("assignee", assigneeID).Msg("[Dry Run] Would add new ticket.")
				}
			}
		} // End of processing loop

		// 6) Save DB if not in dry run mode and modifications occurred
		if !syncDryRun && dbModified {
			if err := db.Save(); err != nil {
				log.Error().Err(err).Msg("Failed to save updated local DB")
				// Return error here as saving failed
				return fmt.Errorf("failed to save local DB changes: %w", err)
			}
			log.Info().Str("path", db.FilePath()).Msg("Local DB saved with updates.")
		} else if syncDryRun {
			log.Info().Msg("Dry run finished. No changes saved to local DB.")
		} else {
			log.Info().Msg("No changes detected requiring DB save.")
		}

		// 7) Print summary
		fmt.Printf("\n--- Sync Summary ---\n")
		fmt.Printf("Total Jira tickets scanned (label: %s): %d\n", syncJiraSearchLabel, len(allIssues))
		fmt.Printf("Tickets added to local DB:              %d\n", addedCount)
		fmt.Printf("Tickets already in local DB:            %d\n", updatedCount)
		fmt.Printf("Assignees updated in local DB:          %d\n", assigneeUpdatedCount)
		fmt.Printf("Tickets skipped (parse error):          %d\n", skippedCount)
		if syncDryRun {
			fmt.Printf("\n** Dry Run Mode: No changes were saved to the local database. **\n")
		} else if dbModified {
			fmt.Printf("\nLocal database updated: %s\n", db.FilePath())
		} else {
			fmt.Printf("\nLocal database is already up-to-date.\n")
		}

		return nil // Success
	},
}

func init() {
	// Use exported DefaultDBPath from localdb package for the default value
	SyncJiraCmd.Flags().StringVar(&syncTestDBPath, "test-db-path", localdb.DefaultDBPath(), "Path to the flaky test JSON database")
	SyncJiraCmd.Flags().StringVar(&syncJiraSearchLabel, "jira-search-label", "flaky_test", "Jira label used to find flaky test tickets") // Default set here
	SyncJiraCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "If true, only show what would be changed without saving")
}

// extractTestName attempts to extract the test name from a ticket summary.
// NOTE: This is fragile and depends on a consistent summary format.
func extractTestName(summary string) string {
	// Expected format variations:
	// "Fix Flaky Test: TestName (X% flake rate)"
	// "Fix Flaky Test: TestName"
	prefix := "Fix Flaky Test: "
	if !strings.HasPrefix(summary, prefix) {
		// Maybe try other prefixes or patterns if needed?
		log.Debug().Str("summary", summary).Msg("Summary does not match expected prefix.")
		return "" // Doesn't match expected format
	}

	// Get the part after the prefix
	testPart := strings.TrimPrefix(summary, prefix)

	// Find the start of the flake rate part " ("
	flakeRateIndex := strings.Index(testPart, " (")
	testName := ""
	if flakeRateIndex != -1 {
		// Flake rate part exists, take the substring before it
		testName = testPart[:flakeRateIndex]
	} else {
		// No flake rate part found, assume the whole remaining part is the test name
		testName = testPart
	}

	// Final trim just in case
	return strings.TrimSpace(testName)
}
