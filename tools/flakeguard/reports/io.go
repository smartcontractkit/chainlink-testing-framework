package reports

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
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
	splunkEvent SplunkEvent
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
func WithSplunk(url, token string, event SplunkEvent) AggregateOption {
	return func(opts *aggregateOptions) {
		opts.splunkURL = url
		opts.splunkToken = token
		opts.splunkEvent = event
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
		if err = decodeField(decoder, &report); err != nil {
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
func decodeField(decoder *json.Decoder, report *TestReport) error {
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
	case "github_workflow_run_url":
		if err := decoder.Decode(&report.GitHubWorkflowRunURL); err != nil {
			return fmt.Errorf("error decoding GitHubWorkflowRunURL: %w", err)
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
	case "id":
		if err := decoder.Decode(&report.ID); err != nil {
			return fmt.Errorf("error decoding ID: %w", err)
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

// aggregate aggregates multiple TestReport objects into a single TestReport as they are received
func aggregate(reportChan <-chan *TestReport, errChan <-chan error, opts *aggregateOptions) (*TestReport, error) {
	var (
		testMap       = make(map[string]TestResult)
		fullReport    = &TestReport{}
		excludedTests = map[string]struct{}{}
		selectedTests = map[string]struct{}{}
		sendToSplunk  = opts.splunkURL != ""
	)

	fullReport.ID = opts.reportID
	for report := range reportChan {
		if fullReport.GoProject == "" {
			fullReport.GoProject = report.GoProject
		} else if fullReport.GoProject != report.GoProject {
			return nil, fmt.Errorf("reports with different Go projects found, expected %s, got %s", fullReport.GoProject, report.GoProject)
		}
		fullReport.TestRunCount += report.TestRunCount
		fullReport.RaceDetection = report.RaceDetection && fullReport.RaceDetection
		for _, test := range report.ExcludedTests {
			excludedTests[test] = struct{}{}
		}
		for _, test := range report.SelectedTests {
			selectedTests[test] = struct{}{}
		}
		for _, result := range report.Results {
			result.ReportID = opts.reportID
			key := result.TestName + "|" + result.TestPackage
			if existing, found := testMap[key]; found {
				existing = mergeTestResults(existing, result)
				testMap[key] = existing
			} else {
				testMap[key] = result
			}
		}
	}

	for err := range errChan {
		return nil, err
	}

	// Finalize excluded and selected tests
	for test := range excludedTests {
		fullReport.ExcludedTests = append(fullReport.ExcludedTests, test)
	}
	for test := range selectedTests {
		fullReport.SelectedTests = append(fullReport.SelectedTests, test)
	}

	// Get the report without any results to send to splunk
	splunkReport := *fullReport

	// Prepare final results
	var (
		aggregatedResults []TestResult
		err               error
	)
	for _, result := range testMap {
		aggregatedResults = append(aggregatedResults, result)
	}

	sortTestResults(aggregatedResults)
	fullReport.Results = aggregatedResults

	if sendToSplunk {
		err = sendDataToSplunk(opts, splunkReport, aggregatedResults...)
	}
	return fullReport, err
}

// sendDataToSplunk sends a truncated TestReport and each individual TestResults to Splunk as events
func sendDataToSplunk(opts *aggregateOptions, report TestReport, results ...TestResult) error {
	start := time.Now()
	client := resty.New().
		SetBaseURL(opts.splunkURL).
		SetAuthScheme("Splunk").
		SetAuthToken(opts.splunkToken).
		SetHeader("Content-Type", "application/json").
		SetLogger(ZerologRestyLogger{})
	client.AddRetryAfterErrorCondition().SetRetryCount(5).SetTimeout(time.Second * 10)

	log.Debug().Str("report id", report.ID).Int("results", len(results)).Msg("Sending aggregated data to Splunk")

	// Send results
	var (
		splunkErrs            = []error{}
		resultsBatchSize      = 10
		resultsBatch          = []SplunkTestResult{}
		successfulResultsSent = 0
	)
	for resultCount, result := range results {
		resultsBatch = append(resultsBatch, SplunkTestResult{
			Event: SplunkTestResultEvent{
				Event: opts.splunkEvent,
				Type:  Result,
				Data:  result,
			},
			SourceType: SplunkSourceType,
			Index:      SplunkIndex,
		})
		// Send results in batches so Splunk doesn't get mad
		if len(resultsBatch) >= resultsBatchSize || resultCount == len(results)-1 {
			batchData, testNames, err := batchSplunkResults(resultsBatch)
			if err != nil {
				return fmt.Errorf("error batching results: %w", err)
			}

			resp, err := client.R().
				SetBody(batchData.String()).
				Post("")
			if err != nil {
				splunkErrs = append(splunkErrs,
					fmt.Errorf("error sending flakeguard results for [%s] to Splunk: %w", strings.Join(testNames, ", "), err),
				)
			}
			if resp.IsError() {
				splunkErrs = append(splunkErrs,
					fmt.Errorf("error sending flakeguard result for [%s] to Splunk: %s", strings.Join(testNames, ", "), resp.String()),
				)
			}
			if err == nil && !resp.IsError() {
				successfulResultsSent += len(resultsBatch)
			}
			resultsBatch = []SplunkTestResult{}
		}
	}

	// Check if errors occurred while uploading results and send report with incomplete flag
	resp, err := client.R().
		SetBody(SplunkTestReport{
			Event: SplunkTestReportEvent{
				Event:      opts.splunkEvent,
				Type:       Report,
				Data:       report,
				Incomplete: len(splunkErrs) > 0,
			},
			SourceType: SplunkSourceType,
			Index:      SplunkIndex,
		}).
		Post("")

	if err != nil {
		splunkErrs = append(splunkErrs, fmt.Errorf("error sending flakeguard report '%s' to Splunk: %w", report.ID, err))
	}
	if resp.IsError() {
		splunkErrs = append(splunkErrs, fmt.Errorf("error sending flakeguard report '%s' to Splunk: %s", report.ID, resp.String()))
	}
	if len(splunkErrs) > 0 {
		log.Error().
			Int("successfully sent", successfulResultsSent).
			Int("total results", len(results)).
			Errs("errors", splunkErrs).
			Str("duration", time.Since(start).String()).
			Msg("Errors occurred while sending test results to Splunk")
	} else {
		log.Debug().
			Int("successfully sent", successfulResultsSent).
			Int("total results", len(results)).
			Str("duration", time.Since(start).String()).
			Msg("All results sent successfully")
	}
	return errors.Join(splunkErrs...)
}

// batchSplunkResults creates a batch of TestResult objects as individual JSON objects
// Splunk doesn't accept JSON arrays, they want individual events as single JSON objects
// https://docs.splunk.com/Documentation/Splunk/9.4.0/Data/FormateventsforHTTPEventCollector
func batchSplunkResults(results []SplunkTestResult) (batchData bytes.Buffer, resultTestNames []string, err error) {
	for _, result := range results {
		data, err := json.Marshal(result)
		if err != nil {
			return batchData, nil, fmt.Errorf("error marshaling result for '%s': %w", result.Event.Data.TestName, err)
		}
		if _, err := batchData.Write(data); err != nil {
			return batchData, nil, fmt.Errorf("error writing result for '%s': %w", result.Event.Data.TestName, err)
		}
		if _, err := batchData.WriteRune('\n'); err != nil {
			return batchData, nil, fmt.Errorf("error writing newline for '%s': %w", result.Event.Data.TestName, err)
		}
		resultTestNames = append(resultTestNames, result.Event.Data.TestName)
	}
	return batchData, resultTestNames, nil
}

// unBatchSplunkResults un-batches a batch of TestResult objects into a slice of TestResult objects
func unBatchSplunkResults(batch []byte) ([]*SplunkTestResult, error) {
	var results []*SplunkTestResult
	scanner := bufio.NewScanner(bufio.NewReader(bytes.NewReader(batch)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue // Skip empty lines
		}
		var result *SplunkTestResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			return results, fmt.Errorf("error unmarshaling result: %w", err)
		}
		results = append(results, result)
	}
	return results, scanner.Err()
}

// aggregateFromReports aggregates multiple TestReport objects into a single TestReport
func aggregateFromReports(opts *aggregateOptions, reports ...*TestReport) (*TestReport, error) {
	reportChan := make(chan *TestReport, len(reports))
	errChan := make(chan error, 1)
	for _, report := range reports {
		reportChan <- report
	}
	close(reportChan)
	close(errChan)
	return aggregate(reportChan, errChan, opts)
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
