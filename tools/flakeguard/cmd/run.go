package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
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
	OmitOutputsOnSuccess           bool
	IgnoreParentFailuresOnSubtests bool
	FailFast                       bool
	GoTestTimeout                  string
	GoTestCount                    *int
}

type summaryAndExit struct {
	buffer bytes.Buffer
	code   int
}

func (s *summaryAndExit) flush() {
	fmt.Print(s.buffer.String())
	os.Exit(s.code)
}

func (s *summaryAndExit) logErrorAndExit(err error, msg string, fields ...map[string]interface{}) {
	l := log.Error().Err(err)
	if len(fields) > 0 {
		l = l.Fields(fields[0])
	}
	l.Msg(msg)
	s.code = ErrorExitCode
	s.flush()
}

func (s *summaryAndExit) logMsgAndExit(level zerolog.Level, msg string, code int, fields ...map[string]interface{}) {
	l := log.WithLevel(level)
	if len(fields) > 0 {
		l = l.Fields(fields[0])
	}
	l.Msg(msg)
	s.code = code
	s.flush()
}

var RunTestsCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests to check if they are flaky",
	Run: func(cmd *cobra.Command, args []string) {
		exitHandler := &summaryAndExit{}

		cfg, err := parseAndValidateFlags(cmd)
		if err != nil {
			exitHandler.logErrorAndExit(err, "Failed to parse or validate flags")
		}

		goProject, err := utils.GetGoProjectName(cfg.ProjectPath)
		if err != nil {
			log.Warn().Err(err).Str("projectPath", cfg.ProjectPath).Msg("Failed to get pretty project path")
		}

		if err := checkDependencies(cfg.ProjectPath); err != nil {
			exitHandler.logErrorAndExit(err, "Error checking project dependencies")
		}

		testPackages, err := determineTestPackages(cfg)
		if err != nil {
			exitHandler.logErrorAndExit(err, "Failed to determine test packages")
		}

		testRunner := initializeRunner(cfg)

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Running initial tests..."
		s.Start()

		var mainResults []reports.TestResult
		var runErr error
		if len(cfg.TestCmds) > 0 {
			s.Suffix = " Running custom test command(s)..."
			mainResults, runErr = testRunner.RunTestCmd(cfg.TestCmds)
		} else {
			s.Suffix = " Running test packages..."
			mainResults, runErr = testRunner.RunTestPackages(testPackages)
		}
		s.Stop()

		if runErr != nil {
			exitHandler.logErrorAndExit(runErr, "Error running initial tests")
		}
		if len(mainResults) == 0 {
			exitHandler.logMsgAndExit(zerolog.ErrorLevel, "No tests were run.", ErrorExitCode)
		}

		mainReport, err := generateMainReport(mainResults, cfg, goProject)
		if err != nil {
			exitHandler.logErrorAndExit(err, "Error creating main test report")
		}
		if cfg.MainResultsPath != "" {
			if err := reports.SaveTestResultsToFile(mainResults, cfg.MainResultsPath); err != nil {
				log.Error().Err(err).Str("path", cfg.MainResultsPath).Msg("Error saving main test results to file")
			} else {
				log.Info().Str("path", cfg.MainResultsPath).Msg("Main test report saved")
			}
		}

		if cfg.RerunFailedCount > 0 {
			// Process the initial run, potentially display logs, and decide if reruns are needed.
			failedTests, proceedToRerun := handleInitialRunResults(exitHandler, mainReport, cfg, goProject)
			if !proceedToRerun {
				// handleInitialRunResults will have set the exit code and printed messages.
				exitHandler.flush() // Exit based on initial run results.
				return
			}

			// Execute the reruns and generate the rerun report.
			rerunResults, rerunReport, err := executeAndReportReruns(exitHandler, testRunner, failedTests, cfg, goProject)
			if err != nil {
				// executeAndReportReruns logs the error and sets the exit code.
				exitHandler.flush()
				return
			}

			// Evaluate the final outcome after reruns.
			evaluateRerunOutcome(exitHandler, rerunReport, rerunResults, cfg)

		} else {
			// Handle the case where reruns are disabled.
			handleNoReruns(exitHandler, mainReport, cfg)
		}

		// If we reach here without an early exit, it implies success (or flaky tests handled by handleNoReruns).
		// The exit code will be set by the specific handlers.
		exitHandler.flush()
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
	cfg.OmitOutputsOnSuccess, _ = cmd.Flags().GetBool("omit-test-outputs-on-success")
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
		return nil, nil // Not needed if running custom commands
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
		cfg.OmitOutputsOnSuccess,
		RawOutputDir,
		nil, // exec
		nil, // parser
	)
}

