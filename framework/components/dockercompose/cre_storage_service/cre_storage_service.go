package cre_storage_service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go/modules/compose"

	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/utils"
)

// Roles the stack's services answer to on the Docker network, as defined by the compose file in cre-storage-service.
const (
	DeployAPIRole      = "cre-deploy-api"
	WorkflowSourceRole = "cre-workflow-source"
	ArtifactAPIRole    = "cre-artifact-api"
	ArtifactStoreRole  = "cre-artifact-store"

	UserServicePort   = "50051"
	NodeServicePort   = "50052"
	ArtifactStorePort = "4566"
)

const (
	// not a prefix of other stacks' names, e.g. cre-storage-service's standalone cre-storage-local, which Stop would remove
	DEFAULT_STACK_NAME = "ctf-cre-storage"

	userServiceName   = "storage-user-service"
	nodeServiceName   = "storage-node-service"
	artifactStoreName = "storage-s3"
)

type Input struct {
	// ComposeFile is local-stack/compose.yaml from cre-storage-service: a local path as file://<path>, or an http(s)://
	// URL readable without credentials.
	ComposeFile string `toml:"compose_file"`
	// Network the stack joins; it must exist. Defaults to framework.DefaultNetworkName.
	Network string `toml:"network"`
	// RegistryRPCURL, RegistryChainSelector and RegistryAddress locate the CapabilitiesRegistry the stack reads tenants'
	// nodes from; the RPC URL must be reachable from Network.
	RegistryRPCURL        string  `toml:"registry_rpc_url"`
	RegistryChainSelector uint64  `toml:"registry_chain_selector"`
	RegistryAddress       string  `toml:"registry_address"`
	DonFamily             string  `toml:"don_family"`
	Output                *Output `toml:"output"`
	UseCache              bool    `toml:"use_cache"`
}

type Output struct {
	// UserService serves cre-deploy-api (gRPC, HMAC auth)
	UserService *Endpoint `toml:"user_service"`
	// NodeService serves cre-workflow-source and cre-artifact-api (gRPC, CSA key auth)
	NodeService *Endpoint `toml:"node_service"`
	// ArtifactStore serves cre-artifact-store (S3); artifact URLs use the internal one
	ArtifactStore *Endpoint `toml:"artifact_store"`
}

// Endpoint is reachable at InternalURL from the stack's network and at ExternalURL from the host. gRPC endpoints are
// host:port, without TLS.
type Endpoint struct {
	InternalURL string `toml:"internal_url"`
	ExternalURL string `toml:"external_url"`
}

// New starts cre-storage-service's local stack from its compose file. Image overrides (CRE_STORAGE_IMAGE,
// CRE_STORAGE_MIGRATION_IMAGE, MOCK_ACCOUNT_REGISTRY_IMAGE) and LOG_LEVEL are taken from the host process env vars.
func New(in *Input) (*Output, error) {
	return NewWithContext(context.Background(), in)
}

