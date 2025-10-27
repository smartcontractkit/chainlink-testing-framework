package chaos

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	f "github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/havoc"

	"github.com/stretchr/testify/require"
)

type K8sChaos struct {
	WaitBeforeStart             string   `toml:"wait_before_start"`
	Namespace                   string   `toml:"namespace"`
	DashboardUUIDs              []string `toml:"dashboard_uuids"`
	ExperimentDuration          string   `toml:"experiment_duration"`
	ExperimentInjectionDuration string   `toml:"experiment_injection_duration"`
	BlockchainHTTPURLs          []string `toml:"blockchain_http_urls"`
	ReorgBelowFinalityThreshold int      `toml:"reorg_below_finality_threshold"`
	ReorgAboveFinalityThreshold int      `toml:"reorg_above_finality_threshold"`
	BlockEvery                  string   `toml:"block_every"`
	RemoveK8sChaos              bool     `toml:"remove_k8s_chaos"`
}

type CfgChaosK8s struct {
	Chaos *K8sChaos `toml:"chaos"`
}

func TestK8sChaos(t *testing.T) {
	config, err := f.Load[CfgChaosK8s](t)
	require.NoError(t, err)
	cfg := config.Chaos

	c, err := havoc.NewChaosMeshClient()
	require.NoError(t, err)
	cr := havoc.NewNamespaceRunner(f.L, c, config.Chaos.RemoveK8sChaos)
	gc := f.NewGrafanaClient(os.Getenv("GRAFANA_URL"), os.Getenv("GRAFANA_TOKEN"))
	rpc0 := rpc.New(cfg.BlockchainHTTPURLs[0], nil)
	rpc1 := rpc.New(cfg.BlockchainHTTPURLs[1], nil)

	gasScheduleFunc := func(t *testing.T, r *rpc.RPCClient, url string, increase *big.Int) {
		startGasPrice := big.NewInt(2e9)
		// ramp
		for i := 0; i < 10; i++ {
			err := r.PrintBlockBaseFee()
			require.NoError(t, err)
			err = r.AnvilSetNextBlockBaseFeePerGas(startGasPrice)
			require.NoError(t, err)
			err = r.AnvilMine([]interface{}{"1"})
			require.NoError(t, err)
			time.Sleep(f.MustParseDuration(cfg.BlockEvery))
			startGasPrice = startGasPrice.Add(startGasPrice, increase)
		}
		// hold
		for i := 0; i < 10; i++ {
			err := r.PrintBlockBaseFee()
			require.NoError(t, err)
			err = r.AnvilSetNextBlockBaseFeePerGas(startGasPrice)
			require.NoError(t, err)
			err = r.AnvilMine([]interface{}{"1"})
			require.NoError(t, err)
			time.Sleep(f.MustParseDuration(cfg.BlockEvery))
		}
		// release
		for i := 0; i < 10; i++ {
			err := r.PrintBlockBaseFee()
			require.NoError(t, err)
			time.Sleep(f.MustParseDuration(cfg.BlockEvery))
		}
	}

	testCases := []struct {
		name     string
		run      func(t *testing.T)
		validate func(t *testing.T)
	}{
		// pod failures
		{
			name: "Fail a single pod",
			run: func(t *testing.T) {
				_, err := cr.RunPodFail(context.Background(),
					havoc.PodFailCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"geth-1337"},
						InjectionDuration: f.MustParseDuration(cfg.ExperimentInjectionDuration),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// network delay
		{
			name: "Introduce latency",
			run: func(t *testing.T) {
				_, err := cr.RunPodDelay(context.Background(),
					havoc.PodDelayCfg{
						Namespace:         cfg.Namespace,
						LabelKey:          "app.kubernetes.io/instance",
						LabelValues:       []string{"geth-1337", "geth-2337"},
						Latency:           400 * time.Millisecond,
						Jitter:            20 * time.Millisecond,
						Correlation:       "0",
						InjectionDuration: f.MustParseDuration(cfg.ExperimentInjectionDuration),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// network partition
		{
			name: "Introduce network partition",
			run: func(t *testing.T) {
				_, err := cr.RunPodPartition(context.Background(),
					havoc.PodPartitionCfg{
						Namespace:         cfg.Namespace,
						LabelFromKey:      "app.kubernetes.io/instance",
						LabelFromValues:   []string{"ccip-0"},
						LabelToKey:        "app.kubernetes.io/instance",
						LabelToValues:     []string{"ccip-1", "ccip-2", "ccip-3"},
						InjectionDuration: f.MustParseDuration(cfg.ExperimentInjectionDuration),
					})
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// reorgs
		{
			name: "Reorg src chain below finality",
			run: func(t *testing.T) {
				err := rpc0.GethSetHead(cfg.ReorgBelowFinalityThreshold)
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		// gas spikes
		{
			name: "Slow spike",
			run: func(t *testing.T) {
				gasScheduleFunc(t, rpc1, cfg.BlockchainHTTPURLs[1], big.NewInt(1e9))
			},
			validate: func(t *testing.T) {},
		},
	}

	startsIn := f.MustParseDuration(cfg.WaitBeforeStart)
	f.L.Info().Msgf("Starting chaos tests in %s", startsIn)
	time.Sleep(startsIn)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			n := time.Now()
			testCase.run(t)
			time.Sleep(f.MustParseDuration(cfg.ExperimentDuration))
			if os.Getenv("GRAFANA_URL") != "" {
				_, _, err := gc.Annotate(f.A(cfg.Namespace, testCase.name, cfg.DashboardUUIDs, havoc.Ptr(n), havoc.Ptr(time.Now())))
				require.NoError(t, err)
			}
			require.NoError(t, err)
			testCase.validate(t)
		})
	}
}
