package environment

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	kubeAppTypes "k8s.io/client-go/kubernetes/typed/apps/v1"
	kubeCoreTypes "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const LatestChainlinkVersion string = "0.10.8"

type Environment interface {
	ChainlinkNodes() []client.Chainlink
	TearDown() error
}

type environment struct {
	name                  string
	kubeClient            *kubernetes.Clientset
	deploymentsClient     kubeAppTypes.DeploymentInterface
	servicesClient        kubeCoreTypes.ServiceInterface
	network               client.BlockchainNetwork
	chainlinkNodes        []client.Chainlink
	forwardedPortChannels []chan os.Signal
}

// NewBasicEnvironment launches a new environment of latest version chainlink nodes, connected to the specified network
func NewBasicEnvironment(environmentName string, nodeCount int, network client.BlockchainNetwork) (Environment, error) {
	kubeClient, err := kubeClient()
	if err != nil {
		return nil, err
	}

	namespace, err := kubeClient.CoreV1().Namespaces().Create(
		context.Background(), &apiv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: environmentName + "-",
			},
		}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("Namespace", namespace.Name).
		Int("Node Count", nodeCount).
		Msg("Deploying K8s environment")
	// Clients for each portion of our K8s cluster
	secretsClient := kubeClient.CoreV1().Secrets(namespace.Name)
	secretsClient.Create(context.Background(), chainlinkNodeSecret(), metav1.CreateOptions{})
	deploymentsClient := kubeClient.AppsV1().Deployments(namespace.Name)
	servicesClient := kubeClient.CoreV1().Services(namespace.Name)

	env := &environment{
		name:              namespace.Name,
		kubeClient:        kubeClient,
		deploymentsClient: deploymentsClient,
		servicesClient:    servicesClient,
		network:           network,
		chainlinkNodes:    []client.Chainlink{},
	}

	// Launch hardhat setup if that's what we're testing on
	if network.ID() == client.EthereumHardhatID {
		connectionString, err := env.deployHardhat()
		if err != nil {
			return nil, err
		}
		network.SetURL(connectionString)
	}

	// Launch chainlink nodes
	for i := 0; i < nodeCount; i++ {
		err = env.deployChainlink()
		if err != nil {
			return nil, err
		}
	}
	// Wait for everything to be up and healthy
	err = waitForHealthyPods(kubeClient, namespace.Name)
	if err != nil {
		return nil, err
	}
	forwardedPorts, hardhatPort, err := env.forwardPorts(namespace.Name)
	if network.ID() == client.EthereumHardhatID { // URL to connect to hardhat from outside the cluster
		network.SetURL("ws://127.0.0.1:" + hardhatPort)
	}

	for _, port := range forwardedPorts {
		if err != nil {
			return nil, err
		}

		cl, err := client.NewChainlink(&client.ChainlinkConfig{
			URL:      "http://127.0.0.1:" + port,
			Email:    "notreal@fakeemail.ch",
			Password: "twochains",
		}, http.DefaultClient)
		log.Info().Str("URL", "http://127.0.0.1:"+port).Msg("Created Chainlink connection")
		if err != nil {
			return nil, err
		}
		env.chainlinkNodes = append(env.chainlinkNodes, cl)
	}

	return env, err
}

// GetChainlinkNodes returns all the chainlink nodes in the launched environment
func (env *environment) ChainlinkNodes() []client.Chainlink {
	return env.chainlinkNodes
}

// TearDown calls delete on all the environment's resources
func (env *environment) TearDown() error {
	// Close all open forwarded ports
	for _, sigChan := range env.forwardedPortChannels {
		sigChan <- syscall.SIGTERM
		close(sigChan)
	}
	err := env.kubeClient.CoreV1().Namespaces().Delete(context.Background(), env.name, metav1.DeleteOptions{})
	log.Info().Str("Name", env.name).Msg("Deleted Environment")
	return err
}

