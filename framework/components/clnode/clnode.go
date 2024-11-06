package clnode

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

const (
	Port    = "6688"
	P2PPort = "6690"
)

var (
	once = &sync.Once{}
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
	Image                   string   `toml:"image" validate:"required"`
	Name                    string   `toml:"name"`
	DockerFilePath          string   `toml:"docker_file"`
	DockerContext           string   `toml:"docker_ctx"`
	DockerImageName         string   `toml:"docker_image_name"`
	PullImage               bool     `toml:"pull_image"`
	CapabilitiesBinaryPaths []string `toml:"capabilities"`
	CapabilityContainerDir  string   `toml:"capabilities_container_dir"`
	TestConfigOverrides     string   `toml:"test_config_overrides"`
	UserConfigOverrides     string   `toml:"user_config_overrides"`
	TestSecretsOverrides    string   `toml:"test_secrets_overrides"`
	UserSecretsOverrides    string   `toml:"user_secrets_overrides"`
	HTTPPort                int      `toml:"port"`
	P2PPort                 int      `toml:"p2p_port"`
}

// Output represents Chainlink node output, nodes and databases connection URLs
type Output struct {
	UseCache   bool             `toml:"use_cache"`
	Node       *NodeOut         `toml:"node"`
	PostgreSQL *postgres.Output `toml:"postgresql"`
}

// NodeOut is CL node container output, URLs to connect
type NodeOut struct {
	HostURL      string `toml:"url"`
	HostP2PURL   string `toml:"p2p_url"`
	DockerURL    string `toml:"docker_internal_url"`
	DockerP2PUrl string `toml:"p2p_docker_internal_url"`
}

// NewNodeWithDB create a new Chainlink node with some image:tag and one or several configs
// see config params: TestConfigOverrides, UserConfigOverrides, etc
func NewNodeWithDB(in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
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
		UseCache:   true,
		Node:       nodeOut,
		PostgreSQL: pgOut,
	}
	in.Out = out
	return out, nil
}

func NewNode(in *Input, pgOut *postgres.Output) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	nodeOut, err := newNode(in, pgOut)
	if err != nil {
		return nil, err
	}
	out := &Output{
		UseCache:   true,
		Node:       nodeOut,
		PostgreSQL: pgOut,
	}
	in.Out = out
	return out, nil
}

func newNode(in *Input, pgOut *postgres.Output) (*NodeOut, error) {
	ctx := context.Background()

	passwordPath, err := WriteTmpFile(DefaultPasswordTxt, "password.txt")
	apiCredentialsPath, err := WriteTmpFile(DefaultAPICredentials, "apicredentials")
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
	overridesFile, err := WriteTmpFile(in.Node.TestConfigOverrides, "overrides.toml")
	if err != nil {
		return nil, err
	}
	secretsOverridesFile, err := WriteTmpFile(in.Node.TestSecretsOverrides, "secrets-overrides.toml")
	if err != nil {
		return nil, err
	}
	userOverridesFile, err := WriteTmpFile(in.Node.UserConfigOverrides, "user-overrides.toml")
	if err != nil {
		return nil, err
	}
	userSecretsOverridesFile, err := WriteTmpFile(in.Node.UserSecretsOverrides, "user-secrets-overrides.toml")
	if err != nil {
		return nil, err
	}

	httpPort := fmt.Sprintf("%s/tcp", Port)
	p2pPort := fmt.Sprintf("%s/udp", P2PPort)
	var containerName string
	if in.Node.Name != "" {
		containerName = in.Node.Name
	} else {
		containerName = framework.DefaultTCName("node")
	}

	req := tc.ContainerRequest{
		AlwaysPullImage: in.Node.PullImage,
		Image:           in.Node.Image,
		Name:            containerName,
		Labels:          framework.DefaultTCLabels(),
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		ExposedPorts: []string{httpPort, p2pPort},
		Entrypoint: []string{
			"/bin/sh", "-c",
			"chainlink -c /config/config -c /config/overrides -c /config/user-overrides -s /config/secrets -s /config/secrets-overrides -s /config/user-secrets-overrides node start -d -p /config/node_password -a /config/apicredentials",
		},
		WaitingFor: wait.ForLog("Listening and serving HTTP").WithStartupTimeout(2 * time.Minute),
	}
	if in.Node.HTTPPort != 0 && in.Node.P2PPort != 0 {
		req.HostConfigModifier = func(h *container.HostConfig) {
			h.PortBindings = nat.PortMap{
				"6688/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%d/tcp", in.Node.HTTPPort),
					},
				},
				"6690/udp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%d/udp", in.Node.P2PPort),
					},
				},
			}
		}
	}
	files := []tc.ContainerFile{
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
	}
	if in.Node.CapabilityContainerDir == "" {
		in.Node.CapabilityContainerDir = "/home/capabilities"
	}
	for _, cp := range in.Node.CapabilitiesBinaryPaths {
		cpPath := filepath.Base(cp)
		framework.L.Info().Any("Path", cpPath).Str("Binary", cpPath).Msg("Copying capability binary")
		files = append(files, tc.ContainerFile{
			HostFilePath:      cp,
			ContainerFilePath: filepath.Join(in.Node.CapabilityContainerDir, cpPath),
			FileMode:          0777,
		})
	}
	req.Files = append(req.Files, files...)
	if req.Image != "" && (in.Node.DockerFilePath != "" || in.Node.DockerContext != "") {
		return nil, errors.New("you provided both 'image' and one of 'docker_file', 'docker_ctx' fields. Please provide either 'image' or params to build a local one")
	}
	if req.Image == "" {
		req.Image, err = framework.RebuildDockerImage(once, in.Node.DockerFilePath, in.Node.DockerContext, in.Node.DockerImageName)
		if err != nil {
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
	var (
		mp    nat.Port
		mpP2P nat.Port
	)
	if in.Node.HTTPPort != 0 && in.Node.P2PPort != 0 {
		mp = nat.Port(fmt.Sprintf("%d/tcp", in.Node.HTTPPort))
		mpP2P = nat.Port(fmt.Sprintf("%d/udp", in.Node.P2PPort))
	} else {
		mp, err = c.MappedPort(ctx, nat.Port(httpPort))
		if err != nil {
			return nil, err
		}
		mpP2P, err = c.MappedPort(ctx, nat.Port(p2pPort))
		if err != nil {
			return nil, err
		}
	}

	return &NodeOut{
		HostURL:      fmt.Sprintf("http://%s:%s", host, mp.Port()),
		HostP2PURL:   fmt.Sprintf("http://%s:%s", host, mpP2P.Port()),
		DockerURL:    fmt.Sprintf("http://%s:%s", containerName, Port),
		DockerP2PUrl: fmt.Sprintf("http://%s:%s", containerName, P2PPort),
	}, nil
}

type DefaultCLNodeConfig struct {
	HTTPPort      string
	SecureCookies bool
}

func generateDefaultConfig(in *Input) (string, error) {
	config := DefaultCLNodeConfig{
		HTTPPort:      Port,
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
	return WriteTmpFile(secretsOverrides, "secrets.toml")
}

func writeDefaultConfig(in *Input) (*os.File, error) {
	cfg, err := generateDefaultConfig(in)
	if err != nil {
		return nil, err
	}
	return WriteTmpFile(cfg, "config.toml")
}

// WriteTmpFile writes the provided data string to a specified filepath and returns the file and any error encountered.
func WriteTmpFile(data, filePath string) (*os.File, error) {
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
