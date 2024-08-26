package seth

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	ErrNoTrace                = "no trace found"
	ErrNoABIMethod            = "no ABI method found"
	ErrNoAbiFound             = "no ABI found in Contract Store"
	ErrNoFourByteFound        = "no method signatures found in tracing data"
	ErrInvalidMethodSignature = "no method signature found or it's not 4 bytes long"
	WrnMissingCallTrace       = "This call was missing from call trace, but it's signature was present in 4bytes trace. Most data is missing; Call order remains unknown"

	FAILED_TO_DECODE = "failed to decode"
	UNKNOWN          = "unknown"
	NO_DATA          = "no data"

	CommentMissingABI = "Call not decoded due to missing ABI instance"
)

type Tracer struct {
	Cfg                      *Config
	rpcClient                *rpc.Client
	traces                   map[string]*Trace
	Addresses                []common.Address
	ContractStore            *ContractStore
	ContractAddressToNameMap ContractMap
	decodedCalls             map[string][]*DecodedCall
	ABIFinder                *ABIFinder
	tracesMutex              *sync.RWMutex
	decodedMutex             *sync.RWMutex
}

func (t *Tracer) getTrace(txHash string) *Trace {
	t.tracesMutex.Lock()
	defer t.tracesMutex.Unlock()
	return t.traces[txHash]
}

func (t *Tracer) addTrace(txHash string, trace *Trace) {
	t.tracesMutex.Lock()
	defer t.tracesMutex.Unlock()
	t.traces[txHash] = trace
}

func (t *Tracer) GetDecodedCalls(txHash string) []*DecodedCall {
	t.decodedMutex.Lock()
	defer t.decodedMutex.Unlock()
	return t.decodedCalls[txHash]
}

func (t *Tracer) GetAllDecodedCalls() map[string][]*DecodedCall {
	t.decodedMutex.Lock()
	defer t.decodedMutex.Unlock()
	return t.decodedCalls
}

func (t *Tracer) AddDecodedCalls(txHash string, calls []*DecodedCall) {
	t.decodedMutex.Lock()
	defer t.decodedMutex.Unlock()
	t.decodedCalls[txHash] = calls
}

type Trace struct {
	TxHash       string
	FourByte     map[string]*TXFourByteMetadataOutput
	CallTrace    *TXCallTraceOutput
	OpCodesTrace map[string]interface{}
}

type TXFourByteMetadataOutput struct {
	CallSize int
	Times    int
}

type TXCallTraceOutput struct {
	Call
	Calls []Call `json:"calls"`
}

func (t *TXCallTraceOutput) AsCall() Call {
	return t.Call
}

type TraceLog struct {
	Address string   `json:"address"`
	Data    string   `json:"data"`
	Topics  []string `json:"topics"`
}

func (t TraceLog) GetTopics() []common.Hash {
	var h []common.Hash
	for _, v := range t.Topics {
		h = append(h, common.HexToHash(v))
	}
	return h
}

func (t TraceLog) GetData() []byte {
	return common.Hex2Bytes(strings.TrimPrefix(t.Data, "0x"))
}

type Call struct {
	From    string     `json:"from"`
	Gas     string     `json:"gas"`
	GasUsed string     `json:"gasUsed"`
	Input   string     `json:"input"`
	Logs    []TraceLog `json:"logs"`
	Output  string     `json:"output"`
	To      string     `json:"to"`
	Type    string     `json:"type"`
	Value   string     `json:"value"`
	Error   string     `json:"error"`
	Calls   []Call     `json:"calls"`
}

func NewTracer(cs *ContractStore, abiFinder *ABIFinder, cfg *Config, contractAddressToNameMap ContractMap, addresses []common.Address) (*Tracer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Network.DialTimeout.Duration())
	defer cancel()
	c, err := rpc.DialOptions(ctx, cfg.FirstNetworkURL(), rpc.WithHeaders(cfg.RPCHeaders))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to '%s' due to: %w", cfg.FirstNetworkURL(), err)
	}
	return &Tracer{
		Cfg:                      cfg,
		rpcClient:                c,
		traces:                   make(map[string]*Trace),
		Addresses:                addresses,
		ContractStore:            cs,
		ContractAddressToNameMap: contractAddressToNameMap,
		decodedCalls:             make(map[string][]*DecodedCall),
		ABIFinder:                abiFinder,
		tracesMutex:              &sync.RWMutex{},
		decodedMutex:             &sync.RWMutex{},
	}, nil
}

