package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/explorer"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	e := explorer.NewExplorer()
	if err := e.Run(); err != nil {
		log.Fatal().Err(err)
	}
}
