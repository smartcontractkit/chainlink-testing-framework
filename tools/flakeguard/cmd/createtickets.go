package cmd

import (
	"context"
	"encoding/csv" // Keep for now, might be needed by TUI model init internally
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/briandowns/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"

	// Import mapping and other utils
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
	testDBPath            string // Now consistent with tickets cmd
	skipExisting          bool
	userMappingPathCT     string // Use CT suffix to avoid clash if both cmds are in same main
	userTestMappingPathCT string // Use CT suffix
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
			jiraSearchLabel = "flaky_test" // Default label
		}

		// 2) Load local DB
		db, err := localdb.LoadDBWithPath(testDBPath)
		if err != nil {
			// Non-fatal: Continue with an empty DB if loading fails
			log.Warn().Err(err).Str("path", testDBPath).Msg("Failed to load local DB; continuing with empty DB.")
			db = localdb.NewDBWithPath(testDBPath) // Create new DB instance with the target path
		}

		// 3) Load Mappings using the new package
		userMap, err := mapping.LoadUserMappings(userMappingPathCT)
		if err != nil {
			// Non-fatal: Log error and continue, Pillar names won't be set automatically
			log.Error().Err(err).Msg("Failed to load user mappings. Pillar names won't be set.")
			userMap = make(map[string]mapping.UserMapping) // Ensure userMap is not nil
		}
		userTestMappings, err := mapping.LoadUserTestMappings(userTestMappingPathCT)
		if err != nil {
			// Non-fatal: Log error and continue, auto-assignment won't work
			log.Error().Err(err).Msg("Failed to load user test mappings. Assignees won't be set automatically from patterns.")
			userTestMappings = []mapping.UserTestMappingWithRegex{} // Ensure not nil
		}

		// 4) Read CSV
		records, err := readFlakyTestsCSV(csvPath) // Assumes readFlakyTestsCSV exists
		if err != nil {
			return fmt.Errorf("error reading CSV file '%s': %w", csvPath, err)
		}
		if len(records) <= 1 { // Check for header only or empty
			log.Warn().Msg("CSV has no data rows.")
			return nil
		}
		originalRecords := records // Keep header + data for final output
		dataRows := records[1:]    // Process only data rows

		// 5) Convert CSV rows -> FlakyTicket objects & Apply Mappings/DB lookups
		var tickets []model.FlakyTicket
		log.Info().Int("rows", len(dataRows)).Msg("Processing CSV data rows...")
		for i, row := range dataRows {
			// Basic validation of row structure could happen here
			if len(row) < 10 { // Ensure enough columns for rowToFlakyTicket
				log.Warn().Int("row_index", i+2).Int("columns", len(row)).Msg("Skipping row: not enough columns.")
				continue
			}

			// Convert row data to a ticket structure
			ft := rowToFlakyTicket(row) // Assumes rowToFlakyTicket is adapted/exists
			ft.RowIndex = i + 2         // Store original CSV row number (1-based index + header)

			if !ft.Valid {
				log.Warn().Str("test", ft.TestName).Str("reason", ft.InvalidReason).Int("row", ft.RowIndex).Msg("Skipping invalid ticket from CSV")
				// We might still want to show invalid tickets in TUI for fixing?
				// For now, let's add them but they won't be creatable.
				// tickets = append(tickets, ft) // Option 1: Include invalid ones
				continue // Option 2: Skip entirely
			}

			// Check local DB first for existing ticket
			if entry, found := db.GetEntry(ft.TestPackage, ft.TestName); found {
				// Entry found in the DB, log it and update the FlakyTicket (ft)
				log.Debug().Str("test", ft.TestName).Str("db_ticket", entry.JiraTicket).Str("db_assignee", entry.AssigneeID).Time("db_skipped_at", entry.SkippedAt).Msg("Found existing entry in local DB")

				// Assign details from the DB entry to the current FlakyTicket (ft)
				if entry.JiraTicket != "" {
					ft.ExistingJiraKey = entry.JiraTicket
					ft.ExistingTicketSource = "localdb" // Mark the source
				}
				// Always assign SkippedAt and AssigneeID from the DB entry if found
				ft.SkippedAt = entry.SkippedAt
				if entry.AssigneeID != "" {
					ft.AssigneeId = entry.AssigneeID
				}
			}

			// If assignee wasn't found in DB, try pattern matching
			if ft.AssigneeId == "" && len(userTestMappings) > 0 {
				testPath := ft.TestPackage // Or combine package + name? Use package based on example patterns.
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
				continue // Skip adding this ticket to the list for the TUI
			}

			tickets = append(tickets, ft)
		} // End CSV processing loop

		if len(tickets) == 0 {
			log.Warn().Msg("No processable tickets found after filtering and validation.")
			return nil
		}

		// 6) Attempt Jira client creation
		client, clientErr := jirautils.GetJiraClient()
		if clientErr != nil {
			log.Warn().Msgf("No valid Jira client: %v\nWill skip searching or creating tickets in Jira.", clientErr)
			client = nil // Ensure client is nil
		}

		// 7) If Jira client exists, search for existing tickets online (for those without a key yet)
		if client != nil {
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Suffix = " Searching Jira for potentially existing tickets..."
			s.Start()
			processedCount := 0
			totalToSearch := 0
			for _, t := range tickets {
				if t.ExistingJiraKey == "" {
					totalToSearch++
				}
			}

			for i := range tickets {
				t := &tickets[i]             // Use pointer to modify slice element
				if t.ExistingJiraKey == "" { // Only search if we don't have a key from DB
					key, searchErr := findExistingTicket(client, jiraSearchLabel, *t) // findExistingTicket needs definition
					processedCount++
					s.Suffix = fmt.Sprintf(" Searching Jira... (%d/%d)", processedCount, totalToSearch)
					if searchErr != nil {
						// Log non-fatal search error
						log.Warn().Err(searchErr).Str("summary", t.Summary).Msg("Jira search failed for test")
					} else if key != "" {
						log.Info().Str("test", t.TestName).Str("found_key", key).Str("label", jiraSearchLabel).Msg("Found existing ticket in Jira via search")
						t.ExistingJiraKey = key
						t.ExistingTicketSource = "jira"
						// Update local DB immediately with this finding
						db.UpsertEntry(t.TestPackage, t.TestName, key, t.SkippedAt, t.AssigneeId) // Persist found key
					}
				}
			}
			s.Stop()
			fmt.Println() // Newline after spinner
		}

		// 8) Create Bubble Tea model
		m := initialCreateModel(tickets, userMap) // Pass userMap for Pillar lookup
		m.DryRun = dryRun
		m.JiraProject = jiraProject
		m.JiraIssueType = jiraIssueType
		m.JiraClient = client
		m.originalRecords = originalRecords // Pass header + all data rows
		m.LocalDB = db                      // Pass DB for TUI actions to update it

		// 9) Run TUI
		finalModel, err := tea.NewProgram(m).Run()
		if err != nil {
			// Log error, but proceed to save DB and write remaining CSV
			log.Error().Err(err).Msg("Error running Bubble Tea program")
		}

		// 10) Process results after TUI exits
		fm, ok := finalModel.(createModel) // Type assertion
		if !ok {
			log.Error().Msg("TUI returned unexpected model type")
			// Still attempt to save DB
			if errDb := db.Save(); errDb != nil {
				log.Error().Err(errDb).Msg("Failed to save local DB after TUI error")
			}
			return fmt.Errorf("TUI model error")
		}

		// 11) Save local DB with any changes made during TUI
		if !fm.DryRun { // Only save if not in dry run
			if err := fm.LocalDB.Save(); err != nil {
				log.Error().Err(err).Msg("Failed to save local DB")
				// Report error but don't necessarily exit
			} else {
				fmt.Printf("Local DB updated: %s\n", fm.LocalDB.FilePath())
			}
		} else {
			log.Info().Msg("Dry Run: Local DB changes were not saved.")
		}

		// Report summary TUI stats
		fmt.Printf("TUI Summary: %d confirmed, %d skipped/existing, %d total processed.\n", fm.confirmed, fm.skipped, fm.confirmed+fm.skipped)

		return nil // Cobra handles printing the error returned by RunE
	},
}

