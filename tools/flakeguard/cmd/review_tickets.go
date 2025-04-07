package cmd

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/jirautils"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/localdb"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/mapping"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/model"
	"github.com/spf13/cobra"
)

// Command flags
var (
	ticketsJSONPath     string
	ticketsDryRun       bool
	missingPillars      bool
	userMappingPath     string
	userTestMappingPath string
)

// Renamed from TicketsCmd to ReviewTicketsCmd
var ReviewTicketsCmd = &cobra.Command{
	Use:   "review-tickets",
	Short: "Review tickets from the local database and sync Jira status",
	Long: `Interactively review tickets stored in the local database (--test-db-path).

Fetches current Pillar Name and Status from associated Jira tickets.
Allows setting the Pillar Name in Jira based on assignee mappings.

Data Source: Reads from the JSON file specified by --test-db-path.
Jira Interaction: Requires JIRA_* environment variables for fetching status/pillar and pillar updates.

Actions:
  [i] set Jira pillar name based on assignee mapping (updates Jira)
  [p] previous ticket
  [n] next ticket
  [q] quit`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := localdb.LoadDBWithPath(ticketsJSONPath)
		if err != nil {
			log.Error().Err(err).Str("path", ticketsJSONPath).Msg("Failed to load local DB")
			return fmt.Errorf("failed to load local DB: %w", err)
		}

		userMap, err := mapping.LoadUserMappings(userMappingPath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load user mappings")
			return err
		}
		_, err = mapping.LoadUserTestMappings(userTestMappingPath)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to load user test mappings, continuing...")
		}

		entries := db.GetAllEntries()
		if len(entries) == 0 {
			log.Info().Msg("No tickets found in local DB")
			return nil
		}
		log.Info().Int("count", len(entries)).Msg("Loaded entries from local DB.")

		tickets := make([]model.FlakyTicket, 0, len(entries))
		for _, entry := range entries {
			ticket := model.FlakyTicket{
				TestPackage:     entry.TestPackage,
				TestName:        entry.TestName,
				ExistingJiraKey: entry.JiraTicket,
				AssigneeId:      entry.AssigneeID,
			}
			if entry.AssigneeID != "" {
				if _, exists := userMap[entry.AssigneeID]; !exists {
					ticket.MissingUserMapping = true
					log.Debug().Str("assignee", entry.AssigneeID).Str("test", entry.TestName).Msg("Assignee from DB not found in user_mapping.json")
				}
			}
			tickets = append(tickets, ticket)
		}

		jiraClient, clientErr := jirautils.GetJiraClient()
		if clientErr != nil {
			log.Warn().Msgf("Jira client not available: %v. Running in offline mode (cannot fetch status/pillar or update).", clientErr)
			jiraClient = nil
		}

		// 6) Fetch pillar names AND STATUS from Jira (if client available and tickets exist)
		if jiraClient != nil && len(tickets) > 0 {
			log.Info().Msg("Attempting to fetch Pillar Names & Status from Jira...")

			var jiraKeysToFetch []string
			keyToIndexMap := make(map[string][]int)
			for i, t := range tickets {
				// Fetch if ticket has a key and we haven't already got Pillar OR Status
				if t.ExistingJiraKey != "" && (t.PillarName == "" || t.JiraStatus == "") {
					if _, exists := keyToIndexMap[t.ExistingJiraKey]; !exists {
						jiraKeysToFetch = append(jiraKeysToFetch, t.ExistingJiraKey)
					}
					keyToIndexMap[t.ExistingJiraKey] = append(keyToIndexMap[t.ExistingJiraKey], i)
				}
			}

			if len(jiraKeysToFetch) > 0 {
				log.Debug().Int("count", len(jiraKeysToFetch)).Msg("Fetching Pillar/Status for unique Jira keys.")
				batchSize := 50
				processedCount := 0 // Track keys processed for spinner
				for i := 0; i < len(jiraKeysToFetch); i += batchSize {
					end := i + batchSize
					if end > len(jiraKeysToFetch) {
						end = len(jiraKeysToFetch)
					}
					batch := jiraKeysToFetch[i:end]
					jql := fmt.Sprintf("key IN (%s)", strings.Join(batch, ","))

					// Request Status field in addition to Pillar Name field
					issues, _, searchErr := jiraClient.Issue.Search(jql, &jira.SearchOptions{
						Fields:     []string{"key", jirautils.PillarCustomFieldID, "status"},
						MaxResults: batchSize,
					})

					if searchErr != nil {
						log.Warn().Err(searchErr).Msgf("Failed to fetch Jira data batch (JQL: %s)", jql)
						continue
					}

					// Update tickets with pillar names and status
					for _, issue := range issues {
						processedCount++
						if indices, found := keyToIndexMap[issue.Key]; found {
							pillarValue := jirautils.ExtractPillarValue(issue) // Use helper
							jiraStatus := ""
							if issue.Fields != nil && issue.Fields.Status != nil {
								jiraStatus = issue.Fields.Status.Name // Get status name
							}

							log.Debug().Str("ticket", issue.Key).Str("pillar", pillarValue).Str("status", jiraStatus).Msg("Data retrieved from Jira.")

							for _, ticketIdx := range indices {
								if ticketIdx < len(tickets) { // Bounds check
									tickets[ticketIdx].PillarName = pillarValue
									tickets[ticketIdx].JiraStatus = jiraStatus
								}
							}
						}
					}
				}
				log.Info().Int("count", processedCount).Msg("Finished fetching Jira data.")
			} else {
				log.Info().Msg("No tickets required fetching data from Jira.")
			}
			fmt.Println()
		}

		if missingPillars {
			filtered := make([]model.FlakyTicket, 0, len(tickets))
			for _, t := range tickets {
				if t.ExistingJiraKey != "" && t.PillarName == "" {
					filtered = append(filtered, t)
				}
			}
			tickets = filtered
			if len(tickets) == 0 {
				log.Info().Msg("No tickets found with missing pillar names after filtering.")
				return nil
			}
			log.Info().Int("count", len(tickets)).Msg("Filtered view to show only tickets missing pillar names.")
		}

		if len(tickets) == 0 {
			log.Info().Msg("No tickets remaining after applying all filters.")
			return nil
		}

		m := initialTicketsModel(tickets, userMap)
		m.JiraClient = jiraClient
		m.LocalDB = db
		m.DryRun = ticketsDryRun

		program := tea.NewProgram(m)
		finalModel, err := program.Run()
		if err != nil {
			log.Error().Err(err).Msg("Error running tickets TUI")
			return fmt.Errorf("error running TUI: %w", err)
		}
		_ = finalModel

		if !ticketsDryRun {
			if db == nil {
				log.Error().Msg("Cannot save DB: DB instance is nil")
			} else if err := db.Save(); err != nil {
				log.Error().Err(err).Msg("Failed to save local DB")
			} else {
				log.Info().Str("path", db.FilePath()).Msg("Local DB saved.")
			}
		} else {
			log.Info().Msg("Dry Run: Skipping save of local DB.")
		}

		log.Info().Msg("Review Tickets command finished.")
		return nil
	},
}

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
	ReviewTicketsCmd.Flags().StringVar(&ticketsJSONPath, "test-db-path", localdb.DefaultDBPath(), "Path to the JSON file for the flaky test database")
	ReviewTicketsCmd.Flags().BoolVar(&ticketsDryRun, "dry-run", false, "Prevent changes to Jira (e.g., pillar updates)")
	ReviewTicketsCmd.Flags().BoolVar(&missingPillars, "missing-pillars", false, "Only show tickets with a Jira Key but no Pillar Name")
	ReviewTicketsCmd.Flags().StringVar(&userMappingPath, "user-mapping-path", "user_mapping.json", "Path to the user mapping JSON (JiraUserID -> PillarName)")
	ReviewTicketsCmd.Flags().StringVar(&userTestMappingPath, "user-test-mapping-path", "user_test_mapping.json", "Path to the user test mapping JSON (Pattern -> JiraUserID)")
}

