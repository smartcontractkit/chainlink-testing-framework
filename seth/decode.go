package seth

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	ErrDecodeInput          = "failed to decode transaction input"
	ErrDecodeOutput         = "failed to decode transaction output"
	ErrDecodeLog            = "failed to decode log"
	ErrDecodedLogNonIndexed = "failed to decode non-indexed log data"
	ErrDecodeILogIndexed    = "failed to decode indexed log data"
	ErrNoTxData             = "no tx data or it's less than 4 bytes"
	ErrRPCJSONCastError     = "failed to cast CallMsg error as rpc.DataError"

	WarnNoContractStore = "ContractStore is nil, use seth.NewContractStore(...) to decode transactions"
)

// DecodedTransaction decoded transaction
type DecodedTransaction struct {
	CommonData
	Index       uint                    `json:"index"`
	Hash        string                  `json:"hash,omitempty"`
	Protected   bool                    `json:"protected,omitempty"`
	Transaction *types.Transaction      `json:"transaction,omitempty"`
	Receipt     *types.Receipt          `json:"receipt,omitempty"`
	Events      []DecodedTransactionLog `json:"events,omitempty"`
}

type CommonData struct {
	CallType        string                 `json:"call_type,omitempty"`
	Signature       string                 `json:"signature"`
	Method          string                 `json:"method"`
	Input           map[string]interface{} `json:"input,omitempty"`
	Output          map[string]interface{} `json:"output,omitempty"`
	NestingLevel    int                    `json:"nesting_level,omitempty"`
	ParentSignature string                 `json:"parent_signature,omitempty"`
	Error           string                 `json:"error,omitempty"`
}

// DecodedCall decoded call
type DecodedCall struct {
	CommonData
	FromAddress string             `json:"from_address,omitempty"`
	ToAddress   string             `json:"to_address,omitempty"`
	From        string             `json:"from,omitempty"`
	To          string             `json:"to,omitempty"`
	Events      []DecodedCommonLog `json:"events,omitempty"`
	Comment     string             `json:"comment,omitempty"`
	Value       int64              `json:"value,omitempty"`
	GasLimit    uint64             `json:"gas_limit,omitempty"`
	GasUsed     uint64             `json:"gas_used,omitempty"`
}

type DecodedCommonLog struct {
	Signature string                 `json:"signature"`
	Address   common.Address         `json:"address"`
	EventData map[string]interface{} `json:"event_data"`
	Topics    []string               `json:"topics,omitempty"`
}

func getDefaultDecodedCall() *DecodedCall {
	return &DecodedCall{
		CommonData: CommonData{
			CallType:  UNKNOWN,
			Signature: UNKNOWN,
			Method:    UNKNOWN,
			Input:     make(map[string]interface{}),
			Output:    make(map[string]interface{}),
		},
		FromAddress: UNKNOWN,
		ToAddress:   UNKNOWN,
		From:        UNKNOWN,
		To:          UNKNOWN,
		Events:      make([]DecodedCommonLog, 0),
	}
}

func (d *DecodedCommonLog) MergeEventData(newEventData map[string]interface{}) {
	if d.EventData == nil {
		d.EventData = make(map[string]interface{})
	}
	for k, v := range newEventData {
		d.EventData[k] = v
	}
}

// DecodedTransactionLog decoded Solidity log(event)
type DecodedTransactionLog struct {
	DecodedCommonLog
	BlockNumber uint64 `json:"block_number"`
	Index       uint   `json:"index"`
	TXHash      string `json:"hash"`
	TXIndex     uint   `json:"tx_index"`
	Removed     bool   `json:"removed"`
	FileTag     string `json:"file_tag,omitempty"`
}

func (d *DecodedTransactionLog) MergeEventData(newEventData map[string]interface{}) {
	if d.EventData == nil {
		d.EventData = make(map[string]interface{})
	}
	for k, v := range newEventData {
		d.EventData[k] = v
	}
}

