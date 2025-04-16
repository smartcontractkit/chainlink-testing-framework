package parser

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/go-test-transform/pkg/transformer"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
)

// Constants related to parsing and transformation outputs
const (
	// RawOutputTransformedDir defines the directory where transformed output files are stored.
	RawOutputTransformedDir = "./flakeguard_raw_output_transformed"
)

// Parser-specific errors
var (
	// ErrBuild indicates a failure during the test build phase. (Exported)
	ErrBuild = errors.New("failed to build test code")
	// errFailedToShowBuild indicates an error occurred while trying to read build failure details. (Internal)
	errFailedToShowBuild = errors.New("flakeguard failed to show build errors")
)

// Parser-specific regexes
var (
	startPanicRe = regexp.MustCompile(`^panic:`)
	startRaceRe  = regexp.MustCompile(`^WARNING: DATA RACE`)
)

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

type Parser interface {
	// ParseFiles takes a list of raw output file paths, processes them (including potential transformation),
	// and returns the aggregated test results and the list of file paths that were actually parsed.
	ParseFiles(rawFilePaths []string, runPrefix string, expectedRuns int, cfg Config) ([]reports.TestResult, []string, error)
}

type Config struct {
	IgnoreParentFailuresOnSubtests bool
	OmitOutputsOnSuccess           bool
}

type defaultParser struct {
	transformedOutputFiles []string
}

func NewParser() Parser {
	return &defaultParser{
		transformedOutputFiles: make([]string, 0),
	}
}

// ParseFiles is the main entry point for the parser.
func (p *defaultParser) ParseFiles(rawFilePaths []string, runPrefix string, expectedRuns int, cfg Config) ([]reports.TestResult, []string, error) {
	var parseFilePaths = rawFilePaths

	if cfg.IgnoreParentFailuresOnSubtests {
		err := p.transformTestOutputFiles(rawFilePaths)
		if err != nil {
			return nil, nil, fmt.Errorf("failed during output transformation: %w", err)
		}
		parseFilePaths = p.transformedOutputFiles
	}

	results, err := p.parseTestResults(parseFilePaths, runPrefix, expectedRuns, cfg)
	if err != nil {
		return nil, parseFilePaths, err // Return paths even on error?
	}

	return results, parseFilePaths, nil
}

// rawEventData stores the original event along with its run ID.
type rawEventData struct {
	RunID string
	Event entry
}

// testProcessingState holds temporary state while processing events for a single test.
type testProcessingState struct {
	result                  *reports.TestResult // Pointer to the result being built
	processedRunIDs         map[string]bool     // runID -> true if terminal action processed
	runOutcome              map[string]string   // runID -> "pass", "fail", "skip"
	panicRaceOutputByRunID  map[string][]string // runID -> []string of panic/race output
	temporaryOutputsByRunID map[string][]string // runID -> []string of normal output
	panicDetectionMode      bool
	raceDetectionMode       bool
	detectedEntries         []entry // Raw entries collected during panic/race
	key                     string  // Test key (pkg/TestName)
	filePath                string  // File path currently being processed (for logging)
}

// parseTestResults orchestrates the multi-pass parsing approach.
func (p *defaultParser) parseTestResults(parseFilePaths []string, runPrefix string, totalExpectedRunsPerTest int, cfg Config) ([]reports.TestResult, error) {
	eventsByTest, pkgOutputs, subTests, panickedPkgs, racedPkgs, err := p.collectAndGroupEvents(parseFilePaths, runPrefix)
	if err != nil {
		if errors.Is(err, ErrBuild) {
			return nil, err
		}
		return nil, fmt.Errorf("error during event collection: %w", err)
	}

	processedTestDetails, err := p.processEventsPerTest(eventsByTest, cfg)
	if err != nil {
		return nil, fmt.Errorf("error during event processing: %w", err)
	}

	finalResults := p.aggregateAndFinalizeResults(processedTestDetails, subTests, panickedPkgs, racedPkgs, pkgOutputs, totalExpectedRunsPerTest, cfg)

	return finalResults, nil
}

