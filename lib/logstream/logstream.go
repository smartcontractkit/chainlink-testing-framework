package logstream

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/testsummary"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/runid"
)

const NO_TEST = "no_test"

// LogNotification notification about log line match for some container
type LogNotification struct {
	Container string
	Prefix    string
	Log       string
}

// LogProducingContainer is a facade that needs to be implemented by any container that wants to be connected to LogStream
type LogProducingContainer interface {
	Name(ctx context.Context) (string, error)
	FollowOutput(consumer tc.LogConsumer)
	StartLogProducer(ctx context.Context, opts ...tc.LogProductionOption) error
	StopLogProducer() error
	GetLogProductionErrorChannel() <-chan error
	IsRunning() bool
	GetContainerID() string
	Terminate(context.Context) error
}

// LogStream is a test helper struct to monitor docker container logs for some patterns
// and push their logs into Loki for further analysis
type LogStream struct {
	testName          string
	log               zerolog.Logger
	loki              *wasp.LokiClient
	containers        []LogProducingContainer
	consumers         map[string]*ContainerLogConsumer
	consumerMutex     sync.RWMutex
	logTargetHandlers map[LogTarget]HandleLogTarget
	enabledLogTargets []LogTarget
	acceptMutex       sync.Mutex
	loggingConfig     config.LoggingConfig
}

// LogContent is a representation of log that will be send to Loki
type LogContent struct {
	TestName      string
	ContainerName string
	Content       []byte
	Time          time.Time
}

type Option func(*LogStream)

// NewLogStream creates a new LogStream instance, with Loki client only if Loki log target is enabled (lazy init)
func NewLogStream(t *testing.T, loggingConfig *config.LoggingConfig, options ...Option) (*LogStream, error) {
	if loggingConfig == nil {
		return nil, errors.New("logging config cannot be nil")
	}

	l := logging.GetLogger(nil, "LOGSTREAM_LOG_LEVEL").With().Str("Component", "LogStream").Logger()
	var testName string
	if t == nil {
		testName = NO_TEST
	} else {
		testName = t.Name()
	}

	logTargets, err := getLogTargetsFromConfig(*loggingConfig)
	if err != nil {
		return nil, err
	}

	runId, err := runid.GetOrGenerateRunId(loggingConfig.RunId)
	if err != nil {
		return nil, err
	}

	loggingConfig.RunId = &runId

	logWatch := &LogStream{
		testName:          testName,
		log:               l,
		consumers:         make(map[string]*ContainerLogConsumer, 0),
		logTargetHandlers: getDefaultLogHandlers(),
		enabledLogTargets: logTargets,
		loggingConfig:     *loggingConfig,
	}

	for _, option := range options {
		option(logWatch)
	}

	if err := logWatch.validateLogTargets(); err != nil {
		return nil, err
	}

	l.Info().Str("Run_id", *logWatch.loggingConfig.RunId).Msg("LogStream initialized")

	return logWatch, nil
}

// validateLogTargets validates that all enabled log targets have a handler and disables handlers that are not enabled
func (m *LogStream) validateLogTargets() error {
	for _, wantedTarget := range m.enabledLogTargets {
		found := false
		for knownTarget := range m.logTargetHandlers {
			if knownTarget == wantedTarget {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("no handler found for log target: %s", wantedTarget)
		}
	}

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
		m.log.Warn().Msg("No log targets enabled. LogStream will not persist any logs")
	}

	return nil
}

// WithCustomLogHandler allows to override default log handler for particular log target
func WithCustomLogHandler(logTarget LogTarget, handler HandleLogTarget) Option {
	return func(lw *LogStream) {
		lw.logTargetHandlers[logTarget] = handler
	}
}

