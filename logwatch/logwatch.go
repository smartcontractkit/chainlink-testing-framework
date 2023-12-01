package logwatch

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/testsummary"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/retries"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/runid"
)

const NO_TEST = "no_test"

// LogNotification notification about log line match for some container
type LogNotification struct {
	Container string
	Prefix    string
	Log       string
}

type LogProducingContainer interface {
	Name(ctx context.Context) (string, error)
	FollowOutput(consumer testcontainers.LogConsumer)
	StartLogProducer(ctx context.Context, timeout time.Duration) error
	StopLogProducer() error
	GetLogProducerErrorChannel() <-chan error
	IsRunning() bool
	GetContainerID() string
}

// LogWatch is a test helper struct to monitor docker container logs for some patterns
// and push their logs into Loki for further analysis
type LogWatch struct {
	testName                     string
	log                          zerolog.Logger
	loki                         *wasp.LokiClient
	patterns                     map[string][]*regexp.Regexp
	notifyTest                   chan *LogNotification
	containers                   []LogProducingContainer
	consumers                    map[string]*ContainerLogConsumer
	logTargetHandlers            map[LogTarget]HandleLogTarget
	enabledLogTargets            []LogTarget
	logListeningDone             chan struct{}
	logProducerTimeout           time.Duration
	logProducerTimeoutRetryLimit int // -1 for infinite retries
	acceptMutex                  sync.Mutex
	runId                        string
}

type LogContent struct {
	TestName      string
	ContainerName string
	Content       []byte
	Time          time.Time
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

	envLogTargets, err := getLogTargetsFromEnv()
	if err != nil {
		return nil, err
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
		enabledLogTargets:            envLogTargets,
		runId:                        fmt.Sprintf("%s-%s", testName, runid.GetOrGenerateRunId()),
	}

	for _, option := range options {
		option(logWatch)
	}

	if err := logWatch.validateLogTargets(); err != nil {
		return nil, err
	}

	for _, handler := range logWatch.logTargetHandlers {
		handler.SetRunId(logWatch.runId)
	}

	l.Info().Str("Run_id", logWatch.runId).Msg("LogWatch initialized")

	return logWatch, nil
}

func (m *LogWatch) validateLogTargets() error {
	// check if all requested log targets are supported
	for _, wantedTarget := range m.enabledLogTargets {
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
		for _, wantedTarget := range m.enabledLogTargets {
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
		m.log.Warn().Msg("No log targets enabled. LogWatch will not persist any logs")
	}

	return nil
}

func WithCustomLogHandler(logTarget LogTarget, handler HandleLogTarget) Option {
	return func(lw *LogWatch) {
		lw.logTargetHandlers[logTarget] = handler
	}
}

func WithLogTarget(logTarget LogTarget) Option {
	return func(lw *LogWatch) {
		lw.enabledLogTargets = append(lw.enabledLogTargets, logTarget)
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
func (m *LogWatch) ConnectContainer(ctx context.Context, container LogProducingContainer, prefix string) error {
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
		cons, err = newContainerLogConsumer(m, name, prefix, enabledLogTargets...)
	} else {
		cons, err = newContainerLogConsumer(m, name, name, enabledLogTargets...)
	}

	if err != nil {
		return err
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
		defer m.log.Info().Str("Container name", name).Msg("Disconnected container logs")
		currentAttempt := 1

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
			case logErr := <-container.GetLogProducerErrorChannel():
				if logErr != nil {
					m.log.Error().
						Err(err).
						Str("Container name", name).
						Msg("Log producer errored")
					if shouldRetry() {
						backoff := retries.Fibonacci(currentAttempt)
						timeout = timeout + time.Duration(backoff)*time.Second
						m.log.Info().
							Str("Prefix", prefix).
							Str("Container name", name).
							Str("Timeout", timeout.String()).
							Msgf("Retrying connection and listening to container logs. Attempt %d/%d", currentAttempt, retryLimit)
						// when log producer starts again it will request all logs again, so we need to remove ones already saved by log watch to avoid duplicates
						// in the unlikely case that log producer fails to start we will copy the messages received so far, so that at least some logs are salvaged
						messagesCopy := append([]string{}, m.consumers[name].Messages...)
						m.consumers[name].Messages = make([]string, 0)
						m.log.Warn().Msgf("Consumer messages: %d", len(m.consumers[name].Messages))

						failedToStart := false
						for container.StartLogProducer(ctx, timeout) != nil {
							if !shouldRetry() {
								failedToStart = true
								break
							}
							m.log.Info().
								Str("Container name", name).
								Msg("Waiting for log producer to stop before restarting it")
							time.Sleep(1 * time.Second)
						}
						if failedToStart {
							m.log.Error().
								Err(err).
								Str("Container name", name).
								Msg("Previously running log producer couldn't be stopped. Used all retry attempts. Won't try again")
							m.consumers[name].Messages = messagesCopy
							return
						}
						m.log.Info().
							Str("Container name", name).
							Msg("Started new log producer")
					} else {
						m.log.Error().
							Err(err).
							Str("Container name", name).
							Msg("Used all attempts to listen to container logs. Won't try again")
						return
					}

					time.Sleep(500 * time.Millisecond)
				}
			case <-done:
				return
			}
		}
	}(m.logListeningDone, m.logProducerTimeout, m.logProducerTimeoutRetryLimit)

	return err
}

