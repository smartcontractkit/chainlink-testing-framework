package networktest

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

type AlpineInput struct {
	Privileged    bool              // Whether to run in privileged mode
	BlockInternet bool              // Whether to block internet access
	Labels        map[string]string // Container labels
}

type AlpineOutput struct{}

// NewNetworkTest creates a minimal Alpine Linux container for network testing
func NewNetworkTest(in AlpineInput) (*AlpineOutput, error) {
	req := testcontainers.ContainerRequest{
		Name:     "networktest",
		Image:    "alpine:latest",
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {"networktest"},
		},
		Labels:     framework.DefaultTCLabels(),
		WaitingFor: wait.ForLog(""),
		Cmd:        []string{"/bin/sh", "-c", "while true; do sleep 30; done;"},
	}

	if in.BlockInternet {
		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.DNS = []string{"127.0.0.1"}
			hc.CapAdd = []string{"NET_ADMIN"}
			if in.Privileged {
				hc.Privileged = true
			}
		}

		req.LifecycleHooks = []testcontainers.ContainerLifecycleHooks{{
			PostStarts: []testcontainers.ContainerHook{
				func(ctx context.Context, c testcontainers.Container) error {
					// Block all internet traffic while allowing local network
					_, _, err := c.Exec(ctx, []string{
						"sh", "-c", `iptables -A OUTPUT -d 10.0.0.0/8 -j ACCEPT &&
		                          iptables -A OUTPUT -d 172.16.0.0/12 -j ACCEPT &&
		                          iptables -A OUTPUT -d 192.168.0.0/16 -j ACCEPT &&
		                          iptables -A OUTPUT -j DROP`,
					})
					return err
				},
			},
		}}
	}

	_, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start alpine container: %w", err)
	}

	return &AlpineOutput{}, nil
}
