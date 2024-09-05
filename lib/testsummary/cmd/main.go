package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/testsummary/cmd/internal"
)

var rootCmd = &cobra.Command{
	Use:   "test-summary",
	Short: "Tests summary printer",
}

func init() {
	rootCmd.AddCommand(internal.PrintKeyCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
