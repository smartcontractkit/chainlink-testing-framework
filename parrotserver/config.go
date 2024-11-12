package main

import (
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
	SaveFile string `json:"save_file"`
	LogLevel string `json:"log_level"`
}

func readConfig() *Config {
	saveFile := os.Getenv("SAVE_FILE")
	if saveFile == "" {
		log.Warn().Msg("SAVE_FILE is not set. Using default file 'save.json'")
		saveFile = "save.json"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		log.Warn().Msg("LOG_LEVEL is not set. Using default level 'debug'")
		logLevel = "debug"
	}
	log.Info().Str("SaveFile", saveFile).Str("LogLevel", logLevel).Msg("Loaded config")
	return &Config{
		SaveFile: saveFile,
		LogLevel: logLevel,
	}
}
