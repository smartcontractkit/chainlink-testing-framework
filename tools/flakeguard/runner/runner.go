package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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

// Runner describes the test run parameters and raw test outputs
type Runner struct {
	ProjectPath          string        // Path to the Go project directory.
	prettyProjectPath    string        // Go project package path, formatted for pretty printing.
	Verbose              bool          // If true, provides detailed logging.
	RunCount             int           // Number of times to run the tests.
	UseRace              bool          // Enable race detector.
	Timeout              time.Duration // Test timeout
	Tags                 []string      // Build tags.
	UseShuffle           bool          // Enable test shuffling. -shuffle=on flag.
	ShuffleSeed          string        // Set seed for test shuffling -shuffle={seed} flag. Must be used with UseShuffle.
	FailFast             bool          // Stop on first test failure.
	SkipTests            []string      // Test names to exclude.
	SelectTests          []string      // Test names to include.
	CollectRawOutput     bool          // Set to true to collect test output for later inspection.
	OmitOutputsOnSuccess bool          // Set to true to omit test outputs on success.
	MaxPassRatio         float64       // Maximum pass ratio threshold for a test to be considered flaky.

	rawOutputs map[string]*bytes.Buffer
}

// RunTestPackages executes the tests for each provided package and aggregates all results.
// It returns all test results and any error encountered during testing.
func (r *Runner) RunTestPackages(packages []string) (*reports.TestReport, error) {
	var jsonFilePaths []string
	for _, p := range packages {
		for i := 0; i < r.RunCount; i++ {
			if r.CollectRawOutput {
				if r.rawOutputs == nil {
					r.rawOutputs = make(map[string]*bytes.Buffer)
				}
				if _, exists := r.rawOutputs[p]; !exists {
					r.rawOutputs[p] = &bytes.Buffer{}
				}
				separator := strings.Repeat("-", 80)
				r.rawOutputs[p].WriteString(fmt.Sprintf("Run %d\n%s\n", i+1, separator))
			}
			jsonFilePath, passed, err := r.runTestPackage(p)
			if err != nil {
				return nil, fmt.Errorf("failed to run tests in package %s: %w", p, err)
			}
			jsonFilePaths = append(jsonFilePaths, jsonFilePath)
			if !passed && r.FailFast {
				break
			}
		}
	}

	results, err := r.parseTestResults(jsonFilePaths)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test results: %w", err)
	}
	report := &reports.TestReport{
		GoProject:     r.prettyProjectPath,
		RaceDetection: r.UseRace,
		ExcludedTests: r.SkipTests,
		SelectedTests: r.SelectTests,
		Results:       results,
		MaxPassRatio:  r.MaxPassRatio,
	}
	report.GenerateSummaryData()

	return report, nil
}

// RunTestCmd runs an arbitrary command testCmd (like ["go", "run", "my_test.go", ...])
// that produces the same JSON lines that 'go test -json' would produce on stdout.
// It captures those lines in a temp file, then parses them for pass/fail/panic/race data.
func (r *Runner) RunTestCmd(testCmd []string) (*reports.TestReport, error) {
	var jsonFilePaths []string

	// Run the command r.RunCount times
	for i := 0; i < r.RunCount; i++ {
		jsonFilePath, passed, err := r.runCmd(testCmd, i)
		if err != nil {
			return nil, fmt.Errorf("failed to run test command: %w", err)
		}
		jsonFilePaths = append(jsonFilePaths, jsonFilePath)
		if !passed && r.FailFast {
			break
		}
	}

	results, err := r.parseTestResults(jsonFilePaths)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test results: %w", err)
	}

	report := &reports.TestReport{
		GoProject:     r.prettyProjectPath,
		RaceDetection: r.UseRace,
		ExcludedTests: r.SkipTests,
		SelectedTests: r.SelectTests,
		Results:       results,
		MaxPassRatio:  r.MaxPassRatio,
	}
	report.GenerateSummaryData()

	return report, nil
}

// RawOutputs retrieves the raw output from the test runs, if CollectRawOutput enabled.
func (r *Runner) RawOutputs() map[string]*bytes.Buffer {
	return r.rawOutputs
}

