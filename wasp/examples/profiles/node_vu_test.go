package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func TestNodeVU(t *testing.T) {
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	labels := map[string]string{
		"branch": "generator_healthcheck",
		"commit": "generator_healthcheck",
	}

	_, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.VU,
			GenName:    "Gamma",
			Schedule:   wasp.Steps(1, 1, 10, 30*time.Second),
			VU:         NewExampleScenario(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.VU,
			GenName:    "Delta",
			Schedule:   wasp.Steps(1, 2, 10, 30*time.Second),
			VU:         NewExampleScenario(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Run(true)
	require.NoError(t, err)
}
