package testsetups

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/testreporters"
)

// KeeperBlockTimeTest builds a test to check that chainlink nodes are able to upkeep a specified amount of Upkeep
// contracts within a certain block time
type KeeperBlockTimeTest struct {
	Inputs KeeperBlockTimeTestInputs

	TestReporter            testreporters.KeeperBlockTimeTestReporter
	keeperConsumerContracts []contracts.KeeperConsumerPerformance
}

// KeeperBlockTimeTestInputs are all the required inputs for a Keeper Block Time Test
type KeeperBlockTimeTestInputs struct {
	NumberOfContracts int           // Number of upkeep contracts
	Timeout           time.Duration // Timeout for the test
	BlockRange        int64
	BlockInterval     int64
	ContractDeployer  contracts.ContractDeployer
	ChainlinkNodes    []client.Chainlink
	Networks          *client.Networks
	LinkTokenContract contracts.LinkToken
}

func NewKeeperBlockTimeTest(inputs KeeperBlockTimeTestInputs) *KeeperBlockTimeTest {
	return &KeeperBlockTimeTest{
		Inputs: inputs,
	}
}

func (k *KeeperBlockTimeTest) Setup() {
	k.ensureInputValues()
	inputs := k.Inputs
	checkGasLimit := uint32(2500000) // Default
	// Deploy Preliminary contracts (Registry, Registrar, and mock feeds)
	ef, err := inputs.ContractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := inputs.ContractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	registry, err := inputs.ContractDeployer.DeployKeeperRegistry(
		&contracts.KeeperRegistryOpts{
			LinkAddr:             inputs.LinkTokenContract.Address(),
			ETHFeedAddr:          ef.Address(),
			GasFeedAddr:          gf.Address(),
			PaymentPremiumPPB:    uint32(200000000),
			BlockCountPerTurn:    big.NewInt(3),
			CheckGasLimit:        checkGasLimit,
			StalenessSeconds:     big.NewInt(90000),
			GasCeilingMultiplier: uint16(1),
			FallbackGasPrice:     big.NewInt(2e11),
			FallbackLinkPrice:    big.NewInt(2e18),
		},
	)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying keeper registry shouldn't fail")
	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = inputs.LinkTokenContract.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(inputs.NumberOfContracts))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")
	// Deploy and configure the UpkeepRegistrar
	registrar, err := inputs.ContractDeployer.DeployUpkeepRegistrationRequests(
		inputs.LinkTokenContract.Address(),
		big.NewInt(0),
	)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying UpkeepRegistrationRequests contract shouldn't fail")
	err = registry.SetRegistrar(registrar.Address())
	Expect(err).ShouldNot(HaveOccurred(), "Registering the registrar address on the registry shouldn't fail")
	err = registrar.SetRegistrarConfig(
		true,
		uint32(999),
		uint16(999),
		registry.Address(),
		big.NewInt(0),
	)
	Expect(err).ShouldNot(HaveOccurred(), "Setting the registrar configuration shouldn't fail")

	// Deploy all the KeeperConsumerPerformance contracts
	for i := 0; i < inputs.NumberOfContracts; i++ {
		// Deploy consumer
		keeperConsumerInstance, err := inputs.ContractDeployer.DeployKeeperConsumerPerformance(
			big.NewInt(inputs.BlockRange),
			big.NewInt(inputs.BlockInterval),
		)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumerPerformance instance %d shouldn't fail", i+1)
		k.keeperConsumerContracts = append(k.keeperConsumerContracts, keeperConsumerInstance)
		Expect(err).ShouldNot(HaveOccurred())
		// err = k.linkTokenContract.Transfer(keeperConsumerInstance.Address(), big.NewInt(1e18))
		// Expect(err).ShouldNot(HaveOccurred(), "Transfering LINK token to KeeperConsumerPerformance instance %d shouldn't fail", i+1)

		// Register Consumer to registrar
		req, err := registrar.EncodeRegisterRequest(
			fmt.Sprintf("upkeep_%d", i),
			[]byte("0x1234"),
			keeperConsumerInstance.Address(),
			checkGasLimit,
			keeperConsumerInstance.Address(),
			[]byte("0x"),
			big.NewInt(1e18),
			0,
		)
		Expect(err).ShouldNot(HaveOccurred(), "Encoding the register request shouldn't fail")
		err = inputs.LinkTokenContract.TransferAndCall(registrar.Address(), big.NewInt(1e18), req)
		Expect(err).ShouldNot(HaveOccurred(), "Error registering the upkeep consumer to the registrar")
	}
	err = inputs.Networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed while waiting for all KeeperConsumerPerformance contracts to deploy")

	// Send keeper jobs to registry and chainlink nodes
	primaryNode := inputs.ChainlinkNodes[0]
	primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
	Expect(err).ShouldNot(HaveOccurred(), "Reading ETH Keys from Chainlink Client shouldn't fail")
	nodeAddresses, err := chainlinkNodeAddresses(inputs.ChainlinkNodes)
	Expect(err).ShouldNot(HaveOccurred(), "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddress)
	}
	err = registry.SetKeepers(nodeAddressesStr, payees)
	Expect(err).ShouldNot(HaveOccurred(), "Setting keepers in the registry shouldn't fail")
	for _, chainlinkNode := range inputs.ChainlinkNodes {
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

	err = inputs.Networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")
}

func (k *KeeperBlockTimeTest) Run() {
	for index, keeperConsumer := range k.keeperConsumerContracts {
		k.Inputs.Networks.Default.AddHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d", index),
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
			k.Inputs.Networks.Default.DeleteHeaderEventSubscription(fmt.Sprintf("Keeper Tracker %d", index))
		}
	}()
	err := k.Inputs.Networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for keeper subscriptions")

	for _, chainlinkNode := range k.Inputs.ChainlinkNodes {
		txData, err := chainlinkNode.ReadTransactionAttempts()
		Expect(err).ShouldNot(HaveOccurred(), "Error retrieving transaction data from chainlink node")
		k.TestReporter.AttemptedChainlinkTransactions = append(k.TestReporter.AttemptedChainlinkTransactions, txData)
	}
}

// ensureValues ensures that all values needed to run the test are present
func (k *KeeperBlockTimeTest) ensureInputValues() {
	inputs := k.Inputs
	Expect(inputs.NumberOfContracts).Should(BeNumerically(">=", 1), "Expecting at least 1 keeper contracts")
	Expect(inputs.Networks).ShouldNot(BeNil(), "Expected for `networks` to be filled out")
	Expect(inputs.Networks.Default).ShouldNot(BeNil(), "Expected there to be a viable Default network")
	if inputs.Timeout == 0 {
		Expect(inputs.BlockRange).Should(BeNumerically(">", 0), "If no `timeout` is provided, a `testBlockRange` is required")
	} else if inputs.BlockRange <= 0 {
		Expect(inputs.Timeout).Should(BeNumerically(">=", 1), "If no `testBlockRange` is provided a `timeout` is required")
	}
	Expect(inputs.ContractDeployer).ShouldNot(BeNil(), "Expected `contractDeployer` to be provided")
	Expect(len(inputs.ChainlinkNodes)).Should(BeNumerically(">=", 2), "Expecting at least 2 chainlink nodes (recommended 6)")
	Expect(inputs.LinkTokenContract).ShouldNot(BeNil(), "Expecting a valid `linkTokenContract`")
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
