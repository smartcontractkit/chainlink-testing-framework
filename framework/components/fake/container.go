package fake

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

type Input struct {
	Image string  `toml:"image"`
	Port  int     `toml:"port" validate:"required"`
	Out   *Output `toml:"out"`
}

type Output struct {
	UseCache      bool   `toml:"use_cache"`
	BaseURLHost   string `toml:"base_url_host"`
	BaseURLDocker string `toml:"base_url_docker"`
}

func defaults(in *Input) {
}

// NewFakeDataProvider creates new fake data provider
func NewFakeDataProvider(in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	ctx := context.Background()
	defaults(in)
	bindPort := fmt.Sprintf("%d/tcp", in.Port)
	req := tc.ContainerRequest{
		Name:  in.Image,
		Image: in.Image,
		//Labels:   framework.DefaultTCLabels(),
		//Networks: []string{framework.DefaultNetworkName},
		//NetworkAliases: map[string][]string{
		//	framework.DefaultNetworkName: {containerName},
		//},
		ExposedPorts: []string{bindPort},
		HostConfigModifier: func(h *container.HostConfig) {
			//h.PortBindings = framework.MapTheSamePort(bindPort)
		},
		Env: map[string]string{
			//"DATABASE_URL": pgOut.JDInternalURL,
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
	out := &Output{
		BaseURLHost:   fmt.Sprintf("http://localhost:%d", in.Port),
		BaseURLDocker: fmt.Sprintf("%s:%d", HostDockerInternal(), in.Port),
	}
	in.Out = out
	//out := &Output{
	//	UseCache:         true,
	//	ContainerName:    containerName,
	//	DBContainerName:  pgOut.ContainerName,
	//	ExternalGRPCUrl:  fmt.Sprintf("%s:%s", host, in.GRPCPort),
	//	InternalGRPCUrl:  fmt.Sprintf("%s:%s", containerName, in.GRPCPort),
	//	ExternalWSRPCUrl: fmt.Sprintf("%s:%s", host, in.WSRPCPort),
	//	InternalWSRPCUrl: fmt.Sprintf("%s:%s", containerName, in.WSRPCPort),
	//}
	//in.Out = out
	return out, nil
}
