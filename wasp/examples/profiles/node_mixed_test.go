package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func TestNodeMixed(t *testing.T) {
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	labels := map[string]string{
		"branch": "generator_healthcheck",
		"commit": "generator_healthcheck",
	}

	thetaSchedule := wasp.Combine(
		wasp.Steps(1, 1, 10, 40*time.Second),
		wasp.Steps(10, -1, 10, 10*time.Second))

	epsilonSchedule := wasp.Combine(
		wasp.Steps(1, 1, 5, 10*time.Second),
		wasp.Plain(5, 40*time.Second))

	_, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.VU,
			GenName:    "Theta",
			Schedule:   thetaSchedule,
			VU:         NewExampleScenario(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.RPS,
			GenName:    "Epsilon",
			Schedule:   epsilonSchedule,
			Gun:        NewExampleHTTPGun(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Run(true)
	require.NoError(t, err)
}