// Deploys a chainlink pod to the cluster
func (env *environment) deployChainlink() error {
	deploymentSpec := newChainlinkDeployment(env.network, LatestChainlinkVersion)
	nodeDeployment, err := env.deploymentsClient.Create(context.Background(), deploymentSpec, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	serviceSpec := newChainlinkService(nodeDeployment.Name)
	_, err = env.servicesClient.Create(context.Background(), serviceSpec, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	log.Info().Str("Name", nodeDeployment.Name).Msg("Deployed chainlink node")
	return err
}

// Deploys a hardhat pod to the cluster
func (env *environment) deployHardhat() (string, error) {
	hardhatDeploySpec, hardhatServiceSpec := newHardhat()
	_, err := env.deploymentsClient.Create(context.Background(), hardhatDeploySpec, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	hardhatService, err := env.servicesClient.Create(context.Background(), hardhatServiceSpec, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	log.Info().Msg("Deployed Hardhat")
	return "ws://" + hardhatService.Spec.ClusterIP + ":8545", err
}

// Builds all that's needed for launching a hardhat network for testing
func newHardhat() (*appsv1.Deployment, *apiv1.Service) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "hardhat-network",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "hardhat-network"},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "hardhat-network"},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "hardhat-network",
							Image: "smartcontract/hardhat-network",
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: 8545,
								},
							},
							ReadinessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(8545),
									},
								},
								PeriodSeconds:       2,
								InitialDelaySeconds: 3,
							},
						},
					},
				},
			},
		},
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "hardhat-network",
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": "hardhat-network",
			},
			Ports: []apiv1.ServicePort{
				{
					Name:       "access",
					Port:       int32(8545),
					TargetPort: intstr.FromInt(8545),
				},
			},
		},
	}

	return deployment, service
}

// Waits for all pods in the namespace to report as running healthy
func waitForHealthyPods(kubeClient *kubernetes.Clientset, namespace string) error {
	start := time.Now()
	log.Info().Str("Name", namespace).Msg("Waiting for environment to be healthy")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			podInterface := kubeClient.CoreV1().Pods(namespace)
			pods, err := podInterface.List(context.Background(), metav1.ListOptions{})
			if err != nil {
				return err
			}
			healthyPodCount := 0
			for _, pod := range pods.Items {
				for _, condition := range pod.Status.Conditions { // Each pod has an unordered list of conditions
					if condition.Type == apiv1.PodReady {
						if condition.Status == apiv1.ConditionTrue {
							healthyPodCount += 1
							break
						}
					}
				}
			}
			if healthyPodCount == len(pods.Items) {
				log.Info().
					Str("Name", namespace).
					Str("Wait Length", time.Since(start).Round(time.Second).String()).
					Msg("Environment healthy")
				return nil
			}
		}
	}
}

// Creates K8s deployment object for chainlink node setup
func newChainlinkDeployment(network client.BlockchainNetwork, chainlinkVersion string) *appsv1.Deployment {
	chainID := network.ChainID()
	chainUrl := network.URL()

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "chainlink-node-",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "chainlink-node",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "chainlink-node",
					},
				},
				Spec: apiv1.PodSpec{
					Volumes: []apiv1.Volume{
						{
							Name: "node-secrets-volume",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: "node-secrets",
								},
							},
						},
					},
					Containers: []apiv1.Container{
						{
							Name:  "node",
							Image: "smartcontract/chainlink:" + chainlinkVersion,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "access",
									ContainerPort: 6688,
								}, {
									Name:          "node",
									ContainerPort: 6060,
								},
							},
							Env: defaultChainlinkEnvVars(chainUrl, chainID.String()),
							Args: []string{
								"node",
								"start",
								"-d",
								"-p",
								"/etc/node-secrets-volume/node-password",
								"-a",
								"/etc/node-secrets-volume/apicredentials",
							},
							LivenessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(6688),
									},
								},
								PeriodSeconds:       10,
								InitialDelaySeconds: 90,
							},
							ReadinessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(6688),
									},
								},
								PeriodSeconds:       5,
								InitialDelaySeconds: 2,
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "node-secrets-volume",
									MountPath: "/etc/node-secrets-volume/",
								},
							},
						}, {
							Name:  "db",
							Image: "postgres:11.6",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "postgres",
									ContainerPort: 5432,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "POSTGRES_DB",
									Value: "chainlink",
								}, {
									Name:  "POSTGRES_PASSWORD",
									Value: "node",
								}, {
									Name:  "PGPASSWORD",
									Value: "node",
								},
							},
							LivenessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									Exec: &apiv1.ExecAction{
										Command: []string{"pg_isready", "-U", "postgres"},
									},
								},
								PeriodSeconds:       60,
								InitialDelaySeconds: 60,
							},
							ReadinessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									Exec: &apiv1.ExecAction{
										Command: []string{"pg_isready", "-U", "postgres"},
									},
								},
								PeriodSeconds:       5,
								InitialDelaySeconds: 2,
							},
						},
					},
				},
			},
		},
	}

	return deployment
}

