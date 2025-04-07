package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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

// Flags for create-tickets
var (
	csvPath               string
	dryRun                bool
	jiraProject           string
	jiraIssueType         string
	jiraSearchLabel       string // defaults to "flaky_test" if empty
	testDBPath            string
	skipExisting          bool
	userMappingPathCT     string
	userTestMappingPathCT string
)

var CreateTicketsCmd = &cobra.Command{
	Use:   "create-tickets",
	Short: "Interactive TUI to confirm and create Jira tickets from CSV",
	Long: `Reads a CSV describing flaky tests, attempts to assign owners based on patterns,
searches for existing tickets, and presents each for confirmation/action in a TUI.

Actions:
  [c] confirm & create ticket (if none exists)
  [n] skip this test
  [e] enter/edit existing Jira key for this test
  [d] delete associated Jira ticket (use with caution!)
  [q] quit

Features:
- Assigns owners based on 'user_test_mapping.json' if no owner is found otherwise.
- Searches Jira (label=flaky_test or --jira-search-label) for existing tickets matching test name.
- Uses local DB (--test-db-path) to remember existing ticket associations.
- Creates tickets in Jira (--dry-run=false) using JIRA_* env vars for auth.
- Includes Pillar Name (from 'user_mapping.json') when creating tickets if assignee is known.
- Outputs remaining unhandled tests to a new CSV file ('remaining_<original_name>.csv').
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Validate input & Set Defaults
		if csvPath == "" {
			return fmt.Errorf("CSV path is required (--csv-path)")
		}
		if jiraProject == "" {
			jiraProject = os.Getenv("JIRA_PROJECT_KEY")
		}
		if jiraProject == "" {
			return fmt.Errorf("Jira project key is required (set --jira-project or JIRA_PROJECT_KEY env)")
		}
		if jiraSearchLabel == "" {
			jiraSearchLabel = "flaky_test"
		}

		db, err := localdb.LoadDBWithPath(testDBPath)
		if err != nil {
			log.Warn().Err(err).Str("path", testDBPath).Msg("Failed to load local DB; continuing with empty DB.")
			db = localdb.NewDBWithPath(testDBPath)
		}

		userMap, err := mapping.LoadUserMappings(userMappingPathCT)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load user mappings. Pillar names won't be set.")
			userMap = make(map[string]mapping.UserMapping) // Ensure userMap is not nil
		}
		userTestMappings, err := mapping.LoadUserTestMappings(userTestMappingPathCT)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load user test mappings. Assignees won't be set automatically from patterns.")
			userTestMappings = []mapping.UserTestMappingWithRegex{} // Ensure not nil
		}

		records, err := readFlakyTestsCSV(csvPath)
		if err != nil {
			return fmt.Errorf("error reading CSV file '%s': %w", csvPath, err)
		}
		if len(records) <= 1 {
			log.Warn().Msg("CSV has no data rows.")
			return nil
		}
		originalRecords := records
		dataRows := records[1:]

		var tickets []model.FlakyTicket
		log.Info().Int("rows", len(dataRows)).Msg("Processing CSV data rows...")
		for i, row := range dataRows {
			if len(row) < 10 {
				log.Warn().Int("row_index", i+2).Int("columns", len(row)).Msg("Skipping row: not enough columns.")
				continue
			}

			ft := rowToFlakyTicket(row)
			ft.RowIndex = i + 2

			if !ft.Valid {
				log.Warn().Str("test", ft.TestName).Str("reason", ft.InvalidReason).Int("row", ft.RowIndex).Msg("Skipping invalid ticket from CSV")
				continue
			}

			// Check local DB first for existing ticket
			if entry, found := db.GetEntry(ft.TestPackage, ft.TestName); found {
				log.Debug().Str("test", ft.TestName).Str("db_ticket", entry.JiraTicket).Str("db_assignee", entry.AssigneeID).Time("db_skipped_at", entry.SkippedAt).Msg("Found existing entry in local DB")

				if entry.JiraTicket != "" {
					ft.ExistingJiraKey = entry.JiraTicket
					ft.ExistingTicketSource = "localdb"
				}
				// Always assign SkippedAt and AssigneeID from the DB entry if found
				ft.SkippedAt = entry.SkippedAt
				if entry.AssigneeID != "" {
					ft.AssigneeId = entry.AssigneeID
				}
			}

			// If assignee wasn't found in DB, try pattern matching
			if ft.AssigneeId == "" && len(userTestMappings) > 0 {
				testPath := ft.TestPackage
				assigneeID, matchErr := mapping.FindAssigneeIDForTest(testPath, userTestMappings)
				if matchErr != nil {
					log.Warn().Err(matchErr).Str("testPath", testPath).Msg("Error during assignee pattern matching")
				} else if assigneeID != "" {
					log.Debug().Str("test", testPath).Str("assignee", assigneeID).Msg("Assignee found via pattern matching")
					ft.AssigneeId = assigneeID
				}
			}

			// Check if the final AssigneeId has a mapping (for Pillar Name later)
			if ft.AssigneeId != "" {
				if _, exists := userMap[ft.AssigneeId]; !exists {
					ft.MissingUserMapping = true
					log.Warn().Str("test", ft.TestName).Str("assignee", ft.AssigneeId).Msg("Assignee ID is set, but not found in user_mapping.json")
				}
			}

			// Decide whether to skip processing based on flag and existing key
			if skipExisting && ft.ExistingJiraKey != "" {
				log.Info().Str("test", ft.TestName).Str("jira_key", ft.ExistingJiraKey).Msg("Skipping test due to --skip-existing flag.")
				continue
			}

			tickets = append(tickets, ft)
		}

		if len(tickets) == 0 {
			log.Warn().Msg("No new tickets to create found after filtering and validation.")
			return nil
		}

		client, clientErr := jirautils.GetJiraClient()
		if clientErr != nil {
			log.Warn().Msgf("No valid Jira client: %v\nWill skip searching or creating tickets in Jira.", clientErr)
			client = nil
		}

		if client != nil {
			processedCount := 0
			totalToSearch := 0
			for _, t := range tickets {
				if t.ExistingJiraKey == "" {
					totalToSearch++
				}
			}

			for i := range tickets {
				t := &tickets[i]
				if t.ExistingJiraKey == "" {
					key, searchErr := findExistingTicket(client, jiraSearchLabel, *t)
					processedCount++
					if searchErr != nil {
						log.Warn().Err(searchErr).Str("summary", t.Summary).Msg("Jira search failed for test")
					} else if key != "" {
						log.Info().Str("test", t.TestName).Str("found_key", key).Str("label", jiraSearchLabel).Msg("Found existing ticket in Jira via search")
						t.ExistingJiraKey = key
						t.ExistingTicketSource = "jira"
						errDb := db.UpsertEntry(t.TestPackage, t.TestName, key, t.SkippedAt, t.AssigneeId)
						if errDb != nil {
							log.Error().Err(errDb).Str("key", key).Msg("Failed to update local DB after finding ticket in Jira!")
						}
					}
				}
			}
		}

		m := initialCreateModel(tickets, userMap)
		m.DryRun = dryRun
		m.JiraProject = jiraProject
		m.JiraIssueType = jiraIssueType
		m.JiraClient = client
		m.originalRecords = originalRecords
		m.LocalDB = db

		finalModel, err := tea.NewProgram(m).Run()
		if err != nil {
			log.Error().Err(err).Msg("Error running Bubble Tea program")
		}

		fm, ok := finalModel.(createModel)
		if !ok {
			log.Error().Msg("TUI returned unexpected model type")
			if db != nil && !dryRun {
				if errDb := db.Save(); errDb != nil {
					log.Error().Err(errDb).Msg("Failed to save local DB after TUI error")
				}
			}
			return fmt.Errorf("TUI model error")
		}

		if !fm.DryRun {
			if db == nil {
				log.Error().Msg("Cannot save DB: DB instance is nil")
			} else if err := db.Save(); err != nil {
				log.Error().Err(err).Msg("Failed to save local DB")
			} else {
				fmt.Printf("Local DB updated: %s\n", db.FilePath())
			}
		} else {
			log.Info().Msg("Dry Run: Local DB changes were not saved.")
		}

		remainingFilePath := generateRemainingCSVPath(csvPath)
		err = writeRemainingCSV(remainingFilePath, fm.originalRecords, fm.tickets, fm.processedIndices)
		if err != nil {
			log.Error().Err(err).Str("path", remainingFilePath).Msg("Failed to write remaining tests CSV")
		} else {
			log.Info().Str("path", remainingFilePath).Msg("Remaining/skipped tests CSV generated.")
		}

		fmt.Printf("TUI Summary: %d confirmed, %d skipped/existing, %d total processed.\n", fm.confirmed, fm.skipped, fm.confirmed+fm.skipped)

		return nil
	},
}

func init() {
	CreateTicketsCmd.Flags().StringVar(&csvPath, "csv-path", "", "Path to CSV file with flaky tests (Required)")
	CreateTicketsCmd.Flags().BoolVar(&dryRun, "dry-run", false, "If true, do not create/modify Jira tickets or save DB changes")
	CreateTicketsCmd.Flags().StringVar(&jiraProject, "jira-project", "", "Jira project key (default: JIRA_PROJECT_KEY env)")
	CreateTicketsCmd.Flags().StringVar(&jiraIssueType, "jira-issue-type", "Task", "Jira issue type for new tickets")
	CreateTicketsCmd.Flags().StringVar(&jiraSearchLabel, "jira-search-label", "flaky_test", "Jira label to filter existing tickets during search") // Default added here
	CreateTicketsCmd.Flags().StringVar(&testDBPath, "test-db-path", localdb.DefaultDBPath(), "Path to the flaky test JSON database")
	CreateTicketsCmd.Flags().BoolVar(&skipExisting, "skip-existing", false, "Skip processing rows in TUI if a Jira ticket is already known (from DB or Jira search)")
	CreateTicketsCmd.Flags().StringVar(&userMappingPathCT, "user-mapping-path", "user_mapping.json", "Path to the user mapping JSON (JiraUserID -> PillarName)")
	CreateTicketsCmd.Flags().StringVar(&userTestMappingPathCT, "user-test-mapping-path", "user_test_mapping.json", "Path to the user test mapping JSON (Pattern -> JiraUserID)")
}

// -------------------------------------------------------------------------------------
// Helper Functions (Need implementations or adjustments)
// -------------------------------------------------------------------------------------

func readFlakyTestsCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open csv file '%s': %w", path, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv data from '%s': %w", path, err)
	}
	return records, nil
}

// rowToFlakyTicket: Converts a CSV row to a FlakyTicket model.
func rowToFlakyTicket(row []string) model.FlakyTicket {
	// Indices: pkg=0, testName=2, flakeRate=7, logs=9 (as per original code)
	const (
		pkgIndex       = 0
		nameIndex      = 2
		flakeRateIndex = 7
		logsIndex      = 9
	)

	if len(row) <= logsIndex {
		log.Error().Int("num_cols", len(row)).Msg("Row has fewer columns than expected in rowToFlakyTicket")
		return model.FlakyTicket{Valid: false, InvalidReason: "Incorrect number of columns"}
	}

	pkg := strings.TrimSpace(row[pkgIndex])
	testName := strings.TrimSpace(row[nameIndex])
	flakeRateStr := strings.TrimSpace(row[flakeRateIndex])
	logs := strings.TrimSpace(row[logsIndex])

	t := model.FlakyTicket{
		TestPackage: pkg,
		TestName:    testName,
		Valid:       true,
	}

	var flakeRate float64
	var parseErr error
	if flakeRateStr == "" || flakeRateStr == "%" {
		log.Warn().Str("test", testName).Str("package", pkg).Msg("Missing Flake Rate. Defaulting to Low priority.")
		t.Priority = "Low"
		flakeRateStr = "0"
	} else {
		// Remove '%' before parsing if necessary
		flakeRateValStr := strings.TrimSuffix(flakeRateStr, "%")
		flakeRate, parseErr = strconv.ParseFloat(flakeRateValStr, 64)
		if parseErr != nil {
			log.Error().Err(parseErr).Str("flake_rate", flakeRateStr).Str("test", testName).Msg("Invalid Flake Rate format.")
			t.Valid = false
			t.InvalidReason = fmt.Sprintf("Invalid Flake Rate: %s", flakeRateStr)
			t.Priority = "Low" // Default priority on error
		} else {
			t.FlakeRate = flakeRate
			switch {
			case flakeRate < 1.0:
				t.Priority = "Low"
			case flakeRate < 3.0:
				t.Priority = "Medium"
			case flakeRate < 5.0:
				t.Priority = "High"
			default:
				t.Priority = "Very High" // >= 5.0
			}
		}
	}
	displayFlakeRate := strings.TrimSuffix(flakeRateStr, "%")

	t.Summary = fmt.Sprintf("Fix Flaky Test: %s (%s%% flake rate)", testName, displayFlakeRate)

	var logSection strings.Builder
	if logs == "" {
		logSection.WriteString("(Logs not available in source data)")
	} else {
		runNumber := 1
		hasLinks := false
		for _, link := range strings.Split(logs, ",") {
			link = strings.TrimSpace(link)
			if link == "" {
				continue
			}
			if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
				logSection.WriteString(fmt.Sprintf("* [Run %d|%s]\n", runNumber, link))
				hasLinks = true
			} else {
				logSection.WriteString(fmt.Sprintf("* Run %d Log Info: %s\n", runNumber, link))
			}
			runNumber++
		}
		if !hasLinks && runNumber == 1 {
			logSection.Reset()
			logSection.WriteString("(No valid log links found in source data)")
		}
	}

	pkgDisplay := pkg
	if pkgDisplay == "" {
		pkgDisplay = "(Package Not Provided)"
	}
	testNameDisplay := testName
	if testNameDisplay == "" {
		testNameDisplay = "(Test Name Not Provided)"
	}

	t.Description = fmt.Sprintf(`h2. Flaky Test Details
* *Test Package:* %s
* *Test Name:* %s
* *Detected Flake Rate:* %s%% (in monitored period)
* *Priority:* %s

h3. Recent Failing Test Run Logs
%s

h3. Action Items
# *Investigate:* Review logs and test code to identify the root cause of the flakiness.
# *Fix:* Implement the necessary code changes or infrastructure adjustments.
# *Verify:* Run the test locally multiple times and monitor in CI to confirm stability.
# *Unskip:* Once confirmed stable, remove any test skip markers to re-enable the test in the CI pipeline.
# *Close Ticket:* Close this ticket once the test is confirmed stable.
# *Guidance:* Refer to the team's [Flaky Test Guide|https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/tools/flakeguard/e2e-flaky-test-guide.md].
`,
		pkgDisplay,
		testNameDisplay,
		displayFlakeRate,
		t.Priority,
		logSection.String(),
	)

	var missing []string
	if pkg == "" {
		missing = append(missing, "Package")
	}
	if testName == "" {
		missing = append(missing, "Test Name")
	}

	if t.Valid && len(missing) > 0 {
		t.Valid = false
		if t.InvalidReason != "" {
			t.InvalidReason += "; "
		}
		t.InvalidReason += fmt.Sprintf("Missing required fields: %s", strings.Join(missing, ", "))
	}

	return t
}

// findExistingTicket: Searches Jira for a ticket based on label and summary match.
func findExistingTicket(client *jira.Client, label string, ticket model.FlakyTicket) (string, error) {
	safeTestName := strings.ReplaceAll(ticket.TestName, `"`, `\"`)

	jql := fmt.Sprintf(`project = "%s" AND labels = "%s" AND summary ~ "\"%s\"" ORDER BY created DESC`, jiraProject, label, safeTestName)
	log.Debug().Str("jql", jql).Msg("Executing Jira search JQL")

	issues, resp, err := client.Issue.SearchWithContext(context.Background(), jql, &jira.SearchOptions{MaxResults: 5})
	if err != nil {
		errMsg := jirautils.ReadJiraErrorResponse(resp)
		log.Error().Err(err).Str("jql", jql).Str("response", errMsg).Msg("Error searching Jira")
		return "", fmt.Errorf("error searching Jira: %w (response: %s)", err, errMsg)
	}

	if len(issues) == 0 {
		log.Debug().Str("testName", ticket.TestName).Str("label", label).Msg("No existing tickets found in Jira via search.")
		return "", nil
	}

	log.Info().Str("testName", ticket.TestName).Str("foundKey", issues[0].Key).Msg("Found potentially matching ticket in Jira.")
	return issues[0].Key, nil
}

