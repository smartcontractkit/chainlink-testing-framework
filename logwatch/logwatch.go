package logwatch

import (
	"context"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/wasp"
	"github.com/testcontainers/testcontainers-go"
	"regexp"
	"strings"
	"testing"
	"time"
)

// LogNotification notification about log line match for some container
type LogNotification struct {
	Container string
	Log       string
}

// LogWatch is a test helper struct to monitor docker container logs for some patterns
// and push their logs into Loki for further analysis
type LogWatch struct {
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
	return &LogWatch{
		log:        utils.GetTestComponentLogger(t, "LogWatch"),
		loki:       loki,
		patterns:   patterns,
		notifyTest: make(chan *LogNotification, 10000),
		consumers:  make(map[string]*ContainerLogConsumer, 0),
	}, nil
}

// Listen listen for the next notification
func (m *LogWatch) Listen() *LogNotification {
	return <-m.notifyTest
}

// ConnectContainer connects consumer to selected container and starts testcontainers.LogProducer
func (m *LogWatch) ConnectContainer(ctx context.Context, container testcontainers.Container, pushToLoki bool) (string, error) {
	name, err := container.Name(ctx)
	if err != nil {
		return "", err
	}
	name = strings.Replace(name, "/", "", 1)
	m.log.Info().Str("name", name).Msg("Connecting container logs")
	cons := NewContainerLogConsumer(m, name, pushToLoki)
	m.consumers[name] = cons
	m.containers = append(m.containers, container)
	container.FollowOutput(cons)
	err = container.StartLogProducer(context.Background())
	if err != nil {
		return "", err
	}
	return name, nil
}

// Shutdown disconnects all containers, stops notifications
func (m *LogWatch) Shutdown() {
	for _, c := range m.containers {
		if c.IsRunning() {
			_ = c.StopLogProducer()
		}
	}
	for cn := range m.consumers {
		m.consumers[cn] = nil
	}
	close(m.notifyTest)
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
	pushToLoki bool
	lw         *LogWatch
	Messages   []string
}

// NewContainerLogConsumer creates new log consumer for a container that
// - signal if log line matches the pattern
// - push all lines to Loki if enabled
func NewContainerLogConsumer(lw *LogWatch, containerName string, pushToLoki bool) *ContainerLogConsumer {
	return &ContainerLogConsumer{
		name:       containerName,
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
		g.lw.notifyTest <- &LogNotification{Container: g.name, Log: string(l.Content)}
	}
	// we can notify more than one time if it matches, but we push only once
	if g.pushToLoki && g.lw.loki != nil {
		_ = g.lw.loki.Handle(model.LabelSet{
			"type":      "log_watch",
			"container": model.LabelValue(g.name),
		}, time.Now(), string(l.Content))
	}
}

// FindMatch check multiple regex patterns for the same string
// can be checked with one regex, made for readability of API
func (g *ContainerLogConsumer) FindMatch(l testcontainers.Log) int {
	matchesPerPattern := 0
	for _, filterRegex := range g.lw.patterns[g.name] {
		if filterRegex.Match(l.Content) {
			g.lw.log.Info().
				Str("Regex", filterRegex.String()).
				Str("String", string(l.Content)).
				Msg("Match found")
			matchesPerPattern++
		}
	}
	return matchesPerPattern
}