func (t *Tracer) TraceGethTX(txHash string, revertErr error) error {
	fourByte, err := t.trace4Byte(txHash)
	if err != nil {
		L.Debug().Err(err).Msg("Failed to trace 4byte signatures. Some tracing data might be missing")
	}
	opCodesTrace, err := t.traceOpCodesTracer(txHash)
	if err != nil {
		L.Debug().Err(err).Msg("Failed to trace opcodes. Some tracing data will be missing")
	}

	callTrace, err := t.traceCallTracer(txHash)
	if err != nil {
		return err
	}

	t.addTrace(txHash, &Trace{
		TxHash:       txHash,
		FourByte:     fourByte,
		CallTrace:    callTrace,
		OpCodesTrace: opCodesTrace,
	})

	decodedCalls, err := t.DecodeTrace(L, *t.getTrace(txHash))
	if err != nil {
		return err
	}

	if len(decodedCalls) != 0 {
		t.printDecodedCallData(L, decodedCalls, revertErr)

		err = t.generateDotGraph(txHash, decodedCalls, revertErr)
		if err != nil {
			return err
		}
	}

	return t.PrintTXTrace(txHash)
}

func (t *Tracer) PrintTXTrace(txHash string) error {
	trace := t.getTrace(txHash)
	if trace == nil {
		return errors.New(ErrNoTrace)
	}
	l := L.With().Str("Transaction", txHash).Logger()
	l.Trace().Interface("4Byte", trace.FourByte).Msg("Calls function signatures (names)")
	l.Trace().Interface("CallTrace", trace.CallTrace).Msg("Full call trace with logs")
	return nil
}

func (t *Tracer) trace4Byte(txHash string) (map[string]*TXFourByteMetadataOutput, error) {
	var trace map[string]int
	if err := t.rpcClient.Call(&trace, "debug_traceTransaction", txHash, map[string]interface{}{"tracer": "4byteTracer"}); err != nil {
		return nil, err
	}
	out := make(map[string]*TXFourByteMetadataOutput)
	for k, v := range trace {
		d := strings.Split(k, "-")
		callParamsSize, err := strconv.Atoi(d[1])
		if err != nil {
			return nil, err
		}
		out[d[0]] = &TXFourByteMetadataOutput{Times: v, CallSize: callParamsSize}
	}
	return out, nil
}

func (t *Tracer) traceCallTracer(txHash string) (*TXCallTraceOutput, error) {
	var trace *TXCallTraceOutput
	if err := t.rpcClient.Call(
		&trace,
		"debug_traceTransaction",
		txHash,
		map[string]interface{}{
			"tracer": "callTracer",
			"tracerConfig": map[string]interface{}{
				"withLog": true,
			},
		}); err != nil {
		return nil, err
	}
	return trace, nil
}

func (t *Tracer) traceOpCodesTracer(txHash string) (map[string]interface{}, error) {
	var trace map[string]interface{}
	if err := t.rpcClient.Call(&trace, "debug_traceTransaction", txHash); err != nil {
		return nil, err
	}
	return trace, nil
}

