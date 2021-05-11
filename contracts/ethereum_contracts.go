package contracts

import (
	"context"
	"integrations-framework/client"
	"integrations-framework/contracts/ethereum"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// EthereumFluxAggregator represents the basic flux aggregation contract
type EthereumFluxAggregator struct {
	client         *client.EthereumClient
	fluxAggregator *ethereum.FluxAggregator
	callerWallet   client.BlockchainWallet
}

// DeployFluxAggregatorContract deploys the Flux Aggregator Contract on an EVM chain
func DeployFluxAggregatorContract(
	ethClient *client.EthereumClient,
	fromWallet client.BlockchainWallet,
) (FluxAggregator, error) {

	_, _, instance, err := ethClient.DeployContract(fromWallet, func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		// Some defaults for deploying this test contract
		linkAddress := common.HexToAddress(ethClient.Network.Config().LinkTokenAddress)
		paymentAmount := big.NewInt(1)
		timeout := uint32(60)
		var validator common.Address
		minSubValue := big.NewInt(1)
		maxSubValue := big.NewInt(10)
		decimals := uint8(18)
		desc := "Test Flux Aggregator"
		return ethereum.DeployFluxAggregator(auth, backend, linkAddress, paymentAmount, timeout, validator,
			minSubValue, maxSubValue, decimals, desc)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFluxAggregator{
		client:         ethClient,
		fluxAggregator: instance.(*ethereum.FluxAggregator),
		callerWallet:   fromWallet,
	}, nil
}

// Description returns the description of the flux aggregator contract
func (f *EthereumFluxAggregator) Description(ctxt context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.callerWallet.Address()),
		Pending: true,
		Context: ctxt,
	}
	return f.fluxAggregator.Description(opts)
}

// EthereumLinkToken represents a LinkToken address
type EthereumLinkToken struct {
	client       *client.EthereumClient
	linkToken    *ethereum.LinkToken
	callerWallet client.BlockchainWallet
}

// DeployLinkTokenContract deploys a Link Token contract to an EVM chain
func DeployLinkTokenContract(ethClient *client.EthereumClient, fromWallet client.BlockchainWallet) (LinkToken, error) {
	// First check if link token is already deployed
	linkTokenAddress := ethClient.Network.Config().LinkTokenAddress
	if linkTokenAddress != "" {
		tokenInstance, err := ethereum.NewLinkToken(common.HexToAddress(linkTokenAddress), ethClient.Client)
		if err != nil {
			return nil, err
		}
		return &EthereumLinkToken{
			client:       ethClient,
			linkToken:    tokenInstance,
			callerWallet: fromWallet,
		}, err
	}

	// Otherwise, deploy a new one
	_, _, instance, err := ethClient.DeployContract(fromWallet, func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployLinkToken(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumLinkToken{
		client:       ethClient,
		linkToken:    instance.(*ethereum.LinkToken),
		callerWallet: fromWallet,
	}, err
}

// Name returns the name of the link token
func (l *EthereumLinkToken) Name(ctxt context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(l.callerWallet.Address()),
		Pending: true,
		Context: ctxt,
	}
	return l.linkToken.Name(opts)
}

type EthereumOffchainAggregator struct {
	client       *client.EthereumClient
	ocr          *ethereum.OffchainAggregator
	callerWallet client.BlockchainWallet
}

func DeployOffChainAggregator(
	ethClient *client.EthereumClient,
	fromWallet client.BlockchainWallet,
) (OffchainAggregator, error) {

	_, _, instance, err := ethClient.DeployContract(fromWallet, func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		// This is complicated, want to wait a bit to clarify approach before implementing this
		return common.Address{}, nil, nil, nil
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOffchainAggregator{
		client:       ethClient,
		ocr:          instance.(*ethereum.OffchainAggregator),
		callerWallet: fromWallet,
	}, err
}

// Link returns the LINK contract address on the EVM chain
func (o *EthereumOffchainAggregator) Link(ctxt context.Context) (common.Address, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.callerWallet.Address()),
		Pending: true,
		Context: ctxt,
	}
	return o.ocr.LINK(opts)
}

// EthereumStorage acts as a conduit for the ethereum version of the storage contract
type EthereumStorage struct {
	client       *client.EthereumClient
	store        *ethereum.Store
	callerWallet client.BlockchainWallet
}

// DeployStorageContract deploys a vanilla storage contract that is a value store
func DeployStorageContract(ethClient *client.EthereumClient, fromWallet client.BlockchainWallet) (Storage, error) {
	_, _, instance, err := ethClient.DeployContract(fromWallet, func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumStorage{
		client:       ethClient,
		store:        instance.(*ethereum.Store),
		callerWallet: fromWallet,
	}, err
}

// Set sets a value in the storage contract
func (e *EthereumStorage) Set(ctxt context.Context, value *big.Int) error {
	opts, err := e.client.TransactionOpts(e.callerWallet, common.Address{}, big.NewInt(0), common.Hash{})
	if err != nil {
		return err
	}

	transaction, err := e.store.Set(opts, value)
	if err != nil {
		return err
	}
	return e.client.WaitForTransaction(transaction.Hash())
}

// Get retrieves a set value from the storage contract
func (e *EthereumStorage) Get(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.callerWallet.Address()),
		Pending: true,
		Context: ctxt,
	}
	return e.store.Get(opts)
}

// EthereumVRF represents a VRF contract
type EthereumVRF struct {
	client       *client.EthereumClient
	vrf          *ethereum.VRF
	callerWallet client.BlockchainWallet
}

// DeployVRFContract deploys a VRF contract
func DeployVRFContract(ethClient *client.EthereumClient, fromWallet client.BlockchainWallet) (VRF, error) {
	_, _, instance, err := ethClient.DeployContract(fromWallet, func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRF(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRF{
		client:       ethClient,
		vrf:          instance.(*ethereum.VRF),
		callerWallet: fromWallet,
	}, err
}

// ProofLength returns the PROOFLENGTH call from the VRF contract
func (v *EthereumVRF) ProofLength(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.callerWallet.Address()),
		Pending: true,
		Context: ctxt,
	}
	return v.vrf.PROOFLENGTH(opts)
}