// -------------------------
// TUI Model and Functions
// -------------------------

type ticketModel struct {
	tickets      []model.FlakyTicket
	index        int
	JiraClient   *jira.Client
	LocalDB      *localdb.DB
	DryRun       bool
	quitting     bool
	infoMessage  string
	errorMessage string
	userMap      map[string]mapping.UserMapping
}

func initialTicketsModel(tickets []model.FlakyTicket, userMap map[string]mapping.UserMapping) ticketModel {
	idx := 0
	if len(tickets) == 0 {
		idx = -1
	}
	return ticketModel{
		tickets: tickets,
		index:   idx,
		userMap: userMap,
	}
}

func (m ticketModel) Init() tea.Cmd {
	return nil
}

func (m ticketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.index == -1 {
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}
		}
		return m, nil
	}
	if m.quitting {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.infoMessage = ""
		m.errorMessage = ""

		key := msg.String()

		switch key {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			log.Info().Msg("Quit signal received.")
			return m, tea.Quit

		case "p": // Previous
			if m.index > 0 {
				m.index--
			} else {
				m.infoMessage = "Already at the first ticket."
			}
			return m, nil

		case "n": // Next
			if m.index < len(m.tickets)-1 {
				m.index++
			} else {
				m.infoMessage = "Already at the last ticket."
			}
			return m, nil

		case "i": // Set Pillar Name based on mapping
			if m.index < 0 || m.index >= len(m.tickets) {
				m.errorMessage = "Internal error: Invalid index for 'i' action."
				return m, nil
			}
			t := &m.tickets[m.index] // Use pointer

			if t.ExistingJiraKey == "" {
				m.errorMessage = "Cannot set pillar: No associated Jira key."
				return m, nil
			}
			if t.AssigneeId == "" {
				m.errorMessage = "Cannot set pillar: Assignee ID not set."
				return m, nil
			}
			if m.JiraClient == nil {
				m.errorMessage = "Cannot set pillar: Jira client unavailable."
				return m, nil
			}
			if m.DryRun {
				m.errorMessage = "Cannot set pillar: Dry Run mode enabled."
				return m, nil
			}

			// Find mapping and target pillar name
			userMapping, exists := m.userMap[t.AssigneeId]
			if !exists {
				m.errorMessage = fmt.Sprintf("Cannot set pillar: No mapping for assignee %s.", t.AssigneeId)
				return m, nil
			}
			targetPillar := userMapping.PillarName
			if targetPillar == "" {
				m.errorMessage = fmt.Sprintf("Cannot set pillar: Pillar name empty in mapping for %s.", t.AssigneeId)
				return m, nil
			}

			// Prevent setting if already set to target? Optional.
			if t.PillarName == targetPillar {
				m.infoMessage = fmt.Sprintf("Pillar name is already '%s'.", targetPillar)
				return m, nil
			}

			// Perform Jira Update (synchronously for immediate feedback)
			m.infoMessage = fmt.Sprintf("Attempting to set Pillar Name to '%s' for %s...", targetPillar, jirautils.GetJiraLink(t.ExistingJiraKey))
			updateErr := jirautils.UpdatePillarName(m.JiraClient, t.ExistingJiraKey, targetPillar)

			if updateErr != nil {
				errMsg := fmt.Sprintf("Failed to update pillar for %s", t.ExistingJiraKey)
				log.Error().Err(updateErr).Str("ticket", t.ExistingJiraKey).Str("pillar", targetPillar).Msg(errMsg)
				m.errorMessage = fmt.Sprintf("%s: %v", errMsg, updateErr)
				m.infoMessage = "" // Clear "Attempting..."
			} else {
				log.Info().Str("ticket", t.ExistingJiraKey).Str("pillar", targetPillar).Msg("Pillar name updated successfully in Jira.")
				m.infoMessage = fmt.Sprintf("Pillar name set to '%s' for %s", targetPillar, jirautils.GetJiraLink(t.ExistingJiraKey))
				// Update the local model state as well so the view refreshes correctly
				t.PillarName = targetPillar
				m.errorMessage = ""
			}
			return m, nil
		}
	}

	return m, nil
}