// DecodeTrace decodes the trace of a transaction including all subcalls. It returns a list of decoded calls.
// Depending on the config it also saves the decoded calls as JSON files.
func (t *Tracer) DecodeTrace(l zerolog.Logger, trace Trace) ([]*DecodedCall, error) {
	var decodedCalls []*DecodedCall

	if t.ContractStore == nil {
		L.Warn().Msg(WarnNoContractStore)
		return []*DecodedCall{}, nil
	}

	// we can still decode the calls without 4byte signatures
	if len(trace.FourByte) == 0 {
		L.Debug().Msg(ErrNoFourByteFound)
	}

	methods := make([]string, 0, len(trace.CallTrace.Calls)+1)

	var getSignature = func(input string) (string, error) {
		if len(input) < 10 {
			err := errors.New(ErrInvalidMethodSignature)
			l.Err(err).
				Str("Input", input).
				Send()
			return "", errors.New(ErrInvalidMethodSignature)
		}

		return input[2:10], nil
	}

	mainSig, err := getSignature(trace.CallTrace.Input)
	if err != nil {
		return nil, err
	}
	methods = append(methods, mainSig)

	var gatherAllMethodsFn func(calls []Call) error
	gatherAllMethodsFn = func(calls []Call) error {
		for _, call := range calls {
			sig, err := getSignature(call.Input)
			if err != nil {
				return err
			}

			methods = append(methods, sig)

			if len(call.Calls) > 0 {
				if err := gatherAllMethodsFn(call.Calls); err != nil {
					return err
				}
			}
		}
		return nil
	}

	err = gatherAllMethodsFn(trace.CallTrace.Calls)
	if err != nil {
		return nil, err
	}

	decodedMainCall, err := t.decodeCall(common.Hex2Bytes(methods[0]), trace.CallTrace.AsCall())
	if err != nil {
		l.Debug().
			Err(err).
			Str("From", decodedMainCall.FromAddress).
			Str("To", decodedMainCall.ToAddress).
			Msg("Failed to decode main call")

		return nil, err
	}

	decodedCalls = append(decodedCalls, decodedMainCall)

	methodCounter := 0
	nestingLevel := 1
	var processCallsFn func(calls []Call, parentSignature string) error
	processCallsFn = func(calls []Call, parentSignature string) error {
		for _, call := range calls {
			methodCounter++
			if methodCounter >= len(methods) {
				return errors.New("method counter exceeds the number of methods. This indicates there's a logical error in tracing. Please reach out to Test Tooling team")
			}

			methodHex := methods[methodCounter]
			methodByte := common.Hex2Bytes(methodHex)
			decodedSubCall, err := t.decodeCall(methodByte, call)
			if err != nil {
				l.Debug().
					Err(err).
					Str("From", call.From).
					Str("To", call.To).
					Msg("Failed to decode sub call")
				decodedCalls = append(decodedCalls, &DecodedCall{
					CommonData: CommonData{Method: FAILED_TO_DECODE,
						Input:  map[string]interface{}{"error": FAILED_TO_DECODE},
						Output: map[string]interface{}{"error": FAILED_TO_DECODE},
					},
					FromAddress: call.From,
					ToAddress:   call.To,
				})
				continue
			}
			decodedSubCall.NestingLevel = nestingLevel
			decodedSubCall.ParentSignature = parentSignature
			decodedCalls = append(decodedCalls, decodedSubCall)

			if len(call.Calls) > 0 {
				nestingLevel++
				if err := processCallsFn(call.Calls, methodHex); err != nil {
					return err
				}
				nestingLevel--
			}
		}
		return nil
	}

	err = processCallsFn(trace.CallTrace.Calls, mainSig)
	if err != nil {
		return nil, err
	}

	missingCalls := t.checkForMissingCalls(trace)
	decodedCalls = append(decodedCalls, missingCalls...)

	t.AddDecodedCalls(trace.TxHash, decodedCalls)
	return decodedCalls, nil
}

