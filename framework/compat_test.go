package framework

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testcase struct {
	name         string
	includes     []string
	excludes     []string
	tags         []string
	expectedTags []string
}

func TestSmokeTagsFilter(t *testing.T) {
	testcases := []testcase{
		{
			name:         "include works - single match",
			includes:     []string{"compat"},
			tags:         []string{"v1.0.0+compat", "v1.0.0", "v2.0.0+compat"},
			expectedTags: []string{"v2.0.0+compat", "v1.0.0+compat"},
		},
		{
			name:         "include works - multiple patterns",
			includes:     []string{"compat", "rc"},
			tags:         []string{"v1.0.0+compat", "v1.0.0-rc1", "v1.0.0", "v2.0.0-beta"},
			expectedTags: []string{"v1.0.0+compat", "v1.0.0-rc1"},
		},
		{
			name:         "exclude works - single pattern",
			excludes:     []string{"rc"},
			tags:         []string{"v1.0.0", "v1.0.0-rc1", "v1.0.0-rc2"},
			expectedTags: []string{"v1.0.0"},
		},
		{
			name:         "exclude works - multiple patterns",
			excludes:     []string{"rc", "beta"},
			tags:         []string{"v1.0.0", "v1.0.0-rc1", "v1.0.0-beta1", "v2.0.0-alpha"},
			expectedTags: []string{"v2.0.0-alpha", "v1.0.0"},
		},
		{
			name:         "include and exclude together",
			includes:     []string{"v1"},
			excludes:     []string{"rc"},
			tags:         []string{"v1.0.0", "v1.0.0-rc1", "v2.0.0"},
			expectedTags: []string{"v1.0.0"},
		},
		{
			name:         "empty include means include all except exclusions",
			includes:     []string{},
			excludes:     []string{"rc"},
			tags:         []string{"v1.0.0", "v1.0.0-rc1", "v2.0.0"},
			expectedTags: []string{"v2.0.0", "v1.0.0"},
		},
		{
			name:         "empty exclusions works",
			includes:     []string{"v1"},
			excludes:     []string{},
			tags:         []string{"v1.0.0", "v1.0.0-rc1", "v2.0.0"},
			expectedTags: []string{"v1.0.0", "v1.0.0-rc1"},
		},
		{
			name:         "both empty returns all tags",
			includes:     []string{},
			excludes:     []string{},
			tags:         []string{"v1.0.0", "v2.0.0", "v3.0.0"},
			expectedTags: []string{"v3.0.0", "v2.0.0", "v1.0.0"},
		},
		{
			name:         "exclude takes precedence over include",
			includes:     []string{"v1"},
			excludes:     []string{"rc"},
			tags:         []string{"v1.0.0", "v1.0.0-rc1", "v1.0.0-rc2", "v2.0.0"},
			expectedTags: []string{"v1.0.0"},
		},
		{
			name:         "no matches returns empty slice",
			includes:     []string{"nonexistent"},
			tags:         []string{"v1.0.0", "v2.0.0"},
			expectedTags: []string{},
		},
		{
			name:         "empty tags returns empty slice",
			includes:     []string{"v1"},
			tags:         []string{},
			expectedTags: []string{},
		},
		{
			name:         "partial matches work correctly",
			includes:     []string{"1.0"},
			tags:         []string{"v1.0.0", "v1.0.1", "v2.0.0"},
			expectedTags: []string{"v1.0.1", "v1.0.0"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tags := SortSemverTags(tc.tags, tc.includes, tc.excludes)
			require.Equal(t, tc.expectedTags, tags, "Test case: %s", tc.name)
		})
	}
}