// Shutdown disconnects all containers, stops notifications
func (m *LogWatch) Shutdown(context context.Context) error {
	defer close(m.logListeningDone)
	var err error
	for _, c := range m.containers {
		singleErr := m.DisconnectContainer(c)
		if singleErr != nil {
			name, _ := c.Name(context)
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

func (m *LogWatch) SaveLogLocationInTestSummary() {
	m.SaveLogTargetsLocations(func(testName string, name string, location interface{}) error {
		return testsummary.AddEntry(testName, name, location)
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
func (m *LogWatch) DisconnectContainer(container LogProducingContainer) error {
	if container.IsRunning() {
		m.log.Info().Str("container", container.GetContainerID()).Msg("Disconnecting container")
		return container.StopLogProducer()
	}

	return nil
}

// ContainerLogs return all logs for the particular container
func (m *LogWatch) ContainerLogs(name string) []string {
	m.acceptMutex.Lock()
	defer m.acceptMutex.Unlock()
	if _, ok := m.consumers[name]; !ok {
		return []string{}
	}

	return m.consumers[name].Messages
}

// AllLogs returns all logs for all containers
func (m *LogWatch) AllLogs() []string {
	m.acceptMutex.Lock()
	defer m.acceptMutex.Unlock()
	logs := make([]string, 0)
	for _, l := range m.consumers {
		logs = append(logs, l.Messages...)
	}
	return logs
}

// PrintAll prints all logs for all containers connected
func (m *LogWatch) PrintAll() {
	m.acceptMutex.Lock()
	defer m.acceptMutex.Unlock()
	for cname, c := range m.consumers {
		for _, msg := range c.Messages {
			m.log.Info().
				Str("Container", cname).
				Str("Msg", msg).
				Send()
		}
	}
}

// FlushLogsToTargets flushes all logs for all consumers (containers) to their targets
func (m *LogWatch) FlushLogsToTargets() error {
	m.acceptMutex.Lock()
	defer m.acceptMutex.Unlock()

	m.log.Info().Msg("Flushing logs to targets")
	for _, consumer := range m.consumers {
		// nothing to do if no log targets are configured
		if len(consumer.logTargets) == 0 {
			continue
		}

		if consumer.tempFile == nil {
			return errors.Errorf("temp file is nil for container %s, this should never happen", consumer.name)
		}

		// do not accept any new logs
		consumer.isDone = true
		// this was done on purpose, so that when we are done flushing all logs we can close the temp file and handle abrupt termination too
		// nolint
		defer consumer.tempFile.Close()

		_, err := consumer.tempFile.Seek(0, 0)
		if err != nil {
			return err
		}

		decoder := gob.NewDecoder(consumer.tempFile)
		counter := 0

		//TODO handle in batches?
		for {
			var log LogContent
			decodeErr := decoder.Decode(&log)
			if decodeErr == nil {
				counter++
				for _, logTarget := range consumer.logTargets {
					if handler, ok := consumer.lw.logTargetHandlers[logTarget]; ok {
						if err := handler.Handle(consumer, log); err != nil {
							m.log.Error().
								Err(err).
								Str("Container", consumer.name).
								Str("log target", string(logTarget)).
								Msg("Failed to handle log target")
						}
					} else {
						m.log.Warn().
							Str("Container", consumer.name).
							Str("log target", string(logTarget)).
							Msg("No handler found for log target")
					}
				}
			} else if errors.Is(decodeErr, io.EOF) {
				m.log.Info().
					Int("Log count", counter).
					Str("Container", consumer.name).
					Msg("Finished flushing logs")
				break
			} else {
				return decodeErr
			}
		}
	}

	m.log.Info().
		Msg("Flushed all logs to targets")

	return nil
}

// ContainerLogConsumer is a container log lines consumer
type ContainerLogConsumer struct {
	name       string
	prefix     string
	logTargets []LogTarget
	lw         *LogWatch
	Messages   []string
	tempFile   *os.File
	encoder    *gob.Encoder
	isDone     bool
	hasErrored bool
}

// newContainerLogConsumer creates new log consumer for a container that
// - signal if log line matches the pattern
// - push all lines to configured log targets
func newContainerLogConsumer(lw *LogWatch, containerName string, prefix string, logTargets ...LogTarget) (*ContainerLogConsumer, error) {
	consumer := &ContainerLogConsumer{
		name:       containerName,
		prefix:     prefix,
		logTargets: logTargets,
		lw:         lw,
		Messages:   make([]string, 0),
		isDone:     false,
		hasErrored: false,
	}

	if len(logTargets) == 0 {
		return consumer, nil
	}

	tempFile, err := os.CreateTemp("", fmt.Sprintf("%s-%s-datafile.gob", containerName, uuid.NewString()[0:8]))
	if err != nil {
		return nil, err
	}

	consumer.tempFile = tempFile
	consumer.encoder = gob.NewEncoder(tempFile)

	return consumer, nil
}

// Accept accepts the log message from particular container and saves it to the temp gob file
func (g *ContainerLogConsumer) Accept(l testcontainers.Log) {
	g.lw.acceptMutex.Lock()
	defer g.lw.acceptMutex.Unlock()

	if g.hasErrored {
		return
	}

	if g.isDone {
		g.lw.log.Error().
			Str("Test", g.lw.testName).
			Str("Container", g.name).
			Str("Log", string(l.Content)).
			Msg("Consumer has finished, but you are still trying to accept logs. This should never happen")
		return
	}

	g.Messages = append(g.Messages, string(l.Content))
	matches := g.FindMatch(l)
	for i := 0; i < matches; i++ {
		g.lw.notifyTest <- &LogNotification{Container: g.name, Prefix: g.prefix, Log: string(l.Content)}
	}

	// if no log targets are configured, we don't need to save the logs
	if len(g.logTargets) == 0 {
		return
	}

	if g.tempFile == nil || g.encoder == nil {
		g.hasErrored = true
		g.lw.log.Error().
			Msg("temp file or encoder is nil, consumer cannot work, this should never happen")
		return
	}

	content := LogContent{
		TestName:      g.lw.testName,
		ContainerName: g.name,
		Content:       l.Content,
		Time:          time.Now(),
	}

	if err := g.streamLogToTempFile(content); err != nil {
		g.lw.log.Error().
			Err(err).
			Str("Container", g.name).
			Msg("Failed to stream log to temp file")
		g.hasErrored = true
		err = g.tempFile.Close()
		if err != nil {
			g.lw.log.Error().
				Err(err).
				Msg("Failed to close temp file")
		}
	}
}

func (g *ContainerLogConsumer) streamLogToTempFile(content LogContent) error {
	if g.encoder == nil {
		return errors.New("encoder is nil, this should never happen")
	}

	return g.encoder.Encode(content)
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