// fibonacci is a helper function for retrying log producer
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// ConnectContainer connects consumer to selected container, starts testcontainers.LogProducer and listens to it's failures in a detached goroutine
func (m *LogStream) ConnectContainer(ctx context.Context, container LogProducingContainer, prefix string) error {
	name, err := container.Name(ctx)
	if err != nil {
		return err
	}
	name = strings.Replace(name, "/", "", 1)
	prefix = strings.Replace(prefix, "/", "", 1)

	if prefix == "" {
		prefix = name
	}

	enabledLogTargets := make([]LogTarget, 0)
	for logTarget := range m.logTargetHandlers {
		enabledLogTargets = append(enabledLogTargets, logTarget)
	}

	cons, err := newContainerLogConsumer(ctx, m, container, prefix, enabledLogTargets...)
	if err != nil {
		return err
	}

	m.log.Trace().
		Str("Prefix", prefix).
		Str("Name", name).
		Str("Timeout", m.loggingConfig.LogStream.LogProducerTimeout.String()).
		Msg("Connecting container logs")
	m.consumerMutex.Lock()
	defer m.consumerMutex.Unlock()
	m.consumers[name] = cons
	m.containers = append(m.containers, container)
	container.FollowOutput(cons)
	err = container.StartLogProducer(ctx, testcontainers.WithLogProductionTimeout(m.loggingConfig.LogStream.LogProducerTimeout.Duration))

	go func(done chan struct{}, timeout time.Duration, retryLimit int) {
		defer m.log.Trace().Str("Container name", name).Msg("Disconnected container logs")
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
			case logErr := <-container.GetLogProductionErrorChannel():
				if logErr != nil {
					m.log.Error().
						Err(err).
						Str("Container name", name).
						Msg("Log producer errored")
					if shouldRetry() {
						backoff := fibonacci(currentAttempt)
						timeout = timeout + time.Duration(backoff)*time.Second
						m.log.Info().
							Str("Prefix", prefix).
							Str("Container name", name).
							Str("Timeout", timeout.String()).
							Msgf("Retrying connection and listening to container logs. Attempt %d/%d", currentAttempt, retryLimit)

						//TODO if there are many failures here we could save the file and restore it content if we fail to
						//create a new temp file; we remove the previous one to avoid log duplication, because new log producer
						//fetches logs from the beginning
						if resetErr := cons.ResetTempFile(); resetErr != nil {
							m.log.Error().
								Err(resetErr).
								Str("Container name", name).
								Msg("Failed to reset temp file. Stopping logging")

							cons.MarkAsErrored()

							return
						}

						startErr := retry.Do(func() error {
							return container.StartLogProducer(ctx, tc.WithLogProductionTimeout(timeout))
						},
							retry.Attempts(uint(retryLimit)), // nolint gosec
							retry.Delay(1*time.Second),
							retry.OnRetry(func(n uint, err error) {
								m.log.Info().
									Str("Container name", name).
									Str("Attempt", fmt.Sprintf("%d/%d", n+1, retryLimit)).
									Msg("Waiting for log producer to stop before restarting it")
							}),
						)

						if startErr != nil {
							m.log.Error().
								Err(err).
								Str("Container name", name).
								Msg("Previously running log producer couldn't be stopped. Used all retry attempts. Won't try again")

							cons.MarkAsErrored()

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

						cons.MarkAsErrored()

						return
					}
				}
			case <-done:
				return
			}
		}
	}(cons.logListeningDone, m.loggingConfig.LogStream.LogProducerTimeout.Duration, int(*m.loggingConfig.LogStream.LogProducerRetryLimit)) // nolint gosec

	return err
}

// GetConsumers returns all consumers
func (m *LogStream) GetConsumers() map[string]*ContainerLogConsumer {
	m.consumerMutex.RLock()
	defer m.consumerMutex.RUnlock()

	return m.consumers
}

// wrapError wraps existing error with new error
func wrapError(existingErr, newErr error) error {
	if existingErr == nil {
		return newErr
	}
	return fmt.Errorf("%w: %w", existingErr, newErr)
}

var noOpPostDisconnectFn = func(m *LogStream) error { return nil }

// Shutdown disconnects all containers and stops all consumers
func (m *LogStream) Shutdown(context context.Context) error {
	return m.shutdownWithFunction(context, noOpPostDisconnectFn)
}

