package actions

import (
	"fmt"
	"math/big"
	"path/filepath"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	// DefaultArtifactsDir default artifacts dir
	DefaultArtifactsDir string = "logs"
)

// GinkgoSuite provides the default setup for running a Ginkgo test suite
func GinkgoSuite() {
	logging.Init()
	gomega.RegisterFailHandler(ginkgo.Fail)
}

// attempts to download the logs of all ephemeral test deployments onto the test runner, also writing a test report
// if one is provided
func WriteTeardownLogs(env *environment.Environment, optionalTestReporter testreporters.TestReporter) error {
	if ginkgo.CurrentSpecReport().Failed() || optionalTestReporter != nil {
		testFilename := strings.Split(ginkgo.CurrentSpecReport().FileName(), ".")[0]
		_, testName := filepath.Split(testFilename)
		logsPath := filepath.Join(DefaultArtifactsDir, fmt.Sprintf("%s-%d", testName, time.Now().Unix()))
		if err := env.Artifacts.DumpTestResult(logsPath, "chainlink"); err != nil {
			log.Warn().Err(err).Msg("Error trying to collect pod logs")
			if kubeerrors.IsForbidden(err) {
				log.Warn().Msg("Unable to gather logs from a remote_test_runner instance. Working on improving this.")
			} else {
				return err
			}
		}
		if err := SendReport(env, logsPath, optionalTestReporter); err != nil {
			log.Warn().Err(err).Msg("Error writing test report")
		}
	}
	return nil
}

// if provided, writes a test report and sends a Slack notification
func SendReport(env *environment.Environment, logsPath string, optionalTestReporter testreporters.TestReporter) error {
	if optionalTestReporter != nil {
		log.Info().Msg("Writing Test Report")
		optionalTestReporter.SetNamespace(env.Cfg.Namespace)
		err := optionalTestReporter.WriteReport(logsPath)
		if err != nil {
			return err
		}
		err = optionalTestReporter.SendSlackNotification(nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// FundAddresses will fund a list of addresses with an amount of native currency
func FundAddresses(blockchain blockchain.EVMClient, amount *big.Float, addresses ...string) error {
	for _, address := range addresses {
		if err := blockchain.Fund(address, amount); err != nil {
			return err
		}
	}
	return blockchain.WaitForEvents()
}

// EncodeOnChainExternalJobID encodes external job uuid to on-chain representation
func EncodeOnChainExternalJobID(jobID uuid.UUID) [32]byte {
	var ji [32]byte
	copy(ji[:], strings.Replace(jobID.String(), "-", "", 4))
	return ji
}
