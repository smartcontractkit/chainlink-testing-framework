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

	"github.com/avast/retry-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/testsummary"
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
	Terminate(context.Context) error
}

// LogWatch is a test helper struct to monitor docker container logs for some patterns
// and push their logs into Loki for further analysis
type LogWatch struct {
	testName                     string
	log                          zerolog.Logger
	loki                         *wasp.LokiClient
	containers                   []LogProducingContainer
	consumers                    map[string]*ContainerLogConsumer
	logTargetHandlers            map[LogTarget]HandleLogTarget
	enabledLogTargets            []LogTarget
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

	runId, err := runid.GetOrGenerateRunId()
	if err != nil {
		return nil, err
	}

	logWatch := &LogWatch{
		testName:                     testName,
		log:                          l,
		consumers:                    make(map[string]*ContainerLogConsumer, 0),
		logTargetHandlers:            getDefaultLogHandlers(),
		logProducerTimeout:           time.Duration(10 * time.Second),
		logProducerTimeoutRetryLimit: 10,
		enabledLogTargets:            envLogTargets,
		runId:                        runId,
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

func WithLogProducerRetryLimit(retryLimit int) Option {
	return func(lw *LogWatch) {
		lw.logProducerTimeoutRetryLimit = retryLimit
	}
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// ConnectContainer connects consumer to selected container and starts testcontainers.LogProducer
func (m *LogWatch) ConnectContainer(ctx context.Context, container LogProducingContainer, prefix string) error {
	name, err := container.Name(ctx)
	if err != nil {
		return err
	}
	name = strings.Replace(name, "/", "", 1)
	prefix = strings.Replace(prefix, "/", "", 1)

	if prefix == "" {
		prefix = name
	}

	if _, ok := m.consumers[name]; ok {
		return errors.Errorf("container %s is already connected", name)
	}

	enabledLogTargets := make([]LogTarget, 0)
	for logTarget := range m.logTargetHandlers {
		enabledLogTargets = append(enabledLogTargets, logTarget)
	}

	cons, err := newContainerLogConsumer(ctx, m, container, prefix, enabledLogTargets...)
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
							return container.StartLogProducer(ctx, timeout)
						},
							retry.Attempts(uint(retryLimit)),
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

					// time.Sleep(500 * time.Millisecond)
				}
			case <-done:
				return
			}
		}
	}(cons.logListeningDone, m.logProducerTimeout, m.logProducerTimeoutRetryLimit)

	return err
}

func (m *LogWatch) GetConsumers() map[string]*ContainerLogConsumer {
	return m.consumers
}

