package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
)

var (
	startPanicRe = regexp.MustCompile(`^panic:`)
	startRaceRe  = regexp.MustCompile(`^WARNING: DATA RACE`)
)

type Runner struct {
	ProjectPath          string   // Path to the Go project directory.
	Verbose              bool     // If true, provides detailed logging.
	RunCount             int      // Number of times to run the tests.
	UseRace              bool     // Enable race detector.
	UseShuffle           bool     // Enable test shuffling. -shuffle=on flag.
	ShuffleSeed          string   // Set seed for test shuffling -shuffle={seed} flag. Must be used with UseShuffle.
	FailFast             bool     // Stop on first test failure.
	SkipTests            []string // Test names to exclude.
	SelectedTestPackages []string // Explicitly selected packages to run.
	CollectRawOutput     bool     // Collect test output for later inspection.
	rawOutputs           map[string]*bytes.Buffer
}

// RunTests executes the tests for each provided package and aggregates all results.
// It returns all test results and any error encountered during testing.
func (r *Runner) RunTests() ([]reports.TestResult, error) {
	var jsonFilePaths []string
	for _, p := range r.SelectedTestPackages {
		for i := 0; i < r.RunCount; i++ {
			if r.CollectRawOutput { // Collect raw output for debugging
				if r.rawOutputs == nil {
					r.rawOutputs = make(map[string]*bytes.Buffer)
				}
				if _, exists := r.rawOutputs[p]; !exists {
					r.rawOutputs[p] = &bytes.Buffer{}
				}
				separator := strings.Repeat("-", 80)
				r.rawOutputs[p].WriteString(fmt.Sprintf("Run %d%s\n", i+1, separator))
			}
			jsonFilePath, passed, err := r.runTests(p)
			if err != nil {
				return nil, fmt.Errorf("failed to run tests in package %s: %w", p, err)
			}
			jsonFilePaths = append(jsonFilePaths, jsonFilePath)
			if !passed && r.FailFast {
				break
			}
		}
	}

	reports, err := parseTestResults(jsonFilePaths)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test results: %w", err)
	}
	for index, report := range reports {
		if report.Panics >= r.RunCount { // This feels hacky, but there aren't any elegant solutions
			reports[index].Failures = 0 // We can sometimes double-count panics as failures
			reports[index].Panics = r.RunCount
			reports[index].Runs = r.RunCount
		}
	}
	return reports, nil
}

// RawOutput retrieves the raw output from the test runs, if CollectRawOutput enabled.
// packageName : raw output
func (r *Runner) RawOutputs() map[string]*bytes.Buffer {
	return r.rawOutputs
}

type exitCoder interface {
	ExitCode() int
}

// runTests runs the tests for a given package and returns the path to the output file.
func (r *Runner) runTests(packageName string) (string, bool, error) {
	args := []string{"test", packageName, "-json", "-count=1"}
	if r.UseRace {
		args = append(args, "-race")
	}
	if r.UseShuffle {
		if r.ShuffleSeed != "" {
			args = append(args, fmt.Sprintf("-shuffle=%s", r.ShuffleSeed))
		} else {
			args = append(args, "-shuffle=on")
		}
	}
	if len(r.SkipTests) > 0 {
		skipPattern := strings.Join(r.SkipTests, "|")
		args = append(args, fmt.Sprintf("-skip=%s", skipPattern))
	}
	args = append(args, "2>/dev/null")

	if r.Verbose {
		log.Printf("Running command: go %s\n", strings.Join(args, " "))
	}

	// Create a temporary file to store the output
	tmpFile, err := os.CreateTemp("", "test-output-*.json")
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Run the command with output directed to the file
	cmd := exec.Command("go", args...)
	cmd.Dir = r.ProjectPath
	if r.CollectRawOutput {
		cmd.Stdout = io.MultiWriter(tmpFile, r.rawOutputs[packageName])
		cmd.Stderr = io.MultiWriter(tmpFile, r.rawOutputs[packageName])
	} else {
		cmd.Stdout = tmpFile
		cmd.Stderr = tmpFile
	}

	err = cmd.Run()
	if err != nil {
		var exErr exitCoder
		// Check if the error is due to a non-zero exit code
		if errors.As(err, &exErr) && exErr.ExitCode() == 0 {
			return "", false, fmt.Errorf("test command failed at %s: %w", packageName, err)
		}
		return tmpFile.Name(), false, nil // Test failed
	}

	return tmpFile.Name(), true, nil // Test succeeded
}

type entry struct {
	Action  string  `json:"Action"`
	Test    string  `json:"Test"`
	Package string  `json:"Package"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"` // Decimal value in seconds
}

func (e entry) String() string {
	return fmt.Sprintf("Action: %s, Test: %s, Package: %s, Output: %s, Elapsed: %f", e.Action, e.Test, e.Package, e.Output, e.Elapsed)
}

