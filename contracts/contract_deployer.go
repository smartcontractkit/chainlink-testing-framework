package contracts

import (
	"context"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts/ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
)

// ContractDeployer is an interface for abstracting the contract deployment methods across network implementations
type ContractDeployer interface {
	Balance() (*big.Float, error)
	DeployStorageContract() (Storage, error)
	DeployAPIConsumer(linkAddr string) (APIConsumer, error)
	DeployOracle(linkAddr string) (Oracle, error)
	DeployReadAccessController() (ReadAccessController, error)
	DeployFlags(rac string) (Flags, error)
	DeployDeviationFlaggingValidator(
		flags string,
		flaggingThreshold *big.Int,
	) (DeviationFlaggingValidator, error)
	DeployFluxAggregatorContract(linkAddr string, fluxOptions FluxAggregatorOptions) (FluxAggregator, error)
	DeployLinkTokenContract() (LinkToken, error)
	DeployOCRv2(
		paymentControllerAddr string,
		requesterControllerAddr string,
		linkTokenAddr string,
	) (OCRv2, error)
	DeployOCRv2AccessController() (OCRv2AccessController, error)
	DeployOffChainAggregator(linkAddr string, offchainOptions OffchainOptions) (OffchainAggregator, error)
	DeployVRFContract() (VRF, error)
	DeployMockETHLINKFeed(answer *big.Int) (MockETHLINKFeed, error)
	DeployMockGasFeed(answer *big.Int) (MockGasFeed, error)
	DeployUpkeepRegistrationRequests(linkAddr string, minLinkJuels *big.Int) (UpkeepRegistrar, error)
	DeployKeeperRegistry(opts *KeeperRegistryOpts) (KeeperRegistry, error)
	DeployKeeperConsumer(updateInterval *big.Int) (KeeperConsumer, error)
	DeployVRFConsumer(linkAddr string, coordinatorAddr string) (VRFConsumer, error)
	DeployVRFCoordinator(linkAddr string, bhsAddr string) (VRFCoordinator, error)
	DeployBlockhashStore() (BlockHashStore, error)
}

// NewContractDeployer returns an instance of a contract deployer based on the client type
func NewContractDeployer(bcClient client.BlockchainClient) (ContractDeployer, error) {
	switch clientImpl := bcClient.Get().(type) {
	case *client.EthereumClient:
		return NewEthereumContractDeployer(clientImpl), nil
	}
	return nil, errors.New("unknown blockchain client implementation")
}

// EthereumContractDeployer provides the implementations for deploying ETH (EVM) based contracts
type EthereumContractDeployer struct {
	eth *client.EthereumClient
}

func (e *EthereumContractDeployer) DeployOCRv2(
	paymentControllerAddr string,
	requesterControllerAddr string,
	linkTokenAddr string,
) (OCRv2, error) {
	panic("implement me")
}

func (e *EthereumContractDeployer) DeployOCRv2AccessController() (OCRv2AccessController, error) {
	panic("implement me")
}