func (t *Tracer) decodeCall(byteSignature []byte, rawCall Call) (*DecodedCall, error) {
	var txInput map[string]interface{}
	var txOutput map[string]interface{}
	var txEvents []DecodedCommonLog

	var generateDuplicatesComment = func(abiResult ABIFinderResult) string {
		var comment string
		if abiResult.DuplicateCount > 0 {
			comment = fmt.Sprintf("potentially inaccurate - method present in %d other contracts", abiResult.DuplicateCount)
		}

		return comment
	}

	defaultCall := getDefaultDecodedCall()

	abiResult, err := t.ABIFinder.FindABIByMethod(rawCall.To, byteSignature)

	defaultCall.CommonData.Signature = common.Bytes2Hex(byteSignature)
	defaultCall.FromAddress = rawCall.From
	defaultCall.ToAddress = rawCall.To
	defaultCall.From = t.getHumanReadableAddressName(rawCall.From)
	defaultCall.To = t.getHumanReadableAddressName(rawCall.To) //somehow mark it with "*"
	defaultCall.Comment = generateDuplicatesComment(abiResult)

	defaultCall.CallType = rawCall.Type
	defaultCall.Error = rawCall.Error

	if rawCall.Value != "" && rawCall.Value != "0x0" {
		decimalValue, err := strconv.ParseInt(strings.TrimPrefix(rawCall.Value, "0x"), 16, 64)
		if err != nil {
			L.Debug().
				Err(err).
				Str("Value", rawCall.Value).
				Msg("Failed to parse value")
		} else {
			defaultCall.Value = decimalValue
		}
	}

	if rawCall.Gas != "" && rawCall.Gas != "0x0" {
		decimalValue, err := strconv.ParseInt(strings.TrimPrefix(rawCall.Gas, "0x"), 16, 64)
		if err != nil {
			L.Debug().
				Err(err).
				Str("Gas", rawCall.Gas).
				Msg("Failed to parse value")
		} else {
			defaultCall.GasLimit = uint64(decimalValue)
		}
	}

	if rawCall.GasUsed != "" && rawCall.GasUsed != "0x0" {
		decimalValue, err := strconv.ParseInt(strings.TrimPrefix(rawCall.GasUsed, "0x"), 16, 64)
		if err != nil {
			L.Debug().
				Err(err).
				Str("GasUsed", rawCall.GasUsed).
				Msg("Failed to parse value")
		} else {
			defaultCall.GasUsed = uint64(decimalValue)
		}
	}

	if err != nil {
		if defaultCall.Comment != "" {
			defaultCall.Comment = fmt.Sprintf("%s; %s", defaultCall.Comment, CommentMissingABI)
		} else {
			defaultCall.Comment = CommentMissingABI
		}
		L.Warn().
			Err(err).
			Str("Method signature", common.Bytes2Hex(byteSignature)).
			Str("Contract", rawCall.To).
			Msg("Method not found in any ABI instance. Unable to provide full tracing information")

		// let's not return the error, as we can still provide some information
		return defaultCall, nil
	}

	defaultCall.Method = abiResult.Method.Sig
	defaultCall.Signature = common.Bytes2Hex(abiResult.Method.ID)

	txInput, err = decodeTxInputs(L, common.Hex2Bytes(strings.TrimPrefix(rawCall.Input, "0x")), abiResult.Method)
	if err != nil {
		L.Debug().Err(err).Msg("Failed to decode inputs")
	} else {
		defaultCall.Input = txInput
	}

	if rawCall.Output != "" {
		output, err := hexutil.Decode(rawCall.Output)
		if err != nil {
			return defaultCall, errors.Wrap(err, ErrDecodeOutput)
		}
		txOutput, err = decodeTxOutputs(L, output, abiResult.Method)
		if err != nil {
			L.Debug().Err(err).Msg("Failed to decode outputs")
		} else {
			defaultCall.Output = txOutput
		}

	}

	txEvents, err = t.decodeContractLogs(L, rawCall.Logs, abiResult.ABI)
	if err != nil {
		L.Debug().Err(err).Msg("Failed to decode logs")
	} else {
		defaultCall.Events = txEvents
	}

	return defaultCall, nil
}

func (t *Tracer) isOwnAddress(addr string) bool {
	for _, a := range t.Addresses {
		if strings.ToLower(a.Hex()) == addr {
			return true
		}
	}

	return false
}