// shutdownWithFunction disconnects all containers and stops all consumers and executes postDisconnectFn after all
// containers are disconnected, but before Loki is shutdown
func (m *LogStream) shutdownWithFunction(context context.Context, postDisconnectFn func(m *LogStream) error) error {
	var wrappedErr error

	var containers []LogProducingContainer
	m.consumerMutex.RLock()
	for _, c := range m.consumers {
		containers = append(containers, c.container)
	}
	m.consumerMutex.RUnlock()

	// first disconnect all containers, so that no new logs are accepted
	for _, container := range containers {
		name, err := container.Name(context)
		if err != nil {
			m.log.Error().
				Err(err).
				Str("Name", name).
				Msg("Failed to get container name")
			wrappedErr = wrapError(wrappedErr, err)

			continue
		}

		if err := m.DisconnectContainer(container); err != nil {
			m.log.Error().
				Err(err).
				Str("Name", name).
				Msg("Failed to disconnect container")

			wrappedErr = wrapError(wrappedErr, err)
		}
	}

	if err := postDisconnectFn(m); err != nil {
		wrappedErr = wrapError(wrappedErr, err)
	}

	if m.loki != nil {
		m.loki.StopNow()
	}

	return wrappedErr
}

// FlushAndShutdown flushes all logs to their targets and shuts down the log stream in a default sequence
func (m *LogStream) FlushAndShutdown() error {
	var logFlushFn = func(m *LogStream) error {
		if err := m.FlushLogsToTargets(); err != nil {
			m.log.Error().
				Err(err).
				Msg("Failed to flush logs to targets")

			return err
		}

		return nil
	}

	return m.shutdownWithFunction(context.Background(), logFlushFn)
}

type LogWriter = func(testName string, name string, location interface{}) error

// PrintLogTargetsLocations prints all log targets locations to stdout
func (m *LogStream) PrintLogTargetsLocations() {
	m.SaveLogTargetsLocations(func(testName string, name string, location interface{}) error {
		m.log.Info().Str("Test", testName).Str("Handler", name).Interface("Location", location).Msg("Log location")
		return nil
	})
}

// SaveLogTargetsLocations saves all log targets locations to test summary
func (m *LogStream) SaveLogLocationInTestSummary() {
	m.SaveLogTargetsLocations(func(testName string, name string, location interface{}) error {
		return testsummary.AddEntry(testName, name, location)
	})
}

// SaveLogTargetsLocations saves all log targets locations to test summary
func (m *LogStream) GetLogLocation() string {
	var logLocation string
	m.SaveLogTargetsLocations(func(testName string, name string, location interface{}) error {
		logLocation = location.(string)
		return nil
	})

	return logLocation
}

// SaveLogTargetsLocations saves all log targets given writer
func (m *LogStream) SaveLogTargetsLocations(writer LogWriter) {
	for _, handler := range m.logTargetHandlers {
		name := string(handler.GetTarget())
		m.consumerMutex.RLock()
		location, err := handler.GetLogLocation(m.consumers)
		m.consumerMutex.RUnlock()
		if err != nil {
			if strings.Contains(err.Error(), ShorteningFailedErr) {
				m.log.Warn().Str("Handler", name).Err(err).Msg("Failed to shorten Grafana URL, won't output any url")
			} else {
				m.log.Error().Str("Handler", name).Err(err).Msg("Failed to get log location")
				continue
			}
		}

		if err := writer(m.testName, name, location); err != nil {
			m.log.Error().Str("Handler", name).Err(err).Msg("Failed to write log location")
		}
	}
}

// Stop stops the consumer and closes listening channel (it won't be accepting any logs from now on)
func (g *ContainerLogConsumer) stop() error {
	if g.isDone {
		return nil
	}

	g.isDone = true
	g.logListeningDone <- struct{}{}
	defer close(g.logListeningDone)

	return nil
}

