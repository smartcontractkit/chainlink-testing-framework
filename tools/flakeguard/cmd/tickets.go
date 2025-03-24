package cmd

import (
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
)

// TicketsCmd is the new CLI command for managing tickets.
var TicketsCmd = &cobra.Command{
	Use:   "tickets",
	Short: "Manage tickets from flaky_test_db.json",
	Long: `Interactively manage your tickets.
	
Actions:
  [s] mark as skipped (and optionally post a comment to the Jira ticket)
  [p] previous ticket
  [n] next ticket
  [q] quit
	
You can later extend this command to support additional actions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Load the local JSON database.
		db, err := localdb.LoadDBWithPath(ticketsJSONPath)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load local DB")
			os.Exit(1)
		}

		// 2) Retrieve all entries from the DB.
		entries := db.GetAllEntries()
		if len(entries) == 0 {
			log.Warn().Msg("No tickets found in local DB")
			return nil
		}

		// Convert entries to model.FlakyTicket
		tickets := make([]model.FlakyTicket, len(entries))
		for i, entry := range entries {
			tickets[i] = model.FlakyTicket{
				TestPackage:     entry.TestPackage,
				TestName:        entry.TestName,
				ExistingJiraKey: entry.JiraTicket,
				IsSkipped:       entry.IsSkipped,
				SkippedAt:       entry.SkippedAt,
			}
		}

		// 3) Setup a Jira client (if available).
		jiraClient, clientErr := jirautils.GetJiraClient()
		if clientErr != nil {
			log.Warn().Msgf("Jira client not available: %v. Running in offline mode.", clientErr)
			jiraClient = nil
		}

		// 4) Initialize the Bubble Tea model.
		m := initialTicketsModel(tickets)
		m.JiraClient = jiraClient
		m.LocalDB = db
		m.JiraComment = jiraComment
		m.DryRun = dryRun

		// 5) Run the TUI.
		finalModel, err := tea.NewProgram(m).Run()
		if err != nil {
			log.Error().Err(err).Msg("Error running tickets TUI")
			os.Exit(1)
		}
		_ = finalModel

		// 6) Save the local DB if any changes were made.
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
}

// -------------------------
// TUI Model and Functions
// -------------------------

// ticketModel represents the state of the TUI.
type ticketModel struct {
	tickets     []model.FlakyTicket
	index       int
	JiraClient  *jira.Client
	LocalDB     localdb.DB
	JiraComment bool
	DryRun      bool
	quitting    bool
}

// initialTicketsModel creates an initial model.
func initialTicketsModel(tickets []model.FlakyTicket) ticketModel {
	return ticketModel{
		tickets: tickets,
		index:   0,
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
			return m, nil
		case "n":
			if m.index < len(m.tickets)-1 {
				m.index++
			}
			return m, nil
		}

		// Action: mark as skipped.
		if msg.String() == "s" {
			t := m.tickets[m.index]
			// Only mark if not already skipped.
			if !t.IsSkipped {
				t.IsSkipped = true
				// Optionally, post a comment to the Jira ticket if not in dry-run.
				if !m.DryRun && m.JiraClient != nil && t.ExistingJiraKey != "" && m.JiraComment {
					comment := fmt.Sprintf("Test marked as skipped on %s.", time.Now().Format(time.RFC822))
					err := jirautils.PostCommentToTicket(m.JiraClient, t.ExistingJiraKey, comment)
					if err != nil {
						log.Error().Err(err).Msgf("Failed to post comment to Jira ticket %s", t.ExistingJiraKey)
					}
				}
				// Update local DB state with the current time.
				m.tickets[m.index] = t
				m.LocalDB.UpdateTicketStatus(t.TestPackage, t.TestName, t.IsSkipped, time.Now())
			}
			return m, nil
		}
	}
	return m, nil
}

// getJiraLink returns the full Jira URL for a given ticket key if JIRA_DOMAIN is set.
func getJiraLink(ticketKey string) string {
	domain := os.Getenv("JIRA_DOMAIN")
	if domain != "" {
		return fmt.Sprintf("https://%s/browse/%s", domain, ticketKey)
	}
	return ticketKey
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
	faintStyle := lipgloss.NewStyle().Faint(true)

	// Build header and ticket details.
	view := headerStyle.Render(fmt.Sprintf("Ticket [%d/%d]", m.index+1, len(m.tickets))) + "\n"
	view += fmt.Sprintf("Test: %s\n", t.TestName)
	view += fmt.Sprintf("Package: %s\n", t.TestPackage)
	if t.ExistingJiraKey != "" {
		view += fmt.Sprintf("Jira: %s\n", getJiraLink(t.ExistingJiraKey))
	}

	// Show status with color.
	if t.IsSkipped {
		view += fmt.Sprintf("Status: %s\n", fmt.Sprintf("Skipped At: %s\n", faintStyle.Render(t.SkippedAt.UTC().Format(time.RFC822))))
	} else {
		view += fmt.Sprintf("Status: %s\n", infoStyle.Render("Not Skipped"))
	}

	view += "\nActions: [s] mark as skipped, [p] previous, [n] next, [q] quit"
	return view
}
