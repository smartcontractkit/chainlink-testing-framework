package reports

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// Data Structures

type TestReport struct {
	GoProject     string
	TestRunCount  int
	RaceDetection bool
	ExcludedTests []string
	SelectedTests []string
	Results       []TestResult
}

type TestResult struct {
	TestName       string
	TestPackage    string
	PackagePanic   bool
	Panic          bool
	Timeout        bool
	Race           bool
	Skipped        bool
	PassRatio      float64
	Runs           int
	Failures       int
	Successes      int
	Skips          int
	Outputs        []string
	Durations      []time.Duration
	PackageOutputs []string
	TestPath       string
	CodeOwners     []string
}

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

func Aggregate(reports ...*TestReport) (*TestReport, error) {
	testMap := make(map[string]TestResult)
	fullReport := &TestReport{}
	excludedTests := map[string]struct{}{}
	selectedTests := map[string]struct{}{}

	for _, report := range reports {
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
			key := result.TestName + "|" + result.TestPackage
			if existing, found := testMap[key]; found {
				existing = mergeTestResults(existing, result)
				testMap[key] = existing
			} else {
				testMap[key] = result
			}
		}
	}

	for test := range excludedTests {
		fullReport.ExcludedTests = append(fullReport.ExcludedTests, test)
	}
	for test := range selectedTests {
		fullReport.SelectedTests = append(fullReport.SelectedTests, test)
	}

	var aggregatedResults []TestResult
	for _, result := range testMap {
		aggregatedResults = append(aggregatedResults, result)
	}

	sortTestResults(aggregatedResults)
	fullReport.Results = aggregatedResults

	return fullReport, nil
}

func mergeTestResults(a, b TestResult) TestResult {
	a.Runs += b.Runs
	a.Durations = append(a.Durations, b.Durations...)
	a.Outputs = append(a.Outputs, b.Outputs...)
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
