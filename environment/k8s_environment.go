package environment

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/integrations-framework/chaos"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/ghodss/yaml"

	"github.com/avast/retry-go"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/transport/spdy"
)

const SelectorLabelKey string = "app"

// K8sEnvSpecs represents a series of environment resources to be deployed. The resources in the array will be
// deployed in the order that they are present in the array.
type K8sEnvSpecs []K8sEnvResource

// K8sEnvSpecInit is the initiator that will return the name of the environment and the specifications to be deployed.
// The name of the environment returned determines the namespace.
type K8sEnvSpecInit func(*config.NetworkConfig) (string, K8sEnvSpecs)

// K8sEnvResource is the interface for deploying a given environment resource. Creating an interface for resource
// deployment allows it to be extended, deploying k8s resources in different ways. For example: K8sManifest deploys
// a single manifest, whereas K8sManifestGroup bundles several K8sManifests to be deployed concurrently.
type K8sEnvResource interface {
	ID() string
	GetConfig() *config.Config
	SetEnvironment(environment *K8sEnvironment) error
	Deploy(values map[string]interface{}) error
	WaitUntilHealthy() error
	ServiceDetails() ([]*ServiceDetails, error)
	SetValue(key string, val interface{})
	Values() map[string]interface{}
	Teardown() error
}

type K8sEnvironment struct {
	// K8s Resources
	k8sClient *kubernetes.Clientset
	k8sConfig *rest.Config

	// Deployment resources
	specs     K8sEnvSpecs
	namespace *coreV1.Namespace

	// Environment resources
	config  *config.Config
	network client.BlockchainNetwork
	chaos   *chaos.Controller
}

