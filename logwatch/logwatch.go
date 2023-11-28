package logwatch

import (
	"context"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/retries"
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
	testName                     string
	log                          zerolog.Logger
	loki                         *wasp.LokiClient
	patterns                     map[string][]*regexp.Regexp
	notifyTest                   chan *LogNotification
	containers                   []testcontainers.Container
	consumers                    map[string]*ContainerLogConsumer
	logTargetHandlers            map[LogTarget]HandleLogTarget
	logListeningDone             chan struct{}
	logProducerTimeout           time.Duration
	logProducerTimeoutRetryLimit int // -1 for infinite retries
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
		testName:                     testName,
		log:                          l,
		patterns:                     patterns,
		notifyTest:                   make(chan *LogNotification, 10000),
		consumers:                    make(map[string]*ContainerLogConsumer, 0),
		logTargetHandlers:            getDefaultLogHandlers(),
		logListeningDone:             make(chan struct{}, 1),
		logProducerTimeout:           time.Duration(10 * time.Second),
		logProducerTimeoutRetryLimit: 10,
	}

	for _, option := range options {
		option(logWatch)
	}

	if err := logWatch.validateLogTargets(); err != nil {
		return nil, err
	}

	return logWatch, nil
}

func (m *LogWatch) validateLogTargets() error {
	envLogTargets, err := getLogTargetsFromEnv()
	if err != nil {
		return err
	}

	// check if all requested log targets are supported
	for _, wantedTarget := range envLogTargets {
		found := false
		for knownTarget := range m.logTargetHandlers {
			if knownTarget == wantedTarget {
				found = true
				break
			}
		}

		if !found {
			return errors.Errorf("no handler found for log target: %s", wantedTarget)
		}
	}

	// deactivate known log targets that are not enabled
	for knownTarget := range m.logTargetHandlers {
		wanted := false
		for _, wantedTarget := range envLogTargets {
			if knownTarget == wantedTarget {
				wanted = true
				break
			}
		}
		if !wanted {
			m.log.Debug().Str("log target", string(knownTarget)).Msg("Log target disabled")
			delete(m.logTargetHandlers, knownTarget)
		}
	}

	if len(m.logTargetHandlers) == 0 {
		m.log.Warn().Msg("No log targets enabled. LogWatch will not do anything")
	}

	return nil
}

func WithCustomLogHandler(logTarget LogTarget, handler HandleLogTarget) Option {
	return func(lw *LogWatch) {
		lw.logTargetHandlers[logTarget] = handler
	}
}

func WithLogProducerTimeout(timeout time.Duration) Option {
	return func(lw *LogWatch) {
		lw.logProducerTimeout = timeout
	}
}

func WithLogProducerTimeoutRetryLimit(retryLimit int) Option {
	return func(lw *LogWatch) {
		lw.logProducerTimeoutRetryLimit = retryLimit
	}
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
func (m *LogWatch) ConnectContainer(ctx context.Context, container testcontainers.Container, prefix string) error {
	name, err := container.Name(ctx)
	if err != nil {
		return err
	}
	name = strings.Replace(name, "/", "", 1)
	prefix = strings.Replace(prefix, "/", "", 1)

	enabledLogTargets := make([]LogTarget, 0)
	for logTarget := range m.logTargetHandlers {
		enabledLogTargets = append(enabledLogTargets, logTarget)
	}

	var cons *ContainerLogConsumer
	if prefix != "" {
		cons = newContainerLogConsumer(m, name, prefix, enabledLogTargets...)
	} else {
		cons = newContainerLogConsumer(m, name, name, enabledLogTargets...)
	}

	m.log.Info().
		Str("Prefix", prefix).
		Str("Name", name).
		Str("Timeout", m.logProducerTimeout.String()).
		Msg("Connecting container logs")
	m.consumers[name] = cons
	m.containers = append(m.containers, container)
	container.FollowOutput(cons)

	err = container.StartLogProducer(ctx, m.logProducerTimeout)

	go func(done chan struct{}, timeout time.Duration, retryLimit int) {
		defer m.log.Info().Str("Container name", name).Msg("Log listener stopped")
		currentAttempt := 0

		var shouldRetry = func() bool {
			if retryLimit == -1 {
				return true
			}

			if currentAttempt < retryLimit {
				currentAttempt++
				return true
			}

			return false
		}

		for {
			select {
			case err = <-container.GetLogProducerErrorChannel():
				if err != nil {
					m.log.Error().
						Str("Name", name).
						Err(err).
						Msg("Log producer errored")
					if shouldRetry() {
						backoff := retries.Fibonacci(currentAttempt)
						timeout = timeout + time.Duration(backoff)*time.Millisecond
						m.log.Info().
							Str("Prefix", prefix).
							Str("Name", name).
							Str("Timeout", timeout.String()).
							Msgf("Retrying connection and listening to container logs. Attempt %d/%d", currentAttempt, retryLimit)
						err = container.StartLogProducer(ctx, timeout)
						if err != nil {
							m.log.Error().Err(err).Msg("Log producer was already running. This should never happen. Exiting")
							return
						}
						m.log.Info().
							Str("Name", name).
							Msg("Started new log producer")
					} else {
						m.log.Error().
							Err(err).
							Str("Name", name).
							Msg("Used all attempts to listen to container logs. Won't try again")
						return
					}
				}
			case <-done:
				return
			}
		}
	}(m.logListeningDone, m.logProducerTimeout, m.logProducerTimeoutRetryLimit)

	return err
}

// Shutdown disconnects all containers, stops notifications
func (m *LogWatch) Shutdown() error {
	defer close(m.logListeningDone)
	var err error
	for _, c := range m.containers {
		singleErr := m.DisconnectContainer(c)
		if singleErr != nil {
			ctx := context.Background()
			name, _ := c.Name(ctx)
			m.log.Error().
				Err(err).
				Str("Name", name).
				Msg("Failed to disconnect container")

			err = errors.Wrap(singleErr, "failed to disconnect container")
		}
	}

	if m.loki != nil {
		m.loki.Stop()
	}

	m.logListeningDone <- struct{}{}

	return err
}

type LogWriter = func(testName string, name string, location interface{}) error

func (m *LogWatch) PrintLogTargetsLocations() {
	m.SaveLogTargetsLocations(func(testName string, name string, location interface{}) error {
		m.log.Info().Str("Test", testName).Str("Handler", name).Interface("Location", location).Msg("Log location")
		return nil
	})
}

func (m *LogWatch) SaveLogTargetsLocations(writer LogWriter) {
	for _, handler := range m.logTargetHandlers {
		name := string(handler.GetTarget())
		location, err := handler.GetLogLocation(m.consumers)
		if err != nil {
			m.log.Error().Str("Handler", name).Err(err).Msg("Failed to get log location")
			continue
		}

		if err := writer(m.testName, name, location); err != nil {
			m.log.Error().Str("Handler", name).Err(err).Msg("Failed to write log location")
		}
	}
}

// DisconnectContainer disconnects the particular container
func (m *LogWatch) DisconnectContainer(container testcontainers.Container) error {
	if container.IsRunning() {
		m.log.Info().Str("container", container.GetContainerID()).Msg("Disconnecting container")
		return container.StopLogProducer()
	}

	return nil
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
			g.lw.log.Warn().Str("log target", string(logTarget)).Msg("No handler found for log target")
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
		if lt == logTarget {
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
