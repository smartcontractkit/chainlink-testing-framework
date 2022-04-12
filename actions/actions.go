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

	"github.com/avast/retry-go"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/helmenv/environment"
	"golang.org/x/sync/errgroup"
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
	// After how many contract actions to wait before starting any more
	// Example: When deploying 1000 contracts, stop every contractDeploymentInterval have been deployed to wait before continuing
	contractDeploymentInterval int = 500
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

// FundChainlinkNodes will fund all of the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodesLink(
	nodes []client.Chainlink,
	blockchain client.BlockchainClient,
	linkToken contracts.LinkToken,
	linkAmount *big.Int,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		err = linkToken.Transfer(toAddress, linkAmount)
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
// specified path. Can also accept a testreporter (if one was used) to log further results
func TeardownSuite(
	env *environment.Environment,
	nets *client.Networks,
	logsFolderPath string,
	chainlinkNodes []client.Chainlink,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
) error {
	if err := writeTeardownLogs(env, optionalTestReporter); err != nil {
		return errors.Wrap(err, "Error dumping environment logs, leaving environment running for manual retrieval")
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

	if nets != nil && chainlinkNodes != nil && len(chainlinkNodes) > 0 {
		if err := returnFunds(chainlinkNodes, nets); err != nil {
			log.Error().Err(err).Str("Namespace", env.Namespace).
				Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
					"Environment is left running so you can try manually!")
			env.Persistent = true
		}
	} else {
		log.Info().Msg("Successfully returned funds from chainlink nodes to default network wallets")
	}
	if nets != nil {
		if err := nets.Teardown(); err != nil {
			return err
		}
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
func TeardownRemoteSuite(
	env *environment.Environment,
	nets *client.Networks,
	chainlinkNodes []client.Chainlink,
	optionalTestReporter testreporters.TestReporter, // Optionally pass in a test reporter to log further metrics
) error {
	err := writeTeardownLogs(env, optionalTestReporter)
	if err != nil {
		return err
	}
	err = returnFunds(chainlinkNodes, nets)
	if err != nil {
		log.Error().Err(err).Str("Namespace", env.Namespace).
			Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
				"Environment is left running so you can try manually!")
	}
	return err
}

// attempts to download the logs of all ephemeral test deployments onto the test runner, also writing a test report
// if one is provided
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

// Returns all the funds from the chainlink nodes to the networks default address
func returnFunds(chainlinkNodes []client.Chainlink, networks *client.Networks) error {
	if networks == nil {
		log.Warn().Msg("No network connections found, unable to return funds from chainlink nodes.")
	}
	log.Info().Msg("Attempting to return Chainlink node funds to default network wallets")
	for _, network := range networks.AllNetworks() {
		if network.GetNetworkType() == client.SimulatedEthNetwork {
			log.Info().Str("Network Name", network.GetNetworkName()).
				Msg("Network is a `eth_simulated` network. Skipping fund return.")
			continue
		}
		addressMap, err := sendFunds(chainlinkNodes, network)
		if err != nil {
			return err
		}

		err = checkFunds(chainlinkNodes, addressMap, strings.ToLower(network.GetDefaultWallet().Address()))
		if err != nil {
			return err
		}
	}

	return nil
}

// Requests that all the chainlink nodes send their funds back to the network's default wallet
// This is surprisingly tricky, and fairly annoying due to Go's lack of syntactic sugar and how chainlink nodes handle txs
func sendFunds(chainlinkNodes []client.Chainlink, network client.BlockchainClient) (map[int]string, error) {
	chainlinkTransactionAddresses := make(map[int]string)
	sendFundsErrGroup := new(errgroup.Group)
	for ni, n := range chainlinkNodes {
		nodeIndex := ni // https://golang.org/doc/faq#closures_and_goroutines
		node := n
		// Send async request to each chainlink node to send a transaction back to the network default wallet
		sendFundsErrGroup.Go(
			func() error {
				primaryEthKeyData, err := node.ReadPrimaryETHKey()
				if err != nil {
					// TODO: Support non-EVM chain fund returns
					if strings.Contains(err.Error(), "No ETH keys present") {
						log.Warn().Msg("Not returning any funds. Only support ETH chains for fund returns at the moment")
						return nil
					}
					return err
				}

				nodeBalanceString := primaryEthKeyData.Attributes.ETHBalance
				if nodeBalanceString != "0" { // If key has a non-zero balance, attempt to transfer it back
					gasCost, err := network.EstimateTransactionGasCost()
					if err != nil {
						return err
					}

					// TODO: Imperfect gas calculation buffer of 50 Gwei. Seems to be the result of differences in chainlink
					// gas handling. Working with core team on a better solution
					gasCost = gasCost.Add(gasCost, big.NewInt(50000000000))
					nodeBalance, _ := big.NewInt(0).SetString(nodeBalanceString, 10)
					transferAmount := nodeBalance.Sub(nodeBalance, gasCost)
					_, err = node.SendNativeToken(transferAmount, primaryEthKeyData.Attributes.Address, network.GetDefaultWallet().Address())
					if err != nil {
						return err
					}
					// Add the address to our map to check for later (hashes aren't returned, sadly)
					chainlinkTransactionAddresses[nodeIndex] = strings.ToLower(primaryEthKeyData.Attributes.Address)
				}
				return nil
			},
		)

	}
	return chainlinkTransactionAddresses, sendFundsErrGroup.Wait()
}

// checks that the funds made it from the chainlink node to the network address
// this turns out to be tricky to do, given how chainlink handles pending transactions, thus the complexity
func checkFunds(chainlinkNodes []client.Chainlink, sentFromAddressesMap map[int]string, toAddress string) error {
	err := retry.Do( // Might take some time for txs to confirm, check up on the nodes a few times
		func() error {
			log.Debug().Msg("Attempting to confirm chainlink nodes transferred back funds")
			transactionErrGroup := new(errgroup.Group)
			for nodeIndex, n := range chainlinkNodes {
				node := n // https://golang.org/doc/faq#closures_and_goroutines
				sentFromAddress, nodeHasFunds := sentFromAddressesMap[nodeIndex]
				// Async check on all the nodes if their transactions are confirmed
				if nodeHasFunds { // Only if the node had funds to begin with
					transactionErrGroup.Go(func() error {
						return confirmTransaction(node, sentFromAddress, toAddress, transactionErrGroup)
					})
				} else {
					log.Debug().Int("Node Number", nodeIndex).Msg("Chainlink node had no funds to return")
				}
			}

			return transactionErrGroup.Wait()
		},
		retry.Delay(time.Second*5),
		retry.MaxDelay(time.Second*5),
		retry.Attempts(20),
	)

	return err
}

// helper to confirm that the latest attempted transaction on the chainlink node with the expected from and to addresses
// has been confirmed
func confirmTransaction(
	chainlinkNode client.Chainlink,
	fromAddress string,
	toAddress string,
	transactionErrGroup *errgroup.Group,
) error {
	transactionAttempts, err := chainlinkNode.ReadTransactionAttempts()
	if err != nil {
		return err
	}
	log.Debug().Str("From", fromAddress).
		Str("To", toAddress).
		Msg("Attempting to confirm node returned funds")
	// Loop through all transactions on the node
	for _, tx := range transactionAttempts.Data {
		if tx.Attributes.From == fromAddress && strings.ToLower(tx.Attributes.To) == toAddress {
			if tx.Attributes.State == "confirmed" {
				return nil
			}
			return fmt.Errorf("Expected transaction to be confirmed. From: %s To: %s State: %s", fromAddress, toAddress, tx.Attributes.State)
		}
	}
	return fmt.Errorf("Did not find expected transaction on node. From: %s To: %s", fromAddress, toAddress)
}
