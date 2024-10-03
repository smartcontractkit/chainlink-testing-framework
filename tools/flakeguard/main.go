package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/cmd"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flakeguard",
	Short: "A tool to find flaky tests",
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

	rootCmd.AddCommand(cmd.FindtestsCmd)
}

func main() {
	rootCmd.Execute()
}
