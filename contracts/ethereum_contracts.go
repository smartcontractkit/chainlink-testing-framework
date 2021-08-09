package contracts

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts/ethereum"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"strings"
)

// EthereumOracle oracle for "directrequest" job tests
type EthereumOracle struct {
	address      *common.Address
	client       *client.EthereumClient
	oracle       *ethereum.Oracle
	callerWallet client.BlockchainWallet
}

func (e *EthereumOracle) Address() string {
	return e.address.Hex()
}

func (e *EthereumOracle) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error {
	return e.client.Fund(fromWallet, e.address.Hex(), ethAmount, linkAmount)
}

// SetFulfillmentPermission sets fulfillment permission for particular address
func (e *EthereumOracle) SetFulfillmentPermission(fromWallet client.BlockchainWallet, address string, allowed bool) error {
	opts, err := e.client.TransactionOpts(fromWallet, *e.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := e.oracle.SetFulfillmentPermission(opts, common.HexToAddress(address), allowed)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx.Hash())
}

// EthereumAPIConsumer API consumer for job type "directrequest" tests
type EthereumAPIConsumer struct {
	address      *common.Address
	client       *client.EthereumClient
	consumer     *ethereum.APIConsumer
	callerWallet client.BlockchainWallet
}

func (e *EthereumAPIConsumer) Address() string {
	return e.address.Hex()
}

func (e *EthereumAPIConsumer) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error {
	return e.client.Fund(fromWallet, e.address.Hex(), ethAmount, linkAmount)
}