// generateRemainingCSVPath creates a path for the output CSV.
func generateRemainingCSVPath(originalPath string) string {
	dir := ""
	filename := originalPath
	lastSlash := strings.LastIndex(originalPath, "/")
	if lastSlash != -1 {
		dir = originalPath[:lastSlash+1]
		filename = originalPath[lastSlash+1:]
	}
	// Handle potential extension
	ext := ""
	lastDot := strings.LastIndex(filename, ".")
	if lastDot > 0 {
		ext = filename[lastDot:]
		filename = filename[:lastDot]
	}
	return fmt.Sprintf("%s%s_remaining%s", dir, filename, ext)
}

// writeRemainingCSV writes the tests that were not confirmed/processed successfully.
func writeRemainingCSV(path string, originalRecords [][]string, processedTickets []model.FlakyTicket, processedIndices map[int]bool) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create remaining CSV file '%s': %w", path, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Write header (first row of original records)
	if len(originalRecords) > 0 {
		if err := w.Write(originalRecords[0]); err != nil {
			return fmt.Errorf("failed to write header to remaining CSV '%s': %w", path, err)
		}
	}

	processedTicketMap := make(map[int]model.FlakyTicket)
	for _, t := range processedTickets {
		processedTicketMap[t.RowIndex] = t
	}

	writtenCount := 0
	for i, originalRow := range originalRecords {
		if i == 0 {
			continue
		}
		rowIndex := i + 1
		writeRow := true

		// Check if this row index was processed by the TUI
		if _, wasProcessed := processedIndices[rowIndex]; wasProcessed {
			if ticket, exists := processedTicketMap[rowIndex]; exists && ticket.Confirmed {
				log.Debug().Int("rowIndex", rowIndex).Str("test", ticket.TestName).Msg("Skipping confirmed ticket from remaining CSV output")
				writeRow = false // Do not write confirmed rows
			}
		} else {
			log.Debug().Int("rowIndex", rowIndex).Msg("Including unprocessed row in remaining CSV output")
		}

		if writeRow {
			if err := w.Write(originalRow); err != nil {
				log.Error().Err(err).Int("rowIndex", rowIndex).Msg("Failed to write row to remaining CSV")
			} else {
				writtenCount++
			}
		}
	}
	log.Info().Int("count", writtenCount).Str("path", path).Msg("Wrote remaining/unconfirmed tests to CSV.")

	return w.Error()
}

