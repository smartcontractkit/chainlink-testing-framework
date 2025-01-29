package examples

import (
	"context"
	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCCIPChaos(t *testing.T) {
	l := defaultLogger()
	c, err := havoc.NewChaosMeshClient()
	require.NoError(t, err)

	namespace := "crib-radek-ccip-v2"

	rebootsChaos, err := podFail(c, l, NodeRebootsConfig{
		Namespace:               namespace,
		Description:             "reboot nodes",
		LabelKey:                "app.kubernetes.io/instance",
		LabelValues:             []string{"ccip-1"},
		ExperimentTotalDuration: 3 * time.Minute,
	})
	require.NoError(t, err)
	latenciesChaos, err := networkDelay(c, l, NodeLatenciesConfig{
		Namespace:               namespace,
		Description:             "network issues",
		Latency:                 5000 * time.Millisecond,
		LatencyDuration:         40 * time.Second,
		FromLabelKey:            "app.kubernetes.io/instance",
		FromLabelValues:         []string{"ccip-1"},
		ToLabelKey:              "app.kubernetes.io/instance",
		ToLabelValues:           []string{"ccip-2", "ccip-3", "ccip-4"},
		ExperimentTotalDuration: 3 * time.Minute,
		//ExperimentCreateDelay:   3 * time.Minute,
	})
	require.NoError(t, err)
	blockchainLatency, err := networkDelay(c, l, NodeLatenciesConfig{
		Namespace:               namespace,
		Description:             "blockchain nodes network issues",
		Latency:                 5000 * time.Millisecond,
		LatencyDuration:         40 * time.Second,
		FromLabelKey:            "app.kubernetes.io/instance",
		FromLabelValues:         []string{"ccip-1", "ccip-2", "ccip-3"},
		ToLabelKey:              "instance",
		ToLabelValues:           []string{"geth-1337"},
		ExperimentTotalDuration: 3 * time.Minute,
		//ExperimentCreateDelay:   3 * time.Minute,
	})
	require.NoError(t, err)

	_ = rebootsChaos
	_ = latenciesChaos
	_ = blockchainLatency

	chaosList := []havoc.ChaosEntity{
		//rebootsChaos,
		//latenciesChaos,
		blockchainLatency,
	}

	for _, chaos := range chaosList {
		chaos.AddListener(havoc.NewChaosLogger(l))
		//chaos.AddListener(havoc.NewSingleLineGrafanaAnnotator(cfg.GrafanaURL, cfg.GrafanaToken, cfg.GrafanaDashboardUID))
		exists, err := havoc.ChaosObjectExists(chaos.GetObject(), c)
		require.NoError(t, err)
		require.False(t, exists, "chaos object already exists: %s. Delete it before starting the test", chaos.GetChaosName())
		chaos.Create(context.Background())
	}
	t.Cleanup(func() {
		for _, chaos := range chaosList {
			chaos.Delete(context.Background())
		}
	})

	// your load test comes here !
	time.Sleep(3*time.Minute + 5*time.Second)
}
