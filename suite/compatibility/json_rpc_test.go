package compatibility

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/utils"
	"github.com/xeipuuv/gojsonschema"
)

// openrpc types declaration
type Result struct {
	Name   string      `json:"name"`
	Schema interface{} `json:"schema"`
}

type Method struct {
	Name    string        `json:"name"`
	Summary string        `json:"summary"`
	Params  []interface{} `json:"params"`
	Result  Result        `json:"result"`
}

type OpenRPCStruct struct {
	Openrpc string      `json:"openrpc"`
	Info    interface{} `json:"info"`
	Methods []Method    `json:"methods"`
}

type RPCMethods map[string][]interface{}

// RPC methods and parameters types declaration
type GetLogs struct {
	FromBlock string `json:"fromBlock"`
	ToBlock   string `json:"toBlock"`
}

type Parameters struct {
	GetBalance            string   `json:"eth_getBalance"`
	GetBlockByNumber      string   `json:"eth_getBlockByNumber"`
	GetCode               string   `json:"eth_getCode"`
	GetLogs               GetLogs  `json:"eth_getLogs"`
	GetTransactionByHash  string   `json:"eth_getTransactionByHash"`
	GetTransactionCount   []string `json:"eth_getTransactionCount"`
	GetTransactionReceipt string   `json:"eth_getTransactionReceipt"`
}

type NetworkParameters struct {
	Name       string     `json:"name"`
	Network    string     `json:"network"`
	ChainID    int        `json:"chain_id"`
	Urls       []string   `json:"urls"`
	Parameters Parameters `json:"parameters"`
}

type NetworksParameters struct {
	SelectedNetworks   []string            `json:"selected_networks"`
	NetworksParameters []NetworkParameters `json:"networks"`
}

func getStringByIndex(values []string, index int) string {
	if values == nil || len(values) <= index {
		return ""
	}
	return values[index]
}

