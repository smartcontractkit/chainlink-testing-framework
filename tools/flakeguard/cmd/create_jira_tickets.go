package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/andygrunwald/go-jira"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	csvPathFlag     string
	dryRunFlag      bool
	jiraProjectKey  string
	jiraIssueType   string
	jiraSearchLabel string // If you want "flaky_test" as a flag or default
)

// CreateTicketsCmd is the Cobra command that runs a Bubble Tea TUI for CSV data
// and optionally creates tickets in Jira.
var CreateTicketsCmd = &cobra.Command{
	Use:   "create-tickets",
	Short: "Interactive TUI to confirm Jira tickets from CSV",
	Long: `Reads a CSV file describing flaky tests and displays each proposed
ticket in a text-based UI. Press 'y' to confirm, 'n' to skip, 'q' to quit.
If --dry-run=false we will attempt to create the ticket in Jira using environment
variables JIRA_DOMAIN, JIRA_EMAIL, JIRA_API_KEY. We also search for existing tickets
(label=flaky_test) with the same summary before prompting to create.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Validate flags
		if csvPathFlag == "" {
			log.Error().Msg("CSV path is required (use --csv-path)")
			os.Exit(1)
		}

		// Load project key from env if not given
		if jiraProjectKey == "" {
			jiraProjectKey = os.Getenv("JIRA_PROJECT_KEY")
		}
		if jiraProjectKey == "" {
			log.Error().Msg("Jira project key is required (set --jira-project or JIRA_PROJECT_KEY env)")
			os.Exit(1)
		}

		if jiraSearchLabel == "" {
			// default label to "flaky_test" if not provided
			jiraSearchLabel = "flaky_test"
		}

		// 2) Try to get a Jira client if possible (even in dry-run, for searching).
		//    If environment vars are missing, we won't do any searches or creation.
		client, clientErr := getJiraClient() // might fail if no env vars
		if clientErr != nil {
			log.Warn().Msgf("No valid Jira client: %v\nWill skip searching or creating tickets in Jira.", clientErr)
			client = nil
		}

		// 3) Parse CSV
		records, err := readCSV(csvPathFlag)
		if err != nil {
			log.Error().Err(err).Msg("Error reading CSV file")
			os.Exit(1)
		}

		//  Skip the first row (column headers).
		//  Make sure we still have data left after skipping.
		if len(records) <= 1 {
			log.Warn().Msg("No data rows found (CSV might be empty or only headers).")
			return nil
		}
		records = records[1:] // remove the header row

		// 4) Convert CSV rows to Tickets
		var tickets []Ticket
		for i, row := range records {
			if len(row) < 10 {
				log.Warn().Msgf("Skipping row %d (not enough columns)", i+1)
				continue
			}
			ticket := rowToTicket(row)
			tickets = append(tickets, ticket)
		}

		if len(tickets) == 0 {
			log.Warn().Msg("No valid tickets found (or CSV is empty).")
			return nil
		}

		// 5) For each valid ticket, optionally do a Jira search to see if a ticket already exists
		if client != nil {
			for i := range tickets {
				t := tickets[i]
				if t.Valid {
					key, err := findExistingTicket(client, jiraSearchLabel, t)
					if err != nil {
						log.Warn().Msgf("Search failed for %q: %v", t.Summary, err)
					} else if key != "" {
						t.ExistingJiraKey = key
					}
				}
			}
		}

		// 6) Prepare Bubble Tea model
		m := initialModel(tickets)
		m.DryRun = dryRunFlag
		m.JiraProject = jiraProjectKey
		m.JiraIssueType = jiraIssueType
		m.JiraClient = client

		// 7) Run Bubble Tea TUI
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			log.Error().Err(err).Msg("Error running Bubble Tea program")
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	CreateTicketsCmd.Flags().StringVar(&csvPathFlag, "csv-path", "", "Path to the CSV file containing ticket data")
	CreateTicketsCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "If true, do not actually create tickets in Jira")
	CreateTicketsCmd.Flags().StringVar(&jiraProjectKey, "jira-project", "", "Jira project key (or env JIRA_PROJECT_KEY)")
	CreateTicketsCmd.Flags().StringVar(&jiraIssueType, "jira-issue-type", "Task", "Type of Jira issue (Task, Bug, etc.)")
	CreateTicketsCmd.Flags().StringVar(&jiraSearchLabel, "jira-search-label", "", "Jira label to filter existing tickets (default: flaky_test)")
}

// -------------------------------------------------------------------------------------
// Ticket Data Model
// -------------------------------------------------------------------------------------

type Ticket struct {
	TestName        string
	Valid           bool
	InvalidReason   string
	Summary         string
	Description     string
	ExistingJiraKey string // if we found an existing ticket with this summary/label
}

// rowToTicket builds a Ticket from one CSV row (your columns).
// Required fields: Package(row[0]), TestName(row[2]), FlakeRate(row[7]), Logs(row[9]).
func rowToTicket(row []string) Ticket {
	pkg := row[0]
	testName := row[2]
	flakeRate := row[7]
	logs := row[9]

	// Check for missing required fields
	var missing []string
	if pkg == "" {
		missing = append(missing, "Package")
	}
	if testName == "" {
		missing = append(missing, "Test Name")
	}
	if flakeRate == "" {
		missing = append(missing, "Flake Rate")
	}
	if logs == "" {
		missing = append(missing, "Test Logs")
	}

	if len(missing) > 0 {
		reason := fmt.Sprintf("Missing required field(s): %s", strings.Join(missing, ", "))
		if testName != "" {
			// Append the test name in brackets for context
			reason += fmt.Sprintf(" [Test Name: %s]", testName)
		}
		return Ticket{
			Valid:         false,
			InvalidReason: reason,
		}
	}

	// Parse logs for multiple links
	var logLines []string
	logLinks := strings.Split(logs, ",")
	runNumber := 1
	for _, link := range logLinks {
		link = strings.TrimSpace(link)
		if link == "" {
			continue
		}
		// Format each log line as a Markdown link: "- [Run 1](http://...)"
		logLines = append(logLines,
			fmt.Sprintf("- [Run %d](%s)", runNumber, link),
		)
		runNumber++
	}
	if len(logLines) == 0 {
		return Ticket{
			Valid:         false,
			InvalidReason: "No valid test logs found after parsing",
		}
	}
	// Join them into one string with newlines
	testLogSection := strings.Join(logLines, "\n")

	// Summary: "Fix Flaky Test: <TestName> (<FlakeRate>% flake rate)"
	summary := fmt.Sprintf("Fix Flaky Test: %s (%s%% flake rate)", testName, flakeRate)

	// Build the description with Markdown headings and bullets
	description := fmt.Sprintf(`
## Test Details:
- **Package:** `+"`%s`"+`
- **Test Name:** `+"`%s`"+`
- **Flake Rate:** %s%% in the last 7 days

### Test Logs:
%s

### Action Items:
1. **Investigate Failed Test Logs:** Thoroughly review the provided logs to identify patterns or common error messages that indicate the root cause.
2. **Fix the Issue:** Analyze and address the underlying problem causing the flakiness.
3. **Rerun Tests Locally:** Execute the test and related changes on a local environment to ensure that the fix stabilizes the test, as well as all other tests that may be affected.
4. **Unskip the Test:** Once confirmed stable, remove any test skip markers to re-enable the test in the CI pipeline.
5. **Reference Guidelines:** Follow the recommendations in the [Flaky Test Guide].
`, pkg, testName, flakeRate, testLogSection)

	return Ticket{
		TestName:    testName,
		Valid:       true,
		Summary:     summary,
		Description: description,
	}
}

// -------------------------------------------------------------------------------------
// Jira Search
// -------------------------------------------------------------------------------------

// findExistingTicket looks for an existing ticket with the given summary, label, and project.
// Returns the Key of the first matching issue or "" if none found.
func findExistingTicket(client *jira.Client, label string, ticket Ticket) (string, error) {
	// Example JQL:
	//   project = MYPROJ AND labels = flaky_test AND summary ~ "Fix Flaky Test: MyTest"
	// We'll do an exact-ish match with double quotes around summary (escape quotes if needed).
	jql := fmt.Sprintf(`labels = "%s" AND summary ~ "%s" order by created DESC`,
		label, ticket.TestName)

	issues, resp, err := client.Issue.SearchWithContext(context.Background(), jql, &jira.SearchOptions{
		MaxResults: 1, // just need the first match
	})
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

type model struct {
	tickets   []Ticket
	index     int  // current ticket index
	confirmed int  // how many confirmed
	skipped   int  // how many skipped
	quitting  bool // whether we're exiting

	DryRun        bool
	JiraProject   string
	JiraIssueType string

	// We'll store the Jira client to create issues
	JiraClient *jira.Client
}

func initialModel(tickets []Ticket) model {
	return model{
		tickets:       tickets,
		index:         0,
		confirmed:     0,
		skipped:       0,
		quitting:      false,
		DryRun:        false,
		JiraProject:   "",
		JiraIssueType: "Task",
		JiraClient:    nil,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.quitting || m.index >= len(m.tickets) {
			return updateQuit(m)
		}
		t := m.tickets[m.index]
		if !t.Valid {
			// For invalid tickets, any key besides 'q' => skip
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				return updateQuit(m)
			default:
				return updateSkip(m)
			}
		} else {
			// If there's an existing ticket, user might still want to create a new one
			switch msg.String() {
			case "y":
				return updateConfirm(m)
			case "n":
				return updateSkip(m)
			case "q", "esc", "ctrl+c":
				return updateQuit(m)
			}
		}
	}
	return m, nil
}

func updateConfirm(m model) (tea.Model, tea.Cmd) {
	current := m.tickets[m.index]
	if !m.DryRun && m.JiraClient != nil {
		issueKey, err := createTicketInJira(m.JiraClient, current.Summary, current.Description, m.JiraProject, m.JiraIssueType)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to create Jira ticket for summary %q", current.Summary)
		} else {
			log.Info().Msgf("Created Jira issue: %s (summary=%s)", issueKey, current.Summary)
		}
	} else {
		log.Info().Msgf("[Dry Run] Would create Jira issue: %s", current.Summary)
	}

	m.confirmed++
	m.index++
	if m.index >= len(m.tickets) {
		m.quitting = true
	}
	return m, nil
}

func updateSkip(m model) (tea.Model, tea.Cmd) {
	m.skipped++
	m.index++
	if m.index >= len(m.tickets) {
		m.quitting = true
	}
	return m, nil
}

func updateQuit(m model) (tea.Model, tea.Cmd) {
	m.quitting = true
	return m, tea.Quit
}

func (m model) View() string {
	if m.quitting || m.index >= len(m.tickets) {
		return finalView(m)
	}

	t := m.tickets[m.index]

	// Some styling
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")) // pink/purple
	summaryStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")) // green
	descHeaderStyle := summaryStyle
	descBodyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")) // light gray
	helpStyle := lipgloss.NewStyle().
		Faint(true)
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9")) // bright red
	existingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")) // yellow

	if !t.Valid {
		header := headerStyle.Render(
			fmt.Sprintf("Ticket #%d of %d (Invalid)", m.index+1, len(m.tickets)),
		)
		errMsg := errorStyle.Render("Cannot create ticket: " + t.InvalidReason)
		hint := helpStyle.Render("\nPress any key to skip, or [q] to quit.\n")
		return fmt.Sprintf("%s\n\n%s\n%s\n", header, errMsg, hint)
	}

	header := headerStyle.Render(
		fmt.Sprintf("Proposed Ticket #%d of %d", m.index+1, len(m.tickets)),
	)
	sum := summaryStyle.Render("Summary:\n") + t.Summary
	descHeader := descHeaderStyle.Render("Description:\n")
	descBody := descBodyStyle.Render(t.Description)

	// If we found an existing ticket, show the user
	var existingLine string
	if t.ExistingJiraKey != "" {
		domain := os.Getenv("JIRA_DOMAIN")
		var link string
		if domain != "" {
			link = fmt.Sprintf("https://%s/browse/%s", domain, t.ExistingJiraKey)
		} else {
			link = t.ExistingJiraKey
		}
		existingLine = existingStyle.Render(
			fmt.Sprintf("\nAn existing ticket already exists: %s", link),
		)
	}

	dryRunLabel := ""
	if m.DryRun || m.JiraClient == nil {
		dryRunLabel = " (DRY RUN)"
	}
	help := helpStyle.Render(fmt.Sprintf("\nPress [y] to confirm%s, [n] to skip, [q] to quit.", dryRunLabel))

	return fmt.Sprintf("%s\n\n%s\n%s%s\n%s\n%s\n",
		header,
		sum,
		descHeader,
		descBody,
		existingLine,
		help,
	)
}

func finalView(m model) string {
	doneStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6")) // cyan
	return doneStyle.Render(fmt.Sprintf(
		"Done! Confirmed %d tickets, skipped %d. Exiting...\n",
		m.confirmed, m.skipped,
	))
}

// -------------------------------------------------------------------------------------
// Actual Jira Creation
// -------------------------------------------------------------------------------------

func createTicketInJira(client *jira.Client, summary, description, projectKey, issueType string) (string, error) {
	i := &jira.Issue{
		Fields: &jira.IssueFields{
			Project:     jira.Project{Key: projectKey},
			Summary:     summary,
			Description: description,
			Type:        jira.IssueType{Name: issueType},
			// Add the flaky_test label automatically if you want:
			Labels: []string{"flaky_test"},
		},
	}

	newIssue, resp, err := client.Issue.CreateWithContext(context.Background(), i)
	if err != nil {
		return "", fmt.Errorf("error creating Jira issue: %w (resp: %v)", err, resp)
	}
	return newIssue.Key, nil
}

// -------------------------------------------------------------------------------------
// Jira Client
// -------------------------------------------------------------------------------------

func getJiraClient() (*jira.Client, error) {
	domain := os.Getenv("JIRA_DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("JIRA_DOMAIN environment variable is not set")
	}

	email := os.Getenv("JIRA_EMAIL")
	if email == "" {
		return nil, fmt.Errorf("JIRA_EMAIL environment variable is not set")
	}

	apiKey := os.Getenv("JIRA_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("JIRA_API_KEY environment variable is not set")
	}

	tp := jira.BasicAuthTransport{
		Username: email,
		Password: apiKey,
	}
	return jira.NewClient(tp.Client(), fmt.Sprintf("https://%s", domain))
}

// readCSV is a helper for reading CSV
func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	return r.ReadAll()
}
