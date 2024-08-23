package seth_test

import (
	"crypto/ecdsa"
	"github.com/barkimedes/go-deepcopy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/seth"
	"github.com/smartcontractkit/seth/test_utils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestContractMapSavesDeployedContractsToFileAndReadsThem(t *testing.T) {
	file, err := os.CreateTemp("", "deployed_contracts.toml")
	require.NoError(t, err, "failed to create temp file")

	err = seth.SaveDeployedContract(file.Name(), "contractName", "0x0DCd1Bf9A1b36cE34237eEaFef220932846BCD82")
	require.NoError(t, err, "failed to save deployed contract")

	contracts, err := seth.LoadDeployedContracts(file.Name())
	require.NoError(t, err, "failed to load deployed contracts")

	require.Equal(t, map[string]string{"0x0DCd1Bf9A1b36cE34237eEaFef220932846BCD82": "contractName"}, contracts)
}

func TestContractMapDoesNotErrorWhenReadingNonExistentFile(t *testing.T) {
	_, err := seth.LoadDeployedContracts("nonexistent.toml")
	require.NoError(t, err, "reading from non-existent file should not error")
}

func TestContractMapErrorsWheneadingInvalidTomlFile(t *testing.T) {
	file, err := os.CreateTemp("", "invalid_contracts.toml")
	require.NoError(t, err, "failed to create temp file")

	_, err = file.WriteString("invalid toml")
	require.NoError(t, err, "failed to write invalid toml")

	_, err = seth.LoadDeployedContracts(file.Name())
	require.Error(t, err, "expected error reading invalid toml file")
}

func TestContractMapErrorsWhenReadingMalformedAddress(t *testing.T) {
	file, err := os.CreateTemp("", "malformed_address.toml")
	require.NoError(t, err, "failed to create temp file")

	err = seth.SaveDeployedContract(file.Name(), "contractName", "malformed")
	require.NoError(t, err, "failed to save deployed contract")

	_, err = seth.LoadDeployedContracts(file.Name())
	require.Error(t, err, "expected error reading malformed address")
	require.Contains(t, err.Error(), "hex string without 0x prefix", "expected error reading malformed address")
}

func TestContractMapNonSimulatedClientSavesAndReadsContractMap(t *testing.T) {
	file, err := os.CreateTemp("", "deployed_contracts.toml")
	require.NoError(t, err, "failed to create temp file")

	client, err := seth.NewClient()
	require.NoError(t, err, "failed to create client")

	client.Cfg.SaveDeployedContractsMap = true
	client.Cfg.ContractMapFile = file.Name()
	// change network name so that is not treated as simulated
	client.Cfg.Network.Name = "geth2"
	data, err := client.DeployContractFromContractStore(client.NewTXOpts(), "NetworkDebugSubContract")
	require.NoError(t, err, "failed to deploy contract")

	cfg := client.Cfg
	newNonSimulatedClient, err := seth.NewClientRaw(cfg, client.Addresses, client.PrivateKeys)
	require.NoError(t, err, "failed to create new client")
	require.Equal(t, 1, newNonSimulatedClient.ContractAddressToNameMap.Size(), "expected contract map to be saved")

	expectedMap := seth.NewContractMap(map[string]string{data.Address.Hex(): "NetworkDebugSubContract"})
	require.Equal(t, expectedMap, newNonSimulatedClient.ContractAddressToNameMap, "expected contract map to be saved")

	cfg.Network.Name = seth.GETH
	newSimulatedClient, err := seth.NewClientRaw(cfg, client.Addresses, client.PrivateKeys)
	require.NoError(t, err, "failed to create new client")
	require.Equal(t, 0, newSimulatedClient.ContractAddressToNameMap.Size(), "expected contract map to be saved")
}

func TestContractMapSimulatedClientDoesntSaveContractMap(t *testing.T) {
	client, err := seth.NewClient()
	require.NoError(t, err, "failed to create client")

	client.Cfg.SaveDeployedContractsMap = true
	_, err = client.DeployContractFromContractStore(client.NewTXOpts(), "NetworkDebugSubContract")
	require.NoError(t, err, "failed to deploy contract")

	_, err = os.Stat(client.Cfg.GenerateContractMapFileName())
	require.Error(t, err, "contract file should not be saved for simulated network")
	require.True(t, os.IsNotExist(err), "contract file should not be saved for simulated network")
}