// -------------------------------------------------------------------------------------
// Bubble Tea Model for create-tickets (createModel)
// -------------------------------------------------------------------------------------

// createModel holds the state for the create-tickets TUI.
type createModel struct {
	tickets          []model.FlakyTicket // The filtered list of tickets to process in TUI
	index            int                 // Current ticket being viewed
	processedIndices map[int]bool        // Tracks RowIndex of tickets user interacted with (confirmed, skipped, edited, deleted)
	confirmed        int                 // Count of confirmed tickets
	skipped          int                 // Count of skipped/existing tickets
	quitting         bool
	DryRun           bool
	JiraProject      string
	JiraIssueType    string
	JiraClient       *jira.Client
	originalRecords  [][]string                     // Needed for writing remaining CSV
	LocalDB          *localdb.DB                    // Use pointer type consistent with RunE
	userMap          map[string]mapping.UserMapping // For Pillar lookup
	mode             string                         // "normal", "promptExisting", "ticketCreated", "confirmDelete"
	inputValue       string                         // For prompt mode
	infoMessage      string                         // Feedback to user
	errorMessage     string                         // Error feedback
}

func initialCreateModel(tickets []model.FlakyTicket, userMap map[string]mapping.UserMapping) createModel {
	idx := 0
	if len(tickets) == 0 {
		idx = -1
	}
	return createModel{
		tickets:          tickets,
		index:            idx,
		processedIndices: make(map[int]bool),
		mode:             "normal",
		userMap:          userMap,
	}
}

