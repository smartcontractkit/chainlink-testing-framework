package runner

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/executor"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/parser"
)

// Runner describes the test run parameters and manages test execution and result parsing.
// It delegates command execution to an Executor and result parsing to a Parser.
type Runner struct {
	// Configuration fields
	ProjectPath       string
	Verbose           bool
	RunCount          int
	GoTestCountFlag   *int
	GoTestRaceFlag    bool
	GoTestTimeoutFlag string
	Tags              []string
	UseShuffle        bool
	ShuffleSeed       string
	FailFast          bool
	SkipTests         []string
	SelectTests       []string
	RawOutputDir      string

	IgnoreParentFailuresOnSubtests bool
	OmitOutputsOnSuccess           bool

	// Dependencies
	exec   executor.Executor // Injected Executor
	parser parser.Parser     // Injected Parser (interface defined in parser.go)

}

// NewRunner creates a new Runner with the default command executor.
func NewRunner(
	projectPath string,
	verbose bool,
	runCount int,
	goTestCountFlag *int,
	goTestRaceFlag bool,
	goTestTimeoutFlag string,
	tags []string,
	useShuffle bool,
	shuffleSeed string,
	failFast bool,
	skipTests []string,
	selectTests []string,
	ignoreParentFailuresOnSubtests bool,
	omitOutputsOnSuccess bool,
	rawOutputDir string,
	exec executor.Executor,
	p parser.Parser,
) *Runner {
	if exec == nil {
		exec = executor.NewCommandExecutor()
	}
	if p == nil {
		p = parser.NewParser() // Use constructor from parser.go
	}
	return &Runner{
		ProjectPath:                    projectPath,
		Verbose:                        verbose,
		RunCount:                       runCount,
		GoTestCountFlag:                goTestCountFlag,
		GoTestRaceFlag:                 goTestRaceFlag,
		GoTestTimeoutFlag:              goTestTimeoutFlag,
		Tags:                           tags,
		UseShuffle:                     useShuffle,
		ShuffleSeed:                    shuffleSeed,
		FailFast:                       failFast,
		SkipTests:                      skipTests,
		SelectTests:                    selectTests,
		IgnoreParentFailuresOnSubtests: ignoreParentFailuresOnSubtests,
		OmitOutputsOnSuccess:           omitOutputsOnSuccess,
		RawOutputDir:                   rawOutputDir,
		exec:                           exec,
		parser:                         p,
	}
}

// Helper function to create executor.Config from Runner fields
func (r *Runner) getExecutorConfig() executor.Config {
	return executor.Config{
		ProjectPath:       r.ProjectPath,
		Verbose:           r.Verbose,
		GoTestCountFlag:   r.GoTestCountFlag,
		GoTestRaceFlag:    r.GoTestRaceFlag,
		GoTestTimeoutFlag: r.GoTestTimeoutFlag,
		Tags:              r.Tags,
		UseShuffle:        r.UseShuffle,
		ShuffleSeed:       r.ShuffleSeed,
		SkipTests:         r.SkipTests,
		SelectTests:       r.SelectTests,
		RawOutputDir:      r.RawOutputDir,
	}
}

// Helper function to create parser.Config from Runner fields
func (r *Runner) getParserConfig() parser.Config {
	return parser.Config{
		IgnoreParentFailuresOnSubtests: r.IgnoreParentFailuresOnSubtests,
		OmitOutputsOnSuccess:           r.OmitOutputsOnSuccess,
	}
}

