package contracts

import (
	"context"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts/celo"

	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	ocrConfigHelper "github.com/smartcontractkit/integrations-framework/libocr/offchainreporting/confighelper"
)

// ContractDeployer is an interface for abstracting the contract deployment methods across network implementations
type ContractDeployer interface {
	Balance(fromWallet client.BlockchainWallet) (*big.Float, error)
	CalculateETHForChainlinkOperations(numberOfOperations int) (*big.Float, error)
	DeployStorageContract(fromWallet client.BlockchainWallet) (Storage, error)
	DeployAPIConsumer(fromWallet client.BlockchainWallet, linkAddr string) (APIConsumer, error)
	DeployOracle(fromWallet client.BlockchainWallet, linkAddr string) (Oracle, error)
	DeployReadAccessController(fromWallet client.BlockchainWallet) (ReadAccessController, error)
	DeployFlags(fromWallet client.BlockchainWallet, rac string) (Flags, error)
	DeployDeviationFlaggingValidator(
		fromWallet client.BlockchainWallet,
		flags string,
		flaggingThreshold *big.Int,
	) (DeviationFlaggingValidator, error)
	DeployFluxAggregatorContract(
		fromWallet client.BlockchainWallet,
		fluxOptions FluxAggregatorOptions,
	) (FluxAggregator, error)
	DeployLinkTokenContract(fromWallet client.BlockchainWallet) (LinkToken, error)
	DeployOCRv2(
		fromWallet client.BlockchainWallet,
		paymentControllerAddr string,
		requesterControllerAddr string,
		linkTokenAddr string,
	) (OCRv2, error)
	DeployOCRv2AccessController(fromWallet client.BlockchainWallet) (OCRv2AccessController, error)
	DeployOffChainAggregator(
		fromWallet client.BlockchainWallet,
		offchainOptions OffchainOptions,
	) (OffchainAggregator, error)
	DeployVRFContract(fromWallet client.BlockchainWallet) (VRF, error)
	DeployMockETHLINKFeed(fromWallet client.BlockchainWallet, answer *big.Int) (MockETHLINKFeed, error)
	DeployMockGasFeed(fromWallet client.BlockchainWallet, answer *big.Int) (MockGasFeed, error)
	DeployUpkeepRegistrationRequests(fromWallet client.BlockchainWallet, linkAddr string, minLinkJuels *big.Int) (UpkeepRegistrar, error)
	DeployKeeperRegistry(
		fromWallet client.BlockchainWallet,
		opts *KeeperRegistryOpts,
	) (KeeperRegistry, error)
	DeployKeeperConsumer(fromWallet client.BlockchainWallet, updateInterval *big.Int) (KeeperConsumer, error)
	DeployVRFConsumer(fromWallet client.BlockchainWallet, linkAddr string, coordinatorAddr string) (VRFConsumer, error)
	DeployVRFCoordinator(fromWallet client.BlockchainWallet, linkAddr string, bhsAddr string) (VRFCoordinator, error)
	DeployBlockhashStore(fromWallet client.BlockchainWallet) (BlockHashStore, error)
}

// NewContractDeployer returns an instance of a contract deployer based on the client type
func NewContractDeployer(bcClient client.BlockchainClient) (ContractDeployer, error) {
	switch clientImpl := bcClient.Get().(type) {
	case *client.CeloClient:
		return NewEthereumContractDeployer(clientImpl), nil
	case *client.CeloClients:
		return NewEthereumContractDeployer(clientImpl.DefaultClient), nil
	}
	return nil, errors.New("unknown blockchain client implementation")
}

// EthereumContractDeployer provides the implementations for deploying ETH (EVM) based contracts
type EthereumContractDeployer struct {
	eth *client.CeloClient
}

func (e *EthereumContractDeployer) DeployOCRv2(
	fromWallet client.BlockchainWallet,
	paymentControllerAddr string,
	requesterControllerAddr string,
	linkTokenAddr string,
) (OCRv2, error) {
	panic("implement me")
}

func (e *EthereumContractDeployer) DeployOCRv2AccessController(fromWallet client.BlockchainWallet) (OCRv2AccessController, error) {
	panic("implement me")
}

// NewEthereumContractDeployer returns an instantiated instance of the ETH contract deployer
func NewEthereumContractDeployer(ethClient *client.CeloClient) *EthereumContractDeployer {
	return &EthereumContractDeployer{
		eth: ethClient,
	}
}

