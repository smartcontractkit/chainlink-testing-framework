package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func TestGrafanaOptsVU(t *testing.T) {
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	labels := map[string]string{
		"branch": "grafana_opts",
		"commit": "grafana_opts",
	}

	_, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.VU,
			GenName:    "Nu",
			Schedule:   wasp.Steps(1, 1, 10, 30*time.Second),
			VU:         NewExampleScenario(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.VU,
			GenName:    "Xi",
			Schedule:   wasp.Steps(1, 2, 10, 30*time.Second),
			VU:         NewExampleScenario(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		WithGrafana(&wasp.GrafanaOpts{
			GrafanaURL:                   os.Getenv("GRAFANA_URL"),
			GrafanaToken:                 os.Getenv("GRAFANA_TOKEN"),
			AnnotateDashboardUID:         "Wasp",
			CheckDashboardAlertsAfterRun: "Wasp",
		}).
		Run(true)
	require.NoError(t, err)
}
