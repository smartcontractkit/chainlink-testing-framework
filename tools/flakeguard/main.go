package main

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
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
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        io.Discard,
		TimeFormat: "15:04:05.00", // hh:mm:ss.ss format
	})
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(cmd.FindTestsCmd)
	rootCmd.AddCommand(cmd.RunTestsCmd)
	rootCmd.AddCommand(cmd.CheckTestOwnersCmd)
	rootCmd.AddCommand(cmd.GenerateTestReportCmd)
	rootCmd.AddCommand(cmd.GenerateReportCmd)
	rootCmd.AddCommand(cmd.GetGHArtifactLinkCmd)
	rootCmd.AddCommand(cmd.SendToSplunkCmd)
	rootCmd.AddCommand(cmd.ReviewTicketsCmd)
	rootCmd.AddCommand(cmd.CreateTicketsCmd)
	rootCmd.AddCommand(cmd.SyncJiraCmd)
}

func main() {
	rootCmd.Execute()
}
