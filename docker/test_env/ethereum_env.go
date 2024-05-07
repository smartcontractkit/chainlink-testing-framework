package test_env

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	toml_utils "github.com/smartcontractkit/chainlink-testing-framework/utils/toml"
)

const (
	CONFIG_ENV_VAR_NAME = "PRIVATE_ETHEREUM_NETWORK_CONFIG_PATH"
)

var (
	ErrMissingConsensusLayer = errors.New("consensus layer is required for PoS")
	ErrTestConfigNotSaved    = errors.New("could not save test env config")
)

var MsgMismatchedExecutionClient = "you provided a custom docker image for %s execution client, but explicitly set a execution client to %s. Make them match or remove one or the other"

type EthereumNetworkBuilder struct {
	t                   *testing.T
	ls                  *logstream.LogStream
	dockerNetworks      []string
	ethereumVersion     config.EthereumVersion
	consensusLayer      *config.ConsensusLayer
	executionLayer      config.ExecutionLayer
	ethereumChainConfig *config.EthereumChainConfig
	existingConfig      *config.EthereumNetworkConfig
	customDockerImages  map[config.ContainerType]string
	addressesToFund     []string
	waitForFinalization bool
	existingFromEnvVar  bool
}

func NewEthereumNetworkBuilder() EthereumNetworkBuilder {
	return EthereumNetworkBuilder{
		dockerNetworks:      []string{},
		waitForFinalization: false,
	}
}

// WithConsensusType sets the consensus type for the network
// Deprecated: use WithEthereumVersion() instead
//
//nolint:staticcheck //ignore SA1019
func (b *EthereumNetworkBuilder) WithConsensusType(consensusType config.ConsensusType) *EthereumNetworkBuilder {
	switch consensusType {
	case config.ConsensusType_PoS:
		b.ethereumVersion = config.EthereumVersion_Eth2
	case config.ConsensusType_PoW:
		b.ethereumVersion = config.EthereumVersion_Eth1
	default:
		panic(fmt.Sprintf("unknown consensus type: %s", consensusType))
	}
	return b
}

func (b *EthereumNetworkBuilder) WithEthereumVersion(ethereumVersion config.EthereumVersion) *EthereumNetworkBuilder {
	b.ethereumVersion = ethereumVersion
	return b
}

func (b *EthereumNetworkBuilder) WithConsensusLayer(consensusLayer config.ConsensusLayer) *EthereumNetworkBuilder {
	b.consensusLayer = &consensusLayer
	return b
}

func (b *EthereumNetworkBuilder) WithExecutionLayer(executionLayer config.ExecutionLayer) *EthereumNetworkBuilder {
	b.executionLayer = executionLayer
	return b
}

func (b *EthereumNetworkBuilder) WithEthereumChainConfig(config config.EthereumChainConfig) *EthereumNetworkBuilder {
	b.ethereumChainConfig = &config
	return b
}

func (b *EthereumNetworkBuilder) WithDockerNetworks(networks []string) *EthereumNetworkBuilder {
	b.dockerNetworks = networks
	return b
}

func (b *EthereumNetworkBuilder) WithExistingConfig(config config.EthereumNetworkConfig) *EthereumNetworkBuilder {
	b.existingConfig = &config
	return b
}

func (b *EthereumNetworkBuilder) WithExistingConfigFromEnvVar() *EthereumNetworkBuilder {
	b.existingFromEnvVar = true
	return b
}

func (b *EthereumNetworkBuilder) WithTest(t *testing.T) *EthereumNetworkBuilder {
	b.t = t
	return b
}

func (b *EthereumNetworkBuilder) WithLogStream(ls *logstream.LogStream) *EthereumNetworkBuilder {
	b.ls = ls
	return b
}

func (b *EthereumNetworkBuilder) WithCustomDockerImages(newImages map[config.ContainerType]string) *EthereumNetworkBuilder {
	b.customDockerImages = newImages
	return b
}

func (b *EthereumNetworkBuilder) WithWaitingForFinalization() *EthereumNetworkBuilder {
	b.waitForFinalization = true
	return b
}