func TestContractMapNewClientIsCreatedEvenIfNoContractMapFileExists(t *testing.T) {
	cfg, err := test_utils.CopyConfig(TestEnv.Client.Cfg)
	require.NoError(t, err, "failed to copy config")

	cfg.SaveDeployedContractsMap = true
	// change network name so that is not treated as simulated
	cfg.Network.Name = "geth2"
	cfg.ContractMapFile = cfg.GenerateContractMapFileName()
	nm, err := seth.NewNonceManager(cfg, TestEnv.Client.Addresses, TestEnv.Client.PrivateKeys)
	require.NoError(t, err, "failed to create nonce manager")

	newClient, err := seth.NewClientRaw(cfg, TestEnv.Client.Addresses, TestEnv.Client.PrivateKeys, seth.WithNonceManager(nm), seth.WithContractStore(TestEnv.Client.ContractStore))
	require.NoError(t, err, "failed to create new client")
	require.Equal(t, 0, newClient.ContractAddressToNameMap.Size(), "expected contract map to be saved")

	_, err = newClient.DeployContractFromContractStore(newClient.NewTXOpts(), "NetworkDebugSubContract")
	require.NoError(t, err, "failed to deploy contract")

	t.Cleanup(func() {
		_ = os.Remove(cfg.ContractMapFile)
	})

	// make sure deployed contract is present in the contract map
	require.Equal(t, 1, newClient.ContractAddressToNameMap.Size(), "expected contract map to be saved")

	// make sure that new client instance loads map from existing file instead of creating a new one
	newClient, err = seth.NewClientRaw(cfg, TestEnv.Client.Addresses, TestEnv.Client.PrivateKeys, seth.WithNonceManager(nm), seth.WithContractStore(TestEnv.Client.ContractStore))
	require.NoError(t, err, "failed to create new client")
	require.Equal(t, 1, newClient.ContractAddressToNameMap.Size(), "expected contract map to be saved")
}

func TestContractMapNewClientIsNotCreatedWhenCorruptedContractMapFileExists(t *testing.T) {
	file, err := os.CreateTemp("", "deployed_contracts.toml")
	require.NoError(t, err, "failed to create temp file")

	err = os.WriteFile(file.Name(), []byte("invalid toml"), 0600)
	require.NoError(t, err, "failed to write invalid toml")

	cfg, err := test_utils.CopyConfig(TestEnv.Client.Cfg)
	require.NoError(t, err, "failed to copy config")
	addresses := deepcopy.MustAnything(TestEnv.Client.Addresses).([]common.Address)
	pks := deepcopy.MustAnything(TestEnv.Client.PrivateKeys).([]*ecdsa.PrivateKey)
	// change network name so that is not treated as simulated
	cfg.Network.Name = "geth2"
	cfg.ContractMapFile = file.Name()
	newClient, err := seth.NewClientRaw(cfg, addresses, pks)
	require.Error(t, err, "succeeded in creation of new client")
	require.Contains(t, err.Error(), seth.ErrReadContractMap, "expected error reading invalid toml")
	require.Nil(t, newClient, "expected new client to be nil")
}

func TestContractMapNewClientIsNotCreatedWhenCorruptedContractMapFileExists_InvalidAddress(t *testing.T) {
	file, err := os.CreateTemp("", "deployed_contracts.toml")
	require.NoError(t, err, "failed to create temp file")

	err = seth.SaveDeployedContract(file.Name(), "contractName", "malformed")
	require.NoError(t, err, "failed to write invalid toml")

	cfg, err := test_utils.CopyConfig(TestEnv.Client.Cfg)
	require.NoError(t, err, "failed to copy config")
	addresses := deepcopy.MustAnything(TestEnv.Client.Addresses).([]common.Address)
	pks := deepcopy.MustAnything(TestEnv.Client.PrivateKeys).([]*ecdsa.PrivateKey)
	// change network name so that is not treated as simulated
	cfg.Network.Name = "geth2"
	cfg.ContractMapFile = file.Name()
	newClient, err := seth.NewClientRaw(cfg, addresses, pks)
	require.Error(t, err, "succeeded in creation of new client")
	require.Contains(t, err.Error(), seth.ErrReadContractMap, "expected error reading invalid contract address")
	require.Nil(t, newClient, "expected new client to be nil")
}
