package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func TestNodeRPS(t *testing.T) {
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	labels := map[string]string{
		"go_test_name": "test1",
		"branch":       "profile-check",
		"commit":       "profile-check",
	}

	_, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			LoadType: wasp.RPS,
			GenName:  "Alpha",
			Schedule: wasp.Combine(
				wasp.Steps(1, 1, 9, 30*time.Second),
				wasp.Plain(10, 30*time.Second),
				wasp.Steps(10, -1, 10, 30*time.Second),
			),
			Gun:        NewExampleHTTPGun(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Add(wasp.NewGenerator(&wasp.Config{
			LoadType: wasp.RPS,
			GenName:  "Beta",
			Schedule: wasp.Combine(
				wasp.Steps(1, 1, 9, 30*time.Second),
				wasp.Plain(10, 30*time.Second),
				wasp.Steps(10, -1, 10, 30*time.Second),
			),
			Gun:        NewExampleHTTPGun(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Run(true)
	require.NoError(t, err)
}
