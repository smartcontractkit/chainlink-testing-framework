package contracts

import (
	"math/big"
)

// ContractDeployer is an interface for abstracting the contract deployment methods across network implementations
type ContractDeployer interface {
	DeployAPIConsumer(linkAddr string) (APIConsumer, error)
	DeployOracle(linkAddr string) (Oracle, error)
	DeployFlags(rac string) (Flags, error)
	DeployDeviationFlaggingValidator(
		flags string,
		flaggingThreshold *big.Int,
	) (DeviationFlaggingValidator, error)
	DeployFluxAggregatorContract(linkAddr string, fluxOptions FluxAggregatorOptions) (FluxAggregator, error)
	DeployLinkTokenContract() (LinkToken, error)
	DeployOffChainAggregator(linkAddr string, offchainOptions OffchainOptions) (OffchainAggregator, error)
	DeployVRFContract() (VRF, error)
	DeployMockETHLINKFeed(answer *big.Int) (MockETHLINKFeed, error)
	DeployMockGasFeed(answer *big.Int) (MockGasFeed, error)
	DeployUpkeepRegistrationRequests(linkAddr string, minLinkJuels *big.Int) (UpkeepRegistrar, error)
	DeployKeeperRegistry(opts *KeeperRegistryOpts) (KeeperRegistry, error)
	DeployKeeperConsumer(updateInterval *big.Int) (KeeperConsumer, error)
	DeployKeeperConsumerPerformance(
		testBlockRange,
		averageCadence,
		checkGasToBurn,
		performGasToBurn *big.Int,
	) (KeeperConsumerPerformance, error)
	DeployUpkeepCounter(testRange *big.Int, interval *big.Int) (UpkeepCounter, error)
	DeployUpkeepPerformCounterRestrictive(testRange *big.Int, averageEligibilityCadence *big.Int) (UpkeepPerformCounterRestrictive, error)
	DeployVRFConsumer(linkAddr string, coordinatorAddr string) (VRFConsumer, error)
	DeployVRFConsumerV2(linkAddr string, coordinatorAddr string) (VRFConsumerV2, error)
	DeployVRFCoordinator(linkAddr string, bhsAddr string) (VRFCoordinator, error)
	DeployVRFCoordinatorV2(linkAddr string, bhsAddr string, linkEthFeedAddr string) (VRFCoordinatorV2, error)
	DeployBlockhashStore() (BlockHashStore, error)
}
