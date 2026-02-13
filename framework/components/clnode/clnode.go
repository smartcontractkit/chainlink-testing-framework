package clnode

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/pods"

	v1 "k8s.io/api/core/v1"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
)

const (
	DefaultHTTPPort        = "6688"
	DefaultP2PPort         = "6690"
	DefaultDebuggerPort    = 40000
	TmpImageName           = "chainlink-tmp:latest"
	CustomPortSeparator    = ":"
	DefaultCapabilitiesDir = "/usr/local/bin"
	ConfigVolumeName       = "clnode-config"
	HomeVolumeName         = "clnode-home"
)

var once = &sync.Once{}

// Input represents Chainlink node input
type Input struct {
	// NoDNS whether to allow DNS in Docker containers or not, useful for isolating containers from network if set to 'false'
	NoDNS bool `toml:"no_dns" comment:"whether to allow DNS in Docker containers or not, useful for isolating containers from network if set to 'false'"`
	// DbInput PostgreSQL database configuration
	DbInput *postgres.Input `toml:"db" validate:"required" comment:"PostgreSQL database configuration"`
	// Node Chainlink node configuration
	Node *NodeInput `toml:"node" validate:"required" comment:"Chainlink node configuration"`
	// Out Chainlink node configuration output
	Out *Output `toml:"out" comment:"Chainlink node configuration output"`
}

// NodeInput is CL nod container inputs
type NodeInput struct {
	// Image Chainlink node Docker image in format $registry:$tag
	Image string `toml:"image" validate:"required" comment:"Chainlink node Docker image in format $registry:$tag"`
	// Name Chainlink node Docker container name
	Name string `toml:"name" comment:"Chainlink node Docker container name"`
	// DockerFilePath Docker file path to rebuild, relative to 'docker_ctx' field path
	DockerFilePath string `toml:"docker_file" comment:"Docker file path to rebuild, relative to 'docker_ctx' field path"`
	// DockerContext Docker build context path
	DockerContext string `toml:"docker_ctx" comment:"Docker build context path"`
	// DockerBuildArgs Docker build args
	DockerBuildArgs map[string]string `toml:"docker_build_args" comment:"Docker build args in format key = value or map format, ex.: \"CL_IS_PROD_BUILD\" = \"false\" "`
	// PullImage whether to pull Docker image or not
	PullImage bool `toml:"pull_image" comment:"Whether to pull Docker image or not"`
	// CapabilitiesBinaryPaths Chainlink CRE capabilities paths for WASM binaries
	CapabilitiesBinaryPaths []string `toml:"capabilities" comment:"Chainlink CRE capabilities paths for WASM binaries"`
	// CapabilityContainerDir path to capabilities inside Docker container (capabilities are copied inside container from local path)
	CapabilityContainerDir string `toml:"capabilities_container_dir" comment:"path to capabilities inside Docker container (capabilities are copied inside container from local path)"`
	// TestConfigOverrides node config overrides field for programmatic usage in tests
	TestConfigOverrides string `toml:"test_config_overrides" comment:"node config overrides field for programmatic usage in tests"`
	// UserConfigOverrides node config overrides field for manual overrides from env.toml configs
	UserConfigOverrides string `toml:"user_config_overrides" comment:"node config overrides field for manual overrides from env.toml configs"`
	// TestSecretsOverrides node secrets config overrides field for programmatic usage in tests
	TestSecretsOverrides string `toml:"test_secrets_overrides" comment:"node secrets config overrides field for programmatic usage in tests"`
	// UserSecretsOverrides node secrets config overrides field for manual overrides from env.toml configs
	UserSecretsOverrides string `toml:"user_secrets_overrides" comment:"node secrets config overrides field for manual overrides from env.toml configs"`
	// HTTPPort Chainlink node API HTTP port
	HTTPPort int `toml:"port" comment:"Chainlink node API HTTP port"`
	// P2PPort Chainlink node P2P port
	P2PPort int `toml:"p2p_port" comment:"Chainlink node P2P port"`
	// CustomPorts Custom ports pairs in format $host_port_number:$docker_port_number
	CustomPorts []string `toml:"custom_ports" comment:"Custom ports pairs in format $host_port_number:$docker_port_number"`
	// DebuggerPort Delve debugger port
	DebuggerPort int `toml:"debugger_port" comment:"Delve debugger port"`
	// ContainerResources Docker container resources
	ContainerResources *framework.ContainerResources `toml:"resources" comment:"Docker container resources"`
	// EnvVars Docker container environment variables
	EnvVars map[string]string `toml:"env_vars" comment:"Docker container environment variables"`
}

