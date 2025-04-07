package cmd

import (
	// Keep if still needed locally, otherwise remove
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/briandowns/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"

	// Import the new mapping package
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/jirautils"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/localdb"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/mapping"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/model"
	"github.com/spf13/cobra"
)

// Command flags (keep existing ones)
var (
	ticketsJSONPath string
	jiraComment     bool
	ticketsDryRun   bool
	hideSkipped     bool
	missingPillars  bool
	// Add flags for mapping files if they weren't already common
	userMappingPath     string
	userTestMappingPath string // Although not directly used for assignment here, load it for consistency/future use
)

var ReviewTicketsCmd = &cobra.Command{
	Use:   "review-tickets",
	Short: "Review tickets from --test-db-path",
	Long: `Interactively review tickets from --test-db-path.
    
Actions:
  [s] mark as skipped (and optionally post a comment to the Jira ticket)
  [u] unskip a ticket
  [i] set pillar name based on user mapping (if assignee exists)
  [p] previous ticket
  [n] next ticket
  [q] quit`, // Added 'p' to description
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Load the local JSON database.
		db, err := localdb.LoadDBWithPath(ticketsJSONPath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load local DB")
			os.Exit(1) // Keep exiting on critical load failure
		}

		// 2) Load Mappings using the new package
		userMap, err := mapping.LoadUserMappings(userMappingPath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load user mappings")
			return err // Return error to cobra
		}
		// Load test mappings even if not used for assignment here, maybe for future validation?
		_, err = mapping.LoadUserTestMappings(userTestMappingPath)
		if err != nil {
			// Non-fatal, just log a warning
			log.Warn().Err(err).Msg("Failed to load user test mappings, continuing...")
		}

		// 3) Retrieve all entries from the DB.
		entries := db.GetAllEntries()
		if len(entries) == 0 {
			log.Info().Msg("No tickets found in local DB") // Changed to Info
			return nil
		}

		// Convert entries to model.FlakyTicket
		tickets := make([]model.FlakyTicket, len(entries))
		for i, entry := range entries {
			tickets[i] = model.FlakyTicket{
				TestPackage:     entry.TestPackage,
				TestName:        entry.TestName,
				ExistingJiraKey: entry.JiraTicket,
				SkippedAt:       entry.SkippedAt,
				AssigneeId:      entry.AssigneeID, // Load AssigneeID from DB
			}

			// Check if the assignee from the DB exists in the user map
			if entry.AssigneeID != "" {
				if _, exists := userMap[entry.AssigneeID]; !exists {
					tickets[i].MissingUserMapping = true
					log.Debug().Str("assignee", entry.AssigneeID).Str("test", entry.TestName).Msg("Assignee from DB not found in user_mapping.json")
				}
			}
		}

		// 4) Filter based on flags (hideSkipped)
		if hideSkipped {
			filtered := make([]model.FlakyTicket, 0, len(tickets))
			for _, t := range tickets {
				if t.SkippedAt.IsZero() {
					filtered = append(filtered, t)
				}
			}
			tickets = filtered
			if len(tickets) == 0 {
				log.Info().Msg("No non-skipped tickets found matching criteria.")
				return nil
			}
		}

		// 5) Setup Jira client
		jiraClient, clientErr := jirautils.GetJiraClient()
		if clientErr != nil {
			log.Warn().Msgf("Jira client not available: %v. Running in offline mode.", clientErr)
			jiraClient = nil // Ensure it's nil if there's an error
		}

		// 6) Fetch pillar names (only if Jira client exists)
		if jiraClient != nil {
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Suffix = " Fetching pillar names from Jira..."
			s.Start()

			// Collect Jira keys that need pillar names
			var jiraKeys []string
			keyToIndex := make(map[string][]int) // Map key to indices in the tickets slice
			for i, t := range tickets {
				if t.ExistingJiraKey != "" && t.PillarName == "" { // Only fetch if not already known (and key exists)
					jiraKeys = append(jiraKeys, t.ExistingJiraKey)
					keyToIndex[t.ExistingJiraKey] = append(keyToIndex[t.ExistingJiraKey], i)
				}
			}
			jiraKeys = uniqueStrings(jiraKeys) // Avoid duplicate JQL queries if multiple local entries point to the same ticket

			if len(jiraKeys) > 0 {
				// Batch processing logic (remains the same)
				batchSize := 50
				for i := 0; i < len(jiraKeys); i += batchSize {
					end := i + batchSize
					if end > len(jiraKeys) {
						end = len(jiraKeys)
					}
					batch := jiraKeys[i:end]
					jql := fmt.Sprintf("key IN (%s)", strings.Join(batch, ","))
					issues, _, err := jiraClient.Issue.Search(jql, &jira.SearchOptions{
						Fields:     []string{"key", "customfield_11016"}, // customfield_11016 is Pillar Name
						MaxResults: batchSize,
					})

					if err != nil {
						log.Warn().Err(err).Msgf("Failed to fetch pillar names for batch of tickets starting at index %d", i)
						continue // Skip this batch on error
					}

					// Update tickets with pillar names
					for _, issue := range issues {
						if indices, found := keyToIndex[issue.Key]; found {
							pillarValue := ""
							if issue.Fields != nil {
								// Safely access the custom field
								if pillarFieldRaw, ok := issue.Fields.Unknowns["customfield_11016"]; ok && pillarFieldRaw != nil {
									if pillarField, ok := pillarFieldRaw.(map[string]interface{}); ok {
										if value, ok := pillarField["value"].(string); ok {
											pillarValue = value
										}
									}
								}
							}
							if pillarValue != "" {
								for _, ticketIdx := range indices {
									if ticketIdx < len(tickets) { // Bounds check
										tickets[ticketIdx].PillarName = pillarValue
										log.Debug().Str("ticket", issue.Key).Str("pillar", pillarValue).Msg("Pillar name fetched from Jira")
									}
								}
							} else {
								log.Debug().Str("ticket", issue.Key).Msg("Pillar name field (customfield_11016) not found or empty in Jira response")
							}
						}
					}
					s.Suffix = fmt.Sprintf(" Fetching pillar names from Jira... (%d/%d)", end, len(jiraKeys))
				}
			} else {
				log.Info().Msg("No tickets require pillar name fetching from Jira.")
			}
			s.Stop()
			fmt.Println() // Add a newline after spinner stops
		}

		// 7) Filter by missing pillars AFTER fetching
		if missingPillars {
			filtered := make([]model.FlakyTicket, 0, len(tickets))
			for _, t := range tickets {
				// A ticket is considered missing pillar if it has a Jira Key but no PillarName fetched/set
				if t.ExistingJiraKey != "" && t.PillarName == "" {
					filtered = append(filtered, t)
				}
			}
			tickets = filtered
			if len(tickets) == 0 {
				log.Info().Msg("No tickets found with missing pillar names.")
				return nil
			}
		}

		// 8) Initialize Bubble Tea model
		m := initialTicketsModel(tickets, userMap) // Pass only userMap, testPatternMap not directly used in TUI actions here
		m.JiraClient = jiraClient
		m.LocalDB = db
		m.JiraComment = jiraComment
		m.DryRun = ticketsDryRun

		// 9) Run TUI
		finalModel, err := tea.NewProgram(m).Run()
		if err != nil {
			log.Error().Err(err).Msg("Error running tickets TUI")
			os.Exit(1)
		}
		_ = finalModel // Use the final model if needed, e.g., for summary stats

		// 10) Save the local DB
		if err := db.Save(); err != nil {
			log.Error().Err(err).Msg("Failed to save local DB")
			// Don't exit here, just log the error
		} else {
			log.Info().Str("path", db.FilePath()).Msg("Local DB saved successfully")
		}
		return nil
	},
}

