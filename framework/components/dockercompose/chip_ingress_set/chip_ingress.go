package chipingressset

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	networkTypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Output struct {
	ChipIngress *ChipIngressOutput `toml:"chip_ingress"`
	RedPanda    *RedPandaOutput    `toml:"redpanda"`
}

type ChipIngressOutput struct {
	GRPCInternalURL string `toml:"grpc_internal_url"`
	GRPCExternalURL string `toml:"grpc_external_url"`
}

type RedPandaOutput struct {
	SchemaRegistryInternalURL string `toml:"schema_registry_internal_url"`
	SchemaRegistryExternalURL string `toml:"schema_registry_external_url"`
	KafkaInternalURL          string `toml:"kafka_internal_url"`
	KafkaExternalURL          string `toml:"kafka_external_url"`
	ConsoleExternalURL        string `toml:"console_external_url"`
}

type Input struct {
	ComposeFile         string   `toml:"compose_file"`
	ExtraDockerNetworks []string `toml:"extra_docker_networks"`
	Output              *Output  `toml:"output"`
	UseCache            bool     `toml:"use_cache"`
}

func defaultChipIngress(in *Input) *Input {
	if in.ComposeFile == "" {
		in.ComposeFile = "./docker-compose.yml"
	}
	return in
}

const (
	DEFAULT_STACK_NAME = "chip-ingress"

	DEFAULT_CHIP_INGRESS_GRPC_PORT    = "50051"
	DEFAULT_CHIP_INGRESS_SERVICE_NAME = "chip-ingress"

	DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT = "18081"
	DEFAULT_RED_PANDA_KAFKA_PORT           = "19092"
	DEFAULT_RED_PANDA_SERVICE_NAME         = "redpanda-0"

	DEFAULT_RED_PANDA_CONSOLE_SERVICE_NAME = "redpanda-console"
	DEFAULT_RED_PANDA_CONSOLE_PORT         = "8080"
)

func New(in *Input) (*Output, error) {
	return NewWithContext(context.Background(), in)
}