// Creates K8s service object for chainlink node setup
func newChainlinkService(deploymentName string) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			Ports: []apiv1.ServicePort{
				{
					Name:       "node-port",
					Port:       int32(6688),
					TargetPort: intstr.FromInt(6688),
				},
			},
			Selector: map[string]string{
				"app": "chainlink-node",
			},
		},
	}
}

// Creates K8s secret object for chainlink node setup
func chainlinkNodeSecret() *apiv1.Secret {
	return &apiv1.Secret{
		Type: apiv1.SecretType("Opaque"),
		Data: map[string][]byte{
			"apicredentials": []byte("notreal@fakeemail.ch\ntwochains"),
			"node-password":  []byte("T.tLHkcmwePT/p,]sYuntjwHKAsrhm#4eRs4LuKHwvHejWYAC2JP4M8HimwgmbaZ"),
			// This bit can probably be removed, waiting on full confirmation
			"0xb90c7E3F7815F59EAD74e7543eB6D9E8538455D6.json": []byte(`{
"address": "b90c7e3f7815f59ead74e7543eb6d9e8538455d6",
"crypto": {
	"cipher": "aes-128-ctr",
	"ciphertext": "e83fe14bcf9197de06d84800c1a76db3945da0e323ec6357d6495581f693b43f",
	"cipherparams": { "iv": "4965208fc86af075261bcea2940f3988" },
	"kdf": "scrypt",
	"kdfparams": {
	"dklen": 32,
	"n": 262144,
	"p": 1,
	"r": 8,
	"salt": "cc07e486400e4b8c86db9b142aeff9151ba214fc1b15cacb3925829e20f6443f"
	},
	"mac": "cab6f449ac715b59f7be31ffe96f9f712e3fb442e0cde619d9cddbe44fa44119"
},
"id": "bf6687ea-3758-4130-843c-b1d16c1be38b",
"version": 3
}`),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-secrets",
		},
	}
}

// Builds the base kubernetes client
func kubeClient() (*kubernetes.Clientset, error) {
	config, err := kubeConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// Uses default loading rules to load the kubernetest config
func kubeConfig() (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	return kubeConfig.ClientConfig()
}

// Forwards all ports needed to connect to the environment, along with the hardhat port
func (env *environment) forwardPorts(namespaceName string) ([]string, string, error) {
	kubeClient, err := kubeClient()
	if err != nil {
		return nil, "", err
	}
	kubeConfig, err := kubeConfig()
	if err != nil {
		return nil, "", err
	}

	// TODO: Figure out a way to do this for services, the obvious way gives inscrutable errors
	// https://gianarb.it/blog/programmatically-kube-port-forward-in-go
	// https://github.com/gianarb/kube-port-forward/issues/3
	podInterface := kubeClient.CoreV1().Pods(namespaceName)
	podList, err := podInterface.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, "", err
	}

	log.Info().Msg("Forwarding Ports")
	forwardedPorts := []string{}
	hardhatPort := ""
	for _, pod := range podList.Items {
		// stopCh control the port forwarding lifecycle. When it gets closed the
		// port forward will terminate
		stopCh := make(chan struct{}, 1)
		// readyCh communicate when the port forward is ready to get traffic
		readyCh := make(chan struct{})
		// stream is used to tell the port forwarder where to place its output or
		// where to expect input if needed. For the port forwarding we just need
		// the output eventually
		stream := genericclioptions.IOStreams{
			In:     os.Stdin,
			Out:    nil, // Changing this the os.Stdout can help with debugging
			ErrOut: os.Stderr,
		}
		// receive the forwarded port when it is generated and ready
		portCh := make(chan string, 1)

		// managing termination signal from the terminal. As you can see the stopCh
		// gets closed to gracefully handle its termination.
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			close(stopCh)
		}()
		env.forwardedPortChannels = append(env.forwardedPortChannels, sigs)

		go func() {
			err := env.forwardPort(portForwardRequest{
				RestConfig: kubeConfig,
				Pod:        pod,
				Streams:    stream,
				StopCh:     stopCh,
				ReadyCh:    readyCh,
				PortCh:     portCh,
			})
			if err != nil {
				log.Err(err).Str("Pod", pod.Name).Msg("Error while forwarding port, tearing down env")
				env.TearDown()
				panic(err)
			}
		}()

		select {
		case forwardedPort := <-portCh:
			if strings.HasPrefix(pod.Name, "hardhat-network") {
				hardhatPort = forwardedPort
			} else {
				forwardedPorts = append(forwardedPorts, forwardedPort)
			}
			log.Info().Str("Pod", pod.Name).Str("Port", forwardedPort).Msg("Forwarded local port")
			break
		}
	}
	return forwardedPorts, hardhatPort, err
}

