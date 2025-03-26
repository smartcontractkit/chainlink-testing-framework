package runner

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/go-test-transform/pkg/transformer"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/testparser"
)

// Runner describes the test run parameters and raw test outputs
type Runner struct {
	ProjectPath                    string   // Path to the Go project directory.
	Verbose                        bool     // If true, provides detailed logging.
	RunCount                       int      // Number of times to run the tests.
	GoTestCountFlag                *int     // Run go test with -count flag.
	GoTestRaceFlag                 bool     // Run go test with -race flag.
	GoTestTimeoutFlag              string   // Run go test with -timeout flag
	Tags                           []string // Build tags.
	UseShuffle                     bool     // Enable test shuffling. -shuffle=on flag.
	ShuffleSeed                    string   // Set seed for test shuffling -shuffle={seed} flag. Must be used with UseShuffle.
	FailFast                       bool     // Stop on first test failure.
	SkipTests                      []string // Test names to exclude.
	SelectTests                    []string // Test names to include.
	CollectRawOutput               bool     // Set to true to collect test output for later inspection.
	OmitOutputsOnSuccess           bool     // Set to true to omit test outputs on success.
	MaxPassRatio                   float64  // Maximum pass ratio threshold for a test to be considered flaky.
	IgnoreParentFailuresOnSubtests bool     // Ignore failures in parent tests when only subtests fail.
	rawOutputs                     map[string]*bytes.Buffer
}

