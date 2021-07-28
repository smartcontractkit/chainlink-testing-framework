package environment

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/smartcontractkit/integrations-framework/config"
	"gopkg.in/yaml.v2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const (
	SelectorLabelKey        = "app"
	BlockchainAppLabelValue = "blockchain"
	ChainlinkAppLabelValue  = "chainlink-node"
	AdapterAppLabelValue    = "external-adapter"
)

type Environment interface {
	ChainlinkNodes() []client.Chainlink
	ChainlinkNodeETHAddresses() ([]common.Address, error)
	Adapter() ExternalAdapter
	BlockchainClient() client.BlockchainClient
	Wallets() client.BlockchainWallets
	ContractDeployer() contracts.ContractDeployer
	LinkContract() contracts.LinkToken

	GetLocalPorts(appLabel string, remotePort uint16) ([]uint16, error)
	GetLocalPort(appLabel string, remotePort uint16) (uint16, error)

	FundAllNodes(fromWallet client.BlockchainWallet, nativeAmount, linkAmount *big.Int) error
	TearDown() error
}

type k8sEnvironment struct {
	// K8s Resources
	kubeClient *kubernetes.Clientset
	kubeConfig *rest.Config

	// Deployment resources
	manifests K8sManifests
	namespace *coreV1.Namespace

	// Environment resources
	network          client.BlockchainNetwork
	chainlinkNodes   []client.Chainlink
	adapter          ExternalAdapter
	blockchainClient client.BlockchainClient
	wallets          client.BlockchainWallets
	contractDeployer contracts.ContractDeployer
	linkContract     contracts.LinkToken

	// Connection resources
	ports        map[string][]portforward.ForwardedPort
	stopChannels []chan struct{}
	mutex        *sync.Mutex
}

// NewK8sEnvironment launches a new k8sEnvironment of latest version chainlink nodes, connected to the specified network
func NewK8sEnvironment(envInit K8sEnvironmentInit, network client.BlockchainNetwork) (Environment, error) {
	k8sConfig, err := k8sConfig()
	if err != nil {
		return nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	env := &k8sEnvironment{
		kubeClient:     k8sClient,
		kubeConfig:     k8sConfig,
		network:        network,
		chainlinkNodes: []client.Chainlink{},
		mutex:          &sync.Mutex{},
	}

	environmentName, k8sManifests := envInit()
	env.manifests = k8sManifests
	namespace, err := env.createNamespace(environmentName)
	if err != nil {
		return nil, err
	}
	env.namespace = namespace

	// TODO: Use templating for this bit
	if err := env.deployChainlinkSecrets(); err != nil {
		return nil, err
	}
	if err := env.deployManifests(); err != nil {
		return nil, err
	}
	if err := env.waitForHealthyPods(); err != nil {
		return nil, err
	}
	ports, err := env.forwardPorts()
	if err != nil {
		return nil, err
	}
	env.ports = ports
	if err := env.initServices(); err != nil {
		return nil, err
	}

	err = env.setEnvTools()
	return env, err
}

// ChainlinkNodes returns all the chainlink nodes in the launched k8sEnvironment
func (env *k8sEnvironment) ChainlinkNodes() []client.Chainlink {
	return env.chainlinkNodes
}

// Adapter returns dummy external adapter that the k8sEnvironment has deployed
func (env *k8sEnvironment) Adapter() ExternalAdapter {
	return env.adapter
}

// ChainlinkNodeETHAddresses returns the primary ETH addresses of all the chainlink nodes in the launched k8sEnvironment
func (env *k8sEnvironment) ChainlinkNodeETHAddresses() ([]common.Address, error) {
	addresses := make([]common.Address, len(env.chainlinkNodes))
	for _, node := range env.chainlinkNodes {
		primaryAddress, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(primaryAddress))
	}
	return addresses, nil
}

// BlockchainClient retrieves the blockchain client for the k8sEnvironment
func (env *k8sEnvironment) BlockchainClient() client.BlockchainClient {
	return env.blockchainClient
}

// Wallets retrieves the configured wallets for the k8sEnvironment
func (env *k8sEnvironment) Wallets() client.BlockchainWallets {
	return env.wallets
}

// ContractDeployer retrieves the deployer that allows further contracts to be deployed to the k8sEnvironment
func (env *k8sEnvironment) ContractDeployer() contracts.ContractDeployer {
	return env.contractDeployer
}

