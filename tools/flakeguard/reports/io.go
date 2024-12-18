package reports

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func LoadAndAggregate(resultsPath string) (*TestReport, error) {
	reportChan := make(chan *TestReport)
	errChan := make(chan error, 1)

	// Start file processing in a goroutine
	go func() {
		defer close(reportChan)
		defer close(errChan)

		err := filepath.Walk(resultsPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error accessing path %s: %w", path, err)
			}
			if !info.IsDir() && filepath.Ext(path) == ".json" {
				log.Printf("Processing file: %s", path)
				processLargeFile(path, reportChan, errChan)
			}
			return nil
		})
		if err != nil {
			errChan <- err
		}
	}()

	// Aggregate results as they are being loaded
	aggregatedReport, err := aggregate(reportChan)
	if err != nil {
		return nil, fmt.Errorf("error aggregating reports: %w", err)
	}
	return aggregatedReport, nil
}

func processLargeFile(filePath string, reportChan chan<- *TestReport, errChan chan<- error) {
	file, err := os.Open(filePath)
	if err != nil {
		errChan <- fmt.Errorf("error opening file %s: %w", filePath, err)
		log.Printf("Error opening file: %s, Error: %v", filePath, err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var report TestReport
	token, err := decoder.Token() // Read opening brace '{'
	if err != nil || token != json.Delim('{') {
		errChan <- fmt.Errorf("error reading JSON object start from file %s: %w", filePath, err)
		log.Printf("Error reading JSON object start from file: %s, Error: %v", filePath, err)
		return
	}

	// Parse fields until we reach the end of the object
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			errChan <- fmt.Errorf("error reading JSON token from file %s: %w", filePath, err)
			log.Printf("Error reading JSON token from file: %s, Error: %v", filePath, err)
			return
		}

		fieldName, ok := token.(string)
		if !ok {
			errChan <- fmt.Errorf("unexpected JSON token in file %s", filePath)
			log.Printf("Unexpected JSON token in file: %s, Token: %v", filePath, token)
			return
		}

		switch fieldName {
		case "go_project":
			if err := decoder.Decode(&report.GoProject); err != nil {
				log.Printf("Error decoding GoProject in file: %s, Error: %v", filePath, err)
				return
			}
		case "head_sha":
			if err := decoder.Decode(&report.HeadSHA); err != nil {
				log.Printf("Error decoding HeadSHA in file: %s, Error: %v", filePath, err)
				return
			}
		case "base_sha":
			if err := decoder.Decode(&report.BaseSHA); err != nil {
				log.Printf("Error decoding BaseSHA in file: %s, Error: %v", filePath, err)
				return
			}
		case "repo_url":
			if err := decoder.Decode(&report.RepoURL); err != nil {
				log.Printf("Error decoding RepoURL in file: %s, Error: %v", filePath, err)
				return
			}
		case "github_workflow_name":
			if err := decoder.Decode(&report.GitHubWorkflowName); err != nil {
				log.Printf("Error decoding GitHubWorkflowName in file: %s, Error: %v", filePath, err)
				return
			}
		case "test_run_count":
			if err := decoder.Decode(&report.TestRunCount); err != nil {
				log.Printf("Error decoding TestRunCount in file: %s, Error: %v", filePath, err)
				return
			}
		case "race_detection":
			if err := decoder.Decode(&report.RaceDetection); err != nil {
				log.Printf("Error decoding RaceDetection in file: %s, Error: %v", filePath, err)
				return
			}
		case "excluded_tests":
			if err := decoder.Decode(&report.ExcludedTests); err != nil {
				log.Printf("Error decoding ExcludedTests in file: %s, Error: %v", filePath, err)
				return
			}
		case "selected_tests":
			if err := decoder.Decode(&report.SelectedTests); err != nil {
				log.Printf("Error decoding SelectedTests in file: %s, Error: %v", filePath, err)
				return
			}
		case "results":
			token, err := decoder.Token() // Read opening bracket '['
			if err != nil || token != json.Delim('[') {
				log.Printf("Error reading Results array start in file: %s, Error: %v", filePath, err)
				return
			}

			for decoder.More() {
				var result TestResult
				if err := decoder.Decode(&result); err != nil {
					log.Printf("Error decoding TestResult in file: %s, Error: %v", filePath, err)
					return
				}
				report.Results = append(report.Results, result)
			}

			if _, err := decoder.Token(); err != nil {
				log.Printf("Error reading Results array end in file: %s, Error: %v", filePath, err)
				return
			}
		default:
			// Skip unknown fields
			var skip interface{}
			if err := decoder.Decode(&skip); err != nil {
				log.Printf("Error skipping unknown field: %s in file: %s, Error: %v", fieldName, filePath, err)
				return
			}
			log.Printf("Skipped unknown field: '%s' in file: %s", fieldName, filePath)
		}
	}

	// Read closing brace '}'
	if _, err := decoder.Token(); err != nil {
		log.Printf("Error reading JSON object end in file: %s, Error: %v", filePath, err)
		return
	}

	reportChan <- &report
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
		r.FailedOutputs = nil
		r.PassedOutputs = nil
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

// SaveReport saves a TestReport to a specified file path in JSON format.
// It ensures the file is created or truncated and handles any errors during
// file operations, providing a reliable way to persist test results.
func SaveReport(fs FileSystem, filePath string, report TestReport) error {
	// Open the file with truncation mode
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("error closing file: %v", cerr)
		}
	}()

	// Use a buffered writer for better performance
	bufferedWriter := bufio.NewWriter(file)
	defer func() {
		if err := bufferedWriter.Flush(); err != nil {
			fmt.Printf("error flushing buffer: %v\n", err)
		}
	}()

	// Create a JSON encoder with the buffered writer
	encoder := json.NewEncoder(bufferedWriter)
	encoder.SetIndent("", "  ")

	// Encode the report
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	return nil
}
