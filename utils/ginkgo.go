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
func GinkgoSuite() {
	RegisterFailHandler(Fail)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	fConf, err := config.LoadFrameworkConfig(filepath.Join(ProjectRoot, "framework.yaml"))
	if err != nil {
		log.Fatal().
			Str("Path", ProjectRoot).
			Msg("Failed to load config")
		return
	}
	log.Logger = log.Logger.Level(zerolog.Level(fConf.Logging.Level))
}
