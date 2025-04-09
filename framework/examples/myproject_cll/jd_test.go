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
	dc, err := framework.NewDockerClient()
	require.NoError(t, err)
	// find what to dump, RDS API here instead?
	_, err = dc.ExecContainer(jdOut.DBContainerName, []string{"pg_dump", "-U", "chainlink", "-h", "localhost", "-p", "5432", "-d", "chainlink", "-F", "c", "-f", "jd.dump"})
	require.NoError(t, err)

	// copy your dump
	err = dc.CopyFile(jdOut.DBContainerName, "jd.dump", "/")
	require.NoError(t, err)

	// restore
	_, err = dc.ExecContainer(jdOut.DBContainerName, []string{"pg_restore", "-U", "chainlink", "-d", "chainlink", "jd.dump"})
	require.NoError(t, err)

	t.Run("test changesets with forked network and real JD state", func(t *testing.T) {
		_ = bcOut
		_ = jdOut
	})
}
