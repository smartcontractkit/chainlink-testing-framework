package benchspy

import (
	"context"
)

type ExecutionEnvironment string

const (
	ExecutionEnvironment_Docker                  ExecutionEnvironment = "docker"
	ExecutionEnvironment_k8sExecutionEnvironment ExecutionEnvironment = "k8s"
)

type ResourceReporter struct {
	// either k8s or docker
	ExecutionEnvironment ExecutionEnvironment `json:"execution_environment"`

	// regex pattern to select the resources we want to fetch
	ResourceSelectionPattern string `json:"resource_selection_pattern"`
}

func (r *ResourceReporter) FetchResources(_ context.Context) error {
	// for k8s we should query Prometheus to get the CPU and mem usage (median, min, max)
	// for Docker we need to wait for @Sergey to add a container that will expose historical CPU and mem usage

	return nil
}