// View function: Displays Jira Status instead of SkippedAt
func (m ticketModel) View() string {
	if m.quitting {
		return "Exiting review...\n"
	}
	if m.index == -1 || len(m.tickets) == 0 {
		return "No tickets loaded or matching filters.\n\n[q] quit\n"
	}
	if m.index >= len(m.tickets) {
		return "Error: Invalid ticket index.\n\n[q] quit\n"
	}

	t := m.tickets[m.index]

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")).PaddingBottom(1)
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81"))
	errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196")).PaddingBottom(1) // Red
	labelStyle := lipgloss.NewStyle().Bold(true).Width(12).Foreground(lipgloss.Color("39"))
	valueStyle := lipgloss.NewStyle()
	statusStyle := lipgloss.NewStyle()
	switch strings.ToLower(t.JiraStatus) {
	case "done", "resolved", "closed":
		statusStyle = statusStyle.Foreground(lipgloss.Color("40")) // Green
	case "in progress", "in review":
		statusStyle = statusStyle.Foreground(lipgloss.Color("208")) // Orange
	case "to do", "backlog", "open":
		statusStyle = statusStyle.Foreground(lipgloss.Color("245")) // Grey
	}

	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208")) // Orange for warnings like missing mapping
	dryRunStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("208")).Background(lipgloss.Color("235")).Padding(0, 1)
	actionHelpStyle := lipgloss.NewStyle().Faint(true).PaddingTop(1)

	var sb strings.Builder

	// Dry Run Banner
	if m.DryRun {
		sb.WriteString(dryRunStyle.Render("DRY RUN MODE") + "\n\n")
	}

	// Error Message Area
	if m.errorMessage != "" {
		sb.WriteString(errorStyle.Render("Error: "+m.errorMessage) + "\n")
	}
	// Info Message Area
	if m.infoMessage != "" {
		sb.WriteString(infoStyle.Render(m.infoMessage) + "\n\n")
	}

	// Header
	sb.WriteString(headerStyle.Render(fmt.Sprintf("Review Ticket [%d / %d]", m.index+1, len(m.tickets))) + "\n")

	// Details Table
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Test Name:"), valueStyle.Render(t.TestName)))
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Package:"), valueStyle.Render(t.TestPackage)))
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Jira Key:"), valueStyle.Render(jirautils.GetJiraLink(t.ExistingJiraKey))))

	// Assignee Info
	assigneeVal := "-"
	if t.AssigneeId != "" {
		assigneeVal = t.AssigneeId
		if t.MissingUserMapping {
			assigneeVal += warningStyle.Render(" (Mapping Missing!)")
		}
	} else {
		assigneeVal = lipgloss.NewStyle().Faint(true).Render("(Not Set)")
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Assignee ID:"), valueStyle.Render(assigneeVal)))

	// Pillar Name
	pillarVal := "-"
	if t.PillarName != "" {
		pillarVal = t.PillarName
	} else if t.ExistingJiraKey != "" {
		// Indicate if fetched or just not set
		if m.JiraClient != nil { // Check if client was available to fetch
			pillarVal = infoStyle.Render("(Not set in Jira)")
		} else {
			pillarVal = infoStyle.Render("(Jira unavailable)")
		}
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Pillar Name:"), valueStyle.Render(pillarVal)))

	// Jira Status
	statusVal := t.JiraStatus
	if statusVal == "" {
		if t.ExistingJiraKey != "" {
			if m.JiraClient != nil {
				statusVal = infoStyle.Render("(Status not fetched)")
			} else {
				statusVal = infoStyle.Render("(Jira unavailable)")
			}
		} else {
			statusVal = "-"
		}
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Jira Status:"), statusStyle.Render(statusVal)))

	// Actions Help Text
	actions := []string{
		"[p]prev", "[n]next",
	}
	if t.ExistingJiraKey != "" &&
		t.AssigneeId != "" &&
		m.JiraClient != nil &&
		!m.DryRun &&
		t.PillarName == "" {
		actions = append(actions, "[i]set_pillar")
	}

	actions = append(actions, "[q]quit")

	sb.WriteString("\n" + lipgloss.NewStyle().Bold(true).Render("Actions:") + "\n" + actionHelpStyle.Render(strings.Join(actions, "  ")))

	return sb.String()
}