// Shutdown disconnects all containers, stops notifications
func (m *LogWatch) Shutdown(context context.Context) error {
	var err error
	for _, c := range m.consumers {
		discErr := m.DisconnectContainer(c)
		if discErr != nil {
			m.log.Error().
				Err(err).
				Str("Name", c.name).
				Msg("Failed to disconnect container")

			if err == nil {
				err = discErr
			} else {
				err = errors.Wrap(err, discErr.Error())
			}
		}
	}

	if m.loki != nil {
		m.loki.Stop()
	}

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
func (m *LogWatch) DisconnectContainer(consumer *ContainerLogConsumer) error {
	if consumer.isDone {
		return nil
	}

	consumer.isDone = true
	consumer.logListeningDone <- struct{}{}
	defer close(consumer.logListeningDone)

	if consumer.container.IsRunning() {
		m.log.Info().Str("container", consumer.container.GetContainerID()).Msg("Disconnecting container")
		return consumer.container.StopLogProducer()
	}

	return nil
}

var noOpConsumerFn = func(consumer *ContainerLogConsumer) error {
	return nil
}

// ContainerLogs return all logs for the particular container
func (m *LogWatch) ContainerLogs(name string) ([]string, error) {
	logs := []string{}
	var getLogsFn = func(consumer *ContainerLogConsumer, log LogContent) error {
		if consumer.name == name {
			logs = append(logs, string(log.Content))
		}
		return nil
	}

	err := m.getAllLogsAndExecute(noOpConsumerFn, getLogsFn, noOpConsumerFn)
	if err != nil {
		return []string{}, err
	}

	return logs, err
}

func (m *LogWatch) getAllLogsAndExecute(preExecuteFn func(consumer *ContainerLogConsumer) error, executeFn func(consumer *ContainerLogConsumer, log LogContent) error, cleanUpFn func(consumer *ContainerLogConsumer) error) error {
	m.acceptMutex.Lock()
	defer m.acceptMutex.Unlock()

	var loopErr error
	var attachError = func(err error) {
		if err == nil {
			return
		}
		if loopErr == nil {
			loopErr = err
		} else {
			loopErr = errors.Wrap(loopErr, err.Error())
		}
	}

	for _, consumer := range m.consumers {
		// nothing to do if no log targets are configured
		if len(consumer.logTargets) == 0 {
			continue
		}

		if consumer.tempFile == nil {
			return errors.Errorf("temp file is nil for container %s, this should never happen", consumer.name)
		}

		preExecuteErr := preExecuteFn(consumer)
		if preExecuteErr != nil {
			m.log.Error().
				Err(preExecuteErr).
				Str("Container", consumer.name).
				Msg("Failed to run pre-execute function")
			attachError(preExecuteErr)
			break
		}

		// set the cursor to the end of the file, when done to resume writing
		//revive:disable
		defer func() {
			_, deferErr := consumer.tempFile.Seek(0, 2)
			attachError(deferErr)
		}()
		//revive:enable

		_, err := consumer.tempFile.Seek(0, 0)
		if err != nil {
			return err
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
				executeErr := executeFn(consumer, log)
				if executeErr != nil {
					m.log.Error().
						Err(executeErr).
						Str("Container", consumer.name).
						Msg("Failed to run execute function")
					attachError(preExecuteErr)
					break LOG_LOOP
				}
			} else if errors.Is(decodeErr, io.EOF) {
				m.log.Info().
					Int("Log count", counter).
					Str("Container", consumer.name).
					Msg("Finished getting logs")
				break
			} else {
				return decodeErr
			}
		}

		c := consumer

		// done on purpose
		//revive:disable
		defer func() {
			attachError(cleanUpFn(c))
		}()
		//revive:enable
	}

	return loopErr
}

// FlushLogsToTargets flushes all logs for all consumers (containers) to their targets
func (m *LogWatch) FlushLogsToTargets() error {
	var preExecuteFn = func(consumer *ContainerLogConsumer) error {
		// do not accept any new logs
		consumer.isDone = true

		return nil
	}
	var flushLogsFn = func(consumer *ContainerLogConsumer, log LogContent) error {
		for _, logTarget := range consumer.logTargets {
			if handler, ok := consumer.lw.logTargetHandlers[logTarget]; ok {
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
			}
		}

		return nil
	}

	var closeTempFileFn = func(consumer *ContainerLogConsumer) error {
		if consumer.tempFile == nil {
			return errors.Errorf("temp file is nil for container %s, this should never happen", consumer.name)
		}

		return consumer.tempFile.Close()
	}

	flushErr := m.getAllLogsAndExecute(preExecuteFn, flushLogsFn, closeTempFileFn)
	if flushErr != nil {
		m.log.Info().
			Msg("Finished flushing logs")
	} else {
		m.log.Info().
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
	lw               *LogWatch
	tempFile         *os.File
	encoder          *gob.Encoder
	isDone           bool
	hasErrored       bool
	logListeningDone chan struct{}
	container        LogProducingContainer
}

// newContainerLogConsumer creates new log consumer for a container that
// - signal if log line matches the pattern
// - push all lines to configured log targets
func newContainerLogConsumer(ctx context.Context, lw *LogWatch, container LogProducingContainer, prefix string, logTargets ...LogTarget) (*ContainerLogConsumer, error) {
	containerName, err := container.Name(ctx)
	if err != nil {
		return nil, err
	}

	containerName = strings.Replace(containerName, "/", "", 1)

	consumer := &ContainerLogConsumer{
		name:             containerName,
		prefix:           prefix,
		logTargets:       logTargets,
		lw:               lw,
		isDone:           false,
		hasErrored:       false,
		logListeningDone: make(chan struct{}, 1),
		container:        container,
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

func (g *ContainerLogConsumer) MarkAsErrored() {
	g.hasErrored = true
	g.isDone = true
	close(g.logListeningDone)
}

func (g *ContainerLogConsumer) GetContainer() LogProducingContainer {
	return g.container
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
			case "in-memory":
				envLogTargets = append(envLogTargets, InMemory)
			default:
				return []LogTarget{}, errors.Errorf("unknown log target: %s", target)
			}
		}

		return envLogTargets, nil
	}

	return []LogTarget{}, nil
}
