package model

import (
	"regexp"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
)

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
	Priority             string
	FlakeRate            float64
	SkippedAt            time.Time // timestamp when the ticket was marked as skipped
	MissingUserMapping   bool      // true if the assignee ID exists but has no mapping in user_mapping.json
	PillarName           string    // pillar name from Jira customfield_11016
	JiraStatus           string    // status from Jira
}

// MapTestPackageToUser maps a test package to a user ID using regex patterns
func MapTestPackageToUser(testPackage string, testPatternMap map[string]string) string {
	// Sort patterns by length (longest first) to ensure most specific match
	patterns := make([]string, 0, len(testPatternMap))
	for pattern := range testPatternMap {
		patterns = append(patterns, pattern)
	}
	sort.Slice(patterns, func(i, j int) bool {
		return len(patterns[i]) > len(patterns[j])
	})

	// Try each pattern
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, testPackage)
		if err != nil {
			log.Error().Err(err).Msgf("Error matching pattern %s against package %s", pattern, testPackage)
			continue
		}
		if matched {
			return testPatternMap[pattern]
		}
	}

	// If no match found, return empty string
	return ""
}
