package logstream_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

/* These tests are for user-facing API */

/* This data is for testing only, won't exist in real deployment */

type testData struct {
	name      string
	repeat    int
	perSecond float64
	streams   []string
}

/* That's your example deployment */

type MyDeployment struct {
	containers []testcontainers.Container
}

func NewDeployment(ctx context.Context, data testData) (*MyDeployment, error) {
	md := &MyDeployment{containers: make([]testcontainers.Container, 0)}
	for i, messages := range data.streams {
		c, err := startTestContainer(ctx, fmt.Sprintf("container-%d", i), messages, data.repeat, data.perSecond, false)
		if err != nil {
			return md, err
		}
		md.containers = append(md.containers, c)
	}
	return md, nil
}

func (m *MyDeployment) Shutdown(ctx context.Context) error {
	for _, c := range m.containers {
		if err := c.Terminate(ctx); err != nil {
			return err
		}
	}
	return nil
}

/* That's what you need to implement to have your logs send to your chosen targets */
func (m *MyDeployment) ConnectLogs(lw *logstream.LogStream) error {
	for _, c := range m.containers {
		if err := lw.ConnectContainer(context.Background(), c, ""); err != nil {
			return err
		}
	}
	return nil
}

/* That's how you use it */

var (
	A = []byte("A\n")
	B = []byte("B\n")
	C = []byte("C\n")
)

func TestFileLoggingTarget(t *testing.T) {
	ctx := context.Background()
	testData := testData{repeat: 10, perSecond: 0.01, streams: []string{"A\nB\nC\nD"}}
	d, err := NewDeployment(ctx, testData)
	// nolint
	defer d.Shutdown(ctx)
	require.NoError(t, err)

	loggingConfig := config.LoggingConfig{}
	loggingConfig.LogStream = &config.LogStreamConfig{
		LogTargets:            []string{"file"},
		LogProducerTimeout:    &blockchain.StrDuration{Duration: 10 * time.Second},
		LogProducerRetryLimit: ptr.Ptr(uint(10)),
	}

	lw, err := logstream.NewLogStream(
		t,
		&loggingConfig,
	)
	require.NoError(t, err, "failed to create logstream")
	err = d.ConnectLogs(lw)
	require.NoError(t, err, "failed to connect logs")

	time.Sleep(2 * time.Second)

	var logFileLocation string

	bufferWriter := func(_ string, _ string, location interface{}) error {
		logFileLocation = location.(string)
		return nil
	}

	err = lw.FlushLogsToTargets()
	require.NoError(t, err, "failed to flush logs to targets")
	lw.SaveLogTargetsLocations(bufferWriter)

	content, err := os.ReadFile(logFileLocation + "/container-0.log")
	require.NoError(t, err, "failed to read log file")

	require.True(t, bytes.Contains(content, A), "A should be present in log file")
	require.True(t, bytes.Contains(content, B), "B should be present in log file")
	require.True(t, bytes.Contains(content, C), "C should be present in log file")

	err = lw.Shutdown(ctx)
	require.NoError(t, err, "failed to shutdown logstream")
}
