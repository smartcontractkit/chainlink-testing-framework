package environment

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/avast/retry-go"
	"github.com/hashicorp/go-multierror"
	"github.com/smartcontractkit/integrations-framework/config"
	"gopkg.in/yaml.v2"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const SelectorLabelKey = "app"

// K8sEnvSpecs represents a series of environment resources to be deployed. The map keys need to be continuous with
// no gaps. For example:
// 0: Hardhat
// 1: Adapter
// 2: Chainlink cluster
type K8sEnvSpecs map[int]K8sEnvResource

// K8sEnvSpecInit is the initiator that will return the name of the environment and the specifications to be deployed.
// The name of the environment returned determines the namespace.
type K8sEnvSpecInit func(*config.NetworkConfig) (string, K8sEnvSpecs)

// K8sEnvResource is the interface for deploying a given environment resource. Creating an interface for resource
// deployment allows it to be extended, deploying k8s resources in different ways. For example: K8sManifest deploys
// a single manifest, whereas K8sManifestGroup bundles several K8sManifests to be deployed concurrently.
type K8sEnvResource interface {
	ID() string
	SetEnvironment(
		k8sClient *kubernetes.Clientset,
		k8sConfig *rest.Config,
		config *config.Config,
		network *config.NetworkConfig,
		namespace *coreV1.Namespace,
	) error
	Deploy(values map[string]interface{}) error
	WaitUntilHealthy() error
	ServiceDetails() ([]*ServiceDetails, error)
	Values() map[string]interface{}
	Teardown() error
}

type k8sEnvironment struct {
	// K8s Resources
	k8sClient *kubernetes.Clientset
	k8sConfig *rest.Config

	// Deployment resources
	specs     K8sEnvSpecs
	namespace *coreV1.Namespace

	// Environment resources
	config  *config.Config
	network client.BlockchainNetwork
}

// NewK8sEnvironment creates and deploys a full ephemeral environment in a k8s cluster. Your current context within
// your kube config will always be used.
func NewK8sEnvironment(
	init K8sEnvSpecInit,
	cfg *config.Config,
	network client.BlockchainNetwork,
) (Environment, error) {
	k8sConfig, err := k8sConfig()
	if err != nil {
		return nil, err
	}
	k8sConfig.QPS = cfg.Kubernetes.QPS
	k8sConfig.Burst = cfg.Kubernetes.Burst

	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	env := &k8sEnvironment{
		k8sClient: k8sClient,
		k8sConfig: k8sConfig,
		config:    cfg,
		network:   network,
	}

	environmentName, deployables := init(network.Config())
	namespace, err := env.createNamespace(environmentName)
	if err != nil {
		return nil, err
	}
	env.namespace = namespace
	env.specs = deployables

	if err := env.deploySpecs(); err != nil {
		return nil, err
	}
	return env, err
}

// ID returns the canonical name of the environment, which in the case of k8s is the namespace
func (env *k8sEnvironment) ID() string {
	if env.namespace != nil {
		return env.namespace.Name
	}
	return ""
}

// GetAllServiceDetails returns all the connectivity details for a deployed service by its remote port within k8s
func (env *k8sEnvironment) GetAllServiceDetails(remotePort uint16) ([]*ServiceDetails, error) {
	var serviceDetails []*ServiceDetails
	var matchedServiceDetails []*ServiceDetails

	for _, spec := range env.specs {
		specServiceDetails, err := spec.ServiceDetails()
		if err != nil {
			return nil, err
		}
		serviceDetails = append(serviceDetails, specServiceDetails...)
	}
	for _, service := range serviceDetails {
		if service.RemoteURL.Port() == fmt.Sprint(remotePort) {
			matchedServiceDetails = append(matchedServiceDetails, service)
		}
	}
	if len(matchedServiceDetails) == 0 {
		return nil, fmt.Errorf("no services with the remote port %d have been deployed", remotePort)
	}
	return matchedServiceDetails, nil
}

