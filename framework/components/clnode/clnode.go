package clnode

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
	"time"
)

const (
	DefaultTestKeystorePassword = "thispasswordislongenough"
)

type Input struct {
	DataProviderURL string          `toml:"data_provider_url" validate:"required"`
	DbInput         *postgres.Input `toml:"db" validate:"required"`
	Node            *NodeInput      `toml:"node" validate:"required"`
	Out             *Output         `toml:"out"`
}

type Output struct {
	Node       *NodeOut         `toml:"node"`
	PostgreSQL *postgres.Output `toml:"postgresql"`
}

type NodeInput struct {
	Image                string `toml:"image" validate:"required"`
	Tag                  string `toml:"tag" validate:"required"`
	Port                 string `toml:"port" validate:"required"`
	TestConfigOverrides  string `toml:"test_config_overrides"`
	UserConfigOverrides  string `toml:"user_config_overrides"`
	TestSecretsOverrides string `toml:"test_secrets_overrides"`
	UserSecretsOverrides string `toml:"user_secrets_overrides"`
}

type NodeOut struct {
	Url               string `toml:"url"`
	DockerInternalURL string `toml:"docker_internal_url"`
}

func NewNode(in *Input) (*Output, error) {
	if in.Out != nil && framework.NoCache() {
		return in.Out, nil
	}
	pgOut, err := postgres.NewPostgreSQL(in.DbInput)
	if err != nil {
		return nil, err
	}
	nodeOut, err := newNode(in, pgOut)
	if err != nil {
		return nil, err
	}
	out := &Output{
		Node:       nodeOut,
		PostgreSQL: pgOut,
	}
	in.Out = out
	return out, nil
}

func newNode(in *Input, pgOut *postgres.Output) (*NodeOut, error) {
	ctx := context.Background()

	passwordPath := filepath.Join(framework.PathCLNode, "password.txt")
	apiCredentialsPath := filepath.Join(framework.PathCLNode, "apicredentials")
	cfgPath, err := writeDefaultConfig(in)
	if err != nil {
		return nil, err
	}
	secretsPath, err := writeDefaultSecrets(pgOut)
	if err != nil {
		return nil, err
	}
	cfgOverridesPath, err := writeTestConfigOverrides(in.Node.TestConfigOverrides)
	if err != nil {
		return nil, err
	}
	cfgSecretsOverridesPath, err := writeTestSecretsOverrides(in.Node.TestSecretsOverrides)
	if err != nil {
		return nil, err
	}
	cfgUserOverridesPath, err := writeUserConfigOverrides(in.Node.UserConfigOverrides)
	if err != nil {
		return nil, err
	}
	cfgUserSecretsOverridesPath, err := writeUserSecretsOverrides(in.Node.UserSecretsOverrides)
	if err != nil {
		return nil, err
	}

	bindPort := fmt.Sprintf("%s/tcp", in.Node.Port)
	containerName := framework.DefaultTCName("clnode")

	req := tc.ContainerRequest{
		Image:    fmt.Sprintf("%s:%s", in.Node.Image, in.Node.Tag),
		Name:     containerName,
		Labels:   framework.DefaultTCLabels(),
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		ExposedPorts: []string{bindPort},
		Entrypoint: []string{
			"/bin/sh", "-c",
			"chainlink -c /config/config -c /config/overrides -c /config/user-overrides -s /config/secrets -s /config/secrets-overrides -s /config/user-secrets-overrides node start -d -p /config/node_password -a /config/apicredentials",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      cfgPath,
				ContainerFilePath: "/config/config",
				FileMode:          0644,
			},
			{
				HostFilePath:      secretsPath,
				ContainerFilePath: "/config/secrets",
				FileMode:          0644,
			},
			{
				HostFilePath:      cfgOverridesPath,
				ContainerFilePath: "/config/overrides",
				FileMode:          0644,
			},
			{
				HostFilePath:      cfgUserOverridesPath,
				ContainerFilePath: "/config/user-overrides",
				FileMode:          0644,
			},
			{
				HostFilePath:      cfgSecretsOverridesPath,
				ContainerFilePath: "/config/secrets-overrides",
				FileMode:          0644,
			},
			{
				HostFilePath:      cfgUserSecretsOverridesPath,
				ContainerFilePath: "/config/user-secrets-overrides",
				FileMode:          0644,
			},
			{
				HostFilePath:      passwordPath,
				ContainerFilePath: "/config/node_password",
				FileMode:          0644,
			},
			{
				HostFilePath:      apiCredentialsPath,
				ContainerFilePath: "/config/apicredentials",
				FileMode:          0644,
			},
		},
		WaitingFor: wait.ForLog("Listening and serving HTTP").WithStartupTimeout(2 * time.Minute),
	}
	// TODO: this is complex, though, desired by developers because of static addresses and fast login
	// TODO: skipping for now
	//if in.HostNetworkEnabled {
	//req.HostConfigModifier = func(hc *container.HostConfig) {
	//	hc.NetworkMode = "host"
	//	hc.PortBindings = framework.MapTheSamePort(bindPort)
	//}
	//}
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
	mp, err := c.MappedPort(ctx, nat.Port(bindPort))
	if err != nil {
		return nil, err
	}

	return &NodeOut{
		DockerInternalURL: fmt.Sprintf("http://%s:%s", containerName, in.Node.Port),
		Url:               fmt.Sprintf("%s:%s", host, mp.Port()),
	}, nil
}
