package billing_platform_service

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/utils"
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
	DSN string
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

	if in.UseCache && in.Output != nil {
		return in.Output, nil
	}

	in = defaultBillingPlatformService(in)
	identifier := framework.DefaultTCName(DEFAULT_STACK_NAME)
	framework.L.Debug().Str("Compose file", in.ComposeFile).
		Msgf("Starting Billing Platform Service stack with identifier %s",
			framework.DefaultTCName(DEFAULT_STACK_NAME))

	cFilePath, fileErr := utils.ComposeFilePath(in.ComposeFile, DEFAULT_STACK_NAME)
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
		connectCtx, connectCancel := context.WithTimeout(ctx, 30*time.Second)
		defer connectCancel()
		if connectErr := utils.ConnectNetwork(connectCtx, 30*time.Second, cli, billingContainer.ID, networkName, identifier); connectErr != nil {
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
	billingExternalHost, billingExternalHostErr := utils.GetContainerHost(ctx, billingContainer)
	if billingExternalHostErr != nil {
		return nil, errors.Wrap(billingExternalHostErr, "failed to get host for Billing Platform Service")
	}

	// get mapped port for billing platform service
	billingExternalPort, billingExternalPortErr := utils.FindMappedPort(ctx, 20*time.Second, billingContainer, DEFAULT_BILLING_PLATFORM_SERVICE_GRPC_PORT)
	if billingExternalPortErr != nil {
		return nil, errors.Wrap(billingExternalPortErr, "failed to get mapped port for Chip Ingress")
	}

	output := &Output{
		BillingPlatformService: &BillingPlatformServiceOutput{
			GRPCInternalURL: fmt.Sprintf("http://%s:%s", DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME, DEFAULT_BILLING_PLATFORM_SERVICE_GRPC_PORT),
			GRPCExternalURL: fmt.Sprintf("http://%s:%s", billingExternalHost, billingExternalPort.Port()),
		},
		Postgres: &PostgresOutput{

		}
	}

	framework.L.Info().Msg("Billing Platform Service stack start")

	return output, nil
}
