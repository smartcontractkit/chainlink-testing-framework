package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/utils"
)

const (
	FlakyTestsExitCode = 1
	ErrorExitCode      = 2
	RawOutputDir       = "./flakeguard_raw_output"
)

// runState holds the configuration and results throughout the run process.
type runState struct {
	cfg         *runConfig
	goProject   string
	testRunner  *runner.Runner
	mainResults []reports.TestResult
	mainReport  *reports.TestReport
}

type runConfig struct {
	ProjectPath                    string
	CodeownersPath                 string
	TestPackagesJson               string
	TestPackages                   []string
	TestCmds                       []string
	RunCount                       int
	RerunFailedCount               int
	Tags                           []string
	UseRace                        bool
	MainResultsPath                string
	RerunResultsPath               string
	MinPassRatio                   float64
	SkipTests                      []string
	SelectTests                    []string
	UseShuffle                     bool
	ShuffleSeed                    string
	IgnoreParentFailuresOnSubtests bool
	FailFast                       bool
	GoTestTimeout                  string
	GoTestCount                    *int
}

// outputManager manages the final output buffer and exit code.
type outputManager struct {
	buffer bytes.Buffer
	code   int
}

// flush prints the buffered output and exits with the stored code.
func (o *outputManager) flush() {
	fmt.Print(o.buffer.String())
	os.Exit(o.code)
}

// logErrorAndExit logs an error, sets the exit code to ErrorExitCode, and flushes.
func (o *outputManager) logErrorAndExit(err error, msg string, fields ...map[string]interface{}) {
	l := log.Error().Err(err)
	if len(fields) > 0 {
		l = l.Fields(fields[0])
	}
	l.Msg(msg)
	fmt.Fprintf(&o.buffer, "[ERROR] %s: %v\n", msg, err)
	o.code = ErrorExitCode
	o.flush()
}

// logMsgAndExit logs a message at a specific level, sets the exit code, and flushes.
func (o *outputManager) logMsgAndExit(level zerolog.Level, msg string, code int, fields ...map[string]interface{}) {
	l := log.WithLevel(level)
	if len(fields) > 0 {
		l = l.Fields(fields[0])
	}
	l.Msg(msg)
	fmt.Fprintf(&o.buffer, "[%s] %s\n", level.String(), msg)
	o.code = code
	o.flush()
}

// info logs an informational message to zerolog and the output buffer.
func (o *outputManager) info(step, totalSteps int, msg string) {
	stepMsg := fmt.Sprintf("(%d/%d) %s", step, totalSteps, msg)
	log.Info().Msg(stepMsg)
	fmt.Fprintf(&o.buffer, "\n[INFO] %s\n", stepMsg)
	fmt.Fprintf(&o.buffer, "%s\n", strings.Repeat("-", len(stepMsg)+7))
}

// detail adds a detail line under the current step in the buffer.
func (o *outputManager) detail(msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	fmt.Fprintf(&o.buffer, "  %s\n", formattedMsg)
	log.Debug().Msg(formattedMsg)
}

// finalStatus adds final status messages (ERROR, WARNING, FAIL) to the buffer.
func (o *outputManager) finalStatus(level zerolog.Level, msg string) {
	log.WithLevel(level).Msg(msg)
	levelStr := strings.ToUpper(level.String())
	if level == zerolog.WarnLevel {
		levelStr = "WARN"
	}
	fmt.Fprintf(&o.buffer, "[%s] %s\n", levelStr, msg)
}

var RunTestsCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests potentially multiple times, check for flakiness, and report results.",
	Long: `Runs tests using 'go test -json'.
Can run tests multiple times and rerun failed tests to detect flakiness.
Provides a structured summary of the execution flow and final results,
followed by detailed logs for all executed tests.

Exit Codes:
  0: Success (all tests passed stability requirements)
  1: Flaky tests found or tests failed persistently after reruns
  2: Error during execution (e.g., setup failure, command error)`,
	Run: func(cmd *cobra.Command, args []string) {
		outputMgr := &outputManager{code: 0} // Default success
		state := &runState{}
		var err error

		// Configuration & Setup
		state.cfg, err = parseAndValidateFlags(cmd)
		if err != nil {
			outputMgr.logErrorAndExit(err, "Failed to parse or validate flags")
		}

		state.goProject, err = utils.GetGoProjectName(state.cfg.ProjectPath)
		if err != nil {
			log.Warn().Err(err).Str("projectPath", state.cfg.ProjectPath).Msg("Failed to get pretty project path for report metadata")
		}

		state.testRunner = initializeRunner(state.cfg)

		totalSteps := 3 // Prep, Initial Run, Final Summary
		if state.cfg.RerunFailedCount > 0 {
			totalSteps++ // Add Retry step
		}

		// Preparation
		outputMgr.info(1, totalSteps, "Preparing environment...")
		if err := checkDependencies(state.cfg.ProjectPath); err != nil {
			outputMgr.detail("Warning: Dependency check ('go mod tidy') failed: %v", err)
		} else {
			outputMgr.detail("Dependency check ('go mod tidy'): OK")
		}
		outputMgr.detail("Preparation complete.")

		// Initial Test Run
		outputMgr.info(2, totalSteps, "Running initial tests...")

		var runErr error
		testPackages, determineErr := determineTestPackages(state.cfg)
		if determineErr != nil {
			outputMgr.logErrorAndExit(determineErr, "Failed to determine test packages")
		}

		if len(state.cfg.TestCmds) > 0 {
			outputMgr.detail("Using custom test command(s)...")
			state.mainResults, runErr = state.testRunner.RunTestCmd(state.cfg.TestCmds)
		} else {
			outputMgr.detail("Running test packages: %s", strings.Join(testPackages, ", "))
			state.mainResults, runErr = state.testRunner.RunTestPackages(testPackages)
		}

		if runErr != nil {
			outputMgr.logErrorAndExit(runErr, "Error running initial tests")
		}
		if len(state.mainResults) == 0 {
			outputMgr.logMsgAndExit(zerolog.ErrorLevel, "No tests were run.", ErrorExitCode)
		}

		state.mainReport, err = generateInitialReport(state.mainResults, state.cfg, state.goProject)
		if err != nil {
			outputMgr.logErrorAndExit(err, "Error creating initial test report")
		}

		if state.cfg.MainResultsPath != "" {
			if err := reports.SaveTestResultsToFile(state.mainResults, state.cfg.MainResultsPath); err != nil {
				log.Error().Err(err).Str("path", state.cfg.MainResultsPath).Msg("Error saving main test results to file")
				outputMgr.detail("Warning: Failed to save initial results to %s", state.cfg.MainResultsPath)
			} else {
				log.Info().Str("path", state.cfg.MainResultsPath).Msg("Main test report saved")
				outputMgr.detail("Initial results saved to: %s", state.cfg.MainResultsPath)
			}
		}

		initialPassed, initialFailed, initialSkipped := countResults(state.mainResults)
		totalInitial := len(state.mainResults) - initialSkipped
		outputMgr.detail("Initial run completed:")
		outputMgr.detail("  - %d total tests run (excluding skipped)", totalInitial)
		outputMgr.detail("  - %d passed", initialPassed)
		outputMgr.detail("  - %d failed", initialFailed)
		if initialSkipped > 0 {
			outputMgr.detail("  - %d skipped", initialSkipped)
		}

		initialFailedTests := reports.FilterTests(state.mainResults, func(tr reports.TestResult) bool {
			return !tr.Skipped && tr.Failures > 0
		})

		// Retry Failed Tests
		persistentlyFailingTests := initialFailedTests
		flakyTests := []reports.TestResult{}
		var rerunReport *reports.TestReport

		if state.cfg.RerunFailedCount > 0 && len(initialFailedTests) > 0 {
			outputMgr.info(3, totalSteps, "Retrying failed tests...")

			if handleCmdLineArgsEdgeCase(outputMgr, initialFailedTests, state.cfg) {
				persistentlyFailingTests = initialFailedTests
			} else {
				suffix := fmt.Sprintf(" Rerunning %d failed test(s) up to %d times...", len(initialFailedTests), state.cfg.RerunFailedCount)
				log.Info().Msg(suffix)

				rerunResults, rerunJSONPaths, rerunErr := state.testRunner.RerunFailedTests(initialFailedTests, state.cfg.RerunFailedCount)

				if rerunErr != nil {
					outputMgr.logErrorAndExit(rerunErr, "Error rerunning failed tests")
				}

				rerunReportVal, err := reports.NewTestReport(rerunResults,
					reports.WithGoProject(state.goProject),
					reports.WithCodeOwnersPath(state.cfg.CodeownersPath),
					reports.WithMaxPassRatio(1),
					reports.WithExcludedTests(state.cfg.SkipTests),
					reports.WithSelectedTests(state.cfg.SelectTests),
					reports.WithJSONOutputPaths(rerunJSONPaths),
				)
				if err != nil {
					outputMgr.logErrorAndExit(err, "Error creating rerun test report")
				}
				rerunReport = &rerunReportVal

				if state.cfg.RerunResultsPath != "" && len(rerunResults) > 0 {
					if err := reports.SaveTestResultsToFile(rerunResults, state.cfg.RerunResultsPath); err != nil {
						log.Error().Err(err).Str("path", state.cfg.RerunResultsPath).Msg("Error saving rerun test results to file")
						outputMgr.detail("Warning: Failed to save rerun results to %s", state.cfg.RerunResultsPath)
					} else {
						log.Info().Str("path", state.cfg.RerunResultsPath).Msg("Rerun test report saved")
						outputMgr.detail("Rerun results saved to: %s", state.cfg.RerunResultsPath)
					}
				}

				persistentlyFailingTests = []reports.TestResult{}
				outputMgr.detail("Retry results:")
				for _, result := range rerunResults {
					if !result.Skipped && result.Successes == 0 {
						persistentlyFailingTests = append(persistentlyFailingTests, result)
						outputMgr.detail("  - %s: still FAIL", result.TestName)
					} else if !result.Skipped && result.Successes > 0 && result.Runs > result.Successes {
						flakyTests = append(flakyTests, result)
						outputMgr.detail("  - %s: now PASS (flaky)", result.TestName)
					} else if !result.Skipped {
						outputMgr.detail("  - %s: now PASS", result.TestName)
					}
				}
			}
		} else if len(initialFailedTests) > 0 {
			outputMgr.detail("No reruns configured or no initial failures to retry.")
			if state.cfg.MinPassRatio < 1.0 {
				for _, test := range initialFailedTests {
					if test.PassRatio >= state.cfg.MinPassRatio {
						flakyTests = append(flakyTests, test)
						persistentlyFailingTests = reports.FilterTests(persistentlyFailingTests, func(pt reports.TestResult) bool {
							return !(pt.TestPackage == test.TestPackage && pt.TestName == test.TestName)
						})
					}
				}
			}
		}

		// Final Summary
		finalStepNum := 3
		if state.cfg.RerunFailedCount > 0 {
			finalStepNum = 4
		}
		outputMgr.info(finalStepNum, totalSteps, "Final summary")

		finalFailCount := len(persistentlyFailingTests)
		finalFlakyCount := len(flakyTests)
		finalPassCount := totalInitial - finalFailCount - finalFlakyCount

		outputMgr.detail("Total tests run: %d", totalInitial)
		outputMgr.detail("  - Final PASS: %d", finalPassCount)
		outputMgr.detail("  - Final FAIL: %d", finalFailCount)
		outputMgr.detail("  - FLAKY:      %d", finalFlakyCount)

		fmt.Fprintln(&outputMgr.buffer)

		if finalFailCount > 0 {
			outputMgr.finalStatus(zerolog.ErrorLevel, fmt.Sprintf("%d stable failing test(s) found", finalFailCount))
			outputMgr.code = FlakyTestsExitCode
		}
		if finalFlakyCount > 0 {
			outputMgr.finalStatus(zerolog.WarnLevel, fmt.Sprintf("%d flaky test(s) found", finalFlakyCount))
			if outputMgr.code == 0 {
				outputMgr.code = FlakyTestsExitCode
			}
		}

		if outputMgr.code == 0 {
			outputMgr.finalStatus(zerolog.InfoLevel, "All tests passed stability requirements.")
		}

		if outputMgr.code == FlakyTestsExitCode {
			outputMgr.finalStatus(zerolog.ErrorLevel, fmt.Sprintf("Exit code = %d (failures or flaky tests detected)", outputMgr.code))
		}

		// Detailed Logs

		fmt.Fprintf(&outputMgr.buffer, "\n%s\n", strings.Repeat("=", 60))
		initialRunHeader := fmt.Sprintf("=== DETAILED LOGS FOR INITIAL RUN (Initial run count: %d) ===", state.cfg.RunCount)
		fmt.Fprintf(&outputMgr.buffer, "%s\n", initialRunHeader)
		fmt.Fprintf(&outputMgr.buffer, "%s\n\n", strings.Repeat("=", len(initialRunHeader)))

		reportToLog := state.mainReport
		if reportToLog != nil && len(reportToLog.Results) > 0 {
			err = reportToLog.PrintGotestsumOutput(&outputMgr.buffer, "testname")
			if err != nil {
				log.Error().Err(err).Msg("Error printing initial run gotestsum output")
				fmt.Fprintf(&outputMgr.buffer, "\n[ERROR] Failed to print detailed initial run logs: %v\n", err)
			}
		} else {
			fmt.Fprintf(&outputMgr.buffer, "No test execution data available for initial run logs.\n")
		}

		if rerunReport != nil && len(rerunReport.Results) > 0 {
			retryHeader := fmt.Sprintf("=== DETAILED LOGS FOR RETRY ATTEMPTS (%d retries per test) ===", state.cfg.RerunFailedCount)
			fmt.Fprintf(&outputMgr.buffer, "\n%s\n", strings.Repeat("=", len(retryHeader)))
			fmt.Fprintf(&outputMgr.buffer, "%s\n", retryHeader)
			fmt.Fprintf(&outputMgr.buffer, "%s\n\n", strings.Repeat("=", len(retryHeader)))
			err = rerunReport.PrintGotestsumOutput(&outputMgr.buffer, "testname")
			if err != nil {
				log.Error().Err(err).Msg("Error printing retry gotestsum output")
				fmt.Fprintf(&outputMgr.buffer, "\n[ERROR] Failed to print detailed retry logs: %v\n", err)
			}
		}

		outputMgr.flush()
	},
}