// GetServiceDetails returns all the connectivity details for a deployed service by its remote port within k8s
func (env *k8sEnvironment) GetServiceDetails(remotePort uint16) (*ServiceDetails, error) {
	if serviceDetails, err := env.GetAllServiceDetails(remotePort); err != nil {
		return nil, err
	} else {
		return serviceDetails[0], err
	}
}

// TearDown cycles through all the specifications and tears down the deployments. This typically entails cleaning
// up port forwarding requests and deleting the namespace that then destroys all definitions.
func (env *k8sEnvironment) TearDown() {
	for _, spec := range env.specs {
		if err := spec.Teardown(); err != nil {
			log.Error().Err(err)
		}
	}

	err := env.k8sClient.CoreV1().Namespaces().Delete(context.Background(), env.namespace.Name, metaV1.DeleteOptions{})
	log.Info().Str("Namespace", env.namespace.Name).Msg("Deleted environment")
	log.Error().Err(err)
}

func (env *k8sEnvironment) deploySpecs() error {
	values := map[string]interface{}{}
	for i := 0; i < len(env.specs); i++ {
		spec, ok := env.specs[i]
		if !ok {
			return fmt.Errorf("specifcation %d wasn't found on deploy, make sure the set are in order", i)
		}
		if err := spec.SetEnvironment(
			env.k8sClient,
			env.k8sConfig,
			env.config,
			env.network.Config(),
			env.namespace,
		); err != nil {
			return err
		}
		values[spec.ID()] = spec.Values()
		if err := spec.Deploy(values); err != nil {
			return err
		}
		if err := spec.WaitUntilHealthy(); err != nil {
			return err
		}
		values[spec.ID()] = spec.Values()
	}
	return nil
}