func (t *Tracer) checkForMissingCalls(trace Trace) []*DecodedCall {
	expected := 0
	for _, v := range trace.FourByte {
		expected += v.Times
	}

	var countAllTracedCallsFn func(calls []Call, previous int) int
	countAllTracedCallsFn = func(call []Call, previous int) int {
		for _, c := range call {
			previous++
			previous = countAllTracedCallsFn(c.Calls, previous)
		}

		return previous
	}

	actual := countAllTracedCallsFn(trace.CallTrace.Calls, 1) // +1 for the main call

	diff := expected - actual
	if diff != 0 {
		L.Debug().
			Int("Debugged calls", actual).
			Int("4byte signatures", len(trace.FourByte)).
			Msgf("Number of calls and signatures does not match. There were %d more call that were't debugged", diff)

		unknownCall := &DecodedCall{
			CommonData: CommonData{Method: NO_DATA,
				Input:  map[string]interface{}{"warning": NO_DATA},
				Output: map[string]interface{}{"warning": NO_DATA},
			},
			FromAddress: UNKNOWN,
			ToAddress:   UNKNOWN,
			Events: []DecodedCommonLog{
				{Signature: NO_DATA, EventData: map[string]interface{}{"warning": NO_DATA}},
			},
		}

		var missingSignatures []string
		var findSignatureFn func(fourByteSign string, calls []Call) bool
		findSignatureFn = func(fourByteSign string, calls []Call) bool {
			for _, c := range calls {
				if strings.Contains(c.Input, fourByteSign) {
					return true
				}

				if findSignatureFn(fourByteSign, c.Calls) {
					return true
				}
			}

			return false
		}
		for k := range trace.FourByte {
			if strings.Contains(trace.CallTrace.Input, k) {
				continue
			}

			found := findSignatureFn(k, trace.CallTrace.Calls)

			if !found {
				missingSignatures = append(missingSignatures, k)
			}
		}

		missedCalls := make([]*DecodedCall, 0, len(missingSignatures))

		for _, missingSig := range missingSignatures {
			byteSignature := common.Hex2Bytes(strings.TrimPrefix(missingSig, "0x"))
			humanName := missingSig

			abiResult, err := t.ABIFinder.FindABIByMethod(UNKNOWN, byteSignature)
			if err != nil {
				L.Info().
					Str("Signature", humanName).
					Msg("Method not found in any ABI instance. Unable to provide any more tracing information")

				missedCalls = append(missedCalls, unknownCall)
				continue
			}

			toAddress := t.ContractAddressToNameMap.GetContractAddress(abiResult.ContractName())
			comment := WrnMissingCallTrace
			if abiResult.DuplicateCount > 0 {
				comment = fmt.Sprintf("%s; Potentially inaccurate - method present in %d other contracts", comment, abiResult.DuplicateCount)
			}

			missedCalls = append(missedCalls, &DecodedCall{
				CommonData: CommonData{
					Signature: humanName,
					Method:    abiResult.Method.Name,
					Input:     map[string]interface{}{"warning": NO_DATA},
					Output:    map[string]interface{}{"warning": NO_DATA},
				},
				FromAddress: UNKNOWN,
				ToAddress:   toAddress,
				To:          abiResult.ContractName(),
				From:        UNKNOWN,
				Comment:     comment,
				Events: []DecodedCommonLog{
					{Signature: NO_DATA, EventData: map[string]interface{}{"warning": NO_DATA}},
				},
			})
		}

		return missedCalls
	}

	return []*DecodedCall{}
}

