package contracts

import (
	"bytes"
	"context"
	"encoding/json"
	"integrations-framework/client"
	"integrations-framework/contracts/ethereum"
	"math/big"
	"os"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// EthereumFluxAggregator represents the basic flux aggregation contract
type EthereumFluxAggregator struct {
	client         *client.EthereumClient
	fluxAggregator *ethereum.FluxAggregator
	callerWallet   client.BlockchainWallet
	address        *common.Address
}

// DeployFluxAggregatorContract deploys the Flux Aggregator Contract on an EVM chain
func DeployFluxAggregatorContract(
	ethClient *client.EthereumClient,
	fromWallet client.BlockchainWallet,
	fluxOptions FluxAggregatorOptions,
) (FluxAggregator, error) {

	address, _, instance, err := ethClient.DeployContract(fromWallet, "Flux Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		linkAddress := common.HexToAddress(ethClient.Network.Config().LinkTokenAddress)
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
		client:         ethClient,
		fluxAggregator: instance.(*ethereum.FluxAggregator),
		callerWallet:   fromWallet,
		address:        address,
	}, nil
}

// Fund sends specified currencies to the contract
func (f *EthereumFluxAggregator) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Int) error {
	return fund(f.client, fromWallet, *f.address, ethAmount, linkAmount)
}

// GetContractData retrieves basic data for the flux aggregator contract
func (f *EthereumFluxAggregator) GetContractData(ctxt context.Context) (*FluxAggregatorData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.callerWallet.Address()),
		Pending: true,
		Context: ctxt,
	}

	allocated, err := f.fluxAggregator.AllocatedFunds(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}

	available, err := f.fluxAggregator.AvailableFunds(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}

	lr, err := f.fluxAggregator.LatestRoundData(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}
	latestRound := RoundData(lr)

	oracles, err := f.fluxAggregator.GetOracles(opts)
	if err != nil {
		return &FluxAggregatorData{}, err
	}

	return &FluxAggregatorData{
		AllocatedFunds:  allocated,
		AvailableFunds:  available,
		LatestRoundData: latestRound,
		Oracles:         oracles,
	}, nil
}

// SetOracles allows the ability to add and/or remove oracles from the contract, and to set admins
func (f *EthereumFluxAggregator) SetOracles(
	fromWallet client.BlockchainWallet,
	toAdd, toRemove, toAdmin []common.Address,
	minSubmissions, maxSubmissions, restartDelay uint32) error {
	opts, err := f.client.TransactionOpts(fromWallet, *f.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}

	tx, err := f.fluxAggregator.ChangeOracles(opts, toRemove, toAdd, toAdmin, minSubmissions, maxSubmissions, restartDelay)
	if err != nil {
		return err
	}
	return f.client.WaitForTransaction(tx.Hash())
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
	address      common.Address
}

