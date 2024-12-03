package reports

import (
	"fmt"
	"path/filepath"
	"strings"
)

// MapTestResultsToPaths maps test results to their corresponding file paths.
func MapTestResultsToPaths(report *TestReport, rootDir string) error {
	// Scan the codebase for test functions
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("error normalizing rootDir: %v", err)
	}

	testFileMap, err := ScanTestFiles(rootDir)
	if err != nil {
		return err
	}

	fmt.Printf("Root Directory: %s\n", rootDir)
	fmt.Printf("Test File Map: %+v\n", testFileMap)

	// Assign file paths to each test result
	for i, result := range report.Results {
		testName := result.TestName
		var filePath string

		if strings.Contains(testName, "/") {
			parentTestName := strings.SplitN(testName, "/", 2)[0]
			if path, exists := testFileMap[parentTestName]; exists {
				filePath = path
			}
		} else if path, exists := testFileMap[testName]; exists {
			filePath = path
		}

		if filePath != "" {
			relFilePath, err := filepath.Rel(rootDir, filePath)
			if err != nil {
				return fmt.Errorf("error getting relative path: %v", err)
			}
			report.Results[i].TestPath = filepath.ToSlash(relFilePath)
		} else {
			fmt.Printf("TestName not mapped: %s\n", testName)
			report.Results[i].TestPath = "NOT FOUND"
		}
	}

	return nil
}