// generateMainReport creates the initial test report from the main run results.
func generateMainReport(results []reports.TestResult, cfg *runConfig, goProject string) (*reports.TestReport, error) {
	// Get the JSON output paths from the raw output directory
	jsonOutputPaths, err := getJSONOutputPaths(RawOutputDir)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get JSON output paths")
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
		return nil, err
	}
	return &reportVal, nil
}

// getJSONOutputPaths returns a list of JSON output files from the given directory
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
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}
	return paths, nil
}

// handleInitialRunResults processes the results of the initial test run.
// It prints summaries and logs for failed tests.
// It returns the list of failed tests and a boolean indicating whether to proceed with reruns.
func handleInitialRunResults(exitHandler *summaryAndExit, mainReport *reports.TestReport, cfg *runConfig, goProject string) ([]reports.TestResult, bool) {
	failedTests := reports.FilterTests(mainReport.Results, func(tr reports.TestResult) bool {
		// Consider a test "failed" initially if it wasn't skipped and didn't pass all its runs.
		return !tr.Skipped && tr.PassRatio < 1.0
	})

	if len(failedTests) == 0 {
		log.Info().Msg("All tests passed the initial run. No tests to rerun.")
		fmt.Fprint(&exitHandler.buffer, "\nFlakeguard Initial Run Summary\n")
		reports.RenderTestReport(&exitHandler.buffer, *mainReport, false, false)
		exitHandler.code = 0 // Success
		return nil, false    // Do not proceed to rerun
	}

	// Print summary of initially failed tests
	fmt.Fprint(&exitHandler.buffer, "\nFailed Tests Summary (Initial Run):\n\n")
	reports.PrintTestResultsTable(&exitHandler.buffer, failedTests, false, false, true, false, false, false)
	fmt.Fprintln(&exitHandler.buffer)

	// Check for the problematic 'go test file.go' case within --test-cmd
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
		// Log warning but don't print logs as context might be unreliable
		warningMsg := "WARNING: Skipping initial failure logs and all reruns because 'go test <file.go>' was detected within --test-cmd. " +
			"Flakeguard cannot reliably determine log context or rerun these tests. " +
			"Results are based on the initial run only. To enable logs and reruns, use 'go test . -run TestPattern' instead of 'go test <file.go>' within your --test-cmd."
		log.Warn().Msg(warningMsg)
		fmt.Fprintf(&exitHandler.buffer, "\n%s\n", warningMsg)
		// Treat this as if reruns were disabled, determining exit code based on MinPassRatio
		handleNoReruns(exitHandler, mainReport, cfg)
		return failedTests, false // Do not proceed to rerun
	}

	// Print logs for initially failed tests (using gotestsum style)
	fmt.Fprint(&exitHandler.buffer, "\nLogs from Initial Run for Failed Tests:\n\n")
	failedOnlyReport, err := reports.NewTestReport(failedTests,
		reports.WithGoProject(goProject), // Include for context
		reports.WithJSONOutputPaths(mainReport.JSONOutputPaths),
		// Potentially add other relevant options if PrintGotestsumOutput needs them
	)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create temporary report for initial failed test logs. Skipping log output.")
	} else {
		// Use "testname" grouping for potentially better focus on individual test failures
		err = failedOnlyReport.PrintGotestsumOutput(&exitHandler.buffer, "testname")
		if err != nil {
			log.Warn().Err(err).Msg("Error printing gotestsum output for initially failed tests")
		}
		fmt.Fprintln(&exitHandler.buffer) // Add a newline for spacing after logs
	}

	log.Info().Int("count", len(failedTests)).Int("rerun_count", cfg.RerunFailedCount).Msg("Proceeding to rerun failed tests...")
	return failedTests, true // Proceed to rerun
}