func (p *defaultParser) collectAndGroupEvents(parseFilePaths []string, runPrefix string) (
	eventsByTest map[string][]rawEventData,
	packageLevelOutputs map[string][]string,
	testsWithSubTests map[string][]string,
	panickedPackages map[string]struct{},
	racePackages map[string]struct{},
	err error,
) {
	eventsByTest = make(map[string][]rawEventData)
	packageLevelOutputs = make(map[string][]string)
	testsWithSubTests = make(map[string][]string)
	panickedPackages = make(map[string]struct{})
	racePackages = make(map[string]struct{})

	runNumber := 0
	for _, filePath := range parseFilePaths {
		runNumber++
		runID := fmt.Sprintf("%s%d", runPrefix, runNumber)
		file, fileErr := os.Open(filePath)
		if fileErr != nil {
			err = fmt.Errorf("failed to open test output file '%s': %w", filePath, fileErr)
			return
		}

		scanner := bufio.NewScanner(file)
		parsingErrorOccurred := false
		for scanner.Scan() {
			lineBytes := scanner.Bytes()
			var entryLine entry
			if jsonErr := json.Unmarshal(lineBytes, &entryLine); jsonErr != nil {
				if !parsingErrorOccurred {
					log.Warn().Str("file", filePath).Err(jsonErr).Str("line_content", scanner.Text()).Msg("Failed to parse JSON line, skipping")
					parsingErrorOccurred = true
				}
				continue
			}
			if entryLine.Action == "build-fail" {
				_, seekErr := file.Seek(0, io.SeekStart)
				if seekErr != nil {
					log.Error().Str("file", filePath).Err(seekErr).Msg("Failed to seek to read build errors")
				}
				buildErrs, readErr := io.ReadAll(file)
				if readErr != nil {
					log.Error().Str("file", filePath).Err(readErr).Msg("Failed to read build errors")
				}
				fmt.Fprintf(os.Stderr, "--- Build Error in %s ---\n%s\n-------------------------\n", filePath, string(buildErrs))
				file.Close()
				err = ErrBuild
				return
			}
			if entryLine.Package != "" {
				if startPanicRe.MatchString(entryLine.Output) {
					panickedPackages[entryLine.Package] = struct{}{}
				}
				if startRaceRe.MatchString(entryLine.Output) {
					racePackages[entryLine.Package] = struct{}{}
				}
				if entryLine.Test != "" {
					key := fmt.Sprintf("%s/%s", entryLine.Package, entryLine.Test)
					ev := rawEventData{RunID: runID, Event: entryLine}
					eventsByTest[key] = append(eventsByTest[key], ev)
					parentTestName, subTestName := parseSubTest(entryLine.Test)
					if subTestName != "" {
						parentTestKey := fmt.Sprintf("%s/%s", entryLine.Package, parentTestName)
						if _, ok := testsWithSubTests[parentTestKey]; !ok {
							testsWithSubTests[parentTestKey] = []string{}
						}
						found := false
						for _, st := range testsWithSubTests[parentTestKey] {
							if st == subTestName {
								found = true
								break
							}
						}
						if !found {
							testsWithSubTests[parentTestKey] = append(testsWithSubTests[parentTestKey], subTestName)
						}
					}
				} else if entryLine.Output != "" {
					if _, exists := packageLevelOutputs[entryLine.Package]; !exists {
						packageLevelOutputs[entryLine.Package] = []string{}
					}
					packageLevelOutputs[entryLine.Package] = append(packageLevelOutputs[entryLine.Package], entryLine.Output)
				}
			}
		}
		if scanErr := scanner.Err(); scanErr != nil {
			file.Close()
			err = fmt.Errorf("scanner error reading file '%s': %w", filePath, scanErr)
			return
		}
		file.Close()
	}
	return
}

func (p *defaultParser) processEventsPerTest(eventsByTest map[string][]rawEventData, cfg Config) (map[string]*reports.TestResult, error) {
	processedTestDetails := make(map[string]*reports.TestResult)
	for key, rawEvents := range eventsByTest {
		if len(rawEvents) == 0 {
			continue
		}
		firstEvent := rawEvents[0].Event
		result := &reports.TestResult{
			TestName:       firstEvent.Test,
			TestPackage:    firstEvent.Package,
			PassedOutputs:  make(map[string][]string),
			FailedOutputs:  make(map[string][]string),
			PackageOutputs: make([]string, 0),
			Durations:      make([]time.Duration, 0),
		}
		state := &testProcessingState{
			result:                  result,
			key:                     key,
			processedRunIDs:         make(map[string]bool),
			runOutcome:              make(map[string]string),
			panicRaceOutputByRunID:  make(map[string][]string),
			temporaryOutputsByRunID: make(map[string][]string),
		}

		for _, rawEv := range rawEvents {
			p.processEvent(state, rawEv)
		}

		p.finalizeOutputs(state, cfg)
		result.Runs = len(state.processedRunIDs)
		processedTestDetails[key] = result
	}
	return processedTestDetails, nil
}

