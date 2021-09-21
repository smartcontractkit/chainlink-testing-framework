package environment

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/chartutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/tools"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ports for common services
const (
	AdapterAPIPort   = 6060
	ChainlinkWebPort = 6688
	ChainlinkP2PPort = 6690
	EVMRPCPort       = 8545
	MinersRPCPort    = 9545
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

// NewExplorerManifest is the k8s manifest that when used will deploy explorer to an environment
// and create access keys for a nodeCount number of times
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
			var podsFullNames []string
			for _, pod := range manifest.pods {
				if strings.Contains(pod.PodName, "explorer") {
					podsFullNames = append(podsFullNames, pod.PodName)
				}
			}
			if len(podsFullNames) == 0 {
				return errors.New("")
			}
			_, _, err := manifest.ExecuteInPod(podsFullNames[0], "explorer",
				[]string{"yarn", "--cwd", "apps/explorer", "admin:seed", "username", "password"})
			if err != nil {
				return err
			}

			keys := TemplateValuesArray{}

			explorerClient, err := GetExplorerClientFromEnv(manifest.env)
			if err != nil {
				return err
			}
			for i := 0; i < nodeCount; i++ {
				credentials, err := explorerClient.PostAdminNodes(fmt.Sprintf("node-%d", i))
				if err != nil {
					return err
				}
				keys.Values = append(keys.Values, credentials)
			}
			manifest.values["keys"] = &keys
			return nil
		},
	}
}

// NewOTPEManifest is the k8s manifest for deploying otpe
func NewOTPEManifest() *K8sManifest {
	return &K8sManifest{
		id:             "otpe",
		DeploymentFile: filepath.Join(tools.ProjectRoot, "/environment/templates/otpe-deployment.yml"),
		ServiceFile:    filepath.Join(tools.ProjectRoot, "/environment/templates/otpe-service.yml"),
	}
}

// NewMockserverConfigHelmChart creates new helm chart for the mockserver configmap
func NewMockserverConfigHelmChart() *HelmChart {
	return &HelmChart{
		id:          "mockserver-config",
		chartPath:   filepath.Join(tools.ProjectRoot, "environment/charts/mockserver-config"),
		releaseName: "mockserver-config",
	}
}

// NewMockserverHelmChart creates new helm chart for the mockserver
func NewMockserverHelmChart() *HelmChart {
	chart := &HelmChart{
		id:          "mockserver",
		chartPath:   filepath.Join(tools.ProjectRoot, "environment/charts/mockserver/mockserver-5.11.1.tgz"),
		releaseName: "mockserver",
		values:      map[string]interface{}{},
		SetValuesHelmFunc: func(manifest *HelmChart) error {
			manifest.values["contractsURL"] = "http://mockserver:1080/contracts.json"
			manifest.values["nodesURL"] = "http://mockserver:1080/nodes.json"
			return nil
		},
	}
	return chart
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
	chainlinkGroup := &K8sManifestGroup{
		id:        "chainlinkCluster",
		manifests: []K8sEnvResource{},
	}
	for i := 0; i < nodeCount; i++ {
		cManifest := NewChainlinkManifest()
		cManifest.id = fmt.Sprintf("%s-%d", cManifest.id, i)
		chainlinkGroup.manifests = append(chainlinkGroup.manifests, cManifest)
	}

	dependencyGroup := getBasicDependencyGroup()
	addPostgresDbsToDependencyGroup(dependencyGroup, nodeCount)
	dependencyGroups := []*K8sManifestGroup{dependencyGroup}
	return addNetworkManifestToDependencyGroup("basic-chainlink", chainlinkGroup, dependencyGroups)
}