// decodeTransaction decodes inputs/outputs/logs of a transaction, if tx have logs with topics all topics are decoded
// if `tracing_enabled` flag is set in Client config it will also automatically trace all transaction calls and print the trace
func (m *Client) decodeTransaction(l zerolog.Logger, tx *types.Transaction, receipt *types.Receipt) (*DecodedTransaction, error) {
	var txInput map[string]interface{}
	var txEvents []DecodedTransactionLog
	txData := tx.Data()
	defaultTxn := &DecodedTransaction{
		Receipt:     receipt,
		Transaction: tx,
		Protected:   tx.Protected(),
		Hash:        tx.Hash().String(),
	}
	// if there is no tx data we have no inputs/outputs/logs
	if len(txData) == 0 || len(txData) < 4 {
		l.Err(errors.New(ErrNoTxData)).Send()
		return defaultTxn, nil
	}
	if m.ContractStore == nil {
		L.Warn().Msg(WarnNoContractStore)
		return defaultTxn, nil
	}

	sig := txData[:4]
	if m.ABIFinder == nil {
		L.Err(errors.New("ABIFInder is nil")).Msg("ABIFinder is required for transaction decoding")
		return defaultTxn, nil
	}

	var address string
	if tx.To() != nil {
		address = tx.To().String()
	} else {
		address = UNKNOWN
	}

	abiResult, err := m.ABIFinder.FindABIByMethod(address, sig)
	if err != nil {
		return defaultTxn, err
	}

	txInput, err = decodeTxInputs(l, txData, abiResult.Method)
	if err != nil {
		return defaultTxn, errors.Wrap(err, ErrDecodeInput)
	}

	var txIndex uint = 0

	if receipt != nil {
		l.Trace().Interface("Receipt", receipt).Msg("TX receipt")
		logsValues := make([]types.Log, 0)
		for _, l := range receipt.Logs {
			logsValues = append(logsValues, *l)
		}
		txEvents, err = m.decodeContractLogs(l, logsValues, abiResult.ABI)
		if err != nil {
			return defaultTxn, err
		}
		txIndex = receipt.TransactionIndex
	}
	ptx := &DecodedTransaction{
		CommonData: CommonData{
			Signature: common.Bytes2Hex(abiResult.Method.ID),
			Method:    abiResult.Method.Sig,
			Input:     txInput,
		},
		Index:       txIndex,
		Receipt:     receipt,
		Transaction: tx,
		Protected:   tx.Protected(),
		Hash:        tx.Hash().String(),
		Events:      txEvents,
	}

	return ptx, nil
}

// printDecodedTXData prints decoded txn data
func (m *Client) printDecodedTXData(l zerolog.Logger, ptx *DecodedTransaction) {
	l.Debug().Str("Method signature", ptx.Signature).Send()
	l.Debug().Str("Method name", ptx.Method).Send()
	if ptx.Input != nil {
		l.Debug().Interface("Inputs", ptx.Input).Send()
	}
	if ptx.Output != nil {
		l.Debug().Interface("Outputs", ptx.Output).Send()
	}
	for _, e := range ptx.Events {
		l.Debug().
			Str("Signature", e.Signature).
			Interface("Log", e.EventData).Send()
	}
}

// DecodeCustomABIErr decodes typed Solidity errors
func (m *Client) DecodeCustomABIErr(txErr error) (string, error) {
	cerr, ok := txErr.(rpc.DataError)
	if !ok {
		return "", errors.New(ErrRPCJSONCastError)
	}
	if m.ContractStore == nil {
		L.Warn().Msg(WarnNoContractStore)
		return "", nil
	}
	if cerr.ErrorData() != nil {
		L.Trace().Msg("Decoding custom ABI error from tx")
		for _, a := range m.ContractStore.ABIs {
			for k, abiError := range a.Errors {
				data, err := hex.DecodeString(cerr.ErrorData().(string)[2:])
				if err != nil {
					return "", err
				}
				if len(data) < 4 {
					return "", err
				}
				if bytes.Equal(data[:4], abiError.ID.Bytes()[:4]) {
					// Found a matching error
					v, err := abiError.Unpack(data)
					if err != nil {
						return "", err
					}
					L.Trace().Interface("Error", k).Interface("Args", v).Msg("Revert Reason")
					return fmt.Sprintf("error type: %s, error values: %v", k, v), nil
				}
			}
		}
	} else {
		L.Warn().Msg("No error data in tx")
	}
	return "", nil
}

