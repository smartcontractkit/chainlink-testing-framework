package model

import "time"

// FlakyTicket represents a ticket for a flaky test.
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
	ExistingTicketSource string // "localdb" or "jira"
	Assignee             string
	IsSkipped            bool
	SkippedAt            time.Time // timestamp when the ticket was marked as skipped
}
