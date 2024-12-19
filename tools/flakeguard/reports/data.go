package reports

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// TestReport reports on the parameters and results of one to many test runs
type TestReport struct {
	ID                 string       `json:"id"`
	GoProject          string       `json:"go_project"`
	HeadSHA            string       `json:"head_sha"`
	BaseSHA            string       `json:"base_sha"`
	RepoURL            string       `json:"repo_url"`
	GitHubWorkflowName string       `json:"github_workflow_name"`
	TestRunCount       int          `json:"test_run_count"`
	RaceDetection      bool         `json:"race_detection"`
	ExcludedTests      []string     `json:"excluded_tests"`
	SelectedTests      []string     `json:"selected_tests"`
	Results            []TestResult `json:"results,omitempty"`
}

// TestResult contains the results and outputs of a single test
type TestResult struct {
	// ReportID is the ID of the report this test result belongs to
	// used mostly for Splunk logging
	ReportID       string              `json:"report_id"`
	TestName       string              `json:"test_name"`
	TestPackage    string              `json:"test_package"`
	PackagePanic   bool                `json:"package_panic"`
	Panic          bool                `json:"panic"`
	Timeout        bool                `json:"timeout"`
	Race           bool                `json:"race"`
	Skipped        bool                `json:"skipped"`
	PassRatio      float64             `json:"pass_ratio"`
	Runs           int                 `json:"runs"`
	Failures       int                 `json:"failures"`
	Successes      int                 `json:"successes"`
	Skips          int                 `json:"skips"`
	Outputs        map[string][]string `json:"-"`              // Temporary storage for outputs during test run
	PassedOutputs  map[string][]string `json:"passed_outputs"` // Outputs for passed runs
	FailedOutputs  map[string][]string `json:"failed_outputs"` // Outputs for failed runs
	Durations      []time.Duration     `json:"durations"`
	PackageOutputs []string            `json:"package_outputs"`
	TestPath       string              `json:"test_path"`
	CodeOwners     []string            `json:"code_owners"`
}

// SummaryData contains aggregated data from a set of test results
type SummaryData struct {
	TotalTests     int     `json:"total_tests"`
	PanickedTests  int     `json:"panicked_tests"`
	RacedTests     int     `json:"raced_tests"`
	FlakyTests     int     `json:"flaky_tests"`
	FlakyTestRatio string  `json:"flaky_test_ratio"`
	TotalRuns      int     `json:"total_runs"`
	PassedRuns     int     `json:"passed_runs"`
	FailedRuns     int     `json:"failed_runs"`
	SkippedRuns    int     `json:"skipped_runs"`
	PassRatio      string  `json:"pass_ratio"`
	MaxPassRatio   float64 `json:"max_pass_ratio"`
}

// SplunkEvent represents a customized splunk event string that helps us distinguish what
// triggered the test to run
type SplunkEvent string

// SplunkType represents what type of data is being sent to Splunk, e.g. a report or a result
type SplunkType string

const (
	Manual      SplunkEvent = "manual"
	Periodic    SplunkEvent = "periodic"
	PullRequest SplunkEvent = "pull_request"

	Report SplunkType = "report"
	Result SplunkType = "result"
)

// SplunkTestReport is the full wrapper structure sent to Splunk for the full test report (sans results)
type SplunkTestReport struct {
	Event      SplunkTestReportEvent `json:"event"`      // https://docs.splunk.com/Splexicon:Event
	SourceType string                `json:"sourcetype"` // https://docs.splunk.com/Splexicon:Sourcetype
}

// SplunkTestReportEvent contains the actual meat of the Splunk test report event
type SplunkTestReportEvent struct {
	Event SplunkEvent `json:"event"`
	Type  SplunkType  `json:"type"`
	Data  TestReport  `json:"data"`
}

// SplunkTestResult is the full wrapper structure sent to Splunk for a single test result
type SplunkTestResult struct {
	Event      SplunkTestResultEvent `json:"event"`      // https://docs.splunk.com/Splexicon:Event
	SourceType string                `json:"sourcetype"` // https://docs.splunk.com/Splexicon:Sourcetype
}

// SplunkTestResultEvent contains the actual meat of the Splunk test result event
type SplunkTestResultEvent struct {
	Event SplunkEvent `json:"event"`
	Type  SplunkType  `json:"type"`
	Data  TestResult  `json:"data"`
}

// Data Processing Functions

func GenerateSummaryData(tests []TestResult, maxPassRatio float64) SummaryData {
	var runs, passes, fails, skips, panickedTests, racedTests, flakyTests, skippedTests int
	for _, result := range tests {
		runs += result.Runs
		passes += result.Successes
		fails += result.Failures
		skips += result.Skips

		// Count tests that were entirely skipped
		if result.Runs == 0 && result.Skipped {
			skippedTests++
		}

		if result.Panic {
			panickedTests++
			flakyTests++
		} else if result.Race {
			racedTests++
			flakyTests++
		} else if !result.Skipped && result.Runs > 0 && result.PassRatio < maxPassRatio {
			flakyTests++
		}
	}

	passPercentage := 100.0
	flakePercentage := 0.0

	if runs > 0 {
		passPercentage = math.Floor((float64(passes)/float64(runs)*100)*100) / 100 // Truncate to 2 decimal places
	}

	totalTests := len(tests)
	if totalTests > 0 {
		flakePercentage = math.Floor((float64(flakyTests)/float64(totalTests)*100)*100) / 100 // Truncate to 2 decimal places
	}

	return SummaryData{
		TotalTests:     totalTests,
		PanickedTests:  panickedTests,
		RacedTests:     racedTests,
		FlakyTests:     flakyTests,
		FlakyTestRatio: fmt.Sprintf("%.2f%%", flakePercentage),
		TotalRuns:      runs,
		PassedRuns:     passes,
		FailedRuns:     fails,
		SkippedRuns:    skips,
		PassRatio:      fmt.Sprintf("%.2f%%", passPercentage),
		MaxPassRatio:   maxPassRatio,
	}
}

