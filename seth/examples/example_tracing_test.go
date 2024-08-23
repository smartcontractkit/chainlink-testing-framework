package seth_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/stretchr/testify/require"
)

// All you need to do is enable automated tracing and then wrap you contract call with `c.Decode`.
// This will automatically trace the transaction and decode it for you
func TestDecodeExample(t *testing.T) {
	contract := setup(t)
	c, err := seth.NewClient()
	require.NoError(t, err, "failed to initialise seth")

	// when this level is set we don't need to call TraceGethTX, because it's called automatically
	c.Cfg.TracingLevel = seth.TracingLevel_All

	_, err = c.Decode(contract.TraceDifferent(c.NewTXOpts(), big.NewInt(1), big.NewInt(2)))
	require.NoError(t, err, "failed to decode transaction")
}
