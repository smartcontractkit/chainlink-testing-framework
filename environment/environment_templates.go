package environment

import (
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
)

const (
	AdapterAPIPort   = 6060
	ChainlinkWebPort = 6688
	ChainlinkP2PPort = 6690
	EVMRPCPort       = 8545
)

// NewAdapterManifest is the k8s manifest that when used will deploy an external adapter to an environment
func NewAdapterManifest(rootPath string) *K8sManifest {
	return &K8sManifest{
		id:             "adapter",
		DeploymentFile: prependRootPath(rootPath, "environment/templates/adapter-deployment.yml"),
		ServiceFile:    prependRootPath(rootPath, "environment/templates/adapter-service.yml"),

		values: map[string]interface{}{
			"apiPort": AdapterAPIPort,
		},

		SetValuesFunc: func(manifest *K8sManifest) error {
			environmentAdapter := &externalAdapter{}
			environmentAdapter.clusterURL = fmt.Sprintf(
				"http://%s:%d",
				manifest.Service.Spec.ClusterIP,
				manifest.Service.Spec.Ports[0].Port,
			)
			environmentAdapter.localURL = fmt.Sprintf(
				"http://127.0.0.1:%d",
				manifest.ports[0].Local,
			)
			manifest.values["clusterURL"] = environmentAdapter.clusterURL
			manifest.values["localURL"] = environmentAdapter.localURL
			return nil
		},
	}
}

// NewChainlinkManifest is the k8s manifest that when used will deploy a chainlink node to an environment
func NewChainlinkManifest(rootPath string) *K8sManifest {
	return &K8sManifest{
		id:             "chainlink",
		DeploymentFile: prependRootPath(rootPath, "environment/templates/chainlink-deployment.yml"),
		ServiceFile:    prependRootPath(rootPath, "environment/templates/chainlink-service.yml"),

		values: map[string]interface{}{
			"webPort": ChainlinkWebPort,
			"p2pPort": ChainlinkP2PPort,
		},

		Secret: &coreV1.Secret{
			ObjectMeta: v1.ObjectMeta{
				GenerateName: "chainlink-",
			},
			Type: "Opaque",
			Data: map[string][]byte{
				"apicredentials": []byte("notreal@fakeemail.ch\ntwochains"),
				"node-password":  []byte("T.tLHkcmwePT/p,]sYuntjwHKAsrhm#4eRs4LuKHwvHejWYAC2JP4M8HimwgmbaZ"),
			},
		},
	}
}

// NewHardhatManifest is the k8s manifest that when used will deploy hardhat to an environment
func NewHardhatManifest(rootPath string) *K8sManifest {
	return &K8sManifest{
		id:             "hardhat",
		DeploymentFile: prependRootPath(rootPath, "environment/templates/hardhat-deployment.yml"),
		ServiceFile:    prependRootPath(rootPath, "environment/templates/hardhat-service.yml"),

		values: map[string]interface{}{
			"rpcPort": EVMRPCPort,
		},

		SetValuesFunc: func(manifest *K8sManifest) error {
			manifest.values["clusterURL"] = fmt.Sprintf(
				"ws://%s:%d",
				manifest.Service.Spec.ClusterIP,
				manifest.Service.Spec.Ports[0].Port,
			)
			manifest.values["localURL"] = fmt.Sprintf("ws://127.0.0.1:%d", manifest.ports[0].Local)
			return nil
		},
	}
}

// NewChainlinkCluster is a basic environment that deploys hardhat with a chainlink cluster and an external adapter
func NewChainlinkCluster(rootPath string, nodeCount int) K8sEnvSpecInit {
	k8sEnvSpecs := K8sEnvSpecs{
		0: NewHardhatManifest(rootPath),
	}
	manifests := []*K8sManifest{NewAdapterManifest(rootPath)}
	for i := 0; i < nodeCount; i++ {
		manifests = append(manifests, NewChainlinkManifest(rootPath))
	}
	k8sEnvSpecs[1] = &K8sManifestGroup{
		id:        "chainlinkCluster",
		manifests: manifests,
	}

	return func() (string, K8sEnvSpecs) {
		return "basic-chainlink", k8sEnvSpecs
	}
}

func prependRootPath(rootPath, templateFile string) string {
	return filepath.Join(rootPath, templateFile)
}
