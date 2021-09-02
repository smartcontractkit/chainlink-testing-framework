package environment

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
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
	ExplorerAPIPort  = 8080
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
			manifest.values["clusterURL"] = fmt.Sprintf(
				"http://%s:%d",
				manifest.Service.Spec.ClusterIP,
				manifest.Service.Spec.Ports[0].Port,
			)
			manifest.values["localURL"] = fmt.Sprintf(
				"http://127.0.0.1:%d",
				manifest.ports[0].Local,
			)
			return nil
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

func NewExplorerManifest(nodeCount int) *K8sManifest {
	return &K8sManifest{
		id:             "explorer",
		DeploymentFile: filepath.Join(tools.ProjectRoot, "/environment/templates/explorer-deployment.yml"),
		ServiceFile:    filepath.Join(tools.ProjectRoot, "/environment/templates/explorer-service.yml"),
		SetValuesFunc: func(manifest *K8sManifest) error {
			manifest.values["clusterURL"] = fmt.Sprintf(
				"ws://%s:8080",
				manifest.Service.Spec.ClusterIP,
			)
			manifest.values["localURL"] = "https://127.0.0.1:8080"
			podsFullNames, err := manifest.GetPodsFullNames("explorer")
			if err != nil {
				return err
			}
			_, _, err = manifest.ExecuteInPod(podsFullNames[0], "explorer", []string{"yarn", "--cwd", "apps/explorer", "admin:seed", "username", "password"})
			if err != nil {
				return err
			}

			accessKeys := TemplateValuesArray{}
			secretKeys := TemplateValuesArray{}

			explorerClient, err := GetExplorerClient(manifest.getServiceDetails)
			if err != nil {
				return err
			}
			for i := 0; i < nodeCount; i++ {
				credentials, err := explorerClient.PostAdminNodes(fmt.Sprintf("node-%d", i))
				if err != nil {
					return err
				}
				accessKeys.Values = append(accessKeys.Values, credentials.AccessKey)
				secretKeys.Values = append(secretKeys.Values, credentials.Secret)

			}
			manifest.values["accessKeys"] = &accessKeys
			manifest.values["secretKeys"] = &secretKeys

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

// NewMixedVersionChainlinkCluster mixes the currently latest chainlink version (as defined by the config file) with
// a number of past stable versions (defined by pastVersionsCount), ensuring that at least one of each is deployed
func NewMixedVersionChainlinkCluster(nodeCount, pastVersionsCount int) K8sEnvSpecInit {
	if nodeCount < 3 {
		log.Warn().
			Int("Provided Node Count", nodeCount).
			Int("Recommended Minimum Node Count", pastVersionsCount+1).
			Msg("You're using less than the recommended number of nodes for a mixed version deployment")
	}

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

	chainlinkGroup := &K8sManifestGroup{
		id:        "chainlinkCluster",
		manifests: []*K8sManifest{},
	}
	for i := 0; i < nodeCount; i++ {
		cManifest := NewChainlinkManifest()
		cManifest.id = fmt.Sprintf("%s-%d", cManifest.id, i)
		cManifest.values["image"] = mixedImages[i%len(mixedImages)]
		cManifest.values["version"] = mixedVersions[i%len(mixedVersions)]
		chainlinkGroup.manifests = append(chainlinkGroup.manifests, cManifest)
	}

	return addDependencyGroup(nodeCount, "mixed-version-chainlink", chainlinkGroup)
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

// NewChainlinkCluster is a basic environment that deploys hardhat with a chainlink cluster and an external adapter
func NewChainlinkCluster(nodeCount int) K8sEnvSpecInit {
	chainlinkGroup := &K8sManifestGroup{
		id:        "chainlinkCluster",
		manifests: []*K8sManifest{},
	}
	for i := 0; i < nodeCount; i++ {
		cManifest := NewChainlinkManifest()
		cManifest.id = fmt.Sprintf("%s-%d", cManifest.id, i)
		chainlinkGroup.manifests = append(chainlinkGroup.manifests, cManifest)
	}

	return addDependencyGroup(nodeCount, "basic-chainlink", chainlinkGroup)
}

// addDependencyGroup add everything that has no dependencies but other pods have
// dependencies on in the first group
func addDependencyGroup(postgresCount int, envName string, chainlinkGroup *K8sManifestGroup) K8sEnvSpecInit {
	group := &K8sManifestGroup{
		id:        "DependencyGroup",
		manifests: []*K8sManifest{NewAdapterManifest()},

		SetValuesFunc: func(mg *K8sManifestGroup) error {
			postgresURLs := TemplateValuesArray{}

			for _, manifest := range mg.manifests {
				if strings.Contains(manifest.id, "postgres") {
					postgresURLs.Values = append(postgresURLs.Values,
						fmt.Sprintf(
							"postgresql://postgres:node@%s:%d",
							manifest.Service.Spec.ClusterIP,
							manifest.Service.Spec.Ports[0].Port,
						))
				}
			}

			mg.values["dbURLs"] = &postgresURLs

			return nil
		},
	}
	for i := 0; i < postgresCount; i++ {
		pManifest := NewPostgresManifest()
		pManifest.id = fmt.Sprintf("%s-%d", pManifest.id, i)
		group.manifests = append(group.manifests, pManifest)
	}

	return func(config *config.NetworkConfig) (string, K8sEnvSpecs) {
		switch config.Name {
		case "Ethereum Geth dev":
			group.manifests = append(
				group.manifests,
				NewGethManifest(),
				NewExplorerManifest(postgresCount))
		case "Ethereum Hardhat":
			group.manifests = append(
				group.manifests,
				NewHardhatManifest())
		case "Ethereum Ganache":
			group.manifests = append(
				group.manifests,
				NewGanacheManifest())
		default: // no simulated chain
			group.manifests = append(
				group.manifests,
				NewExplorerManifest(postgresCount))
		}
		if len(chainlinkGroup.manifests) > 0 {
			return envName, K8sEnvSpecs{group, chainlinkGroup}
		}
		return envName, K8sEnvSpecs{group}
	}
}
