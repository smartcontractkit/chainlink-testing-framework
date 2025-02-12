package test_env

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/stretchr/testify/require"
)

func TestParrot(t *testing.T) {
	t.Parallel()

	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)

	p := NewParrot([]string{network.Name}).WithTestInstance(t)
	err = p.StartContainer()
	require.NoError(t, err)
}