func (b *EthereumNetworkBuilder) buildNetworkConfig() EthereumNetwork {
	n := EthereumNetwork{
		EthereumNetworkConfig: config.EthereumNetworkConfig{
			EthereumVersion: &b.ethereumVersion,
			ExecutionLayer:  &b.executionLayer,
			ConsensusLayer:  b.consensusLayer,
		},
	}

	if b.existingConfig != nil && len(b.existingConfig.Containers) > 0 {
		n.isRecreated = true
		n.Containers = b.existingConfig.Containers
		n.GeneratedDataHostDir = b.existingConfig.GeneratedDataHostDir
		n.ValKeysDir = b.existingConfig.ValKeysDir
	}

	n.DockerNetworkNames = b.dockerNetworks
	n.WaitForFinalization = &b.waitForFinalization
	n.EthereumNetworkConfig.EthereumChainConfig = b.ethereumChainConfig
	n.EthereumNetworkConfig.CustomDockerImages = b.customDockerImages
	n.t = b.t
	n.ls = b.ls

	return n
}

func (b *EthereumNetworkBuilder) Build() (EthereumNetwork, error) {
	if b.existingFromEnvVar {
		path := os.Getenv(CONFIG_ENV_VAR_NAME)
		if path == "" {
			return EthereumNetwork{}, fmt.Errorf("environment variable %s is not set, but build from env var was requested", CONFIG_ENV_VAR_NAME)
		}

		config, err := NewPrivateChainEnvConfigFromFile(path)
		if err != nil {
			return EthereumNetwork{}, err
		}

		config.isRecreated = true

		return config, nil
	}

	if !b.importExistingConfig() {
		if b.ethereumChainConfig == nil {
			defaultConfig := config.GetDefaultChainConfig()
			b.ethereumChainConfig = &defaultConfig
		} else {
			b.ethereumChainConfig.FillInMissingValuesWithDefault()
		}

		b.ethereumChainConfig.GenerateGenesisTimestamp()
	}

	err := b.autoFill()
	if err != nil {
		return EthereumNetwork{}, err
	}

	err = b.validate()
	if err != nil {
		return EthereumNetwork{}, err
	}

	network := b.buildNetworkConfig()

	return network, network.Validate()
}

func (b *EthereumNetworkBuilder) importExistingConfig() bool {
	if b.existingConfig == nil {
		return false
	}

	if b.existingConfig.EthereumVersion != nil {
		b.ethereumVersion = *b.existingConfig.EthereumVersion
	}

	if b.existingConfig.ConsensusLayer != nil {
		b.consensusLayer = b.existingConfig.ConsensusLayer
	}

	if b.existingConfig.ExecutionLayer != nil {
		b.executionLayer = *b.existingConfig.ExecutionLayer
	}

	if len(b.existingConfig.DockerNetworkNames) > 0 {
		b.dockerNetworks = b.existingConfig.DockerNetworkNames
	}
	b.ethereumChainConfig = b.existingConfig.EthereumChainConfig
	b.customDockerImages = b.existingConfig.CustomDockerImages

	return true
}

func (b *EthereumNetworkBuilder) validate() error {
	if b.ethereumVersion == "" {
		return config.ErrMissingEthereumVersion
	}

	if b.executionLayer == "" {
		return config.ErrMissingExecutionLayer
	}

	//nolint:staticcheck //ignore SA1019
	if (b.ethereumVersion == config.EthereumVersion_Eth2 || b.ethereumVersion == config.EthereumVersion_Eth2_Legacy) && b.consensusLayer == nil {
		return ErrMissingConsensusLayer
	}

	err := b.validateCustomDockerImages()
	if err != nil {
		return err
	}

	for _, addr := range b.addressesToFund {
		if !common.IsHexAddress(addr) {
			return fmt.Errorf("address %s is not a valid hex address", addr)
		}
	}

	if b.ethereumChainConfig == nil {
		return errors.New("ethereum chain config is required")
	}

	return b.ethereumChainConfig.Validate(logging.GetTestLogger(nil), &b.ethereumVersion)
}

func (b *EthereumNetworkBuilder) validateCustomDockerImages() error {
	if len(b.customDockerImages) > 0 {
		if image, ok := b.customDockerImages[config.ContainerType_ExecutionLayer]; ok {

			isSupported, reason, err := IsDockerImageVersionSupported(image)
			if err != nil {
				return err
			}

			if !isSupported {
				return fmt.Errorf("docker image %s is not supported, due to: %s", image, reason)
			}

			executionLayer, err := GetExecutionLayerFromDockerImage(image)
			if err != nil {
				return err
			}

			if executionLayer != b.executionLayer {
				return fmt.Errorf(MsgMismatchedExecutionClient, executionLayer, b.executionLayer)
			}
		}
	}

	return nil
}