// Output represents Chainlink node output, nodes and databases connection URLs
type Output struct {
	// UseCache Whether to respect caching or not, if cache = true component won't be deployed again
	UseCache bool `toml:"use_cache" comment:"Whether to respect caching or not, if cache = true component won't be deployed again"`
	// Node Chainlink node config output
	Node *NodeOut `toml:"node" comment:"Chainlink node config output"`
	// PostgreSQL PostgreSQL config output
	PostgreSQL *postgres.Output `toml:"postgresql" comment:"PostgreSQL config output"`
}

// NodeOut is CL node container output, URLs to connect
type NodeOut struct {
	// APIAuthUser user name for basic login/password authorization in Chainlink node
	APIAuthUser string `toml:"api_auth_user" comment:"User name for basic login/password authorization in Chainlink node"`
	// APIAuthPassword password for basic login/password authorization in Chainlink node
	APIAuthPassword string `toml:"api_auth_password" comment:"Password for basic login/password authorization in Chainlink node"`
	// ContainerName node Docker contaienr name
	ContainerName string `toml:"container_name" comment:"Node Docker contaner name"`
	// ExternalURL node external API HTTP URL
	ExternalURL string `toml:"url" comment:"Node external API HTTP URL"`
	// InternalURL node internal API HTTP URL
	InternalURL string `toml:"internal_url" comment:"Node internal API HTTP URL"`
	// InternalP2PUrl node internal P2P URL
	InternalP2PUrl string `toml:"p2p_internal_url" comment:"Node internal P2P URL"`
	// InternalIP node internal IP
	InternalIP string `toml:"internal_ip" comment:"Node internal IP"`
	// K8sService is a Kubernetes service spec used to connect locally
	K8sService *v1.Service `toml:"k8s_service" comment:"Kubernetes service spec used to connect locally"`
}

// NewNodeWithDB create a new Chainlink node with some image:tag and one or several configs
// see config params: TestConfigOverrides, UserConfigOverrides, etc
func NewNodeWithDB(in *Input) (*Output, error) {
	return NewNodeWithDBAndContext(context.Background(), in)
}

// NewNodeWithDBAndContext create a new Chainlink node with some image:tag and one or several configs
// see config params: TestConfigOverrides, UserConfigOverrides, etc
func NewNodeWithDBAndContext(ctx context.Context, in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	pgOut, err := postgres.NewWithContext(ctx, in.DbInput)
	if err != nil {
		return nil, err
	}
	nodeOut, err := newNode(ctx, in, pgOut)
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
	return NewNodeWithContext(context.Background(), in, pgOut)
}

func NewNodeWithContext(ctx context.Context, in *Input, pgOut *postgres.Output) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	nodeOut, err := newNode(ctx, in, pgOut)
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

func generateEntryPoint() []string {
	entrypoint := []string{
		"/bin/sh", "-c",
	}
	if os.Getenv("CTF_CLNODE_DLV") == "true" {
		entrypoint = append(entrypoint, "dlv exec /usr/local/bin/chainlink --continue --listen=0.0.0.0:40000 --headless=true --api-version=2 --accept-multiclient -- -c /config/config -c /config/overrides -c /config/user-overrides -s /config/secrets -s /config/secrets-overrides -s /config/user-secrets-overrides node start -d -p /config/node_password -a /config/apicredentials")
	} else {
		entrypoint = append(entrypoint, "chainlink -c /config/config -c /config/overrides -c /config/user-overrides -s /config/secrets -s /config/secrets-overrides -s /config/user-secrets-overrides node start -d -p /config/node_password -a /config/apicredentials")
	}
	return entrypoint
}

