package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner"
	"github.com/spf13/cobra"
)

var RunTestsCmd = &cobra.Command{
	Use:   "run-tests",
	Short: "Run tests to find flaky ones",
	Run: func(cmd *cobra.Command, args []string) {
		repoPath, _ := cmd.Flags().GetString("repo")
		testPathsJson, _ := cmd.Flags().GetString("test-paths-json")

		var testPaths []string
		if err := json.Unmarshal([]byte(testPathsJson), &testPaths); err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding test paths JSON: %s\n", err)
			os.Exit(1)
		}

		runner := runner.Runner{
			Verbose: true,
			Dir:     repoPath,
		}

		if err := runner.RunTests(testPaths); err != nil {
			fmt.Fprintf(os.Stderr, "Error running tests: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RunTestsCmd.Flags().StringP("repo", "r", ".", "Path to the Git repository")
	RunTestsCmd.Flags().String("test-paths-json", "", "JSON-encoded string of test paths")
}
