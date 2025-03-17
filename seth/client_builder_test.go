package seth_test

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/pelletier/go-toml/v2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	link_token "github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/link"
)

func TestConfig_MinimalBuilder(t *testing.T) {
	builder := seth.NewClientBuilder()

	client, err := builder.WithRpcUrl(os.Getenv("SETH_URL")).
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

	firstNetwork := &seth.Network{
		Name:                           "First",
		EIP1559DynamicFees:             true,
		TxnTimeout:                     seth.MustMakeDuration(5 * time.Minute),
		DialTimeout:                    seth.MustMakeDuration(seth.DefaultDialTimeout),
		TransferGasFee:                 seth.DefaultTransferGasFee,
		GasPriceEstimationEnabled:      true,
		GasPriceEstimationBlocks:       200,
		GasPriceEstimationTxPriority:   seth.Priority_Standard,
		GasPrice:                       seth.DefaultGasPrice,
		GasFeeCap:                      seth.DefaultGasFeeCap,
		GasTipCap:                      seth.DefaultGasTipCap,
		GasPriceEstimationAttemptCount: seth.DefaultGasPriceEstimationsAttemptCount,
	}

	secondNetwork := &seth.Network{
		Name:                           "Second",
		EIP1559DynamicFees:             true,
		TxnTimeout:                     seth.MustMakeDuration(5 * time.Minute),
		DialTimeout:                    seth.MustMakeDuration(seth.DefaultDialTimeout),
		TransferGasFee:                 seth.DefaultTransferGasFee,
		GasPriceEstimationEnabled:      true,
		GasPriceEstimationBlocks:       200,
		GasPriceEstimationTxPriority:   seth.Priority_Standard,
		GasPrice:                       seth.DefaultGasPrice,
		GasFeeCap:                      seth.DefaultGasFeeCap,
		GasTipCap:                      seth.DefaultGasTipCap,
		GasPriceEstimationAttemptCount: seth.DefaultGasPriceEstimationsAttemptCount,
		URLs:                           []string{os.Getenv("SETH_URL")},
	}

	client, err := builder.
		// network
		WithNetworks([]*seth.Network{firstNetwork, secondNetwork}).
		UseNetworkWithName("Second").
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
	require.Equal(t, 2, len(client.Cfg.Networks), "expected 2 networks")
	require.Equal(t, "Second", client.Cfg.Network.Name, "expected network to be set")

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
	if strings.EqualFold(os.Getenv("SETH_NETWORK"), "anvil") {
		t.Skip("skipping test in anvil network")
	}
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

	// remove default network
	networks := []*seth.Network{}
	for _, network := range sethConfig.Networks {
		if network.Name != seth.DefaultNetworkName {
			networks = append(networks, network)
		}
	}

	sethConfig.Networks = networks

	_, err = seth.NewClientBuilderWithConfig(&sethConfig).
		UseNetworkWithChainId(225).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	expectedError := `errors occurred during building the config:
network with chainId '225' not found
at least one method that required network to be set was called, but network is nil
you need to set the Network`

	require.Error(t, err, "succeeded to create client")
	require.Equal(t, expectedError, err.Error(), "expected error message")
}

func TestConfig_ModifyExistingConfigWithBuilder_UnknownChainId_UseDefault(t *testing.T) {
	configPath := os.Getenv(seth.CONFIG_FILE_ENV_VAR)
	require.NotEmpty(t, configPath, "expected config file path to be set")

	d, err := os.ReadFile(configPath)
	require.NoError(t, err, "failed to read config file")

	var sethConfig seth.Config
	err = toml.Unmarshal(d, &sethConfig)
	require.NoError(t, err, "failed to unmarshal config file")

	_, err = seth.NewClientBuilderWithConfig(&sethConfig).
		UseNetworkWithChainId(225).
		WithRpcUrl(os.Getenv("SETH_URL")).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	expectedError := `errors occurred during building the config:
network with chainId '225' not found
at least one method that required network to be set was called, but network is nil
at least one method that required network to be set was called, but network is nil
you need to set the Network`

	require.Error(t, err, "succeeded to create client")
	require.Equal(t, expectedError, err.Error(), "expected error message")
}