func (m createModel) Init() tea.Cmd {
	return nil
}

// Main Update routing based on mode
func (m createModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.index == -1 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}
		}
		return m, nil
	}

	// Always handle quit globally, regardless of mode
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	switch m.mode {
	case "promptExisting":
		return updatePromptExisting(m, msg)
	case "ticketCreated":
		return updateTicketCreated(m, msg)
	case "confirmDelete":
		return updateConfirmDelete(m, msg)
	case "normal":
		fallthrough
	default:
		return updateNormalMode(m, msg)
	}
}

// updateNormalMode handles keypresses in the default view.
func updateNormalMode(m createModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.index >= len(m.tickets) || m.index < 0 {
		log.Error().Int("index", m.index).Int("len", len(m.tickets)).Msg("Invalid index in updateNormalMode")
		m.quitting = true
		return m, tea.Quit
	}

	t := &m.tickets[m.index]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.infoMessage = ""
		m.errorMessage = ""

		switch msg.String() {
		case "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "e": // Enter/Edit Existing Ticket ID
			m.mode = "promptExisting"
			m.inputValue = t.ExistingJiraKey // Pre-fill
			return m, nil

		case "d": // Delete existing ticket (transition to confirm)
			if t.ExistingJiraKey != "" && m.JiraClient != nil && !m.DryRun {
				m.mode = "confirmDelete"
				m.errorMessage = fmt.Sprintf("Confirm delete %s? This cannot be undone.", jirautils.GetJiraLink(t.ExistingJiraKey))
			} else if t.ExistingJiraKey == "" {
				m.errorMessage = "No Jira ticket associated to delete."
			} else if m.DryRun {
				m.errorMessage = "Cannot delete tickets in Dry Run mode."
			} else {
				m.errorMessage = "Cannot delete tickets: Jira client unavailable."
			}
			return m, nil

		case "c": // Confirm/Create Ticket
			if !t.Valid {
				m.errorMessage = "Cannot create: Invalid data. Press [n] to skip."
				return m, nil
			}
			if t.ExistingJiraKey != "" {
				m.errorMessage = fmt.Sprintf("Cannot create: Ticket %s exists. Press [n] skip/[e] edit.", t.ExistingJiraKey)
				return m, nil
			}
			return updateConfirm(m)

		case "n": // Skip / Next
			return updateSkip(m)

		default: // Ignore other keys
			return m, nil
		}
	}
	return m, nil
}