func (b *EthereumNetworkBuilder) autoFill() error {
	err := b.setExecutionLayerBasedOnCustomDocker()
	if err != nil {
		return err
	}

	err = b.fetchLatestReleaseVersionIfNeed()
	if err != nil {
		return err
	}

	if b.ethereumVersion == "" {
		if err := b.trySettingEthereumVersionBasedOnCustomImage(); err != nil {
			return err
		}
	}

	//nolint:staticcheck //ignore SA1019
	if (b.ethereumVersion == config.EthereumVersion_Eth2_Legacy || b.ethereumVersion == config.EthereumVersion_Eth2) && b.consensusLayer == nil {
		b.consensusLayer = &config.ConsensusLayer_Prysm
	}

	//nolint:staticcheck //ignore SA1019
	if b.ethereumVersion == config.EthereumVersion_Eth1_Legacy {
		b.ethereumVersion = config.EthereumVersion_Eth1
	}

	//nolint:staticcheck //ignore SA1019
	if b.ethereumVersion == config.EthereumVersion_Eth2_Legacy {
		b.ethereumVersion = config.EthereumVersion_Eth2
	}

	return nil
}

func (b *EthereumNetworkBuilder) setExecutionLayerBasedOnCustomDocker() error {
	if b.executionLayer == "" && len(b.customDockerImages) > 0 {
		if image, ok := b.customDockerImages[config.ContainerType_ExecutionLayer]; ok {
			var err error
			b.executionLayer, err = GetExecutionLayerFromDockerImage(image)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *EthereumNetworkBuilder) fetchLatestReleaseVersionIfNeed() error {
	if image, ok := b.customDockerImages[config.ContainerType_ExecutionLayer]; ok {
		var err error
		b.customDockerImages[config.ContainerType_ExecutionLayer], err = FetchLatestEthereumClientDockerImageVersionIfNeed(image)
		if err != nil {
			return err
		}

	}

	return nil
}

func (b *EthereumNetworkBuilder) trySettingEthereumVersionBasedOnCustomImage() error {
	var dockerImageToUse string

	count := 0

	// if we are using custom docker image for execution client, extract it
	for t, customImage := range b.customDockerImages {
		if t == config.ContainerType_ExecutionLayer {
			dockerImageToUse = customImage
			count++
		}
	}

	if count > 1 {
		return errors.New("multiple custom docker images for execution layer provided, but only one is allowed")
	}

	if dockerImageToUse == "" {
		return errors.New("couldn't determine ethereum version as no custom docker image for execution layer was provided")
	}

	ethereumVersion, err := GetEthereumVersionFromImage(b.executionLayer, dockerImageToUse)
	if err != nil {
		return err
	}

	b.ethereumVersion = ethereumVersion

	return nil
}

type EthereumNetwork struct {
	config.EthereumNetworkConfig
	isRecreated bool
	t           *testing.T
	ls          *logstream.LogStream
}

func (en *EthereumNetwork) Start() (blockchain.EVMNetwork, RpcProvider, error) {
	switch *en.EthereumVersion {
	//nolint:staticcheck //ignore SA1019
	case config.EthereumVersion_Eth1, config.EthereumVersion_Eth1_Legacy:
		return en.startEth1()
	//nolint:staticcheck //ignore SA1019
	case config.EthereumVersion_Eth2_Legacy, config.EthereumVersion_Eth2:
		return en.startEth2()
	default:
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unknown ethereum version: %s", *en.EthereumVersion)
	}
}

func (en *EthereumNetwork) startEth2() (blockchain.EVMNetwork, RpcProvider, error) {
	rpcProvider := NewRPCProvider([]string{}, []string{}, []string{}, []string{})

	var net blockchain.EVMNetwork

	if *en.ConsensusLayer != config.ConsensusLayer_Prysm {
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf("unsupported consensus layer: %s. Use 'prysm'", *en.ConsensusLayer)
	}

	dockerNetworks, err := en.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to create docker networks")
	}
	generatedDataHostDir, valKeysDir, err := en.generateGenesisAndFoldersIfNeeded()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to create host directories")
	}

	opts := en.getExecutionLayerEnvComponentOpts()

	var client ExecutionClient
	var clientErr error
	switch *en.ExecutionLayer {
	case config.ExecutionLayer_Geth:
		client, clientErr = NewGethEth2(dockerNetworks, en.EthereumChainConfig, generatedDataHostDir, config.ConsensusLayer_Prysm, opts...)
	case config.ExecutionLayer_Nethermind:
		client, clientErr = NewNethermindEth2(dockerNetworks, generatedDataHostDir, config.ConsensusLayer_Prysm, opts...)
	case config.ExecutionLayer_Erigon:
		client, clientErr = NewErigonEth2(dockerNetworks, en.EthereumChainConfig, generatedDataHostDir, config.ConsensusLayer_Prysm, opts...)
	case config.ExecutionLayer_Besu:
		client, clientErr = NewBesuEth2(dockerNetworks, en.EthereumChainConfig, generatedDataHostDir, config.ConsensusLayer_Prysm, opts...)
	default:
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf(MsgUnsupportedExecutionLayer, *en.ExecutionLayer)
	}

	if clientErr != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(clientErr, "failed to create  %s execution client instance", *en.ExecutionLayer)
	}

	client.WithTestInstance(en.t)

	net, err = client.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to start %s execution client", *en.ExecutionLayer)
	}

	beacon, err := NewPrysmBeaconChain(dockerNetworks, en.EthereumChainConfig, generatedDataHostDir, client.GetInternalExecutionURL(), append(en.getImageOverride(config.ContainerType_ValKeysGenerator), en.setExistingContainerName(config.ContainerType_ConsensusLayer))...)
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to create beacon chain instance")
	}

	beacon.WithTestInstance(en.t)
	err = beacon.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to start beacon chain")
	}

	validator, err := NewPrysmValidator(dockerNetworks, en.EthereumChainConfig, generatedDataHostDir, valKeysDir, beacon.
		InternalBeaconRpcProvider, append(en.getImageOverride(config.ContainerType_ValKeysGenerator), en.setExistingContainerName(config.ContainerType_ConsensusValidator))...)
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to create validator instance")
	}

	validator.WithTestInstance(en.t)
	err = validator.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to start validator")
	}

	err = client.WaitUntilChainIsReady(testcontext.Get(en.t), en.EthereumChainConfig.GetDefaultWaitDuration())
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to wait for chain to be ready")
	}

	en.DockerNetworkNames = dockerNetworks
	net = en.getFinalEvmNetworkConfig(net)

	logger := logging.GetTestLogger(en.t)
	if en.WaitForFinalization != nil && *en.WaitForFinalization {
		evmClient, err := blockchain.NewEVMClientFromNetwork(net, logger)
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to create evm client")
		}

		err = waitForChainToFinaliseAnEpoch(logger, evmClient, en.EthereumChainConfig.GetDefaultFinalizationWaitDuration())
		if err != nil {
			return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to wait for chain to finalize first epoch")
		}
	} else {
		logger.Info().Msg("Not waiting for chain to finalize first epoch")
	}

	containers := config.EthereumNetworkContainers{
		{
			ContainerName: client.GetContainerName(),
			ContainerType: config.ContainerType_ExecutionLayer,
			Container:     client.GetContainer(),
		},
		{
			ContainerName: beacon.ContainerName,
			ContainerType: config.ContainerType_ConsensusLayer,
			Container:     &beacon.Container,
		},
		{
			ContainerName: validator.ContainerName,
			ContainerType: config.ContainerType_ConsensusValidator,
			Container:     &validator.Container,
		},
	}

	en.Containers = append(en.Containers, containers...)

	rpcProvider.privateHttpUrls = append(rpcProvider.privateHttpUrls, client.GetInternalHttpUrl())
	rpcProvider.privateWsUrls = append(rpcProvider.privateWsUrls, client.GetInternalWsUrl())
	rpcProvider.publiclHttpUrls = append(rpcProvider.publiclHttpUrls, client.GetExternalHttpUrl())
	rpcProvider.publicWsUrls = append(rpcProvider.publicWsUrls, client.GetExternalWsUrl())

	return net, rpcProvider, nil
}