// DefaultFluxAggregatorOptions produces some basic defaults for a flux aggregator contract
func DefaultFluxAggregatorOptions() FluxAggregatorOptions {
	return FluxAggregatorOptions{
		PaymentAmount: big.NewInt(1),
		Timeout:       uint32(30),
		MinSubValue:   big.NewInt(0),
		MaxSubValue:   big.NewInt(1000000000000),
		Decimals:      uint8(0),
		Description:   "Test Flux Aggregator",
	}
}

// DeployReadAccessController deploys read/write access controller contract
func (e *EthereumContractDeployer) DeployReadAccessController(
	fromWallet client.BlockchainWallet,
) (ReadAccessController, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "Read Access Controller", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeploySimpleReadAccessController(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumReadAccessController{
		client:       e.eth,
		rac:          instance.(*celo.SimpleReadAccessController),
		callerWallet: fromWallet,
		address:      address,
	}, nil
}

// DeployFlags deploys flags contract
func (e *EthereumContractDeployer) DeployFlags(
	fromWallet client.BlockchainWallet,
	rac string,
) (Flags, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "Flags", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		racAddr := common.HexToAddress(rac)
		return celo.DeployFlags(auth, backend, racAddr)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFlags{
		client:       e.eth,
		flags:        instance.(*celo.Flags),
		callerWallet: fromWallet,
		address:      address,
	}, nil
}

// DeployDeviationFlaggingValidator deploys deviation flagging validator contract
func (e *EthereumContractDeployer) DeployDeviationFlaggingValidator(
	fromWallet client.BlockchainWallet,
	flags string,
	flaggingThreshold *big.Int,
) (DeviationFlaggingValidator, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "Deviation flagging validator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		flagAddr := common.HexToAddress(flags)
		return celo.DeployDeviationFlaggingValidator(auth, backend, flagAddr, flaggingThreshold)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumDeviationFlaggingValidator{
		client:       e.eth,
		dfv:          instance.(*celo.DeviationFlaggingValidator),
		callerWallet: fromWallet,
		address:      address,
	}, nil
}

// DeployFluxAggregatorContract deploys the Flux Aggregator Contract on an EVM chain
func (e *EthereumContractDeployer) DeployFluxAggregatorContract(
	fromWallet client.BlockchainWallet,
	fluxOptions FluxAggregatorOptions,
) (FluxAggregator, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "Flux Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		linkAddress := common.HexToAddress(e.eth.Network.Config().LinkTokenAddress)
		return celo.DeployFluxAggregator(auth,
			backend,
			linkAddress,
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
		client:         e.eth,
		fluxAggregator: instance.(*celo.FluxAggregator),
		callerWallet:   fromWallet,
		address:        address,
	}, nil
}

// DeployLinkTokenContract deploys a Link Token contract to an EVM chain
func (e *EthereumContractDeployer) DeployLinkTokenContract(fromWallet client.BlockchainWallet) (LinkToken, error) {
	linkTokenAddress, _, instance, err := e.eth.DeployContract(fromWallet, "LINK Token", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployLinkToken(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	// Set config address
	e.eth.Network.Config().LinkTokenAddress = linkTokenAddress.Hex()

	return &EthereumLinkToken{
		client:       e.eth,
		linkToken:    instance.(*celo.LinkToken),
		callerWallet: fromWallet,
		address:      *linkTokenAddress,
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

// DefaultOffChainAggregatorConfig returns some base defaults for configuring an OCR contract on a very fast simulated
// geth network
func DefaultOffChainAggregatorConfig(numberNodes int) OffChainAggregatorConfig {
	s := []int{}
	for i := 0; i < numberNodes; i++ {
		s = append(s, 1)
	}
	if numberNodes <= 3 {
		log.Warn().
			Int("Number Chainlink Nodes", numberNodes).
			Msg("You likely need more chainlink nodes to properly configure OCR, try 5 or more.")
	}
	return OffChainAggregatorConfig{
		AlphaPPB:         1,
		DeltaC:           time.Minute * 10,
		DeltaGrace:       time.Second,
		DeltaProgress:    time.Second * 30,
		DeltaStage:       time.Second * 10,
		DeltaResend:      time.Second * 10,
		DeltaRound:       time.Second * 20,
		RMax:             4,
		S:                s,
		N:                numberNodes,
		F:                int(math.Max(1, float64(numberNodes/3-1))),
		OracleIdentities: []ocrConfigHelper.OracleIdentityExtra{},
	}
}

// OptimismOffChainAggregatorConfig returns some base defaults for configuring an OCR contract on Optimism chain
// I suspect most of the tuning left to try is in here.
func OptimismOffChainAggregatorConfig(numberNodes int) OffChainAggregatorConfig {
	s := []int{1}
	// First node's stage already inputted as a 1 in line above, so numberNodes-1.
	for i := 0; i < numberNodes-1; i++ {
		s = append(s, 2)
	}
	if numberNodes < 4 {
		log.Warn().
			Int("Number Chainlink Nodes", numberNodes).
			Msg("You likely need more chainlink nodes to properly configure OCR, try 5 or more total (one for bootstrap).")
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
func (e *EthereumContractDeployer) DeployOffChainAggregator(
	fromWallet client.BlockchainWallet,
	offchainOptions OffchainOptions,
) (OffchainAggregator, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "OffChain Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		linkAddress := common.HexToAddress(e.eth.Network.Config().LinkTokenAddress)
		return celo.DeployOffchainAggregator(auth,
			backend,
			offchainOptions.MaximumGasPrice,
			offchainOptions.ReasonableGasPrice,
			offchainOptions.MicroLinkPerEth,
			offchainOptions.LinkGweiPerObservation,
			offchainOptions.LinkGweiPerTransmission,
			linkAddress,
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
		client:       e.eth,
		ocr:          instance.(*celo.OffchainAggregator),
		callerWallet: fromWallet,
		address:      address,
	}, err
}

// Balance get deployer wallet balance
func (e *EthereumContractDeployer) Balance(fromWallet client.BlockchainWallet) (*big.Float, error) {
	balance, err := e.eth.Client.PendingBalanceAt(context.Background(), common.HexToAddress(fromWallet.Address()))
	if err != nil {
		return nil, err
	}
	bf := new(big.Float).SetInt(balance)
	return big.NewFloat(1).Quo(bf, client.OneEth), nil
}

// CalculateETHForChainlinkOperations calculates required amount of ETH for amountOfOperations Chainlink operations
// based on the network's suggested gas price and the chainlink gas limit. This is fairly imperfect and should be used
// as only a rough, upper-end estimate instead of an exact calculation.
// See https://ethereum.org/en/developers/docs/gas/#post-london for info on how gas calculation works
func (e *EthereumContractDeployer) CalculateETHForChainlinkOperations(amountOfOperations int) (*big.Float, error) {
	bigAmountOfOperations := big.NewInt(int64(amountOfOperations))
	gasPriceInWei, err := e.eth.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// https://ethereum.stackexchange.com/questions/19665/how-to-calculate-transaction-fee
	// total gas limit = chainlink gas limit + gas limit buffer
	gasLimit := e.eth.Network.Config().GasEstimationBuffer + e.eth.Network.Config().ChainlinkGasLimit
	// gas cost for TX = total gas limit * estimated gas price
	gasCostPerOperationWei := big.NewInt(1).Mul(big.NewInt(1).SetUint64(gasLimit), gasPriceInWei)
	gasCostPerOperationWeiFloat := big.NewFloat(1).SetInt(gasCostPerOperationWei)
	gasCostPerOperationETH := big.NewFloat(1).Quo(gasCostPerOperationWeiFloat, client.OneEth)
	// total Wei needed for all TXs = total value for TX * number of TXs
	totalWeiForAllOperations := big.NewInt(1).Mul(gasCostPerOperationWei, bigAmountOfOperations)
	totalWeiForAllOperationsFloat := big.NewFloat(1).SetInt(totalWeiForAllOperations)
	totalEthForAllOperations := big.NewFloat(1).Quo(totalWeiForAllOperationsFloat, client.OneEth)

	log.Debug().
		Int("Number of Operations", amountOfOperations).
		Uint64("Gas Limit per Operation", gasLimit).
		Str("Value per Operation (ETH)", gasCostPerOperationETH.String()).
		Str("Total (ETH)", totalEthForAllOperations.String()).
		Msg("Calculated ETH for Chainlink Operations")

	return totalEthForAllOperations, nil
}

// DeployStorageContract deploys a vanilla storage contract that is a value store
func (e *EthereumContractDeployer) DeployStorageContract(fromWallet client.BlockchainWallet) (Storage, error) {
	_, _, instance, err := e.eth.DeployContract(fromWallet, "Storage", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumStorage{
		client:       e.eth,
		store:        instance.(*celo.Store),
		callerWallet: fromWallet,
	}, err
}

// DeployAPIConsumer deploys api consumer for oracle
func (e *EthereumContractDeployer) DeployAPIConsumer(fromWallet client.BlockchainWallet, linkAddr string) (APIConsumer, error) {
	addr, _, instance, err := e.eth.DeployContract(fromWallet, "APIConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployAPIConsumer(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumAPIConsumer{
		address:      addr,
		client:       e.eth,
		consumer:     instance.(*celo.APIConsumer),
		callerWallet: fromWallet,
	}, err
}

// DeployOracle deploys oracle for consumer test
func (e *EthereumContractDeployer) DeployOracle(fromWallet client.BlockchainWallet, linkAddr string) (Oracle, error) {
	addr, _, instance, err := e.eth.DeployContract(fromWallet, "Oracle", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployOracle(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOracle{
		address:      addr,
		client:       e.eth,
		oracle:       instance.(*celo.Oracle),
		callerWallet: fromWallet,
	}, err
}

// DeployVRFContract deploy VRF contract
func (e *EthereumContractDeployer) DeployVRFContract(fromWallet client.BlockchainWallet) (VRF, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "VRF", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployVRF(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRF{
		client:       e.eth,
		vrf:          instance.(*celo.VRF),
		callerWallet: fromWallet,
		address:      address,
	}, err
}

func (e *EthereumContractDeployer) DeployMockETHLINKFeed(fromWallet client.BlockchainWallet, answer *big.Int) (MockETHLINKFeed, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "MockETHLINKFeed", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployMockETHLINKAggregator(auth, backend, answer)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockETHLINKFeed{
		client:       e.eth,
		feed:         instance.(*celo.MockETHLINKAggregator),
		callerWallet: fromWallet,
		address:      address,
	}, err
}

func (e *EthereumContractDeployer) DeployMockGasFeed(fromWallet client.BlockchainWallet, answer *big.Int) (MockGasFeed, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "MockGasFeed", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployMockGASAggregator(auth, backend, answer)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockGASFeed{
		client:       e.eth,
		feed:         instance.(*celo.MockGASAggregator),
		callerWallet: fromWallet,
		address:      address,
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepRegistrationRequests(fromWallet client.BlockchainWallet, linkAddr string, minLinkJuels *big.Int) (UpkeepRegistrar, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "UpkeepRegistrationRequests", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployUpkeepRegistrationRequests(auth, backend, common.HexToAddress(linkAddr), minLinkJuels)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumUpkeepRegistrationRequests{
		client:       e.eth,
		registrar:    instance.(*celo.UpkeepRegistrationRequests),
		callerWallet: fromWallet,
		address:      address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperRegistry(
	fromWallet client.BlockchainWallet,
	opts *KeeperRegistryOpts,
) (KeeperRegistry, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "KeeperRegistry", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployKeeperRegistry(
			auth,
			backend,
			common.HexToAddress(opts.LinkAddr),
			common.HexToAddress(opts.ETHFeedAddr),
			common.HexToAddress(opts.GasFeedAddr),
			opts.PaymentPremiumPPB,
			opts.BlockCountPerTurn,
			opts.CheckGasLimit,
			opts.StalenessSeconds,
			opts.GasCeilingMultiplier,
			opts.FallbackGasPrice,
			opts.FallbackLinkPrice,
		)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperRegistry{
		client:       e.eth,
		registry:     instance.(*celo.KeeperRegistry),
		callerWallet: fromWallet,
		address:      address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperConsumer(fromWallet client.BlockchainWallet, updateInterval *big.Int) (KeeperConsumer, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "KeeperConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployKeeperConsumer(auth, backend, updateInterval)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperConsumer{
		client:       e.eth,
		consumer:     instance.(*celo.KeeperConsumer),
		callerWallet: fromWallet,
		address:      address,
	}, err
}

// DeployBlockhashStore deploys blockhash store used with VRF contract
func (e *EthereumContractDeployer) DeployBlockhashStore(fromWallet client.BlockchainWallet) (BlockHashStore, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "BlockhashStore", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployBlockhashStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumBlockhashStore{
		client:         e.eth,
		blockHashStore: instance.(*celo.BlockhashStore),
		callerWallet:   fromWallet,
		address:        address,
	}, err
}

// DeployVRFCoordinator deploys VRF coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinator(fromWallet client.BlockchainWallet, linkAddr string, bhsAddr string) (VRFCoordinator, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "VRFCoordinator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployVRFCoordinator(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinator{
		client:       e.eth,
		coordinator:  instance.(*celo.VRFCoordinator),
		callerWallet: fromWallet,
		address:      address,
	}, err
}

// DeployVRFConsumer deploys VRF consumer contract
func (e *EthereumContractDeployer) DeployVRFConsumer(fromWallet client.BlockchainWallet, linkAddr string, coordinatorAddr string) (VRFConsumer, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "VRFConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return celo.DeployVRFConsumer(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumer{
		client:       e.eth,
		consumer:     instance.(*celo.VRFConsumer),
		callerWallet: fromWallet,
		address:      address,
	}, err
}