// LinkContract retrieves the deployed link contract for the k8sEnvironment
func (env *k8sEnvironment) LinkContract() contracts.LinkToken {
	return env.linkContract
}

// FundAllNodes funds all chainlink nodes in the k8sEnvironment with the specified wallet for the specified amounts
func (env *k8sEnvironment) FundAllNodes(fromWallet client.BlockchainWallet, nativeAmount, linkAmount *big.Int) error {
	for _, cl := range env.chainlinkNodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		err = env.blockchainClient.Fund(fromWallet, toAddress, nativeAmount, linkAmount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (env *k8sEnvironment) GetLocalPorts(appLabel string, remotePort uint16) ([]uint16, error) {
	var localPorts []uint16
	ports, ok := env.ports[appLabel]
	if !ok {
		return nil, fmt.Errorf("app label doesn't exist in the deployment")
	}
	for _, port := range ports {
		if port.Remote == remotePort {
			localPorts = append(localPorts, port.Local)
		}
	}
	if len(localPorts) == 0 {
		return nil, fmt.Errorf("remote port doesn't exist in the deployment")
	}
	return localPorts, nil
}

func (env *k8sEnvironment) GetLocalPort(appLabel string, remotePort uint16) (uint16, error) {
	if ports, err := env.GetLocalPorts(appLabel, remotePort); err != nil {
		return 0, err
	} else {
		return ports[0], nil
	}
}

// TearDown calls delete on all the k8sEnvironment's resources
func (env *k8sEnvironment) TearDown() error {
	env.mutex.Lock()
	defer env.mutex.Unlock()

	for _, stopChan := range env.stopChannels {
		stopChan <- struct{}{}
	}

	err := env.kubeClient.CoreV1().Namespaces().Delete(context.Background(), env.namespace.Name, metaV1.DeleteOptions{})
	log.Info().Str("Namespace", env.namespace.Name).Msg("Deleted environment")
	return err
}

func (env *k8sEnvironment) createNamespace(namespace string) (*coreV1.Namespace, error) {
	createdNamespace, err := env.kubeClient.CoreV1().Namespaces().Create(
		context.Background(),
		&coreV1.Namespace{
			ObjectMeta: metaV1.ObjectMeta{
				GenerateName: namespace + "-",
			},
		},
		metaV1.CreateOptions{},
	)
	if err == nil {
		log.Info().Str("Namespace", createdNamespace.Name).Msg("Created namespace")
	}
	return createdNamespace, err
}

func (env *k8sEnvironment) deployManifests() error {
	k8sSecrets := env.kubeClient.CoreV1().Secrets(env.namespace.Name)
	k8sDeployments := env.kubeClient.AppsV1().Deployments(env.namespace.Name)
	k8sServices := env.kubeClient.CoreV1().Services(env.namespace.Name)
	deployedManifests := map[string]*K8sManifest{}

	for k, manifest := range env.manifests {
		if err := manifest.Parse(env.network, deployedManifests); err != nil {
			return err
		}

		if len(manifest.SecretFile) > 0 {
			if secret, err := k8sSecrets.Create(
				context.Background(),
				manifest.Secret,
				metaV1.CreateOptions{},
			); err != nil {
				return fmt.Errorf("failed to deploy %s in k8s: %v", manifest.SecretFile, err)
			} else {
				env.manifests[k].Secret = secret
			}
		}

		if len(manifest.DeploymentFile) > 0 {
			if deployment, err := k8sDeployments.Create(
				context.Background(),
				manifest.Deployment,
				metaV1.CreateOptions{},
			); err != nil {
				return fmt.Errorf("failed to deploy %s in k8s: %v", manifest.DeploymentFile, err)
			} else {
				env.manifests[k].Deployment = deployment
			}
		}

		if len(manifest.ServiceFile) > 0 {
			if service, err := k8sServices.Create(
				context.Background(),
				manifest.Service,
				metaV1.CreateOptions{},
			); err != nil {
				return fmt.Errorf("failed to create service %s in k8s: %v", manifest.ServiceFile, err)
			} else {
				env.manifests[k].Service = service
			}
		}

		deployedManifests[manifest.Type] = manifest
	}
	return nil
}

func (env *k8sEnvironment) initServices() error {
	for _, manifest := range env.manifests {
		if err := manifest.CallbackFunc(env); err != nil {
			return err
		}
	}
	return nil
}

func (env *k8sEnvironment) findServicesBySelector(k, v string) ([]*coreV1.Service, error) {
	var services []*coreV1.Service
	for _, manifest := range env.manifests {
		for sk, sv := range manifest.Service.Spec.Selector {
			if sk == k && sv == v {
				services = append(services, manifest.Service)
			}
		}
	}
	if len(services) > 0 {
		return services, nil
	}
	return nil, fmt.Errorf(
		"service by label `%s: %s` is not found, make sure the label exists in the templates",
		k,
		v,
	)
}

func (env *k8sEnvironment) setEnvTools() error {
	blockchainClient, err := client.NewBlockchainClient(env.network)
	if err != nil {
		return err
	}
	wallets, err := env.network.Wallets()
	if err != nil {
		return err
	}
	contractDeployer, err := contracts.NewContractDeployer(blockchainClient)
	if err != nil {
		return err
	}
	linkContract, err := contractDeployer.DeployLinkTokenContract(wallets.Default())
	if err != nil {
		return err
	}
	env.blockchainClient = blockchainClient
	env.wallets = wallets
	env.contractDeployer = contractDeployer
	env.linkContract = linkContract
	return err
}

// Waits for all pods in the namespace to report as running healthy
func (env *k8sEnvironment) waitForHealthyPods() error {
	namespace := env.namespace.Name

	start := time.Now()
	log.Info().Str("Namespace", namespace).Msg("Waiting for environment to be healthy")
	ticker := time.NewTicker(time.Millisecond * 500)

	defer ticker.Stop()
	for range ticker.C {
		podInterface := env.kubeClient.CoreV1().Pods(namespace)
		pods, err := podInterface.List(context.Background(), metaV1.ListOptions{})
		if err != nil {
			return err
		}
		healthyPodCount := 0
		for _, pod := range pods.Items {
			for _, condition := range pod.Status.Conditions { // Each pod has an unordered list of conditions
				if condition.Type == coreV1.PodReady {
					if condition.Status == coreV1.ConditionTrue {
						healthyPodCount += 1
						break
					}
				}
			}
		}
		if healthyPodCount == len(pods.Items) {
			log.Info().
				Str("Namespace", namespace).
				Str("Wait Length", time.Since(start).Round(time.Second).String()).
				Msg("Environment healthy")
			return nil
		}
	}
	return nil
}

// Forwards all ports needed to connect to the k8sEnvironment
func (env *k8sEnvironment) forwardPorts() (map[string][]portforward.ForwardedPort, error) {
	forwardedPorts := map[string][]portforward.ForwardedPort{}

	k8sPods := env.kubeClient.CoreV1().Pods(env.namespace.Name)
	podList, err := k8sPods.List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range podList.Items {
		log.Info().Str("Pod", pod.Name).Msg("Forwarding ports")

		ports, err := env.forwardPort(&pod)
		if err != nil {
			return nil, err
		}
		label, ok := pod.Labels["app"]
		if !ok {
			return nil, fmt.Errorf("pod %s cannot be port forwarded as it doesn't have the 'app' label", pod.Name)
		}
		if _, ok := forwardedPorts[label]; ok {
			forwardedPorts[label] = append(forwardedPorts[label], ports...)
		} else {
			forwardedPorts[label] = ports
		}

		logger := log.Info().Str("Pod", pod.Name)
		for i, port := range ports {
			logger.Str(fmt.Sprintf("Port %d", i), fmt.Sprintf("%d:%d", port.Local, port.Remote))
		}
		logger.Msg("Forwarded ports")
	}

	return forwardedPorts, err
}

func (env *k8sEnvironment) forwardPort(pod *coreV1.Pod) ([]portforward.ForwardedPort, error) {
	roundTripper, upgrader, err := spdy.RoundTripperFor(env.kubeConfig)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", env.namespace.Name, pod.Name)
	hostIP := strings.TrimLeft(env.kubeConfig.Host, "htps:/")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)

	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	var ports []string
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			ports = append(ports, fmt.Sprintf("0:%d", port.ContainerPort))
		}
	}

	forwarder, err := portforward.New(dialer, ports, stopChan, readyChan, out, errOut)
	if err != nil {
		return nil, err
	}
	go env.doPortForward(pod, forwarder)

	<-readyChan
	if len(errOut.String()) > 0 {
		return nil, fmt.Errorf("error on forwarding k8s port: %v", errOut.String())
	}
	if len(out.String()) > 0 {
		log.Debug().Str("Pod", pod.Name).Msgf("Debug message on port forward: %s", out.String())
	}

	env.mutex.Lock()
	env.stopChannels = append(env.stopChannels, stopChan)
	env.mutex.Unlock()

	return forwarder.GetPorts()
}

