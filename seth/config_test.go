package seth_test

import (
	"crypto/ecdsa"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	link_token "github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/link"
)

func TestConfig_DefaultClient(t *testing.T) {
	client, err := seth.DefaultClient(os.Getenv("SETH_URL"), []string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"})
	require.NoError(t, err, "failed to create client with default config")
	require.Equal(t, 1, len(client.PrivateKeys), "expected 1 private key")

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
}

func TestConfig_Default_TwoPks(t *testing.T) {
	client, err := seth.DefaultClient(os.Getenv("SETH_URL"), []string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"})
	require.NoError(t, err, "failed to create client with default config")
	require.Equal(t, 2, len(client.PrivateKeys), "expected 2 private keys")

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

func TestConfig_ReadOnly_WithPk(t *testing.T) {
	cfg := seth.Config{
		ReadOnly: true,
		Network: &seth.Network{
			Name: "some_other",
			URLs: []string{os.Getenv("SETH_URL")},
		},
	}

	addrs := []common.Address{common.HexToAddress("0xb794f5ea0ba39494ce839613fffba74279579268")}

	_, err := seth.NewClientRaw(&cfg, addrs, nil)
	require.Error(t, err, "succeeded in creating client")
	require.Equal(t, seth.ErrReadOnlyWithPrivateKeys, err.Error(), "expected different error message")

	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	require.NoError(t, err, "failed to parse private key")

	pks := []*ecdsa.PrivateKey{privateKey}
	_, err = seth.NewClientRaw(&cfg, nil, pks)
	require.Error(t, err, "succeeded in creating client")
	require.Equal(t, seth.ErrReadOnlyWithPrivateKeys, err.Error(), "expected different error message")

	_, err = seth.NewClientRaw(&cfg, addrs, pks)
	require.Error(t, err, "succeeded in creating client")
	require.Equal(t, seth.ErrReadOnlyWithPrivateKeys, err.Error(), "expected different error message")
}

func TestConfig_ReadOnly_GasBumping(t *testing.T) {
	cfg := seth.Config{
		ReadOnly: true,
		Network: &seth.Network{
			Name:        "some_other",
			URLs:        []string{os.Getenv("SETH_URL")},
			DialTimeout: &seth.Duration{D: 10 * time.Second},
		},
		GasBump: &seth.GasBumpConfig{
			Retries: uint(1),
		},
	}

	_, err := seth.NewClientRaw(&cfg, nil, nil)
	require.Error(t, err, "succeeded in creating client")
	require.Equal(t, seth.ErrReadOnlyGasBumping, err.Error(), "expected different error message")
}

func TestConfig_ReadOnly_RpcHealth(t *testing.T) {
	cfg := seth.Config{
		ReadOnly:              true,
		CheckRpcHealthOnStart: true,
		Network: &seth.Network{
			Name:        "some_other",
			URLs:        []string{os.Getenv("SETH_URL")},
			DialTimeout: &seth.Duration{D: 10 * time.Second},
		},
	}

	_, err := seth.NewClientRaw(&cfg, nil, nil)
	require.Error(t, err, "succeeded in creating client")
	require.Equal(t, seth.ErrReadOnlyRpcHealth, err.Error(), "expected different error message")
}

func TestConfig_ReadOnly_PendingNonce(t *testing.T) {
	cfg := seth.Config{
		ReadOnly:                      true,
		PendingNonceProtectionEnabled: true,
		Network: &seth.Network{
			Name:        "some_other",
			URLs:        []string{os.Getenv("SETH_URL")},
			DialTimeout: &seth.Duration{D: 10 * time.Second},
		},
	}

	_, err := seth.NewClientRaw(&cfg, nil, nil)
	require.Error(t, err, "succeeded in creating client")
	require.Equal(t, seth.ErrReadOnlyPendingNonce, err.Error(), "expected different error message")
}

func TestConfig_ReadOnly_EphemeralKeys(t *testing.T) {
	ten := int64(10)
	cfg := seth.Config{
		ReadOnly:       true,
		EphemeralAddrs: &ten,
		Network: &seth.Network{
			Name:        "some_other",
			URLs:        []string{os.Getenv("SETH_URL")},
			DialTimeout: &seth.Duration{D: 10 * time.Second},
		},
	}

	_, err := seth.NewClientRaw(&cfg, nil, nil)
	require.Error(t, err, "succeeded in creating client")
	require.Equal(t, seth.ErrNoPksEphemeralMode, err.Error(), "expected different error message")
}
