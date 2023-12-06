package logstream

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/smartcontractkit/wasp"
)

type LogTarget string

const (
	Loki     LogTarget = "loki"
	File     LogTarget = "file"
	InMemory LogTarget = "in-memory"
)

type HandleLogTarget interface {
	Handle(*ContainerLogConsumer, LogContent) error
	GetLogLocation(map[string]*ContainerLogConsumer) (string, error)
	GetTarget() LogTarget
	SetRunId(string)
	GetRunId() string
}

func getDefaultLogHandlers() map[LogTarget]HandleLogTarget {
	handlers := make(map[LogTarget]HandleLogTarget)
	handlers[Loki] = &LokiLogHandler{}
	handlers[File] = &FileLogHandler{}
	handlers[InMemory] = &InMemoryLogHandler{}

	return handlers
}

// FileLogHandler saves logs to local files
type FileLogHandler struct {
	logFolder         string
	shouldSkipLogging bool
	runId             string
}

func (h *FileLogHandler) Handle(c *ContainerLogConsumer, content LogContent) error {
	if h.shouldSkipLogging {
		return nil
	}

	folder, err := h.getOrCreateLogFolder(content.TestName)
	if err != nil {
		h.shouldSkipLogging = true

		return errors.Wrap(err, "failed to create logs folder. File logging stopped")
	}

	logFileName := filepath.Join(folder, fmt.Sprintf("%s.log", content.ContainerName))
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		h.shouldSkipLogging = true

		return errors.Wrap(err, "failed to open log file. File logging stopped")
	}

	defer logFile.Close()

	if _, err := logFile.WriteString(string(content.Content)); err != nil {
		h.shouldSkipLogging = true

		return errors.Wrap(err, "failed to write to log file. File logging stopped")
	}

	return nil
}

func (h FileLogHandler) GetLogLocation(_ map[string]*ContainerLogConsumer) (string, error) {
	return h.logFolder, nil
}

func (h *FileLogHandler) getOrCreateLogFolder(testname string) (string, error) {
	if h.logFolder == "" {
		folder := fmt.Sprintf("./logs/%s-%s-%s", testname, time.Now().Format("2006-01-02T15-04-05"), h.runId)
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			return "", err
		}
		h.logFolder = folder
	}

	return h.logFolder, nil
}

func (h FileLogHandler) GetTarget() LogTarget {
	return File
}

func (h *FileLogHandler) SetRunId(runId string) {
	h.runId = runId
}

func (h *FileLogHandler) GetRunId() string {
	return h.runId
}

// LokiLogHandler sends logs to Loki
type LokiLogHandler struct {
	grafanaUrl        string
	shouldSkipLogging bool
	runId             string
}

func (h *LokiLogHandler) Handle(c *ContainerLogConsumer, content LogContent) error {
	if h.shouldSkipLogging {
		c.lw.log.Warn().Str("Test", content.TestName).Msg("Skipping pushing logs to Loki for this test")
		return nil
	}

	if c.lw.loki == nil {
		loki, err := wasp.NewLokiClient(wasp.NewEnvLokiConfig())
		if err != nil {
			c.lw.log.Error().Err(err).Msg("Failed to create Loki client")
			h.shouldSkipLogging = true

			return err
		}
		c.lw.loki = loki
	}

	err := c.lw.loki.Handle(model.LabelSet{
		"type":         "log_stream",
		"test":         model.LabelValue(content.TestName),
		"container_id": model.LabelValue(content.ContainerName),
		"run_id":       model.LabelValue(h.runId),
	}, content.Time, string(content.Content))

	return err
}

func (h *LokiLogHandler) GetLogLocation(consumers map[string]*ContainerLogConsumer) (string, error) {
	if h.grafanaUrl != "" {
		return h.grafanaUrl, nil
	}

	grafanaBaseUrl := os.Getenv("GRAFANA_URL")
	if grafanaBaseUrl == "" {
		return "", errors.New("GRAFANA_URL env var is not set")
	}

	grafanaBaseUrl = strings.TrimSuffix(grafanaBaseUrl, "/")

	rangeFrom := time.Now()
	rangeTo := time.Now().Add(time.Minute) //just to make sure we get the last message

	var sb strings.Builder
	sb.WriteString(grafanaBaseUrl)
	sb.WriteString("/d/ddf75041-1e39-42af-aa46-361fe4c36e9e/ci-e2e-tests-logs?orgId=1&")
	sb.WriteString(fmt.Sprintf("var-run_id=%s", h.runId))

	if len(consumers) == 0 {
		return "", errors.New("no Loki consumers found")
	}

	for _, c := range consumers {
		if c.hasLogTarget(Loki) {
			sb.WriteString(fmt.Sprintf("&var-container_id=%s", c.name))
		}

		if c.GetStartTime().Before(rangeFrom) {
			rangeFrom = c.GetStartTime()
		}
	}

	sb.WriteString(fmt.Sprintf("&from=%d&to=%d", rangeFrom.UnixMilli(), rangeTo.UnixMilli()))
	h.grafanaUrl = sb.String()

	return h.grafanaUrl, nil
}

func (h LokiLogHandler) GetTarget() LogTarget {
	return Loki
}

func (h *LokiLogHandler) SetRunId(runId string) {
	h.runId = runId
}

func (h *LokiLogHandler) GetRunId() string {
	return h.runId
}

// InMemoryLogHandler stores logs in memory
type InMemoryLogHandler struct {
	logs  map[string][]LogContent
	runId string
}

func (h *InMemoryLogHandler) Handle(c *ContainerLogConsumer, content LogContent) error {
	if h.logs == nil {
		h.logs = make(map[string][]LogContent)
	}

	if _, ok := h.logs[content.ContainerName]; !ok {
		h.logs[content.ContainerName] = make([]LogContent, 0)
	} else {
		h.logs[content.ContainerName] = append(h.logs[content.ContainerName], content)
	}

	return nil
}

func (h InMemoryLogHandler) GetLogLocation(_ map[string]*ContainerLogConsumer) (string, error) {
	return "", nil
}

func (h InMemoryLogHandler) GetTarget() LogTarget {
	return InMemory
}

func (h *InMemoryLogHandler) SetRunId(runId string) {
	h.runId = runId
}

func (h *InMemoryLogHandler) GetRunId() string {
	return h.runId
}