// executeCommand is our unified helper that runs an exec.Cmd, captures its JSON output
// in a temporary file, and returns that file's path along with a pass/fail status.
func (r *Runner) executeCommand(cmd *exec.Cmd, tempFilePattern, outputKey string) (outputPath string, passed bool, err error) {
	tmpFile, err := os.CreateTemp("", tempFilePattern)
	if err != nil {
		return "", false, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	if r.CollectRawOutput {
		if r.rawOutputs == nil {
			r.rawOutputs = make(map[string]*bytes.Buffer)
		}
		if _, exists := r.rawOutputs[outputKey]; !exists {
			r.rawOutputs[outputKey] = &bytes.Buffer{}
		}
		cmd.Stdout = io.MultiWriter(tmpFile, r.rawOutputs[outputKey])
	} else {
		cmd.Stdout = tmpFile
	}
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	outputPath = tmpFile.Name()

	type exitCoder interface {
		ExitCode() int
	}
	var ec exitCoder
	if errors.As(err, &ec) {
		passed = ec.ExitCode() == 0
		err = nil
		return outputPath, passed, nil
	} else if err != nil {
		return outputPath, false, fmt.Errorf("error running command: %w", err)
	}

	passed = true
	return outputPath, true, nil
}

// RunTestPackages executes the tests for each provided package and aggregates all results.
func (r *Runner) RunTestPackages(packages []string) ([]reports.TestResult, error) {
	var jsonFilePaths []string
	// Initial runs.
	for _, p := range packages {
		for i := 0; i < r.RunCount; i++ {
			if r.CollectRawOutput { // Collect raw output for debugging.
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

	// Pre-process outputs if needed.
	if r.IgnoreParentFailuresOnSubtests {
		transformedPaths, err := r.transformTestOutputFiles(jsonFilePaths)
		if err != nil {
			return nil, err
		}
		jsonFilePaths = transformedPaths
	}

	results, err := testparser.ParseTestResults(jsonFilePaths, "run", r.RunCount, testparser.ParseOptions{
		OmitOutputsOnSuccess: r.OmitOutputsOnSuccess,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse test results: %w", err)
	}

	return results, nil
}

// RunTestCmd runs an arbitrary command testCmd that produces JSON output similar to 'go test -json'.
func (r *Runner) RunTestCmd(testCmd []string) ([]reports.TestResult, error) {
	var jsonOutputPaths []string

	for i := 0; i < r.RunCount; i++ {
		jsonOutputPath, passed, err := r.runCmd(testCmd, i)
		if err != nil {
			return nil, fmt.Errorf("failed to run test command: %w", err)
		}
		jsonOutputPaths = append(jsonOutputPaths, jsonOutputPath)
		if !passed && r.FailFast {
			break
		}
	}

	if r.IgnoreParentFailuresOnSubtests {
		transformedPaths, err := r.transformTestOutputFiles(jsonOutputPaths)
		if err != nil {
			return nil, err
		}
		jsonOutputPaths = transformedPaths
	}

	results, err := testparser.ParseTestResults(jsonOutputPaths, "run", r.RunCount, testparser.ParseOptions{
		OmitOutputsOnSuccess: r.OmitOutputsOnSuccess,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse test results: %w", err)
	}

	return results, nil
}

func (r *Runner) RerunFailedTests(failedTests []reports.TestResult, rerunCount int) ([]reports.TestResult, []string, error) {
	// Group the provided failed tests by package for more efficient reruns.
	failingTestsByPackage := make(map[string][]string)
	for _, tr := range failedTests {
		failingTestsByPackage[tr.TestPackage] = append(failingTestsByPackage[tr.TestPackage], tr.TestName)
	}

	if r.Verbose {
		log.Info().Msgf("Rerunning failing tests grouped by package: %v", failingTestsByPackage)
	}

	var rerunJsonOutputPaths []string

	// Rerun each failing test package up to rerunCount times.
	for i := 0; i < rerunCount; i++ {
		for pkg, tests := range failingTestsByPackage {
			// Build regex pattern to match all failing tests in this package.
			testPattern := fmt.Sprintf("^(%s)$", strings.Join(tests, "|"))

			cmdArgs := []string{
				"test", pkg, "-count=1", "-run", testPattern, "-json",
			}
			if r.GoTestRaceFlag {
				cmdArgs = append(cmdArgs, "-race")
			}
			if r.GoTestTimeoutFlag != "" {
				cmdArgs = append(cmdArgs, fmt.Sprintf("-timeout=%s", r.GoTestTimeoutFlag))
			}
			if len(r.Tags) > 0 {
				cmdArgs = append(cmdArgs, fmt.Sprintf("-tags=%s", strings.Join(r.Tags, ",")))
			}
			if r.Verbose {
				cmdArgs = append(cmdArgs, "-v")
				log.Info().Msgf("Rerun iteration %d for package %s: %v", i+1, pkg, cmdArgs)
			}

			cmd := exec.Command("go", cmdArgs...)
			cmd.Dir = r.ProjectPath

			jsonOutputPath, _, err := r.executeCommand(cmd, "test-output-*.json", pkg)
			if err != nil {
				return nil, nil, fmt.Errorf("error on rerunCmd for package %s: %w", pkg, err)
			}
			rerunJsonOutputPaths = append(rerunJsonOutputPaths, jsonOutputPath)
		}
	}

	if r.IgnoreParentFailuresOnSubtests {
		transformedPaths, err := r.transformTestOutputFiles(rerunJsonOutputPaths)
		if err != nil {
			return nil, rerunJsonOutputPaths, err
		}
		rerunJsonOutputPaths = transformedPaths
	}

	rerunResults, err := testparser.ParseTestResults(rerunJsonOutputPaths, "rerun", rerunCount, testparser.ParseOptions{
		OmitOutputsOnSuccess: r.OmitOutputsOnSuccess,
	})
	if err != nil {
		return nil, rerunJsonOutputPaths, fmt.Errorf("failed to parse rerun results: %w", err)
	}

	return rerunResults, rerunJsonOutputPaths, nil
}

// RawOutputs retrieves the raw output from the test runs, if CollectRawOutput enabled.
func (r *Runner) RawOutputs() map[string]*bytes.Buffer {
	return r.rawOutputs
}

// runTestPackage runs the tests for a given package and returns the path to the output file.
func (r *Runner) runTestPackage(packageName string) (string, bool, error) {
	args := []string{"test", packageName, "-json"}
	if r.GoTestCountFlag != nil {
		args = append(args, fmt.Sprintf("-count=%d", *r.GoTestCountFlag))
	}
	if r.GoTestRaceFlag {
		args = append(args, "-race")
	}
	if r.GoTestTimeoutFlag != "" {
		args = append(args, fmt.Sprintf("-timeout=%s", r.GoTestTimeoutFlag))
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
		selectPattern := strings.Join(r.SelectTests, "$|^")
		args = append(args, fmt.Sprintf("-run=^%s$", selectPattern))
	}

	if r.Verbose {
		log.Info().Str("command", fmt.Sprintf("go %s\n", strings.Join(args, " "))).Msg("Running command")
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = r.ProjectPath

	return r.executeCommand(cmd, "test-output-*.json", packageName)
}

// runCmd runs the user-supplied command once and returns the temp file path, pass status, and any error.
func (r *Runner) runCmd(testCmd []string, runIndex int) (string, bool, error) {
	cmd := exec.Command(testCmd[0], testCmd[1:]...) //nolint:gosec
	cmd.Dir = r.ProjectPath

	outputKey := fmt.Sprintf("customCmd-run%d", runIndex+1)
	tempPattern := fmt.Sprintf("test-output-cmd-run%d-*.json", runIndex+1)

	return r.executeCommand(cmd, tempPattern, outputKey)
}

// transformTestOutputFiles transforms the test output JSON files to ignore parent failures when only subtests fail.
// It returns the paths to the transformed files.
func (r *Runner) transformTestOutputFiles(filePaths []string) ([]string, error) {
	transformedPaths := make([]string, len(filePaths))
	for i, origPath := range filePaths {
		inFile, err := os.Open(origPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open original file %s: %w", origPath, err)
		}
		// Create a temporary file for the transformed output.
		outFile, err := os.CreateTemp("", "transformed-output-*.json")
		if err != nil {
			inFile.Close()
			return nil, fmt.Errorf("failed to create transformed temp file: %w", err)
		}
		// Transform the JSON output.
		err = transformer.TransformJSON(inFile, outFile, transformer.NewOptions(true))
		inFile.Close()
		outFile.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to transform output file %s: %v", origPath, err)
		}
		transformedPaths[i] = outFile.Name()
	}
	return transformedPaths, nil
}
