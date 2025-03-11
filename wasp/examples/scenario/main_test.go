package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func TestScenario(t *testing.T) {
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	_, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T: t,
			Labels: map[string]string{
				"branch": "generator_healthcheck",
				"commit": "generator_healthcheck",
			},
			LoadType: wasp.VU,
			VU:       NewExampleScenario(srv.URL()),
			Schedule: wasp.Combine(
				wasp.Plain(1, 30*time.Second),
				wasp.Plain(2, 30*time.Second),
				wasp.Plain(3, 30*time.Second),
			),
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).Run(true)
	require.NoError(t, err)
}
