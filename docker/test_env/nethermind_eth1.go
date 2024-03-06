package test_env

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/templates"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

func NewNethermindEth1(networks []string, chainConfg *EthereumChainConfig, opts ...EnvComponentOption) (*Nethermind, error) {
	parts := strings.Split(defaultNethermindEth1Image, ":")
	g := &Nethermind{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "nethermind-eth1", uuid.NewString()[0:8]),
			Networks:         networks,
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
		},
		chainConfg: chainConfg,
		l:          logging.GetTestLogger(nil),
	}
	g.SetDefaultHooks()
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}

	g.EnvComponent.ContainerName = fmt.Sprintf("%s-%s-%s", "nethermind-eth1", strings.Replace(g.ContainerVersion, ".", "_", -1), uuid.NewString()[0:8])
	// if the internal docker repo is set then add it to the version
	g.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(g.EnvComponent.ContainerImage)
	return g, nil
}

func (g *Nethermind) getEth1ContainerRequest() (*tc.ContainerRequest, error) {
	keystoreDir, err := os.MkdirTemp("", "keystore")
	if err != nil {
		return nil, err
	}

	generatedData, err := generateKeystoreAndExtraData(keystoreDir)
	if err != nil {
		return nil, err
	}

	genesisJsonStr, err := templates.NethermindPoAGenesisJsonTemplate{
		ChainId:     fmt.Sprintf("%d", g.chainConfg.ChainID),
		AccountAddr: RootFundingAddr,
		ExtraData:   fmt.Sprintf("0x%s", hex.EncodeToString(generatedData.extraData)),
	}.String()
	if err != nil {
		return nil, err
	}
	genesisFile, err := os.CreateTemp("", "genesis_json")
	if err != nil {
		return nil, err
	}
	_, err = genesisFile.WriteString(genesisJsonStr)
	if err != nil {
		return nil, err
	}

	// create empty cfg file since if we don't pass any
	// default mainnet.cfg will be used
	noneCfg, err := os.CreateTemp("", "none.cfg")
	if err != nil {
		return nil, err
	}

	_, err = noneCfg.WriteString("{}")
	if err != nil {
		return nil, err
	}

	passFile, err := os.CreateTemp("", "password.txt")
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(passFile.Name(), []byte(""), 0600)
	if err != nil {
		return nil, err
	}

	rootFile, err := os.CreateTemp(keystoreDir, RootFundingAddr)
	if err != nil {
		return nil, err
	}
	_, err = rootFile.WriteString(RootFundingWallet)
	if err != nil {
		return nil, err
	}

	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		Image:           g.GetImageWithVersion(),
		Networks:        g.Networks,
		AlwaysPullImage: true,
		// ImagePlatform: "linux/x86_64",  //don't even try this on Apple Silicon, the node won't start due to .NET error
		ExposedPorts: []string{NatPortFormat(DEFAULT_EVM_NODE_HTTP_PORT), NatPortFormat(DEFAULT_EVM_NODE_WS_PORT)},
		WaitingFor: tcwait.ForAll(
			tcwait.ForLog("Nethermind initialization completed").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
		),
		Cmd: []string{
			"--config=/none.cfg",
			"--Init.ChainSpecPath=/chainspec.json",
			"--Init.DiscoveryEnabled=false",
			"--Init.WebSocketsEnabled=true",
			fmt.Sprintf("--JsonRpc.WebSocketsPort=%s", DEFAULT_EVM_NODE_WS_PORT),
			"--JsonRpc.Enabled=true",
			"--JsonRpc.EnabledModules=net,consensus,eth,subscribe,web3,admin,trace,txpool",
			"--JsonRpc.Host=0.0.0.0",
			fmt.Sprintf("--JsonRpc.Port=%s", DEFAULT_EVM_NODE_HTTP_PORT),
			"--KeyStore.KeyStoreDirectory=/keystore",
			fmt.Sprintf("--KeyStore.BlockAuthorAccount=%s", generatedData.minerAccount.Address.Hex()),
			fmt.Sprintf("--KeyStore.UnlockAccounts=%s", generatedData.minerAccount.Address.Hex()),
			"--KeyStore.PasswordFiles=/password.txt",
			"--Mining.Enabled=true",
			"--Init.PeerManagerEnabled=false",
			"--HealthChecks.Enabled=true", // default slug /health
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      genesisFile.Name(),
				ContainerFilePath: "/chainspec.json",
				FileMode:          0644,
			},
			{
				HostFilePath:      noneCfg.Name(),
				ContainerFilePath: "/none.cfg",
				FileMode:          0644,
			},
			{
				HostFilePath:      passFile.Name(),
				ContainerFilePath: "/password.txt",
				FileMode:          0644,
			},
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   keystoreDir,
				Target:   "/keystore",
				ReadOnly: false,
			})
		},
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: g.PostStartsHooks,
				PostStops:  g.PostStopsHooks,
			},
		},
	}, nil
}
