package billing_platform_service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/utils"
)

const DefaultPostgresDSN = "postgres://postgres:postgres@postgres:5432/billing_platform"

type Output struct {
	BillingPlatformService *BillingPlatformServiceOutput
	Postgres               *PostgresOutput
}

type BillingPlatformServiceOutput struct {
	BillingGRPCInternalURL   string
	BillingGRPCExternalURL   string
	CreditGRPCInternalURL    string
	CreditGRPCExternalURL    string
	OwnershipGRPCInternalURL string
	OwnershipGRPCExternalURL string
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

	DEFAULT_BILLING_PLATFORM_SERVICE_BILLING_GRPC_PORT   = "2222"
	DEFAULT_BILLING_PLATFORM_SERVICE_CREDIT_GRPC_PORT    = "2223"
	DEFAULT_BILLING_PLATFORM_SERVICE_OWNERSHIP_GRPC_PORT = "2257"
	DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME        = "billing-platform-service"
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

	// Start the stackwith all environment variables from the host process
	// set development defaults for necessary environment variables and allow them to be overridden by the host process
	envVars := make(map[string]string)

	envVars["MAINNET_WORKFLOW_REGISTRY_CHAIN_SELECTOR"] = "7759470850252068959"                          // Anvil Devnet
	envVars["MAINNET_WORKFLOW_REGISTRY_CONTRACT_ADDRESS"] = "0xA15BB66138824a1c7167f5E85b957d04Dd34E468" // Deployed via Linking integration tests
	envVars["MAINNET_WORKFLOW_REGISTRY_RPC_URL"] = "http://anvil:8545"                                   // Anvil inside Docker
	envVars["MAINNET_WORKFLOW_REGISTRY_FINALITY_DEPTH"] = "0"                                            // Instant finality on devnet
	envVars["TESTNET_WORKFLOW_REGISTRY_CHAIN_SELECTOR"] = "10344971235874465080"                         // Base Sepolia
	envVars["TESTNET_WORKFLOW_REGISTRY_CONTRACT_ADDRESS"] = "0xED1D0d87706a466151d67A6a06d69534C97BE66F" // Used for Billing integration tests
	envVars["TESTNET_WORKFLOW_REGISTRY_RPC_URL"] = "http://anvil:8545"                                   // Anvil inside Docker
	envVars["TESTNET_WORKFLOW_REGISTRY_FINALITY_DEPTH"] = "10"                                           // Arbitrary value, adjust as needed
	envVars["KMS_PROOF_SIGNING_KEY_ID"] = "00000000-0000-0000-0000-000000000001"                         // provisioned via LocalStack
	envVars["VERIFIER_INITIAL_INTERVAL"] = "0s"                                                          // reduced to force verifier to start immediately in integration tests
	envVars["VERIFIER_MAXIMUM_INTERVAL"] = "1s"                                                          // reduced to force verifier to start immediately in integration tests
	envVars["LINKING_REQUEST_COOLDOWN"] = "0s"                                                           // reduced to force consequtive linking requests to be processed immediately in integration tests

	envVars["MAINNET_CAPABILITIES_REGISTRY_CHAIN_SELECTOR"] = "10344971235874465080"                         // Base Sepolia
	envVars["MAINNET_CAPABILITIES_REGISTRY_CONTRACT_ADDRESS"] = "0x4c0a7d8f1b2e3c5f6a9b8e2d3c4f5e6b7a8b9c0d" // dummy address
	envVars["MAINNET_CAPABILITIES_REGISTRY_RPC_URL"] = "http://anvil:8545"                                   // Anvil RPC URL
	envVars["MAINNET_CAPABILITIES_REGISTRY_FINALITY_DEPTH"] = "10"                                           // Arbitrary value, adjust as needed
	envVars["TESTNET_CAPABILITIES_REGISTRY_CHAIN_SELECTOR"] = "10344971235874465080"                         // Base Sepolia
	envVars["TESTNET_CAPABILITIES_REGISTRY_CONTRACT_ADDRESS"] = "0x4c0a7d8f1b2e3c5f6a9b8e2d3c4f5e6b7a8b9c0d" // dummy address
	envVars["TESTNET_CAPABILITIES_REGISTRY_RPC_URL"] = "http://anvil:8545"                                   // Anvil RPC URL
	envVars["TESTNET_CAPABILITIES_REGISTRY_FINALITY_DEPTH"] = "10"                                           // Arbitrary value, adjust as needed

	envVars["STREAMS_API_URL"] = ""
	envVars["STREAMS_API_KEY"] = ""
	envVars["STREAMS_API_SECRET"] = ""

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
		return nil, errors.Wrap(upErr, "failed to start stack for Billing Platform Service")
	}

	stack.WaitForService(DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME,
		wait.ForAll(
			wait.ForLog("GRPC server is live").WithPollInterval(200*time.Millisecond),
			wait.ForListeningPort(DEFAULT_BILLING_PLATFORM_SERVICE_BILLING_GRPC_PORT),
			wait.ForListeningPort(DEFAULT_BILLING_PLATFORM_SERVICE_CREDIT_GRPC_PORT),
			wait.ForListeningPort(DEFAULT_BILLING_PLATFORM_SERVICE_OWNERSHIP_GRPC_PORT),
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

	// get mapped ports for billing platform service
	serviceOutput, err := getExternalPorts(ctx, billingExternalHost, billingContainer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get mapped port for Billing Platform Service")
	}

	output := &Output{
		BillingPlatformService: serviceOutput,
		Postgres: &PostgresOutput{
			DSN: DefaultPostgresDSN,
		},
	}

	framework.L.Info().Msg("Billing Platform Service stack start")

	return output, nil
}

func getExternalPorts(ctx context.Context, billingExternalHost string, billingContainer *testcontainers.DockerContainer) (*BillingPlatformServiceOutput, error) {
	ports := map[string]nat.Port{
		"billing":   DEFAULT_BILLING_PLATFORM_SERVICE_BILLING_GRPC_PORT,
		"credit":    DEFAULT_BILLING_PLATFORM_SERVICE_CREDIT_GRPC_PORT,
		"ownership": DEFAULT_BILLING_PLATFORM_SERVICE_OWNERSHIP_GRPC_PORT,
	}

	output := BillingPlatformServiceOutput{}

	for name, defaultPort := range ports {
		externalPort, err := utils.FindMappedPort(ctx, 20*time.Second, billingContainer, defaultPort)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get mapped port for Billing Platform Service")
		}

		internal := fmt.Sprintf("http://%s:%s", DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME, defaultPort)
		external := fmt.Sprintf("http://%s:%s", billingExternalHost, externalPort.Port())

		switch name {
		case "billing":
			output.BillingGRPCInternalURL = internal
			output.BillingGRPCExternalURL = external
		case "credit":
			output.CreditGRPCInternalURL = internal
			output.CreditGRPCExternalURL = external
		case "ownership":
			output.OwnershipGRPCInternalURL = internal
			output.OwnershipGRPCExternalURL = external
		}
	}

	return &output, nil
}
