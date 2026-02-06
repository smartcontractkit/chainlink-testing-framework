package pods

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/smartcontractkit/pods/imports/k8s"
)

// Config describes Pods library configuration
type Config struct {
	Namespace *string
	Pods      []*PodConfig
}

// PodConfig describes particular Pod configuration
type PodConfig struct {
	StatefulSet bool
	// Name is a pod name
	Name *string
	// Replicas amount of replicase for a pod
	Replicas *float64
	// Labels are K8s labels added to a pod
	Labels map[string]string
	// Annotations are K8s annotations added to a pod
	Annotations map[string]string
	// Image docker image URI in format $repo/$image_name:$tag, ex. "public.ecr.aws/chainlink/chainlink:v2.17.0"
	Image *string
	// Env represents container environment variables
	Env *[]*k8s.EnvVar
	// Command is a container command to run on start
	Command *string
	// Ports is a list of $svc:$container ports, ex.: ["8080:80", "9090:90"]
	Ports []string
	// ConfigMap is a map of files in ConfigMap, ex.: "config.toml": `some_toml`
	// ConfigMap key should be used in ConfigMapMountPath with a path to mount the file
	ConfigMap map[string]*string
	// ConfigMapMountPath mounts files with paths, ex.: "config.toml": "/config.toml"
	ConfigMapMountPath map[string]*string
	// Secrets is a map of files in K8s Secret, ex. "secrets.toml": `some_secret`
	// Secrets key should be used in SecretsMountPath with a path to mount the secret
	Secrets map[string]*string
	// SecretsMountPath mounts secrets with paths, ex.: "secrets.toml": "/secrets.toml"
	SecretsMountPath map[string]*string
	// ReadinessProbe is container readiness probe definition
	ReadinessProbe *k8s.Probe
	// Requests is K8s resources requests on CPU/Mem, see Resources func and examples in tests
	Requests map[string]k8s.Quantity
	// Limits is K8s resources limits on CPU/Mem, see Resources func and examples in tests
	Limits map[string]k8s.Quantity
	// ContainerSecurityContext is a container security context
	ContainerSecurityContext *k8s.SecurityContext
	// PodSecurityContext is a Pod security context
	PodSecurityContext *k8s.PodSecurityContext
	// VolumeClaimTemplates is a list K8s persistent volume claim templates
	// mostly used with databases, see SizedVolumeClaim used with PostgreSQL
	// if one template is present we deploy a StatefulSet
	VolumeClaimTemplates []*k8s.KubePersistentVolumeClaimProps
}

type Pods struct {
	*API
	cfg      *Config
	app      cdk8s.App
	chart    cdk8s.Chart
	manifest *string
}

// New creates and generates a new K8s YAML manifest
func New(cfg *Config) *Pods {
	api, err := NewAPI(*cfg.Namespace)
	if err != nil {
		return nil
	}
	p := &Pods{
		cfg: cfg,
		API: api,
	}
	return p
}