// updatePromptExisting handles the mode where user enters a Jira key.
func updatePromptExisting(m createModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.index < 0 || m.index >= len(m.tickets) {
		log.Warn().Int("index", m.index).Int("len", len(m.tickets)).Msg("Invalid index in updatePromptExisting")
		m.errorMessage = "Internal error: invalid index"
		m.mode = "normal"
		return m, tea.ClearScreen
	}
	t := &m.tickets[m.index]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			newKey := strings.TrimSpace(m.inputValue)
			if newKey != "" && !regexp.MustCompile(`^[A-Z][A-Z0-9]+-\d+$`).MatchString(newKey) {
				m.errorMessage = "Invalid Jira key format (e.g., PROJ-123). Press Esc or Enter empty key to cancel."
				m.infoMessage = ""
				return m, nil
			}

			err := m.LocalDB.UpsertEntry(t.TestPackage, t.TestName, newKey, t.SkippedAt, t.AssigneeId)
			if err != nil {
				log.Error().Err(err).Str("test", t.TestName).Msg("Failed to update local DB entry from prompt")
				m.errorMessage = fmt.Sprintf("Error saving to local DB: %v", err)
				m.infoMessage = ""
				return m, nil
			}

			if newKey == "" {
				log.Debug().Str("test", t.TestName).Msg("Cleared Jira association for test.")
				t.ExistingJiraKey = ""
				t.ExistingTicketSource = ""
			} else {
				log.Debug().Str("test", t.TestName).Str("newKey", newKey).Msg("Manually set Jira key for test.")
				t.ExistingJiraKey = newKey
				t.ExistingTicketSource = "manual"
			}

			m.processedIndices[t.RowIndex] = true
			m.skipped++
			m.mode = "normal"
			m.inputValue = ""
			m.index++
			if m.index >= len(m.tickets) {
				m.quitting = true
			}

			m.infoMessage = ""
			m.errorMessage = ""

			var cmd tea.Cmd
			if m.quitting {
				cmd = tea.Quit
			} else {
				cmd = tea.ClearScreen
			}
			return m, cmd

		case tea.KeyEsc: // Cancel prompt
			m.mode = "normal"
			m.inputValue = ""
			m.errorMessage = ""
			m.infoMessage = "Edit cancelled."
			return m, tea.ClearScreen

		case tea.KeyBackspace:
			if len(m.inputValue) > 0 {
				m.inputValue = m.inputValue[:len(m.inputValue)-1]
			}
			return m, nil

		case tea.KeyRunes:
			if !strings.ContainsAny(string(msg.Runes), "\n\t\r") {
				m.inputValue += string(msg.Runes)
			}
			return m, nil
		}
	}

	return m, nil
}