// processEvent is the main dispatcher for processing a single event for a test.
func (p *defaultParser) processEvent(state *testProcessingState, rawEv rawEventData) {
	runID := rawEv.RunID
	event := rawEv.Event

	// 1. Handle Output / Panic/Race Start Detection
	if event.Output != "" {
		panicRaceStarted := p.handleOutputEvent(state, event, runID)
		if panicRaceStarted || state.panicDetectionMode || state.raceDetectionMode {
			if state.panicDetectionMode || state.raceDetectionMode {
				state.detectedEntries = append(state.detectedEntries, event)
			}
			return
		}
	}

	// 2. Handle Panic/Race Termination
	p.handlePanicRaceTermination(state, event, runID)

	// 3. Handle Terminal Actions (only if not already processed by panic/race)
	terminalAction := event.Action == "pass" || event.Action == "fail" || event.Action == "skip"
	if terminalAction && !state.processedRunIDs[runID] {
		p.handleTerminalAction(state, event, runID)
	}
}

// handleOutputEvent handles output collection and panic/race start detection.
// Returns true if panic/race mode started.
func (p *defaultParser) handleOutputEvent(state *testProcessingState, event entry, runID string) (panicRaceStarted bool) {
	if state.panicDetectionMode || state.raceDetectionMode {
		return false
	}

	if startPanicRe.MatchString(event.Output) {
		state.detectedEntries = append(state.detectedEntries, event)
		state.panicDetectionMode = true
		return true
	}
	if startRaceRe.MatchString(event.Output) {
		state.detectedEntries = append(state.detectedEntries, event)
		state.raceDetectionMode = true
		return true
	}

	if state.temporaryOutputsByRunID[runID] == nil {
		state.temporaryOutputsByRunID[runID] = []string{}
	}
	state.temporaryOutputsByRunID[runID] = append(state.temporaryOutputsByRunID[runID], event.Output)
	return false
}

// handlePanicRaceTermination processes the end of a panic/race block.
func (p *defaultParser) handlePanicRaceTermination(state *testProcessingState, event entry, runID string) {
	terminalAction := event.Action == "pass" || event.Action == "fail" || event.Action == "skip"
	if !(state.panicDetectionMode || state.raceDetectionMode) || !terminalAction {
		return
	}

	var outputs []string
	for _, de := range state.detectedEntries {
		outputs = append(outputs, de.Output)
	}
	outputStr := strings.Join(outputs, "\n")
	currentPackage := event.Package
	if currentPackage == "" && len(state.detectedEntries) > 0 {
		currentPackage = state.detectedEntries[0].Package
	}

	attributedTestName := event.Test
	var isTimeout bool
	var attrErr error

	if currentPackage == "" {
		log.Error().Str("file", state.filePath).Msg("Cannot attribute panic/race: Package context is missing.")
	} else {
		if state.panicDetectionMode {
			attributedTestName, isTimeout, attrErr = AttributePanicToTest(outputs)
			if attrErr != nil {
				log.Error().Str("test", state.key).Err(attrErr).Str("output", outputStr).Msg("Panic attribution failed")
			}
			state.result.Panic = true
			state.result.Timeout = isTimeout
			if state.panicRaceOutputByRunID[runID] == nil {
				state.panicRaceOutputByRunID[runID] = []string{}
			}
			state.panicRaceOutputByRunID[runID] = append(state.panicRaceOutputByRunID[runID], "--- PANIC DETECTED ---")
			state.panicRaceOutputByRunID[runID] = append(state.panicRaceOutputByRunID[runID], outputs...)
			state.panicRaceOutputByRunID[runID] = append(state.panicRaceOutputByRunID[runID], "--- END PANIC ---")
		} else { // raceDetectionMode
			attributedTestName, attrErr = AttributeRaceToTest(outputs)
			if attrErr != nil {
				log.Warn().Str("test", state.key).Err(attrErr).Str("output", outputStr).Msg("Race attribution failed")
			}
			state.result.Race = true
			if state.panicRaceOutputByRunID[runID] == nil {
				state.panicRaceOutputByRunID[runID] = []string{}
			}
			state.panicRaceOutputByRunID[runID] = append(state.panicRaceOutputByRunID[runID], "--- RACE DETECTED ---")
			state.panicRaceOutputByRunID[runID] = append(state.panicRaceOutputByRunID[runID], outputs...)
			state.panicRaceOutputByRunID[runID] = append(state.panicRaceOutputByRunID[runID], "--- END RACE ---")
		}
		if attributedTestName != state.result.TestName {
			log.Warn().Str("event_test", state.result.TestName).Str("attributed_test", attributedTestName).Msg("Panic/Race attribution mismatch")
		}

		// Mark run as processed (failed) if not already done
		if !state.processedRunIDs[runID] {
			state.result.Failures++
			state.processedRunIDs[runID] = true
			state.runOutcome[runID] = "fail"
		}
	}

	// Reset state
	state.detectedEntries = []entry{}
	state.panicDetectionMode = false
	state.raceDetectionMode = false
}

