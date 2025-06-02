package golang

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock version of hasTests function to simulate various scenarios
func mockHasTests(pkgName string) (bool, error) {
	switch pkgName {
	case "pkgWithTests":
		return true, nil
	case "pkgWithoutTests":
		return false, nil
	case "pkgWithError":
		return false, errors.New("test error")
	default:
		return false, nil
	}
}

func TestFilterPackagesWithTests(t *testing.T) {
	t.Parallel()

	// Replace hasTests with mock function
	originalHasTests := hasTests
	hasTests = mockHasTests
	defer func() { hasTests = originalHasTests }() // Restore original function after test

	t.Run("should return packages that contain tests", func(t *testing.T) {
		pkgs := []string{"pkgWithTests", "pkgWithoutTests", "pkgWithError"}
		expected := []string{"pkgWithTests"}

		result := FilterPackagesWithTests(pkgs)

		assert.Equal(t, expected, result, "Expected packages with tests only")
	})

	t.Run("should return an empty slice when all packages have no tests", func(t *testing.T) {
		pkgs := []string{"pkgWithoutTests"}
		expected := []string{}

		result := FilterPackagesWithTests(pkgs)

		assert.Equal(t, expected, result, "Expected empty slice for packages without tests")
	})

	t.Run("should handle error scenarios gracefully", func(t *testing.T) {
		pkgs := []string{"pkgWithError"}
		expected := []string{}

		result := FilterPackagesWithTests(pkgs)

		assert.Equal(t, expected, result, "Expected empty slice for packages with errors")
	})
}

func TestSkipTests(t *testing.T) {
	t.Parallel()

	// Create a temp dir copying over the testdata directory
	tempDir, err := os.MkdirTemp("./", "testdata-*")
	require.NoError(t, err, "Failed to create temp dir")
	t.Cleanup(func() {
		if !t.Failed() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Fatalf("Failed to remove temp dir: %v", err)
			}
			return
		}
		t.Logf("Skipping cleanup of temp dir %s because test failed, leaving it for debugging", tempDir)
	})

	// Copy the testdata directory to the temp dir for running the test on
	err = os.CopyFS(tempDir, os.DirFS("testdata"))
	require.NoError(t, err, "Error copying testdata to temp dir")

	var (
		expectedRan     = []string{"TestPackAFail", "TestPackBFail", "TestPackAPass", "TestPackBPass", "TestPackAFailTrick", "TestPackBFailTrick", "TestPackAAlreadySkipped", "TestPackBAlreadySkipped"}
		expectedPassed  = []string{"TestPackAPass", "TestPackBPass", "TestPackAFailTrick", "TestPackBFailTrick"}
		expectedSkipped = []string{"TestPackAFail", "TestPackBFail", "TestPackAAlreadySkipped", "TestPackBAlreadySkipped"}
		testsToSkip     = []*SkipTest{
			{Package: "github.com/owner/repo/testdata/package_a", Name: "TestPackAFail"},
			{Package: "github.com/owner/repo/testdata/package_b", Name: "TestPackBFail"},
		}
	)
	slices.Sort(expectedRan)
	slices.Sort(expectedPassed)
	slices.Sort(expectedSkipped)

	err = SkipTests(tempDir, testsToSkip)
	assert.NoError(t, err)

	for _, test := range testsToSkip {
		if test.Name == "TestPackAFail" || test.Name == "TestPackBFail" {
			assert.True(t, test.Skipped, "Expected flaky test to be skipped")
		} else if test.Name == "TestPackAAlreadySkipped" || test.Name == "TestPackBAlreadySkipped" {
			assert.False(t, test.Skipped, "Expected already skipped test to not be skipped again")
		}
	}

	// Run the tests to make sure only the tests that should be skipped are skipped
	cmd := exec.Command("go", "test", "-json", "-v", "./...")
	cmd.Dir = tempDir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Error running tests after modifying test files")
	err = os.WriteFile(filepath.Join(tempDir, "test_output.json"), out, 0644)
	require.NoError(t, err, "Error writing test output to file")

	ran, passed, skipped, err := parseTestOutput(t, filepath.Join(tempDir, "test_output.json"))
	require.NoError(t, err, "Error parsing test output")

	assert.EqualValues(t, expectedRan, ran, "Expected certain tests to be run")
	assert.EqualValues(t, expectedPassed, passed, "Expected certain tests to pass")
	assert.EqualValues(t, expectedSkipped, skipped, "Expected certain tests to be skipped")
}

type goTestOutput struct {
	Package string `json:"package"`
	Test    string `json:"test"`
	Action  string `json:"action"`
	Output  string `json:"output"`
}

// Don't use the other flakeguard parsing as they're too complex for this, and we shouldn't test them here
func parseTestOutput(t *testing.T, jsonFile string) (ran, passed, skipped []string, err error) {
	t.Helper()

	f, err := os.Open(jsonFile)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	seen := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var event goTestOutput
		line := scanner.Bytes()
		if err := json.Unmarshal(line, &event); err != nil {
			continue // skip lines that aren't test events
		}
		if event.Test == "" {
			continue
		}
		switch event.Action {
		case "run":
			if !seen[event.Test] {
				ran = append(ran, event.Test)
				seen[event.Test] = true
			}
		case "pass":
			passed = append(passed, event.Test)
		case "skip":
			skipped = append(skipped, event.Test)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, nil, err
	}
	slices.Sort(ran)
	slices.Sort(passed)
	slices.Sort(skipped)
	return ran, passed, skipped, nil
}
