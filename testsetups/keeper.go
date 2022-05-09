package testsetups

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"
	"time"

	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/helmenv/environment"
)

// KeeperBlockTimeTest builds a test to check that chainlink nodes are able to upkeep a specified amount of Upkeep
// contracts within a certain block time
type KeeperBlockTimeTest struct {
	Inputs       KeeperBlockTimeTestInputs
	TestReporter testreporters.KeeperBlockTimeTestReporter

	keeperRegistry          contracts.KeeperRegistry
	keeperConsumerContracts []contracts.KeeperConsumerPerformance
	mockServer              *client.MockserverClient

	env            *environment.Environment
	chainlinkNodes []client.Chainlink
	networks       *blockchain.Networks
	defaultNetwork blockchain.EVMClient
}

// KeeperBlockTimeTestInputs are all the required inputs for a Keeper Block Time Test
type KeeperBlockTimeTestInputs struct {
	NumberOfContracts      int                               // Number of upkeep contracts
	KeeperRegistrySettings *contracts.KeeperRegistrySettings // Settings of each keeper contract
	Timeout                time.Duration                     // Timeout for the test
	BlockRange             int64                             // How many blocks to run the test for
	BlockInterval          int64                             // Interval of blocks that upkeeps are expected to be performed
	CheckGasToBurn         int64                             // How much gas should be burned on checkUpkeep() calls
	PerformGasToBurn       int64                             // How much gas should be burned on performUpkeep() calls
	ChainlinkNodeFunding   *big.Float                        // Amount of ETH to fund each chainlink node with
}

// NewKeeperBlockTimeTest prepares a new keeper block time test to be run
func NewKeeperBlockTimeTest(inputs KeeperBlockTimeTestInputs) *KeeperBlockTimeTest {
	return &KeeperBlockTimeTest{
		Inputs: inputs,
	}
}

// Setup prepares contracts for the test
func (k *KeeperBlockTimeTest) Setup(env *environment.Environment) {
	startTime := time.Now()
	k.ensureInputValues()
	k.env = env
	inputs := k.Inputs
	var err error

	// Connect to networks and prepare for contract deployment
	networkRegistry := blockchain.NewSoakNetworkRegistry()
	k.networks, err = networkRegistry.GetNetworks(k.env)
	Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
	k.defaultNetwork = k.networks.Default
	contractDeployer, err := contracts.NewContractDeployer(k.defaultNetwork)
	Expect(err).ShouldNot(HaveOccurred(), "Building a new contract deployer shouldn't fail")
	k.chainlinkNodes, err = client.ConnectChainlinkNodesSoak(k.env)
	Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
	k.defaultNetwork.ParallelTransactions(true)

	// Fund chainlink nodes
	err = actions.FundChainlinkNodes(k.chainlinkNodes, k.defaultNetwork, k.Inputs.ChainlinkNodeFunding)
	Expect(err).ShouldNot(HaveOccurred(), "Funding Chainlink nodes shouldn't fail")
	linkToken, err := contractDeployer.DeployLinkTokenContract()
	Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
	err = k.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for LINK Contract deployment")

	k.keeperRegistry, k.keeperConsumerContracts = actions.DeployPerformanceKeeperContracts(
		inputs.NumberOfContracts,
		linkToken,
		contractDeployer,
		k.networks,
		k.Inputs.KeeperRegistrySettings,
		inputs.BlockRange,
		inputs.BlockInterval,
		inputs.CheckGasToBurn,
		inputs.PerformGasToBurn,
	)

	// Send keeper jobs to registry and chainlink nodes
	actions.CreateKeeperJobs(k.chainlinkNodes, k.keeperRegistry)

	log.Info().Str("Setup Time", time.Since(startTime).String()).Msg("Finished Keeper Block Time Test Setup")
}

// Run runs the keeper block time test
func (k *KeeperBlockTimeTest) Run() {
	startTime := time.Now()

	for index, keeperConsumer := range k.keeperConsumerContracts {
		k.defaultNetwork.AddHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d", index),
			contracts.NewKeeperConsumerPerformanceRoundConfirmer(
				keeperConsumer,
				k.Inputs.BlockInterval,
				k.Inputs.BlockRange,
				&k.TestReporter,
			),
		)
	}
	defer func() { // Cleanup the subscriptions
		for index := range k.keeperConsumerContracts {
			k.defaultNetwork.DeleteHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d", index))
		}
	}()
	err := k.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for keeper subscriptions")

	for _, chainlinkNode := range k.chainlinkNodes {
		txData, err := chainlinkNode.ReadTransactionAttempts()
		Expect(err).ShouldNot(HaveOccurred(), "Error retrieving transaction data from chainlink node")
		k.TestReporter.AttemptedChainlinkTransactions = append(k.TestReporter.AttemptedChainlinkTransactions, txData)
	}

	log.Info().Str("Run Time", time.Since(startTime).String()).Msg("Finished Keeper Block Time Test")
}

// Networks returns the networks that the test is running on
func (k *KeeperBlockTimeTest) TearDownVals() (*environment.Environment, *blockchain.Networks, []client.Chainlink, testreporters.TestReporter) {
	return k.env, k.networks, k.chainlinkNodes, &k.TestReporter
}

// ensureValues ensures that all values needed to run the test are present
func (k *KeeperBlockTimeTest) ensureInputValues() {
	inputs := k.Inputs
	Expect(inputs.NumberOfContracts).Should(BeNumerically(">=", 1), "Expecting at least 1 keeper contracts")
	if inputs.Timeout == 0 {
		Expect(inputs.BlockRange).Should(BeNumerically(">", 0), "If no `timeout` is provided, a `testBlockRange` is required")
	} else if inputs.BlockRange <= 0 {
		Expect(inputs.Timeout).Should(BeNumerically(">=", 1), "If no `testBlockRange` is provided a `timeout` is required")
	}
	Expect(inputs.KeeperRegistrySettings).ShouldNot(BeNil(), "You need to set KeeperRegistrySettings")
	Expect(k.Inputs.ChainlinkNodeFunding).ShouldNot(BeNil(), "You need to set a funding amount for chainlink nodes")
	clFunds, _ := k.Inputs.ChainlinkNodeFunding.Float64()
	Expect(clFunds).Should(BeNumerically(">=", 0), "Expecting Chainlink node funding to be more than 0 ETH")
	Expect(inputs.CheckGasToBurn).Should(BeNumerically(">", 0), "You need to set an expected amount of gas to burn on checkUpkeep()")
	Expect(inputs.KeeperRegistrySettings.CheckGasLimit).Should(BeNumerically(">=", inputs.CheckGasToBurn),
		"CheckGasLimit should be >= CheckGasToBurn")
	Expect(inputs.PerformGasToBurn).Should(BeNumerically(">", 0), "You need to set an expected amount of gas to burn on performUpkeep()")
}