// runTestPackage runs the tests for a given package and returns the path to the output file.
func (r *Runner) runTestPackage(packageName string) (string, bool, error) {
	args := []string{"test", packageName, "-json", "-count=1"}
	if r.UseRace {
		args = append(args, "-race")
	}
	if r.Timeout > 0 {
		args = append(args, fmt.Sprintf("-timeout=%s", r.Timeout.String()))
	}
	if len(r.Tags) > 0 {
		args = append(args, fmt.Sprintf("-tags=%s", strings.Join(r.Tags, ",")))
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
	if len(r.SelectTests) > 0 {
		// Turn each test into a pattern "TestA$|^TestB$|^TestC$" for -run
		selectPattern := strings.Join(r.SelectTests, "$|^")
		args = append(args, fmt.Sprintf("-run=^%s$", selectPattern))
	}

	if r.Verbose {
		log.Info().Str("command", fmt.Sprintf("go %s\n", strings.Join(args, " "))).Msg("Running command")
	}

	// Create a temporary file to store the output
	tmpFile, err := os.CreateTemp("", "test-output-*.json")
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Prettify the project path (for final reports)
	r.prettyProjectPath, err = prettyProjectPath(r.ProjectPath)
	if err != nil {
		r.prettyProjectPath = r.ProjectPath
		log.Warn().Err(err).
			Str("projectPath", r.ProjectPath).
			Msg("Failed to get pretty project path")
	}

	// Run the command with output directed to the file
	cmd := exec.Command("go", args...)
	cmd.Dir = r.ProjectPath
	if r.CollectRawOutput {
		cmd.Stdout = io.MultiWriter(tmpFile, r.rawOutputs[packageName])
	} else {
		cmd.Stdout = tmpFile
	}
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		// Some errors are just non-zero exit codes
		type exitCoder interface {
			ExitCode() int
		}
		var exErr exitCoder
		if errors.As(err, &exErr) {
			// If code != 0 => test failed
			if exErr.ExitCode() != 0 {
				return tmpFile.Name(), false, nil
			}
			// If exit code is 0, that's unusual with an error; treat it as real error
			return "", false, fmt.Errorf("test command failed at %s: %w", packageName, err)
		}
		// real error
		return "", false, fmt.Errorf("test command at %s gave error: %w", packageName, err)
	}

	return tmpFile.Name(), true, nil // Test succeeded
}

// runCmd runs the user-supplied command once, captures its JSON output,
// and returns the temp file path, whether the test passed, and an error if any.
func (r *Runner) runCmd(testCmd []string, runIndex int) (tempFilePath string, passed bool, err error) {
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("test-output-cmd-run%d-*.json", runIndex+1))
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	if r.Verbose {
		log.Info().
			Msgf("Running custom test command (%d/%d): %s",
				runIndex+1, r.RunCount, strings.Join(testCmd, " "))
	}

	cmd := exec.Command(testCmd[0], testCmd[1:]...)
	cmd.Dir = r.ProjectPath

	// If collecting raw output, write to both file & buffer
	if r.CollectRawOutput {
		if r.rawOutputs == nil {
			r.rawOutputs = make(map[string]*bytes.Buffer)
		}
		key := fmt.Sprintf("customCmd-run%d", runIndex+1)
		if _, exists := r.rawOutputs[key]; !exists {
			r.rawOutputs[key] = &bytes.Buffer{}
		}
		cmd.Stdout = io.MultiWriter(tmpFile, r.rawOutputs[key])
	} else {
		cmd.Stdout = tmpFile
	}
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	tempFilePath = tmpFile.Name()

	// Interpret error/exit code for pass/fail
	type exitCoder interface {
		ExitCode() int
	}
	var ec exitCoder
	if errors.As(err, &ec) {
		// Non-zero exit code => test failure
		if ec.ExitCode() != 0 {
			passed = false
			err = nil // We’ll treat non-zero as not an actual error, but a test fail
			return
		}
		// Zero exit code => pass
		passed = true
		err = nil
		return
	} else if err != nil {
		// Some other error that doesn't implement ExitCode => real error
		return "", false, fmt.Errorf("error running test command: %w", err)
	}

	passed = true
	return
}