func (e *EthereumAPIConsumer) Data(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	data, err := e.consumer.Data(opts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateRequestTo creates request to an oracle for particular jobID with params
func (e *EthereumAPIConsumer) CreateRequestTo(
	fromWallet client.BlockchainWallet,
	oracleAddr string,
	jobID [32]byte,
	payment *big.Int,
	url string,
	path string,
	times *big.Int,
) error {
	opts, err := e.client.TransactionOpts(fromWallet, *e.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := e.consumer.CreateRequestTo(opts, common.HexToAddress(oracleAddr), jobID, payment, url, path, times)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx.Hash())
}

// EthereumFluxAggregator represents the basic flux aggregation contract
type EthereumFluxAggregator struct {
	client         *client.EthereumClient
	fluxAggregator *ethereum.FluxAggregator
	callerWallet   client.BlockchainWallet
	address        *common.Address
}

func (f *EthereumFluxAggregator) Address() string {
	return f.address.Hex()
}

// Fund sends specified currencies to the contract
func (f *EthereumFluxAggregator) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error {
	return f.client.Fund(fromWallet, f.address.Hex(), ethAmount, linkAmount)
}

func (f *EthereumFluxAggregator) UpdateAvailableFunds(ctx context.Context, fromWallet client.BlockchainWallet) error {
	opts, err := f.client.TransactionOpts(fromWallet, *f.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.UpdateAvailableFunds(opts)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx.Hash())
}

func (f *EthereumFluxAggregator) PaymentAmount(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	payment, err := f.fluxAggregator.PaymentAmount(opts)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (f *EthereumFluxAggregator) RequestNewRound(ctx context.Context, fromWallet client.BlockchainWallet) error {
	opts, err := f.client.TransactionOpts(fromWallet, *f.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.RequestNewRound(opts)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx.Hash())
}

func (f *EthereumFluxAggregator) SetRequesterPermissions(ctx context.Context, fromWallet client.BlockchainWallet, addr common.Address, authorized bool, roundsDelay uint32) error {
	opts, err := f.client.TransactionOpts(fromWallet, *f.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.SetRequesterPermissions(opts, addr, authorized, roundsDelay)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx.Hash())
}

func (f *EthereumFluxAggregator) GetOracles(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	addresses, err := f.fluxAggregator.GetOracles(opts)
	if err != nil {
		return nil, err
	}
	var oracleAddrs []string
	for _, o := range addresses {
		oracleAddrs = append(oracleAddrs, o.Hex())
	}
	return oracleAddrs, nil
}

func (f *EthereumFluxAggregator) LatestRound(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	rID, err := f.fluxAggregator.LatestRound(opts)
	if err != nil {
		return nil, err
	}
	return rID, nil
}

// AwaitNextRoundFinalized awaits for the next round to be finalized
func (f *EthereumFluxAggregator) AwaitNextRoundFinalized(ctx context.Context) error {
	lr, err := f.LatestRound(ctx)
	if err != nil {
		return err
	}
	log.Info().Int64("Round", lr.Int64()).Msg("Awaiting next round after")
	if err := retry.Do(func() error {
		newRound, err := f.LatestRound(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to get round in retry loop")
		}
		if newRound.Cmp(lr) <= 0 {
			return errors.New("awaiting new round")
		}
		return nil
	}, retry.Attempts(60)); err != nil {
		return err
	}
	return nil
}

func (f *EthereumFluxAggregator) WithdrawPayment(
	ctx context.Context,
	caller client.BlockchainWallet,
	from common.Address,
	to common.Address,
	amount *big.Int) error {
	opts, err := f.client.TransactionOpts(caller, *f.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.WithdrawPayment(opts, from, to, amount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx.Hash())
}

func (f *EthereumFluxAggregator) WithdrawablePayment(ctx context.Context, addr common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	balance, err := f.fluxAggregator.WithdrawablePayment(opts, addr)
	if err != nil {
		return nil, err
	}
	return balance, nil
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
	o SetOraclesOptions) error {
	opts, err := f.client.TransactionOpts(fromWallet, *f.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}

	tx, err := f.fluxAggregator.ChangeOracles(opts, o.RemoveList, o.AddList, o.AdminList, o.MinSubmissions, o.MaxSubmissions, o.RestartDelayRounds)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx.Hash())
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

// Fund the LINK Token contract with ETH to distribute the token
func (l *EthereumLinkToken) Fund(fromWallet client.BlockchainWallet, ethAmount *big.Float) error {
	return l.client.Fund(fromWallet, l.address.Hex(), ethAmount, nil)
}

func (l *EthereumLinkToken) BalanceOf(ctx context.Context, addr common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(l.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	balance, err := l.linkToken.BalanceOf(opts, addr)
	if err != nil {
		return nil, err
	}
	return balance, nil
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

func (l *EthereumLinkToken) Address() string {
	return l.address.Hex()
}

func (l *EthereumLinkToken) Approve(fromWallet client.BlockchainWallet, to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(fromWallet, l.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := l.linkToken.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx.Hash())
}

func (l *EthereumLinkToken) Transfer(fromWallet client.BlockchainWallet, to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(fromWallet, l.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := l.linkToken.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx.Hash())
}

func (l *EthereumLinkToken) TransferAndCall(fromWallet client.BlockchainWallet, to string, amount *big.Int, data []byte) error {
	opts, err := l.client.TransactionOpts(fromWallet, l.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := l.linkToken.TransferAndCall(opts, common.HexToAddress(to), amount, data)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx.Hash())
}

// EthereumOffchainAggregator represents the offchain aggregation contract
type EthereumOffchainAggregator struct {
	client       *client.EthereumClient
	ocr          *ethereum.OffchainAggregator
	callerWallet client.BlockchainWallet
	address      *common.Address
}

// Fund sends specified currencies to the contract
func (o *EthereumOffchainAggregator) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error {
	return o.client.Fund(fromWallet, o.address.Hex(), ethAmount, linkAmount)
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
	return o.client.ProcessTransaction(tx.Hash())
}

// SetConfig sets offchain reporting protocol configuration including participating oracles
func (o *EthereumOffchainAggregator) SetConfig(
	fromWallet client.BlockchainWallet,
	chainlinkNodes []client.Chainlink,
	ocrConfig OffChainAggregatorConfig,
) error {
	// Gather necessary addresses and keys from our chainlink nodes to properly configure the OCR contract
	log.Info().Str("Contract Address", o.address.Hex()).Msg("Configuring OCR Contract")
	for _, node := range chainlinkNodes {
		ocrKeys, err := node.ReadOCRKeys()
		if err != nil {
			return err
		}
		primaryOCRKey := ocrKeys.Data[0]
		primaryEthKey, err := node.PrimaryEthAddress()
		if err != nil {
			return err
		}
		p2pKeys, err := node.ReadP2PKeys()
		if err != nil {
			return err
		}
		primaryP2PKey := p2pKeys.Data[0]

		// Need to convert the key representations
		var onChainSigningAddress [20]byte
		var configPublicKey [32]byte
		offchainSigningAddress, err := hex.DecodeString(primaryOCRKey.Attributes.OffChainPublicKey)
		if err != nil {
			return err
		}
		decodeConfigKey, err := hex.DecodeString(primaryOCRKey.Attributes.ConfigPublicKey)
		if err != nil {
			return err
		}

		// https://stackoverflow.com/questions/8032170/how-to-assign-string-to-bytes-array
		copy(onChainSigningAddress[:], common.HexToAddress(primaryOCRKey.Attributes.OnChainSigningAddress).Bytes())
		copy(configPublicKey[:], decodeConfigKey)

		oracleIdentity := ocrConfigHelper.OracleIdentity{
			TransmitAddress:       common.HexToAddress(primaryEthKey),
			OnChainSigningAddress: onChainSigningAddress,
			PeerID:                primaryP2PKey.Attributes.PeerID,
			OffchainPublicKey:     offchainSigningAddress,
		}
		oracleIdentityExtra := ocrConfigHelper.OracleIdentityExtra{
			OracleIdentity:                  oracleIdentity,
			SharedSecretEncryptionPublicKey: ocrTypes.SharedSecretEncryptionPublicKey(configPublicKey),
		}

		ocrConfig.OracleIdentities = append(ocrConfig.OracleIdentities, oracleIdentityExtra)
	}

	signers, transmitters, threshold, encodedConfigVersion, encodedConfig, err := ocrConfigHelper.ContractSetConfigArgs(
		ocrConfig.DeltaProgress,
		ocrConfig.DeltaResend,
		ocrConfig.DeltaRound,
		ocrConfig.DeltaGrace,
		ocrConfig.DeltaC,
		ocrConfig.AlphaPPB,
		ocrConfig.DeltaStage,
		ocrConfig.RMax,
		ocrConfig.S,
		ocrConfig.OracleIdentities,
		ocrConfig.F,
	)
	if err != nil {
		return err
	}

	// Set Payees
	opts, err := o.client.TransactionOpts(fromWallet, *o.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := o.ocr.SetPayees(opts, transmitters, transmitters)
	if err != nil {
		return err
	}
	if err := o.client.ProcessTransaction(tx.Hash()); err != nil {
		return err
	}

	// Set Config
	opts, err = o.client.TransactionOpts(fromWallet, *o.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err = o.ocr.SetConfig(opts, signers, transmitters, threshold, encodedConfigVersion, encodedConfig)
	if err != nil {
		return err
	}
	return o.client.ProcessTransaction(tx.Hash())
}

// RequestNewRound requests the OCR contract to create a new round
func (o *EthereumOffchainAggregator) RequestNewRound(fromWallet client.BlockchainWallet) error {
	opts, err := o.client.TransactionOpts(fromWallet, *o.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := o.ocr.RequestNewRound(opts)
	if err != nil {
		return err
	}
	log.Info().Str("Contract Address", o.address.Hex()).Msg("New OCR round requested")

	return o.client.ProcessTransaction(tx.Hash())
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

func (o *EthereumOffchainAggregator) Address() string {
	return o.address.Hex()
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
	return e.client.ProcessTransaction(transaction.Hash())
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

// Fund sends specified currencies to the contract
func (v *EthereumVRF) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error {
	return v.client.Fund(fromWallet, v.address.Hex(), ethAmount, linkAmount)
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

// EthereumMockETHLINKFeed represents mocked ETH/LINK feed contract
type EthereumMockETHLINKFeed struct {
	client       *client.EthereumClient
	feed         *ethereum.MockETHLINKAggregator
	callerWallet client.BlockchainWallet
	address      *common.Address
}

func (v *EthereumMockETHLINKFeed) Address() string {
	return v.address.Hex()
}

// EthereumMockGASFeed represents mocked Gas feed contract
type EthereumMockGASFeed struct {
	client       *client.EthereumClient
	feed         *ethereum.MockGASAggregator
	callerWallet client.BlockchainWallet
	address      *common.Address
}

func (v *EthereumMockGASFeed) Address() string {
	return v.address.Hex()
}

// EthereumKeeperRegistry represents keeper registry contract
type EthereumKeeperRegistry struct {
	client       *client.EthereumClient
	registry     *ethereum.KeeperRegistry
	callerWallet client.BlockchainWallet
	address      *common.Address
}

func (v *EthereumKeeperRegistry) Address() string {
	return v.address.Hex()
}

func (v *EthereumKeeperRegistry) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error {
	return v.client.Fund(fromWallet, v.address.Hex(), ethAmount, linkAmount)
}

func (v *EthereumKeeperRegistry) SetRegistrar(fromWallet client.BlockchainWallet, registrarAddr string) error {
	opts, err := v.client.TransactionOpts(fromWallet, *v.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := v.registry.SetRegistrar(opts, common.HexToAddress(registrarAddr))
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx.Hash())
}

// AddUpkeepFunds adds link for particular upkeep id
func (v *EthereumKeeperRegistry) AddUpkeepFunds(fromWallet client.BlockchainWallet, id *big.Int, amount *big.Int) error {
	opts, err := v.client.TransactionOpts(fromWallet, *v.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := v.registry.AddFunds(opts, id, amount)
	if err != nil {
		return err
	}
	if err := v.client.ProcessTransaction(tx.Hash()); err != nil {
		return err
	}
	return nil
}

// GetUpkeepInfo gets upkeep info
func (v *EthereumKeeperRegistry) GetUpkeepInfo(ctx context.Context, id *big.Int) (*UpkeepInfo, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	uk, err := v.registry.GetUpkeep(opts, id)
	if err != nil {
		return nil, err
	}
	return &UpkeepInfo{
		Target:              uk.Target.Hex(),
		ExecuteGas:          uk.ExecuteGas,
		CheckData:           uk.CheckData,
		Balance:             uk.Balance,
		LastKeeper:          uk.LastKeeper.Hex(),
		Admin:               uk.Admin.Hex(),
		MaxValidBlocknumber: uk.MaxValidBlocknumber,
	}, nil
}

func (v *EthereumKeeperRegistry) GetKeeperInfo(ctx context.Context, keeperAddr string) (*KeeperInfo, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	info, err := v.registry.GetKeeperInfo(opts, common.HexToAddress(keeperAddr))
	if err != nil {
		return nil, err
	}
	return &KeeperInfo{
		Payee:   info.Payee.Hex(),
		Active:  info.Active,
		Balance: info.Balance,
	}, nil
}

func (v *EthereumKeeperRegistry) SetKeepers(fromWallet client.BlockchainWallet, keepers []string, payees []string) error {
	opts, err := v.client.TransactionOpts(fromWallet, *v.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	keepersAddresses := make([]common.Address, 0)
	for _, k := range keepers {
		keepersAddresses = append(keepersAddresses, common.HexToAddress(k))
	}
	payeesAddresses := make([]common.Address, 0)
	for _, p := range payees {
		payeesAddresses = append(payeesAddresses, common.HexToAddress(p))
	}
	tx, err := v.registry.SetKeepers(opts, keepersAddresses, payeesAddresses)
	if err != nil {
		return err
	}
	if err := v.client.ProcessTransaction(tx.Hash()); err != nil {
		return err
	}
	return nil
}

// RegisterUpkeep registers contract to perform upkeep
func (v *EthereumKeeperRegistry) RegisterUpkeep(fromWallet client.BlockchainWallet, target string, gasLimit uint32, admin string, checkData []byte) error {
	opts, err := v.client.TransactionOpts(fromWallet, *v.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := v.registry.RegisterUpkeep(opts, common.HexToAddress(target), gasLimit, common.HexToAddress(admin), checkData)
	if err != nil {
		return err
	}
	if err := v.client.ProcessTransaction(tx.Hash()); err != nil {
		return err
	}
	return nil
}

// GetKeeperList get list of all registered keeper addresses
func (v *EthereumKeeperRegistry) GetKeeperList(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	list, err := v.registry.GetKeeperList(opts)
	if err != nil {
		return []string{}, err
	}
	addrs := make([]string, 0)
	for _, ca := range list {
		addrs = append(addrs, ca.Hex())
	}
	return addrs, nil
}

// EthereumKeeperConsumer represents keeper consumer (upkeep) contract
type EthereumKeeperConsumer struct {
	client       *client.EthereumClient
	consumer     *ethereum.KeeperConsumer
	callerWallet client.BlockchainWallet
	address      *common.Address
}

func (v *EthereumKeeperConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumKeeperConsumer) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error {
	return v.client.Fund(fromWallet, v.address.Hex(), ethAmount, linkAmount)
}

func (v *EthereumKeeperConsumer) Counter(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	cnt, err := v.consumer.Counter(opts)
	if err != nil {
		return nil, err
	}
	return cnt, nil
}

// EthereumUpkeepRegistrationRequests keeper contract to register upkeeps
type EthereumUpkeepRegistrationRequests struct {
	client       *client.EthereumClient
	registrar    *ethereum.UpkeepRegistrationRequests
	callerWallet client.BlockchainWallet
	address      *common.Address
}

func (v *EthereumUpkeepRegistrationRequests) Address() string {
	return v.address.Hex()
}

// SetRegistrarConfig sets registrar config, allowing auto register or pending requests for manual registration
func (v *EthereumUpkeepRegistrationRequests) SetRegistrarConfig(
	fromWallet client.BlockchainWallet,
	autoRegister bool,
	windowSizeBlocks uint32,
	allowedPerWindow uint16,
	registryAddr string,
	minLinkJuels *big.Int,
) error {
	opts, err := v.client.TransactionOpts(fromWallet, *v.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := v.registrar.SetRegistrationConfig(opts, autoRegister, windowSizeBlocks, allowedPerWindow, common.HexToAddress(registryAddr), minLinkJuels)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx.Hash())
}

func (v *EthereumUpkeepRegistrationRequests) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error {
	return v.client.Fund(fromWallet, v.address.Hex(), ethAmount, linkAmount)
}

// EncodeRegisterRequest encodes register request to call it through link token TransferAndCall
func (v *EthereumUpkeepRegistrationRequests) EncodeRegisterRequest(
	name string,
	email []byte,
	upkeepAddr string,
	gasLimit uint32,
	adminAddr string,
	checkData []byte,
	amount *big.Int,
	source uint8,
) ([]byte, error) {
	registryABI, err := abi.JSON(strings.NewReader(ethereum.UpkeepRegistrationRequestsABI))
	if err != nil {
		return nil, err
	}
	req, err := registryABI.Pack(
		"register",
		name,
		email,
		common.HexToAddress(upkeepAddr),
		gasLimit,
		common.HexToAddress(adminAddr),
		checkData,
		amount,
		source,
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// EthereumBlockhashStore represents a blockhash store for VRF contract
type EthereumBlockhashStore struct {
	address        *common.Address
	client         *client.EthereumClient
	blockHashStore *ethereum.BlockhashStore
	callerWallet   client.BlockchainWallet
}

func (v *EthereumBlockhashStore) Address() string {
	return v.address.Hex()
}

// EthereumVRFCoordinator represents VRF coordinator contract
type EthereumVRFCoordinator struct {
	address      *common.Address
	client       *client.EthereumClient
	coordinator  *ethereum.VRFCoordinator
	callerWallet client.BlockchainWallet
}

func (v *EthereumVRFCoordinator) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFCoordinator) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	hash, err := v.coordinator.HashOfKey(opts, pubKey)
	if err != nil {
		return [32]byte{}, err
	}
	return hash, nil
}

func (v *EthereumVRFCoordinator) RegisterProvingKey(
	fromWallet client.BlockchainWallet,
	fee *big.Int,
	oracleAddr string,
	publicProvingKey [2]*big.Int,
	jobID [32]byte,
) error {
	opts, err := v.client.TransactionOpts(fromWallet, *v.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := v.coordinator.RegisterProvingKey(opts, fee, common.HexToAddress(oracleAddr), publicProvingKey, jobID)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx.Hash())
}

// EthereumVRFConsumer represents VRF consumer contract
type EthereumVRFConsumer struct {
	address      *common.Address
	client       *client.EthereumClient
	consumer     *ethereum.VRFConsumer
	callerWallet client.BlockchainWallet
}

func (v *EthereumVRFConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFConsumer) Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error {
	return v.client.Fund(fromWallet, v.address.Hex(), ethAmount, linkAmount)
}

func (v *EthereumVRFConsumer) RequestRandomness(fromWallet client.BlockchainWallet, hash [32]byte, fee *big.Int) error {
	opts, err := v.client.TransactionOpts(fromWallet, *v.address, big.NewInt(0), nil)
	if err != nil {
		return err
	}
	tx, err := v.consumer.TestRequestRandomness(opts, hash, fee)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx.Hash())
}

func (v *EthereumVRFConsumer) RandomnessOutput(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.callerWallet.Address()),
		Pending: true,
		Context: ctx,
	}
	out, err := v.consumer.RandomnessOutput(opts)
	if err != nil {
		return nil, err
	}
	return out, nil
}