// NewEthereumContractDeployer returns an instantiated instance of the ETH contract deployer
func NewEthereumContractDeployer(ethClient *client.EthereumClient) *EthereumContractDeployer {
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
func (e *EthereumContractDeployer) DeployReadAccessController() (ReadAccessController, error) {
	address, _, instance, err := e.eth.DeployContract("Read Access Controller", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeploySimpleReadAccessController(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumReadAccessController{
		client:  e.eth,
		rac:     instance.(*ethereum.SimpleReadAccessController),
		address: address,
	}, nil
}

// DeployFlags deploys flags contract
func (e *EthereumContractDeployer) DeployFlags(
	rac string,
) (Flags, error) {
	address, _, instance, err := e.eth.DeployContract("Flags", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		racAddr := common.HexToAddress(rac)
		return ethereum.DeployFlags(auth, backend, racAddr)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFlags{
		client:  e.eth,
		flags:   instance.(*ethereum.Flags),
		address: address,
	}, nil
}

// DeployDeviationFlaggingValidator deploys deviation flagging validator contract
func (e *EthereumContractDeployer) DeployDeviationFlaggingValidator(
	flags string,
	flaggingThreshold *big.Int,
) (DeviationFlaggingValidator, error) {
	address, _, instance, err := e.eth.DeployContract("Deviation flagging validator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		flagAddr := common.HexToAddress(flags)
		return ethereum.DeployDeviationFlaggingValidator(auth, backend, flagAddr, flaggingThreshold)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumDeviationFlaggingValidator{
		client:  e.eth,
		dfv:     instance.(*ethereum.DeviationFlaggingValidator),
		address: address,
	}, nil
}

// DeployFluxAggregatorContract deploys the Flux Aggregator Contract on an EVM chain
func (e *EthereumContractDeployer) DeployFluxAggregatorContract(
	linkAddr string,
	fluxOptions FluxAggregatorOptions,
) (FluxAggregator, error) {
	address, _, instance, err := e.eth.DeployContract("Flux Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		la := common.HexToAddress(linkAddr)
		return ethereum.DeployFluxAggregator(auth,
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
		client:         e.eth,
		fluxAggregator: instance.(*ethereum.FluxAggregator),
		address:        address,
	}, nil
}

// DeployLinkTokenContract deploys a Link Token contract to an EVM chain
func (e *EthereumContractDeployer) DeployLinkTokenContract() (LinkToken, error) {
	linkTokenAddress, _, instance, err := e.eth.DeployContract("LINK Token", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployLinkToken(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	//Set config address
	//e.eth.NetworkConfig.Config().LinkTokenAddress = linkTokenAddress.Hex()

	return &EthereumLinkToken{
		client:   e.eth,
		instance: instance.(*ethereum.LinkToken),
		address:  *linkTokenAddress,
	}, err
}

// DefaultOffChainAggregatorOptions returns some base defaults for deploying an OCR contract
func DefaultOffChainAggregatorOptions() OffchainOptions {
	return OffchainOptions{
		MaximumGasPrice:         uint32(500000000),
		ReasonableGasPrice:      uint32(28000),
		MicroLinkPerEth:         uint32(500),
		LinkGweiPerObservation:  uint32(500),
		LinkGweiPerTransmission: uint32(500),
		MinimumAnswer:           big.NewInt(1),
		MaximumAnswer:           big.NewInt(5000000000),
		Decimals:                8,
		Description:             "Test OCR",
	}
}

// DefaultOffChainAggregatorConfig returns some base defaults for configuring an OCR contract
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

// DeployOffChainAggregator deploys the offchain aggregation contract to the EVM chain
func (e *EthereumContractDeployer) DeployOffChainAggregator(
	linkAddr string,
	offchainOptions OffchainOptions,
) (OffchainAggregator, error) {
	address, _, instance, err := e.eth.DeployContract("OffChain Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		la := common.HexToAddress(linkAddr)
		return ethereum.DeployOffchainAggregator(auth,
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
		client:  e.eth,
		ocr:     instance.(*ethereum.OffchainAggregator),
		address: address,
	}, err
}

// Balance get deployer wallet balance
func (e *EthereumContractDeployer) Balance() (*big.Float, error) {
	balance, err := e.eth.Client.PendingBalanceAt(context.Background(), common.HexToAddress(e.eth.DefaultWallet.Address()))
	if err != nil {
		return nil, err
	}
	bf := new(big.Float).SetInt(balance)
	return big.NewFloat(1).Quo(bf, client.OneEth), nil
}

//
//// CalculateTXSCost calculates required amount of ETH for N transactions based on network suggested gas price and tx gas limit
//func (e *EthereumContractDeployer) CalculateTXSCost(fromWallet client.BlockchainWallet, networkConfig *config.NetworkConfig, txs int64) (*big.Float, error) {
//	txsLimit := networkConfig.TransactionLimit
//	gasPrice, err := e.eth.Client.SuggestGasPrice(context.Background())
//	if err != nil {
//		return nil, err
//	}
//	gpFloat := big.NewFloat(1).SetInt(gasPrice)
//	oneGWei := big.NewFloat(1).SetInt(client.OneGWei)
//	gpGWei := big.NewFloat(1).Quo(gpFloat, oneGWei)
//	log.Debug().Str("Gas price (GWei)", gpGWei.String()).Msg("Suggested gas price")
//	txl := big.NewFloat(1).SetUint64(txsLimit)
//	oneTx := big.NewFloat(1).Mul(txl, gpFloat)
//	transactions := big.NewFloat(1).SetInt64(txs)
//	totalWei := big.NewFloat(1).Mul(oneTx, transactions)
//	totalETH := big.NewFloat(1).Quo(totalWei, client.OneEth)
//	log.Debug().Str("ETH", totalETH.String()).Int64("TXs", txs).Msg("Calculated required ETH")
//	return totalETH, nil
//}

// DeployStorageContract deploys a vanilla storage contract that is a value store
func (e *EthereumContractDeployer) DeployStorageContract() (Storage, error) {
	_, _, instance, err := e.eth.DeployContract("Storage", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumStorage{
		client: e.eth,
		store:  instance.(*ethereum.Store),
	}, err
}

// DeployAPIConsumer deploys api consumer for oracle
func (e *EthereumContractDeployer) DeployAPIConsumer(linkAddr string) (APIConsumer, error) {
	addr, _, instance, err := e.eth.DeployContract("APIConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployAPIConsumer(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumAPIConsumer{
		address:  addr,
		client:   e.eth,
		consumer: instance.(*ethereum.APIConsumer),
	}, err
}

// DeployOracle deploys oracle for consumer test
func (e *EthereumContractDeployer) DeployOracle(linkAddr string) (Oracle, error) {
	addr, _, instance, err := e.eth.DeployContract("Oracle", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployOracle(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOracle{
		address: addr,
		client:  e.eth,
		oracle:  instance.(*ethereum.Oracle),
	}, err
}

// DeployVRFContract deploy VRF contract
func (e *EthereumContractDeployer) DeployVRFContract() (VRF, error) {
	address, _, instance, err := e.eth.DeployContract("VRF", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRF(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRF{
		client:  e.eth,
		vrf:     instance.(*ethereum.VRF),
		address: address,
	}, err
}

func (e *EthereumContractDeployer) DeployMockETHLINKFeed(answer *big.Int) (MockETHLINKFeed, error) {
	address, _, instance, err := e.eth.DeployContract("MockETHLINKFeed", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployMockETHLINKAggregator(auth, backend, answer)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockETHLINKFeed{
		client:  e.eth,
		feed:    instance.(*ethereum.MockETHLINKAggregator),
		address: address,
	}, err
}

func (e *EthereumContractDeployer) DeployMockGasFeed(answer *big.Int) (MockGasFeed, error) {
	address, _, instance, err := e.eth.DeployContract("MockGasFeed", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployMockGASAggregator(auth, backend, answer)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockGASFeed{
		client:  e.eth,
		feed:    instance.(*ethereum.MockGASAggregator),
		address: address,
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepRegistrationRequests(linkAddr string, minLinkJuels *big.Int) (UpkeepRegistrar, error) {
	address, _, instance, err := e.eth.DeployContract("UpkeepRegistrationRequests", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployUpkeepRegistrationRequests(auth, backend, common.HexToAddress(linkAddr), minLinkJuels)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumUpkeepRegistrationRequests{
		client:    e.eth,
		registrar: instance.(*ethereum.UpkeepRegistrationRequests),
		address:   address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperRegistry(
	opts *KeeperRegistryOpts,
) (KeeperRegistry, error) {
	address, _, instance, err := e.eth.DeployContract("KeeperRegistry", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployKeeperRegistry(
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
		client:   e.eth,
		registry: instance.(*ethereum.KeeperRegistry),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperConsumer(updateInterval *big.Int) (KeeperConsumer, error) {
	address, _, instance, err := e.eth.DeployContract("KeeperConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployKeeperConsumer(auth, backend, updateInterval)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperConsumer{
		client:   e.eth,
		consumer: instance.(*ethereum.KeeperConsumer),
		address:  address,
	}, err
}

// DeployBlockhashStore deploys blockhash store used with VRF contract
func (e *EthereumContractDeployer) DeployBlockhashStore() (BlockHashStore, error) {
	address, _, instance, err := e.eth.DeployContract("BlockhashStore", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployBlockhashStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumBlockhashStore{
		client:         e.eth,
		blockHashStore: instance.(*ethereum.BlockhashStore),
		address:        address,
	}, err
}

// DeployVRFCoordinator deploys VRF coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinator(linkAddr string, bhsAddr string) (VRFCoordinator, error) {
	address, _, instance, err := e.eth.DeployContract("VRFCoordinator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRFCoordinator(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinator{
		client:      e.eth,
		coordinator: instance.(*ethereum.VRFCoordinator),
		address:     address,
	}, err
}

// DeployVRFConsumer deploys VRF consumer contract
func (e *EthereumContractDeployer) DeployVRFConsumer(linkAddr string, coordinatorAddr string) (VRFConsumer, error) {
	address, _, instance, err := e.eth.DeployContract("VRFConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRFConsumer(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumer{
		client:   e.eth,
		consumer: instance.(*ethereum.VRFConsumer),
		address:  address,
	}, err
}
