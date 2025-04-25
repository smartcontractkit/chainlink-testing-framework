package main_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
)

func TestChaosSample(t *testing.T) {
	l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	c, err := havoc.NewChaosMeshClient()
	require.NoError(t, err)
	r := havoc.NewNamespaceRunner(l, c, false)
	dur := 1 * time.Minute

	// choose any experiment type
	_, err = r.RunPodDelay(context.Background(),
		havoc.PodDelayCfg{
			// fill your target here
			Namespace:         "crib-aw-remote",
			LabelKey:          "app.kubernetes.io/instance",
			LabelValues:       []string{"ccip-15"},
			Latency:           200 * time.Millisecond,
			Jitter:            200 * time.Millisecond,
			Correlation:       "0",
			InjectionDuration: dur,
		})
	require.NoError(t, err)
	time.Sleep(dur)
}