// executeAndReportReruns performs the rerun of failed tests and generates a report.
func executeAndReportReruns(exitHandler *summaryAndExit, testRunner *runner.Runner, failedTests []reports.TestResult, cfg *runConfig, goProject string) ([]reports.TestResult, *reports.TestReport, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Rerunning failed tests..."
	s.Start()
	rerunResults, rerunJsonOutputPaths, err := testRunner.RerunFailedTests(failedTests, cfg.RerunFailedCount)
	s.Stop()

	if err != nil {
		exitHandler.logErrorAndExit(err, "Error rerunning failed tests")
		return nil, nil, err // Return error to signal failure
	}

	rerunReport, err := reports.NewTestReport(rerunResults,
		reports.WithGoProject(goProject),
		reports.WithCodeOwnersPath(cfg.CodeownersPath),
		reports.WithMaxPassRatio(1), // Pass ratio for reruns is effectively 1 (did it pass at least once?)
		reports.WithExcludedTests(cfg.SkipTests),
		reports.WithSelectedTests(cfg.SelectTests),
		reports.WithJSONOutputPaths(rerunJsonOutputPaths),
	)
	if err != nil {
		exitHandler.logErrorAndExit(err, "Error creating rerun test report")
		return rerunResults, nil, err // Return error
	}

	fmt.Fprint(&exitHandler.buffer, "Tests After Rerun:")
	reports.PrintTestResultsTable(&exitHandler.buffer, rerunResults, false, false, true, true, true, true)
	fmt.Fprintln(&exitHandler.buffer)

	// Save the rerun test report to file
	if cfg.RerunResultsPath != "" && len(rerunResults) > 0 {
		if err := reports.SaveTestResultsToFile(rerunResults, cfg.RerunResultsPath); err != nil {
			log.Error().Err(err).Str("path", cfg.RerunResultsPath).Msg("Error saving rerun test results to file")
			// Don't treat this as a fatal error for the overall run
		} else {
			log.Info().Str("path", cfg.RerunResultsPath).Msg("Rerun test report saved")
		}
	}

	return rerunResults, &rerunReport, nil // Success
}

// evaluateRerunOutcome determines the final exit code based on persistently failing tests after reruns.
func evaluateRerunOutcome(exitHandler *summaryAndExit, rerunReport *reports.TestReport, rerunResults []reports.TestResult, cfg *runConfig) {
	// Filter tests that still failed after reruns (0 successes)
	failedAfterRerun := reports.FilterTests(rerunResults, func(tr reports.TestResult) bool {
		return !tr.Skipped && tr.Successes == 0
	})

	if len(failedAfterRerun) > 0 {
		fmt.Fprint(&exitHandler.buffer, "Persistently Failing Test Logs (After Reruns):")
		err := rerunReport.PrintGotestsumOutput(&exitHandler.buffer, "testname") // Use testname grouping
		if err != nil {
			log.Error().Err(err).Msg("Error printing gotestsum output for persistently failing tests")
		}

		exitHandler.logMsgAndExit(zerolog.ErrorLevel, "Some tests are still failing after multiple reruns with no successful attempts.", ErrorExitCode, map[string]interface{}{
			"persistently_failing_count": len(failedAfterRerun),
			"rerun_attempts":             cfg.RerunFailedCount,
		})
	} else {
		log.Info().Msg("All initially failing tests passed at least once after reruns.")
		exitHandler.code = 0 // Success
		// No need to call flush here, the main function will do it.
	}
}

// handleNoReruns determines the outcome when reruns are disabled or skipped.
func handleNoReruns(exitHandler *summaryAndExit, mainReport *reports.TestReport, cfg *runConfig) {
	flakyTests := reports.FilterTests(mainReport.Results, func(tr reports.TestResult) bool {
		// A test is flaky if it wasn't skipped and its pass ratio is below the threshold.
		return !tr.Skipped && tr.PassRatio < cfg.MinPassRatio
	})

	// Print the final summary report only if it wasn't printed already (e.g., in the no-initial-failures case)
	// We check if the buffer already contains the "Flakeguard Initial Run Summary" header
	if !bytes.Contains(exitHandler.buffer.Bytes(), []byte("Flakeguard Initial Run Summary")) {
		fmt.Fprint(&exitHandler.buffer, "Flakeguard Summary (No Reruns Performed)")
		reports.RenderTestReport(&exitHandler.buffer, *mainReport, false, false)
	}

	if len(flakyTests) > 0 {
		// Create a new report with only flaky tests to get their gotestsum output
		// We need to ensure JSON paths are associated correctly for log retrieval.
		flakyReport, err := reports.NewTestReport(flakyTests,
			reports.WithJSONOutputPaths(mainReport.JSONOutputPaths),
			// Add other necessary options if PrintGotestsumOutput depends on them
		)
		if err != nil {
			log.Error().Err(err).Msg("Error creating flaky tests report for log printing")
		} else {
			fmt.Fprint(&exitHandler.buffer, "Flaky Test Logs (Based on Initial Run):")
			err := flakyReport.PrintGotestsumOutput(&exitHandler.buffer, "testname") // Use testname grouping
			if err != nil {
				log.Error().Err(err).Msg("Error printing gotestsum output for flaky tests")
			}
		}

		exitHandler.logMsgAndExit(zerolog.InfoLevel, "Found flaky tests based on initial run.", FlakyTestsExitCode, map[string]interface{}{
			"flaky_count":         len(flakyTests),
			"stability_threshold": fmt.Sprintf("%.0f%%", cfg.MinPassRatio*100),
		})
	} else {
		// If no tests were flaky according to the threshold, it's a success.
		// The exit code might already be 0 if set by handleInitialRunResults.
		if exitHandler.code != 0 { // Only log success if not already handled (e.g., by the command-line-args warning path)
			log.Info().Msg("All tests passed stability requirements based on the initial run.")
			exitHandler.code = 0
		}
		// No need to call flush here, the main function will do it.
	}
}

