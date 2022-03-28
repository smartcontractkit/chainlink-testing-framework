// Package actions enables common chainlink interactions
package actions

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/helmenv/environment"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/testreporters"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/client"
)

const (
	// DefaultArtifactsDir default artifacts dir
	DefaultArtifactsDir string = "logs"
)

// GinkgoSuite provides the default setup for running a Ginkgo test suite
func GinkgoSuite(frameworkConfigFileLocation string) {
	LoadConfigs(frameworkConfigFileLocation)
	gomega.RegisterFailHandler(ginkgo.Fail)
}

func LoadConfigs(frameworkConfigFileLocation string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	absoluteConfigFileLocation, err := filepath.Abs(frameworkConfigFileLocation)
	if err != nil {
		log.Fatal().
			Str("Path", frameworkConfigFileLocation).
			Msg("Unable to resolve path to an absolute path")
		return
	}

	frameworkConfig := filepath.Join(absoluteConfigFileLocation, "framework.yaml")
	if os.Getenv("FRAMEWORK_CONFIG_FILE") != "" {
		frameworkConfig = os.Getenv("FRAMEWORK_CONFIG_FILE")
	}
	fConf, err := config.LoadFrameworkConfig(frameworkConfig)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("Path", absoluteConfigFileLocation).
			Msg("Failed to load config")
		return
	}
	log.Logger = log.Logger.Level(zerolog.Level(fConf.Logging.Level))

	networksConfig := filepath.Join(absoluteConfigFileLocation, "networks.yaml")
	if os.Getenv("NETWORKS_CONFIG_FILE") != "" {
		networksConfig = os.Getenv("NETWORKS_CONFIG_FILE")
	}
	_, err = config.LoadNetworksConfig(networksConfig)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("Path", absoluteConfigFileLocation).
			Msg("Failed to load config")
		return
	}
}

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodes(
	nodes []client.Chainlink,
	blockchain client.BlockchainClient,
	amount *big.Float,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		err = blockchain.Fund(toAddress, amount)
		if err != nil {
			return err
		}
	}
	return blockchain.WaitForEvents()
}

// FundAddresses will fund a list of addresses with an amount of native currency
func FundAddresses(blockchain client.BlockchainClient, amount *big.Float, addresses ...string) error {
	for _, address := range addresses {
		if err := blockchain.Fund(address, amount); err != nil {
			return err
		}
	}
	return blockchain.WaitForEvents()
}

// ChainlinkNodeAddresses will return all the on-chain wallet addresses for a set of Chainlink nodes
func ChainlinkNodeAddresses(nodes []client.Chainlink) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, node := range nodes {
		primaryAddress, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(primaryAddress))
	}
	return addresses, nil
}

// SetChainlinkAPIPageSize specifies the page size from the Chainlink API, useful for high volume testing
func SetChainlinkAPIPageSize(nodes []client.Chainlink, pageSize int) {
	for _, n := range nodes {
		n.SetPageSize(pageSize)
	}
}

// EncodeOnChainExternalJobID encodes external job uuid to on-chain representation
func EncodeOnChainExternalJobID(jobID uuid.UUID) [32]byte {
	var ji [32]byte
	copy(ji[:], strings.Replace(jobID.String(), "-", "", 4))
	return ji
}

// ExtractRequestIDFromJobRun extracts RequestID from job runs response
func ExtractRequestIDFromJobRun(jobDecodeData client.RunsResponseData) ([]byte, error) {
	var taskRun client.TaskRun
	for _, tr := range jobDecodeData.Attributes.TaskRuns {
		if tr.Type == "ethabidecodelog" {
			taskRun = tr
		}
	}
	var decodeLogTaskRun *client.DecodeLogTaskRun
	if err := json.Unmarshal([]byte(taskRun.Output), &decodeLogTaskRun); err != nil {
		return nil, err
	}
	rqInts := decodeLogTaskRun.RequestID
	return rqInts, nil
}

// EncodeOnChainVRFProvingKey encodes uncompressed public VRF key to on-chain representation
func EncodeOnChainVRFProvingKey(vrfKey client.VRFKey) ([2]*big.Int, error) {
	uncompressed := vrfKey.Data.Attributes.Uncompressed
	provingKey := [2]*big.Int{}
	var set1 bool
	var set2 bool
	// strip 0x to convert to int
	provingKey[0], set1 = new(big.Int).SetString(uncompressed[2:66], 16)
	if !set1 {
		return [2]*big.Int{}, fmt.Errorf("can not convert VRF key to *big.Int")
	}
	provingKey[1], set2 = new(big.Int).SetString(uncompressed[66:], 16)
	if !set2 {
		return [2]*big.Int{}, fmt.Errorf("can not convert VRF key to *big.Int")
	}
	return provingKey, nil
}