// NewChainlinkClusterForAlertsTesting is a basic environment that deploys a chainlink cluster with dependencies
// for testing alerts
func NewChainlinkClusterForAlertsTesting(nodeCount int) K8sEnvSpecInit {
	chainlinkGroup := &K8sManifestGroup{
		id:        "chainlinkCluster",
		manifests: []K8sEnvResource{},
	}
	for i := 0; i < nodeCount; i++ {
		cManifest := NewChainlinkManifest()
		cManifest.id = fmt.Sprintf("%s-%d", cManifest.id, i)
		chainlinkGroup.manifests = append(chainlinkGroup.manifests, cManifest)
	}

	kafkaDependecyGroup := &K8sManifestGroup{
		id:        "KafkaGroup",
		manifests: []K8sEnvResource{NewKafkaHelmChart()},
	}

	dependencyGroup := getBasicDependencyGroup()
	addServicesForTestingAlertsToDependencyGroup(dependencyGroup, nodeCount)
	addPostgresDbsToDependencyGroup(dependencyGroup, nodeCount)
	dependencyGroups := []*K8sManifestGroup{kafkaDependecyGroup, dependencyGroup}

	return addNetworkManifestToDependencyGroup("basic-chainlink", chainlinkGroup, dependencyGroups)
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
		manifests: []K8sEnvResource{},
	}
	for i := 0; i < nodeCount; i++ {
		cManifest := NewChainlinkManifest()
		cManifest.id = fmt.Sprintf("%s-%d", cManifest.id, i)
		cManifest.values["image"] = mixedImages[i%len(mixedImages)]
		cManifest.values["version"] = mixedVersions[i%len(mixedVersions)]
		chainlinkGroup.manifests = append(chainlinkGroup.manifests, cManifest)
	}

	dependencyGroup := getBasicDependencyGroup()
	addPostgresDbsToDependencyGroup(dependencyGroup, nodeCount)
	dependencyGroups := []*K8sManifestGroup{dependencyGroup}
	return addNetworkManifestToDependencyGroup("mixed-version-chainlink", chainlinkGroup, dependencyGroups)
}

// NewGethReorgHelmChart creates new helm chart for multi-node Geth network
func NewGethReorgHelmChart() *HelmChart {
	return &HelmChart{
		id:          "evm",
		chartPath:   filepath.Join(tools.ProjectRoot, "environment/charts/geth-reorg"),
		releaseName: "reorg-1",
		SetValuesHelmFunc: func(k *HelmChart) error {
			details, err := k.ServiceDetails()
			if err != nil {
				return err
			}
			for _, d := range details {
				if d.RemoteURL.Port() == strconv.Itoa(EVMRPCPort) {
					k.values["clusterURL"] = strings.Replace(d.RemoteURL.String(), "http", "ws", -1)
					k.values["localURL"] = strings.Replace(d.LocalURL.String(), "http", "ws", -1)
				}
			}
			k.values["rpcPort"] = EVMRPCPort
			return nil
		},
	}
}

