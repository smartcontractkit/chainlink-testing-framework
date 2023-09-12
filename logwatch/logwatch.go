package logwatch

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/wasp"
	"github.com/testcontainers/testcontainers-go"
)

// LogNotification notification about log line match for some container
type LogNotification struct {
	Container string
	Prefix    string
	Log       string
}

// LogWatch is a test helper struct to monitor docker container logs for some patterns
// and push their logs into Loki for further analysis
type LogWatch struct {
	t          *testing.T
	log        zerolog.Logger
	loki       *wasp.LokiClient
	patterns   map[string][]*regexp.Regexp
	notifyTest chan *LogNotification
	containers []testcontainers.Container
	consumers  map[string]*ContainerLogConsumer
}

// NewLogWatch creates a new LogWatch instance, with a Loki client
func NewLogWatch(t *testing.T, patterns map[string][]*regexp.Regexp) (*LogWatch, error) {
	loki, err := wasp.NewLokiClient(wasp.NewEnvLokiConfig())
	if err != nil {
		return nil, err
	}
	l := logging.GetLogger(t, "LOGWATCH_LOG_LEVEL").With().Str("Component", "LogWatch").Logger()
	return &LogWatch{
		t:          t,
		log:        l,
		loki:       loki,
		patterns:   patterns,
		notifyTest: make(chan *LogNotification, 10000),
		consumers:  make(map[string]*ContainerLogConsumer, 0),
	}, nil
}

// Listen listen for the next notification
func (m *LogWatch) Listen() *LogNotification {
	msg := <-m.notifyTest
	m.log.Warn().
		Str("Container", msg.Container).
		Str("Line", msg.Log).
		Msg("Received notification from container")
	return msg
}

// OnMatch calling your testing hook on first match
func (m *LogWatch) OnMatch(f func(ln *LogNotification)) {
	go func() {
		for {
			msg := <-m.notifyTest
			m.log.Warn().
				Str("Container", msg.Container).
				Str("Line", msg.Log).
				Msg("Received notification from container")
			f(msg)
		}
	}()
}

// ConnectContainer connects consumer to selected container and starts testcontainers.LogProducer
func (m *LogWatch) ConnectContainer(ctx context.Context, container testcontainers.Container, prefix string, pushToLoki bool) error {
	name, err := container.Name(ctx)
	if err != nil {
		return err
	}
	name = strings.Replace(name, "/", "", 1)
	prefix = strings.Replace(prefix, "/", "", 1)
	var cons *ContainerLogConsumer
	if prefix != "" {
		cons = newContainerLogConsumer(m, name, prefix, pushToLoki)
	} else {
		cons = newContainerLogConsumer(m, name, name, pushToLoki)
	}
	m.log.Info().
		Str("Prefix", prefix).
		Str("Name", name).
		Msg("Connecting container logs")
	m.consumers[name] = cons
	m.containers = append(m.containers, container)
	container.FollowOutput(cons)
	return container.StartLogProducer(ctx)
}

// Shutdown disconnects all containers, stops notifications
func (m *LogWatch) Shutdown() {
	m.loki.Stop()
}

// DisconnectContainer disconnects the particular container
func (m *LogWatch) DisconnectContainer(container testcontainers.Container) {
	if container.IsRunning() {
		_ = container.StopLogProducer()
	}
}

// ContainerLogs return all logs for the particular container
func (m *LogWatch) ContainerLogs(name string) []string {
	return m.consumers[name].Messages
}

// AllLogs returns all logs for all containers
func (m *LogWatch) AllLogs() []string {
	logs := make([]string, 0)
	for _, l := range m.consumers {
		logs = append(logs, l.Messages...)
	}
	return logs
}

// PrintAll prints all logs for all containers connected
func (m *LogWatch) PrintAll() {
	for cname, c := range m.consumers {
		for _, msg := range c.Messages {
			m.log.Info().
				Str("Container", cname).
				Str("Msg", msg).
				Send()
		}
	}
}

// ContainerLogConsumer is a container log lines consumer
type ContainerLogConsumer struct {
	name       string
	prefix     string
	pushToLoki bool
	lw         *LogWatch
	Messages   []string
}

// newContainerLogConsumer creates new log consumer for a container that
// - signal if log line matches the pattern
// - push all lines to Loki if enabled
func newContainerLogConsumer(lw *LogWatch, containerName string, prefix string, pushToLoki bool) *ContainerLogConsumer {
	return &ContainerLogConsumer{
		name:       containerName,
		prefix:     prefix,
		pushToLoki: pushToLoki,
		lw:         lw,
		Messages:   make([]string, 0),
	}
}

// Accept accepts the log message from particular container
func (g *ContainerLogConsumer) Accept(l testcontainers.Log) {
	g.Messages = append(g.Messages, string(l.Content))
	matches := g.FindMatch(l)
	for i := 0; i < matches; i++ {
		g.lw.notifyTest <- &LogNotification{Container: g.name, Prefix: g.prefix, Log: string(l.Content)}
	}
	var testName string
	if g.lw.t == nil {
		testName = "no_test"
	} else {
		testName = g.lw.t.Name()
	}
	// we can notify more than one time if it matches, but we push only once
	if g.pushToLoki && g.lw.loki != nil {
		_ = g.lw.loki.Handle(model.LabelSet{
			"type":      "log_watch",
			"test":      model.LabelValue(testName),
			"container": model.LabelValue(g.name),
		}, time.Now(), string(l.Content))
	}
}

// FindMatch check multiple regex patterns for the same string
// can be checked with one regex, made for readability of user-facing API
func (g *ContainerLogConsumer) FindMatch(l testcontainers.Log) int {
	matchesPerPattern := 0
	if g.prefix == "" {
		g.prefix = g.name
	}
	for _, filterRegex := range g.lw.patterns[g.name] {
		if filterRegex.Match(l.Content) {
			g.lw.log.Info().
				Str("Container", g.name).
				Str("Regex", filterRegex.String()).
				Str("String", string(l.Content)).
				Msg("Match found")
			matchesPerPattern++
		}
	}
	return matchesPerPattern
}
