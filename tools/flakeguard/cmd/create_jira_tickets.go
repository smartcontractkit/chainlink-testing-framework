package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// CreateTicketsCmd is the command that runs a Bubble Tea TUI for your CSV data.
var CreateTicketsCmd = &cobra.Command{
	Use:   "create-tickets",
	Short: "Interactive TUI to confirm Jira tickets from CSV",
	Long: `Reads a CSV file describing flaky tests and displays each proposed
ticket in a text-based UI. Press 'y' to confirm, 'n' to skip, 'q' to quit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		csvPath, _ := cmd.Flags().GetString("csv-path")
		if csvPath == "" {
			log.Error().Msg("CSV path is required")
			os.Exit(1)
		}

		records, err := readCSV(csvPath)
		if err != nil {
			log.Error().Err(err).Msg("Error reading CSV file")
			os.Exit(1)
		}

		// Convert CSV rows to Ticket objects (skipping invalid rows).
		var tickets []Ticket
		for i, row := range records {
			if len(row) < 10 {
				log.Warn().Msgf("Skipping row %d (not enough columns)", i+1)
				continue
			}
			tickets = append(tickets, rowToTicket(row))
		}

		if len(tickets) == 0 {
			log.Warn().Msg("No valid tickets found in CSV.")
			return nil
		}

		// Create and run Bubble Tea program
		p := tea.NewProgram(initialModel(tickets))
		if _, err := p.Run(); err != nil {
			log.Error().Err(err).Msg("Error running Bubble Tea program")
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	CreateTicketsCmd.Flags().String("csv-path", "", "Path to the CSV file containing ticket data")
	// Add this command to your root command in your main.go or wherever you define your CLI.
	// Example: rootCmd.AddCommand(BubbleTeaJiraCmd)
}

// Ticket is a simple struct holding Summary and Description for each row.
type Ticket struct {
	Summary     string
	Description string
}

// readCSV reads the entire CSV into a [][]string.
func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	return r.ReadAll()
}

// rowToTicket maps one CSV row to a Ticket, embedding your “flaky test” info.
func rowToTicket(row []string) Ticket {
	// CSV columns (based on your format):
	// 0: Package
	// 1: Test Path
	// 2: Test
	// 3: Runs
	// 4: Successes
	// 5: Failures
	// 6: Skips
	// 7: Overall Flake Rate
	// 8: Code Owners
	// 9: Failed Logs

	packageName := row[0]
	testName := row[2]
	flakeRate := row[7]
	codeOwners := row[8]
	failedLogs := row[9]

	if codeOwners == "" {
		codeOwners = "N/A"
	}
	if failedLogs == "" {
		failedLogs = "N/A"
	}

	summary := fmt.Sprintf("[Flaky Test] %s (Rate: %s%%)", testName, flakeRate)

	description := fmt.Sprintf(`Test Details:

Package: %s
Test Name: %s
Flake Rate: %s%% in the last 7 days
Code Owners: %s

Test Logs:
%s

Action Items:

Investigate Failed Test Logs: Thoroughly review the provided logs to identify patterns or common error messages that indicate the root cause.

Fix the Issue: Analyze and address the underlying problem causing the flakiness.

Rerun Tests Locally: Execute the test and related changes on a local environment to ensure that the fix stabilizes the test, as well as all other tests that may be affected.

Unskip the Test: Once confirmed stable, remove any test skip markers to re-enable the test in the CI pipeline.

Reference Guidelines: Follow the recommendations in the Flaky Test Guide.
`, packageName, testName, flakeRate, codeOwners, failedLogs)

	return Ticket{
		Summary:     summary,
		Description: description,
	}
}

// ------------------ Bubble Tea Model ------------------

type model struct {
	tickets   []Ticket
	index     int  // current ticket index
	confirmed int  // how many confirmed
	skipped   int  // how many skipped
	quitting  bool // whether we're exiting
}

// initialModel sets up the Bubble Tea model with the tickets to display.
func initialModel(tickets []Ticket) model {
	return model{
		tickets:   tickets,
		index:     0,
		confirmed: 0,
		skipped:   0,
		quitting:  false,
	}
}

// Custom messages if we needed them (not strictly required here).
type confirmMsg struct{}
type skipMsg struct{}
type quitMsg struct{}

// Init is called when the program starts (we have no initial I/O to perform).
func (m model) Init() tea.Cmd {
	return nil
}

// Update is called on every event (keyboard or otherwise).
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y": // confirm
			return updateConfirm(m)
		case "n": // skip
			return updateSkip(m)
		case "q", "ctrl+c", "esc": // quit
			return updateQuit(m)
		}
	}
	return m, nil
}

func updateConfirm(m model) (tea.Model, tea.Cmd) {
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

// View renders our UI: if done, show summary; otherwise show current ticket.
func (m model) View() string {
	if m.quitting || m.index >= len(m.tickets) {
		return finalView(m)
	}

	t := m.tickets[m.index]

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")) // pinkish
	summaryStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")) // green
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")) // light gray
	helpStyle := lipgloss.NewStyle().
		Faint(true)

	header := headerStyle.Render(
		fmt.Sprintf("Proposed Ticket #%d of %d", m.index+1, len(m.tickets)),
	)
	sum := summaryStyle.Render("Summary: ") + t.Summary
	desc := descStyle.Render("Description:\n" + t.Description)
	help := helpStyle.Render("\nPress [y] to confirm, [n] to skip, [q] to quit.")

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n", header, sum, desc, help)
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