// Helper function for unique strings (used for Jira keys)
func uniqueStrings(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	j := 0
	for _, v := range input {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		input[j] = v
		j++
	}
	return input[:j]
}

func init() {
	ReviewTicketsCmd.Flags().StringVar(&ticketsJSONPath, "test-db-path", localdb.DefaultDBPath(), "Path to the JSON file containing tickets") // Use default path function
	ReviewTicketsCmd.Flags().BoolVar(&jiraComment, "jira-comment", true, "If true, post a comment to the Jira ticket when marking as skipped/unskipped")
	ReviewTicketsCmd.Flags().BoolVar(&ticketsDryRun, "dry-run", false, "If true, do not modify Jira tickets (comments, pillars) or save DB") // Updated help text
	ReviewTicketsCmd.Flags().BoolVar(&hideSkipped, "hide-skipped", false, "If true, do not show tests already marked as skipped")
	ReviewTicketsCmd.Flags().BoolVar(&missingPillars, "missing-pillars", false, "If true, only show tickets that have a Jira Key but no Pillar Name") // Updated help text
	// Make mapping paths consistent flags
	ReviewTicketsCmd.Flags().StringVar(&userMappingPath, "user-mapping-path", "user_mapping.json", "Path to the JSON file containing user mapping (JiraUserID -> PillarName)")
	ReviewTicketsCmd.Flags().StringVar(&userTestMappingPath, "user-test-mapping-path", "user_test_mapping.json", "Path to the JSON file containing user test mapping (Pattern -> JiraUserID)")
}