// handleTerminalAction processes pass/fail/skip actions.
func (p *defaultParser) handleTerminalAction(state *testProcessingState, event entry, runID string) {
	switch event.Action {
	case "pass":
		state.result.Successes++
		state.runOutcome[runID] = "pass"
	case "fail":
		state.result.Failures++
		state.runOutcome[runID] = "fail"
	case "skip":
		state.result.Skips++
		state.result.Skipped = true
		state.runOutcome[runID] = "skip"
		delete(state.temporaryOutputsByRunID, runID)
	}
	state.processedRunIDs[runID] = true

	if event.Action == "pass" || event.Action == "fail" {
		duration, parseErr := time.ParseDuration(strconv.FormatFloat(event.Elapsed, 'f', -1, 64) + "s")
		if parseErr == nil {
			state.result.Durations = append(state.result.Durations, duration)
		} else {
			log.Warn().Str("test", state.key).Float64("elapsed", event.Elapsed).Err(parseErr).Msg("Failed to parse duration")
		}
	}
}

// finalizeOutputs moves collected temporary outputs to the correct final map based on run outcome.
func (p *defaultParser) finalizeOutputs(state *testProcessingState, cfg Config) {
	for runID, outcome := range state.runOutcome {
		normalOutputs := state.temporaryOutputsByRunID[runID]
		panicOrRaceOutputs := state.panicRaceOutputByRunID[runID]

		if outcome == "pass" {
			if !cfg.OmitOutputsOnSuccess {
				if len(normalOutputs) > 0 {
					if state.result.PassedOutputs[runID] == nil {
						state.result.PassedOutputs[runID] = []string{}
					}
					state.result.PassedOutputs[runID] = append(state.result.PassedOutputs[runID], normalOutputs...)
				}
			}
		} else if outcome == "fail" {
			if len(panicOrRaceOutputs) > 0 || len(normalOutputs) > 0 {
				if state.result.FailedOutputs[runID] == nil {
					state.result.FailedOutputs[runID] = []string{}
				}
			}
			if len(panicOrRaceOutputs) > 0 {
				state.result.FailedOutputs[runID] = append(state.result.FailedOutputs[runID], panicOrRaceOutputs...)
			}
			if len(normalOutputs) > 0 {
				state.result.FailedOutputs[runID] = append(state.result.FailedOutputs[runID], normalOutputs...)
			}
			if len(state.result.FailedOutputs[runID]) == 0 {
				state.result.FailedOutputs[runID] = []string{"--- TEST FAILED (no specific output captured) ---"}
			}
		}
	}
}

