package examples

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/examples/example_components/onchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type CfgQuickDeploy struct {
	ContractsSrc  *onchain.Input    `toml:"contracts_src" validate:"required"`
	BlockchainSrc *blockchain.Input `toml:"blockchain_src" validate:"required"`
	NodeSet       *ns.Input         `toml:"nodeset" validate:"required"`
}

func randAddr() (string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", err
	}
	publicKey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(publicKey).Hex()
	return address, nil
}

func TestQuickDeploy(t *testing.T) {
	in, err := framework.Load[CfgQuickDeploy](t)
	require.NoError(t, err)

	// spin up 2 anvils
	bcSrc, err := blockchain.NewBlockchainNetwork(in.BlockchainSrc)
	require.NoError(t, err)

	_, err = ns.NewSharedDBNodeSet(in.NodeSet, bcSrc)
	require.NoError(t, err)

	// deploy all the contracts and start periodic mining in tests later
	in.ContractsSrc.URL = bcSrc.Nodes[0].ExternalWSUrl
	c, err := ethclient.Dial(bcSrc.Nodes[0].ExternalWSUrl)
	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		ra, err := randAddr()
		require.NoError(t, err)
		err = simple_node_set.SendETH(c, blockchain.DefaultAnvilPrivateKey, ra, big.NewFloat(0.1))
		require.NoError(t, err)
	}
	// start periodic mining so nodes can receive heads (async)
	miner := rpc.NewRemoteAnvilMiner(bcSrc.Nodes[0].ExternalHTTPUrl, nil)
	miner.MinePeriodically(5 * time.Second)

	t.Run("quickly deploy contracts then test with some block speed", func(t *testing.T) {
		// test your on-chain + off-chain
	})
}
