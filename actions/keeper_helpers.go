package actions

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"

	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
)

// DeployKeeperContracts deploys a number of basic keeper contracts with an update interval of 5
func DeployKeeperContracts(
	numberOfContracts int,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	chainlinkNodes []client.Chainlink,
	networks *client.Networks,
) (contracts.KeeperRegistry, []contracts.KeeperConsumer) {
	defaultNetwork := networks.Default
	keeperConsumerContracts := make([]contracts.KeeperConsumer, 0)

	defaultRegistrySettings := &contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    uint32(200000000),
		BlockCountPerTurn:    big.NewInt(3),
		CheckGasLimit:        uint32(2500000),
		StalenessSeconds:     big.NewInt(90000),
		GasCeilingMultiplier: uint16(1),
		FallbackGasPrice:     big.NewInt(2e11),
		FallbackLinkPrice:    big.NewInt(2e18),
	}

	registry, registrar := prepKeeperDeployments(
		numberOfContracts,
		linkToken,
		contractDeployer,
		defaultNetwork,
		defaultRegistrySettings,
	)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumer(big.NewInt(5))
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		keeperConsumerContracts = append(keeperConsumerContracts, keeperConsumerInstance)
		log.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
		if (contractCount+1)%contractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = defaultNetwork.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumer deployments")
		}
	}
	err := defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	// Register contracts to registrar
	// TODO: Could be simplified with generics?
	for contractCount, keeperConsumerInstance := range keeperConsumerContracts {
		req, err := registrar.EncodeRegisterRequest(
			fmt.Sprintf("upkeep_%d", contractCount+1),
			[]byte("0x1234"),
			keeperConsumerInstance.Address(),
			defaultRegistrySettings.CheckGasLimit,
			keeperConsumerInstance.Address(),
			[]byte("0x"),
			big.NewInt(9e18),
			0,
		)
		Expect(err).ShouldNot(HaveOccurred(), "Encoding the register request shouldn't fail")
		err = linkToken.TransferAndCall(registrar.Address(), big.NewInt(9e18), req)
		Expect(err).ShouldNot(HaveOccurred(), "Error registering the upkeep consumer to the registrar")
		log.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Registered Keeper Consumer Contract")
		if (contractCount+1)%contractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = defaultNetwork.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait after registering upkeep consumers")
		}
	}
	err = defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed while waiting for all consumer contracts to be registered to registrar")
	log.Info().Msg("Successfully registered all Keeper Consumer Contracts")

	return registry, keeperConsumerContracts
}

// DeployPerformanceKeeperContracts deploys a set amount of keeper performance contracts registered to a single registry
func DeployPerformanceKeeperContracts(
	numberOfContracts int,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	chainlinkNodes []client.Chainlink,
	networks *client.Networks,
	keeperContractSettings *contracts.KeeperRegistrySettings,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
) (contracts.KeeperRegistry, []contracts.KeeperConsumerPerformance) {
	defaultNetwork := networks.Default
	keeperConsumerContracts := make([]contracts.KeeperConsumerPerformance, 0)

	registry, registrar := prepKeeperDeployments(
		numberOfContracts,
		linkToken,
		contractDeployer,
		defaultNetwork,
		keeperContractSettings,
	)

	linkFunds := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(blockRange/blockInterval))
	// Deploy all the KeeperConsumerPerformance contracts
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumerPerformance(
			big.NewInt(blockRange),
			big.NewInt(blockInterval),
			big.NewInt(checkGasToBurn),
			big.NewInt(performGasToBurn),
		)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumerPerformance instance %d shouldn't fail", contractCount+1)
		keeperConsumerContracts = append(keeperConsumerContracts, keeperConsumerInstance)
		log.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Performance Contract")
		if (contractCount+1)%contractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = defaultNetwork.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumerPerformance deployments")
		}
	}
	err := defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	for contractCount, keeperConsumerInstance := range keeperConsumerContracts {
		// Register Consumer to registrar
		req, err := registrar.EncodeRegisterRequest(
			fmt.Sprintf("upkeep_%d", contractCount+1),
			[]byte("0x1234"),
			keeperConsumerInstance.Address(),
			keeperContractSettings.CheckGasLimit,
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
			Int("Out Of", numberOfContracts).
			Msg("Registered Keeper Performance Contract")
		if (contractCount+1)%500 == 0 { // For large amounts of contract deployments, space things out some
			err = defaultNetwork.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait after registering upkeep consumers")
		}
	}
	err = defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed while waiting for all consumer contracts to be registered to registrar")
	log.Info().Msg("Successfully registered all Keeper Consumer Contracts")

	return registry, keeperConsumerContracts
}

