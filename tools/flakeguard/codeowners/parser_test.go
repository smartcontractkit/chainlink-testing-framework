package codeowners

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindOwners(t *testing.T) {
	tests := []struct {
		name           string
		filePath       string
		patterns       []PatternOwner
		expectedOwners []string
	}{
		{
			name:     "Exact match",
			filePath: "core/services/job/job_test.go",
			patterns: []PatternOwner{
				{Pattern: "/core/services/job", Owners: []string{"@team1", "@team2"}}, // Leading /
			},
			expectedOwners: []string{"@team1", "@team2"},
		},
		{
			name:     "Wildcard match",
			filePath: "core/services/ocr_test.go",
			patterns: []PatternOwner{
				{Pattern: "/core/services/ocr*", Owners: []string{"@ocr-team"}}, // Leading /
			},
			expectedOwners: []string{"@ocr-team"},
		},
		{
			name:     "No match",
			filePath: "core/services/unknown/unknown_test.go",
			patterns: []PatternOwner{
				{Pattern: "/core/services/job", Owners: []string{"@team1"}},     // Leading /
				{Pattern: "/core/services/ocr*", Owners: []string{"@ocr-team"}}, // Leading /
			},
			expectedOwners: nil,
		},
		{
			name:     "Multiple matches, last wins",
			filePath: "core/services/ocr_test.go",
			patterns: []PatternOwner{
				{Pattern: "/core/services/*", Owners: []string{"@general-team"}}, // Leading /
				{Pattern: "/core/services/ocr*", Owners: []string{"@ocr-team"}},  // Leading /
			},
			expectedOwners: []string{"@ocr-team"},
		},
		{
			name:     "Directory match",
			filePath: "core/services/job/subdir/job_test.go",
			patterns: []PatternOwner{
				{Pattern: "/core/services/job", Owners: []string{"@team1"}}, // Leading /
			},
			expectedOwners: []string{"@team1"},
		},
		{
			name:     "Leading slash directory match",
			filePath: "core/capabilities/compute/compute_test.go",
			patterns: []PatternOwner{
				{
					Pattern: "/core/capabilities/",
					Owners:  []string{"@smartcontractkit/keystone", "@smartcontractkit/capabilities-team"},
				},
			},
			expectedOwners: []string{"@smartcontractkit/keystone", "@smartcontractkit/capabilities-team"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			owners := FindOwners(test.filePath, test.patterns)
			assert.Equal(t, test.expectedOwners, owners)
		})
	}
}