// DeployLinkTokenContract deploys a Link Token contract to an EVM chain
func DeployLinkTokenContract(ethClient *client.EthereumClient, fromWallet client.BlockchainWallet) (LinkToken, error) {
	linkTokenAddress, _, instance, err := ethClient.DeployContract(fromWallet, "LINK Token", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployLinkToken(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	// Set config address
	ethClient.Network.Config().LinkTokenAddress = linkTokenAddress.Hex()

	return &EthereumLinkToken{
		client:       ethClient,
		linkToken:    instance.(*ethereum.LinkToken),
		callerWallet: fromWallet,
		address:      *linkTokenAddress,
	}, err
}

// Address of the the link token address
func (l *EthereumLinkToken) Address() string {
	return l.address.Hex()
}

// Fund the LINK Token contract with ETH to distribute the token
func (l *EthereumLinkToken) Fund(fromWallet client.BlockchainWallet, ethAmount *big.Int) error {
	return fund(l.client, fromWallet, l.address, ethAmount, nil)
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

// EthereumOffchainAggregator represents the offchain aggregation contract
type EthereumOffchainAggregator struct {
	client       *client.EthereumClient
	ocr          *ethereum.OffchainAggregator
	callerWallet client.BlockchainWallet
	address      *common.Address
}

// DeployOffChainAggregator deploys the offchain aggregation contract to the EVM chain, using supplied chainlink nodes
// for setting its configuration
func DeployOffChainAggregator(
	ethClient *client.EthereumClient,
	fromWallet client.BlockchainWallet,
) (OffchainAggregator, error) {
	address, _, instance, err := ethClient.DeployContract(fromWallet, "OffChain Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		linkAddress := common.HexToAddress(ethClient.Network.Config().LinkTokenAddress)
		// Defaults
		offchainOptions := OffchainOptions{
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
		client:       ethClient,
		ocr:          instance.(*ethereum.OffchainAggregator),
		callerWallet: fromWallet,
		address:      address,
	}, err
}

// Address of the the ocr contract
func (o *EthereumOffchainAggregator) Address() string {
	return o.address.Hex()
}

// Fund sends specified currencies to the contract
func (o *EthereumOffchainAggregator) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Int) error {
	return fund(o.client, fromWallet, *o.address, ethAmount, linkAmount)
}

// GetContractData retrieves basic data for the offchain aggregator contract
func (o *EthereumOffchainAggregator) GetContractData(ctxt context.Context) (*OffchainAggregatorData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.callerWallet.Address()),
		Pending: true,
		Context: ctxt,
	}

	lr, err := o.ocr.LatestRoundData(opts)
	if err != nil {
		return &OffchainAggregatorData{}, err
	}
	latestRound := RoundData(lr)

	return &OffchainAggregatorData{
		LatestRoundData: latestRound,
	}, nil
}

// SetPayees sets wallets for the contract to pay out to?
func (o *EthereumOffchainAggregator) SetPayees(
	fromWallet client.BlockchainWallet,
	transmitters, payees []common.Address,
) error {
	opts, err := o.client.TransactionOpts(fromWallet, *o.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}

	tx, err := o.ocr.SetPayees(opts, transmitters, payees)
	if err != nil {
		return err
	}
	return o.client.WaitForTransaction(tx.Hash())
}

// SetConfig sets offchain reporting protocol configuration including participating oracles
func (o *EthereumOffchainAggregator) SetConfig(
	fromWallet client.BlockchainWallet,
	chainlinkNodes []client.Chainlink,
) error {
	ocrConfig := OffChainAggregatorConfig{
		AlphaPPB:         1,
		DeltaC:           "120s",
		DeltaGrace:       "1s",
		DeltaProgress:    "30s",
		DeltaStage:       "3s",
		DeltaResend:      "10s",
		DeltaRound:       "20s",
		RMax:             4,
		N:                5,
		F:                1,
		OracleIdentities: []OracleIdentity{},
	}

	// Gather necessary addresses and keys from our chainlink nodes to properly configure the OCR contract
	for _, node := range chainlinkNodes {
		ocrKeys, err := node.ReadOCRKeys()
		if err != nil {
			return err
		}
		ethKeys, err := node.ReadETHKeys()
		if err != nil {
			return err
		}
		p2pKeys, err := node.ReadP2PKeys()
		if err != nil {
			return err
		}

		oracleIdentity := OracleIdentity{
			TransmitAddress:       ethKeys.Data[0].Attributes.Address,
			OnchainSigningAddress: ocrKeys.Data[0].Attributes.OnChainSigningAddress,
			PeerID:                p2pKeys.Data[0].Attributes.PeerID,
			OffchainPublicKey:     "0x" + ocrKeys.Data[0].Attributes.OffChainPublicKey,
			ConfigPublicKey:       "0x" + ocrKeys.Data[0].Attributes.ConfigPublicKey,
		}
		ocrConfig.OracleIdentities = append(ocrConfig.OracleIdentities, oracleIdentity)
	}

	p, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return err
	}
	projectDir := strings.TrimSpace(string(p))

	// Build a config file for Lorenz's prototype script to run from
	confWrap := map[string]interface{}{
		o.address.Hex(): ocrConfig,
	}

	confJson, err := json.MarshalIndent(confWrap, "", " ")
	if err != nil {
		return err
	}

	confFile, err := os.OpenFile(projectDir+"/config/ocr_config.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = confFile.Write(confJson)
	if err != nil {
		return err
	}
	err = confFile.Close()
	if err != nil {
		return err
	}

	// Run the script to configure the contract
	log.Info().Str("CMD", "./prototype -mode=configure -rpcurl="+o.client.Network.URL()+" -configfile="+
		projectDir+"/config/ocr_config.json").
		Msg("Running command")
	cmd := exec.Command(
		"./prototype", "-mode=configure",
		"-rpcurl="+o.client.Network.URL(),
		"-configfile="+projectDir+"/config/ocr_config.json",
	)
	cmd.Dir = projectDir + "/tools"
	var e bytes.Buffer
	cmd.Stderr = &e
	err = cmd.Run()
	if err != nil {
		log.Error().Str("ERROR", e.String()).Msg("STDERR")
	}
	return err
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

// GetLatestAnswer returns the latest answer from the OCR contract
func (o *EthereumOffchainAggregator) GetLatestAnswer(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.callerWallet.Address()),
		Pending: true,
		Context: ctxt,
	}
	return o.ocr.LatestAnswer(opts)
}

