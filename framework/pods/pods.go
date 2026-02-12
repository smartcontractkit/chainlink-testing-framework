package pods

import (
	"context"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/yaml"
)

const (
	ManifestsDir       = "pods-out"
	K8sNamespaceEnvVar = "KUBERNETES_NAMESPACE"
)

// Client is a global K8s client that we use for all deployments
var Client *API

// Config describes Pods library configuration
type Config struct {
	Namespace string
	Pods      []*PodConfig
}

// PodConfig describes particular Pod configuration
type PodConfig struct {
	StatefulSet bool
	// Name is a pod name
	Name *string
	// Replicas amount of replicas for a pod
	Replicas *int32
	// Labels are K8s labels added to a pod
	Labels map[string]string
	// Annotations are K8s annotations added to a pod
	Annotations map[string]string
	// Image docker image URI in format $repo/$image_name:$tag, ex. "public.ecr.aws/chainlink/chainlink:v2.17.0"
	Image *string
	// Env represents container environment variables
	Env []corev1.EnvVar
	// Command is a container command to run on start
	Command *string
	// Ports is a list of $svc:$container ports, ex.: ["8080:80", "9090:90"]
	Ports []string
	// ConfigMap is a map of files in ConfigMap, ex.: "config.toml": `some_toml`
	ConfigMap map[string]string
	// ConfigMapMountPath mounts files with paths, ex.: "config.toml": "/config.toml"
	ConfigMapMountPath map[string]string
	// Secrets is a map of files in K8s Secret, ex. "secrets.toml": `some_secret`
	Secrets map[string]string
	// SecretsMountPath mounts secrets with paths, ex.: "secrets.toml": "/secrets.toml"
	SecretsMountPath map[string]string
	// ReadinessProbe is container readiness probe definition
	ReadinessProbe *corev1.Probe
	// Requests is K8s resources requests on CPU/Mem
	Requests map[string]string
	// Limits is K8s resources limits on CPU/Mem
	Limits map[string]string
	// ContainerSecurityContext is a container security context
	ContainerSecurityContext *corev1.SecurityContext
	// PodSecurityContext is a Pod security context
	PodSecurityContext *corev1.PodSecurityContext
	// VolumeClaimTemplates is a list K8s persistent volume claim templates
	VolumeClaimTemplates []corev1.PersistentVolumeClaim
}

// App is an application context with a generated Kubernetes manifest
type App struct {
	cfg      *Config
	objects  []any
	svcObj   *corev1.Service
	manifest string
}

// K8sEnabled is a flag that means Kubernetes in enabled, used in framework components
func K8sEnabled() bool {
	return os.Getenv(K8sNamespaceEnvVar) != ""
}

// Run generates and applies a new K8s YAML manifest
func Run(ctx context.Context, cfg *Config) (string, *corev1.Service, error) {
	var err error
	if cfg.Namespace == "" {
		cfg.Namespace = os.Getenv(K8sNamespaceEnvVar)
	}
	Client, err = NewAPI(cfg.Namespace)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create K8s client: %w", err)
	}
	if Client != nil {
		if err := Client.CreateNamespace(ctx, cfg.Namespace); err != nil {
			return "", nil, fmt.Errorf("failed to create namespace: %s, %w", cfg.Namespace, err)
		}
	}
	p := &App{
		cfg: cfg,
	}
	if err := p.generate(); err != nil {
		return p.Manifest(), nil, err
	}
	svc, err := p.apply()
	if err != nil {
		return p.Manifest(), svc, err
	}
	return p.Manifest(), svc, nil
}

// GetConnectionDetails returns connection details needed to forward ports
func (n *App) GetConnectionDetails() *corev1.Service {
	return n.svcObj
}

