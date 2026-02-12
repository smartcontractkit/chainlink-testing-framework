package pods

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const (
	// ClientReadyTimeout is a timeout client (tests) will wait until abandoning attempts
	ClientReadyTimeout = 1 * time.Minute
	// RetryDelay is a delay before retrying forwarding
	RetryDelay = 1 * time.Second
	// K8sFunctionCallTimeout is a common K8s API timeout we use in functions
	// function may contain multiple calls
	K8sFunctionCallTimeout = 2 * time.Minute
)

// PortForwardConfig represents a single port forward configuration
type PortForwardConfig struct {
	ServiceName string
	LocalPort   int
	ServicePort int
	Namespace   string
}

func (c PortForwardConfig) validate() error {
	if c.Namespace == "" {
		return fmt.Errorf("empty K8s namespace")
	}
	if c.LocalPort == 0 {
		return fmt.Errorf("empty local port")
	}
	if c.ServicePort == 0 {
		return fmt.Errorf("empty service port")
	}
	if c.ServiceName == "" {
		return fmt.Errorf("empty service name")
	}
	return nil
}

// PortForwardManager manages multiple port forwards
type PortForwardManager struct {
	cs     *kubernetes.Clientset
	config *rest.Config
	// key in format serviceName:localPort
	forwards map[string]*forwardInfo
	mu       sync.Mutex
}

// forwardInfo holds information about a running port forward
type forwardInfo struct {
	stopChan chan struct{}
	// signalClientReadyChan is used to signal the caller when first connection is established
	signalClientReadyChan chan struct{}
	cleanup               func()
}

// NewForwarder creates a new manager for multiple port forwards
func NewForwarder(api *API) *PortForwardManager {
	return &PortForwardManager{
		cs:       api.ClientSet,
		config:   api.RESTConfig,
		forwards: make(map[string]*forwardInfo),
	}
}

// startForwardService starts forwarding a single service port with retry logic
func (m *PortForwardManager) startForwardService(cfg PortForwardConfig) {
	key := fmt.Sprintf("%s:%d", cfg.ServiceName, cfg.LocalPort)
	stopChan, readyChan := make(chan struct{}), make(chan struct{})

	m.mu.Lock()
	m.forwards[key] = &forwardInfo{
		stopChan:              stopChan,
		signalClientReadyChan: readyChan,
		cleanup:               func() { close(stopChan) },
	}
	m.mu.Unlock()

	go m.forwardAndRetry(cfg, stopChan, readyChan)
}

// forwardAndRetry continuously attempts to forward the port with retries
func (m *PortForwardManager) forwardAndRetry(cfg PortForwardConfig, stopChan <-chan struct{}, readyChan chan struct{}) {
	key := fmt.Sprintf("%s:%d", cfg.ServiceName, cfg.LocalPort)
	consecutiveFailures := 0

	for {
		select {
		case <-stopChan:
			L.Info().Msgf("Stopped retry loop for %s", key)
			return
		default:
			L.Info().
				Str("ServiceName", cfg.ServiceName).
				Int("LocalPort", cfg.LocalPort).
				Msg("Starting forwarder")
			err := m.attemptForward(cfg, readyChan)

			if err != nil {
				// Connection failed or broke - retry
				consecutiveFailures++
				L.Debug().
					Err(err).
					Str("Key", key).
					Int("Attempt", consecutiveFailures).
					Msg("Port forward failed")
				L.Debug().Msgf("Retrying %s in %v", key, RetryDelay)
				select {
				case <-stopChan:
					return
				case <-time.After(RetryDelay):
					continue
				}
			} else {
				// attemptForward returned nil = clean stop requested
				L.Info().Msgf("Port forward %s stopped cleanly", key)
				return
			}
		}
	}
}

