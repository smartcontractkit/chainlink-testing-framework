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
	ParseFiles(rawFilePaths []string, runPrefix string, expectedRuns int, ignoreParentFailures bool, omitSuccessOutputs bool) ([]reports.TestResult, []string, error)
}

// Config holds configuration relevant to the parser.
type Config struct {
	IgnoreParentFailuresOnSubtests bool
	OmitOutputsOnSuccess           bool
	// Potentially add MaxPassRatio here if calculation moves to parser
}

// defaultParser implements the Parser interface.
type defaultParser struct {
	// config Config // Embed or pass config to methods
	transformedOutputFiles []string // State for transformed files
}

// NewParser creates a new default parser.
func NewParser() Parser {
	return &defaultParser{
		transformedOutputFiles: make([]string, 0),
	}
}

// ParseFiles is the main entry point for the parser.
// It orchestrates transformation (if needed) and parsing of multiple files.
func (p *defaultParser) ParseFiles(rawFilePaths []string, runPrefix string, expectedRuns int, ignoreParentFailures bool, omitSuccessOutputs bool) ([]reports.TestResult, []string, error) {
	var parseFilePaths = rawFilePaths

	// If the option is enabled, transform each JSON output file before parsing.
	if ignoreParentFailures {
		err := p.transformTestOutputFiles(rawFilePaths)
		if err != nil {
			return nil, nil, fmt.Errorf("failed during output transformation: %w", err)
		}
		parseFilePaths = p.transformedOutputFiles
	}

	// Now parse the selected files (raw or transformed)
	results, err := p.parseTestResults(parseFilePaths, runPrefix, expectedRuns, omitSuccessOutputs)
	if err != nil {
		return nil, parseFilePaths, err // Return paths even on error?
	}

	return results, parseFilePaths, nil
}

