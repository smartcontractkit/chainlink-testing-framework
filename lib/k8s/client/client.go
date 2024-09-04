package client

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/cmd/cp"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

const (
	TempDebugManifest    = "tmp-manifest-%s.yaml"
	K8sStatePollInterval = 10 * time.Second
	JobFinalizedTimeout  = 2 * time.Minute
	AppLabel             = "app"
)

// K8sClient high level k8s client
type K8sClient struct {
	ClientSet  *kubernetes.Clientset
	RESTConfig *rest.Config
}

// GetLocalK8sDeps get local k8s context config
func GetLocalK8sDeps() (*kubernetes.Clientset, *rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	k8sConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, nil, err
	}
	return k8sClient, k8sConfig, nil
}

// NewK8sClient creates a new k8s client with a REST config
func NewK8sClient() (*K8sClient, error) {
	cs, cfg, err := GetLocalK8sDeps()
	if err != nil {
		return nil, err
	}
	return &K8sClient{
		ClientSet:  cs,
		RESTConfig: cfg,
	}, nil
}

// ListPods lists pods for a namespace and selector
func (m *K8sClient) ListPods(namespace, selector string) (*v1.PodList, error) {
	pods, err := m.ClientSet.CoreV1().Pods(namespace).List(context.Background(), metaV1.ListOptions{LabelSelector: selector})
	sort.Slice(pods.Items, func(i, j int) bool {
		return pods.Items[i].CreationTimestamp.Before(pods.Items[j].CreationTimestamp.DeepCopy())
	})
	return pods.DeepCopy(), err
}

// ListPods lists services for a namespace and selector
func (m *K8sClient) ListServices(namespace, selector string) (*v1.ServiceList, error) {
	services, err := m.ClientSet.CoreV1().Services(namespace).List(context.Background(), metaV1.ListOptions{LabelSelector: selector})
	return services.DeepCopy(), err
}

// ListNamespaces lists k8s namespaces
func (m *K8sClient) ListNamespaces(selector string) (*v1.NamespaceList, error) {
	return m.ClientSet.CoreV1().Namespaces().List(context.Background(), metaV1.ListOptions{LabelSelector: selector})
}

// AddLabel adds a new label to a group of pods defined by selector
func (m *K8sClient) AddLabel(namespace string, selector string, label string) error {
	podList, err := m.ListPods(namespace, selector)
	if err != nil {
		return err
	}
	l := strings.Split(label, "=")
	if len(l) != 2 {
		return fmt.Errorf("labels must be in format key=value")
	}
	for _, pod := range podList.Items {
		labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, l[0], l[1])
		_, err := m.ClientSet.CoreV1().Pods(namespace).Patch(
			context.Background(),
			pod.GetName(),
			types.JSONPatchType,
			[]byte(labelPatch),
			metaV1.PatchOptions{},
		)
		if err != nil {
			return fmt.Errorf("failed to update labels %s for pod %s err: %w", labelPatch, pod.Name, err)
		}
	}
	log.Debug().Str("Selector", selector).Str("Label", label).Msg("Updated label")
	return nil
}