// attemptForward establishes and monitors a single port forward connection
func (m *PortForwardManager) attemptForward(cfg PortForwardConfig, signalReadyChan chan struct{}) error {
	namespace := cfg.Namespace
	if namespace == "" {
		namespace = "default"
	}

	// Get target pod for forwarding
	targetPod, targetPort, err := m.getTargetPodAndPort(cfg, namespace)
	if err != nil {
		return fmt.Errorf("failed to get target: %w", err)
	}

	l := L.With().
		Str("ServiceName", cfg.ServiceName).
		Int("LocalPort", cfg.LocalPort).Logger()

	l.Info().Msgf("Forwarding service %s:%d -> pod %s:%d -> localhost:%d",
		cfg.ServiceName, cfg.ServicePort, targetPod.Name, targetPort, cfg.LocalPort)

	// run port forward
	stopChan := make(chan struct{})
	readyChan := make(chan struct{})
	errChan := make(chan error, 1)

	restClient := m.cs.CoreV1().RESTClient()
	url := restClient.Post().
		Resource("pods").
		Namespace(namespace).
		Name(targetPod.Name).
		SubResource("portforward").
		URL()

	transport, upgrader, err := spdy.RoundTripperFor(m.config)
	if err != nil {
		return fmt.Errorf("failed to create SPDY transport: %w", err)
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", url)
	ports := []string{fmt.Sprintf("%d:%d", cfg.LocalPort, targetPort)}

	pf, err := portforward.New(dialer, ports, stopChan, readyChan, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create port forwarder: %w", err)
	}

	go func() {
		errChan <- pf.ForwardPorts()
	}()

	// monitor the connection
	select {
	case <-readyChan:
		l.Info().
			Msg("🟢 established connection for")
		signalReadyChan <- struct{}{}
		// error is received when someone is trying to call the local port
		// and connection is broken, it is not proactively checked
		select {
		case err := <-errChan:
			// Connection terminated (error or network issue)
			close(stopChan)
			if err != nil {
				l.Warn().
					Err(err).Msg("🟡 lost connection")
			}
			return err
		case <-stopChan:
			l.Info().Msg("🔴 Stop method was called, stopping forwarder")
			return nil
		}

	case err := <-errChan:
		close(stopChan)
		l.Error().
			Err(err).Msg("🔴 Port forward FAILED to establish for %s:%d")
		return fmt.Errorf("failed to establish: %w", err)

	case <-time.After(10 * time.Second):
		close(stopChan)
		l.Error().Msg("⏱️Port forward timeout")
		return fmt.Errorf("timeout establishing connection")
	}
}

// getTargetPodAndPort finds the target pod and resolves the port
func (m *PortForwardManager) getTargetPodAndPort(cfg PortForwardConfig, namespace string) (*corev1.Pod, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), K8sFunctionCallTimeout)
	defer cancel()
	service, err := m.cs.CoreV1().Services(namespace).Get(ctx, cfg.ServiceName, metav1.GetOptions{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get service: %w", err)
	}

	if len(service.Spec.Selector) == 0 {
		return nil, 0, fmt.Errorf("service %s has no selector", cfg.ServiceName)
	}
	pods, err := m.cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.FormatLabels(service.Spec.Selector),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find pods: %w", err)
	}

	if len(pods.Items) == 0 {
		return nil, 0, fmt.Errorf("no pods found for service %s", cfg.ServiceName)
	}

	// Find running pod
	var targetPod *corev1.Pod
	for i := range pods.Items {
		if pods.Items[i].Status.Phase == corev1.PodRunning {
			targetPod = &pods.Items[i]
			break
		}
	}
	if targetPod == nil {
		return nil, 0, fmt.Errorf("no running pods found for service %s", cfg.ServiceName)
	}
	targetPort, err := m.resolveServicePort(service, cfg.ServicePort, targetPod)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to resolve port: %w", err)
	}
	return targetPod, targetPort, nil
}

// Forward starts multiple port forwards concurrently
func (m *PortForwardManager) Forward(configs []PortForwardConfig) error {
	for _, cfg := range configs {
		if err := cfg.validate(); err != nil {
			return err
		}
		m.startForwardService(cfg)
	}
	// this code is used so library users can block until first successful connection
	eg := &errgroup.Group{}
	for _, fwd := range m.forwards {
		eg.Go(func() error {
			L.Info().Msg("Awaiting for first established connection")
			select {
			case <-fwd.signalClientReadyChan:
				L.Info().Msg("Connection established")
				return nil
			case <-time.After(2 * time.Minute):
				return fmt.Errorf("failed to forward ports until deadline")
			}
		})
	}
	return eg.Wait()
}

// StopAll stops all active port forwards
func (m *PortForwardManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, info := range m.forwards {
		info.cleanup()
		delete(m.forwards, key)
		L.Info().Msgf("Stopped port forwarding: %s", key)
	}
}

// Stop stops a specific port forward
func (m *PortForwardManager) Stop(serviceName string, localPort int) {
	key := fmt.Sprintf("%s:%d", serviceName, localPort)

	m.mu.Lock()
	defer m.mu.Unlock()

	if info, exists := m.forwards[key]; exists {
		info.cleanup()
		delete(m.forwards, key)
		L.Info().Msgf("Stopped port forward: %s", key)
	}
}

// List returns all active port forwards
func (m *PortForwardManager) List() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []string
	for key := range m.forwards {
		result = append(result, key)
	}
	return result
}

// resolveServicePort resolves service port to container port
func (m *PortForwardManager) resolveServicePort(service *corev1.Service, servicePort int, pod *corev1.Pod) (int, error) {
	for _, port := range service.Spec.Ports {
		if int(port.Port) == servicePort {
			if port.TargetPort.IntValue() != 0 {
				return port.TargetPort.IntValue(), nil
			}
			if port.TargetPort.StrVal != "" {
				return m.findNamedPort(pod, port.TargetPort.StrVal)
			}
			return servicePort, nil
		}
	}
	return 0, fmt.Errorf("service port %d not found", servicePort)
}

// findNamedPort finds a named port in a pod's containers
func (m *PortForwardManager) findNamedPort(pod *corev1.Pod, portName string) (int, error) {
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			if port.Name == portName {
				return int(port.ContainerPort), nil
			}
		}
	}
	return 0, fmt.Errorf("named port %s not found in pod %s", portName, pod.Name)
}