func (env *k8sEnvironment) createNamespace(namespace string) (*coreV1.Namespace, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createdNamespace, err := env.k8sClient.CoreV1().Namespaces().Create(
		ctx,
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

type k8sTemplateData struct {
	Config   *config.Config
	Network  *config.NetworkConfig
	Values   map[string]interface{}
	Manifest *K8sManifest
}

type k8sSetValuesFunc func(*K8sManifest) error

// K8sManifest represents a manifest of k8s definitions to be deployed. It implements the K8sEnvResource interface
// to allow the deployment of the definitions into a cluster. It consists of a k8s secret, deployment and service
// but can be expanded to allow more definitions if needed, or extended with another interface to expand on its
// functionality.
type K8sManifest struct {
	id string

	// Manifest properties
	DeploymentFile string
	ServiceFile    string
	SecretFile     string
	ConfigMapFile  string
	Deployment     *appsV1.Deployment
	Service        *coreV1.Service
	ConfigMap      *coreV1.ConfigMap
	Secret         *coreV1.Secret
	SetValuesFunc  k8sSetValuesFunc

	// Deployment properties
	ports        []portforward.ForwardedPort
	values       map[string]interface{}
	stopChannels []chan struct{}

	// Environment properties
	k8sClient *kubernetes.Clientset
	k8sConfig *rest.Config
	config    *config.Config
	network   *config.NetworkConfig
	namespace *coreV1.Namespace
}

// SetEnvironment is the K8sEnvResource implementation that sets the current cluster and config to be used on deploy
func (m *K8sManifest) SetEnvironment(
	k8sClient *kubernetes.Clientset,
	k8sConfig *rest.Config,
	config *config.Config,
	network *config.NetworkConfig,
	namespace *coreV1.Namespace,
) error {
	m.k8sClient = k8sClient
	m.k8sConfig = k8sConfig
	m.config = config
	m.network = network
	m.namespace = namespace
	return nil
}

// ID returns the identifier for the manifest. The ID is important as the manifest will automatically add labels and
// service selectors to link deployments to their manifests.
func (m *K8sManifest) ID() string {
	return m.id
}

// Deploy will create the definitions for each manifest on the k8s cluster
func (m *K8sManifest) Deploy(values map[string]interface{}) error {
	if err := m.createConfigMap(values); err != nil {
		return err
	}
	if err := m.createSecret(values); err != nil {
		return err
	}
	if err := m.createDeployment(values); err != nil {
		return err
	}
	if err := m.createService(values); err != nil {
		return err
	}
	return nil
}

// WaitUntilHealthy will wait until all pods that are created from a given manifest are healthy. Once healthy, it will
// then forward all ports that are exposed within the service and callback to set values.
func (m *K8sManifest) WaitUntilHealthy() error {
	k8sPods := m.k8sClient.CoreV1().Pods(m.namespace.Name)

	// Have a retry mechanism here as if the k8s cluster is under strain, then the pods will not
	// appear instantly after deployment
	var pods *coreV1.PodList
	err := retry.Do(
		func() error {
			labelSelector := fmt.Sprintf("%s=%s", SelectorLabelKey, m.id)
			localPods, localErr := k8sPods.List(context.Background(), metaV1.ListOptions{
				LabelSelector: labelSelector,
			})

			if localErr != nil {
				return fmt.Errorf("unable to fetch pods after deployment: %v", localErr)
			} else if len(localPods.Items) == 0 {
				return fmt.Errorf("no pods returned for manifest %s after deploying", m.id)
			}
			pods = localPods

			return nil
		},
		retry.Delay(time.Millisecond*500),
	)
	if err != nil {
		return err
	}

	if err := waitForHealthyPods(m.k8sClient, m.namespace, pods); err != nil {
		return err
	}

	for _, p := range pods.Items {
		ports, err := m.forwardPodPorts(&p)
		if err != nil {
			return fmt.Errorf("unable to forward ports: %v", err)
		}
		m.ports = append(m.ports, ports...)
		log.Info().Str("Manifest ID", m.id).Interface("Ports", ports).Msg("Forwarded ports")
	}

	if m.SetValuesFunc != nil {
		return m.setValues()
	}
	return nil
}

// ServiceDetails returns the connectivity details for a deployed service
func (m *K8sManifest) ServiceDetails() ([]*ServiceDetails, error) {
	var serviceDetails []*ServiceDetails
	for _, port := range m.ports {
		remoteURL, err := url.Parse(fmt.Sprintf("http://%s:%d", m.Service.Spec.ClusterIP, port.Remote))
		if err != nil {
			return serviceDetails, err
		}
		localURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port.Local))
		if err != nil {
			return serviceDetails, err
		}
		serviceDetails = append(serviceDetails, &ServiceDetails{
			RemoteURL:  remoteURL,
			LocalURL:   localURL,
		})
	}
	return serviceDetails, nil
}

// Values returns all the values to be exposed in the definition templates
func (m *K8sManifest) Values() map[string]interface{} {
	return m.values
}

// Teardown sends a message to the port forwarding channels to stop the request
func (m *K8sManifest) Teardown() error {
	for _, stopChan := range m.stopChannels {
		stopChan <- struct{}{}
	}
	return nil
}

func (m *K8sManifest) createSecret(values map[string]interface{}) error {
	k8sSecrets := m.k8sClient.CoreV1().Secrets(m.namespace.Name)

	if err := m.parseSecret(m.config, m.network, values); err != nil {
		return err
	}

	if m.Secret != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if secret, err := k8sSecrets.Create(
			ctx,
			m.Secret,
			metaV1.CreateOptions{},
		); err != nil {
			return fmt.Errorf("failed to deploy %s in k8s: %v", m.SecretFile, err)
		} else {
			m.Secret = secret
		}
	}
	return nil
}

func (m *K8sManifest) createDeployment(values map[string]interface{}) error {
	k8sDeployments := m.k8sClient.AppsV1().Deployments(m.namespace.Name)

	if err := m.parseDeployment(m.config, m.network, values); err != nil {
		return err
	}

	if m.Deployment != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if deployment, err := k8sDeployments.Create(
			ctx,
			m.Deployment,
			metaV1.CreateOptions{},
		); err != nil {
			return fmt.Errorf("failed to deploy %s in k8s: %v", m.SecretFile, err)
		} else {
			m.Deployment = deployment
		}
	}
	return nil
}

