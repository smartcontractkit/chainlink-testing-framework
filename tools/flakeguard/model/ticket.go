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
	AssigneeId           string
	AssigneeName         string
	Priority             string
	FlakeRate            float64
	SkippedAt            time.Time // timestamp when the ticket was marked as skipped
}