// updateConfirm handles the 'c' key press: creates ticket and changes mode.
func updateConfirm(m createModel) (tea.Model, tea.Cmd) {
	i := m.index
	t := &m.tickets[i]

	if !t.Valid || t.ExistingJiraKey != "" {
		m.errorMessage = "Cannot create ticket (invalid or already exists)."
		m.mode = "normal"
		return m, nil
	}

	m.processedIndices[t.RowIndex] = true

	pillarName := ""
	assigneeForJira := ""
	if t.AssigneeId != "" {
		assigneeForJira = t.AssigneeId
		if userMapping, exists := m.userMap[t.AssigneeId]; exists {
			pillarName = userMapping.PillarName
		} else {
			log.Warn().Str("assignee", t.AssigneeId).Msg("Assignee ID present but no matching entry in user_mapping.json found.")
		}
	}

	// Create Jira ticket
	if !m.DryRun && m.JiraClient != nil {
		log.Info().Str("summary", t.Summary).Str("project", m.JiraProject).Str("assignee", assigneeForJira).Str("pillar", pillarName).Msg("Attempting to create Jira ticket...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		issueKey, err := jirautils.CreateTicketInJira(ctx, m.JiraClient, t.Summary, t.Description, m.JiraProject, m.JiraIssueType, assigneeForJira, t.Priority, []string{jiraSearchLabel}, pillarName)

		if err != nil {
			errMsg := fmt.Sprintf("Failed to create Jira ticket for %q", t.Summary)
			log.Error().Err(err).Msg(errMsg)
			m.errorMessage = fmt.Sprintf("%s: %v", errMsg, err)
			m.mode = "normal"
			m.infoMessage = ""
			return m, tea.ClearScreen
		}
		log.Info().Str("key", issueKey).Str("summary", t.Summary).Msg("Successfully created Jira ticket")
		t.Confirmed = true
		t.ExistingJiraKey = issueKey
		t.ExistingTicketSource = "jira-created"
		errDb := m.LocalDB.UpsertEntry(t.TestPackage, t.TestName, issueKey, t.SkippedAt, t.AssigneeId)
		if errDb != nil {
			log.Error().Err(errDb).Str("key", issueKey).Msg("Failed to update local DB after Jira creation!")
			m.errorMessage = "WARN: Jira ticket created but failed to update local DB!"
		}
		m.confirmed++
		m.mode = "ticketCreated"
		m.errorMessage = ""
		m.infoMessage = ""
		return m, tea.ClearScreen
	} else {
		// --- Dry Run or No Client Case ---
		t.Confirmed = true
		t.ExistingJiraKey = "DRYRUN-" + strconv.Itoa(1000+i)
		t.ExistingTicketSource = "dryrun-created"
		m.confirmed++
		m.mode = "ticketCreated"
		m.errorMessage = ""
		m.infoMessage = ""
		return m, tea.ClearScreen
	}
}

// updateTicketCreated handles the confirmation screen after a ticket is made.
func updateTicketCreated(m createModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			m.mode = "normal"
			m.index++
			if m.index >= len(m.tickets) {
				m.quitting = true
			}
			m.infoMessage = ""
			m.errorMessage = ""

			var cmd tea.Cmd
			if m.quitting {
				cmd = tea.Quit
			} else {
				cmd = tea.ClearScreen
			}
			return m, cmd
		}
	}
	return m, nil
}

// updateSkip handles the 'n' key press: increments skip count and advances.
func updateSkip(m createModel) (tea.Model, tea.Cmd) {
	if m.index < 0 || m.index >= len(m.tickets) {
		log.Warn().Int("index", m.index).Int("len", len(m.tickets)).Msg("Invalid index in updateSkip")
		m.errorMessage = "Internal error: invalid index"
		return m, nil
	}

	t := &m.tickets[m.index]
	m.processedIndices[t.RowIndex] = true

	m.skipped++
	m.index++
	m.errorMessage = ""
	m.infoMessage = ""

	if m.index >= len(m.tickets) {
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

// updateConfirmDelete handles the confirmation prompt for deleting a ticket.
func updateConfirmDelete(m createModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.index < 0 || m.index >= len(m.tickets) {
		m.mode = "normal"
		return m, tea.ClearScreen
	}
	t := &m.tickets[m.index] // Use pointer

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "y": // Yes, delete the ticket
			originalKey := t.ExistingJiraKey
			m.infoMessage = fmt.Sprintf("Attempting to delete Jira ticket %s...", originalKey)
			m.errorMessage = ""

			err := jirautils.DeleteTicketInJira(m.JiraClient, originalKey)

			if err != nil {
				errMsg := fmt.Sprintf("Failed to delete Jira ticket %s", originalKey)
				log.Error().Err(err).Str("key", originalKey).Msg(errMsg)
				m.errorMessage = fmt.Sprintf("%s: %v", errMsg, err)
				m.infoMessage = ""
			} else {
				log.Info().Str("key", originalKey).Msg("Successfully deleted Jira ticket")
				m.infoMessage = fmt.Sprintf("Deleted Jira ticket %s.", originalKey)
				t.ExistingJiraKey = ""
				t.ExistingTicketSource = ""
				errDb := m.LocalDB.UpsertEntry(t.TestPackage, t.TestName, "", t.SkippedAt, t.AssigneeId)
				if errDb != nil {
					log.Error().Err(errDb).Str("test", t.TestName).Msg("Failed to update local DB after deleting Jira key")
					m.errorMessage = fmt.Sprintf("WARN: Jira ticket deleted but failed to update local DB! Error: %v", errDb)
				}
				m.processedIndices[t.RowIndex] = true
			}
			m.mode = "normal"
			return m, tea.ClearScreen

		case "n", "esc":
			m.infoMessage = "Delete cancelled."
			m.errorMessage = ""
			m.mode = "normal"
			return m, tea.ClearScreen

		default:
			return m, nil
		}
	}
	return m, nil
}

