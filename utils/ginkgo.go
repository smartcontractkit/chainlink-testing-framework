package utils

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
)

// GinkgoSuite provides the default setup for running a Ginkgo test suite
func GinkgoSuite() {
	RegisterFailHandler(Fail)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	conf, err := config.NewConfig(ProjectRoot)
	if err != nil {
		log.Panic().Msgf("Failed to load config with project root: %s", ProjectRoot)
		return
	}
	log.Logger = log.Logger.Level(zerolog.Level(conf.Logging.Level))
}
