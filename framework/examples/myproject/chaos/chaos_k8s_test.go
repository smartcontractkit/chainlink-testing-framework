package chaos

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/examples/chaos/template"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func defaultLogger() zerolog.Logger {
	return log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
}

func TestK8sChaos(t *testing.T) {
	l := defaultLogger()
	c, err := havoc.NewChaosMeshClient()
	require.NoError(t, err)

	namespace := "default"
	experimentInterval := 3 * time.Minute
	injectionDuration := 90 * time.Second

	reorgBelowFinalityDepth := 10
	reorgAboveFinalityDepth := 50
	srcURL := os.Getenv("CCIP_SRC_CHAIN_HTTP_URL")
	dstURL := os.Getenv("CCIP_DST_CHAIN_HTTP_URL")

	testCases := []struct {
		name     string
		run      func(t *testing.T)
		validate func(t *testing.T)
	}{
		// pod failures
		{
			name: "Fail src chain",
			run: func(t *testing.T) {
				src, err := template.PodFail(c, l, template.PodFailCfg{
					Namespace:         namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"blockchain-1"},
					InjectionDuration: injectionDuration,
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail dst chain",
			run: func(t *testing.T) {
				dst, err := template.PodFail(c, l, template.PodFailCfg{
					Namespace:         namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"blockchain-2"},
					InjectionDuration: injectionDuration,
				})
				require.NoError(t, err)
				dst.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail one node",
			run: func(t *testing.T) {
				node1, err := template.PodFail(c, l, template.PodFailCfg{
					Namespace:         namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"ccip-1"},
					InjectionDuration: injectionDuration,
				})
				require.NoError(t, err)
				node1.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Fail two nodes",
			run: func(t *testing.T) {
				node1, err := template.PodFail(c, l, template.PodFailCfg{
					Namespace:         namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"ccip-1", "ccip-2"},
					InjectionDuration: injectionDuration,
				})
				require.NoError(t, err)
				node1.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		// network delay
		{
			name: "Slow src chain",
			run: func(t *testing.T) {
				src, err := template.PodDelay(c, l, template.PodDelayCfg{
					Namespace:         namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"blockchain-1"},
					Latency:           200 * time.Millisecond,
					Jitter:            200 * time.Millisecond,
					Correlation:       "0",
					InjectionDuration: injectionDuration,
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Slow dst chain",
			run: func(t *testing.T) {
				src, err := template.PodDelay(c, l, template.PodDelayCfg{
					Namespace:         namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"blockchain-2"},
					Latency:           200 * time.Millisecond,
					Jitter:            200 * time.Millisecond,
					Correlation:       "0",
					InjectionDuration: injectionDuration,
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "One slow node",
			run: func(t *testing.T) {
				src, err := template.PodDelay(c, l, template.PodDelayCfg{
					Namespace:         namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"ccip-1"},
					Latency:           200 * time.Millisecond,
					Jitter:            200 * time.Millisecond,
					Correlation:       "0",
					InjectionDuration: injectionDuration,
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Two slow nodes",
			run: func(t *testing.T) {
				src, err := template.PodDelay(c, l, template.PodDelayCfg{
					Namespace:         namespace,
					LabelKey:          "app.kubernetes.io/instance",
					LabelValues:       []string{"ccip-1", "ccip-2"},
					Latency:           200 * time.Millisecond,
					Jitter:            200 * time.Millisecond,
					Correlation:       "0",
					InjectionDuration: injectionDuration,
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "One node partition",
			run: func(t *testing.T) {
				src, err := template.PodPartition(c, l, template.PodPartitionCfg{
					Namespace:         namespace,
					LabelFromKey:      "app.kubernetes.io/instance",
					LabelFromValues:   []string{"ccip-1"},
					LabelToKey:        "app.kubernetes.io/instance",
					LabelToValues:     []string{"ccip-2", "ccip-3", "ccip-4"},
					InjectionDuration: injectionDuration,
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Two nodes partition",
			run: func(t *testing.T) {
				src, err := template.PodPartition(c, l, template.PodPartitionCfg{
					Namespace:         namespace,
					LabelFromKey:      "app.kubernetes.io/instance",
					LabelFromValues:   []string{"ccip-1", "ccip-2"},
					LabelToKey:        "app.kubernetes.io/instance",
					LabelToValues:     []string{"ccip-3", "ccip-4"},
					InjectionDuration: injectionDuration,
				})
				require.NoError(t, err)
				src.Create(context.Background())
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Reorg src chain below finality",
			run: func(t *testing.T) {
				r := rpc.New(srcURL, nil)
				err := r.GethSetHead(reorgBelowFinalityDepth)
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Reorg dst chain below finality",
			run: func(t *testing.T) {
				r := rpc.New(dstURL, nil)
				err := r.GethSetHead(reorgBelowFinalityDepth)
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Reorg src chain above finality",
			run: func(t *testing.T) {
				r := rpc.New(srcURL, nil)
				err := r.GethSetHead(reorgAboveFinalityDepth)
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
		{
			name: "Reorg dst chain above finality",
			run: func(t *testing.T) {
				r := rpc.New(dstURL, nil)
				err := r.GethSetHead(reorgAboveFinalityDepth)
				require.NoError(t, err)
			},
			validate: func(t *testing.T) {},
		},
	}

	// Start WASP load test here, apply average load profile that you expect in production!
	// Configure timeouts and validate all the test cases until the test ends

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
			time.Sleep(experimentInterval)
			tc.validate(t)
		})
	}
}