func getStringValue(value string, fallback string) string {
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getRPCMethods(parameters Parameters) RPCMethods {
	return RPCMethods{
		"eth_chainId":          []interface{}{},
		"eth_getBalance":       []interface{}{getStringValue(parameters.GetBalance, "0x0000000000000000000000000000000000000000")},
		"eth_getBlockByNumber": []interface{}{getStringValue(parameters.GetBlockByNumber, "0x333333"), false},
		"eth_getCode":          []interface{}{getStringValue(parameters.GetCode, "0x0000000000000000000000000000000000000000")},
		"eth_getLogs": []interface{}{map[string]interface{}{
			"fromBlock": getStringValue(parameters.GetLogs.FromBlock, "0x444444"),
			"toBlock":   getStringValue(parameters.GetLogs.ToBlock, "0x444444"),
		}},
		"eth_gasPrice":             []interface{}{},
		"eth_getTransactionByHash": []interface{}{getStringValue(parameters.GetTransactionByHash, "0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0")},
		"eth_getTransactionCount": []interface{}{
			getStringValue(getStringByIndex(parameters.GetTransactionCount, 0), "0x0000000000000000000000000000000000000000"),
			getStringValue(getStringByIndex(parameters.GetTransactionCount, 1), "0x444444"),
		},
		"eth_getTransactionReceipt": []interface{}{getStringValue(parameters.GetTransactionReceipt, "0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0")},
	}
}

func Subscription(rpcClient *rpc.Client, chainId int) (block *types.Header, error error) {
	headerChannel := make(chan *types.Header)
	subscription, err := rpcClient.EthSubscribe(context.Background(), headerChannel, "newHeads")

	defer subscription.Unsubscribe()

	if err != nil {
		log.Error().
			Int("ChainID", chainId).
			Msgf("Could not subscribe to new block headers: %s", err.Error())
	}

	log.Info().
		Int("ChainID", chainId).
		Msg("Subscribed to new block headers")

	select {
	case err := <-subscription.Err():
		return nil, err
	case header := <-headerChannel:
		return header, nil
	case <-time.After(120 * time.Second):
		return nil, fmt.Errorf("timed out waiting for block header")
	}
}

var _ = Describe("JSON RPC compatibility @json_rpc", func() {
	var (
		networksParameters NetworksParameters
		selectedNetworks   []string
		openrpcSchema      OpenRPCStruct
		methods            []Method
	)

	rpcClientsByNetwork := make(map[string][]*rpc.Client)
	parametersByNetwork := make(map[string]NetworkParameters)

	// read openrpc JSON schema
	openrpcJSON, err := os.Open(filepath.Join(utils.TestSuiteRoot, "compatibility", "openrpc.json"))
	Expect(err).ShouldNot(HaveOccurred())
	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(openrpcJSON)

	openrpcJSONBytes, _ := ioutil.ReadAll(openrpcJSON)
	err = json.Unmarshal(openrpcJSONBytes, &openrpcSchema)
	Expect(err).ShouldNot(HaveOccurred())

	methods = openrpcSchema.Methods

	// read PRC methods parameters
	rpcMethodsParametersJSON, err := os.Open(filepath.Join(utils.TestSuiteRoot, "compatibility", "values.json"))
	Expect(err).ShouldNot(HaveOccurred())
	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(rpcMethodsParametersJSON)

	rpcMethodsParametersJSONBytes, _ := ioutil.ReadAll(rpcMethodsParametersJSON)
	err = json.Unmarshal(rpcMethodsParametersJSONBytes, &networksParameters)
	Expect(err).ShouldNot(HaveOccurred())

	selectedNetworks = networksParameters.SelectedNetworks
	for _, networkParameters := range networksParameters.NetworksParameters {
		parametersByNetwork[networkParameters.Name] = networkParameters
	}

	BeforeEach(func() {
		By("Getting RPC clients", func() {
			for _, networkName := range selectedNetworks {
				var rpcClients []*rpc.Client

				urls := parametersByNetwork[networkName].Urls
				Expect(urls).ShouldNot(BeNil(), "Network web socket urls should be set: %v", networkName)

				for _, url := range urls {
					rpcClient, err := rpc.Dial(fmt.Sprintf("%v", url))
					Expect(err).ShouldNot(HaveOccurred())

					rpcClients = append(rpcClients, rpcClient)
				}

				rpcClientsByNetwork[networkName] = rpcClients
			}
		})
	})

	Describe("Test JSON RPC GET-methods and validate results", func() {
		It("OCR test GET Methods", func() {
			for networkName, rpcClients := range rpcClientsByNetwork {
				network := parametersByNetwork[networkName].Network
				chainId := parametersByNetwork[networkName].ChainID

				log.Info().
					Str("Network", network).
					Int("ChainID", chainId).
					Msg("Starting JSON RPC compatibility test")

				for rpcMethod, rpcMethodParameters := range getRPCMethods(parametersByNetwork[networkName].Parameters) {
					log.Info().
						Str("Network", network).
						Int("ChainID", chainId).
						Str("Method", rpcMethod).
						Msg("Testing RPC method call")

					var method Method
					for _, value := range methods {
						if value.Name == rpcMethod {
							method = value
							break
						}
					}

					schemaLoader := gojsonschema.NewGoLoader(method.Result.Schema)
					for _, rpcClient := range rpcClients {
						var rpcCallResult interface{}
						err := rpcClient.CallContext(context.Background(), &rpcCallResult, rpcMethod, rpcMethodParameters...)
						if err != nil {
							log.Error().
								Str("Network", network).
								Int("ChainID", chainId).
								Str("Method", rpcMethod).
								Msgf("Error while calling RPC method: %s", err.Error())
							break
						}

						log.Info().
							Str("Network", network).
							Int("ChainID", chainId).
							Str("Method", rpcMethod).
							Msgf("RPC call result: %v", rpcCallResult)

						if schemaLoader.JsonSource() == nil {
							log.Info().
								Str("Network", network).
								Int("ChainID", chainId).
								Str("Method", rpcMethod).
								Msg("Schema loader is empty, nothing to validate")
							break
						}

						validationResult, err := gojsonschema.Validate(schemaLoader, gojsonschema.NewGoLoader(rpcCallResult))
						if err != nil {
							log.Error().
								Str("Network", network).
								Int("ChainID", chainId).
								Str("Method", rpcMethod).
								Msgf("Error during RPC call result schema validation: %s", err.Error())
							break
						}

						if validationResult.Valid() {
							log.Info().
								Str("Network", network).
								Int("ChainID", chainId).
								Str("Method", rpcMethod).
								Msg("RPC call result schema is valid")
						} else {
							log.Error().
								Str("Network", network).
								Int("ChainID", chainId).
								Str("Method", rpcMethod).
								Msg("RPC call result schema is not valid. See errors:")
							for _, desc := range validationResult.Errors() {
								log.Error().
									Msgf("- %s", desc)
							}
						}
					}
				}
			}
		})
	})

	Describe("Test JSON RPC subscription", func() {
		It("Subscribe to block headers", func() {
			for networkName, rpcClients := range rpcClientsByNetwork {
				network := parametersByNetwork[networkName].Network
				chainId := parametersByNetwork[networkName].ChainID

				log.Info().
					Str("Network", network).
					Int("ChainID", chainId).
					Msg("Starting block header subscription compatibility test")

				for _, rpcClient := range rpcClients {
					header, err := Subscription(rpcClient, chainId)
					if err != nil {
						log.Error().
							Str("Network", network).
							Int("ChainID", chainId).
							Msgf("Error occurred while processing subscription: %s", err.Error())
						continue
					}

					log.Info().
						Str("Network", network).
						Int("ChainID", chainId).
						Msgf("New block header successfully received: block number %v, hash %v", header.Number, header.Hash())
				}
			}
		})
	})
})
