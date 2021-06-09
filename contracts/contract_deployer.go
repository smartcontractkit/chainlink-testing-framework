package contracts

import (
	"errors"
	"integrations-framework/client"
	"integrations-framework/contracts/ethereum"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
)

// ContractDeployer is an interface for abstracting the contract deployment methods across network implementations
type ContractDeployer interface {
	DeployStorageContract(fromWallet client.BlockchainWallet) (Storage, error)
	DeployFluxAggregatorContract(
		fromWallet client.BlockchainWallet,
		fluxOptions FluxAggregatorOptions,
	) (FluxAggregator, error)
	DeployLinkTokenContract(fromWallet client.BlockchainWallet) (LinkToken, error)
	DeployOffChainAggregator(
		fromWallet client.BlockchainWallet,
		offchainOptions OffchainOptions,
	) (OffchainAggregator, error)
	DeployVRFContract(fromWallet client.BlockchainWallet) (VRF, error)
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
		Timeout:       uint32(5),
		MinSubValue:   big.NewInt(1),
		MaxSubValue:   big.NewInt(10),
		Decimals:      uint8(8),
		Description:   "Hardhat Flux Aggregator",
	}
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
		return ethereum.DeployFluxAggregator(auth,
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
		fluxAggregator: instance.(*ethereum.FluxAggregator),
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
		return ethereum.DeployLinkToken(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	// Set config address
	e.eth.Network.Config().LinkTokenAddress = linkTokenAddress.Hex()

	return &EthereumLinkToken{
		client:       e.eth,
		linkToken:    instance.(*ethereum.LinkToken),
		callerWallet: fromWallet,
		address:      *linkTokenAddress,
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
		MaximumAnswer:           big.NewInt(5000),
		Decimals:                8,
		Description:             "Test OCR",
	}
}

// DefaultOffChainAggregatorConfig returns some base defaults for configuring an OCR contract
func DefaultOffChainAggregatorConfig() OffChainAggregatorConfig {
	return OffChainAggregatorConfig{
		AlphaPPB:         1,
		DeltaC:           time.Second * 15,
		DeltaGrace:       time.Second,
		DeltaProgress:    time.Second * 30,
		DeltaStage:       time.Second * 3,
		DeltaResend:      time.Second * 5,
		DeltaRound:       time.Second * 10,
		RMax:             4,
		S:                []int{1, 1, 1, 1, 1},
		N:                5,
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
		return ethereum.DeployOffchainAggregator(auth,
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
		ocr:          instance.(*ethereum.OffchainAggregator),
		callerWallet: fromWallet,
		address:      address,
	}, err
}

// DeployStorageContract deploys a vanilla storage contract that is a value store
func (e *EthereumContractDeployer) DeployStorageContract(fromWallet client.BlockchainWallet) (Storage, error) {
	_, _, instance, err := e.eth.DeployContract(fromWallet, "Storage", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumStorage{
		client:       e.eth,
		store:        instance.(*ethereum.Store),
		callerWallet: fromWallet,
	}, err
}

func (e *EthereumContractDeployer) DeployVRFContract(fromWallet client.BlockchainWallet) (VRF, error) {
	address, _, instance, err := e.eth.DeployContract(fromWallet, "VRF", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRF(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRF{
		client:       e.eth,
		vrf:          instance.(*ethereum.VRF),
		callerWallet: fromWallet,
		address:      address,
	}, err
}
