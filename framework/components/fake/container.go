package fake

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	v1 "k8s.io/api/core/v1"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/pods"
)

type Input struct {
	Image string  `toml:"image" comment:"Fake service image, usually can be found in our ECR with $project-fakes name"`
	Port  int     `toml:"port" validate:"required" comment:"The port which Docker container is exposing"`
	Out   *Output `toml:"out" comment:"Fakes service config output"`
}

type Output struct {
	UseCache      bool   `toml:"use_cache" comment:"Whether to respect caching or not, if cache = true component won't be deployed again"`
	BaseURLHost   string `toml:"base_url_host" comment:"Base URL which can be used when running locally"`
	BaseURLDocker string `toml:"base_url_docker" comment:"Base URL to reach fakes service from other Docker containers"`
	// K8sService is a Kubernetes service spec used to connect locally
	K8sService *v1.Service `toml:"k8s_service" comment:"Kubernetes service spec used to connect locally"`
}

// NewDockerFakeDataProvider creates new fake data provider in Docker using testcontainers-go
func NewDockerFakeDataProvider(in *Input) (*Output, error) {
	return NewWithContext(context.Background(), in)
}

// NewWithContext creates new fake data provider in Docker using testcontainers-go
func NewWithContext(ctx context.Context, in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	bindPort := fmt.Sprintf("%d/tcp", in.Port)
	containerName := framework.DefaultTCName("fake")
	if pods.K8sEnabled() {
		_, svc, err := pods.Run(ctx, &pods.Config{
			Pods: []*pods.PodConfig{
				{
					Name:     pods.Ptr(containerName),
					Image:    &in.Image,
					Ports:    []string{fmt.Sprintf("%d:%d", in.Port, in.Port)},
					Requests: pods.ResourcesSmall(),
					Limits:   pods.ResourcesSmall(),
					ContainerSecurityContext: &v1.SecurityContext{
						RunAsUser:  pods.Ptr[int64](999),
						RunAsGroup: pods.Ptr[int64](999),
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}
		in.Out = &Output{
			K8sService:    svc,
			BaseURLHost:   fmt.Sprintf("http://%s:%d", "localhost", in.Port),
			BaseURLDocker: fmt.Sprintf("http://%s:%d", fmt.Sprintf("%s-svc", containerName), in.Port),
		}
		return in.Out, nil
	}
	req := tc.ContainerRequest{
		Name:     containerName,
		Image:    in.Image,
		Labels:   framework.DefaultTCLabels(),
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		ExposedPorts: []string{bindPort},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(bindPort)
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(nat.Port(fmt.Sprintf("%d/tcp", in.Port))),
		),
	}
	_, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	in.Out = &Output{
		BaseURLHost:   fmt.Sprintf("http://localhost:%d", in.Port),
		BaseURLDocker: fmt.Sprintf("http://%s:%d", containerName, in.Port),
	}
	return in.Out, nil
}