// entry is the raw JSON line item we unmarshal from `go test -json`.
type entry struct {
	Action  string  `json:"Action"`
	Test    string  `json:"Test"`
	Package string  `json:"Package"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"` // in seconds
}

// parseTestResults orchestrates reading and parsing multiple JSON output files.
func (r *Runner) parseTestResults(filePaths []string) ([]reports.TestResult, error) {
	// We’ll store test results keyed by "package/test"
	testResultsMap := make(map[string]*reports.TestResult)

	// Track packages that had panics or data races
	panickedPackages := make(map[string]struct{})
	racePackages := make(map[string]struct{})

	// Keep track of package-level outputs
	packageLevelOutputs := make(map[string][]string)

	// Keep track of parent test => subtest names
	testsWithSubTests := make(map[string][]string)

	expectedRuns := r.RunCount
	for i, filePath := range filePaths {
		runNumber := i + 1
		runID := fmt.Sprintf("run%d", runNumber)

		err := r.parseFileLines(
			filePath,
			runNumber,
			testResultsMap,
			panickedPackages,
			racePackages,
			packageLevelOutputs,
			testsWithSubTests,
			runID,
		)
		if err != nil {
			return nil, err
		}
	}

	// Once all files are parsed, fix up parent/child (subtest) relationships:
	zeroOutParentFailsIfSubtestOnlyFails(testResultsMap, testsWithSubTests)

	// Finalize results array
	var results []reports.TestResult
	for _, result := range testResultsMap {
		// If a package had any panic, mark all tests in that package with PackagePanic
		if _, hasPanic := panickedPackages[result.TestPackage]; hasPanic {
			result.PackagePanic = true
		}
		// Omit success outputs if configured
		if r.OmitOutputsOnSuccess {
			result.PassedOutputs = map[string][]string{}
			result.Outputs = map[string][]string{}
		}
		results = append(results, *result)
	}

	// Clean up runs vs. failures if panics introduced double counts
	for i := range results {
		if results[i].Runs > expectedRuns {
			// If it panicked, we typically see double counts. Cap at expectedRuns
			if results[i].Panic {
				results[i].Failures = expectedRuns
				results[i].Runs = expectedRuns
			}
			// If there's some other discrepancy, it's suspicious but we won't do extra logic here
		}
	}

	return results, nil
}

// parseFileLines parses the lines in a single JSON output file, updating shared data structures.
func (r *Runner) parseFileLines(
	filePath string,
	runNumber int,
	testResultsMap map[string]*reports.TestResult,
	panickedPackages map[string]struct{},
	racePackages map[string]struct{},
	packageLevelOutputs map[string][]string,
	testsWithSubTests map[string][]string,
	runID string,
) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open test output file %q: %w", filePath, err)
	}
	defer file.Close()
	defer func() {
		_ = os.Remove(filePath) // Best effort cleanup
	}()

	scanner := bufio.NewScanner(file)

	// We collect panic or race lines until we see a "fail" action
	var (
		collectingPanicOutput bool
		collectingRaceOutput  bool
		detectedEntries       []entry
		precedingLines        []string
		followingLines        []string
	)

	for scanner.Scan() {
		line := scanner.Text()

		// Keep some context lines in case JSON unmarshal fails
		precedingLines = append(precedingLines, line)
		if len(precedingLines) > 15 {
			precedingLines = precedingLines[1:]
		}

		var e entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			// Collect 15 lines after the error for more context
			for scanner.Scan() && len(followingLines) < 15 {
				followingLines = append(followingLines, scanner.Text())
			}
			context := append(precedingLines, followingLines...)
			return fmt.Errorf(
				"failed to parse JSON test output near lines:\n%s\nerror: %w",
				strings.Join(context, "\n"), err,
			)
		}

		// If we are currently collecting panic or race lines
		if collectingPanicOutput || collectingRaceOutput {
			detectedEntries = append(detectedEntries, e)
			// Wait until we see a "fail" action to finalize
			if e.Action == "fail" {
				if collectingPanicOutput {
					panickedPackages[e.Package] = struct{}{}
					err := finishPanicOrRaceCollection(
						detectedEntries, testResultsMap, runID, true,
					)
					if err != nil {
						return err
					}
					collectingPanicOutput = false
				} else {
					racePackages[e.Package] = struct{}{}
					err := finishPanicOrRaceCollection(
						detectedEntries, testResultsMap, runID, false,
					)
					if err != nil {
						return err
					}
					collectingRaceOutput = false
				}
				detectedEntries = nil
			}
			continue
		}

		// Not currently collecting panic/race lines:
		switch {
		case startPanicRe.MatchString(e.Output):
			// We found start of a panic
			panickedPackages[e.Package] = struct{}{}
			collectingPanicOutput = true
			detectedEntries = append(detectedEntries, e)
			continue
		case startRaceRe.MatchString(e.Output):
			// We found start of a data race
			racePackages[e.Package] = struct{}{}
			collectingRaceOutput = true
			detectedEntries = append(detectedEntries, e)
			continue
		}

		// Normal line
		r.handleNormalLine(e, runID, testResultsMap, packageLevelOutputs, testsWithSubTests)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error reading file %s: %w", filePath, err)
	}
	return nil
}

