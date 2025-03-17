package wasp

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	K8sStatePollInterval = 1 * time.Second
)

// K8sClient high level k8s client
type K8sClient struct {
	ClientSet  *kubernetes.Clientset
	RESTConfig *rest.Config
}

// GetLocalK8sDeps retrieves the local Kubernetes Clientset and REST configuration.
// It is used to initialize a Kubernetes client for interacting with the cluster.
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

// NewK8sClient initializes and returns a new K8sClient for interacting with the local Kubernetes cluster.
// It is used to perform operations such as synchronizing groups and managing cluster profiles.
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

// jobPods returns a list of pods in the specified namespace matching the sync label.
// It is used to track and manage job-related pods within Kubernetes environments.
func (m *K8sClient) jobPods(ctx context.Context, nsName, syncLabel string) (*v1.PodList, error) {
	return m.ClientSet.CoreV1().Pods(nsName).List(ctx, metaV1.ListOptions{LabelSelector: syncSelector(syncLabel)})
}

// jobs retrieves the list of Kubernetes jobs within the specified namespace
// that match the provided synchronization label.
// It returns a JobList and an error if the operation fails.
func (m *K8sClient) jobs(ctx context.Context, nsName, syncLabel string) (*batchV1.JobList, error) {
	return m.ClientSet.BatchV1().Jobs(nsName).List(ctx, metaV1.ListOptions{LabelSelector: syncSelector(syncLabel)})
}

// syncSelector formats a sync label into a label selector string.
// It is used to filter Kubernetes jobs and pods based on the specified synchronization label.
func syncSelector(s string) string {
	return fmt.Sprintf("sync=%s", s)
}

// removeJobs deletes all jobs in the given JobList within the specified namespace.
// It is used to clean up job resources after they have completed or failed.
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

// waitSyncGroup waits until the specified namespace has jobNum pods with the given syncLabel running.
// It ensures that all required pods are synchronized and operational before proceeding.
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

// TrackJobs monitors Kubernetes jobs in the specified namespace and label selector until the desired number succeed or a failure occurs.
// It optionally removes jobs upon completion based on the keepJobs flag.
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
					successfulJobs++
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
