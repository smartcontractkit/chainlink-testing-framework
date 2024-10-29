package reports

type TestResult struct {
	TestName    string
	TestPackage string
	PassRatio   float64
	Runs        int
	Outputs     []string // Stores outputs for a test
}

// FilterFailedTests returns a slice of TestResult where the pass ratio is below the specified threshold.
func FilterFailedTests(results []TestResult, threshold float64) []TestResult {
	var failedTests []TestResult
	for _, result := range results {
		if result.PassRatio < threshold {
			failedTests = append(failedTests, result)
		}
	}
	return failedTests
}
