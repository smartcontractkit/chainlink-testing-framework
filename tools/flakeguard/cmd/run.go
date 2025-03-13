package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner"
	"github.com/spf13/cobra"
)

const (
	// FlakyTestsExitCode indicates that Flakeguard ran correctly and was able to identify flaky tests
	FlakyTestsExitCode = 1
	// ErrorExitCode indicates that Flakeguard ran into an error and was not able to complete operation
	ErrorExitCode = 2
)

var RunTestsCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests to check if they are flaky",
	Run: func(cmd *cobra.Command, args []string) {
		// Create a buffer to accumulate all summary output.
		var summaryBuffer bytes.Buffer

		// Helper function to flush the summary buffer and exit.
		flushSummaryAndExit := func(code int) {
			fmt.Print(summaryBuffer.String())
			os.Exit(code)
		}

		// Retrieve flags
		projectPath, _ := cmd.Flags().GetString("project-path")
		testPackagesJson, _ := cmd.Flags().GetString("test-packages-json")
		testPackagesArg, _ := cmd.Flags().GetStringSlice("test-packages")
		testCmdStrings, _ := cmd.Flags().GetStringArray("test-cmd")
		runCount, _ := cmd.Flags().GetInt("run-count")
		rerunFailedCount, _ := cmd.Flags().GetInt("rerun-failed-count")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		tags, _ := cmd.Flags().GetStringArray("tags")
		useRace, _ := cmd.Flags().GetBool("race")
		mainReportPath, _ := cmd.Flags().GetString("main-report-path")
		rerunReportPath, _ := cmd.Flags().GetString("rerun-report-path")
		minPassRatio, _ := cmd.Flags().GetFloat64("min-pass-ratio")
		// For backward compatibility, check if max-pass-ratio was used
		maxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio")
		maxPassRatioSpecified := cmd.Flags().Changed("max-pass-ratio")
		skipTests, _ := cmd.Flags().GetStringSlice("skip-tests")
		selectTests, _ := cmd.Flags().GetStringSlice("select-tests")
		useShuffle, _ := cmd.Flags().GetBool("shuffle")
		shuffleSeed, _ := cmd.Flags().GetString("shuffle-seed")
		omitOutputsOnSuccess, _ := cmd.Flags().GetBool("omit-test-outputs-on-success")
		ignoreParentFailuresOnSubtests, _ := cmd.Flags().GetBool("ignore-parent-failures-on-subtests")
		failFast, _ := cmd.Flags().GetBool("fail-fast")

		// Retrieve go-test-count flag as a pointer if explicitly provided.
		var goTestCount *int
		if cmd.Flags().Changed("go-test-count") {
			v, err := cmd.Flags().GetInt("go-test-count")
			if err != nil {
				log.Error().Err(err).Msg("Error retrieving flag go-test-count")
				flushSummaryAndExit(ErrorExitCode)
			}
			goTestCount = &v
		}

		// Handle the compatibility between min/max pass ratio
		passRatioThreshold := minPassRatio
		if maxPassRatioSpecified && maxPassRatio != 1.0 {
			// If max-pass-ratio was explicitly set, use it (convert to min-pass-ratio)
			log.Warn().Msg("--max-pass-ratio is deprecated, please use --min-pass-ratio instead")
			passRatioThreshold = maxPassRatio
		}

		// Validate pass ratio
		if passRatioThreshold < 0 || passRatioThreshold > 1 {
			log.Error().Float64("pass ratio", passRatioThreshold).Msg("Error: pass ratio must be between 0 and 1")
			flushSummaryAndExit(ErrorExitCode)
		}

		// Check if project dependencies are correctly set up
		if err := checkDependencies(projectPath); err != nil {
			log.Error().Err(err).Msg("Error checking project dependencies")
			flushSummaryAndExit(ErrorExitCode)
		}

		// Determine test packages
		var testPackages []string
		if len(testCmdStrings) == 0 {
			if testPackagesJson != "" {
				if err := json.Unmarshal([]byte(testPackagesJson), &testPackages); err != nil {
					log.Error().Err(err).Msg("Error decoding test packages JSON")
					flushSummaryAndExit(ErrorExitCode)
				}
			} else if len(testPackagesArg) > 0 {
				testPackages = testPackagesArg
			} else {
				log.Error().Msg("Error: must specify either --test-packages-json or --test-packages")
				flushSummaryAndExit(ErrorExitCode)
			}
		}

		// Initialize the runner
		testRunner := runner.Runner{
			ProjectPath:                    projectPath,
			Verbose:                        true,
			RunCount:                       runCount,
			Timeout:                        timeout,
			Tags:                           tags,
			GoTestCountFlag:                goTestCount,
			GoTestRaceFlag:                 useRace,
			SkipTests:                      skipTests,
			SelectTests:                    selectTests,
			UseShuffle:                     useShuffle,
			ShuffleSeed:                    shuffleSeed,
			OmitOutputsOnSuccess:           omitOutputsOnSuccess,
			MaxPassRatio:                   passRatioThreshold, // Use the calculated threshold
			IgnoreParentFailuresOnSubtests: ignoreParentFailuresOnSubtests,
			RerunCount:                     rerunFailedCount,
			FailFast:                       failFast,
		}

		var (
			mainReport  *reports.TestReport // Main test report
			rerunReport *reports.TestReport // Test report after rerunning failed tests
		)

		// Run the tests
		var err error
		if len(testCmdStrings) > 0 {
			mainReport, err = testRunner.RunTestCmd(testCmdStrings)
			if err != nil {
				log.Fatal().Err(err).Msg("Error running custom test command")
				flushSummaryAndExit(ErrorExitCode)
			}
		} else {
			mainReport, err = testRunner.RunTestPackages(testPackages)
			if err != nil {
				log.Fatal().Err(err).Msg("Error running test packages")
				flushSummaryAndExit(ErrorExitCode)
			}
		}

		// Save the main test report to file
		if mainReportPath != "" && len(mainReport.Results) > 0 {
			if err := mainReport.SaveToFile(mainReportPath); err != nil {
				log.Error().Err(err).Msg("Error saving test results to file")
				flushSummaryAndExit(ErrorExitCode)
			}
			log.Info().Str("path", mainReportPath).Msg("Main test report saved")
		}

		if len(mainReport.Results) == 0 {
			log.Warn().Msg("No tests were run for the specified packages")
			flushSummaryAndExit(0)
		}

		// Accumulate main summary into summaryBuffer
		fmt.Fprint(&summaryBuffer, "\nFlakeguard Main Summary:\n")
		fmt.Fprintf(&summaryBuffer, "-------------------------\n")
		reports.RenderResults(&summaryBuffer, *mainReport, false, false)
		fmt.Fprintln(&summaryBuffer)

		// Rerun failed tests
		if rerunFailedCount > 0 {
			failedTests := reports.FilterTests(mainReport.Results, func(tr reports.TestResult) bool {
				return !tr.Skipped && tr.PassRatio < 1.0
			})

			if len(failedTests) == 0 {
				log.Info().Msg("No tests to rerun. All tests passed")
				flushSummaryAndExit(0)
			}

			rerunReport, err = testRunner.RerunFailedTests(failedTests)
			if err != nil {
				log.Fatal().Err(err).Msg("Error rerunning failed tests")
				flushSummaryAndExit(ErrorExitCode)
			}

			fmt.Fprint(&summaryBuffer, "\nAll Tests That Were Rerun:\n")
			fmt.Fprintf(&summaryBuffer, "---------------------------\n")
			reports.PrintTestResultsTable(&summaryBuffer, rerunReport.Results, false, false)
			fmt.Fprintln(&summaryBuffer)

			// Save the rerun test report to file
			if rerunReportPath != "" && len(rerunReport.Results) > 0 {
				if err := rerunReport.SaveToFile(rerunReportPath); err != nil {
					log.Error().Err(err).Msg("Error saving test results to file")
					flushSummaryAndExit(ErrorExitCode)
				}
				log.Info().Str("path", rerunReportPath).Msg("Rerun test report saved")
			}

			// Filter tests that failed after reruns
			failedAfterRerun := reports.FilterTests(rerunReport.Results, func(tr reports.TestResult) bool {
				return !tr.Skipped && tr.Successes == 0
			})

			if len(failedAfterRerun) > 0 {
				fmt.Fprint(&summaryBuffer, "\nTests That Have 0 Success Runs:\n")
				fmt.Fprintf(&summaryBuffer, "--------------------------------\n")
				reports.PrintTestResultsTable(&summaryBuffer, failedAfterRerun, false, false)
				fmt.Fprintln(&summaryBuffer)

				fmt.Fprint(&summaryBuffer, "\nLogs From All Reruns:\n")
				fmt.Fprintf(&summaryBuffer, "----------------------\n")
				err := rerunReport.PrintGotestsumOutput(&summaryBuffer, "pkgname")
				if err != nil {
					log.Error().Err(err).Msg("Error printing gotestsum output")
				}

				log.Error().
					Int("noSuccessTests", len(failedAfterRerun)).
					Int("reruns", rerunFailedCount).
					Msg("Some tests are still failing after multiple reruns with no successful attempts.")
				flushSummaryAndExit(ErrorExitCode)
			} else {
				log.Info().Msg("All tests passed at least once after reruns")
				flushSummaryAndExit(0)
			}
		} else {
			// Filter flaky tests using FilterTests
			flakyTests := reports.FilterTests(mainReport.Results, func(tr reports.TestResult) bool {
				return !tr.Skipped && tr.PassRatio < passRatioThreshold
			})

			if len(flakyTests) > 0 {
				log.Info().
					Int("count", len(flakyTests)).
					Str("stability threshold", fmt.Sprintf("%.0f%%", passRatioThreshold*100)).
					Msg("Found flaky tests")
				flushSummaryAndExit(FlakyTestsExitCode)
			} else {
				log.Info().Msg("All tests passed stability requirements")
			}
		}

		flushSummaryAndExit(0)
	},
}

