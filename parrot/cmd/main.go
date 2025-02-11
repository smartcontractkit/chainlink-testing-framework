package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/parrot"
	"github.com/spf13/cobra"
)

func main() {
	var (
		port      int
		debug     bool
		trace     bool
		silent    bool
		json      bool
		recorders []string
	)

	rootCmd := &cobra.Command{
		Use:   "parrot",
		Short: "A server that can register and parrot back dynamic requests",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Check environment variables if flags are not set
			if !cmd.Flags().Changed("port") {
				if envPort, err := strconv.Atoi(os.Getenv("PARROT_PORT")); err == nil {
					port = envPort
				}
			}
			if !cmd.Flags().Changed("debug") {
				debug = os.Getenv("PARROT_DEBUG") == "true"
			}
			if !cmd.Flags().Changed("trace") {
				trace = os.Getenv("PARROT_TRACE") == "true"
			}
			if !cmd.Flags().Changed("silent") {
				silent = os.Getenv("PARROT_SILENT") == "true"
			}
			if !cmd.Flags().Changed("json") {
				json = os.Getenv("PARROT_JSON") == "true"
			}
			if !cmd.Flags().Changed("recorders") {
				if envRecorders := os.Getenv("PARROT_RECORDERS"); envRecorders != "" {
					recorders = strings.Split(envRecorders, ",")
				}
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			options := []parrot.ServerOption{parrot.WithPort(port)}
			logLevel := zerolog.InfoLevel
			if debug {
				logLevel = zerolog.DebugLevel
			}
			if trace {
				logLevel = zerolog.TraceLevel
			}
			if silent {
				logLevel = zerolog.Disabled
			}
			options = append(options, parrot.WithLogLevel(logLevel))
			if json {
				options = append(options, parrot.WithJSONLogs())
			}
			options = append(options, parrot.WithRecorders(recorders...))

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			p, err := parrot.NewServer(options...)
			if err != nil {
				return err
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c
			err = p.Shutdown(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Error putting parrot to sleep")
			}
			return nil
		},
	}

	rootCmd.Flags().IntVarP(&port, "port", "p", 0, "Port to run the parrot on (env: PARROT_PORT)")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug output (env: PARROT_DEBUG)")
	rootCmd.Flags().BoolVarP(&trace, "trace", "t", false, "Enable trace and debug output (env: PARROT_TRACE)")
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "Disable all output (env: PARROT_SILENT)")
	rootCmd.Flags().BoolVarP(&json, "json", "j", false, "Output logs in JSON format (env: PARROT_JSON)")
	rootCmd.Flags().StringSliceVarP(&recorders, "recorders", "r", nil, "Existing recorders to use (env: PARROT_RECORDERS)")

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("error executing command")
		os.Exit(1)
	}
}
