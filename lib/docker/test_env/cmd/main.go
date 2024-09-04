package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/testcontainers/testcontainers-go"
	"golang.org/x/net/context"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env/cmd/internal"
)

var rootCmd = &cobra.Command{
	Use:   "test-envs",
	Short: "CTF Test Environments Tool",
}

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	internal.StartTestEnvCmd.SetContext(ctx)

	rootCmd.AddCommand(internal.StartTestEnvCmd)

	// Set default log level for non-testcontainer code
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Discard testcontainers logs
	testcontainers.Logger = log.New(io.Discard, "", log.LstdFlags)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
