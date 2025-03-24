package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/briandowns/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/jirautils"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/localdb"
	"github.com/spf13/cobra"
)

var (
	csvPath             string
	dryRun              bool
	jiraProject         string
	jiraIssueType       string
	jiraSearchLabel     string // defaults to "flaky_test" if empty
	flakyTestJSONDBPath string
	assigneeMappingPath string
	skipExisting        bool
)

// AssigneeMapping holds a regex pattern and its corresponding assignee.
type AssigneeMapping struct {
	Pattern  string `json:"pattern"`
	Assignee string `json:"assignee"`
}

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
- Optionally, an assignee mapping file (JSON) can be provided to set the ticket’s assignee
  based on the test package. The mapping supports regex.
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
		db, err := localdb.LoadDBWithPath(flakyTestJSONDBPath)
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
		var tickets []FlakyTicket
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

		// Load assignee mapping (if provided)
		var mappings []AssigneeMapping
		if assigneeMappingPath != "" {
			mappingData, err := os.ReadFile(assigneeMappingPath)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to read assignee mapping file; proceeding without assignee mapping.")
			} else {
				if err := json.Unmarshal(mappingData, &mappings); err != nil {
					log.Warn().Err(err).Msg("Failed to unmarshal assignee mapping; proceeding without assignee mapping.")
				} else {
					// Apply mapping: iterate over tickets and assign based on regex match.
					for i := range tickets {
						for _, mapping := range mappings {
							re, err := regexp.Compile(mapping.Pattern)
							if err != nil {
								log.Warn().Msgf("Invalid regex pattern %q: %v", mapping.Pattern, err)
								continue
							}
							if re.MatchString(tickets[i].TestPackage) {
								tickets[i].Assignee = mapping.Assignee
								break // use first matching mapping
							}
						}
					}
				}
			}
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
		fm := finalModel.(model)

		// 9) Save local DB with any new knowledge
		if err := fm.LocalDB.Save(); err != nil {
			log.Error().Err(err).Msg("Failed to save local DB")
		} else {
			fmt.Printf("Local DB has been updated at: %s\n", fm.LocalDB.FilePath())
		}

		// 10) Write remaining CSV
		remainingCSVPath := makeRemainingCSVPath(csvPath)
		if err := writeRemainingTicketsCSV(remainingCSVPath, fm); err != nil {
			log.Error().Err(err).Msgf("Failed to write updated CSV to %s", remainingCSVPath)
		} else {
			fmt.Printf("Remaining tickets have been written to: %s\n", remainingCSVPath)
		}

		return nil
	},
}

func init() {
	CreateTicketsCmd.Flags().StringVar(&csvPath, "csv-path", "", "Path to the CSV file containing flaky test data")
	CreateTicketsCmd.Flags().BoolVar(&dryRun, "dry-run", false, "If true, do not actually create tickets in Jira")
	CreateTicketsCmd.Flags().StringVar(&jiraProject, "jira-project", "", "Jira project key (or env JIRA_PROJECT_KEY)")
	CreateTicketsCmd.Flags().StringVar(&jiraIssueType, "jira-issue-type", "Task", "Type of Jira issue (Task, Bug, etc.)")
	CreateTicketsCmd.Flags().StringVar(&jiraSearchLabel, "jira-search-label", "", "Jira label to filter existing tickets (default: flaky_test)")
	CreateTicketsCmd.Flags().StringVar(&flakyTestJSONDBPath, "flaky-test-json-db-path", "", "Path to the flaky test JSON database (default: ~/.flaky_tes_db.json)")
	CreateTicketsCmd.Flags().StringVar(&assigneeMappingPath, "assignee-mapping", "", "Path to JSON file with assignee mapping (supports regex)")
	CreateTicketsCmd.Flags().BoolVar(&skipExisting, "skip-existing", false, "Skip processing tickets that already have a Jira ticket ID")
}

// -------------------------------------------------------------------------------------
// FlakyTicket Data Model
// -------------------------------------------------------------------------------------

type FlakyTicket struct {
	RowIndex             int
	Confirmed            bool
	Valid                bool
	InvalidReason        string
	TestName             string
	TestPackage          string
	Summary              string
	Description          string
	ExistingJiraKey      string
	ExistingTicketSource string // "localdb" or "jira" (if found)
	Assignee             string
}

// rowToFlakyTicket: build a ticket from one CSV row (index assumptions: pkg=0, testName=2, flakeRate=7, logs=9).
func rowToFlakyTicket(row []string) FlakyTicket {
	pkg := strings.TrimSpace(row[0])
	testName := strings.TrimSpace(row[2])
	flakeRate := strings.TrimSpace(row[7])
	logs := strings.TrimSpace(row[9])

	summary := fmt.Sprintf("Fix Flaky Test: %s (%s%% flake rate)", testName, flakeRate)

	// Parse logs
	var logSection string
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
			// Build a Jira wiki bullet for each log link
			lines = append(lines, fmt.Sprintf("* [Run %d|%s]", runNumber, link))
			runNumber++
		}
		if len(lines) == 0 {
			logSection = "(Logs not available)"
		} else {
			logSection = strings.Join(lines, "\n")
		}
	}

	// Use Jira Wiki Markup rather than Markdown
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
		flakeRate,
		logSection,
	)

	t := FlakyTicket{
		TestPackage: pkg,
		TestName:    testName,
		Summary:     summary,
		Description: desc,
		Valid:       true,
	}

	// check required fields
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
		missing = append(missing, "Logs")
	}

	if len(missing) > 0 {
		t.Valid = false
		t.InvalidReason = fmt.Sprintf("Missing required: %s", strings.Join(missing, ", "))
	}
	return t
}

