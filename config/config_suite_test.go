package config_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	junitReporter := reporters.NewJUnitReporter("../logs/tests-config.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "NetworksConfig Suite", []Reporter{junitReporter})
}
