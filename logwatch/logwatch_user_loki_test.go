package logwatch_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
)

/* These tests are for testing Loki format, they rarely change so you can run them manually */

func TestExampleLokiStreaming(t *testing.T) {
	t.Skip("uncomment and run manually")
	os.Setenv("LOGWATCH_LOG_TARGETS", "loki")
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
			d, err := NewDeployment(tc)
			// nolint
			defer d.Shutdown()
			require.NoError(t, err)
			lw, err := logwatch.NewLogWatch(t, nil)
			require.NoError(t, err)
			err = d.ConnectLogs(lw)
			require.NoError(t, err)
			time.Sleep(5 * time.Second)
		})
	}
}