// CallMsgFromTx creates ethereum.CallMsg from tx, used in simulated calls
func (m *Client) CallMsgFromTx(tx *types.Transaction) (ethereum.CallMsg, error) {
	signer := types.LatestSignerForChainID(tx.ChainId())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		L.Warn().Err(err).Msg("Failed to get sender from tx")
		return ethereum.CallMsg{}, err
	}

	if tx.Type() == types.LegacyTxType {
		return ethereum.CallMsg{
			From:     sender,
			To:       tx.To(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice(),
			Value:    tx.Value(),
			Data:     tx.Data(),
		}, nil
	}
	return ethereum.CallMsg{
		From:      sender,
		To:        tx.To(),
		Gas:       tx.Gas(),
		GasFeeCap: tx.GasFeeCap(),
		GasTipCap: tx.GasTipCap(),
		Value:     tx.Value(),
		Data:      tx.Data(),
	}, nil
}

func (m *Client) DownloadContractAndGetPragma(address common.Address, block *big.Int) (Pragma, error) {
	bytecode, err := m.Client.CodeAt(context.Background(), address, block)
	if err != nil {
		return Pragma{}, errors.Wrap(err, "failed to get contract code")
	}

	pragma, err := DecodePragmaVersion(common.Bytes2Hex(bytecode))
	if err != nil {
		return Pragma{}, err
	}

	return pragma, nil
}

// callAndGetRevertReason executes transaction locally and gets revert reason
func (m *Client) callAndGetRevertReason(tx *types.Transaction, rc *types.Receipt) error {
	L.Trace().Msg("Decoding revert error")
	// bind should support custom errors decoding soon, not yet merged
	// https://github.com/ethereum/go-ethereum/issues/26823
	// there are 2 types of possible errors, plain old assert/revert string
	// or new ABI encoded errors, first we try to find ABI one
	// if there is no match we print the error from CallMsg call
	msg, err := m.CallMsgFromTx(tx)
	if err != nil {
		L.Warn().Err(err).Msg("Failed to get call msg from tx. We won't be able to decode revert reason.")
		return nil
	}
	_, plainStringErr := m.Client.CallContract(context.Background(), msg, rc.BlockNumber)

	decodedABIErrString, err := m.DecodeCustomABIErr(plainStringErr)
	if err != nil {
		return err
	}
	if decodedABIErrString != "" {
		return errors.New(decodedABIErrString)
	}

	if plainStringErr != nil {
		L.Warn().Msg("Failed to decode revert reason")

		if plainStringErr.Error() == "execution reverted" && tx != nil && rc != nil {
			if tx.To() != nil {
				pragma, err := m.DownloadContractAndGetPragma(*tx.To(), rc.BlockNumber)
				if err == nil {
					if DoesPragmaSupportCustomRevert(pragma) {
						L.Warn().Str("Pragma", fmt.Sprint(pragma)).Msg("Custom revert reason is supported by pragma, but we could not decode it. This might be a bug in Seth. Please contact the Test Tooling team.")
					} else {
						L.Info().Str("Pragma", fmt.Sprint(pragma)).Msg("Custom revert reason is not supported by pragma version (must be >= 0.8.4). There's nothing more we can do to get custom revert reason.")
					}
				} else {
					L.Warn().Err(err).Msg("Failed to decode pragma version. Contract either uses very old version or was compiled without metadata. We won't be able to decode revert reason.")
				}
			} else {
				L.Warn().Msg("Transaction has no recipient address. Most likely it's a contract creation transaction. We don't support decoding revert reasons for contract creation transactions yet.")
			}
		}

		return plainStringErr
	}
	return nil
}