func init() {
	CreateTicketsCmd.Flags().StringVar(&csvPath, "csv-path", "", "Path to CSV file with flaky tests (Required)")
	CreateTicketsCmd.Flags().BoolVar(&dryRun, "dry-run", false, "If true, do not create/modify Jira tickets or save DB changes")
	CreateTicketsCmd.Flags().StringVar(&jiraProject, "jira-project", "", "Jira project key (default: JIRA_PROJECT_KEY env)")
	CreateTicketsCmd.Flags().StringVar(&jiraIssueType, "jira-issue-type", "Task", "Jira issue type for new tickets")
	CreateTicketsCmd.Flags().StringVar(&jiraSearchLabel, "jira-search-label", "flaky_test", "Jira label to filter existing tickets during search") // Default added here
	// Use consistent flag name and default path function
	CreateTicketsCmd.Flags().StringVar(&testDBPath, "test-db-path", localdb.DefaultDBPath(), "Path to the flaky test JSON database")
	CreateTicketsCmd.Flags().BoolVar(&skipExisting, "skip-existing", false, "Skip processing rows in TUI if a Jira ticket is already known (from DB or Jira search)")
	// Add mapping paths using suffix convention
	CreateTicketsCmd.Flags().StringVar(&userMappingPathCT, "user-mapping-path", "user_mapping.json", "Path to the user mapping JSON (JiraUserID -> PillarName)")
	CreateTicketsCmd.Flags().StringVar(&userTestMappingPathCT, "user-test-mapping-path", "user_test_mapping.json", "Path to the user test mapping JSON (Pattern -> JiraUserID)")
}

