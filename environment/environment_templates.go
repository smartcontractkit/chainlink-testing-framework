package environment

import (
	"bufio"
	"context"
	"fmt"
	"github.com/smartcontractkit/integrations-framework/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/chartutil"

	"github.com/google/go-github/github"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ports for common services
const (
	AdapterAPIPort    uint16 = 6060
	ChainlinkWebPort  uint16 = 6688
	ChainlinkP2PPort  uint16 = 6690
	ExplorerAPIPort   uint16 = 8080
	PrometheusAPIPort uint16 = 9090
	MockserverAPIPort uint16 = 1080
	KafkaRestAPIPort  uint16 = 8082
)

// Ethereum ports
const (
	DefaultEVMRPCPort uint16 = 8545
	HardhatRPCPort    uint16 = 8545
	GethRPCPort       uint16 = 8546
	GanacheRPCPort    uint16 = 8547
	MinersRPCPort     uint16 = 9545
)

var (
	defaultChainlinkAuthMap = map[string][]byte{
		"apicredentials": []byte("notreal@fakeemail.ch\ntwochains"),
		"node-password":  []byte("T.tLHkcmwePT/p,]sYuntjwHKAsrhm#4eRs4LuKHwvHejWYAC2JP4M8HimwgmbaZ"),
	}
)

// ChartDeploymentConfig chart deployment configs
type ChartDeploymentConfig struct {
	Name          string
	Path          string
	Values        map[string]interface{}
	SetValuesFunc SetValuesHelmFunc
}

// NewExternalCharts create new charts from chart configs
func NewExternalCharts(chartConfigs []ChartDeploymentConfig) K8sEnvSpecInit {
	envSpecs := make([]K8sEnvResource, 0)
	for _, cfg := range chartConfigs {
		envSpecs = append(envSpecs, &K8sManifestGroup{
			id:        "ExternalDependencyGroup",
			manifests: []K8sEnvResource{NewExternalChart(cfg)},
		})
	}
	return func(networks ...client.BlockchainNetwork) K8sEnvSpecs {
		return envSpecs
	}
}

// NewChainlinkCustomNetworksCluster is a basic environment that deploys headless Chainlink cluster with a custom helm networks
func NewChainlinkCustomNetworksCluster(nodeCount int, networkDeploymentConfigs []ChartDeploymentConfig) K8sEnvSpecInit {
	chainlinkGroup := &K8sManifestGroup{
		id:        "chainlinkCluster",
		manifests: []K8sEnvResource{},
	}
	for i := 0; i < nodeCount; i++ {
		cManifest := NewHeadlessChainlinkManifest(i)
		cManifest.id = fmt.Sprintf("%s-%d", cManifest.id, i)
		chainlinkGroup.manifests = append(chainlinkGroup.manifests, cManifest)
	}
	dependencyGroup := getBasicDependencyGroup()
	addPostgresDbsToDependencyGroup(dependencyGroup, nodeCount)
	return func(networks ...client.BlockchainNetwork) K8sEnvSpecs {
		var specs K8sEnvSpecs
		for _, dc := range networkDeploymentConfigs {
			specs = append(specs, NewExternalChart(dc))
		}
		specs = append(specs, dependencyGroup, chainlinkGroup)
		return specs
	}
}

// NewAdapterManifest is the k8s manifest that when used will deploy an external adapter to an environment
func NewAdapterManifest() *K8sManifest {
	return &K8sManifest{
		id:             "adapter",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/adapter-deployment.yml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/adapter-service.yml"),

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

// NewHeadlessChainlinkManifest is the k8s manifest for Chainlink node without network, using only EI to get network data
func NewHeadlessChainlinkManifest(idx int) *K8sManifest {
	return &K8sManifest{
		id:             "chainlink",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/chainlink/chainlink-deployment.yml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/chainlink/chainlink-service.yml"),

		values: map[string]interface{}{
			"idx":                         idx,
			"webPort":                     ChainlinkWebPort,
			"p2pPort":                     ChainlinkP2PPort,
			"eth_disabled":                true,
			"feature_external_initiators": true,
		},

		Secret: &coreV1.Secret{
			ObjectMeta: v1.ObjectMeta{
				GenerateName: "chainlink-",
			},
			Type: "Opaque",
			Data: defaultChainlinkAuthMap,
		},
	}
}

// NewChainlinkManifest is the k8s manifest that when used will deploy a chainlink node to an environment
func NewChainlinkManifest(idx int) *K8sManifest {
	return &K8sManifest{
		id:             "chainlink",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/chainlink/chainlink-deployment.yml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/chainlink/chainlink-service.yml"),

		values: map[string]interface{}{
			"idx":                         idx,
			"webPort":                     ChainlinkWebPort,
			"p2pPort":                     ChainlinkP2PPort,
			"eth_disabled":                false,
			"feature_external_initiators": false,
		},

		Secret: &coreV1.Secret{
			ObjectMeta: v1.ObjectMeta{
				GenerateName: "chainlink-",
			},
			Type: "Opaque",
			Data: defaultChainlinkAuthMap,
		},
	}
}

// NewPostgresManifest is the k8s manifest that when used will deploy a postgres db to an environment
func NewPostgresManifest() *K8sManifest {
	return &K8sManifest{
		id:             "postgres",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/postgres/postgres-deployment.yml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/postgres/postgres-service.yml"),

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

// NewExplorerManifest is the k8s manifest that when used will deploy explorer to an environment
// and create access keys for a nodeCount number of times
func NewExplorerManifest(nodeCount int) *K8sManifest {
	return &K8sManifest{
		id:             "explorer",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/explorer-deployment.yml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/explorer-service.yml"),
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
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/otpe-deployment.yml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/otpe-service.yml"),
		SetValuesFunc: func(manifest *K8sManifest) error {
			manifest.values["clusterURL"] = fmt.Sprintf(
				"%s:%d",
				manifest.Service.Spec.ClusterIP,
				manifest.Service.Spec.Ports[0].Port,
			)
			return nil
		},
	}
}

// NewMockserverConfigHelmChart creates new helm chart for the mockserver configmap
func NewMockserverConfigHelmChart() *HelmChart {
	return &HelmChart{
		id:          "mockserver-config",
		chartPath:   filepath.Join(utils.ProjectRoot, "environment/charts/mockserver-config"),
		releaseName: "mockserver-config",
	}
}

// NewMockserverHelmChart creates new helm chart for the mockserver
func NewMockserverHelmChart() *HelmChart {
	chart := &HelmChart{
		id:          "mockserver",
		chartPath:   filepath.Join(utils.ProjectRoot, "environment/charts/mockserver/mockserver-5.11.1.tgz"),
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

// NewPrometheusManifest creates new k8s manifest for prometheus
// It receives a map of strings to *os.File which it uses in the following way:
// The string in the map is the template value that will be used the prometheus-config-map.yml file.
// The *os.File contains the rules yaml file, before being added to the values map of the K8sManifest.
// Every line of the file is appended 4 spaces, this is done so after the file is templated to the
// prometheus-config-map.yml file, the yml will be formatted correctly.
func NewPrometheusManifest(rules map[string]*os.File) *K8sManifest {
	vals := map[string]interface{}{}
	for val, file := range rules {
		scanner := bufio.NewScanner(file)
		var txtlines []string
		for scanner.Scan() {
			txtlines = append(txtlines, scanner.Text())
		}
		for index, line := range txtlines {
			txtlines[index] = fmt.Sprintf("    %s", line)
		}
		vals[val] = strings.Join(txtlines, "\n")
	}

	return &K8sManifest{
		id:             "prometheus",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/prometheus/prometheus-deployment.yml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/prometheus/prometheus-service.yml"),
		ConfigMapFile:  filepath.Join(utils.ProjectRoot, "/environment/templates/prometheus/prometheus-config-map.yml"),

		values: vals,
	}
}

// NewGethManifest is the k8s manifest that when used will deploy geth to an environment
func NewGethManifest(networkCount int, network *config.NetworkConfig) *K8sManifest {
	network.Name = fmt.Sprintf("ethereum-geth-%d", networkCount)
	network.RPCPort = GetFreePort()
	return &K8sManifest{
		id:             network.Name,
		DeploymentFile: filepath.Join(utils.ProjectRoot, "environment/templates/geth-deployment.yml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "environment/templates/geth-service.yml"),
		ConfigMapFile:  filepath.Join(utils.ProjectRoot, "environment/templates/geth-config-map.yml"),
		Network:        network,
		values: map[string]interface{}{
			"rpcPort": network.RPCPort,
		},

		SetValuesFunc: func(manifest *K8sManifest) error {
			network.ClusterURL = fmt.Sprintf(
				"ws://%s:%d",
				manifest.Service.Spec.ClusterIP,
				manifest.Service.Spec.Ports[0].Port,
			)
			network.LocalURL = fmt.Sprintf("ws://127.0.0.1:%d", manifest.ports[0].Local)
			manifest.values["clusterURL"] = network.ClusterURL
			manifest.values["localURL"] = network.LocalURL
			return nil
		},
	}
}

// NewExternalChart creates new helm chart
func NewExternalChart(cfg ChartDeploymentConfig) *HelmChart {
	return &HelmChart{
		id:                fmt.Sprintf("%s-%d", cfg.Name, cfg.Values["idx"]),
		chartPath:         cfg.Path,
		releaseName:       fmt.Sprintf("%s-%d", cfg.Name, cfg.Values["idx"]),
		values:            cfg.Values,
		SetValuesHelmFunc: cfg.SetValuesFunc,
	}
}

// NewGethReorgHelmChart creates new helm chart for multi-node Geth network
func NewGethReorgHelmChart(networkCount int, network *config.NetworkConfig) *HelmChart {
	network.Name = fmt.Sprintf("ethereum-geth-reorg-%d", networkCount)
	network.RPCPort = GetFreePort()
	return &HelmChart{
		id:          network.Name,
		chartPath:   filepath.Join(utils.ProjectRoot, "environment/charts/geth-reorg"),
		releaseName: "reorg-1",
		network:     network,
		values: map[string]interface{}{
			"rpcPort": network.RPCPort,
		},
		SetValuesHelmFunc: func(k *HelmChart) error {
			details, err := k.ServiceDetails()
			if err != nil {
				return err
			}
			for _, d := range details {
				if d.RemoteURL.Port() == strconv.Itoa(int(GethRPCPort)) {
					network.ClusterURL = strings.Replace(d.RemoteURL.String(), "http", "ws", -1)
					network.LocalURL = strings.Replace(d.LocalURL.String(), "http", "ws", -1)
				}
			}
			k.values["rpcPort"] = GetFreePort()
			return nil
		},
	}
}

// NewHardhatManifest is the k8s manifest that when used will deploy hardhat to an environment
func NewHardhatManifest(networkCount int, network *config.NetworkConfig) *K8sManifest {
	network.Name = fmt.Sprintf("ethereum-hardhat-%d", networkCount)
	network.RPCPort = GetFreePort()
	return &K8sManifest{
		id:             network.Name,
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/hardhat-deployment.yml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/hardhat-service.yml"),
		ConfigMapFile:  filepath.Join(utils.ProjectRoot, "/environment/templates/hardhat-config-map.yml"),
		Network:        network,
		values: map[string]interface{}{
			"rpcPort": network.RPCPort,
		},

		SetValuesFunc: func(manifest *K8sManifest) error {
			network.ClusterURL = fmt.Sprintf(
				"ws://%s:%d",
				manifest.Service.Spec.ClusterIP,
				manifest.Service.Spec.Ports[0].Port,
			)
			network.LocalURL = fmt.Sprintf("ws://127.0.0.1:%d", manifest.ports[0].Local)
			manifest.values["clusterURL"] = network.ClusterURL
			manifest.values["localURL"] = network.LocalURL
			return nil
		},
	}
}

// NewGanacheManifest is the k8s manifest that when used will deploy ganache to an environment
func NewGanacheManifest(networkCount int, network *config.NetworkConfig) *K8sManifest {
	network.Name = fmt.Sprintf("ethereum-ganache-%d", networkCount)
	network.RPCPort = GetFreePort()
	return &K8sManifest{
		id:             network.Name,
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/ganache-deployment.yml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/ganache-service.yml"),
		Network:        network,
		values: map[string]interface{}{
			"rpcPort": network.RPCPort,
		},

		SetValuesFunc: func(manifest *K8sManifest) error {
			network.ClusterURL = fmt.Sprintf(
				"ws://%s:%d",
				manifest.Service.Spec.ClusterIP,
				manifest.Service.Spec.Ports[0].Port,
			)
			network.LocalURL = fmt.Sprintf("ws://127.0.0.1:%d", manifest.ports[0].Local)
			manifest.values["clusterURL"] = network.ClusterURL
			manifest.values["localURL"] = network.LocalURL
			return nil
		},
	}
}

// NewAtlasEvmBlocksManifest is the k8s manifest that when used will deploy atlas-evm-blocks to an env
func NewAtlasEvmBlocksManifest() *K8sManifest {
	return &K8sManifest{
		id:             "atlas_evm_blocks",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/atlas-evm/atlas-evm-blocks-deployment.yaml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/atlas-evm/atlas-evm-blocks-service.yaml"),
	}
}

// NewAtlasEvmEventsManifest is the k8s manifest that when used will deploy atlas-evm-events to an env
func NewAtlasEvmEventsManifest() *K8sManifest {
	return &K8sManifest{
		id:             "atlas_evm_events",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/atlas-evm/atlas-evm-events-deployment.yaml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/atlas-evm/atlas-evm-events-service.yaml"),
	}
}

// NewAtlasEvmReceiptsManifest is the k8s manifest that when used will deploy atlas-evm-receipts to an env
func NewAtlasEvmReceiptsManifest() *K8sManifest {
	return &K8sManifest{
		id:             "atlas_evm_receipts",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/atlas-evm/atlas-evm-receipts-deployment.yaml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/atlas-evm/atlas-evm-receipts-service.yaml"),
	}
}

// NewSchemaRegistryManifest is the k8s manifest that when used will deploy schema registry to an env
// Confluent Schema Registry provides a serving layer for your metadata. It provides a RESTful interface for storing
// and retrieving your Avro®, JSON Schema, and Protobuf schemas. In Atlas it stores the schemas for different
// components like atlas-evm-blocks, atlas-evm-events etc.
func NewSchemaRegistryManifest() *K8sManifest {
	return &K8sManifest{
		id:             "schema_registry",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/schema-registry/schema-registry-deployment.yaml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/schema-registry/schema-registry-service.yaml"),
		SetValuesFunc: func(manifest *K8sManifest) error {
			manifest.values["clusterURL"] = fmt.Sprintf(
				"http://%s:%d",
				manifest.Service.Spec.ClusterIP,
				manifest.Service.Spec.Ports[0].Port,
			)
			return nil
		},
	}
}

// NewKafkaRestManifest is the k8s manifest that when used will deploy kafka rest to an env
// this is used to retrieve kafka info through REST
func NewKafkaRestManifest() *K8sManifest {
	return &K8sManifest{
		id:             "kafka_rest",
		DeploymentFile: filepath.Join(utils.ProjectRoot, "/environment/templates/kafka-rest/kafka-rest-deployment.yaml"),
		ServiceFile:    filepath.Join(utils.ProjectRoot, "/environment/templates/kafka-rest/kafka-rest-service.yaml"),
	}
}

// NewChainlinkCluster is a basic environment that deploys hardhat with a chainlink cluster and an external adapter
func NewChainlinkCluster(nodeCount int) K8sEnvSpecInit {
	mockserverConfigDependencyGroup := &K8sManifestGroup{
		id:        "MockserverConfigDependencyGroup",
		manifests: []K8sEnvResource{NewMockserverConfigHelmChart()},
	}

	mockserverDependencyGroup := &K8sManifestGroup{
		id:        "MockserverDependencyGroup",
		manifests: []K8sEnvResource{NewMockserverHelmChart()},
	}

	chainlinkGroup := &K8sManifestGroup{
		id:        "chainlinkCluster",
		manifests: []K8sEnvResource{},
	}
	for i := 0; i < nodeCount; i++ {
		cManifest := NewChainlinkManifest(i)
		cManifest.id = fmt.Sprintf("%s-%d", cManifest.id, i)
		chainlinkGroup.manifests = append(chainlinkGroup.manifests, cManifest)
	}

	dependencyGroup := getBasicDependencyGroup()
	addPostgresDbsToDependencyGroup(dependencyGroup, nodeCount)
	dependencyGroups := []*K8sManifestGroup{mockserverConfigDependencyGroup, mockserverDependencyGroup, dependencyGroup}
	return addNetworkManifestToDependencyGroup(chainlinkGroup, dependencyGroups)
}

// NewChainlinkClusterForObservabilityTesting is a basic environment that deploys a chainlink cluster with dependencies
// for testing observability
func NewChainlinkClusterForObservabilityTesting(nodeCount int) K8sEnvSpecInit {
	mockserverConfigDependencyGroup := &K8sManifestGroup{
		id:        "MockserverConfigDependencyGroup",
		manifests: []K8sEnvResource{NewMockserverConfigHelmChart()},
	}

	mockserverDependencyGroup := &K8sManifestGroup{
		id:        "MockserverDependencyGroup",
		manifests: []K8sEnvResource{NewMockserverHelmChart()},
	}

	chainlinkGroup := &K8sManifestGroup{
		id:        "chainlinkCluster",
		manifests: []K8sEnvResource{},
	}
	for i := 0; i < nodeCount; i++ {
		cManifest := NewChainlinkManifest(i)
		cManifest.id = fmt.Sprintf("%s-%d", cManifest.id, i)
		chainlinkGroup.manifests = append(chainlinkGroup.manifests, cManifest)
	}

	kafkaDependecyGroup := &K8sManifestGroup{
		id:        "KafkaGroup",
		manifests: []K8sEnvResource{NewKafkaHelmChart()},
	}

	dependencyGroup := getBasicDependencyGroup()
	dependencyGroup.manifests = append(dependencyGroup.manifests, NewExplorerManifest(nodeCount))
	addPostgresDbsToDependencyGroup(dependencyGroup, nodeCount)
	dependencyGroups := []*K8sManifestGroup{mockserverConfigDependencyGroup, mockserverDependencyGroup, kafkaDependecyGroup, dependencyGroup}

	return addNetworkManifestToDependencyGroup(chainlinkGroup, dependencyGroups)
}

// NewChainlinkClusterForAtlasTesting is a basic environment that deploys a chainlink cluster with dependencies
// for testing Atlas
func NewChainlinkClusterForAtlasTesting(nodeCount int) K8sEnvSpecInit {
	mockserverConfigDependencyGroup := &K8sManifestGroup{
		id:        "MockserverConfigDependencyGroup",
		manifests: []K8sEnvResource{NewMockserverConfigHelmChart()},
	}

	mockserverDependencyGroup := &K8sManifestGroup{
		id:        "MockserverDependencyGroup",
		manifests: []K8sEnvResource{NewMockserverHelmChart()},
	}

	chainlinkGroup := &K8sManifestGroup{
		id:        "chainlinkCluster",
		manifests: []K8sEnvResource{},
	}
	for i := 0; i < nodeCount; i++ {
		cManifest := NewChainlinkManifest(i)
		cManifest.id = fmt.Sprintf("%s-%d", cManifest.id, i)
		chainlinkGroup.manifests = append(chainlinkGroup.manifests, cManifest)
	}

	kafkaDependecyGroup := &K8sManifestGroup{
		id:        "KafkaGroup",
		manifests: []K8sEnvResource{NewKafkaHelmChart()},
	}

	schemaRegistryDependencyGroup := &K8sManifestGroup{
		id:        "SchemaRegistryGroup",
		manifests: []K8sEnvResource{NewSchemaRegistryManifest()},
	}

	kafkaRestDependencyGroup := &K8sManifestGroup{
		id:        "KafkaRestGroup",
		manifests: []K8sEnvResource{NewKafkaRestManifest()},
	}

	dependencyGroup := getBasicDependencyGroup()
	addPostgresDbsToDependencyGroup(dependencyGroup, nodeCount)
	dependencyGroups := []*K8sManifestGroup{
		mockserverConfigDependencyGroup,
		mockserverDependencyGroup,
		kafkaDependecyGroup,
		schemaRegistryDependencyGroup,
		kafkaRestDependencyGroup,
		dependencyGroup,
	}

	return addNetworkManifestToDependencyGroup(chainlinkGroup, dependencyGroups)
}

// NewMixedVersionChainlinkCluster mixes the currently latest chainlink version (as defined by the config file) with
// a number of past stable versions (defined by pastVersionsCount), ensuring that at least one of each is deployed
func NewMixedVersionChainlinkCluster(nodeCount, pastVersionsCount int) K8sEnvSpecInit {
	mockserverConfigDependencyGroup := &K8sManifestGroup{
		id:        "MockserverConfigDependencyGroup",
		manifests: []K8sEnvResource{NewMockserverConfigHelmChart()},
	}

	mockserverDependencyGroup := &K8sManifestGroup{
		id:        "MockserverDependencyGroup",
		manifests: []K8sEnvResource{NewMockserverHelmChart()},
	}

	if nodeCount < pastVersionsCount+1 {
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
		cManifest := NewChainlinkManifest(i)
		cManifest.id = fmt.Sprintf("%s-%d", cManifest.id, i)
		cManifest.values["image"] = mixedImages[i%len(mixedImages)]
		cManifest.values["version"] = mixedVersions[i%len(mixedVersions)]
		chainlinkGroup.manifests = append(chainlinkGroup.manifests, cManifest)
	}

	dependencyGroup := getBasicDependencyGroup()
	addPostgresDbsToDependencyGroup(dependencyGroup, nodeCount)
	dependencyGroups := []*K8sManifestGroup{mockserverConfigDependencyGroup, mockserverDependencyGroup, dependencyGroup}
	return addNetworkManifestToDependencyGroup(chainlinkGroup, dependencyGroups)
}

// NewKafkaHelmChart creates new helm chart for kafka
func NewKafkaHelmChart() *HelmChart {
	valuesFilePath := filepath.Join(utils.ProjectRoot, "environment/charts/kafka/overrideValues.yaml")
	overrideValues, err := chartutil.ReadValuesFile(valuesFilePath)
	if err != nil {
		return nil
	}

	chart := &HelmChart{
		id:          "kafka",
		chartPath:   filepath.Join(utils.ProjectRoot, "environment/charts/kafka/kafka-14.1.0.tgz"),
		releaseName: "kafka",
		values:      map[string]interface{}{},
		SetValuesHelmFunc: func(manifest *HelmChart) error {
			manifest.values["clusterURL"] = "kafka:9092"
			manifest.values["zookeeperURL"] = "kafka-zookeeper:2181"
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
// an array of all groups, this should be called as the last function when creating deployments
func addNetworkManifestToDependencyGroup(chainlinkGroup *K8sManifestGroup, dependencyGroups []*K8sManifestGroup) K8sEnvSpecInit {
	return func(networks ...client.BlockchainNetwork) K8sEnvSpecs {
		var specs K8sEnvSpecs
		indexOfLastElementInDependencyGroups := len(dependencyGroups) - 1
		networkCounts := map[string]int{
			"Ethereum Geth":    0,
			"Ethereum Hardhat": 0,
			"Ethereum Ganache": 0,
		}
		for _, network := range networks {
			switch network.Config().Name {
			case "Ethereum Geth reorg":
				dependencyGroups[indexOfLastElementInDependencyGroups].manifests = append(
					dependencyGroups[indexOfLastElementInDependencyGroups].manifests,
					NewGethReorgHelmChart(networkCounts["Ethereum Geth"], network.Config()))
				networkCounts["Ethereum Geth"] += 1
			case "Ethereum Geth dev":
				dependencyGroups[indexOfLastElementInDependencyGroups].manifests = append(
					dependencyGroups[indexOfLastElementInDependencyGroups].manifests,
					NewGethManifest(networkCounts["Ethereum Geth"], network.Config()))
				networkCounts["Ethereum Geth"] += 1
			case "Ethereum Hardhat":
				dependencyGroups[indexOfLastElementInDependencyGroups].manifests = append(
					dependencyGroups[indexOfLastElementInDependencyGroups].manifests,
					NewHardhatManifest(networkCounts[network.Config().Name], network.Config()))
				networkCounts[network.Config().Name] += 1
			case "Ethereum Ganache":
				dependencyGroups[indexOfLastElementInDependencyGroups].manifests = append(
					dependencyGroups[indexOfLastElementInDependencyGroups].manifests,
					NewGanacheManifest(networkCounts[network.Config().Name], network.Config()))
				networkCounts[network.Config().Name] += 1
			default:
				network.SetClusterURL(network.URLs()[0])
				network.SetLocalURL(network.URLs()[0])
			}
		}

		for _, group := range dependencyGroups {
			specs = append(specs, group)
		}

		if len(chainlinkGroup.manifests) > 0 {
			specs = append(specs, chainlinkGroup)
			return specs
		}
		return specs
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

// OtpeGroup contains manifests for otpe
func OtpeGroup() K8sEnvSpecInit {
	return func(networks ...client.BlockchainNetwork) K8sEnvSpecs {
		var specs K8sEnvSpecs
		otpeDependencyGroup := &K8sManifestGroup{
			id:        "OTPEDependencyGroup",
			manifests: []K8sEnvResource{NewOTPEManifest()},
		}
		specs = append(specs, otpeDependencyGroup)
		return specs
	}
}

// PrometheusGroup contains manifests for prometheus
func PrometheusGroup(rules map[string]*os.File) K8sEnvSpecInit {
	return func(_ ...client.BlockchainNetwork) K8sEnvSpecs {
		var specs K8sEnvSpecs
		prometheusDependencyGroup := &K8sManifestGroup{
			id:        "PrometheusDependencyGroup",
			manifests: []K8sEnvResource{NewPrometheusManifest(rules)},
		}
		specs = append(specs, prometheusDependencyGroup)
		return specs
	}
}

// AtlasEvmBlocksGroup contains manifests for atlas-evm-blocks
func AtlasEvmBlocksGroup() K8sEnvSpecInit {
	return func(networks ...client.BlockchainNetwork) K8sEnvSpecs {
		var specs K8sEnvSpecs
		atlasEvmBlocksDependencyGroup := &K8sManifestGroup{
			id:        "AtlasEvmBlocksGroup",
			manifests: []K8sEnvResource{NewAtlasEvmBlocksManifest()},
		}
		specs = append(specs, atlasEvmBlocksDependencyGroup)
		return specs
	}
}

// AtlasEvmEventsAndReceiptsGroup contains manifests for atlas-evm-events and atlas-evm-receipts
func AtlasEvmEventsAndReceiptsGroup() K8sEnvSpecInit {
	return func(networks ...client.BlockchainNetwork) K8sEnvSpecs {
		var specs K8sEnvSpecs

		atlasEvmEventsAndReceiptsDependencyGroup := &K8sManifestGroup{
			id:        "AtlasEvmEventsAndReceiptsGroup",
			manifests: []K8sEnvResource{NewAtlasEvmEventsManifest(), NewAtlasEvmReceiptsManifest()},
		}

		specs = append(specs, atlasEvmEventsAndReceiptsDependencyGroup)
		return specs
	}
}
