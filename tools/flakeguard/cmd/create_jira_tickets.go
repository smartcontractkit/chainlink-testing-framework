package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/briandowns/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/jirautils"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/localdb"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/model"
	"github.com/spf13/cobra"
)

var (
	csvPath         string
	dryRun          bool
	jiraProject     string
	jiraIssueType   string
	jiraSearchLabel string // defaults to "flaky_test" if empty
	testDBPath      string
	skipExisting    bool
)

var CreateTicketsCmd = &cobra.Command{
	Use:   "create-tickets",
	Short: "Interactive TUI to confirm and create Jira tickets from CSV",
	Long: `Reads a CSV describing flaky tests and displays each proposed
ticket in a text-based UI. Press 'y' to confirm creation, 'n' to skip,
'e' if you know of an existing Jira ticket, or 'q' to quit.

- If --dry-run=false, we attempt to create the ticket in Jira (using
  environment variables JIRA_DOMAIN, JIRA_EMAIL, JIRA_API_KEY).
- We also search for existing tickets (label=flaky_test) with a matching
  test name before prompting creation.
- A local JSON "database" (via internal/localdb) remembers any tickets
  already mapped to tests, so you won't be prompted again in the future.
- After the TUI ends, a new CSV is produced, omitting any confirmed rows.
  The original CSV remains untouched.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Validate input
		if csvPath == "" {
			log.Error().Msg("CSV path is required (--csv-path)")
			os.Exit(1)
		}
		if jiraProject == "" {
			jiraProject = os.Getenv("JIRA_PROJECT_KEY")
		}
		if jiraProject == "" {
			log.Error().Msg("Jira project key is required (set --jira-project or JIRA_PROJECT_KEY env)")
			os.Exit(1)
		}
		if jiraSearchLabel == "" {
			jiraSearchLabel = "flaky_test"
		}

		// 2) Load local DB (test -> known Jira ticket)
		db, err := localdb.LoadDBWithPath(testDBPath)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to load local DB; continuing with empty DB.")
			db = localdb.NewDB()
		}

		// 3) Read CSV
		records, err := readFlakyTestsCSV(csvPath)
		if err != nil {
			log.Error().Err(err).Msg("Error reading CSV file")
			os.Exit(1)
		}
		if len(records) == 0 {
			log.Warn().Msg("CSV is empty!")
			return nil
		}

		originalRecords := records
		if len(records) <= 1 {
			log.Warn().Msg("No data rows found (CSV might have only a header).")
			return nil
		}
		dataRows := records[1:] // skip the header row

		// 4) Convert CSV rows -> FlakyTicket objects
		var tickets []model.FlakyTicket
		for i, row := range dataRows {
			if len(row) < 10 {
				log.Warn().Msgf("Skipping row %d (not enough columns)", i+1)
				continue
			}
			ft := rowToFlakyTicket(row)
			ft.RowIndex = i + 1

			// Check local DB for known Jira ticket
			if ticketID, found := db.Get(ft.TestPackage, ft.TestName); found {
				ft.ExistingJiraKey = ticketID
				ft.ExistingTicketSource = "localdb"
			}

			// Skip processing if flag is set and a Jira ticket ID exists
			if skipExisting && ft.ExistingJiraKey != "" {
				continue
			}

			tickets = append(tickets, ft)
		}
		if len(tickets) == 0 {
			log.Warn().Msg("No valid tickets found in CSV.")
			return nil
		}

		// 5) Attempt Jira client creation
		client, clientErr := jirautils.GetJiraClient()
		if clientErr != nil {
			log.Warn().Msgf("No valid Jira client: %v\nWill skip searching or creating tickets in Jira.", clientErr)
			client = nil
		}

		// 6) If we have a Jira client, do label-based search for existing tickets
		if client != nil {
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Suffix = " Searching for existing jira tickets..."
			s.Start()

			for i := range tickets {
				t := &tickets[i]
				if t.ExistingJiraKey == "" {
					key, err := findExistingTicket(client, jiraSearchLabel, *t)
					if err != nil {
						log.Warn().Msgf("Search failed for %q: %v", t.Summary, err)
					} else if key != "" {
						t.ExistingJiraKey = key
						t.ExistingTicketSource = "jira"
						db.Set(t.TestPackage, t.TestName, key)
					}
				}
			}
			s.Stop()
		}

		// 7) Create Bubble Tea model
		m := initialModel(tickets)
		m.DryRun = dryRun
		m.JiraProject = jiraProject
		m.JiraIssueType = jiraIssueType
		m.JiraClient = client
		m.originalRecords = originalRecords
		m.LocalDB = db

		// 8) Run TUI
		finalModel, err := tea.NewProgram(m).Run()
		if err != nil {
			log.Error().Err(err).Msg("Error running Bubble Tea program")
			os.Exit(1)
		}
		fm := finalModel.(tmodel)

		// 9) Save local DB with any new knowledge
		if err := fm.LocalDB.Save(); err != nil {
			log.Error().Err(err).Msg("Failed to save local DB")
		} else {
			fmt.Printf("Local DB has been updated at: %s\n", fm.LocalDB.FilePath())
		}

		return nil
	},
}

func init() {
	CreateTicketsCmd.Flags().StringVar(&csvPath, "csv-path", "", "Path to CSV file with flaky tests")
	CreateTicketsCmd.Flags().BoolVar(&dryRun, "dry-run", false, "If true, do not create tickets in Jira")
	CreateTicketsCmd.Flags().StringVar(&jiraProject, "jira-project", "", "Jira project key (default: JIRA_PROJECT_KEY env)")
	CreateTicketsCmd.Flags().StringVar(&jiraIssueType, "jira-issue-type", "Task", "Jira issue type")
	CreateTicketsCmd.Flags().StringVar(&jiraSearchLabel, "jira-search-label", "", "Jira label to filter existing tickets (default: flaky_test)")
	CreateTicketsCmd.Flags().StringVar(&testDBPath, "test-db-path", "", "Path to the flaky test JSON database (default: ~/.flaky_tes_db.json)")
	CreateTicketsCmd.Flags().BoolVar(&skipExisting, "skip-existing", false, "Skip processing tickets that already have a Jira ticket ID")
	InitCommonFlags(CreateTicketsCmd)
}

// -------------------------------------------------------------------------------------
// FlakyTicket Data Model
// -------------------------------------------------------------------------------------

// rowToFlakyTicket: build a ticket from one CSV row (index assumptions: pkg=0, testName=2, flakeRate=7, logs=9).
func rowToFlakyTicket(row []string) model.FlakyTicket {
	pkg := strings.TrimSpace(row[0])
	testName := strings.TrimSpace(row[2])
	flakeRateStr := strings.TrimSpace(row[7]) // Keep the string for display
	logs := strings.TrimSpace(row[9])

	t := model.FlakyTicket{
		TestPackage: pkg,
		TestName:    testName,
		// Summary and Description will be set later
		Valid: true, // Assume valid initially
	}

	// Priority Calculation
	var flakeRate float64
	var parseErr error
	if flakeRateStr == "" {
		log.Warn().Msgf("Missing Flake Rate for test %q (%s). Defaulting to Low priority.", testName, pkg)
		t.Priority = "Low" // Default priority if empty
		flakeRateStr = "0" // Use "0" for summary display if empty
	} else {
		flakeRate, parseErr = strconv.ParseFloat(flakeRateStr, 64)
		if parseErr != nil {
			log.Error().Err(parseErr).Msgf("Invalid Flake Rate '%s' for test %q (%s).", flakeRateStr, testName, pkg)
			t.Valid = false
			t.InvalidReason = fmt.Sprintf("Invalid Flake Rate: %s", flakeRateStr)
		} else {
			t.FlakeRate = flakeRate // Store numeric value
			// Determine priority based on numeric value
			switch {
			case flakeRate < 1.0:
				t.Priority = "Low"
			case flakeRate >= 1.0 && flakeRate < 3.0:
				t.Priority = "Medium"
			case flakeRate >= 3.0 && flakeRate < 5.0:
				t.Priority = "High"
			default: // >= 5.0
				t.Priority = "Very High"
			}
		}
	}

	// Use flakeRateStr for display in summary/description to keep the % format as originally intended
	summary := fmt.Sprintf("Fix Flaky Test: %s (%s%% flake rate)", testName, flakeRateStr)
	t.Summary = summary

	// Parse logs (same as before)
	var logSection string
	// ... (keep existing log parsing logic) ...
	if logs == "" {
		logSection = "(Logs not available)"
	} else {
		var lines []string
		runNumber := 1
		for _, link := range strings.Split(logs, ",") {
			link = strings.TrimSpace(link)
			if link == "" {
				continue
			}
			lines = append(lines, fmt.Sprintf("* [Run %d|%s]", runNumber, link))
			runNumber++
		}
		if len(lines) == 0 {
			logSection = "(Logs not available)"
		} else {
			logSection = strings.Join(lines, "\n")
		}
	}

	// Use Jira Wiki Markup (same as before)
	desc := fmt.Sprintf(`h2. Test Details:
* *Package:* %s
* *Test Name:* %s
* *Flake Rate:* %s%% in the last 7 days

h3. Test Logs:
%s

h3. Action Items:
# *Investigate:* Review logs to find the root cause.
# *Fix:* Address the underlying problem causing flakiness.
# *Rerun Locally:* Confirm the fix stabilizes the test.
# *Unskip:* Re-enable test in the CI pipeline once stable.
# *Ref:* [Follow guidelines in the Flaky Test Guide|https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/tools/flakeguard/e2e-flaky-test-guide.md].
`,
		pkg,
		testName,
		flakeRateStr,
		logSection,
	)
	t.Description = desc

	// Check required fields
	var missing []string
	if pkg == "" {
		missing = append(missing, "Package")
	}
	if testName == "" {
		missing = append(missing, "Test Name")
	}
	if flakeRateStr == "" && parseErr != nil { // Check if original was empty AND parsing failed (edge case)
		missing = append(missing, "Flake Rate (missing or invalid)")
	} else if flakeRateStr == "" {
		missing = append(missing, "Flake Rate (missing)")
	} else if parseErr != nil {
		missing = append(missing, "Flake Rate (invalid format)")
	}
	if logs == "" {
		missing = append(missing, "Logs")
	}

	// Only overwrite Valid/InvalidReason if previously valid
	if t.Valid && len(missing) > 0 {
		t.Valid = false
		t.InvalidReason = fmt.Sprintf("Missing/Invalid required fields: %s", strings.Join(missing, ", "))
	}

	return t
}

// -------------------------------------------------------------------------------------
// Jira Search
// -------------------------------------------------------------------------------------

func findExistingTicket(client *jira.Client, label string, ticket model.FlakyTicket) (string, error) {
	jql := fmt.Sprintf(`labels = "%s" AND summary ~ "%s" order by created DESC`, label, ticket.TestName)
	issues, resp, err := client.Issue.SearchWithContext(context.Background(), jql, &jira.SearchOptions{MaxResults: 1})
	if err != nil {
		return "", fmt.Errorf("error searching Jira: %w (resp: %v)", err, resp)
	}
	if len(issues) == 0 {
		return "", nil
	}
	return issues[0].Key, nil
}

// -------------------------------------------------------------------------------------
// Bubble Tea Model
// -------------------------------------------------------------------------------------

type tmodel struct {
	tickets         []model.FlakyTicket
	index           int
	confirmed       int
	skipped         int
	quitting        bool
	DryRun          bool
	JiraProject     string
	JiraIssueType   string
	JiraClient      *jira.Client
	originalRecords [][]string

	LocalDB localdb.DB // reference to our local DB

	mode       string // "normal", "promptExisting", or "ticketCreated"
	inputValue string // user-typed input for existing ticket
}

func initialModel(tickets []model.FlakyTicket) tmodel {
	return tmodel{
		tickets: tickets,
		index:   0,
		mode:    "normal",
	}
}

func (m tmodel) Init() tea.Cmd {
	return nil
}

func (m tmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode == "promptExisting" {
			return updatePromptExisting(m, msg)
		}
		if m.mode == "ticketCreated" {
			return updateTicketCreated(m, msg)
		}
		return updateNormalMode(m, msg)
	default:
		return m, nil
	}
}

func updateNormalMode(m tmodel, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.quitting || m.index >= len(m.tickets) {
		return updateQuit(m)
	}
	t := m.tickets[m.index]

	// Always allow 'q' (quit), 'e' (enter existing ticket), and 'd' (delete existing ticket)
	switch msg.String() {
	case "q", "esc", "ctrl+c":
		return updateQuit(m)
	case "e":
		m.mode = "promptExisting"
		m.inputValue = ""
		return m, nil
	case "d":
		// Only allow removal if an existing ticket is present
		if t.ExistingJiraKey != "" && m.JiraClient != nil {
			err := jirautils.DeleteTicketInJira(m.JiraClient, t.ExistingJiraKey)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to delete ticket %s", t.ExistingJiraKey)
			} else {
				t.ExistingJiraKey = ""
				t.ExistingTicketSource = ""
				m.tickets[m.index] = t
				m.LocalDB.Set(t.TestPackage, t.TestName, "")
			}
		}
		return m, nil
	}

	// If ticket is valid, allow create ('c') and skip ('n')
	if t.Valid {
		switch msg.String() {
		case "c":
			return updateConfirm(m)
		case "n":
			return updateSkip(m)
		}
		return m, nil
	}

	// For invalid tickets, default to skipping if any other key is pressed.
	return updateSkip(m)
}

func updatePromptExisting(m tmodel, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyRunes:
		m.inputValue += string(msg.Runes)
	case tea.KeyBackspace:
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}
	case tea.KeyEnter:
		t := m.tickets[m.index]
		t.ExistingJiraKey = m.inputValue
		t.ExistingTicketSource = "localdb"
		m.tickets[m.index] = t
		m.mode = "normal"
		m.inputValue = ""
		return updateSkip(m)
	case tea.KeyEsc:
		m.mode = "normal"
		m.inputValue = ""
	}
	return m, nil
}

func updateConfirm(m tmodel) (tea.Model, tea.Cmd) {
	i := m.index
	t := m.tickets[i]

	// Attempt Jira creation if not dry-run and we have a client.
	if !m.DryRun && m.JiraClient != nil {
		issueKey, err := jirautils.CreateTicketInJira(m.JiraClient, t.Summary, t.Description, m.JiraProject, m.JiraIssueType, t.AssigneeId, t.Priority)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to create Jira ticket: %s", t.Summary)
		} else {
			log.Info().Msgf("Created Jira ticket: %s (summary=%q)", issueKey, t.Summary)
			t.Confirmed = true
			t.ExistingJiraKey = issueKey
			t.ExistingTicketSource = "jira"
			m.LocalDB.Set(t.TestPackage, t.TestName, issueKey)
		}
	} else {
		t.Confirmed = true
		// Set a dummy ticket key for testing purposes in dry-run mode.
		t.ExistingJiraKey = "DRYRUN-1234"
	}
	m.tickets[i] = t
	m.confirmed++
	m.mode = "ticketCreated"
	return m, nil
}

func updateTicketCreated(m tmodel, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "e": // Still allow 'e' to go to the prompt mode
		m.mode = "promptExisting"
		m.inputValue = ""
		return m, nil
	case "q", "esc", "ctrl+c": // Still allow quitting
		return updateQuit(m)
	default: // Make continuing the default action for *any other key*
		m.mode = "normal"
		m.index++
		if m.index >= len(m.tickets) {
			m.quitting = true // Ensure it quits if this was the last ticket
		}
		return m, nil
	}
}

func updateSkip(m tmodel) (tea.Model, tea.Cmd) {
	m.skipped++
	m.index++
	if m.index >= len(m.tickets) {
		m.quitting = true
	}
	return m, nil
}

func updateQuit(m tmodel) (tea.Model, tea.Cmd) {
	m.quitting = true
	return m, tea.Quit
}

// View logic
func (m tmodel) View() string {
	if m.quitting || m.index >= len(m.tickets) {
		return finalView(m)
	}

	if m.mode == "ticketCreated" {
		domain := os.Getenv("JIRA_DOMAIN")
		ticketURL := m.tickets[m.index].ExistingJiraKey
		if domain != "" {
			ticketURL = fmt.Sprintf("https://%s/browse/%s", domain, m.tickets[m.index].ExistingJiraKey)
		}
		ticketStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Align(lipgloss.Left)

		return fmt.Sprintf(
			"\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s",
			ticketStyle.Render("Ticket created!"),
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render("Summary:"),
			m.tickets[m.index].Summary,
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Render("URL:"),
			ticketURL,
			"Press any key to continue...",
		)
	}

	if m.mode == "promptExisting" {
		return fmt.Sprintf(
			"Enter existing Jira ticket ID for test %q:\n\n%s\n\n(Press Enter to confirm, Esc to cancel)",
			m.tickets[m.index].TestName,
			m.inputValue,
		)
	}

	bodyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7")) // For general text
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	summaryStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	descHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	descBodyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	faintStyle := lipgloss.NewStyle().Faint(true)
	errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	existingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	lowPriorityStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))    // Green
	mediumPriorityStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // Blue/Purple
	highPriorityStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))  // Orange
	veryHighPriorityStyle := errorStyle                                         // Use red for highest

	t := m.tickets[m.index]

	var header string
	if t.Valid {
		header = headerStyle.Render(fmt.Sprintf("Proposed Ticket #%d of %d", m.index+1, len(m.tickets)))
	} else {
		header = headerStyle.Render(fmt.Sprintf("Ticket #%d of %d (Invalid)", m.index+1, len(m.tickets)))
	}

	// Assignee
	var assigneeLine string
	var assigneeDisplayValue string
	if t.AssigneeId != "" {
		assigneeDisplayValue = fmt.Sprintf("%s (%s)", t.AssigneeId, t.AssigneeId)
	} else {
		assigneeDisplayValue = t.AssigneeId
	}
	assigneeLine = summaryStyle.Render("Assignee:") + "\n" + bodyStyle.Render(assigneeDisplayValue)

	sum := summaryStyle.Render("Summary:")
	sumBody := descBodyStyle.Render(t.Summary)
	descHeader := descHeaderStyle.Render("\nDescription:")
	descBody := descBodyStyle.Render(t.Description)

	existingLine := ""
	if t.ExistingJiraKey != "" {
		prefix := "Existing ticket found"
		switch t.ExistingTicketSource {
		case "localdb":
			prefix = "Existing ticket found in local db"
		case "jira":
			prefix = "Existing ticket found in jira"
		}
		domain := os.Getenv("JIRA_DOMAIN")
		link := t.ExistingJiraKey
		if domain != "" {
			link = fmt.Sprintf("https://%s/browse/%s", domain, t.ExistingJiraKey)
		}
		existingLine = existingStyle.Render(fmt.Sprintf("\n%s: %s", prefix, link))
	}

	invalidLine := ""
	if !t.Valid {
		invalidLine = errorStyle.Render(fmt.Sprintf("\nCannot create ticket: %s", t.InvalidReason))
	}

	var priorityLine string
	if t.Priority != "" {
		var style lipgloss.Style
		switch t.Priority {
		case "Low":
			style = lowPriorityStyle
		case "Medium":
			style = mediumPriorityStyle
		case "High":
			style = highPriorityStyle
		case "Very High": // Match the name used in calculation
			style = veryHighPriorityStyle
		default:
			style = bodyStyle // Default style if unknown
		}
		priorityLine = summaryStyle.Render("Priority:") + "\n" + style.Render(t.Priority)
	}

	var helpLine string
	if !t.Valid {
		if t.ExistingJiraKey != "" {
			helpLine = faintStyle.Render("\n[n] to next, [e] to SET existing ticket id, [q] to quit.")
		} else {
			helpLine = faintStyle.Render("\n[n] to next, [e] to SET existing ticket id, [q] to quit.")
		}
	} else {
		if t.ExistingJiraKey != "" {
			helpLine = fmt.Sprintf("\n[n] to next, [e] to update ticket id, %s to remove ticket, [q] to quit.",
				redStyle.Render("[d]"))
		} else {
			dryRunLabel := ""
			if m.DryRun || m.JiraClient == nil {
				dryRunLabel = " (DRY RUN)"
			}
			helpLine = faintStyle.Render(fmt.Sprintf("\nPress [c] to create NEW ticket%s, [n] to skip, [e] to SET existing ticket id, [q] to quit.", dryRunLabel))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		assigneeLine,
		"",
		priorityLine,
		"",
		sum,
		sumBody,
		descHeader,
		descBody,
		existingLine,
		invalidLine,
		"",
		helpLine,
	)
}

func finalView(m tmodel) string {
	doneStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	return doneStyle.Render(fmt.Sprintf(
		"Done! Confirmed %d tickets, skipped %d. Press any key to exit...\n",
		m.confirmed, m.skipped,
	))
}

// -------------------------------------------------------------------------------------
// CSV Reading / Writing
// -------------------------------------------------------------------------------------

func readFlakyTestsCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	return r.ReadAll()
}
