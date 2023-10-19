package logwatch

import (
	"context"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

const NO_TEST = "no_test"

// LogNotification notification about log line match for some container
type LogNotification struct {
	Container string
	Prefix    string
	Log       string
}

// LogWatch is a test helper struct to monitor docker container logs for some patterns
// and push their logs into Loki for further analysis
type LogWatch struct {
	testName          string
	log               zerolog.Logger
	loki              *wasp.LokiClient
	patterns          map[string][]*regexp.Regexp
	notifyTest        chan *LogNotification
	containers        []testcontainers.Container
	consumers         map[string]*ContainerLogConsumer
	logTargetHandlers map[LogTarget]HandleLogTarget
}

type LogContent struct {
	TestName      string
	ContainerName string
	Content       []byte
}

type Option func(*LogWatch)

// NewLogWatch creates a new LogWatch instance, with Loki client only if Loki log target is enabled (lazy init)
func NewLogWatch(t *testing.T, patterns map[string][]*regexp.Regexp, options ...Option) (*LogWatch, error) {
	l := logging.GetLogger(nil, "LOGWATCH_LOG_LEVEL").With().Str("Component", "LogWatch").Logger()
	var testName string
	if t == nil {
		testName = NO_TEST
	} else {
		testName = t.Name()
	}

	logWatch := &LogWatch{
		testName:          testName,
		log:               l,
		patterns:          patterns,
		notifyTest:        make(chan *LogNotification, 10000),
		consumers:         make(map[string]*ContainerLogConsumer, 0),
		logTargetHandlers: getDefaultLogHandlers(),
	}

	for _, option := range options {
		option(logWatch)
	}

	if err := logWatch.validateLogTargets(); err != nil {
		return nil, err
	}

	return logWatch, nil
}

func (l *LogWatch) validateLogTargets() error {
	envLogTargets, err := getLogTargetsFromEnv()
	if err != nil {
		return err
	}

	// check if all requested log targets are supported
	for _, wantedTarget := range envLogTargets {
		found := false
		for knownTargets := range l.logTargetHandlers {
			if knownTargets == wantedTarget {
				found = true
				break
			}
		}

		if !found {
			return errors.Errorf("no handler found for log target: %d", wantedTarget)
		}
	}

	// deactivate known log targets that are not enabled
	for knownTarget := range l.logTargetHandlers {
		wanted := false
		for _, wantedTarget := range envLogTargets {
			if knownTarget == wantedTarget {
				wanted = true
				break
			}
		}
		if !wanted {
			l.log.Debug().Int("handler id", int(knownTarget)).Msg("Log target disabled")
			delete(l.logTargetHandlers, knownTarget)
		}
	}

	if len(l.logTargetHandlers) == 0 {
		l.log.Warn().Msg("No log targets enabled. LogWatch will not do anything")
	}

	return nil
}

func WithCustomLogHandler(logTarget LogTarget, handler HandleLogTarget) Option {
	return func(lw *LogWatch) {
		lw.logTargetHandlers[logTarget] = handler
	}
}

// Listen listen for the next notification
func (l *LogWatch) Listen() *LogNotification {
	msg := <-l.notifyTest
	l.log.Warn().
		Str("Container", msg.Container).
		Str("Line", msg.Log).
		Msg("Received notification from container")
	return msg
}

// OnMatch calling your testing hook on first match
func (l *LogWatch) OnMatch(f func(ln *LogNotification)) {
	go func() {
		for {
			msg := <-l.notifyTest
			l.log.Warn().
				Str("Container", msg.Container).
				Str("Line", msg.Log).
				Msg("Received notification from container")
			f(msg)
		}
	}()
}

// ConnectContainer connects consumer to selected container and starts testcontainers.LogProducer
func (l *LogWatch) ConnectContainer(ctx context.Context, container testcontainers.Container, prefix string) error {
	name, err := container.Name(ctx)
	if err != nil {
		return err
	}
	name = strings.Replace(name, "/", "", 1)
	prefix = strings.Replace(prefix, "/", "", 1)

	enabledLogTargets := make([]LogTarget, 0)
	for logTarget := range l.logTargetHandlers {
		enabledLogTargets = append(enabledLogTargets, logTarget)
	}

	var cons *ContainerLogConsumer
	if prefix != "" {
		cons = newContainerLogConsumer(l, name, prefix, enabledLogTargets...)
	} else {
		cons = newContainerLogConsumer(l, name, name, enabledLogTargets...)
	}

	l.log.Info().
		Str("Prefix", prefix).
		Str("Name", name).
		Msg("Connecting container logs")
	l.consumers[name] = cons
	l.containers = append(l.containers, container)
	container.FollowOutput(cons)
	return container.StartLogProducer(ctx)
}

// Shutdown disconnects all containers, stops notifications
func (m *LogWatch) Shutdown() {
	for _, c := range m.containers {
		m.DisconnectContainer(c)
	}

	if m.loki != nil {
		m.loki.Stop()
	}
}

func (m *LogWatch) PrintLogTargetsLocations() {
	for _, handler := range m.logTargetHandlers {
		handler.PrintLogLocation(m)
	}
}

// DisconnectContainer disconnects the particular container
func (m *LogWatch) DisconnectContainer(container testcontainers.Container) {
	if container.IsRunning() {
		m.log.Info().Str("container", container.GetContainerID()).Msg("Disconnecting container")
		_ = container.StopLogProducer()
	}
}

// ContainerLogs return all logs for the particular container
func (m *LogWatch) ContainerLogs(name string) []string {
	if _, ok := m.consumers[name]; !ok {
		return []string{}
	}

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
	logTargets []LogTarget
	lw         *LogWatch
	Messages   []string
}

// newContainerLogConsumer creates new log consumer for a container that
// - signal if log line matches the pattern
// - push all lines to configured log targets
func newContainerLogConsumer(lw *LogWatch, containerName string, prefix string, logTargets ...LogTarget) *ContainerLogConsumer {
	return &ContainerLogConsumer{
		name:       containerName,
		prefix:     prefix,
		logTargets: logTargets,
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

	content := LogContent{
		TestName:      g.lw.testName,
		ContainerName: g.name,
		Content:       l.Content,
	}

	for _, logTarget := range g.logTargets {
		if handler, ok := g.lw.logTargetHandlers[logTarget]; ok {
			if err := handler.Handle(g, content); err != nil {
				g.lw.log.Error().Err(err).Msg("Failed to handle log target")
			}
		} else {
			g.lw.log.Warn().Int("handler id", int(logTarget)).Msg("No handler found for log target")
		}
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

func (g *ContainerLogConsumer) hasLogTarget(logTarget LogTarget) bool {
	for _, lt := range g.logTargets {
		if lt&logTarget != 0 {
			return true
		}
	}

	return false
}

func getLogTargetsFromEnv() ([]LogTarget, error) {
	envLogTargetsValue := os.Getenv("LOGWATCH_LOG_TARGETS")
	if envLogTargetsValue != "" {
		envLogTargets := make([]LogTarget, 0)
		for _, target := range strings.Split(envLogTargetsValue, ",") {
			switch strings.TrimSpace(strings.ToLower(target)) {
			case "loki":
				envLogTargets = append(envLogTargets, Loki)
			case "file":
				envLogTargets = append(envLogTargets, File)
			default:
				return []LogTarget{}, errors.Errorf("unknown log target: %s", target)
			}
		}

		return envLogTargets, nil
	}

	return []LogTarget{}, nil
}