// parseAndValidateFlags parses flags from the cobra command, validates them, and returns a runConfig.
func parseAndValidateFlags(cmd *cobra.Command) (*runConfig, error) {
	cfg := &runConfig{}
	var err error

	cfg.ProjectPath, _ = cmd.Flags().GetString("project-path")
	cfg.CodeownersPath, _ = cmd.Flags().GetString("codeowners-path")
	cfg.TestPackagesJson, _ = cmd.Flags().GetString("test-packages-json")
	cfg.TestPackages, _ = cmd.Flags().GetStringSlice("test-packages")
	cfg.TestCmds, _ = cmd.Flags().GetStringArray("test-cmd")
	cfg.RunCount, _ = cmd.Flags().GetInt("run-count")
	cfg.RerunFailedCount, _ = cmd.Flags().GetInt("rerun-failed-count")
	cfg.Tags, _ = cmd.Flags().GetStringArray("tags")
	cfg.UseRace, _ = cmd.Flags().GetBool("race")
	cfg.MainResultsPath, _ = cmd.Flags().GetString("main-results-path")
	cfg.RerunResultsPath, _ = cmd.Flags().GetString("rerun-results-path")
	cfg.SkipTests, _ = cmd.Flags().GetStringSlice("skip-tests")
	cfg.SelectTests, _ = cmd.Flags().GetStringSlice("select-tests")
	cfg.UseShuffle, _ = cmd.Flags().GetBool("shuffle")
	cfg.ShuffleSeed, _ = cmd.Flags().GetString("shuffle-seed")
	cfg.IgnoreParentFailuresOnSubtests, _ = cmd.Flags().GetBool("ignore-parent-failures-on-subtests")
	cfg.FailFast, _ = cmd.Flags().GetBool("fail-fast")
	cfg.GoTestTimeout, _ = cmd.Flags().GetString("go-test-timeout")

	cfg.ProjectPath, err = utils.ResolveFullPath(cfg.ProjectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve full path for project path '%s': %w", cfg.ProjectPath, err)
	}
	if cfg.MainResultsPath != "" {
		cfg.MainResultsPath, err = utils.ResolveFullPath(cfg.MainResultsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve full path for main results path '%s': %w", cfg.MainResultsPath, err)
		}
	}
	if cfg.RerunResultsPath != "" {
		cfg.RerunResultsPath, err = utils.ResolveFullPath(cfg.RerunResultsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve full path for rerun results path '%s': %w", cfg.RerunResultsPath, err)
		}
	}

	if cmd.Flags().Changed("go-test-count") {
		v, err := cmd.Flags().GetInt("go-test-count")
		if err != nil {
			return nil, fmt.Errorf("error retrieving flag go-test-count: %w", err)
		}
		cfg.GoTestCount = &v
	}

	minPassRatio, _ := cmd.Flags().GetFloat64("min-pass-ratio")
	maxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio")
	maxPassRatioSpecified := cmd.Flags().Changed("max-pass-ratio")

	cfg.MinPassRatio = minPassRatio
	if maxPassRatioSpecified && maxPassRatio != 1.0 {
		log.Warn().Msg("--max-pass-ratio is deprecated, please use --min-pass-ratio instead. Using max-pass-ratio value for now.")
		cfg.MinPassRatio = maxPassRatio // Use the deprecated value if specified
	}

	if cfg.MinPassRatio < 0 || cfg.MinPassRatio > 1 {
		return nil, fmt.Errorf("pass ratio must be between 0 and 1, got: %.2f", cfg.MinPassRatio)
	}

	return cfg, nil
}

