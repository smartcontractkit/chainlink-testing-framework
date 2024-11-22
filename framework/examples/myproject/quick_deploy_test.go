package examples

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/examples/example_components/onchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

type CfgQuickDeploy struct {
	ContractsSrc  *onchain.Input    `toml:"contracts_src" validate:"required"`
	BlockchainSrc *blockchain.Input `toml:"blockchain_src" validate:"required"`
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

	// deploy 2 example product contracts
	// you can replace it with chainlink-deployments
	in.ContractsSrc.URL = bcSrc.Nodes[0].HostWSUrl
	c, err := ethclient.Dial(bcSrc.Nodes[0].HostWSUrl)
	require.NoError(t, err)

	t.Run("test some contracts with fork", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			ra, err := randAddr()
			require.NoError(t, err)
			err = simple_node_set.SendETH(c, blockchain.DefaultAnvilPrivateKey, ra, big.NewFloat(0.1))
			require.NoError(t, err)
		}
		t.Log("now mining")

		miner := rpc.NewRemoteAnvilMiner(bcSrc.Nodes[0].HostHTTPUrl, nil)
		miner.MinePeriodically(200 * time.Millisecond)
		time.Sleep(10 * time.Second)
	})
}
