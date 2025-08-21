package examples

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
)

type CfgJD struct {
	Blockchain *blockchain.Input `toml:"blockchain" validate:"required"`
	JD         *jd.Input         `toml:"jd" validate:"required"`
}

func TestJD(t *testing.T) {
	in, err := framework.Load[CfgJD](t)
	require.NoError(t, err)

	bcOut, err := blockchain.NewBlockchainNetwork(in.Blockchain)
	require.NoError(t, err)

	jdOut, err := jd.NewJD(in.JD)
	require.NoError(t, err)

	t.Run("test changesets with forked network/JD state", func(t *testing.T) {
		_ = bcOut
		_ = jdOut
	})
}