// RunTestPackages executes the tests for each provided package and aggregates all results.
func (r *Runner) RunTestPackages(packages []string) ([]reports.TestResult, error) {
	rawOutputFiles := make([]string, 0) // Collect output file paths for this run
	execCfg := r.getExecutorConfig()

	for _, p := range packages {
		for runIdx := 0; runIdx < r.RunCount; runIdx++ {
			// Delegate execution to the executor
			jsonFilePath, passed, err := r.exec.RunTestPackage(execCfg, p, runIdx)
			if err != nil {
				// Handle executor errors (e.g., command not found, setup issues)
				return nil, fmt.Errorf("executor failed for package %s on run %d: %w", p, runIdx, err)
			}
			if jsonFilePath != "" { // Append path even if tests failed (passed == false)
				rawOutputFiles = append(rawOutputFiles, jsonFilePath)
			}
			if !passed && r.FailFast {
				log.Warn().Msgf("FailFast enabled: Stopping run after failure in package %s", p)
				goto ParseResults // Exit outer loop early
			}
		}
	}

ParseResults:
	// Delegate parsing to the parser
	if len(rawOutputFiles) == 0 {
		log.Warn().Msg("No output files were generated, likely due to FailFast or an early error.")
		return []reports.TestResult{}, nil // Return empty results
	}

	log.Info().Int("file_count", len(rawOutputFiles)).Msg("Parsing output files")
	// Create parser config and pass it
	parserCfg := r.getParserConfig()
	// Ignore the returned file paths here, as they aren't used in this flow
	results, _, err := r.parser.ParseFiles(rawOutputFiles, "run", len(rawOutputFiles), parserCfg)
	if err != nil {
		// Check if it's a build error from the parser
		if errors.Is(err, parser.ErrBuild) { // Updated check
			// No extra wrapping needed if buildErr already provides enough context
			return nil, err
		}
		return nil, fmt.Errorf("failed to parse test results: %w", err)
	}

	return results, nil
}

// RunTestCmd runs an arbitrary command testCmd that produces Go test JSON output.
func (r *Runner) RunTestCmd(testCmd []string) ([]reports.TestResult, error) {
	rawOutputFiles := make([]string, 0) // Reset output files for this run
	execCfg := r.getExecutorConfig()

	for i := 0; i < r.RunCount; i++ {
		// Delegate execution to the executor
		jsonOutputPath, passed, err := r.exec.RunCmd(execCfg, testCmd, i)
		if err != nil {
			// Handle executor errors
			return nil, fmt.Errorf("executor failed for custom command on run %d: %w", i, err)
		}
		if jsonOutputPath != "" {
			rawOutputFiles = append(rawOutputFiles, jsonOutputPath)
		}
		if !passed && r.FailFast {
			log.Warn().Msgf("FailFast enabled: Stopping run after custom command failure")
			break // Exit loop early
		}
	}

	// Delegate parsing to the parser
	if len(rawOutputFiles) == 0 {
		log.Warn().Msg("No output files were generated for custom command, likely due to FailFast or an early error.")
		return []reports.TestResult{}, nil
	}

	log.Info().Int("file_count", len(rawOutputFiles)).Msg("Parsing output files from custom command")
	// Create parser config and pass it
	parserCfg := r.getParserConfig()
	// Ignore the returned file paths here as well
	results, _, err := r.parser.ParseFiles(rawOutputFiles, "run", len(rawOutputFiles), parserCfg)
	if err != nil {
		if errors.Is(err, parser.ErrBuild) { // Updated check
			return nil, err
		}
		return nil, fmt.Errorf("failed to parse test results from custom command: %w", err)
	}

	return results, nil
}

