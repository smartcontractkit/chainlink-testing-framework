package reports

type TestResult struct {
	TestName    string
	TestPackage string
	PassRatio   float64
	Skipped     bool // Indicates if the test was skipped
	Runs        int
	Outputs     []string  // Stores outputs for a test
	Durations   []float64 // Stores elapsed time in seconds for each run of the test
}

// FilterFailedTests returns a slice of TestResult where the pass ratio is below the specified threshold.
func FilterFailedTests(results []TestResult, threshold float64) []TestResult {
	var failedTests []TestResult
	for _, result := range results {
		if !result.Skipped && result.PassRatio < threshold {
			failedTests = append(failedTests, result)
		}
	}
	return failedTests
}

// FilterPassedTests returns a slice of TestResult where the tests passed and were not skipped.
func FilterPassedTests(results []TestResult, threshold float64) []TestResult {
	var passedTests []TestResult
	for _, result := range results {
		if !result.Skipped && result.PassRatio >= threshold {
			passedTests = append(passedTests, result)
		}
	}
	return passedTests
}

// FilterSkippedTests returns a slice of TestResult where the tests were skipped.
func FilterSkippedTests(results []TestResult) []TestResult {
	var skippedTests []TestResult
	for _, result := range results {
		if result.Skipped {
			skippedTests = append(skippedTests, result)
		}
	}
	return skippedTests
}
