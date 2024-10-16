package clnode

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type NetworkConfig struct {
	MinIncomingConfirmations int
	MinContractPayment       string
	ChainID                  string
	EVMNodes                 []*EVMNode
}

type EVMNode struct {
	Name     string
	WsUrl    string
	HttpUrl  string
	SendOnly bool
	Order    int
}

// NewNetworkCfg generate new network configuration from blockchain.Output
// EVMNodes is used to set priority and primary/secondary for particular nodes
func NewNetworkCfg(in *NetworkConfig, out *blockchain.Output) (string, error) {
	if len(out.Nodes) != len(in.EVMNodes) {
		return "", fmt.Errorf("error configuring network, requested %d nodes, has %d blockchain outputs", len(in.EVMNodes), len(out.Nodes))
	}
	for i, n := range out.Nodes {
		in.EVMNodes[i].Name = fmt.Sprintf("node-%s-%d", uuid.NewString()[0:5], i)
		in.EVMNodes[i].WsUrl = n.WSUrl
		in.EVMNodes[i].HttpUrl = n.HTTPUrl
	}
	in.ChainID = out.ChainID
	resultCfg, err := framework.RenderTemplate(`
	[[EVM]]
	ChainID = '{{.ChainID}}'
	MinIncomingConfirmations = {{.MinIncomingConfirmations}}
	MinContractPayment = '{{.MinContractPayment}}'

	{{range .EVMNodes}}
	[[EVM.Nodes]]
	Name = '{{.Name}}'
	WsUrl = '{{.WsUrl}}'
	HttpUrl = '{{.HttpUrl}}'
	SendOnly = {{.SendOnly}}
	Order = {{.Order}}
	{{end}}
	`, in)
	if err != nil {
		return "", err
	}
	fmt.Println("Configuring networks for CL node based on blockchain outputs:")
	fmt.Println(resultCfg)
	return resultCfg, nil
}
