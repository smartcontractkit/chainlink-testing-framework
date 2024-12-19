package reports

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// TestReport reports on the parameters and results of one to many test runs
type TestReport struct {
	ID                   string       `json:"id"`
	GoProject            string       `json:"go_project"`
	HeadSHA              string       `json:"head_sha"`
	BaseSHA              string       `json:"base_sha"`
	RepoURL              string       `json:"repo_url"`
	GitHubWorkflowName   string       `json:"github_workflow_name"`
	GitHubWorkflowRunURL string       `json:"github_workflow_run_url"`
	TestRunCount         int          `json:"test_run_count"`
	RaceDetection        bool         `json:"race_detection"`
	ExcludedTests        []string     `json:"excluded_tests"`
	SelectedTests        []string     `json:"selected_tests"`
	Results              []TestResult `json:"results,omitempty"`
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
// triggered the test to run. This is a custom field, different from the Splunk event field.
type SplunkEvent string

// SplunkType represents what type of data is being sent to Splunk, e.g. a report or a result.
// This is a custom field to help us distinguish what kind of data we're sending.
type SplunkType string

const (
	Manual      SplunkEvent = "manual"
	Scheduled   SplunkEvent = "scheduled"
	PullRequest SplunkEvent = "pull_request"

	Report SplunkType = "report"
	Result SplunkType = "result"

	// https://docs.splunk.com/Splexicon:Sourcetype
	SplunkSourceType = "flakeguard_json"
	// https://docs.splunk.com/Splexicon:Index
	SplunkIndex = "github_flakeguard_runs"
)

// SplunkTestReport is the full wrapper structure sent to Splunk for the full test report (sans results)
type SplunkTestReport struct {
	Event      SplunkTestReportEvent `json:"event"`      // https://docs.splunk.com/Splexicon:Event
	SourceType string                `json:"sourcetype"` // https://docs.splunk.com/Splexicon:Sourcetype
	Index      string                `json:"index"`      // https://docs.splunk.com/Splexicon:Index
}

// SplunkTestReportEvent contains the actual meat of the Splunk test report event
type SplunkTestReportEvent struct {
	Event SplunkEvent `json:"event"`
	Type  SplunkType  `json:"type"`
	Data  TestReport  `json:"data"`
	// Incomplete indicates that there were issues uploading test results and the report is incomplete
	Incomplete bool `json:"incomplete"`
}

// SplunkTestResult is the full wrapper structure sent to Splunk for a single test result
type SplunkTestResult struct {
	Event      SplunkTestResultEvent `json:"event"`      // https://docs.splunk.com/Splexicon:Event
	SourceType string                `json:"sourcetype"` // https://docs.splunk.com/Splexicon:Sourcetype
	Index      string                `json:"index"`      // https://docs.splunk.com/Splexicon:Index
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
