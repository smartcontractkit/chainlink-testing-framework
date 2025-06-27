package chipingressset

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	ChipIngress *ChipIngressOutput
	RedPanda    *RedPandaOutput
}

type ChipIngressOutput struct {
	GRPCInternalURL string
	GRPCExternalURL string
}

type RedPandaOutput struct {
	SchemaRegistryInternalURL string
	SchemaRegistryExternalURL string
	KafkaInternalURL          string
	KafkaExternalURL          string
	ConsoleExternalURL        string
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	upErr := stack.
		WithEnv(map[string]string{
			"BASIC_AUTH_ENABLED": "false",
			"BASIC_AUTH_PREFIX":  "",
		}).
		Up(ctx)

	if upErr != nil {
		return nil, errors.Wrap(upErr, "failed to start stack for Chip Ingress")
	}

	stack.WaitForService(DEFAULT_CHIP_INGRESS_SERVICE_NAME,
		wait.ForAll(
			wait.ForLog("GRPC server is live").WithPollInterval(200*time.Millisecond),
			wait.ForListeningPort(DEFAULT_CHIP_INGRESS_GRPC_PORT),
		).WithDeadline(1*time.Minute),
	).WaitForService(DEFAULT_RED_PANDA_SERVICE_NAME,
		wait.ForListeningPort(DEFAULT_RED_PANDA_KAFKA_PORT).WithStartupTimeout(1*time.Minute),
	).WaitForService(DEFAULT_RED_PANDA_CONSOLE_SERVICE_NAME,
		wait.ForListeningPort(DEFAULT_RED_PANDA_CONSOLE_PORT).WithStartupTimeout(1*time.Minute),
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

	timeout := time.After(1 * time.Minute)
	tick := time.Tick(500 * time.Millisecond)

	// so let's try to connect to a Docker network a couple of times, there must be a race condition in Docker
	// and even when network sandbox has been created and container is running, this call can still fail
	// retrying is simpler than trying to figure out how to correctly wait for the network sandbox to be ready
	var connectNetwork = func(networkName string) error {
		for {
			select {
			case <-timeout:
				return fmt.Errorf("timeout while trying to connect chip-ingress to default network")
			case <-tick:
				if networkErr := cli.NetworkConnect(
					ctx,
					networkName,
					chipIngressContainer.ID,
					&networkTypes.EndpointSettings{
						Aliases: []string{identifier},
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

	networks := []string{framework.DefaultNetworkName}
	networks = append(networks, in.ExtraDockerNetworks...)

	for _, networkName := range networks {
		framework.L.Debug().Msgf("Connecting chip-ingress to %s network", networkName)
		if connectErr := connectNetwork(networkName); connectErr != nil {
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
	}

	// get hosts and ports for chip-ingress and redpanda
	chipIngressExternalHost, chipIngressExternalHostErr := chipIngressContainer.Host(ctx)
	if chipIngressExternalHostErr != nil {
		return nil, errors.Wrap(chipIngressExternalHostErr, "failed to get host for Chip Ingress")
	}
	chipIngressExternalPort, chipIngressExternalPortErr := chipIngressContainer.MappedPort(ctx, DEFAULT_CHIP_INGRESS_GRPC_PORT)
	if chipIngressExternalPortErr != nil {
		return nil, errors.Wrap(chipIngressExternalPortErr, "failed to get mapped port for Chip Ingress")
	}

	redpandaContainer, redpandaErr := stack.ServiceContainer(ctx, DEFAULT_RED_PANDA_SERVICE_NAME)
	if redpandaErr != nil {
		return nil, errors.Wrap(redpandaErr, "failed to get redpanda container")
	}

	redpandaExternalHost, redpandaExternalHostErr := redpandaContainer.Host(ctx)
	if redpandaExternalHostErr != nil {
		return nil, errors.Wrap(redpandaExternalHostErr, "failed to get host for Red Panda")
	}
	redpandaExternalKafkaPort, redpandaExternalKafkaPortErr := redpandaContainer.MappedPort(ctx, DEFAULT_RED_PANDA_KAFKA_PORT)
	if redpandaExternalKafkaPortErr != nil {
		return nil, errors.Wrap(redpandaExternalKafkaPortErr, "failed to get mapped port for Red Panda")
	}
	redpandaExternalSchemaRegistryPort, redpandaExternalSchemaRegistryPortErr := redpandaContainer.MappedPort(ctx, DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT)
	if redpandaExternalSchemaRegistryPortErr != nil {
		return nil, errors.Wrap(redpandaExternalSchemaRegistryPortErr, "failed to get mapped port for Red Panda")
	}

	redpandaConsoleContainer, redpandaConsoleErr := stack.ServiceContainer(ctx, DEFAULT_RED_PANDA_CONSOLE_SERVICE_NAME)
	if redpandaConsoleErr != nil {
		return nil, errors.Wrap(redpandaConsoleErr, "failed to get redpanda-console container")
	}

	redpandaExternalConsoleHost, redpandaExternalConsoleHostErr := redpandaConsoleContainer.Host(ctx)
	if redpandaExternalConsoleHostErr != nil {
		return nil, errors.Wrap(redpandaExternalConsoleHostErr, "failed to get host for Red Panda Console")
	}
	redpandaExternalConsolePort, redpandaExternalConsolePortErr := redpandaConsoleContainer.MappedPort(ctx, DEFAULT_RED_PANDA_CONSOLE_PORT)
	if redpandaExternalConsolePortErr != nil {
		return nil, errors.Wrap(redpandaExternalConsolePortErr, "failed to get mapped port for Red Panda Console")
	}

	output := &Output{
		ChipIngress: &ChipIngressOutput{
			GRPCInternalURL: fmt.Sprintf("http://%s:%s", DEFAULT_CHIP_INGRESS_SERVICE_NAME, DEFAULT_CHIP_INGRESS_GRPC_PORT),
			GRPCExternalURL: fmt.Sprintf("http://%s:%s", chipIngressExternalHost, chipIngressExternalPort.Port()),
		},
		RedPanda: &RedPandaOutput{
			SchemaRegistryInternalURL: fmt.Sprintf("http://%s:%s", DEFAULT_RED_PANDA_SERVICE_NAME, DEFAULT_RED_PANDA_SCHEMA_REGISTRY_PORT),
			SchemaRegistryExternalURL: fmt.Sprintf("http://%s:%s", redpandaExternalHost, redpandaExternalSchemaRegistryPort.Port()),
			KafkaInternalURL:          fmt.Sprintf("%s:%s", DEFAULT_RED_PANDA_SERVICE_NAME, DEFAULT_RED_PANDA_KAFKA_PORT),
			KafkaExternalURL:          fmt.Sprintf("%s:%s", redpandaExternalHost, redpandaExternalKafkaPort.Port()),
			ConsoleExternalURL:        fmt.Sprintf("http://%s:%s", redpandaExternalConsoleHost, redpandaExternalConsolePort.Port()),
		},
	}

	framework.L.Info().Msg("Chip Ingress stack start")

	return output, nil
}

func composeFilePath(rawFilePath string) (string, error) {
	// if it's not a URL, return it as is and assume it's a local file
	if !strings.HasPrefix(rawFilePath, "http") {
		return rawFilePath, nil
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