func (env *k8sEnvironment) doPortForward(pod *coreV1.Pod, forwarder *portforward.PortForwarder) {
	if err := forwarder.ForwardPorts(); err != nil {
		log.Error().Str("Pod", pod.Name).Err(err)
	}
}

func (env *k8sEnvironment) deployChainlinkSecrets() error {
	k8sSecrets := env.kubeClient.CoreV1().Secrets(env.namespace.Name)

	secret := &coreV1.Secret{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "node-secrets",
		},
		Type: coreV1.SecretType("Opaque"),
		StringData: map[string]string{
			"apicredentials": "notreal@fakeemail.ch\ntwochains",
			"node-password":  "T.tLHkcmwePT/p,]sYuntjwHKAsrhm#4eRs4LuKHwvHejWYAC2JP4M8HimwgmbaZ",
		},
	}

	_, err := k8sSecrets.Create(context.Background(), secret, metaV1.CreateOptions{})
	return err
}

type K8sEnvironmentInit func() (string, K8sManifests)

type K8sManifests map[int]*K8sManifest

type k8sCallback func(env *k8sEnvironment) error

type K8sManifest struct {
	Type           string
	DeploymentFile string
	ServiceFile    string
	SecretFile     string
	Deployment     *appsV1.Deployment
	Service        *coreV1.Service
	Secret         *coreV1.Secret
	Config         *config.Config

	CallbackFunc k8sCallback
}