// -------------------------------------------------------------------------------------
// Helper Functions (Need implementations or adjustments)
// -------------------------------------------------------------------------------------

// readFlakyTestsCSV: Reads the specified CSV file.
// (Implementation assumed to exist as before)
func readFlakyTestsCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open csv file '%s': %w", path, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	// Configure reader if needed (e.g., r.Comma = ';')
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv data from '%s': %w", path, err)
	}
	return records, nil
}

// rowToFlakyTicket: Converts a CSV row to a FlakyTicket model.
// (Implementation assumed to exist, ensure it sets TestPackage, TestName, parses FlakeRate for Priority, generates Summary/Description)
// IMPORTANT: Ensure this function is robust and handles potential errors (e.g., parsing flake rate).
func rowToFlakyTicket(row []string) model.FlakyTicket {
	// Simplified example - adapt based on your actual CSV structure and needs
	// Indices: pkg=0, testName=2, flakeRate=7, logs=9 (as per original code)
	const (
		pkgIndex       = 0
		nameIndex      = 2
		flakeRateIndex = 7
		logsIndex      = 9
	)

	// Basic bounds check (already done partially in RunE)
	if len(row) <= logsIndex {
		log.Error().Int("num_cols", len(row)).Msg("Row has fewer columns than expected in rowToFlakyTicket")
		// Return an invalid ticket
		return model.FlakyTicket{Valid: false, InvalidReason: "Incorrect number of columns"}
	}

	pkg := strings.TrimSpace(row[pkgIndex])
	testName := strings.TrimSpace(row[nameIndex])
	flakeRateStr := strings.TrimSpace(row[flakeRateIndex])
	logs := strings.TrimSpace(row[logsIndex])

	t := model.FlakyTicket{
		TestPackage: pkg,
		TestName:    testName,
		Valid:       true, // Assume valid initially
	}

	// --- Priority Calculation ---
	var flakeRate float64
	var parseErr error
	if flakeRateStr == "" || flakeRateStr == "%" { // Handle empty or just "%"
		log.Warn().Str("test", testName).Str("package", pkg).Msg("Missing Flake Rate. Defaulting to Low priority.")
		t.Priority = "Low"
		flakeRateStr = "0" // Use "0" for display
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
			t.FlakeRate = flakeRate // Store numeric value if needed elsewhere
			// Determine priority
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
	// Use original (or '0') flakeRateStr for display to keep formatting
	displayFlakeRate := strings.TrimSuffix(flakeRateStr, "%") // Ensure % isn't doubled up

	// --- Summary and Description ---
	t.Summary = fmt.Sprintf("Fix Flaky Test: %s (%s%% flake rate)", testName, displayFlakeRate)

	// Parse logs into Jira list format
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
			// Basic check if it looks like a URL
			if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
				logSection.WriteString(fmt.Sprintf("* [Run %d|%s]\n", runNumber, link))
				hasLinks = true
			} else {
				// If not a link, just list it
				logSection.WriteString(fmt.Sprintf("* Run %d Log Info: %s\n", runNumber, link))
			}
			runNumber++
		}
		if !hasLinks && runNumber == 1 { // No links were processed
			logSection.Reset() // Clear the builder
			logSection.WriteString("(No valid log links found in source data)")
		}
	}

	// Jira Wiki Markup Description
	// Ensure placeholders are correctly handled if values are empty
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
# *Close Ticket:* Close this ticket once the test is confirmed stable.
# *Guidance:* Refer to the team's [Flaky Test Guide|https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/tools/flakeguard/e2e-flaky-test-guide.md] (Update link if needed).
`,
		pkgDisplay,
		testNameDisplay,
		displayFlakeRate,
		t.Priority, // Include calculated priority
		logSection.String(),
	)

	// --- Final Validation Checks ---
	var missing []string
	if pkg == "" {
		missing = append(missing, "Package")
	}
	if testName == "" {
		missing = append(missing, "Test Name")
	}
	// Flake rate validity already checked

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
// (Implementation assumed to exist as before)
func findExistingTicket(client *jira.Client, label string, ticket model.FlakyTicket) (string, error) {
	// Escape potentially problematic characters in test name for JQL query if necessary
	// Basic escaping for quotes might be needed, more complex names might require more robust escaping.
	safeTestName := strings.ReplaceAll(ticket.TestName, `"`, `\"`)
	// Consider truncating if names are very long and might exceed JQL limits
	// maxLen := 200
	// if len(safeTestName) > maxLen { safeTestName = safeTestName[:maxLen] }

	// Using ~ for contains, might be too broad. Consider exact match if needed, but flake rates change summary.
	// Search for summary containing the test name AND the label. Order by creation date descending to get the latest.
	jql := fmt.Sprintf(`project = "%s" AND labels = "%s" AND summary ~ "\"%s\"" ORDER BY created DESC`, jiraProject, label, safeTestName)
	log.Debug().Str("jql", jql).Msg("Executing Jira search JQL")

	issues, resp, err := client.Issue.SearchWithContext(context.Background(), jql, &jira.SearchOptions{MaxResults: 5}) // Fetch a few to check relevance? Or just 1?
	if err != nil {
		// Check for specific Jira errors if possible (e.g., invalid JQL)
		errMsg := fmt.Sprintf("error searching Jira with JQL '%s'", jql)
		log.Error().Err(err).Interface("response", resp).Msg(errMsg)
		return "", fmt.Errorf("%s: %w", errMsg, err)
	}

	if len(issues) == 0 {
		log.Debug().Str("testName", ticket.TestName).Str("label", label).Msg("No existing tickets found in Jira via search.")
		return "", nil // No matching ticket found
	}

	// Optional: Add more checks here. Does the found ticket *really* match?
	// E.g., check if the summary *starts* with "Fix Flaky Test: TestName" or similar stricter pattern.
	// For now, return the key of the most recently created matching ticket.
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
	if lastDot > 0 { // Ensure dot is not the first char
		ext = filename[lastDot:]
		filename = filename[:lastDot]
	}
	// Append _remaining before the extension
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

	// Create a map of processed tickets by their original RowIndex for quick lookup
	// We need to write the *original* rows for tests that were skipped or had errors.
	processedTicketMap := make(map[int]model.FlakyTicket)
	for _, t := range processedTickets {
		processedTicketMap[t.RowIndex] = t
	}

	// Iterate through original data rows (excluding header)
	for i, originalRow := range originalRecords {
		if i == 0 {
			continue
		} // Skip header row
		rowIndex := i + 1 // 1-based index matching ft.RowIndex

		// Check if this row index was processed by the TUI
		if _, wasProcessed := processedIndices[rowIndex]; wasProcessed {
			// Now check if the processed ticket corresponding to this row was *confirmed*
			if ticket, exists := processedTicketMap[rowIndex]; exists && ticket.Confirmed {
				// It was processed AND confirmed, so DO NOT write it to remaining CSV
				log.Debug().Int("rowIndex", rowIndex).Str("test", ticket.TestName).Msg("Skipping confirmed ticket from remaining CSV output")
				continue
			}
		}
		// If it wasn't processed OR it was processed but NOT confirmed (skipped, error, etc.), write the original row.
		if err := w.Write(originalRow); err != nil {
			// Log error for the specific row but continue trying to write others
			log.Error().Err(err).Int("rowIndex", rowIndex).Msg("Failed to write row to remaining CSV")
		}

	}

	return w.Error() // Return any error encountered during flushing
}