func init() {
	RunTestsCmd.Flags().StringP("project-path", "r", ".", "The path to the Go project. Default is the current directory. Useful for subprojects")
	RunTestsCmd.Flags().String("test-packages-json", "", "JSON-encoded string of test packages")
	RunTestsCmd.Flags().StringSlice("test-packages", nil, "Comma-separated list of test packages to run")
	RunTestsCmd.Flags().StringArray("test-cmd", nil,
		"Optional custom test command (e.g. 'go test -json github.com/smartcontractkit/chainlink/integration-tests/smoke -v -run TestForwarderOCR2Basic'), which must produce go test -json output.",
	)
	RunTestsCmd.Flags().Bool("run-all-packages", false, "Run all test packages in the project. This flag overrides --test-packages and --test-packages-json")
	RunTestsCmd.Flags().IntP("run-count", "c", 1, "Number of times to run the tests")
	RunTestsCmd.Flags().Duration("timeout", 0, "Passed on to the 'go test' command as the -timeout flag")
	RunTestsCmd.Flags().StringArray("tags", nil, "Passed on to the 'go test' command as the -tags flag")
	RunTestsCmd.Flags().Int("go-test-count", -1, "go test -count flag value. By default -count flag is not passed to go test")
	RunTestsCmd.Flags().Bool("race", false, "Enable the race detector")
	RunTestsCmd.Flags().Bool("shuffle", false, "Enable test shuffling")
	RunTestsCmd.Flags().String("shuffle-seed", "", "Set seed for test shuffling. Must be used with --shuffle")
	RunTestsCmd.Flags().Bool("fail-fast", false, "Stop on the first test failure")
	RunTestsCmd.Flags().String("main-report-path", "", "Path to the main test report in JSON format")
	RunTestsCmd.Flags().String("rerun-report-path", "", "Path to the rerun test report in JSON format")
	RunTestsCmd.Flags().StringSlice("skip-tests", nil, "Comma-separated list of test names to skip from running")
	RunTestsCmd.Flags().StringSlice("select-tests", nil, "Comma-separated list of test names to specifically run")

	// Add the min-pass-ratio flag (new recommended approach)
	RunTestsCmd.Flags().Float64("min-pass-ratio", 1.0, "The minimum pass ratio required for a test to be considered stable (0.0-1.0)")

	// Keep max-pass-ratio for backward compatibility but mark as deprecated
	RunTestsCmd.Flags().Float64("max-pass-ratio", 1.0, "DEPRECATED: Use min-pass-ratio instead")
	RunTestsCmd.Flags().MarkDeprecated("max-pass-ratio", "use min-pass-ratio instead")

	RunTestsCmd.Flags().Bool("omit-test-outputs-on-success", true, "Omit test outputs and package outputs for tests that pass")
	RunTestsCmd.Flags().Bool("ignore-parent-failures-on-subtests", false, "Ignore failures in parent tests when only subtests fail")

	// Add rerun failed tests flag
	RunTestsCmd.Flags().Int("rerun-failed-count", 0, "Number of times to rerun tests that did not get 100 percent pass ratio (0 disables reruns)")
}

func checkDependencies(projectPath string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dependency check failed: %v\n%s\nPlease run 'go mod tidy' to fix missing or unused dependencies", err, out.String())
	}

	return nil
}
