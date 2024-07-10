package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "e2e_tests_tool",
	Short: "A tool to manage E2E tests on Github CI",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(envresolveCmd)
	rootCmd.AddCommand(checkTestsCmd)
	rootCmd.AddCommand(filterCmd)
	rootCmd.AddCommand(testConfigCmd)
	testConfigCmd.AddCommand(maskTestConfigCmd)
	testConfigCmd.AddCommand(createTestConfigCmd)
	testConfigCmd.AddCommand(overrideTestConfigCmd)
}