// -------------------------------------------------------------------------------------
// Bubble Tea Model for create-tickets (createModel)
// -------------------------------------------------------------------------------------

// createModel holds the state for the create-tickets TUI.
type createModel struct {
	tickets          []model.FlakyTicket // The filtered list of tickets to process in TUI
	index            int                 // Current ticket being viewed
	processedIndices map[int]bool        // Tracks RowIndex of tickets user interacted with (confirmed, skipped, edited)
	confirmed        int                 // Count of confirmed tickets
	skipped          int                 // Count of skipped/existing tickets
	quitting         bool
	DryRun           bool
	JiraProject      string
	JiraIssueType    string
	JiraClient       *jira.Client
	originalRecords  [][]string                     // Needed for writing remaining CSV
	LocalDB          *localdb.DB                    // Live DB instance for updates
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
	return nil // No initial commands needed
}

// Main Update routing based on mode
func (m createModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle empty state
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

	// Mode-specific updates
	switch m.mode {
	case "promptExisting":
		return updatePromptExisting(m, msg)
	case "ticketCreated":
		return updateTicketCreated(m, msg)
	case "confirmDelete":
		return updateConfirmDelete(m, msg)
	case "normal":
		fallthrough // Default mode if not others
	default:
		return updateNormalMode(m, msg)
	}
}

