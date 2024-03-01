package test_env

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

func TestEth1WithBesu(t *testing.T) {
	l := logging.GetTestLogger(t)

	builder := NewEthereumNetworkBuilder()
	cfg, err := builder.
		WithConsensusType(ConsensusType_PoW).
		WithExecutionLayer(ExecutionLayer_Besu).
		Build()
	require.NoError(t, err, "Builder validation failed")

	net, _, err := cfg.Start()
	require.NoError(t, err, "Couldn't start PoW network")

	c, err := blockchain.ConnectEVMClient(net, l)
	require.NoError(t, err, "Couldn't connect to the evm client")

	address := common.HexToAddress("0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")
	err = sendAndCompareBalances(testcontext.Get(t), c, address)
	require.NoError(t, err, fmt.Sprintf("balance wasn't correctly updated when %s network", net.Name))

	err = c.Close()
	require.NoError(t, err, "Couldn't close the client")
}