func TestConfig_LegacyGas_No_Estimations(t *testing.T) {
	builder := seth.NewClientBuilder()

	client, err := builder.
		// network
		WithNetworkName("my network").
		WithRpcUrl(os.Getenv("SETH_URL")).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		// Gas price and estimations
		WithLegacyGasPrice(710_000_000).
		WithGasPriceEstimations(false, 0, "", 0).
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
		WithRpcUrl(os.Getenv("SETH_URL")).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		// Gas price and estimations
		WithEIP1559DynamicFees(true).
		WithDynamicGasPrices(120_000_000_000, 44_000_000_000).
		WithGasPriceEstimations(true, 10, seth.Priority_Fast, seth.DefaultGasPriceEstimationsAttemptCount).
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
}

func TestConfig_NoPrivateKeys_RpcHealthEnabled(t *testing.T) {
	builder := seth.NewClientBuilder()

	_, err := builder.
		// network
		WithNetworkName("my network").
		WithRpcUrl(os.Getenv("SETH_URL")).
		// Gas price and estimations
		WithEIP1559DynamicFees(true).
		WithDynamicGasPrices(120_000_000_000, 44_000_000_000).
		WithGasPriceEstimations(false, 10, seth.Priority_Fast, seth.DefaultGasPriceEstimationsAttemptCount).
		Build()

	require.Error(t, err, "succeeded in building the client")
	require.Contains(t, err.Error(), seth.NoPkForRpcHealthCheckErr, "expected error message")
}

func TestConfig_NoPrivateKeys_PendingNonce(t *testing.T) {
	builder := seth.NewClientBuilder()

	_, err := builder.
		// network
		WithNetworkName("my network").
		WithRpcUrl(os.Getenv("SETH_URL")).
		// Gas price and estimations
		WithEIP1559DynamicFees(true).
		WithDynamicGasPrices(120_000_000_000, 44_000_000_000).
		WithGasPriceEstimations(false, 10, seth.Priority_Fast, 2).
		WithProtections(true, false, seth.MustMakeDuration(1*time.Minute)).
		Build()

	require.Error(t, err, "succeeded in building the client")
	require.Contains(t, err.Error(), seth.NoPkForNonceProtection, "expected error message")
}

func TestConfig_NoPrivateKeys_EphemeralKeys(t *testing.T) {
	builder := seth.NewClientBuilder()

	_, err := builder.
		// network
		WithNetworkName("my network").
		WithRpcUrl(os.Getenv("SETH_URL")).
		WithEphemeralAddresses(10, 1000).
		// Gas price and estimations
		WithEIP1559DynamicFees(true).
		WithDynamicGasPrices(120_000_000_000, 44_000_000_000).
		WithGasPriceEstimations(false, 10, seth.Priority_Fast, 2).
		WithProtections(false, false, seth.MustMakeDuration(1*time.Minute)).
		Build()

	require.Error(t, err, "succeeded in building the client")
	require.Contains(t, err.Error(), seth.NoPkForEphemeralKeys, "expected error message")
}

func TestConfig_NoPrivateKeys_GasEstimations(t *testing.T) {
	builder := seth.NewClientBuilder()

	_, err := builder.
		WithNetworkName("my network").
		WithRpcUrl(os.Getenv("SETH_URL")).
		WithGasPriceEstimations(true, 10, seth.Priority_Fast, 2).
		WithProtections(false, false, seth.MustMakeDuration(1*time.Minute)).
		Build()

	require.Error(t, err, "succeeded in building the client")
	require.Contains(t, err.Error(), seth.NoPkForGasPriceEstimation, "expected error message")
}

