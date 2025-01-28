package main

import (
	"context"
	"os"
	"os/signal"
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

			p, err := parrot.Wake(options...)
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
			p.WaitShutdown()
			return nil
		},
	}

	rootCmd.Flags().IntVarP(&port, "port", "p", 0, "Port to run the parrot on")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")
	rootCmd.Flags().BoolVarP(&trace, "trace", "t", false, "Enable trace and debug output")
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "Disable all output")
	rootCmd.Flags().BoolVarP(&json, "json", "j", false, "Output logs in JSON format")
	rootCmd.Flags().StringSliceVarP(&recorders, "recorders", "r", nil, "Existing recorders to use")

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("error executing command")
		os.Exit(1)
	}
}
