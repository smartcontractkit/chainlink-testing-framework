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
	// Keeping this consistent with the original runner for now.
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

// entry represents a single line of the go test -json output.
// Moved from runner.go
type entry struct {
	Action  string  `json:"Action"`
	Test    string  `json:"Test"`
	Package string  `json:"Package"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"` // Decimal value in seconds
}

// String provides a string representation of an entry.
// Moved from runner.go
func (e entry) String() string {
	return fmt.Sprintf("Action: %s, Test: %s, Package: %s, Output: %s, Elapsed: %f", e.Action, e.Test, e.Package, e.Output, e.Elapsed)
}

// Parser defines the interface for parsing go test output.
type Parser interface {
	// ParseFiles takes a list of raw output file paths, processes them (including potential transformation),
	// and returns the aggregated test results and the list of file paths that were actually parsed.
	ParseFiles(rawFilePaths []string, runPrefix string, expectedRuns int, cfg Config) ([]reports.TestResult, []string, error)
}

// Config holds configuration relevant to the parser.
type Config struct {
	IgnoreParentFailuresOnSubtests bool
	OmitOutputsOnSuccess           bool
	// Potentially add MaxPassRatio here if calculation moves to parser
}

// defaultParser implements the Parser interface.
type defaultParser struct {
	transformedOutputFiles []string // State for transformed files
}

// NewParser creates a new default parser.
func NewParser() Parser {
	return &defaultParser{
		transformedOutputFiles: make([]string, 0),
	}
}

// ParseFiles is the main entry point for the parser.
func (p *defaultParser) ParseFiles(rawFilePaths []string, runPrefix string, expectedRuns int, cfg Config) ([]reports.TestResult, []string, error) {
	var parseFilePaths = rawFilePaths

	// Use cfg for transformation decision
	if cfg.IgnoreParentFailuresOnSubtests {
		err := p.transformTestOutputFiles(rawFilePaths)
		if err != nil {
			return nil, nil, fmt.Errorf("failed during output transformation: %w", err)
		}
		parseFilePaths = p.transformedOutputFiles
	}

	// Pass cfg down to parseTestResults
	results, err := p.parseTestResults(parseFilePaths, runPrefix, expectedRuns, cfg)
	if err != nil {
		return nil, parseFilePaths, err // Return paths even on error?
	}

	return results, parseFilePaths, nil
}