func FilterResults(report *TestReport, maxPassRatio float64) *TestReport {
	filteredResults := FilterTests(report.Results, func(tr TestResult) bool {
		return !tr.Skipped && tr.PassRatio < maxPassRatio
	})
	report.Results = filteredResults
	return report
}

func FilterTests(results []TestResult, predicate func(TestResult) bool) []TestResult {
	var filtered []TestResult
	for _, result := range results {
		if predicate(result) {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func aggregate(reportChan <-chan *TestReport, errChan <-chan error, opts *aggregateOptions) (*TestReport, error) {
	testMap := make(map[string]TestResult)
	fullReport := &TestReport{}
	excludedTests := map[string]struct{}{}
	selectedTests := map[string]struct{}{}

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

	// Send report to Splunk before adding test results
	if opts.splunkURL != "" {
		err := sendReportToSplunk(opts.splunkURL, opts.splunkToken, opts.splunkEvent, fullReport)
		if err != nil {
			log.Error().Err(err).Msg("Error sending report to Splunk")
		} else {
			log.Debug().Str("event", string(opts.splunkEvent)).Msg("Report sent to Splunk")
		}
	}

	// Prepare final results
	eg := errgroup.Group{}
	var aggregatedResults []TestResult
	for _, r := range testMap {
		result := r
		aggregatedResults = append(aggregatedResults, result)
		if opts.splunkURL != "" {
			eg.Go(func() error {
				return sendResultsToSplunk(opts.splunkURL, opts.splunkToken, opts.splunkEvent, result)
			})
		}
	}

	sortTestResults(aggregatedResults)
	fullReport.Results = aggregatedResults

	if splunkErr := eg.Wait(); splunkErr != nil {
		log.Error().Err(splunkErr).Msg("Error sending results to Splunk")
	} else {
		log.Debug().Str("event", string(opts.splunkEvent)).Msg("All results sent to Splunk")
	}
	return fullReport, nil
}

// sendReportToSplunk sends meta test report data to Splunk
func sendReportToSplunk(url, token string, event SplunkEvent, report *TestReport) error {
	client := resty.New()
	client.AddRetryAfterErrorCondition().SetRetryCount(3).SetRetryWaitTime(5 * time.Second)
	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Splunk %s", token)).
		SetHeader("Content-Type", "application/json").
		SetBody(SplunkTestReport{
			Event: SplunkTestReportEvent{
				Event: event,
				Type:  Report,
				Data:  *report,
			},
			SourceType: "flakeguard_json",
		}).
		Post(url)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("error sending report to Splunk: %s", resp.String())
	}
	return nil
}

func sendResultsToSplunk(url, token string, event SplunkEvent, results ...TestResult) error {
	client := resty.New()
	client.AddRetryAfterErrorCondition().SetRetryCount(3).SetRetryWaitTime(5 * time.Second)
	eg := errgroup.Group{}
	for _, r := range results {
		result := r
		eg.Go(func() error {
			resp, err := client.R().
				SetHeader("Authorization", fmt.Sprintf("Splunk %s", token)).
				SetHeader("Content-Type", "application/json").
				SetBody(SplunkTestResult{
					Event: SplunkTestResultEvent{
						Event: event,
						Type:  Result,
						Data:  result,
					},
					SourceType: "flakeguard_json",
				}).
				Post(url)
			if err != nil {
				return err
			}
			if resp.IsError() {
				return fmt.Errorf("error sending result to Splunk: %s", resp.String())
			}
			return nil
		})
	}

	return eg.Wait()
}

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

func mergeTestResults(a, b TestResult) TestResult {
	a.Runs += b.Runs
	a.Durations = append(a.Durations, b.Durations...)
	if a.PassedOutputs == nil {
		a.PassedOutputs = make(map[string][]string)
	}
	if a.FailedOutputs == nil {
		a.FailedOutputs = make(map[string][]string)
	}
	for runID, outputs := range b.PassedOutputs {
		a.PassedOutputs[runID] = append(a.PassedOutputs[runID], outputs...)
	}
	for runID, outputs := range b.FailedOutputs {
		a.FailedOutputs[runID] = append(a.FailedOutputs[runID], outputs...)
	}
	a.PackageOutputs = append(a.PackageOutputs, b.PackageOutputs...)
	a.Successes += b.Successes
	a.Failures += b.Failures
	a.Panic = a.Panic || b.Panic
	a.Race = a.Race || b.Race
	a.Skips += b.Skips
	a.Skipped = a.Skipped && b.Skipped

	if a.Runs > 0 {
		a.PassRatio = float64(a.Successes) / float64(a.Runs)
	} else {
		a.PassRatio = -1.0 // Indicate undefined pass ratio for skipped tests
	}

	return a
}

func sortTestResults(results []TestResult) {
	sort.Slice(results, func(i, j int) bool {
		if results[i].TestPackage != results[j].TestPackage {
			return results[i].TestPackage < results[j].TestPackage
		}
		iParts := strings.Split(results[i].TestName, "/")
		jParts := strings.Split(results[j].TestName, "/")
		for k := 0; k < len(iParts) && k < len(jParts); k++ {
			if iParts[k] != jParts[k] {
				return iParts[k] < jParts[k]
			}
		}
		if len(iParts) != len(jParts) {
			return len(iParts) < len(jParts)
		}
		return results[i].PassRatio < results[j].PassRatio
	})
}

func avgDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}