// init sets up the cobra command flags.
func init() {
	RunTestsCmd.Flags().StringP("project-path", "r", ".", "The path to the Go project. Default is the current directory. Useful for subprojects")
	RunTestsCmd.Flags().StringP("codeowners-path", "", "", "Path to the CODEOWNERS file")
	RunTestsCmd.Flags().String("test-packages-json", "", "JSON-encoded string of test packages")
	RunTestsCmd.Flags().StringSlice("test-packages", nil, "Comma-separated list of test packages to run")
	RunTestsCmd.Flags().StringArray("test-cmd", nil,
		"Optional custom test command(s) (e.g. 'go test -json ./... -v'), which must produce 'go test -json' output. "+
			"Avoid 'go test <file.go>' syntax as it prevents reliable log output and reruns. Use 'go test . -run TestName' instead. "+
			"Can be specified multiple times.",
	)
	RunTestsCmd.Flags().StringSlice("skip-tests", nil, "Comma-separated list of test names (regex supported by `go test -skip`) to skip")
	RunTestsCmd.Flags().StringSlice("select-tests", nil, "Comma-separated list of test names (regex supported by `go test -run`) to specifically run")
	RunTestsCmd.Flags().IntP("run-count", "c", 1, "Number of times to run the tests (for main run)")
	RunTestsCmd.Flags().Int("rerun-failed-count", 0, "Number of times to rerun tests that did not achieve 100% pass rate in the main run (0 disables reruns)")
	RunTestsCmd.Flags().StringArray("tags", nil, "Passed on to the 'go test' command as the -tags flag")
	RunTestsCmd.Flags().String("go-test-timeout", "", "Passed on to the 'go test' command as the -timeout flag (e.g., '30m')")
	RunTestsCmd.Flags().Int("go-test-count", -1, "Passes the '-count' flag directly to 'go test'. Default (-1) omits the flag.")
	RunTestsCmd.Flags().Bool("race", false, "Enable the race detector (-race flag for 'go test')")
	RunTestsCmd.Flags().Bool("shuffle", false, "Enable test shuffling ('go test -shuffle=on')")
	RunTestsCmd.Flags().String("shuffle-seed", "", "Set seed for test shuffling. Requires --shuffle. ('go test -shuffle=on -shuffle.seed=...')")
	RunTestsCmd.Flags().Bool("fail-fast", false, "Stop test execution on the first failure (-failfast flag for 'go test')")
	RunTestsCmd.Flags().String("main-results-path", "", "Path to save the main test results (JSON format)")
	RunTestsCmd.Flags().String("rerun-results-path", "", "Path to save the rerun test results (JSON format)")
	RunTestsCmd.Flags().Bool("omit-test-outputs-on-success", true, "Omit test outputs and package outputs for tests that pass all runs")
	RunTestsCmd.Flags().Bool("ignore-parent-failures-on-subtests", false, "Ignore failures in parent tests when only subtests fail (affects parsing)")
	RunTestsCmd.Flags().Float64("min-pass-ratio", 1.0, "The minimum pass ratio (0.0-1.0) required for a test in the main run to be considered stable.")
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
	cmd.Stderr = &out // Capture stderr as well

	if err := cmd.Run(); err != nil {
		// Don't block execution, just warn, as sometimes tidy fails for unrelated reasons
		log.Warn().Err(err).Str("output", out.String()).Msg("Dependency check ('go mod tidy') failed. Continuing execution, but dependencies might be inconsistent.")
		// return fmt.Errorf("dependency check ('go mod tidy') failed: %w
		return fmt.Errorf("dependency check ('go mod tidy') failed: %w\n%s", err, out.String())
	} else {
		log.Debug().Msg("'go mod tidy' completed successfully.")
	}
	return nil
}
