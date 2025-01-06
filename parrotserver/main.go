package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var config *Config

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05.00", // hh:mm:ss.ss format
	})
	config = readConfig()
	level := zerolog.DebugLevel
	if config.LogLevel != "" {
		l, err := zerolog.ParseLevel(config.LogLevel)
		if err != nil {
			fmt.Printf("Invalid log level '%s'\n", config.LogLevel)
			os.Exit(1)
		}
		level = l
	}
	log.Logger = log.Logger.Level(level).With().Timestamp().Logger()
}

func start() int {
	defer func() {
		if err := save(); err != nil {
			log.Error().Err(err).Msg("Failed to save configuration")
		}
	}()

	http.HandleFunc("/register", registerRouteHandler)
	http.HandleFunc("/", dynamicHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Info().Int("port", 8080).Interface("config", config).Msg("Parrot server started")
	if err := server.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("Server stopped")
		return 1
	}
	return 0
}

func main() {
	os.Exit(start())
}
