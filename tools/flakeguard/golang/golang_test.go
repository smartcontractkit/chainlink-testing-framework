package golang

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
