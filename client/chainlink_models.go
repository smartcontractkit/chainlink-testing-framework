package client

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// EIServiceConfig represents External Initiator service config
type EIServiceConfig struct {
	URL string
}

// ChainlinkConfig represents the variables needed to connect to a Chainlink node
type ChainlinkConfig struct {
	URL      string
	Email    string
	Password string
	RemoteIP string
}

// ResponseSlice is the generic model that can be used for all Chainlink API responses that are an slice
type ResponseSlice struct {
	Data []map[string]interface{}
}

// Response is the generic model that can be used for all Chainlink API responses
type Response struct {
	Data map[string]interface{}
}

// JobRunsResponse job runs
type JobRunsResponse struct {
	Data []RunsResponseData `json:"data"`
	Meta RunsMetaResponse   `json:"meta"`
}

// RunsResponseData runs response data
type RunsResponseData struct {
	Type       string                 `json:"type"`
	ID         string                 `json:"id"`
	Attributes RunsAttributesResponse `json:"attributes"`
}

// RunsAttributesResponse runs attributes
type RunsAttributesResponse struct {
	Meta       interface{}   `json:"meta"`
	Errors     []interface{} `json:"errors"`
	Inputs     RunInputs     `json:"inputs"`
	TaskRuns   []TaskRun     `json:"taskRuns"`
	CreatedAt  time.Time     `json:"createdAt"`
	FinishedAt time.Time     `json:"finishedAt"`
}

//DecodeLogTaskRun is "ethabidecodelog" task run info,
// also used for "RequestID" tracing in perf tests
type DecodeLogTaskRun struct {
	Fee       int    `json:"fee"`
	JobID     []int  `json:"jobID"`
	KeyHash   []int  `json:"keyHash"`
	RequestID []byte `json:"requestID"`
	Sender    string `json:"sender"`
}

//TaskRun is pipeline task run info
type TaskRun struct {
	Type       string      `json:"type"`
	CreatedAt  time.Time   `json:"createdAt"`
	FinishedAt time.Time   `json:"finishedAt"`
	Output     string      `json:"output"`
	Error      interface{} `json:"error"`
	DotID      string      `json:"dotId"`
}

// RunInputs run inputs (value)
type RunInputs struct {
	Parse int `json:"parse"`
}

// RunsMetaResponse runs meta
type RunsMetaResponse struct {
	Count int `json:"count"`
}

// BridgeType is the model that represents the bridge when read or created on a Chainlink node
type BridgeType struct {
	Data BridgeTypeData `json:"data"`
}

// BridgeTypeData is the model that represents the bridge when read or created on a Chainlink node
type BridgeTypeData struct {
	Attributes BridgeTypeAttributes `json:"attributes"`
}

// BridgeTypeAttributes is the model that represents the bridge when read or created on a Chainlink node
type BridgeTypeAttributes struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	RequestData string `json:"requestData,omitempty"`
}