// generate provides a simplified template that is focused on deploying K8s Pods
func (n *App) generate() error {
	for _, podConfig := range n.cfg.Pods {
		podName := *podConfig.Name
		namespace := n.cfg.Namespace

		// Define resources
		if podConfig.Requests == nil {
			podConfig.Requests = ResourcesMedium()
		}
		if podConfig.Limits == nil {
			podConfig.Limits = ResourcesMedium()
		}

		// Define labels
		labels := map[string]string{"app": podName, "generated-by": "ctfv2"}
		maps.Copy(labels, podConfig.Labels)

		// Define annotations
		annotations := map[string]string{}
		maps.Copy(annotations, podConfig.Annotations)

		// Create ConfigMap if provided
		if len(podConfig.ConfigMap) > 0 {
			configMap := &corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-configmap", podName),
					Namespace: namespace,
				},
				Data: podConfig.ConfigMap,
			}
			n.objects = append(n.objects, configMap)
		}

		// Create Secret if provided
		if len(podConfig.Secrets) > 0 {
			secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-secret", podName),
					Namespace: namespace,
				},
				StringData: podConfig.Secrets,
			}
			n.objects = append(n.objects, secret)
		}

		// Define volumes and volume mounts
		var volumes []corev1.Volume
		var volumeMounts []corev1.VolumeMount

		// Prepare ConfigMap volumes
		idx := 0
		for _, fileName := range SortedKeys(podConfig.ConfigMapMountPath) {
			mountPath := podConfig.ConfigMapMountPath[fileName]
			volumes = append(volumes, corev1.Volume{
				Name: fmt.Sprintf("%s-configmap-volume-%d", podName, idx),
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: fmt.Sprintf("%s-configmap", podName),
						},
						Items: []corev1.KeyToPath{
							{
								Key:  fileName,
								Path: fileName,
							},
						},
					},
				},
			})
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      fmt.Sprintf("%s-configmap-volume-%d", podName, idx),
				MountPath: mountPath,
				SubPath:   fileName,
			})
			idx++
		}

		// Prepare secrets volumes
		idx = 0
		for _, fileName := range SortedKeys(podConfig.SecretsMountPath) {
			mountPath := podConfig.SecretsMountPath[fileName]
			volumes = append(volumes, corev1.Volume{
				Name: fmt.Sprintf("%s-secret-volume-%d", podName, idx),
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: fmt.Sprintf("%s-secret", podName),
						Items: []corev1.KeyToPath{
							{
								Key:  fileName,
								Path: fileName,
							},
						},
					},
				},
			})
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      fmt.Sprintf("%s-secret-volume-%d", podName, idx),
				MountPath: mountPath,
				SubPath:   fileName,
			})
			idx++
		}

		// Parse port mappings for the container
		var containerPorts []corev1.ContainerPort
		for i, portMapping := range podConfig.Ports {
			parts := strings.Split(portMapping, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid port mapping: %s, should be \"$svc_port:$container_port\"", portMapping)
			}

			containerPort, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid container port number: %s", parts[1])
			}

			containerPorts = append(containerPorts, corev1.ContainerPort{
				Name:          fmt.Sprintf("port-%d", i),
				ContainerPort: int32(containerPort),
			})
		}

		// Convert resources to Kubernetes format
		resourceRequests := corev1.ResourceList{}
		for k, v := range podConfig.Requests {
			quantity, err := resource.ParseQuantity(v)
			if err != nil {
				return fmt.Errorf("invalid resource request %s=%s: %w", k, v, err)
			}
			switch k {
			case "cpu":
				resourceRequests[corev1.ResourceCPU] = quantity
			case "memory":
				resourceRequests[corev1.ResourceMemory] = quantity
			}
		}

		resourceLimits := corev1.ResourceList{}
		for k, v := range podConfig.Limits {
			quantity, err := resource.ParseQuantity(v)
			if err != nil {
				return fmt.Errorf("invalid resource limit %s=%s: %w", k, v, err)
			}
			switch k {
			case "cpu":
				resourceLimits[corev1.ResourceCPU] = quantity
			case "memory":
				resourceLimits[corev1.ResourceMemory] = quantity
			}
		}

		container := corev1.Container{
			Name:  fmt.Sprintf("%s-container", podName),
			Image: *podConfig.Image,
			Env:   podConfig.Env,
			Ports: containerPorts,
			Resources: corev1.ResourceRequirements{
				Requests: resourceRequests,
				Limits:   resourceLimits,
			},
			VolumeMounts:    volumeMounts,
			SecurityContext: podConfig.ContainerSecurityContext,
			ReadinessProbe:  podConfig.ReadinessProbe,
		}

		// Transform container command
		if podConfig.Command != nil {
			container.Command = strings.Split(*podConfig.Command, " ")
		}

		// Override replicas
		replicas := int32(1)
		if podConfig.Replicas != nil {
			replicas = *podConfig.Replicas
		}

		// Create StatefulSet if volume claims present or StatefulSet flag is true
		if len(podConfig.VolumeClaimTemplates) > 0 || podConfig.StatefulSet {
			statefulSet := &appsv1.StatefulSet{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps/v1",
					Kind:       "StatefulSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-statefulset", podName),
					Namespace: namespace,
				},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: fmt.Sprintf("%s-svc", podName),
					Replicas:    &replicas,
					Selector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name:        fmt.Sprintf("%s-pod", podName),
							Labels:      labels,
							Annotations: annotations,
							Namespace:   namespace,
						},
						Spec: corev1.PodSpec{
							SecurityContext: podConfig.PodSecurityContext,
							Containers:      []corev1.Container{container},
							Volumes:         volumes,
						},
					},
					VolumeClaimTemplates: podConfig.VolumeClaimTemplates,
				},
			}
			n.objects = append(n.objects, statefulSet)
		} else {
			// Create Deployment
			deployment := &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-deployment", podName),
					Namespace: namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &replicas,
					Selector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      labels,
							Annotations: annotations,
							Namespace:   namespace,
						},
						Spec: corev1.PodSpec{
							SecurityContext: podConfig.PodSecurityContext,
							Containers:      []corev1.Container{container},
							Volumes:         volumes,
						},
					},
				},
			}
			n.objects = append(n.objects, deployment)
		}

		// Parse port mappings for the service
		var servicePorts []corev1.ServicePort
		for i, portMapping := range podConfig.Ports {
			parts := strings.Split(portMapping, ":")
			if len(parts) != 2 {
				L.Fatal().Msgf("Invalid port mapping: %s", portMapping)
			}

			port, err := strconv.ParseInt(parts[0], 10, 32)
			if err != nil {
				L.Fatal().Msgf("Invalid port number: %s", parts[0])
			}

			targetPort, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				L.Fatal().Msgf("Invalid container port number: %s", parts[1])
			}

			servicePorts = append(servicePorts, corev1.ServicePort{
				Name:       fmt.Sprintf("port-%d", i),
				Port:       int32(port),
				TargetPort: intstr.FromInt(int(targetPort)),
			})
		}

		if len(servicePorts) > 0 {
			// Create the Service
			service := &corev1.Service{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Service",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-svc", podName),
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Type:     corev1.ServiceTypeLoadBalancer,
					Ports:    servicePorts,
					Selector: labels,
				},
			}
			n.svcObj = service
			n.objects = append(n.objects, service)
		}
	}

	// Generate YAML from all objects
	yamlDocs := make([]string, 0, len(n.objects))
	for _, obj := range n.objects {
		yamlBytes, err := yaml.Marshal(obj)
		if err != nil {
			return fmt.Errorf("failed to marshal object to YAML: %w", err)
		}
		yamlDocs = append(yamlDocs, string(yamlBytes))
	}

	n.manifest = strings.Join(yamlDocs, "---\n")
	L.Debug().Msgf("Generated YAML:\n%s", n.manifest)
	return nil
}

