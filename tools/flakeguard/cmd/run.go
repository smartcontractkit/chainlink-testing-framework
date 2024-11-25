package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/reports"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner"
	"github.com/spf13/cobra"
)

var RunTestsCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests to check if they are flaky",
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, _ := cmd.Flags().GetString("project-path")
		testPackagesJson, _ := cmd.Flags().GetString("test-packages-json")
		testPackagesArg, _ := cmd.Flags().GetStringSlice("test-packages")
		runCount, _ := cmd.Flags().GetInt("run-count")
		useRace, _ := cmd.Flags().GetBool("race")
		outputPath, _ := cmd.Flags().GetString("output-json")
		maxPassRatio, _ := cmd.Flags().GetFloat64("max-pass-ratio")
		skipTests, _ := cmd.Flags().GetStringSlice("skip-tests")
		printFailedTests, _ := cmd.Flags().GetBool("print-failed-tests")
		useShuffle, _ := cmd.Flags().GetBool("shuffle")
		shuffleSeed, _ := cmd.Flags().GetString("shuffle-seed")

		// Check if project dependencies are correctly set up
		if err := checkDependencies(projectPath); err != nil {
			log.Fatalf("Error: %v", err)
		}

		var testPackages []string
		if testPackagesJson != "" {
			if err := json.Unmarshal([]byte(testPackagesJson), &testPackages); err != nil {
				log.Fatalf("Error decoding test packages JSON: %v", err)
			}
		} else if len(testPackagesArg) > 0 {
			testPackages = testPackagesArg
		} else {
			log.Fatalf("Error: must specify either --test-packages-json or --test-packages")
		}

		runner := runner.Runner{
			ProjectPath:          projectPath,
			Verbose:              true,
			RunCount:             runCount,
			UseRace:              useRace,
			SkipTests:            skipTests,
			SelectedTestPackages: testPackages,
			UseShuffle:           useShuffle,
			ShuffleSeed:          shuffleSeed,
		}

		testReport, err := runner.RunTests()
		if err != nil {
			fmt.Printf("Error running tests: %v\n", err)
			os.Exit(1)
		}

		// Print all failed tests including flaky tests
		if printFailedTests {
			fmt.Printf("PassRatio threshold for flaky tests: %.2f\n", maxPassRatio)
			reports.PrintTests(os.Stdout, testReport.Results, maxPassRatio)
		}

		// Save the test results in JSON format
		if outputPath != "" && len(testReport.Results) > 0 {
			jsonData, err := json.MarshalIndent(testReport, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling test results to JSON: %v", err)
			}
			if err := os.WriteFile(outputPath, jsonData, 0644); err != nil { //nolint:gosec
				log.Fatalf("Error writing test results to file: %v", err)
			}
			fmt.Printf("All test results saved to %s\n", outputPath)
		}

		flakyTests := reports.FilterFlakyTests(testReport.Results, maxPassRatio)
		if len(flakyTests) > 0 {
			// Exit with error code if there are flaky tests
			os.Exit(1)
		} else if len(testReport.Results) == 0 {
			fmt.Printf("No tests were run for the specified packages.\n")
		}
	},
}

func init() {
	RunTestsCmd.Flags().StringP("project-path", "r", ".", "The path to the Go project. Default is the current directory. Useful for subprojects")
	RunTestsCmd.Flags().String("test-packages-json", "", "JSON-encoded string of test packages")
	RunTestsCmd.Flags().StringSlice("test-packages", nil, "Comma-separated list of test packages to run")
	RunTestsCmd.Flags().Bool("run-all-packages", false, "Run all test packages in the project. This flag overrides --test-packages and --test-packages-json")
	RunTestsCmd.Flags().IntP("run-count", "c", 1, "Number of times to run the tests")
	RunTestsCmd.Flags().Bool("race", false, "Enable the race detector")
	RunTestsCmd.Flags().Bool("shuffle", false, "Enable test shuffling")
	RunTestsCmd.Flags().String("shuffle-seed", "", "Set seed for test shuffling. Must be used with --shuffle")
	RunTestsCmd.Flags().Bool("fail-fast", false, "Stop on the first test failure")
	RunTestsCmd.Flags().String("output-json", "", "Path to output the test results in JSON format")
	RunTestsCmd.Flags().StringSlice("skip-tests", nil, "Comma-separated list of test names to skip from running")
	RunTestsCmd.Flags().Bool("print-failed-tests", true, "Print failed test results to the console")
	RunTestsCmd.Flags().Float64("max-pass-ratio", 1.0, "The maximum (non-inclusive) pass ratio threshold for a test to be considered a failure. Any tests below this pass rate will be considered flaky.")
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