// DisconnectContainer disconnects particular container
func (m *LogStream) DisconnectContainer(container LogProducingContainer) error {
	var err error

	if container.IsRunning() {
		m.log.Trace().Str("container", container.GetContainerID()).Msg("Disconnecting container")
		err = container.StopLogProducer()
	}

	consumerFound := false
	m.consumerMutex.RLock()
	for _, consumer := range m.consumers {
		if consumer.container.GetContainerID() == container.GetContainerID() {
			consumerFound = true
			if stopErr := consumer.stop(); err != nil {
				m.log.Error().
					Err(stopErr).
					Str("Name", consumer.name).
					Msg("Failed to stop consumer")
				err = wrapError(err, stopErr)
			}
			break
		}
	}
	m.consumerMutex.RUnlock()

	if !consumerFound {
		m.log.Warn().
			Str("container ID", container.GetContainerID()).
			Msg("No consume found for container")
	}

	return err
}

// ContainerLogs return all logs for particular container
func (m *LogStream) ContainerLogs(name string) ([]string, error) {
	logs := []string{}
	var getLogsFn = func(consumer *ContainerLogConsumer, log LogContent) error {
		if consumer.name == name {
			logs = append(logs, string(log.Content))
		}
		return nil
	}

	err := m.GetAllLogsAndConsume(NoOpConsumerFn, getLogsFn, NoOpConsumerFn)
	if err != nil {
		return []string{}, err
	}

	return logs, err
}

type ConsumerConsumingFn = func(consumer *ContainerLogConsumer) error
type ConsumerLogConsumingFn = func(consumer *ContainerLogConsumer, log LogContent) error

// NoOpConsumerFn is a no-op consumer function
func NoOpConsumerFn(consumer *ContainerLogConsumer) error {
	return nil
}

// GetAllLogsAndConsume gets all logs for all consumers (containers) and consumes them using consumeLogFn
func (m *LogStream) GetAllLogsAndConsume(preExecuteFn ConsumerConsumingFn, consumeLogFn ConsumerLogConsumingFn, postExecuteFn ConsumerConsumingFn) (loopErr error) {
	m.acceptMutex.Lock()
	defer m.acceptMutex.Unlock()

	var attachError = func(err error) {
		if err == nil {
			return
		}
		if loopErr == nil {
			loopErr = err
		} else {
			loopErr = wrapError(loopErr, err)
		}
	}

	m.consumerMutex.RLock()
	defer m.consumerMutex.RUnlock()
	for _, consumer := range m.consumers {
		// nothing to do if no log targets are configured
		if len(consumer.logTargets) == 0 {
			continue
		}

		if consumer.tempFile == nil {
			attachError(fmt.Errorf("temp file is nil for container %s, this should never happen", consumer.name))
			return
		}

		preExecuteErr := preExecuteFn(consumer)
		if preExecuteErr != nil {
			m.log.Error().
				Err(preExecuteErr).
				Str("Container", consumer.name).
				Msg("Failed to run pre-execute function")
			attachError(preExecuteErr)
			continue
		}

		// set the cursor to the end of the file, when done to resume writing, unless it was closed
		//revive:disable
		defer func() {
			if !consumer.isDone {
				_, deferErr := consumer.tempFile.Seek(0, 2)
				attachError(deferErr)
			}
		}()
		//revive:enable

		_, seekErr := consumer.tempFile.Seek(0, 0)
		if seekErr != nil {
			attachError(seekErr)
			return
		}

		decoder := gob.NewDecoder(consumer.tempFile)
		counter := 0

		//TODO handle in batches?
	LOG_LOOP:
		for {
			var log LogContent
			decodeErr := decoder.Decode(&log)
			if decodeErr == nil {
				counter++
				consumeErr := consumeLogFn(consumer, log)
				if consumeErr != nil {
					m.log.Error().
						Err(consumeErr).
						Str("Container", consumer.name).
						Msg("Failed to consume log")
					attachError(consumeErr)
					break LOG_LOOP
				}
			} else if errors.Is(decodeErr, io.EOF) {
				m.log.Debug().
					Int("Log count", counter).
					Str("Container", consumer.name).
					Msg("Finished collecting logs")
				break
			} else {
				m.log.Error().
					Err(decodeErr).
					Str("Container", consumer.name).
					Msg("Failed to decode log")
				attachError(decodeErr)
				return
			}
		}

		postExecuteErr := postExecuteFn(consumer)
		if postExecuteErr != nil {
			m.log.Error().
				Err(postExecuteErr).
				Str("Container", consumer.name).
				Msg("Failed to run post-execute function")
			attachError(postExecuteErr)
			continue
		}
	}

	return
}