// View logic for create-tickets TUI
func (m createModel) View() string {
	// Handle quitting state
	if m.quitting {
		return "Processing complete. Exiting...\n"
	}

	// Handle empty state
	if m.index == -1 || len(m.tickets) == 0 {
		return "No tickets to process.\n\n[q] quit\n"
	}

	// Handle potential out-of-bounds index defensively
	if m.index >= len(m.tickets) {
		return "Processing complete. Exiting...\n"
	}

	// Render based on current mode
	switch m.mode {
	case "ticketCreated":
		return viewTicketCreated(m)
	case "promptExisting":
		return viewPromptExisting(m)
	case "confirmDelete":
		return viewConfirmDelete(m)
	case "normal":
		fallthrough
	default:
		return viewNormal(m)
	}
}

// viewNormal renders the main ticket display.
func viewNormal(m createModel) string {
	t := m.tickets[m.index] // Get current ticket

	// --- Styles ---
	bodyStyle := lipgloss.NewStyle().PaddingBottom(1)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")).PaddingBottom(1) // Purple
	labelStyle := lipgloss.NewStyle().Bold(true).Width(15).Foreground(lipgloss.Color("39"))         // Blue
	valueStyle := lipgloss.NewStyle()
	summaryStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))       // Green (for summary label)
	descHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))    // Yellow (for desc label)
	descBodyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("248")).PaddingLeft(2) // Light grey, indented
	faintStyle := lipgloss.NewStyle().Faint(true)
	errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))   // Red
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))            // Orange
	existingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("40")).Bold(true) // Green Bold (for existing key)
	dryRunStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("208")).Background(lipgloss.Color("235")).Padding(0, 1)
	helpStyle := lipgloss.NewStyle().Faint(true).PaddingTop(1)
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")) // Cyan/Blueish (for info messages)

	// --- Build View ---
	var sb strings.Builder

	// Dry Run Banner
	if m.DryRun {
		sb.WriteString(dryRunStyle.Render("DRY RUN MODE") + "\n\n")
	}

	// Error Message Area (Top)
	if m.errorMessage != "" {
		sb.WriteString(errorStyle.Render("Error: "+m.errorMessage) + "\n\n")
	}
	// Info Message Area (Top, below error)
	if m.infoMessage != "" {
		sb.WriteString(infoStyle.Render(m.infoMessage) + "\n\n")
	}

	// Header
	validStatus := ""
	if !t.Valid {
		validStatus = errorStyle.Render(" (Invalid Data!)")
	}
	sb.WriteString(headerStyle.Render(fmt.Sprintf("Review Ticket [%d / %d]%s", m.index+1, len(m.tickets), validStatus)))
	sb.WriteString("\n")

	// Details Section
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Test Name:"), valueStyle.Render(t.TestName)))
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Package:"), valueStyle.Render(t.TestPackage)))

	// Priority
	priorityVal := t.Priority
	if priorityVal == "" {
		priorityVal = "(Not Set)"
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Priority:"), valueStyle.Render(priorityVal)))

	// Assignee & Pillar
	assigneeVal := "-"
	pillarVal := "-"
	if t.AssigneeId != "" {
		assigneeVal = t.AssigneeId
		if t.MissingUserMapping {
			assigneeVal += warningStyle.Render(" (Mapping Missing!)")
			pillarVal = warningStyle.Render("(Mapping Missing!)")
		} else if userMapEntry, exists := m.userMap[t.AssigneeId]; exists {
			pillarVal = userMapEntry.PillarName
			if pillarVal == "" {
				pillarVal = "(Pillar Not Set in Map)"
			}
		} else {
			pillarVal = errorStyle.Render("(Error: Map Lookup Failed)")
		}
	} else {
		assigneeVal = faintStyle.Render("(Not Assigned)")
		pillarVal = faintStyle.Render("(N/A - No Assignee)")
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Assigned To:"), valueStyle.Render(assigneeVal)))
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Pillar Name:"), valueStyle.Render(pillarVal)))

	// Existing Jira Key
	existingLine := faintStyle.Render("(None Found)")
	if t.ExistingJiraKey != "" {
		sourceInfo := ""
		if t.ExistingTicketSource != "" {
			sourceInfo = fmt.Sprintf(" (from %s)", t.ExistingTicketSource)
		}
		link := jirautils.GetJiraLink(t.ExistingJiraKey)
		existingLine = existingStyle.Render(fmt.Sprintf("%s%s", link, sourceInfo))
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Existing Jira:"), existingLine))

	// Summary & Description (Collapsible or truncated?)
	sb.WriteString("\n")
	sb.WriteString(summaryStyle.Render("Proposed Summary:") + "\n")
	sb.WriteString(valueStyle.Render(t.Summary) + "\n\n")

	sb.WriteString(descHeaderStyle.Render("Proposed Description:") + "\n")
	sb.WriteString(descBodyStyle.Render(t.Description) + "\n")

	// Help Line / Actions
	sb.WriteString(helpStyle.Render(buildHelpLine(m)))

	return bodyStyle.Render(sb.String())
}

// viewTicketCreated renders the confirmation screen after creating a ticket.
func viewTicketCreated(m createModel) string {
	if m.index < 0 || m.index >= len(m.tickets) {
		return "Error: Invalid index for ticket confirmation.\n"
	}
	t := m.tickets[m.index]
	ticketStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)   // Green Bold
	urlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Underline(true) // Blue Underline
	helpStyle := lipgloss.NewStyle().Faint(true).PaddingTop(1)
	labelStyle := lipgloss.NewStyle().Width(10)

	var sb strings.Builder

	sb.WriteString(ticketStyle.Render("Ticket processed successfully!"))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("%s%s\n", labelStyle.Render("Summary:"), t.Summary))
	sb.WriteString(fmt.Sprintf("%s%s\n", labelStyle.Render("Jira Key:"), urlStyle.Render(jirautils.GetJiraLink(t.ExistingJiraKey))))
	if t.AssigneeId != "" {
		sb.WriteString(fmt.Sprintf("%s%s\n", labelStyle.Render("Assignee:"), t.AssigneeId))
	}
	if pillar, ok := m.userMap[t.AssigneeId]; ok && pillar.PillarName != "" {
		sb.WriteString(fmt.Sprintf("%s%s\n", labelStyle.Render("Pillar:"), pillar.PillarName))
	} else if t.AssigneeId != "" {
		sb.WriteString(fmt.Sprintf("%s%s\n", labelStyle.Render("Pillar:"), "(Not Set)"))
	}

	sb.WriteString(helpStyle.Render("\nPress any key to continue to the next test, or [q] to quit."))

	return sb.String()
}

