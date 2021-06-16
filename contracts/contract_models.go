// Package contracts handles deployment, management, and interactions of smart contracts on various chains
package contracts

import (
	"bytes"
	"context"
	"html/template"
	"math/big"
	"time"

	"github.com/smartcontractkit/integrations-framework/client"

	"github.com/ethereum/go-ethereum/common"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
)

type FluxAggregatorOptions struct {
	PaymentAmount *big.Int       // The amount of LINK paid to each oracle per submission, in wei (units of 10⁻¹⁸ LINK)
	Timeout       uint32         // The number of seconds after the previous round that are allowed to lapse before allowing an oracle to skip an unfinished round
	Validator     common.Address // An optional contract address for validating external validation of answers
	MinSubValue   *big.Int       // An immutable check for a lower bound of what submission values are accepted from an oracle
	MaxSubValue   *big.Int       // An immutable check for an upper bound of what submission values are accepted from an oracle
	Decimals      uint8          // The number of decimals to offset the answer by
	Description   string         // A short description of what is being reported
}

type FluxAggregatorData struct {
	AllocatedFunds  *big.Int         // The amount of payment yet to be withdrawn by oracles
	AvailableFunds  *big.Int         // The amount of future funding available to oracles
	LatestRoundData RoundData        // Data about the latest round
	Oracles         []common.Address // Addresses of oracles on the contract
}

type FluxAggregator interface {
	Fund(fromWallet client.BlockchainWallet, ethAmount *big.Int, linkAmount *big.Int) error
	GetContractData(ctxt context.Context) (*FluxAggregatorData, error)
	SetOracles(
		client.BlockchainWallet,
		[]common.Address,
		[]common.Address,
		[]common.Address,
		uint32,
		uint32,
		uint32,
	) error
	Description(ctxt context.Context) (string, error)
}

type LinkToken interface {
	Address() string
	Fund(fromWallet client.BlockchainWallet, ethAmount *big.Int) error
	Name(context.Context) (string, error)
}

type OffchainOptions struct {
	MaximumGasPrice           uint32         // The highest gas price for which transmitter will be compensated
	ReasonableGasPrice        uint32         // The transmitter will receive reward for gas prices under this value
	MicroLinkPerEth           uint32         // The reimbursement per ETH of gas cost, in 1e-6LINK units
	LinkGweiPerObservation    uint32         // The reward to the oracle for contributing an observation to a successfully transmitted report, in 1e-9LINK units
	LinkGweiPerTransmission   uint32         // The reward to the transmitter of a successful report, in 1e-9LINK units
	MinimumAnswer             *big.Int       // The lowest answer the median of a report is allowed to be
	MaximumAnswer             *big.Int       // The highest answer the median of a report is allowed to be
	BillingAccessController   common.Address // The access controller for billing admin functions
	RequesterAccessController common.Address // The access controller for requesting new rounds
	Decimals                  uint8          // Answers are stored in fixed-point format, with this many digits of precision
	Description               string         // A short description of what is being reported
}

// https://uploads-ssl.webflow.com/5f6b7190899f41fb70882d08/603651a1101106649eef6a53_chainlink-ocr-protocol-paper-02-24-20.pdf
type OffChainAggregatorConfig struct {
	DeltaProgress    time.Duration // The duration in which a leader must achieve progress or be replaced
	DeltaResend      time.Duration // The interval at which nodes resend NEWEPOCH messages
	DeltaRound       time.Duration // The duration after which a new round is started
	DeltaGrace       time.Duration // The duration of the grace period during which delayed oracles can still submit observations
	DeltaC           time.Duration // Limits how often updates are transmitted to the contract as long as the median isn’t changing by more then AlphaPPB
	AlphaPPB         uint64        // Allows larger changes of the median to be reported immediately, bypassing DeltaC
	DeltaStage       time.Duration // Used to stagger stages of the transmission protocol. Multiple Ethereum blocks must be mineable in this period
	RMax             uint8         // The maximum number of rounds in an epoch
	S                []int         // Transmission Schedule
	F                int           // The allowed number of "bad" oracles
	N                int           // The number of oracles
	OracleIdentities []ocrConfigHelper.OracleIdentityExtra
}

type OffChainAggregatorSpec struct {
	ContractAddress    string // Address of the OCR contract
	P2PId              string // This node's P2P ID
	BootstrapP2PId     string // The P2P ID of the bootstrap node
	KeyBundleId        string // ID of the ETH key bundle of this chainlink node
	TransmitterAddress string // Primary ETH address of this chainlink node
}

type OffChainAggregatorBootstrapSpec struct {
	ContractAddress string // Address of the OCR contract
	P2PId           string // This node's P2P ID
}

func TemplatizeOCRJobSpec(spec OffChainAggregatorSpec) (string, error) {
	ocrJobSpecTemplateString := `type = "offchainreporting"
schemaVersion = 1
contractAddress = "{{.ContractAddress}}"
p2pPeerID = "{{.P2PId}}"
p2pBootstrapPeers = [
		"/dns4/chainlink-node-1/tcp/6690/p2p/{{.BootstrapP2PId}}"  
]
isBootstrapPeer = false
keyBundleID = "{{.KeyBundleId}}"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "{{.TransmitterAddress}}"
observationTimeout = "10s"
blockchainTimeout  = "20s"
contractConfigTrackerSubscribeInterval = "2m"
contractConfigTrackerPollInterval = "1m"
contractConfigConfirmations = 3
observationSource = """
	fetch    [type=http method=POST url="http://host.docker.internal:6644/five" requestData="{}"];
	parse    [type=jsonparse path="data,result"];    
	fetch -> parse;
	"""`
	var buf bytes.Buffer
	tmpl, err := template.New("OCR Job Spec Template").Parse(ocrJobSpecTemplateString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, spec)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

func TemplatizeOCRBootsrapSpec(spec OffChainAggregatorBootstrapSpec) (string, error) {
	ocrBootstrapSpecTemplateString := `blockchainTimeout = "20s"
contractAddress = "{{.ContractAddress}}"
contractConfigConfirmations = 3
contractConfigTrackerPollInterval = "1m"
contractConfigTrackerSubscribeInterval = "2m"
isBootstrapPeer = true
p2pBootstrapPeers = []
p2pPeerID = "{{.P2PId}}"
schemaVersion = 1
type = "offchainreporting"`
	var buf bytes.Buffer
	tmpl, err := template.New("OCR Bootstrap Spec Template").Parse(ocrBootstrapSpecTemplateString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, spec)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

type OffchainAggregatorData struct {
	LatestRoundData RoundData // Data about the latest round
}

type OffchainAggregator interface {
	Address() string
	Fund(client.BlockchainWallet, *big.Int, *big.Int) error
	GetContractData(ctxt context.Context) (*OffchainAggregatorData, error)
	SetConfig(
		fromWallet client.BlockchainWallet,
		chainlinkNodes []client.Chainlink,
		ocrConfig OffChainAggregatorConfig,
	) error
	SetPayees(client.BlockchainWallet, []common.Address, []common.Address) error
	RequestNewRound(fromWallet client.BlockchainWallet) error
	Link(ctxt context.Context) (common.Address, error)
	GetLatestAnswer(ctxt context.Context) (*big.Int, error)
	GetLatestRound(ctxt context.Context) (*RoundData, error)
}

type Storage interface {
	Get(ctxt context.Context) (*big.Int, error)
	Set(*big.Int) error
}

type VRF interface {
	Fund(client.BlockchainWallet, *big.Int, *big.Int) error
	ProofLength(context.Context) (*big.Int, error)
}

type RoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
