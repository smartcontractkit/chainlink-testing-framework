package comparator

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/pkg/errors"
	tc "github.com/testcontainers/testcontainers-go"
	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ExecutionEnvironment string

const (
	ExecutionEnvironment_Docker                  ExecutionEnvironment = "docker"
	ExecutionEnvironment_k8sExecutionEnvironment                      = "k8s"
)

type ResourceReporter struct {
	// either k8s or docker
	ExecutionEnvironment ExecutionEnvironment `json:"execution_environment"`

	// AUT metrics
	PodsResources      map[string]*PodResources    `json:"pods_resources"`
	ContainerResources map[string]*DockerResources `json:"container_resources"`
	// regex pattern to select the resources we want to fetch
	ResourceSelectionPattern string `json:"resource_selection_pattern"`
}

type DockerResources struct {
	NanoCPUs   int64
	CpuShares  int64
	Memory     int64
	MemorySwap int64
}

type PodResources struct {
	RequestsCPU    int64
	RequestsMemory int64
	LimitsCPU      int64
	LimitsMemory   int64
}

func (r *ResourceReporter) FetchResources() error {
	if r.ExecutionEnvironment == ExecutionEnvironment_Docker {
		err := r.fetchDockerResources()
		if err != nil {
			return err
		}
	} else {
		err := r.fetchK8sResources()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ResourceReporter) fetchK8sResources() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrapf(err, "failed to get in-cluster config, are you sure this is running in a k8s cluster?")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrapf(err, "failed to create k8s clientset")
	}

	namespaceFile := "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	namespace, err := os.ReadFile(namespaceFile)
	if err != nil {
		return errors.Wrapf(err, "failed to read namespace file %s", namespaceFile)
	}

	pods, err := clientset.CoreV1().Pods(string(namespace)).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	r.PodsResources = make(map[string]*PodResources)

	for _, pod := range pods.Items {
		r.PodsResources[pod.Name] = &PodResources{
			RequestsCPU:    pod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue(),
			RequestsMemory: pod.Spec.Containers[0].Resources.Requests.Memory().Value(),
			LimitsCPU:      pod.Spec.Containers[0].Resources.Limits.Cpu().MilliValue(),
			LimitsMemory:   pod.Spec.Containers[0].Resources.Limits.Memory().Value(),
		}
	}

	return nil
}

func (r *ResourceReporter) fetchDockerResources() error {
	provider, err := tc.NewDockerProvider()
	if err != nil {
		return fmt.Errorf("failed to create Docker provider: %w", err)
	}

	containers, err := provider.Client().ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list Docker containers: %w", err)
	}

	eg := &errgroup.Group{}
	pattern := regexp.MustCompile(r.ResourceSelectionPattern)

	var dockerResources = make(map[string]*DockerResources)

	for _, containerInfo := range containers {
		eg.Go(func() error {
			containerName := containerInfo.Names[0]
			if !pattern.Match([]byte(containerName)) {
				return nil
			}

			ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
			info, err := provider.Client().ContainerInspect(ctx, containerInfo.ID)
			if err != nil {
				cancelFn()
				return errors.Wrapf(err, "failed to inspect container %s", containerName)
			}

			cancelFn()
			dockerResources[containerName] = &DockerResources{
				NanoCPUs:   info.HostConfig.NanoCPUs,
				CpuShares:  info.HostConfig.CPUShares,
				Memory:     info.HostConfig.Memory,
				MemorySwap: info.HostConfig.MemorySwap,
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrapf(err, "failed to fetch Docker resources")
	}

	r.ContainerResources = dockerResources

	return nil
}

func (r *ResourceReporter) CompareResources(other *ResourceReporter) error {
	if r.ExecutionEnvironment != other.ExecutionEnvironment {
		return fmt.Errorf("execution environments are different. Expected %s, got %s", r.ExecutionEnvironment, other.ExecutionEnvironment)
	}

	if r.ExecutionEnvironment == ExecutionEnvironment_Docker {
		return r.compareDockerResources(other.ContainerResources)
	}

	return r.comparePodResources(other.PodsResources)
}

func (r *ResourceReporter) comparePodResources(other map[string]*PodResources) error {
	this := r.PodsResources
	if len(this) != len(other) {
		return fmt.Errorf("pod resources count is different. Expected %d, got %d", len(this), len(other))
	}

	for name1, res1 := range this {
		if res2, ok := other[name1]; !ok {
			return fmt.Errorf("pod resource %s is missing from the other report", name1)
		} else {
			if res1 == nil {
				return fmt.Errorf("pod resource %s is nil in the current report", name1)
			}
			if res2 == nil {
				return fmt.Errorf("pod resource %s is nil in the other report", name1)
			}
			if *res1 != *res2 {
				return fmt.Errorf("pod resource %s is different. Expected %v, got %v", name1, res1, res2)
			}
		}
	}

	for name2 := range other {
		if _, ok := this[name2]; !ok {
			return fmt.Errorf("pod resource %s is missing from the current report", name2)
		}
	}

	return nil
}

func (r *ResourceReporter) compareDockerResources(other map[string]*DockerResources) error {
	this := r.ContainerResources
	if len(this) != len(other) {
		return fmt.Errorf("container resources count is different. Expected %d, got %d", len(this), len(other))
	}

	for name1, res1 := range this {
		if res2, ok := other[name1]; !ok {
			return fmt.Errorf("container resource %s is missing from the other report", name1)
		} else {
			if res1 == nil {
				return fmt.Errorf("container resource %s is nil in the current report", name1)
			}
			if res2 == nil {
				return fmt.Errorf("container resource %s is nil in the other report", name1)
			}
			if *res1 != *res2 {
				return fmt.Errorf("container resource %s is different. Expected %v, got %v", name1, res1, res2)
			}
		}
	}

	for name2 := range other {
		if _, ok := this[name2]; !ok {
			return fmt.Errorf("container resource %s is missing from the current report", name2)
		}
	}

	return nil
}
