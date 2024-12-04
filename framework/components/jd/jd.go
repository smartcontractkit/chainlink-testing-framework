package jd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"os"
)

const (
	TmpImageName            = "jd-local"
	GRPCPort         string = "42242"
	CSAEncryptionKey string = "!PASsword000!"
	WSRPCPort        string = "8080"
)

type Input struct {
	Image            string  `toml:"image"`
	GRPCPort         string  `toml:"grpc_port"`
	WSRPCPort        string  `toml:"wsrpc_port"`
	DBURL            string  `toml:"db_url"`
	CSAEncryptionKey string  `toml:"csa_encryption_key"`
	DockerFilePath   string  `toml:"docker_file"`
	DockerContext    string  `toml:"docker_ctx"`
	Out              *Output `toml:"out"`
}

type Output struct {
	UseCache       bool   `toml:"use_cache"`
	HostGRPCUrl    string `toml:"grpc_url"`
	DockerGRPCUrl  string `toml:"docker_internal_grpc_url"`
	HostWSRPCUrl   string `toml:"wsrpc_url"`
	DockerWSRPCUrl string `toml:"docker_internal_wsrpc_url"`
}

func defaults(in *Input) {
	if in.GRPCPort == "" {
		in.GRPCPort = GRPCPort
	}
	if in.WSRPCPort == "" {
		in.WSRPCPort = WSRPCPort
	}
	if in.CSAEncryptionKey == "" {
		in.CSAEncryptionKey = CSAEncryptionKey
	}
}

func NewJD(in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	ctx := context.Background()
	defaults(in)
	jdImg := os.Getenv("CTF_JD_IMAGE")
	if jdImg != "" {
		in.Image = jdImg
	}
	containerName := framework.DefaultTCName("jd")
	bindPort := fmt.Sprintf("%s/tcp", in.GRPCPort)
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
		Env: map[string]string{
			"DATABASE_URL":              in.DBURL,
			"PORT":                      in.GRPCPort,
			"NODE_RPC_PORT":             in.WSRPCPort,
			"CSA_KEY_ENCRYPTION_SECRET": in.CSAEncryptionKey,
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(nat.Port(fmt.Sprintf("%s/tcp", in.GRPCPort))),
		),
	}
	if req.Image == "" {
		req.Image = TmpImageName
		if err := framework.BuildImage(in.DockerContext, in.DockerFilePath, req.Image); err != nil {
			return nil, err
		}
		req.KeepImage = false
	}
	c, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	host, err := framework.GetHost(c)
	if err != nil {
		return nil, err
	}
	out := &Output{
		UseCache:       true,
		HostGRPCUrl:    fmt.Sprintf("http://%s:%s", host, in.GRPCPort),
		DockerGRPCUrl:  fmt.Sprintf("http://%s:%s", containerName, in.GRPCPort),
		HostWSRPCUrl:   fmt.Sprintf("ws://%s:%s", host, in.WSRPCPort),
		DockerWSRPCUrl: fmt.Sprintf("ws://%s:%s", containerName, in.WSRPCPort),
	}
	in.Out = out
	return out, nil
}