// updateNormalMode handles keypresses in the default view.
func updateNormalMode(m createModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	// Ensure we don't go out of bounds
	if m.index >= len(m.tickets) || m.index < 0 {
		// This case should ideally lead to quitting state handled elsewhere or logged
		log.Error().Int("index", m.index).Int("len", len(m.tickets)).Msg("Invalid index in updateNormalMode")
		m.quitting = true
		return m, tea.Quit
	}

	t := &m.tickets[m.index] // Use pointer

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Clear any feedback messages from the previous action/state
		// before processing the new keypress for the current item.
		m.infoMessage = ""
		m.errorMessage = ""

		switch msg.String() {
		case "q", "esc": // Quit
			m.quitting = true
			return m, tea.Quit

		case "e": // Enter/Edit Existing Ticket ID
			m.mode = "promptExisting"
			m.inputValue = t.ExistingJiraKey // Pre-fill
			return m, nil

		case "d": // Delete existing ticket (transition to confirm)
			// Logic to check conditions and set mode/errorMessage...
			if t.ExistingJiraKey != "" && m.JiraClient != nil && !m.DryRun {
				m.mode = "confirmDelete"
				// Set a specific message for confirmation, errorMessage might be okay here
				m.errorMessage = fmt.Sprintf("Confirm delete %s? This cannot be undone.", jirautils.GetJiraLink(t.ExistingJiraKey)) // Reusing error for prompt
			} else if t.ExistingJiraKey == "" {
				m.errorMessage = "No Jira ticket associated to delete."
			} else if m.DryRun {
				m.errorMessage = "Cannot delete tickets in Dry Run mode."
			} else {
				m.errorMessage = "Cannot delete tickets: Jira client unavailable."
			}
			return m, nil

		case "c": // Confirm/Create Ticket
			// Logic to check conditions...
			if !t.Valid {
				m.errorMessage = "Cannot create: Invalid data. Press [n] to skip."
				return m, nil
			}
			if t.ExistingJiraKey != "" {
				m.errorMessage = fmt.Sprintf("Cannot create: Ticket %s exists. Press [n] skip/[e] edit.", t.ExistingJiraKey)
				return m, nil
			}
			// Call creation logic which might set messages
			return updateConfirm(m)

		case "n": // Skip / Next
			// updateSkip now handles setting its own (optional) message or clearing errors
			return updateSkip(m)

		default: // Ignore other keys
			return m, nil
		}
	}
	// If not a KeyMsg
	return m, nil
}

// updatePromptExisting handles the mode where user enters a Jira key.
func updatePromptExisting(m createModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	// Ensure index is valid before processing
	if m.index < 0 || m.index >= len(m.tickets) {
		// ... (handle invalid index) ...
		return m, nil
	}
	t := &m.tickets[m.index] // Use pointer

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			newKey := strings.TrimSpace(m.inputValue)
			// --- Input Validation ---
			if newKey != "" && !regexp.MustCompile(`^[A-Z][A-Z0-9]+-\d+$`).MatchString(newKey) {
				m.errorMessage = "Invalid Jira key format (e.g., PROJ-123). Press Esc or Enter empty key to cancel."
				m.infoMessage = ""
				return m, nil
			}

			// --- Update DB ---
			err := m.LocalDB.UpsertEntry(t.TestPackage, t.TestName, newKey, t.SkippedAt, t.AssigneeId)
			if err != nil {
				log.Error().Err(err).Str("test", t.TestName).Msg("Failed to update local DB entry from prompt")
				m.errorMessage = fmt.Sprintf("Error saving to local DB: %v", err)
				m.infoMessage = ""
				return m, nil
			}

			// --- Log Confirmation & Update Model State ---
			if newKey == "" {
				log.Debug().Str("test", t.TestName).Msg("Cleared Jira association for test.")
				t.ExistingJiraKey = ""
				t.ExistingTicketSource = ""
			} else {
				log.Debug().Str("test", t.TestName).Str("newKey", newKey).Msg("Manually set Jira key for test.")
				t.ExistingJiraKey = newKey
				t.ExistingTicketSource = "manual"
			}

			// --- Transition and Advance ---
			m.processedIndices[t.RowIndex] = true
			m.skipped++
			m.mode = "normal"
			m.inputValue = ""
			m.index++
			if m.index >= len(m.tickets) {
				m.quitting = true
			}

			// --- Clear messages ---
			m.infoMessage = ""
			m.errorMessage = ""

			var cmd tea.Cmd // Declare cmd variable
			if m.quitting {
				cmd = tea.Quit // If quitting, Quit command takes precedence
			} else {
				// Force a screen clear before rendering the next view
				cmd = tea.ClearScreen
			}
			return m, cmd // Return model and the command

		case tea.KeyEsc: // Cancel prompt
			m.mode = "normal"
			m.inputValue = ""
			m.errorMessage = ""
			m.infoMessage = "Edit cancelled."
			return m, tea.ClearScreen // Return ClearScreen command

		case tea.KeyBackspace:
			// ... (no changes needed) ...
			if len(m.inputValue) > 0 {
				m.inputValue = m.inputValue[:len(m.inputValue)-1]
			}
			return m, nil

		case tea.KeyRunes:
			// ... (no changes needed) ...
			if !strings.ContainsAny(string(msg.Runes), "\n\t\r") {
				m.inputValue += string(msg.Runes)
			}
			return m, nil
		}
	}
	// If msg type is not KeyMsg or key didn't match
	return m, nil
}

