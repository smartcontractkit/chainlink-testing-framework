package seth_test

import (
	"context"
	"os"
	"testing"

	"github.com/smartcontractkit/seth"
	sethcmd "github.com/smartcontractkit/seth/cmd"
	"github.com/stretchr/testify/require"
)

func TestCLITracing(t *testing.T) {
	c := newClientWithContractMapFromEnv(t)
	SkipAnvil(t, c)

	file, err := os.CreateTemp("", "reverted_transactions.json")
	require.NoError(t, err, "should have created temp file")

	tx, txErr := TestEnv.DebugContract.AlwaysRevertsCustomError(c.NewTXOpts())
	require.NoError(t, txErr, "transaction should have reverted")

	_, err = c.WaitMined(context.Background(), seth.L, c.Client, tx)
	require.NoError(t, err, "should have waited for transaction to be mined")

	err = seth.CreateOrAppendToJsonArray(file.Name(), tx.Hash().Hex())
	require.NoError(t, err, "should have written to file")

	_ = os.Setenv(seth.CONFIG_FILE_ENV_VAR, "seth.toml")
	err = sethcmd.RunCLI([]string{"seth", "-n", "Geth", "trace", "-f", file.Name()})
	require.NoError(t, err, "should have traced transactions")
}
