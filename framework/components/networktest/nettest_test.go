package networktest_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/networktest"
)

func TestNetworkIsolationRules(t *testing.T) {
	t.Run("can work without DNS isolation", func(t *testing.T) {
		name := "nettest1"
		err := networktest.NewNetworkTest(networktest.Input{NoDNS: false, Name: name})
		require.NoError(t, err)
		dc, err := framework.NewDockerClient()
		require.NoError(t, err)
		sOut, err := dc.ExecContainer(name, []string{"ping", "-c", "1", "google.com"})
		require.NoError(t, err)
		require.NotContains(t, sOut, "bad address")
		fmt.Println(sOut)
	})
	t.Run("DNS isolation works", func(t *testing.T) {
		name := "nettest2"
		err := networktest.NewNetworkTest(networktest.Input{NoDNS: true, Name: name})
		require.NoError(t, err)
		dc, err := framework.NewDockerClient()
		require.NoError(t, err)
		sOut, err := dc.ExecContainer(name, []string{"ping", "-c", "1", "google.com"})
		require.NoError(t, err)
		require.Contains(t, sOut, "bad address")
		fmt.Println(sOut)
	})
}
