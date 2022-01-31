package gauntlet_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	RunSpecs(t, "Gauntlet")
}