// Session is the form structure used for authenticating
type Session struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// VRFKeyAttributes is the model that represents the created VRF key attributes when read
type VRFKeyAttributes struct {
	Compressed   string      `json:"compressed"`
	Uncompressed string      `json:"uncompressed"`
	Hash         string      `json:"hash"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	DeletedAt    interface{} `json:"deletedAt"`
}

// VRFKeyData is the model that represents the created VRF key's data when read
type VRFKeyData struct {
	Type       string           `json:"type"`
	ID         string           `json:"id"`
	Attributes VRFKeyAttributes `json:"attributes"`
}

// VRFKey is the model that represents the created VRF key when read
type VRFKey struct {
	Data VRFKeyData `json:"data"`
}

// VRFKeys is the model that represents the created VRF keys when read
type VRFKeys struct {
	Data []VRFKey `json:"data"`
}

// OCRKeys is the model that represents the created OCR keys when read
type OCRKeys struct {
	Data []OCRKeyData `json:"data"`
}

// OCRKey is the model that represents the created OCR keys when read
type OCRKey struct {
	Data OCRKeyData `json:"data"`
}

// OCRKeyData is the model that represents the created OCR keys when read
type OCRKeyData struct {
	Attributes OCRKeyAttributes `json:"attributes"`
	ID         string           `json:"id"`
}

// OCRKeyAttributes is the model that represents the created OCR keys when read
type OCRKeyAttributes struct {
	ConfigPublicKey       string `json:"configPublicKey"`
	OffChainPublicKey     string `json:"offChainPublicKey"`
	OnChainSigningAddress string `json:"onChainSigningAddress"`
}

// OCR2Keys is the model that represents the created OCR2 keys when read
type OCR2Keys struct {
	Data []OCR2KeyData `json:"data"`
}

// OCR2Key is the model that represents the created OCR2 keys when read
type OCR2Key struct {
	Data OCR2KeyData `json:"data"`
}

// OCR2KeyData is the model that represents the created OCR2 keys when read
type OCR2KeyData struct {
	Type       string            `json:"type"`
	Attributes OCR2KeyAttributes `json:"attributes"`
	ID         string            `json:"id"`
}

// OCR2KeyAttributes is the model that represents the created OCR2 keys when read
type OCR2KeyAttributes struct {
	ChainType         string `json:"chainType"`
	ConfigPublicKey   string `json:"configPublicKey"`
	OffChainPublicKey string `json:"offchainPublicKey"`
	OnChainPublicKey  string `json:"onchainPublicKey"`
}

// P2PKeys is the model that represents the created P2P keys when read
type P2PKeys struct {
	Data []P2PKeyData `json:"data"`
}

// P2PKey is the model that represents the created P2P keys when read
type P2PKey struct {
	Data P2PKeyData `json:"data"`
}

// P2PKeyData is the model that represents the created P2P keys when read
type P2PKeyData struct {
	Attributes P2PKeyAttributes `json:"attributes"`
}

// P2PKeyAttributes is the model that represents the created P2P keys when read
type P2PKeyAttributes struct {
	ID        int    `json:"id"`
	PeerID    string `json:"peerId"`
	PublicKey string `json:"publicKey"`
}

// ETHKeys is the model that represents the created ETH keys when read
type ETHKeys struct {
	Data []ETHKeyData `json:"data"`
}

// ETHKey is the model that represents the created ETH keys when read
type ETHKey struct {
	Data ETHKeyData `json:"data"`
}

// ETHKeyData is the model that represents the created ETH keys when read
type ETHKeyData struct {
	Attributes ETHKeyAttributes `json:"attributes"`
}

// ETHKeyAttributes is the model that represents the created ETH keys when read
type ETHKeyAttributes struct {
	Address string `json:"address"`
}

// TxKeys is the model that represents the created keys when read
type TxKeys struct {
	Data []TxKeyData `json:"data"`
}

// TxKey is the model that represents the created keys when read
type TxKey struct {
	Data TxKeyData `json:"data"`
}

// TxKeyData is the model that represents the created keys when read
type TxKeyData struct {
	Type       string          `json:"type"`
	ID         string          `json:"id"`
	Attributes TxKeyAttributes `json:"attributes"`
}

// TxKeyAttributes is the model that represents the created keys when read
type TxKeyAttributes struct {
	PublicKey string `json:"publicKey"`
}

// EIAttributes is the model that represents the EI keys when created and read
type EIAttributes struct {
	Name              string `json:"name,omitempty"`
	URL               string `json:"url,omitempty"`
	IncomingAccessKey string `json:"incomingAccessKey,omitempty"`
	AccessKey         string `json:"accessKey,omitempty"`
	Secret            string `json:"incomingSecret,omitempty"`
	OutgoingToken     string `json:"outgoingToken,omitempty"`
	OutgoingSecret    string `json:"outgoingSecret,omitempty"`
}

// EIKeys is the model that represents the EI configs when read
type EIKeys struct {
	Data []EIKey `json:"data"`
}

// EIKeyCreate is the model that represents the EI config when created
type EIKeyCreate struct {
	Data EIKey `json:"data"`
}

// EIKey is the model that represents the EI configs when read
type EIKey struct {
	Attributes EIAttributes `json:"attributes"`
}

// SpecForm is the form used when creating a v2 job spec, containing the TOML of the v2 job
type SpecForm struct {
	TOML string `json:"toml"`
}

// Spec represents a job specification that contains information about the job spec
type Spec struct {
	Data SpecData `json:"data"`
}

// SpecData contains the ID of the job spec
type SpecData struct {
	ID string `json:"id"`
}

// JobForm is the form used when creating a v2 job spec, containing the TOML of the v2 job
type JobForm struct {
	TOML string `json:"toml"`
}

// Job contains the job data for a given job
type Job struct {
	Data JobData `json:"data"`
}

// JobData contains the ID for a given job
type JobData struct {
	ID string `json:"id"`
}

// JobSpec represents the different possible job types that chainlink nodes can handle
type JobSpec interface {
	Type() string
	// String Returns TOML representation of the job
	String() (string, error)
}

// CronJobSpec represents a cron job spec
type CronJobSpec struct {
	Schedule          string `toml:"schedule"`          // CRON job style schedule string
	ObservationSource string `toml:"observationSource"` // List of commands for the chainlink node
}

// Type is cron
func (c *CronJobSpec) Type() string { return "cron" }

// String representation of the job
func (c *CronJobSpec) String() (string, error) {
	cronJobTemplateString := `type     = "cron"
schemaVersion     = 1
schedule          = "{{.Schedule}}"
observationSource = """
{{.ObservationSource}}
"""`
	return marshallTemplate(c, "CRON Job", cronJobTemplateString)
}

// PipelineSpec common API call pipeline
type PipelineSpec struct {
	BridgeTypeAttributes BridgeTypeAttributes
	DataPath             string
}

// Type is common_pipeline
func (d *PipelineSpec) Type() string {
	return "common_pipeline"
}

// String representation of the pipeline
func (d *PipelineSpec) String() (string, error) {
	sourceString := `
		fetch [type=bridge name="{{.BridgeTypeAttributes.Name}}" requestData="{{.BridgeTypeAttributes.RequestData}}"];
		parse [type=jsonparse path="{{.DataPath}}"];
		fetch -> parse;`
	return marshallTemplate(d, "API call pipeline template", sourceString)
}

// VRFTxPipelineSpec VRF request with tx callback
type VRFTxPipelineSpec struct {
	Address string
}

// Type returns the type of the pipeline
func (d *VRFTxPipelineSpec) Type() string {
	return "vrf_pipeline"
}

// String representation of the pipeline
func (d *VRFTxPipelineSpec) String() (string, error) {
	sourceString := `
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
              publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="{{.Address}}"
            data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx`
	return marshallTemplate(d, "VRF pipeline template", sourceString)
}

// DirectRequestTxPipelineSpec oracle request with tx callback
type DirectRequestTxPipelineSpec struct {
	BridgeTypeAttributes BridgeTypeAttributes
	DataPath             string
}

// Type returns the type of the pipeline
func (d *DirectRequestTxPipelineSpec) Type() string {
	return "directrequest_pipeline"
}

// String representation of the pipeline
func (d *DirectRequestTxPipelineSpec) String() (string, error) {
	sourceString := `
            decode_log   [type=ethabidecodelog
                         abi="OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)"
                         data="$(jobRun.logData)"
                         topics="$(jobRun.logTopics)"]
			encode_tx  [type=ethabiencode
                        abi="fulfill(bytes32 _requestId, uint256 _data)"
                        data=<{
                          "_requestId": $(decode_log.requestId),
                          "_data": $(parse)
                         }>
                       ]
			fetch  [type=bridge name="{{.BridgeTypeAttributes.Name}}" requestData="{{.BridgeTypeAttributes.RequestData}}"];
			parse  [type=jsonparse path="{{.DataPath}}"]
            submit [type=ethtx to="$(decode_log.requester)" data="$(encode_tx)"]
			decode_log -> fetch -> parse -> encode_tx -> submit`
	return marshallTemplate(d, "Direct request pipeline template", sourceString)
}

// DirectRequestJobSpec represents a direct request spec
type DirectRequestJobSpec struct {
	Name              string `toml:"name"`
	ContractAddress   string `toml:"contractAddress"`
	ExternalJobID     string `toml:"externalJobID"`
	ObservationSource string `toml:"observationSource"` // List of commands for the chainlink node
}

// Type returns the type of the pipeline
func (d *DirectRequestJobSpec) Type() string { return "directrequest" }

// String representation of the pipeline
func (d *DirectRequestJobSpec) String() (string, error) {
	directRequestTemplateString := `type     = "directrequest"
schemaVersion     = 1
name              = "{{.Name}}"
maxTaskDuration   = "60s"
contractAddress   = "{{.ContractAddress}}"
externalJobID     = "{{.ExternalJobID}}"
observationSource = """
{{.ObservationSource}}
"""`
	return marshallTemplate(d, "Direct Request Job", directRequestTemplateString)
}

// FluxMonitorJobSpec represents a flux monitor spec
type FluxMonitorJobSpec struct {
	Name              string        `toml:"name"`
	ContractAddress   string        `toml:"contractAddress"`   // Address of the Flux Monitor script
	Precision         int           `toml:"precision"`         // Optional
	Threshold         float32       `toml:"threshold"`         // Optional
	AbsoluteThreshold float32       `toml:"absoluteThreshold"` // Optional
	IdleTimerPeriod   time.Duration `toml:"idleTimerPeriod"`   // Optional
	IdleTimerDisabled bool          `toml:"idleTimerDisabled"` // Optional
	PollTimerPeriod   time.Duration `toml:"pollTimerPeriod"`   // Optional
	PollTimerDisabled bool          `toml:"pollTimerDisabled"` // Optional
	MaxTaskDuration   time.Duration `toml:"maxTaskDuration"`   // Optional
	ObservationSource string        `toml:"observationSource"` // List of commands for the chainlink node
}

// Type returns the type of the job
func (f *FluxMonitorJobSpec) Type() string { return "fluxmonitor" }

// String representation of the job
func (f *FluxMonitorJobSpec) String() (string, error) {
	fluxMonitorTemplateString := `type              = "fluxmonitor"
schemaVersion     = 1
name              = "{{.Name}}"
contractAddress   = "{{.ContractAddress}}"
precision         ={{if not .Precision}} 0 {{else}} {{.Precision}} {{end}}
threshold         ={{if not .Threshold}} 0.5 {{else}} {{.Threshold}} {{end}}
absoluteThreshold ={{if not .AbsoluteThreshold}} 0.1 {{else}} {{.AbsoluteThreshold}} {{end}}

idleTimerPeriod   ={{if not .IdleTimerPeriod}} "1ms" {{else}} "{{.IdleTimerPeriod}}" {{end}}
idleTimerDisabled ={{if not .IdleTimerDisabled}} false {{else}} {{.IdleTimerDisabled}} {{end}}

pollTimerPeriod   ={{if not .PollTimerPeriod}} "1m" {{else}} "{{.PollTimerPeriod}}" {{end}}
pollTimerDisabled ={{if not .PollTimerDisabled}} false {{else}} {{.PollTimerDisabled}} {{end}}

maxTaskDuration = {{if not .Precision}} "180s" {{else}} {{.Precision}} {{end}}

observationSource = """
{{.ObservationSource}}
"""`
	return marshallTemplate(f, "Flux Monitor Job", fluxMonitorTemplateString)
}

// KeeperJobSpec represents a V2 keeper spec
type KeeperJobSpec struct {
	Name                     string `toml:"name"`
	ContractAddress          string `toml:"contractAddress"`
	FromAddress              string `toml:"fromAddress"` // Hex representation of the from address
	MinIncomingConfirmations int    `toml:"minIncomingConfirmations"`
	ObservationSource        string `toml:"observationSource"`
}

// Type returns the type of the job
func (k *KeeperJobSpec) Type() string { return "keeper" }

// String representation of the job
func (k *KeeperJobSpec) String() (string, error) {
	keeperTemplateString := `
type                     = "keeper"
schemaVersion            = 3
name                     = "{{.Name}}"
contractAddress          = "{{.ContractAddress}}"
fromAddress              = "{{.FromAddress}}"
minIncomingConfirmations = {{.MinIncomingConfirmations}}

observationSource        = """
{{.ObservationSource}}
"""`
	return marshallTemplate(k, "Keeper Job", keeperTemplateString)
}

// OCRBootstrapJobSpec represents the spec for bootstrapping an OCR job, given to one node that then must be linked
// back to by others by OCRTaskJobSpecs
type OCRBootstrapJobSpec struct {
	Name                     string        `toml:"name"`
	BlockChainTimeout        time.Duration `toml:"blockchainTimeout"`                      // Optional
	ContractConfirmations    int           `toml:"contractConfigConfirmations"`            // Optional
	TrackerPollInterval      time.Duration `toml:"contractConfigTrackerPollInterval"`      // Optional
	TrackerSubscribeInterval time.Duration `toml:"contractConfigTrackerSubscribeInterval"` // Optional
	ContractAddress          string        `toml:"contractAddress"`                        // Address of the OCR contract
	IsBootstrapPeer          bool          `toml:"isBootstrapPeer"`                        // Typically true
	P2PPeerID                string        `toml:"p2pPeerID"`                              // This node's P2P ID
}

// Type returns the type of the job
func (o *OCRBootstrapJobSpec) Type() string { return "offchainreporting" }

// String representation of the job
func (o *OCRBootstrapJobSpec) String() (string, error) {
	ocrTemplateString := `type = "offchainreporting"
schemaVersion                          = 1
blockchainTimeout                      ={{if not .BlockChainTimeout}} "20s" {{else}} {{.BlockChainTimeout}} {{end}}
contractConfigConfirmations            ={{if not .ContractConfirmations}} 3 {{else}} {{.ContractConfirmations}} {{end}}
contractConfigTrackerPollInterval      ={{if not .TrackerPollInterval}} "1m" {{else}} {{.TrackerPollInterval}} {{end}}
contractConfigTrackerSubscribeInterval ={{if not .TrackerSubscribeInterval}} "2m" {{else}} {{.TrackerSubscribeInterval}} {{end}}
contractAddress                        = "{{.ContractAddress}}"
p2pBootstrapPeers                      = []
isBootstrapPeer                        = {{.IsBootstrapPeer}}
p2pPeerID                              = "{{.P2PPeerID}}"`
	return marshallTemplate(o, "OCR Bootstrap Job", ocrTemplateString)
}

// OCRTaskJobSpec represents an OCR job that is given to other nodes, meant to communicate with the bootstrap node,
// and provide their answers
type OCRTaskJobSpec struct {
	Name                     string        `toml:"name"`
	BlockChainTimeout        time.Duration `toml:"blockchainTimeout"`                      // Optional
	ContractConfirmations    int           `toml:"contractConfigConfirmations"`            // Optional
	TrackerPollInterval      time.Duration `toml:"contractConfigTrackerPollInterval"`      // Optional
	TrackerSubscribeInterval time.Duration `toml:"contractConfigTrackerSubscribeInterval"` // Optional
	ContractAddress          string        `toml:"contractAddress"`                        // Address of the OCR contract
	P2PBootstrapPeers        []Chainlink   `toml:"p2pBootstrapPeers"`                      // P2P ID of the bootstrap node
	IsBootstrapPeer          bool          `toml:"isBootstrapPeer"`                        // Typically false
	P2PPeerID                string        `toml:"p2pPeerID"`                              // This node's P2P ID
	KeyBundleID              string        `toml:"keyBundleID"`                            // ID of this node's OCR key bundle
	MonitoringEndpoint       string        `toml:"monitoringEndpoint"`                     // Typically "chain.link:4321"
	TransmitterAddress       string        `toml:"transmitterAddress"`                     // ETH address this node will use to transmit its answer
	ObservationSource        string        `toml:"observationSource"`                      // List of commands for the chainlink node
}

// P2PData holds the remote ip and the peer id and port
type P2PData struct {
	RemoteIP   string
	RemotePort string
	PeerID     string
}

// Type returns the type of the job
func (o *OCRTaskJobSpec) Type() string { return "offchainreporting" }

// String representation of the job
func (o *OCRTaskJobSpec) String() (string, error) {
	// Pre-process P2P data for easier templating
	peers := []P2PData{}
	for _, peer := range o.P2PBootstrapPeers {
		p2pKeys, err := peer.ReadP2PKeys()
		if err != nil {
			return "", err
		}
		peers = append(peers, P2PData{
			RemoteIP: peer.RemoteIP(),
			PeerID:   p2pKeys.Data[0].Attributes.PeerID,
		})
	}
	specWrap := struct {
		Name                     string
		BlockChainTimeout        time.Duration
		ContractConfirmations    int
		TrackerPollInterval      time.Duration
		TrackerSubscribeInterval time.Duration
		ContractAddress          string
		P2PBootstrapPeers        []P2PData
		IsBootstrapPeer          bool
		P2PPeerID                string
		KeyBundleID              string
		MonitoringEndpoint       string
		TransmitterAddress       string
		ObservationSource        string
	}{
		Name:                     o.Name,
		BlockChainTimeout:        o.BlockChainTimeout,
		ContractConfirmations:    o.ContractConfirmations,
		TrackerPollInterval:      o.TrackerPollInterval,
		TrackerSubscribeInterval: o.TrackerSubscribeInterval,
		ContractAddress:          o.ContractAddress,
		P2PBootstrapPeers:        peers,
		IsBootstrapPeer:          o.IsBootstrapPeer,
		P2PPeerID:                o.P2PPeerID,
		KeyBundleID:              o.KeyBundleID,
		MonitoringEndpoint:       o.MonitoringEndpoint,
		TransmitterAddress:       o.TransmitterAddress,
		ObservationSource:        o.ObservationSource,
	}
	// Results in /dns4//tcp/6690/p2p/12D3KooWAuC9xXBnadsYJpqzZZoB4rMRWqRGpxCrr2mjS7zCoAdN\
	ocrTemplateString := `type = "offchainreporting"
schemaVersion                          = 1
blockchainTimeout                      ={{if not .BlockChainTimeout}} "20s" {{else}} {{.BlockChainTimeout}} {{end}}
contractConfigConfirmations            ={{if not .ContractConfirmations}} 3 {{else}} {{.ContractConfirmations}} {{end}}
contractConfigTrackerPollInterval      ={{if not .TrackerPollInterval}} "1m" {{else}} {{.TrackerPollInterval}} {{end}}
contractConfigTrackerSubscribeInterval ={{if not .TrackerSubscribeInterval}} "2m" {{else}} {{.TrackerSubscribeInterval}} {{end}}
contractAddress                        = "{{.ContractAddress}}"
{{if .P2PBootstrapPeers}}
p2pBootstrapPeers                      = [
  {{range $peer := .P2PBootstrapPeers}}
  "/dns4/{{$peer.RemoteIP}}/tcp/6690/p2p/{{$peer.PeerID}}",
  {{end}}
]
{{else}}
p2pBootstrapPeers                      = []
{{end}}
isBootstrapPeer                        = {{.IsBootstrapPeer}}
p2pPeerID                              = "{{.P2PPeerID}}"
keyBundleID                            = "{{.KeyBundleID}}"
monitoringEndpoint                     ={{if not .MonitoringEndpoint}} "chain.link:4321" {{else}} "{{.MonitoringEndpoint}}" {{end}}
transmitterAddress                     = "{{.TransmitterAddress}}"
observationSource                      = """
{{.ObservationSource}}
"""`

	return marshallTemplate(specWrap, "OCR Job", ocrTemplateString)
}

// OCR2TaskJobSpec represents an OCR2 job that is given to other nodes, meant to communicate with the bootstrap node,
// and provide their answers
type OCR2TaskJobSpec struct {
	Name                     string            `toml:"name"`
	ContractID               string            `toml:"contractID"`                             // Address of the OCR contract/account(s)
	Relay                    string            `toml:"relay"`                                  // Name of blockchain relay to use
	RelayConfig              map[string]string `toml:"relayConfig"`                            // Relay spec object in stringified form
	P2PPeerID                string            `toml:"p2pPeerID"`                              // This node's P2P ID
	P2PBootstrapPeers        []P2PData         `toml:"p2pBootstrapPeers"`                      // P2P ID of the bootstrap node
	IsBootstrapPeer          bool              `toml:"isBootstrapPeer"`                        // Typically false
	OCRKeyBundleID           string            `toml:"ocrKeyBundleID"`                         // ID of this node's OCR key bundle
	MonitoringEndpoint       string            `toml:"monitoringEndpoint"`                     // Typically "chain.link:4321"
	TransmitterID            string            `toml:"transmitterID"`                          // ID of address this node will use to transmit
	BlockChainTimeout        time.Duration     `toml:"blockchainTimeout"`                      // Optional
	TrackerSubscribeInterval time.Duration     `toml:"contractConfigTrackerSubscribeInterval"` // Optional
	TrackerPollInterval      time.Duration     `toml:"contractConfigTrackerPollInterval"`      // Optional
	ContractConfirmations    int               `toml:"contractConfigConfirmations"`            // Optional
	ObservationSource        string            `toml:"observationSource"`                      // List of commands for the chainlink node
	JuelsPerFeeCoinSource    string            `toml:"juelsPerFeeCoinSource"`                  // List of commands to fetch JuelsPerFeeCoin value (used to calculate ocr payments)
}

// Type returns the type of the job
func (o *OCR2TaskJobSpec) Type() string { return "offchainreporting2" }

// String representation of the job
func (o *OCR2TaskJobSpec) String() (string, error) {
	ocr2TemplateString := `type = "offchainreporting2"
schemaVersion                          = 1
blockchainTimeout                      ={{if not .BlockChainTimeout}} "20s" {{else}} "{{.BlockChainTimeout}}" {{end}}
contractConfigConfirmations            ={{if not .ContractConfirmations}} 3 {{else}} {{.ContractConfirmations}} {{end}}
contractConfigTrackerPollInterval      ={{if not .TrackerPollInterval}} "1m" {{else}} "{{.TrackerPollInterval}}" {{end}}
contractConfigTrackerSubscribeInterval ={{if not .TrackerSubscribeInterval}} "2m" {{else}} "{{.TrackerSubscribeInterval}}" {{end}}
name 																	 = "{{.Name}}"
relay																	 = "{{.Relay}}"
contractID		                         = "{{.ContractID}}"
{{if .P2PBootstrapPeers}}
p2pBootstrapPeers                      = [
  {{range $peer := .P2PBootstrapPeers}}
  "{{$peer.PeerID}}@{{$peer.RemoteIP}}:{{if $peer.RemotePort}}{{$peer.RemotePort}}{{else}}6690{{end}}",
  {{end}}
]
{{else}}
p2pBootstrapPeers                      = []
{{end}}
isBootstrapPeer                        = {{.IsBootstrapPeer}}
p2pPeerID                              = "{{.P2PPeerID}}"
ocrKeyBundleID                         = "{{.OCRKeyBundleID}}"
monitoringEndpoint                     ={{if not .MonitoringEndpoint}} "chain.link:4321" {{else}} "{{.MonitoringEndpoint}}" {{end}}
transmitterID                     		 = "{{.TransmitterID}}"
{{if .IsBootstrapPeer}}
{{else}}
observationSource                      = """
{{.ObservationSource}}
"""
juelsPerFeeCoinSource                  = """
{{.JuelsPerFeeCoinSource}}
"""
{{end}}

[relayConfig]
{{range $key, $value := .RelayConfig}}
{{$key}} = "{{$value}}"
{{end}}`

	return marshallTemplate(o, "OCR2 Job", ocr2TemplateString)
}

// VRFJobSpec represents a VRF job
type VRFJobSpec struct {
	Name                     string `toml:"name"`
	CoordinatorAddress       string `toml:"coordinatorAddress"` // Address of the VRF Coordinator contract
	PublicKey                string `toml:"publicKey"`          // Public key of the proving key
	ExternalJobID            string `toml:"externalJobID"`
	ObservationSource        string `toml:"observationSource"` // List of commands for the chainlink node
	MinIncomingConfirmations int    `toml:"minIncomingConfirmations"`
}

// Type returns the type of the job
func (v *VRFJobSpec) Type() string { return "vrf" }

// String representation of the job
func (v *VRFJobSpec) String() (string, error) {
	vrfTemplateString := `
type                     = "vrf"
schemaVersion            = 1
name                     = "{{.Name}}"
coordinatorAddress       = "{{.CoordinatorAddress}}"
minIncomingConfirmations = {{.MinIncomingConfirmations}}
publicKey                = "{{.PublicKey}}"
externalJobID            = "{{.ExternalJobID}}"
observationSource = """
{{.ObservationSource}}
"""
`
	return marshallTemplate(v, "VRF Job", vrfTemplateString)
}

// WebhookJobSpec reprsents a webhook job
type WebhookJobSpec struct {
	Name              string `toml:"name"`
	Initiator         string `toml:"initiator"`         // External initiator name
	InitiatorSpec     string `toml:"initiatorSpec"`     // External initiator spec object in stringified form
	ObservationSource string `toml:"observationSource"` // List of commands for the chainlink node
}

// Type returns the type of the job
func (w *WebhookJobSpec) Type() string { return "webhook" }

// String representation of the job
func (w *WebhookJobSpec) String() (string, error) {
	webHookTemplateString := `type = "webhook"
schemaVersion      = 1
name               = "{{.Name}}"
externalInitiators = [
	{ name = "{{.Initiator}}", spec = "{{.InitiatorSpec}}"}
]
observationSource = """
{{.ObservationSource}}
"""`
	return marshallTemplate(w, "Webhook Job", webHookTemplateString)
}

// ObservationSourceSpecHTTP creates a http GET task spec for json data
func ObservationSourceSpecHTTP(url string) string {
	return fmt.Sprintf(`
		fetch [type=http method=GET url="%s"];
		parse [type=jsonparse path="data,result"];
		fetch -> parse;`, url)
}

// ObservationSourceSpecBridge creates a bridge task spec for json data
func ObservationSourceSpecBridge(bta BridgeTypeAttributes) string {
	return fmt.Sprintf(`
		fetch [type=bridge name="%s" requestData="%s"];
		parse [type=jsonparse path="data,result"];
		fetch -> parse;`, bta.Name, bta.RequestData)
}

// ObservationSourceKeeperDefault is a basic keeper default that checks and performs upkeep of the contract address
func ObservationSourceKeeperDefault() string {
	return `encode_check_upkeep_tx   [type=ethabiencode
                          abi="checkUpkeep(uint256 id, address from)"
                          data="{\\"id\\":$(jobSpec.upkeepID),\\"from\\":$(jobSpec.fromAddress)}"]
check_upkeep_tx          [type=ethcall
                          failEarly=true
                          extractRevertReason=true
                          evmChainID="$(jobSpec.evmChainID)"
                          contract="$(jobSpec.contractAddress)"
                          gas="$(jobSpec.checkUpkeepGasLimit)"
                          gasPrice="$(jobSpec.gasPrice)"
                          gasTipCap="$(jobSpec.gasTipCap)"
                          gasFeeCap="$(jobSpec.gasFeeCap)"
                          data="$(encode_check_upkeep_tx)"]
decode_check_upkeep_tx   [type=ethabidecode
                          abi="bytes memory performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth"]
encode_perform_upkeep_tx [type=ethabiencode
                          abi="performUpkeep(uint256 id, bytes calldata performData)"
                          data="{\\"id\\": $(jobSpec.upkeepID),\\"performData\\":$(decode_check_upkeep_tx.performData)}"]
perform_upkeep_tx        [type=ethtx
                          minConfirmations=0
                          to="$(jobSpec.contractAddress)"
                          from="[$(jobSpec.fromAddress)]"
                          evmChainID="$(jobSpec.evmChainID)"
                          data="$(encode_perform_upkeep_tx)"
                          gasLimit="$(jobSpec.performUpkeepGasLimit)"
                          txMeta="{\\"jobID\\":$(jobSpec.jobID)}"]
encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> perform_upkeep_tx`
}

// marshallTemplate Helper to marshall templates
func marshallTemplate(jobSpec interface{}, name, templateString string) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New(name).Parse(templateString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, jobSpec)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}
