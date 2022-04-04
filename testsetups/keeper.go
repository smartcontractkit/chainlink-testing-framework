package testsetups

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/testreporters"
)

// KeeperBlockTimeTest builds a test to check that chainlink nodes are able to upkeep a specified amount of Upkeep
// contracts within a certain block time
type KeeperBlockTimeTest struct {
	Inputs       KeeperBlockTimeTestInputs
	TestReporter testreporters.KeeperBlockTimeTestReporter

	keeperConsumerContracts []contracts.KeeperConsumerPerformance
	mockServer              *client.MockserverClient

	env            *environment.Environment
	chainlinkNodes []client.Chainlink
	networks       *client.Networks
	defaultNetwork client.BlockchainClient
}

// KeeperBlockTimeTestInputs are all the required inputs for a Keeper Block Time Test
type KeeperBlockTimeTestInputs struct {
	NumberOfContracts      int                     // Number of upkeep contracts
	KeeperContractSettings *KeeperContractSettings // Settings of each keeper contract
	Timeout                time.Duration           // Timeout for the test
	BlockRange             int64                   // How many blocks to run the test for
	BlockInterval          int64                   // Interval of blocks that upkeeps are expected to be performed
	ChainlinkNodeFunding   *big.Float              // Amount of ETH to fund each chainlink node with
	CheckGasLimit          uint32                  // Max amount of gas that checkUpkeep uses for off-chain computation
}

// KeeperContractSettings represents the fine tuning settings for each upkeep contract
type KeeperContractSettings struct {
	PaymentPremiumPPB    uint32   // payment premium rate oracles receive on top of being reimbursed for gas, measured in parts per billion
	BlockCountPerTurn    *big.Int // number of blocks each oracle has during their turn to perform upkeep before it will be the next keeper's turn to submit
	CheckGasLimit        uint32   // gas limit when checking for upkeep
	StalenessSeconds     *big.Int // number of seconds that is allowed for feed data to be stale before switching to the fallback pricing
	GasCeilingMultiplier uint16   // multiplier to apply to the fast gas feed price when calculating the payment ceiling for keepers
	FallbackGasPrice     *big.Int // gas price used if the gas price feed is stale
	FallbackLinkPrice    *big.Int // LINK price used if the LINK price feed is stale
}

// NewKeeperBlockTimeTest prepares a new keeper block time test to be run
func NewKeeperBlockTimeTest(inputs KeeperBlockTimeTestInputs) *KeeperBlockTimeTest {
	return &KeeperBlockTimeTest{
		Inputs: inputs,
	}
}

