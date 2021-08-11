package performance_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestContracts(t *testing.T) {
	RegisterFailHandler(Fail)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Performance suite", []Reporter{junitReporter})
}
