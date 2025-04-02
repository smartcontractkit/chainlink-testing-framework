package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/jirautils"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/localdb"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/model"
	"github.com/spf13/cobra"
)

// Command flags
var (
	ticketsJSONPath string
	jiraComment     bool // if true, post a comment when marking as skipped
	ticketsDryRun   bool // if true, do not send anything to Jira
	hideSkipped     bool // if true, do not show skipped tests
)

// TicketsCmd is the new CLI command for managing tickets.
var TicketsCmd = &cobra.Command{
	Use:   "tickets",
	Short: "Manage tickets from flaky_test_db.json",
	Long: `Interactively manage your tickets.
	
Actions:
  [s] mark as skipped (and optionally post a comment to the Jira ticket)
  [u] unskip a ticket
  [i] set pillar name based on user mapping
  [n] next ticket
  [q] quit`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Load the local JSON database.
		db, err := localdb.LoadDBWithPath(ticketsJSONPath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load local DB")
			os.Exit(1)
		}

		// 3) Load user mapping
		var userMap map[string]UserMapping
		var testPatternMap map[string]UserTestMapping

		// Load user mapping file
		userMappingData, err := os.ReadFile(userMappingPath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read user mapping file")
			return err
		}
		var userMappings []UserMapping
		if err := json.Unmarshal(userMappingData, &userMappings); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal user mapping")
			return err
		}
		// Convert array to map
		userMap = make(map[string]UserMapping)
		for _, user := range userMappings {
			userMap[user.JiraUserID] = user
		}

		// Load user test mapping file
		userTestMappingData, err := os.ReadFile(userTestMappingPath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read user test mapping file")
			return err
		}
		var userTestMappings []UserTestMapping
		if err := json.Unmarshal(userTestMappingData, &userTestMappings); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal user test mapping")
			return err
		}
		// Convert array to map
		patternToUserID := make(map[string]string)
		for _, mapping := range userTestMappings {
			patternToUserID[mapping.Pattern] = mapping.JiraUserID
		}

		// 4) Retrieve all entries from the DB.
		entries := db.GetAllEntries()
		if len(entries) == 0 {
			log.Warn().Msg("No tickets found in local DB")
			return nil
		}

		// Convert entries to model.FlakyTicket (using SkippedAt only).
		tickets := make([]model.FlakyTicket, len(entries))
		for i, entry := range entries {
			tickets[i] = model.FlakyTicket{
				TestPackage:     entry.TestPackage,
				TestName:        entry.TestName,
				ExistingJiraKey: entry.JiraTicket,
				SkippedAt:       entry.SkippedAt,
			}

			// Map user ID based on test package pattern
			if userID := model.MapTestPackageToUser(entry.TestPackage, patternToUserID); userID != "" {
				if userMapping, exists := userMap[userID]; exists {
					tickets[i].AssigneeId = userID
					tickets[i].AssigneeName = userMapping.UserName
				}
			}
		}

		// If the hideSkipped flag is set, filter out tickets with a non-zero SkippedAt.
		if hideSkipped {
			filtered := make([]model.FlakyTicket, 0, len(tickets))
			for _, t := range tickets {
				if t.SkippedAt.IsZero() {
					filtered = append(filtered, t)
				}
			}
			tickets = filtered
		}

		// 5) Setup a Jira client (if available).
		jiraClient, clientErr := jirautils.GetJiraClient()
		if clientErr != nil {
			log.Warn().Msgf("Jira client not available: %v. Running in offline mode.", clientErr)
			jiraClient = nil
		}

		// 6) Initialize the Bubble Tea model.
		m := initialTicketsModel(tickets, userMap, testPatternMap)
		m.JiraClient = jiraClient
		m.LocalDB = db
		m.JiraComment = jiraComment
		m.DryRun = ticketsDryRun

		// 7) Run the TUI.
		finalModel, err := tea.NewProgram(m).Run()
		if err != nil {
			log.Error().Err(err).Msg("Error running tickets TUI")
			os.Exit(1)
		}
		_ = finalModel

		// 8) Save the local DB if any changes were made.
		if err := db.Save(); err != nil {
			log.Error().Err(err).Msg("Failed to save local DB")
		}
		return nil
	},
}

