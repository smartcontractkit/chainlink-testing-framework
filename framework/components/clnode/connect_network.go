package clnode

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

// EVMNetworkConfig is CL node network configuration
type EVMNetworkConfig struct {
	MinIncomingConfirmations int
	MinContractPayment       string
	ChainID                  string
	EVMNodes                 []*EVMNode
}

// EVMNode is CL node EVM node
type EVMNode struct {
	Name     string
	WsUrl    string
	HttpUrl  string
	SendOnly bool
	Order    int
}

// NewNetworkCfg generate new network configuration from blockchain.Output
// EVMNodes is used to set priority and primary/secondary for particular nodes
func NewNetworkCfg(in *EVMNetworkConfig, out *blockchain.Output) (string, error) {
	if len(out.Nodes) != len(in.EVMNodes) {
		return "", fmt.Errorf("error configuring network, requested %d nodes, has %d blockchain outputs", len(in.EVMNodes), len(out.Nodes))
	}
	for i, n := range out.Nodes {
		in.EVMNodes[i].Name = fmt.Sprintf("node-%s-%d", uuid.NewString()[0:5], i)
		in.EVMNodes[i].WsUrl = n.InternalWSUrl
		in.EVMNodes[i].HttpUrl = n.InternalHTTPUrl
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
	return resultCfg, nil
}

// NewNetworkCfgOneNetworkAllNodes is simplified CL node network configuration where we
// add all the nodes are from the same network
func NewNetworkCfgOneNetworkAllNodes(out *blockchain.Output) (string, error) {
	resultCfg, err := framework.RenderTemplate(`
	[[EVM]]
	ChainID = '{{.ChainID}}'
	MinIncomingConfirmations = 1
	MinContractPayment = '0.0001 link'

	{{range .Nodes}}
	[[EVM.Nodes]]
	Name = 'default'
	WsUrl = '{{.InternalWSUrl}}'
	HttpUrl = '{{.InternalHTTPUrl}}'
	{{end}}
	`, out)
	if err != nil {
		return "", err
	}
	return resultCfg, nil
}
