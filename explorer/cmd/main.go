package main

import (
	"explorer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	e := explorer.NewExplorer()
	if err := e.Run(); err != nil {
		log.Fatal().Err(err)
	}
}