// -------------------------------------------------------------------------------------
// Jira Search
// -------------------------------------------------------------------------------------

func findExistingTicket(client *jira.Client, label string, ticket FlakyTicket) (string, error) {
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

type model struct {
	tickets         []FlakyTicket
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

func initialModel(tickets []FlakyTicket) model {
	return model{
		tickets: tickets,
		index:   0,
		mode:    "normal",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func updateNormalMode(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

func updatePromptExisting(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		m.LocalDB.Set(t.TestPackage, t.TestName, t.ExistingJiraKey)
		m.mode = "normal"
		m.inputValue = ""
		return updateSkip(m)
	case tea.KeyEsc:
		m.mode = "normal"
		m.inputValue = ""
	}
	return m, nil
}

func updateConfirm(m model) (tea.Model, tea.Cmd) {
	i := m.index
	t := m.tickets[i]

	// Attempt Jira creation if not dry-run and we have a client.
	// Pass the assignee (if any) to the CreateTicketInJira function.
	if !m.DryRun && m.JiraClient != nil {
		issueKey, err := jirautils.CreateTicketInJira(m.JiraClient, t.Summary, t.Description, m.JiraProject, m.JiraIssueType, t.Assignee)
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
		log.Info().Msgf("[Dry Run] Would create Jira issue: %q", t.Summary)
		t.Confirmed = true
	}
	m.tickets[i] = t
	m.confirmed++
	m.mode = "ticketCreated"
	return m, nil
}

func updateTicketCreated(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "n":
		m.mode = "normal"
		m.index++
		if m.index >= len(m.tickets) {
			m.quitting = true
		}
		return m, nil
	case "e":
		m.mode = "promptExisting"
		m.inputValue = ""
		return m, nil
	case "q", "esc", "ctrl+c":
		return updateQuit(m)
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

// View logic
func (m model) View() string {
	if m.quitting || m.index >= len(m.tickets) {
		return finalView(m)
	}

	if m.mode == "promptExisting" {
		return fmt.Sprintf(
			"Enter existing Jira ticket ID for test %q:\n\n%s\n\n(Press Enter to confirm, Esc to cancel)",
			m.tickets[m.index].TestName,
			m.inputValue,
		)
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	summaryStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	descHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	descBodyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	faintStyle := lipgloss.NewStyle().Faint(true)
	errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	existingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	t := m.tickets[m.index]

	var header string
	if t.Valid {
		header = headerStyle.Render(fmt.Sprintf("Proposed Ticket #%d of %d", m.index+1, len(m.tickets)))
	} else {
		header = headerStyle.Render(fmt.Sprintf("Ticket #%d of %d (Invalid)", m.index+1, len(m.tickets)))
	}

	// New: Assignee line above Summary
	var assigneeLine string
	if t.Assignee != "" {
		assigneeLine = summaryStyle.Render(fmt.Sprintf("Assignee: %s", t.Assignee))
	}

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

	var helpLine string
	if !t.Valid {
		if t.ExistingJiraKey != "" {
			helpLine = faintStyle.Render("\n[n] to next, [e] to update ticket id, [q] to quit.")
		} else {
			helpLine = faintStyle.Render("\n[n] to next, [e] to add existing ticket ID, [q] to quit.")
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
			helpLine = faintStyle.Render(fmt.Sprintf("\nPress [c] to create ticket%s, [n] to skip, [e] to enter existing ticket, [q] to quit.", dryRunLabel))
		}
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n%s\n%s%s%s\n%s\n",
		header,
		assigneeLine,
		sum,
		sumBody,
		descHeader,
		descBody,
		existingLine,
		invalidLine,
		helpLine,
	)
}

func finalView(m model) string {
	doneStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	return doneStyle.Render(fmt.Sprintf(
		"Done! Confirmed %d tickets, skipped %d. Exiting...\n",
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

func writeRemainingTicketsCSV(newPath string, m model) error {
	confirmedRows := make(map[int]bool)
	for _, t := range m.tickets {
		if t.Confirmed || t.ExistingJiraKey != "" {
			confirmedRows[t.RowIndex] = true
		}
	}
	var newRecords [][]string
	orig := m.originalRecords
	if len(orig) > 0 {
		newRecords = append(newRecords, orig[0])
	}
	for i := 1; i < len(orig); i++ {
		if !confirmedRows[i] {
			newRecords = append(newRecords, orig[i])
		}
	}
	f, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	if err := w.WriteAll(newRecords); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

func makeRemainingCSVPath(originalPath string) string {
	ext := filepath.Ext(originalPath)
	base := strings.TrimSuffix(filepath.Base(originalPath), ext)
	dir := filepath.Dir(originalPath)
	newName := base + ".remaining" + ext
	return filepath.Join(dir, newName)
}
