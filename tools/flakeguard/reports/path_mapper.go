package reports

import (
	"fmt"
	"path/filepath"
	"strings"
)

// MapTestResultsToPaths maps test results to their corresponding file paths.
func MapTestResultsToPaths(report *TestReport, rootDir string) error {
	// Scan the codebase for test functions
	testFileMap, err := ScanTestFiles(rootDir)
	if err != nil {
		return err
	}

	// Assign file paths to each test result
	for i, result := range report.Results {
		testName := result.TestName
		var filePath string

		// Handle subtests
		if strings.Contains(testName, "/") {
			parentTestName := strings.SplitN(testName, "/", 2)[0] // Extract parent test
			if path, exists := testFileMap[parentTestName]; exists {
				filePath = path
			}
		} else if path, exists := testFileMap[testName]; exists {
			// Handle normal tests
			filePath = path
		}

		// Normalize filePath to be relative to the project root
		if filePath != "" {
			relFilePath, err := filepath.Rel(rootDir, filePath)
			if err != nil {
				return fmt.Errorf("error getting relative path: %v", err)
			}
			report.Results[i].TestPath = relFilePath
		} else {
			// Log or mark tests not found in the codebase
			report.Results[i].TestPath = "NOT FOUND"
		}
	}

	return nil
}