func (en *EthereumNetwork) startEth1() (blockchain.EVMNetwork, RpcProvider, error) {
	var net blockchain.EVMNetwork
	rpcProvider := RpcProvider{
		privateHttpUrls: []string{},
		privateWsUrls:   []string{},
		publiclHttpUrls: []string{},
		publicWsUrls:    []string{},
	}

	dockerNetworks, err := en.getOrCreateDockerNetworks()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to create docker networks")
	}

	opts := en.getExecutionLayerEnvComponentOpts()

	var client ExecutionClient
	var clientErr error
	switch *en.ExecutionLayer {
	case config.ExecutionLayer_Geth:
		client = NewGethEth1(dockerNetworks, en.EthereumChainConfig, opts...)
	case config.ExecutionLayer_Besu:
		client, clientErr = NewBesuEth1(dockerNetworks, en.EthereumChainConfig, opts...)
	case config.ExecutionLayer_Erigon:
		client, clientErr = NewErigonEth1(dockerNetworks, en.EthereumChainConfig, opts...)
	case config.ExecutionLayer_Nethermind:
		client, clientErr = NewNethermindEth1(dockerNetworks, en.EthereumChainConfig, opts...)
	default:
		return blockchain.EVMNetwork{}, RpcProvider{}, fmt.Errorf(MsgUnsupportedExecutionLayer, *en.ExecutionLayer)
	}

	if clientErr != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(clientErr, "failed to create  %s execution client instance", *en.ExecutionLayer)
	}

	client.WithTestInstance(en.t)

	net, err = client.StartContainer()
	if err != nil {
		return blockchain.EVMNetwork{}, RpcProvider{}, errors.Wrapf(err, "failed to start %s execution client", *en.ExecutionLayer)
	}

	containers := config.EthereumNetworkContainers{
		{
			ContainerName: client.GetContainerName(),
			ContainerType: config.ContainerType_ExecutionLayer,
			Container:     client.GetContainer(),
		},
	}

	en.Containers = append(en.Containers, containers...)
	rpcProvider.privateHttpUrls = append(rpcProvider.privateHttpUrls, client.GetInternalHttpUrl())
	rpcProvider.privateWsUrls = append(rpcProvider.privateWsUrls, client.GetInternalWsUrl())
	rpcProvider.publiclHttpUrls = append(rpcProvider.publiclHttpUrls, client.GetExternalHttpUrl())
	rpcProvider.publicWsUrls = append(rpcProvider.publicWsUrls, client.GetExternalWsUrl())

	en.DockerNetworkNames = dockerNetworks
	net.ChainID = int64(en.EthereumChainConfig.ChainID)

	return net, rpcProvider, nil
}