// NewK8sEnvironment creates and deploys a full ephemeral environment in a k8s cluster. Your current context within
// your kube config will always be used.
func NewK8sEnvironment(
	init K8sEnvSpecInit,
	cfg *config.Config,
	network client.BlockchainNetwork,
) (Environment, error) {
	k8sConfig, err := K8sConfig()
	if err != nil {
		return nil, err
	}
	k8sConfig.QPS = cfg.Kubernetes.QPS
	k8sConfig.Burst = cfg.Kubernetes.Burst

	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	env := &K8sEnvironment{
		k8sClient: k8sClient,
		k8sConfig: k8sConfig,
		config:    cfg,
		network:   network,
	}
	log.Info().Str("Host", k8sConfig.Host).Msg("Using Kubernetes cluster")

	environmentName, deployables := init(network.Config())
	namespace, err := env.createNamespace(environmentName)
	if err != nil {
		return nil, err
	}
	env.namespace = namespace
	env.specs = deployables
	cc, err := chaos.NewController(&chaos.Config{
		Client:    k8sClient,
		Namespace: namespace.Name,
	})
	if err != nil {
		return nil, err
	}
	env.chaos = cc

	ctx, ctxCancel := context.WithTimeout(context.Background(), env.config.Kubernetes.DeploymentTimeout)
	defer ctxCancel()

	errChan := make(chan error)
	go env.deploySpecs(errChan)

deploymentLoop:
	for {
		select {
		case err, open := <-errChan:
			if err != nil {
				return nil, err
			} else if !open {
				break deploymentLoop
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("error while waiting for deployment: %v", ctx.Err())
		}
	}
	return env, err
}

// ID returns the canonical name of the environment, which in the case of k8s is the namespace
func (env K8sEnvironment) ID() string {
	if env.namespace != nil {
		return env.namespace.Name
	}
	return ""
}

// GetAllServiceDetails returns all the connectivity details for a deployed service by its remote port within k8s
func (env K8sEnvironment) GetAllServiceDetails(remotePort uint16) ([]*ServiceDetails, error) {
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
func (env K8sEnvironment) GetServiceDetails(remotePort uint16) (*ServiceDetails, error) {
	if serviceDetails, err := env.GetAllServiceDetails(remotePort); err != nil {
		return nil, err
	} else {
		return serviceDetails[0], err
	}
}

// WriteArtifacts dumps pod logs and DB info within the environment into local log files,
// used near exclusively on test failure
func (env K8sEnvironment) WriteArtifacts(testLogFolder string) {
	// Get logs from K8s pods
	podsClient := env.k8sClient.CoreV1().Pods(env.namespace.Name)
	podsList, err := podsClient.List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		log.Err(err).Str("Env Name", env.namespace.Name).Msg("Error retrieving pod list from K8s environment")
	}

	// Each Pod gets a folder
	for _, pod := range podsList.Items {
		podName := pod.Labels[SelectorLabelKey]
		podFolder := filepath.Join(testLogFolder, podName)
		if _, err := os.Stat(podFolder); os.IsNotExist(err) {
			if err = os.Mkdir(podFolder, 0755); err != nil {
				log.Err(err).Str("Folder Name", podFolder).Msg("Error creating logs directory")
			}
		}
		err = env.writeDatabaseContents(pod, podFolder)
		if err != nil {
			log.Err(err).Str("Namespace", env.ID()).Str("Pod", pod.Name).Msg("Error fetching DB contents for pod")
		}
		err = writeLogsForPod(podsClient, pod, podFolder)
		if err != nil {
			log.Err(err).Str("Namespace", env.ID()).Str("Pod", pod.Name).Msg("Error writing logs for pod")
		}
	}
}

// ApplyChaos applies chaos experiment in the environment namespace
func (env K8sEnvironment) ApplyChaos(exp chaos.Experimentable) (string, error) {
	name, err := env.chaos.Run(exp)
	if err != nil {
		return name, err
	}
	return name, nil
}

// StopChaos stops experiment by name
func (env K8sEnvironment) StopChaos(name string) error {
	if err := env.chaos.Stop(name); err != nil {
		return err
	}
	return nil
}

// StopAllChaos stops all chaos experiments
func (env K8sEnvironment) StopAllChaos() error {
	if err := env.chaos.StopAll(); err != nil {
		return err
	}
	return nil
}

// TearDown cycles through all the specifications and tears down the deployments. This typically entails cleaning
// up port forwarding requests and deleting the namespace that then destroys all definitions.
func (env K8sEnvironment) TearDown() {
	for _, spec := range env.specs {
		if err := spec.Teardown(); err != nil {
			log.Error().Err(err)
		}
	}

	err := env.k8sClient.CoreV1().Namespaces().Delete(context.Background(), env.namespace.Name, metaV1.DeleteOptions{})
	log.Info().Str("Namespace", env.namespace.Name).Msg("Deleted env")
	log.Error().Err(err)
}

// Collects the contents of DB containers and writes them to local log files
func (env *K8sEnvironment) writeDatabaseContents(pod coreV1.Pod, podFolder string) error {
	for _, container := range pod.Spec.Containers {
		if strings.Contains(container.Image, "postgres") { // If there's a postgres image, dump its DB
			dumpContents, err := env.dumpDB(pod, container)
			if err != nil {
				return err
			}

			// Write pg_dump
			logFile, err := os.Create(filepath.Join(podFolder, fmt.Sprintf("%s_dump.sql", container.Name)))
			if err != nil {
				return err
			}
			_, err = logFile.WriteString(dumpContents)
			if err != nil {
				return err
			}

			if err = logFile.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Dumps db contents to a log file
func (env *K8sEnvironment) dumpDB(pod coreV1.Pod, container coreV1.Container) (string, error) {
	postRequestBase := env.k8sClient.CoreV1().RESTClient().Post().
		Namespace(pod.Namespace).Resource("pods").Name(pod.Name).SubResource("exec")
	exportDBRequest := postRequestBase.VersionedParams(
		&coreV1.PodExecOptions{
			Container: container.Name,
			Command:   []string{"/bin/sh", "-c", "pg_dump", "chainlink"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(env.k8sConfig, "POST", exportDBRequest.URL())
	if err != nil {
		return "", err
	}
	outBuff, errBuff := &bytes.Buffer{}, &bytes.Buffer{}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  &bytes.Reader{},
		Stdout: outBuff,
		Stderr: errBuff,
		Tty:    false,
	})
	if err != nil || errBuff.Len() > 0 {
		return "", fmt.Errorf("Error in dumping DB contents | STDOUT: %v | STDERR: %v", outBuff.String(),
			errBuff.String())
	}
	return outBuff.String(), err
}

func (env K8sEnvironment) GetPrivateKeyFromSecret(namespace string, privateKey string) (string, error) {
	res, err := env.k8sClient.CoreV1().Secrets(namespace).Get(context.Background(), "private-keys", metaV1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(res.Data[privateKey]), nil
}

// Writes logs for each container in a pod
func writeLogsForPod(podsClient v1.PodInterface, pod coreV1.Pod, podFolder string) error {
	for _, container := range pod.Spec.Containers {
		logFile, err := os.Create(filepath.Join(podFolder, container.Name) + ".log")
		if err != nil {
			return err
		}

		podLogRequest := podsClient.GetLogs(pod.Name, &coreV1.PodLogOptions{Container: container.Name})
		podLogs, err := podLogRequest.Stream(context.Background())
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			return err
		}
		_, err = logFile.Write(buf.Bytes())
		if err != nil {
			return err
		}

		if err = logFile.Close(); err != nil {
			return err
		}
		if err = podLogs.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (env *K8sEnvironment) deploySpecs(errChan chan<- error) {
	values := map[string]interface{}{}
	for i := 0; i < len(env.specs); i++ {
		spec := env.specs[i]
		if err := spec.SetEnvironment(env); err != nil {
			errChan <- err
			return
		}
		values[spec.ID()] = spec.Values()
		if err := spec.Deploy(values); err != nil {
			errChan <- err
			return
		}
		if err := spec.WaitUntilHealthy(); err != nil {
			errChan <- err
			return
		}
		values[spec.ID()] = spec.Values()
	}
	close(errChan)
}

func (env *K8sEnvironment) createNamespace(namespace string) (*coreV1.Namespace, error) {
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
type manifestGroupSetValuesFunc func(group *K8sManifestGroup) error

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
	pods         []PodForwardedInfo

	// Environment
	env *K8sEnvironment
}

// SetEnvironment is the K8sEnvResource implementation that sets the current cluster and config to be used on deploy
func (m *K8sManifest) SetEnvironment(
	environment *K8sEnvironment,
) error {
	m.env = environment
	return nil
}

// ID returns the identifier for the manifest. The ID is important as the manifest will automatically add labels and
// service selectors to link deployments to their manifests.
func (m *K8sManifest) ID() string {
	return m.id
}

func (m *K8sManifest) SetValue(key string, val interface{}) {
	m.values[key] = val
}

func (m *K8sManifest) GetConfig() *config.Config {
	return m.env.config
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
	k8sPods := m.env.k8sClient.CoreV1().Pods(m.env.namespace.Name)

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

	if err := waitForHealthyPods(m.env.k8sClient, m.env.namespace, pods); err != nil {
		return err
	}

	for _, p := range pods.Items {
		ports, err := forwardPodPorts(&p, m.env.k8sConfig, m.env.namespace.Name, m.stopChannels)
		if err != nil {
			return fmt.Errorf("unable to forward ports: %v", err)
		}
		m.ports = append(m.ports, ports...)
		log.Info().Str("Manifest ID", m.id).Interface("Ports", ports).Msg("Forwarded ports")
		m.pods = append(m.pods, PodForwardedInfo{
			PodIP:          p.Status.PodIP,
			ForwardedPorts: ports,
			PodName:        p.Name,
		})
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
			RemoteURL: remoteURL,
			LocalURL:  localURL,
		})
	}
	return serviceDetails, nil
}

// ExecuteInPod is similar to kubectl exec
func (m *K8sManifest) ExecuteInPod(podName string, containerName string, command []string) ([]byte, []byte, error) {
	set := labels.Set(m.Service.Spec.Selector)
	listOptions := metaV1.ListOptions{LabelSelector: set.AsSelector().String()}

	v1Interface := m.env.k8sClient.CoreV1()
	pods, err := v1Interface.Pods(m.env.namespace.Name).List(context.Background(), listOptions)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	var filteredPods []string

	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, podName) {
			filteredPods = append(filteredPods, pod.Name)
		}
	}

	pod, err := v1Interface.Pods(m.env.namespace.Name).Get(context.Background(), filteredPods[0], metaV1.GetOptions{})
	if err != nil {
		return []byte{}, []byte{}, err
	}

	req := m.env.k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")
	req.VersionedParams(&coreV1.PodExecOptions{
		Container: containerName,
		Command:   command,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(m.env.k8sConfig, "POST", req.URL())
	if err != nil {
		return []byte{}, []byte{}, err
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return []byte{}, []byte{}, err
	}
	return stdout.Bytes(), stderr.Bytes(), nil
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
	k8sSecrets := m.env.k8sClient.CoreV1().Secrets(m.env.namespace.Name)

	if err := m.parseSecret(m.env.config, m.env.network.Config(), values); err != nil {
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

var deploymentMutex sync.Mutex

func (m *K8sManifest) createDeployment(values map[string]interface{}) error {
	k8sDeployments := m.env.k8sClient.AppsV1().Deployments(m.env.namespace.Name)

	if image, ok := m.Values()["image"]; ok {
		deploymentMutex.Lock()
		m.env.config.Apps.Chainlink.Image = image.(string)
		m.env.config.Apps.Chainlink.Version = m.Values()["version"].(string)
	}
	if err := m.parseDeployment(m.env.config, m.env.network.Config(), values); err != nil {
		return err
	}
	if _, ok := m.Values()["image"]; ok {
		deploymentMutex.Unlock()
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
	k8sServices := m.env.k8sClient.CoreV1().Services(m.env.namespace.Name)

	if err := m.parseService(m.env.config, m.env.network.Config(), values); err != nil {
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
	cm := m.env.k8sClient.CoreV1().ConfigMaps(m.env.namespace.Name)
	if err := m.parseConfigMap(m.env.config, m.env.network.Config(), values); err != nil {
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

// TemplateValuesArray is used in the next template go func
// It's goal is to store an array of objects
// The only function it has, next, returns the first object from they array,
// and then removing that object from the array
type TemplateValuesArray struct {
	Values []interface{}
	mu     sync.Mutex
}

func (t *TemplateValuesArray) next() (interface{}, error) {
	if len(t.Values) > 0 {
		valueToReturn := t.Values[0]
		t.Values = t.Values[1:]
		return valueToReturn, nil
	} else {
		return nil, errors.New("No more Values in the array")
	}
}

func next(array *TemplateValuesArray) (interface{}, error) {
	array.mu.Lock()
	val, err := array.next()
	array.mu.Unlock()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func present(name string, data map[string]interface{}) bool {
	_, ok := data[name]
	return ok
}

func (m *K8sManifest) parse(path string, obj interface{}, data interface{}) error {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read k8s file: %v", err)
	}

	var funcs = template.FuncMap{"next": next, "present": present}

	tpl, err := template.New(path).Funcs(funcs).Parse(string(fileBytes))
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

func forwardPodPorts(pod *coreV1.Pod, k8sConfig *rest.Config, nsName string, stopChans []chan struct{}) ([]portforward.ForwardedPort, error) {
	roundTripper, upgrader, err := spdy.RoundTripperFor(k8sConfig)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", nsName, pod.Name)
	hostIP := strings.TrimLeft(k8sConfig.Host, "htps:/")
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

	//nolint
	stopChans = append(stopChans, stopChan)

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
	id            string
	manifests     []K8sEnvResource
	SetValuesFunc manifestGroupSetValuesFunc
	values        map[string]interface{}
}

// ID returns the identifier of the manifest group
func (mg *K8sManifestGroup) ID() string {
	return mg.id
}

func (mg *K8sManifestGroup) SetValue(key string, val interface{}) {
}

func (mg *K8sManifestGroup) GetConfig() *config.Config {
	return nil
}

// SetEnvironment initiates the k8s cluster and config within all the nested manifests
func (mg *K8sManifestGroup) SetEnvironment(env *K8sEnvironment) error {
	for _, m := range mg.manifests {
		if err := m.SetEnvironment(env); err != nil {
			return err
		}
	}
	return nil
}

// Deploy concurrently creates all of the definitions on the k8s cluster
func (mg *K8sManifestGroup) Deploy(values map[string]interface{}) error {
	var errGroup error
	wg := mg.waitGroup()

	originalImage := mg.manifests[0].GetConfig().Apps.Chainlink.Image
	originalVersion := mg.manifests[0].GetConfig().Apps.Chainlink.Version
	// Deploy manifests
	for i := 0; i < len(mg.manifests); i++ {
		m := mg.manifests[i]
		if manifestImage, ok := m.Values()["image"]; ok { // Check if manifest has specified image
			if manifestImage == "" { // Blank means the default from the config file
				m.SetValue("image", originalImage)
				m.SetValue("version", originalVersion)
			}
		}

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

	idMap := map[string]K8sEnvResource{}
	for _, manifest := range mg.manifests {
		idMap[manifest.ID()] = manifest
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

	if mg.SetValuesFunc != nil {
		err := mg.setValues()
		if err != nil {
			errGroup = multierror.Append(errGroup, err)
		}
	}

	return errGroup
}

func (mg *K8sManifestGroup) setValues() error {
	if mg.values == nil {
		mg.values = map[string]interface{}{}
	}
	if err := mg.SetValuesFunc(mg); err != nil {
		return err
	}
	return nil
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

	for key, value := range mg.values {
		values[key] = value
	}

	for _, m := range mg.manifests {
		id := strings.Split(m.ID(), "-")
		if len(id) > 1 {
			values[strings.Join(id, "_")] = m.Values()
		}
		if _, ok := mg.values[id[0]]; !ok {
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

// K8sConfig loads new default k8s config from filesystem
func K8sConfig() (*rest.Config, error) {
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
