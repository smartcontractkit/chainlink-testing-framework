package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/utils"
)

// Exit codes for the application
const (
	// FlakyTestsExitCode indicates that Flakeguard ran correctly and was able to identify flaky tests
	FlakyTestsExitCode = 1
	// ErrorExitCode indicates that Flakeguard ran into an error and was not able to complete operation
	ErrorExitCode = 2
)

// runConfig holds the configuration parameters for the run command.
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

// summaryAndExit encapsulates the summary buffer and exit logic.
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
		exitHandler := &summaryAndExit{} // Initialize the exit handler

		cfg, err := parseAndValidateFlags(cmd)
		if err != nil {
			exitHandler.logErrorAndExit(err, "Failed to parse or validate flags")
		}

		goProject, err := utils.GetGoProjectName(cfg.ProjectPath)
		if err != nil {
			// Log warning but continue, as it's not critical for core functionality
			log.Warn().Err(err).Str("projectPath", cfg.ProjectPath).Msg("Failed to get pretty project path")
		}

		// Check if project dependencies are correctly set up
		if err := checkDependencies(cfg.ProjectPath); err != nil {
			exitHandler.logErrorAndExit(err, "Error checking project dependencies")
		}

		// Determine test packages
		testPackages, err := determineTestPackages(cfg)
		if err != nil {
			exitHandler.logErrorAndExit(err, "Failed to determine test packages")
		}

		// Initialize the runner
		testRunner := initializeRunner(cfg)

		// --- Main Test Execution ---
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Running tests..."
		s.Start()

		var mainResults []reports.TestResult
		var runErr error
		if len(cfg.TestCmds) > 0 {
			s.Suffix = " Running custom test command..."
			mainResults, runErr = testRunner.RunTestCmd(cfg.TestCmds)
		} else {
			s.Suffix = " Running test packages..."
			mainResults, runErr = testRunner.RunTestPackages(testPackages)
		}
		s.Stop()

		if runErr != nil {
			exitHandler.logErrorAndExit(runErr, "Error running tests")
		}
		if len(mainResults) == 0 {
			// Use logMsgAndExit for non-error exits
			exitHandler.logMsgAndExit(zerolog.ErrorLevel, "No tests were run.", ErrorExitCode)
		}

		// --- Report Generation and Saving ---
		mainReport, err := generateMainReport(mainResults, cfg, goProject)
		if err != nil {
			exitHandler.logErrorAndExit(err, "Error creating main test report")
		}
		if cfg.MainResultsPath != "" {
			if err := reports.SaveTestResultsToFile(mainResults, cfg.MainResultsPath); err != nil {
				// Log error but continue processing if possible
				log.Error().Err(err).Str("path", cfg.MainResultsPath).Msg("Error saving main test results to file")
			} else {
				log.Info().Str("path", cfg.MainResultsPath).Msg("Main test report saved")
			}
		}

		// --- Rerun Logic ---
		if cfg.RerunFailedCount > 0 {
			handleReruns(exitHandler, testRunner, mainReport, cfg, goProject)
		} else {
			handleNoReruns(exitHandler, mainReport, cfg)
		}

		// If we haven't exited yet, exit successfully
		exitHandler.code = 0
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

	// Resolve paths
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

	// Handle go-test-count flag
	if cmd.Flags().Changed("go-test-count") {
		v, err := cmd.Flags().GetInt("go-test-count")
		if err != nil {
			return nil, fmt.Errorf("error retrieving flag go-test-count: %w", err)
		}
		cfg.GoTestCount = &v
	}

	// Handle pass ratio compatibility
	minPassRatio, _ := cmd.Flags().GetFloat64("min-pass-ratio")
	maxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio") // Deprecated
	maxPassRatioSpecified := cmd.Flags().Changed("max-pass-ratio")

	cfg.MinPassRatio = minPassRatio
	if maxPassRatioSpecified && maxPassRatio != 1.0 {
		log.Warn().Msg("--max-pass-ratio is deprecated, please use --min-pass-ratio instead. Using max-pass-ratio value for now.")
		cfg.MinPassRatio = maxPassRatio // Use the deprecated value if specified
	}

	// Validate pass ratio
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
		// TODO: Consider adding support for --run-all-packages flag if needed
		return nil, fmt.Errorf("must specify either --test-packages-json, --test-packages, or --test-cmd")
	}
	return testPackages, nil
}

// initializeRunner creates and configures a new test runner.
func initializeRunner(cfg *runConfig) *runner.Runner {
	return runner.NewRunner(
		cfg.ProjectPath,
		true, // Assuming verbose is always true for now, consider adding a flag if needed
		// Runner specific config
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
		// Parser specific config
		cfg.IgnoreParentFailuresOnSubtests,
		cfg.OmitOutputsOnSuccess,
		// Dependencies (pass nil to get defaults)
		nil, // exec
		nil, // parser
	)
}