func init() {
	TicketsCmd.Flags().StringVar(&ticketsJSONPath, "test-db-path", "flaky_test_db.json", "Path to the JSON file containing tickets")
	TicketsCmd.Flags().BoolVar(&jiraComment, "jira-comment", true, "If true, post a comment to the Jira ticket when marking as skipped")
	TicketsCmd.Flags().BoolVar(&ticketsDryRun, "dry-run", false, "If true, do not send anything to Jira")
	TicketsCmd.Flags().BoolVar(&hideSkipped, "hide-skipped", false, "If true, do not show skipped tests")
	TicketsCmd.Flags().StringVar(&userMappingPath, "user-mapping-path", "user_mapping.json", "Path to the JSON file containing user mapping")
	TicketsCmd.Flags().StringVar(&userTestMappingPath, "user-test-mapping-path", "user_test_mapping.json", "Path to the JSON file containing user test mapping")
	InitCommonFlags(TicketsCmd)
}

// -------------------------
// TUI Model and Functions
// -------------------------

// ticketModel represents the state of the TUI.
type ticketModel struct {
	tickets        []model.FlakyTicket
	index          int
	JiraClient     *jira.Client
	LocalDB        localdb.DB
	JiraComment    bool
	DryRun         bool
	quitting       bool
	infoMessage    string
	userMap        map[string]UserMapping     // map of JiraUserID to UserMapping
	testPatternMap map[string]UserTestMapping // map of JiraUserID to UserTestMapping
}

// initialTicketsModel creates an initial model.
func initialTicketsModel(tickets []model.FlakyTicket, userMap map[string]UserMapping, testPatternMap map[string]UserTestMapping) ticketModel {
	return ticketModel{
		tickets:        tickets,
		index:          0,
		userMap:        userMap,
		testPatternMap: testPatternMap,
	}
}

// Init is part of the Bubble Tea model interface.
func (m ticketModel) Init() tea.Cmd {
	return nil
}