func CreateKeeperJobs(chainlinkNodes []client.Chainlink, keeperRegistry contracts.KeeperRegistry) {
	// Send keeper jobs to registry and chainlink nodes
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
	err = keeperRegistry.SetKeepers(nodeAddressesStr, payees)
	Expect(err).ShouldNot(HaveOccurred(), "Setting keepers in the registry shouldn't fail")

	for _, chainlinkNode := range chainlinkNodes {
		chainlinkNodeAddress, err := chainlinkNode.PrimaryEthAddress()
		Expect(err).ShouldNot(HaveOccurred(), "Error retrieving chainlink node address")
		_, err = chainlinkNode.CreateJob(&client.KeeperJobSpec{
			Name:                     fmt.Sprintf("keeper-test-%s", keeperRegistry.Address()),
			ContractAddress:          keeperRegistry.Address(),
			FromAddress:              chainlinkNodeAddress,
			MinIncomingConfirmations: 1,
			ObservationSource:        client.ObservationSourceKeeperDefault(),
		})
		Expect(err).ShouldNot(HaveOccurred(), "Creating KeeperV2 Job shouldn't fail")
	}
}

func prepKeeperDeployments(
	numberOfContracts int,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	defaultNetwork client.BlockchainClient,
	keeperContractSettings *contracts.KeeperRegistrySettings,
) (contracts.KeeperRegistry, contracts.UpkeepRegistrar) {
	// Deploy Preliminary contracts (Registry, Registrar, and mock feeds)
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	err = defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for mock feeds to deploy")

	registry, err := contractDeployer.DeployKeeperRegistry(
		&contracts.KeeperRegistryOpts{
			LinkAddr:             linkToken.Address(),
			ETHFeedAddr:          ef.Address(),
			GasFeedAddr:          gf.Address(),
			PaymentPremiumPPB:    keeperContractSettings.PaymentPremiumPPB,
			BlockCountPerTurn:    keeperContractSettings.BlockCountPerTurn,
			CheckGasLimit:        keeperContractSettings.CheckGasLimit,
			StalenessSeconds:     keeperContractSettings.StalenessSeconds,
			GasCeilingMultiplier: keeperContractSettings.GasCeilingMultiplier,
			FallbackGasPrice:     keeperContractSettings.FallbackGasPrice,
			FallbackLinkPrice:    keeperContractSettings.FallbackLinkPrice,
		},
	)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying keeper registry shouldn't fail")
	err = defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for keeper registry to deploy")

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfContracts))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")
	// Deploy and configure the UpkeepRegistrar
	registrar, err := contractDeployer.DeployUpkeepRegistrationRequests(
		linkToken.Address(),
		big.NewInt(0),
	)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying UpkeepRegistrationRequests contract shouldn't fail")
	err = defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registrar to deploy")
	err = registry.SetRegistrar(registrar.Address())
	Expect(err).ShouldNot(HaveOccurred(), "Registering the registrar address on the registry shouldn't fail")
	err = defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registry to set registrar")

	err = registrar.SetRegistrarConfig(
		true,
		uint32(6000000),
		uint16(numberOfContracts),
		registry.Address(),
		big.NewInt(0),
	)
	Expect(err).ShouldNot(HaveOccurred(), "Setting the registrar configuration shouldn't fail")
	err = defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registrar and supporting contract deployments")

	return registry, registrar
}
