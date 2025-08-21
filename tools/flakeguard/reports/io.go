package reports

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
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

// LoadAndAggregate reads all JSON files in `resultsDir` and aggregates them
// into a single slice of `TestResult`.
func LoadAndAggregate(resultsDir string) ([]TestResult, error) {
	if _, err := os.Stat(resultsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("results directory does not exist: %s", resultsDir)
	}

	resultsChan := make(chan []TestResult)
	errChan := make(chan error, 1)

	// Walk the directory in a separate goroutine
	go func() {
		defer close(resultsChan)
		defer close(errChan)

		err := filepath.Walk(resultsDir, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("error accessing path %s: %w", path, walkErr)
			}

			// If it's a .json file, parse it, then send the slice to resultsChan
			if !info.IsDir() && filepath.Ext(path) == ".json" {
				parsed, parseErr := processFile(path)
				if parseErr != nil {
					// If we can't parse this file, return an error to halt the walk
					return fmt.Errorf("error processing file '%s': %w", path, parseErr)
				}
				// Send the parsed results from this file into the aggregator
				resultsChan <- parsed
			}

			return nil
		})
		if err != nil {
			errChan <- err
		}
	}()

	// Aggregate all the test results from resultsChan
	return aggregate(resultsChan, errChan)
}

// processFile reads a large JSON report file and creates TestReport objects in a memory-efficient way.
func processFile(filePath string) ([]TestResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	var results []TestResult
	if err := json.NewDecoder(file).Decode(&results); err != nil {
		return nil, fmt.Errorf("error decoding JSON array in %s: %w", filePath, err)
	}

	return results, nil
}

// LoadReport reads a JSON file and returns a TestReport pointer
func LoadReport(filePath string) (*TestReport, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}
	var report TestReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON from file %s: %w", filePath, err)
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
			log.Error().Err(err).Msg("Error flushing buffer")
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

// aggregate listens for slices of TestResult on resultsChan,
// merges them, and returns a single slice of aggregated results once
// the channel is closed. It also checks for errors on errChan.
func aggregate(
	resultsChan <-chan []TestResult,
	errChan <-chan error,
) ([]TestResult, error) {
	// Maps each unique "Package|Name" to a single TestResult
	testMap := make(map[string]TestResult)

	// Consume data from the results channel
	for results := range resultsChan {
		for _, r := range results {
			key := r.TestPackage + "|" + r.TestName
			if existing, found := testMap[key]; found {
				// Merge your results (runs, failures, pass ratios, etc.)
				testMap[key] = mergeTestResults(existing, r)
			} else {
				testMap[key] = r
			}
		}
	}

	// Collect any errors from errChan
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	// Convert the map back to a slice
	aggregated := make([]TestResult, 0, len(testMap))
	for _, r := range testMap {
		aggregated = append(aggregated, r)
	}

	// Sort if desired
	sortTestResults(aggregated)

	return aggregated, nil
}

// ZerologRestyLogger wraps zerolog for Resty's logging interface
type ZerologRestyLogger struct{}

// Errorf logs errors using zerolog's global logger
func (ZerologRestyLogger) Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

// Warnf logs warnings using zerolog's global logger
func (ZerologRestyLogger) Warnf(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}

// Debugf logs debug messages using zerolog's global logger
func (ZerologRestyLogger) Debugf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}