func (n *App) apply() (*corev1.Service, error) {
	if os.Getenv("SNAPSHOT_TESTS") == "true" {
		return nil, nil
	}
	if n.manifest == "" {
		return nil, fmt.Errorf("manifest is empty, nothing to generate")
	}
	_ = os.Mkdir(ManifestsDir, os.ModePerm)
	// write generate manifest
	manifestFile := filepath.Join(ManifestsDir, fmt.Sprintf("pods-%s.tmp.yml", uuid.NewString()[0:5]))
	err := os.WriteFile(manifestFile, []byte(n.manifest), 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to write manifest to file: %w", err)
	}
	// apply the manifest
	cmd := exec.Command("kubectl", "apply", "-f", manifestFile, "--wait=true")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to apply manifest: %v\nOutput: %s", err, string(output))
	}
	L.Info().Str("Manifest", manifestFile).Msg("Manifest applied successfully")
	return n.svcObj, Connect(n.svcObj)
}

// Connect connects service to localhost, the same method is used internally
// by environment and externally by tests
// 'blocking' means it'd wait until first successful port connection
func Connect(svc *corev1.Service) error {
	ns := os.Getenv(K8sNamespaceEnvVar)
	if ns == "" {
		return fmt.Errorf("empty namespace")
	}
	var err error
	if Client == nil {
		Client, err = NewAPI(ns)
		if err != nil {
			return err
		}
	}
	f := NewForwarder(Client)
	forwardConfigs := make([]PortForwardConfig, 0)
	for _, p := range svc.Spec.Ports {
		forwardConfigs = append(forwardConfigs, PortForwardConfig{
			Namespace:   ns,
			ServiceName: svc.Name,
			LocalPort:   int(p.Port),
			ServicePort: int(p.Port),
		})
	}
	return f.Forward(forwardConfigs)
}

// WaitReady waits for all pods to be in status ready
func WaitReady(ctx context.Context, t time.Duration) error {
	_, err := Client.waitAllPodsReady(ctx, t)
	return err
}

// Manifest returns current generated YAML manifest
func (n *App) Manifest() string {
	return n.manifest
}

// NewKubernetesClient creates a new Kubernetes client
func NewKubernetesClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return clientset, nil
}
