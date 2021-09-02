package environment

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/chaos"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"net/http"
	"net/url"
)

// Environment is the interface that represents a deployed environment, whether locally or on remote machines
type Environment interface {
	ID() string

	GetAllServiceDetails(remotePort uint16) ([]*ServiceDetails, error)
	GetServiceDetails(remotePort uint16) (*ServiceDetails, error)
	GetPrivateKeyFromSecret(namespace string, privateKey string) (string, error)

	WriteArtifacts(testLogFolder string)
	ApplyChaos(exp chaos.Experimentable) (string, error)
	StopChaos(name string) error
	StopAllChaos() error
	TearDown()
}

// ServiceDetails contains all of the connectivity properties about a given deployed service
type ServiceDetails struct {
	RemoteURL *url.URL
	LocalURL  *url.URL
}

// GetExplorerMockClient will return all instantiated Explorer mock client for a given environment
func GetExplorerMockClient(env Environment) (*client.ExplorerClient, error) {
	sd, err := env.GetServiceDetails(ExplorerWSPort)
	if err != nil {
		return nil, err
	}
	forwardedURL := sd.LocalURL.String()
	return client.NewExplorerMockClient(&config.ExplorerMockConfig{
		URL: forwardedURL,
	}), nil
}

// GetChainlinkClients will return all instantiated Chainlink clients for a given environment
func GetChainlinkClients(env Environment) ([]client.Chainlink, error) {
	var clients []client.Chainlink

	sd, err := env.GetAllServiceDetails(ChainlinkWebPort)
	if err != nil {
		return nil, err
	}
	for _, service := range sd {
		linkClient, err := client.NewChainlink(&client.ChainlinkConfig{
			URL:      service.LocalURL.String(),
			Email:    "notreal@fakeemail.ch",
			Password: "twochains",
			RemoteIP: service.RemoteURL.Hostname(),
		}, http.DefaultClient)
		if err != nil {
			return nil, err
		}
		clients = append(clients, linkClient)
	}
	return clients, nil
}

// ExternalAdapter represents a dummy external adapter within the K8sEnvironment
type ExternalAdapter interface {
	TriggerValueChange(i int) (int, error)
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

func (ex *externalAdapter) TriggerValueChange(i int) (int, error) {
	log.Info().Int("Iteration", i).Msg("Triggering new round")
	if i%2 == 0 {
		err := ex.SetVariable(5)
		if err != nil {
			return 0, err
		}
		return 5, nil
	} else {
		err := ex.SetVariable(6)
		if err != nil {
			return 0, err
		}
		return 6, nil
	}
}

// GetExternalAdapter will return a deployed external adapter on an environment
func GetExternalAdapter(env Environment) (ExternalAdapter, error) {
	sd, err := env.GetServiceDetails(AdapterAPIPort)
	if err != nil {
		return nil, err
	}
	return &externalAdapter{
		localURL:   sd.LocalURL.String(),
		clusterURL: sd.RemoteURL.String(),
	}, nil
}

// NewBlockchainClient will return an instantiated blockchain client and switch the URL depending if there's one
// deployed into the environment. If there's no deployed blockchain in the environment, the URL from the network
// config will be used
func NewBlockchainClient(env Environment, network client.BlockchainNetwork) (client.BlockchainClient, error) {
	sd, err := env.GetServiceDetails(EVMRPCPort)
	if err == nil {
		network.SetURL(fmt.Sprintf("ws://%s", sd.LocalURL.Host))
	}

	network.Config().PrivateKeyStore, err = NewPrivateKeyStoreFromEnv(env, network.Config())
	if err != nil {
		return nil, err
	}

	return client.NewBlockchainClient(network)
}

// NewPrivateKeyStoreFromEnv returns a keystore looking either in a cluster secret or directly from the config
func NewPrivateKeyStoreFromEnv(env Environment, network *config.NetworkConfig) (config.PrivateKeyStore, error) {
	var localKeysAndSecretKeys []string

	if network.SecretPrivateKeys {
		for _, key := range network.PrivateKeys {
			secretKey, err := env.GetPrivateKeyFromSecret(network.NamespaceForSecret, key)
			if err != nil {
				return nil, err
			}
			localKeysAndSecretKeys = append(localKeysAndSecretKeys, secretKey)
		}
	} else {
		localKeysAndSecretKeys = network.PrivateKeys
	}

	return &config.LocalStore{RawKeys: localKeysAndSecretKeys}, nil
}
