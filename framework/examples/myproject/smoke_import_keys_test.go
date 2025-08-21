package examples

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type CfgImport struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	NodeSet     *ns.Input         `toml:"nodeset" validate:"required"`
}

func generateP2PKeys(pwd string, n int) ([][]byte, error) {
	encryptedP2PKeyJSONs := make([][]byte, 0)
	for i := 0; i < n; i++ {
		key, err := p2pkey.NewV2()
		if err != nil {
			return nil, err
		}
		d, err := key.ToEncryptedJSON(pwd, utils.DefaultScryptParams)
		if err != nil {
			return nil, err
		}
		encryptedP2PKeyJSONs = append(encryptedP2PKeyJSONs, d)
	}
	return encryptedP2PKeyJSONs, nil
}

func generateEVMKeys(pwd string, n int) ([][]byte, error) {
	encryptedEVMKeyJSONs := make([][]byte, 0)
	for i := 0; i < n; i++ {
		key, _, err := clclient.NewETHKey(pwd)
		if err != nil {
			return nil, err
		}
		encryptedEVMKeyJSONs = append(encryptedEVMKeyJSONs, key)
	}
	return encryptedEVMKeyJSONs, nil
}

func TestImportSmoke(t *testing.T) {
	in, err := framework.Load[CfgImport](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc)
	require.NoError(t, err)
	c, err := clclient.New(out.CLNodes)
	require.NoError(t, err)

	// EVM keys import
	evmKeys, err := generateEVMKeys("", len(out.CLNodes))
	require.NoError(t, err)
	err = clclient.ImportEVMKeys(c, evmKeys, bc.ChainID)
	require.NoError(t, err)
	// p2p keys import
	p2pKeys, err := generateP2PKeys("", len(out.CLNodes))
	require.NoError(t, err)
	err = clclient.ImportP2PKeys(c, p2pKeys)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.ExternalURL)
		}
	})
}