func (m *K8sManifest) createService(values map[string]interface{}) error {
	k8sServices := m.k8sClient.CoreV1().Services(m.namespace.Name)

	if err := m.parseService(m.config, m.network, values); err != nil {
		return err
	}

	if m.Service != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if service, err := k8sServices.Create(
			ctx,
			m.Service,
			metaV1.CreateOptions{},
		); err != nil {
			return fmt.Errorf("failed to deploy %s in k8s: %v", m.SecretFile, err)
		} else {
			m.Service = service
		}
	}
	return nil
}

func (m *K8sManifest) createConfigMap(values map[string]interface{}) error {
	cm := m.k8sClient.CoreV1().ConfigMaps(m.namespace.Name)
	if err := m.parseConfigMap(m.config, m.network, values); err != nil {
		return err
	}
	if m.ConfigMap != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if configMap, err := cm.Create(ctx, m.ConfigMap, metaV1.CreateOptions{}); err != nil {
			return fmt.Errorf("failed to deploy %s in k8s: %v", m.ConfigMap, err)
		} else {
			m.ConfigMap = configMap
		}
	}
	return nil
}

func (m *K8sManifest) parseConfigMap(
	cfg *config.Config,
	network *config.NetworkConfig,
	values map[string]interface{},
) error {
	if len(m.ConfigMapFile) > 0 && m.ConfigMap == nil {
		m.ConfigMap = &coreV1.ConfigMap{}
		if err := m.parse(m.ConfigMapFile, m.ConfigMap, m.initTemplateData(cfg, network, values)); err != nil {
			return err
		}
	}
	return nil
}

func (m *K8sManifest) parseSecret(
	cfg *config.Config,
	network *config.NetworkConfig,
	values map[string]interface{},
) error {
	if len(m.SecretFile) > 0 && m.Secret == nil {
		m.Secret = &coreV1.Secret{}
		if err := m.parse(
			m.SecretFile,
			m.Secret,
			m.initTemplateData(cfg, network, values),
		); err != nil {
			return err
		}
	}
	return nil
}

func (m *K8sManifest) parseDeployment(
	cfg *config.Config,
	network *config.NetworkConfig,
	values map[string]interface{},
) error {
	if len(m.DeploymentFile) > 0 && m.Deployment == nil {
		m.Deployment = &appsV1.Deployment{}
		if err := m.parse(
			m.DeploymentFile,
			m.Deployment,
			m.initTemplateData(cfg, network, values),
		); err != nil {
			return err
		}
	}
	return nil
}

func (m *K8sManifest) parseService(
	cfg *config.Config,
	network *config.NetworkConfig,
	values map[string]interface{},
) error {
	if len(m.ServiceFile) > 0 && m.Service == nil {
		m.Service = &coreV1.Service{}
		if err := m.parse(
			m.ServiceFile,
			m.Service,
			m.initTemplateData(cfg, network, values),
		); err != nil {
			return err
		}
	}
	return nil
}

func (m *K8sManifest) initTemplateData(
	cfg *config.Config,
	network *config.NetworkConfig,
	values map[string]interface{},
) k8sTemplateData {
	return k8sTemplateData{
		Config:   cfg,
		Network:  network,
		Values:   values,
		Manifest: m,
	}
}

func (m *K8sManifest) setLabels() {
	if m.Deployment != nil {
		if m.Deployment.Spec.Selector == nil {
			m.Deployment.Spec.Selector = &metaV1.LabelSelector{}
		}
		if m.Deployment.Spec.Selector.MatchLabels == nil {
			m.Deployment.Spec.Selector.MatchLabels = map[string]string{}
		}
		m.Deployment.Spec.Selector.MatchLabels[SelectorLabelKey] = m.id
		if m.Deployment.Spec.Template.ObjectMeta.Labels == nil {
			m.Deployment.Spec.Template.ObjectMeta.Labels = map[string]string{}
		}
		m.Deployment.Spec.Template.ObjectMeta.Labels[SelectorLabelKey] = m.id
	}
	if m.Service != nil {
		m.Service.Spec.Selector[SelectorLabelKey] = m.id
		if m.Service.Spec.Selector == nil {
			m.Service.Spec.Selector = map[string]string{}
		}
	}
}