// natPortsToK8sFormat transforms nat.PortMap
// to Pods port pair format: $external_port:$internal_port
func natPortsToK8sFormat(nat nat.PortMap) []string {
	out := make([]string, 0)
	for port, portBinding := range nat {
		for _, b := range portBinding {
			out = append(out, fmt.Sprintf("%s:%s", b.HostPort, strconv.Itoa(port.Int())))
		}
	}
	return out
}

// generatePortBindings generates exposed ports and port bindings
// exposes default CL node port
// exposes custom_ports in format "host:docker" or map 1-to-1 if only "host" port is provided
func generatePortBindings(in *Input) ([]string, nat.PortMap, error) {
	httpPort := fmt.Sprintf("%s/tcp", DefaultHTTPPort)
	p2pPort := fmt.Sprintf("%s/udp", DefaultP2PPort)
	exposedPorts := []string{httpPort, p2pPort}
	portBindings := nat.PortMap{
		nat.Port(httpPort): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(in.Node.HTTPPort),
			},
		},
		nat.Port(p2pPort): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(in.Node.P2PPort),
			},
		},
	}
	if os.Getenv("CTF_CLNODE_DLV") == "true" {
		innerDebuggerPort := fmt.Sprintf("%d/tcp", DefaultDebuggerPort)
		portBindings[nat.Port(innerDebuggerPort)] = append(portBindings[nat.Port(innerDebuggerPort)], nat.PortBinding{
			HostIP:   "0.0.0.0",
			HostPort: strconv.Itoa(in.Node.DebuggerPort),
		})
		exposedPorts = append(exposedPorts, strconv.Itoa(DefaultDebuggerPort))
	}
	customPorts := make([]string, 0)
	for _, p := range in.Node.CustomPorts {
		if strings.Contains(p, CustomPortSeparator) {
			pp := strings.Split(p, CustomPortSeparator)
			if len(pp) != 2 {
				return nil, nil, fmt.Errorf("custom_ports has ':' but you must provide both ports, you provided: %s", pp)
			}
			customPorts = append(customPorts, fmt.Sprintf("%s/tcp", pp[1]))

			dockerPort := nat.Port(fmt.Sprintf("%s/tcp", pp[1]))
			hostPort := pp[0]
			portBindings[dockerPort] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			}
		} else {
			customPorts = append(customPorts, fmt.Sprintf("%s/tcp", p))

			dockerPort := nat.Port(fmt.Sprintf("%s/tcp", p))
			hostPort := p
			portBindings[dockerPort] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			}
		}
	}
	exposedPorts = append(exposedPorts, customPorts...)
	return exposedPorts, portBindings, nil
}