// finishPanicOrRaceCollection finalizes the logic once we detect the "fail" action for a panic or race scenario.
func finishPanicOrRaceCollection(
	detectedEntries []entry,
	testDetails map[string]*reports.TestResult,
	runID string,
	isPanic bool,
) error {
	if isPanic {
		panicTestName, timeout, err := attributePanicToTest(
			detectedEntries[0].Package,
			detectedEntries,
		)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s/%s", detectedEntries[0].Package, panicTestName)

		res := getOrCreateTestResult(testDetails, key)
		res.Panic = true
		res.Timeout = timeout
		res.Failures++
		res.Runs++

		// All lines are failed outputs
		for _, det := range detectedEntries {
			if det.Test == "" {
				res.PackageOutputs = append(res.PackageOutputs, det.Output)
			} else {
				res.FailedOutputs[runID] = append(res.FailedOutputs[runID], det.Output)
			}
		}
	} else {
		// race
		raceTestName, err := attributeRaceToTest(
			detectedEntries[0].Package,
			detectedEntries,
		)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s/%s", detectedEntries[0].Package, raceTestName)

		res := getOrCreateTestResult(testDetails, key)
		res.Race = true
		res.Failures++
		res.Runs++

		for _, det := range detectedEntries {
			if det.Test == "" {
				res.PackageOutputs = append(res.PackageOutputs, det.Output)
			} else {
				res.FailedOutputs[runID] = append(res.FailedOutputs[runID], det.Output)
			}
		}
	}
	return nil
}

// handleNormalLine processes a single line that is not in panic/race mode.
func (r *Runner) handleNormalLine(
	e entry,
	runID string,
	testDetails map[string]*reports.TestResult,
	packageLevelOutputs map[string][]string,
	testsWithSubTests map[string][]string,
) {
	// If it's package-level with no "Test" field, store package-level output
	if e.Test == "" {
		if e.Package != "" && e.Output != "" {
			packageLevelOutputs[e.Package] = append(packageLevelOutputs[e.Package], e.Output)
		}
		return
	}

	// It's a test line
	key := fmt.Sprintf("%s/%s", e.Package, e.Test)
	res := getOrCreateTestResult(testDetails, key)

	// Check for subtest
	parentTestName, subTestName := parseSubTest(e.Test)
	if subTestName != "" {
		parentKey := fmt.Sprintf("%s/%s", e.Package, parentTestName)
		testsWithSubTests[parentKey] = append(testsWithSubTests[parentKey], subTestName)
	}

	// Append the raw line output if Action == "output"
	if e.Action == "output" && e.Output != "" {
		if res.Outputs == nil {
			res.Outputs = make(map[string][]string)
		}
		res.Outputs[runID] = append(res.Outputs[runID], e.Output)
		return
	}

	// If pass/fail/skip, update counters & durations
	switch e.Action {
	case "pass":
		// Mark pass
		dur := parseElapsedDuration(e.Elapsed)
		res.Durations = append(res.Durations, dur)
		res.Successes++
		res.Runs = res.Successes + res.Failures
		res.PassRatio = float64(res.Successes) / float64(res.Runs)

		// Move outputs from "Outputs" to "PassedOutputs"
		if existing := res.Outputs[runID]; len(existing) > 0 {
			if res.PassedOutputs == nil {
				res.PassedOutputs = make(map[string][]string)
			}
			res.PassedOutputs[runID] = existing
			delete(res.Outputs, runID) // Clear ephemeral storage
		}

	case "fail":
		dur := parseElapsedDuration(e.Elapsed)
		res.Durations = append(res.Durations, dur)
		res.Failures++
		res.Runs = res.Successes + res.Failures
		res.PassRatio = float64(res.Successes) / float64(res.Runs)

		// Move outputs from "Outputs" to "FailedOutputs"
		if existing := res.Outputs[runID]; len(existing) > 0 {
			if res.FailedOutputs == nil {
				res.FailedOutputs = make(map[string][]string)
			}
			res.FailedOutputs[runID] = existing
			delete(res.Outputs, runID)
		}

	case "skip":
		res.Skipped = true
		res.Skips++
	}
}

