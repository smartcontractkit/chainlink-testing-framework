package clnode

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"text/template"
	"time"
)

// Input represents Chainlink node input
type Input struct {
	DataProviderURL string          `toml:"data_provider_url" validate:"required"`
	DbInput         *postgres.Input `toml:"db" validate:"required"`
	Node            *NodeInput      `toml:"node" validate:"required"`
	Out             *Output         `toml:"out"`
}

// NodeInput is CL nod container inputs
type NodeInput struct {
	Image                string `toml:"image" validate:"required"`
	Tag                  string `toml:"tag" validate:"required"`
	Port                 string `toml:"port" validate:"required"`
	TestConfigOverrides  string `toml:"test_config_overrides"`
	UserConfigOverrides  string `toml:"user_config_overrides"`
	TestSecretsOverrides string `toml:"test_secrets_overrides"`
	UserSecretsOverrides string `toml:"user_secrets_overrides"`
}

// Output represents Chainlink node output, nodes and databases connection URLs
type Output struct {
	Node       *NodeOut         `toml:"node"`
	PostgreSQL *postgres.Output `toml:"postgresql"`
}

// NodeOut is CL node container output, URLs to connect
type NodeOut struct {
	Url               string `toml:"url"`
	DockerInternalURL string `toml:"docker_internal_url"`
}

// NewNode create a new Chainlink node with some image:tag and one or several configs
// see config params: TestConfigOverrides, UserConfigOverrides, etc
func NewNode(in *Input) (*Output, error) {
	if in.Out != nil && framework.UseCache() {
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

	passwordPath, err := writeToFile(DefaultPasswordTxt, "password.txt")
	apiCredentialsPath, err := writeToFile(DefaultAPICredentials, "apicredentials")
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	cfgPath, err := writeDefaultConfig(in)
	if err != nil {
		return nil, err
	}
	secretsPath, err := writeDefaultSecrets(pgOut)
	if err != nil {
		return nil, err
	}
	overridesFile, err := writeToFile(in.Node.TestConfigOverrides, "overrides.toml")
	if err != nil {
		return nil, err
	}
	secretsOverridesFile, err := writeToFile(in.Node.TestSecretsOverrides, "secrets-overrides.toml")
	if err != nil {
		return nil, err
	}
	userOverridesFile, err := writeToFile(in.Node.UserConfigOverrides, "user-overrides.toml")
	if err != nil {
		return nil, err
	}
	userSecretsOverridesFile, err := writeToFile(in.Node.UserSecretsOverrides, "user-secrets-overrides.toml")
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
				HostFilePath:      cfgPath.Name(),
				ContainerFilePath: "/config/config",
				FileMode:          0644,
			},
			{
				HostFilePath:      secretsPath.Name(),
				ContainerFilePath: "/config/secrets",
				FileMode:          0644,
			},
			{
				HostFilePath:      overridesFile.Name(),
				ContainerFilePath: "/config/overrides",
				FileMode:          0644,
			},
			{
				HostFilePath:      userOverridesFile.Name(),
				ContainerFilePath: "/config/user-overrides",
				FileMode:          0644,
			},
			{
				HostFilePath:      secretsOverridesFile.Name(),
				ContainerFilePath: "/config/secrets-overrides",
				FileMode:          0644,
			},
			{
				HostFilePath:      userSecretsOverridesFile.Name(),
				ContainerFilePath: "/config/user-secrets-overrides",
				FileMode:          0644,
			},
			{
				HostFilePath:      passwordPath.Name(),
				ContainerFilePath: "/config/node_password",
				FileMode:          0644,
			},
			{
				HostFilePath:      apiCredentialsPath.Name(),
				ContainerFilePath: "/config/apicredentials",
				FileMode:          0644,
			},
		},
		WaitingFor: wait.ForLog("Listening and serving HTTP").WithStartupTimeout(2 * time.Minute),
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
	mp, err := c.MappedPort(ctx, nat.Port(bindPort))
	if err != nil {
		return nil, err
	}

	return &NodeOut{
		DockerInternalURL: fmt.Sprintf("http://%s:%s", containerName, in.Node.Port),
		Url:               fmt.Sprintf("%s:%s", host, mp.Port()),
	}, nil
}

type DefaultCLNodeConfig struct {
	HTTPPort      string
	SecureCookies bool
}

func generateDefaultConfig(in *Input) (string, error) {
	config := DefaultCLNodeConfig{
		HTTPPort:      in.Node.Port,
		SecureCookies: false,
	}
	tmpl, err := template.New("toml").Parse(defaultConfigTmpl)
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	err = tmpl.Execute(&output, config)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}

type DefaultSecretsConfig struct {
	DatabaseURL string
	Keystore    string
}

func generateSecretsConfig(connString, password string) (string, error) {
	// Create the configuration with example values
	config := DefaultSecretsConfig{
		DatabaseURL: connString,
		Keystore:    password,
	}
	tmpl, err := template.New("toml").Parse(dbTmpl)
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	err = tmpl.Execute(&output, config)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}

func writeDefaultSecrets(pgOut *postgres.Output) (*os.File, error) {
	secretsOverrides, err := generateSecretsConfig(pgOut.DockerInternalURL, DefaultTestKeystorePassword)
	if err != nil {
		return nil, err
	}
	return writeToFile(secretsOverrides, "secrets.toml")
}

func writeDefaultConfig(in *Input) (*os.File, error) {
	cfg, err := generateDefaultConfig(in)
	if err != nil {
		return nil, err
	}
	return writeToFile(cfg, "config.toml")
}

// writeToFile writes the provided data string to a specified filepath and returns the file and any error encountered.
func writeToFile(data, filePath string) (*os.File, error) {
	file, err := os.CreateTemp("", filePath)
	if err != nil {
		return nil, err
	}
	_, err = file.WriteString(data)
	if err != nil {
		return nil, err
	}
	return file, nil
}