func NewWithContext(ctx context.Context, in *Input) (*Output, error) {
	if in == nil {
		return nil, errors.New("input is nil")
	}

	if in.UseCache {
		if in.Output != nil {
			return in.Output, nil
		}
	}

	in = defaultChipIngress(in)
	identifier := framework.DefaultTCName(DEFAULT_STACK_NAME)
	framework.L.Debug().Str("Compose file", in.ComposeFile).Msgf("Starting Chip Ingress stack with identifier %s", framework.DefaultTCName(DEFAULT_STACK_NAME))

	composeFilePath, fileErr := composeFilePath(in.ComposeFile)
	if fileErr != nil {
		return nil, errors.Wrap(fileErr, "failed to get compose file path")
	}

	stack, stackErr := compose.NewDockerComposeWith(
		compose.WithStackFiles(composeFilePath),
		compose.StackIdentifier(identifier),
	)
	if stackErr != nil {
		return nil, errors.Wrap(stackErr, "failed to create compose stack for Chip Ingress")
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Start the stackwith all environment variables from the host process
	// set BASIC_AUTH_ENABLED and BASIC_AUTH_PREFIX to false and empty string and allow them to be overridden by the host process
	envVars := make(map[string]string)
	envVars["BASIC_AUTH_ENABLED"] = "false"
	envVars["BASIC_AUTH_PREFIX"] = ""

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envVars[pair[0]] = pair[1]
		}
	}

	upErr := stack.
		WithEnv(envVars).
		Up(ctx)

	if upErr != nil {
		return nil, errors.Wrap(upErr, "failed to start stack for Chip Ingress")
	}

	stack.WaitForService(DEFAULT_CHIP_INGRESS_SERVICE_NAME,
		wait.ForAll(
			wait.ForLog("GRPC server is live").WithPollInterval(100*time.Millisecond),
			wait.ForListeningPort(DEFAULT_CHIP_INGRESS_GRPC_PORT).WithPollInterval(100*time.Millisecond),
			wait.NewHostPortStrategy(DEFAULT_CHIP_INGRESS_GRPC_PORT).WithPollInterval(100*time.Millisecond),
		).WithDeadline(2*time.Minute),
	).WaitForService(DEFAULT_RED_PANDA_SERVICE_NAME,
		wait.ForAll(
			wait.ForListeningPort(DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT).WithPollInterval(100*time.Millisecond),
			wait.NewHostPortStrategy(DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT).WithPollInterval(100*time.Millisecond),
			wait.ForHTTP("/status/ready").WithPort(DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT).WithPollInterval(100*time.Millisecond),
		).WithDeadline(2*time.Minute),
	).WaitForService(DEFAULT_RED_PANDA_CONSOLE_SERVICE_NAME,
		wait.ForAll(
			wait.ForListeningPort(DEFAULT_RED_PANDA_CONSOLE_PORT).WithPollInterval(100*time.Millisecond),
			wait.NewHostPortStrategy(DEFAULT_RED_PANDA_CONSOLE_PORT).WithPollInterval(100*time.Millisecond),
		).WithDeadline(2*time.Minute),
	)

	chipIngressContainer, ingressErr := stack.ServiceContainer(ctx, DEFAULT_CHIP_INGRESS_SERVICE_NAME)
	if ingressErr != nil {
		return nil, errors.Wrap(ingressErr, "failed to get chip-ingress container")
	}

	cli, cliErr := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if cliErr != nil {
		return nil, errors.Wrap(cliErr, "failed to create docker client")
	}
	defer cli.Close()

	// so let's try to connect to a Docker network a couple of times, there must be a race condition in Docker
	// and even when network sandbox has been created and container is running, this call can still fail
	// retrying is simpler than trying to figure out how to correctly wait for the network sandbox to be ready
	networks := []string{framework.DefaultNetworkName}
	networks = append(networks, in.ExtraDockerNetworks...)

	for _, networkName := range networks {
		framework.L.Debug().Msgf("Connecting chip-ingress to %s network", networkName)
		connectContex, connectCancel := context.WithTimeout(ctx, 30*time.Second)
		defer connectCancel()
		if connectErr := connectNetwork(connectContex, 30*time.Second, cli, chipIngressContainer.ID, networkName, identifier); connectErr != nil {
			return nil, errors.Wrapf(connectErr, "failed to connect chip-ingress to %s network", networkName)
		}
		// verify that the container is connected to framework's network
		inspected, inspectErr := cli.ContainerInspect(ctx, chipIngressContainer.ID)
		if inspectErr != nil {
			return nil, errors.Wrapf(inspectErr, "failed to inspect container %s", chipIngressContainer.ID)
		}

		_, ok := inspected.NetworkSettings.Networks[networkName]
		if !ok {
			return nil, fmt.Errorf("container %s is NOT on network %s", chipIngressContainer.ID, networkName)
		}

		framework.L.Debug().Msgf("Container %s is connected to network %s", chipIngressContainer.ID, networkName)
	}

	// get hosts and ports for chip-ingress and redpanda
	chipIngressExternalHost, chipIngressExternalHostErr := chipIngressContainer.Host(ctx)
	if chipIngressExternalHostErr != nil {
		return nil, errors.Wrap(chipIngressExternalHostErr, "failed to get host for Chip Ingress")
	}

	redpandaContainer, redpandaErr := stack.ServiceContainer(ctx, DEFAULT_RED_PANDA_SERVICE_NAME)
	if redpandaErr != nil {
		return nil, errors.Wrap(redpandaErr, "failed to get redpanda container")
	}

	redpandaExternalHost, redpandaExternalHostErr := redpandaContainer.Host(ctx)
	if redpandaExternalHostErr != nil {
		return nil, errors.Wrap(redpandaExternalHostErr, "failed to get host for Red Panda")
	}

	redpandaConsoleContainer, redpandaConsoleErr := stack.ServiceContainer(ctx, DEFAULT_RED_PANDA_CONSOLE_SERVICE_NAME)
	if redpandaConsoleErr != nil {
		return nil, errors.Wrap(redpandaConsoleErr, "failed to get redpanda-console container")
	}
	redpandaExternalConsoleHost, redpandaExternalConsoleHostErr := redpandaConsoleContainer.Host(ctx)
	if redpandaExternalConsoleHostErr != nil {
		return nil, errors.Wrap(redpandaExternalConsoleHostErr, "failed to get host for Red Panda Console")
	}

	output := &Output{
		ChipIngress: &ChipIngressOutput{
			GRPCInternalURL: fmt.Sprintf("http://%s:%s", DEFAULT_CHIP_INGRESS_SERVICE_NAME, DEFAULT_CHIP_INGRESS_GRPC_PORT),
			GRPCExternalURL: fmt.Sprintf("http://%s:%s", chipIngressExternalHost, DEFAULT_CHIP_INGRESS_GRPC_PORT),
		},
		RedPanda: &RedPandaOutput{
			SchemaRegistryInternalURL: fmt.Sprintf("http://%s:%s", DEFAULT_RED_PANDA_SERVICE_NAME, DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT),
			SchemaRegistryExternalURL: fmt.Sprintf("http://%s:%s", redpandaExternalHost, DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT),
			KafkaInternalURL:          fmt.Sprintf("%s:%s", DEFAULT_RED_PANDA_SERVICE_NAME, DEFAULT_RED_PANDA_KAFKA_PORT),
			KafkaExternalURL:          fmt.Sprintf("%s:%s", redpandaExternalHost, DEFAULT_RED_PANDA_KAFKA_PORT),
			ConsoleExternalURL:        fmt.Sprintf("http://%s:%s", redpandaExternalConsoleHost, DEFAULT_RED_PANDA_CONSOLE_PORT),
		},
	}

	in.Output = output
	in.UseCache = true
	framework.L.Info().Msg("Chip Ingress stack started")

	return output, checkSchemaRegistryReadiness(ctx, 2*time.Minute, 300*time.Millisecond, output.RedPanda.SchemaRegistryExternalURL, 3)
}

