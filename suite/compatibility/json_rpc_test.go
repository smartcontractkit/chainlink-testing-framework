package compatibility

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/integrations-framework/utils"
)

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

type RPCMethodCalls map[string][]interface{}

var rpcMethodCalls RPCMethodCalls

var _ = Describe("Json RPC compatibility @json_rpc", func() {
	var (
		openrpcSchema OpenRPCStruct
		methods       []Method
	)
	rpcClientsByChain := make(map[int][]*rpc.Client)

	jsonFile, err := os.Open(filepath.Join(utils.ProjectRoot, "openrpc.json"))
	Expect(err).ShouldNot(HaveOccurred())
	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &openrpcSchema)
	Expect(err).ShouldNot(HaveOccurred())

	methods = openrpcSchema.Methods

	rpcMethodCalls = RPCMethodCalls{
		"eth_getBlockByNumber": []interface{}{"0x000000"},
		"eth_chainId":          []interface{}{},
		"eth_gasPrice":         []interface{}{},
		"eth_getBalance":       []interface{}{"0x0000000000000000000000000000000000000000"},
		"eth_getCode":          []interface{}{"0x0000000000000000000000000000000000000000"},
		"eth_getLogs": []interface{}{map[string]interface{}{
			"fromBlock": "0x000000",
			"toBlock":   "0x000010",
			"address":   "0x0000000000000000000000000000000000000000",
		}},
	}

	BeforeEach(func() {
		By("Getting RPC clients", func() {
			nc, err := config.LoadNetworksConfig(filepath.Join(utils.ProjectRoot, "networks.yaml"))
			Expect(err).ShouldNot(HaveOccurred())
			for _, networkName := range nc.SelectedNetworks {
				networkSettings, ok := nc.NetworkSettings[networkName]
				Equal(ok)
				chainId := networkSettings["chain_id"].(int)
				urls := networkSettings["urls"]
				var rpcClients []*rpc.Client
				for _, url := range urls.([]interface{}) {
					rpcClient, err := rpc.Dial(fmt.Sprintf("%v", url))
					Expect(err).ShouldNot(HaveOccurred())
					rpcClients = append(rpcClients, rpcClient)
				}
				rpcClientsByChain[chainId] = rpcClients
			}
		})
	})

	Describe("Test GET RPC methods and validate results", func() {
		It("OCR test GET Methods", func() {
			for chainId, rpcClients := range rpcClientsByChain {
				fmt.Printf("Starting tests for chain ID %d\n", chainId)
				for rpcMethod, rpcMethodParameters := range rpcMethodCalls {
					fmt.Printf("\n\nMethod: %s", rpcMethod)
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
							fmt.Printf("Error during %s RPC call: %s\n", rpcMethod, err.Error())
							break
						}
						fmt.Printf("RPC call %s result: %v\n", rpcMethod, rpcCallResult)

						if schemaLoader == nil {
							fmt.Printf("Schema loader is empty, nothing to validate")
							break
						}

						validationResult, err := gojsonschema.Validate(schemaLoader, gojsonschema.NewGoLoader(rpcCallResult))
						if err != nil {
							fmt.Printf("Error during %s validation: %s\n", rpcMethod, err.Error())
							break
						}

						if validationResult.Valid() {
							fmt.Printf("Method result schema is valid\n")
						} else {
							fmt.Printf("Method result schema is not valid. See errors :\n")
							for _, desc := range validationResult.Errors() {
								fmt.Printf("- %s\n", desc)
							}
						}
					}
				}
			}
		})
	})
})
