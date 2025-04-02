package crib

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

const (
	/*
	 These are constants for simulated CRIB that should never change
	 Ideally, they should be placed into CRIB repository, however, for simplicity we keep all environment connectors in CTF
	 CRIB: https://github.com/smartcontractkit/crib/tree/main/core
	 Core Chart: https://github.com/smartcontractkit/infra-charts/tree/main/chainlink-cluster
	*/
	MockserverCRIBTemplate        = "https://%s-mockserver%s"
	InternalNodeDNSTemplate       = "app-node%d"
	IngressNetworkWSURLTemplate   = "wss://%s-geth-%d-ws%s"
	IngressNetworkHTTPURLTemplate = "https://%s-geth-%d-http%s"
	// DefaultSimulatedPrivateKey is a first key used for Geth/Hardhat/Anvil
	DefaultSimulatedPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	// DefaultSimulatedNetworkChainID is a default chainID we use for Geth/Hardhat/Anvil
	DefaultSimulatedNetworkChainID = 1337
)

// ConnectionVars common K8s connection vars
type ConnectionVars struct {
	IngressSuffix string
	Namespace     string
	Network       string
	Nodes         int
}

// CoreDONConnectionConfig Chainlink DON connection config
type CoreDONConnectionConfig struct {
	*ConnectionVars
	PrivateKeys     []string
	NodeURLs        []string
	NodeInternalDNS []string
	ChainID         int64
	NetworkWSURL    string
	NetworkHTTPURL  string
	MockserverURL   string
}

// CoreDONSimulatedConnection returns all vars required to connect to core DON Simulated CRIB
// connects in CI via GAP if GAP_URL is provided
func CoreDONSimulatedConnection() (*CoreDONConnectionConfig, error) {
	vars, err := ReadCRIBVars()
	if err != nil {
		return nil, err
	}
	var conn *CoreDONConnectionConfig
	clNodeURLs := make([]string, 0)
	clNodesInternalDNS := make([]string, 0)
	for i := 1; i <= vars.Nodes; i++ {
		clNodesInternalDNS = append(clNodesInternalDNS, fmt.Sprintf(InternalNodeDNSTemplate, i))
	}
	conn = &CoreDONConnectionConfig{
		ConnectionVars:  vars,
		PrivateKeys:     []string{DefaultSimulatedPrivateKey},
		NodeURLs:        clNodeURLs,
		NodeInternalDNS: clNodesInternalDNS,
		ChainID:         DefaultSimulatedNetworkChainID,
	}
	// GAP connection
	gapEnabled := os.Getenv("GAP_ENABLED")
	if gapEnabled == "true" {
		logging.L.Info().Msg("Connecting to CRIB using GAP")
		for i := 1; i <= vars.Nodes; i++ {
			conn.NodeURLs = append(conn.NodeURLs, fmt.Sprintf("https://gap-%s-node%d.public%s", vars.Namespace, i, vars.IngressSuffix))
		}
		conn.NetworkWSURL = fmt.Sprintf("wss://gap-%s-geth-%d-ws.public%s", vars.Namespace, DefaultSimulatedNetworkChainID, vars.IngressSuffix)
		conn.NetworkHTTPURL = fmt.Sprintf("https://gap-%s-geth-%d-http.public%s", vars.Namespace, DefaultSimulatedNetworkChainID, vars.IngressSuffix)
		conn.MockserverURL = fmt.Sprintf("https://gap-%s-mockserver.public%s", vars.Namespace, vars.IngressSuffix)
	} else {
		logging.L.Info().Msg("Connecting to CRIB locally")
		for i := 1; i <= vars.Nodes; i++ {
			conn.NodeURLs = append(conn.NodeURLs, fmt.Sprintf("https://%s-node%d%s", vars.Namespace, i, vars.IngressSuffix))
		}
		conn.NetworkWSURL = fmt.Sprintf(IngressNetworkWSURLTemplate, vars.Namespace, DefaultSimulatedNetworkChainID, vars.IngressSuffix)
		conn.NetworkHTTPURL = fmt.Sprintf(IngressNetworkHTTPURLTemplate, vars.Namespace, DefaultSimulatedNetworkChainID, vars.IngressSuffix)
		conn.MockserverURL = fmt.Sprintf(MockserverCRIBTemplate, vars.Namespace, vars.IngressSuffix)

	}

	logging.L.Debug().Any("ConnectionInfo", conn).Msg("CRIB connection info")
	return conn, nil

}

// ReadCRIBVars read CRIB environment variables
func ReadCRIBVars() (*ConnectionVars, error) {
	ingressSuffix := os.Getenv("K8S_STAGING_INGRESS_SUFFIX")
	if ingressSuffix == "" {
		return nil, errors.New("K8S_STAGING_INGRESS_SUFFIX must be set to connect to k8s ingresses")
	}
	cribNamespace := os.Getenv("CRIB_NAMESPACE")
	if cribNamespace == "" {
		return nil, errors.New("CRIB_NAMESPACE must be set to connect")
	}
	cribNetwork := os.Getenv("CRIB_NETWORK")
	if cribNetwork == "" {
		return nil, errors.New("CRIB_NETWORK must be set to connect, only 'geth' is supported for now")
	}
	cribNodes := os.Getenv("CRIB_NODES")
	nodes, err := strconv.Atoi(cribNodes)
	if err != nil {
		return nil, errors.New("CRIB_NODES must be a number, 5-19 nodes")
	}
	if nodes < 2 {
		return nil, fmt.Errorf("not enough chainlink nodes, need at least 2")
	}
	return &ConnectionVars{
		IngressSuffix: ingressSuffix,
		Namespace:     cribNamespace,
		Network:       cribNetwork,
		Nodes:         nodes,
	}, nil
}
