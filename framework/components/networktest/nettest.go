package networktest

/*
This component exists purely for Docker debug purposes
*/

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

type Input struct {
	Name          string
	NoDNS         bool
	CustomNetwork bool
}

type Output struct{}

// NewNetworkTest creates a minimal Alpine Linux container for network testing
func NewNetworkTest(in Input) error {
	req := testcontainers.ContainerRequest{
		Name:     in.Name,
		Image:    "alpine:latest",
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {"networktest"},
		},
		Labels:     framework.DefaultTCLabels(),
		WaitingFor: wait.ForLog(""),
		Cmd:        []string{"/bin/sh", "-c", "while true; do sleep 30; done;"},
	}
	req.HostConfigModifier = func(hc *container.HostConfig) {
		// Remove external DNS
		framework.NoDNS(in.NoDNS, hc)
	}

	_, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start alpine container: %w", err)
	}
	return nil
}
