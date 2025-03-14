package reports

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

// SendTestReportToSplunk sends a truncated TestReport and each individual TestResults to Splunk as events
func SendTestReportToSplunk(splunkURL, splunkToken, splunkEvent string, report TestReport) error {
	start := time.Now()
	results := report.Results
	report.Results = nil // Don't send results to Splunk, doing that individually
	// Dry-run mode for example runs
	isExampleRun := strings.Contains(splunkURL, "splunk.example.com")

	client := resty.New().
		SetBaseURL(splunkURL).
		SetAuthScheme("Splunk").
		SetAuthToken(splunkToken).
		SetHeader("Content-Type", "application/json").
		SetLogger(ZerologRestyLogger{})

	log.Debug().Str("report id", report.ID).Int("results", len(results)).Msg("Sending aggregated data to Splunk")

	const (
		resultsBatchSize             = 10
		splunkSizeLimitBytes         = 100_000_000 // 100MB. Actual limit is over 800MB, but that's excessive
		exampleSplunkReportFileName  = "example_results/example_splunk_report.json"
		exampleSplunkResultsFileName = "example_results/example_splunk_results_batch_%d.json"
	)

	var (
		splunkErrs            = []error{}
		resultsBatch          = []SplunkTestResult{}
		successfulResultsSent = 0
		batchNum              = 1
	)

	for resultCount, result := range results {
		// No need to send log outputs to Splunk
		result.FailedOutputs = nil
		result.PassedOutputs = nil
		result.PackageOutputs = nil

		resultsBatch = append(resultsBatch, SplunkTestResult{
			Event: SplunkTestResultEvent{
				Event: splunkEvent,
				Type:  Result,
				Data:  result,
			},
			SourceType: SplunkSourceType,
			Index:      SplunkIndex,
		})

		if len(resultsBatch) >= resultsBatchSize ||
			resultCount == len(results)-1 ||
			binary.Size(resultsBatch) >= splunkSizeLimitBytes {

			batchData, testNames, err := batchSplunkResults(resultsBatch)
			if err != nil {
				return fmt.Errorf("error batching results: %w", err)
			}

			if isExampleRun {
				exampleSplunkResultsFileName := fmt.Sprintf(exampleSplunkResultsFileName, batchNum)
				exampleSplunkResultsFile, err := os.Create(exampleSplunkResultsFileName)
				if err != nil {
					return fmt.Errorf("error creating example Splunk results file: %w", err)
				}
				for _, result := range resultsBatch {
					jsonResult, err := json.Marshal(result)
					if err != nil {
						return fmt.Errorf("error marshaling result for '%s' to json: %w", result.Event.Data.TestName, err)
					}
					_, err = exampleSplunkResultsFile.Write(jsonResult)
					if err != nil {
						return fmt.Errorf("error writing result for '%s' to file: %w", result.Event.Data.TestName, err)
					}
				}
				err = exampleSplunkResultsFile.Close()
				if err != nil {
					return fmt.Errorf("error closing example Splunk results file: %w", err)
				}
			} else {
				resp, err := client.R().SetBody(batchData.String()).Post("")
				if err != nil {
					splunkErrs = append(splunkErrs,
						fmt.Errorf("error sending results for [%s] to Splunk: %w", strings.Join(testNames, ", "), err),
					)
				}
				if resp.IsError() {
					splunkErrs = append(splunkErrs,
						fmt.Errorf("error sending result for [%s] to Splunk: %s", strings.Join(testNames, ", "), resp.String()),
					)
				}
				if err == nil && !resp.IsError() {
					successfulResultsSent += len(resultsBatch)
				}
			}
			resultsBatch = []SplunkTestResult{}
			batchNum++
		}
	}

	if isExampleRun {
		log.Info().Msg("Example Run. See 'example_results/splunk_results' for the results that would be sent to splunk")
	}

	// Sanitize report fields
	report.BranchName = strings.TrimSpace(report.BranchName)

	reportData := SplunkTestReport{
		Event: SplunkTestReportEvent{
			Event:      splunkEvent,
			Type:       Report,
			Data:       report,
			Incomplete: len(splunkErrs) > 0,
		},
		SourceType: SplunkSourceType,
		Index:      SplunkIndex,
	}

	if isExampleRun {
		exampleSplunkReportFile, err := os.Create(exampleSplunkReportFileName)
		if err != nil {
			return fmt.Errorf("error creating example Splunk report file: %w", err)
		}
		jsonReport, err := json.Marshal(reportData)
		if err != nil {
			return fmt.Errorf("error marshaling report: %w", err)
		}
		_, err = exampleSplunkReportFile.Write(jsonReport)
		if err != nil {
			return fmt.Errorf("error writing report: %w", err)
		}
		log.Info().Msgf("Example Run. See '%s' for the results that would be sent to splunk", exampleSplunkReportFileName)
	} else {
		resp, err := client.R().SetBody(reportData).Post("")
		if err != nil {
			splunkErrs = append(splunkErrs, fmt.Errorf("error sending report '%s' to Splunk: %w", report.ID, err))
		}
		if resp.IsError() {
			splunkErrs = append(splunkErrs, fmt.Errorf("error sending report '%s' to Splunk: %s", report.ID, resp.String()))
		}
	}

	if len(splunkErrs) > 0 {
		log.Error().
			Int("successfully sent", successfulResultsSent).
			Int("total results", len(results)).
			Errs("errors", splunkErrs).
			Str("report id", report.ID).
			Str("duration", time.Since(start).String()).
			Msg("Errors occurred while sending test results to Splunk")
	} else {
		log.Debug().
			Int("successfully sent", successfulResultsSent).
			Int("total results", len(results)).
			Int("result batches", batchNum).
			Str("duration", time.Since(start).String()).
			Str("report id", report.ID).
			Msg("All results sent successfully to Splunk")
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
	results := make([]*SplunkTestResult, 0, bytes.Count(batch, []byte{'\n'}))
	scanner := bufio.NewScanner(bytes.NewReader(batch))

	maxCapacity := 1024 * 1024 // 1 MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var pool sync.Pool
	pool.New = func() interface{} { return new(SplunkTestResult) }

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue // Skip empty lines
		}

		result := pool.Get().(*SplunkTestResult)
		if err := json.Unmarshal(line, result); err != nil {
			return results, fmt.Errorf("error unmarshalling result: %w", err)
		}
		results = append(results, result)
	}

	if err := scanner.Err(); err != nil {
		return results, fmt.Errorf("error scanning: %w", err)
	}

	return results, nil
}