// parseTestResults reads the test output Go test json output files and returns processed TestResults.
func (p *defaultParser) parseTestResults(parseFilePaths []string, runPrefix string, totalExpectedRunsPerTest int, cfg Config) ([]reports.TestResult, error) {
	var (
		testDetails         = make(map[string]*reports.TestResult) // Holds cumulative results
		panickedPackages    = map[string]struct{}{}
		racePackages        = map[string]struct{}{}
		packageLevelOutputs = map[string][]string{}
		testsWithSubTests   = map[string][]string{}
		processedRunIDs     = make(map[string]map[string]bool) // map[testKey][runID] -> true if terminal action processed
	)

	runNumber := 0
	for _, filePath := range parseFilePaths {
		runNumber++
		runID := fmt.Sprintf("%s%d", runPrefix, runNumber)
		panicDetectionMode := false
		raceDetectionMode := false
		detectedEntries := []entry{}

		file, err := os.Open(filePath)
		if err != nil {
			// Try to provide more context about which file failed
			log.Error().Str("file", filePath).Err(err).Msg("Failed to open test output file")
			// Consider returning partial results or stopping? For now, return error.
			return nil, fmt.Errorf("failed to open test output file '%s': %w", filePath, err)
		}

		scanner := bufio.NewScanner(file)
		var precedingLines []string   // Store preceding lines for context
		parsingErrorOccurred := false // Flag to track if we already logged a JSON parsing error for this file

		for scanner.Scan() {
			line := scanner.Text()
			precedingLines = append(precedingLines, line)

			// Limit precedingLines to the last 15 lines
			if len(precedingLines) > 15 {
				precedingLines = precedingLines[len(precedingLines)-15:] // More efficient slice operation
			}

			var entryLine entry
			if err := json.Unmarshal(scanner.Bytes(), &entryLine); err != nil {
				// Log the error only once per file to avoid spamming
				if !parsingErrorOccurred {
					log.Warn().Str("file", filePath).Err(err).Str("line_content", line).Msg("Failed to parse JSON line in test output, subsequent lines might also fail")
					parsingErrorOccurred = true // Set flag
				}
				continue // Skip processing this line
			}
			if entryLine.Action == "build-fail" {
				// Reset file reader to beginning to capture full build error context
				_, seekErr := file.Seek(0, io.SeekStart) // Use io.SeekStart
				if seekErr != nil {
					log.Error().Str("file", filePath).Err(seekErr).Msg("Failed to seek beginning of file to read build errors")
					file.Close() // Close before returning
					return nil, fmt.Errorf("%w from file '%s': %w", errFailedToShowBuild, filePath, ErrBuild)
				}
				// Print all build errors
				buildErrs, readErr := io.ReadAll(file)
				file.Close() // Close file as we are returning
				if readErr != nil {
					log.Error().Str("file", filePath).Err(readErr).Msg("Failed to read build errors after seeking")
					return nil, fmt.Errorf("%w from file '%s': %w", errFailedToShowBuild, filePath, ErrBuild)
				}
				// Output build errors for user visibility
				fmt.Fprintf(os.Stderr, "--- Build Error in %s ---\n%s\n-------------------------\n", filePath, string(buildErrs))
				return nil, ErrBuild // Return the sentinel buildErr
			}

			var result *reports.TestResult
			key := ""
			if entryLine.Test != "" {
				key = fmt.Sprintf("%s/%s", entryLine.Package, entryLine.Test)
				parentTestName, subTestName := parseSubTest(entryLine.Test)
				if subTestName != "" {
					parentTestKey := fmt.Sprintf("%s/%s", entryLine.Package, parentTestName)
					if _, ok := testsWithSubTests[parentTestKey]; !ok {
						testsWithSubTests[parentTestKey] = []string{}
					}
					// Avoid adding duplicate subtest names
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
				// Initialize result if first time seeing this test
				if _, exists := testDetails[key]; !exists {
					testDetails[key] = &reports.TestResult{
						TestName:       entryLine.Test,
						TestPackage:    entryLine.Package,
						PassedOutputs:  make(map[string][]string),
						FailedOutputs:  make(map[string][]string),
						Outputs:        make(map[string][]string),
						PackageOutputs: make([]string, 0),
						Durations:      make([]time.Duration, 0),
					}
				}
				result = testDetails[key]
				if processedRunIDs[key] == nil {
					processedRunIDs[key] = make(map[string]bool)
				}
			}

			// --- Stage 1: Collect Output / Detect Panic/Race Start ---
			if entryLine.Output != "" {
				if panicDetectionMode || raceDetectionMode {
					detectedEntries = append(detectedEntries, entryLine)
					continue // Don't process output further if collecting for panic/race
				} else if startPanicRe.MatchString(entryLine.Output) {
					if entryLine.Package != "" {
						panickedPackages[entryLine.Package] = struct{}{}
					}
					detectedEntries = append(detectedEntries, entryLine)
					panicDetectionMode = true
					continue
				} else if startRaceRe.MatchString(entryLine.Output) {
					if entryLine.Package != "" {
						racePackages[entryLine.Package] = struct{}{}
					}
					detectedEntries = append(detectedEntries, entryLine)
					raceDetectionMode = true
					continue
				} else if result != nil { // Regular test output
					if result.Outputs[runID] == nil {
						result.Outputs[runID] = []string{}
					}
					result.Outputs[runID] = append(result.Outputs[runID], entryLine.Output)
				} else if entryLine.Package != "" { // Package output
					if _, exists := packageLevelOutputs[entryLine.Package]; !exists {
						packageLevelOutputs[entryLine.Package] = []string{}
					}
					packageLevelOutputs[entryLine.Package] = append(packageLevelOutputs[entryLine.Package], entryLine.Output)
				}
			}

			// --- Stage 2: Process Panic/Race Termination & Attribution ---
			terminalAction := entryLine.Action == "pass" || entryLine.Action == "fail" || entryLine.Action == "skip"
			if (panicDetectionMode || raceDetectionMode) && terminalAction {
				var outputs []string
				for _, entry := range detectedEntries {
					outputs = append(outputs, entry.Output)
				}
				outputStr := strings.Join(outputs, "\n")
				currentPackage := entryLine.Package // Use package from the terminating line
				if currentPackage == "" && len(detectedEntries) > 0 {
					currentPackage = detectedEntries[0].Package
				}

				var attributedTestKey string
				var attributedTestName string
				var isTimeout bool
				var isPanic bool = panicDetectionMode
				var isRace bool = raceDetectionMode

				if currentPackage != "" {
					if isPanic {
						panicTest, timeout, attrErr := AttributePanicToTest(outputs)
						if attrErr != nil {
							log.Error().Str("file", filePath).Str("package", currentPackage).Err(attrErr).Str("output_snippet", outputStr).Msg("Unable to attribute panic")
							panicTest = fmt.Sprintf("UnableToAttributePanicInPackage_%s", currentPackage)
						}
						attributedTestName = panicTest
						isTimeout = timeout
					} else { // isRace
						raceTest, attrErr := AttributeRaceToTest(outputs)
						if attrErr != nil {
							log.Warn().Str("file", filePath).Str("package", currentPackage).Err(attrErr).Str("output_snippet", outputStr).Msg("Unable to attribute race")
							raceTest = fmt.Sprintf("UnableToAttributeRaceInPackage_%s", currentPackage)
						}
						attributedTestName = raceTest
					}
					attributedTestKey = fmt.Sprintf("%s/%s", currentPackage, attributedTestName)
					attrResult, exists := testDetails[attributedTestKey]
					if !exists {
						testDetails[attributedTestKey] = &reports.TestResult{
							TestName:       attributedTestName,
							TestPackage:    currentPackage,
							PassedOutputs:  make(map[string][]string),
							FailedOutputs:  make(map[string][]string),
							Outputs:        make(map[string][]string),
							PackageOutputs: make([]string, 0),
							Durations:      make([]time.Duration, 0),
						}
						attrResult = testDetails[attributedTestKey]
					}
					if processedRunIDs[attributedTestKey] == nil {
						processedRunIDs[attributedTestKey] = make(map[string]bool)
					}
					if attrResult.FailedOutputs == nil {
						attrResult.FailedOutputs = make(map[string][]string)
					}

					attrResult.Panic = attrResult.Panic || isPanic // Persist flags
					attrResult.Race = attrResult.Race || isRace
					attrResult.Timeout = attrResult.Timeout || isTimeout

					// Mark run processed (as failed)
					if !processedRunIDs[attributedTestKey][runID] {
						attrResult.Failures++
						// Do NOT increment attrResult.Runs here, use processedRunIDs length at the end
						processedRunIDs[attributedTestKey][runID] = true
					}
					// Prepend panic/race info to FailedOutputs
					marker := ""
					if isPanic {
						marker = "--- PANIC DETECTED ---"
					}
					if isRace {
						marker = "--- RACE DETECTED ---"
					}
					existingOutput := attrResult.FailedOutputs[runID]
					attrResult.FailedOutputs[runID] = []string{marker}
					attrResult.FailedOutputs[runID] = append(attrResult.FailedOutputs[runID], outputs...)
					if isPanic {
						attrResult.FailedOutputs[runID] = append(attrResult.FailedOutputs[runID], "--- END PANIC ---")
					}
					if isRace {
						attrResult.FailedOutputs[runID] = append(attrResult.FailedOutputs[runID], "--- END RACE ---")
					}
					attrResult.FailedOutputs[runID] = append(attrResult.FailedOutputs[runID], existingOutput...) // Append any previously moved output
				} else {
					log.Error().Str("file", filePath).Msg("Cannot attribute panic/race: Package context is missing.")
				}
				// Reset detection state
				detectedEntries = []entry{}
				panicDetectionMode = false
				raceDetectionMode = false
			}

			// --- Stage 3: Process Terminal Actions (Pass/Fail/Skip) & Move Output ---
			if result != nil && terminalAction {
				processed := processedRunIDs[key][runID] // Re-check if panic/race already processed this runID for this test

				// Record duration first if applicable and not processed
				if (entryLine.Action == "pass" || entryLine.Action == "fail") && !processed {
					duration, parseErr := time.ParseDuration(strconv.FormatFloat(entryLine.Elapsed, 'f', -1, 64) + "s")
					if parseErr == nil {
						result.Durations = append(result.Durations, duration)
					} else { /* log error */
					}
				}

				switch entryLine.Action {
				case "pass":
					if !processed {
						result.Successes++
						processedRunIDs[key][runID] = true
					}
					// Move output AFTER processing state
					if result.PassedOutputs == nil {
						result.PassedOutputs = make(map[string][]string)
					}
					if outputs, ok := result.Outputs[runID]; ok {
						result.PassedOutputs[runID] = append(result.PassedOutputs[runID], outputs...)
						delete(result.Outputs, runID)
					}
				case "fail":
					if !processed {
						result.Failures++
						processedRunIDs[key][runID] = true
					}
					// Move output AFTER processing state
					if result.FailedOutputs == nil {
						result.FailedOutputs = make(map[string][]string)
					}
					if outputs, ok := result.Outputs[runID]; ok {
						result.FailedOutputs[runID] = append(result.FailedOutputs[runID], outputs...)
						delete(result.Outputs, runID)
					}
					// Add placeholder only if no output was moved AND no panic/race marker exists (already checked via !processed)
					if len(result.FailedOutputs[runID]) == 0 {
						result.FailedOutputs[runID] = []string{"--- TEST FAILED (no specific output captured) ---"}
					}
				case "skip":
					if !processed {
						result.Skips++
						processedRunIDs[key][runID] = true
					}
					result.Skipped = true
					delete(result.Outputs, runID) // Discard collected output for skips
				}
			} // end processing terminal action
		} // end scanner loop

		if err := scanner.Err(); err != nil {
			file.Close() // Close file before returning error
			log.Error().Str("file", filePath).Err(err).Msg("Error reading test output file")
			return nil, fmt.Errorf("reading test output file '%s': %w", filePath, err)
		}
		if err = file.Close(); err != nil {
			log.Warn().Err(err).Str("file", filePath).Msg("Failed to close file after processing")
		}
	} // end file loop

	// --- Post-processing --- (Panic Inheritance and Final Aggregation)
	var finalResults []reports.TestResult
	// 1. Panic Inheritance (Bubble down)
	// Bubble panics down from parent tests to subtests
	for parentTestKey, subTests := range testsWithSubTests {
		if parentTestResult, exists := testDetails[parentTestKey]; exists {
			if parentTestResult.Panic {
				for _, subTest := range subTests {
					subTestKey := fmt.Sprintf("%s/%s", parentTestKey, subTest) // Correct key construction
					if subTestResult, exists := testDetails[subTestKey]; exists {
						if !subTestResult.Skipped { // Don't mark skipped subtests as panicked
							subTestResult.Panic = true
							for runID := range subTestResult.FailedOutputs {
								subTestResult.FailedOutputs[runID] = append([]string{"--- ATTRIBUTED PANIC (from parent test) ---"}, subTestResult.FailedOutputs[runID]...)
							}
							if subTestResult.Failures == 0 && subTestResult.Successes > 0 {
								log.Warn().Str("subtest", subTestKey).Msg("Marking subtest as failed due to parent panic.")
								subTestResult.Failures = subTestResult.Successes
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
					} else {
						// log.Debug().Str("expected_subtest", subTestKey).Str("parent_test", parentTestKey).Msg("Expected subtest not found during panic bubbling.")
					}
				}
			}
		} else {
			log.Warn().Str("parent_test", parentTestKey).Msg("Parent test mentioned in subtest key not found in results during panic bubbling.")
		}
	}

	// 2. Final Calculation and Result List Generation
	for key, result := range testDetails {
		// Calculate final Runs based on processed actions for this test
		finalRuns := 0
		if runsMap, ok := processedRunIDs[key]; ok {
			finalRuns = len(runsMap)
		}
		result.Runs = finalRuns // Assign the accurately counted runs

		// Apply Run Count Correction only if necessary
		if !result.Skipped && result.Runs > totalExpectedRunsPerTest {
			log.Warn().Str("test", key).Int("actualRuns", result.Runs).Int("expectedRuns", totalExpectedRunsPerTest).Msg("Correcting run count exceeding expected total runs")
			targetRuns := totalExpectedRunsPerTest
			// Recalculate/cap Successes and Failures based on targetRuns
			if result.Panic || result.Race {
				newFailures := result.Failures
				if newFailures == 0 {
					newFailures = 1
				} // Panic/race is at least 1 failure
				if newFailures > targetRuns {
					newFailures = targetRuns
				}
				newSuccesses := targetRuns - newFailures
				if newSuccesses < 0 {
					newSuccesses = 0
				}
				result.Successes = newSuccesses
				result.Failures = newFailures
			} else { // Scale proportionally
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
			result.Runs = targetRuns // Cap the final run count
		}

		// Final PassRatio calculation
		if !result.Skipped {
			if result.Runs > 0 {
				result.PassRatio = float64(result.Successes) / float64(result.Runs)
			} else {
				result.PassRatio = 0.0 // No runs, not skipped -> 0% pass
			}
		} else {
			result.PassRatio = 1.0 // Skipped -> 100% (or undefined)
			if result.Runs != 0 {
				log.Warn().Str("test", key).Int("runs", result.Runs).Msg("Skipped test has non-zero run count, resetting runs to 0")
				result.Runs = 0 // Skipped tests should have 0 runs
			}
		}

		// Apply package-level flags/outputs
		if _, panicked := panickedPackages[result.TestPackage]; panicked && !result.Skipped {
			result.PackagePanic = true
		}
		if _, raced := racePackages[result.TestPackage]; raced && !result.Skipped {
			// result.PackageRace = true // Uncomment when field exists
		}
		if outputs, exists := packageLevelOutputs[result.TestPackage]; exists {
			result.PackageOutputs = outputs
		}
		if cfg.OmitOutputsOnSuccess {
			result.PassedOutputs = make(map[string][]string)
			result.Outputs = make(map[string][]string)
		}
		// Filter out results with no accurate runs and not skipped
		if result.Runs > 0 || result.Skipped {
			finalResults = append(finalResults, *result)
		}
	}
	return finalResults, nil
}

// transformTestOutputFiles transforms the test output JSON files to ignore parent failures when only subtests fail.
// Moved from runner.go
func (p *defaultParser) transformTestOutputFiles(filePaths []string) error {
	// Clear previous transformed files
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

		// The transformer option is set to ignore parent failures when only subtests fail.
		transformErr := transformer.TransformJSON(inFile, outFile, transformer.NewOptions(true))

		closeErrIn := inFile.Close()
		closeErrOut := outFile.Close()

		if transformErr != nil {
			// Attempt to remove partially written/failed transformed file
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

// parseSubTest checks if a test name is a subtest and returns the parent and sub names.
// Moved back into parser package and kept unexported.
func parseSubTest(testName string) (parentTestName, subTestName string) {
	parts := strings.SplitN(testName, "/", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
