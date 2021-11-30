package contracts

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts/ethereum"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

// EthereumOracle oracle for "directrequest" job tests
type EthereumOracle struct {
	address *common.Address
	client  *client.EthereumClient
	oracle  *ethereum.Oracle
}

func (e *EthereumOracle) Address() string {
	return e.address.Hex()
}

func (e *EthereumOracle) Fund(ethAmount *big.Float) error {
	return e.client.Fund(e.address.Hex(), ethAmount)
}

// SetFulfillmentPermission sets fulfillment permission for particular address
func (e *EthereumOracle) SetFulfillmentPermission(address string, allowed bool) error {
	opts, err := e.client.TransactionOpts(e.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := e.oracle.SetFulfillmentPermission(opts, common.HexToAddress(address), allowed)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// EthereumAPIConsumer API consumer for job type "directrequest" tests
type EthereumAPIConsumer struct {
	address  *common.Address
	client   *client.EthereumClient
	consumer *ethereum.APIConsumer
}

func (e *EthereumAPIConsumer) Address() string {
	return e.address.Hex()
}

func (e *EthereumAPIConsumer) RoundID(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.DefaultWallet.Address()),
		Context: ctx,
	}
	return e.consumer.CurrentRoundID(opts)
}

func (e *EthereumAPIConsumer) Fund(ethAmount *big.Float) error {
	return e.client.Fund(e.address.Hex(), ethAmount)
}

