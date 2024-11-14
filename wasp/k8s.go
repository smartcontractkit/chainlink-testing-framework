package wasp

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

const (
	K8sStatePollInterval = 1 * time.Second
)

// K8sClient high level k8s client
type K8sClient struct {
	ClientSet  *kubernetes.Clientset
	RESTConfig *rest.Config
}

// GetLocalK8sDeps retrieves the local Kubernetes clientset and REST configuration.
// It loads the default kubeconfig and initializes a clientset based on the configuration.
// Returns the clientset, the REST config, and any error encountered during the process.
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

// NewK8sClient creates a new K8sClient for interacting with the Kubernetes cluster.
// It returns a pointer to the initialized K8sClient.
// If the client cannot be initialized, the program will terminate.
func NewK8sClient() *K8sClient {
	cs, cfg, err := GetLocalK8sDeps()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return &K8sClient{
		ClientSet:  cs,
		RESTConfig: cfg,
	}
}

// jobPods retrieves the Pods in the specified namespace that match the given sync label.
// It returns a PodList containing the matching Pods and any error encountered during the retrieval.
func (m *K8sClient) jobPods(ctx context.Context, nsName, syncLabel string) (*v1.PodList, error) {
	return m.ClientSet.CoreV1().Pods(nsName).List(ctx, metaV1.ListOptions{LabelSelector: syncSelector(syncLabel)})
}

// jobs fetches the Kubernetes JobList in the specified namespace filtered by the syncLabel.
// It uses the provided context for the request and returns the JobList or an error if the retrieval fails.
func (m *K8sClient) jobs(ctx context.Context, nsName, syncLabel string) (*batchV1.JobList, error) {
	return m.ClientSet.BatchV1().Jobs(nsName).List(ctx, metaV1.ListOptions{LabelSelector: syncSelector(syncLabel)})
}

// syncSelector returns a label selector formatted as "sync=<s>".
// It is used to filter Kubernetes jobs and pods based on the provided sync label.
func syncSelector(s string) string {
	return fmt.Sprintf("sync=%s", s)
}

// removeJobs deletes all jobs in the provided JobList within the specified namespace.
// It applies a foreground deletion propagation policy to ensure proper cleanup.
// The function returns an error if any job deletion fails.
func (m *K8sClient) removeJobs(ctx context.Context, nsName string, jobs *batchV1.JobList) error {
	log.Info().Msg("Removing jobs")
	for _, j := range jobs.Items {
		dp := metaV1.DeletePropagationForeground
		if err := m.ClientSet.BatchV1().Jobs(nsName).Delete(ctx, j.Name, metaV1.DeleteOptions{
			PropagationPolicy: &dp,
		}); err != nil {
			return err
		}
	}
	return nil
}

// waitSyncGroup blocks until jobNum pods with the specified syncLabel in the namespace nsName are in the Running state.
// It periodically polls the Kubernetes API and logs the synchronization status.
// Returns nil when all pods are running, or an error if polling fails.
func (m *K8sClient) waitSyncGroup(ctx context.Context, nsName string, syncLabel string, jobNum int) error {
outer:
	for {
		time.Sleep(K8sStatePollInterval)
		log.Info().Str("SyncLabel", syncLabel).Msg("Awaiting group sync")
		pods, err := m.jobPods(ctx, nsName, syncLabel)
		if err != nil {
			return err
		}
		if len(pods.Items) != jobNum {
			log.Info().Str("SyncLabel", syncLabel).Msg("Awaiting pods")
			continue
		}
		for _, p := range pods.Items {
			if p.Status.Phase != v1.PodRunning {
				continue outer
			}
		}
		return nil
	}
}

// TrackJobs monitors Kubernetes jobs in the specified namespace using the given label selector.
// It waits until the number of job pods matches jobNum and all jobs have succeeded.
// If keepJobs is false, it removes the jobs after completion.
// The function continues tracking until all jobs succeed, a job fails, or the provided context is canceled.
// It returns an error if any job fails or if monitoring is interrupted by the context.
func (m *K8sClient) TrackJobs(ctx context.Context, nsName, syncLabel string, jobNum int, keepJobs bool) error {
	log.Debug().Str("LabelSelector", syncSelector(syncLabel)).Msg("Searching for jobs/pods")
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Cluster context finished")
			return nil
		default:
			time.Sleep(K8sStatePollInterval)
			jobs, err := m.jobs(ctx, nsName, syncLabel)
			if err != nil {
				return err
			}
			jobPods, err := m.jobPods(ctx, nsName, syncLabel)
			if err != nil {
				return err
			}
			if len(jobPods.Items) != jobNum {
				log.Info().Int("JobPods", jobNum).Msg("Awaiting job pods")
				continue
			}
			for _, jp := range jobPods.Items {
				log.Debug().Interface("Phase", jp.Status.Phase).Msg("Job status")
			}
			var successfulJobs int
			for _, j := range jobs.Items {
				log.Debug().Interface("Status", j.Status).Str("Name", j.Name).Msg("Pod status")
				if j.Status.Failed > 0 {
					log.Warn().Str("Name", j.Name).Msg("Job has failed")
					if !keepJobs {
						if err := m.removeJobs(ctx, nsName, jobs); err != nil {
							return err
						}
					}
					return fmt.Errorf("job %s has failed", j.Name)
				}
				if j.Status.Succeeded > 0 {
					successfulJobs += 1
				}
			}
			if successfulJobs == jobNum {
				log.Info().Msg("Test ended")
				if !keepJobs {
					return m.removeJobs(ctx, nsName, jobs)
				}
				return nil
			}
		}
	}
}