func composeFilePath(rawFilePath string) (string, error) {
	// if it's not a URL, return it as is and assume it's a local file
	if !strings.HasPrefix(rawFilePath, "http") {
		if !strings.HasPrefix(rawFilePath, "file://") {
			return "", fmt.Errorf("docker compose URI must start either with 'file://', 'http://' or 'https://', but '%s' was found", rawFilePath)
		}
		return strings.TrimPrefix(rawFilePath, "file://"), nil
	}

	resp, respErr := http.Get(rawFilePath)
	if respErr != nil {
		return "", errors.Wrap(respErr, "failed to download docker-compose file")
	}
	defer resp.Body.Close()

	tempFile, tempErr := os.CreateTemp("", "chip-ingress-docker-compose-*.yml")
	if tempErr != nil {
		return "", errors.Wrap(tempErr, "failed to create temp file")
	}
	defer tempFile.Close()

	_, copyErr := io.Copy(tempFile, resp.Body)
	if copyErr != nil {
		tempFile.Close()
		return "", errors.Wrap(copyErr, "failed to write compose file")
	}

	return tempFile.Name(), nil
}

func connectNetwork(connCtx context.Context, timeout time.Duration, dockerClient *client.Client, containerID, networkName, stackIdentifier string) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	networkCtx, networkCancel := context.WithTimeout(connCtx, timeout)
	defer networkCancel()

	for {
		select {
		case <-networkCtx.Done():
			return fmt.Errorf("timeout while trying to connect chip-ingress to default network")
		case <-ticker.C:
			if networkErr := dockerClient.NetworkConnect(
				connCtx,
				networkName,
				containerID,
				&networkTypes.EndpointSettings{
					Aliases: []string{stackIdentifier},
				},
			); networkErr != nil && !strings.Contains(networkErr.Error(), "already exists in network") {
				framework.L.Trace().Msgf("failed to connect to default network: %v", networkErr)
				continue
			}
			framework.L.Trace().Msgf("connected to %s network", networkName)
			return nil
		}
	}
}

// checkSchemaRegistryReadiness verifies that the Schema Registry answers 2xx on GET /subjects
// for minSuccessCount *consecutive* attempts, polling every `interval`, with an overall `timeout`.
func checkSchemaRegistryReadiness(ctx context.Context, timeout, interval time.Duration, registryURL string, minSuccessCount int) error {
	if minSuccessCount < 1 {
		minSuccessCount = 1
	}
	u, uErr := url.Parse(registryURL)
	if uErr != nil {
		return fmt.Errorf("parse registry URL: %w", uErr)
	}
	var pErr error
	u.Path, pErr = url.JoinPath(u.Path, "/subjects") // keeps existing base path, adds /subjects
	if pErr != nil {
		return fmt.Errorf("join /subjects path: %w", pErr)
	}

	// Fresh connection per request; prefer IPv4 to avoid ::1 races.
	tr := &http.Transport{
		DisableKeepAlives: true,
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			d := &net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}
			return d.DialContext(ctx, "tcp4", addr)
		},
		ForceAttemptHTTP2: false, // optional; stick to HTTP/1.1 for simplicity
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second, // per-request timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	t := time.NewTicker(interval)
	defer t.Stop()

	consecutive := 0
	var lastErr error

	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		// small belt-and-suspenders to ensure no reuse even if transport changes
		req.Close = true

		resp, err := client.Do(req)
		if err == nil {
			// Always drain & close.
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}

		if err == nil && resp.StatusCode/100 == 2 {
			framework.L.Debug().Msgf("schema registry ready check succeeded with status %d (%d/%d)", resp.StatusCode, consecutive+1, minSuccessCount)
			consecutive++
			if consecutive >= minSuccessCount {
				framework.L.Debug().Msg("schema registry is ready")
				return nil
			}
		} else {
			consecutive = 0
			if err != nil {
				framework.L.Debug().Msgf("schema registry ready check failed with error %v (need %d/%d consecutive successes)", err, consecutive, minSuccessCount)
				lastErr = fmt.Errorf("GET /subjects failed: %w", err)
			} else {
				framework.L.Debug().Msgf("schema registry ready check failed with error %v and status code %d (need %d/%d consecutive successes)", err, resp.StatusCode, consecutive, minSuccessCount)
				lastErr = fmt.Errorf("GET /subjects status %d", resp.StatusCode)
			}
		}

		select {
		case <-ctx.Done():
			if lastErr == nil {
				lastErr = ctx.Err()
			}
			return fmt.Errorf("schema registry not ready after %s; needed %d consecutive successes (got %d): %w",
				timeout, minSuccessCount, consecutive, lastErr)
		case <-t.C:
			framework.L.Debug().Msg("schema registry not ready yet, retrying...")
			// poll again
		}
	}
}
