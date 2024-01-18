package logstream_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

/* These tests are for testing Loki format, they rarely change so you can run them manually */

func TestExampleLokiStreaming(t *testing.T) {
	t.Skip("uncomment and run manually")
	tests := []testData{
		{
			name:      "stream all container logs to Loki, subtest 1",
			repeat:    1,
			perSecond: 0.01,
			streams:   []string{"A\nB\nC\nD", "E\nF\nG\nH"},
		},
		{
			name:      "stream all container logs to Loki, subtest 2",
			repeat:    1,
			perSecond: 0.01,
			streams:   []string{"1\n2\n3\n4", "5\n6\n7\n8"},
		},
		{
			name:      "nobody expects the spammish repetition",
			repeat:    1000,
			perSecond: 0.0001,
			streams: []string{
				"nobody expects the spanish inquisition",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			d, err := NewDeployment(ctx, tc)
			// nolint
			defer d.Shutdown(ctx)
			require.NoError(t, err)

			loggingConfig := config.LoggingConfig{}
			loggingConfig.LogStream = &config.LogStreamConfig{
				LogTargets:            []string{"loki"},
				LogProducerTimeout:    &blockchain.StrDuration{Duration: 10 * time.Second},
				LogProducerRetryLimit: ptr.Ptr(uint(10)),
			}
			loggingConfig.Loki = &config.LokiConfig{
				TenantId: ptr.Ptr("CHANGE-ME"),
				Endpoint: ptr.Ptr("CHANGE-ME"),
			}
			loggingConfig.Grafana = &config.GrafanaConfig{
				BaseUrl:      ptr.Ptr("CHANGE-ME"),
				DashboardUrl: ptr.Ptr("CHANGE-ME"),
				BearerToken:  ptr.Ptr("CHANGE-ME"),
			}

			lw, err := logstream.NewLogStream(t, &loggingConfig)
			require.NoError(t, err)
			for _, c := range d.containers {
				err = lw.ConnectContainer(ctx, c, "")
				require.NoError(t, err)
			}
			time.Sleep(5 * time.Second)

			// we don't want them to keep logging after we have stopped log stream by flushing logs
			for _, c := range d.containers {
				err = lw.DisconnectContainer(c)
				require.NoError(t, err)
			}

			err = lw.FlushLogsToTargets()
			require.NoError(t, err)
			lw.PrintLogTargetsLocations()
		})
	}
}