func TestConfig_NoPrivateKeys_TxOpts(t *testing.T) {
	builder := seth.NewClientBuilder()

	client, err := builder.
		// network
		WithNetworkName("my network").
		WithRpcUrl(os.Getenv("SETH_URL")).
		// Gas price and estimations
		WithEIP1559DynamicFees(true).
		WithDynamicGasPrices(120_000_000_000, 44_000_000_000).
		WithGasPriceEstimations(false, 10, seth.Priority_Fast, 2).
		WithProtections(false, false, seth.MustMakeDuration(1*time.Minute)).
		Build()

	require.NoError(t, err, "failed to the client")
	require.Equal(t, 0, len(client.PrivateKeys), "expected 0 private keys")

	_ = client.NewTXOpts()
	require.Equal(t, 1, len(client.Errors), "expected 1 error")
	require.Equal(t, "no private keys were loaded, but keyNum 0 was requested", client.Errors[0].Error(), "expected error message")
}

func TestConfig_NoPrivateKeys_Tracing(t *testing.T) {
	if strings.EqualFold(os.Getenv("SETH_NETWORK"), "anvil") {
		t.Skip("skipping tracing test in anvil network")
	}
	builder := seth.NewClientBuilder()

	client, err := builder.
		WithNetworkName("my network").
		WithRpcUrl(os.Getenv("SETH_URL")).
		WithEIP1559DynamicFees(true).
		WithDynamicGasPrices(120_000_000_000, 44_000_000_000).
		WithGasPriceEstimations(false, 10, seth.Priority_Fast, 2).
		WithProtections(false, false, seth.MustMakeDuration(1*time.Minute)).
		WithGethWrappersFolders([]string{"./contracts/bind"}).
		Build()

	require.NoError(t, err, "failed to the client")
	require.Equal(t, 0, len(client.PrivateKeys), "expected 0 private keys")

	ethClient, err := ethclient.Dial(os.Getenv("SETH_URL"))
	require.NoError(t, err, "failed to dial eth client")

	pk, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	require.NoError(t, err, "failed to parse private key")

	opts, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(1337))
	require.NoError(t, err, "failed to create transactor")

	addr, tx, instance, err := link_token.DeployLinkToken(opts, ethClient)
	require.NoError(t, err, "failed to deploy LINK contract")

	// it's a deployment transaction, we don't know yet how to decode it
	_, decodeErr := client.DecodeTx(tx)
	require.NoError(t, decodeErr, "failed to decode transaction")

	publicKeyECDSA, ok := pk.Public().(*ecdsa.PublicKey)
	require.True(t, ok, "failed to cast public key to ECDSA")
	pubKeyAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	tx, err = instance.GrantMintRole(opts, pubKeyAddress)
	decoded, decodeErr := client.Decode(tx, err)
	require.NoError(t, decodeErr, "failed to decode transaction")
	require.NotNil(t, decoded, "expected decoded call")
	require.Equal(t, "c2e3273d", decoded.Signature, "signature mismatch")
	require.Equal(t, "grantMintRole(address)", decoded.Method, "method mismatch")
	require.Equal(t, common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"), decoded.Input["minter"], "minter mismatch")
	require.Equal(t, 1, len(decoded.Events), "expected 1 event")
	require.Equal(t, "MintAccessGranted(address)", decoded.Events[0].Signature, "event signature mismatch")
	require.Equal(t, addr, decoded.Events[0].Address, "event address mismatch")
	require.Equal(t, common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"), decoded.Events[0].EventData["minter"], "event minter mismatch")
}

func TestConfig_ReadOnlyMode(t *testing.T) {
	builder := seth.NewClientBuilder()

	client, err := builder.
		WithNetworkName("my network").
		WithRpcUrl(os.Getenv("SETH_URL")).
		WithEphemeralAddresses(10, 1000).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		WithGasPriceEstimations(true, 10, seth.Priority_Fast, 2).
		WithReadOnlyMode().
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 0, len(client.PrivateKeys), "expected 0 private keys")
	require.Equal(t, 0, len(client.Addresses), "expected 0 addresses")
	require.False(t, client.Cfg.CheckRpcHealthOnStart, "expected rpc health check to be disabled")
	require.False(t, client.Cfg.PendingNonceProtectionEnabled, "expected pending nonce protection to be disabled")
	require.False(t, client.Cfg.Network.GasPriceEstimationEnabled, "expected gas price estimations to be disabled")
}