// Starts a go routine needed to actively forward a port
func (env *environment) forwardPort(req portForwardRequest) error {
	portForwardRequest := []string{"0:6688"} // Default for chainlink nodes
	if strings.HasPrefix(req.Pod.Name, "hardhat-network") {
		portForwardRequest = []string{"0:8545"}
	}
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", req.Pod.Namespace, req.Pod.Name)
	hostIP := strings.TrimLeft(req.RestConfig.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(req.RestConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, portForwardRequest, req.StopCh, req.ReadyCh, req.Streams.Out, req.Streams.ErrOut)
	if err != nil {
		return err
	}
	go func() {
		err = fw.ForwardPorts()
		if err != nil {
			log.Err(err).Str("Pod", req.Pod.Name).Msg("Error while forwarding port, tearing down env")
			env.TearDown()
			panic(err)
		}
	}()
	select {
	case <-req.ReadyCh:
		break
	}
	forwarded, err := fw.GetPorts()
	if err != nil {
		return err
	}
	req.PortCh <- fmt.Sprint(forwarded[0].Local)
	return err
}

// All info needed to forward a K8s port
type portForwardRequest struct {
	// RestConfig is the kubernetes config
	RestConfig *rest.Config
	// Pod is the selected pod for this port forwarding
	Pod apiv1.Pod
	// Steams configures where to write or read input from
	Streams genericclioptions.IOStreams
	// StopCh is the channel used to manage the port forward lifecycle
	StopCh <-chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	ReadyCh chan struct{}
	// PortCh receives from when the ready channel is done, and contains the local port that has been forwarded
	PortCh chan string
}

func defaultChainlinkEnvVars(ethUrl, ethChainID string) []apiv1.EnvVar {
	return []apiv1.EnvVar{
		{
			Name:  "DATABASE_URL",
			Value: "postgresql://postgres:node@127.0.0.1:5432/chainlink?sslmode=disable",
		}, {
			Name:  "DATABASE_NAME",
			Value: "chainlink",
		}, {
			Name:  "ETH_URL",
			Value: ethUrl,
		}, {
			Name:  "ETH_CHAIN_ID",
			Value: ethChainID,
		}, {
			Name:  "ALLOW_ORIGINS",
			Value: "*",
		}, {
			Name:  "CHAINLINK_DEV",
			Value: "true",
		}, {
			Name:  "CHAINLINK_PGPASSWORD",
			Value: "node",
		}, {
			Name:  "CHAINLINK_PORT",
			Value: "6688",
		}, {
			Name:  "CHAINLINK_TLS_PORT",
			Value: "0",
		}, {
			Name:  "DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS",
			Value: "true",
		}, {
			Name:  "ENABLE_BULLETPROOF_TX_MANAGER",
			Value: "true",
		}, {
			Name:  "FEATURE_OFFCHAIN_REPORTING",
			Value: "true",
		}, {
			Name:  "JSON_CONSOLE",
			Value: "false",
		}, {
			Name:  "LOG_LEVEL",
			Value: "info",
		}, {
			Name:  "MAX_EXPORT_HTML_THREADS",
			Value: "2",
		}, {
			Name:  "MINIMUM_CONTRACT_PAYMENT",
			Value: "0",
		}, {
			Name:  "OCR_TRACE_LOGGING",
			Value: "true",
		}, {
			Name:  "P2P_LISTEN_IP",
			Value: "0.0.0.0",
		}, {
			Name:  "P2P_LISTEN_PORT",
			Value: "6690",
		}, {
			Name:  "ROOT",
			Value: "./clroot",
		}, {
			Name:  "SECURE_COOKIES",
			Value: "false",
		},
	}
}
