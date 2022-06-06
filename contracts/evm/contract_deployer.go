package evm

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
)

// ContractDeployer provides the implementations for deploying ETH (EVM) based contracts
type ContractDeployer struct {
	client blockchain.EVMClient
}

// DeployFlags deploys flags contract
func (e *ContractDeployer) DeployFlags(
	rac string,
) (Flags, error) {
	address, _, instance, err := e.client.DeployContract("Flags", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		racAddr := common.HexToAddress(rac)
		return DeployFlags(auth, backend, racAddr)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFlags{
		client:  e.client,
		flags:   instance.(*Flags),
		address: address,
	}, nil
}

// DeployDeviationFlaggingValidator deploys deviation flagging validator contract
func (e *ContractDeployer) DeployDeviationFlaggingValidator(
	flags string,
	flaggingThreshold *big.Int,
) (DeviationFlaggingValidator, error) {
	address, _, instance, err := e.client.DeployContract("Deviation flagging validator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		flagAddr := common.HexToAddress(flags)
		return DeployDeviationFlaggingValidator(auth, backend, flagAddr, flaggingThreshold)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumDeviationFlaggingValidator{
		client:  e.client,
		dfv:     instance.(*DeviationFlaggingValidator),
		address: address,
	}, nil
}

// DeployFluxAggregatorContract deploys the Flux Aggregator Contract on an EVM chain
func (e *ContractDeployer) DeployFluxAggregatorContract(
	linkAddr string,
	fluxOptions FluxAggregatorOptions,
) (FluxAggregator, error) {
	address, _, instance, err := e.client.DeployContract("Flux Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		la := common.HexToAddress(linkAddr)
		return DeployFluxAggregator(auth,
			backend,
			la,
			fluxOptions.PaymentAmount,
			fluxOptions.Timeout,
			fluxOptions.Validator,
			fluxOptions.MinSubValue,
			fluxOptions.MaxSubValue,
			fluxOptions.Decimals,
			fluxOptions.Description)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFluxAggregator{
		client:         e.client,
		fluxAggregator: instance.(*FluxAggregator),
		address:        address,
	}, nil
}

// DeployLinkTokenContract deploys a Link Token contract to an EVM chain
func (e *ContractDeployer) DeployLinkTokenContract() (LinkToken, error) {
	linkTokenAddress, _, instance, err := e.client.DeployContract("LINK Token", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployLinkToken(auth, backend)
	})
	if err != nil {
		return nil, err
	}

	return &EthereumLinkToken{
		client:   e.client,
		instance: instance.(*LinkToken),
		address:  *linkTokenAddress,
	}, err
}

// DefaultOffChainAggregatorOptions returns some base defaults for deploying an OCR contract
func DefaultOffChainAggregatorOptions() OffchainOptions {
	return OffchainOptions{
		MaximumGasPrice:         uint32(3000),
		ReasonableGasPrice:      uint32(10),
		MicroLinkPerEth:         uint32(500),
		LinkGweiPerObservation:  uint32(500),
		LinkGweiPerTransmission: uint32(500),
		MinimumAnswer:           big.NewInt(1),
		MaximumAnswer:           big.NewInt(50000000000000000),
		Decimals:                8,
		Description:             "Test OCR",
	}
}

// DefaultOffChainAggregatorConfig returns some base defaults for configuring an OCR contract
func DefaultOffChainAggregatorConfig(numberNodes int) OffChainAggregatorConfig {
	if numberNodes <= 4 {
		log.Err(fmt.Errorf("Insufficient number of nodes (%d) supplied for OCR, need at least 5", numberNodes)).
			Int("Number Chainlink Nodes", numberNodes).
			Msg("You likely need more chainlink nodes to properly configure OCR, try 5 or more.")
	}
	s := []int{1}
	// First node's stage already inputted as a 1 in line above, so numberNodes-1.
	for i := 0; i < numberNodes-1; i++ {
		s = append(s, 2)
	}
	return OffChainAggregatorConfig{
		AlphaPPB:         1,
		DeltaC:           time.Minute * 60,
		DeltaGrace:       time.Second * 12,
		DeltaProgress:    time.Second * 35,
		DeltaStage:       time.Second * 60,
		DeltaResend:      time.Second * 17,
		DeltaRound:       time.Second * 30,
		RMax:             6,
		S:                s,
		N:                numberNodes,
		F:                1,
		OracleIdentities: []ocrConfigHelper.OracleIdentityExtra{},
	}
}

// DeployOffChainAggregator deploys the offchain aggregation contract to the EVM chain
func (e *ContractDeployer) DeployOffChainAggregator(
	linkAddr string,
	offchainOptions OffchainOptions,
) (OffchainAggregator, error) {
	address, _, instance, err := e.client.DeployContract("OffChain Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		la := common.HexToAddress(linkAddr)
		return DeployOffchainAggregator(auth,
			backend,
			offchainOptions.MaximumGasPrice,
			offchainOptions.ReasonableGasPrice,
			offchainOptions.MicroLinkPerEth,
			offchainOptions.LinkGweiPerObservation,
			offchainOptions.LinkGweiPerTransmission,
			la,
			offchainOptions.MinimumAnswer,
			offchainOptions.MaximumAnswer,
			offchainOptions.BillingAccessController,
			offchainOptions.RequesterAccessController,
			offchainOptions.Decimals,
			offchainOptions.Description)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOffchainAggregator{
		client:  e.client,
		ocr:     instance.(*OffchainAggregator),
		address: address,
	}, err
}

// DeployAPIConsumer deploys api consumer for oracle
func (e *ContractDeployer) DeployAPIConsumer(linkAddr string) (APIConsumer, error) {
	addr, _, instance, err := e.client.DeployContract("APIConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployAPIConsumer(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumAPIConsumer{
		address:  addr,
		client:   e.client,
		consumer: instance.(*APIConsumer),
	}, err
}

// DeployOracle deploys oracle for consumer test
func (e *ContractDeployer) DeployOracle(linkAddr string) (Oracle, error) {
	addr, _, instance, err := e.client.DeployContract("Oracle", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployOracle(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOracle{
		address: addr,
		client:  e.client,
		oracle:  instance.(*Oracle),
	}, err
}

// DeployVRFContract deploy VRF contract
func (e *ContractDeployer) DeployVRFContract() (VRF, error) {
	address, _, instance, err := e.client.DeployContract("VRF", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployVRF(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRF{
		client:  e.client,
		vrf:     instance.(*VRF),
		address: address,
	}, err
}

func (e *ContractDeployer) DeployMockETHLINKFeed(answer *big.Int) (MockETHLINKFeed, error) {
	address, _, instance, err := e.client.DeployContract("MockETHLINKAggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployMockV3AggregatorContract(auth, backend, 18, answer)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockETHLINKFeed{
		client:  e.client,
		feed:    instance.(*MockV3AggregatorContract),
		address: address,
	}, err
}

func (e *ContractDeployer) DeployMockGasFeed(answer *big.Int) (MockGasFeed, error) {
	address, _, instance, err := e.client.DeployContract("MockGasFeed", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployMockGASAggregator(auth, backend, answer)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockGASFeed{
		client:  e.client,
		feed:    instance.(*MockGASAggregator),
		address: address,
	}, err
}

func (e *ContractDeployer) DeployUpkeepRegistrationRequests(linkAddr string, minLinkJuels *big.Int) (UpkeepRegistrar, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepRegistrationRequests", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployUpkeepRegistrationRequests(auth, backend, common.HexToAddress(linkAddr), minLinkJuels)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumUpkeepRegistrationRequests{
		client:    e.client,
		registrar: instance.(*UpkeepRegistrationRequests),
		address:   address,
	}, err
}

func (e *ContractDeployer) DeployKeeperRegistry(
	opts *KeeperRegistryOpts,
) (KeeperRegistry, error) {
	switch opts.RegistryVersion {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		address, _, instance, err := e.client.DeployContract("KeeperRegistry1_1", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return DeployKeeperRegistry11(
				auth,
				backend,
				common.HexToAddress(opts.LinkAddr),
				common.HexToAddress(opts.ETHFeedAddr),
				common.HexToAddress(opts.GasFeedAddr),
				opts.Settings.PaymentPremiumPPB,
				opts.Settings.FlatFeeMicroLINK,
				opts.Settings.BlockCountPerTurn,
				opts.Settings.CheckGasLimit,
				opts.Settings.StalenessSeconds,
				opts.Settings.GasCeilingMultiplier,
				opts.Settings.FallbackGasPrice,
				opts.Settings.FallbackLinkPrice,
			)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			client:      e.client,
			version:     RegistryVersion_1_1,
			registry1_1: instance.(*KeeperRegistry11),
			registry1_2: nil,
			address:     address,
		}, err
	case RegistryVersion_1_2:
		address, _, instance, err := e.client.DeployContract("KeeperRegistry", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return DeployKeeperRegistry(
				auth,
				backend,
				common.HexToAddress(opts.LinkAddr),
				common.HexToAddress(opts.ETHFeedAddr),
				common.HexToAddress(opts.GasFeedAddr),
				Config{
					PaymentPremiumPPB:    opts.Settings.PaymentPremiumPPB,
					FlatFeeMicroLink:     opts.Settings.FlatFeeMicroLINK,
					BlockCountPerTurn:    opts.Settings.BlockCountPerTurn,
					CheckGasLimit:        opts.Settings.CheckGasLimit,
					StalenessSeconds:     opts.Settings.StalenessSeconds,
					GasCeilingMultiplier: opts.Settings.GasCeilingMultiplier,
					MinUpkeepSpend:       opts.Settings.MinUpkeepSpend,
					MaxPerformGas:        opts.Settings.MaxPerformGas,
					FallbackGasPrice:     opts.Settings.FallbackGasPrice,
					FallbackLinkPrice:    opts.Settings.FallbackLinkPrice,
					Transcoder:           common.HexToAddress(opts.TranscoderAddr),
					Registrar:            common.HexToAddress(opts.RegistrarAddr),
				},
			)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			client:      e.client,
			version:     RegistryVersion_1_2,
			registry1_1: nil,
			registry1_2: instance.(*KeeperRegistry),
			address:     address,
		}, err

	default:
		return nil, fmt.Errorf("keeper registry version %d is not supported", opts.RegistryVersion)
	}
}

func (e *ContractDeployer) DeployKeeperConsumer(updateInterval *big.Int) (KeeperConsumer, error) {
	address, _, instance, err := e.client.DeployContract("KeeperConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployKeeperConsumer(auth, backend, updateInterval)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperConsumer{
		client:   e.client,
		consumer: instance.(*KeeperConsumer),
		address:  address,
	}, err
}

func (e *ContractDeployer) DeployUpkeepCounter(testRange *big.Int, interval *big.Int) (UpkeepCounter, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepCounter", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployUpkeepCounter(auth, backend, testRange, interval)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumUpkeepCounter{
		client:   e.client,
		consumer: instance.(*UpkeepCounter),
		address:  address,
	}, err
}

func (e *ContractDeployer) DeployUpkeepPerformCounterRestrictive(testRange *big.Int, averageEligibilityCadence *big.Int) (UpkeepPerformCounterRestrictive, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepPerformCounterRestrictive", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployUpkeepPerformCounterRestrictive(auth, backend, testRange, averageEligibilityCadence)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumUpkeepPerformCounterRestrictive{
		client:   e.client,
		consumer: instance.(*UpkeepPerformCounterRestrictive),
		address:  address,
	}, err
}

func (e *ContractDeployer) DeployKeeperConsumerPerformance(
	testBlockRange,
	averageCadence,
	checkGasToBurn,
	performGasToBurn *big.Int,
) (KeeperConsumerPerformance, error) {
	address, _, instance, err := e.client.DeployContract("KeeperConsumerPerformance", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployKeeperConsumerPerformance(
			auth,
			backend,
			testBlockRange,
			averageCadence,
			checkGasToBurn,
			performGasToBurn,
		)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperConsumerPerformance{
		client:   e.client,
		consumer: instance.(*KeeperConsumerPerformance),
		address:  address,
	}, err
}

// DeployBlockhashStore deploys blockhash store used with VRF contract
func (e *ContractDeployer) DeployBlockhashStore() (BlockHashStore, error) {
	address, _, instance, err := e.client.DeployContract("BlockhashStore", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployBlockhashStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumBlockhashStore{
		client:         e.client,
		blockHashStore: instance.(*BlockhashStore),
		address:        address,
	}, err
}

// DeployVRFCoordinatorV2 deploys VRFV2 coordinator contract
func (e *ContractDeployer) DeployVRFCoordinatorV2(linkAddr string, bhsAddr string, linkEthFeedAddr string) (VRFCoordinatorV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployVRFCoordinatorV2(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr), common.HexToAddress(linkEthFeedAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV2{
		client:      e.client,
		coordinator: instance.(*VRFCoordinatorV2),
		address:     address,
	}, err
}

// DeployVRFCoordinator deploys VRF coordinator contract
func (e *ContractDeployer) DeployVRFCoordinator(linkAddr string, bhsAddr string) (VRFCoordinator, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployVRFCoordinator(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinator{
		client:      e.client,
		coordinator: instance.(*VRFCoordinator),
		address:     address,
	}, err
}

// DeployVRFConsumer deploys VRF consumer contract
func (e *ContractDeployer) DeployVRFConsumer(linkAddr string, coordinatorAddr string) (VRFConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployVRFConsumer(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumer{
		client:   e.client,
		consumer: instance.(*VRFConsumer),
		address:  address,
	}, err
}

// DeployVRFConsumerV2 deploys VRFv@ consumer contract
func (e *ContractDeployer) DeployVRFConsumerV2(linkAddr string, coordinatorAddr string) (VRFConsumerV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFConsumerV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return DeployVRFConsumerV2(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumerV2{
		client:   e.client,
		consumer: instance.(*VRFConsumerV2),
		address:  address,
	}, err
}