// -------------------------
// TUI Model and Functions
// -------------------------

// ticketModel represents the state of the TUI.
type ticketModel struct {
	tickets     []model.FlakyTicket
	index       int
	JiraClient  *jira.Client
	LocalDB     *localdb.DB
	JiraComment bool
	DryRun      bool
	quitting    bool
	infoMessage string
	userMap     map[string]mapping.UserMapping // Use mapping.UserMapping
	// testPatternMap not needed directly here anymore
}

// initialTicketsModel creates an initial model.
// Use the UserMapping type from the mapping package.
func initialTicketsModel(tickets []model.FlakyTicket, userMap map[string]mapping.UserMapping) ticketModel {
	// Ensure index is valid if tickets slice is empty
	idx := 0
	if len(tickets) == 0 {
		idx = -1 // Or handle appropriately in View/Update
	}
	return ticketModel{
		tickets: tickets,
		index:   idx,
		userMap: userMap,
	}
}

// Init is part of the Bubble Tea model interface.
func (m ticketModel) Init() tea.Cmd {
	return nil
}

// Update processes keypresses.
func (m ticketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle case where there are no tickets
	if m.index == -1 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			default:
				return m, nil // Ignore other keys if no tickets
			}
		default:
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "p": // Previous
			if m.index > 0 {
				m.index--
				m.infoMessage = "" // Clear message on navigation
			}
			return m, nil
		case "n": // Next
			if m.index < len(m.tickets)-1 {
				m.index++
				m.infoMessage = "" // Clear message on navigation
			}
			return m, nil
		case "s": // Skip
			t := &m.tickets[m.index] // Use pointer to modify in place
			if t.SkippedAt.IsZero() {
				now := time.Now()
				t.SkippedAt = now
				m.infoMessage = fmt.Sprintf("Marked as skipped at %s", now.UTC().Format(time.RFC822))
				// Update DB immediately (even in dry-run for TUI state, but save is conditional)
				if err := m.LocalDB.UpsertEntry(t.TestPackage, t.TestName, t.ExistingJiraKey, t.SkippedAt, t.AssigneeId); err != nil {
					log.Error().Err(err).Msg("Failed to update skip status in local DB state")
					m.infoMessage = "Error updating skip status in DB"
				}

				if !m.DryRun && m.JiraClient != nil && t.ExistingJiraKey != "" && m.JiraComment {
					comment := fmt.Sprintf("Test %s/%s marked as skipped via flakeguard on %s.", t.TestPackage, t.TestName, now.Format(time.RFC822))
					go func(key, comment string) { // Post comment in background to avoid blocking TUI
						err := jirautils.PostCommentToTicket(m.JiraClient, key, comment)
						if err != nil {
							// How to signal back to TUI? Could use a channel, but for now just log.
							log.Error().Err(err).Str("ticket", key).Msg("Failed to post skip comment to Jira")
						} else {
							log.Info().Str("ticket", key).Msg("Skip comment posted to Jira")
							// We can't easily update m.infoMessage from here without channels/Cmds
						}
					}(t.ExistingJiraKey, comment)
					m.infoMessage += fmt.Sprintf(" (Posting comment to %s...)", jirautils.GetJiraLink(t.ExistingJiraKey))
				} else if m.DryRun {
					m.infoMessage += " (Dry Run - Jira comment not sent)"
				}
			} else {
				m.infoMessage = "Already skipped."
			}
			return m, nil
		case "u": // Unskip
			t := &m.tickets[m.index] // Use pointer
			if !t.SkippedAt.IsZero() {
				unskippedAt := t.SkippedAt // Store old time for comment
				t.SkippedAt = time.Time{}  // Zero value means not skipped
				m.infoMessage = fmt.Sprintf("Marked as unskipped (was skipped at %s)", unskippedAt.UTC().Format(time.RFC822))
				// Update DB
				if err := m.LocalDB.UpsertEntry(t.TestPackage, t.TestName, t.ExistingJiraKey, t.SkippedAt, t.AssigneeId); err != nil {
					log.Error().Err(err).Msg("Failed to update unskip status in local DB state")
					m.infoMessage = "Error updating unskip status in DB"
				}

				if !m.DryRun && m.JiraClient != nil && t.ExistingJiraKey != "" && m.JiraComment {
					now := time.Now()
					comment := fmt.Sprintf("Test %s/%s marked as unskipped via flakeguard on %s (was previously skipped).", t.TestPackage, t.TestName, now.Format(time.RFC822))
					go func(key, comment string) { // Post comment in background
						err := jirautils.PostCommentToTicket(m.JiraClient, key, comment)
						if err != nil {
							log.Error().Err(err).Str("ticket", key).Msg("Failed to post unskip comment to Jira")
						} else {
							log.Info().Str("ticket", key).Msg("Unskip comment posted to Jira")
						}
					}(t.ExistingJiraKey, comment)
					m.infoMessage += fmt.Sprintf(" (Posting comment to %s...)", jirautils.GetJiraLink(t.ExistingJiraKey))
				} else if m.DryRun {
					m.infoMessage += " (Dry Run - Jira comment not sent)"
				}
			} else {
				m.infoMessage = "Not currently skipped."
			}
			return m, nil
		case "i": // Set Pillar Name based on mapping
			t := &m.tickets[m.index] // Use pointer
			if t.ExistingJiraKey == "" {
				m.infoMessage = "Cannot set pillar name: No associated Jira ticket key."
				return m, nil
			}
			if t.AssigneeId == "" {
				m.infoMessage = "Cannot set pillar name: Assignee ID is not set for this test."
				return m, nil
			}
			if m.JiraClient == nil {
				m.infoMessage = "Cannot set pillar name: Jira client is not available."
				return m, nil
			}
			if m.DryRun {
				m.infoMessage = "Cannot set pillar name: Running in Dry Run mode."
				return m, nil
			}

			// Find the user mapping for the ticket's assignee
			userMapping, exists := m.userMap[t.AssigneeId]
			if !exists {
				m.infoMessage = fmt.Sprintf("Cannot set pillar name: No user mapping found for Assignee ID %s.", t.AssigneeId)
				return m, nil
			}
			if userMapping.PillarName == "" {
				m.infoMessage = fmt.Sprintf("Cannot set pillar name: Pillar name is empty in mapping for Assignee ID %s.", t.AssigneeId)
				return m, nil
			}

			// Update the Jira ticket
			targetPillar := userMapping.PillarName
			m.infoMessage = fmt.Sprintf("Attempting to set Pillar Name to '%s' for %s...", targetPillar, jirautils.GetJiraLink(t.ExistingJiraKey))

			// Perform Jira update in background? Or block TUI? Let's block for immediate feedback.
			issueUpdate := &jira.Issue{
				Key: t.ExistingJiraKey,
				Fields: &jira.IssueFields{
					Unknowns: map[string]interface{}{
						"customfield_11016": map[string]interface{}{ // Pillar Name field ID
							"value": targetPillar,
						},
					},
				},
			}
			// Use UpdateIssue instead of Issue.Update for more flexibility if needed
			updatedIssue, resp, err := m.JiraClient.Issue.Update(issueUpdate)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to update pillar name for %s", t.ExistingJiraKey)
				log.Error().Err(err).Interface("response", resp).Msg(errMsg)
				m.infoMessage = fmt.Sprintf("%s: %v", errMsg, err)
			} else {
				log.Info().Str("ticket", updatedIssue.Key).Str("pillar", targetPillar).Msg("Pillar name updated successfully")
				m.infoMessage = fmt.Sprintf("Pillar name set to '%s' for %s", targetPillar, jirautils.GetJiraLink(updatedIssue.Key))
				// Update the local model state as well
				t.PillarName = targetPillar
			}
			return m, nil // Return after processing 'i'
		}
	}
	// Default: return current model if no key matched
	return m, nil
}