// generateMainReport creates the initial test report from the main run results.
func generateMainReport(results []reports.TestResult, cfg *runConfig, goProject string) (*reports.TestReport, error) {
	// Assuming NewTestReport returns (TestReport, error) based on the linter error
	reportVal, err := reports.NewTestReport(results,
		reports.WithGoProject(goProject),
		reports.WithCodeOwnersPath(cfg.CodeownersPath),
		reports.WithMaxPassRatio(cfg.MinPassRatio), // Use the validated MinPassRatio
		reports.WithGoRaceDetection(cfg.UseRace),
		reports.WithExcludedTests(cfg.SkipTests),
		reports.WithSelectedTests(cfg.SelectTests),
	)
	if err != nil {
		return nil, err // Propagate the error
	}
	return &reportVal, nil // Return the address of the struct value
}

// handleReruns manages the process of rerunning failed tests and reporting results.
func handleReruns(exitHandler *summaryAndExit, testRunner *runner.Runner, mainReport *reports.TestReport, cfg *runConfig, goProject string) {
	failedTests := reports.FilterTests(mainReport.Results, func(tr reports.TestResult) bool {
		return !tr.Skipped && tr.PassRatio < 1.0 // Rerun only tests that failed completely or partially in the main run
	})

	if len(failedTests) == 0 {
		log.Info().Msg("All tests passed the initial run. No tests to rerun.")
		// Use explicit newlines for clarity and to avoid linter issues
		fmt.Fprint(&exitHandler.buffer, "\nFlakeguard Summary\n")
		reports.RenderTestReport(&exitHandler.buffer, *mainReport, false, false)
		exitHandler.code = 0
		exitHandler.flush()
		return // Exit successfully
	}

	// Use explicit newlines
	fmt.Fprint(&exitHandler.buffer, "\nFailed Tests On The First Run:\n\n")
	reports.PrintTestResultsTable(&exitHandler.buffer, failedTests, false, false, true, false, false, false)
	fmt.Fprintln(&exitHandler.buffer)

	log.Info().Int("count", len(failedTests)).Int("rerun_count", cfg.RerunFailedCount).Msg("Rerunning failed tests...")

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Rerunning failed tests..."
	s.Start()
	rerunResults, rerunJsonOutputPaths, err := testRunner.RerunFailedTests(failedTests, cfg.RerunFailedCount)
	s.Stop()

	if err != nil {
		exitHandler.logErrorAndExit(err, "Error rerunning failed tests")
	}

	rerunReport, err := reports.NewTestReport(rerunResults,
		reports.WithGoProject(goProject),
		reports.WithCodeOwnersPath(cfg.CodeownersPath),
		reports.WithMaxPassRatio(1), // Rerun target is always 100% pass
		reports.WithExcludedTests(cfg.SkipTests),
		reports.WithSelectedTests(cfg.SelectTests),
		reports.WithJSONOutputPaths(rerunJsonOutputPaths),
	)
	if err != nil {
		exitHandler.logErrorAndExit(err, "Error creating rerun test report")
	}

	// Use explicit newlines
	fmt.Fprint(&exitHandler.buffer, "\nTests After Rerun:\n\n")
	reports.PrintTestResultsTable(&exitHandler.buffer, rerunResults, false, false, true, true, true, true)
	fmt.Fprintln(&exitHandler.buffer)

	// Save the rerun test report to file
	if cfg.RerunResultsPath != "" && len(rerunResults) > 0 {
		if err := reports.SaveTestResultsToFile(rerunResults, cfg.RerunResultsPath); err != nil {
			log.Error().Err(err).Str("path", cfg.RerunResultsPath).Msg("Error saving rerun test results to file")
			// Continue even if saving fails
		} else {
			log.Info().Str("path", cfg.RerunResultsPath).Msg("Rerun test report saved")
		}
	}

	// Filter tests that still failed after reruns (0 successes)
	failedAfterRerun := reports.FilterTests(rerunResults, func(tr reports.TestResult) bool {
		return !tr.Skipped && tr.Successes == 0
	})

	if len(failedAfterRerun) > 0 {
		// Use explicit newlines
		fmt.Fprint(&exitHandler.buffer, "\nPersistently Failing Test Logs:\n\n")
		// Attempt to print logs only for tests that never passed
		err := rerunReport.PrintGotestsumOutput(&exitHandler.buffer, "pkgname")
		if err != nil {
			log.Error().Err(err).Msg("Error printing gotestsum output for persistently failing tests")
			// Try printing all logs as a fallback - already done above
		}

		exitHandler.logMsgAndExit(zerolog.ErrorLevel, "Some tests are still failing after multiple reruns with no successful attempts.", ErrorExitCode, map[string]interface{}{
			"persistently_failing_count": len(failedAfterRerun),
			"rerun_attempts":             cfg.RerunFailedCount,
		})
	} else {
		// All originally failing tests passed at least once during reruns
		log.Info().Msg("All initially failing tests passed at least once after reruns.")
		exitHandler.code = 0 // Success
		// Summary already includes rerun info printed above
		exitHandler.flush()
	}
}