// updateConfirm handles the 'c' key press: creates ticket and changes mode.
func updateConfirm(m createModel) (tea.Model, tea.Cmd) {
	i := m.index
	t := &m.tickets[i] // Use pointer

	// Double check conditions
	if !t.Valid || t.ExistingJiraKey != "" {
		m.errorMessage = "Cannot create ticket (invalid or already exists)."
		m.mode = "normal"
		return m, nil // Return nil command, stay in normal mode
	}

	m.processedIndices[t.RowIndex] = true
	// Setting infoMessage here is probably not useful as it gets cleared immediately
	// or overwritten by the confirmation screen.
	// m.infoMessage = fmt.Sprintf("Processing creation for: %s...", t.Summary)

	// --- Prepare Jira Creation ---
	pillarName := ""
	assigneeForJira := ""
	if t.AssigneeId != "" {
		assigneeForJira = t.AssigneeId
		if userMapping, exists := m.userMap[t.AssigneeId]; exists {
			pillarName = userMapping.PillarName
			log.Debug().Str("assignee", t.AssigneeId).Str("pillar", pillarName).Msg("Found pillar name for assignee")
		} else {
			log.Warn().Str("assignee", t.AssigneeId).Msg("Assignee ID present but no matching entry in user_mapping.json found.")
		}
	}

	// --- Attempt Jira Creation (if not dry run and client exists) ---
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
			m.infoMessage = ""        // Clear info message on error
			return m, tea.ClearScreen // Clear screen before showing normal view with error
		}
		// --- Success Case ---
		log.Info().Str("key", issueKey).Str("summary", t.Summary).Msg("Successfully created Jira ticket")
		t.Confirmed = true
		t.ExistingJiraKey = issueKey
		t.ExistingTicketSource = "jira-created"
		// Update local DB immediately
		errDb := m.LocalDB.UpsertEntry(t.TestPackage, t.TestName, issueKey, t.SkippedAt, t.AssigneeId)
		if errDb != nil {
			log.Error().Err(errDb).Str("key", issueKey).Msg("Failed to update local DB after Jira creation!")
			// Should we revert or just warn? Warn for now.
			m.errorMessage = "WARN: Jira ticket created but failed to update local DB!"
		}
		m.confirmed++
		m.mode = "ticketCreated" // Set mode for confirmation view
		// Clear any previous messages before showing confirmation
		m.errorMessage = ""
		m.infoMessage = ""

		// --- FIX: Return ClearScreen command ---
		return m, tea.ClearScreen // Force clear before showing viewTicketCreated
	} else {
		// --- Dry Run or No Client Case ---
		t.Confirmed = true
		t.ExistingJiraKey = "DRYRUN-" + strconv.Itoa(1000+i)
		t.ExistingTicketSource = "dryrun-created"
		m.confirmed++
		m.mode = "ticketCreated" // Set mode for confirmation view
		// Clear any previous messages before showing confirmation
		m.errorMessage = ""
		m.infoMessage = "" // Info message will be generated by viewTicketCreated

		// --- FIX: Return ClearScreen command ---
		return m, tea.ClearScreen // Force clear before showing viewTicketCreated
	}
}

