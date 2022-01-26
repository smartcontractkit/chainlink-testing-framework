package utils

//revive:disable:dot-imports
import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
)

// GinkgoSuite provides the default setup for running a Ginkgo test suite
func GinkgoSuite(frameworkConfigFileLocation string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	RegisterFailHandler(Fail)
	absoluteConfigFileLocation, err := filepath.Abs(frameworkConfigFileLocation)
	if err != nil {
		log.Fatal().
			Str("Path", frameworkConfigFileLocation).
			Msg("Unable to resolve path to an absolute path")
		return
	}

	fConf, err := config.LoadFrameworkConfig(filepath.Join(absoluteConfigFileLocation, "framework.yaml"))
	if err != nil {
		log.Fatal().
			Str("Path", absoluteConfigFileLocation).
			Msg("Failed to load config")
		return
	}
	log.Logger = log.Logger.Level(zerolog.Level(fConf.Logging.Level))

	_, err = config.LoadNetworksConfig(filepath.Join(absoluteConfigFileLocation, "networks.yaml"))
	if err != nil {
		log.Fatal().
			Str("Path", absoluteConfigFileLocation).
			Msg("Failed to load config")
		return
	}
}