// handleNoReruns determines the outcome when reruns are disabled.
func handleNoReruns(exitHandler *summaryAndExit, mainReport *reports.TestReport, cfg *runConfig) {
	// Filter flaky tests based on the MinPassRatio from the initial run
	flakyTests := reports.FilterTests(mainReport.Results, func(tr reports.TestResult) bool {
		// A test is flaky if it wasn't skipped and its pass ratio is below the configured threshold
		return !tr.Skipped && tr.PassRatio < cfg.MinPassRatio
	})

	// Use explicit newlines
	fmt.Fprint(&exitHandler.buffer, "\nFlakeguard Summary\n")
	reports.RenderTestReport(&exitHandler.buffer, *mainReport, false, false)

	if len(flakyTests) > 0 {
		exitHandler.logMsgAndExit(zerolog.InfoLevel, "Found flaky tests.", FlakyTestsExitCode, map[string]interface{}{
			"flaky_count":         len(flakyTests),
			"stability_threshold": fmt.Sprintf("%.0f%%", cfg.MinPassRatio*100),
		})
	} else {
		log.Info().Msg("All tests passed stability requirements.")
		exitHandler.code = 0 // Success
		// Summary already printed above
		exitHandler.flush()
	}
}

// init sets up the cobra command flags.
func init() {
	// General execution flags
	RunTestsCmd.Flags().StringP("project-path", "r", ".", "The path to the Go project. Default is the current directory. Useful for subprojects")
	RunTestsCmd.Flags().StringP("codeowners-path", "", "", "Path to the CODEOWNERS file")

	// Test selection flags
	RunTestsCmd.Flags().String("test-packages-json", "", "JSON-encoded string of test packages")
	RunTestsCmd.Flags().StringSlice("test-packages", nil, "Comma-separated list of test packages to run")
	RunTestsCmd.Flags().StringArray("test-cmd", nil,
		"Optional custom test command(s) (e.g. 'go test -json ./... -v'), which must produce 'go test -json' output. Can be specified multiple times.",
	)
	RunTestsCmd.Flags().StringSlice("skip-tests", nil, "Comma-separated list of test names (regex supported by `go test -skip`) to skip")
	RunTestsCmd.Flags().StringSlice("select-tests", nil, "Comma-separated list of test names (regex supported by `go test -run`) to specifically run")

	// Test execution control flags
	RunTestsCmd.Flags().IntP("run-count", "c", 1, "Number of times to run the tests (for main run)")
	RunTestsCmd.Flags().Int("rerun-failed-count", 0, "Number of times to rerun tests that did not achieve 100% pass rate in the main run (0 disables reruns)")
	RunTestsCmd.Flags().StringArray("tags", nil, "Passed on to the 'go test' command as the -tags flag")
	RunTestsCmd.Flags().String("go-test-timeout", "", "Passed on to the 'go test' command as the -timeout flag (e.g., '30m')")
	RunTestsCmd.Flags().Int("go-test-count", -1, "Passes the '-count' flag directly to 'go test'. Default (-1) omits the flag.")
	RunTestsCmd.Flags().Bool("race", false, "Enable the race detector (-race flag for 'go test')")
	RunTestsCmd.Flags().Bool("shuffle", false, "Enable test shuffling ('go test -shuffle=on')")
	RunTestsCmd.Flags().String("shuffle-seed", "", "Set seed for test shuffling. Requires --shuffle. ('go test -shuffle=on -shuffle.seed=...')")
	RunTestsCmd.Flags().Bool("fail-fast", false, "Stop test execution on the first failure (-failfast flag for 'go test')")

	// Reporting and output flags
	RunTestsCmd.Flags().String("main-results-path", "", "Path to save the main test results (JSON format)")
	RunTestsCmd.Flags().String("rerun-results-path", "", "Path to save the rerun test results (JSON format)")
	RunTestsCmd.Flags().Bool("omit-test-outputs-on-success", true, "Omit test outputs and package outputs for tests that pass all runs")
	RunTestsCmd.Flags().Bool("ignore-parent-failures-on-subtests", false, "Ignore failures in parent tests when only subtests fail (affects parsing)")

	// Flake detection flags
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
		// Try to provide a more helpful error message including the output
		return fmt.Errorf("dependency check ('go mod tidy') failed: %w\n%s", err, out.String())
	}
	log.Debug().Msg("'go mod tidy' completed successfully.")
	return nil
}