func (m *K8sManifest) setValues() error {
	if m.values == nil {
		m.values = map[string]interface{}{}
	}
	if err := m.SetValuesFunc(m); err != nil {
		return err
	}
	return nil
}

func (m *K8sManifest) parse(path string, obj interface{}, data interface{}) error {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read k8s file: %v", err)
	}
	tpl, err := template.New(path).Parse(string(fileBytes))
	if err != nil {
		return fmt.Errorf("failed to read k8s template file %s: %v", path, err)
	}
	var tplBuffer bytes.Buffer
	if err := tpl.Execute(&tplBuffer, data); err != nil {
		return fmt.Errorf("failed to execute k8s template file %s: %v", path, err)
	}
	if err := yaml.Unmarshal(tplBuffer.Bytes(), obj); err != nil {
		return fmt.Errorf("failed to unmarshall k8s template file %s: %v", path, err)
	}
	m.setLabels()
	return nil
}

func (m *K8sManifest) forwardPodPorts(pod *coreV1.Pod) ([]portforward.ForwardedPort, error) {
	roundTripper, upgrader, err := spdy.RoundTripperFor(m.k8sConfig)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", m.namespace.Name, pod.Name)
	hostIP := strings.TrimLeft(m.k8sConfig.Host, "htps:/")
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
	if len(ports) == 0 {
		return []portforward.ForwardedPort{}, nil
	}

	log.Debug().Str("Pod", pod.Name).Interface("Ports", ports).Msg("Attempting to forward ports")

	forwarder, err := portforward.New(dialer, ports, stopChan, readyChan, out, errOut)
	if err != nil {
		return nil, err
	}
	go func() {
		if err := forwarder.ForwardPorts(); err != nil {
			log.Error().Str("Pod", pod.Name).Err(err)
		}
	}()

	log.Debug().Str("Pod", pod.Name).Msg("Waiting on pods forwarded ports to be ready")
	<-readyChan
	if len(errOut.String()) > 0 {
		return nil, fmt.Errorf("error on forwarding k8s port: %v", errOut.String())
	}
	if len(out.String()) > 0 {
		msg := strings.ReplaceAll(out.String(), "\n", " ")
		log.Debug().Str("Pod", pod.Name).Msgf("%s", msg)
	}

	m.stopChannels = append(m.stopChannels, stopChan)

	return forwarder.GetPorts()
}

// K8sManifestGroup is an implementation of K8sEnvResource that allows a group of manifests to be deployed
// concurrently on the cluster. This is important for services that don't have dependencies on each other.
// For example: a Chainlink node doesn't depend on other Chainlink nodes on deploy and an adapter doesn't depend on
// Chainlink nodes on deploy, only later on within the test lifecycle which means they can be included within a single
// group.
// Whereas, Chainlink does depend on a deployed Geth, Hardhat, Ganache on deploy so they cannot be included in the group
// as Chainlink definition needs to know the cluster IP of the deployment for it to boot.
type K8sManifestGroup struct {
	id        string
	manifests []*K8sManifest
}

// ID returns the identifier of the manifest group
func (mg *K8sManifestGroup) ID() string {
	return mg.id
}

// SetEnvironment initiates the k8s cluster and config within all the nested manifests
func (mg *K8sManifestGroup) SetEnvironment(
	k8sClient *kubernetes.Clientset,
	k8sConfig *rest.Config,
	config *config.Config,
	network *config.NetworkConfig,
	namespace *coreV1.Namespace,
) error {
	for _, m := range mg.manifests {
		if err := m.SetEnvironment(k8sClient, k8sConfig, config, network, namespace); err != nil {
			return err
		}
	}
	return nil
}

// Deploy concurrency creates all of the definitions on the k8s cluster
func (mg *K8sManifestGroup) Deploy(values map[string]interface{}) error {
	var errGroup error

	wg := mg.waitGroup()
	for _, manifest := range mg.manifests {
		m := manifest
		go func() {
			defer wg.Done()
			if err := m.Deploy(values); err != nil {
				errGroup = multierror.Append(errGroup, err)
			}
		}()
	}
	wg.Wait()
	return errGroup
}

