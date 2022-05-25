package actions

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
)

var ZeroAddress = common.Address{}

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
		})
		Expect(err).ShouldNot(HaveOccurred(), "Creating KeeperV2 Job shouldn't fail")
	}
}

// DeployKeeperContracts deploys keeper registry and a number of basic upkeep contracts with an update interval of 5
func DeployKeeperContracts(
	registryVersion ethereum.KeeperRegistryVersion,
	registrySettings contracts.KeeperRegistrySettings,
	numberOfUpkeeps int,
	upkeepGasLimit uint32,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	networks *blockchain.Networks,
) (contracts.KeeperRegistry, []contracts.KeeperConsumer, []*big.Int) {
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for mock feeds to deploy")

	registry := DeployKeeperRegistry(contractDeployer, networks,
		&contracts.KeeperRegistryOpts{
			RegistryVersion: registryVersion,
			LinkAddr:        linkToken.Address(),
			ETHFeedAddr:     ef.Address(),
			GasFeedAddr:     gf.Address(),
			TranscoderAddr:  ZeroAddress.Hex(),
			RegistrarAddr:   ZeroAddress.Hex(),
			Settings:        registrySettings,
		},
	)

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfUpkeeps))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoRegister:     true,
		WindowSizeBlocks: uint32(6000000),
		AllowedPerWindow: uint16(numberOfUpkeeps),
		RegistryAddr:     registry.Address(),
		MinLinkJuels:     big.NewInt(0),
	}
	registrar := DeployKeeperRegistrar(linkToken, registrarSettings, contractDeployer, networks, registry)

	upkeeps := DeployKeeperConsumers(contractDeployer, networks, numberOfUpkeeps)
	upkeepsAddresses := []string{}
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(linkToken, big.NewInt(9e18), networks, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses)

	return registry, upkeeps, upkeepIds
}

// DeployPerformanceKeeperContracts deploys a set amount of keeper performance contracts registered to a single registry
func DeployPerformanceKeeperContracts(
	registryVersion ethereum.KeeperRegistryVersion,
	numberOfContracts int,
	upkeepGasLimit uint32,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	networks *blockchain.Networks,
	registrySettings *contracts.KeeperRegistrySettings,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
) (contracts.KeeperRegistry, []contracts.KeeperConsumerPerformance, []*big.Int) {
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for mock feeds to deploy")

	registry := DeployKeeperRegistry(contractDeployer, networks,
		&contracts.KeeperRegistryOpts{
			RegistryVersion: registryVersion,
			LinkAddr:        linkToken.Address(),
			ETHFeedAddr:     ef.Address(),
			GasFeedAddr:     gf.Address(),
			TranscoderAddr:  ZeroAddress.Hex(),
			RegistrarAddr:   ZeroAddress.Hex(),
			Settings:        *registrySettings,
		},
	)

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfContracts))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoRegister:     true,
		WindowSizeBlocks: uint32(6000000),
		AllowedPerWindow: uint16(numberOfContracts),
		RegistryAddr:     registry.Address(),
		MinLinkJuels:     big.NewInt(0),
	}
	registrar := DeployKeeperRegistrar(linkToken, registrarSettings, contractDeployer, networks, registry)

	upkeeps := DeployKeeperConsumersPerformance(contractDeployer, networks, numberOfContracts, blockRange, blockInterval, checkGasToBurn, performGasToBurn)

	upkeepsAddresses := []string{}
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	linkFunds := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(blockRange/blockInterval))

	upkeepIds := RegisterUpkeepContracts(linkToken, linkFunds, networks, upkeepGasLimit, registry, registrar, numberOfContracts, upkeepsAddresses)

	return registry, upkeeps, upkeepIds
}

func DeployKeeperRegistry(
	contractDeployer contracts.ContractDeployer,
	networks *blockchain.Networks,
	registryOpts *contracts.KeeperRegistryOpts,
) contracts.KeeperRegistry {
	registry, err := contractDeployer.DeployKeeperRegistry(
		registryOpts,
	)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying keeper registry shouldn't fail")
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for keeper registry to deploy")

	return registry
}

func DeployKeeperRegistrar(
	linkToken contracts.LinkToken,
	registrarSettings contracts.KeeperRegistrarSettings,
	contractDeployer contracts.ContractDeployer,
	networks *blockchain.Networks,
	registry contracts.KeeperRegistry,
) contracts.UpkeepRegistrar {
	//#### Deploy and configure the UpkeepRegistrar
	var err error
	registrar, err := contractDeployer.DeployUpkeepRegistrationRequests(
		linkToken.Address(),
		big.NewInt(0),
	)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying UpkeepRegistrationRequests contract shouldn't fail")
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registrar to deploy")
	err = registry.SetRegistrar(registrar.Address())
	Expect(err).ShouldNot(HaveOccurred(), "Registering the registrar address on the registry shouldn't fail")
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registry to set registrar")

	err = registrar.SetRegistrarConfig(
		registrarSettings.AutoRegister,
		registrarSettings.WindowSizeBlocks,
		registrarSettings.AllowedPerWindow,
		registrarSettings.RegistryAddr,
		registrarSettings.MinLinkJuels,
	)
	Expect(err).ShouldNot(HaveOccurred(), "Setting the registrar configuration shouldn't fail")
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registrar and supporting contract deployments")

	return registrar
}