// viewPromptExisting renders the input prompt for the Jira key.
func viewPromptExisting(m createModel) string {
	if m.index < 0 || m.index >= len(m.tickets) {
		return "Error: Invalid index for prompt.\n"
	}
	t := m.tickets[m.index]
	promptStyle := lipgloss.NewStyle().Bold(true)
	inputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82")).BorderStyle(lipgloss.NormalBorder()).Padding(0, 1).Width(40) // Added width
	helpStyle := lipgloss.NewStyle().Faint(true).PaddingTop(1)
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).PaddingBottom(1) // Red

	var sb strings.Builder

	if m.errorMessage != "" {
		sb.WriteString(errorStyle.Render("Error: "+m.errorMessage) + "\n")
	}

	sb.WriteString(promptStyle.Render(fmt.Sprintf("Enter existing Jira Key for: %s", t.TestName)))
	sb.WriteString("\n\n")
	sb.WriteString(inputStyle.Render(m.inputValue + "_")) // Show cursor
	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("(Press Enter to confirm, Esc to cancel)"))

	return sb.String()
}

// viewConfirmDelete renders the 'Are you sure?' prompt for deletion.
func viewConfirmDelete(m createModel) string {
	if m.index < 0 || m.index >= len(m.tickets) {
		return "Error: Invalid index for delete confirmation.\n"
	}
	t := m.tickets[m.index]
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))        // Red Bold
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).PaddingBottom(1) // Orange warning
	helpStyle := lipgloss.NewStyle().Faint(true).PaddingTop(1)

	var sb strings.Builder

	sb.WriteString(promptStyle.Render(fmt.Sprintf("Permanently delete Jira ticket %s associated with test %s?", t.ExistingJiraKey, t.TestName)))
	sb.WriteString("\n\n")
	if m.errorMessage != "" {
		sb.WriteString(warningStyle.Render(m.errorMessage) + "\n")
	}

	sb.WriteString(helpStyle.Render("[y] Yes, delete ticket  |  [n/esc] No, cancel"))

	return sb.String()
}

// buildHelpLine generates the dynamic help text based on the ticket state.
func buildHelpLine(m createModel) string {
	if m.index < 0 || m.index >= len(m.tickets) {
		return "[q] Quit"
	}
	t := m.tickets[m.index]
	var actions []string

	if t.Valid && t.ExistingJiraKey == "" {
		createLabel := "[c] Create"
		if m.DryRun || m.JiraClient == nil {
			createLabel += " (DryRun/Offline)"
		}
		actions = append(actions, createLabel)
	}

	actions = append(actions, "[n] Skip")
	actions = append(actions, "[e] Edit/Set Key")

	if t.ExistingJiraKey != "" && !m.DryRun && m.JiraClient != nil {
		actions = append(actions, lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("[d] Delete Key")) // Red 'd'
	}

	actions = append(actions, "[q] Quit")

	return "Actions: " + strings.Join(actions, " | ")
}
