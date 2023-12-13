package logstream_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
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
				LogTargets: []string{"loki"},
			}

			lw, err := logstream.NewLogStream(t, &loggingConfig)
			require.NoError(t, err)
			err = d.ConnectLogs(lw)
			require.NoError(t, err)
			time.Sleep(5 * time.Second)
		})
	}
}