// decodeTxInputs decoded tx inputs
func decodeTxInputs(l zerolog.Logger, txData []byte, method *abi.Method) (map[string]interface{}, error) {
	l.Trace().Msg("Parsing tx inputs")
	if (len(txData)) < 4 {
		return nil, errors.New(ErrNoTxData)
	}

	inputMap := make(map[string]interface{})
	payload := txData[4:]
	if len(payload) == 0 || len(method.Inputs) == 0 {
		return nil, nil
	}
	err := method.Inputs.UnpackIntoMap(inputMap, payload)
	if err != nil {
		return nil, err
	}
	l.Trace().Interface("Inputs", inputMap).Msg("Transaction inputs")
	return inputMap, nil
}

// decodeTxOutputs decoded tx outputs
func decodeTxOutputs(l zerolog.Logger, payload []byte, method *abi.Method) (map[string]interface{}, error) {
	l.Trace().Msg("Parsing tx outputs")
	outputMap := make(map[string]interface{})
	// unpack method outputs
	if len(payload) == 0 {
		return nil, nil
	}
	// TODO: is it possible to have both anonymous and non-anonymous fields in solidity?
	if len(method.Outputs) > 0 && method.Outputs[0].Name == "" {
		vals, err := method.Outputs.UnpackValues(payload)
		if err != nil {
			return nil, err
		}
		for i, v := range vals {
			outputMap[strconv.Itoa(i)] = v
		}
	} else {
		err := method.Outputs.UnpackIntoMap(outputMap, payload)
		if err != nil {
			return nil, errors.Wrap(err, ErrDecodeOutput)
		}
	}
	l.Trace().Interface("Outputs", outputMap).Msg("Transaction outputs")
	return outputMap, nil
}

type DecodableLog interface {
	GetTopics() []common.Hash
	GetData() []byte
}

// decodeEventFromLog parses log body and additional topic fields
func decodeEventFromLog(
	l zerolog.Logger,
	a abi.ABI,
	eventABISpec abi.Event,
	lo DecodableLog,
) (map[string]interface{}, map[string]interface{}, error) {
	eventsMap := make(map[string]interface{})
	topicsMap := make(map[string]interface{})
	// if no data event has only indexed fields
	if len(lo.GetData()) != 0 {
		err := a.UnpackIntoMap(eventsMap, eventABISpec.Name, lo.GetData())
		if err != nil {
			return nil, nil, errors.Wrap(err, ErrDecodedLogNonIndexed)
		}
		l.Trace().Interface("Non-indexed", eventsMap).Send()
	}
	// might have up to 3 additional indexed fields
	if len(lo.GetTopics()) > 1 {
		topics := lo.GetTopics()[1:]
		var indexed []abi.Argument
		indexedTopics := make([]common.Hash, 0)
		for idx, topic := range topics {
			arg := eventABISpec.Inputs[idx]
			if arg.Indexed {
				indexed = append(indexed, arg)
				indexedTopics = append(indexedTopics, topic)
			}
		}
		l.Trace().Int("Topics", len(lo.GetTopics()[1:])).Int("Arguments", len(indexed)).Send()
		l.Trace().Interface("AllTopics", lo.GetTopics()).Send()
		l.Trace().Interface("HashOfName", eventABISpec.ID.Hex()).Send()
		l.Trace().Interface("FirstTopic", lo.GetTopics()[0]).Send()
		l.Trace().Interface("Topics", lo.GetTopics()[1:]).Send()
		l.Trace().Interface("Arguments", eventABISpec.Inputs).Send()
		l.Trace().Interface("Indexed", indexed).Send()
		err := abi.ParseTopicsIntoMap(topicsMap, indexed, indexedTopics)
		if err != nil {
			return nil, nil, errors.Wrap(err, ErrDecodeILogIndexed)
		}
		l.Trace().Interface("Indexed", topicsMap).Send()
	}
	return eventsMap, topicsMap, nil
}

type LogWithEventData interface {
	MergeEventData(map[string]interface{})
}

// decodedLogFromMaps creates DecodedLog from events and topics maps
func decodedLogFromMaps(log LogWithEventData, eventsMap map[string]interface{}, topicsMap map[string]interface{}) LogWithEventData {
	newMap := map[string]interface{}{}
	for k, v := range eventsMap {
		newMap[k] = v
	}
	for k, v := range topicsMap {
		newMap[k] = v
	}

	log.MergeEventData(newMap)

	return log
}