func (en *EthereumNetwork) getOrCreateDockerNetworks() ([]string, error) {
	if len(en.DockerNetworkNames) != 0 {
		return en.DockerNetworkNames, nil
	}

	network, err := docker.CreateNetwork(logging.GetTestLogger(en.t))
	if err != nil {
		return []string{}, err
	}

	return []string{network.Name}, nil
}

func (en *EthereumNetwork) generateGenesisAndFoldersIfNeeded() (generatedDataHostDir string, valKeysDir string, err error) {
	// create host directories and run genesis containers only if we are NOT recreating existing containers
	if !en.isRecreated {
		generatedDataHostDir, valKeysDir, err = createHostDirectories()

		en.GeneratedDataHostDir = &generatedDataHostDir
		en.ValKeysDir = &valKeysDir

		if err != nil {
			return
		}

		var valKeysGenerator *ValKeysGenerator
		valKeysGenerator, err = NewValKeysGeneretor(en.EthereumChainConfig, valKeysDir, en.getImageOverride(config.ContainerType_ValKeysGenerator)...)
		if err != nil {
			err = errors.Wrap(err, "failed to start val keys generator")
			return
		}
		valKeysGenerator.WithTestInstance(en.t)

		err = valKeysGenerator.StartContainer()
		if err != nil {
			err = errors.Wrap(err, "failed to start val keys generator")
			return
		}

		var genesis *EthGenesisGeneretor
		genesis, err = NewEthGenesisGenerator(*en.EthereumChainConfig, generatedDataHostDir, en.getImageOverride(config.ContainerType_GenesisGenerator)...)
		if err != nil {
			err = errors.Wrap(err, "failed to start genesis generator")
			return
		}

		genesis.WithTestInstance(en.t)

		err = genesis.StartContainer()
		if err != nil {
			return
		}

		initHelper := NewInitHelper(*en.EthereumChainConfig, generatedDataHostDir).WithTestInstance(en.t)
		err = initHelper.StartContainer()
		if err != nil {
			err = errors.Wrap(err, "failed to start init helper")
			return
		}
	} else {
		// we don't set actual values to not increase complexity, as they do not matter for containers that are already running
		if en.GeneratedDataHostDir == nil {
			generatedDataHostDir = ""
		}

		if en.ValKeysDir == nil {
			valKeysDir = ""
		}

		generatedDataHostDir = *en.GeneratedDataHostDir
		valKeysDir = *en.ValKeysDir
	}

	return
}