// updateTicketCreated handles the confirmation screen after a ticket is made.
func updateTicketCreated(m createModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	// This mode just waits for *any* key (except quit signals) to advance
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c": // Allow quitting from confirmation
			m.quitting = true
			return m, tea.Quit
		default: // Any other key advances
			m.mode = "normal" // Back to normal mode
			m.index++         // Move to next ticket
			if m.index >= len(m.tickets) {
				m.quitting = true // Set quitting flag if it was the last ticket
			}
			m.infoMessage = "" // Clear confirmation message
			m.errorMessage = ""
			// If quitting, send Quit command, otherwise nil
			if m.quitting {
				return m, tea.Quit
			}
			return m, nil
		}
	}
	return m, nil
}

// updateSkip handles the 'n' key press: increments skip count and advances.
func updateSkip(m createModel) (tea.Model, tea.Cmd) {
	// Mark the current ticket as processed (interacted with)
	t := &m.tickets[m.index]
	m.processedIndices[t.RowIndex] = true

	m.skipped++
	m.index++
	m.errorMessage = ""

	if m.index >= len(m.tickets) {
		m.quitting = true
		return m, tea.Quit // Quit after skipping the last item
	}
	return m, nil
}

// updateConfirmDelete handles the confirmation prompt for deleting a ticket.
func updateConfirmDelete(m createModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	t := &m.tickets[m.index] // Use pointer

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "y": // Yes, delete the ticket
			m.infoMessage = fmt.Sprintf("Attempting to delete Jira ticket %s...", t.ExistingJiraKey)
			m.errorMessage = ""

			// Call Jira delete function (needs implementation in jirautils)
			err := jirautils.DeleteTicketInJira(m.JiraClient, t.ExistingJiraKey)

			if err != nil {
				errMsg := fmt.Sprintf("Failed to delete Jira ticket %s", t.ExistingJiraKey)
				log.Error().Err(err).Msg(errMsg)
				m.errorMessage = fmt.Sprintf("%s: %v", errMsg, err)
				m.mode = "normal" // Go back to normal mode to show error
			} else {
				log.Info().Str("key", t.ExistingJiraKey).Msg("Successfully deleted Jira ticket")
				m.infoMessage = fmt.Sprintf("Deleted Jira ticket %s.", t.ExistingJiraKey)
				// Clear association in the model and DB
				t.ExistingJiraKey = ""
				t.ExistingTicketSource = ""
				m.LocalDB.UpsertEntry(t.TestPackage, t.TestName, "", t.SkippedAt, t.AssigneeId) // Remove key from DB
				m.processedIndices[t.RowIndex] = true                                           // Mark as processed
				// Decide if deleting counts as skip or something else? Let's count it as skipped for now.
				// m.skipped++ // Or have a deleted counter?

				// Should we advance after delete? Let's stay on the current item, now without a key.
				m.mode = "normal"
			}
			return m, tea.ClearScreen

		case "n", "esc": // No, cancel delete
			m.infoMessage = "Delete cancelled."
			m.errorMessage = ""
			m.mode = "normal"
			return m, nil

		default: // Any other key - ignore
			return m, nil
		}
	}
	return m, nil
}

// updateQuit handles the final state before exiting.
// Currently not explicitly called, quitting is handled by returning tea.Quit.
// Could be used for cleanup if needed.
// func updateQuit(m createModel) (tea.Model, tea.Cmd) {
//     m.quitting = true
//     return m, tea.Quit
// }