// attributePanicToTest tries to figure out which test caused the panic.
// If it finds "TestFoo" in the stack, returns that test name. Otherwise an error.
// We also detect if it's a "timeout" type panic.
func attributePanicToTest(
	panicPackage string,
	panicEntries []entry,
) (test string, timeout bool, err error) {
	// We'll just look at the base of the package path, since that's what the trace usually shows
	regexSanitize := filepath.Base(panicPackage)
	reTestName := regexp.MustCompile(fmt.Sprintf(`%s\.(Test[^\.\(]+)`, regexSanitize))
	reTimeout := regexp.MustCompile(`(Test.*?)\W+\(.*\)`)

	var lines []string
	for _, line := range panicEntries {
		lines = append(lines, line.Output)
		// Look for something like "mypkg.TestSomething"
		if m := reTestName.FindStringSubmatch(line.Output); len(m) > 1 {
			return strings.TrimSpace(m[1]), false, nil
		}
		// Look for possible test name in parentheses (could be a timeout trace)
		if m := reTimeout.FindStringSubmatch(line.Output); len(m) > 1 {
			return strings.TrimSpace(m[1]), true, nil
		}
	}

	return "", false, fmt.Errorf(
		"failed to attribute panic to test in package %s using regex %q.\nEntries:\n%s",
		panicPackage, reTestName.String(), strings.Join(lines, "\n"),
	)
}

// attributeRaceToTest tries to figure out which test triggered a race.
func attributeRaceToTest(
	racePackage string,
	raceEntries []entry,
) (string, error) {
	regexSanitize := filepath.Base(racePackage)
	reTestName := regexp.MustCompile(fmt.Sprintf(`%s\.(Test[^\.\(]+)`, regexSanitize))

	var lines []string
	for _, line := range raceEntries {
		lines = append(lines, line.Output)
		if m := reTestName.FindStringSubmatch(line.Output); len(m) > 1 {
			return strings.TrimSpace(m[1]), nil
		}
	}
	return "", fmt.Errorf(
		"failed to attribute race to test in package %s using regex %q.\nEntries:\n%s",
		racePackage, reTestName.String(), strings.Join(lines, "\n"),
	)
}

// parseSubTest checks if a test name is a subtest and returns the parent/sub names.
func parseSubTest(testName string) (parentTestName, subTestName string) {
	parts := strings.SplitN(testName, "/", 2)
	if len(parts) < 2 {
		return testName, ""
	}
	return parts[0], parts[1]
}

