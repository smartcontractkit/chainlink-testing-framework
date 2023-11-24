package test_env

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

type LogRegexStrategy struct {
	timeout      *time.Duration
	Pattern      *regexp.Regexp
	Occurrence   int
	PollInterval time.Duration
}

func NewLogRegexStrategy(pattern *regexp.Regexp) *LogRegexStrategy {
	return &LogRegexStrategy{
		Pattern:      pattern,
		Occurrence:   1,
		PollInterval: defaultPollInterval(),
	}
}

func (ws *LogRegexStrategy) WithStartupTimeout(timeout time.Duration) *LogRegexStrategy {
	ws.timeout = &timeout
	return ws
}

// WithPollInterval can be used to override the default polling interval of 100 milliseconds
func (ws *LogRegexStrategy) WithPollInterval(pollInterval time.Duration) *LogRegexStrategy {
	ws.PollInterval = pollInterval
	return ws
}

func (ws *LogRegexStrategy) WithOccurrence(o int) *LogRegexStrategy {
	// the number of occurrence needs to be positive
	if o <= 0 {
		o = 1
	}
	ws.Occurrence = o
	return ws
}

// WaitUntilReady implements Strategy.WaitUntilReady
func (ws *LogRegexStrategy) WaitUntilReady(ctx context.Context, target tcwait.StrategyTarget) (err error) {
	timeout := defaultStartupTimeout()
	if ws.timeout != nil {
		timeout = *ws.timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	length := 0

LOOP:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			state, err := target.State(ctx)
			if err != nil {
				return err
			}
			if !state.Running {
				return fmt.Errorf("container is not running %s", state.Status)
			}

			reader, err := target.Logs(ctx)
			if err != nil {
				time.Sleep(ws.PollInterval)
				continue
			}

			b, err := io.ReadAll(reader)
			if err != nil {
				time.Sleep(ws.PollInterval)
				continue
			}

			logs := string(b)
			if length == len(logs) && err != nil {
				return err
			} else if len(ws.Pattern.FindAllString(logs, -1)) >= ws.Occurrence {
				break LOOP
			} else {
				length = len(logs)
				time.Sleep(ws.PollInterval)
				continue
			}
		}
	}

	return nil
}

func defaultStartupTimeout() time.Duration {
	return 60 * time.Second
}

func defaultPollInterval() time.Duration {
	return 100 * time.Millisecond
}