func TestConfig_SimulatedBackend(t *testing.T) {
	backend, cancelFn := StartSimulatedBackend([]common.Address{common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")})
	t.Cleanup(func() {
		cancelFn()
	})

	builder := seth.NewClientBuilder()

	client, err := builder.
		WithNetworkName("simulated").
		WithEthClient(backend.Client()).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")
	require.Equal(t, 1, len(client.Addresses), "expected 1 addresse")
	require.IsType(t, backend.Client(), client.Client, "expected simulated client")
}

func TestConfig_SimulatedBackend_ContractDeploymentHooks(t *testing.T) {
	backend, cancelFn := StartSimulatedBackend([]common.Address{common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")})
	t.Cleanup(func() {
		cancelFn()
	})

	builder := seth.NewClientBuilder()

	wasContractPreHookCalled := false
	wasContractPostHookCalled := false
	order := []string{}

	hooks := seth.Hooks{
		ContractDeployment: seth.ContractDeploymentHooks{
			Pre: func(_ *bind.TransactOpts, _ string, _ abi.ABI, _ []byte, _ ...interface{}) error {
				wasContractPreHookCalled = true
				order = append(order, "pre")
				return nil
			},
			Post: func(_ *seth.Client, _ *types.Transaction) error {
				backend.Commit()
				wasContractPostHookCalled = true
				order = append(order, "post")
				return nil
			},
		},
	}

	client, err := builder.
		WithNetworkName("simulated").
		WithHooks(hooks).
		WithEthClient(backend.Client()).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")
	require.Equal(t, 1, len(client.Addresses), "expected 1 addresse")
	require.IsType(t, backend.Client(), client.Client, "expected simulated client")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
	require.True(t, wasContractPreHookCalled, "expected contract deployment pre hook to be called")
	require.True(t, wasContractPostHookCalled, "expected contract deployment post hook to be called")
	require.Equal(t, []string{"pre", "post"}, order, "expected order to be preserved")
}

func TestConfig_SimulatedBackend_ContractDeploymentHooks_PreError(t *testing.T) {
	backend, cancelFn := StartSimulatedBackend([]common.Address{common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")})
	t.Cleanup(func() {
		cancelFn()
	})

	builder := seth.NewClientBuilder()

	wasContractPreHookCalled := false
	wasContractPostHookCalled := false
	order := []string{}

	hooks := seth.Hooks{
		ContractDeployment: seth.ContractDeploymentHooks{
			Pre: func(_ *bind.TransactOpts, _ string, _ abi.ABI, _ []byte, _ ...interface{}) error {
				wasContractPreHookCalled = true
				order = append(order, "pre")
				return errors.New("pre hook error")
			},
			Post: func(_ *seth.Client, _ *types.Transaction) error {
				backend.Commit()
				wasContractPostHookCalled = true
				order = append(order, "post")
				return nil
			},
		},
	}

	client, err := builder.
		WithNetworkName("simulated").
		WithHooks(hooks).
		WithEthClient(backend.Client()).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")
	require.Equal(t, 1, len(client.Addresses), "expected 1 addresse")
	require.IsType(t, backend.Client(), client.Client, "expected simulated client")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.Error(t, err, "succeeded in deploying LINK contract")
	require.Contains(t, err.Error(), "pre hook error", "expected error message")
	require.True(t, wasContractPreHookCalled, "expected contract deployment pre hook to be called")
	require.False(t, wasContractPostHookCalled, "did no expect contract deployment post hook to be called")
	require.Equal(t, []string{"pre"}, order, "expected order to be preserved")
}

func TestConfig_SimulatedBackend_ContractDeploymentHooks_PostError(t *testing.T) {
	backend, cancelFn := StartSimulatedBackend([]common.Address{common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")})
	t.Cleanup(func() {
		cancelFn()
	})

	builder := seth.NewClientBuilder()

	wasContractPreHookCalled := false
	wasContractPostHookCalled := false
	order := []string{}

	hooks := seth.Hooks{
		ContractDeployment: seth.ContractDeploymentHooks{
			Pre: func(_ *bind.TransactOpts, _ string, _ abi.ABI, _ []byte, _ ...interface{}) error {
				wasContractPreHookCalled = true
				order = append(order, "pre")
				return nil
			},
			Post: func(_ *seth.Client, _ *types.Transaction) error {
				backend.Commit()
				wasContractPostHookCalled = true
				order = append(order, "post")
				return errors.New("post hook error")
			},
		},
	}

	client, err := builder.
		WithNetworkName("simulated").
		WithHooks(hooks).
		WithEthClient(backend.Client()).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")
	require.Equal(t, 1, len(client.Addresses), "expected 1 addresse")
	require.IsType(t, backend.Client(), client.Client, "expected simulated client")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.Error(t, err, "succeeded in deploying LINK contract")
	require.Contains(t, err.Error(), "post hook error", "expected error message")
	require.True(t, wasContractPreHookCalled, "expected contract deployment pre hook to be called")
	require.True(t, wasContractPostHookCalled, "expected contract deployment post hook to be called")
	require.Equal(t, []string{"pre", "post"}, order, "expected order to be preserved")
}

func TestConfig_SimulatedBackend_TxDecodingHooks(t *testing.T) {
	backend, cancelFn := StartSimulatedBackend([]common.Address{common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")})
	t.Cleanup(func() {
		cancelFn()
	})

	builder := seth.NewClientBuilder()

	wasTxDecodingPreHookCalled := false
	wasTxDecodingPostHookCalled := false
	order := []string{}

	hooks := seth.Hooks{
		TxDecoding: seth.TxDecodingHooks{
			Pre: func(_ *seth.Client) error {
				backend.Commit()
				wasTxDecodingPreHookCalled = true
				order = append(order, "pre")
				return nil
			},
			Post: func(_ *seth.Client, _ *seth.DecodedTransaction, _ error) error {
				wasTxDecodingPostHookCalled = true
				order = append(order, "post")
				return nil
			},
		},
	}

	client, err := builder.
		WithNetworkName("simulated").
		WithHooks(hooks).
		WithEthClient(backend.Client()).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")
	require.Equal(t, 1, len(client.Addresses), "expected 1 address")
	require.IsType(t, backend.Client(), client.Client, "expected simulated client")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	data, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")

	instance, err := link_token.NewLinkToken(data.Address, client.Client)
	require.NoError(t, err, "failed to get LINK instance")

	_, err = client.Decode(instance.GrantMintAndBurnRoles(client.NewTXOpts(), common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")))
	require.NoError(t, err, "failed to decode transaction")

	require.True(t, wasTxDecodingPreHookCalled, "expected tx decoding pre hook to be called")
	require.True(t, wasTxDecodingPostHookCalled, "expected tx decoding post hook to be called")
	require.Equal(t, []string{"pre", "post"}, order, "expected order to be preserved")
}

func TestConfig_SimulatedBackend_TxDecodingHooks_PreError(t *testing.T) {
	backend, cancelFn := StartSimulatedBackend([]common.Address{common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")})
	t.Cleanup(func() {
		cancelFn()
	})

	builder := seth.NewClientBuilder()

	wasTxDecodingPreHookCalled := false
	wasTxDecodingPostHookCalled := false
	order := []string{}

	hooks := seth.Hooks{
		TxDecoding: seth.TxDecodingHooks{
			Pre: func(_ *seth.Client) error {
				backend.Commit()
				wasTxDecodingPreHookCalled = true
				order = append(order, "pre")
				return errors.New("pre hook error")
			},
			Post: func(_ *seth.Client, _ *seth.DecodedTransaction, _ error) error {
				wasTxDecodingPostHookCalled = true
				order = append(order, "post")
				return nil
			},
		},
	}

	client, err := builder.
		WithNetworkName("simulated").
		WithHooks(hooks).
		WithEthClient(backend.Client()).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")
	require.Equal(t, 1, len(client.Addresses), "expected 1 address")
	require.IsType(t, backend.Client(), client.Client, "expected simulated client")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	data, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")

	instance, err := link_token.NewLinkToken(data.Address, client.Client)
	require.NoError(t, err, "failed to get LINK instance")

	_, err = client.Decode(instance.GrantMintAndBurnRoles(client.NewTXOpts(), common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")))
	require.Error(t, err, "succeeded in decoding transaction")
	require.Contains(t, err.Error(), "pre hook error", "expected error message")
	require.True(t, wasTxDecodingPreHookCalled, "expected tx decoding pre hook to be called")
	require.False(t, wasTxDecodingPostHookCalled, "expected tx decoding post hook to be called")
	require.Equal(t, []string{"pre"}, order, "expected order to be preserved")
}

func TestConfig_SimulatedBackend_TxDecodingHooks_PostError(t *testing.T) {
	backend, cancelFn := StartSimulatedBackend([]common.Address{common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")})
	t.Cleanup(func() {
		cancelFn()
	})

	builder := seth.NewClientBuilder()

	wasTxDecodingPreHookCalled := false
	wasTxDecodingPostHookCalled := false
	order := []string{}

	hooks := seth.Hooks{
		TxDecoding: seth.TxDecodingHooks{
			Pre: func(_ *seth.Client) error {
				backend.Commit()
				wasTxDecodingPreHookCalled = true
				order = append(order, "pre")
				return nil
			},
			Post: func(_ *seth.Client, _ *seth.DecodedTransaction, _ error) error {
				wasTxDecodingPostHookCalled = true
				order = append(order, "post")
				return errors.New("post hook error")
			},
		},
	}

	client, err := builder.
		WithNetworkName("simulated").
		WithHooks(hooks).
		WithEthClient(backend.Client()).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")
	require.Equal(t, 1, len(client.Addresses), "expected 1 address")
	require.IsType(t, backend.Client(), client.Client, "expected simulated client")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	data, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")

	instance, err := link_token.NewLinkToken(data.Address, client.Client)
	require.NoError(t, err, "failed to get LINK instance")

	_, err = client.Decode(instance.GrantMintAndBurnRoles(client.NewTXOpts(), common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")))
	require.Error(t, err, "succeeded in decoding transaction")
	require.Contains(t, err.Error(), "post hook error", "expected error message")
	require.True(t, wasTxDecodingPreHookCalled, "expected tx decoding pre hook to be called")
	require.True(t, wasTxDecodingPostHookCalled, "expected tx decoding post hook to be called")
	require.Equal(t, []string{"pre", "post"}, order, "expected order to be preserved")
}

func TestConfig_EthClient_DoesntAllowRpcUrl(t *testing.T) {
	backend, cancelFn := StartSimulatedBackend([]common.Address{common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")})
	t.Cleanup(func() {
		cancelFn()
	})

	builder := seth.NewClientBuilder()

	client, err := builder.
		WithNetworkName("simulated").
		WithRpcUrl("ws://localhost:8546").
		WithEthClient(backend.Client()).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.Error(t, err, "failed to build client")
	require.Contains(t, err.Error(), seth.EthClientAndUrlsSet, "expected error message")
	require.Nil(t, client, "expected client to be nil")
}

func TestConfig_EthClient(t *testing.T) {
	builder := seth.NewClientBuilder()

	ethclient, err := ethclient.Dial(os.Getenv("SETH_URL"))
	require.NoError(t, err, "failed to dial eth client")

	client, err := builder.
		WithNetworkName("my network").
		WithEthClient(ethclient).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.NoError(t, err, "failed to build client")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")
	require.Equal(t, 1, len(client.Addresses), "expected 1 address")
	require.IsType(t, ethclient, client.Client, "expected real client")
}

func TestConfig_UnknownNetwork(t *testing.T) {
	builder := seth.NewClientBuilder()

	client, err := builder.
		UseNetworkWithName("my network").
		WithNetworks([]*seth.Network{}).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		Build()

	require.Error(t, err, "succeeded in building the client")
	require.Contains(t, err.Error(), "network with name 'my network' not found", "expected error message")
	require.Nil(t, client, "expected client to be nil")
}