// NewKafkaHelmChart creates new helm chart for kafka
func NewKafkaHelmChart() *HelmChart {
	valuesFilePath := filepath.Join(tools.ProjectRoot, "environment/charts/kafka/overrideValues.yaml")
	overrideValues, err := chartutil.ReadValuesFile(valuesFilePath)
	if err != nil {
		return nil
	}

	chart := &HelmChart{
		id:          "kafka",
		chartPath:   filepath.Join(tools.ProjectRoot, "environment/charts/kafka/kafka-14.1.0.tgz"),
		releaseName: "kafka",
		values:      map[string]interface{}{},
		SetValuesHelmFunc: func(manifest *HelmChart) error {
			manifest.values["clusterURL"] = "kafka:9092"
			return nil
		},
	}

	for index, element := range overrideValues {
		chart.values[index] = element
	}

	return chart
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

// getBasicDependencyGroup returns a manifest group containing the basic setup for a chainlink deployment
func getBasicDependencyGroup() *K8sManifestGroup {
	group := &K8sManifestGroup{
		id:        "DependencyGroup",
		manifests: []K8sEnvResource{NewAdapterManifest()},

		SetValuesFunc: func(mg *K8sManifestGroup) error {
			postgresURLs := TemplateValuesArray{}

			for _, manifest := range mg.manifests {
				if strings.Contains(manifest.ID(), "postgres") {
					postgresURLs.Values = append(postgresURLs.Values, manifest.Values()["clusterURL"])
				}
			}

			mg.values["dbURLs"] = &postgresURLs

			return nil
		},
	}
	return group
}

// addNetworkManifestToDependencyGroup adds the correct network to the dependency group and returns
// an array of all groups, this should be called as the last function when creating deploys
func addNetworkManifestToDependencyGroup(envName string, chainlinkGroup *K8sManifestGroup, dependencyGroups []*K8sManifestGroup) K8sEnvSpecInit {
	return func(config *config.NetworkConfig) (string, K8sEnvSpecs) {
		var specs K8sEnvSpecs
		indexOfLastElementInDependencyGroups := len(dependencyGroups) - 1
		switch config.Name {
		case "Ethereum Geth reorg":
			dependencyGroups[indexOfLastElementInDependencyGroups].manifests = append(
				dependencyGroups[indexOfLastElementInDependencyGroups].manifests,
				NewGethReorgHelmChart())
		case "Ethereum Geth dev":
			dependencyGroups[indexOfLastElementInDependencyGroups].manifests = append(
				dependencyGroups[indexOfLastElementInDependencyGroups].manifests,
				NewGethManifest())
		case "Ethereum Hardhat":
			dependencyGroups[indexOfLastElementInDependencyGroups].manifests = append(
				dependencyGroups[indexOfLastElementInDependencyGroups].manifests,
				NewHardhatManifest())
		case "Ethereum Ganache":
			dependencyGroups[indexOfLastElementInDependencyGroups].manifests = append(
				dependencyGroups[indexOfLastElementInDependencyGroups].manifests,
				NewGanacheManifest())
		default: // no simulated chain
		}
		for _, group := range dependencyGroups {
			specs = append(specs, group)
		}

		if len(chainlinkGroup.manifests) > 0 {
			specs = append(specs, chainlinkGroup)
			return envName, specs
		}
		return envName, specs
	}
}

// addPostgresDbsToDependencyGroup adds a postgresCount number of postgres dbs to the dependency group
func addPostgresDbsToDependencyGroup(dependencyGroup *K8sManifestGroup, postgresCount int) {
	for i := 0; i < postgresCount; i++ {
		pManifest := NewPostgresManifest()
		pManifest.id = fmt.Sprintf("%s-%d", pManifest.id, i)
		dependencyGroup.manifests = append(dependencyGroup.manifests, pManifest)
	}
}

// addServicesForTestingAlertsToDependencyGroup adds services necessary for testing alerts to the dependency group
func addServicesForTestingAlertsToDependencyGroup(dependencyGroup *K8sManifestGroup, nodeCount int) {
	dependencyGroup.manifests = append(dependencyGroup.manifests, NewExplorerManifest(nodeCount))
}

// OtpeGroup contains manifests for mockserver, mockserver-config, and otpe
func OtpeGroup() K8sEnvSpecInit {
	return func(config *config.NetworkConfig) (string, K8sEnvSpecs) {
		var specs K8sEnvSpecs
		mockserverConfigDependencyGroup := &K8sManifestGroup{
			id:        "MockserverConfigDependencyGroup",
			manifests: []K8sEnvResource{NewMockserverConfigHelmChart()},
		}

		mockserverDependencyGroup := &K8sManifestGroup{
			id:        "MockserverDependencyGroup",
			manifests: []K8sEnvResource{NewMockserverHelmChart()},
		}

		specs = append(specs, mockserverConfigDependencyGroup)
		specs = append(specs, mockserverDependencyGroup)

		otpeDependencyGroup := &K8sManifestGroup{
			id:        "OTPEDependencyGroup",
			manifests: []K8sEnvResource{NewOTPEManifest()},
		}

		specs = append(specs, otpeDependencyGroup)

		return "envName", specs
	}
}
