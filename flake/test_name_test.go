package flake

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

func TestDuplicateTestNamesAreFound(t *testing.T) {
	fullList, err := GetDuplicateTestNames(filepath.Join(utils.ProjectRoot, "flake"))
	require.NoError(t, err)
	require.True(t, len(fullList) > 1, "Should have more tests than just the one duplicate")
	require.Equal(t, 2, fullList["TestDuplicateTestNamesAreFound"], "The duplicate should exist twice")
}

func TestCompareDuplicatesGood(t *testing.T) {
	allTests, err := GetDuplicateTestNames(filepath.Join(utils.ProjectRoot, "flake"))
	require.NoError(t, err)
	flakyTests, err := ReadFlakyTests(filepath.Join(utils.ProjectRoot, "flake", "test_data", "flaky_test_good.json"))
	require.NoError(t, err)
	err = CompareDuplicateTestNamesToFlakeTestNames(flakyTests, allTests)
	require.NoError(t, err, "Should not have any duplicates")
}

func TestCompareDuplicatesBad(t *testing.T) {
	allTests, err := GetDuplicateTestNames(filepath.Join(utils.ProjectRoot, "flake"))
	require.NoError(t, err)
	flakyTests, err := ReadFlakyTests(filepath.Join(utils.ProjectRoot, "flake", "test_data", "flaky_test_bad.json"))
	require.NoError(t, err)
	err = CompareDuplicateTestNamesToFlakeTestNames(flakyTests, allTests)
	require.Error(t, err, "Should not have any duplicates")
	require.True(t, strings.Contains(err.Error(), "test name TestDuplicateTestNamesAreFound is a duplicate test name"), fmt.Sprintf("Error message should contain the expected error: %s", err.Error()))
}

func TestCompareDuplicatesNotExist(t *testing.T) {
	allTests, err := GetDuplicateTestNames(filepath.Join(utils.ProjectRoot, "flake"))
	require.NoError(t, err)
	flakyTests, err := ReadFlakyTests(filepath.Join(utils.ProjectRoot, "flake", "test_data", "flaky_test_not_exist.json"))
	require.NoError(t, err)
	err = CompareDuplicateTestNamesToFlakeTestNames(flakyTests, allTests)
	require.Error(t, err, "Should not have any duplicates")
	require.True(t, strings.Contains(err.Error(), "test name TestNotExist in flaky test file does not exist in project"), fmt.Sprintf("Error message should contain the expected error: %s", err.Error()))
}