// WaitUntilHealthy will wait until all of the manifests in the group are considered healthy.
// To avoid it duplicating checks for multiple of the same manifest in a group, it will first create a unique map
// of manifest IDs so checks aren't performed multiple times.
func (mg *K8sManifestGroup) WaitUntilHealthy() error {
	var errGroup error

	idMap := map[string]*K8sManifest{}
	for _, manifest := range mg.manifests {
		idMap[manifest.id] = manifest
	}
	wg := sync.WaitGroup{}
	wg.Add(len(idMap))

	for _, manifest := range idMap {
		m := manifest
		go func() {
			defer wg.Done()
			if err := m.WaitUntilHealthy(); err != nil {
				errGroup = multierror.Append(errGroup, err)
			}
		}()
	}
	wg.Wait()
	return errGroup
}

// ServiceDetails will return all the details of the services within a group
func (mg *K8sManifestGroup) ServiceDetails() ([]*ServiceDetails, error) {
	var serviceDetails []*ServiceDetails
	for _, m := range mg.manifests {
		if manifestServiceDetails, err := m.ServiceDetails(); err != nil {
			return nil, err
		} else {
			serviceDetails = append(serviceDetails, manifestServiceDetails...)
		}
	}
	return serviceDetails, nil
}

// Values will return all of the defined values to be exposed in the template definitions.
// Due to there possibly being multiple of the same manifest in the group, it will create keys of each manifest id,
// also keys with each manifest followed by its index. For example:
// values["adapter"].apiPort
// values["chainlink"].webPort
// values["chainlink_0"].webPort
// values["chainlink_1"].webPort
func (mg *K8sManifestGroup) Values() map[string]interface{} {
	values := map[string]interface{}{}
	for _, m := range mg.manifests {
		id := strings.Split(m.id, "-")
		if len(id) > 1 {
			values[strings.Join(id, "_")] = m.Values()
		}
		if _, ok := values[id[0]]; !ok {
			values[id[0]] = m.Values()
		}
	}
	return values
}

// Teardown will iterate through each manifest and tear it down
func (mg *K8sManifestGroup) Teardown() error {
	for _, m := range mg.manifests {
		if err := m.Teardown(); err != nil {
			return err
		}
	}
	return nil
}

func (mg *K8sManifestGroup) waitGroup() *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(len(mg.manifests))
	return &wg
}

func k8sConfig() (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	return kubeConfig.ClientConfig()
}

func waitForHealthyPods(k8sClient *kubernetes.Clientset, namespace *coreV1.Namespace, pods *coreV1.PodList) error {
	wg := sync.WaitGroup{}
	wg.Add(len(pods.Items))

	var errGroup error

	for _, pod := range pods.Items {
		p := pod.Name
		go func() {
			defer wg.Done()
			if err := waitForHealthyPod(k8sClient, namespace, p); err != nil {
				errGroup = multierror.Append(errGroup, err)
			}
		}()
	}
	wg.Wait()

	return errGroup
}

func waitForHealthyPod(k8sClient *kubernetes.Clientset, namespace *coreV1.Namespace, podName string) error {
	k8sPods := k8sClient.CoreV1().Pods(namespace.Name)

	log.Info().Str("Pod", podName).Msg("Waiting for pod to be healthy")
	ticker := time.NewTicker(time.Millisecond * 500)

	defer ticker.Stop()
ticker:
	for range ticker.C {
		pod, err := k8sPods.Get(context.Background(), podName, metaV1.GetOptions{})
		if err != nil {
			return err
		}
		for _, condition := range pod.Status.Conditions { // Each pod has an unordered list of conditions
			if condition.Type == coreV1.PodReady {
				if condition.Status == coreV1.ConditionTrue {
					log.Info().Str("Pod", podName).Msg("Pod healthy")
					break ticker
				}
			}
		}
	}
	return nil
}
