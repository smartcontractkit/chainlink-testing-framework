package billing_platform_service

import (
	"context"
	"fmt"
	"os"
	"strconv"
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
	"github.com/smartcontractkit/freeport"
)

const DefaultPostgresDSN = "postgres://postgres:postgres@postgres:5432/billing_platform?sslmode=disable"

type Output struct {
	BillingPlatformService *BillingPlatformServiceOutput
	Postgres               *PostgresOutput
}

type BillingPlatformServiceOutput struct {
	BillingGRPCInternalURL string
	BillingGRPCExternalURL string
	CreditGRPCInternalURL  string
	CreditGRPCExternalURL  string
}

type PostgresOutput struct {
	DSN string
}

type Input struct {
	ComposeFile                 string   `toml:"compose_file"`
	ExtraDockerNetworks         []string `toml:"extra_docker_networks"`
	Output                      *Output  `toml:"output"`
	UseCache                    bool     `toml:"use_cache"`
	ChainSelector               uint64   `toml:"chain_selector"`
	StreamsAPIURL               string   `toml:"streams_api_url"`
	StreamsAPIKey               string   `toml:"streams_api_key"`
	StreamsAPISecret            string   `toml:"streams_api_secret"`
	RPCURL                      string   `toml:"rpc_url"`
	WorkflowRegistryAddress     string   `toml:"workflow_registry_address"`
	CapabilitiesRegistryAddress string   `toml:"capabilities_registry_address"`
	WorkflowOwners              []string `toml:"workflow_owners"`
}

func defaultBillingPlatformService(in *Input) *Input {
	if in.ComposeFile == "" {
		in.ComposeFile = "./docker-compose.yml"
	}
	return in
}

const (
	DEFAULT_BILLING_PLATFORM_SERVICE_BILLING_GRPC_PORT = "2222"
	DEFAULT_BILLING_PLATFORM_SERVICE_CREDIT_GRPC_PORT  = "2223"
	DEFAULT_POSTGRES_PORT                              = "5432"
	DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME      = "billing-platform-service"
	DEFAULT_POSTGRES_SERVICE_NAME                      = "postgres"
)

