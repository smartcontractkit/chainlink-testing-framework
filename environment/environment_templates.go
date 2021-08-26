package environment

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/tools"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AdapterAPIPort   = 6060
	ChainlinkWebPort = 6688
	ChainlinkP2PPort = 6690
	EVMRPCPort       = 8545
	ExplorerWSPort   = 4321
)

// NewAdapterManifest is the k8s manifest that when used will deploy an external adapter to an environment
func NewAdapterManifest() *K8sManifest {
	return &K8sManifest{
		id:             "adapter",
		DeploymentFile: filepath.Join(tools.ProjectRoot, "/environment/templates/adapter-deployment.yml"),
		ServiceFile:    filepath.Join(tools.ProjectRoot, "/environment/templates/adapter-service.yml"),

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
func NewChainlinkManifest() *K8sManifest {
	return &K8sManifest{
		id:             "chainlink",
		DeploymentFile: filepath.Join(tools.ProjectRoot, "/environment/templates/chainlink/chainlink-deployment.yml"),
		ServiceFile:    filepath.Join(tools.ProjectRoot, "/environment/templates/chainlink/chainlink-service.yml"),

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

// NewPostgresManifest is the k8s manifest that when used will deploy a postgres db to an environment
func NewPostgresManifest() *K8sManifest {
	return &K8sManifest{
		id:             "postgres",
		DeploymentFile: filepath.Join(tools.ProjectRoot, "/environment/templates/postgres/postgres-deployment.yml"),
		ServiceFile:    filepath.Join(tools.ProjectRoot, "/environment/templates/postgres/postgres-service.yml"),

		SetValuesFunc: func(manifest *K8sManifest) error {
			manifest.values["clusterURL"] = fmt.Sprintf(
				"postgresql://postgres:node@%s:%d",
				manifest.Service.Spec.ClusterIP,
				manifest.Service.Spec.Ports[0].Port,
			)
			manifest.values["localURL"] = fmt.Sprintf("postgresql://postgres:node@127.0.0.1:%d", manifest.ports[0].Local)
			return nil
		},
	}
}

// NewGethManifest is the k8s manifest that when used will deploy geth to an environment
func NewGethManifest() *K8sManifest {
	return &K8sManifest{
		id:             "evm",
		DeploymentFile: filepath.Join(tools.ProjectRoot, "environment/templates/geth-deployment.yml"),
		ServiceFile:    filepath.Join(tools.ProjectRoot, "environment/templates/geth-service.yml"),
		ConfigMapFile:  filepath.Join(tools.ProjectRoot, "environment/templates/geth-config-map.yml"),

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

// NewExplorerManifest is the k8s manifest that when used will deploy explorer mock service
func NewExplorerManifest() *K8sManifest {
	return &K8sManifest{
		id:             "explorer",
		DeploymentFile: filepath.Join(tools.ProjectRoot, "/environment/templates/explorer-deployment.yml"),
		ServiceFile:    filepath.Join(tools.ProjectRoot, "/environment/templates/explorer-service.yml"),
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

// NewHardhatManifest is the k8s manifest that when used will deploy hardhat to an environment
func NewHardhatManifest() *K8sManifest {
	return &K8sManifest{
		id:             "evm",
		DeploymentFile: filepath.Join(tools.ProjectRoot, "/environment/templates/hardhat-deployment.yml"),
		ServiceFile:    filepath.Join(tools.ProjectRoot, "/environment/templates/hardhat-service.yml"),
		ConfigMapFile:  filepath.Join(tools.ProjectRoot, "/environment/templates/hardhat-config-map.yml"),

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

// NewGanacheManifest is the k8s manifest that when used will deploy ganache to an environment
func NewGanacheManifest() *K8sManifest {
	return &K8sManifest{
		id:             "evm",
		DeploymentFile: filepath.Join(tools.ProjectRoot, "/environment/templates/ganache-deployment.yml"),
		ServiceFile:    filepath.Join(tools.ProjectRoot, "/environment/templates/ganache-service.yml"),

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
func NewChainlinkCluster(nodeCount int) K8sEnvSpecInit {
	// 1 adapter per set of nodes
	manifests := K8sEnvSpecs{NewAdapterManifest()}
	for i := 0; i < nodeCount; i++ {
		postgresManifest := indexedGroupManifestHelper(fmt.Sprintf("%d", i), NewPostgresManifest(), "postgresNode", "", "")
		manifests = append(manifests, postgresManifest)

		chainlinkManifest := indexedGroupManifestHelper(fmt.Sprintf("%d", i), NewChainlinkManifest(), "chainlinkCluster", "", "")
		manifests = append(manifests, chainlinkManifest)
	}

	return newChainlinkManifest(manifests, "tate-chainlink")
}

// NewMixedVersionChainlinkCluster mixes the currently latest chainlink version (as defined by the config file) with
// a number of past stable versions (defined by pastVersionsCount), ensuring that at least one of each is deployed
func NewMixedVersionChainlinkCluster(nodeCount, pastVersionsCount int) K8sEnvSpecInit {
	if nodeCount < 3 {
		log.Warn().
			Int("Provided Node Count", nodeCount).
			Int("Recommended Minimum Node Count", pastVersionsCount+1).
			Msg("You're using less than the recommended number of nodes for a mixed version deployment")
	}
	manifests := K8sEnvSpecs{NewAdapterManifest()}

	ecrImage := "public.ecr.aws/chainlink/chainlink"
	mixedImages := []string{""}
	for i := 0; i < pastVersionsCount; i++ {
		mixedImages = append(mixedImages, ecrImage)
	}

	retrievedVersions, err := getMixedVersions(pastVersionsCount)
	if err != nil {
		log.Err(err).Msg("Error retrieving versions from github")
	}
	mixedVersions := append([]string{""}, retrievedVersions...)

	for i := 0; i < nodeCount; i++ {
		postgresManifest := indexedGroupManifestHelper(fmt.Sprintf("%d", i), NewPostgresManifest(), "postgresNode", "", "")
		manifests = append(manifests, postgresManifest)

		chainlinkManifest := indexedGroupManifestHelper(fmt.Sprintf("%d", i), NewChainlinkManifest(), "chainlinkCluster", mixedImages[i%len(mixedImages)], mixedVersions[i%len(mixedVersions)])
		manifests = append(manifests, chainlinkManifest)
	}

	return newChainlinkManifest(manifests, "mixed-version-chainlink")
}

// indexedGroupManifestHelper helps build out a K8sManifestGroup for images that need multiple deploys in the same
// namespace with and index added to the identifier to separate the service selectors.
func indexedGroupManifestHelper(index string, manifest *K8sManifest, id, image, version string) *K8sManifestGroup {
	manifest.id = fmt.Sprintf("%s-%s", manifest.id, index)
	if len(image) > 0 {
		manifest.values["image"] = image
	}
	if len(version) > 0 {
		manifest.values["version"] = version
	}
	array := []*K8sManifest{manifest}
	group := &K8sManifestGroup{
		id:        id,
		manifests: array,
	}
	return group
}

// Queries github for the latest major release versions
func getMixedVersions(versionCount int) ([]string, error) {
	githubClient := github.NewClient(nil)
	releases, _, err := githubClient.Repositories.ListReleases(
		context.Background(),
		"smartcontractkit",
		"chainlink",
		&github.ListOptions{},
	)
	if err != nil {
		return []string{}, err
	}
	mixedVersions := []string{}
	for i := 0; i < versionCount; i++ {
		mixedVersions = append(mixedVersions, strings.TrimLeft(*releases[i].TagName, "v"))
	}
	return mixedVersions, nil
}

// Builds possible chainlink manifests
func newChainlinkManifest(chainlinkCluster K8sEnvSpecs, envName string) K8sEnvSpecInit {
	envWithHardhat := K8sEnvSpecs{
		NewHardhatManifest(),
	}
	envWithHardhat = append(envWithHardhat, chainlinkCluster...)
	envWithGanache := K8sEnvSpecs{
		NewGanacheManifest(),
	}
	envWithGanache = append(envWithGanache, chainlinkCluster...)
	envWithGeth := K8sEnvSpecs{
		NewGethManifest(),
		NewExplorerManifest(),
		1: chainlinkCluster,
	}
	envWithGeth = append(envWithGeth, chainlinkCluster...)
	envNoSimulatedChain := K8sEnvSpecs{}
	envNoSimulatedChain = append(envNoSimulatedChain, chainlinkCluster...)

	return func(config *config.NetworkConfig) (string, K8sEnvSpecs) {
		switch config.Name {
		case "Ethereum Geth dev":
			return envName, envWithGeth
		case "Ethereum Hardhat":
			return envName, envWithHardhat
		case "Ethereum Ganache":
			return envName, envWithGanache
		default:
			return envName, envNoSimulatedChain
		}
	}
}