func newNode(ctx context.Context, in *Input, pgOut *postgres.Output) (*NodeOut, error) {
	passwordPath, err := WriteTmpFile(DefaultPasswordTxt, "password.txt")
	if err != nil {
		return nil, err
	}
	apiCredentialsPath, err := WriteTmpFile(DefaultAPICredentials, "apicredentials")
	if err != nil {
		return nil, err
	}
	cfg, err := generateDefaultConfig()
	if err != nil {
		return nil, err
	}
	cfgPath, err := WriteTmpFile(cfg, "config.toml")
	if err != nil {
		return nil, err
	}

	secretsData, err := generateSecretsConfig(pgOut.InternalURL, DefaultTestKeystorePassword)
	if err != nil {
		return nil, err
	}
	secretsPath, err := WriteTmpFile(secretsData, "secrets.toml")
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

	var containerName string
	if in.Node.Name != "" {
		containerName = in.Node.Name
	} else {
		containerName = framework.DefaultTCName("node")
	}

	exposedPorts, portBindings, err := generatePortBindings(in)
	if err != nil {
		return nil, err
	}

	defaultHTTPPortInt, err := strconv.Atoi(DefaultHTTPPort)
	if err != nil {
		return nil, err
	}

	// k8s deployment
	if pods.K8sEnabled() {
		_, svc, err := pods.Run(ctx, &pods.Config{
			Pods: []*pods.PodConfig{
				{
					Name:     pods.Ptr(containerName),
					Image:    pods.Ptr(in.Node.Image),
					Env:      pods.EnvsFromMap(in.Node.EnvVars),
					Requests: pods.ResourcesMedium(),
					Limits:   pods.ResourcesMedium(),
					Ports:    natPortsToK8sFormat(portBindings),
					ContainerSecurityContext: &v1.SecurityContext{
						// these are specific things we need for staging cluster
						RunAsNonRoot: pods.Ptr(true),
						RunAsUser:    pods.Ptr[int64](14933),
						RunAsGroup:   pods.Ptr[int64](999),
					},
					ReadinessProbe: pods.TCPReadyProbe(defaultHTTPPortInt),
					ConfigMap: map[string]string{
						"config.toml":         cfg,
						"overrides.toml":      in.Node.TestConfigOverrides,
						"user-overrides.toml": in.Node.UserConfigOverrides,
						"node_password":       DefaultPasswordTxt,
						"apicredentials": fmt.Sprintf(`%s
			%s`, DefaultAPIUser, DefaultAPIPassword),
					},
					ConfigMapMountPath: map[string]string{
						"config.toml":         "/config/config",
						"overrides.toml":      "/config/overrides",
						"user-overrides.toml": "/config/user-overrides",
						"node_password":       "/config/node_password",
						"apicredentials":      "/config/apicredentials",
					},
					Secrets: map[string]string{
						"secrets.toml":                secretsData,
						"secrets-overrides.toml":      in.Node.TestSecretsOverrides,
						"secrets-user-overrides.toml": in.Node.UserSecretsOverrides,
					},
					SecretsMountPath: map[string]string{
						"secrets.toml":                "/config/secrets",
						"secrets-overrides.toml":      "/config/secrets-overrides",
						"secrets-user-overrides.toml": "/config/user-secrets-overrides",
					},
					Command: pods.Ptr("chainlink -c /config/config -c /config/overrides -c /config/user-overrides -s /config/secrets -s /config/secrets-overrides -s /config/user-secrets-overrides node start -d -p /config/node_password -a /config/apicredentials"),
				},
			},
		})
		if err != nil {
			return nil, err
		}
		return &NodeOut{
			APIAuthUser:     DefaultAPIUser,
			APIAuthPassword: DefaultAPIPassword,
			ContainerName:   containerName,
			ExternalURL:     fmt.Sprintf("http://%s:%d", "localhost", in.Node.HTTPPort),
			InternalURL:     fmt.Sprintf("http://%s:%s", containerName, DefaultHTTPPort),
			InternalP2PUrl:  fmt.Sprintf("http://%s:%s", containerName, DefaultP2PPort),
			K8sService:      svc,
		}, nil
	}
	// local deployment
	req := tc.ContainerRequest{
		AlwaysPullImage: in.Node.PullImage,
		Image:           in.Node.Image,
		Name:            containerName,
		Labels:          framework.DefaultTCLabels(),
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		Env:          in.Node.EnvVars,
		ExposedPorts: exposedPorts,
		Entrypoint:   generateEntryPoint(),
		WaitingFor: wait.ForHTTP("/").
			WithPort(DefaultHTTPPort).
			WithStartupTimeout(3 * time.Minute).
			WithPollInterval(200 * time.Millisecond),
		Mounts: tc.ContainerMounts{
			{
				// various configuration files
				Source: tc.GenericVolumeMountSource{
					Name: ConfigVolumeName + "-" + in.Node.Name,
				},
				Target: "/config",
			},
			{
				// kv store of the OCR jobs and other state files are stored
				// in the user's home instead of the DB
				Source: tc.GenericVolumeMountSource{
					Name: HomeVolumeName + "-" + in.Node.Name,
				},
				Target: "/home/chainlink",
			},
		},
	}
	if in.Node.HTTPPort != 0 && in.Node.P2PPort != 0 {
		req.HostConfigModifier = func(h *container.HostConfig) {
			framework.NoDNS(in.NoDNS, h)
			h.PortBindings = portBindings
			framework.ResourceLimitsFunc(h, in.Node.ContainerResources)
		}
	}
	files := []tc.ContainerFile{
		{
			HostFilePath:      cfgPath.Name(),
			ContainerFilePath: "/config/config",
			FileMode:          0o644,
		},
		{
			HostFilePath:      secretsPath.Name(),
			ContainerFilePath: "/config/secrets",
			FileMode:          0o644,
		},
		{
			HostFilePath:      overridesFile.Name(),
			ContainerFilePath: "/config/overrides",
			FileMode:          0o644,
		},
		{
			HostFilePath:      userOverridesFile.Name(),
			ContainerFilePath: "/config/user-overrides",
			FileMode:          0o644,
		},
		{
			HostFilePath:      secretsOverridesFile.Name(),
			ContainerFilePath: "/config/secrets-overrides",
			FileMode:          0o644,
		},
		{
			HostFilePath:      userSecretsOverridesFile.Name(),
			ContainerFilePath: "/config/user-secrets-overrides",
			FileMode:          0o644,
		},
		{
			HostFilePath:      passwordPath.Name(),
			ContainerFilePath: "/config/node_password",
			FileMode:          0o644,
		},
		{
			HostFilePath:      apiCredentialsPath.Name(),
			ContainerFilePath: "/config/apicredentials",
			FileMode:          0o644,
		},
	}
	if in.Node.CapabilityContainerDir == "" {
		in.Node.CapabilityContainerDir = DefaultCapabilitiesDir
	}
	for _, cp := range in.Node.CapabilitiesBinaryPaths {
		cpPath := filepath.Base(cp)
		framework.L.Info().Any("Path", cpPath).Str("Binary", cpPath).Msg("Copying capability binary")
		files = append(files, tc.ContainerFile{
			HostFilePath:      cp,
			ContainerFilePath: filepath.Join(in.Node.CapabilityContainerDir, cpPath),
			FileMode:          0o777,
		})
	}
	req.Files = append(req.Files, files...)
	if req.Image != "" && (in.Node.DockerFilePath != "" || in.Node.DockerContext != "") {
		return nil, errors.New("you provided both 'image' and one of 'docker_file', 'docker_ctx' fields. Please provide either 'image' or params to build a local one")
	}
	if req.Image == "" {
		req.Image = TmpImageName
		if err := framework.BuildImageOnce(once, in.Node.DockerContext, in.Node.DockerFilePath, req.Image, in.Node.DockerBuildArgs); err != nil {
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
	ip, err := c.ContainerIP(ctx)
	if err != nil {
		return nil, err
	}
	host, err := framework.GetHostWithContext(ctx, c)
	if err != nil {
		return nil, err
	}

	mp := nat.Port(fmt.Sprintf("%d/tcp", in.Node.HTTPPort))

	return &NodeOut{
		APIAuthUser:     DefaultAPIUser,
		APIAuthPassword: DefaultAPIPassword,
		ContainerName:   containerName,
		ExternalURL:     fmt.Sprintf("http://%s:%s", host, mp.Port()),
		InternalURL:     fmt.Sprintf("http://%s:%s", containerName, DefaultHTTPPort),
		InternalP2PUrl:  fmt.Sprintf("http://%s:%s", containerName, DefaultP2PPort),
		InternalIP:      ip,
	}, nil
}

type DefaultCLNodeConfig struct {
	HTTPPort      string
	SecureCookies bool
}

func generateDefaultConfig() (string, error) {
	config := DefaultCLNodeConfig{
		HTTPPort:      DefaultHTTPPort,
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