// GetLatestRound returns data from the latest round
func (o *EthereumOffchainAggregator) GetLatestRound(ctxt context.Context) (*RoundData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.callerWallet.Address()),
		Pending: true,
		Context: ctxt,
	}

	roundData, err := o.ocr.LatestRoundData(opts)
	if err != nil {
		return nil, err
	}

	return &RoundData{
		RoundId:         roundData.RoundId,
		Answer:          roundData.Answer,
		AnsweredInRound: roundData.AnsweredInRound,
		StartedAt:       roundData.StartedAt,
		UpdatedAt:       roundData.UpdatedAt,
	}, err
}

// EthereumStorage acts as a conduit for the ethereum version of the storage contract
type EthereumStorage struct {
	client       *client.EthereumClient
	store        *ethereum.Store
	callerWallet client.BlockchainWallet
}

// DeployStorageContract deploys a vanilla storage contract that is a value store
func DeployStorageContract(ethClient *client.EthereumClient, fromWallet client.BlockchainWallet) (Storage, error) {
	_, _, instance, err := ethClient.DeployContract(fromWallet, "Storage", func(
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
func (e *EthereumStorage) Set(value *big.Int) error {
	opts, err := e.client.TransactionOpts(e.callerWallet, common.Address{}, big.NewInt(0), nil)
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
	address      *common.Address
}

// DeployVRFContract deploys a VRF contract
func DeployVRFContract(ethClient *client.EthereumClient, fromWallet client.BlockchainWallet) (VRF, error) {
	address, _, instance, err := ethClient.DeployContract(fromWallet, "VRF", func(
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
		address:      address,
	}, err
}

// Fund sends specified currencies to the contract
func (v *EthereumVRF) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Int) error {
	return fund(v.client, fromWallet, *v.address, ethAmount, linkAmount)
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

func fund(
	ethClient *client.EthereumClient,
	fromWallet client.BlockchainWallet,
	toAddress common.Address,
	ethAmount, linkAmount *big.Int,
) error {

	// Send ETH if not 0
	if ethAmount != nil && big.NewInt(0).Cmp(ethAmount) != 0 {
		log.Info().
			Str("Token", "ETH").
			Str("From", fromWallet.Address()).
			Str("To", toAddress.Hex()).
			Str("Amount", ethAmount.String()).
			Msg("Funding Contract")
		_, err := ethClient.SendTransaction(fromWallet, toAddress, ethAmount, nil)
		if err != nil {
			return err
		}
	}

	// Send LINK if not 0
	if linkAmount != nil && big.NewInt(0).Cmp(linkAmount) != 0 {
		// Prepare data field for token tx
		log.Info().
			Str("Token", "LINK").
			Str("From", fromWallet.Address()).
			Str("To", toAddress.Hex()).
			Str("Amount", linkAmount.String()).
			Msg("Funding Contract")
		linkAddress := common.HexToAddress(ethClient.Network.Config().LinkTokenAddress)
		linkInstance, err := ethereum.NewLinkToken(linkAddress, ethClient.Client)
		if err != nil {
			return err
		}
		opts, err := ethClient.TransactionOpts(fromWallet, toAddress, nil, nil)
		if err != nil {
			return err
		}
		tx, err := linkInstance.Transfer(opts, toAddress, linkAmount)
		if err != nil {
			return err
		}

		err = ethClient.WaitForTransaction(tx.Hash())
		if err != nil {
			return err
		}
	}
	return nil
}