func (en *EthereumNetwork) getFinalEvmNetworkConfig(net blockchain.EVMNetwork) blockchain.EVMNetwork {
	net.ChainID = int64(en.EthereumChainConfig.ChainID)
	// use a higher value than the default, because eth2 is slower than dev-mode eth1
	net.Timeout = blockchain.StrDuration{Duration: time.Duration(4 * time.Minute)}
	net.FinalityTag = true
	net.FinalityDepth = 0

	if *en.ExecutionLayer == config.ExecutionLayer_Besu {
		// Besu doesn't support "eth_maxPriorityFeePerGas" https://github.com/hyperledger/besu/issues/5658
		// And if gas is too low, then transaction doesn't get to prioritized pool and is not a candidate for inclusion in the next block
		net.GasEstimationBuffer = 10_000_000_000
	} else {
		net.SupportsEIP1559 = true
	}

	return net
}

func (en *EthereumNetwork) setExistingContainerName(ct config.ContainerType) EnvComponentOption {
	if !en.isRecreated {
		return func(c *EnvComponent) {}
	}

	for _, container := range en.Containers {
		if container.ContainerType == ct {
			return func(c *EnvComponent) {
				if container.ContainerName != "" {
					c.ContainerName = container.ContainerName
					c.WasRecreated = true
				}
			}
		}
	}

	return func(c *EnvComponent) {}
}

func (en *EthereumNetwork) getImageOverride(ct config.ContainerType) []EnvComponentOption {
	options := []EnvComponentOption{}
	if image, ok := en.CustomDockerImages[ct]; ok {
		options = append(options, WithContainerImageWithVersion(image))
	}
	return options
}

func (en *EthereumNetwork) Save() error {
	name := fmt.Sprintf("ethereum_network_%s", uuid.NewString()[0:8])
	confPath, err := toml_utils.SaveStructAsToml(en, ".private_chains", name)
	if err != nil {
		return ErrTestConfigNotSaved
	}

	log := logging.GetTestLogger(en.t)
	log.Info().Msgf("Saved private Ethereum Network config. To reuse in e2e tests, set: %s=%s", CONFIG_ENV_VAR_NAME, confPath)

	return nil
}

func (en *EthereumNetwork) getExecutionLayerEnvComponentOpts() []EnvComponentOption {
	opts := []EnvComponentOption{}
	opts = append(opts, en.getImageOverride(config.ContainerType_ExecutionLayer)...)
	opts = append(opts, en.setExistingContainerName(config.ContainerType_ExecutionLayer))
	opts = append(opts, WithLogStream(en.ls))

	if en.NodeLogLevel != nil && *en.NodeLogLevel != "" {
		opts = append(opts, WithLogLevel(strings.ToLower(*en.NodeLogLevel)))
	}

	return opts
}

// RpcProvider holds all necessary URLs to connect to a simulated chain or a real RPC provider connected to a live chain
// maybe only store ports here and depending on where the test is executed return different URLs?
// maybe 3 different constructors for each "perspective"? (docker, k8s with local runner, k8s with remote runner)
type RpcProvider struct {
	privateHttpUrls []string
	privateWsUrls   []string
	publiclHttpUrls []string
	publicWsUrls    []string
}

// NewRPCProvider creates a new RpcProvider, and should only be used for custom network connections e.g. to a live testnet chain
func NewRPCProvider(
	privateHttpUrls,
	privateWsUrls,
	publiclHttpUrls,
	publicWsUrls []string,
) RpcProvider {
	return RpcProvider{
		privateHttpUrls: privateHttpUrls,
		privateWsUrls:   privateWsUrls,
		publiclHttpUrls: publiclHttpUrls,
		publicWsUrls:    publicWsUrls,
	}
}

func (s *RpcProvider) PrivateHttpUrls() []string {
	return s.privateHttpUrls
}

func (s *RpcProvider) PrivateWsUrsl() []string {
	return s.privateWsUrls
}

func (s *RpcProvider) PublicHttpUrls() []string {
	return s.publiclHttpUrls
}

func (s *RpcProvider) PublicWsUrls() []string {
	return s.publicWsUrls
}

func createHostDirectories() (string, string, error) {
	customConfigDataDir, err := os.MkdirTemp("", "custom_config_data")
	if err != nil {
		return "", "", err
	}

	valKeysDir, err := os.MkdirTemp("", "val_keys")
	if err != nil {
		return "", "", err
	}

	return customConfigDataDir, valKeysDir, nil
}

func NewPrivateChainEnvConfigFromFile(path string) (EthereumNetwork, error) {
	c := EthereumNetwork{}
	err := toml_utils.OpenTomlFileAsStruct(path, &c)
	return c, err
}
