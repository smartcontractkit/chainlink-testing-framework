package benchmark

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/stretchr/testify/require"
)

func TestKeeperBenchmark(t *testing.T) {
	err := actions.RunBenchmarkTest("@benchmark-keeper", "chainlink-benchmark-keeper", 6)
	require.NoError(t, err, "Failed to run the test")
}