// determineTestPackages decides which test packages to run based on the config.
func determineTestPackages(cfg *runConfig) ([]string, error) {
	if len(cfg.TestCmds) > 0 {
		return nil, nil
	}

	var testPackages []string
	if cfg.TestPackagesJson != "" {
		if err := json.Unmarshal([]byte(cfg.TestPackagesJson), &testPackages); err != nil {
			return nil, fmt.Errorf("error decoding test packages JSON: %w", err)
		}
	} else if len(cfg.TestPackages) > 0 {
		testPackages = cfg.TestPackages
	} else {
		return nil, fmt.Errorf("must specify either --test-packages-json, --test-packages, or --test-cmd")
	}
	return testPackages, nil
}

// initializeRunner creates and configures a new test runner.
func initializeRunner(cfg *runConfig) *runner.Runner {
	// Force OmitOutputsOnSuccess to false because we are printing all logs at the end
	omitOutputs := false
	return runner.NewRunner(
		cfg.ProjectPath,
		true,
		cfg.RunCount,
		cfg.GoTestCount,
		cfg.UseRace,
		cfg.GoTestTimeout,
		cfg.Tags,
		cfg.UseShuffle,
		cfg.ShuffleSeed,
		cfg.FailFast,
		cfg.SkipTests,
		cfg.SelectTests,
		cfg.IgnoreParentFailuresOnSubtests,
		omitOutputs,
		RawOutputDir,
		nil, // exec
		nil, // parser
	)
}

