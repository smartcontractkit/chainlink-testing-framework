package environment

import (
	"fmt"
	"github.com/smartcontractkit/integrations-framework/client"
	"net/http"
	"net/url"
)

// Environment is the interface that represents a deployed environment, whether locally or on remote machines
type Environment interface {
	GetLocalPorts(remotePort uint16) ([]uint16, error)
	GetLocalPort(remotePort uint16) (uint16, error)
	GetRemoteURL() (*url.URL, error)

	TearDown()
}

// GetChainlinkClients will return all instantiated Chainlink clients for a given environment
func GetChainlinkClients(env Environment) ([]client.Chainlink, error) {
	var clients []client.Chainlink

	ports, err := env.GetLocalPorts(ChainlinkWebPort)
	if err != nil {
		return nil, err
	}
	for _, port := range ports {
		linkClient, err := client.NewChainlink(&client.ChainlinkConfig{
			URL:      fmt.Sprintf("http://127.0.0.1:%d", port),
			Email:    "notreal@fakeemail.ch",
			Password: "twochains",
		}, http.DefaultClient)
		if err != nil {
			return nil, err
		}
		clients = append(clients, linkClient)
	}
	return clients, nil
}

// ExternalAdapter represents a dummy external adapter within the k8sEnvironment
type ExternalAdapter interface {
	LocalURL() string
	ClusterURL() string
	SetVariable(variable int) error
}

type externalAdapter struct {
	// LocalURL communicates with the dummy adapter from outside the cluster
	localURL string
	// ClusterURL communicates with the dummy adapter from within the cluster
	clusterURL string
}

// LocalURL is used for communication with the dummy adapter from outside the cluster
func (ex *externalAdapter) LocalURL() string {
	return ex.localURL
}

// ClusterURL is used for communication with the dummy adapter from within the cluster
func (ex *externalAdapter) ClusterURL() string {
	return ex.clusterURL
}

// SetVariable set the variable that's retrieved by the `/variable` call on the dummy adapter
func (ex *externalAdapter) SetVariable(variable int) error {
	_, err := http.Post(
		fmt.Sprintf("%s/set_variable?var=%d", ex.localURL, variable),
		"application/json",
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetExternalAdapter will return a deployed external adapter on an environment
func GetExternalAdapter(env Environment) (ExternalAdapter, error) {
	u, err := env.GetRemoteURL()
	if err != nil {
		return nil, err
	}
	port, err := env.GetLocalPort(AdapterAPIPort)
	if err != nil {
		return nil, err
	}
	return &externalAdapter{
		localURL:   fmt.Sprintf("http://127.0.0.1:%d", port),
		clusterURL: fmt.Sprintf("http://%s:%d", u.Host, AdapterAPIPort),
	}, nil
}

// NewBlockchainClient will return an instantiated blockchain client and switch the URL depending if there's one
// deployed into the environment. If there's no deployed blockchain in the environment, the URL from the network
// config will be used
func NewBlockchainClient(env Environment, network client.BlockchainNetwork) (client.BlockchainClient, error) {
	port, err := env.GetLocalPort(EVMRPCPort)
	if err == nil {
		network.SetURL(fmt.Sprintf("ws://127.0.0.1:%d", port))
	}
	return client.NewBlockchainClient(network)
}