func (k8s *K8sManifest) Parse(network client.BlockchainNetwork, previousManifests map[string]*K8sManifest) error {
	if len(k8s.SecretFile) > 0 {
		if k8s.Secret == nil {
			k8s.Secret = &coreV1.Secret{}
		}
		if err := k8s.parse(k8s.SecretFile, k8s.Secret, network, previousManifests); err != nil {
			return err
		}
	}
	if len(k8s.ServiceFile) > 0 {
		if k8s.Service == nil {
			k8s.Service = &coreV1.Service{}
		}
		if err := k8s.parse(k8s.ServiceFile, k8s.Service, network, previousManifests); err != nil {
			return err
		}
	}
	if len(k8s.DeploymentFile) > 0 {
		if k8s.Deployment == nil {
			k8s.Deployment = &appsV1.Deployment{}
		}
		if err := k8s.parse(k8s.DeploymentFile, k8s.Deployment, network, previousManifests); err != nil {
			return err
		}
	}
	return nil
}

func (k8s *K8sManifest) parse(
	path string, obj interface{},
	network client.BlockchainNetwork,
	previousManifests map[string]*K8sManifest,
) error {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read k8s file: %v", err)
	}
	tpl, err := template.New(path).Parse(string(fileBytes))
	if err != nil {
		return fmt.Errorf("failed to read k8s template file %s: %v", path, err)
	}

	chainlinkConnectDetails := struct {
		Version  string
		ChainID  string
		ChainURL string
	}{
		Version:  "0.10.10", // TODO: Make this variable
		ChainID:  network.ChainID().String(),
		ChainURL: network.URL(),
	}
	if chainManifest, ok := previousManifests[BlockchainAppLabelValue]; ok {
		chainlinkConnectDetails.ChainURL = "ws://" + chainManifest.Service.Spec.ClusterIP + ":8545"
	}
	var tplBuffer bytes.Buffer
	if err := tpl.Execute(&tplBuffer, chainlinkConnectDetails); err != nil {
		return fmt.Errorf("failed to execute k8s template file %s: %v", path, err)
	}
	if err := yaml.Unmarshal(tplBuffer.Bytes(), obj); err != nil {
		return fmt.Errorf("failed to unmarshall k8s template file %s: %v", path, err)
	}
	return nil
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

func k8sConfig() (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	return kubeConfig.ClientConfig()
}