func NewWithContext(ctx context.Context, in *Input) (*Output, error) {
	if in == nil {
		return nil, errors.New("input is nil")
	}
	if in.UseCache && in.Output != nil {
		return in.Output, nil
	}
	if err := validate(in); err != nil {
		return nil, err
	}

	composeFile, err := composeFilePath(in.ComposeFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get compose file path")
	}

	identifier := framework.DefaultTCName(DEFAULT_STACK_NAME)
	stack, err := compose.NewDockerComposeWith(
		compose.WithStackFiles(composeFile),
		compose.StackIdentifier(identifier),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compose stack for CRE Storage Service")
	}

	hostPorts, err := freeport.Take(3)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get free ports for CRE Storage Service")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// down removes the stack after a failed start, without ctx's deadline or cancellation, which may be what failed it
	down := func() {
		downCtx, downCancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
		defer downCancel()
		if downErr := stack.Down(downCtx, compose.RemoveOrphans(true), compose.RemoveVolumes(true)); downErr != nil {
			framework.L.Warn().Err(downErr).Msgf("failed to remove CRE Storage Service stack %s", identifier)
		}
	}

	framework.L.Info().Str("Compose file", composeFile).Msgf("Starting CRE Storage Service stack %s", identifier)
	if err := stack.WithEnv(env(in, os.Environ(), hostPorts)).Up(ctx, compose.Wait(true)); err != nil {
		down()
		return nil, errors.Wrap(err, "failed to start stack for CRE Storage Service")
	}

	external := func(service, port string) (string, error) {
		c, err := stack.ServiceContainer(ctx, service)
		if err != nil {
			return "", errors.Wrapf(err, "failed to get %s container", service)
		}
		host, err := utils.GetContainerHost(ctx, c)
		if err != nil {
			return "", errors.Wrapf(err, "failed to get host for %s", service)
		}
		mapped, err := utils.FindMappedPort(ctx, 20*time.Second, c, port+"/tcp")
		if err != nil {
			return "", errors.Wrapf(err, "failed to get mapped port for %s", service)
		}
		return fmt.Sprintf("%s:%s", host, mapped.Port()), nil
	}

	userExternal, err := external(userServiceName, UserServicePort)
	if err != nil {
		down()
		return nil, err
	}
	nodeExternal, err := external(nodeServiceName, NodeServicePort)
	if err != nil {
		down()
		return nil, err
	}
	storeExternal, err := external(artifactStoreName, ArtifactStorePort)
	if err != nil {
		down()
		return nil, err
	}

	in.Output = &Output{
		UserService:   &Endpoint{InternalURL: DeployAPIRole + ":" + UserServicePort, ExternalURL: userExternal},
		NodeService:   &Endpoint{InternalURL: WorkflowSourceRole + ":" + NodeServicePort, ExternalURL: nodeExternal},
		ArtifactStore: &Endpoint{InternalURL: "http://" + ArtifactStoreRole + ":" + ArtifactStorePort, ExternalURL: "http://" + storeExternal},
	}
	framework.L.Info().Msg("CRE Storage Service stack started")
	return in.Output, nil
}

// Stop removes every CRE Storage Service stack started by New, with its volumes.
func Stop() error {
	return framework.RemoveTestStack(DEFAULT_STACK_NAME)
}

// composeFilePath returns a local path to the compose file, downloading it if it's a URL.
func composeFilePath(uri string) (string, error) {
	if path, ok := strings.CutPrefix(uri, "file://"); ok {
		return path, nil
	}
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		return utils.ComposeFilePath(uri, DEFAULT_STACK_NAME)
	}
	return "", fmt.Errorf("compose_file must start with 'file://', 'http://' or 'https://', but '%s' was found", uri)
}

func validate(in *Input) error {
	var missing []string
	if in.ComposeFile == "" {
		missing = append(missing, "compose_file")
	}
	if in.RegistryRPCURL == "" {
		missing = append(missing, "registry_rpc_url")
	}
	if in.RegistryChainSelector == 0 {
		missing = append(missing, "registry_chain_selector")
	}
	if in.RegistryAddress == "" {
		missing = append(missing, "registry_address")
	}
	if in.DonFamily == "" {
		missing = append(missing, "don_family")
	}
	if len(missing) > 0 {
		return fmt.Errorf("CRE Storage Service input is missing %s", strings.Join(missing, ", "))
	}
	return nil
}

// env passes the host process env vars through, so they can override images and log level, and sets the inputs and
// host ports over them.
func env(in *Input, hostEnv []string, hostPorts []int) map[string]string {
	vars := make(map[string]string)
	for _, kv := range hostEnv {
		if k, v, ok := strings.Cut(kv, "="); ok {
			vars[k] = v
		}
	}

	network := in.Network
	if network == "" {
		network = framework.DefaultNetworkName
	}
	vars["CRE_NETWORK"] = network
	vars["CRE_NETWORK_EXTERNAL"] = "true"
	vars["CRE_LOCAL_REGISTRY_RPC_URL"] = in.RegistryRPCURL
	vars["CRE_LOCAL_REGISTRY_CHAIN_SELECTOR"] = strconv.FormatUint(in.RegistryChainSelector, 10)
	vars["CRE_LOCAL_REGISTRY_ADDRESS"] = in.RegistryAddress
	vars["CRE_LOCAL_DON_FAMILY"] = in.DonFamily
	vars["USER_SERVICE_HOST_PORT"] = strconv.Itoa(hostPorts[0])
	vars["NODE_SERVICE_HOST_PORT"] = strconv.Itoa(hostPorts[1])
	vars["S3_HOST_PORT"] = strconv.Itoa(hostPorts[2])
	return vars
}