func RegisterUpkeepContracts(
	linkToken contracts.LinkToken,
	linkFunds *big.Int,
	networks *blockchain.Networks,
	upkeepGasLimit uint32,
	registry contracts.KeeperRegistry,
	registrar contracts.UpkeepRegistrar,
	numberOfContracts int,
	upkeepAdresses []string,
) []*big.Int {
	registrationTxHashes := make([]common.Hash, 0)
	upkeepIds := make([]*big.Int, 0)
	for contractCount, upkeepAddress := range upkeepAdresses {
		req, err := registrar.EncodeRegisterRequest(
			fmt.Sprintf("upkeep_%d", contractCount+1),
			[]byte("0x1234"),
			upkeepAddress,
			upkeepGasLimit,
			networks.Default.GetDefaultWallet().Address(), // upkeep Admin
			[]byte("0x"),
			linkFunds,
			0,
		)
		Expect(err).ShouldNot(HaveOccurred(), "Encoding the register request shouldn't fail")
		tx, err := linkToken.TransferAndCall(registrar.Address(), linkFunds, req)
		Expect(err).ShouldNot(HaveOccurred(), "Error registering the upkeep consumer to the registrar")
		log.Debug().
			Str("Contract Address", upkeepAddress).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Str("TxHash", tx.Hash().String()).
			Msg("Registered Keeper Consumer Contract")
		registrationTxHashes = append(registrationTxHashes, tx.Hash())
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait after registering upkeep consumers")
		}
	}
	err := networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed while waiting for all consumer contracts to be registered to registrar")

	// Fetch the upkeep IDs
	for _, txHash := range registrationTxHashes {
		receipt, err := networks.Default.GetTxReceipt(txHash)
		Expect(err).ShouldNot(HaveOccurred(), "Registration tx should be completed")
		var upkeepId *big.Int
		for _, rawLog := range receipt.Logs {
			parsedUpkeepId, err := registry.ParseUpkeepIdFromRegisteredLog(rawLog)
			if err == nil {
				upkeepId = parsedUpkeepId
				break
			}
		}
		Expect(upkeepId).ShouldNot(BeNil(), "Upkeep ID should be found after registration")
		log.Debug().
			Str("TxHash", txHash.String()).
			Str("Upkeep ID", upkeepId.String()).
			Msg("Found upkeepId in tx hash")
		upkeepIds = append(upkeepIds, upkeepId)
	}
	log.Info().Msg("Successfully registered all Keeper Consumer Contracts")
	return upkeepIds
}

func DeployKeeperConsumers(
	contractDeployer contracts.ContractDeployer,
	networks *blockchain.Networks,
	numberOfContracts int,
) []contracts.KeeperConsumer {
	keeperConsumerContracts := make([]contracts.KeeperConsumer, 0)

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
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumer deployments")
		}
	}
	err := networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return keeperConsumerContracts
}

func DeployKeeperConsumersPerformance(
	contractDeployer contracts.ContractDeployer,
	networks *blockchain.Networks,
	numberOfContracts int,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
) []contracts.KeeperConsumerPerformance {
	upkeeps := make([]contracts.KeeperConsumerPerformance, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumerPerformance(
			big.NewInt(blockRange),
			big.NewInt(blockInterval),
			big.NewInt(checkGasToBurn),
			big.NewInt(performGasToBurn),
		)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumerPerformance instance %d shouldn't fail", contractCount+1)
		upkeeps = append(upkeeps, keeperConsumerInstance)
		log.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Performance Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumerPerformance deployments")
		}
	}
	err := networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeeps
}

func DeployUpkeepCounters(
	contractDeployer contracts.ContractDeployer,
	networks *blockchain.Networks,
	numberOfContracts int,
	testRange *big.Int,
	interval *big.Int,
) []contracts.UpkeepCounter {
	upkeepCounters := make([]contracts.UpkeepCounter, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		upkeepCounter, err := contractDeployer.DeployUpkeepCounter(testRange, interval)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		upkeepCounters = append(upkeepCounters, upkeepCounter)
		log.Debug().
			Str("Contract Address", upkeepCounter.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumer deployments")
		}
	}
	err := networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeepCounters
}

func DeployUpkeepPerformCounterRestrictive(
	contractDeployer contracts.ContractDeployer,
	networks *blockchain.Networks,
	numberOfContracts int,
	testRange *big.Int,
	averageEligibilityCadence *big.Int,
) []contracts.UpkeepPerformCounterRestrictive {
	upkeepCounters := make([]contracts.UpkeepPerformCounterRestrictive, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		upkeepCounter, err := contractDeployer.DeployUpkeepPerformCounterRestrictive(testRange, averageEligibilityCadence)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		upkeepCounters = append(upkeepCounters, upkeepCounter)
		log.Debug().
			Str("Contract Address", upkeepCounter.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumer deployments")
		}
	}
	err := networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeepCounters
}
