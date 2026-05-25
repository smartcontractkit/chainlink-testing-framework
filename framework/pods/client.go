package pods

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// API is a struct that provides methods to interact with Kubernetes clusters.
type API struct {
	ClientSet  *kubernetes.Clientset
	RESTConfig *rest.Config
	namespace  string
}

// NewAPI creates a new instance of K8s API.
// It takes the kubeconfig path and namespace as parameters.
func NewAPI(namespace string) (*API, error) {
	if os.Getenv("SNAPSHOT_TESTS") == "true" { // coverage-ignore
		L.Warn().Msg("Snapshot tests mode, skipping connecting to Kubernetes API!")
		return nil, nil
	}
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	config, err := kubeConfig.ClientConfig()
	if err != nil { // coverage-ignore
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil { // coverage-ignore
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}
	return &API{
		ClientSet:  clientset,
		RESTConfig: config,
		namespace:  namespace,
	}, nil
}

// GetPods returns a list of Pods in the specified namespace.
func (k *API) GetPods(ctx context.Context) (*corev1.PodList, error) {
	pods, err := k.ClientSet.CoreV1().Pods(k.namespace).List(ctx, metav1.ListOptions{})
	if err != nil { // coverage-ignore
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}
	return pods, nil
}

// AllPodsReady checks if all Pods in the namespace are ready.
// A Pod is considered ready if all its containers are ready and the Pod's phase is "Running".
func (k *API) AllPodsReady(ctx context.Context) (bool, error) {
	pods, err := k.GetPods(ctx)
	if err != nil { // coverage-ignore
		return false, fmt.Errorf("failed to get pods: %v", err)
	}
	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning {
			return false, nil
		}
		L.Debug().Str("Pod", pod.Name).Str("Status", string(pod.Status.Phase)).Msg("Pod status")
		for _, containerStatus := range pod.Status.ContainerStatuses {
			L.Debug().Str("Pod", pod.Name).Str("Status", containerStatus.State.String()).Msg("Pod status")
			if !containerStatus.Ready {
				return false, nil
			}
		}
	}
	return true, nil
}

// CreateNamespace creates a new Kubernetes namespace with the specified name.
func (k *API) CreateNamespace(ctx context.Context, name string) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := k.ClientSet.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil { // coverage-ignore
		if strings.Contains(err.Error(), "already exists") {
			L.Debug().Str("Namespace", name).Msg("Namespace already exists, proceeding..")
			return nil
		}
		return fmt.Errorf("failed to create namespace: %v", err)
	}
	return nil
}

// RemoveNamespace deletes a Kubernetes namespace with the specified name.
func (k *API) RemoveNamespace(name string) error {
	err := k.ClientSet.CoreV1().Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil { // coverage-ignore
		return fmt.Errorf("failed to delete namespace: %v", err)
	}
	return nil
}

// waitAllPodsReady waits until all Pods in the namespace are ready or the timeout is reached.
// It retries the check periodically until the condition is met or the timeout occurs.
func (k *API) waitAllPodsReady(ctx context.Context, timeout time.Duration) (bool, error) {
	L.Info().Str("Namespace", k.namespace).Msg("Waiting for all pods to be ready")
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	retryInterval := 3 * time.Second
	time.Sleep(retryInterval)
	for {
		select {
		case <-ctx.Done():
			// coverage-ignore
			return false, fmt.Errorf("timeout reached while waiting for Pods to be ready")
		default:
			ready, err := k.AllPodsReady(ctx)
			if err != nil { // coverage-ignore
				return false, fmt.Errorf("failed to check Pod readiness: %v", err)
			}
			if ready {
				return true, nil
			}
			L.Debug().Msg("Checking if all pods are ready")
			time.Sleep(retryInterval)
		}
	}
}