// zeroOutParentFailsIfSubtestOnlyFails scans through parent tests in testDetails.
// If a parent test has failures *only* referencing subtests (common in Go's default behavior),
// we clear out the parent's "fail" so that only the subtest is marked as failing.
//
// This ensures that a parent test that merely has a subtest fail is not itself flagged
// as a separate test failure, which is what you requested to change.
func zeroOutParentFailsIfSubtestOnlyFails(
	testDetails map[string]*reports.TestResult,
	testsWithSubTests map[string][]string,
) {
	for parentKey, subTests := range testsWithSubTests {
		parentRes, ok := testDetails[parentKey]
		if !ok || parentRes.Failures == 0 {
			continue
		}

		parentHasOwnFailure := false
		for _, failLines := range parentRes.FailedOutputs {
			for _, line := range failLines {
				// If the line doesn't reference the parent at all, skip it.
				if !strings.Contains(line, parentRes.TestName) {
					continue
				}
				// If it's just the summary line (e.g. `--- FAIL: TestParent (0.00s)`),
				// or it references `TestParent/Subtest` or deeper, we skip it
				isSubtestLine := false
				if strings.HasPrefix(line, "--- FAIL: "+parentRes.TestName+" (") {
					isSubtestLine = true
				} else {
					// Check if line references any of the parent's subtests
					for _, sub := range subTests {
						expected := parentRes.TestName + "/" + sub
						if strings.Contains(line, expected) {
							isSubtestLine = true
							break
						}
					}
				}
				if !isSubtestLine {
					parentHasOwnFailure = true
					break
				}
			}
			if parentHasOwnFailure {
				break
			}
		}

		if !parentHasOwnFailure {
			// The parent's only failures are caused by subtests => un-fail the parent
			parentRes.Runs -= parentRes.Failures // remove parent's failures from "runs"
			parentRes.Failures = 0
			if parentRes.Runs < parentRes.Successes {
				parentRes.Successes = parentRes.Runs
			}
			if parentRes.Runs > 0 {
				parentRes.PassRatio = float64(parentRes.Successes) / float64(parentRes.Runs)
			} else {
				// If parent has no runs left, we consider it effectively 100% pass
				parentRes.PassRatio = 1
			}
			// Clear the fail outputs
			parentRes.FailedOutputs = map[string][]string{}
		}
	}
}

// getOrCreateTestResult is a small helper to create/lookup a TestResult struct in the map.
func getOrCreateTestResult(
	testResultsMap map[string]*reports.TestResult,
	key string,
) *reports.TestResult {
	parts := strings.Split(key, "/")
	// Find the first part that looks like a test function.
	var idx = -1
	for i, part := range parts {
		if strings.HasPrefix(part, "Test") {
			idx = i
			break
		}
	}
	// If no part starts with "Test", treat the whole key as the test name.
	if idx == -1 {
		idx = len(parts)
	}
	testPackage := strings.Join(parts[:idx], "/")
	testName := ""
	if idx < len(parts) {
		testName = parts[idx]
		// if idx+1 < len(parts) {
		// subTests = parts[idx+1:]
		// }
	}

	// Create the TestResult if it doesn't exist.
	if _, exists := testResultsMap[key]; !exists {
		testResultsMap[key] = &reports.TestResult{
			TestName:       testName,
			TestPackage:    testPackage,
			PassedOutputs:  make(map[string][]string),
			FailedOutputs:  make(map[string][]string),
			Outputs:        make(map[string][]string),
			PackageOutputs: []string{},
			// If you update the reports.TestResult struct, you could add:
			// SubTests:       subTests,
		}
	}
	return testResultsMap[key]
}

// parseElapsedDuration converts the float "Elapsed" value into a time.Duration safely.
func parseElapsedDuration(elapsedSeconds float64) time.Duration {
	dur, err := time.ParseDuration(strconv.FormatFloat(elapsedSeconds, 'f', -1, 64) + "s")
	if err != nil {
		return 0
	}
	return dur
}

// prettyProjectPath returns the project path formatted for pretty printing in results.
func prettyProjectPath(projectPath string) (string, error) {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}
	dir := absPath

	// Walk upward to find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found in or above %s", projectPath)
		}
		dir = parent
	}

	// Read go.mod to extract the module path
	goModPath := filepath.Join(dir, "go.mod")
	goModData, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}
	moduleLinePrefix := "module "
	for _, line := range strings.Split(string(goModData), "\n") {
		if strings.HasPrefix(line, moduleLinePrefix) {
			goProject := strings.TrimSpace(strings.TrimPrefix(line, moduleLinePrefix))
			relativePath := strings.TrimPrefix(projectPath, dir)
			relativePath = strings.TrimLeft(relativePath, string(os.PathSeparator))
			if relativePath == "" {
				return goProject, nil
			}
			return filepath.ToSlash(filepath.Join(goProject, relativePath)), nil
		}
	}

	return "", fmt.Errorf("module path not found in go.mod")
}