// parseSubTest checks if a test name is a subtest and returns the parent and sub names.
func parseSubTest(testName string) (parentTestName, subTestName string) {
	parts := strings.SplitN(testName, "/", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

// transformTestOutputFiles transforms the test output JSON files to ignore parent failures when only subtests fail.
func (p *defaultParser) transformTestOutputFiles(filePaths []string) error {
	p.transformedOutputFiles = make([]string, 0, len(filePaths))
	err := os.MkdirAll(RawOutputTransformedDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create transformed output directory '%s': %w", RawOutputTransformedDir, err)
	}
	log.Info().Int("count", len(filePaths)).Msg("Starting transformation of output files")
	for i, origPath := range filePaths {
		inFile, err := os.Open(origPath)
		if err != nil {
			return fmt.Errorf("failed to open original file %s for transformation: %w", origPath, err)
		}

		baseName := filepath.Base(origPath)
		outBaseName := fmt.Sprintf("transformed-%d-%s", i, strings.TrimSuffix(baseName, filepath.Ext(baseName)))
		outPath := filepath.Join(RawOutputTransformedDir, outBaseName+".json")

		outFile, err := os.Create(outPath)
		if err != nil {
			inFile.Close()
			return fmt.Errorf("failed to create transformed file '%s': %w", outPath, err)
		}

		transformErr := transformer.TransformJSON(inFile, outFile, transformer.NewOptions(true))

		closeErrIn := inFile.Close()
		closeErrOut := outFile.Close()

		if transformErr != nil {
			if removeErr := os.Remove(outPath); removeErr != nil {
				log.Warn().Str("file", outPath).Err(removeErr).Msg("Failed to remove incomplete transformed file after error")
			}
			return fmt.Errorf("failed to transform output file %s to %s: %w", origPath, outPath, transformErr)
		}
		if closeErrIn != nil {
			log.Warn().Str("file", origPath).Err(closeErrIn).Msg("Error closing input file after transformation")
		}
		if closeErrOut != nil {
			log.Warn().Str("file", outPath).Err(closeErrOut).Msg("Error closing output file after transformation")
		}

		p.transformedOutputFiles = append(p.transformedOutputFiles, outPath)
	}
	log.Info().Int("count", len(p.transformedOutputFiles)).Msg("Finished transforming output files")
	return nil
}

func (p *defaultParser) aggregateAndFinalizeResults(
	processedTestDetails map[string]*reports.TestResult,
	testsWithSubTests map[string][]string,
	panickedPackages map[string]struct{},
	racePackages map[string]struct{},
	packageLevelOutputs map[string][]string,
	totalExpectedRunsPerTest int,
	cfg Config,
) []reports.TestResult {
	finalResults := make([]reports.TestResult, 0, len(processedTestDetails))

	// Panic Inheritance
	for parentTestKey, subTests := range testsWithSubTests {
		if parentTestResult, exists := processedTestDetails[parentTestKey]; exists {
			if parentTestResult.Panic {
				for _, subTestName := range subTests {
					subTestKey := fmt.Sprintf("%s/%s", parentTestKey, subTestName)
					if subTestResult, subExists := processedTestDetails[subTestKey]; subExists {
						if !subTestResult.Skipped {
							subTestResult.Panic = true
							if subTestResult.Failures == 0 && subTestResult.Successes > 0 {
								log.Warn().Str("subtest", subTestKey).Msg("Marking subtest as failed due to parent panic.")
								subTestResult.Failures += subTestResult.Successes
								subTestResult.Successes = 0
								for runID, outputs := range subTestResult.PassedOutputs {
									if subTestResult.FailedOutputs == nil {
										subTestResult.FailedOutputs = make(map[string][]string)
									}
									subTestResult.FailedOutputs[runID] = append(subTestResult.FailedOutputs[runID], outputs...)
								}
								subTestResult.PassedOutputs = make(map[string][]string)
							}
						}
					}
				}
			}
		}
	}

	// Final Calculation, Correction, Filtering
	for key, result := range processedTestDetails {
		if !result.Skipped && result.Runs > totalExpectedRunsPerTest {
			log.Warn().Str("test", key).Int("actualRuns", result.Runs).Int("expectedRuns", totalExpectedRunsPerTest).Msg("Correcting run count exceeding expected total runs")
			targetRuns := totalExpectedRunsPerTest
			if result.Panic || result.Race {
				newFailures := result.Failures
				if newFailures == 0 {
					newFailures = 1
				}
				if newFailures > targetRuns {
					newFailures = targetRuns
				}
				newSuccesses := targetRuns - newFailures
				if newSuccesses < 0 {
					newSuccesses = 0
				}
				result.Successes = newSuccesses
				result.Failures = newFailures
			} else {
				if result.Runs > 0 {
					newSuccesses := int(float64(result.Successes*targetRuns) / float64(result.Runs))
					newFailures := targetRuns - newSuccesses
					if newFailures < 0 {
						newFailures = 0
						newSuccesses = targetRuns
					}
					result.Successes = newSuccesses
					result.Failures = newFailures
				} else {
					result.Successes = 0
					result.Failures = targetRuns
				}
			}
			result.Runs = targetRuns
		}

		if !result.Skipped {
			if result.Runs > 0 {
				result.PassRatio = float64(result.Successes) / float64(result.Runs)
			} else {
				result.PassRatio = 0.0
			}
		} else {
			result.PassRatio = 1.0
			if result.Runs != 0 {
				result.Runs = 0
			}
		}

		// Apply package-level flags/outputs
		if _, panicked := panickedPackages[result.TestPackage]; panicked && !result.Skipped {
			result.PackagePanic = true
		}
		if outputs, exists := packageLevelOutputs[result.TestPackage]; exists {
			result.PackageOutputs = outputs
		}

		// Filter out results with no runs and not skipped
		if result.Runs > 0 || result.Skipped {
			if cfg.OmitOutputsOnSuccess {
				result.PassedOutputs = make(map[string][]string)
			}
			finalResults = append(finalResults, *result)
		}
	}

	sort.Slice(finalResults, func(i, j int) bool {
		if finalResults[i].TestPackage != finalResults[j].TestPackage {
			return finalResults[i].TestPackage < finalResults[j].TestPackage
		}
		return finalResults[i].TestName < finalResults[j].TestName
	})

	return finalResults
}