// FlushLogsToTargets flushes all logs for all consumers (containers) to their targets
func (m *LogStream) FlushLogsToTargets() error {
	var preExecuteFn = func(consumer *ContainerLogConsumer) error {
		// do not accept any new logs
		consumer.isDone = true

		for _, handler := range m.logTargetHandlers {
			consumer.ls.log.Debug().
				Str("container name", consumer.name).
				Str("Handler", string(handler.GetTarget())).
				Msg("Initializing log target handler")

			if err := handler.Init(consumer); err != nil {
				return err
			}
		}

		return nil
	}
	var flushLogsFn = func(consumer *ContainerLogConsumer, log LogContent) error {
		for _, logTarget := range consumer.logTargets {
			if handler, ok := consumer.ls.logTargetHandlers[logTarget]; ok {
				if err := handler.Handle(consumer, log); err != nil {
					m.log.Error().
						Err(err).
						Str("Container", consumer.name).
						Str("log target", string(logTarget)).
						Msg("Failed to handle log target. Aborting")
					return err
				}
			} else {
				m.log.Warn().
					Str("Container", consumer.name).
					Str("log target", string(logTarget)).
					Msg("No handler found for log target. Aborting")

				return fmt.Errorf("no handler found for log target: %s", logTarget)
			}
		}

		return nil
	}

	var postExecuteFn = func(consumer *ContainerLogConsumer) error {
		for _, handler := range m.logTargetHandlers {
			consumer.ls.log.Debug().
				Str("container name", consumer.name).
				Str("Handler", string(handler.GetTarget())).
				Msg("Tearing down log target handler")

			if err := handler.Teardown(); err != nil {
				return err
			}
		}

		return nil
	}

	flushErr := m.GetAllLogsAndConsume(preExecuteFn, flushLogsFn, postExecuteFn)
	if flushErr == nil {
		m.log.Debug().
			Msg("Finished flushing logs")
	} else {
		m.log.Error().
			Err(flushErr).
			Msg("Failed to flush logs")
	}

	return flushErr
}

// ContainerLogConsumer is a container log lines consumer
type ContainerLogConsumer struct {
	name             string
	prefix           string
	logTargets       []LogTarget
	ls               *LogStream
	tempFile         *os.File
	encoder          *gob.Encoder
	isDone           bool
	hasErrored       bool
	logListeningDone chan struct{}
	container        LogProducingContainer
	firstLogTs       time.Time
}

