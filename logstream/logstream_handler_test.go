package logstream_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

type MockedLogHandler struct {
	logs        []logstream.LogContent
	Target      logstream.LogTarget
	executionId string
}

func (m *MockedLogHandler) Handle(consumer *logstream.ContainerLogConsumer, content logstream.LogContent) error {
	m.logs = append(m.logs, content)
	return nil
}

func (m *MockedLogHandler) GetLogLocation(consumers map[string]*logstream.ContainerLogConsumer) (string, error) {
	return "", nil
}

func (m *MockedLogHandler) GetTarget() logstream.LogTarget {
	return m.Target
}

func (m *MockedLogHandler) SetRunId(executionId string) {
	m.executionId = executionId
}

func (m *MockedLogHandler) GetRunId() string {
	return m.executionId
}

func (m *MockedLogHandler) Init(_ *logstream.ContainerLogConsumer) error {
	return nil
}

func (m *MockedLogHandler) Teardown() error {
	return nil
}

func TestMultipleMockedLoggingTargets(t *testing.T) {
	ctx := context.Background()
	testData := testData{repeat: 10, perSecond: 0.01, streams: []string{"A\nB\nC\nD"}}
	d, err := NewDeployment(ctx, testData)
	// nolint
	defer d.Shutdown(ctx)
	require.NoError(t, err)
	mockedFileHandler := &MockedLogHandler{Target: logstream.File}
	mockedLokiHandler := &MockedLogHandler{Target: logstream.Loki}

	loggingConfig := config.LoggingConfig{}
	loggingConfig.LogStream = &config.LogStreamConfig{
		LogTargets:            []string{"loki", "file"},
		LogProducerTimeout:    &blockchain.StrDuration{Duration: 10 * time.Second},
		LogProducerRetryLimit: ptr.Ptr(uint(10)),
	}

	lw, err := logstream.NewLogStream(
		t,
		&loggingConfig,
		logstream.WithCustomLogHandler(logstream.File, mockedFileHandler),
		logstream.WithCustomLogHandler(logstream.Loki, mockedLokiHandler),
	)
	require.NoError(t, err, "failed to create logstream")
	err = d.ConnectLogs(lw)
	require.NoError(t, err, "failed to connect logs")

	time.Sleep(2 * time.Second)
	err = lw.FlushLogsToTargets()
	require.NoError(t, err, "failed to flush logs to targets")

	assertMockedHandlerHasLogs(t, mockedFileHandler)
	assertMockedHandlerHasLogs(t, mockedLokiHandler)

	err = lw.Shutdown(ctx)
	require.NoError(t, err, "failed to shutdown logstream")
}

func TestOneMockedLoggingTarget(t *testing.T) {
	ctx := context.Background()
	testData := testData{repeat: 10, perSecond: 0.01, streams: []string{"A\nB\nC\nD"}}
	d, err := NewDeployment(ctx, testData)
	// nolint
	defer d.Shutdown(ctx)
	require.NoError(t, err)
	mockedLokiHandler := &MockedLogHandler{Target: logstream.Loki}

	loggingConfig := config.LoggingConfig{}
	loggingConfig.LogStream = &config.LogStreamConfig{
		LogTargets:            []string{"loki"},
		LogProducerTimeout:    &blockchain.StrDuration{Duration: 10 * time.Second},
		LogProducerRetryLimit: ptr.Ptr(uint(10)),
	}

	lw, err := logstream.NewLogStream(
		t,
		&loggingConfig,
		logstream.WithCustomLogHandler(logstream.Loki, mockedLokiHandler),
	)
	require.NoError(t, err, "failed to create logstream")
	err = d.ConnectLogs(lw)
	require.NoError(t, err, "failed to connect logs")

	time.Sleep(2 * time.Second)
	err = lw.FlushLogsToTargets()
	require.NoError(t, err, "failed to flush logs to targets")

	assertMockedHandlerHasLogs(t, mockedLokiHandler)

	err = lw.Shutdown(ctx)
	require.NoError(t, err, "failed to shutdown logstream")
}

func assertMockedHandlerHasLogs(t *testing.T, handler *MockedLogHandler) {
	matches := make(map[string]int)
	matches["A"] = 0
	matches["B"] = 0
	matches["C"] = 0

	for _, log := range handler.logs {
		require.Equal(t, log.TestName, t.Name())
		require.Equal(t, log.ContainerName, "container-0")

		if bytes.Equal(log.Content, A) {
			matches["A"]++
		}

		if bytes.Equal(log.Content, B) {
			matches["B"]++
		}

		if bytes.Equal(log.Content, C) {
			matches["C"]++
		}
	}

	require.Greater(t, matches["A"], 0, "A should be present at least once in handler for %s", handler.Target)
	require.Greater(t, matches["B"], 0, "B should be matched at least once in handler for %s", handler.Target)
	require.Greater(t, matches["C"], 0, "C should be matched at least once in handler for %s", handler.Target)
}