// generateInitialReport creates the initial test report from the main run results.
func generateInitialReport(results []reports.TestResult, cfg *runConfig, goProject string) (*reports.TestReport, error) {
	// Get the JSON output paths from the raw output directory
	jsonOutputPaths, err := getJSONOutputPaths(RawOutputDir)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get JSON output paths for initial report")
	}

	reportVal, err := reports.NewTestReport(results,
		reports.WithGoProject(goProject),
		reports.WithCodeOwnersPath(cfg.CodeownersPath),
		reports.WithMaxPassRatio(cfg.MinPassRatio),
		reports.WithGoRaceDetection(cfg.UseRace),
		reports.WithExcludedTests(cfg.SkipTests),
		reports.WithSelectedTests(cfg.SelectTests),
		reports.WithJSONOutputPaths(jsonOutputPaths),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating main test report: %w", err)
	}
	return &reportVal, nil
}

// getJSONOutputPaths returns a list of absolute paths for JSON output files from the given directory.
func getJSONOutputPaths(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var paths []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			absPath, err := filepath.Abs(filepath.Join(dir, file.Name()))
			if err != nil {
				log.Warn().Err(err).Str("file", file.Name()).Msg("Failed to get absolute path for JSON output file")
				continue
			}
			paths = append(paths, absPath)
		}
	}
	return paths, nil
}

// countResults counts the number of passed, failed, and skipped tests.
func countResults(results []reports.TestResult) (passed, failed, skipped int) {
	for _, r := range results {
		if r.Skipped {
			skipped++
		} else if r.Failures == 0 && r.Runs > 0 {
			passed++
		} else if r.Failures > 0 {
			failed++
		}
	}
	return
}

// handleCmdLineArgsEdgeCase checks for and handles the 'go test file.go' edge case.
// Returns true if the edge case was detected and handled, false otherwise.
func handleCmdLineArgsEdgeCase(outputMgr *outputManager, failedTests []reports.TestResult, cfg *runConfig) bool {
	foundCommandLineArgs := false
	if len(cfg.TestCmds) > 0 {
		for _, test := range failedTests {
			if test.TestPackage == "command-line-arguments" {
				foundCommandLineArgs = true
				break
			}
		}
	}

	if foundCommandLineArgs {
		warningMsg := "WARNING: Skipping reruns because 'go test <file.go>' was detected within --test-cmd. " +
			"Flakeguard cannot reliably rerun these tests. " +
			"Final results will be based on the initial run only. " +
			"To enable reruns, use 'go test . -run TestPattern' instead of 'go test <file.go>' within your --test-cmd."
		log.Warn().Msg(warningMsg)
		outputMgr.detail("%s", warningMsg)
		return true
	}
	return false
}

