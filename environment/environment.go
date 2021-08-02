package environment

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/smartcontractkit/integrations-framework/client"
	"k8s.io/client-go/tools/portforward"
)

// Environment is the interface that represents a deployed environment, whether locally or on remote machines
type Environment interface {
	ID() string
	GetServices() []*ServiceDetails
	GetService(remotPort uint16) *ServiceDetails
	GetChainlinkServices() []*ServiceDetails
	GetAdapterService() *ServiceDetails

	TearDown()
}

// GetChainlinkClients will return all instantiated Chainlink clients for a given environment
func GetChainlinkClients(env Environment) ([]client.Chainlink, error) {
	var clients []client.Chainlink
	services := env.GetChainlinkServices()

	for _, service := range services {
		var localWebPort uint16
		for _, port := range service.Ports {
			if port.Remote == ChainlinkWebPort {
				localWebPort = port.Local
				break
			}
		}

		linkClient, err := client.NewChainlink(&client.ChainlinkConfig{
			URL:      fmt.Sprintf("http://127.0.0.1:%d", localWebPort),
			Email:    "notreal@fakeemail.ch",
			Password: "twochains",
			RemoteIP: service.RemoteIP,
		}, http.DefaultClient)
		if err != nil {
			return nil, err
		}
		clients = append(clients, linkClient)
	}
	return clients, nil
}

// ServiceDetails contains info on deployed services
type ServiceDetails struct {
	RemoteIP string
	Ports    []portforward.ForwardedPort
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
	adapterService := env.GetAdapterService()
	for _, port := range adapterService.Ports {
		if port.Remote == AdapterAPIPort {
			return &externalAdapter{
				localURL:   fmt.Sprintf("http://127.0.0.1:%d", port.Local),
				clusterURL: adapterService.RemoteIP,
			}, nil
		}
	}
	return nil, errors.New("No adapter found in environment")
}

// NewBlockchainClient will return an instantiated blockchain client and switch the URL depending if there's one
// deployed into the environment. If there's no deployed blockchain in the environment, the URL from the network
// config will be used
func NewBlockchainClient(env Environment, network client.BlockchainNetwork) (client.BlockchainClient, error) {
	service := env.GetService(EVMRPCPort)
	for _, port := range service.Ports {
		if port.Remote == EVMRPCPort {
			network.SetURL(fmt.Sprintf("ws://127.0.0.1:%d", service.Ports))
		}
	}
	return client.NewBlockchainClient(network)
}
