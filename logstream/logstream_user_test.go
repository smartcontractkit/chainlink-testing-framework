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

	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
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
	lw, err := logstream.NewLogStream(
		t,
		nil,
		logstream.WithLogTarget(logstream.File),
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

func TestMultipleMockedLoggingTargets(t *testing.T) {
	ctx := context.Background()
	testData := testData{repeat: 10, perSecond: 0.01, streams: []string{"A\nB\nC\nD"}}
	d, err := NewDeployment(ctx, testData)
	// nolint
	defer d.Shutdown(ctx)
	require.NoError(t, err)
	mockedFileHandler := &MockedLogHandler{Target: logstream.File}
	mockedLokiHanlder := &MockedLogHandler{Target: logstream.Loki}
	lw, err := logstream.NewLogStream(
		t,
		nil,
		logstream.WithCustomLogHandler(logstream.File, mockedFileHandler),
		logstream.WithCustomLogHandler(logstream.Loki, mockedLokiHanlder),
		logstream.WithLogTarget(logstream.Loki),
		logstream.WithLogTarget(logstream.File),
	)
	require.NoError(t, err, "failed to create logstream")
	err = d.ConnectLogs(lw)
	require.NoError(t, err, "failed to connect logs")

	time.Sleep(2 * time.Second)
	err = lw.FlushLogsToTargets()
	require.NoError(t, err, "failed to flush logs to targets")

	assertMockedHandlerHasLogs(t, mockedFileHandler)
	assertMockedHandlerHasLogs(t, mockedLokiHanlder)

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
	mockedLokiHanlder := &MockedLogHandler{Target: logstream.Loki}
	lw, err := logstream.NewLogStream(
		t,
		nil,
		logstream.WithCustomLogHandler(logstream.Loki, mockedLokiHanlder),
		logstream.WithLogTarget(logstream.Loki),
	)
	require.NoError(t, err, "failed to create logstream")
	err = d.ConnectLogs(lw)
	require.NoError(t, err, "failed to connect logs")

	time.Sleep(2 * time.Second)
	err = lw.FlushLogsToTargets()
	require.NoError(t, err, "failed to flush logs to targets")

	assertMockedHandlerHasLogs(t, mockedLokiHanlder)

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
