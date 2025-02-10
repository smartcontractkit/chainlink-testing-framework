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

var RunTestsCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests to check if they are flaky",
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve flags
		projectPath, _ := cmd.Flags().GetString("project-path")
		testPackagesJson, _ := cmd.Flags().GetString("test-packages-json")
		testPackagesArg, _ := cmd.Flags().GetStringSlice("test-packages")
		testCmdStrings, _ := cmd.Flags().GetStringArray("test-cmd")
		runCount, _ := cmd.Flags().GetInt("run-count")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		tags, _ := cmd.Flags().GetStringArray("tags")
		useRace, _ := cmd.Flags().GetBool("race")
		outputPath, _ := cmd.Flags().GetString("output-json")
		maxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio")
		skipTests, _ := cmd.Flags().GetStringSlice("skip-tests")
		selectTests, _ := cmd.Flags().GetStringSlice("select-tests")
		useShuffle, _ := cmd.Flags().GetBool("shuffle")
		shuffleSeed, _ := cmd.Flags().GetString("shuffle-seed")
		omitOutputsOnSuccess, _ := cmd.Flags().GetBool("omit-test-outputs-on-success")

		// Check if project dependencies are correctly set up
		if err := checkDependencies(projectPath); err != nil {
			log.Fatal().Err(err).Msg("Error checking project dependencies")
		}

		// Determine test packages
		var testPackages []string
		if len(testCmdStrings) == 0 {
			// No custom command -> parse packages
			if testPackagesJson != "" {
				if err := json.Unmarshal([]byte(testPackagesJson), &testPackages); err != nil {
					log.Fatal().Err(err).Msg("Error decoding test packages JSON")
				}
			} else if len(testPackagesArg) > 0 {
				testPackages = testPackagesArg
			} else {
				log.Fatal().Msg("Error: must specify either --test-packages-json or --test-packages (or use --test-cmd).")
			}
		}

		// Initialize the runner
		testRunner := runner.Runner{
			ProjectPath:          projectPath,
			Verbose:              true,
			RunCount:             runCount,
			Timeout:              timeout,
			Tags:                 tags,
			UseRace:              useRace,
			SkipTests:            skipTests,
			SelectTests:          selectTests,
			SelectedTestPackages: testPackages,
			UseShuffle:           useShuffle,
			ShuffleSeed:          shuffleSeed,
			OmitOutputsOnSuccess: omitOutputsOnSuccess,
		}

		// Run the tests
		var (
			testReport *reports.TestReport
			err        error
		)

		if len(testCmdStrings) > 0 {
			testReport, err = testRunner.RunTestsByCmd(testCmdStrings)
			if err != nil {
				log.Fatal().Err(err).Msg("Error running custom test command")
			}
		} else {
			// Otherwise, use the normal go test approach
			testReport, err = testRunner.RunTests()
			if err != nil {
				log.Fatal().Err(err).Msg("Error running tests")
			}
		}

		// Save the test results in JSON format
		if outputPath != "" && len(testReport.Results) > 0 {
			jsonData, err := json.MarshalIndent(testReport, "", "  ")
			if err != nil {
				log.Fatal().Err(err).Msg("Error marshaling test results to JSON")
			}
			if err := os.WriteFile(outputPath, jsonData, 0600); err != nil {
				log.Fatal().Err(err).Msg("Error writing test results to file")
			}
			log.Info().Str("path", outputPath).Msg("Test results saved")
		}

		if len(testReport.Results) == 0 {
			log.Warn().Msg("No tests were run for the specified packages")
			return
		}

		// Filter flaky tests using FilterTests
		flakyTests := reports.FilterTests(testReport.Results, func(tr reports.TestResult) bool {
			return !tr.Skipped && tr.PassRatio < maxPassRatio
		})

		if len(flakyTests) > 0 {
			log.Info().Int("count", len(flakyTests)).Str("pass ratio threshold", fmt.Sprintf("%.2f%%", maxPassRatio*100)).Msg("Found flaky tests")
			fmt.Printf("\nFlakeguard Summary\n")
			reports.RenderResults(os.Stdout, flakyTests, maxPassRatio, false)
			// Exit with error code if there are flaky tests
			os.Exit(1)
		}
	},
}

func init() {
	RunTestsCmd.Flags().StringP("project-path", "r", ".", "The path to the Go project. Default is the current directory. Useful for subprojects")
	RunTestsCmd.Flags().String("test-packages-json", "", "JSON-encoded string of test packages")
	RunTestsCmd.Flags().StringSlice("test-packages", nil, "Comma-separated list of test packages to run")
	RunTestsCmd.Flags().StringArray("test-cmd", nil,
		"Optional custom test command (e.g. 'go run e2e_test.go -someflag'), which must produce go test -json output.",
	)
	RunTestsCmd.Flags().Bool("run-all-packages", false, "Run all test packages in the project. This flag overrides --test-packages and --test-packages-json")
	RunTestsCmd.Flags().IntP("run-count", "c", 1, "Number of times to run the tests")
	RunTestsCmd.Flags().Duration("timeout", 0, "Passed on to the 'go test' command as the -timeout flag")
	RunTestsCmd.Flags().StringArray("tags", nil, "Passed on to the 'go test' command as the -tags flag")
	RunTestsCmd.Flags().Bool("race", false, "Enable the race detector")
	RunTestsCmd.Flags().Bool("shuffle", false, "Enable test shuffling")
	RunTestsCmd.Flags().String("shuffle-seed", "", "Set seed for test shuffling. Must be used with --shuffle")
	RunTestsCmd.Flags().Bool("fail-fast", false, "Stop on the first test failure")
	RunTestsCmd.Flags().String("output-json", "", "Path to output the test results in JSON format")
	RunTestsCmd.Flags().StringSlice("skip-tests", nil, "Comma-separated list of test names to skip from running")
	RunTestsCmd.Flags().StringSlice("select-tests", nil, "Comma-separated list of test names to specifically run")
	RunTestsCmd.Flags().Float64("max-pass-ratio", 1.0, "The maximum pass ratio threshold for a test to be considered flaky. Any tests below this pass rate will be considered flaky.")
	RunTestsCmd.Flags().Bool("omit-test-outputs-on-success", true, "Omit test outputs and package outputs for tests that pass")
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
