package environment

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
)

var Adapter = &K8sManifest{
	Type:           "adapter",
	DeploymentFile: "templates/adapter-deployment.yml",
	ServiceFile:    "templates/adapter-service.yml",

	CallbackFunc: func(env *k8sEnvironment) error {
		environmentAdapter := &externalAdapter{}
		if services, err := env.findServicesBySelector(SelectorLabelKey, AdapterAppLabelValue); err != nil {
			return err
		} else {
			service := services[0]
			environmentAdapter.clusterURL = fmt.Sprintf(
				"http://%s:%d",
				service.Spec.ClusterIP,
				service.Spec.Ports[0].Port,
			)
			adapterPorts := env.ports[AdapterAppLabelValue]
			environmentAdapter.localURL = fmt.Sprintf(
				"http://127.0.0.1:%d",
				adapterPorts[0].Local,
			)
		}
		return nil
	},
}

var Chainlink = &K8sManifest{
	Type:           "chainlink",
	DeploymentFile: "templates/chainlink-deployment.yml",
	ServiceFile:    "templates/chainlink-service.yml",

	CallbackFunc: func(env *k8sEnvironment) error {
		for _, ports := range env.ports[ChainlinkAppLabelValue] {
			if ports.Remote == 6688 {
				cl, err := client.NewChainlink(&client.ChainlinkConfig{
					URL:      fmt.Sprintf("http://127.0.0.1:%d", ports.Local),
					Email:    "notreal@fakeemail.ch",
					Password: "twochains",
				}, http.DefaultClient)
				if err != nil {
					return err
				}
				log.Logger.Info().Str("url", fmt.Sprintf("http://127.0.0.1:%d", ports.Local)).Msg("Adding chainlink node to env")
				env.chainlinkNodes = append(env.chainlinkNodes, cl)
			}
		}
		return nil
	},
}

var Hardhat = &K8sManifest{
	Type:           "hardhat",
	DeploymentFile: "templates/hardhat-deployment.yml",
	ServiceFile:    "templates/hardhat-service.yml",

	CallbackFunc: func(env *k8sEnvironment) error {
		if _, err := env.findServicesBySelector(SelectorLabelKey, BlockchainAppLabelValue); err != nil {
			return err
		} else {
			port := env.ports[BlockchainAppLabelValue]
			env.network.SetURL(fmt.Sprintf("ws://127.0.0.1:%d", port[0].Local))
		}
		return nil
	},
}

func BasicChainlinkEnvironment(nodeCount int) K8sEnvironmentInit {
	k8sManifests := K8sManifests{
		0: Hardhat,
		1: Adapter,
	}
	for i := 0; i < nodeCount; i++ {
		k8sManifests[i+2] = Chainlink
	}
	return func() (string, K8sManifests) {
		return "basic-chainlink", k8sManifests
	}
}