// View logic for create-tickets TUI
func (m createModel) View() string {
	// Handle quitting state
	if m.quitting {
		// The finalView function is called outside the TUI loop after it exits.
		// This View just needs to handle the transition.
		return "Processing complete. Exiting...\n"
	}

	// Handle empty state
	if m.index == -1 || len(m.tickets) == 0 {
		return "No tickets to process.\n\n[q] quit\n"
	}

	// Handle potential out-of-bounds index defensively
	if m.index >= len(m.tickets) {
		// This indicates the TUI should be quitting, return the quitting message
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
		sb.WriteString(errorStyle.Render(m.errorMessage) + "\n\n")
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
	sb.WriteString("\n") // Add space after header

	// Details Section
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Test Name:"), valueStyle.Render(t.TestName)))
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Package:"), valueStyle.Render(t.TestPackage)))

	// Priority
	priorityVal := t.Priority
	if priorityVal == "" {
		priorityVal = "(Not Set)"
	}
	// Add color based on priority?
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Priority:"), valueStyle.Render(priorityVal))) // Maybe color code this

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
			// This case shouldn't happen if MissingUserMapping logic is correct
			pillarVal = errorStyle.Render("(Error: Map Lookup Failed)")
		}
	} else {
		assigneeVal = faintStyle.Render("(Not Assigned)")
		pillarVal = faintStyle.Render("(N/A - No Assignee)")
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Assigned To:"), valueStyle.Render(assigneeVal)))
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Pillar Name:"), valueStyle.Render(pillarVal))) // Display derived pillar

	// Existing Jira Key
	existingLine := faintStyle.Render("(None Found)")
	if t.ExistingJiraKey != "" {
		sourceInfo := ""
		if t.ExistingTicketSource != "" {
			sourceInfo = fmt.Sprintf(" (from %s)", t.ExistingTicketSource)
		}
		link := jirautils.GetJiraLink(t.ExistingJiraKey) // Get link if possible
		existingLine = existingStyle.Render(fmt.Sprintf("%s%s", link, sourceInfo))
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Existing Jira:"), existingLine))

	// Summary & Description (Collapsible or truncated?)
	sb.WriteString("\n") // Space before Summary/Desc
	sb.WriteString(summaryStyle.Render("Proposed Summary:") + "\n")
	sb.WriteString(valueStyle.Render(t.Summary) + "\n\n") // Add space after summary

	sb.WriteString(descHeaderStyle.Render("Proposed Description:") + "\n")
	// Limit description length? Add scroll? For now, just print.
	sb.WriteString(descBodyStyle.Render(t.Description) + "\n")

	// Help Line / Actions
	sb.WriteString(helpStyle.Render(buildHelpLine(m)))

	return bodyStyle.Render(sb.String())
}

// viewTicketCreated renders the confirmation screen after creating a ticket.
func viewTicketCreated(m createModel) string {
	t := m.tickets[m.index]
	ticketStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)   // Green Bold
	urlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Underline(true) // Blue Underline
	helpStyle := lipgloss.NewStyle().Faint(true).PaddingTop(1)

	var sb strings.Builder

	sb.WriteString(ticketStyle.Render("Ticket processed successfully!"))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("  Summary: %s\n", t.Summary))
	sb.WriteString(fmt.Sprintf("  Jira Key: %s\n", urlStyle.Render(jirautils.GetJiraLink(t.ExistingJiraKey)))) // Show link
	if t.AssigneeId != "" {
		sb.WriteString(fmt.Sprintf("  Assignee: %s\n", t.AssigneeId))
	}
	if pillar, ok := m.userMap[t.AssigneeId]; ok && pillar.PillarName != "" {
		sb.WriteString(fmt.Sprintf("  Pillar: %s\n", pillar.PillarName))
	}

	sb.WriteString(helpStyle.Render("\nPress any key to continue to the next test, or [q] to quit."))

	return sb.String()
}

// viewPromptExisting renders the input prompt for the Jira key.
func viewPromptExisting(m createModel) string {
	t := m.tickets[m.index]
	promptStyle := lipgloss.NewStyle().Bold(true)
	inputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82")).BorderStyle(lipgloss.NormalBorder()).Padding(0, 1)
	helpStyle := lipgloss.NewStyle().Faint(true).PaddingTop(1)
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).PaddingBottom(1) // Red

	var sb strings.Builder

	if m.errorMessage != "" {
		sb.WriteString(errorStyle.Render(m.errorMessage) + "\n")
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
	t := m.tickets[m.index]
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196")) // Red Bold
	helpStyle := lipgloss.NewStyle().Faint(true).PaddingTop(1)

	var sb strings.Builder

	sb.WriteString(promptStyle.Render(fmt.Sprintf("Permanently delete Jira ticket %s associated with test %s?", t.ExistingJiraKey, t.TestName)))
	sb.WriteString("\n\n")
	// Display the error message which contains the warning
	if m.errorMessage != "" {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render(m.errorMessage) + "\n\n") // Orange warning
	}

	sb.WriteString(helpStyle.Render("[y] Yes, delete ticket  |  [n/esc] No, cancel"))

	return sb.String()
}

// buildHelpLine generates the dynamic help text based on the ticket state.
func buildHelpLine(m createModel) string {
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

// finalView is called *after* the TUI program exits.
// It's not part of the BubbleTea View() interface directly.
func finalView(m createModel) string {
	// This function is likely called from RunE after tea.NewProgram().Run() finishes.
	// We already print summary stats and remaining file path in RunE.
	// So this function might not be strictly necessary anymore, or can be simpler.
	doneStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("82")) // Green
	return doneStyle.Render(fmt.Sprintf(
		"Processing finished. Confirmed: %d, Skipped/Existing: %d.",
		m.confirmed, m.skipped,
	))
	// The detailed summary is printed in RunE.
}
