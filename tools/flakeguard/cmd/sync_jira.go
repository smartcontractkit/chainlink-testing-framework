package cmd

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
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
	syncJiraSearchLabels []string
	syncTestDBPath       string
	syncDryRun           bool
)

var SyncJiraCmd = &cobra.Command{
	Use:   "sync-jira",
	Short: "Sync Jira tickets with local database",
	Long: `Scans Jira for flaky test tickets and ensures they exist in the local database.

This command performs the following actions:
1. Searches Jira for all tickets matching the specified labels (default: flaky_test, flaky_test).
2. Fetches ticket summary and assignee information.
3. Compares the found Jira tickets against the local database (by Jira Key).
4. Adds entries to the local database for any Jira tickets not found locally.
   - Note: TestPackage will be empty for newly added entries as it cannot be reliably determined from Jira.
5. Updates the Assignee ID in the local database if it differs from the one in Jira.
6. Use --dry-run to preview changes without modifying the local database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Strs("labels", syncJiraSearchLabels).Msg("Using Jira search labels")

		db, err := localdb.LoadDBWithPath(syncTestDBPath)
		if err != nil {
			log.Error().Err(err).Str("path", syncTestDBPath).Msg("Failed to load or decode local DB")
			return fmt.Errorf("failed to load local DB from %s: %w", syncTestDBPath, err)
		}

		client, err := jirautils.GetJiraClient()
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Jira client")
			return err
		}
		log.Info().Msg("Jira client created successfully.")

		var (
			jql                = fmt.Sprintf(`labels IN (%s) ORDER BY created DESC`, strings.Join(syncJiraSearchLabels, ","))
			startAt            int
			allIssues          []jira.Issue
			totalTicketsInJira int
		)

		log.Info().Str("jql", jql).Msg("Searching Jira for tickets...")
		for {
			issues, resp, searchErr := client.Issue.SearchWithContext(context.Background(), jql, &jira.SearchOptions{
				StartAt:    startAt,
				MaxResults: 50, // Fetch in batches
				Fields:     []string{"summary", "assignee", "status", "created", "description"},
			})
			if searchErr != nil || resp.StatusCode != http.StatusOK {
				errMsg := jirautils.ReadJiraErrorResponse(resp)
				log.Error().
					Err(searchErr).
					Int("status", resp.StatusCode).
					Str("jql", jql).
					Str("response", errMsg).
					Msg("Error searching Jira")
				return fmt.Errorf("error searching Jira: %w (response: %s)", searchErr, errMsg)
			}
			totalTicketsInJira = resp.Total

			if len(issues) == 0 {
				break
			}
			allIssues = append(allIssues, issues...)
			startAt += len(issues)
			log.Debug().Int("fetched", len(issues)).Int("total_fetched", startAt).Int("jira_total", totalTicketsInJira).Msg("Fetched batch of issues from Jira")

			if totalTicketsInJira > 0 && startAt >= totalTicketsInJira {
				break
			}
		}
		if len(allIssues) == 0 {
			log.Warn().Msg("No matching Jira tickets found")
		} else {
			log.Info().Int("count", len(allIssues)).Msg("Fetched all matching Jira tickets")
		}

		var (
			addedCount           int
			updatedCount         int
			closedCount          int
			skippedCount         int
			assigneeUpdatedCount int
			dbModified           bool
		)

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
		log.Debug().Int("count", len(entryMapByJiraKey)).Msg("Read existing DB entries by Jira Key")

		for _, issue := range allIssues {
			summary := issue.Fields.Summary
			testName, testPackage := extractTestNameAndPackage(summary, issue.Fields.Description)
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

			issueClosed := issue.Fields.Status.Name == "Closed"
			// Check if this ticket (by Jira Key) is already in our local DB map
			if entry, exists := entryMapByJiraKey[issue.Key]; exists { // Ticket exists in DB
				// Check if the ticket is closed and remove it as an active entry
				if issueClosed && entry.JiraTicket != "" {
					closedCount++
					log.Info().Str("key", issue.Key).Str("jira_ticket", entry.JiraTicket).Msg("Ticket is closed, marking as inactive.")
					if !syncDryRun {
						err := db.UpsertEntry(entry.TestPackage, entry.TestName, "", entry.SkippedAt, entry.AssigneeID)
						if err != nil {
							log.Error().Err(err).Str("key", issue.Key).Str("jira_ticket", entry.JiraTicket).Msg("Failed to mark ticket as inactive in local DB")
						} else {
							dbModified = true
							log.Info().Str("key", issue.Key).Str("jira_ticket", entry.JiraTicket).Msg("Successfully marked ticket as inactive in local DB.")
						}
					} else {
						log.Info().Str("key", issue.Key).Str("jira_ticket", entry.JiraTicket).Msg("[Dry Run] Would mark ticket as inactive.")
					}
					continue
				}

				updatedCount++
				// Check if the assignee needs updating
				if entry.AssigneeID != assigneeID {
					log.Info().Str("key", issue.Key).Str("old_assignee", entry.AssigneeID).Str("new_assignee", assigneeID).Msg("Assignee mismatch found.")
					if !syncDryRun {
						errUpsert := db.UpsertEntry(entry.TestPackage, entry.TestName, entry.JiraTicket, entry.SkippedAt, assigneeID) // Pass new assigneeID
						if errUpsert != nil {
							log.Error().Err(errUpsert).Str("key", issue.Key).Msg("Failed to update assignee in local DB")
						} else {
							assigneeUpdatedCount++
							dbModified = true
							log.Info().Str("key", issue.Key).Str("new_assignee", assigneeID).Msg("Successfully updated assignee in local DB.")
						}
					} else {
						// In dry run, just log the potential update and increment counter
						assigneeUpdatedCount++
						log.Info().Str("key", issue.Key).Str("new_assignee", assigneeID).Msg("[Dry Run] Would update assignee.")
					}
				} else {
					log.Debug().Str("key", issue.Key).Msg("Existing ticket found in DB, assignee matches.")
				}
			} else if !issueClosed {
				// Ticket NOT found in DB, add it
				log.Info().Str("key", issue.Key).Str("test", testName).Str("assignee", assigneeID).Msg("New ticket found in Jira, adding to DB.")
				if !syncDryRun {
					errUpsert := db.UpsertEntry(testPackage, testName, issue.Key, time.Time(issue.Fields.Created), assigneeID)
					if errUpsert != nil {
						log.Error().Err(errUpsert).Str("key", issue.Key).Msg("Failed to add new ticket to local DB")
					} else {
						addedCount++
						dbModified = true
						log.Info().Str("key", issue.Key).Msg("Successfully added new ticket to local DB.")
					}
				} else {
					addedCount++
					log.Info().Str("key", issue.Key).Str("test", testName).Str("assignee", assigneeID).Msg("[Dry Run] Would add new ticket.")
				}
			}
		}

		if !syncDryRun && dbModified {
			if err := db.Save(); err != nil {
				log.Error().Err(err).Msg("Failed to save updated local DB")
				return fmt.Errorf("failed to save local DB changes: %w", err)
			}
			log.Info().Str("path", db.FilePath()).Msg("Local DB saved with updates.")
		} else if syncDryRun {
			log.Info().Msg("Dry run finished. No changes saved to local DB.")
		} else {
			log.Info().Msg("No changes detected requiring DB save.")
		}

		fmt.Printf("\n--- Sync Summary ---\n")
		fmt.Printf("Total Jira tickets scanned:     %d\n", len(allIssues))
		fmt.Printf("Tickets added to local DB:      %d\n", addedCount)
		fmt.Printf("Tickets already in local DB:    %d\n", updatedCount)
		fmt.Printf("Tickets found closed:           %d\n", closedCount)
		fmt.Printf("Assignees updated in local DB:  %d\n", assigneeUpdatedCount)
		fmt.Printf("Tickets skipped (parse error):  %d\n", skippedCount)
		if syncDryRun {
			fmt.Printf("\n** Dry Run Mode: No changes were saved to the local database. **\n")
		} else if dbModified {
			fmt.Printf("\nLocal database updated: %s\n", db.FilePath())
		} else {
			fmt.Printf("\nLocal database is already up-to-date.\n")
		}

		return nil
	},
}

func init() {
	SyncJiraCmd.Flags().StringVar(&syncTestDBPath, "test-db-path", localdb.DefaultDBPath(), "Path to the flaky test JSON database")
	SyncJiraCmd.Flags().StringSliceVar(&syncJiraSearchLabels, "jira-search-labels", []string{"flaky_test", "flakey_test"}, "Jira labels used to find flaky test tickets")
	SyncJiraCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "If true, only show what would be changed without saving")
}

var (
	testNameRegex    = regexp.MustCompile(`(Test[^\s]+)`)
	testPackageRegex = regexp.MustCompile(`([a-zA-Z0-9.-]+(?:/[a-zA-Z0-9._-]+)+)`)
)

func extractTestNameAndPackage(summary, description string) (testName, testPackage string) {
	summary = strings.TrimSpace(summary)
	summary = strings.TrimPrefix(summary, "Fix Flaky Test:")
	testName = testNameRegex.FindString(summary)
	if testName == "" {
		testName = testNameRegex.FindString(description)
	}
	testPackage = testPackageRegex.FindString(description)
	if testPackage == "" {
		testPackage = testPackageRegex.FindString(summary)
	}
	return strings.TrimSpace(testName), strings.TrimSpace(testPackage)
}