// Update processes keypresses.
func (m ticketModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Always allow quitting.
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

		// Navigation: previous and next.
		switch msg.String() {
		case "p":
			if m.index > 0 {
				m.index--
			}
			// Clear info message on navigation.
			m.infoMessage = ""
			return m, nil
		case "n":
			if m.index < len(m.tickets)-1 {
				m.index++
			}
			// Clear info message on navigation.
			m.infoMessage = ""
			return m, nil
		}

		// Action: mark as skipped.
		if msg.String() == "s" {
			t := m.tickets[m.index]
			// Only mark as skipped if not already skipped.
			if t.SkippedAt.IsZero() {
				t.SkippedAt = time.Now()
				// Optionally, post a comment to the Jira ticket if not in dry-run.
				if !m.DryRun && m.JiraClient != nil && t.ExistingJiraKey != "" && m.JiraComment {
					comment := fmt.Sprintf("Test marked as skipped on %s.", time.Now().Format(time.RFC822))
					err := jirautils.PostCommentToTicket(m.JiraClient, t.ExistingJiraKey, comment)
					if err != nil {
						log.Error().Err(err).Msgf("Failed to post comment to Jira ticket %s", t.ExistingJiraKey)
						m.infoMessage = fmt.Sprintf("Failed to post comment to Jira ticket %s", t.ExistingJiraKey)
					} else {
						m.infoMessage = fmt.Sprintf("Skip comment posted to Jira ticket: %s", jirautils.GetJiraLink(t.ExistingJiraKey))
					}
				}
				m.tickets[m.index] = t
				m.LocalDB.UpdateTicketStatus(t.TestPackage, t.TestName, t.SkippedAt)
			}
			return m, nil
		}

		// Action: unskip a ticket.
		if msg.String() == "u" {
			t := m.tickets[m.index]
			// Only unskip if the ticket is currently marked as skipped.
			if !t.SkippedAt.IsZero() {
				t.SkippedAt = time.Time{} // reset to zero value
				// Optionally, post a comment to the Jira ticket if not in dry-run.
				if !m.DryRun && m.JiraClient != nil && t.ExistingJiraKey != "" && m.JiraComment {
					comment := fmt.Sprintf("Test unskipped on %s.", time.Now().Format(time.RFC822))
					err := jirautils.PostCommentToTicket(m.JiraClient, t.ExistingJiraKey, comment)
					if err != nil {
						log.Error().Err(err).Msgf("Failed to post unskip comment to Jira ticket %s", t.ExistingJiraKey)
						m.infoMessage = fmt.Sprintf("Failed to post unskip comment to Jira ticket %s", t.ExistingJiraKey)
					} else {
						m.infoMessage = fmt.Sprintf("Unskip comment posted to Jira ticket: %s", jirautils.GetJiraLink(t.ExistingJiraKey))
					}
				}
				m.tickets[m.index] = t
				m.LocalDB.UpdateTicketStatus(t.TestPackage, t.TestName, t.SkippedAt)
			}
			return m, nil
		}

		// Action: set pillar name.
		if msg.String() == "i" {
			t := m.tickets[m.index]
			if t.ExistingJiraKey != "" {
				// Find the user mapping for the ticket's assignee
				if userMapping, exists := m.userMap[t.AssigneeId]; exists && !m.DryRun && m.JiraClient != nil {
					// Update the Jira ticket with the pillar name
					issue := &jira.Issue{
						Key: t.ExistingJiraKey,
						Fields: &jira.IssueFields{
							Unknowns: map[string]interface{}{
								"customfield_11016": map[string]interface{}{
									"value": userMapping.PillarName,
								},
							},
						},
					}
					_, _, err := m.JiraClient.Issue.Update(issue)
					if err != nil {
						log.Error().Err(err).Msgf("Failed to update pillar name for Jira ticket %s", t.ExistingJiraKey)
						m.infoMessage = fmt.Sprintf("Failed to update pillar name for Jira ticket %s", t.ExistingJiraKey)
					} else {
						m.infoMessage = fmt.Sprintf("Pillar name set to %s for ticket: %s", userMapping.PillarName, jirautils.GetJiraLink(t.ExistingJiraKey))
					}
				}
			}
			return m, nil
		}
	}
	return m, nil
}

// View renders the current ticket and available actions.
func (m ticketModel) View() string {
	if m.quitting {
		return "Exiting tickets manager..."
	}
	t := m.tickets[m.index]

	// Define styles.
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
	dryRunStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("208"))
	actionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))

	var view string

	// Show a dry-run indicator at the top.
	if m.DryRun {
		view += dryRunStyle.Render("DRY RUN MODE") + "\n\n"
	}

	// Build header and ticket details.
	view += headerStyle.Render(fmt.Sprintf("Ticket [%d/%d]\n", m.index+1, len(m.tickets))) + "\n"
	view += fmt.Sprintf("%s %s\n", labelStyle.Render("Test:"), t.TestName)
	view += fmt.Sprintf("%s %s\n", labelStyle.Render("Package:"), t.TestPackage)
	if t.ExistingJiraKey != "" {
		view += fmt.Sprintf("%s %s\n", labelStyle.Render("Jira:"), jirautils.GetJiraLink(t.ExistingJiraKey))
	}

	// Show status with color.
	if !t.SkippedAt.IsZero() {
		view += fmt.Sprintf("%s %s\n", labelStyle.Render("Status:"), fmt.Sprintf("skipped at: %s", t.SkippedAt.UTC().Format(time.RFC822)))
	} else {
		view += fmt.Sprintf("%s %s\n", labelStyle.Render("Status:"), infoStyle.Render("not skipped"))
	}

	// Display any info message.
	if m.infoMessage != "" {
		view += "\n" + infoStyle.Render(m.infoMessage) + "\n"
	}

	view += "\n" + actionStyle.Render("Actions:") + " [s] mark as skipped, [u] unskip, [i] set pillar name, [n] next, [q] quit"
	return view
}