func (m *K8sClient) LabelChaosGroup(namespace string, labelPrefix string, startInstance int, endInstance int, group string) error {
	for i := startInstance; i <= endInstance; i++ {
		err := m.AddLabel(namespace, fmt.Sprintf("%s%d", labelPrefix, i), fmt.Sprintf("%s=1", group))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *K8sClient) LabelChaosGroupByLabels(namespace string, labels map[string]string, group string) error {
	labelSelector := ""
	for key, value := range labels {
		if labelSelector == "" {
			labelSelector = fmt.Sprintf("%s=%s", key, value)
		} else {
			labelSelector = fmt.Sprintf("%s, %s=%s", labelSelector, key, value)
		}
	}
	podList, err := m.ListPods(namespace, labelSelector)
	if err != nil {
		return err
	}
	for _, pod := range podList.Items {
		err = m.AddPodLabel(namespace, pod, group, "1")
		if err != nil {
			return err
		}
	}
	return nil
}

// AddPodsLabels adds map of labels to all pods in list
func (m *K8sClient) AddPodsLabels(namespace string, podList *v1.PodList, labels map[string]string) error {
	for _, pod := range podList.Items {
		for k, v := range labels {
			err := m.AddPodLabel(namespace, pod, k, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AddPodsAnnotations adds map of annotations to all pods in list
func (m *K8sClient) AddPodsAnnotations(namespace string, podList *v1.PodList, annotations map[string]string) error {
	// when applying annotations the key doesn't like `/` characters here but everywhere else it does
	// replacing it here with ~1
	fixedAnnotations := make(map[string]string)
	for k, v := range annotations {
		fixedAnnotations[strings.ReplaceAll(k, "/", "~1")] = v
	}
	for _, pod := range podList.Items {
		for k, v := range fixedAnnotations {
			err := m.AddPodAnnotation(namespace, pod, k, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// UniqueLabels gets all unique application labels
func (m *K8sClient) UniqueLabels(namespace string, selector string) ([]string, error) {
	uniqueLabels := make([]string, 0)
	isUnique := make(map[string]bool)
	podList, err := m.ListPods(namespace, selector)
	if err != nil {
		return nil, err
	}
	for _, p := range podList.Items {
		appLabel := p.Labels[AppLabel]
		if _, ok := isUnique[appLabel]; !ok {
			uniqueLabels = append(uniqueLabels, appLabel)
		}
	}
	log.Info().
		Interface("Apps", uniqueLabels).
		Int("Count", len(uniqueLabels)).
		Msg("Apps found")
	return uniqueLabels, nil
}

// AddPodLabel adds a label to a pod
func (m *K8sClient) AddPodLabel(namespace string, pod v1.Pod, key, value string) error {
	labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, key, value)
	_, err := m.ClientSet.CoreV1().Pods(namespace).Patch(
		context.Background(), pod.GetName(), types.JSONPatchType, []byte(labelPatch), metaV1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

// AddPodAnnotation adds an annotation to a pod
func (m *K8sClient) AddPodAnnotation(namespace string, pod v1.Pod, key, value string) error {
	labelPatch := fmt.Sprintf(`[{"op":"add","path":"/metadata/annotations/%s","value":"%s" }]`, key, value)
	_, err := m.ClientSet.CoreV1().Pods(namespace).Patch(
		context.Background(), pod.GetName(), types.JSONPatchType, []byte(labelPatch), metaV1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

// EnumerateInstances enumerate pods with instance label
func (m *K8sClient) EnumerateInstances(namespace string, selector string) error {
	podList, err := m.ListPods(namespace, selector)
	if err != nil {
		return err
	}

	for id, pod := range podList.Items {
		// skip if already labeled with instance
		existingLabels := pod.Labels
		_, exists := existingLabels["instance"]
		if exists {
			continue
		}
		if err := m.AddPodLabel(namespace, pod, "instance", strconv.Itoa(id)); err != nil {
			return err
		}
	}
	return nil
}

// waitForPodsExist waits for all the expected number of pods to exist
func (m *K8sClient) waitForPodsExist(ns string, expectedPodCount int) error {
	log.Debug().Int("ExpectedCount", expectedPodCount).Msg("Waiting for pods to exist")
	var exitErr error
	timeout := 15 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		// nolint:contextcheck
		apps, err2 := m.UniqueLabels(ns, AppLabel)
		if err2 != nil {
			exitErr = err2
			return false, nil
		}
		if len(apps) >= expectedPodCount {
			exitErr = nil
			return true, nil
		}
		return false, nil
	}); err != nil {
		return err
	}

	return exitErr
}

// WaitPodsReady waits until all pods are ready
func (m *K8sClient) WaitPodsReady(ns string, rcd *ReadyCheckData, expectedPodCount int) error {
	// Wait for pods to exist
	err := m.waitForPodsExist(ns, expectedPodCount)
	if err != nil {
		return err
	}

	log.Info().Msg("Waiting for pods to be ready")
	ticker := time.NewTicker(K8sStatePollInterval)
	defer ticker.Stop()
	timeout := time.NewTimer(rcd.Timeout)
	readyCount := 0
	defer timeout.Stop()
	for {
		select {
		case <-timeout.C:
			return fmt.Errorf("waitcontainersready, no pods in '%s' with selector '%s' after timeout '%s'",
				ns, rcd.ReadinessProbeCheckSelector, rcd.Timeout)
		case <-ticker.C:
			podList, err := m.ListPods(ns, rcd.ReadinessProbeCheckSelector)
			if err != nil {
				return err
			}
			if len(podList.Items) == 0 && expectedPodCount > 0 {
				log.Debug().
					Str("Namespace", ns).
					Str("Selector", rcd.ReadinessProbeCheckSelector).
					Msg("No pods found with selector")
				continue
			}
			log.Debug().Interface("Pods", podNames(podList)).Msg("Waiting for pods readiness probes")
			allReady := true
			for _, pod := range podList.Items {
				if pod.Status.Phase == "Succeeded" {
					log.Debug().Str("Pod", pod.Name).Msg("Pod is in Succeeded state")
					continue
				} else if pod.Status.Phase != v1.PodRunning {
					log.Debug().Str("Pod", pod.Name).Str("Phase", string(pod.Status.Phase)).Msg("Pod is not running")
					allReady = false
					break
				}
				for _, c := range pod.Status.Conditions {
					if c.Type == v1.ContainersReady && c.Status != "True" {
						log.Debug().Str("Text", c.Message).Msg("Pod condition message")
						allReady = false
					}
				}
			}

			if allReady {
				readyCount++
				// wait for it to be ready 3 times since there is no good way to know if an old pod
				// was present but not yet decommisiond during a rollout
				// usually there is just a very small blip that we can run into this and this will
				// prevent that from happening
				if readyCount == 3 {
					return nil
				}
			}
		}
	}
}

// NamespaceExists check if namespace exists
func (m *K8sClient) NamespaceExists(namespace string) bool {
	if _, err := m.ClientSet.CoreV1().Namespaces().Get(context.Background(), namespace, metaV1.GetOptions{}); err != nil {
		return false
	}
	return true
}

// RemoveNamespace removes namespace
func (m *K8sClient) RemoveNamespace(namespace string) error {
	log.Info().Str("Namespace", namespace).Msg("Removing namespace")
	return m.ClientSet.CoreV1().Namespaces().Delete(context.Background(), namespace, metaV1.DeleteOptions{})
}

// CopyFromPod copies src from a particular namespace:pod:container to local.
func (m *K8sClient) CopyFromPod(ctx context.Context, namespace, podName, containername, srcPath, destPath string) error {
	scmd := fmt.Sprintf("kubectl cp %s/%s:%s -c %s %s", namespace, podName, srcPath, containername, destPath)
	log.Info().Str("Command", scmd).Msg("Waiting to finish copying")
	return ExecCmdWithContext(ctx, scmd)
}

// RolloutStatefulSets applies "rollout statefulset" to all existing statefulsets in that namespace
func (m *K8sClient) RolloutStatefulSets(ctx context.Context, namespace string) error {
	stsClient := m.ClientSet.AppsV1().StatefulSets(namespace)
	sts, err := stsClient.List(ctx, metaV1.ListOptions{})
	if err != nil {
		return err
	}
	for _, s := range sts.Items {
		cmd := fmt.Sprintf("kubectl rollout restart statefulset %s --namespace %s", s.Name, namespace)
		log.Info().Str("Command", cmd).Msg("Applying StatefulSet rollout")
		if err := ExecCmdWithContext(ctx, cmd); err != nil {
			return err
		}
	}
	// wait for the statefulsets to be ready in a separate loop otherwise this can take a long time
	for _, s := range sts.Items {
		// wait for the rollout to be complete
		scmd := fmt.Sprintf("kubectl rollout status statefulset %s --namespace %s", s.Name, namespace)
		log.Info().Str("Command", scmd).Msg("Waiting for StatefulSet rollout to finish")
		if err := ExecCmdWithContext(ctx, scmd); err != nil {
			return err
		}
	}
	return nil
}

// RolloutRestartBySelector rollouts and restarts object by selector
func (m *K8sClient) RolloutRestartBySelector(ctx context.Context, namespace, resource, selector string) error {
	cmd := fmt.Sprintf("kubectl --namespace %s rollout restart -l %s %s", namespace, selector, resource)
	log.Info().Str("Command", cmd).Msg("rollout restart by selector")
	if err := ExecCmdWithContext(ctx, cmd); err != nil {
		return err
	}
	// wait for the rollout to be complete
	waitCmd := fmt.Sprintf("kubectl --namespace %s rollout status -l %s %s", namespace, selector, resource)
	log.Info().Str("Command", waitCmd).Msg("Waiting for StatefulSet rollout to finish")
	return ExecCmdWithContext(ctx, waitCmd)
}

// ReadyCheckData data to check if selected pods are running and all containers are ready ( readiness check ) are ready
type ReadyCheckData struct {
	ReadinessProbeCheckSelector string
	Timeout                     time.Duration
}

// WaitForJob wait for job execution, follow logs and returns an error if job failed
func (m *K8sClient) WaitForJob(namespaceName string, jobName string, fundReturnStatus func(string)) error {
	cmd := fmt.Sprintf("kubectl --namespace %s logs --follow job/%s", namespaceName, jobName)
	log.Info().Str("Job", jobName).Str("cmd", cmd).Msg("Waiting for job to complete")
	ctx := context.Background()
	if err := ExecCmdWithOptions(ctx, cmd, fundReturnStatus); err != nil {
		return err
	}
	var exitErr error
	ctx, cancel := context.WithTimeout(ctx, JobFinalizedTimeout)
	defer cancel()
	if err := wait.PollUntilContextTimeout(ctx, K8sStatePollInterval, JobFinalizedTimeout, true, func(ctx context.Context) (bool, error) {
		job, err := m.ClientSet.BatchV1().Jobs(namespaceName).Get(ctx, jobName, metaV1.GetOptions{})
		if err != nil {
			exitErr = err
		}
		if int(job.Status.Failed) > 0 {
			exitErr = fmt.Errorf("job failed")
			return true, nil
		}
		if int(job.Status.Succeeded) > 0 {
			exitErr = nil
			return true, nil
		}
		return false, nil
	}); err != nil {
		return err
	}
	return exitErr
}

func (m *K8sClient) WaitForDeploymentsAvailable(ctx context.Context, namespace string) error {
	deployments, err := m.ClientSet.AppsV1().Deployments(namespace).List(ctx, metaV1.ListOptions{})
	if err != nil {
		return err
	}
	log.Debug().Int("Number", len(deployments.Items)).Msg("Deployments found")
	for _, d := range deployments.Items {
		log.Debug().Str("status", d.Status.String()).Msg("Deployment info")
		waitCmd := fmt.Sprintf("kubectl rollout status -n %s deployment/%s", namespace, d.Name)
		log.Debug().Str("cmd", waitCmd).Msg("wait for deployment to be available")
		if err := ExecCmdWithContext(ctx, waitCmd); err != nil {
			return err
		}
	}
	return nil
}

// Apply applying a manifest to a currently connected k8s context
func (m *K8sClient) Apply(ctx context.Context, manifest, namespace string, waitForDeployment bool) error {
	manifestFile := fmt.Sprintf(TempDebugManifest, uuid.NewString())
	log.Info().Str("File", manifestFile).Msg("Applying manifest")
	if err := os.WriteFile(manifestFile, []byte(manifest), os.ModePerm); err != nil {
		return err
	}
	cmd := fmt.Sprintf("kubectl apply -f %s", manifestFile)
	log.Debug().Str("cmd", cmd).Msg("Apply command")
	if err := ExecCmdWithContext(ctx, cmd); err != nil {
		return err
	}
	if waitForDeployment {
		return m.WaitForDeploymentsAvailable(ctx, namespace)
	}
	return nil
}

// DeleteResource deletes resource
func (m *K8sClient) DeleteResource(namespace string, resource string, instance string) error {
	return ExecCmd(fmt.Sprintf("kubectl delete %s %s --namespace %s", resource, instance, namespace))
}

// Create creating a manifest to a currently connected k8s context
func (m *K8sClient) Create(manifest string) error {
	manifestFile := fmt.Sprintf(TempDebugManifest, uuid.NewString())
	log.Info().Str("File", manifestFile).Msg("Creating manifest")
	if err := os.WriteFile(manifestFile, []byte(manifest), os.ModePerm); err != nil {
		return err
	}
	cmd := fmt.Sprintf("kubectl create -f %s", manifestFile)
	return ExecCmd(cmd)
}

// DryRun generates manifest and writes it in a file
func (m *K8sClient) DryRun(manifest string) error {
	manifestFile := fmt.Sprintf(TempDebugManifest, uuid.NewString())
	log.Info().Str("File", manifestFile).Msg("Creating manifest")
	return os.WriteFile(manifestFile, []byte(manifest), os.ModePerm)
}

// CopyToPod copies src to a particular container. Destination should be in the form of a proper K8s destination path
// NAMESPACE/POD_NAME:folder/FILE_NAME
func (m *K8sClient) CopyToPod(namespace, src, destination, containername string) (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer, error) {
	m.RESTConfig.APIPath = "/api"
	m.RESTConfig.GroupVersion = &schema.GroupVersion{Version: "v1"} // this targets the core api groups so the url path will be /api/v1
	m.RESTConfig.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	ioStreams, in, out, errOut := genericclioptions.NewTestIOStreams()

	copyOptions := cp.NewCopyOptions(ioStreams)
	configFlags := genericclioptions.NewConfigFlags(false)
	f := cmdutil.NewFactory(configFlags)
	cmd := cp.NewCmdCp(f, ioStreams)
	err := copyOptions.Complete(f, cmd, []string{src, destination})
	if err != nil {
		return nil, nil, nil, err
	}
	copyOptions.Clientset = m.ClientSet
	copyOptions.ClientConfig = m.RESTConfig
	copyOptions.Container = containername
	copyOptions.Namespace = namespace

	formatted, err := regexp.MatchString(".*?\\/.*?\\:.*", destination)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not parse the pod destination: %w", err)
	}
	if !formatted {
		return nil, nil, nil, fmt.Errorf("pod destination string improperly formatted, see reference 'NAMESPACE/POD_NAME:folder/FILE_NAME'")
	}

	log.Info().
		Str("Namespace", namespace).
		Str("Source", src).
		Str("Destination", destination).
		Str("Container", containername).
		Msg("Uploading file to pod")
	err = copyOptions.Run()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not run copy operation: %w", err)
	}
	return in, out, errOut, nil
}

// ExecuteInPod is similar to kubectl exec
func (m *K8sClient) ExecuteInPod(namespace, podName, containerName string, command []string) ([]byte, []byte, error) {
	log.Info().Interface("Command", command).Msg("Executing command in pod")
	req := m.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Container: containerName,
		Command:   command,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(m.RESTConfig, "POST", req.URL())
	if err != nil {
		return []byte{}, []byte{}, err
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	return stdout.Bytes(), stderr.Bytes(), err
}

func podNames(podItems *v1.PodList) []string {
	on := make([]string, 0)
	for _, p := range podItems.Items {
		on = append(on, p.Name)
	}
	return on
}