// View renders the current ticket and available actions.
func (m ticketModel) View() string {
	if m.quitting {
		// Consider saving DB on quit signal if not already handled by RunE final save
		return "Exiting tickets manager...\n"
	}
	if m.index == -1 || len(m.tickets) == 0 {
		return "No tickets loaded or matching filters.\n\n[q] quit\n"
	}
	if m.index >= len(m.tickets) {
		// Should not happen if navigation logic is correct, but handle defensively
		return "Error: Invalid ticket index.\n"
	}

	t := m.tickets[m.index]

	// Define styles.
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))          // Magenta/Purple
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))                      // Grey
	labelStyle := lipgloss.NewStyle().Bold(true).Width(12).Foreground(lipgloss.Color("39")) // Blue
	valueStyle := lipgloss.NewStyle()
	skippedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208")) // Orange
	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("40"))   // Green
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))   // Red
	dryRunStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("208")).Background(lipgloss.Color("235")).Padding(0, 1)
	actionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	helpStyle := lipgloss.NewStyle().Faint(true)

	var sb strings.Builder

	// Dry Run Banner
	if m.DryRun {
		sb.WriteString(dryRunStyle.Render("DRY RUN MODE") + "\n\n")
	}

	// Header
	sb.WriteString(headerStyle.Render(fmt.Sprintf("Ticket [%d / %d]", m.index+1, len(m.tickets))) + "\n\n")

	// Details Table
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Test Name:"), valueStyle.Render(t.TestName)))
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Package:"), valueStyle.Render(t.TestPackage)))

	if t.ExistingJiraKey != "" {
		sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Jira Key:"), valueStyle.Render(jirautils.GetJiraLink(t.ExistingJiraKey))))
	} else {
		sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Jira Key:"), valueStyle.Render("-")))
	}

	// Assignee Info
	assigneeVal := "-"
	if t.AssigneeId != "" {
		assigneeVal = t.AssigneeId
		// Check if mapping exists
		if _, exists := m.userMap[t.AssigneeId]; !exists {
			assigneeVal += errorStyle.Render(" (Mapping Missing!)")
		} else {
			// Optionally show Pillar name from map if userMap is available
			// assigneeVal += fmt.Sprintf(" (Pillar: %s)", m.userMap[t.AssigneeId].PillarName)
		}
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Assignee ID:"), valueStyle.Render(assigneeVal)))

	// Pillar Name (fetched from Jira or set)
	pillarVal := "-"
	if t.PillarName != "" {
		pillarVal = t.PillarName
	} else if t.ExistingJiraKey != "" {
		pillarVal = infoStyle.Render("(Not set in Jira)")
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Pillar Name:"), valueStyle.Render(pillarVal)))

	// Status
	statusLabel := labelStyle.Render("Status:")
	var statusValue string
	if !t.SkippedAt.IsZero() {
		statusValue = skippedStyle.Render(fmt.Sprintf("Skipped @ %s", t.SkippedAt.UTC().Format(time.RFC822)))
	} else {
		statusValue = activeStyle.Render("Active (Not Skipped)")
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", statusLabel, statusValue))

	// Info Message Area
	if m.infoMessage != "" {
		// Determine style based on content (simple check)
		infoMsgStyle := infoStyle
		lowerMsg := strings.ToLower(m.infoMessage)
		if strings.Contains(lowerMsg, "fail") || strings.Contains(lowerMsg, "error") {
			infoMsgStyle = errorStyle
		} else if strings.Contains(lowerMsg, "success") || strings.Contains(lowerMsg, "set to") || strings.Contains(lowerMsg, "posted") {
			infoMsgStyle = activeStyle // Use Green for success messages too
		}
		sb.WriteString("\n" + infoMsgStyle.Render(m.infoMessage) + "\n")
	}

	// Actions
	actions := []string{
		"[p]prev", "[n]next",
	}
	if t.SkippedAt.IsZero() {
		actions = append(actions, "[s]skip")
	} else {
		actions = append(actions, "[u]unskip")
	}
	if t.ExistingJiraKey != "" && t.AssigneeId != "" && m.JiraClient != nil && !m.DryRun {
		actions = append(actions, "[i]set_pillar")
	}
	actions = append(actions, "[q]quit")

	sb.WriteString("\n" + actionStyle.Render("Actions:") + "\n" + helpStyle.Render(strings.Join(actions, "  ")))

	return sb.String()
}
