package actions

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
)

// This actions file often returns functions, rather than just values. These are used as common test helpers, and are
// handy to have returning as functions so that Ginkgo can use them in an aesthetically pleasing way.

// DeployKeeperConsumerPerformanceContracts deploys and funds the following
// 1 KeeperRegistry contract
// 1 each of a mock ETH/LINK and GAS feeds
// A numberOfContracts amount of KeeperConsumerPerformance contracts
// The KeeperConsumerPerformance contracts are all registered with the single KeeperRegistry, and then standard keeper
// jobs are sent to all chainlink nodes to upkeep all the consumer contracts
func DeployKeeperConsumerPerformanceContracts(
	numberOfContracts int,
	testBlockRange *big.Int,
	averageCadence *big.Int,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	chainlinkNodes []client.Chainlink,
	networks *client.Networks,
) []contracts.KeeperConsumerPerformance {
	checkGasLimit := uint32(2500000) // Default
	// Deploy Preliminary contracts (Registry, Registrar, and mock feeds)
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	registry, err := contractDeployer.DeployKeeperRegistry(
		&contracts.KeeperRegistryOpts{
			LinkAddr:             linkTokenContract.Address(),
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
	err = linkTokenContract.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfContracts))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")
	// Deploy and configure the UpkeepRegistrar
	registrar, err := contractDeployer.DeployUpkeepRegistrationRequests(
		linkTokenContract.Address(),
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
	var keeperConsumerInstances []contracts.KeeperConsumerPerformance
	for i := 0; i < numberOfContracts; i++ {
		// Deploy consumer
		keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumerPerformance(
			testBlockRange,
			averageCadence,
		)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumerPerformance instance %d shouldn't fail", i+1)
		keeperConsumerInstances = append(keeperConsumerInstances, keeperConsumerInstance)
		Expect(err).ShouldNot(HaveOccurred())
		err = linkTokenContract.Transfer(keeperConsumerInstance.Address(), big.NewInt(1e18))
		Expect(err).ShouldNot(HaveOccurred(), "Transfering LINK token to KeeperConsumerPerformance instance %d shouldn't fail", i+1)

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
		err = linkTokenContract.TransferAndCall(registrar.Address(), big.NewInt(1e18), req)
		Expect(err).ShouldNot(HaveOccurred(), "Error registering the upkeep consumer to the registrar")
	}
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed while waiting for all KeeperConsumerPerformance contracts to deploy")

	// Send keeper jobs to the keeper registry
	primaryNode := chainlinkNodes[0]
	primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
	Expect(err).ShouldNot(HaveOccurred(), "Reading ETH Keys from Chainlink Client shouldn't fail")
	nodeAddresses, err := ChainlinkNodeAddresses(chainlinkNodes)
	Expect(err).ShouldNot(HaveOccurred(), "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddress)
	}
	err = registry.SetKeepers(nodeAddressesStr, payees)
	Expect(err).ShouldNot(HaveOccurred(), "Setting keepers in the registry shouldn't fail")
	_, err = primaryNode.CreateJob(&client.KeeperJobSpec{
		Name:                     "keeper-test-job",
		ContractAddress:          registry.Address(),
		FromAddress:              primaryNodeAddress,
		MinIncomingConfirmations: 1,
		ObservationSource:        client.ObservationSourceKeeperDefault(),
	})
	Expect(err).ShouldNot(HaveOccurred(), "Creating KeeperV2 Job shouldn't fail")
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")

	return keeperConsumerInstances
}
