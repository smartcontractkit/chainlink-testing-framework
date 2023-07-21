package logwatch_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
)

/* These tests are for testing Loki format, they rarely change so you can run them manually */

func TestExampleLokiStreaming(t *testing.T) {
	t.Skip("uncomment and run manually")
	t.Run("stream all container logs to Loki, subtest 1", func(t *testing.T) {
		testData := testData{repeat: 1, perSecond: 0.01, streams: []string{"A\nB\nC\nD", "E\nF\nG\nH"}}
		d, err := NewDeployment(testData)
		// nolint
		defer d.Shutdown()
		require.NoError(t, err)
		lw, err := logwatch.NewLogWatch(t, nil)
		require.NoError(t, err)
		err = d.ConnectLogs(lw, true)
		require.NoError(t, err)
		time.Sleep(5 * time.Second)
	})
	t.Run("stream all container logs to Loki, subtest 2", func(t *testing.T) {
		testData := testData{repeat: 1, perSecond: 0.01, streams: []string{"1\n2\n3\n4", "5\n6\n7\n8"}}
		d, err := NewDeployment(testData)
		// nolint
		defer d.Shutdown()
		require.NoError(t, err)
		lw, err := logwatch.NewLogWatch(t, nil)
		require.NoError(t, err)
		err = d.ConnectLogs(lw, true)
		require.NoError(t, err)
		time.Sleep(5 * time.Second)
	})
	t.Run("nobody expects the spammish repetition", func(t *testing.T) {
		testData := testData{
			repeat:    1000,
			perSecond: 0.0001,
			streams:   []string{"nobody expects the spanish inquisition"}}
		d, err := NewDeployment(testData)
		// nolint
		defer d.Shutdown()
		require.NoError(t, err)
		lw, err := logwatch.NewLogWatch(t, nil)
		require.NoError(t, err)
		err = d.ConnectLogs(lw, true)
		require.NoError(t, err)
		time.Sleep(5 * time.Second)
	})
}
