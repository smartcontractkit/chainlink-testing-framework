package seth_test

import (
	"os"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	link_token "github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/link"
)

func TestConfig_DefaultClient(t *testing.T) {
	client, err := seth.DefaultClient("ws://localhost:8546", []string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"})
	require.NoError(t, err, "failed to create client with default config")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
}

func TestConfig_Default_TwoPks(t *testing.T) {
	client, err := seth.DefaultClient("ws://localhost:8546", []string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"})
	require.NoError(t, err, "failed to create client with default config")
	require.Equal(t, 2, len(client.PrivateKeys), "expected 2 private keys")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
}

func TestConfig_MinimalBuilder(t *testing.T) {
	builder := seth.NewClientBuilder()

	client, err := builder.WithRpcUrl("ws://localhost:8546").
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()
	require.NoError(t, err, "failed to build client")

	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
}

func TestConfig_MaximalBuilder(t *testing.T) {
	builder := seth.NewClientBuilder()

	client, err := builder.
		// network
		WithNetworkName("my network").
		WithRpcUrl("ws://localhost:8546").
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		WithRpcDialTimeout(10*time.Second).
		WithTransactionTimeout(1*time.Minute).
		// addresses
		WithEphemeralAddresses(10, 10).
		// tracing
		WithTracing(seth.TracingLevel_All, []string{seth.TraceOutput_Console}).
		// protections
		WithProtections(true, true, seth.MustMakeDuration(2*time.Minute)).
		// artifacts folder
		WithArtifactsFolder("some_folder").
		// geth wrappers folders
		WithGethWrappersFolders([]string{"./contracts/bind"}).
		// nonce manager
		WithNonceManager(10, 3, 60, 5).
		Build()

	require.NoError(t, err, "failed to create client")
	require.Equal(t, 11, len(client.PrivateKeys), "expected 11 private keys")

	t.Cleanup(func() {
		err = seth.ReturnFunds(client, client.Addresses[0].Hex())
		require.NoError(t, err, "failed to return funds")
	})

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
}

func TestConfig_ModifyExistingConfigWithBuilder(t *testing.T) {
	configPath := os.Getenv(seth.CONFIG_FILE_ENV_VAR)
	require.NotEmpty(t, configPath, "expected config file path to be set")

	d, err := os.ReadFile(configPath)
	require.NoError(t, err, "failed to read config file")

	var sethConfig seth.Config
	err = toml.Unmarshal(d, &sethConfig)
	require.NoError(t, err, "failed to unmarshal config file")

	client, err := seth.NewClientBuilderWithConfig(&sethConfig).
		UseNetworkWithName(seth.GETH).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.NoError(t, err, "failed to create client")
	require.Equal(t, 2, len(client.PrivateKeys), "expected 11 private keys")
	require.NotNil(t, client.Cfg.Network, "expected network to be set")
	require.Equal(t, uint64(1337), client.Cfg.Network.ChainID, "expected chain ID to be set")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
}

func TestConfig_ModifyExistingConfigWithBuilder_UnknownChainId(t *testing.T) {
	configPath := os.Getenv(seth.CONFIG_FILE_ENV_VAR)
	require.NotEmpty(t, configPath, "expected config file path to be set")

	d, err := os.ReadFile(configPath)
	require.NoError(t, err, "failed to read config file")

	var sethConfig seth.Config
	err = toml.Unmarshal(d, &sethConfig)
	require.NoError(t, err, "failed to unmarshal config file")

	_, err = seth.NewClientBuilderWithConfig(&sethConfig).
		UseNetworkWithChainId(225).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	expectedError := `errors occurred during building the config:
network with chainId '225' not found
at least one method that required network to be set was called, but network is nil`

	require.Error(t, err, "succeeded to create client")
	require.Equal(t, expectedError, err.Error(), "expected error message")
}

func TestConfig_LegacyGas_No_Estimations(t *testing.T) {
	builder := seth.NewClientBuilder()

	client, err := builder.
		// network
		WithNetworkName("my network").
		WithRpcUrl("ws://localhost:8546").
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		// Gas price and estimations
		WithLegacyGasPrice(710_000_000).
		WithGasPriceEstimations(false, 0, "").
		Build()
	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
}

func TestConfig_Eip1559Gas_With_Estimations(t *testing.T) {
	builder := seth.NewClientBuilder()

	client, err := builder.
		// network
		WithNetworkName("my network").
		WithRpcUrl("ws://localhost:8546").
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		// Gas price and estimations
		WithEIP1559DynamicFees(true).
		WithDynamicGasPrices(120_000_000_000, 44_000_000_000).
		WithGasPriceEstimations(true, 10, seth.Priority_Fast).
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
}

func TestConfigAppendPkToEmptyNetwork(t *testing.T) {
	networkName := "network"
	cfg := &seth.Config{
		Network: &seth.Network{
			Name: networkName,
		},
	}

	added := cfg.AppendPksToNetwork([]string{"pk"}, networkName)
	require.True(t, added, "should have added pk to network")
	require.Equal(t, []string{"pk"}, cfg.Network.PrivateKeys, "network should have 1 pk")
}

func TestConfigAppendPkToEmptySharedNetwork(t *testing.T) {
	networkName := "network"
	network := &seth.Network{
		Name: networkName,
	}
	cfg := &seth.Config{
		Network:  network,
		Networks: []*seth.Network{network},
	}

	added := cfg.AppendPksToNetwork([]string{"pk"}, networkName)
	require.True(t, added, "should have added pk to network")
	require.Equal(t, []string{"pk"}, cfg.Network.PrivateKeys, "network should have 1 pk")
	require.Equal(t, []string{"pk"}, cfg.Networks[0].PrivateKeys, "network should have 1 pk")
}

func TestConfigAppendPkToNetworkWithPk(t *testing.T) {
	networkName := "network"
	cfg := &seth.Config{
		Network: &seth.Network{
			Name:        networkName,
			PrivateKeys: []string{"pk1"},
		},
	}

	added := cfg.AppendPksToNetwork([]string{"pk2"}, networkName)
	require.True(t, added, "should have added pk to network")
	require.Equal(t, []string{"pk1", "pk2"}, cfg.Network.PrivateKeys, "network should have 2 pks")
}

func TestConfigAppendPkToMissingNetwork(t *testing.T) {
	networkName := "network"
	cfg := &seth.Config{
		Network: &seth.Network{
			Name: "some_other",
		},
	}

	added := cfg.AppendPksToNetwork([]string{"pk"}, networkName)
	require.False(t, added, "should have not added pk to network")
	require.Equal(t, 0, len(cfg.Network.PrivateKeys), "network should have 0 pks")
}

func TestConfigAppendPkToInactiveNetwork(t *testing.T) {
	networkName := "network"
	cfg := &seth.Config{
		Network: &seth.Network{
			Name: "some_other",
		},
		Networks: []*seth.Network{
			{
				Name: "some_other",
			},
			{
				Name: networkName,
			},
		},
	}

	added := cfg.AppendPksToNetwork([]string{"pk"}, networkName)
	require.True(t, added, "should have added pk to network")
	require.Equal(t, 0, len(cfg.Network.PrivateKeys), "network should have 0 pks")
	require.Equal(t, 0, len(cfg.Networks[0].PrivateKeys), "network should have 0 pks")
	require.Equal(t, []string{"pk"}, cfg.Networks[1].PrivateKeys, "network should have 1 pk")
}
