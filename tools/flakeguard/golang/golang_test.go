package golang

import (
	"errors"
	"os"
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
	testsToSkip := []SkipTest{
		{Package: "package_a", Name: "TestPackAFail"},
		{Package: "package_b", Name: "TestPackBFail"},
	}
	err = SkipTests(tempDir, testsToSkip)
	assert.NoError(t, err)
}
