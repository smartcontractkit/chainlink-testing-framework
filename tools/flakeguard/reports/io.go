package reports

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileSystem interface and implementations

type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (io.WriteCloser, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

type OSFileSystem struct{}

func (OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (OSFileSystem) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

func (OSFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// LoadReports reads JSON files from a directory and returns a slice of TestReport pointers
func LoadReports(resultsPath string) ([]*TestReport, error) {
	var testReports []*TestReport
	err := filepath.Walk(resultsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return fmt.Errorf("error reading file %s: %w", path, readErr)
			}
			var report TestReport
			if jsonErr := json.Unmarshal(data, &report); jsonErr != nil {
				return fmt.Errorf("error unmarshaling JSON from file %s: %w", path, jsonErr)
			}
			testReports = append(testReports, &report)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return testReports, nil
}

// LoadReport reads a JSON file and returns a TestReport pointer
func LoadReport(filePath string) (*TestReport, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}
	var report TestReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON from file %s: %w", filePath, err)
	}
	return &report, nil
}

func SaveSummaryAsJSON(fs FileSystem, path string, summary SummaryData) error {
	file, err := fs.Create(path)
	if err != nil {
		return fmt.Errorf("error creating JSON summary file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		return fmt.Errorf("error writing JSON summary: %w", err)
	}
	return nil
}

func SaveReportNoLogs(fs FileSystem, filePath string, report TestReport) error {
	var filteredResults []TestResult
	for _, r := range report.Results {
		r.Outputs = nil
		r.PackageOutputs = nil
		filteredResults = append(filteredResults, r)
	}
	report.Results = filteredResults

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling results: %v", err)
	}
	return fs.WriteFile(filePath, data, 0644)
}

func SaveReport(fs FileSystem, filePath string, report TestReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling outputs: %v", err)
	}
	return fs.WriteFile(filePath, data, 0644)
}