// newContainerLogConsumer creates new log consumer for a container that saves logs to a temp file
func newContainerLogConsumer(ctx context.Context, lw *LogStream, container LogProducingContainer, prefix string, logTargets ...LogTarget) (*ContainerLogConsumer, error) {
	containerName, err := container.Name(ctx)
	if err != nil {
		return nil, err
	}

	containerName = strings.Replace(containerName, "/", "", 1)

	consumer := &ContainerLogConsumer{
		name:             containerName,
		prefix:           prefix,
		logTargets:       logTargets,
		ls:               lw,
		isDone:           false,
		hasErrored:       false,
		logListeningDone: make(chan struct{}, 1),
		container:        container,
		firstLogTs:       time.Now(),
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

// GetStartTime returns the time of the first log line
func (g *ContainerLogConsumer) GetStartTime() time.Time {
	return g.firstLogTs
}

// ResetTempFile resets the temp file and gob encoder
func (g *ContainerLogConsumer) ResetTempFile() error {
	if g.tempFile != nil {
		if err := g.tempFile.Close(); err != nil {
			return err
		}
	}

	tempFile, err := os.CreateTemp("", fmt.Sprintf("%s-%s-datafile.gob", g.name, uuid.NewString()[0:8]))
	if err != nil {
		return err
	}

	g.tempFile = tempFile
	g.encoder = gob.NewEncoder(tempFile)

	return nil
}

// MarkAsErrored marks the consumer as errored (which makes it stop accepting logs)
func (g *ContainerLogConsumer) MarkAsErrored() {
	g.hasErrored = true
	g.isDone = true
	close(g.logListeningDone)
}

// GetContainer returns the container that this consumer is connected to
func (g *ContainerLogConsumer) GetContainer() LogProducingContainer {
	return g.container
}

// Accept accepts the log message from particular container and saves it to the temp gob file
func (g *ContainerLogConsumer) Accept(l tc.Log) {
	g.ls.acceptMutex.Lock()
	defer g.ls.acceptMutex.Unlock()

	if g.hasErrored {
		return
	}

	if g.isDone {
		g.ls.log.Error().
			Str("Test", g.ls.testName).
			Str("Container", g.name).
			Str("Log", string(l.Content)).
			Msg("Consumer has finished, but you are still trying to accept logs. This should never happen")
		return
	}

	// if no log targets are configured, we don't need to save the logs
	if len(g.logTargets) == 0 {
		return
	}

	if g.tempFile == nil || g.encoder == nil {
		g.ls.log.Error().
			Str("Container", g.name).
			Msg("temp file or encoder is nil, consumer cannot work, this should never happen")
		g.MarkAsErrored()

		return
	}

	var logMsg struct {
		Ts string `json:"ts"`
	}

	// if we cannot unmarshal it, ignore it
	if err := json.Unmarshal(l.Content, &logMsg); err == nil {
		maybeFirstTs, err := time.Parse(time.RFC3339, logMsg.Ts)
		// if it's not a valid timestamp, ignore it
		if err == nil && maybeFirstTs.Before(g.firstLogTs) {
			g.firstLogTs = maybeFirstTs
		}
	}

	content := LogContent{
		TestName:      g.ls.testName,
		ContainerName: g.name,
		Content:       l.Content,
		Time:          time.Now(),
	}

	if err := g.streamLogToTempFile(content); err != nil {
		g.ls.log.Error().
			Err(err).
			Str("Container", g.name).
			Msg("Failed to stream log to temp file")
		g.hasErrored = true
		err = g.tempFile.Close()
		if err != nil {
			g.ls.log.Error().
				Err(err).
				Msg("Failed to close temp file")
		}
	}
}

// streamLogToTempFile streams log to temp file
func (g *ContainerLogConsumer) streamLogToTempFile(content LogContent) error {
	if g.encoder == nil {
		return errors.New("encoder is nil, this should never happen")
	}

	return g.encoder.Encode(content)
}

// hasLogTarget checks if the consumer has a particular log target
func (g *ContainerLogConsumer) hasLogTarget(logTarget LogTarget) bool {
	for _, lt := range g.logTargets {
		if lt == logTarget {
			return true
		}
	}

	return false
}

// getLogTargetsFromConfig gets log targets from logging config or returns 'file' log targets if none are configured
func getLogTargetsFromConfig(config config.LoggingConfig) ([]LogTarget, error) {
	if config.LogStream != nil && len(config.LogStream.LogTargets) > 0 {
		logTargets := make([]LogTarget, 0)
		for _, target := range config.LogStream.LogTargets {
			switch strings.TrimSpace(strings.ToLower(target)) {
			case "loki":
				logTargets = append(logTargets, Loki)
			case "file":
				logTargets = append(logTargets, File)
			case "in-memory":
				logTargets = append(logTargets, InMemory)
			default:
				return []LogTarget{}, fmt.Errorf("unknown log target: %s", target)
			}
		}

		return logTargets, nil
	}

	return []LogTarget{"file"}, nil
}
