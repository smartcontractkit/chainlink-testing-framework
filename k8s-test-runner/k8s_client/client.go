package k8s_client

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	exec "github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/exec"
)

const (
	K8sStatePollInterval = 3 * time.Second
)

// High level k8s client
type Client struct {
	ClientSet       *kubernetes.Clientset
	RESTConfig      *rest.Config
	callRetryPolicy wait.Backoff
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
		callRetryPolicy: wait.Backoff{
			Steps:    5,                // Max number of retries.
			Duration: 10 * time.Second, // Initial delay before the first retry.
			Factor:   2.0,              // Multiplier factor for subsequent delays.
			Jitter:   0.1,              // Random jitter for the delay.
		},
	}
}

// k8sOperation defines the function type for Kubernetes operations that need retries.
type k8sOperation func() error

// retryK8sCall attempts the provided Kubernetes operation with retries.
func retryK8sCall(operation k8sOperation, retryPolicy wait.Backoff) error {
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
		return fmt.Errorf("after %d attempts, last error: %s", retryPolicy.Steps, lastError)
	}
	return nil
}

func (m *Client) ListPods(ctx context.Context, namespace, syncLabel string) (*v1.PodList, error) {
	var pods *v1.PodList
	timeout := int64(30)
	labelSelector := syncSelector(syncLabel)
	call := func() error {
		var err error
		pods, err = m.ClientSet.CoreV1().Pods(namespace).List(ctx, metaV1.ListOptions{
			LabelSelector:  labelSelector,
			TimeoutSeconds: &timeout,
		})
		return err
	}

	err := retryK8sCall(call, m.callRetryPolicy)
	if err != nil {
		// Wrap and return any error encountered during the retry operation
		return nil, fmt.Errorf(`failed to call CoreV1().Pods().List(), namespace: %s, labelSelector: %s, timeout: %d: %w`, namespace, labelSelector, timeout, err)
	}

	// At this point, `pods` should be populated successfully
	return pods, nil
}

func (m *Client) ListJobs(ctx context.Context, namespace, syncLabel string) (*batchV1.JobList, error) {
	var jobs *batchV1.JobList

	timeout := int64(30)
	labelSelector := syncSelector(syncLabel)
	call := func() error {
		var err error
		jobs, err = m.ClientSet.BatchV1().Jobs(namespace).List(ctx, metaV1.ListOptions{
			LabelSelector:  labelSelector,
			TimeoutSeconds: &timeout,
		})
		return err
	}

	err := retryK8sCall(call, m.callRetryPolicy)
	if err != nil {
		// Wrap and return any error encountered during the retry operation
		return nil, fmt.Errorf(`failed to call BatchV1().Jobs().List(), namespace: %s, labelSelector: %s, timeout: %d: %w`, namespace, labelSelector, timeout, err)
	}

	// At this point, `jobs` should be populated successfully
	return jobs, nil
}

func (m *Client) GetPodLogs(ctx context.Context, namespace string, pods []v1.Pod) (map[string]string, error) {
	podLogs := make(map[string]string)

	call := func() error {
		for _, pod := range pods {
			req := m.ClientSet.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{})
			podLog, err := req.Stream(ctx)
			if err != nil {
				return fmt.Errorf("failed to open log stream for pod %s: %w", pod.Name, err)
			}

			logs, err := io.ReadAll(podLog)
			if closeErr := podLog.Close(); closeErr != nil {
				log.Debug().Err(closeErr).Msg("Failed to close log stream")
			}
			if err != nil {
				return fmt.Errorf("failed to read log for pod %s: %w", pod.Name, err)
			}

			podLogs[pod.Name] = string(logs)
		}

		return nil // Success
	}

	err := retryK8sCall(call, m.callRetryPolicy)
	if err != nil {
		// Wrap and return any error encountered during the retry operation
		return nil, fmt.Errorf("failed to retrieve pod logs: %w", err)
	}

	return podLogs, nil
}

func syncSelector(s string) string {
	return fmt.Sprintf("sync=%s", s)
}

func (m *Client) RemoveJobs(ctx context.Context, nsName string, syncLabelValue string) error {
	jobs, err := m.ListJobs(ctx, nsName, syncLabelValue)
	if err != nil {
		return err
	}

	log.Info().Interface("jobs", jobs).Msg("Removing jobs")

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

// Use polling to wait for jobs to complete.
// `watch` does not work in CI (GAP), so we have to poll.
func (m *Client) WaitUntilJobsComplete(ctx context.Context, namespace, syncLabelValue string, expectedJobsCount int) error {
	labelSelector := syncSelector(syncLabelValue)
	completedJobs := make(map[string]bool)
	pollingInterval := time.Second * 5

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Context canceled")
			return ctx.Err()
		default:
			jobs, err := m.ClientSet.BatchV1().Jobs(namespace).List(ctx, metaV1.ListOptions{
				LabelSelector:  labelSelector,
				TimeoutSeconds: ptr.Int64(30), // query timeout
			})

			if err != nil {
				log.Error().Err(err).Str("labelSelector", labelSelector).Msg("Failed to list jobs, will retry...")

				cmd := fmt.Sprintf("kubectl get jobs -l sync=%s -n wasp --v=7", syncLabelValue)
				log.Info().Str("cmd", cmd).Msg("Running CLI command with verbose output to debug...")
				_ = exec.Cmd(cmd)

				time.Sleep(pollingInterval)
				continue
			}

			for _, job := range jobs.Items {
				if job.Status.Succeeded > 0 {
					completedJobs[job.Name] = true
					log.Info().Str("job", job.Name).Msg("Job succeeded")
				} else if job.Status.Failed > 0 {
					completedJobs[job.Name] = true
					log.Info().Str("job", job.Name).Msg("Job failed")
					return errors.Errorf("job %s failed", job.Name)
				}
			}

			if len(completedJobs) >= expectedJobsCount {
				log.Info().Msgf("All %d jobs completed", expectedJobsCount)
				return nil
			} else {
				log.Info().Msgf("Waiting for %d job(s) to complete...", expectedJobsCount-len(completedJobs))
			}

			time.Sleep(pollingInterval)
		}
	}
}

func (m *Client) PrintPodLogs(ctx context.Context, namespace string, syncLabelValue string) {
	pods, err := m.ListPods(ctx, namespace, syncLabelValue)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list pods")
		return
	}

	logs, err := m.GetPodLogs(ctx, namespace, pods.Items)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get pod logs")
	} else {
		for k, v := range logs {
			log.Info().Str("Pod", k).Msg("Pod logs")
			fmt.Println(v)
		}
	}
}

func (m *Client) LogNamespaceEvents(ctx context.Context, namespace string) {
	events, err := m.ClientSet.CoreV1().Events(namespace).List(ctx, metaV1.ListOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get namespace events")
	} else {
		log.Info().Msg("Namespace events:")
		for _, e := range events.Items {
			log.Info().Interface("Event", e).Msg("Event")
		}
	}
}