// New starts a Billing Platform Service stack using docker-compose. Various env vars are set to sensible defaults and
// input values, but can be overridden by the host process env vars if needed.
//
// Import env vars that can be set to override defaults:
//   - TEST_OWNERS = comma separated list of workflow owners
//   - STREAMS_API_URL = URL for the Streams API; can use a mock server if needed
//   - STREAMS_API_KEY = API key if using a staging or prod Streams API
//   - STREAMS_API_SECRET = API secret if using a staging or prod Streams API
func New(in *Input) (*Output, error) {
	if in == nil {
		return nil, errors.New("input is nil")
	}

	if in.UseCache && in.Output != nil {
		return in.Output, nil
	}

	in = defaultBillingPlatformService(in)
	identifier := framework.DefaultTCName(DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME)
	framework.L.Debug().Str("Compose file", in.ComposeFile).
		Msgf("Starting Billing Platform Service stack with identifier %s",
			framework.DefaultTCName(DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME))

	cFilePath, fileErr := utils.ComposeFilePath(in.ComposeFile, DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME)
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

	envVars["MAINNET_WORKFLOW_REGISTRY_CHAIN_SELECTOR"] = strconv.FormatUint(in.ChainSelector, 10)
	envVars["MAINNET_WORKFLOW_REGISTRY_CONTRACT_ADDRESS"] = in.WorkflowRegistryAddress
	envVars["MAINNET_WORKFLOW_REGISTRY_RPC_URL"] = in.RPCURL
	envVars["MAINNET_WORKFLOW_REGISTRY_FINALITY_DEPTH"] = "0"                                     // Instant finality on devnet
	envVars["KMS_PROOF_SIGNING_KEY_ID"] = "00000000-0000-0000-0000-000000000001"                  // provisioned via LocalStack
	envVars["VERIFIER_INITIAL_INTERVAL"] = "0s"                                                   // reduced to force verifier to start immediately in integration tests
	envVars["VERIFIER_MAXIMUM_INTERVAL"] = "1s"                                                   // reduced to force verifier to start immediately in integration tests
	envVars["LINKING_REQUEST_COOLDOWN"] = "0s"                                                    // reduced to force consequtive linking requests to be processed immediately in integration tests
	envVars["ETH_FEED_ID"] = "0x000359843a543ee2fe414dc14c7e7920ef10f4372990b79d6361cdc0dd1ba782" // set as default eth feed ID

	envVars["MAINNET_CAPABILITIES_REGISTRY_CHAIN_SELECTOR"] = strconv.FormatUint(in.ChainSelector, 10)
	envVars["MAINNET_CAPABILITIES_REGISTRY_CONTRACT_ADDRESS"] = in.CapabilitiesRegistryAddress
	envVars["MAINNET_CAPABILITIES_REGISTRY_RPC_URL"] = in.RPCURL
	envVars["MAINNET_CAPABILITIES_REGISTRY_FINALITY_DEPTH"] = "10" // Arbitrary value, adjust as needed

	envVars["TEST_OWNERS"] = strings.Join(in.WorkflowOwners, ",")
	envVars["STREAMS_API_URL"] = in.StreamsAPIURL
	envVars["STREAMS_API_KEY"] = in.StreamsAPIKey
	envVars["STREAMS_API_SECRET"] = in.StreamsAPISecret

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envVars[pair[0]] = pair[1]
		}
	}

	// set these env vars after reading env vars from host
	port, err := freeport.Take(1)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get free port for Billing Platform Service postgres")
	}

	envVars["POSTGRES_PORT"] = strconv.FormatInt(int64(port[0]), 10)
	envVars["DEFAULT_DSN"] = DefaultPostgresDSN

	upErr := stack.
		WithEnv(envVars).
		Up(ctx)

	if upErr != nil {
		return nil, errors.Wrap(upErr, "failed to start stack for Billing Platform Service")
	}

	stack.WaitForService(DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME,
		wait.ForAll(
			wait.ForLog("GRPC server is live").WithPollInterval(200*time.Millisecond),
			wait.ForListeningPort(nat.Port(DEFAULT_BILLING_PLATFORM_SERVICE_BILLING_GRPC_PORT)),
			wait.ForListeningPort(nat.Port(DEFAULT_BILLING_PLATFORM_SERVICE_CREDIT_GRPC_PORT)),
		).WithDeadline(1*time.Minute),
	)

	billingContainer, billingErr := stack.ServiceContainer(ctx, DEFAULT_BILLING_PLATFORM_SERVICE_SERVICE_NAME)
	if billingErr != nil {
		return nil, errors.Wrap(billingErr, "failed to get billing-platform-service container")
	}

	postgresContainer, postgresErr := stack.ServiceContainer(ctx, DEFAULT_POSTGRES_SERVICE_NAME)
	if postgresErr != nil {
		return nil, errors.Wrap(postgresErr, "failed to get postgres container")
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

	// get hosts for billing platform service
	postgresExternalHost, postgresExternalHostErr := utils.GetContainerHost(ctx, postgresContainer)
	if postgresExternalHostErr != nil {
		return nil, errors.Wrap(postgresExternalHostErr, "failed to get host for postgres")
	}

	// get mapped ports for billing platform service
	serviceOutput, err := getExternalPorts(ctx, billingExternalHost, billingContainer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get mapped port for Billing Platform Service")
	}

	externalPostgresPort, err := utils.FindMappedPort(ctx, 20*time.Second, postgresContainer, nat.Port(DEFAULT_POSTGRES_PORT+"/tcp"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get mapped port for postgres")
	}

	output := &Output{
		BillingPlatformService: serviceOutput,
		Postgres: &PostgresOutput{
			DSN: fmt.Sprintf("postgres://postgres:postgres@%s:%s/billing_platform", postgresExternalHost, externalPostgresPort.Port()),
		},
	}

	framework.L.Info().Msg("Billing Platform Service stack start")

	return output, nil
}

func getExternalPorts(ctx context.Context, billingExternalHost string, billingContainer *testcontainers.DockerContainer) (*BillingPlatformServiceOutput, error) {
	ports := map[string]nat.Port{
		"billing": DEFAULT_BILLING_PLATFORM_SERVICE_BILLING_GRPC_PORT,
		"credit":  DEFAULT_BILLING_PLATFORM_SERVICE_CREDIT_GRPC_PORT,
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
		}
	}

	return &output, nil
}