// init sets up the cobra command flags.
func init() {
	RunTestsCmd.Flags().StringP("project-path", "r", ".", "The path to the Go project. Default is the current directory. Useful for subprojects")
	RunTestsCmd.Flags().StringP("codeowners-path", "", "", "Path to the CODEOWNERS file")
	RunTestsCmd.Flags().String("test-packages-json", "", "JSON-encoded string of test packages")
	RunTestsCmd.Flags().StringSlice("test-packages", nil, "Comma-separated list of test packages to run")
	RunTestsCmd.Flags().StringArray("test-cmd", nil,
		"Optional custom test command(s) (e.g. 'go test -json ./... -v'), which must produce 'go test -json' output. "+
			"Avoid 'go test <file.go>' syntax as it prevents reliable reruns. Use 'go test . -run TestName' instead. "+
			"Can be specified multiple times.",
	)
	RunTestsCmd.Flags().StringSlice("skip-tests", nil, "Comma-separated list of test names (regex supported by `go test -skip`) to skip")
	RunTestsCmd.Flags().StringSlice("select-tests", nil, "Comma-separated list of test names (regex supported by `go test -run`) to specifically run")
	RunTestsCmd.Flags().IntP("run-count", "c", 1, "Number of times to run the tests (for main run)")
	RunTestsCmd.Flags().Int("rerun-failed-count", 0, "Number of times to rerun tests that failed the main run (0 disables reruns)")
	RunTestsCmd.Flags().StringArray("tags", nil, "Passed on to the 'go test' command as the -tags flag")
	RunTestsCmd.Flags().String("go-test-timeout", "", "Passed on to the 'go test' command as the -timeout flag (e.g., '30m')")
	RunTestsCmd.Flags().Int("go-test-count", -1, "Passes the '-count' flag directly to 'go test'. Default (-1) omits the flag.")
	RunTestsCmd.Flags().Bool("race", false, "Enable the race detector (-race flag for 'go test')")
	RunTestsCmd.Flags().Bool("shuffle", false, "Enable test shuffling ('go test -shuffle=on')")
	RunTestsCmd.Flags().String("shuffle-seed", "", "Set seed for test shuffling. Requires --shuffle. ('go test -shuffle=on -shuffle.seed=...')")
	RunTestsCmd.Flags().Bool("fail-fast", false, "Stop test execution on the first failure (-failfast flag for 'go test')")
	RunTestsCmd.Flags().String("main-results-path", "", "Path to save the main test results (JSON format)")
	RunTestsCmd.Flags().String("rerun-results-path", "", "Path to save the rerun test results (JSON format)")
	RunTestsCmd.Flags().Bool("omit-test-outputs-on-success", true, "DEPRECATED: No longer used, as all logs are shown at the end.")
	_ = RunTestsCmd.Flags().MarkDeprecated("omit-test-outputs-on-success", "no longer used, as all logs are shown at the end.")
	RunTestsCmd.Flags().Bool("ignore-parent-failures-on-subtests", false, "Ignore failures in parent tests when only subtests fail (affects parsing)")
	RunTestsCmd.Flags().Float64("min-pass-ratio", 1.0, "The minimum pass ratio (0.0-1.0) required for a test in the main run to be considered stable (relevant only if reruns are disabled).")
	RunTestsCmd.Flags().Float64("max-pass-ratio", 1.0, "DEPRECATED: Use --min-pass-ratio instead. This flag will be removed in a future version.")
	_ = RunTestsCmd.Flags().MarkDeprecated("max-pass-ratio", "use --min-pass-ratio instead")
}

// checkDependencies runs 'go mod tidy' to ensure dependencies are correct.
func checkDependencies(projectPath string) error {
	log.Debug().Str("path", projectPath).Msg("Running 'go mod tidy' to check dependencies...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		log.Warn().Err(err).Str("output", out.String()).Msg("Dependency check ('go mod tidy') failed. Continuing execution, but dependencies might be inconsistent.")
		return fmt.Errorf("dependency check ('go mod tidy') failed: %w - %s", err, out.String())
	} else {
		log.Debug().Msg("'go mod tidy' completed successfully.")
	}
	return nil
}