// Generate provides a simplified Docker Compose like API for K8s and generates YAML a manifest to deploy
func (n *Pods) Generate() error {
	n.app = cdk8s.NewApp(nil)
	n.chart = cdk8s.NewChart(n.app, S("pods-chart"), nil)
	for _, podConfig := range n.cfg.Pods {
		podName := *podConfig.Name
		namespace := n.cfg.Namespace

		// Define resources
		if podConfig.Requests == nil {
			podConfig.Requests = ResourcesSmall()
		}
		if podConfig.Limits == nil {
			podConfig.Limits = ResourcesSmall()
		}

		// Define labels
		labels := map[string]*string{"app": S(podName), "generated-by": S("cdk8s")}
		for k, v := range podConfig.Labels {
			labels[k] = S(v)
		}

		// Define annotations
		annotations := map[string]*string{}
		for k, v := range podConfig.Annotations {
			annotations[k] = S(v)
		}

		// Create ConfigMaps if provided
		if len(podConfig.ConfigMap) > 0 {
			k8s.NewKubeConfigMap(n.chart, S(fmt.Sprintf("%s-configmap", podName)), &k8s.KubeConfigMapProps{
				Metadata: &k8s.ObjectMeta{
					Name:      S(fmt.Sprintf("%s-configmap", podName)),
					Namespace: namespace,
				},
				Data: &podConfig.ConfigMap,
			})
		}

		// Create Secrets if provided
		if len(podConfig.Secrets) > 0 {
			k8s.NewKubeSecret(n.chart, S(fmt.Sprintf("%s-secret", podName)), &k8s.KubeSecretProps{
				Metadata: &k8s.ObjectMeta{
					Name:      S(fmt.Sprintf("%s-secret", podName)),
					Namespace: namespace,
				},
				StringData: &podConfig.Secrets,
			})
		}

		// Define volumes and volume mounts
		var volumes []*k8s.Volume
		var volumeMounts []*k8s.VolumeMount

		// Prepare ConfigMap volumes
		idx := 0
		for _, fileName := range SortedKeys(podConfig.ConfigMapMountPath) {
			mountPath := podConfig.ConfigMapMountPath[fileName]
			volumes = append(volumes, &k8s.Volume{
				Name: S(fmt.Sprintf("%s-configmap-volume-%d", podName, idx)),
				ConfigMap: &k8s.ConfigMapVolumeSource{
					Name: S(fmt.Sprintf("%s-configmap", podName)),
					Items: &[]*k8s.KeyToPath{
						{
							Key:  S(fileName),
							Path: S(fileName),
						},
					},
				},
			})
			volumeMounts = append(volumeMounts, &k8s.VolumeMount{
				Name:      S(fmt.Sprintf("%s-configmap-volume-%d", podName, idx)),
				MountPath: mountPath,
				SubPath:   S(fileName),
			})
			idx++
		}

		// Prepare secrets volumes
		idx = 0
		for _, fileName := range SortedKeys(podConfig.SecretsMountPath) {
			mountPath := podConfig.SecretsMountPath[fileName]
			volumes = append(volumes, &k8s.Volume{
				Name: S(fmt.Sprintf("%s-secret-volume-%d", podName, idx)),
				Secret: &k8s.SecretVolumeSource{
					SecretName: S(fmt.Sprintf("%s-secret", podName)),
					Items: &[]*k8s.KeyToPath{
						{
							Key:  S(fileName),
							Path: S(fileName),
						},
					},
				},
			})
			volumeMounts = append(volumeMounts, &k8s.VolumeMount{
				Name:      S(fmt.Sprintf("%s-secret-volume-%d", podName, idx)),
				MountPath: mountPath,
				SubPath:   S(fileName),
			})
			idx++
		}

		// Parse port mappings for the container
		var containerPorts []*k8s.ContainerPort
		for i, portMapping := range podConfig.Ports {
			parts := strings.Split(portMapping, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid port mapping: %s, should be \"$svc_port:$container_port\"", portMapping)
			}

			containerPort, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return fmt.Errorf("invalid container port number: %s", parts[1])
			}

			containerPorts = append(containerPorts, &k8s.ContainerPort{
				Name: S(fmt.Sprintf("port-%d", i)),
				// Use HostPort here to enable node exposed port
				ContainerPort: &containerPort,
			})
		}

		container := &k8s.Container{
			Name:            S(fmt.Sprintf("%s-container", podName)),
			Image:           podConfig.Image,
			Env:             podConfig.Env,
			Ports:           &containerPorts,
			Resources:       &k8s.ResourceRequirements{Limits: &podConfig.Limits, Requests: &podConfig.Requests},
			VolumeMounts:    &volumeMounts,
			SecurityContext: podConfig.ContainerSecurityContext,
			ReadinessProbe:  podConfig.ReadinessProbe,
		}

		// Transform container command
		if podConfig.Command != nil {
			cmds := strings.Split(*podConfig.Command, " ")
			L.Info().Msg(fmt.Sprintf("commands: %s", strings.Join(cmds, " ")))
			command := make([]*string, 0)
			for _, cmd := range cmds {
				command = append(command, S(cmd))
			}
			container.Command = &command
			L.Debug().Str("Cmd", *podConfig.Command).Msg("Container command")
		}

		// Override replicas
		replicas := I(1)
		if podConfig.Replicas != nil {
			replicas = podConfig.Replicas
		}

		// Create Deployment or StatefulSet if any volume claim is present
		if len(podConfig.VolumeClaimTemplates) > 0 || podConfig.StatefulSet {
			k8s.NewKubeStatefulSet(n.chart, S(fmt.Sprintf("%s-statefulset", podName)), &k8s.KubeStatefulSetProps{
				Metadata: &k8s.ObjectMeta{
					Name:      S(fmt.Sprintf("%s-statefulset", podName)),
					Namespace: namespace,
				},
				Spec: &k8s.StatefulSetSpec{
					ServiceName: S(fmt.Sprintf("%s-svc", podName)),
					Replicas:    replicas,
					Selector: &k8s.LabelSelector{
						MatchLabels: &labels,
					},
					Template: &k8s.PodTemplateSpec{
						Metadata: &k8s.ObjectMeta{
							Name:        S(fmt.Sprintf("%s-pp", podName)),
							Labels:      &labels,
							Annotations: &annotations,
							Namespace:   namespace,
						},
						Spec: &k8s.PodSpec{
							SecurityContext: podConfig.PodSecurityContext,
							Containers:      &[]*k8s.Container{container},
							Volumes:         &volumes,
						},
					},
					VolumeClaimTemplates: &podConfig.VolumeClaimTemplates,
				},
			})
		} else {
			k8s.NewKubeDeployment(n.chart, S(fmt.Sprintf("%s-deployment", podName)), &k8s.KubeDeploymentProps{
				Metadata: &k8s.ObjectMeta{
					Name:      S(fmt.Sprintf("%s-deployment", podName)),
					Namespace: namespace,
				},
				Spec: &k8s.DeploymentSpec{
					Replicas: podConfig.Replicas,
					Selector: &k8s.LabelSelector{
						MatchLabels: &labels,
					},
					Template: &k8s.PodTemplateSpec{
						Metadata: &k8s.ObjectMeta{
							Labels:      &labels,
							Annotations: &annotations,
							Namespace:   namespace,
						},
						Spec: &k8s.PodSpec{
							SecurityContext: podConfig.PodSecurityContext,
							Containers:      &[]*k8s.Container{container},
							Volumes:         &volumes,
						},
					},
				},
			})
		}

		// Parse port mappings for the service
		var servicePorts []*k8s.ServicePort
		for i, portMapping := range podConfig.Ports {
			parts := strings.Split(portMapping, ":")
			if len(parts) != 2 {
				log.Fatalf("Invalid port mapping: %s", portMapping)
			}

			port, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				log.Fatalf("Invalid port number: %s", parts[0])
			}

			containerPort, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				log.Fatalf("Invalid container port number: %s", parts[1])
			}

			servicePorts = append(servicePorts, &k8s.ServicePort{
				Name:       S(fmt.Sprintf("port-%d", i)),
				Port:       &port,
				TargetPort: k8s.IntOrString_FromNumber(&containerPort),
			})
		}

		if len(servicePorts) > 0 {
			// Create the KubeService with the parsed ports
			k8s.NewKubeService(n.chart, S(fmt.Sprintf("%s-svc", podName)), &k8s.KubeServiceProps{
				Metadata: &k8s.ObjectMeta{
					Name:      S(fmt.Sprintf("%s-svc", podName)),
					Namespace: namespace,
				},
				Spec: &k8s.ServiceSpec{
					Type:     S("LoadBalancer"),
					Ports:    &servicePorts,
					Selector: &labels,
				},
			})
		}
	}
	yaml := n.app.SynthYaml()
	L.Info().Msg(*yaml)
	n.manifest = yaml
	return nil
}

func (n *Pods) apply() error {
	if n.manifest == nil {
		return fmt.Errorf("manifest is not generated. Call Generate() first")
	}
	tmpFile := "pods.tmp.yml"
	err := os.WriteFile(tmpFile, []byte(*n.manifest), 0600)
	if err != nil {
		return fmt.Errorf("failed to write manifest to file: %v", err)
	}
	cmd := exec.Command("kubectl", "apply", "-f", tmpFile, "--wait=true")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply manifest: %v\nOutput: %s", err, string(output))
	}
	L.Info().Msg("Manifest applied successfully:")
	L.Info().Msg(string(output))
	return nil
}

// Apply runs "kubectl apply" with custom readiness assertions
func (n *Pods) Apply() error {
	if err := n.Generate(); err != nil {
		return err
	}
	if err := n.apply(); err != nil {
		return err
	}
	_, err := n.WaitForAllPodsReady(context.Background(), 3*time.Minute)
	return err
}

func (n *Pods) ResetPodsConfig() {
	n.cfg.Pods = make([]*PodConfig, 0)
}

// Manifest returns current generated YAML manifest
func (n *Pods) Manifest() *string {
	return n.manifest
}

func ApplyFlag() bool { return os.Getenv("APPLY") == "true" }