// Setup prepares contracts for the test
func (k *KeeperBlockTimeTest) Setup(env *environment.Environment) {
	k.ensureInputValues()
	k.env = env
	inputs := k.Inputs
	var err error

	// Connect to networks and prepare for contract deployment
	networkRegistry := client.NewSoakNetworkRegistry()
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
	// Edge case where simulated networks need some funds at the 0x0 address in order for keeper reads to work
	if k.defaultNetwork.GetNetworkType() == "eth_simulated" {
		err = actions.FundAddresses(k.defaultNetwork, big.NewFloat(1), "0x0")
		Expect(err).ShouldNot(HaveOccurred())
	}
	linkToken, err := contractDeployer.DeployLinkTokenContract()
	Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
	err = k.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for LINK Contract deployment")

	// Deploy Preliminary contracts (Registry, Registrar, and mock feeds)
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	err = k.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for mock feeds to deploy")

	registry, err := contractDeployer.DeployKeeperRegistry(
		&contracts.KeeperRegistryOpts{
			LinkAddr:             linkToken.Address(),
			ETHFeedAddr:          ef.Address(),
			GasFeedAddr:          gf.Address(),
			PaymentPremiumPPB:    k.Inputs.KeeperContractSettings.PaymentPremiumPPB,
			BlockCountPerTurn:    k.Inputs.KeeperContractSettings.BlockCountPerTurn,
			CheckGasLimit:        k.Inputs.KeeperContractSettings.CheckGasLimit,
			StalenessSeconds:     k.Inputs.KeeperContractSettings.StalenessSeconds,
			GasCeilingMultiplier: k.Inputs.KeeperContractSettings.GasCeilingMultiplier,
			FallbackGasPrice:     k.Inputs.KeeperContractSettings.FallbackGasPrice,
			FallbackLinkPrice:    k.Inputs.KeeperContractSettings.FallbackLinkPrice,
		},
	)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying keeper registry shouldn't fail")
	err = k.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for keeper registry to deploy")

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(inputs.NumberOfContracts))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")
	// Deploy and configure the UpkeepRegistrar
	registrar, err := contractDeployer.DeployUpkeepRegistrationRequests(
		linkToken.Address(),
		big.NewInt(0),
	)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying UpkeepRegistrationRequests contract shouldn't fail")
	err = k.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registrar to deploy")
	err = registry.SetRegistrar(registrar.Address())
	Expect(err).ShouldNot(HaveOccurred(), "Registering the registrar address on the registry shouldn't fail")
	err = k.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registry to set registrar")

	err = registrar.SetRegistrarConfig(
		true,
		uint32(6000000),
		uint16(k.Inputs.NumberOfContracts),
		registry.Address(),
		big.NewInt(0),
	)
	Expect(err).ShouldNot(HaveOccurred(), "Setting the registrar configuration shouldn't fail")
	err = k.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registrar and supporting contract deployments")

	linkFunds := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(k.Inputs.BlockRange/k.Inputs.BlockInterval))
	// Deploy all the KeeperConsumerPerformance contracts
	for contractCount := 0; contractCount < inputs.NumberOfContracts; contractCount++ {
		// Deploy consumer
		keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumerPerformance(
			big.NewInt(inputs.BlockRange),
			big.NewInt(inputs.BlockInterval),
		)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumerPerformance instance %d shouldn't fail", contractCount+1)
		k.keeperConsumerContracts = append(k.keeperConsumerContracts, keeperConsumerInstance)
		log.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", inputs.NumberOfContracts).
			Msg("Deployed Keeper Performance Contract")
		if contractCount+1%100 == 0 { // For large amounts of contract deployments, space things out some
			k.defaultNetwork.WaitForEvents()
		}
	}
	err = k.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	for contractCount, keeperConsumerInstance := range k.keeperConsumerContracts {
		// Register Consumer to registrar
		req, err := registrar.EncodeRegisterRequest(
			fmt.Sprintf("upkeep_%d", contractCount+1),
			[]byte("0x1234"),
			keeperConsumerInstance.Address(),
			k.Inputs.KeeperContractSettings.CheckGasLimit,
			keeperConsumerInstance.Address(),
			[]byte("0x"),
			linkFunds,
			0,
		)
		Expect(err).ShouldNot(HaveOccurred(), "Encoding the register request shouldn't fail")
		err = linkToken.TransferAndCall(registrar.Address(), linkFunds, req)
		Expect(err).ShouldNot(HaveOccurred(), "Error registering the upkeep consumer to the registrar")
		log.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", inputs.NumberOfContracts).
			Msg("Registered Keeper Performance Contract")
		if (contractCount+1)%100 == 0 { // For large amounts of contract deployments, space things out some
			k.defaultNetwork.WaitForEvents()
		}
	}
	err = k.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed while waiting for all consumer contracts to be registered to registrar")
	log.Info().Msg("Successfully registered all Keeper Consumer Contracts")

	// Send keeper jobs to registry and chainlink nodes
	primaryNode := k.chainlinkNodes[0]
	primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
	Expect(err).ShouldNot(HaveOccurred(), "Reading ETH Keys from Chainlink Client shouldn't fail")
	nodeAddresses, err := chainlinkNodeAddresses(k.chainlinkNodes)
	Expect(err).ShouldNot(HaveOccurred(), "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddress)
	}
	err = registry.SetKeepers(nodeAddressesStr, payees)
	Expect(err).ShouldNot(HaveOccurred(), "Setting keepers in the registry shouldn't fail")
	for _, chainlinkNode := range k.chainlinkNodes {
		chainlinkNodeAddress, err := chainlinkNode.PrimaryEthAddress()
		Expect(err).ShouldNot(HaveOccurred(), "Error retrieving chainlink node address")
		_, err = chainlinkNode.CreateJob(&client.KeeperJobSpec{
			Name:                     fmt.Sprintf("keeper-test-%s", registry.Address()),
			ContractAddress:          registry.Address(),
			FromAddress:              chainlinkNodeAddress,
			MinIncomingConfirmations: 1,
			ObservationSource:        client.ObservationSourceKeeperDefault(),
		})
		Expect(err).ShouldNot(HaveOccurred(), "Creating KeeperV2 Job shouldn't fail")
	}
}

// Run runs the keeper block time test
func (k *KeeperBlockTimeTest) Run() {
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
}

// Networks returns the networks that the test is running on
func (k *KeeperBlockTimeTest) TearDownVals() (*environment.Environment, *client.Networks, []client.Chainlink, testreporters.TestReporter) {
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
	Expect(inputs.KeeperContractSettings).ShouldNot(BeNil(), "You need to set KeeperContractSettings")
	Expect(k.Inputs.ChainlinkNodeFunding).ShouldNot(BeNil(), "You need to set a funding amount for chainlink nodes")
	clFunds, _ := k.Inputs.ChainlinkNodeFunding.Float64()
	Expect(clFunds).Should(BeNumerically(">=", 0), "Expecting Chainlink node funding to be more than 0 ETH")
}

// chainlinkNodeAddresses will return all the on-chain wallet addresses for a set of Chainlink nodes
func chainlinkNodeAddresses(nodes []client.Chainlink) ([]common.Address, error) {
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
