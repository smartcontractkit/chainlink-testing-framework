package billing_platform_service

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
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Output struct {
	BillingPlatformService *BillingPlatformServiceOutput
	Postgres               *PostgresOutput
}

type BillingPlatformServiceOutput struct {
	GRPCInternalURL string
	GRPCExternalURL string
}

type PostgresOutput struct {
}

type Input struct {
	ComposeFile         string   `toml:"compose_file"`
	ExtraDockerNetworks []string `toml:"extra_docker_networks"`
	Output              *Output  `toml:"output"`
	UseCache            bool     `toml:"use_cache"`
}

func defaultBillingPlatformService(in *Input) *Input {
	if in.ComposeFile == "" {
		in.ComposeFile = "./docker-compose.yml"
	}
	return in
}

const (
	DEFAULT_STACK_NAME = "billing-platform-service"

	DEFAULT_BILLING_PLATFORM_SERVICE_GRPC_PORT    = "2022"
	DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME = "billing-platform-service"
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

	in = defaultBillingPlatformService(in)
	identifier := framework.DefaultTCName(DEFAULT_STACK_NAME)
	framework.L.Debug().Str("Compose file", in.ComposeFile).
		Msgf("Starting Billing Platform Service stack with identifier %s",
			framework.DefaultTCName(DEFAULT_STACK_NAME))

	cFilePath, fileErr := composeFilePath(in.ComposeFile)
	if fileErr != nil {
		return nil, errors.Wrap(fileErr, "failed to get compose file path")
	}

	stack, stackErr := compose.NewDockerComposeWith(
		compose.WithStackFiles(cFilePath),
		compose.StackIdentifier(identifier),
	)
	if stackErr != nil {
		return nil, errors.Wrap(stackErr, "failed to create compose stack for Billing Platform Service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	upErr := stack.Up(ctx)

	if upErr != nil {
		return nil, errors.Wrap(upErr, "failed to start stack for Billing Platform Service")
	}

	stack.WaitForService(DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME,
		wait.ForAll(
			wait.ForLog("GRPC server is live").WithPollInterval(200*time.Millisecond),
			wait.ForListeningPort(DEFAULT_BILLING_PLATFORM_SERVICE_GRPC_PORT),
		).WithDeadline(1*time.Minute),
	)

	billingContainer, billingErr := stack.ServiceContainer(ctx, DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME)
	if billingErr != nil {
		return nil, errors.Wrap(billingErr, "failed to get billing-platform-service container")
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
		framework.L.Debug().Msgf("Connecting billing-platform-service to %s network", networkName)
		connectContex, connectCancel := context.WithTimeout(ctx, 30*time.Second)
		defer connectCancel()
		if connectErr := connectNetwork(connectContex, 30*time.Second, cli, billingContainer.ID, networkName, identifier); connectErr != nil {
			return nil, errors.Wrapf(connectErr, "failed to connect billing-platform-service to %s network", networkName)
		}
		// verify that the container is connected to framework's network
		inspected, inspectErr := cli.ContainerInspect(ctx, billingContainer.ID)
		if inspectErr != nil {
			return nil, errors.Wrapf(inspectErr, "failed to inspect container %s", billingContainer.ID)
		}

		_, ok := inspected.NetworkSettings.Networks[networkName]
		if !ok {
			return nil, fmt.Errorf("container %s is NOT on network %s", billingContainer.ID, networkName)
		}

		framework.L.Debug().Msgf("Container %s is connected to network %s", billingContainer.ID, networkName)
	}

	// get hosts for billing platform service
	billingExternalHost, billingExternalHostErr := billingContainer.Host(ctx)
	if billingExternalHostErr != nil {
		return nil, errors.Wrap(billingExternalHostErr, "failed to get host for Billing Platform Service")
	}

	// get mapped port for billing platform service
	billingExternalPort, billingExternalPortErr := findMappedPort(ctx, 20*time.Second, billingContainer, DEFAULT_BILLING_PLATFORM_SERVICE_GRPC_PORT)
	if billingExternalPortErr != nil {
		return nil, errors.Wrap(billingExternalPortErr, "failed to get mapped port for Chip Ingress")
	}

	output := &Output{
		BillingPlatformService: &BillingPlatformServiceOutput{
			GRPCInternalURL: fmt.Sprintf("http://%s:%s", DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME, DEFAULT_BILLING_PLATFORM_SERVICE_GRPC_PORT),
			GRPCExternalURL: fmt.Sprintf("http://%s:%s", billingExternalHost, billingExternalPort.Port()),
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

	tempFile, tempErr := os.CreateTemp("", "billing-platform-service-docker-compose-*.yml")
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

func findMappedPort(ctx context.Context, timeout time.Duration, container *testcontainers.DockerContainer, port nat.Port) (nat.Port, error) {
	forCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tickerInterval := 5 * time.Second
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-forCtx.Done():
			return "", fmt.Errorf("timeout while waiting for mapped port for %s", port)
		case <-ticker.C:
			portCtx, portCancel := context.WithTimeout(ctx, tickerInterval)
			defer portCancel()
			mappedPort, mappedPortErr := container.MappedPort(portCtx, port)
			if mappedPortErr != nil {
				return "", errors.Wrapf(mappedPortErr, "failed to get mapped port for %s", port)
			}
			if mappedPort.Port() == "" {
				return "", fmt.Errorf("mapped port for %s is empty", port)
			}
			return mappedPort, nil
		}
	}
}

func connectNetwork(connCtx context.Context, timeout time.Duration, dockerClient *client.Client, containerID, networkName, stackIdentifier string) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	networkCtx, networkCancel := context.WithTimeout(connCtx, timeout)
	defer networkCancel()

	for {
		select {
		case <-networkCtx.Done():
			return fmt.Errorf("timeout while trying to connect billing-platform-service to default network")
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