// parseTestResults reads the test output files and returns the parsed test results.
func parseTestResults(filePaths []string) ([]reports.TestResult, error) {
	var (
		testDetails         = make(map[string]*reports.TestResult) // Holds run, pass counts, and other details for each test
		panickedPackages    = map[string]struct{}{}                // Packages with tests that panicked
		racePackages        = map[string]struct{}{}                // Packages with tests that raced
		packageLevelOutputs = map[string][]string{}                // Package-level outputs
		panicDetectionMode  = false
		raceDetectionMode   = false
		detectedEntries     = []entry{} // race or panic entries
	)

	// Process each file
	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open test output file: %w", err)
		}

		scanner := bufio.NewScanner(file)
		var precedingLines []string // Store preceding lines for context
		var followingLines []string // To collect lines after an error

		for scanner.Scan() {
			line := scanner.Text()
			precedingLines = append(precedingLines, line)

			// Limit precedingLines to the last 15 lines
			if len(precedingLines) > 15 {
				precedingLines = precedingLines[1:]
			}

			var entryLine entry
			if err := json.Unmarshal(scanner.Bytes(), &entryLine); err != nil {
				// Collect 15 lines after the error for more context
				for scanner.Scan() && len(followingLines) < 15 {
					followingLines = append(followingLines, scanner.Text())
				}

				// Combine precedingLines and followingLines to provide 15 lines before and after
				context := append(precedingLines, followingLines...)
				return nil, fmt.Errorf("failed to parse json test output near lines:\n%s\nerror: %w", strings.Join(context, "\n"), err)
			}

			// Only create TestResult for test-level entries
			var result *reports.TestResult
			if entryLine.Test != "" {
				// Determine the key
				key := fmt.Sprintf("%s/%s", entryLine.Package, entryLine.Test)

				if _, exists := testDetails[key]; !exists {
					testDetails[key] = &reports.TestResult{
						TestName:       entryLine.Test,
						TestPackage:    entryLine.Package,
						PassRatio:      0,
						Outputs:        []string{},
						PackageOutputs: []string{},
					}
				}
				result = testDetails[key]
			}

			// TODO: This is a bit of a logical mess, probably worth a refactor
			if entryLine.Output != "" {
				if panicDetectionMode || raceDetectionMode { // currently collecting panic or race output
					detectedEntries = append(detectedEntries, entryLine)
					continue
				} else if startPanicRe.MatchString(entryLine.Output) { // found a panic, start collecting output
					panickedPackages[entryLine.Package] = struct{}{}
					detectedEntries = append(detectedEntries, entryLine)
					panicDetectionMode = true
					continue // Don't process this entry further
				} else if startRaceRe.MatchString(entryLine.Output) {
					racePackages[entryLine.Package] = struct{}{}
					detectedEntries = append(detectedEntries, entryLine)
					raceDetectionMode = true
					continue // Don't process this entry further
				} else if entryLine.Test == "" {
					if _, exists := packageLevelOutputs[entryLine.Package]; !exists {
						packageLevelOutputs[entryLine.Package] = []string{}
					}
					packageLevelOutputs[entryLine.Package] = append(packageLevelOutputs[entryLine.Package], entryLine.Output)
				} else if entryLine.Test != "" {
					result.Outputs = append(result.Outputs, entryLine.Output)
				}
			}

			if (panicDetectionMode || raceDetectionMode) && entryLine.Action == "fail" { // End of panic or race output
				if panicDetectionMode {
					panicTest, err := attributePanicToTest(entryLine.Package, detectedEntries)
					if err != nil {
						return nil, err
					}
					panicTestKey := fmt.Sprintf("%s/%s", entryLine.Package, panicTest)
					testDetails[panicTestKey].Panicked = true
					testDetails[panicTestKey].Panics++
					testDetails[panicTestKey].Runs++
					// TODO: durations and panics are weird in the same way as Runs: lots of double-counting
					// duration, err := time.ParseDuration(strconv.FormatFloat(entryLine.Elapsed, 'f', -1, 64) + "s")
					// if err != nil {
					// 	return nil, fmt.Errorf("failed to parse duration: %w", err)
					// }
					// testDetails[raceTestKey].Durations = append(testDetails[raceTestKey].Durations, duration)
					testDetails[panicTestKey].Outputs = append(testDetails[panicTestKey].Outputs, entryLine.Output)
					for _, entry := range detectedEntries {
						if entry.Test == "" {
							testDetails[panicTestKey].PackageOutputs = append(testDetails[panicTestKey].PackageOutputs, entry.Output)
						} else {
							testDetails[panicTestKey].Outputs = append(testDetails[panicTestKey].Outputs, entry.Output)
						}
					}
				} else if raceDetectionMode {
					raceTest, err := attributeRaceToTest(entryLine.Package, detectedEntries)
					if err != nil {
						return nil, err
					}
					raceTestKey := fmt.Sprintf("%s/%s", entryLine.Package, raceTest)
					testDetails[raceTestKey].Races++
					testDetails[raceTestKey].Runs++
					// TODO: durations and races are weird in the same way as Runs: lots of double-counting
					// duration, err := time.ParseDuration(strconv.FormatFloat(entryLine.Elapsed, 'f', -1, 64) + "s")
					// if err != nil {
					// 	return nil, fmt.Errorf("failed to parse duration: %w", err)
					// }
					// testDetails[raceTestKey].Durations = append(testDetails[raceTestKey].Durations, duration)
					testDetails[raceTestKey].Outputs = append(testDetails[raceTestKey].Outputs, entryLine.Output)
					for _, entry := range detectedEntries {
						if entry.Test == "" {
							testDetails[raceTestKey].PackageOutputs = append(testDetails[raceTestKey].PackageOutputs, entry.Output)
						} else {
							testDetails[raceTestKey].Outputs = append(testDetails[raceTestKey].Outputs, entry.Output)
						}
					}
				}

				detectedEntries = []entry{}
				panicDetectionMode = false
				raceDetectionMode = false
				continue
			}

			switch entryLine.Action {
			case "pass":
				if entryLine.Test != "" {
					duration, err := time.ParseDuration(strconv.FormatFloat(entryLine.Elapsed, 'f', -1, 64) + "s")
					if err != nil {
						return nil, fmt.Errorf("failed to parse duration: %w", err)
					}
					result.Durations = append(result.Durations, duration)
					result.Successes++
					result.Runs++
				}
			case "fail":
				if entryLine.Test != "" {
					duration, err := time.ParseDuration(strconv.FormatFloat(entryLine.Elapsed, 'f', -1, 64) + "s")
					if err != nil {
						return nil, fmt.Errorf("failed to parse duration: %w", err)
					}
					result.Durations = append(result.Durations, duration)
					result.Failures++
					result.Runs++
				}
			case "skip":
				if entryLine.Test != "" {
					duration, err := time.ParseDuration(strconv.FormatFloat(entryLine.Elapsed, 'f', -1, 64) + "s")
					if err != nil {
						return nil, fmt.Errorf("failed to parse duration: %w", err)
					}
					result.Durations = append(result.Durations, duration)
					result.Skipped = true
					result.Skips++
					result.Runs++
				}
			}
			if entryLine.Test != "" && result.Runs > 0 {
				result.PassRatio = float64(result.Successes) / float64(result.Runs)
				result.PassRatioPercentage = fmt.Sprintf("%.0f%%", result.PassRatio*100)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading test output file: %w", err)
		}
		// Clean up file after parsing
		if err = file.Close(); err != nil {
			log.Printf("WARN: failed to close file: %v", err)
		}
		if err = os.Remove(filePath); err != nil {
			log.Printf("WARN: failed to delete file: %v", err)
		}
	}

	var results []reports.TestResult
	for _, result := range testDetails {
		if _, panicked := panickedPackages[result.TestPackage]; panicked {
			result.PackagePanicked = true
		}
		if outputs, exists := packageLevelOutputs[result.TestPackage]; exists {
			result.PackageOutputs = outputs
		}
		results = append(results, *result)
	}

	return results, nil
}