// RerunFailedTests reruns specific tests that failed in previous runs using the Executor and Parser.
func (r *Runner) RerunFailedTests(failedTests []reports.TestResult, rerunCount int) ([]reports.TestResult, []string, error) {
	if len(failedTests) == 0 || rerunCount <= 0 {
		log.Info().Msg("No failed tests provided or rerun count is zero. Skipping reruns.")
		return []reports.TestResult{}, []string{}, nil // Nothing to rerun
	}

	// Use a map for efficient lookup and update of currently failing tests
	currentlyFailing := make(map[string]map[string]struct{}) // pkg -> testName -> exists
	for _, tr := range failedTests {
		if tr.TestPackage == "" || tr.TestName == "" {
			log.Warn().Interface("test_result", tr).Msg("Skipping rerun for test result with missing package or name")
			continue
		}
		if _, ok := currentlyFailing[tr.TestPackage]; !ok {
			currentlyFailing[tr.TestPackage] = make(map[string]struct{})
		}
		currentlyFailing[tr.TestPackage][tr.TestName] = struct{}{}
	}

	if len(currentlyFailing) == 0 {
		log.Warn().Msg("No valid failed tests found to rerun after filtering.")
		return []reports.TestResult{}, []string{}, nil
	}

	if r.Verbose {
		log.Info().Int("packages", len(currentlyFailing)).Int("rerun_count", rerunCount).Msg("Starting test reruns for failed tests")
	}

	rerunOutputFiles := make([]string, 0)
	baseExecCfg := r.getExecutorConfig()

	// 2. Iterate Rerun Count
	for i := 0; i < rerunCount; i++ {
		if len(currentlyFailing) == 0 {
			log.Info().Int("iteration", i).Msg("All previously failing tests passed in reruns. Stopping reruns early.")
			break // Stop if no more tests are failing
		}

		if r.Verbose {
			log.Info().Int("iteration", i+1).Int("total", rerunCount).Int("tests_to_rerun", countMapKeys(currentlyFailing)).Msg("Running rerun iteration")
		}

		failingThisIteration := make(map[string]map[string]struct{}) // Track tests still failing *after this iteration*

		// 3. Execute Rerun per Package for currently failing tests
		for pkg, testsMap := range currentlyFailing {
			if len(testsMap) == 0 {
				continue
			}

			testsToRun := make([]string, 0, len(testsMap))
			for testName := range testsMap {
				testsToRun = append(testsToRun, testName)
			}

			// Escape test names for regex and join with |
			escapedTests := make([]string, len(testsToRun))
			for j, testName := range testsToRun {
				escapedTests[j] = regexp.QuoteMeta(testName)
			}
			testPattern := fmt.Sprintf("^(?:%s)$", strings.Join(escapedTests, "|"))

			// Create specific executor config for this rerun invocation
			rerunExecCfg := baseExecCfg // Copy base config
			one := 1
			rerunExecCfg.GoTestCountFlag = &one              // Force -count=1 for rerun
			rerunExecCfg.SelectTests = []string{testPattern} // Target specific tests via -run
			rerunExecCfg.SkipTests = nil                     // Ensure no tests are skipped via -skip

			if r.Verbose {
				log.Info().Str("package", pkg).Str("pattern", testPattern).Int("rerun_iter", i+1).Msg("Executing package rerun")
			}

			jsonOutputPath, passed, err := r.exec.RunTestPackage(rerunExecCfg, pkg, i)
			if err != nil {
				// If execution fails for a package, return the error immediately
				log.Error().Err(err).Str("package", pkg).Int("rerun_iteration", i+1).Msg("Error executing rerun command for package")
				return nil, nil, fmt.Errorf("error on rerun execution for package %s: %w", pkg, err)
			}
			if jsonOutputPath != "" {
				rerunOutputFiles = append(rerunOutputFiles, jsonOutputPath)
			}

			// If the command failed (exit code != 0), keep all tests from this package in the failing list for the next iteration.
			// Otherwise (passed=true), assume tests in this run passed and remove them from the failing list.
			if !passed {
				if _, ok := failingThisIteration[pkg]; !ok {
					failingThisIteration[pkg] = make(map[string]struct{})
				}
				for testName := range testsMap {
					failingThisIteration[pkg][testName] = struct{}{}
				}
			}
		} // end loop over packages for this iteration

		// Update the set of failing tests for the next iteration
		currentlyFailing = failingThisIteration
	} // end loop over rerunCount

	// 4. Parse Rerun Outputs
	if len(rerunOutputFiles) == 0 {
		log.Warn().Msg("No output files were generated during reruns (possibly due to execution errors).")
		return []reports.TestResult{}, []string{}, nil
	}

	log.Info().Int("file_count", len(rerunOutputFiles)).Msg("Parsing rerun output files")
	// Create parser config and pass it
	parserCfg := r.getParserConfig()
	// For parsing reruns, the effective number of runs *per test included in the output* is `rerunCount`.
	// The parser's `expectedRuns` helps adjust for potential overcounting within each file, using `rerunCount` seems correct here.
	rerunResults, parsedFilePaths, err := r.parser.ParseFiles(rerunOutputFiles, "rerun", rerunCount, parserCfg)
	if err != nil {
		// Check for build error specifically?
		if errors.Is(err, parser.ErrBuild) { // Updated check
			log.Error().Err(err).Msg("Build error occurred unexpectedly during test reruns")
			// Fallthrough to return wrapped error
		}
		// Return the file paths even if parsing failed? No, the report wouldn't be useful.
		return nil, nil, fmt.Errorf("failed to parse rerun results: %w", err)
	}

	// 5. Return Results
	log.Info().Int("result_count", len(rerunResults)).Msg("Finished parsing rerun results")
	// Return the parsed results AND the list of files parsed.
	// Note: The function signature needs to change back to return []string
	return rerunResults, parsedFilePaths, nil
}

// Helper function to count keys in the nested map for logging
func countMapKeys(m map[string]map[string]struct{}) int {
	count := 0
	for _, subMap := range m {
		count += len(subMap)
	}
	return count
}