func (e *EthereumAPIConsumer) WatchPerfEvents(ctx context.Context, eventChan chan<- *PerfEvent) error {
	ethEventChan := make(chan *ethereum.APIConsumerPerfMetricsEvent)
	sub, err := e.consumer.WatchPerfMetricsEvent(&bind.WatchOpts{}, ethEventChan)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case event := <-ethEventChan:
			eventChan <- &PerfEvent{
				Contract:       e,
				RequestID:      event.RequestId,
				Round:          event.RoundID,
				BlockTimestamp: event.Timestamp,
			}
		case err := <-sub.Err():
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (e *EthereumAPIConsumer) Data(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.DefaultWallet.Address()),
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
	oracleAddr string,
	jobID [32]byte,
	payment *big.Int,
	url string,
	path string,
	times *big.Int,
) error {
	opts, err := e.client.TransactionOpts(e.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := e.consumer.CreateRequestTo(opts, common.HexToAddress(oracleAddr), jobID, payment, url, path, times)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// EthereumFluxAggregator represents the basic flux aggregation contract
type EthereumFluxAggregator struct {
	client         *client.EthereumClient
	fluxAggregator *ethereum.FluxAggregator
	address        *common.Address
}

func (f *EthereumFluxAggregator) Address() string {
	return f.address.Hex()
}

// Fund sends specified currencies to the contract
func (f *EthereumFluxAggregator) Fund(ethAmount *big.Float) error {
	return f.client.Fund(f.address.Hex(), ethAmount)
}

func (f *EthereumFluxAggregator) UpdateAvailableFunds() error {
	opts, err := f.client.TransactionOpts(f.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.UpdateAvailableFunds(opts)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFluxAggregator) PaymentAmount(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.DefaultWallet.Address()),
		Context: ctx,
	}
	payment, err := f.fluxAggregator.PaymentAmount(opts)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (f *EthereumFluxAggregator) RequestNewRound(ctx context.Context) error {
	opts, err := f.client.TransactionOpts(f.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.RequestNewRound(opts)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// WatchSubmissionReceived subscribes to any submissions on a flux feed
func (f *EthereumFluxAggregator) WatchSubmissionReceived(ctx context.Context, eventChan chan<- *SubmissionEvent) error {
	ethEventChan := make(chan *ethereum.FluxAggregatorSubmissionReceived)
	sub, err := f.fluxAggregator.WatchSubmissionReceived(&bind.WatchOpts{}, ethEventChan, nil, nil, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-ethEventChan:
			eventChan <- &SubmissionEvent{
				Contract:    event.Raw.Address,
				Submission:  event.Submission,
				Round:       event.Round,
				BlockNumber: event.Raw.BlockNumber,
				Oracle:      event.Oracle,
			}
		case err := <-sub.Err():
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (f *EthereumFluxAggregator) SetRequesterPermissions(ctx context.Context, addr common.Address, authorized bool, roundsDelay uint32) error {
	opts, err := f.client.TransactionOpts(f.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.SetRequesterPermissions(opts, addr, authorized, roundsDelay)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFluxAggregator) GetOracles(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.DefaultWallet.Address()),
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

func (f *EthereumFluxAggregator) LatestRoundID(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.DefaultWallet.Address()),
		Context: ctx,
	}
	rID, err := f.fluxAggregator.LatestRound(opts)
	if err != nil {
		return nil, err
	}
	return rID, nil
}

func (f *EthereumFluxAggregator) WithdrawPayment(
	ctx context.Context,
	from common.Address,
	to common.Address,
	amount *big.Int) error {
	opts, err := f.client.TransactionOpts(f.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := f.fluxAggregator.WithdrawPayment(opts, from, to, amount)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

func (f *EthereumFluxAggregator) WithdrawablePayment(ctx context.Context, addr common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.DefaultWallet.Address()),
		Context: ctx,
	}
	balance, err := f.fluxAggregator.WithdrawablePayment(opts, addr)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (f *EthereumFluxAggregator) LatestRoundData(ctx context.Context) (RoundData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.DefaultWallet.Address()),
		Context: ctx,
	}
	lr, err := f.fluxAggregator.LatestRoundData(opts)
	if err != nil {
		return RoundData{}, err
	}
	return lr, nil
}

// GetContractData retrieves basic data for the flux aggregator contract
func (f *EthereumFluxAggregator) GetContractData(ctx context.Context) (*FluxAggregatorData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.DefaultWallet.Address()),
		Context: ctx,
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
func (f *EthereumFluxAggregator) SetOracles(o FluxAggregatorSetOraclesOptions) error {
	opts, err := f.client.TransactionOpts(f.client.DefaultWallet)
	if err != nil {
		return err
	}

	tx, err := f.fluxAggregator.ChangeOracles(opts, o.RemoveList, o.AddList, o.AdminList, o.MinSubmissions, o.MaxSubmissions, o.RestartDelayRounds)
	if err != nil {
		return err
	}
	return f.client.ProcessTransaction(tx)
}

// Description returns the description of the flux aggregator contract
func (f *EthereumFluxAggregator) Description(ctxt context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(f.client.DefaultWallet.Address()),
		Context: ctxt,
	}
	return f.fluxAggregator.Description(opts)
}

// FluxAggregatorRoundConfirmer is a header subscription that awaits for a certain flux round to be completed
type FluxAggregatorRoundConfirmer struct {
	fluxInstance FluxAggregator
	roundID      *big.Int
	doneChan     chan struct{}
	context      context.Context
	cancel       context.CancelFunc
	done         bool
}

// NewFluxAggregatorRoundConfirmer provides a new instance of a FluxAggregatorRoundConfirmer
func NewFluxAggregatorRoundConfirmer(
	contract FluxAggregator,
	roundID *big.Int,
	timeout time.Duration,
) *FluxAggregatorRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &FluxAggregatorRoundConfirmer{
		fluxInstance: contract,
		roundID:      roundID,
		doneChan:     make(chan struct{}),
		context:      ctx,
		cancel:       ctxCancel,
	}
}

// ReceiveBlock will query the latest FluxAggregator round and check to see whether the round has confirmed
func (f *FluxAggregatorRoundConfirmer) ReceiveBlock(block client.NodeBlock) error {
	if block.Block == nil {
		return nil
	}
	if f.done {
		return nil
	}
	lr, err := f.fluxInstance.LatestRoundID(context.Background())
	if err != nil {
		return err
	}
	fluxLog := log.Debug().
		Str("Contract Address", f.fluxInstance.Address()).
		Int64("Current Round", lr.Int64()).
		Int64("Waiting for Round", f.roundID.Int64()).
		Uint64("Block Number", block.NumberU64())
	if lr.Cmp(f.roundID) >= 0 {
		fluxLog.Msg("FluxAggregator round completed")
		f.done = true
		f.doneChan <- struct{}{}
	} else {
		fluxLog.Msg("Waiting for FluxAggregator round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (f *FluxAggregatorRoundConfirmer) Wait() error {
	for {
		select {
		case <-f.doneChan:
			f.cancel()
			return nil
		case <-f.context.Done():
			return fmt.Errorf("timeout waiting for flux round to confirm: %d", f.roundID)
		}
	}
}

// VRFConsumerRoundConfirmer is a header subscription that awaits for a certain VRF round to be completed
type VRFConsumerRoundConfirmer struct {
	consumer VRFConsumer
	roundID  *big.Int
	doneChan chan struct{}
	context  context.Context
	cancel   context.CancelFunc
	done     bool
}

// NewVRFConsumerRoundConfirmer provides a new instance of a NewVRFConsumerRoundConfirmer
func NewVRFConsumerRoundConfirmer(
	contract VRFConsumer,
	roundID *big.Int,
	timeout time.Duration,
) *VRFConsumerRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &VRFConsumerRoundConfirmer{
		consumer: contract,
		roundID:  roundID,
		doneChan: make(chan struct{}),
		context:  ctx,
		cancel:   ctxCancel,
	}
}

// ReceiveBlock will query the latest VRFConsumer round and check to see whether the round has confirmed
func (f *VRFConsumerRoundConfirmer) ReceiveBlock(block client.NodeBlock) error {
	if f.done {
		return nil
	}
	roundID, err := f.consumer.CurrentRoundID(context.Background())
	if err != nil {
		return err
	}
	l := log.Debug().
		Str("Contract Address", f.consumer.Address()).
		Int64("Waiting for Round", f.roundID.Int64()).
		Int64("Current round ID", roundID.Int64()).
		Uint64("Block Number", block.NumberU64())
	if roundID.Int64() == f.roundID.Int64() {
		randomness, err := f.consumer.RandomnessOutput(context.Background())
		if err != nil {
			return err
		}
		l.Uint64("Randomness", randomness.Uint64()).
			Msg("VRFConsumer round completed")
		f.done = true
		f.doneChan <- struct{}{}
	} else {
		l.Msg("Waiting for VRFConsumer round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (f *VRFConsumerRoundConfirmer) Wait() error {
	for {
		select {
		case <-f.doneChan:
			f.cancel()
			return nil
		case <-f.context.Done():
			return fmt.Errorf("timeout waiting for VRFConsumer round to confirm: %d", f.roundID)
		}
	}
}

// EthereumLinkToken represents a LinkToken address
type EthereumLinkToken struct {
	client   *client.EthereumClient
	instance *ethereum.LinkToken
	address  common.Address
}

func (l *EthereumLinkToken) Deploy() (LinkToken, error) {
	opts, err := l.client.TransactionOpts(l.client.DefaultWallet)
	if err != nil {
		return nil, err
	}
	contractAddress, tx, contractInstance, err := ethereum.DeployLinkToken(opts, l.client.Client)
	if err != nil {
		return nil, err
	}
	if err := l.client.ProcessTransaction(tx); err != nil {
		return nil, err
	}
	log.Info().
		Str("Contract Address", contractAddress.Hex()).
		Str("Contract Name", "Link Token").
		Str("From", l.client.DefaultWallet.Address()).
		Str("Gas Cost", tx.Cost().String()).
		Str("NetworkName", l.client.NetworkConfig.Name).
		Msg("Deployed contract")
	return &EthereumLinkToken{client: l.client, instance: contractInstance, address: contractAddress}, nil
}

// Fund the LINK Token contract with ETH to distribute the token
func (l *EthereumLinkToken) Fund(ethAmount *big.Float) error {
	return l.client.Fund(l.address.Hex(), ethAmount)
}

func (l *EthereumLinkToken) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(l.client.DefaultWallet.Address()),
		Context: ctx,
	}
	balance, err := l.instance.BalanceOf(opts, common.HexToAddress(addr))
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// Name returns the name of the link token
func (l *EthereumLinkToken) Name(ctxt context.Context) (string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(l.client.DefaultWallet.Address()),
		Context: ctxt,
	}
	return l.instance.Name(opts)
}

func (l *EthereumLinkToken) Address() string {
	return l.address.Hex()
}

func (l *EthereumLinkToken) Approve(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := l.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

func (l *EthereumLinkToken) Transfer(to string, amount *big.Int) error {
	log.Info().
		Str("From", l.client.DefaultWallet.Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Msg("Transferring LINK")
	opts, err := l.client.TransactionOpts(l.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := l.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

func (l *EthereumLinkToken) TransferAndCall(to string, amount *big.Int, data []byte) error {
	opts, err := l.client.TransactionOpts(l.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := l.instance.TransferAndCall(opts, common.HexToAddress(to), amount, data)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

// EthereumOffchainAggregator represents the offchain aggregation contract
type EthereumOffchainAggregator struct {
	client  *client.EthereumClient
	ocr     *ethereum.OffchainAggregator
	address *common.Address
}

// Fund sends specified currencies to the contract
func (o *EthereumOffchainAggregator) Fund(ethAmount *big.Float) error {
	return o.client.Fund(o.address.Hex(), ethAmount)
}

// GetContractData retrieves basic data for the offchain aggregator contract
func (o *EthereumOffchainAggregator) GetContractData(ctxt context.Context) (*OffchainAggregatorData, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.client.DefaultWallet.Address()),
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
	transmitters, payees []string,
) error {
	opts, err := o.client.TransactionOpts(o.client.DefaultWallet)
	if err != nil {
		return err
	}
	transmittersAddr := make([]common.Address, 0)
	for _, tr := range transmitters {
		transmittersAddr = append(transmittersAddr, common.HexToAddress(tr))
	}
	payeesAddr := make([]common.Address, 0)
	for _, p := range payees {
		transmittersAddr = append(transmittersAddr, common.HexToAddress(p))
	}

	tx, err := o.ocr.SetPayees(opts, transmittersAddr, payeesAddr)
	if err != nil {
		return err
	}
	return o.client.ProcessTransaction(tx)
}

// SetConfig sets offchain reporting protocol configuration including participating oracles
func (o *EthereumOffchainAggregator) SetConfig(
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
	opts, err := o.client.TransactionOpts(o.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := o.ocr.SetPayees(opts, transmitters, transmitters)
	if err != nil {
		return err
	}
	if err := o.client.ProcessTransaction(tx); err != nil {
		return err
	}

	// Set Config
	opts, err = o.client.TransactionOpts(o.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err = o.ocr.SetConfig(opts, signers, transmitters, threshold, encodedConfigVersion, encodedConfig)
	if err != nil {
		return err
	}
	return o.client.ProcessTransaction(tx)
}

// RequestNewRound requests the OCR contract to create a new round
func (o *EthereumOffchainAggregator) RequestNewRound() error {
	opts, err := o.client.TransactionOpts(o.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := o.ocr.RequestNewRound(opts)
	if err != nil {
		return err
	}
	log.Info().Str("Contract Address", o.address.Hex()).Msg("New OCR round requested")

	return o.client.ProcessTransaction(tx)
}

// GetLatestAnswer returns the latest answer from the OCR contract
func (o *EthereumOffchainAggregator) GetLatestAnswer(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(o.client.DefaultWallet.Address()),
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
		From:    common.HexToAddress(o.client.DefaultWallet.Address()),
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

// RunlogRoundConfirmer is a header subscription that awaits for a certain Runlog round to be completed
type RunlogRoundConfirmer struct {
	consumer APIConsumer
	roundID  *big.Int
	doneChan chan struct{}
	context  context.Context
	cancel   context.CancelFunc
}

// NewRunlogRoundConfirmer provides a new instance of a RunlogRoundConfirmer
func NewRunlogRoundConfirmer(
	contract APIConsumer,
	roundID *big.Int,
	timeout time.Duration,
) *RunlogRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &RunlogRoundConfirmer{
		consumer: contract,
		roundID:  roundID,
		doneChan: make(chan struct{}),
		context:  ctx,
		cancel:   ctxCancel,
	}
}

// ReceiveBlock will query the latest Runlog round and check to see whether the round has confirmed
func (o *RunlogRoundConfirmer) ReceiveBlock(_ client.NodeBlock) error {
	currentRoundID, err := o.consumer.RoundID(context.Background())
	if err != nil {
		return err
	}
	ocrLog := log.Info().
		Str("Contract Address", o.consumer.Address()).
		Int64("Current Round", currentRoundID.Int64()).
		Int64("Waiting for Round", o.roundID.Int64())
	if currentRoundID.Cmp(o.roundID) >= 0 {
		ocrLog.Msg("Runlog round completed")
		o.doneChan <- struct{}{}
	} else {
		ocrLog.Msg("Waiting for Runlog round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *RunlogRoundConfirmer) Wait() error {
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for OCR round to confirm: %d", o.roundID)
		}
	}
}

// OffchainAggregatorRoundConfirmer is a header subscription that awaits for a certain OCR round to be completed
type OffchainAggregatorRoundConfirmer struct {
	ocrInstance OffchainAggregator
	roundID     *big.Int
	doneChan    chan struct{}
	context     context.Context
	cancel      context.CancelFunc
}

// NewOffchainAggregatorRoundConfirmer provides a new instance of a OffchainAggregatorRoundConfirmer
func NewOffchainAggregatorRoundConfirmer(
	contract OffchainAggregator,
	roundID *big.Int,
	timeout time.Duration,
) *OffchainAggregatorRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &OffchainAggregatorRoundConfirmer{
		ocrInstance: contract,
		roundID:     roundID,
		doneChan:    make(chan struct{}),
		context:     ctx,
		cancel:      ctxCancel,
	}
}

// ReceiveBlock will query the latest OffchainAggregator round and check to see whether the round has confirmed
func (o *OffchainAggregatorRoundConfirmer) ReceiveBlock(_ client.NodeBlock) error {
	lr, err := o.ocrInstance.GetLatestRound(context.Background())
	if err != nil {
		return err
	}
	currRound := lr.RoundId
	ocrLog := log.Info().
		Str("Contract Address", o.ocrInstance.Address()).
		Int64("Current Round", currRound.Int64()).
		Int64("Waiting for Round", o.roundID.Int64())
	if currRound.Cmp(o.roundID) >= 0 {
		ocrLog.Msg("OCR round completed")
		o.doneChan <- struct{}{}
	} else {
		ocrLog.Msg("Waiting for OCR round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *OffchainAggregatorRoundConfirmer) Wait() error {
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for OCR round to confirm: %d", o.roundID)
		}
	}
}

// KeeperConsumerRoundConfirmer is a header subscription that awaits for a round of upkeeps
type KeeperConsumerRoundConfirmer struct {
	instance     KeeperConsumer
	upkeepsValue int
	doneChan     chan struct{}
	context      context.Context
	cancel       context.CancelFunc
}

// NewKeeperConsumerRoundConfirmer provides a new instance of a KeeperConsumerRoundConfirmer
func NewKeeperConsumerRoundConfirmer(
	contract KeeperConsumer,
	counterValue int,
	timeout time.Duration,
) *KeeperConsumerRoundConfirmer {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	return &KeeperConsumerRoundConfirmer{
		instance:     contract,
		upkeepsValue: counterValue,
		doneChan:     make(chan struct{}),
		context:      ctx,
		cancel:       ctxCancel,
	}
}

// ReceiveBlock will query the latest Keeper round and check to see whether the round has confirmed
func (o *KeeperConsumerRoundConfirmer) ReceiveBlock(_ client.NodeBlock) error {
	upkeeps, err := o.instance.Counter(context.Background())
	if err != nil {
		return err
	}
	l := log.Info().
		Str("Contract Address", o.instance.Address()).
		Int64("Upkeeps", upkeeps.Int64()).
		Int("Required upkeeps", o.upkeepsValue)
	if upkeeps.Int64() == int64(o.upkeepsValue) {
		l.Msg("Upkeep completed")
		o.doneChan <- struct{}{}
	} else {
		l.Msg("Waiting for upkeep round")
	}
	return nil
}

// Wait is a blocking function that will wait until the round has confirmed, and timeout if the deadline has passed
func (o *KeeperConsumerRoundConfirmer) Wait() error {
	for {
		select {
		case <-o.doneChan:
			o.cancel()
			return nil
		case <-o.context.Done():
			return fmt.Errorf("timeout waiting for upkeeps to confirm: %d", o.upkeepsValue)
		}
	}
}

// EthereumStorage acts as a conduit for the ethereum version of the storage contract
type EthereumStorage struct {
	client *client.EthereumClient
	store  *ethereum.Store
}

// Set sets a value in the storage contract
func (e *EthereumStorage) Set(value *big.Int) error {
	opts, err := e.client.TransactionOpts(e.client.DefaultWallet)
	if err != nil {
		return err
	}

	tx, err := e.store.Set(opts, value)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// Get retrieves a set value from the storage contract
func (e *EthereumStorage) Get(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.DefaultWallet.Address()),
		Context: ctxt,
	}
	return e.store.Get(opts)
}

// EthereumVRF represents a VRF contract
type EthereumVRF struct {
	client  *client.EthereumClient
	vrf     *ethereum.VRF
	address *common.Address
}

// Fund sends specified currencies to the contract
func (v *EthereumVRF) Fund(ethAmount *big.Float) error {
	return v.client.Fund(v.address.Hex(), ethAmount)
}

// ProofLength returns the PROOFLENGTH call from the VRF contract
func (v *EthereumVRF) ProofLength(ctxt context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.DefaultWallet.Address()),
		Context: ctxt,
	}
	return v.vrf.PROOFLENGTH(opts)
}

// EthereumMockETHLINKFeed represents mocked ETH/LINK feed contract
type EthereumMockETHLINKFeed struct {
	client  *client.EthereumClient
	feed    *ethereum.MockETHLINKAggregator
	address *common.Address
}

func (v *EthereumMockETHLINKFeed) Address() string {
	return v.address.Hex()
}

// EthereumMockGASFeed represents mocked Gas feed contract
type EthereumMockGASFeed struct {
	client  *client.EthereumClient
	feed    *ethereum.MockGASAggregator
	address *common.Address
}

func (v *EthereumMockGASFeed) Address() string {
	return v.address.Hex()
}

// EthereumKeeperRegistry represents keeper registry contract
type EthereumKeeperRegistry struct {
	client   *client.EthereumClient
	registry *ethereum.KeeperRegistry
	address  *common.Address
}

func (v *EthereumKeeperRegistry) Address() string {
	return v.address.Hex()
}

func (v *EthereumKeeperRegistry) Fund(ethAmount *big.Float) error {
	return v.client.Fund(v.address.Hex(), ethAmount)
}

func (v *EthereumKeeperRegistry) SetRegistrar(registrarAddr string) error {
	opts, err := v.client.TransactionOpts(v.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := v.registry.SetRegistrar(opts, common.HexToAddress(registrarAddr))
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// AddUpkeepFunds adds link for particular upkeep id
func (v *EthereumKeeperRegistry) AddUpkeepFunds(id *big.Int, amount *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := v.registry.AddFunds(opts, id, amount)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// GetUpkeepInfo gets upkeep info
func (v *EthereumKeeperRegistry) GetUpkeepInfo(ctx context.Context, id *big.Int) (*UpkeepInfo, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.DefaultWallet.Address()),
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
		From:    common.HexToAddress(v.client.DefaultWallet.Address()),
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

func (v *EthereumKeeperRegistry) SetKeepers(keepers []string, payees []string) error {
	opts, err := v.client.TransactionOpts(v.client.DefaultWallet)
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
	return v.client.ProcessTransaction(tx)
}

// RegisterUpkeep registers contract to perform upkeep
func (v *EthereumKeeperRegistry) RegisterUpkeep(target string, gasLimit uint32, admin string, checkData []byte) error {
	opts, err := v.client.TransactionOpts(v.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := v.registry.RegisterUpkeep(opts, common.HexToAddress(target), gasLimit, common.HexToAddress(admin), checkData)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// GetKeeperList get list of all registered keeper addresses
func (v *EthereumKeeperRegistry) GetKeeperList(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.DefaultWallet.Address()),
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
	client   *client.EthereumClient
	consumer *ethereum.KeeperConsumer
	address  *common.Address
}

func (v *EthereumKeeperConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumKeeperConsumer) Fund(ethAmount *big.Float) error {
	return v.client.Fund(v.address.Hex(), ethAmount)
}

func (v *EthereumKeeperConsumer) Counter(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.DefaultWallet.Address()),
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
	client    *client.EthereumClient
	registrar *ethereum.UpkeepRegistrationRequests
	address   *common.Address
}

func (v *EthereumUpkeepRegistrationRequests) Address() string {
	return v.address.Hex()
}

// SetRegistrarConfig sets registrar config, allowing auto register or pending requests for manual registration
func (v *EthereumUpkeepRegistrationRequests) SetRegistrarConfig(
	autoRegister bool,
	windowSizeBlocks uint32,
	allowedPerWindow uint16,
	registryAddr string,
	minLinkJuels *big.Int,
) error {
	opts, err := v.client.TransactionOpts(v.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := v.registrar.SetRegistrationConfig(opts, autoRegister, windowSizeBlocks, allowedPerWindow, common.HexToAddress(registryAddr), minLinkJuels)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

func (v *EthereumUpkeepRegistrationRequests) Fund(ethAmount *big.Float) error {
	return v.client.Fund(v.address.Hex(), ethAmount)
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
}

func (v *EthereumBlockhashStore) Address() string {
	return v.address.Hex()
}

// EthereumVRFCoordinator represents VRF coordinator contract
type EthereumVRFCoordinator struct {
	address     *common.Address
	client      *client.EthereumClient
	coordinator *ethereum.VRFCoordinator
}

func (v *EthereumVRFCoordinator) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFCoordinator) HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.DefaultWallet.Address()),
		Context: ctx,
	}
	hash, err := v.coordinator.HashOfKey(opts, pubKey)
	if err != nil {
		return [32]byte{}, err
	}
	return hash, nil
}

func (v *EthereumVRFCoordinator) RegisterProvingKey(
	fee *big.Int,
	oracleAddr string,
	publicProvingKey [2]*big.Int,
	jobID [32]byte,
) error {
	opts, err := v.client.TransactionOpts(v.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := v.coordinator.RegisterProvingKey(opts, fee, common.HexToAddress(oracleAddr), publicProvingKey, jobID)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// EthereumVRFConsumer represents VRF consumer contract
type EthereumVRFConsumer struct {
	address  *common.Address
	client   *client.EthereumClient
	consumer *ethereum.VRFConsumer
}

func (v *EthereumVRFConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumVRFConsumer) Fund(ethAmount *big.Float) error {
	return v.client.Fund(v.address.Hex(), ethAmount)
}

func (v *EthereumVRFConsumer) RequestRandomness(hash [32]byte, fee *big.Int) error {
	opts, err := v.client.TransactionOpts(v.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := v.consumer.TestRequestRandomness(opts, hash, fee)
	if err != nil {
		return err
	}
	return v.client.ProcessTransaction(tx)
}

// CurrentRoundID helper roundID counter in consumer to check when all randomness requests are finished
func (v *EthereumVRFConsumer) CurrentRoundID(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.DefaultWallet.Address()),
		Context: ctx,
	}
	return v.consumer.CurrentRoundID(opts)
}

func (v *EthereumVRFConsumer) WatchPerfEvents(ctx context.Context, eventChan chan<- *PerfEvent) error {
	ethEventChan := make(chan *ethereum.VRFConsumerPerfMetricsEvent)
	sub, err := v.consumer.WatchPerfMetricsEvent(&bind.WatchOpts{}, ethEventChan)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case event := <-ethEventChan:
			eventChan <- &PerfEvent{
				Contract:       v,
				RequestID:      event.RequestId,
				Round:          event.RoundID,
				BlockTimestamp: event.Timestamp,
			}
		case err := <-sub.Err():
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (v *EthereumVRFConsumer) RandomnessOutput(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(v.client.DefaultWallet.Address()),
		Context: ctx,
	}
	out, err := v.consumer.RandomnessOutput(opts)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EthereumReadAccessController represents read access controller contract
type EthereumReadAccessController struct {
	client  *client.EthereumClient
	rac     *ethereum.SimpleReadAccessController
	address *common.Address
}

// AddAccess grants access to particular address to raise a flag
func (e *EthereumReadAccessController) AddAccess(addr string) error {
	opts, err := e.client.TransactionOpts(e.client.DefaultWallet)
	if err != nil {
		return err
	}
	log.Debug().Str("Address", addr).Msg("Adding access for address")
	tx, err := e.rac.AddAccess(opts, common.HexToAddress(addr))
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

// DisableAccessCheck disables all access checks
func (e *EthereumReadAccessController) DisableAccessCheck() error {
	opts, err := e.client.TransactionOpts(e.client.DefaultWallet)
	if err != nil {
		return err
	}
	tx, err := e.rac.DisableAccessCheck(opts)
	if err != nil {
		return err
	}
	return e.client.ProcessTransaction(tx)
}

func (e *EthereumReadAccessController) Address() string {
	return e.address.Hex()
}

// EthereumFlags represents flags contract
type EthereumFlags struct {
	client  *client.EthereumClient
	flags   *ethereum.Flags
	address *common.Address
}

func (e *EthereumFlags) Address() string {
	return e.address.Hex()
}

// GetFlag returns boolean if a flag was set for particular address
func (e *EthereumFlags) GetFlag(ctx context.Context, addr string) (bool, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(e.client.DefaultWallet.Address()),
		Context: ctx,
	}
	flag, err := e.flags.GetFlag(opts, common.HexToAddress(addr))
	if err != nil {
		return false, err
	}
	return flag, nil
}

// EthereumDeviationFlaggingValidator represents deviation flagging validator contract
type EthereumDeviationFlaggingValidator struct {
	client  *client.EthereumClient
	dfv     *ethereum.DeviationFlaggingValidator
	address *common.Address
}

func (e *EthereumDeviationFlaggingValidator) Address() string {
	return e.address.Hex()
}
