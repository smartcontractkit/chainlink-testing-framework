package testparser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
)

var (
	startPanicRe = regexp.MustCompile(`^panic:`)
	startRaceRe  = regexp.MustCompile(`^WARNING: DATA RACE`)
)

// ParseOptions holds options that control how test results are parsed.
type ParseOptions struct {
	OmitOutputsOnSuccess bool
}

// entry represents a single JSON record from go test -json output.
type entry struct {
	Action  string  `json:"Action"`
	Test    string  `json:"Test"`
	Package string  `json:"Package"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"`
}

func (e entry) String() string {
	return fmt.Sprintf("Action: %s, Test: %s, Package: %s, Output: %s, Elapsed: %f", e.Action, e.Test, e.Package, e.Output, e.Elapsed)
}

// ParseTestResults parses the test output JSON files and produces test results.
// Note that any pre-processing (such as transforming the files to ignore parent failures)
// must be performed before calling this function.
func ParseTestResults(jsonOutputPaths []string, runPrefix string, runCount int, options ParseOptions) ([]reports.TestResult, error) {
	var (
		testDetails         = make(map[string]*reports.TestResult)
		panickedPackages    = map[string]struct{}{}
		packageLevelOutputs = make(map[string][]string)
		testsWithSubTests   = make(map[string][]string)
		panicDetectionMode  = false
		raceDetectionMode   = false
		detectedEntries     []entry
		expectedRuns        = runCount
	)

	runNumber := 0
	// Process each JSON output file.
	for _, filePath := range jsonOutputPaths {
		runNumber++
		runID := fmt.Sprintf("%s%d", runPrefix, runNumber)
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open test output file: %w", err)
		}

		scanner := bufio.NewScanner(file)
		var precedingLines []string // context for error reporting
		var followingLines []string

		for scanner.Scan() {
			line := scanner.Text()
			precedingLines = append(precedingLines, line)
			if len(precedingLines) > 15 {
				precedingLines = precedingLines[1:]
			}

			var entryLine entry
			if err := json.Unmarshal(scanner.Bytes(), &entryLine); err != nil {
				// Gather extra context for error reporting.
				for scanner.Scan() && len(followingLines) < 15 {
					followingLines = append(followingLines, scanner.Text())
				}
				context := append(precedingLines, followingLines...)
				return nil, fmt.Errorf("failed to parse json test output near lines:\n%s\nerror: %w", strings.Join(context, "\n"), err)
			}

			var result *reports.TestResult
			if entryLine.Test != "" {
				// Build a key with package and test name.
				key := fmt.Sprintf("%s/%s", entryLine.Package, entryLine.Test)
				parentTestName, subTestName := parseSubTest(entryLine.Test)
				if subTestName != "" {
					parentTestKey := fmt.Sprintf("%s/%s", entryLine.Package, parentTestName)
					testsWithSubTests[parentTestKey] = append(testsWithSubTests[parentTestKey], subTestName)
				}
				if _, exists := testDetails[key]; !exists {
					testDetails[key] = &reports.TestResult{
						TestName:       entryLine.Test,
						TestPackage:    entryLine.Package,
						PassRatio:      0,
						PassedOutputs:  make(map[string][]string),
						FailedOutputs:  make(map[string][]string),
						PackageOutputs: []string{},
					}
				}
				result = testDetails[key]
			}

			// Process output field and handle panic/race detection.
			if entryLine.Output != "" {
				if panicDetectionMode || raceDetectionMode {
					detectedEntries = append(detectedEntries, entryLine)
					continue
				} else if startPanicRe.MatchString(entryLine.Output) {
					panickedPackages[entryLine.Package] = struct{}{}
					detectedEntries = append(detectedEntries, entryLine)
					panicDetectionMode = true
					continue
				} else if startRaceRe.MatchString(entryLine.Output) {
					detectedEntries = append(detectedEntries, entryLine)
					raceDetectionMode = true
					continue
				} else if entryLine.Test != "" && entryLine.Action == "output" {
					if result.Outputs == nil {
						result.Outputs = make(map[string][]string)
					}
					result.Outputs[runID] = append(result.Outputs[runID], entryLine.Output)
				} else if entryLine.Test == "" {
					packageLevelOutputs[entryLine.Package] = append(packageLevelOutputs[entryLine.Package], entryLine.Output)
				} else {
					switch entryLine.Action {
					case "pass":
						result.PassedOutputs[runID] = append(result.PassedOutputs[runID], entryLine.Output)
					case "fail":
						result.FailedOutputs[runID] = append(result.FailedOutputs[runID], entryLine.Output)
					}
				}
			}

			// If in panic or race detection mode, wait for a "fail" action to close the block.
			if (panicDetectionMode || raceDetectionMode) && entryLine.Action == "fail" {
				if panicDetectionMode {
					var outputs []string
					for _, entry := range detectedEntries {
						outputs = append(outputs, entry.Output)
					}
					panicTest, timeout, err := attributePanicToTest(outputs)
					if err != nil {
						log.Warn().Err(err).Msg("Unable to attribute panic to a test")
						panicTest = "UnableToAttributePanicTestPleaseInvestigate"
					}
					panicTestKey := fmt.Sprintf("%s/%s", entryLine.Package, panicTest)
					result, exists := testDetails[panicTestKey]
					if !exists {
						result = &reports.TestResult{
							TestName:       panicTest,
							TestPackage:    entryLine.Package,
							PassRatio:      0,
							PassedOutputs:  make(map[string][]string),
							FailedOutputs:  make(map[string][]string),
							PackageOutputs: []string{},
						}
						testDetails[panicTestKey] = result
					}
					result.Panic = true
					result.Timeout = timeout
					result.Failures++
					result.Runs++
					for _, entry := range detectedEntries {
						if entry.Test == "" {
							result.PackageOutputs = append(result.PackageOutputs, entry.Output)
						} else {
							result.FailedOutputs[runID] = append(result.FailedOutputs[runID], entry.Output)
						}
					}
				} else if raceDetectionMode {
					raceTest, err := attributeRaceToTest(entryLine.Package, detectedEntries)
					if err != nil {
						return nil, err
					}
					raceTestKey := fmt.Sprintf("%s/%s", entryLine.Package, raceTest)
					result, exists := testDetails[raceTestKey]
					if !exists {
						result = &reports.TestResult{
							TestName:       raceTest,
							TestPackage:    entryLine.Package,
							PassRatio:      0,
							PassedOutputs:  make(map[string][]string),
							FailedOutputs:  make(map[string][]string),
							PackageOutputs: []string{},
						}
						testDetails[raceTestKey] = result
					}
					result.Race = true
					result.Failures++
					result.Runs++
					for _, entry := range detectedEntries {
						if entry.Test == "" {
							result.PackageOutputs = append(result.PackageOutputs, entry.Output)
						} else {
							result.FailedOutputs[runID] = append(result.FailedOutputs[runID], entry.Output)
						}
					}
				}
				detectedEntries = []entry{}
				panicDetectionMode = false
				raceDetectionMode = false
				continue
			}

			// Process pass, fail, and skip actions.
			switch entryLine.Action {
			case "pass":
				if entryLine.Test != "" {
					duration, err := time.ParseDuration(strconv.FormatFloat(entryLine.Elapsed, 'f', -1, 64) + "s")
					if err != nil {
						return nil, fmt.Errorf("failed to parse duration: %w", err)
					}
					result.Durations = append(result.Durations, duration)
					result.Successes++
					if result.PassedOutputs == nil {
						result.PassedOutputs = make(map[string][]string)
					}
					result.PassedOutputs[runID] = result.Outputs[runID]
					delete(result.Outputs, runID)
				}
			case "fail":
				if entryLine.Test != "" {
					duration, err := time.ParseDuration(strconv.FormatFloat(entryLine.Elapsed, 'f', -1, 64) + "s")
					if err != nil {
						return nil, fmt.Errorf("failed to parse duration: %w", err)
					}
					result.Durations = append(result.Durations, duration)
					result.Failures++
					if result.FailedOutputs == nil {
						result.FailedOutputs = make(map[string][]string)
					}
					result.FailedOutputs[runID] = result.Outputs[runID]
					delete(result.Outputs, runID)
				}
			case "skip":
				if entryLine.Test != "" {
					result.Skipped = true
					result.Skips++
				}
			}
			if entryLine.Test != "" {
				result.Runs = result.Successes + result.Failures
				if result.Runs > 0 {
					result.PassRatio = float64(result.Successes) / float64(result.Runs)
				} else {
					result.PassRatio = 1
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading test output file: %w", err)
		}
		if err = file.Close(); err != nil {
			log.Warn().Err(err).Str("file", filePath).Msg("failed to close file")
		}
	}

	// Propagate panic status from parent tests to subtests.
	for parentTestKey, subTests := range testsWithSubTests {
		if parentTestResult, exists := testDetails[parentTestKey]; exists {
			if parentTestResult.Panic {
				for _, subTest := range subTests {
					subTestKey := fmt.Sprintf("%s/%s/%s", parentTestResult.TestPackage, parentTestResult.TestName, subTest)
					if subTestResult, exists := testDetails[subTestKey]; exists {
						if subTestResult.Failures > 0 {
							subTestResult.Panic = true
							if subTestResult.FailedOutputs == nil {
								subTestResult.FailedOutputs = make(map[string][]string)
							}
							for runID := range subTestResult.FailedOutputs {
								subTestResult.FailedOutputs[runID] = append(subTestResult.FailedOutputs[runID], "Panic in parent test")
							}
						}
					} else {
						log.Warn().Str("expected subtest", subTestKey).Str("parent test", parentTestKey).Msg("expected subtest not found in parent test")
					}
				}
			}
		} else {
			log.Warn().Str("parent test", parentTestKey).Msg("expected parent test not found")
		}
	}

	var results []reports.TestResult
	for _, result := range testDetails {
		// Correct for possible double-counting caused by panics.
		if result.Runs > expectedRuns {
			if result.Panic {
				result.Failures = expectedRuns
				result.Runs = expectedRuns
			} else {
				log.Warn().Str("test", result.TestName).Int("actual runs", result.Runs).Int("expected runs", expectedRuns).Msg("unexpected test runs")
			}
		}
		if outputs, exists := packageLevelOutputs[result.TestPackage]; exists {
			result.PackageOutputs = outputs
		}
		results = append(results, *result)
	}

	if options.OmitOutputsOnSuccess {
		for i := range results {
			results[i].PassedOutputs = make(map[string][]string)
			results[i].Outputs = make(map[string][]string)
		}
	}

	return results, nil
}

// attributePanicToTest extracts the test function name causing a panic.
func attributePanicToTest(outputs []string) (test string, timeout bool, err error) {
	testNameRe := regexp.MustCompile(`(?:.*\.)?(Test[A-Z]\w+)(?:\.[^(]+)?\s*\(`)
	timeoutRe := regexp.MustCompile(`(?i)(timeout|timedout|timed\s*out)`)
	for _, o := range outputs {
		if matches := testNameRe.FindStringSubmatch(o); len(matches) > 1 {
			testName := strings.TrimSpace(matches[1])
			if timeoutRe.MatchString(o) {
				return testName, true, nil
			}
			return testName, false, nil
		}
	}
	return "", false, fmt.Errorf("failed to attribute panic to test using regex '%s' on these strings:\n\n%s", testNameRe.String(), strings.Join(outputs, ""))
}

// attributeRaceToTest extracts the test function name causing a race condition.
func attributeRaceToTest(racePackage string, raceEntries []entry) (string, error) {
	regexSanitizeRacePackage := filepath.Base(racePackage)
	raceAttributionRe := regexp.MustCompile(fmt.Sprintf(`%s\.(Test[^\.\(]+)`, regexSanitizeRacePackage))
	var entriesOutputs []string
	for _, entry := range raceEntries {
		entriesOutputs = append(entriesOutputs, entry.Output)
		if matches := raceAttributionRe.FindStringSubmatch(entry.Output); len(matches) > 1 {
			testName := strings.TrimSpace(matches[1])
			return testName, nil
		}
	}
	return "", fmt.Errorf("failed to attribute race to test using regex %s on these strings:\n%s", raceAttributionRe.String(), strings.Join(entriesOutputs, ""))
}

// parseSubTest splits a test name into parent and subtest names.
func parseSubTest(testName string) (parentTestName, subTestName string) {
	parts := strings.SplitN(testName, "/", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