func (t *Tracer) SaveDecodedCallsAsJson(dirname string) error {
	for txHash, calls := range t.GetAllDecodedCalls() {
		_, err := saveAsJson(calls, dirname, txHash)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Tracer) decodeContractLogs(l zerolog.Logger, logs []TraceLog, a abi.ABI) ([]DecodedCommonLog, error) {
	l.Trace().Msg("Decoding events")
	var eventsParsed []DecodedCommonLog
	for _, lo := range logs {
		for _, evSpec := range a.Events {
			if evSpec.ID.Hex() == lo.Topics[0] {
				l.Trace().Str("Name", evSpec.RawName).Str("Signature", evSpec.Sig).Msg("Unpacking event")
				eventsMap, topicsMap, err := decodeEventFromLog(l, a, evSpec, lo)
				if err != nil {
					return nil, errors.Wrap(err, ErrDecodeLog)
				}
				parsedEvent := decodedLogFromMaps(&DecodedCommonLog{}, eventsMap, topicsMap)
				if decodedLog, ok := parsedEvent.(*DecodedCommonLog); ok {
					decodedLog.Signature = evSpec.Sig
					t.mergeLogMeta(decodedLog, lo)
					eventsParsed = append(eventsParsed, *decodedLog)
					l.Trace().Interface("Log", parsedEvent).Msg("Transaction log")
				} else {
					l.Trace().
						Str("Actual type", fmt.Sprintf("%T", decodedLog)).
						Msg("Failed to cast decoded event to DecodedCommonLog")
				}
			}
		}
	}
	return eventsParsed, nil
}

// mergeLogMeta add metadata from log
func (t *Tracer) mergeLogMeta(pe *DecodedCommonLog, l TraceLog) {
	pe.Address = common.HexToAddress(l.Address)
	pe.Topics = l.Topics
}

func (t *Tracer) getHumanReadableAddressName(address string) string {
	if t.ContractAddressToNameMap.IsKnownAddress(address) {
		address = t.ContractAddressToNameMap.GetContractName(address)
	} else if t.isOwnAddress(address) {
		address = "you"
	} else {
		address = "unknown"
	}

	return address
}

// printDecodedCallData prints decoded txn data
func (t *Tracer) printDecodedCallData(l zerolog.Logger, calls []*DecodedCall, revertErr error) {
	if !t.Cfg.hasOutput(TraceOutput_Console) {
		return
	}
	var getIndentation = func(dc *DecodedCall) string {
		var indentation string
		for i := 0; i < dc.NestingLevel; i++ {
			indentation += "  "
		}
		return indentation
	}

	L.Debug().
		Msg("----------- Decoding transaction trace started -----------")

	for i, dc := range calls {
		indentation := getIndentation(dc)

		l.Debug().Str(fmt.Sprintf("%s- Call", indentation), fmt.Sprintf("%s -> %s", dc.FromAddress, dc.ToAddress)).Send()
		l.Debug().Str(fmt.Sprintf("%s- From -> To", indentation), fmt.Sprintf("%s -> %s", dc.From, dc.To)).Send()
		l.Debug().Str(fmt.Sprintf("%s- Call Type", indentation), dc.CallType).Send()

		l.Debug().Str(fmt.Sprintf("%s- Method signature", indentation), dc.Signature).Send()
		l.Debug().Str(fmt.Sprintf("%s- Method name", indentation), dc.Method).Send()
		l.Debug().Str(fmt.Sprintf("%s- Gas used/limit", indentation), fmt.Sprintf("%d/%d", dc.GasUsed, dc.GasLimit)).Send()
		l.Debug().Str(fmt.Sprintf("%s- Gas left", indentation), fmt.Sprintf("%d", dc.GasLimit-dc.GasUsed)).Send()
		if dc.Comment != "" {
			l.Debug().Str(fmt.Sprintf("%s- Comment", indentation), dc.Comment).Send()
		}
		if dc.Input != nil {
			l.Debug().Interface(fmt.Sprintf("%s- Inputs", indentation), dc.Input).Send()
		}
		if dc.Output != nil {
			l.Debug().Interface(fmt.Sprintf("%s- Outputs", indentation), dc.Output).Send()
		}
		for _, e := range dc.Events {
			l.Debug().
				Str("Signature", e.Signature).
				Interface(fmt.Sprintf("%s- Log", indentation), e.EventData).Send()
		}

		if revertErr != nil && dc.Error != "" {
			l.Error().Str(fmt.Sprintf("%s- Revert", indentation), revertErr.Error()).Send()
		}

		if i < len(calls)-1 {
			l.Debug().Msg("")
		}
	}

	L.Debug().
		Msg("----------- Decoding transaction trace started -----------")

	if revertErr != nil {
		L.Error().Err(revertErr).Msg("Transaction reverted")
	}
}