// properly attributes panics to the test that caused them
// Go JSON output gets confused, especially when tests are run in parallel
func attributePanicToTest(panicPackage string, panicEntries []entry) (string, error) {
	regexSanitizePanicPackage := filepath.Base(panicPackage)
	panicAttributionRe := regexp.MustCompile(fmt.Sprintf(`%s\.(Test.*?)\(.*\)`, regexSanitizePanicPackage))
	for _, entry := range panicEntries {
		if matches := panicAttributionRe.FindStringSubmatch(entry.Output); len(matches) > 1 {
			return matches[1], nil
		}
	}
	return "", fmt.Errorf("failed to attribute panic to test, using regex: %s", panicAttributionRe.String())
}

// properly attributes races to the test that caused them
// Go JSON output gets confused, especially when tests are run in parallel
func attributeRaceToTest(racePackage string, raceEntries []entry) (string, error) {
	regexSanitizeRacePackage := filepath.Base(racePackage)
	raceAttributionRe := regexp.MustCompile(fmt.Sprintf(`%s\.(Test[^\.]+?)\(.*\)`, regexSanitizeRacePackage))
	for _, entry := range raceEntries {
		if matches := raceAttributionRe.FindStringSubmatch(entry.Output); len(matches) > 1 {
			return matches[1], nil
		}
	}
	return "", fmt.Errorf("failed to attribute race to test, using regex: %s", raceAttributionRe.String())
}
