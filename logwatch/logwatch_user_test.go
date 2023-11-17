package logwatch_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
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

/* That's what you need to implement to have your logs in Loki */

func (m *MyDeployment) ConnectLogs(ctx context.Context, lw *logwatch.LogWatch, pushToLoki bool) error {
	for _, c := range m.containers {
		if err := lw.ConnectContainer(ctx, c, "", pushToLoki); err != nil {
			return err
		}
	}
	return nil
}

/* That's how you use it */

func TestExampleUserInteraction(t *testing.T) {
	ctx := testcontext.Get(t)
	t.Run("sync API, block, receive one message", func(t *testing.T) {
		td := testData{repeat: 10, perSecond: 0.01, streams: []string{"A\nB\nC\nD"}}
		d, err := NewDeployment(ctx, td)
		// nolint
		defer d.Shutdown(ctx)
		require.NoError(t, err)
		lw, err := logwatch.NewLogWatch(
			t,
			map[string][]*regexp.Regexp{
				"container-0": {
					regexp.MustCompile("A"),
				},
			},
		)
		require.NoError(t, err)
		err = d.ConnectLogs(ctx, lw, false)
		require.NoError(t, err)
		match := lw.Listen()
		require.NotEmpty(t, match)
	})
	t.Run("async API, execute some logic on match", func(t *testing.T) {
		td := testData{repeat: 10, perSecond: 0.01, streams: []string{"A\nB\nC\nD\n", "E\nF\nG\nH\n"}}
		notifications := 0
		d, err := NewDeployment(ctx, td)
		// nolint
		defer d.Shutdown(ctx)
		require.NoError(t, err)
		lw, err := logwatch.NewLogWatch(
			t,
			map[string][]*regexp.Regexp{
				"container-0": {
					regexp.MustCompile("A"),
				},
				"container-1": {
					regexp.MustCompile("E"),
				},
			},
		)
		require.NoError(t, err)
		lw.OnMatch(func(ln *logwatch.LogNotification) { notifications++ })
		err = d.ConnectLogs(ctx, lw, false)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		require.Equal(t, td.repeat*len(td.streams), notifications)
	})
}
