package reports

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
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

type aggregateOptions struct {
	reportID    string
	splunkURL   string
	splunkToken string
}

// AggregateOption is a functional option for configuring the aggregation process.
type AggregateOption func(*aggregateOptions)

// WithReportID explicitly sets the report ID for the aggregated report.
func WithReportID(reportID string) AggregateOption {
	return func(opts *aggregateOptions) {
		opts.reportID = reportID
	}
}

// WithSplunk also sends the aggregation to a Splunk instance as events.
func WithSplunk(splunkURL, splunkToken string) AggregateOption {
	return func(opts *aggregateOptions) {
		opts.splunkURL = splunkURL
		opts.splunkToken = splunkToken
	}
}

// LoadAndAggregate reads all JSON files in a directory and aggregates the results into a single TestReport.
func LoadAndAggregate(resultsPath string, options ...AggregateOption) (*TestReport, error) {
	if _, err := os.Stat(resultsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("results directory does not exist: %s", resultsPath)
	}

	// Apply options
	opts := aggregateOptions{}
	for _, opt := range options {
		opt(&opts)
	}

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
				log.Debug().Str("path", path).Msg("Processing file")
				err = processLargeFile(path, reportChan)
				if err != nil {
					return fmt.Errorf("error processing file '%s': %w", path, err)
				}
			}
			return nil
		})
		if err != nil {
			errChan <- err
		}
	}()

	if opts.reportID == "" {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("error generating UUID: %w", err)
		}
		opts.reportID = uuid.String()
	}

	// Aggregate results as they are being loaded
	aggregatedReport, err := aggregate(reportChan, errChan, &opts)
	if err != nil {
		return nil, fmt.Errorf("error aggregating reports: %w", err)
	}
	return aggregatedReport, nil
}

// processLargeFile reads a large JSON report file and creates TestReport objects in a memory-efficient way.
func processLargeFile(filePath string, reportChan chan<- *TestReport) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var report TestReport
	token, err := decoder.Token() // Read opening brace '{'
	if err != nil || token != json.Delim('{') {
		return fmt.Errorf("error reading JSON object start from file: %w", err)
	}

	// Parse fields until we reach the end of the object
	for decoder.More() {
		if err = decodeField(decoder, report); err != nil {
			return fmt.Errorf("error decoding field: %w", err)
		}
	}

	// Read closing brace '}'
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("error reading JSON object end from file: %w", err)
	}

	reportChan <- &report
	return nil
}

// decodeField reads a JSON field from the decoder and populates the corresponding field in the TestReport.
func decodeField(decoder *json.Decoder, report TestReport) error {
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("error reading JSON token: %w", err)
	}

	fieldName, ok := token.(string)
	if !ok {
		return fmt.Errorf("unexpected JSON token")
	}

	switch fieldName {
	case "go_project":
		if err := decoder.Decode(&report.GoProject); err != nil {
			return fmt.Errorf("error decoding GoProject: %w", err)
		}
	case "head_sha":
		if err := decoder.Decode(&report.HeadSHA); err != nil {
			return fmt.Errorf("error decoding HeadSHA: %w", err)
		}
	case "base_sha":
		if err := decoder.Decode(&report.BaseSHA); err != nil {
			return fmt.Errorf("error decoding BaseSHA: %w", err)
		}
	case "repo_url":
		if err := decoder.Decode(&report.RepoURL); err != nil {
			return fmt.Errorf("error decoding RepoURL: %w", err)
		}
	case "github_workflow_name":
		if err := decoder.Decode(&report.GitHubWorkflowName); err != nil {
			return fmt.Errorf("error decoding GitHubWorkflowName: %w", err)
		}
	case "test_run_count":
		if err := decoder.Decode(&report.TestRunCount); err != nil {
			return fmt.Errorf("error decoding TestRunCount: %w", err)
		}
	case "race_detection":
		if err := decoder.Decode(&report.RaceDetection); err != nil {
			return fmt.Errorf("error decoding RaceDetection: %w", err)
		}
	case "excluded_tests":
		if err := decoder.Decode(&report.ExcludedTests); err != nil {
			return fmt.Errorf("error decoding ExcludedTests: %w", err)
		}
	case "selected_tests":
		if err := decoder.Decode(&report.SelectedTests); err != nil {
			return fmt.Errorf("error decoding SelectedTests: %w", err)
		}
	case "report_id":
		if err := decoder.Decode(&report.ReportID); err != nil {
			return fmt.Errorf("error decoding ReportID: %w", err)
		}
	case "results":
		token, err := decoder.Token() // Read opening bracket '['
		if err != nil || token != json.Delim('[') {
			return fmt.Errorf("error reading Results array start: %w", err)
		}

		for decoder.More() {
			var result TestResult
			if err := decoder.Decode(&result); err != nil {
				return fmt.Errorf("error decoding TestResult: %w", err)
			}
			report.Results = append(report.Results, result)
		}

		if _, err := decoder.Token(); err != nil {
			return fmt.Errorf("error reading results array end: %w", err)
		}
	default:
		// Skip unknown fields
		var skip any
		if err := decoder.Decode(&skip); err != nil {
			return fmt.Errorf("error skipping unknown field: %w", err)
		}
		log.Warn().Str("field", fieldName).Msg("Skipped unknown field, check the test report struct to see if it's been properly updated")
	}
	return nil
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