// GetMockserverInitializerDataForOTPE creates mocked weiwatchers data needed for otpe
func GetMockserverInitializerDataForOTPE(
	OCRInstances []contracts.OffchainAggregator,
	chainlinkNodes []client.Chainlink,
) (interface{}, error) {
	var contractsInfo []client.ContractInfoJSON

	for index, OCRInstance := range OCRInstances {
		contractInfo := client.ContractInfoJSON{
			ContractVersion: 4,
			Path:            fmt.Sprintf("contract_%d", index),
			Status:          "live",
			ContractAddress: OCRInstance.Address(),
		}

		contractsInfo = append(contractsInfo, contractInfo)
	}

	contractsInitializer := client.HttpInitializer{
		Request:  client.HttpRequest{Path: "/contracts.json"},
		Response: client.HttpResponse{Body: contractsInfo},
	}

	var nodesInfo []client.NodeInfoJSON

	for _, chainlink := range chainlinkNodes {
		ocrKeys, err := chainlink.ReadOCRKeys()
		if err != nil {
			return nil, err
		}
		nodeInfo := client.NodeInfoJSON{
			NodeAddress: []string{ocrKeys.Data[0].Attributes.OnChainSigningAddress},
			ID:          ocrKeys.Data[0].ID,
		}
		nodesInfo = append(nodesInfo, nodeInfo)
	}

	nodesInitializer := client.HttpInitializer{
		Request:  client.HttpRequest{Path: "/nodes.json"},
		Response: client.HttpResponse{Body: nodesInfo},
	}
	initializers := []client.HttpInitializer{contractsInitializer, nodesInitializer}
	return initializers, nil
}

// TeardownSuite tears down networks/clients and environment and creates a logs folder for failed tests in the
// specified path. Can also accept a testsetup (if one was used) to log results
func TeardownSuite(
	env *environment.Environment,
	nets *client.Networks,
	logsFolderPath string,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
) error {
	if err := writeTeardownLogs(env, optionalTestReporter); err != nil {
		return errors.Wrap(err, "Error dumping environment logs, leaving environment running for manual retrieval")
	}
	if nets != nil {
		if err := nets.Teardown(); err != nil {
			return err
		}
	}
	switch strings.ToUpper(config.ProjectFrameworkSettings.KeepEnvironments) {
	case "ALWAYS":
		env.Persistent = true
	case "ONFAIL":
		if ginkgo.CurrentSpecReport().Failed() {
			env.Persistent = true
		}
	case "NEVER":
		env.Persistent = false
	default:
		log.Warn().Str("Invalid Keep Value", config.ProjectFrameworkSettings.KeepEnvironments).
			Msg("Invalid 'keep_environments' value, see the 'framework.yaml' file")
	}
	if !env.Config.Persistent {
		if err := env.Teardown(); err != nil {
			return err
		}
	}
	return nil
}

// TeardownRemoteSuite is used when running a test within a remote-test-runner, like for long-running performance and
// soak tests
func TeardownRemoteSuite(env *environment.Environment, optionalTestReporter testreporters.TestReporter) error {
	return writeTeardownLogs(env, optionalTestReporter)
}

func writeTeardownLogs(env *environment.Environment, optionalTestReporter testreporters.TestReporter) error {
	if ginkgo.CurrentSpecReport().Failed() || optionalTestReporter != nil {
		testFilename := strings.Split(ginkgo.CurrentSpecReport().FileName(), ".")[0]
		_, testName := filepath.Split(testFilename)
		logsPath := filepath.Join(config.ProjectConfigDirectory, DefaultArtifactsDir, fmt.Sprintf("%s-%d", testName, time.Now().Unix()))
		if err := env.Artifacts.DumpTestResult(logsPath, "chainlink"); err != nil {
			log.Warn().Err(err).Msg("Error trying to collect pod logs")
			if kubeerrors.IsForbidden(err) {
				log.Warn().Msg("Unable to gather logs from a remote_test_runner instance. Working on improving this.")
			} else {
				return err
			}
		}
		if optionalTestReporter != nil {
			log.Info().Msg("Writing Test Report")
			optionalTestReporter.SetNamespace(env.Namespace)
			err := optionalTestReporter.WriteReport(logsPath)
			if err != nil {
				return err
			}
			err = optionalTestReporter.SendSlackNotification(nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