// parseTestResults reads the test output Go test json output files and returns processed TestResults.
// This is the core logic moved from the original Runner.parseTestResults.
// It now takes file paths directly.
func (p *defaultParser) parseTestResults(parseFilePaths []string, runPrefix string, expectedRuns int, omitSuccessOutputs bool) ([]reports.TestResult, error) {
	var (
		testDetails         = make(map[string]*reports.TestResult) // Holds run, pass counts, and other details for each test
		panickedPackages    = map[string]struct{}{}                // Packages with tests that panicked
		racePackages        = map[string]struct{}{}                // Packages with tests that raced
		packageLevelOutputs = map[string][]string{}                // Package-level outputs
		testsWithSubTests   = map[string][]string{}                // Parent tests that have subtests
		// Note: parseSubTest is now expected to be in the runner package or a shared utils package
	)

	runNumber := 0
	// Process each file
	for _, filePath := range parseFilePaths {
		runNumber++
		runID := fmt.Sprintf("%s%d", runPrefix, runNumber) // Generate RunID based on file index

		// --- Per-file state ---
		panicDetectionMode := false
		raceDetectionMode := false
		detectedEntries := []entry{} // race or panic entries
		// ---------------------

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
			if entryLine.Test != "" {
				// If it's a subtest, associate it with its parent for easier processing of panics later
				key := fmt.Sprintf("%s/%s", entryLine.Package, entryLine.Test)
				// Call the utility function (assuming it's accessible, might need import alias or move)
				// For now, assume it's callable directly if parser stays in runner package temporarily, or needs adjustment.
				// Let's assume we need to call a utility function `parseSubTest`.
				parentTestName, subTestName := parseSubTest(entryLine.Test) // This needs to resolve
				if subTestName != "" {
					parentTestKey := fmt.Sprintf("%s/%s", entryLine.Package, parentTestName)
					// Ensure slice exists before appending
					if _, ok := testsWithSubTests[parentTestKey]; !ok {
						testsWithSubTests[parentTestKey] = make([]string, 0, 1) // Initialize with capacity
					}
					testsWithSubTests[parentTestKey] = append(testsWithSubTests[parentTestKey], subTestName)
				}

				if _, exists := testDetails[key]; !exists {
					// Initialize new test result
					testDetails[key] = &reports.TestResult{
						TestName:       entryLine.Test,
						TestPackage:    entryLine.Package,
						PassRatio:      0,
						PassedOutputs:  make(map[string][]string),
						FailedOutputs:  make(map[string][]string),
						PackageOutputs: make([]string, 0),
						Outputs:        make(map[string][]string),
						Durations:      make([]time.Duration, 0),
					}
				}
				result = testDetails[key]
			}

			if entryLine.Output != "" {
				if panicDetectionMode || raceDetectionMode { // currently collecting panic or race output
					detectedEntries = append(detectedEntries, entryLine)
					continue
				} else if startPanicRe.MatchString(entryLine.Output) { // found a panic, start collecting output
					if entryLine.Package == "" {
						log.Warn().Str("file", filePath).Str("output", entryLine.Output).Msg("Detected panic pattern but package is empty, cannot reliably track package panic state.")
					} else {
						panickedPackages[entryLine.Package] = struct{}{}
					}
					detectedEntries = append(detectedEntries, entryLine)
					panicDetectionMode = true
					continue // Don't process this entry further
				} else if startRaceRe.MatchString(entryLine.Output) {
					if entryLine.Package == "" {
						log.Warn().Str("file", filePath).Str("output", entryLine.Output).Msg("Detected race pattern but package is empty, cannot reliably track package race state.")
					} else {
						racePackages[entryLine.Package] = struct{}{}
					}
					detectedEntries = append(detectedEntries, entryLine)
					raceDetectionMode = true
					continue // Don't process this entry further
				} else if entryLine.Test != "" && entryLine.Action == "output" {
					// Collect outputs temporarily, they will be moved based on pass/fail/skip status
					if result != nil { // Ensure result exists (it should if Test is not empty)
						if result.Outputs == nil {
							result.Outputs = make(map[string][]string)
						}
						result.Outputs[runID] = append(result.Outputs[runID], entryLine.Output)
					} else {
						log.Warn().Str("file", filePath).Str("package", entryLine.Package).Str("test", entryLine.Test).Msg("Received output for test, but test details struct not found.")
					}
				} else if entryLine.Test == "" {
					// Package level output
					if entryLine.Package == "" {
						log.Warn().Str("file", filePath).Str("output", entryLine.Output).Msg("Received package-level output but package name is empty.")
					} else {
						if _, exists := packageLevelOutputs[entryLine.Package]; !exists {
							packageLevelOutputs[entryLine.Package] = []string{}
						}
						packageLevelOutputs[entryLine.Package] = append(packageLevelOutputs[entryLine.Package], entryLine.Output)
					}
				} else {
					// This case should ideally not be hit if the logic above is correct for handling outputs
					log.Warn().Str("file", filePath).Interface("entry", entryLine).Msg("Unhandled output entry type")
				}
			}

			// Check for end of panic/race sequence
			if (panicDetectionMode || raceDetectionMode) && (entryLine.Action == "fail" || entryLine.Action == "pass" || entryLine.Action == "skip") { // End of panic or race output
				if entryLine.Test == "" || entryLine.Package == "" {
					log.Warn().Str("file", filePath).Interface("entry", entryLine).Bool("is_panic", panicDetectionMode).Bool("is_race", raceDetectionMode).Msg("Detected end of panic/race sequence, but entry lacks Test or Package info. Attribution might be incorrect.")
				}

				var outputs []string
				for _, entry := range detectedEntries {
					outputs = append(outputs, entry.Output)
				}
				outputStr := strings.Join(outputs, "\n") // Join for easier logging/error messages

				currentPackage := entryLine.Package // Use package from the fail/pass/skip line as context
				if currentPackage == "" && len(detectedEntries) > 0 {
					currentPackage = detectedEntries[0].Package // Fallback to package from first detected entry
					log.Warn().Str("file", filePath).Str("fallback_package", currentPackage).Msg("Used fallback package for panic/race attribution")
				}

				if currentPackage == "" {
					log.Error().Str("file", filePath).Bool("is_panic", panicDetectionMode).Bool("is_race", raceDetectionMode).Msg("Cannot attribute panic/race: Package context is missing.")
				} else {
					if panicDetectionMode {
						// Call attribution function (now external)
						panicTest, timeout, err := attributePanicToTest(outputs)
						if err != nil {
							log.Error().Str("file", filePath).Str("package", currentPackage).Err(err).Str("output_snippet", outputStr).Msg("Unable to attribute panic to a test")
							panicTest = fmt.Sprintf("UnableToAttributePanicInPackage_%s", currentPackage) // Create a placeholder name
						}
						panicTestKey := fmt.Sprintf("%s/%s", currentPackage, panicTest)

						result, exists := testDetails[panicTestKey]
						if !exists {
							result = &reports.TestResult{
								TestName:       panicTest,
								TestPackage:    currentPackage,
								PassedOutputs:  make(map[string][]string),
								FailedOutputs:  make(map[string][]string),
								PackageOutputs: make([]string, 0),
								Outputs:        make(map[string][]string),
								Durations:      make([]time.Duration, 0),
							}
							testDetails[panicTestKey] = result
						}

						result.Panic = true
						result.Timeout = timeout
						result.Failures++
						result.Runs++

						if result.FailedOutputs == nil {
							result.FailedOutputs = make(map[string][]string)
						}
						result.FailedOutputs[runID] = append(result.FailedOutputs[runID], "--- PANIC DETECTED ---")
						result.FailedOutputs[runID] = append(result.FailedOutputs[runID], outputs...)
						result.FailedOutputs[runID] = append(result.FailedOutputs[runID], "--- END PANIC ---")

					} else if raceDetectionMode {
						// Call attribution function (now external)
						raceTest, err := attributeRaceToTest(outputs)
						if err != nil {
							log.Warn().Str("file", filePath).Str("package", currentPackage).Err(err).Str("output_snippet", outputStr).Msg("Unable to attribute race to a test")
							raceTest = fmt.Sprintf("UnableToAttributeRaceInPackage_%s", currentPackage) // Create placeholder
						}
						raceTestKey := fmt.Sprintf("%s/%s", currentPackage, raceTest)

						result, exists := testDetails[raceTestKey]
						if !exists {
							result = &reports.TestResult{
								TestName:       raceTest,
								TestPackage:    currentPackage,
								PassedOutputs:  make(map[string][]string),
								FailedOutputs:  make(map[string][]string),
								PackageOutputs: make([]string, 0),
								Outputs:        make(map[string][]string),
								Durations:      make([]time.Duration, 0),
							}
							testDetails[raceTestKey] = result
						}

						result.Race = true
						result.Failures++
						result.Runs++

						if result.FailedOutputs == nil {
							result.FailedOutputs = make(map[string][]string)
						}
						result.FailedOutputs[runID] = append(result.FailedOutputs[runID], "--- RACE DETECTED ---")
						result.FailedOutputs[runID] = append(result.FailedOutputs[runID], outputs...)
						result.FailedOutputs[runID] = append(result.FailedOutputs[runID], "--- END RACE ---")
					}
				}
				// Reset detection state
				detectedEntries = []entry{}
				panicDetectionMode = false
				raceDetectionMode = false
				// Continue processing the current 'fail'/'pass'/'skip' entry normally below
			}

			// Process the primary action (pass, fail, skip)
			if result == nil && entryLine.Test != "" {
				key := fmt.Sprintf("%s/%s", entryLine.Package, entryLine.Test)
				log.Warn().Str("key", key).Str("action", entryLine.Action).Msg("Test result struct was nil when processing action, creating it now.")
				result = &reports.TestResult{
					TestName:       entryLine.Test,
					TestPackage:    entryLine.Package,
					PassedOutputs:  make(map[string][]string),
					FailedOutputs:  make(map[string][]string),
					PackageOutputs: make([]string, 0),
					Outputs:        make(map[string][]string),
					Durations:      make([]time.Duration, 0),
				}
				testDetails[key] = result
			}

			if result != nil {
				var duration time.Duration
				var parseErr error
				if entryLine.Action == "pass" || entryLine.Action == "fail" {
					duration, parseErr = time.ParseDuration(strconv.FormatFloat(entryLine.Elapsed, 'f', -1, 64) + "s")
					if parseErr != nil {
						log.Warn().Str("file", filePath).Str("test", entryLine.Test).Float64("elapsed", entryLine.Elapsed).Err(parseErr).Msg("Failed to parse duration from test result")
					} else {
						result.Durations = append(result.Durations, duration)
					}
				}

				switch entryLine.Action {
				case "pass":
					result.Successes++
					if outputs, ok := result.Outputs[runID]; ok {
						if result.PassedOutputs == nil {
							result.PassedOutputs = make(map[string][]string)
						}
						result.PassedOutputs[runID] = outputs
						delete(result.Outputs, runID)
					}
				case "fail":
					_, panicRaceFailure := result.FailedOutputs[runID]
					if !panicRaceFailure {
						result.Failures++
					}
					if outputs, ok := result.Outputs[runID]; ok {
						if result.FailedOutputs == nil {
							result.FailedOutputs = make(map[string][]string)
						}
						result.FailedOutputs[runID] = append(result.FailedOutputs[runID], outputs...)
						delete(result.Outputs, runID)
					} else if !panicRaceFailure {
						if result.FailedOutputs == nil {
							result.FailedOutputs = make(map[string][]string)
						}
						// Ensure entry exists even with no output
						if _, ok := result.FailedOutputs[runID]; !ok {
							result.FailedOutputs[runID] = []string{"--- TEST FAILED (no specific output captured) ---"}
						}
					}
				case "skip":
					result.Skipped = true
					result.Skips++
					delete(result.Outputs, runID)
				case "output":
					// Handled earlier
				}

				// Update Runs count based on Successes and Failures (consistent with original logic for now)
				result.Runs = result.Successes + result.Failures
				if result.Runs > 0 {
					result.PassRatio = float64(result.Successes) / float64(result.Runs)
				} else if result.Skips > 0 {
					result.PassRatio = 1
				} else {
					result.PassRatio = 0 // Or 1? Default to 0 if no runs/skips recorded.
				}
			} // end if result != nil
		} // end scanner.Scan()

		if err := scanner.Err(); err != nil {
			file.Close() // Close file before returning error
			log.Error().Str("file", filePath).Err(err).Msg("Error reading test output file")
			return nil, fmt.Errorf("reading test output file '%s': %w", filePath, err)
		}
		if err = file.Close(); err != nil {
			log.Warn().Err(err).Str("file", filePath).Msg("Failed to close file after processing")
		}
	} // end loop over files

	// --- Post-processing after all files are parsed ---

	var results []reports.TestResult
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

	// Final aggregation and adjustments
	for key, result := range testDetails {
		if !result.Skipped {
			actualRunIDs := make(map[string]struct{})
			for runID := range result.PassedOutputs {
				actualRunIDs[runID] = struct{}{}
			}
			for runID := range result.FailedOutputs {
				actualRunIDs[runID] = struct{}{}
			}
			effectiveRunCount := len(actualRunIDs)

			if result.Runs > expectedRuns || (effectiveRunCount > 0 && result.Runs > effectiveRunCount) {
				targetRuns := expectedRuns
				if effectiveRunCount > 0 && effectiveRunCount < targetRuns {
					targetRuns = effectiveRunCount
				}

				if result.Panic || result.Race {
					if result.Runs != targetRuns {
						log.Warn().Str("test", key).Int("recorded_runs", result.Runs).Int("target_runs", targetRuns).Msg("Adjusting run/failure count for panicked/raced test.")
						result.Failures = targetRuns
						result.Successes = 0
						result.Runs = targetRuns
					}
				} else if result.Runs > targetRuns {
					log.Warn().Str("test", key).Int("recorded_runs", result.Runs).Int("target_runs", targetRuns).Msg("Adjusting run count for test with excessive runs.")
					if result.Runs > 0 {
						result.Successes = int(float64(result.Successes*targetRuns) / float64(result.Runs))
						result.Failures = targetRuns - result.Successes
						// Ensure failures aren't negative if rounding caused issues
						if result.Failures < 0 {
							log.Warn().Str("test", key).Int("success", result.Successes).Int("failures", result.Failures).Msg("Correcting negative failures after scaling")
							result.Successes = targetRuns // Assign all to success if scaling failed badly
							result.Failures = 0
						}
					} else {
						result.Failures = targetRuns
						result.Successes = 0
					}
					result.Runs = targetRuns
				}
			}
			// Recalculate PassRatio after adjustments
			if result.Runs > 0 {
				result.PassRatio = float64(result.Successes) / float64(result.Runs)
			} else {
				result.PassRatio = 1 // Skipped or 0 runs defaults to 100% pass
			}
		} else {
			result.PassRatio = 1 // Ensure skipped tests have PassRatio of 1
		}

		if _, panicked := panickedPackages[result.TestPackage]; panicked {
			if !result.Skipped {
				result.PackagePanic = true
			}
		}
		if _, raced := racePackages[result.TestPackage]; raced {
			if !result.Skipped {
				// TODO: Add PackageRace field to reports.TestResult struct
				// result.PackageRace = true
			}
		}
		if outputs, exists := packageLevelOutputs[result.TestPackage]; exists {
			result.PackageOutputs = outputs
		}

		if omitSuccessOutputs {
			result.PassedOutputs = make(map[string][]string)
			result.Outputs = make(map[string][]string)
		}

		results = append(results, *result)
	}

	return results, nil
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

// parseSubTest needs to be accessible here. It should be moved to a shared utility location
// or this parser needs to be part of the runner package still.
// Assuming it's moved to runner/utils.go and this file becomes part of the runner package
// OR this file becomes parser package and imports runner (circular dependency risk) or utils.
// Let's assume utils for now.
// This function is currently duplicated - needs cleanup.
func parseSubTest(testName string) (parentTestName, subTestName string) {
	parts := strings.SplitN(testName, "/", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
