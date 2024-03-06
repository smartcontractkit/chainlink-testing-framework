package k8s_client

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	K8sStatePollInterval = 3 * time.Second
)

// Client high level k8s client
type Client struct {
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
func NewClient() *Client {
	cs, cfg, err := GetLocalK8sDeps()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return &Client{
		ClientSet:  cs,
		RESTConfig: cfg,
	}
}

// k8sOperation defines the function type for Kubernetes operations that need retries.
type k8sOperation func() error

// retryK8sCall attempts the provided Kubernetes operation with retries.
func retryK8sCall(operation k8sOperation, maxRetries int) error {
	retryPolicy := wait.Backoff{
		Steps:    maxRetries,       // Max number of retries.
		Duration: 10 * time.Second, // Initial delay before the first retry.
		Factor:   2.0,              // Multiplier factor for subsequent delays.
		Jitter:   0.1,              // Random jitter for the delay.
	}

	var lastError error
	err := wait.ExponentialBackoff(retryPolicy, func() (bool, error) {
		lastError = operation()
		if lastError != nil {
			log.Warn().Err(lastError).Msg("Error encountered during K8s call, will retry")
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return fmt.Errorf("after %d attempts, last error: %s", maxRetries, lastError)
	}
	return nil
}

func (m *Client) ListPods(ctx context.Context, namespace, syncLabel string) (*v1.PodList, error) {
	var pods *v1.PodList
	maxRetries := 5 // Maximum number of retries

	timeout := int64(30)
	opts := metaV1.ListOptions{
		LabelSelector:  syncSelector(syncLabel), // Using syncSelector to format label selector
		TimeoutSeconds: &timeout,
	}
	operation := func() error {
		var err error
		pods, err = m.ClientSet.CoreV1().Pods(namespace).List(ctx, opts)
		return err
	}

	err := retryK8sCall(operation, maxRetries)
	if err != nil {
		// Wrap and return any error encountered during the retry operation
		return nil, fmt.Errorf("failed to call CoreV1().Pods(%s).List(%v): %w", namespace, opts, err)
	}

	// At this point, `pods` should be populated successfully
	return pods, nil
}

func (m *Client) ListJobs(ctx context.Context, namespace, syncLabel string) (*batchV1.JobList, error) {
	var jobs *batchV1.JobList
	maxRetries := 5 // Maximum number of retries

	timeout := int64(30)
	opts := metaV1.ListOptions{
		LabelSelector:  syncSelector(syncLabel), // Assuming syncSelector is a function that formats the label selector
		TimeoutSeconds: &timeout,
	}
	call := func() error {
		var err error
		jobs, err = m.ClientSet.BatchV1().Jobs(namespace).List(ctx, opts)
		return err
	}

	err := retryK8sCall(call, maxRetries)
	if err != nil {
		// Wrap and return any error encountered during the retry operation
		return nil, fmt.Errorf("failed to call BatchV1().Jobs(%s).List(%v): %w", namespace, opts, err)
	}

	// At this point, `jobs` should be populated successfully
	return jobs, nil
}

func (m *Client) GetPodLogs(ctx context.Context, nsName, syncLabel string) (map[string]string, error) {
	podLogs := make(map[string]string)
	maxRetries := 5 // Maximum number of retries

	timeout := int64(30)
	opts := metaV1.ListOptions{
		LabelSelector:  syncSelector(syncLabel),
		TimeoutSeconds: &timeout,
	}
	operation := func() error {
		pods, err := m.ClientSet.CoreV1().Pods(nsName).List(ctx, opts)
		if err != nil {
			return err // Return the error to the retry mechanism
		}

		for _, pod := range pods.Items {
			req := m.ClientSet.CoreV1().Pods(nsName).GetLogs(pod.Name, &v1.PodLogOptions{})
			podLog, err := req.Stream(ctx)
			if err != nil {
				return fmt.Errorf("failed to open log stream for pod %s: %w", pod.Name, err)
			}
			defer podLog.Close()

			logs, err := io.ReadAll(podLog)
			if err != nil {
				return fmt.Errorf("failed to read log for pod %s: %w", pod.Name, err)
			}

			podLogs[pod.Name] = string(logs)
		}

		return nil // Success
	}

	err := retryK8sCall(operation, maxRetries)
	if err != nil {
		// Wrap and return any error encountered during the retry operation
		return nil, fmt.Errorf("failed to retrieve pod logs: %w", err)
	}

	return podLogs, nil
}

func syncSelector(s string) string {
	return fmt.Sprintf("sync=%s", s)
}

func (m *Client) removeJobs(ctx context.Context, nsName string, jobs *batchV1.JobList) error {
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

func (m *Client) waitSyncGroup(ctx context.Context, nsName string, syncLabel string, jobNum int) error {
outer:
	for {
		time.Sleep(K8sStatePollInterval)
		log.Info().Str("SyncLabel", syncLabel).Msg("Awaiting group sync")
		pods, err := m.ListPods(ctx, nsName, syncLabel)
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

// TrackJobs tracks both jobs and their pods until they succeed or fail
func (m *Client) TrackJobs(ctx context.Context, nsName, syncLabel string, jobNum int, keepJobs bool) error {
	log.Debug().Str("LabelSelector", syncSelector(syncLabel)).Msg("Searching for jobs/pods")
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Cluster context finished")
			return nil
		default:
			time.Sleep(K8sStatePollInterval)
			jobs, err := m.ListJobs(ctx, nsName, syncLabel)
			if err != nil {
				return errors.Wrapf(err, "failed to get jobs")
			}
			jobPods, err := m.ListPods(ctx, nsName, syncLabel)
			if err != nil {
				return errors.Wrapf(err, "failed to get job pods")
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
					logs, err := m.GetPodLogs(ctx, nsName, syncLabel)
					if err != nil {
						log.Warn().Err(err).Msg("Failed to get pod logs")
					} else {
						for k, v := range logs {
							log.Info().Str("Pod", k).Str("Logs", v).Msg("Pod logs")
						}
					}
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
