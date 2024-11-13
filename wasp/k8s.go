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

// GetLocalK8sDeps retrieves a Kubernetes client set and its associated REST configuration 
// for interacting with a local Kubernetes cluster. 
// It returns a pointer to the kubernetes.Clientset, a pointer to the rest.Config, 
// and an error if any issues occur during the retrieval process. 
// If successful, the client set can be used to perform operations on the Kubernetes API.
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

// NewK8sClient initializes a new K8sClient instance by retrieving the necessary Kubernetes dependencies 
// from the local environment. It returns a pointer to the K8sClient, which contains the ClientSet and 
// RESTConfig required for interacting with the Kubernetes API. If there is an error while fetching 
// the dependencies, the function logs the error and terminates the program.
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

// jobPods retrieves a list of pods in the specified namespace that match the given synchronization label. 
// It returns a pointer to a v1.PodList containing the matching pods and an error if the operation fails. 
// If no pods are found, the returned PodList will be empty, and the error will be nil.
func (m *K8sClient) jobPods(ctx context.Context, nsName, syncLabel string) (*v1.PodList, error) {
	return m.ClientSet.CoreV1().Pods(nsName).List(ctx, metaV1.ListOptions{LabelSelector: syncSelector(syncLabel)})
}

// jobs retrieves a list of Kubernetes jobs in the specified namespace that match the given label selector. 
// It takes a context for cancellation and a namespace name along with a synchronization label as parameters. 
// The function returns a pointer to a batchV1.JobList containing the jobs that match the criteria, 
// or an error if the retrieval fails.
func (m *K8sClient) jobs(ctx context.Context, nsName, syncLabel string) (*batchV1.JobList, error) {
	return m.ClientSet.BatchV1().Jobs(nsName).List(ctx, metaV1.ListOptions{LabelSelector: syncSelector(syncLabel)})
}

// syncSelector formats a string to create a label selector for synchronization purposes. 
// It returns a string in the format "sync=<input>", where <input> is the provided string argument. 
// This formatted string can be used in Kubernetes API calls to filter resources based on the specified synchronization label.
func syncSelector(s string) string {
	return fmt.Sprintf("sync=%s", s)
}

// removeJobs deletes the specified jobs in the given namespace. 
// It takes a context for cancellation and a namespace name along with a list of jobs to be removed. 
// The function logs the removal process and ensures that each job is deleted with a foreground deletion policy. 
// If any error occurs during the deletion of a job, it returns the error. 
// If all jobs are successfully deleted, it returns nil.
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

// waitSyncGroup blocks until the specified number of pods with the given sync label are in the Running state. 
// It periodically checks the status of the pods in the specified namespace and logs the progress. 
// If an error occurs while retrieving the pods, it returns the error. 
// The function will return nil once all pods are confirmed to be running, or it will continue to wait if the conditions are not met.
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

// TrackJobs monitors the status of Kubernetes jobs and their associated pods in a specified namespace. 
// It continuously checks for the specified number of job pods and logs their statuses. 
// If any job fails, it returns an error and can optionally remove the jobs based on the keepJobs parameter. 
// The function will exit gracefully if the provided context is done, returning nil. 
// If all jobs succeed, it will either remove the jobs or return nil based on the keepJobs flag.
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
